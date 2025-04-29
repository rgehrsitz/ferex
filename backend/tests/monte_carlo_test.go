package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"testing"
)

func TestRunMonteCarlo_BasicSanity(t *testing.T) {
	input := models.MonteCarloInput{
		NumSimulations:   1000,
		Years:            30,
		InitialBalance:   1000000,
		AnnualWithdrawal: 40000,
		ExpectedReturn:   0.06,
		ReturnStdDev:     0.12,
		InflationMean:    0.025,
		InflationStdDev:  0.01,
		Seed:             42,
	}

	result := calculation.RunMonteCarlo(input)

	if result.SuccessRate < 0.0 || result.SuccessRate > 1.0 {
		t.Errorf("SuccessRate out of bounds: got %f", result.SuccessRate)
	}

	if len(result.YearlyBalances) != input.NumSimulations {
		t.Errorf("Expected %d simulations, got %d", input.NumSimulations, len(result.YearlyBalances))
	}

	if len(result.DepletionProbabilities) != input.Years {
		t.Errorf("Expected %d years of depletion probabilities, got %d", input.Years, len(result.DepletionProbabilities))
	}

	if result.Percentiles[50] <= 0 {
		t.Errorf("Median ending balance should be positive, got %f", result.Percentiles[50])
	}
}

func TestRunMonteCarlo_Deterministic_NoVolatility(t *testing.T) {
	input := models.MonteCarloInput{
		NumSimulations:   100,
		Years:            10,
		InitialBalance:   100000,
		AnnualWithdrawal: 10000,
		ExpectedReturn:   0.05,
		ReturnStdDev:     0.0,
		InflationMean:    0.0,
		InflationStdDev:  0.0,
		Seed:             1,
	}

	result := calculation.RunMonteCarlo(input)

	// All simulations should be identical
	firstSim := result.YearlyBalances[0]
	for i := 1; i < len(result.YearlyBalances); i++ {
		for y := 0; y < len(firstSim); y++ {
			if result.YearlyBalances[i][y] != firstSim[y] {
				t.Errorf("Simulation %d year %d differs: got %f, want %f", i, y, result.YearlyBalances[i][y], firstSim[y])
			}
		}
	}
}

func TestRunMonteCarlo_Depletion(t *testing.T) {
	input := models.MonteCarloInput{
		NumSimulations:   100,
		Years:            20,
		InitialBalance:   50000,
		AnnualWithdrawal: 50000,
		ExpectedReturn:   0.0,
		ReturnStdDev:     0.0,
		InflationMean:    0.0,
		InflationStdDev:  0.0,
		Seed:             7,
	}

	result := calculation.RunMonteCarlo(input)

	if result.SuccessRate != 0.0 {
		t.Errorf("Expected 0%% success rate, got %f", result.SuccessRate)
	}

	for i, bal := range result.YearlyBalances {
		if bal[0] != 0 {
			t.Errorf("Simulation %d should be depleted in year 1, got %f", i, bal[0])
		}
	}
}
