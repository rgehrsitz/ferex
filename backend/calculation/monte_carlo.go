package calculation

import (
	"ferex/backend/models"
	"math"
	"math/rand"
)

// RunMonteCarlo runs a Monte Carlo simulation for retirement withdrawals
func RunMonteCarlo(input models.MonteCarloInput) models.MonteCarloResult {
	if input.NumSimulations <= 0 {
		input.NumSimulations = 1000
	}
	if input.Years <= 0 {
		input.Years = 30
	}
	if input.Seed != 0 {
		rand.Seed(input.Seed)
	}

	sims := input.NumSimulations
	years := input.Years
	balances := make([][]float64, sims)
	finalBalances := make([]float64, sims)
	depletionByYear := make([]int, years)

	for i := 0; i < sims; i++ {
		bal := input.InitialBalance
		balances[i] = make([]float64, years)
		depleted := false
		for y := 0; y < years; y++ {
			// Random return for this year
			ret := rand.NormFloat64()*input.ReturnStdDev + input.ExpectedReturn
			inf := rand.NormFloat64()*input.InflationStdDev + input.InflationMean
			if !depleted {
				bal = bal * (1 + ret)
				bal -= input.AnnualWithdrawal * math.Pow(1+inf, float64(y))
				if bal < 0 {
					bal = 0
					depleted = true
					depletionByYear[y]++
				}
			}
			balances[i][y] = bal
		}
		finalBalances[i] = bal
	}

	// Calculate percentiles
	percentiles := map[int]float64{}
	for _, p := range []int{10, 25, 50, 75, 90} {
		percentiles[p] = Percentile(finalBalances, float64(p))
	}

	// Depletion probability by year
	depletionProb := make([]float64, years)
	for y := 0; y < years; y++ {
		depletionProb[y] = float64(depletionByYear[y]) / float64(sims)
	}

	success := 0
	for _, bal := range finalBalances {
		if bal > 0 {
			success++
		}
	}

	return models.MonteCarloResult{
		SuccessRate:            float64(success) / float64(sims),
		Percentiles:            percentiles,
		YearlyBalances:         balances,
		DepletionProbabilities: depletionProb,
	}
}

// Percentile returns the p-th percentile of a slice
func Percentile(data []float64, p float64) float64 {
	if len(data) == 0 {
		return 0
	}
	copyData := make([]float64, len(data))
	copy(copyData, data)
	sortFloat64s(copyData)
	idx := int((p/100.0)*float64(len(copyData)-1) + 0.5)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(copyData) {
		idx = len(copyData) - 1
	}
	return copyData[idx]
}

// sortFloat64s sorts a slice of float64s in ascending order
func sortFloat64s(a []float64) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j-1] > a[j]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
