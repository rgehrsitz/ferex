package core

// TSPPlan models both the accumulation and decumulation assumptions for the
// Thrift Savings Plan. It separates Traditional and Roth balances so tax logic
// can determine taxable versus tax-free withdrawals.
type TSPPlan struct {
	TraditionalBalance  float64           `json:"traditional_balance"`
	RothBalance         float64           `json:"roth_balance"`
	ExpectedReturn      float64           `json:"expected_return"`
	ExpectedVolatility  float64           `json:"expected_volatility"`
	ContributionPercent float64           `json:"contribution_percent"`
	AgencyMatchPercent  float64           `json:"agency_match_percent"`
	Withdrawal          TSPWithdrawalPlan `json:"withdrawal"`
	RMDStartAge         int               `json:"rmd_start_age"`
	Notes               string            `json:"notes,omitempty"`
}

// TSPWithdrawalPlan defines how the participant intends to take distributions
// once retired. It mirrors the combinations supported by TSP: fixed amount or
// life expectancy installments, optional partial lump sums, and the ability to
// steer the source of funds.
type TSPWithdrawalPlan struct {
	Strategy       WithdrawalStrategy `json:"strategy"`
	FixedAmount    float64            `json:"fixed_amount"`
	Frequency      PaymentFrequency   `json:"frequency"`
	Source         WithdrawalSource   `json:"source"`
	PartialLumpSum float64            `json:"partial_lump_sum"`
	LumpSumAge     int                `json:"lump_sum_age"`
}

// Validate confirms the strategy parameters are consistent.
func (p TSPPlan) Validate() *ValidationError {
	var verr ValidationError
	if p.TraditionalBalance < 0 {
		verr.Add("traditional_balance", "cannot be negative")
	}
	if p.RothBalance < 0 {
		verr.Add("roth_balance", "cannot be negative")
	}
	if p.Withdrawal.Strategy == "" {
		verr.Add("withdrawal.strategy", "strategy is required")
	}
	if p.Withdrawal.Strategy == WithdrawalStrategyFixedAmount && p.Withdrawal.FixedAmount <= 0 {
		verr.Add("withdrawal.fixed_amount", "fixed amount must be positive")
	}
	if p.Withdrawal.Frequency == "" {
		verr.Add("withdrawal.frequency", "frequency is required")
	}
	if verr.Empty() {
		return nil
	}
	return &verr
}
