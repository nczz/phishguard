package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard response helpers — use these in all new handlers.
// Existing handlers will be migrated gradually.

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

func msg(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, gin.H{"error": message})
}

func notFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}

func serverError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
