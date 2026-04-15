package repo

import (
	"github.com/phishguard/phishguard/internal/model"
	"gorm.io/gorm"
)

type ScenarioRepo struct{ DB *gorm.DB }

func (r *ScenarioRepo) Create(s *model.Scenario) error {
	return r.DB.Create(s).Error
}

func (r *ScenarioRepo) FindAllByTenant(tenantID int64) ([]model.Scenario, error) {
	var list []model.Scenario
	err := r.DB.Preload("Template").Preload("Page").
		Where("tenant_id = ? OR tenant_id IS NULL", tenantID).Find(&list).Error
	return list, err
}

func (r *ScenarioRepo) FindByID(tenantID, id int64) (*model.Scenario, error) {
	var s model.Scenario
	err := r.DB.Preload("Template").Preload("Page").
		Where("id = ? AND (tenant_id = ? OR tenant_id IS NULL)", id, tenantID).First(&s).Error
	return &s, err
}

func (r *ScenarioRepo) Update(tenantID int64, s *model.Scenario) error {
	return r.DB.Model(s).Where("id = ? AND tenant_id = ?", s.ID, tenantID).
		Select("name", "category", "difficulty", "language", "template_id", "page_id", "education_html", "is_active", "updated_at").
		Updates(s).Error
}

func (r *ScenarioRepo) Delete(tenantID, id int64) error {
	return r.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Scenario{}).Error
}

func (r *ScenarioRepo) FindActiveByTenant(tenantID int64) ([]model.Scenario, error) {
	var list []model.Scenario
	err := r.DB.Preload("Template").Preload("Page").
		Where("(tenant_id = ? OR tenant_id IS NULL) AND is_active = ?", tenantID, true).Find(&list).Error
	return list, err
}
