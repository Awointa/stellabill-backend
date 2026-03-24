package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/config"
	"stellarbill-backend/internal/handlers"
	"stellarbill-backend/internal/routes"
	"stellarbill-backend/internal/services"
)

func main() {
	cfg := config.Load()
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router, err := newRouter()
	if err != nil {
		log.Fatal(err)
	}

	addr := ":" + cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("Stellarbill backend listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func newRouter() (*gin.Engine, error) {
	handler, err := handlers.New(handlers.Dependencies{
		HealthService:       services.NewStaticHealthService("stellarbill-backend"),
		PlanService:         services.NewStaticPlanService(nil),
		SubscriptionService: services.NewPlaceholderSubscriptionService(),
	})
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	if err := routes.Register(router, routes.Dependencies{Handler: handler}); err != nil {
		return nil, err
	}

	return router, nil
}
