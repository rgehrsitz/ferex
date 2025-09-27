package models

// TSPFundAllocation defines the allocation across G, F, C, S, I funds.
// Corresponds to TypeScript TSPFundAllocation.
// This will be used by the new ScenarioVariant.
// Note: The existing TSPFundAllocationPercentages is similar but might be used by older calculation logic.
// We prefer this new one for direct mapping from ScenarioInput.
type TSPFundAllocation struct {
	G         *float64 `json:"G,omitempty"`
	F         *float64 `json:"F,omitempty"`
	C         *float64 `json:"C,omitempty"`
	S         *float64 `json:"S,omitempty"`
	I         *float64 `json:"I,omitempty"`
	LFundName *string  `json:"LFundName,omitempty"` // e.g., "L2050", "LIncome"
}

// TSPIndividualReturnAssumptions holds expected returns for each fund.
// Corresponds to TypeScript TSPReturnAssumptions.
// This will be used by the new ScenarioVariant.
type TSPIndividualReturnAssumptions struct {
	G *float64 `json:"G,omitempty"`
	F *float64 `json:"F,omitempty"`
	C *float64 `json:"C,omitempty"`
	S *float64 `json:"S,omitempty"`
	I *float64 `json:"I,omitempty"`
}

// TSPVolatilityRates holds standard deviations for each fund.
// Corresponds to TypeScript TSPVolatility / ScenarioInput.volatilityTSPFunds.
// This will be used by the new ScenarioVariant.
type TSPVolatilityRates struct {
	GStdDev *float64 `json:"gStdDev,omitempty"`
	FStdDev *float64 `json:"fStdDev,omitempty"`
	CStdDev *float64 `json:"cStdDev,omitempty"`
	SStdDev *float64 `json:"sStdDev,omitempty"`
	IStdDev *float64 `json:"iStdDev,omitempty"`
}

// TSPFundAllocationPercentages defines the allocation across G, F, C, S, I funds.
// Percentages should sum to 100 if fully allocated.
// Values are percentages (e.g., 20.0 for 20%).
// LFundName can be used if an L Fund is selected, implying an underlying mix.
// If LFundName is present, individual fund percentages might be derived or ignored based on context.
// TSPFundAllocationPercentages defines the allocation across G, F, C, S, I funds.
// Percentages should sum to 100 if fully allocated.
// Values are percentages (e.g., 20.0 for 20%).
// LFundName can be used if an L Fund is selected, implying an underlying mix.
// If LFundName is present, individual fund percentages might be derived or ignored based on context.
// UPDATED: Fields changed to pointers to support nullability and omitempty.
// This struct might be refactored or replaced by TSPFundAllocation if calculation logic is updated.
type TSPFundAllocationPercentages struct {
	G         *float64 `json:"g,omitempty"`
	F         *float64 `json:"f,omitempty"`
	C         *float64 `json:"c,omitempty"`
	S         *float64 `json:"s,omitempty"`
	I         *float64 `json:"i,omitempty"`
	LFundName *string  `json:"lFundName,omitempty"` // e.g., "L2050", "LIncome"
}

// TSPReturnAssumptions holds expected returns for each fund or an overall rate.
// If individual fund returns are provided, they might be used for a weighted average.
// All rates are annual (e.g., 0.06 for 6%).
// TSPReturnAssumptions holds expected returns for each fund or an overall rate.
// If individual fund returns are provided, they might be used for a weighted average.
// All rates are annual (e.g., 0.06 for 6%).
// UPDATED: Fields changed to pointers to support nullability and omitempty.
// This struct might be refactored if calculation logic is updated to use TSPIndividualReturnAssumptions and TSPVolatilityRates.
type TSPReturnAssumptions struct {
	G                 *float64 `json:"g,omitempty"`
	F                 *float64 `json:"f,omitempty"`
	C                 *float64 `json:"c,omitempty"`
	S                 *float64 `json:"s,omitempty"`
	I                 *float64 `json:"i,omitempty"`
	Overall           *float64 `json:"overall,omitempty"`           // Used if per-fund rates are not provided
	UseOverallForPre  bool     `json:"useOverallForPre,omitempty"`  // Flag to use Overall for pre-retirement growth
	UseOverallForPost bool     `json:"useOverallForPost,omitempty"` // Flag to use Overall for post-retirement growth
	VolatilityG       *float64 `json:"volatilityG,omitempty"`       // Standard deviation for G fund
	VolatilityF       *float64 `json:"volatilityF,omitempty"`       // Standard deviation for F fund
	VolatilityC       *float64 `json:"volatilityC,omitempty"`       // Standard deviation for C fund
	VolatilityS       *float64 `json:"volatilityS,omitempty"`       // Standard deviation for S fund
	VolatilityI       *float64 `json:"volatilityI,omitempty"`       // Standard deviation for I fund
}

// TSPCalculationInput holds all user-provided and derived data for TSP projection and withdrawal modeling.
type TSPCalculationInput struct {
	// Current State
	BaseSalaryForContributions float64                      `json:"baseSalaryForContributions"` // Usually High-3 or current salary for contribution % calc
	CurrentAgeYears            int                          `json:"currentAgeYears"`
	BirthYear                  int                          `json:"birthYear"` // Employee's year of birth, used for RMD calculations
	RetirementAgeYears         int                          `json:"retirementAgeYears"`
	CurrentTraditionalBalance  float64                      `json:"currentTraditionalBalance"`
	CurrentRothBalance         float64                      `json:"currentRothBalance"`
	CurrentLoanBalance         float64                      `json:"currentLoanBalance,omitempty"`
	CurrentAllocation          TSPFundAllocationPercentages `json:"currentAllocation"`

	// Contributions Pre-Retirement
	EmployeeContributionAmount           float64                      `json:"employeeContributionAmount,omitempty"`     // If fixed amount
	EmployeeContributionPercentage       float64                      `json:"employeeContributionPercentage,omitempty"` // If percentage of base salary
	IsContributionPercentage             bool                         `json:"isContributionPercentage"`
	ContributeUntilRetirement            bool                         `json:"contributeUntilRetirement"`     // True if contributing up to retirementAgeYears
	ContributionStopAge                  int                          `json:"contributionStopAge,omitempty"` // If ContributeUntilRetirement is false
	CatchUpContributionsEligible         bool                         `json:"catchUpContributionsEligible"`
	TraditionalContributionAllocationPct float64                      `json:"traditionalContributionAllocationPct"` // % of total contribution to Traditional TSP (0-100)
	RothContributionAllocationPct        float64                      `json:"rothContributionAllocationPct"`        // % of total contribution to Roth TSP (0-100)
	ContributionFundAllocation           TSPFundAllocationPercentages `json:"contributionFundAllocation"`

	// Growth Assumptions
	UserReturnAssumptions       TSPReturnAssumptions         `json:"userReturnAssumptions"`
	ExpenseRatio                float64                      `json:"expenseRatio,omitempty"` // Annual TSP expense ratio (e.g., 0.00051 for 0.051%)
	ExpectedAnnualInflationRate float64                      `json:"expectedAnnualInflationRate"`
	FutureAllocationStrategy    string                       `json:"futureAllocationStrategy,omitempty"` // "MaintainCurrent", "MatchContributionAllocation", "UsePostRetirementAllocation"
	PostRetirementAllocation    TSPFundAllocationPercentages `json:"postRetirementAllocation,omitempty"` // Allocation to use after retirement

	// Withdrawals Post-Retirement
	WithdrawalStrategy         string  `json:"withdrawalStrategy"` // "None", "FixedAmountYearly", "FixedAmountMonthly", "PercentageOfBalanceYearly", "IRSMinimumRequiredDistribution"
	WithdrawalFixedAmountValue float64 `json:"withdrawalFixedAmountValue,omitempty"`
	WithdrawalPercentageValue  float64 `json:"withdrawalPercentageValue,omitempty"`
	WithdrawalStartDate        string  `json:"withdrawalStartDate"`          // "Retirement", "SpecificAge"
	WithdrawalStartAge         int     `json:"withdrawalStartAge,omitempty"` // if WithdrawalStartDate is "SpecificAge"
	WithdrawalOrder            string  `json:"withdrawalOrder,omitempty"`    // "ProRata", "TraditionalFirst", "RothFirst"
	YearsToProjectWithdrawals  int     `json:"yearsToProjectWithdrawals"`    // How many years of withdrawals to simulate (e.g., until age 95 or 100)
}

// YearlyTSPWithdrawalDetail provides a snapshot of TSP account for a given year during withdrawal phase.
type YearlyTSPWithdrawalDetail struct {
	Year                  int     `json:"year"` // Calendar year
	Age                   int     `json:"age"`  // Age at end of year
	BeginningBalanceTotal float64 `json:"beginningBalanceTotal"`
	BeginningBalanceTrad  float64 `json:"beginningBalanceTrad"`
	BeginningBalanceRoth  float64 `json:"beginningBalanceRoth"`
	GrowthAmount          float64 `json:"growthAmount"`
	TotalWithdrawn        float64 `json:"totalWithdrawn"`
	TraditionalWithdrawn  float64 `json:"traditionalWithdrawn"`
	RothWithdrawn         float64 `json:"rothWithdrawn"`
	EndingBalanceTotal    float64 `json:"endingBalanceTotal"`
	EndingBalanceTrad     float64 `json:"endingBalanceTrad"`
	EndingBalanceRoth     float64 `json:"endingBalanceRoth"`
	RmdAmount             float64 `json:"rmdAmount,omitempty"` // Calculated RMD for the year, if applicable
	Notes                 string  `json:"notes,omitempty"`     // Additional information about this year's withdrawal
}

// TSPCalculationResult holds the output of the TSP projection and withdrawal modeling.
type TSPCalculationResult struct {
	ProjectedTraditionalBalanceAtRetirement float64                     `json:"projectedTraditionalBalanceAtRetirement"`
	ProjectedRothBalanceAtRetirement        float64                     `json:"projectedRothBalanceAtRetirement"`
	TotalProjectedBalanceAtRetirement       float64                     `json:"totalProjectedBalanceAtRetirement"`
	WithdrawalSchedule                      []YearlyTSPWithdrawalDetail `json:"withdrawalSchedule,omitempty"`
	Notes                                   string                      `json:"notes,omitempty"`
}
