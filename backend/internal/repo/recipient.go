package repo

import (
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

type RecipientRepo struct{ DB *gorm.DB }

func (r *RecipientRepo) CreateGroup(g *model.RecipientGroup) error {
	return r.DB.Create(g).Error
}

func (r *RecipientRepo) FindGroupsByTenant(tenantID int64) ([]model.RecipientGroup, error) {
	var groups []model.RecipientGroup
	err := r.DB.Debug().Preload("Recipients").Where("tenant_id = ?", tenantID).Find(&groups).Error
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
	err := r.DB.Where("tenant_id = ? AND group_id IN ? AND is_active = ?", tenantID, groupIDs, true).Find(&recipients).Error
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

func (r *RecipientRepo) UpsertRecipients(tenantID, groupID int64, recipients []model.Recipient, sync bool) (created, updated, deactivated int64, err error) {
	importedEmails := make(map[string]bool, len(recipients))
	for _, rec := range recipients {
		importedEmails[rec.Email] = true
		var existing model.Recipient
		result := r.DB.Where("tenant_id = ? AND email = ?", tenantID, rec.Email).First(&existing)
		if result.Error != nil {
			rec.IsActive = true
			if err = r.DB.Create(&rec).Error; err != nil {
				return
			}
			created++
		} else {
			err = r.DB.Model(&existing).Updates(map[string]interface{}{
				"group_id":   groupID,
				"first_name": rec.FirstName,
				"last_name":  rec.LastName,
				"department": rec.Department,
				"gender":     rec.Gender,
				"position":   rec.Position,
				"is_active":  true,
			}).Error
			if err != nil {
				return
			}
			updated++
		}
	}

	if sync {
		res := r.DB.Model(&model.Recipient{}).
			Where("tenant_id = ? AND group_id = ? AND is_active = ? AND email NOT IN ?", tenantID, groupID, true, mapKeys(importedEmails)).
			Update("is_active", false)
		if res.Error != nil {
			err = res.Error
			return
		}
		deactivated = res.RowsAffected
	}
	return
}

func mapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (r *RecipientRepo) UpdateRecipient(tenantID, id int64, email, firstName, lastName, department, gender, position string) error {
	return r.DB.Model(&model.Recipient{}).Where("tenant_id = ? AND id = ?", tenantID, id).Updates(map[string]interface{}{
		"email": email, "first_name": firstName, "last_name": lastName,
		"department": department, "gender": gender, "position": position,
	}).Error
}

func (r *RecipientRepo) DeleteRecipient(tenantID, id int64) error {
	return r.DB.Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.Recipient{}).Error
}

func (r *RecipientRepo) BatchDelete(tenantID int64, ids []int64) error {
	return r.DB.Where("tenant_id = ? AND id IN ?", tenantID, ids).Delete(&model.Recipient{}).Error
}

func (r *RecipientRepo) BatchSetActive(tenantID int64, ids []int64, active bool) error {
	return r.DB.Model(&model.Recipient{}).Where("tenant_id = ? AND id IN ?", tenantID, ids).Update("is_active", active).Error
}
