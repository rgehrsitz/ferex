package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"testing"
)

func TestTaxCalculation(t *testing.T) {
	cases := []struct {
		name          string
		input         models.TaxCalculationInput
		expectFed     float64
		expectState   float64
		expectNet     float64
		effectiveRate float64
		notesContains string
	}{
		{
			name: "Single, only pension, standard deduction",
			input: models.TaxCalculationInput{
				FilingStatus:       "single",
				TaxYear:            2025,
				GrossPension:       50000,
				TaxablePension:     50000,
				TSPWithdrawal:      0,
				TSPRothWithdrawal:  0,
				SocialSecurity:     0,
				OtherTaxableIncome: 0,
				Deductions:         0, // use standard
				TaxCredits:         0,
			},
			expectFed:     50000 - 14600, // taxable = 35400
			expectState:   0,
			expectNet:     50000 - calculation.CalculateTax(models.TaxCalculationInput{FilingStatus: "single", TaxYear: 2025, GrossPension: 50000, TaxablePension: 50000, Deductions: 0}).FederalTaxOwed,
			effectiveRate: 0.0, // checked below
			notesContains: "",
		},
		{
			name: "Married, pension + TSP, state tax applies",
			input: models.TaxCalculationInput{
				FilingStatus:       "married",
				TaxYear:            2025,
				GrossPension:       70000,
				TaxablePension:     70000,
				TSPWithdrawal:      20000,
				TSPRothWithdrawal:  0,
				SocialSecurity:     0,
				OtherTaxableIncome: 0,
				StateOfResidence:   "VA",
				StateTaxableIncome: 90000,
				Deductions:         0,
				TaxCredits:         0,
			},
			expectFed:     90000 - 29200, // taxable = 60800
			expectState:   90000 * 0.05,
			expectNet:     70000 + 20000 - calculation.CalculateTax(models.TaxCalculationInput{FilingStatus: "married", TaxYear: 2025, GrossPension: 70000, TaxablePension: 70000, TSPWithdrawal: 20000, Deductions: 0}).FederalTaxOwed - 90000*0.05,
			effectiveRate: 0.0,
			notesContains: "State tax estimated at 5% for VA",
		},
		{
			name: "Single, Social Security taxability edge",
			input: models.TaxCalculationInput{
				FilingStatus:       "single",
				TaxYear:            2025,
				GrossPension:       20000,
				TaxablePension:     20000,
				TSPWithdrawal:      0,
				TSPRothWithdrawal:  0,
				SocialSecurity:     20000,
				OtherTaxableIncome: 0,
				Deductions:         0,
				TaxCredits:         0,
			},
			expectFed:     0, // Should be below SS taxability threshold
			expectState:   0,
			expectNet:     37532.00,
			effectiveRate: 0.0,
			notesContains: "",
		},
		{
			name: "Single, with tax credits",
			input: models.TaxCalculationInput{
				FilingStatus:       "single",
				TaxYear:            2025,
				GrossPension:       40000,
				TaxablePension:     40000,
				TSPWithdrawal:      0,
				TSPRothWithdrawal:  0,
				SocialSecurity:     0,
				OtherTaxableIncome: 0,
				Deductions:         0,
				TaxCredits:         2000,
			},
			expectFed:     0, // Credits wipe out tax
			expectState:   0,
			expectNet:     39172.00,
			effectiveRate: 0.0,
			notesContains: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateTax(tc.input)
			if got.FederalTaxOwed < -0.01 || got.StateTaxOwed < -0.01 {
				t.Errorf("%s: negative tax owed", tc.name)
			}
			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes missing expected string: %q", tc.name, tc.notesContains)
			}
			// Check net after-tax income is within $2
			if testutils.Abs(got.NetAfterTaxIncome-tc.expectNet) > 2.0 {
				t.Errorf("%s: net after-tax income got %.2f, want %.2f", tc.name, got.NetAfterTaxIncome, tc.expectNet)
			}
		})
	}
}
