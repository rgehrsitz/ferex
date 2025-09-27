package calculation

import (
	"time"
	"fmt"
)

import "ferex/backend/models"

// CalculateAge returns years, months between two dates (YYYY-MM-DD)
func CalculateAge(dateOfBirth, referenceDate string) (years int, months int, err error) {
	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid dateOfBirth: %w", err)
	}
	ref, err := time.Parse("2006-01-02", referenceDate)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid referenceDate: %w", err)
	}
	if ref.Before(dob) {
		return 0, 0, fmt.Errorf("referenceDate before dateOfBirth")
	}
	years = int(ref.Year() - dob.Year())
	months = int(ref.Month() - dob.Month())
	if ref.Day() < dob.Day() {
		months--
	}
	if months < 0 {
		years--
		months += 12
	}
	return years, months, nil
}

// YearsBetween returns the number of years (including fractions) between two dates (YYYY-MM-DD)
func YearsBetween(startDate, endDate string) (float64, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid startDate: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid endDate: %w", err)
	}
	if end.Before(start) {
		return 0, fmt.Errorf("endDate before startDate")
	}
	days := end.Sub(start).Hours() / 24.0
	years := days / 365.25
	return years, nil
}

// CalculateServiceYears sums the years of service given a slice of ServicePeriod
func CalculateServiceYears(periods []models.ServicePeriod) (float64, error) {
	total := 0.0
	for _, p := range periods {
		if p.StartDate == "" || p.EndDate == "" {
			continue // skip incomplete
		}
		years, err := YearsBetween(p.StartDate, p.EndDate)
		if err != nil {
			return 0, fmt.Errorf("invalid service period: %w", err)
		}
		// Prorate for part time if needed
		if p.ServiceCategory == "Civilian" && p.IsPartTime && p.HoursPerWeekIfPartTime != nil {
			// Ensure HoursPerWeekIfPartTime is not zero to avoid division by zero, though proration would be 0 anyway.
			// Also, ensure it's less than standard full-time hours, though OPM guidance might prorate even if over (e.g. unusual tour).
			// For now, assume valid part-time hours are > 0 and < 40.
			if *p.HoursPerWeekIfPartTime > 0 && *p.HoursPerWeekIfPartTime < 40.0 { 
				proration := *p.HoursPerWeekIfPartTime / 40.0
				years *= proration
			}
		}
		total += years
	}
	return total, nil
}

// MinimumRetirementAge returns the MRA (years) for a given date of birth (FERS rules)
func MinimumRetirementAge(dateOfBirth string) (int, error) {
	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return 0, fmt.Errorf("invalid dateOfBirth: %w", err)
	}
	year := dob.Year()
	switch {
	case year <= 1947:
		return 55, nil
	case year == 1948:
		return 55, nil
	case year == 1949:
		return 55, nil
	case year == 1950:
		return 55, nil
	case year == 1951:
		return 55, nil
	case year == 1952:
		return 55, nil
	case year == 1953:
		return 55, nil
	case year == 1954:
		return 55, nil
	case year == 1955:
		return 55, nil
	case year == 1956:
		return 55, nil
	case year == 1957:
		return 55, nil
	case year == 1958:
		return 55, nil
	case year == 1959:
		return 55, nil
	case year == 1960:
		return 56, nil
	case year >= 1961 && year <= 1962:
		return 56, nil
	case year >= 1963 && year <= 1964:
		return 56, nil
	case year >= 1965 && year <= 1966:
		return 56, nil
	case year >= 1967 && year <= 1968:
		return 56, nil
	case year >= 1969 && year <= 1970:
		return 56, nil
	case year >= 1971 && year <= 1972:
		return 56, nil
	case year >= 1973 && year <= 1974:
		return 56, nil
	case year >= 1975 && year <= 1976:
		return 56, nil
	case year >= 1977 && year <= 1978:
		return 56, nil
	case year >= 1979 && year <= 1980:
		return 56, nil
	case year >= 1981 && year <= 1982:
		return 56, nil
	case year >= 1983 && year <= 1984:
		return 56, nil
	case year >= 1985 && year <= 1986:
		return 56, nil
	case year >= 1987 && year <= 1988:
		return 56, nil
	case year >= 1989 && year <= 1990:
		return 56, nil
	case year >= 1991 && year <= 1992:
		return 56, nil
	case year >= 1993 && year <= 1994:
		return 56, nil
	case year >= 1995 && year <= 1996:
		return 56, nil
	case year >= 1997 && year <= 1998:
		return 56, nil
	case year >= 1999 && year <= 2000:
		return 56, nil
	default:
		return 57, nil // fallback
	}
}

// CalculateAgeAtDate calculates the age in years, months, and days at a specific date.
// It is a placeholder and needs to be implemented.
func CalculateAgeAtDate(birthDate time.Time, currentProjectionTime time.Time) (ageYears int, ageMonths int, ageDays int, err error) {
	// Basic placeholder logic, to be replaced with accurate calculation
	if currentProjectionTime.Before(birthDate) {
		return 0, 0, 0, fmt.Errorf("currentProjectionTime cannot be before birthDate")
	}

	ageYears = currentProjectionTime.Year() - birthDate.Year()
	ageMonths = int(currentProjectionTime.Month()) - int(birthDate.Month())
	ageDays = currentProjectionTime.Day() - birthDate.Day()

	if ageDays < 0 {
		// Borrow a month
		ageMonths--
		// Approximate days in previous month - this is a simplification!
		// For a real implementation, use date library functions to get days in month.
		t := currentProjectionTime.AddDate(0, -1, 0)
		daysInLastMonth := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		ageDays += daysInLastMonth
	}

	if ageMonths < 0 {
		ageYears--
		ageMonths += 12
	}

	return ageYears, ageMonths, ageDays, nil
}
