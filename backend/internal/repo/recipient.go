package repo

import (
	"github.com/phishguard/phishguard/internal/model"
	"gorm.io/gorm"
)

type RecipientRepo struct{ DB *gorm.DB }

func (r *RecipientRepo) CreateGroup(g *model.RecipientGroup) error {
	return r.DB.Create(g).Error
}

func (r *RecipientRepo) FindGroupsByTenant(tenantID int64) ([]model.RecipientGroup, error) {
	var groups []model.RecipientGroup
	err := r.DB.Preload("Recipients").Where("tenant_id = ?", tenantID).Find(&groups).Error
	return groups, err
}

func (r *RecipientRepo) FindGroupByID(tenantID, id int64) (*model.RecipientGroup, error) {
	var g model.RecipientGroup
	err := r.DB.Preload("Recipients").Where("tenant_id = ? AND id = ?", tenantID, id).First(&g).Error
	return &g, err
}

func (r *RecipientRepo) DeleteGroup(tenantID, id int64) error {
	return r.DB.Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.RecipientGroup{}).Error
}

func (r *RecipientRepo) BulkCreateRecipients(recipients []model.Recipient) error {
	return r.DB.CreateInBatches(recipients, 500).Error
}

func (r *RecipientRepo) FindByGroupIDs(tenantID int64, groupIDs []int64) ([]model.Recipient, error) {
	var recipients []model.Recipient
	err := r.DB.Where("tenant_id = ? AND group_id IN ?", tenantID, groupIDs).Find(&recipients).Error
	return recipients, err
}

func (r *RecipientRepo) FindByDepartments(tenantID int64, departments []string) ([]model.Recipient, error) {
	var recipients []model.Recipient
	err := r.DB.Where("tenant_id = ? AND department IN ?", tenantID, departments).Find(&recipients).Error
	return recipients, err
}

func (r *RecipientRepo) CountByTenant(tenantID int64) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Recipient{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}
