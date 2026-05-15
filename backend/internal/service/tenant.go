package service

import (
	"github.com/nczz/phishguard/internal/model"
	"github.com/nczz/phishguard/internal/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PlatformStats struct {
	TotalTenants  int64 `json:"total_tenants"`
	ActiveTenants int64 `json:"active_tenants"`
}

type TenantService struct {
	TenantRepo *repo.TenantRepo
	UserRepo   *repo.UserRepo
	DB         *gorm.DB
}

func (s *TenantService) Create(name, slug, plan string) (*model.Tenant, error) {
	if plan == "" {
		plan = "free"
	}
	t := &model.Tenant{Name: name, Slug: slug, Plan: plan, IsActive: true}
	return t, s.TenantRepo.Create(t)
}

func (s *TenantService) CreateWithAdmin(name, slug, plan, adminEmail, adminPassword string) (*model.Tenant, error) {
	if plan == "" {
		plan = "free"
	}
	var tenant *model.Tenant
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		t := &model.Tenant{Name: name, Slug: slug, Plan: plan, IsActive: true}
		if err := tx.Create(t).Error; err != nil {
			return err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := tx.Create(&model.User{
			TenantID:     &t.ID,
			Email:        adminEmail,
			Name:         "Tenant Admin",
			PasswordHash: string(hash),
			Role:         "tenant_admin",
			IsActive:     true,
		}).Error; err != nil {
			return err
		}
		tenant = t
		return SeedTenantData(tx, t.ID)
	})
	return tenant, err
}

func (s *TenantService) GetDashboardStats() (*PlatformStats, error) {
	total, err := s.TenantRepo.CountAll()
	if err != nil {
		return nil, err
	}
	tenants, err := s.TenantRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var active int64
	for _, t := range tenants {
		if t.IsActive {
			active++
		}
	}
	return &PlatformStats{TotalTenants: total, ActiveTenants: active}, nil
}
