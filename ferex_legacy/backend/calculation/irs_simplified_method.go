package calculation

import "math"

// CalculateIRSSimplifiedMethodExclusion computes the annual tax-free portion of a pension
// using the IRS Simplified Method. See IRS Pub 575 Table 1 for number of expected payments.
// Inputs:
// - totalContributions: the employee's after-tax contributions
// - annuitantAge: age at annuity start (retirement)
// - hasSurvivor: whether a survivor benefit is elected
// Returns: the annual exclusion amount.
func CalculateIRSSimplifiedMethodExclusion(totalContributions float64, annuitantAge int, hasSurvivor bool) float64 {
	numPayments := getExpectedPaymentsIRS(annuitantAge, hasSurvivor)
	if numPayments == 0 {
		return 0
	}
	return math.Round((totalContributions/float64(numPayments))*100) / 100
}

// getExpectedPaymentsIRS returns the number of expected monthly payments per IRS Table 1.
func getExpectedPaymentsIRS(age int, hasSurvivor bool) int {
	// Table 1 from IRS Pub 575 (2024):
	// https://www.irs.gov/pub/irs-pdf/p575.pdf
	// This is a simplified version for common ages; extend as needed.
	if hasSurvivor {
		switch {
		case age >= 70:
			return 310
		case age >= 66:
			return 320
		case age >= 61:
			return 340
		case age >= 55:
			return 360
		case age >= 50:
			return 380
		default:
			return 410
		}
	} else {
		switch {
		case age >= 70:
			return 140
		case age >= 66:
			return 150
		case age >= 61:
			return 170
		case age >= 55:
			return 190
		case age >= 50:
			return 210
		default:
			return 230
		}
	}
}
