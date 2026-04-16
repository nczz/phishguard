package repo

import (
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

type TenantRepo struct{ DB *gorm.DB }

func (r *TenantRepo) Create(t *model.Tenant) error {
	return r.DB.Create(t).Error
}

func (r *TenantRepo) FindAll() ([]model.Tenant, error) {
	var tenants []model.Tenant
	err := r.DB.Find(&tenants).Error
	return tenants, err
}

func (r *TenantRepo) FindByID(id int64) (*model.Tenant, error) {
	var t model.Tenant
	err := r.DB.First(&t, id).Error
	return &t, err
}

func (r *TenantRepo) FindBySlug(slug string) (*model.Tenant, error) {
	var t model.Tenant
	err := r.DB.Where("slug = ?", slug).First(&t).Error
	return &t, err
}

func (r *TenantRepo) Update(t *model.Tenant) error {
	return r.DB.Save(t).Error
}

func (r *TenantRepo) CountAll() (int64, error) {
	var count int64
	err := r.DB.Model(&model.Tenant{}).Count(&count).Error
	return count, err
}
