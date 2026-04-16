package repo

import (
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

type TemplateRepo struct{ DB *gorm.DB }

func (r *TemplateRepo) Create(t *model.EmailTemplate) error {
	return r.DB.Create(t).Error
}

func (r *TemplateRepo) FindAllByTenant(tenantID int64) ([]model.EmailTemplate, error) {
	var list []model.EmailTemplate
	err := r.DB.Where("tenant_id = ? OR tenant_id IS NULL", tenantID).Find(&list).Error
	return list, err
}

func (r *TemplateRepo) FindByID(tenantID, id int64) (*model.EmailTemplate, error) {
	var t model.EmailTemplate
	err := r.DB.Where("id = ? AND (tenant_id = ? OR tenant_id IS NULL)", id, tenantID).First(&t).Error
	return &t, err
}

func (r *TemplateRepo) Update(tenantID int64, t *model.EmailTemplate) error {
	return r.DB.Model(t).Where("id = ? AND tenant_id = ?", t.ID, tenantID).
		Select("name", "subject", "html_body", "text_body", "category", "language", "updated_at").
		Updates(t).Error
}

func (r *TemplateRepo) Delete(tenantID, id int64) error {
	return r.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.EmailTemplate{}).Error
}
