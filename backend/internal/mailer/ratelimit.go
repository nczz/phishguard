package mailer

import (
	"sync"
	"time"
)

// RateLimiter enforces per-domain and per-provider sending limits
// aligned with official provider policies.
type RateLimiter struct {
	mu            sync.Mutex
	domainCounts  map[string]*windowCounter
	providerCount *windowCounter
	config        RateLimitConfig
}

type RateLimitConfig struct {
	MaxPerDomainPerHour int
	MaxPerSecond        int
	MinInterval         time.Duration
	BurstSize           int // allow short bursts up to this count within 1 second
}

type windowCounter struct {
	counts []time.Time
}

// ProviderLimits — aligned with official documentation:
//
// AWS SES (production):
//   - Default sending rate: 14 emails/sec (account-level, can be increased)
//   - We use 12/sec to maintain 85% headroom
//   - ThrottlingException on exceed → must back off
//   - Ref: https://docs.aws.amazon.com/ses/latest/dg/manage-sending-quotas.html
//
// AWS SES (sandbox):
//   - 1 email/sec, 200 emails/day
//   - Only verified addresses
//
// Mailgun:
//   - API rate: 500 requests per 10 seconds (50/sec)
//   - New/unverified accounts: 100 messages/hour
//   - Verified accounts: Mailgun queues internally, no hard per-second limit
//   - HTTP 429 on rate limit exceed
//   - Ref: https://documentation.mailgun.com/docs/mailgun/api-reference/send/mailgun/metrics/rate-limits-and-quotas
//
// SMTP (generic):
//   - Most providers: 2-5 emails/sec before throttling
//   - Gmail SMTP: ~2000/day, ~20/sec burst
//   - Microsoft 365: 30 messages/min
//   - Conservative defaults for unknown servers
var ProviderLimits = map[string]RateLimitConfig{
	"ses": {
		MaxPerDomainPerHour: 500,  // SES doesn't limit per-domain, but good practice
		MaxPerSecond:        12,   // SES production default is 14/sec, keep 85% margin
		MinInterval:         80 * time.Millisecond,
		BurstSize:           14,   // allow brief burst to SES limit
	},
	"ses_sandbox": {
		MaxPerDomainPerHour: 50,
		MaxPerSecond:        1,    // SES sandbox: strictly 1/sec
		MinInterval:         1000 * time.Millisecond,
		BurstSize:           1,
	},
	"mailgun": {
		MaxPerDomainPerHour: 500,  // Mailgun queues internally
		MaxPerSecond:        40,   // API limit is 50/sec, keep 80% margin
		MinInterval:         25 * time.Millisecond,
		BurstSize:           50,
	},
	"mailgun_unverified": {
		MaxPerDomainPerHour: 100,  // Mailgun unverified: 100/hour total
		MaxPerSecond:        2,
		MinInterval:         500 * time.Millisecond,
		BurstSize:           2,
	},
	"smtp": {
		MaxPerDomainPerHour: 100,  // Conservative for unknown SMTP
		MaxPerSecond:        3,
		MinInterval:         350 * time.Millisecond,
		BurstSize:           5,
	},
}

func NewRateLimiter(providerType string) *RateLimiter {
	cfg, ok := ProviderLimits[providerType]
	if !ok {
		cfg = ProviderLimits["smtp"]
	}
	return &RateLimiter{
		domainCounts:  make(map[string]*windowCounter),
		providerCount: &windowCounter{},
		config:        cfg,
	}
}

// Wait blocks until it's safe to send to the given recipient domain.
// Returns false if domain hourly limit reached (caller should defer to next cycle).
func (rl *RateLimiter) Wait(recipientDomain string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// 1. Check per-domain hourly limit
	dc, ok := rl.domainCounts[recipientDomain]
	if !ok {
		dc = &windowCounter{}
		rl.domainCounts[recipientDomain] = dc
	}
	dc.prune(now.Add(-1 * time.Hour))
	if rl.config.MaxPerDomainPerHour > 0 && len(dc.counts) >= rl.config.MaxPerDomainPerHour {
		return false
	}

	// 2. Enforce per-second rate with burst allowance
	rl.providerCount.prune(now.Add(-1 * time.Second))
	if len(rl.providerCount.counts) >= rl.config.BurstSize {
		// At burst limit — wait until oldest entry in window expires
		oldest := rl.providerCount.counts[0]
		waitUntil := oldest.Add(1 * time.Second)
		rl.mu.Unlock()
		time.Sleep(time.Until(waitUntil))
		rl.mu.Lock()
		now = time.Now()
		rl.providerCount.prune(now.Add(-1 * time.Second))
	}

	// 3. Enforce minimum interval between sends
	if len(rl.providerCount.counts) > 0 {
		last := rl.providerCount.counts[len(rl.providerCount.counts)-1]
		elapsed := now.Sub(last)
		if elapsed < rl.config.MinInterval {
			rl.mu.Unlock()
			time.Sleep(rl.config.MinInterval - elapsed)
			rl.mu.Lock()
		}
	}

	// Record this send
	sendTime := time.Now()
	dc.counts = append(dc.counts, sendTime)
	rl.providerCount.counts = append(rl.providerCount.counts, sendTime)
	return true
}

func (rl *RateLimiter) GetDomainCount(domain string) int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	dc, ok := rl.domainCounts[domain]
	if !ok {
		return 0
	}
	dc.prune(time.Now().Add(-1 * time.Hour))
	return len(dc.counts)
}

func (wc *windowCounter) prune(before time.Time) {
	i := 0
	for i < len(wc.counts) && wc.counts[i].Before(before) {
		i++
	}
	wc.counts = wc.counts[i:]
}
