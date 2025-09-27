package core

import (
	"strconv"
	"time"
)

// ServicePeriod captures a span of employment or military service that counts
// toward retirement eligibility and/or computation. Additional metadata (like
// whether a deposit was paid) supports accurate pension modeling and SRS rules.
type ServicePeriod struct {
	StartDate    time.Time   `json:"start_date"`
	EndDate      time.Time   `json:"end_date"`
	Type         ServiceType `json:"type"`
	HoursPerWeek float64     `json:"hours_per_week"`
	DepositPaid  bool        `json:"deposit_paid"`
	Notes        string      `json:"notes,omitempty"`
}

// CompensationEntry stores a period of basic pay used for the high-3
// calculation as well as optional overtime/bonus data for planning context.
type CompensationEntry struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	PayType   PayType   `json:"pay_type"`
	Amount    float64   `json:"amount"`
}

// SocialSecurityRecord holds the assumptions required to estimate SSA benefits
// at various claiming ages.
type SocialSecurityRecord struct {
	Age62Estimate           float64   `json:"age_62_estimate"`
	FullRetirementAgeMonths int       `json:"fra_months"`
	EarningsHistory         []float64 `json:"earnings_history,omitempty"`
	LastUpdated             time.Time `json:"last_updated"`
}

// UserProfile represents the foundational data that seeds scenarios. It is a
// master record that aligns with the "My Profile" screen in the UI vision.
type UserProfile struct {
	ID               string               `json:"id"`
	DisplayName      string               `json:"display_name"`
	BirthDate        time.Time            `json:"birth_date"`
	HireDate         time.Time            `json:"hire_date"`
	RetirementSystem RetirementSystem     `json:"retirement_system"`
	ServiceHistory   []ServicePeriod      `json:"service_history"`
	Compensation     []CompensationEntry  `json:"compensation"`
	SocialSecurity   SocialSecurityRecord `json:"social_security"`
	Notes            string               `json:"notes,omitempty"`
}

// Clone produces a deep copy so callers can manipulate data without mutating
// the original profile reference.
func (p UserProfile) Clone() UserProfile {
	cloned := p
	cloned.ServiceHistory = append([]ServicePeriod(nil), p.ServiceHistory...)
	cloned.Compensation = append([]CompensationEntry(nil), p.Compensation...)
	cloned.SocialSecurity.EarningsHistory = append([]float64(nil), p.SocialSecurity.EarningsHistory...)
	return cloned
}

// Validate performs basic consistency checks aligned with the legacy
// application rules (e.g., IDs required, service dates ordered).
func (p UserProfile) Validate() *ValidationError {
	var verr ValidationError
	if p.ID == "" {
		verr.Add("id", "profile id is required")
	}
	if p.DisplayName == "" {
		verr.Add("display_name", "display name is required")
	}
	if err := p.RetirementSystem.Validate(); err != nil {
		verr.Add("retirement_system", err.Error())
	}
	for i, period := range p.ServiceHistory {
		if period.EndDate.Before(period.StartDate) {
			verr.Add("service_history", "period "+period.StartDate.Format(time.DateOnly)+" has an end date before start")
		}
		if period.Type == "" {
			verr.Add("service_history", "period index "+strconv.Itoa(i)+" missing type")
		}
	}
	for i, comp := range p.Compensation {
		if comp.EndDate.Before(comp.StartDate) {
			verr.Add("compensation", "entry index "+strconv.Itoa(i)+" end date before start")
		}
		if comp.Amount <= 0 {
			verr.Add("compensation", "entry index "+strconv.Itoa(i)+" amount must be positive")
		}
	}
	if verr.Empty() {
		return nil
	}
	return &verr
}
