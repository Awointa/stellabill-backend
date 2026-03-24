package handlers

import (
	"errors"

	"stellarbill-backend/internal/services"
)

type Dependencies struct {
	HealthService       services.HealthService
	PlanService         services.PlanService
	SubscriptionService services.SubscriptionService
}

type Handler struct {
	healthService       services.HealthService
	planService         services.PlanService
	subscriptionService services.SubscriptionService
}

func New(deps Dependencies) (*Handler, error) {
	switch {
	case deps.HealthService == nil:
		return nil, errors.New("health service is required")
	case deps.PlanService == nil:
		return nil, errors.New("plan service is required")
	case deps.SubscriptionService == nil:
		return nil, errors.New("subscription service is required")
	}

	return &Handler{
		healthService:       deps.HealthService,
		planService:         deps.PlanService,
		subscriptionService: deps.SubscriptionService,
	}, nil
}
