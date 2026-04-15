package repo

import (
	"context"

	"github.com/phishguard/phishguard/internal/middleware"
	"github.com/phishguard/phishguard/internal/model"
	"gorm.io/gorm"
)

type AuditRepo struct{ DB *gorm.DB }

func (r *AuditRepo) Create(log *model.AuditLog) error {
	return r.DB.Create(log).Error
}

func (r *AuditRepo) FindByTenant(tenantID int64, limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	q := r.DB.Where("tenant_id = ?", tenantID)
	if err := q.Model(&model.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *AuditRepo) FindAll(limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	if err := r.DB.Model(&model.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.DB.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

// DBAuditLogger implements middleware.AuditLogger.
type DBAuditLogger struct{ Repo *AuditRepo }

func (l *DBAuditLogger) Log(_ context.Context, entry middleware.AuditEntry) error {
	return l.Repo.Create(&model.AuditLog{
		TenantID:   entry.TenantID,
		UserID:     entry.UserID,
		UserEmail:  entry.UserEmail,
		Role:       entry.Role,
		Action:     entry.Action,
		Resource:   entry.Resource,
		ResourceID: entry.ResourceID,
		Detail:     entry.Detail,
		IPAddress:  entry.IPAddress,
		UserAgent:  entry.UserAgent,
	})
}
