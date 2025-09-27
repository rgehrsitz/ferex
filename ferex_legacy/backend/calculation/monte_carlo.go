package calculation

import (
	"math"
	"math/rand"
	"time"
	"ferex/backend/models"
)

// RunMonteCarloSimulation runs a Monte Carlo simulation for TSP or annuity projections.
// It supports both pre-tax (Traditional) and Roth balances, as well as annuity payouts.
func RunMonteCarloSimulation(input models.MonteCarloInput) models.MonteCarloResult {
	if input.NumSimulations <= 0 {
		input.NumSimulations = 1000
	}
	if input.Years <= 0 {
		input.Years = 30
	}
	
	// Seed RNG
	seed := input.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))

	yearlyBalances := make([][]float64, input.NumSimulations)
	depletionProbabilities := make([]float64, input.Years)
	endingBalances := make([]float64, input.NumSimulations)
	depletedCountByYear := make([]int, input.Years)
	successfulRuns := 0

	for sim := 0; sim < input.NumSimulations; sim++ {
		balance := input.InitialBalance
		balances := make([]float64, input.Years)
		depleted := false
		for year := 0; year < input.Years; year++ {
			// Draw annual return and inflation
			annualReturn := rng.NormFloat64()*input.ReturnStdDev + input.ExpectedReturn
			annualInflation := rng.NormFloat64()*input.InflationStdDev + input.InflationMean
			// Grow balance
			balance *= (1 + annualReturn)
			// Withdraw (inflation-adjusted)
			withdrawal := input.AnnualWithdrawal * math.Pow(1+annualInflation, float64(year))
			balance -= withdrawal
			if balance < 0 {
				balance = 0
				depleted = true
			}
			balances[year] = balance
			if depleted && depletedCountByYear[year] == 0 {
				depletedCountByYear[year]++
			}
		}
		if balances[input.Years-1] > 0 {
			successfulRuns++
		}
		yearlyBalances[sim] = balances
		endingBalances[sim] = balances[input.Years-1]
	}
	// Calculate percentiles
	percentiles := map[int]float64{}
	percentileList := []int{10, 25, 50, 75, 90}
	copyBalances := make([]float64, len(endingBalances))
	copy(copyBalances, endingBalances)
	// Sort for percentile calculation
	for _, p := range percentileList {
		percentiles[p] = percentile(copyBalances, p)
	}
	// Calculate depletion probabilities
	for year := 0; year < input.Years; year++ {
		depletionProbabilities[year] = float64(depletedCountByYear[year]) / float64(input.NumSimulations)
	}
	return models.MonteCarloResult{
		SuccessRate:           float64(successfulRuns) / float64(input.NumSimulations),
		Percentiles:           percentiles,
		YearlyBalances:        yearlyBalances,
		DepletionProbabilities: depletionProbabilities,
	}
}

// percentile returns the p-th percentile (e.g., 50 for median) of a float64 slice
func percentile(data []float64, p int) float64 {
	if len(data) == 0 {
		return 0
	}
	copyData := make([]float64, len(data))
	copy(copyData, data)
	// Sort ascending
	for i := 1; i < len(copyData); i++ {
		key := copyData[i]
		j := i - 1
		for j >= 0 && copyData[j] > key {
			copyData[j+1] = copyData[j]
			j--
		}
		copyData[j+1] = key
	}
	pos := float64(p) / 100 * float64(len(copyData)-1)
	lower := int(math.Floor(pos))
	upper := int(math.Ceil(pos))
	if lower == upper {
		return copyData[lower]
	}
	return copyData[lower] + (copyData[upper]-copyData[lower])*(pos-float64(lower))
}
