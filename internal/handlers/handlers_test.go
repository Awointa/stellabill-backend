package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/services"
)

type stubHealthService struct {
	status services.HealthStatus
}

func (s stubHealthService) GetStatus(context.Context) services.HealthStatus {
	return s.status
}

type stubPlanService struct {
	plans []services.Plan
	err   error
}

func (s stubPlanService) ListPlans(context.Context) ([]services.Plan, error) {
	return s.plans, s.err
}

type stubSubscriptionService struct {
	subscriptions []services.Subscription
	subscription  services.Subscription
	listErr       error
	getErr        error
	lastID        string
}

func (s *stubSubscriptionService) ListSubscriptions(context.Context) ([]services.Subscription, error) {
	return s.subscriptions, s.listErr
}

func (s *stubSubscriptionService) GetSubscription(_ context.Context, id string) (services.Subscription, error) {
	s.lastID = id
	return s.subscription, s.getErr
}

func TestNewRequiresDependencies(t *testing.T) {
	t.Parallel()

	base := Dependencies{
		HealthService:       stubHealthService{},
		PlanService:         stubPlanService{},
		SubscriptionService: &stubSubscriptionService{},
	}

	cases := []struct {
		name string
		deps Dependencies
		want string
	}{
		{
			name: "missing health service",
			deps: Dependencies{
				PlanService:         base.PlanService,
				SubscriptionService: base.SubscriptionService,
			},
			want: "health service is required",
		},
		{
			name: "missing plan service",
			deps: Dependencies{
				HealthService:       base.HealthService,
				SubscriptionService: base.SubscriptionService,
			},
			want: "plan service is required",
		},
		{
			name: "missing subscription service",
			deps: Dependencies{
				HealthService: base.HealthService,
				PlanService:   base.PlanService,
			},
			want: "subscription service is required",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler, err := New(tc.deps)
			if err == nil || err.Error() != tc.want {
				t.Fatalf("expected error %q, got handler=%v err=%v", tc.want, handler, err)
			}
		})
	}
}

func TestHealthUsesInjectedService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, err := New(Dependencies{
		HealthService: stubHealthService{
			status: services.HealthStatus{
				Status:  "green",
				Service: "mock-health",
			},
		},
		PlanService:         stubPlanService{},
		SubscriptionService: &stubSubscriptionService{},
	})
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler.Health(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body services.HealthStatus
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Status != "green" || body.Service != "mock-health" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestListPlansUsesInjectedService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, err := New(Dependencies{
		HealthService: stubHealthService{},
		PlanService: stubPlanService{
			plans: []services.Plan{
				{ID: "plan_pro", Name: "Pro", Amount: "2500", Currency: "USD", Interval: "month"},
			},
		},
		SubscriptionService: &stubSubscriptionService{},
	})
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler.ListPlans(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string][]services.Plan
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got := len(body["plans"]); got != 1 || body["plans"][0].ID != "plan_pro" {
		t.Fatalf("unexpected plans: %+v", body["plans"])
	}
}

func TestGetSubscriptionMapsServiceErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("not found", func(t *testing.T) {
		subscriptions := &stubSubscriptionService{getErr: services.ErrSubscriptionNotFound}
		handler, err := New(Dependencies{
			HealthService:       stubHealthService{},
			PlanService:         stubPlanService{},
			SubscriptionService: subscriptions,
		})
		if err != nil {
			t.Fatalf("new handler: %v", err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "sub_missing"}}

		handler.GetSubscription(c)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
		if subscriptions.lastID != "sub_missing" {
			t.Fatalf("expected requested id to reach service, got %q", subscriptions.lastID)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		handler, err := New(Dependencies{
			HealthService:       stubHealthService{},
			PlanService:         stubPlanService{},
			SubscriptionService: &stubSubscriptionService{getErr: errors.New("db offline")},
		})
		if err != nil {
			t.Fatalf("new handler: %v", err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "sub_123"}}

		handler.GetSubscription(c)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
