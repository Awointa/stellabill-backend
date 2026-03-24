package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()

	if cfg.Env != "development" {
		t.Errorf("Env: got %q, want %q", cfg.Env, "development")
	}
	if cfg.Port != "8080" {
		t.Errorf("Port: got %q, want %q", cfg.Port, "8080")
	}
	if cfg.DBConn != "postgres://localhost/stellarbill?sslmode=disable" {
		t.Errorf("DBConn: got %q", cfg.DBConn)
	}
	if cfg.JWTSecret != "change-me-in-production" {
		t.Errorf("JWTSecret: got %q", cfg.JWTSecret)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://prod-host/db")
	t.Setenv("JWT_SECRET", "supersecret")

	cfg := Load()

	if cfg.Env != "production" {
		t.Errorf("Env: got %q, want %q", cfg.Env, "production")
	}
	if cfg.Port != "9090" {
		t.Errorf("Port: got %q, want %q", cfg.Port, "9090")
	}
	if cfg.DBConn != "postgres://prod-host/db" {
		t.Errorf("DBConn: got %q", cfg.DBConn)
	}
	if cfg.JWTSecret != "supersecret" {
		t.Errorf("JWTSecret: got %q", cfg.JWTSecret)
	}
}
