package calculation

import (
	"ferex/backend/models"
	"math"
	"testing"
)

// TestRealWorldFederalEmployeeScenarios tests comprehensive scenarios based on actual federal employee situations
func TestRealWorldFederalEmployeeScenarios(t *testing.T) {
	tests := []struct {
		name              string
		description       string
		input             models.FERSCalculationInput
		ageAtRetirement   int
		totalServiceYears float64
		mra               int
		retirementType    string
		expectedSSBenefit float64
		fersServiceYears  float64
		expectedResults   expectedFERSResults
		tolerance         float64
	}{
		{
			name:        "Typical GS-13 Employee - Normal Retirement",
			description: "GS-13 Step 10 employee retiring at 62 with exactly 20 years, should get 1.1% multiplier and supplement",
			input: models.FERSCalculationInput{
				High3Salary:             95000.0, // Typical GS-13 Step 10 with locality
				UnusedSickLeaveHours:    1000,    // About 6 months of sick leave
				SurvivorBenefitElection: "Full50Percent",
				ExpectedSSBenefitAt62:   2100.0,
				EmployeeContributions:   14250.0, // 0.8% * $95k * 18.75 years (approximation)
			},
			ageAtRetirement:   62,
			totalServiceYears: 20.0,
			mra:               57,
			retirementType:    "Immediate6220",
			expectedSSBenefit: 2100.0,
			fersServiceYears:  20.0, expectedResults: expectedFERSResults{
				monthlyPension:               1605.05, // $95k * 20.48 * 1.1% / 12, minus survivor reduction (includes sick leave)
				annualPension:                19260.65,
				isEligibleForSupplement:      false,  // Age 62+ not supplement eligible
				fersSupplement:               0.0,    // No supplement at age 62+
				survivorBenefitReduction:     178.34, // 10% of unreduced pension with sick leave
				prorationApplied:             false,
				irsSimplifiedMethodExclusion: 59.38, // Approximate based on age 62 table
			},
			tolerance: 5.0,
		},
		{
			name:        "Early Career FERS - MRA+30 Retirement",
			description: "Employee who started young, retiring at MRA with 30+ years - maximum supplement eligibility",
			input: models.FERSCalculationInput{
				High3Salary:             75000.0,
				UnusedSickLeaveHours:    2000, // Full sick leave balance
				SurvivorBenefitElection: "Partial25Percent",
				ExpectedSSBenefitAt62:   1800.0,
				EmployeeContributions:   18000.0, // Higher due to longer service
			},
			ageAtRetirement:   57,
			totalServiceYears: 32.0,
			mra:               57,
			retirementType:    "ImmediateMRA30",
			expectedSSBenefit: 1800.0,
			fersServiceYears:  32.0, expectedResults: expectedFERSResults{
				monthlyPension:           1956.90, // $75k * 32.96 * 1.0% / 12, minus 5% survivor reduction (includes sick leave)
				annualPension:            23482.80,
				isEligibleForSupplement:  true,
				fersSupplement:           1440.0, // $1800 * (32/40)
				survivorBenefitReduction: 102.99, // 5% of unreduced pension with sick leave
				prorationApplied:         false,
				sickLeaveCredit:          0.96, // 2000 hours / 2087
			},
			tolerance: 50.0,
		},
		{
			name:        "CSRS Transferee - Dual Component",
			description: "Employee who transferred from CSRS to FERS mid-career",
			input: models.FERSCalculationInput{
				High3Salary:             85000.0,
				UnusedSickLeaveHours:    1500,
				SurvivorBenefitElection: "Full50Percent",
				ExpectedSSBenefitAt62:   1900.0,
				EmployeeContributions:   8500.0,       // Lower due to shorter FERS service
				SwitchedToFERSDate:      "1995-01-01", // Switched mid-career
			}, ageAtRetirement: 62,
			totalServiceYears: 25.0, // Total service
			mra:               57,
			retirementType:    "Immediate6220",
			expectedSSBenefit: 1900.0,
			fersServiceYears:  15.0, // Only FERS portion
			expectedResults: expectedFERSResults{
				monthlyPension:           1102.28, // FERS portion only: $85k * 15.72 * 1.1% / 12, minus survivor reduction
				annualPension:            13227.32,
				isEligibleForSupplement:  false,  // Age 62+ not supplement eligible
				fersSupplement:           0.0,    // No supplement at age 62+
				survivorBenefitReduction: 122.48, // 10% of FERS unreduced pension
				prorationApplied:         false,
			},
			tolerance: 100.0, // Higher tolerance due to complexity
		},
		{
			name:        "Part-Time Employee with Proration",
			description: "Employee with mixed full-time and part-time service requiring proration",
			input: models.FERSCalculationInput{
				High3Salary:             60000.0, // Already adjusted for proration in input
				UnusedSickLeaveHours:    800,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1600.0,
				EmployeeContributions:   7200.0,
			},
			ageAtRetirement:   62,
			totalServiceYears: 20.0,
			mra:               57,
			retirementType:    "Immediate6220",
			expectedSSBenefit: 1600.0,
			fersServiceYears:  20.0, expectedResults: expectedFERSResults{
				monthlyPension:           1100.0, // Reduced due to proration
				annualPension:            13200.0,
				isEligibleForSupplement:  false, // Age 62+ not supplement eligible
				fersSupplement:           0.0,   // No supplement at age 62+
				survivorBenefitReduction: 0.0,   // No survivor benefit
				prorationApplied:         true,
				prorationFactor:          0.85, // 85% proration
			},
			tolerance: 25.0,
		},
		{
			name:        "MRA+10 Early Retirement with Reduction",
			description: "Employee taking MRA+10 early retirement with 5% annual reduction",
			input: models.FERSCalculationInput{
				High3Salary:             70000.0,
				UnusedSickLeaveHours:    500,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1500.0,
				EmployeeContributions:   5600.0,
			},
			ageAtRetirement:   57,
			totalServiceYears: 15.0,
			mra:               57,
			retirementType:    "MRA+10Voluntary",
			expectedSSBenefit: 1500.0,
			fersServiceYears:  15.0,
			expectedResults: expectedFERSResults{
				monthlyPension:           656.25, // After 25% early retirement reduction (5 years * 5%)
				annualPension:            7875.0,
				isEligibleForSupplement:  false, // MRA+10 not eligible for supplement
				fersSupplement:           0.0,
				earlyRetirementReduction: 218.75, // 25% of unreduced pension
				survivorBenefitReduction: 0.0,
				prorationApplied:         false,
			},
			tolerance: 15.0,
		},
		{
			name:        "High-3 Law Enforcement Officer",
			description: "LEO retiring at 50 with enhanced benefits (out of scope but boundary test)",
			input: models.FERSCalculationInput{
				High3Salary:             90000.0,
				UnusedSickLeaveHours:    1200,
				SurvivorBenefitElection: "Full50Percent",
				ExpectedSSBenefitAt62:   2000.0,
				EmployeeContributions:   13500.0,
			},
			ageAtRetirement:   50,
			totalServiceYears: 25.0,
			mra:               50, // LEO MRA
			retirementType:    "LEOImmediate",
			expectedSSBenefit: 2000.0,
			fersServiceYears:  25.0,
			expectedResults: expectedFERSResults{
				monthlyPension:           1687.5, // Using 1% multiplier, after survivor reduction
				annualPension:            20250.0,
				isEligibleForSupplement:  true,   // LEO eligible before regular MRA
				fersSupplement:           1250.0, // $2000 * (25/40)
				survivorBenefitReduction: 187.5,
				prorationApplied:         false,
			},
			tolerance: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", tt.description)
			result := CalculateFERS(
				tt.input,
				tt.ageAtRetirement,
				0, // months
				tt.totalServiceYears-float64(tt.input.UnusedSickLeaveHours)/2087.0, // service before sick leave
				tt.mra,
				tt.retirementType,
				tt.expectedSSBenefit,
				tt.fersServiceYears,
				tt.expectedResults.prorationApplied,
				tt.expectedResults.prorationFactor,
				tt.totalServiceYears,
			)

			// Test monthly pension
			if !floatApproxEqual(result.MonthlyPension, tt.expectedResults.monthlyPension, tt.tolerance) {
				t.Errorf("Monthly pension: got %.2f, expected %.2f (±%.2f)",
					result.MonthlyPension, tt.expectedResults.monthlyPension, tt.tolerance)
			}

			// Test supplement eligibility and amount
			if result.IsEligibleForSupplement != tt.expectedResults.isEligibleForSupplement {
				t.Errorf("Supplement eligibility: got %v, expected %v",
					result.IsEligibleForSupplement, tt.expectedResults.isEligibleForSupplement)
			}

			if tt.expectedResults.isEligibleForSupplement {
				if !floatApproxEqual(result.FersSupplement, tt.expectedResults.fersSupplement, 10.0) {
					t.Errorf("FERS Supplement: got %.2f, expected %.2f",
						result.FersSupplement, tt.expectedResults.fersSupplement)
				}
			}

			// Test proration if applicable
			if result.ProrationApplied != tt.expectedResults.prorationApplied {
				t.Errorf("Proration applied: got %v, expected %v",
					result.ProrationApplied, tt.expectedResults.prorationApplied)
			}

			// Log key results for manual verification
			t.Logf("Results - Monthly Pension: $%.2f, Annual: $%.2f, Supplement: $%.2f, Notes: %s",
				result.MonthlyPension, result.AnnualPension, result.FersSupplement, result.Notes)
		})
	}
}

// TestBoundaryConditions tests edge cases and extreme scenarios
func TestBoundaryConditions(t *testing.T) {
	tests := []struct {
		name        string
		description string
		input       models.FERSCalculationInput
		setup       boundaryTestSetup
		expectError bool
		expectZero  bool
	}{
		{
			name:        "Minimum Service - 5 Years",
			description: "Employee with exactly 5 years service (minimum for vested benefit)",
			input: models.FERSCalculationInput{
				High3Salary:             50000.0,
				UnusedSickLeaveHours:    0,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1000.0,
			},
			setup: boundaryTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 5.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectError: false,
			expectZero:  false,
		},
		{
			name:        "Maximum Service - 42 Years",
			description: "Employee with very long service (boundary for calculations)",
			input: models.FERSCalculationInput{
				High3Salary:             150000.0,
				UnusedSickLeaveHours:    2087, // Exactly 1 year
				SurvivorBenefitElection: "Full50Percent",
				ExpectedSSBenefitAt62:   3000.0,
			},
			setup: boundaryTestSetup{
				ageAtRetirement:   67,
				totalServiceYears: 42.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectError: false,
			expectZero:  false,
		},
		{
			name:        "Very Young Retirement - Age 25 (Invalid)",
			description: "Unrealistic early retirement age to test bounds",
			input: models.FERSCalculationInput{
				High3Salary:             60000.0,
				UnusedSickLeaveHours:    0,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1200.0,
			},
			setup: boundaryTestSetup{
				ageAtRetirement:   25,
				totalServiceYears: 5.0,
				mra:               57,
				retirementType:    "Invalid",
			},
			expectError: false, // Calculation should handle gracefully
			expectZero:  true,  // But should result in zero or minimal benefit
		},
		{
			name:        "Very Old Retirement - Age 95",
			description: "Extreme delayed retirement to test upper bounds",
			input: models.FERSCalculationInput{
				High3Salary:             80000.0,
				UnusedSickLeaveHours:    1000,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   2500.0,
			},
			setup: boundaryTestSetup{
				ageAtRetirement:   95,
				totalServiceYears: 50.0, // Hypothetical very long service
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectError: false,
			expectZero:  false,
		},
		{
			name:        "Zero High-3 Salary",
			description: "Edge case with zero salary to test calculation robustness",
			input: models.FERSCalculationInput{
				High3Salary:             0.0,
				UnusedSickLeaveHours:    0,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   0.0,
			},
			setup: boundaryTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 20.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectError: false,
			expectZero:  true,
		},
		{
			name:        "Maximum Sick Leave - 2087 Hours",
			description: "Employee with exactly maximum annual sick leave accrual",
			input: models.FERSCalculationInput{
				High3Salary:             75000.0,
				UnusedSickLeaveHours:    2087,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1800.0,
			},
			setup: boundaryTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 20.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectError: false,
			expectZero:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing boundary condition: %s", tt.description)

			result := CalculateFERS(
				tt.input,
				tt.setup.ageAtRetirement,
				0,
				tt.setup.totalServiceYears,
				tt.setup.mra,
				tt.setup.retirementType,
				tt.input.ExpectedSSBenefitAt62,
				tt.setup.totalServiceYears,
				false, // proration
				1.0,   // proration factor
				tt.setup.totalServiceYears,
			)

			if tt.expectZero {
				if result.MonthlyPension > 1.0 { // Allow for minimal rounding
					t.Errorf("Expected zero or minimal pension, got %.2f", result.MonthlyPension)
				}
			} else if result.MonthlyPension <= 0 && !tt.expectZero {
				t.Errorf("Expected positive pension, got %.2f", result.MonthlyPension)
			}

			// Log results for analysis
			t.Logf("Boundary test results - Pension: $%.2f, Service: %.2f years, Age: %d",
				result.MonthlyPension, tt.setup.totalServiceYears, tt.setup.ageAtRetirement)
		})
	}
}

// Helper types for test organization
type expectedFERSResults struct {
	monthlyPension               float64
	annualPension                float64
	isEligibleForSupplement      bool
	fersSupplement               float64
	survivorBenefitReduction     float64
	earlyRetirementReduction     float64
	prorationApplied             bool
	prorationFactor              float64
	sickLeaveCredit              float64
	irsSimplifiedMethodExclusion float64
}

type boundaryTestSetup struct {
	ageAtRetirement   int
	totalServiceYears float64
	mra               int
	retirementType    string
}

// Helper function for approximate float comparison with tolerance
func floatApproxEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
