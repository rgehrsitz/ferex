package file

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/ferexapp/ferex/internal/core"
)

// ScenarioStore persists scenarios as JSON files on disk with a .ferex
// extension. Each scenario is stored independently to simplify version control
// and manual editing.
type ScenarioStore struct {
	root string
	mu   sync.Mutex
}

// NewScenarioStore ensures the root directory exists and returns a file-backed
// store.
func NewScenarioStore(root string) (*ScenarioStore, error) {
	if root == "" {
		return nil, errors.New("root path is required")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &ScenarioStore{root: root}, nil
}

func (s *ScenarioStore) filePath(id string) string {
	return filepath.Join(s.root, id+".ferex")
}

// Save writes the scenario to disk as indented JSON.
func (s *ScenarioStore) Save(_ context.Context, scenario core.Scenario) error {
	if s == nil {
		return errors.New("file store is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.CreateTemp(s.root, scenario.ID+"-*.tmp")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&scenario); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(f.Name())
		return err
	}
	return os.Rename(f.Name(), s.filePath(scenario.ID))
}

// Get loads and unmarshals the scenario file.
func (s *ScenarioStore) Get(_ context.Context, id string) (core.Scenario, error) {
	if s == nil {
		return core.Scenario{}, errors.New("file store is nil")
	}
	path := s.filePath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return core.Scenario{}, errors.New("scenario not found")
		}
		return core.Scenario{}, err
	}
	var scenario core.Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return core.Scenario{}, err
	}
	return scenario, nil
}

// List enumerates all scenario metadata.
func (s *ScenarioStore) List(_ context.Context) ([]core.ScenarioMetadata, error) {
	if s == nil {
		return nil, errors.New("file store is nil")
	}
	entries, err := os.ReadDir(s.root)
	if err != nil {
		return nil, err
	}
	var items []core.ScenarioMetadata
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".ferex" {
			continue
		}
		path := filepath.Join(s.root, entry.Name())
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		decoder := json.NewDecoder(f)
		var scenario core.Scenario
		if err := decoder.Decode(&scenario); err != nil && err != io.EOF {
			f.Close()
			return nil, err
		}
		f.Close()
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

// Delete removes the scenario file.
func (s *ScenarioStore) Delete(_ context.Context, id string) error {
	if s == nil {
		return errors.New("file store is nil")
	}
	path := s.filePath(id)
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
