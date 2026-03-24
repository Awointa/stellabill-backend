package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newRouter() *gin.Engine {
	r := gin.New()
	Register(r)
	return r
}

func TestRoutes_Health(t *testing.T) {
	r := newRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if w.Code != http.StatusOK {
		t.Errorf("health: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutes_Plans(t *testing.T) {
	r := newRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/plans", nil))
	if w.Code != http.StatusOK {
		t.Errorf("plans: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutes_Subscriptions(t *testing.T) {
	r := newRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/subscriptions", nil))
	if w.Code != http.StatusOK {
		t.Errorf("subscriptions: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutes_SubscriptionByID(t *testing.T) {
	r := newRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/subscriptions/abc", nil))
	if w.Code != http.StatusOK {
		t.Errorf("subscription by id: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutes_CORS_Headers(t *testing.T) {
	r := newRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("CORS origin header: got %q, want %q", got, "*")
	}
}

func TestRoutes_CORS_Preflight(t *testing.T) {
	r := newRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodOptions, "/api/health", nil))
	if w.Code != http.StatusNoContent {
		t.Errorf("preflight: got %d, want %d", w.Code, http.StatusNoContent)
	}
}
