package repo

import (
	"time"

	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

type CampaignRepo struct{ DB *gorm.DB }

func (r *CampaignRepo) Create(c *model.Campaign) error {
	return r.DB.Create(c).Error
}

func (r *CampaignRepo) FindAllByTenant(tenantID int64) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	err := r.DB.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&campaigns).Error
	return campaigns, err
}

func (r *CampaignRepo) FindByID(tenantID, id int64) (*model.Campaign, error) {
	var c model.Campaign
	err := r.DB.Preload("Results.Recipient").Preload("Groups").
		Where("tenant_id = ? AND id = ?", tenantID, id).First(&c).Error
	return &c, err
}

func (r *CampaignRepo) Update(c *model.Campaign) error {
	return r.DB.Save(c).Error
}

func (r *CampaignRepo) Delete(tenantID, id int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("campaign_id IN (SELECT id FROM results WHERE campaign_id = ? AND tenant_id = ?)", id, tenantID).
			Delete(&model.Event{}).Error; err != nil {
			return err
		}
		if err := tx.Where("campaign_id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Result{}).Error; err != nil {
			return err
		}
		if err := tx.Where("campaign_id = ?", id).Delete(&model.CampaignGroup{}).Error; err != nil {
			return err
		}
		return tx.Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.Campaign{}).Error
	})
}

func (r *CampaignRepo) FindByStatus(status string) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	err := r.DB.Where("status = ?", status).Find(&campaigns).Error
	return campaigns, err
}

func (r *CampaignRepo) CreateCampaignGroups(groups []model.CampaignGroup) error {
	return r.DB.Create(&groups).Error
}

func (r *CampaignRepo) CountByTenantThisYear(tenantID int64) (int64, error) {
	var count int64
	startOfYear := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.DB.Model(&model.Campaign{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, startOfYear).Count(&count).Error
	return count, err
}
