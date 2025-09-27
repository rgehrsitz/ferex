package calculation

import (
	"testing"
	"ferex/backend/models"
)

func TestCalculateAge(t *testing.T) {
	y, m, err := CalculateAge("1960-06-15", "2025-05-20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if y != 64 || m != 11 {
		t.Errorf("expected 64y 11m, got %dy %dm", y, m)
	}
}

func TestYearsBetween(t *testing.T) {
	years, err := YearsBetween("2000-01-01", "2025-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if years < 24.99 || years > 25.01 {
		t.Errorf("expected ~25 years, got %f", years)
	}
}

func TestCalculateServiceYears_FullAndPartTime(t *testing.T) {
	periods := []models.ServicePeriod{
		{
			StartDate: "2000-01-01", EndDate: "2010-01-01", 
			ServiceCategory: "Civilian", 
			IsPartTime: false, 
			// Other fields like CivilianServiceType, DepositRedepositPaymentStatus, SystemDuringService can be added if needed for test completeness
			// but are not strictly required for this specific function's logic being tested.
		},
		{
			StartDate: "2010-01-01", EndDate: "2015-01-01", 
			ServiceCategory: "Civilian", 
			IsPartTime: true, 
			HoursPerWeekIfPartTime: floatPtr(20),
		},
	}
	total, err := CalculateServiceYears(periods)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// First period: 10y, second: 5y at 0.5 = 2.5y, total 12.5y
	if total < 12.49 || total > 12.51 {
		t.Errorf("expected ~12.5 years, got %f", total)
	}
}

func TestMinimumRetirementAge(t *testing.T) {
	mra, err := MinimumRetirementAge("1962-07-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mra != 56 {
		t.Errorf("expected MRA 56, got %d", mra)
	}
}

func floatPtr(f float64) *float64 { return &f }
