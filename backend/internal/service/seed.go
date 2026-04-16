package service

import (
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

func SeedTenantData(db *gorm.DB, tenantID int64) error {
	templates := seedTemplates(tenantID)
	for i := range templates {
		if err := db.Create(&templates[i]).Error; err != nil {
			return err
		}
	}

	pages := seedLandingPages(tenantID)
	for i := range pages {
		if err := db.Create(&pages[i]).Error; err != nil {
			return err
		}
	}

	eduHTML := seedEducationHTML()

	scenarios := []model.Scenario{
		{TenantID: &tenantID, Name: "密碼到期通知", Category: "password_reset", Difficulty: 2, Language: "zh-TW",
			TemplateID: &templates[0].ID, PageID: &pages[0].ID, EducationHTML: eduHTML, IsActive: true},
		{TenantID: &tenantID, Name: "包裹到貨通知", Category: "package", Difficulty: 1, Language: "zh-TW",
			TemplateID: &templates[1].ID, PageID: &pages[1].ID, EducationHTML: eduHTML, IsActive: true},
		{TenantID: &tenantID, Name: "薪資單確認", Category: "hr_notice", Difficulty: 2, Language: "zh-TW",
			TemplateID: &templates[2].ID, PageID: &pages[2].ID, EducationHTML: eduHTML, IsActive: true},
		{TenantID: &tenantID, Name: "資安警告通知", Category: "it_alert", Difficulty: 3, Language: "zh-TW",
			TemplateID: &templates[3].ID, PageID: &pages[3].ID, EducationHTML: eduHTML, IsActive: true},
		{TenantID: &tenantID, Name: "發票確認通知", Category: "invoice", Difficulty: 2, Language: "zh-TW",
			TemplateID: &templates[4].ID, PageID: &pages[4].ID, EducationHTML: eduHTML, IsActive: true},
	}
	for i := range scenarios {
		if err := db.Create(&scenarios[i]).Error; err != nil {
			return err
		}
	}

	group := model.RecipientGroup{TenantID: tenantID, Name: "範例員工"}
	if err := db.Create(&group).Error; err != nil {
		return err
	}
	return db.CreateInBatches([]model.Recipient{
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-wang@example.com", FirstName: "小明", LastName: "王", Department: "業務部", Gender: "男", Position: "業務經理"},
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-chen@example.com", FirstName: "小華", LastName: "陳", Department: "財務部", Gender: "女", Position: "會計"},
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-lin@example.com", FirstName: "小美", LastName: "林", Department: "研發部", Gender: "女", Position: "工程師"},
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-zhang@example.com", FirstName: "大偉", LastName: "張", Department: "行政部", Gender: "男", Position: "行政專員"},
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-li@example.com", FirstName: "雅婷", LastName: "李", Department: "業務部", Gender: "不指定", Position: "業務代表"},
	}, 100).Error
}
