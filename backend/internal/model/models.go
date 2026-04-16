package model

import "time"

// Campaign statuses
const (
	CampaignStatusDraft     = "draft"
	CampaignStatusScheduled = "scheduled"
	CampaignStatusSending   = "sending"
	CampaignStatusSent      = "sent"
	CampaignStatusCompleted = "completed"
)

// Event types
const (
	EventSent      = "sent"
	EventOpened    = "opened"
	EventClicked   = "clicked"
	EventSubmitted = "submitted"
	EventReported  = "reported"
	EventError     = "error"
)

type Tenant struct {
	ID                  int64  `gorm:"primaryKey" json:"id"`
	Name                string `gorm:"size:255;not null" json:"name"`
	Slug                string `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Plan                string `gorm:"size:50;not null;default:free" json:"plan"`
	MaxRecipients       int    `gorm:"not null;default:100" json:"max_recipients"`
	MaxCampaignsPerYear *int   `json:"max_campaigns_per_year"`
	MaxEmailsPerMonth   *int   `json:"max_emails_per_month"`
	IsActive            bool   `gorm:"not null;default:true" json:"is_active"`
	Config              string `gorm:"type:text" json:"config"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type User struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	TenantID     *int64     `gorm:"index" json:"tenant_id"`
	Email        string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Name         string     `gorm:"size:255;not null" json:"name"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Role         string     `gorm:"size:50;not null;default:viewer" json:"role"`
	IsActive     bool       `gorm:"not null;default:true" json:"is_active"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type EmailTemplate struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  *int64    `gorm:"index" json:"tenant_id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Subject   string    `gorm:"size:500;not null" json:"subject"`
	HTMLBody  string    `gorm:"type:mediumtext" json:"html_body"`
	TextBody  string    `gorm:"type:text" json:"text_body"`
	Category  string    `gorm:"size:100" json:"category"`
	Language  string    `gorm:"size:10;not null;default:zh-TW" json:"language"`
	CreatedBy *int64    `json:"created_by"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (EmailTemplate) TableName() string { return "email_templates" }

type LandingPage struct {
	ID                 int64     `gorm:"primaryKey" json:"id"`
	TenantID           *int64    `gorm:"index" json:"tenant_id"`
	Name               string    `gorm:"size:255;not null" json:"name"`
	HTML               string    `gorm:"type:mediumtext" json:"html"`
	CaptureCredentials bool      `gorm:"not null;default:false" json:"capture_credentials"`
	CaptureFields      string    `gorm:"type:text" json:"capture_fields"`
	RedirectURL        string    `gorm:"size:500" json:"redirect_url"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Scenario struct {
	ID            int64          `gorm:"primaryKey" json:"id"`
	TenantID      *int64         `gorm:"index" json:"tenant_id"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	Category      string         `gorm:"size:100" json:"category"`
	Difficulty    int            `gorm:"not null;default:1" json:"difficulty"`
	Language      string         `gorm:"size:10;not null;default:zh-TW" json:"language"`
	TemplateID    *int64         `json:"template_id"`
	PageID        *int64         `json:"page_id"`
	EducationHTML string         `gorm:"type:mediumtext" json:"education_html"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	Template      *EmailTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Page          *LandingPage   `gorm:"foreignKey:PageID" json:"page,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type RecipientGroup struct {
	ID         int64       `gorm:"primaryKey" json:"id"`
	TenantID   int64       `gorm:"not null;index" json:"tenant_id"`
	Name       string      `gorm:"size:255;not null" json:"name"`
	Recipients []Recipient `gorm:"foreignKey:GroupID" json:"recipients"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type Recipient struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TenantID   int64     `gorm:"not null;index" json:"tenant_id"`
	GroupID    int64     `gorm:"not null;index" json:"group_id"`
	Email      string    `gorm:"size:255;not null" json:"email"`
	FirstName  string    `gorm:"size:100" json:"first_name"`
	LastName   string    `gorm:"size:100" json:"last_name"`
	Department string    `gorm:"size:100" json:"department"`
	Gender     string    `gorm:"size:10" json:"gender"`
	Position   string    `gorm:"size:100" json:"position"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type SMTPProfile struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      int64     `gorm:"not null;index" json:"tenant_id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	MailerType    string    `gorm:"size:20;not null;default:smtp" json:"mailer_type"`
	Host          string    `gorm:"size:255" json:"host"`
	Port          *int      `json:"port"`
	Username      string    `gorm:"size:255" json:"username"`
	PasswordEnc   []byte    `gorm:"type:varbinary(512)" json:"-"`
	FromAddress   string    `gorm:"size:255" json:"from_address"`
	FromName      string    `gorm:"size:255" json:"from_name"`
	TLSRequired   bool      `gorm:"not null;default:true" json:"tls_required"`
	MailgunDomain string    `gorm:"size:255" json:"mailgun_domain"`
	MailgunAPIKey []byte    `gorm:"type:varbinary(512)" json:"-"`
	SESRegion     string    `gorm:"size:50" json:"ses_region"`
	SESAccessKey  []byte    `gorm:"type:varbinary(512)" json:"-"`
	SESSecretKey  []byte    `gorm:"type:varbinary(512)" json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Campaign struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	TenantID      int64           `gorm:"not null;index" json:"tenant_id"`
	Name          string          `gorm:"size:255;not null" json:"name"`
	Status        string          `gorm:"size:20;not null;default:draft" json:"status"`
	ScenarioID    *int64          `json:"scenario_id"`
	TemplateID    *int64          `json:"template_id"`
	PageID        *int64          `json:"page_id"`
	SMTPProfileID int64           `gorm:"not null" json:"smtp_profile_id"`
	PhishURL      string          `gorm:"size:500" json:"phish_url"`
	SelectionMode string          `gorm:"size:20;not null;default:all" json:"selection_mode"`
	SamplePercent int             `gorm:"not null;default:100" json:"sample_percent"`
	Departments   []string        `gorm:"serializer:json;type:text" json:"departments"`
	LaunchedAt       *time.Time      `json:"launched_at"`
	SendBy           *time.Time      `json:"send_by"`
	ScheduleStart    *time.Time      `json:"schedule_start"`
	WorkingHoursOnly bool            `gorm:"not null;default:false" json:"working_hours_only"`
	SkipWeekends     bool            `gorm:"not null;default:false" json:"skip_weekends"`
	Timezone         string          `gorm:"size:50;not null;default:UTC" json:"timezone"`
	CompletedAt      *time.Time      `json:"completed_at"`
	CreatedBy     *int64          `json:"created_by"`
	Results       []Result        `gorm:"foreignKey:CampaignID" json:"results,omitempty"`
	Groups        []CampaignGroup `gorm:"foreignKey:CampaignID" json:"groups,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type CampaignGroup struct {
	CampaignID int64 `gorm:"primaryKey" json:"campaign_id"`
	GroupID    int64 `gorm:"primaryKey" json:"group_id"`
}

type Result struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	CampaignID  int64      `gorm:"not null;index" json:"campaign_id"`
	TenantID    int64      `gorm:"not null;index" json:"tenant_id"`
	RecipientID int64      `gorm:"not null" json:"recipient_id"`
	RID         string     `gorm:"column:rid;size:36;uniqueIndex;not null" json:"rid"`
	Status      string     `gorm:"size:20;not null;default:pending" json:"status"`
	SendDate    *time.Time `json:"send_date"`
	SentAt      *time.Time `json:"sent_at"`
	OpenedAt    *time.Time `json:"opened_at"`
	ClickedAt   *time.Time `json:"clicked_at"`
	SubmittedAt *time.Time `json:"submitted_at"`
	ReportedAt  *time.Time `json:"reported_at"`
	ErrorDetail string     `gorm:"type:text" json:"error_detail"`
	Recipient   *Recipient `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`
	Events      []Event    `gorm:"foreignKey:ResultID" json:"events,omitempty"`
}

type Event struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	ResultID   int64     `gorm:"not null;index" json:"result_id"`
	CampaignID int64     `gorm:"not null;index" json:"campaign_id"`
	EventType  string    `gorm:"size:20;not null" json:"event_type"`
	IPAddress  string    `gorm:"size:45" json:"ip_address"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	Detail     string    `gorm:"type:text" json:"detail"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AutoTestConfig struct {
	ID            int64      `gorm:"primaryKey" json:"id"`
	TenantID      int64      `gorm:"uniqueIndex;not null" json:"tenant_id"`
	IsEnabled     bool       `gorm:"not null;default:false" json:"is_enabled"`
	Frequency     string     `gorm:"size:20;not null;default:monthly" json:"frequency"`
	TargetMode    string     `gorm:"size:20;not null;default:all" json:"target_mode"`
	SamplePercent int        `gorm:"not null;default:100" json:"sample_percent"`
	Difficulty    int        `gorm:"not null;default:1" json:"difficulty"`
	NotifyEmails  string     `gorm:"type:text" json:"notify_emails"`
	NextRunAt     *time.Time `json:"next_run_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Subscription struct {
	ID                  int64     `gorm:"primaryKey" json:"id"`
	TenantID            int64     `gorm:"not null;index" json:"tenant_id"`
	Plan                string    `gorm:"size:50;not null" json:"plan"`
	MaxRecipients       int       `gorm:"not null" json:"max_recipients"`
	MaxCampaignsPerYear *int      `json:"max_campaigns_per_year"`
	Features            string    `gorm:"type:text" json:"features"`
	PriceYearly         int       `gorm:"not null;default:0" json:"price_yearly"`
	StartsAt            time.Time `gorm:"not null" json:"starts_at"`
	ExpiresAt           time.Time `gorm:"not null" json:"expires_at"`
	AutoRenew           bool      `gorm:"not null;default:true" json:"auto_renew"`
	Status              string    `gorm:"size:20;not null;default:active" json:"status"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UsageRecord struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	TenantID int64  `gorm:"not null;index" json:"tenant_id"`
	Period   string `gorm:"size:10;not null" json:"period"`
	Metric   string `gorm:"size:50;not null" json:"metric"`
	Value    int64  `gorm:"not null;default:0" json:"value"`
}

type AuditLog struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TenantID   *int64    `gorm:"index" json:"tenant_id"`
	UserID     int64     `gorm:"not null" json:"user_id"`
	UserEmail  string    `gorm:"size:255;not null" json:"user_email"`
	Role       string    `gorm:"size:50" json:"role"`
	Action     string    `gorm:"size:50;not null" json:"action"`
	Resource   string    `gorm:"size:50;not null" json:"resource"`
	ResourceID *int64    `json:"resource_id"`
	Detail     string    `gorm:"type:text" json:"detail"`
	IPAddress  string    `gorm:"size:45" json:"ip_address"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
