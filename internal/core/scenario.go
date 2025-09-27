package core

import "time"

// IncomeStream represents an arbitrary source of income that feeds analysis
// reports. Pension, SRS, Social Security, and TSP withdrawals will all be
// normalized into this structure for reporting consistency.
type IncomeStream struct {
	Kind      CashflowKind     `json:"kind"`
	Name      string           `json:"name"`
	Amount    float64          `json:"amount"`
	Frequency PaymentFrequency `json:"frequency"`
	StartDate time.Time        `json:"start_date"`
	EndDate   time.Time        `json:"end_date"`
	Taxable   bool             `json:"taxable"`
}

// Scenario encapsulates a retirement plan variant. It embeds the normalized
// profile data so calculations operate on a self-contained aggregate while
// retaining a reference back to the source profile for auditing.
type Scenario struct {
	ID                     string           `json:"id"`
	Name                   string           `json:"name"`
	ProfileRef             string           `json:"profile_ref"`
	ProfileSnapshot        UserProfile      `json:"profile_snapshot"`
	RetirementDate         time.Time        `json:"retirement_date"`
	SurvivorPercent        float64          `json:"survivor_percent"`
	IncludeSupplement      bool             `json:"include_supplement"`
	PostRetirementEarnings float64          `json:"post_retirement_earnings"`
	TSPPlan                TSPPlan          `json:"tsp_plan"`
	TaxSettings            TaxSettings      `json:"tax_settings"`
	OtherIncome            []IncomeStream   `json:"other_income"`
	MonteCarlo             SimulationConfig `json:"monte_carlo"`
	CreatedAt              time.Time        `json:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at"`
	Notes                  string           `json:"notes,omitempty"`
}

// ScenarioMetadata exposes lightweight fields for listings without pulling the
// full snapshot payload.
type ScenarioMetadata struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	ProfileRef     string    `json:"profile_ref"`
	UpdatedAt      time.Time `json:"updated_at"`
	RetirementDate time.Time `json:"retirement_date"`
}

// Validate checks that the scenario is internally consistent before running
// calculations or persisting to disk.
func (s Scenario) Validate() *ValidationError {
	var verr ValidationError
	if s.ID == "" {
		verr.Add("id", "scenario id is required")
	}
	if s.Name == "" {
		verr.Add("name", "scenario name is required")
	}
	if s.ProfileRef == "" {
		verr.Add("profile_ref", "profile reference is required")
	}
	if !s.ProfileSnapshot.BirthDate.IsZero() && s.RetirementDate.Before(s.ProfileSnapshot.BirthDate) {
		verr.Add("retirement_date", "retirement cannot be before birth date")
	}
	if profileErr := s.ProfileSnapshot.Validate(); profileErr != nil {
		verr.Merge(profileErr)
	}
	if tspErr := s.TSPPlan.Validate(); tspErr != nil {
		verr.Merge(tspErr)
	}
	if taxErr := s.TaxSettings.Validate(); taxErr != nil {
		verr.Merge(taxErr)
	}
	if verr.Empty() {
		return nil
	}
	return &verr
}
