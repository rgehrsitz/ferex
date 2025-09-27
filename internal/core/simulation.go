package core

// SimulationConfig stores Monte Carlo parameters for scenario analysis. It
// provides defaults aligned with the risk & Monte Carlo UI controls in the spec.
type SimulationConfig struct {
	HorizonYears      int       `json:"horizon_years"`
	Trials            int       `json:"trials"`
	ConfidenceLevels  []float64 `json:"confidence_levels"`
	Seed              int64     `json:"seed"`
	InflationMean     float64   `json:"inflation_mean"`
	InflationStdDev   float64   `json:"inflation_std_dev"`
	ReturnStdOverride float64   `json:"return_std_override"`
}
