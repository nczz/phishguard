package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TenantMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cl := GetClaims(c)
		if cl == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}

		if cl.Role == "platform_admin" {
			if h := c.GetHeader("X-Tenant-ID"); h != "" {
				tid, err := strconv.ParseInt(h, 10, 64)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid X-Tenant-ID"})
					return
				}
				if !tenantIsActive(db, tid) {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "tenant is inactive or not found"})
					return
				}
				c.Set("tenant_id", &tid)
			}
			// no header → tenant_id stays nil (platform-level operation)
		} else {
			if cl.TenantID == nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no tenant assigned"})
				return
			}
			if !tenantIsActive(db, *cl.TenantID) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "tenant is inactive or not found"})
				return
			}
			c.Set("tenant_id", cl.TenantID)
		}

		c.Next()
	}
}

func tenantIsActive(db *gorm.DB, tenantID int64) bool {
	var active bool
	err := db.Raw("SELECT is_active FROM tenants WHERE id = ?", tenantID).Scan(&active).Error
	return err == nil && active
}

func GetContextTenantID(c *gin.Context) *int64 {
	v, ok := c.Get("tenant_id")
	if !ok {
		return nil
	}
	tid, _ := v.(*int64)
	return tid
}

func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetContextTenantID(c) == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
			return
		}
		c.Next()
	}
}
