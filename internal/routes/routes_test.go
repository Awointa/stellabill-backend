package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/handlers"
	"stellarbill-backend/internal/services"
)

func TestRegisterValidatesDependencies(t *testing.T) {
	t.Parallel()

	handler, err := handlers.New(handlers.Dependencies{
		HealthService:       services.NewStaticHealthService("test-service"),
		PlanService:         services.NewStaticPlanService(nil),
		SubscriptionService: services.NewPlaceholderSubscriptionService(),
	})
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	cases := []struct {
		name string
		r    *gin.Engine
		deps Dependencies
		want string
	}{
		{
			name: "missing router",
			deps: Dependencies{Handler: handler},
			want: "router is required",
		},
		{
			name: "missing handler",
			r:    gin.New(),
			want: "handler is required",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := Register(tc.r, tc.deps); err == nil || err.Error() != tc.want {
				t.Fatalf("expected error %q, got %v", tc.want, err)
			}
		})
	}
}

func TestRegisterUsesInjectedHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, err := handlers.New(handlers.Dependencies{
		HealthService: services.NewStaticHealthService("route-mock"),
		PlanService: services.NewStaticPlanService([]services.Plan{
			{ID: "starter", Name: "Starter", Amount: "0", Currency: "USD", Interval: "month"},
		}),
		SubscriptionService: services.NewPlaceholderSubscriptionService(),
	})
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	router := gin.New()
	if err := Register(router, Dependencies{Handler: handler}); err != nil {
		t.Fatalf("register: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/plans", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	var body map[string][]services.Plan
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got := len(body["plans"]); got != 1 || body["plans"][0].ID != "starter" {
		t.Fatalf("unexpected body: %+v", body)
	}
}
