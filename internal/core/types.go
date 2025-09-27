package core

import "fmt"

// RetirementSystem identifies the federal retirement system rules that apply to a
// profile and scenario. The MVP focuses on standard FERS and CSRS employees with
// CSRS Offset treated as a hybrid case. Special categories (e.g. LEO/FF/ATC) are
// out of scope for v1 but the enum leaves room for future expansion.
type RetirementSystem string

const (
	RetirementSystemFERS       RetirementSystem = "FERS"
	RetirementSystemCSRS       RetirementSystem = "CSRS"
	RetirementSystemCSRSOffset RetirementSystem = "CSRS_OFFSET"
)

// Validate ensures the retirement system is one of the supported values.
func (r RetirementSystem) Validate() error {
	switch r {
	case RetirementSystemFERS, RetirementSystemCSRS, RetirementSystemCSRSOffset:
		return nil
	default:
		return fmt.Errorf("unsupported retirement system %q", string(r))
	}
}

// ServiceType captures the nature of a service period so that eligibility,
// deposits, and proration rules can be modeled correctly.
type ServiceType string

const (
	ServiceTypeFERSCIVIL       ServiceType = "FERS_CIVILIAN"
	ServiceTypeCSRSCIVIL       ServiceType = "CSRS_CIVILIAN"
	ServiceTypeMilitary        ServiceType = "MILITARY"
	ServiceTypePartTime        ServiceType = "PART_TIME"
	ServiceTypeLeaveWithoutPay ServiceType = "LWOP"
)

// PayType indicates the compensation basis for high-3 calculations.
type PayType string

const (
	PayTypeBaseSalary PayType = "BASE"
	PayTypeOvertime   PayType = "OVERTIME"
	PayTypeBonus      PayType = "BONUS"
)

// PaymentFrequency identifies how often an income stream pays out. This is used
// for both pension payments and TSP withdrawals.
type PaymentFrequency string

const (
	PaymentFrequencyMonthly   PaymentFrequency = "MONTHLY"
	PaymentFrequencyQuarterly PaymentFrequency = "QUARTERLY"
	PaymentFrequencyAnnual    PaymentFrequency = "ANNUAL"
	PaymentFrequencySingle    PaymentFrequency = "ONCE"
)

// TaxFilingStatus enumerates the primary filing statuses supported by the tax
// module so that marginal brackets can be applied appropriately.
type TaxFilingStatus string

const (
	TaxFilingSingle          TaxFilingStatus = "SINGLE"
	TaxFilingMarriedJoint    TaxFilingStatus = "MARRIED_JOINT"
	TaxFilingMarriedSeparate TaxFilingStatus = "MARRIED_SEPARATE"
	TaxFilingHeadOfHousehold TaxFilingStatus = "HEAD_OF_HOUSEHOLD"
)

// WithdrawalStrategy identifies how a TSP installment withdrawal is determined.
type WithdrawalStrategy string

const (
	WithdrawalStrategyFixedAmount    WithdrawalStrategy = "FIXED_AMOUNT"
	WithdrawalStrategyLifeExpectancy WithdrawalStrategy = "LIFE_EXPECTANCY"
)

// WithdrawalSource determines how TSP withdrawals are sourced across Traditional
// and Roth balances.
type WithdrawalSource string

const (
	WithdrawalSourceProportional     WithdrawalSource = "PROPORTIONAL"
	WithdrawalSourceTraditionalFirst WithdrawalSource = "TRADITIONAL_FIRST"
	WithdrawalSourceRothFirst        WithdrawalSource = "ROTH_FIRST"
)

// CashflowKind differentiates the modeled cashflows so reports and analysis can
// aggregate them consistently.
type CashflowKind string

const (
	CashflowKindPension        CashflowKind = "PENSION"
	CashflowKindSupplement     CashflowKind = "SUPPLEMENT"
	CashflowKindSocialSecurity CashflowKind = "SOCIAL_SECURITY"
	CashflowKindTSPWithdrawal  CashflowKind = "TSP_WITHDRAWAL"
	CashflowKindOtherIncome    CashflowKind = "OTHER_INCOME"
)
