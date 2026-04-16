-- PhishGuard Complete Schema (fresh install)
-- MySQL 8.0+

CREATE TABLE IF NOT EXISTS tenants (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    slug            VARCHAR(100) NOT NULL UNIQUE,
    plan            VARCHAR(20) NOT NULL DEFAULT 'pro',
    max_recipients  INT NOT NULL DEFAULT 1000,
    max_campaigns_per_year INT DEFAULT NULL,
    max_emails_per_month INT DEFAULT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    config          JSON,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS users (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT DEFAULT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    name            VARCHAR(200) NOT NULL DEFAULT '',
    password_hash   VARCHAR(255) NOT NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'operator',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    last_login      DATETIME DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_tenant (tenant_id),
    CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS email_templates (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT DEFAULT NULL,
    name            VARCHAR(200) NOT NULL,
    subject         VARCHAR(500) NOT NULL,
    html_body       LONGTEXT NOT NULL,
    text_body       TEXT,
    category        VARCHAR(50) DEFAULT NULL,
    language        VARCHAR(10) NOT NULL DEFAULT 'zh-TW',
    created_by      BIGINT DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_templates_tenant (tenant_id),
    CONSTRAINT fk_templates_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_templates_creator FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS landing_pages (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT DEFAULT NULL,
    name            VARCHAR(200) NOT NULL,
    html            LONGTEXT NOT NULL,
    capture_credentials BOOLEAN NOT NULL DEFAULT FALSE,
    capture_fields  JSON,
    redirect_url    VARCHAR(500) DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_pages_tenant (tenant_id),
    CONSTRAINT fk_pages_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS scenarios (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT DEFAULT NULL,
    name            VARCHAR(200) NOT NULL,
    category        VARCHAR(50) NOT NULL,
    difficulty      TINYINT NOT NULL DEFAULT 2,
    language        VARCHAR(10) NOT NULL DEFAULT 'zh-TW',
    template_id     BIGINT NOT NULL,
    page_id         BIGINT NOT NULL,
    education_html  LONGTEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_scenarios_tenant (tenant_id),
    CONSTRAINT fk_scenarios_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_scenarios_template FOREIGN KEY (template_id) REFERENCES email_templates(id),
    CONSTRAINT fk_scenarios_page FOREIGN KEY (page_id) REFERENCES landing_pages(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS recipient_groups (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    name            VARCHAR(200) NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_groups_tenant (tenant_id),
    CONSTRAINT fk_groups_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS recipients (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    group_id        BIGINT NOT NULL,
    email           VARCHAR(255) NOT NULL,
    first_name      VARCHAR(100) DEFAULT '',
    last_name       VARCHAR(100) DEFAULT '',
    department      VARCHAR(100) DEFAULT '',
    gender          VARCHAR(10) DEFAULT '',
    position        VARCHAR(100) DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_recipients_tenant (tenant_id),
    INDEX idx_recipients_group (group_id),
    INDEX idx_recipients_dept (tenant_id, department),
    CONSTRAINT fk_recipients_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_recipients_group FOREIGN KEY (group_id) REFERENCES recipient_groups(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS smtp_profiles (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    name            VARCHAR(200) NOT NULL,
    mailer_type     VARCHAR(20) NOT NULL DEFAULT 'smtp',
    host            VARCHAR(255) DEFAULT NULL,
    port            INT DEFAULT NULL,
    username        VARCHAR(255) DEFAULT NULL,
    password_enc    VARBINARY(512) DEFAULT NULL,
    from_address    VARCHAR(255) NOT NULL,
    from_name       VARCHAR(200) DEFAULT '',
    tls_required    BOOLEAN NOT NULL DEFAULT TRUE,
    mailgun_domain  VARCHAR(255) DEFAULT NULL,
    mailgun_api_key VARBINARY(512) DEFAULT NULL,
    ses_region      VARCHAR(50) DEFAULT NULL,
    ses_access_key  VARBINARY(512) DEFAULT NULL,
    ses_secret_key  VARBINARY(512) DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_smtp_tenant (tenant_id),
    CONSTRAINT fk_smtp_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS campaigns (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    name            VARCHAR(200) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    scenario_id     BIGINT DEFAULT NULL,
    template_id     BIGINT DEFAULT NULL,
    page_id         BIGINT DEFAULT NULL,
    smtp_profile_id BIGINT NOT NULL,
    phish_url       VARCHAR(500) NOT NULL,
    selection_mode  VARCHAR(20) NOT NULL DEFAULT 'all',
    sample_percent  INT NOT NULL DEFAULT 100,
    departments     TEXT,
    launched_at     DATETIME DEFAULT NULL,
    send_by         DATETIME DEFAULT NULL,
    schedule_start  DATETIME DEFAULT NULL,
    working_hours_only BOOLEAN NOT NULL DEFAULT FALSE,
    skip_weekends   BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at    DATETIME DEFAULT NULL,
    created_by      BIGINT DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_campaigns_tenant (tenant_id),
    INDEX idx_campaigns_status (tenant_id, status),
    CONSTRAINT fk_campaigns_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_campaigns_scenario FOREIGN KEY (scenario_id) REFERENCES scenarios(id),
    CONSTRAINT fk_campaigns_template FOREIGN KEY (template_id) REFERENCES email_templates(id),
    CONSTRAINT fk_campaigns_page FOREIGN KEY (page_id) REFERENCES landing_pages(id),
    CONSTRAINT fk_campaigns_smtp FOREIGN KEY (smtp_profile_id) REFERENCES smtp_profiles(id),
    CONSTRAINT fk_campaigns_creator FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS campaign_groups (
    campaign_id     BIGINT NOT NULL,
    group_id        BIGINT NOT NULL,
    PRIMARY KEY (campaign_id, group_id),
    CONSTRAINT fk_cg_campaign FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE,
    CONSTRAINT fk_cg_group FOREIGN KEY (group_id) REFERENCES recipient_groups(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS results (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    campaign_id     BIGINT NOT NULL,
    tenant_id       BIGINT NOT NULL,
    recipient_id    BIGINT NOT NULL,
    rid             VARCHAR(36) NOT NULL UNIQUE,
    status          VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    send_date       DATETIME DEFAULT NULL,
    sent_at         DATETIME DEFAULT NULL,
    opened_at       DATETIME DEFAULT NULL,
    clicked_at      DATETIME DEFAULT NULL,
    submitted_at    DATETIME DEFAULT NULL,
    reported_at     DATETIME DEFAULT NULL,
    error_detail    TEXT DEFAULT NULL,
    INDEX idx_results_rid (rid),
    INDEX idx_results_campaign (campaign_id),
    INDEX idx_results_tenant (tenant_id),
    CONSTRAINT fk_results_campaign FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE,
    CONSTRAINT fk_results_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_results_recipient FOREIGN KEY (recipient_id) REFERENCES recipients(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS events (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    result_id       BIGINT NOT NULL,
    campaign_id     BIGINT NOT NULL,
    event_type      VARCHAR(20) NOT NULL,
    ip_address      VARCHAR(45) DEFAULT NULL,
    user_agent      TEXT DEFAULT NULL,
    detail          JSON,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_events_result (result_id),
    INDEX idx_events_campaign (campaign_id),
    CONSTRAINT fk_events_result FOREIGN KEY (result_id) REFERENCES results(id) ON DELETE CASCADE,
    CONSTRAINT fk_events_campaign FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS auto_test_configs (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT NOT NULL UNIQUE,
    is_enabled      BOOLEAN NOT NULL DEFAULT FALSE,
    frequency       VARCHAR(20) NOT NULL DEFAULT 'quarterly',
    target_mode     VARCHAR(20) NOT NULL DEFAULT 'random',
    sample_percent  INT NOT NULL DEFAULT 30,
    difficulty      INT NOT NULL DEFAULT 1,
    notify_emails   JSON,
    next_run_at     DATETIME DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_autotest_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS subscriptions (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    plan            VARCHAR(20) NOT NULL,
    max_recipients  INT NOT NULL,
    max_campaigns_per_year INT DEFAULT NULL,
    features        JSON NOT NULL,
    price_yearly    INT NOT NULL DEFAULT 0,
    starts_at       DATE NOT NULL,
    expires_at      DATE NOT NULL,
    auto_renew      BOOLEAN NOT NULL DEFAULT FALSE,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_subs_tenant (tenant_id),
    CONSTRAINT fk_subs_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS usage_records (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    period          VARCHAR(7) NOT NULL,
    metric          VARCHAR(30) NOT NULL,
    value           BIGINT NOT NULL DEFAULT 0,
    UNIQUE KEY uk_usage (tenant_id, period, metric),
    CONSTRAINT fk_usage_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS audit_logs (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id       BIGINT DEFAULT NULL,
    user_id         BIGINT NOT NULL,
    user_email      VARCHAR(255) NOT NULL,
    role            VARCHAR(20) NOT NULL,
    action          VARCHAR(50) NOT NULL,
    resource        VARCHAR(50) NOT NULL,
    resource_id     BIGINT DEFAULT NULL,
    detail          JSON,
    ip_address      VARCHAR(45) DEFAULT NULL,
    user_agent      TEXT DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_audit_tenant_time (tenant_id, created_at),
    INDEX idx_audit_user (user_id),
    INDEX idx_audit_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
