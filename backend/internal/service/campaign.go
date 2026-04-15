package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/phishguard/phishguard/internal/model"
	"github.com/phishguard/phishguard/internal/repo"
)

type CreateCampaignRequest struct {
	Name          string     `json:"name"`
	ScenarioID    *int64     `json:"scenario_id"`
	TemplateID    *int64     `json:"template_id"`
	PageID        *int64     `json:"page_id"`
	SMTPProfileID int64      `json:"smtp_profile_id"`
	GroupIDs      []int64    `json:"group_ids"`
	PhishURL      string     `json:"phish_url"`
	SendBy        *time.Time `json:"send_by"`
	SelectionMode string     `json:"selection_mode"` // all, random, department
	SamplePercent int        `json:"sample_percent"`
	Departments   []string   `json:"departments"`
}

type CampaignService struct {
	CampaignRepo  *repo.CampaignRepo
	ResultRepo    *repo.ResultRepo
	RecipientRepo *repo.RecipientRepo
	ScenarioRepo  *repo.ScenarioRepo
}

func (s *CampaignService) CreateCampaign(tenantID int64, req *CreateCampaignRequest) (*model.Campaign, error) {
	templateID := req.TemplateID
	pageID := req.PageID

	if req.ScenarioID != nil {
		sc, err := s.ScenarioRepo.FindByID(tenantID, *req.ScenarioID)
		if err != nil {
			return nil, fmt.Errorf("scenario lookup: %w", err)
		}
		templateID = sc.TemplateID
		pageID = sc.PageID
	}

	c := &model.Campaign{
		TenantID:      tenantID,
		Name:          req.Name,
		Status:        model.CampaignStatusDraft,
		ScenarioID:    req.ScenarioID,
		TemplateID:    templateID,
		PageID:        pageID,
		SMTPProfileID: req.SMTPProfileID,
		PhishURL:      req.PhishURL,
		SendBy:        req.SendBy,
		SelectionMode: req.SelectionMode,
		SamplePercent: req.SamplePercent,
		Departments:   req.Departments,
	}
	if err := s.CampaignRepo.Create(c); err != nil {
		return nil, fmt.Errorf("create campaign: %w", err)
	}

	if len(req.GroupIDs) > 0 {
		groups := make([]model.CampaignGroup, len(req.GroupIDs))
		for i, gid := range req.GroupIDs {
			groups[i] = model.CampaignGroup{CampaignID: c.ID, GroupID: gid}
		}
		if err := s.CampaignRepo.CreateCampaignGroups(groups); err != nil {
			return nil, fmt.Errorf("create campaign groups: %w", err)
		}
	}

	return c, nil
}

func (s *CampaignService) LaunchCampaign(tenantID, campaignID int64) error {
	c, err := s.CampaignRepo.FindByID(tenantID, campaignID)
	if err != nil {
		return fmt.Errorf("find campaign: %w", err)
	}
	if c.Status != model.CampaignStatusDraft {
		return fmt.Errorf("campaign status is %s, expected draft", c.Status)
	}

	// Get recipients from campaign groups
	groupIDs := make([]int64, len(c.Groups))
	for i, g := range c.Groups {
		groupIDs[i] = g.GroupID
	}
	recipients, err := s.RecipientRepo.FindByGroupIDs(tenantID, groupIDs)
	if err != nil {
		return fmt.Errorf("find recipients: %w", err)
	}

	// Apply selection mode
	switch c.SelectionMode {
	case "random":
		rand.Shuffle(len(recipients), func(i, j int) { recipients[i], recipients[j] = recipients[j], recipients[i] })
		n := len(recipients) * c.SamplePercent / 100
		if n < 1 {
			n = 1
		}
		recipients = recipients[:n]
	case "department":
		recipients, err = s.RecipientRepo.FindByDepartments(tenantID, c.Departments)
		if err != nil {
			return fmt.Errorf("find by departments: %w", err)
		}
	}

	// Cooldown: exclude recipients tested in last 30 days
	filtered := make([]model.Recipient, 0, len(recipients))
	for _, r := range recipients {
		recent, _ := s.ResultRepo.FindRecentByRecipientEmail(tenantID, r.Email, 30)
		if len(recent) == 0 {
			filtered = append(filtered, r)
		}
	}
	recipients = filtered

	if len(recipients) == 0 {
		return fmt.Errorf("no recipients selected")
	}

	// Build results with spread send dates
	now := time.Now()
	results := make([]model.Result, len(recipients))
	for i, r := range recipients {
		sendDate := now
		if c.SendBy != nil && c.SendBy.After(now) && len(recipients) > 1 {
			interval := c.SendBy.Sub(now)
			sendDate = now.Add(interval * time.Duration(i) / time.Duration(len(recipients)-1))
		}
		results[i] = model.Result{
			CampaignID:  campaignID,
			TenantID:    tenantID,
			RecipientID: r.ID,
			RID:         uuid.New().String(),
			Status:      model.CampaignStatusScheduled,
			SendDate:    &sendDate,
		}
	}

	if err := s.ResultRepo.BulkCreate(results); err != nil {
		return fmt.Errorf("bulk create results: %w", err)
	}

	c.Status = model.CampaignStatusSending
	c.LaunchedAt = &now
	return s.CampaignRepo.Update(c)
}
