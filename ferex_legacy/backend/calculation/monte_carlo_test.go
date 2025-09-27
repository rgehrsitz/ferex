package calculation

import (
	"testing"
	"ferex/backend/models"
)

func TestRunMonteCarloSimulation_Basic(t *testing.T) {
	input := models.MonteCarloInput{
		NumSimulations:   1000,
		Years:            30,
		InitialBalance:   100000.0,
		AnnualWithdrawal: 4000.0,
		ExpectedReturn:   0.06,
		ReturnStdDev:     0.12,
		InflationMean:    0.025,
		InflationStdDev:  0.01,
		Seed:             42,
	}
	result := RunMonteCarloSimulation(input)

	if result.SuccessRate < 0.5 || result.SuccessRate > 1.0 {
		t.Errorf("unexpected success rate: got %v, want between 0.5 and 1.0", result.SuccessRate)
	}
	if len(result.YearlyBalances) != input.NumSimulations {
		t.Errorf("expected %d simulations, got %d", input.NumSimulations, len(result.YearlyBalances))
	}
	if len(result.DepletionProbabilities) != input.Years {
		t.Errorf("expected %d years of depletion probabilities, got %d", input.Years, len(result.DepletionProbabilities))
	}
	if len(result.Percentiles) == 0 {
		t.Error("expected percentiles to be calculated")
	}
	// Check that the 50th percentile is reasonable
	p50 := result.Percentiles[50]
	if p50 < 50000 || p50 > 500000 {
		t.Errorf("unexpected 50th percentile ending balance: got %v", p50)
	}
}
