package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	TenantID *int64 `json:"tenant_id,omitempty"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		token, err := jwt.ParseWithClaims(h[7:], &Claims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("claims", token.Claims.(*Claims))
		c.Next()
	}
}

func GetClaims(c *gin.Context) *Claims {
	v, _ := c.Get("claims")
	cl, _ := v.(*Claims)
	return cl
}

func GetTenantID(c *gin.Context) *int64 {
	if cl := GetClaims(c); cl != nil {
		return cl.TenantID
	}
	return nil
}

func GetUserID(c *gin.Context) int64 {
	if cl := GetClaims(c); cl != nil {
		return cl.UserID
	}
	return 0
}

func GetRole(c *gin.Context) string {
	if cl := GetClaims(c); cl != nil {
		return cl.Role
	}
	return ""
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := GetRole(c)
		for _, allowed := range roles {
			if r == allowed {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}

func GenerateToken(secret string, userID int64, tenantID *int64, role, email string) (string, error) {
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
