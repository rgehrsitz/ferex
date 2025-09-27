package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ferexapp/ferex/internal/core"
)

// Clock abstracts time retrieval to simplify testing.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now().UTC() }

// ScenarioStore defines the persistence contract for scenarios regardless of
// the backing medium (in-memory, filesystem, database, etc.).
type ScenarioStore interface {
	Save(ctx context.Context, scenario core.Scenario) error
	Get(ctx context.Context, id string) (core.Scenario, error)
	List(ctx context.Context) ([]core.ScenarioMetadata, error)
	Delete(ctx context.Context, id string) error
}

// ScenarioManager orchestrates creation and updates of scenarios.
type ScenarioManager struct {
	store ScenarioStore
	clock Clock
}

// NewScenarioManager wires the manager with the provided store. If clock is
// nil, a real clock is used.
func NewScenarioManager(store ScenarioStore, clock Clock) *ScenarioManager {
	if clock == nil {
		clock = realClock{}
	}
	return &ScenarioManager{store: store, clock: clock}
}

// CreateScenario seeds a new scenario from the supplied profile. It clones the
// profile data to ensure subsequent profile edits do not implicitly mutate the
// scenario snapshot.
func (m *ScenarioManager) CreateScenario(ctx context.Context, profile core.UserProfile, name string, configure func(*core.Scenario)) (core.Scenario, error) {
	if m == nil {
		return core.Scenario{}, errors.New("scenario manager is nil")
	}
	if err := profile.Validate(); err != nil {
		return core.Scenario{}, err
	}
	scenario := core.Scenario{
		ID:                uuid.NewString(),
		Name:              name,
		ProfileRef:        profile.ID,
		ProfileSnapshot:   profile.Clone(),
		CreatedAt:         m.clock.Now(),
		UpdatedAt:         m.clock.Now(),
		IncludeSupplement: true,
		MonteCarlo: core.SimulationConfig{
			HorizonYears:     40,
			Trials:           1000,
			ConfidenceLevels: []float64{0.1, 0.5, 0.9},
		},
	}
	if configure != nil {
		configure(&scenario)
	}
	if err := m.saveValidated(ctx, scenario); err != nil {
		return core.Scenario{}, err
	}
	return scenario, nil
}

// UpdateScenario persists a modified scenario after revalidating inputs.
func (m *ScenarioManager) UpdateScenario(ctx context.Context, scenario core.Scenario) error {
	if m == nil {
		return errors.New("scenario manager is nil")
	}
	scenario.UpdatedAt = m.clock.Now()
	return m.saveValidated(ctx, scenario)
}

// saveValidated ensures the scenario passes validation before delegating to the
// store.
func (m *ScenarioManager) saveValidated(ctx context.Context, scenario core.Scenario) error {
	if err := scenario.Validate(); err != nil {
		return err
	}
	return m.store.Save(ctx, scenario)
}

// GetScenario retrieves a scenario by ID.
func (m *ScenarioManager) GetScenario(ctx context.Context, id string) (core.Scenario, error) {
	if m == nil {
		return core.Scenario{}, errors.New("scenario manager is nil")
	}
	return m.store.Get(ctx, id)
}

// ListScenarios returns lightweight metadata for selectors and dashboards.
func (m *ScenarioManager) ListScenarios(ctx context.Context) ([]core.ScenarioMetadata, error) {
	if m == nil {
		return nil, errors.New("scenario manager is nil")
	}
	return m.store.List(ctx)
}

// DeleteScenario removes a scenario from storage.
func (m *ScenarioManager) DeleteScenario(ctx context.Context, id string) error {
	if m == nil {
		return errors.New("scenario manager is nil")
	}
	return m.store.Delete(ctx, id)
}
