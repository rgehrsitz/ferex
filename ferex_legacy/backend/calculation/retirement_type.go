package calculation

import (
	"ferex/backend/models"
)

// DetermineRetirementType determines the FERS retirement type based on age, service years, and MRA
// Returns one of: "ImmediateMRA30", "Immediate60_20", "MRA+10Voluntary", "Deferred", "Disability", "Other"
func DetermineRetirementType(ageYears int, serviceYears float64, mraYears int, plannedRetirementDate, dateOfBirth string) string {
	if ageYears >= mraYears && serviceYears >= 30 {
		return "ImmediateMRA30"
	}
	if ageYears >= 60 && serviceYears >= 20 {
		return "Immediate60_20"
	}
	if ageYears >= mraYears && serviceYears >= 10 && serviceYears < 30 {
		return "MRA+10Voluntary"
	}
	// TODO: Add logic for Deferred, Disability, VERA, etc. as needed
	return "Other"
}

// ProrationFactor returns the weighted average proration factor across all service periods
// If all service is full-time, returns 1.0
func ProrationFactor(periods []models.ServicePeriod) float64 {
	totalYears := 0.0
	weightedProration := 0.0
	
	for _, p := range periods {
		years, err := YearsBetween(p.StartDate, p.EndDate)
		if err != nil {
			continue
		}
		
		var proration float64
		if p.ServiceCategory == "Civilian" && p.IsPartTime && p.HoursPerWeekIfPartTime != nil {
			// Assume valid part-time hours are > 0 and < 40.0 for proration.
			if *p.HoursPerWeekIfPartTime > 0 && *p.HoursPerWeekIfPartTime < 40.0 {
				proration = *p.HoursPerWeekIfPartTime / 40.0
			} else {
				// If HoursPerWeekIfPartTime is 0, 40, or nil (despite IsPartTime being true), treat as full-time or non-prorated.
				// This case might need more specific handling based on business rules for invalid part-time data.
				proration = 1.0 
			}
		} else {
			proration = 1.0 // Full-time or non-Civilian part-time
		}
		
		totalYears += years
		weightedProration += years * proration
	}
	
	if totalYears == 0 {
		return 1.0
	}
	
	return weightedProration / totalYears
}
