package repo

import (
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

type SMTPRepo struct{ DB *gorm.DB }

func (r *SMTPRepo) Create(p *model.SMTPProfile) error {
	return r.DB.Create(p).Error
}

func (r *SMTPRepo) FindAllByTenant(tenantID int64) ([]model.SMTPProfile, error) {
	var list []model.SMTPProfile
	err := r.DB.Where("tenant_id = ?", tenantID).Find(&list).Error
	return list, err
}

func (r *SMTPRepo) FindByID(tenantID, id int64) (*model.SMTPProfile, error) {
	var p model.SMTPProfile
	err := r.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&p).Error
	return &p, err
}

func (r *SMTPRepo) Update(tenantID int64, p *model.SMTPProfile) error {
	return r.DB.Where("id = ? AND tenant_id = ?", p.ID, tenantID).Save(p).Error
}

func (r *SMTPRepo) Delete(tenantID, id int64) error {
	return r.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.SMTPProfile{}).Error
}
