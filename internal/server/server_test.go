package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handleHealth(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSolveEndpoint(t *testing.T) {
	body := `{"fuel":"CH4","equivalence_ratio":1.0,"inlet_temperature":298.15}`
	req := httptest.NewRequest(http.MethodPost, "/api/solve", strings.NewReader(body))
	w := httptest.NewRecorder()
	handleSolve(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "adiabatic_flame_t") {
		t.Fatalf("expected adiabatic_flame_t in response, got %s", w.Body.String())
	}
}

func TestSolveInvalidFuel(t *testing.T) {
	body := `{"fuel":"","equivalence_ratio":1.0,"inlet_temperature":298.15}`
	req := httptest.NewRequest(http.MethodPost, "/api/solve", strings.NewReader(body))
	w := httptest.NewRecorder()
	handleSolve(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSolveMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/solve", nil)
	w := httptest.NewRecorder()
	handleSolve(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestSweepEndpoint(t *testing.T) {
	body := `{"fuel":"CH4","inlet_temperature":298.15}`
	req := httptest.NewRequest(http.MethodPost, "/api/sweep", strings.NewReader(body))
	w := httptest.NewRecorder()
	handleSweep(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "sweep") {
		t.Fatalf("expected sweep in response, got %s", w.Body.String())
	}
}
