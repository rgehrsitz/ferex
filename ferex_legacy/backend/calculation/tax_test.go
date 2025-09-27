package calculation // Changed from 'tests' to 'calculation'

import (
	// "ferex/backend/calculation" // Removed as test is in the same package
	"ferex/backend/models"         // Assuming 'ferex' is the module name
	"ferex/backend/testutils"    // Assuming 'ferex' is the module name
	"testing"
)

func TestTaxCalculation(t *testing.T) {
	cases := []struct {
		name          string
		input         models.TaxCalculationInput
		expectFed     float64 // Note: expectFed in legacy test was comparing against taxable income, not actual tax.
		expectState   float64 // This will be re-evaluated.
		expectNet     float64
		effectiveRate float64
		notesContains string
	}{
		{
			name: "Single, only pension, standard deduction",
			input: models.TaxCalculationInput{
				FilingStatus:       "Single",
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
			// expectFed needs to be the actual tax. For Taxable income 35400 (50000-14600):
			// 11000*0.10 = 1100
			// (35400-11000)*0.12 = 24400*0.12 = 2928
			// Total = 1100 + 2928 = 4028
			expectFed:     4028, 
			expectState:   0,
			expectNet:     50000 - 4028, // 45972
			effectiveRate: float64(4028) / float64(50000),
			notesContains: "",
		},
		{
			name: "Married, pension + TSP, state tax applies",
			input: models.TaxCalculationInput{
				FilingStatus:       "MarriedFilingJointly",
				TaxYear:            2025,
				GrossPension:       70000,
				TaxablePension:     70000,
				TSPWithdrawal:      20000,
				TSPRothWithdrawal:  0,
				SocialSecurity:     0,
				OtherTaxableIncome: 0,
				StateOfResidence:   "Virginia",
				StateTaxableIncome: 90000, // Assuming total income is state taxable for simplicity here
				Deductions:         0,
				TaxCredits:         0,
			},
			// Taxable income: 90000 (Gross) - 29200 (Std Ded) = 60800
			// 22000*0.10 = 2200
			// (60800-22000)*0.12 = 38800*0.12 = 4656
			// Total Fed = 2200 + 4656 = 6856
			expectFed:     6856,
			expectState:   90000 * 0.0575, // 5175 (Virginia rate)
			expectNet:     90000 - 6856 - 5175, // 77969
			effectiveRate: (float64(6856) + float64(5175)) / float64(90000),
			notesContains: "State tax calculated for Virginia",
		},
		{
			name: "Single, Social Security taxability threshold (AGI 20k, SS 20k -> Prov 30k -> 50% SS taxable)",
			input: models.TaxCalculationInput{
				FilingStatus:       "Single",
				TaxYear:            2025,
				GrossPension:       20000, // Part of AGI
				TaxablePension:     20000,
				TSPWithdrawal:      0,
				TSPRothWithdrawal:  0,
				SocialSecurity:     20000,
				OtherTaxableIncome: 0,
				Deductions:         0,
				TaxCredits:         0,
			},
			// AGI = 20000. Provisional = 20000 + 0.5*20000 = 30000.
			// For single, if provisional > 25k, min(50% of SS, 50% of (provisional-25k))
			// SS Taxable = min(0.5*20000, 0.5*(30000-25000)) = min(10000, 2500) = 2500
			// Total Taxable Income = AGI + SS_Taxable - Deduction = 20000 + 2500 - 14600 = 7900
			// Fed Tax: 7900 * 0.10 = 790 (all in 10% bracket)
			expectFed:     790,
			expectState:   0,
			// Total Gross Income = 20000 (pension) + 20000 (SS) = 40000
			expectNet:     40000 - 790, // 39210
			effectiveRate: float64(790) / float64(40000),
			notesContains: "",
		},
		{
			name: "Single, with tax credits enough to wipe out tax",
			input: models.TaxCalculationInput{
				FilingStatus:       "Single",
				TaxYear:            2025,
				GrossPension:       40000,
				TaxablePension:     40000,
				TSPWithdrawal:      0,
				TSPRothWithdrawal:  0,
				SocialSecurity:     0,
				OtherTaxableIncome: 0,
				Deductions:         0, // std deduction = 14600
				TaxCredits:         2000,
			},
			// Taxable Income = 40000 - 14600 = 25400
			// Fed Tax before credits: 
			// 11000 * 0.10 = 1100
			// (25400-11000)*0.12 = 14400 * 0.12 = 1728
			// Total = 1100 + 1728 = 2828
			// Fed Tax after credits = 2828 - 2000 = 828
			expectFed:     828, 
			expectState:   0,
			expectNet:     40000 - 828, // 39172
			effectiveRate: float64(828) / float64(40000),
			notesContains: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateTax(tc.input) // Direct call as it's in the same package

			if testutils.Abs(got.FederalTaxOwed-tc.expectFed) > 0.01 {
				t.Errorf("%s: FederalTaxOwed got %.2f, want %.2f", tc.name, got.FederalTaxOwed, tc.expectFed)
			}
			if testutils.Abs(got.StateTaxOwed-tc.expectState) > 0.01 {
				t.Errorf("%s: StateTaxOwed got %.2f, want %.2f", tc.name, got.StateTaxOwed, tc.expectState)
			}
			if testutils.Abs(got.NetAfterTaxIncome-tc.expectNet) > 0.01 {
				t.Errorf("%s: NetAfterTaxIncome got %.2f, want %.2f", tc.name, got.NetAfterTaxIncome, tc.expectNet)
			}
			if totalIncome := tc.input.GrossPension + tc.input.TSPWithdrawal + tc.input.TSPRothWithdrawal + tc.input.SocialSecurity + tc.input.OtherTaxableIncome; totalIncome > 0 {
				expectedEffectiveRate := (tc.expectFed + tc.expectState) / totalIncome
				if testutils.Abs(got.EffectiveTaxRate-expectedEffectiveRate) > 0.0001 { // Effective rate comparison
					t.Errorf("%s: EffectiveTaxRate got %.4f, want %.4f", tc.name, got.EffectiveTaxRate, expectedEffectiveRate)
				}
			} else if got.EffectiveTaxRate != 0 { // Handle zero total income case
				t.Errorf("%s: EffectiveTaxRate got %.4f, want 0.0000 for zero total income", tc.name, got.EffectiveTaxRate)
			}

			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes missing expected string: %q, got: %q", tc.name, tc.notesContains, got.Notes)
			}
		})
	}
}
