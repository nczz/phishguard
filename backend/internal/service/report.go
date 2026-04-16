package service

import (
	"github.com/nczz/phishguard/internal/repo"
)

type OverviewStats struct {
	TotalCampaigns  int64   `json:"total_campaigns"`
	TotalRecipients int64   `json:"total_recipients"`
	AvgClickRate    float64 `json:"avg_click_rate"`
}

type ReportService struct {
	ResultRepo *repo.ResultRepo
}

func (s *ReportService) GetCampaignFunnel(tenantID, campaignID int64) (*repo.FunnelStats, error) {
	return s.ResultRepo.GetFunnelStats(tenantID, campaignID)
}

func (s *ReportService) GetDepartmentStats(tenantID, campaignID int64) ([]repo.DepartmentStat, error) {
	return s.ResultRepo.GetDepartmentStats(tenantID, campaignID)
}

func (s *ReportService) GetOverview(tenantID int64, campaignIDs []int64) (*OverviewStats, error) {
	overview := &OverviewStats{TotalCampaigns: int64(len(campaignIDs))}
	var totalClickRate float64

	for _, cid := range campaignIDs {
		stats, err := s.ResultRepo.GetFunnelStats(tenantID, cid)
		if err != nil {
			return nil, err
		}
		overview.TotalRecipients += stats.Total
		if stats.Total > 0 {
			totalClickRate += float64(stats.Clicked) / float64(stats.Total)
		}
	}

	if overview.TotalCampaigns > 0 {
		overview.AvgClickRate = totalClickRate / float64(overview.TotalCampaigns)
	}

	return overview, nil
}
