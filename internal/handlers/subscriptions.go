package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/services"
)

func (h *Handler) ListSubscriptions(c *gin.Context) {
	ctx := context.Background()
	if c.Request != nil {
		ctx = c.Request.Context()
	}

	subscriptions, err := h.subscriptionService.ListSubscriptions(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load subscriptions"})
		return
	}

	if subscriptions == nil {
		subscriptions = []services.Subscription{}
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

func (h *Handler) GetSubscription(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscription id required"})
		return
	}

	ctx := context.Background()
	if c.Request != nil {
		ctx = c.Request.Context()
	}

	subscription, err := h.subscriptionService.GetSubscription(ctx, id)
	if err != nil {
		if errors.Is(err, services.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load subscription"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}
