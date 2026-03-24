package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/services"
)

func (h *Handler) ListPlans(c *gin.Context) {
	ctx := context.Background()
	if c.Request != nil {
		ctx = c.Request.Context()
	}

	plans, err := h.planService.ListPlans(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load plans"})
		return
	}

	if plans == nil {
		plans = []services.Plan{}
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}
