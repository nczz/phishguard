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
	Name          string `json:"name"`
	ScenarioID    *int64 `json:"scenario_id"`
	TemplateID    *int64 `json:"template_id"`
	PageID        *int64 `json:"page_id"`
	SMTPProfileID int64  `json:"smtp_profile_id"`
	GroupIDs      []int64  `json:"group_ids"`
	PhishURL      string   `json:"phish_url"`
	SelectionMode string   `json:"selection_mode"`
	SamplePercent int      `json:"sample_percent"`
	Departments   []string `json:"departments"`

	// Schedule
	SendMode         string `json:"send_mode"`          // immediate / scheduled
	ScheduleStart    string `json:"schedule_start"`      // RFC3339, empty = now
	ScheduleEnd      string `json:"schedule_end"`        // RFC3339, required for scheduled
	WorkingHoursOnly bool   `json:"working_hours_only"`  // 09:00-17:00 only
	SkipWeekends     bool   `json:"skip_weekends"`       // skip Sat/Sun
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

	// Parse schedule
	var schedStart *time.Time
	if req.ScheduleStart != "" {
		t, _ := time.Parse(time.RFC3339, req.ScheduleStart)
		if !t.IsZero() {
			schedStart = &t
		}
	}
	var sendBy *time.Time
	if req.ScheduleEnd != "" {
		t, _ := time.Parse(time.RFC3339, req.ScheduleEnd)
		if !t.IsZero() {
			sendBy = &t
		}
	}

	c := &model.Campaign{
		TenantID:         tenantID,
		Name:             req.Name,
		Status:           model.CampaignStatusDraft,
		ScenarioID:       req.ScenarioID,
		TemplateID:       templateID,
		PageID:           pageID,
		SMTPProfileID:    req.SMTPProfileID,
		PhishURL:         req.PhishURL,
		SendBy:           sendBy,
		ScheduleStart:    schedStart,
		WorkingHoursOnly: req.WorkingHoursOnly,
		SkipWeekends:     req.SkipWeekends,
		SelectionMode:    req.SelectionMode,
		SamplePercent:    req.SamplePercent,
		Departments:      req.Departments,
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

	// Build results with scheduled send dates
	now := time.Now()
	rand.Shuffle(len(recipients), func(i, j int) { recipients[i], recipients[j] = recipients[j], recipients[i] })

	// Determine schedule window
	schedStart := now
	if c.ScheduleStart != nil && c.ScheduleStart.After(now) {
		schedStart = *c.ScheduleStart
	}
	schedEnd := schedStart // immediate: all at start
	if c.SendBy != nil && c.SendBy.After(schedStart) {
		schedEnd = *c.SendBy
	}

	slots := generateTimeSlots(schedStart, schedEnd, c.WorkingHoursOnly, c.SkipWeekends, len(recipients))

	results := make([]model.Result, len(recipients))
	for i, r := range recipients {
		results[i] = model.Result{
			CampaignID:  campaignID,
			TenantID:    tenantID,
			RecipientID: r.ID,
			RID:         uuid.New().String(),
			Status:      model.CampaignStatusScheduled,
			SendDate:    &slots[i],
		}
	}

	if err := s.ResultRepo.BulkCreate(results); err != nil {
		return fmt.Errorf("bulk create results: %w", err)
	}

	c.Status = model.CampaignStatusSending
	c.LaunchedAt = &now
	return s.CampaignRepo.Update(c)
}

// generateTimeSlots distributes N send times across a window,
// respecting working hours and weekend restrictions.
func generateTimeSlots(start, end time.Time, workingHoursOnly, skipWeekends bool, count int) []time.Time {
	if count == 0 {
		return nil
	}

	// If start == end (immediate), return all at start
	if !end.After(start) {
		slots := make([]time.Time, count)
		for i := range slots {
			slots[i] = start
		}
		return slots
	}

	// Collect all valid minute-slots in the window
	validMinutes := []time.Time{}
	cursor := start.Truncate(time.Minute)
	for !cursor.After(end) {
		if isValidSendTime(cursor, workingHoursOnly, skipWeekends) {
			validMinutes = append(validMinutes, cursor)
		}
		cursor = cursor.Add(time.Minute)
	}

	if len(validMinutes) == 0 {
		// No valid slots found — fall back to immediate
		slots := make([]time.Time, count)
		for i := range slots {
			slots[i] = start
		}
		return slots
	}

	// Distribute recipients evenly across valid minutes with jitter
	slots := make([]time.Time, count)
	for i := 0; i < count; i++ {
		idx := i * len(validMinutes) / count
		if idx >= len(validMinutes) {
			idx = len(validMinutes) - 1
		}
		// Add random jitter within the minute (0-59 seconds)
		jitter := time.Duration(rand.Intn(60)) * time.Second
		slots[i] = validMinutes[idx].Add(jitter)
	}
	return slots
}

func isValidSendTime(t time.Time, workingHoursOnly, skipWeekends bool) bool {
	if skipWeekends {
		wd := t.Weekday()
		if wd == time.Saturday || wd == time.Sunday {
			return false
		}
	}
	if workingHoursOnly {
		hour := t.Hour()
		if hour < 9 || hour >= 17 {
			return false
		}
	}
	return true
}
