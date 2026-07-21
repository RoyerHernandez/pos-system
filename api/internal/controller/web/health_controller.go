package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// HealthController handles the health check endpoint.
type HealthController struct {
	db *sqlx.DB
}

// NewHealthController creates a new HealthController.
func NewHealthController(db *sqlx.DB) (*HealthController, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &HealthController{db: db}, nil
}

// Health returns the application and database status.
func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	dbStatus := "connected"
	if err := c.db.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"db":     dbStatus,
	})
}
