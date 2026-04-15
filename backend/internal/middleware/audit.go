package middleware

import (
	"context"
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

func AuditMiddleware(logger AuditLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
			return
		}

		cl := GetClaims(c)
		if cl == nil {
			return
		}

		action, resource := parseAction(method, c.FullPath())

		entry := AuditEntry{
			TenantID:  GetContextTenantID(c),
			UserID:    cl.UserID,
			UserEmail: cl.Email,
			Role:      cl.Role,
			Action:    action,
			Resource:  resource,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		}

		go logger.Log(context.Background(), entry)
	}
}

// parseAction maps HTTP method + route path to action and resource names.
// e.g. POST /api/campaigns → ("campaign.create", "campaign")
//      DELETE /api/campaigns/:id → ("campaign.delete", "campaign")
func parseAction(method, path string) (action, resource string) {
	// strip /api/ prefix and split
	path = strings.TrimPrefix(path, "/api/")
	parts := strings.Split(path, "/")

	// find the first non-parameter segment as resource
	resource = "unknown"
	for _, p := range parts {
		if p != "" && !strings.HasPrefix(p, ":") {
			resource = p
			break
		}
	}

	// singularize: strip trailing 's'
	singular := strings.TrimSuffix(resource, "s")
	if singular == "" {
		singular = resource
	}

	var verb string
	switch method {
	case "POST":
		verb = "create"
	case "PUT", "PATCH":
		verb = "update"
	case "DELETE":
		verb = "delete"
	}

	return singular + "." + verb, singular
}
