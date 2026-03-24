package services

import (
	"context"
	"errors"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type HealthStatus struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type Plan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Interval    string `json:"interval"`
	Description string `json:"description,omitempty"`
}

type Subscription struct {
	ID          string `json:"id"`
	PlanID      string `json:"plan_id,omitempty"`
	Customer    string `json:"customer,omitempty"`
	Status      string `json:"status"`
	Amount      string `json:"amount,omitempty"`
	Interval    string `json:"interval,omitempty"`
	NextBilling string `json:"next_billing,omitempty"`
}

type HealthService interface {
	GetStatus(ctx context.Context) HealthStatus
}

type PlanService interface {
	ListPlans(ctx context.Context) ([]Plan, error)
}

type SubscriptionService interface {
	ListSubscriptions(ctx context.Context) ([]Subscription, error)
	GetSubscription(ctx context.Context, id string) (Subscription, error)
}

type StaticHealthService struct {
	serviceName string
}

func NewStaticHealthService(serviceName string) *StaticHealthService {
	if serviceName == "" {
		serviceName = "stellarbill-backend"
	}

	return &StaticHealthService{serviceName: serviceName}
}

func (s *StaticHealthService) GetStatus(context.Context) HealthStatus {
	return HealthStatus{
		Status:  "ok",
		Service: s.serviceName,
	}
}

type StaticPlanService struct {
	plans []Plan
}

func NewStaticPlanService(plans []Plan) *StaticPlanService {
	if plans == nil {
		plans = []Plan{}
	}

	return &StaticPlanService{plans: plans}
}

func (s *StaticPlanService) ListPlans(context.Context) ([]Plan, error) {
	return append([]Plan(nil), s.plans...), nil
}

type PlaceholderSubscriptionService struct{}

func NewPlaceholderSubscriptionService() *PlaceholderSubscriptionService {
	return &PlaceholderSubscriptionService{}
}

func (s *PlaceholderSubscriptionService) ListSubscriptions(context.Context) ([]Subscription, error) {
	return []Subscription{}, nil
}

func (s *PlaceholderSubscriptionService) GetSubscription(_ context.Context, id string) (Subscription, error) {
	return Subscription{
		ID:     id,
		Status: "placeholder",
	}, nil
}
