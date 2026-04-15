package service

import "github.com/phishguard/phishguard/internal/model"

// PlanDefaults defines what each plan includes
type PlanConfig struct {
	MaxRecipients       int  `json:"max_recipients"`
	MaxCampaignsPerYear int  `json:"max_campaigns_per_year"` // 0 = unlimited
	CustomTemplate      bool `json:"custom_template"`
	AutoTest            bool `json:"auto_test"`
	DepartmentReport    bool `json:"department_report"`
	CustomDomain        bool `json:"custom_domain"`
	AuditLog            bool `json:"audit_log"`
}

var PlanDefaults = map[string]PlanConfig{
	"free": {
		MaxRecipients:       50,
		MaxCampaignsPerYear: 4,
		CustomTemplate:      false,
		AutoTest:            false,
		DepartmentReport:    false,
		CustomDomain:        false,
		AuditLog:            false,
	},
	"pro": {
		MaxRecipients:       1000,
		MaxCampaignsPerYear: 0, // unlimited
		CustomTemplate:      true,
		AutoTest:            true,
		DepartmentReport:    true,
		CustomDomain:        false,
		AuditLog:            true,
	},
	"enterprise": {
		MaxRecipients:       0, // unlimited
		MaxCampaignsPerYear: 0,
		CustomTemplate:      true,
		AutoTest:            true,
		DepartmentReport:    true,
		CustomDomain:        true,
		AuditLog:            true,
	},
}

func GetPlanConfig(plan string) PlanConfig {
	if c, ok := PlanDefaults[plan]; ok {
		return c
	}
	return PlanDefaults["free"]
}

// GetEffectiveLimits returns the actual limits for a tenant (plan defaults overridden by tenant-level settings)
func GetEffectiveLimits(t *model.Tenant) PlanConfig {
	pc := GetPlanConfig(t.Plan)
	// Tenant-level overrides (if admin has set custom values)
	if t.MaxRecipients > 0 {
		pc.MaxRecipients = t.MaxRecipients
	}
	if t.MaxCampaignsPerYear != nil && *t.MaxCampaignsPerYear >= 0 {
		pc.MaxCampaignsPerYear = *t.MaxCampaignsPerYear
	}
	return pc
}
