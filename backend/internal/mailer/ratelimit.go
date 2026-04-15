package mailer

import (
	"sync"
	"time"
)

// RateLimiter enforces per-domain and per-provider sending limits
type RateLimiter struct {
	mu            sync.Mutex
	domainCounts  map[string]*windowCounter // per recipient domain
	providerCount *windowCounter            // global provider limit
	config        RateLimitConfig
}

type RateLimitConfig struct {
	// Per recipient domain limits
	MaxPerDomainPerHour int           // e.g. 100 emails/hour to same domain
	// Provider limits
	MaxPerSecond        int           // e.g. SES=14/sec, Mailgun=100/sec, SMTP=5/sec
	MinInterval         time.Duration // minimum gap between sends (fallback)
}

type windowCounter struct {
	counts []time.Time
}

// Provider-specific defaults
var ProviderLimits = map[string]RateLimitConfig{
	"ses": {
		MaxPerDomainPerHour: 200,
		MaxPerSecond:        10,  // SES default is 14/sec, keep margin
		MinInterval:         100 * time.Millisecond,
	},
	"mailgun": {
		MaxPerDomainPerHour: 300,
		MaxPerSecond:        50,  // Mailgun is generous
		MinInterval:         20 * time.Millisecond,
	},
	"smtp": {
		MaxPerDomainPerHour: 100, // conservative for unknown SMTP
		MaxPerSecond:        3,
		MinInterval:         350 * time.Millisecond,
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
// Returns true if OK to send, false if domain hourly limit reached (skip for now).
func (rl *RateLimiter) Wait(recipientDomain string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Check per-domain hourly limit
	dc, ok := rl.domainCounts[recipientDomain]
	if !ok {
		dc = &windowCounter{}
		rl.domainCounts[recipientDomain] = dc
	}
	dc.prune(now.Add(-1 * time.Hour))
	if rl.config.MaxPerDomainPerHour > 0 && len(dc.counts) >= rl.config.MaxPerDomainPerHour {
		return false // domain limit reached, caller should defer
	}

	// Enforce per-second rate: check how many sent in last second
	rl.providerCount.prune(now.Add(-1 * time.Second))
	if rl.config.MaxPerSecond > 0 && len(rl.providerCount.counts) >= rl.config.MaxPerSecond {
		// Wait until the oldest one in the window expires
		oldest := rl.providerCount.counts[0]
		waitUntil := oldest.Add(1 * time.Second)
		rl.mu.Unlock()
		time.Sleep(time.Until(waitUntil))
		rl.mu.Lock()
	}

	// Enforce minimum interval
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

// DomainLimitReached returns emails that were skipped due to domain limits
// so the worker can retry them later
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
