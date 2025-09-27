package calculation

import (
	"ferex/backend/models"
	"testing"
)

func TestDetermineRetirementType(t *testing.T) {
	cases := []struct {
		ageYears     int
		serviceYears float64
		mraYears     int
		expect       string
	}{
		{57, 31, 56, "ImmediateMRA30"},
		{61, 21, 56, "Immediate60_20"},
		{57, 15, 56, "MRA+10Voluntary"},
		{54, 29, 56, "Other"},
	}
	for _, c := range cases {
		got := DetermineRetirementType(c.ageYears, c.serviceYears, c.mraYears, "2025-05-20", "1968-01-01")
		if got != c.expect {
			t.Errorf("age %d, svc %.1f, mra %d: want %s, got %s", c.ageYears, c.serviceYears, c.mraYears, c.expect, got)
		}
	}
}

func TestProrationFactor(t *testing.T) {
	periods := []models.ServicePeriod{
		{
			StartDate: "2000-01-01", EndDate: "2010-01-01", 
			ServiceCategory: "Civilian", 
			IsPartTime: false,
		},
		{
			StartDate: "2010-01-01", EndDate: "2015-01-01", 
			ServiceCategory: "Civilian", 
			IsPartTime: true, 
			HoursPerWeekIfPartTime: floatPtr(20),
		},
	}
	factor := ProrationFactor(periods)
	if factor < 0.82 || factor > 0.84 {
		t.Errorf("expected proration ~0.83 for 10y full-time + 5y half-time, got %f", factor)
	}

	// All part-time
	periods = []models.ServicePeriod{
		{
			StartDate: "2000-01-01", EndDate: "2010-01-01", 
			ServiceCategory: "Civilian", 
			IsPartTime: true, 
			HoursPerWeekIfPartTime: floatPtr(20),
		},
	}
	factor = ProrationFactor(periods)
	if factor < 0.49 || factor > 0.51 {
		t.Errorf("expected proration ~0.5 for all half-time, got %f", factor)
	}
}
