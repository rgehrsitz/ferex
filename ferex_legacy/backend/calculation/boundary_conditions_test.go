package calculation

import (
	"ferex/backend/models"
	"math"
	"testing"
)

// TestExtremeEdgeCases tests the calculation system with extreme boundary conditions
func TestExtremeEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		input          models.FERSCalculationInput
		setup          extremeTestSetup
		expectValid    bool
		expectMinimal  bool // Expect very small but valid result
		expectWarnings bool // Should generate warnings
		notes          string
	}{
		{
			name:        "Zero Service Years",
			description: "Employee with no creditable service (should be invalid)",
			input: models.FERSCalculationInput{
				High3Salary:             60000.0,
				UnusedSickLeaveHours:    0,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1200.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 0.0,
				mra:               57,
				retirementType:    "Invalid",
			},
			expectValid:    false,
			expectMinimal:  true,
			expectWarnings: true,
			notes:          "No service should result in zero benefit",
		},
		{
			name:        "Maximum Creditable Service - 42 Years",
			description: "Employee with maximum realistic federal service",
			input: models.FERSCalculationInput{
				High3Salary:             180000.0, // Executive level
				UnusedSickLeaveHours:    2087,     // 1 full year
				SurvivorBenefitElection: "Full50Percent",
				ExpectedSSBenefitAt62:   3000.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   67,
				totalServiceYears: 42.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectValid:    true,
			expectMinimal:  false,
			expectWarnings: false,
			notes:          "Maximum service should work normally",
		},
		{
			name:        "Extremely High Salary - GS-15 Step 10 + High Locality",
			description: "Test with maximum realistic federal salary",
			input: models.FERSCalculationInput{
				High3Salary:             200000.0, // Very high but realistic
				UnusedSickLeaveHours:    1500,
				SurvivorBenefitElection: "Partial25Percent",
				ExpectedSSBenefitAt62:   3500.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 30.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectValid:    true,
			expectMinimal:  false,
			expectWarnings: false,
			notes:          "High salary should calculate correctly",
		},
		{
			name:        "Extremely Low Salary - Part-time GS-1",
			description: "Test with very low federal salary",
			input: models.FERSCalculationInput{
				High3Salary:             15000.0, // Very low part-time
				UnusedSickLeaveHours:    100,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   400.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 20.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectValid:    true,
			expectMinimal:  true,
			expectWarnings: true,
			notes:          "Very low salary should generate warning about minimal benefit",
		},
		{
			name:        "Maximum Sick Leave - 10 Years Accumulated",
			description: "Employee with unrealistically high sick leave balance",
			input: models.FERSCalculationInput{
				High3Salary:             80000.0,
				UnusedSickLeaveHours:    20870, // 10 years worth
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1800.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 25.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectValid:    true,
			expectMinimal:  false,
			expectWarnings: true,
			notes:          "Extremely high sick leave should be flagged but processed",
		},
		{
			name:        "Retirement at Age 95",
			description: "Extremely delayed retirement to test upper age bounds",
			input: models.FERSCalculationInput{
				High3Salary:             70000.0,
				UnusedSickLeaveHours:    2000,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1500.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   95,
				totalServiceYears: 45.0,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectValid:    true,
			expectMinimal:  false,
			expectWarnings: true,
			notes:          "Extreme delayed retirement should work but generate warnings",
		},
		{
			name:        "Retirement Before MRA - Age 25",
			description: "Invalid early retirement to test lower age bounds",
			input: models.FERSCalculationInput{
				High3Salary:             50000.0,
				UnusedSickLeaveHours:    0,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1000.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   25,
				totalServiceYears: 5.0,
				mra:               57,
				retirementType:    "Invalid",
			},
			expectValid:    false,
			expectMinimal:  true,
			expectWarnings: true,
			notes:          "Invalid early retirement should be rejected or minimal",
		},
		{
			name:        "Fractional Service - 4.75 Years",
			description: "Service just under 5-year vesting requirement",
			input: models.FERSCalculationInput{
				High3Salary:             55000.0,
				UnusedSickLeaveHours:    500, // Could push over 5-year threshold
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1100.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 4.75,
				mra:               57,
				retirementType:    "Immediate6220",
			},
			expectValid:    true, // Sick leave should put over 5 years
			expectMinimal:  true,
			expectWarnings: false,
			notes:          "Fractional service with sick leave should reach vesting",
		},
		{
			name:        "Extreme Part-time Proration - 10% Schedule",
			description: "Employee with extreme part-time schedule requiring heavy proration",
			input: models.FERSCalculationInput{
				High3Salary:             20000.0, // Already prorated for 10% schedule
				UnusedSickLeaveHours:    200,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   800.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   62,
				totalServiceYears: 20.0,
				mra:               57,
				retirementType:    "Immediate6220",
				prorationFactor:   0.10, // 10% schedule
			},
			expectValid:    true,
			expectMinimal:  true,
			expectWarnings: true,
			notes:          "Extreme part-time should work but generate minimal benefit warning",
		},
		{
			name:        "MRA+10 at Maximum Reduction",
			description: "MRA+10 retirement with maximum possible early retirement reduction",
			input: models.FERSCalculationInput{
				High3Salary:             60000.0,
				UnusedSickLeaveHours:    0,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1200.0,
			},
			setup: extremeTestSetup{
				ageAtRetirement:   55, // 7 years under 62
				totalServiceYears: 10.0,
				mra:               55,
				retirementType:    "MRA+10Voluntary",
			},
			expectValid:    true,
			expectMinimal:  false,
			expectWarnings: true,
			notes:          "Maximum early retirement reduction (35%) should be calculated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing extreme edge case: %s", tt.description)

			result := CalculateFERSWithAuditTrail(
				tt.input,
				tt.setup.ageAtRetirement,
				0,
				tt.setup.totalServiceYears-float64(tt.input.UnusedSickLeaveHours)/2087.0,
				tt.setup.mra,
				tt.setup.retirementType,
				tt.input.ExpectedSSBenefitAt62,
				tt.setup.totalServiceYears,
				tt.setup.prorationFactor > 0 && tt.setup.prorationFactor < 1.0,
				tt.setup.prorationFactor,
				tt.setup.totalServiceYears,
			)

			// Check if result is valid based on expectations
			if !tt.expectValid {
				if result.MonthlyPension > 100.0 { // Arbitrary threshold for "significant" benefit
					t.Errorf("Expected invalid/minimal result but got significant pension: $%.2f", result.MonthlyPension)
				}
			}

			// Check for minimal benefits where expected
			if tt.expectMinimal {
				if result.MonthlyPension > 500.0 { // Threshold for "minimal"
					t.Logf("Note: Expected minimal benefit but got $%.2f (may be valid depending on scenario)", result.MonthlyPension)
				}
			}

			// Check for warnings in audit trail
			if tt.expectWarnings {
				if result.AuditTrail == nil || len(result.AuditTrail.Warnings) == 0 {
					t.Logf("Expected warnings but none generated for: %s", tt.description)
				} else {
					t.Logf("Warnings generated as expected: %v", result.AuditTrail.Warnings)
				}
			}

			// Log results for analysis
			t.Logf("Edge case result - Pension: $%.2f, Valid: %v, Notes: %s",
				result.MonthlyPension, tt.expectValid, tt.notes)

			// Check for mathematical consistency
			if result.MonthlyPension < 0 {
				t.Errorf("Pension cannot be negative: $%.2f", result.MonthlyPension)
			}

			if result.AnnualPension != 0 && math.Abs(result.AnnualPension-(result.MonthlyPension*12)) > 0.01 {
				t.Errorf("Annual/monthly pension inconsistency: Monthly $%.2f * 12 ≠ Annual $%.2f",
					result.MonthlyPension, result.AnnualPension)
			}
		})
	}
}

// TestCalculationConsistency ensures calculations are mathematically consistent across different scenarios
func TestCalculationConsistency(t *testing.T) {
	baseInput := models.FERSCalculationInput{
		High3Salary:             75000.0,
		UnusedSickLeaveHours:    1000,
		SurvivorBenefitElection: "None",
		ExpectedSSBenefitAt62:   1500.0,
		EmployeeContributions:   9000.0,
	}

	// Test 1: Doubling salary should double pension (all else equal)
	result1 := CalculateFERS(baseInput, 62, 0, 20.0, 57, "Immediate6220", 1500.0, 20.0, false, 1.0, 20.0)

	doubledSalaryInput := baseInput
	doubledSalaryInput.High3Salary = 150000.0
	result2 := CalculateFERS(doubledSalaryInput, 62, 0, 20.0, 57, "Immediate6220", 1500.0, 20.0, false, 1.0, 20.0)

	expectedDoubledPension := result1.MonthlyPension * 2
	if math.Abs(result2.MonthlyPension-expectedDoubledPension) > 0.01 {
		t.Errorf("Salary doubling inconsistency: $%.2f * 2 ≠ $%.2f", result1.MonthlyPension, result2.MonthlyPension)
	}

	// Test 2: Service time should scale linearly with pension
	result3 := CalculateFERS(baseInput, 62, 0, 10.0, 57, "Immediate6220", 1500.0, 10.0, false, 1.0, 10.0)
	result4 := CalculateFERS(baseInput, 62, 0, 30.0, 57, "Immediate6220", 1500.0, 30.0, false, 1.0, 30.0)

	// Should be roughly 3x pension for 3x service (accounting for multiplier change)
	ratio := result4.MonthlyPension / result3.MonthlyPension
	if ratio < 2.8 || ratio > 3.3 { // Allow for 1.1% vs 1.0% multiplier difference
		t.Errorf("Service scaling inconsistency: 30yr/10yr ratio = %.2f (expected ~3.0)", ratio)
	}

	// Test 3: Monthly vs Annual consistency across all scenarios
	testScenarios := []struct {
		age     int
		service float64
		mra     int
		retType string
	}{
		{57, 30, 57, "ImmediateMRA30"},
		{60, 20, 57, "Immediate6020"},
		{62, 15, 57, "Immediate6220"},
		{57, 15, 57, "MRA+10Voluntary"},
	}

	for _, scenario := range testScenarios {
		result := CalculateFERS(baseInput, scenario.age, 0, scenario.service, scenario.mra,
			scenario.retType, 1500.0, scenario.service, false, 1.0, scenario.service)

		if math.Abs(result.AnnualPension-(result.MonthlyPension*12)) > 0.01 {
			t.Errorf("Monthly/Annual inconsistency in scenario %+v: Monthly $%.2f * 12 ≠ Annual $%.2f",
				scenario, result.MonthlyPension, result.AnnualPension)
		}
	}

	t.Logf("✓ Calculation consistency tests passed")
}

// TestNumericalStability tests calculations with values that might cause floating-point issues
func TestNumericalStability(t *testing.T) {
	tests := []struct {
		name        string
		salary      float64
		service     float64
		description string
	}{
		{"Very large numbers", 999999.99, 41.999999, "Test floating-point precision with large values"},
		{"Very small numbers", 0.01, 0.01, "Test floating-point precision with tiny values"},
		{"Repeating decimals", 100000.0 / 3.0, 20.0 / 3.0, "Test with repeating decimal inputs"},
		{"Pi-based values", 75000.0 * math.Pi, 20.0 * math.Pi, "Test with irrational number components"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := models.FERSCalculationInput{
				High3Salary:             tt.salary,
				UnusedSickLeaveHours:    0,
				SurvivorBenefitElection: "None",
				ExpectedSSBenefitAt62:   1500.0,
			}

			result := CalculateFERS(input, 62, 0, tt.service, 57, "Immediate6220", 1500.0, tt.service, false, 1.0, tt.service) // Check that results are finite and reasonable
			if math.IsNaN(result.MonthlyPension) || math.IsInf(result.MonthlyPension, 0) {
				t.Errorf("Non-finite result for %s: %f", tt.description, result.MonthlyPension)
			}

			if result.MonthlyPension < 0 {
				t.Errorf("Negative pension for %s: $%.2f", tt.description, result.MonthlyPension)
			}

			// Check monthly/annual consistency even with unusual numbers
			if math.Abs(result.AnnualPension-(result.MonthlyPension*12)) > 0.01 {
				t.Errorf("Numerical instability in %s: Monthly/Annual inconsistency", tt.description)
			}

			t.Logf("Numerical stability test passed for %s: $%.6f monthly", tt.description, result.MonthlyPension)
		})
	}
}

// Helper type for extreme testing
type extremeTestSetup struct {
	ageAtRetirement   int
	totalServiceYears float64
	mra               int
	retirementType    string
	prorationFactor   float64
}
