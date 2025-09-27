package models

// ServicePeriod corresponds to TypeScript ServicePeriod interface
type ServicePeriod struct {
	ID                            string   `json:"id"` // auto-generated unique identifier
	ServiceCategory               string   `json:"serviceCategory"`               // 'Civilian' | 'Military'
	CivilianServiceType           *string  `json:"civilianServiceType,omitempty"` // Applicable if serviceCategory is 'Civilian'. e.g., 'RegularDeductionFERS', etc.
	MilitaryServiceType           *string  `json:"militaryServiceType,omitempty"` // Applicable if serviceCategory is 'Military'. e.g., 'Post1956MilitaryServiceDepositPaid', etc.
	DepositRedepositPaymentStatus string   `json:"depositRedepositPaymentStatus"` // 'PaidInFull' | 'OwedOrPartiallyPaid' | 'NotApplicable' | 'AwaitingDetermination'
	SystemDuringService           string   `json:"systemDuringService"`           // 'FERS' | 'CSRS' | 'CSRS_Offset' | 'SocialSecurityOnly' | 'None' | 'Other'
	StartDate                     string   `json:"startDate"`                     // YYYY-MM-DD format
	EndDate                       string   `json:"endDate"`                       // YYYY-MM-DD format
	IsPartTime                    bool     `json:"isPartTime"`                    // Indicates if the service was part-time
	HoursPerWeekIfPartTime        *float64 `json:"hoursPerWeekIfPartTime,omitempty"` // Hours worked per week if isPartTime is true
	Notes                         *string  `json:"notes,omitempty"`                 // Optional user-provided notes
}

// LWOPPeriod corresponds to TypeScript LWOPPeriod interface
type LWOPPeriod struct {
	ID        string `json:"id"`
	StartDate string `json:"startDate"` // YYYY-MM-DD
	EndDate   string `json:"endDate"`   // YYYY-MM-DD
	Type      string `json:"type"`      // 'PersonalNonMilitary' | 'MilitaryLWOP' | 'OWCP'
}

// InsurableInterestDetails corresponds to TypeScript InsurableInterestDetails interface
type InsurableInterestDetails struct {
	RelationshipToEmployee string `json:"relationshipToEmployee"` // e.g., Parent, Sibling, Financially Dependent Child
	DateOfBirth            string `json:"dateOfBirth"`            // YYYY-MM-DD
}

// TaxableAccountAssetAllocation corresponds to TypeScript TaxableAccountAssetAllocation interface
type TaxableAccountAssetAllocation struct {
	StocksPercent *float64 `json:"stocksPercent,omitempty"`
	BondsPercent  *float64 `json:"bondsPercent,omitempty"`
	CashPercent   *float64 `json:"cashPercent,omitempty"`
}

// OtherRecurringIncomeSource corresponds to TypeScript OtherRecurringIncomeSource interface
type OtherRecurringIncomeSource struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Amount              *float64 `json:"amount,omitempty"` // Changed to omitempty as amount can be null in TS
	Frequency           string   `json:"frequency"`             // 'Monthly' | 'Annually'
	StartDate           string   `json:"startDate"`             // YYYY-MM-DD
	EndDate             *string  `json:"endDate,omitempty"`     // YYYY-MM-DD, optional
	IsInflationAdjusted *bool    `json:"isInflationAdjusted,omitempty"`
	AnnualIncreaseRate  *float64 `json:"annualIncreaseRate,omitempty"` // e.g., 0.02 for 2%
	SubjectToFICA       *bool    `json:"subjectToFICA,omitempty"`
}

// OneTimeIncomeEvent corresponds to TypeScript OneTimeIncomeEvent interface
type OneTimeIncomeEvent struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Amount *float64 `json:"amount,omitempty"` // Changed to omitempty
	Date   string   `json:"date"` // YYYY-MM-DD
	Type   string   `json:"type"` // 'Inheritance' | 'SaleOfAsset' | 'Bonus' | 'Other'
}

// OneTimeExpenseEvent corresponds to TypeScript OneTimeExpenseEvent interface
type OneTimeExpenseEvent struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Amount *float64 `json:"amount,omitempty"` // Changed to omitempty
	Date   string   `json:"date"` // YYYY-MM-DD
	Type   string   `json:"type"` // 'HomePurchase' | 'VehiclePurchase' | 'Education' | 'Travel' | 'Medical' | 'Other'
}
