package models

// MonteCarloInput defines the parameters for the simulation
// Add fields as needed for future extensibility
// (e.g., asset allocation, longevity, etc.)
type MonteCarloInput struct {
	NumSimulations      int     // Number of simulation runs
	Years               int     // Years to simulate
	InitialBalance      float64 // Starting account balance
	AnnualWithdrawal    float64 // Planned annual withdrawal
	ExpectedReturn      float64 // Mean annual return (e.g., 0.06 for 6%)
	ReturnStdDev        float64 // Std dev of returns (e.g., 0.12 for 12%)
	InflationMean       float64 // Mean annual inflation (e.g., 0.025 for 2.5%)
	InflationStdDev     float64 // Std dev of inflation
	Seed                int64   // Optional: for deterministic tests
}

// MonteCarloResult provides summary statistics and simulation output
type MonteCarloResult struct {
	SuccessRate             float64              // % of simulations with money left at end
	Percentiles             map[int]float64      // e.g., 10, 25, 50, 75, 90 percentiles of ending balance
	YearlyBalances          [][]float64          // [simulation][year] balances
	DepletionProbabilities  []float64            // Probability of depletion by year
}
