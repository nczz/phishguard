package service

import "github.com/phishguard/phishguard/internal/model"

type PlanConfig struct {
	MaxRecipients       int  `json:"max_recipients"`
	MaxCampaignsPerYear int  `json:"max_campaigns_per_year"`
	MaxEmailsPerMonth   int  `json:"max_emails_per_month"`
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
		MaxEmailsPerMonth:   200,
		CustomTemplate:      false,
		AutoTest:            false,
		DepartmentReport:    false,
		CustomDomain:        false,
		AuditLog:            false,
	},
	"pro": {
		MaxRecipients:       1000,
		MaxCampaignsPerYear: 0,
		MaxEmailsPerMonth:   10000,
		CustomTemplate:      true,
		AutoTest:            true,
		DepartmentReport:    true,
		CustomDomain:        false,
		AuditLog:            true,
	},
	"enterprise": {
		MaxRecipients:       0,
		MaxCampaignsPerYear: 0,
		MaxEmailsPerMonth:   0,
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

func GetEffectiveLimits(t *model.Tenant) PlanConfig {
	pc := GetPlanConfig(t.Plan)
	if t.MaxRecipients > 0 {
		pc.MaxRecipients = t.MaxRecipients
	}
	if t.MaxCampaignsPerYear != nil && *t.MaxCampaignsPerYear >= 0 {
		pc.MaxCampaignsPerYear = *t.MaxCampaignsPerYear
	}
	if t.MaxEmailsPerMonth != nil && *t.MaxEmailsPerMonth >= 0 {
		pc.MaxEmailsPerMonth = *t.MaxEmailsPerMonth
	}
	return pc
}
