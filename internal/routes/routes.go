package routes

import (
	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/handlers"
	"stellarbill-backend/internal/middleware"
)

func Register(r *gin.Engine) {
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)
		
		// Feature-flagged endpoints
		api.GET("/subscriptions", 
			middleware.FeatureFlagWithDefault("subscriptions_enabled", true),
			handlers.ListSubscriptions)
		api.GET("/subscriptions/:id", 
			middleware.FeatureFlagWithDefault("subscriptions_enabled", true),
			handlers.GetSubscription)
		api.GET("/plans", 
			middleware.FeatureFlagWithDefault("plans_enabled", true),
			handlers.ListPlans)
		
		// Example of new feature that can be toggled
		api.GET("/billing/new-flow", 
			middleware.FeatureFlag("new_billing_flow"),
			func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "New billing flow is enabled"})
			})
		
		// Example requiring multiple flags
		api.GET("/analytics/advanced", 
			middleware.RequireAllFeatureFlags("advanced_analytics", "subscriptions_enabled"),
			func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Advanced analytics available"})
			})
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
