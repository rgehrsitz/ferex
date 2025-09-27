package config

import (
	"encoding/json"
	"ferex/backend/models"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// FerexConfig holds all configurable Ferex parameters
type FerexConfig struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	
	TSP        TSPConfig        `json:"tsp"`
	FederalTax FederalTaxConfig `json:"federalTax"`
	StateTax   StateTaxConfig   `json:"stateTax"`
}

// TSPConfig holds all configurable TSP parameters
type TSPConfig struct {
	// Contribution limits
	ContributionLimits struct {
		MaxAgencyMatchPercentage              float64 `json:"maxAgencyMatchPercentage"`              // Max 4% matching
		AgencyAutomaticContributionPercentage float64 `json:"agencyAutomaticContributionPercentage"` // Automatic 1% agency contribution
		StandardEmployeeContributionForMaxMatch float64 `json:"standardEmployeeContributionForMaxMatch"` // Employee needs to contribute 5% to get full 4% match
		CatchUpContributionLimit              float64 `json:"catchUpContributionLimit"`              // Annual catch-up contribution limit
		MaxRegularContributionLimit           float64 `json:"maxRegularContributionLimit"`           // Annual regular contribution limit
		Year                                  int     `json:"year"`                                  // Year these limits apply to
	} `json:"contributionLimits"`

	// L-Fund compositions
	LFundBaseCompositions map[string]models.TSPFundAllocationPercentages `json:"lFundBaseCompositions"` // Initial allocations for L-Funds
	LFundTargetComposition models.TSPFundAllocationPercentages          `json:"lFundTargetComposition"` // Target allocation when L-Funds reach target date

	// RMD parameters
	UniformLifetimeTable map[string]float64 `json:"uniformLifetimeTable"` // RMD factors by age
	RMDStartAgeRules     []struct {
		BirthYearStart int `json:"birthYearStart"` // Inclusive
		BirthYearEnd   int `json:"birthYearEnd"`   // Inclusive
		StartAge       int `json:"startAge"`       // RMD start age for this birth year range
	} `json:"rmdStartAgeRules"`
}

// FederalTaxConfig holds federal tax configuration
type FederalTaxConfig struct {
	TaxYear                       int                                    `json:"taxYear"`
	Brackets                     map[string][]TaxBracket                `json:"brackets"`
	StandardDeduction            map[string]float64                     `json:"standardDeduction"`
	AdditionalStandardDeduction  AdditionalStandardDeductionConfig      `json:"additionalStandardDeduction"`
	SocialSecurityTaxabilityThresholds map[string]SSTaxabilityThreshold `json:"socialSecurityTaxabilityThresholds"`
}

// TaxBracket represents a single tax bracket
type TaxBracket struct {
	Rate float64 `json:"rate"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

// AdditionalStandardDeductionConfig holds additional deduction amounts
type AdditionalStandardDeductionConfig struct {
	Age65OrOlder map[string]float64 `json:"age65OrOlder"`
	Blind        map[string]float64 `json:"blind"`
}

// SSTaxabilityThreshold holds Social Security taxability thresholds
type SSTaxabilityThreshold struct {
	Tier1 float64 `json:"tier1"`
	Tier2 float64 `json:"tier2"`
}

// StateTaxConfig holds state tax configuration
type StateTaxConfig struct {
	States map[string]StateInfo `json:"states"`
}

// StateInfo holds tax information for a specific state
type StateInfo struct {
	Abbreviation       string  `json:"abbreviation"`
	HasStateTax        bool    `json:"hasStateTax"`
	FederalPensionTax  string  `json:"federalPensionTax"`  // "exempt", "taxed", "partial"
	TSPTax             string  `json:"tspTax"`             // "exempt", "taxed", "partial"
	SocialSecurityTax  string  `json:"socialSecurityTax"`  // "exempt", "taxed", "partial"
	EffectiveRate      float64 `json:"effectiveRate"`
	Notes              string  `json:"notes"`
}

var (
	ferexConfig     *FerexConfig
	ferexConfigOnce sync.Once
	ferexConfigLock sync.RWMutex
)

// GetFerexConfig returns the Ferex configuration, loading it from disk if necessary
func GetFerexConfig() (*FerexConfig, error) {
	var initErr error
	ferexConfigOnce.Do(func() {
		ferexConfig, initErr = loadFerexConfig()
	})

	if initErr != nil {
		return nil, initErr
	}

	ferexConfigLock.RLock()
	defer ferexConfigLock.RUnlock()
	return ferexConfig, nil
}

// ReloadFerexConfig forces a reload of the Ferex configuration from disk
func ReloadFerexConfig() (*FerexConfig, error) {
	ferexConfigLock.Lock()
	defer ferexConfigLock.Unlock()

	var err error
	ferexConfig, err = loadFerexConfig()
	return ferexConfig, err
}

// loadFerexConfig loads the Ferex configuration from the config file, trying multiple standard locations.
func loadFerexConfig() (*FerexConfig, error) {
	pathsToTry := []string{
		filepath.Join("config", "ferex_config.json"),           // project root
		filepath.Join("..", "config", "ferex_config.json"),     // when running from backend/
		filepath.Join("..", "..", "config", "ferex_config.json"), // when running from backend/calculation/
		filepath.Join(".", "ferex_config.json"),                // current dir (fallback)
	}

	var lastErr error
	for _, configPath := range pathsToTry {
		if _, err := os.Stat(configPath); err == nil {
			// Found the file, try to load it
			data, err := ioutil.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read Ferex config at %s: %w", configPath, err)
			}
			var cfg FerexConfig
			if err := json.Unmarshal(data, &cfg); err != nil {
				return nil, fmt.Errorf("failed to parse Ferex config at %s: %w", configPath, err)
			}
			fmt.Printf("Loaded Ferex config from: %s\n", configPath)
			return &cfg, nil
		} else {
			lastErr = err
		}
	}
	return nil, fmt.Errorf("Ferex configuration file not found in any known location. Last error: %v", lastErr)
}

// SaveFerexConfig saves the Ferex configuration to the config file
func SaveFerexConfig(config *FerexConfig) error {
	ferexConfigLock.Lock()
	defer ferexConfigLock.Unlock()

	configPath := filepath.Join("config", "ferex_config.json")

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling Ferex configuration: %v", err)
	}

	if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing Ferex configuration file: %v", err)
	}

	ferexConfig = config
	return nil
}

// GetUniformLifetimeTableAsMap converts the string keys in the JSON to integer keys for use in the application
func (c *TSPConfig) GetUniformLifetimeTableAsMap() map[int]float64 {
	result := make(map[int]float64)
	for ageStr, factor := range c.UniformLifetimeTable {
		var age int
		if _, err := fmt.Sscanf(ageStr, "%d", &age); err == nil {
			result[age] = factor
		}
	}
	return result
}

// GetRMDStartAge returns the RMD start age for a given birth year
func (c *TSPConfig) GetRMDStartAge(birthYear int) int {
	for _, rule := range c.RMDStartAgeRules {
		if birthYear >= rule.BirthYearStart && birthYear <= rule.BirthYearEnd {
			return rule.StartAge
		}
	}
	// Default to the highest age if no rule matches (conservative approach)
	highestAge := 0
	for _, rule := range c.RMDStartAgeRules {
		if rule.StartAge > highestAge {
			highestAge = rule.StartAge
		}
	}
	return highestAge
}

// Legacy TSP functions for backward compatibility
func GetTSPConfig() (*TSPConfig, error) {
	cfg, err := GetFerexConfig()
	if err != nil {
		return nil, err
	}
	return &cfg.TSP, nil
}

func ReloadTSPConfig() (*TSPConfig, error) {
	cfg, err := ReloadFerexConfig()
	if err != nil {
		return nil, err
	}
	return &cfg.TSP, nil
}