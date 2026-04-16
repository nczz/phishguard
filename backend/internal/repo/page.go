package repo

import (
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

type PageRepo struct{ DB *gorm.DB }

func (r *PageRepo) Create(p *model.LandingPage) error {
	return r.DB.Create(p).Error
}

func (r *PageRepo) FindAllByTenant(tenantID int64) ([]model.LandingPage, error) {
	var list []model.LandingPage
	err := r.DB.Where("tenant_id = ? OR tenant_id IS NULL", tenantID).Find(&list).Error
	return list, err
}

func (r *PageRepo) FindByID(tenantID, id int64) (*model.LandingPage, error) {
	var p model.LandingPage
	err := r.DB.Where("id = ? AND (tenant_id = ? OR tenant_id IS NULL)", id, tenantID).First(&p).Error
	return &p, err
}

func (r *PageRepo) Update(tenantID int64, p *model.LandingPage) error {
	return r.DB.Model(p).Where("id = ? AND tenant_id = ?", p.ID, tenantID).
		Select("name", "html", "capture_credentials", "capture_fields", "redirect_url", "updated_at").
		Updates(p).Error
}

func (r *PageRepo) Delete(tenantID, id int64) error {
	return r.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.LandingPage{}).Error
}
