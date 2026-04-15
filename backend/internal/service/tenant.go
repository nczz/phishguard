package service

import (
	"github.com/phishguard/phishguard/internal/model"
	"github.com/phishguard/phishguard/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type PlatformStats struct {
	TotalTenants  int64 `json:"total_tenants"`
	ActiveTenants int64 `json:"active_tenants"`
}

type TenantService struct {
	TenantRepo *repo.TenantRepo
	UserRepo   *repo.UserRepo
}

func (s *TenantService) Create(name, slug, plan string) (*model.Tenant, error) {
	t := &model.Tenant{Name: name, Slug: slug, Plan: plan, IsActive: true, MaxRecipients: 100}
	return t, s.TenantRepo.Create(t)
}

func (s *TenantService) CreateWithAdmin(name, slug, plan, adminEmail, adminPassword string) (*model.Tenant, error) {
	t, err := s.Create(name, slug, plan)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	err = s.UserRepo.Create(&model.User{
		TenantID:     &t.ID,
		Email:        adminEmail,
		Name:         "Tenant Admin",
		PasswordHash: string(hash),
		Role:         "tenant_admin",
		IsActive:     true,
	})
	return t, err
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
