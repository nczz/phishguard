package repo

import (
	"time"

	"github.com/phishguard/phishguard/internal/model"
	"gorm.io/gorm"
)

type FunnelStats struct {
	Total     int64
	Sent      int64
	Opened    int64
	Clicked   int64
	Submitted int64
	Reported  int64
}

type DepartmentStat struct {
	Department string
	Total      int64
	Clicked    int64
}

type ResultRepo struct{ DB *gorm.DB }

func (r *ResultRepo) BulkCreate(results []model.Result) error {
	return r.DB.CreateInBatches(results, 500).Error
}

func (r *ResultRepo) FindByRID(rid string) (*model.Result, error) {
	var res model.Result
	err := r.DB.Where("rid = ?", rid).First(&res).Error
	return &res, err
}

func (r *ResultRepo) FindByCampaign(tenantID, campaignID int64) ([]model.Result, error) {
	var results []model.Result
	err := r.DB.Preload("Recipient").
		Where("tenant_id = ? AND campaign_id = ?", tenantID, campaignID).Find(&results).Error
	return results, err
}

func (r *ResultRepo) UpdateStatus(id int64, status string) error {
	return r.DB.Model(&model.Result{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ResultRepo) CreateEvent(e *model.Event) error {
	return r.DB.Create(e).Error
}

func (r *ResultRepo) GetFunnelStats(tenantID, campaignID int64) (*FunnelStats, error) {
	var s FunnelStats
	err := r.DB.Model(&model.Result{}).
		Select(`COUNT(*) AS total,
			SUM(CASE WHEN sent_at IS NOT NULL THEN 1 ELSE 0 END) AS sent,
			SUM(CASE WHEN opened_at IS NOT NULL THEN 1 ELSE 0 END) AS opened,
			SUM(CASE WHEN clicked_at IS NOT NULL THEN 1 ELSE 0 END) AS clicked,
			SUM(CASE WHEN submitted_at IS NOT NULL THEN 1 ELSE 0 END) AS submitted,
			SUM(CASE WHEN reported_at IS NOT NULL THEN 1 ELSE 0 END) AS reported`).
		Where("tenant_id = ? AND campaign_id = ?", tenantID, campaignID).
		Scan(&s).Error
	return &s, err
}

func (r *ResultRepo) GetDepartmentStats(tenantID, campaignID int64) ([]DepartmentStat, error) {
	var stats []DepartmentStat
	err := r.DB.Model(&model.Result{}).
		Select(`rec.department AS department, COUNT(*) AS total,
			SUM(CASE WHEN results.clicked_at IS NOT NULL THEN 1 ELSE 0 END) AS clicked`).
		Joins("JOIN recipients rec ON rec.id = results.recipient_id").
		Where("results.tenant_id = ? AND results.campaign_id = ?", tenantID, campaignID).
		Group("rec.department").
		Scan(&stats).Error
	return stats, err
}

func (r *ResultRepo) FindRecentByRecipientEmail(tenantID int64, email string, days int) ([]model.Result, error) {
	var results []model.Result
	since := time.Now().AddDate(0, 0, -days)
	err := r.DB.Joins("JOIN recipients rec ON rec.id = results.recipient_id").
		Where("results.tenant_id = ? AND rec.email = ? AND results.created_at >= ?", tenantID, email, since).
		Find(&results).Error
	return results, err
}
