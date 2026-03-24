package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// --------------------
// MOCK DB
// --------------------

type MockDB struct {
	ShouldFail   bool
	ShouldTimeout bool
}

func (m *MockDB) PingContext(ctx context.Context) error {
	if m.ShouldTimeout {
		<-ctx.Done()
		return ctx.Err()
	}
	if m.ShouldFail {
		return errors.New("db failure")
	}
	return nil
}

// --------------------
// TESTS
// --------------------

func TestLivenessHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/live", LivenessHandler)

	req, _ := http.NewRequest("GET", "/live", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestReadiness_Healthy(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test")

	db := &MockDB{}

	router := gin.Default()
	router.GET("/ready", ReadinessHandler(db))

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestReadiness_DBDown(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test")

	db := &MockDB{ShouldFail: true}

	router := gin.Default()
	router.GET("/ready", ReadinessHandler(db))

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestReadiness_DBTimeout(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test")

	db := &MockDB{ShouldTimeout: true}

	router := gin.Default()
	router.GET("/ready", ReadinessHandler(db))

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestReadiness_NoDBConfigured(t *testing.T) {
	os.Unsetenv("DATABASE_URL")

	router := gin.Default()
	router.GET("/ready", ReadinessHandler(nil))

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}