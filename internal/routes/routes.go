package routes

import (
	"errors"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/handlers"
)

type Dependencies struct {
	Handler *handlers.Handler
}

func Register(r *gin.Engine, deps Dependencies) error {
	switch {
	case r == nil:
		return errors.New("router is required")
	case deps.Handler == nil:
		return errors.New("handler is required")
	}

	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.GET("/health", deps.Handler.Health)
		api.GET("/subscriptions", deps.Handler.ListSubscriptions)
		api.GET("/subscriptions/:id", deps.Handler.GetSubscription)
		api.GET("/plans", deps.Handler.ListPlans)
	}

	return nil
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
