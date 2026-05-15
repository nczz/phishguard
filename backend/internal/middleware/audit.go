package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuditEntry struct {
	TenantID   *int64 `json:"tenant_id"`
	UserID     int64  `json:"user_id"`
	UserEmail  string `json:"user_email"`
	Role       string `json:"role"`
	Action     string `json:"action"`
	Resource   string `json:"resource"`
	ResourceID *int64 `json:"resource_id"`
	Detail     string `json:"detail"`
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
}

type AuditLogger interface {
	Log(ctx context.Context, entry AuditEntry) error
}

// sensitiveKeys are fields that should be redacted from audit detail.
func AuditMiddleware(logger AuditLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// Buffer request body for write operations
		var bodyBytes []byte
		if method == "POST" || method == "PUT" || method == "PATCH" {
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		c.Next()

		if method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
			return
		}

		cl := GetClaims(c)
		if cl == nil {
			return
		}

		action, resource := parseAction(method, c.FullPath())
		resourceID := extractResourceID(c)

		entry := AuditEntry{
			TenantID:   GetContextTenantID(c),
			UserID:     cl.UserID,
			UserEmail:  cl.Email,
			Role:       cl.Role,
			Action:     action,
			Resource:   resource,
			ResourceID: resourceID,
			Detail:     summarizeBody(bodyBytes),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		go func() {
			if err := logger.Log(context.Background(), entry); err != nil {
				log.Printf("[AUDIT] failed to write audit log: %v", err)
			}
		}()
	}
}

// summarizeBody extracts a short summary from the request body, redacting sensitive fields.
// Returns valid JSON string or empty string (repo layer handles NULL).
func summarizeBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return ""
	}
	summary := make(map[string]any, len(raw))
	for k, v := range raw {
		if isSensitiveAuditKey(k) {
			summary[k] = "***"
			continue
		}
		switch val := v.(type) {
		case string:
			if len(val) > 100 {
				summary[k] = val[:100] + "..."
			} else {
				summary[k] = val
			}
		default:
			summary[k] = v
		}
	}
	out, _ := json.Marshal(summary)
	if len(out) > 1000 {
		out = out[:1000]
	}
	return string(out)
}

func isSensitiveAuditKey(key string) bool {
	k := strings.ToLower(key)
	sensitiveParts := []string{
		"password",
		"api_key",
		"secret",
		"access_key",
		"token",
		"jwt",
		"encrypt_key",
	}
	for _, part := range sensitiveParts {
		if strings.Contains(k, part) {
			return true
		}
	}
	return false
}

// extractResourceID extracts numeric ID from URL path params.
func extractResourceID(c *gin.Context) *int64 {
	for _, p := range c.Params {
		if p.Key == "id" || p.Key == "uid" {
			if id, err := strconv.ParseInt(p.Value, 10, 64); err == nil {
				return &id
			}
		}
	}
	return nil
}

// parseAction maps HTTP method + route path to action and resource names.
func parseAction(method, path string) (action, resource string) {
	path = strings.TrimPrefix(path, "/api/")
	path = strings.TrimPrefix(path, "admin/")
	parts := strings.Split(path, "/")

	resource = "unknown"
	for _, p := range parts {
		if p != "" && !strings.HasPrefix(p, ":") {
			resource = p
			break
		}
	}

	singular := strings.TrimSuffix(resource, "s")
	if singular == "" {
		singular = resource
	}

	// Detect sub-action (e.g. /campaigns/:id/launch → campaign.launch)
	var verb string
	lastPart := parts[len(parts)-1]
	if !strings.HasPrefix(lastPart, ":") && lastPart != resource {
		verb = lastPart
	} else {
		switch method {
		case "POST":
			verb = "create"
		case "PUT", "PATCH":
			verb = "update"
		case "DELETE":
			verb = "delete"
		}
	}

	return singular + "." + verb, singular
}
