package calculation

import (
	"math"
	"testing"

	"ferex/backend/models"
)

func TestCheckFersSupplementEligibility(t *testing.T) {
	tests := []struct {
		name              string
		ageAtRetirement   int
		mra               int
		totalServiceYears float64
		retirementType    string // Added to match function signature
		want              bool
	}{
		// Eligible Scenarios - Assuming appropriate retirement types
		{"MRA with 30 years", 57, 57, 30, "ImmediateMRA30", true},
		{"MRA with >30 years", 57, 57, 35, "ImmediateMRA30", true},
		{"Age 60 with 20 years", 60, 57, 20, "Immediate6020", true},
		{"Age 60 with >20 years", 60, 57, 25, "Immediate6020", true},
		{"Age 61 with 20 years (still <62)", 61, 57, 20, "Immediate6020", true},

		// Ineligible Scenarios
		{"Age 62 (supplement ends)", 62, 57, 30, "ImmediateMRA30", false},
		{"Age >62 (supplement ends)", 63, 57, 30, "ImmediateMRA30", false},
		{"Below MRA", 56, 57, 30, "ImmediateMRA30", false},                   // Retirement type doesn't make it eligible if basic age/service fails
		{"MRA with <30 years (not 60)", 57, 57, 25, "ImmediateMRA30", false}, // Fails MRA+30
		{"Age 60 with <20 years", 60, 57, 19, "Immediate6020", false},        // Fails 60+20
		{"Insufficient age and service", 50, 57, 10, "ImmediateMRA30", false},
		{"MRA with 10 years (MRA+10 rule - NOT eligible for supplement)", 57, 57, 10, "MRA+10Voluntary", false}, // MRA+10 is NOT supplement eligible per OPM rules
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkFersSupplementEligibility(tt.ageAtRetirement, tt.mra, tt.retirementType, tt.totalServiceYears); got != tt.want {
				t.Errorf("checkFersSupplementEligibility(%d, %d, %s, %.1f) = %v, want %v", tt.ageAtRetirement, tt.mra, tt.retirementType, tt.totalServiceYears, got, tt.want)
			}
		})
	}
}

const epsilon = 1e-6 // Tolerance for float comparisons

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestCalculateFERS(t *testing.T) {
	tests := []struct {
		name                                       string
		input                                      models.FERSCalculationInput
		ageAtRetirementYears                       int
		ageAtRetirementMonths                      int // Added, can be 0 for simplicity
		totalCreditableServiceYearsBeforeSickLeave float64
		totalCombinedActualServiceYears            float64
		mraYears                                   int
		retirementType                             string // Added
		expectedSSBenefitAt62                      float64
		yearsActuallyWorkedUnderFERS               float64
		prorationAppliedByCaller                   bool    // New field for testing
		prorationFactorUsedByCaller                float64 // New field for testing
		wantResult                                 models.FERSCalculationResult
	}{
		{
			name: "Basic 1% scenario, supplement eligible",
			input: models.FERSCalculationInput{
				High3Salary:          100000,
				UnusedSickLeaveHours: 0,
				ServicePeriods:       nil, // IsPartTime and PartTimeProrationFactor removed
				// Other FERSCalculationInput fields like DateOfBirth, SCD can be omitted if not directly used by CalculateFERS for these tests
			},
			prorationAppliedByCaller:                   false,
			prorationFactorUsedByCaller:                1.0,
			ageAtRetirementYears:                       60,
			ageAtRetirementMonths:                      0,
			totalCreditableServiceYearsBeforeSickLeave: 20,
			mraYears:                     57,
			retirementType:               "Immediate6020", // Eligible for supplement if age < 62
			expectedSSBenefitAt62:        1500,
			yearsActuallyWorkedUnderFERS: 20,
			wantResult: models.FERSCalculationResult{
				TotalServiceYears:       20.0,
				BasicAnnuity:            20000.0, // 100000 * 20 * 0.01
				ProratedPension:         20000.0, // Equals BasicAnnuity if no proration
				IsEligibleForSupplement: true,
				FersSupplement:          750.0, // 1500 * (20/40)
				AnnualPension:           20000.0,
				MonthlyPension:          20000.0 / 12.0,
			},
		},
		{
			name: "1.1% multiplier, supplement ineligible (age)",
			input: models.FERSCalculationInput{
				High3Salary:          100000,
				UnusedSickLeaveHours: 0,
				ServicePeriods:       nil,
			},
			prorationAppliedByCaller:                   false,
			prorationFactorUsedByCaller:                1.0,
			ageAtRetirementYears:                       62,
			ageAtRetirementMonths:                      0,
			totalCreditableServiceYearsBeforeSickLeave: 20,
			mraYears:                     57,
			retirementType:               "ImmediateMRA30", // Example, actual type might vary but age 62 makes supplement ineligible
			expectedSSBenefitAt62:        1500,
			yearsActuallyWorkedUnderFERS: 20,
			wantResult: models.FERSCalculationResult{
				TotalServiceYears:       20.0,
				BasicAnnuity:            22000.0, // 100000 * 20 * 0.011
				ProratedPension:         22000.0, // Equals BasicAnnuity if no proration
				IsEligibleForSupplement: false,
				FersSupplement:          0.0,
				AnnualPension:           22000.0,
				MonthlyPension:          22000.0 / 12.0,
			},
		},
		{
			name: "With Sick Leave (2087 hrs = 1 yr), 1% mult, supplement eligible",
			input: models.FERSCalculationInput{
				High3Salary:          100000,
				UnusedSickLeaveHours: 2087,
				ServicePeriods:       nil,
			},
			prorationAppliedByCaller:                   false,
			prorationFactorUsedByCaller:                1.0,
			ageAtRetirementYears:                       57,
			ageAtRetirementMonths:                      0,
			totalCreditableServiceYearsBeforeSickLeave: 29, // Actual service
			mraYears:                     57,
			retirementType:               "ImmediateMRA30",
			expectedSSBenefitAt62:        1500,
			yearsActuallyWorkedUnderFERS: 29, // Actual FERS service years for supplement
			wantResult: models.FERSCalculationResult{
				TotalServiceYears:       30.0,    // 29 + (2087/2087.0)
				BasicAnnuity:            30000.0, // 100000 * 30 * 0.01
				ProratedPension:         30000.0, // Equals BasicAnnuity if no proration
				IsEligibleForSupplement: true,
				FersSupplement:          1087.5, // 1500 * (29/40)
				AnnualPension:           30000.0,
				MonthlyPension:          30000.0 / 12.0,
			},
		},
		{
			name: "Part-time proration (0.5 factor)",
			input: models.FERSCalculationInput{
				High3Salary:          50000, // Pre-adjusted: 100000 * 0.5
				UnusedSickLeaveHours: 0,
				ServicePeriods:       nil,
			},
			ageAtRetirementYears:                       60,
			ageAtRetirementMonths:                      0,
			totalCreditableServiceYearsBeforeSickLeave: 20,
			mraYears:                     57,
			retirementType:               "Immediate6020",
			expectedSSBenefitAt62:        1500,
			yearsActuallyWorkedUnderFERS: 20,
			prorationAppliedByCaller:     true,
			prorationFactorUsedByCaller:  0.5,
			wantResult: models.FERSCalculationResult{
				TotalServiceYears:       20.0,
				BasicAnnuity:            10000.0, // Based on pre-adjusted High3: 50000 * 20 * 0.01
				ProrationApplied:        true,
				ProratedPension:         10000.0, // Equals BasicAnnuity as High3 was pre-adjusted
				IsEligibleForSupplement: true,
				FersSupplement:          750.0, // Supplement not typically prorated by this factor
				AnnualPension:           10000.0,
				MonthlyPension:          10000.0 / 12.0,
				// Note: 'Notes' field comparison is omitted for brevity but might be needed for full validation
			},
		},
		{
			name: "MRA+10_Voluntary,_NOT_Supplement_Eligible_With_Reduction",
			input: models.FERSCalculationInput{
				High3Salary:          80000,
				UnusedSickLeaveHours: 0,
				ServicePeriods:       nil,
			},
			prorationAppliedByCaller:                   false,
			prorationFactorUsedByCaller:                1.0,
			ageAtRetirementYears:                       57,
			ageAtRetirementMonths:                      0,
			totalCreditableServiceYearsBeforeSickLeave: 15,
			mraYears:                     57,
			retirementType:               "MRA+10Voluntary", // Example for non-supplement eligible scenario based on rules
			expectedSSBenefitAt62:        1200,
			yearsActuallyWorkedUnderFERS: 15,
			wantResult: models.FERSCalculationResult{
				TotalServiceYears:        15.0,
				BasicAnnuity:             12000.0, // 80000 * 15 * 0.01
				ProratedPension:          12000.0, // Equals BasicAnnuity if no proration
				EarlyRetirementReduction: 3000.0,  // 12000 * 0.25 (5 years under 62 * 5% per year)				IsEligibleForSupplement: false, // MRA+10 at 57 is NOT supplement eligible per OPM rules
				FersSupplement:           0.0,     // No supplement for MRA+10
				AnnualPension:            9000.0,  // BasicAnnuity - EarlyRetirementReduction
				MonthlyPension:           9000.0 / 12.0,
			},
		},
		{
			name: "Deferred_Retirement_at_62",
			input: models.FERSCalculationInput{
				High3Salary:          100000,
				UnusedSickLeaveHours: 0,
				ServicePeriods:       nil,
			},
			prorationAppliedByCaller:                   false,
			prorationFactorUsedByCaller:                1.0,
			ageAtRetirementYears:                       62,
			ageAtRetirementMonths:                      0,
			totalCreditableServiceYearsBeforeSickLeave: 15,
			mraYears:                     57, // MRA is not directly relevant for claiming age, but good for context
			retirementType:               "Deferred",
			expectedSSBenefitAt62:        0, // No supplement for deferred
			yearsActuallyWorkedUnderFERS: 15,
			wantResult: models.FERSCalculationResult{
				TotalServiceYears:       15.0,
				BasicAnnuity:            15000.0, // 100000 * 15 * 0.01 (age 62, <20 years = 1.0%)
				ProratedPension:         15000.0, // Equals BasicAnnuity
				IsEligibleForSupplement: false,
				FersSupplement:          0.0,
				AnnualPension:           15000.0,
				MonthlyPension:          15000.0 / 12.0,
			},
		},
		{
			name: "MRA+10_Voluntary,_Early_Reduction",
			input: models.FERSCalculationInput{
				High3Salary:             100000,
				UnusedSickLeaveHours:    0,
				ServicePeriods:          nil,
				SurvivorBenefitElection: "None", // Assume no survivor cost for this specific test focus
			},
			prorationAppliedByCaller:                   false,
			prorationFactorUsedByCaller:                1.0,
			ageAtRetirementYears:                       57, // MRA
			ageAtRetirementMonths:                      0,
			totalCreditableServiceYearsBeforeSickLeave: 15, // Between 10 and 20 years
			mraYears:                     57,
			retirementType:               "MRA+10Voluntary",
			expectedSSBenefitAt62:        1500, // NOT eligible for supplement
			yearsActuallyWorkedUnderFERS: 15,
			wantResult: models.FERSCalculationResult{
				TotalServiceYears:        15.0,
				BasicAnnuity:             15000.0,              // 100000 * 15 * 0.01
				ProratedPension:          15000.0,              // Equals BasicAnnuity
				EarlyRetirementReduction: 15000.0 * 0.05 * 5.0, // 5 years under 62 (62-57=5), 5% reduction per year				IsEligibleForSupplement: false,
				FersSupplement:           0.0,                  // No supplement for MRA+10
				AnnualPension:            15000.0 - (15000.0 * 0.05 * 5.0),
				MonthlyPension:           (15000.0 - (15000.0 * 0.05 * 5.0)) / 12.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := CalculateFERS(tt.input, tt.ageAtRetirementYears, tt.ageAtRetirementMonths, tt.totalCreditableServiceYearsBeforeSickLeave, tt.mraYears, tt.retirementType, tt.expectedSSBenefitAt62, tt.yearsActuallyWorkedUnderFERS, tt.prorationAppliedByCaller, tt.prorationFactorUsedByCaller, tt.totalCreditableServiceYearsBeforeSickLeave)

			if !floatEquals(gotResult.TotalServiceYears, tt.wantResult.TotalServiceYears) {
				t.Errorf("CalculateFERS().TotalServiceYears = %v, want %v", gotResult.TotalServiceYears, tt.wantResult.TotalServiceYears)
			}
			if !floatEquals(gotResult.BasicAnnuity, tt.wantResult.BasicAnnuity) {
				t.Errorf("CalculateFERS().BasicAnnuity = %v, want %v", gotResult.BasicAnnuity, tt.wantResult.BasicAnnuity)
			}
			if gotResult.ProrationApplied != tt.wantResult.ProrationApplied {
				t.Errorf("CalculateFERS().ProrationApplied = %v, want %v", gotResult.ProrationApplied, tt.wantResult.ProrationApplied)
			}
			if !floatEquals(gotResult.ProratedPension, tt.wantResult.ProratedPension) {
				t.Errorf("CalculateFERS().ProratedPension = %v, want %v", gotResult.ProratedPension, tt.wantResult.ProratedPension)
			}
			if gotResult.IsEligibleForSupplement != tt.wantResult.IsEligibleForSupplement {
				t.Errorf("CalculateFERS().IsEligibleForSupplement = %v, want %v", gotResult.IsEligibleForSupplement, tt.wantResult.IsEligibleForSupplement)
			}
			if !floatEquals(gotResult.FersSupplement, tt.wantResult.FersSupplement) {
				t.Errorf("CalculateFERS().FersSupplement = %v, want %v", gotResult.FersSupplement, tt.wantResult.FersSupplement)
			}
			if !floatEquals(gotResult.AnnualPension, tt.wantResult.AnnualPension) {
				t.Errorf("CalculateFERS().AnnualPension = %v, want %v", gotResult.AnnualPension, tt.wantResult.AnnualPension)
			}
			if !floatEquals(gotResult.MonthlyPension, tt.wantResult.MonthlyPension) {
				t.Errorf("CalculateFERS().MonthlyPension = %v, want %v", gotResult.MonthlyPension, tt.wantResult.MonthlyPension)
			}
		})
	}
}
