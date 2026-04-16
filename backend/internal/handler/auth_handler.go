package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
)

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, user, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (h *Handler) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)
	user, err := h.UserRepo.FindByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
