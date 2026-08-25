package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"adiab-flame/internal/flame"
)

type Config struct {
	Addr string
}

func ListenAndServe(cfg Config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/solve", handleSolve)
	mux.HandleFunc("/api/sweep", handleSweep)
	mux.HandleFunc("/health", handleHealth)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

type solveRequest struct {
	Fuel             string  `json:"fuel"`
	EquivalenceRatio float64 `json:"equivalence_ratio"`
	InletTemperature float64 `json:"inlet_temperature"`
	PressureAtm      float64 `json:"pressure_atm"`
	Dissociation     bool    `json:"dissociation"`
}

func handleSolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("read body: %v", err))
		return
	}
	var req solveRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON: %v", err))
		return
	}
	if req.Fuel == "" {
		writeError(w, http.StatusBadRequest, "fuel is required")
		return
	}
	if req.EquivalenceRatio <= 0 {
		writeError(w, http.StatusBadRequest, "equivalence_ratio must be positive")
		return
	}
	if req.InletTemperature <= 0 {
		writeError(w, http.StatusBadRequest, "inlet_temperature must be positive")
		return
	}
	cfg := flame.NewConfig(req.Fuel, req.EquivalenceRatio, req.InletTemperature)
	if req.PressureAtm > 0 {
		cfg.PressureAtm = req.PressureAtm
	}
	cfg.Dissociation = req.Dissociation

	solver := flame.NewSolver()
	result, err := solver.Solve(cfg)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("solve: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"fuel":               result.Fuel,
		"equivalence_ratio":  result.EquivalenceRatio,
		"inlet_temperature":  result.InletTemperature,
		"adiabatic_flame_t":  result.Tad,
		"converged":          result.Converged,
		"iterations":         result.Iterations,
		"max_residual":       result.MaxResidual(),
		"enthalpy_residual":  result.EnthalpyResidual(),
	})
}

func handleSweep(w http.ResponseWriter, r *http.Request) {
	if err := abortSweepContext(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("read body: %v", err))
		return
	}
	var req solveRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON: %v", err))
		return
	}
	if req.Fuel == "" {
		writeError(w, http.StatusBadRequest, "fuel is required")
		return
	}
	if req.InletTemperature <= 0 {
		req.InletTemperature = 298.15
	}

	type point struct {
		Phi float64 `json:"phi"`
		Tad float64 `json:"tad"`
	}
	solver := flame.NewSolver()
	points := make([]point, 0, 20)
	for phi := 0.5; phi <= 2.0; phi += 0.1 {
		cfg := flame.NewConfig(req.Fuel, phi, req.InletTemperature)
		cfg.Dissociation = req.Dissociation
		result, err := solver.Solve(cfg)
		if err != nil {
			continue
		}
		points = append(points, point{Phi: phi, Tad: result.Tad})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"fuel":  req.Fuel,
		"sweep": points,
	})
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
