package memory

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/ferexapp/ferex/internal/core"
)

// ScenarioStore provides an in-memory implementation suitable for early
// development and testing. It satisfies the app.ScenarioStore interface.
type ScenarioStore struct {
	mu        sync.RWMutex
	scenarios map[string]core.Scenario
}

// NewScenarioStore constructs an empty in-memory store.
func NewScenarioStore() *ScenarioStore {
	return &ScenarioStore{scenarios: make(map[string]core.Scenario)}
}

// Save persists or updates the scenario.
func (s *ScenarioStore) Save(_ context.Context, scenario core.Scenario) error {
	if s == nil {
		return errors.New("memory store is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scenarios[scenario.ID] = scenario
	return nil
}

// Get retrieves a scenario by ID.
func (s *ScenarioStore) Get(_ context.Context, id string) (core.Scenario, error) {
	if s == nil {
		return core.Scenario{}, errors.New("memory store is nil")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	scenario, ok := s.scenarios[id]
	if !ok {
		return core.Scenario{}, errors.New("scenario not found")
	}
	return scenario, nil
}

// List returns sorted metadata by last updated timestamp descending.
func (s *ScenarioStore) List(_ context.Context) ([]core.ScenarioMetadata, error) {
	if s == nil {
		return nil, errors.New("memory store is nil")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]core.ScenarioMetadata, 0, len(s.scenarios))
	for _, scenario := range s.scenarios {
		items = append(items, core.ScenarioMetadata{
			ID:             scenario.ID,
			Name:           scenario.Name,
			ProfileRef:     scenario.ProfileRef,
			UpdatedAt:      scenario.UpdatedAt,
			RetirementDate: scenario.RetirementDate,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return items, nil
}

// Delete removes the scenario from the store.
func (s *ScenarioStore) Delete(_ context.Context, id string) error {
	if s == nil {
		return errors.New("memory store is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.scenarios, id)
	return nil
}
