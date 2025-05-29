package test

import (
	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/health/transport/http"
	"github.com/alexfalkowski/go-service/v2/time"
)

// RegisterHealth for test.
func RegisterHealth(health, live, ready *subscriber.Observer) {
	params := http.RegisterParams{
		Health:    &http.HealthObserver{Observer: health},
		Liveness:  &http.LivenessObserver{Observer: live},
		Readiness: &http.ReadinessObserver{Observer: ready},
	}

	http.Register(params)
}

// NewDBChecker for test.
func NewDBChecker(world *World) (*checker.DBChecker, error) {
	db, err := world.OpenDatabase()
	if err != nil {
		return nil, err
	}

	return checker.NewDBChecker(db, 1*time.Second), nil
}
