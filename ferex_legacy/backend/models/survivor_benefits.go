package models

// SurvivorBenefitFERSInput corresponds to TypeScript SurvivorBenefitFERSInput interface
type SurvivorBenefitFERSInput struct {
	SpouseElection                *string                   `json:"spouseElection,omitempty"`
	FormerSpouseElection          *string                   `json:"formerSpouseElection,omitempty"`
	FormerSpouseConsent           *string                   `json:"formerSpouseConsent,omitempty"`
	CurrentSpouseConsentForFormer *string                   `json:"currentSpouseConsentForFormer,omitempty"`
	CurrentSpouseWaiverForSelf    *string                   `json:"currentSpouseWaiverForSelf,omitempty"`
	InsurableInterestDetails      *InsurableInterestDetails `json:"insurableInterestDetails,omitempty"`
}

// SurvivorBenefitCSRSInput corresponds to TypeScript SurvivorBenefitCSRSInput interface
type SurvivorBenefitCSRSInput struct {
	Election                      *string                   `json:"election,omitempty"`
	CustomBaseAmountForPartial    *float64                  `json:"customBaseAmountForPartial,omitempty"`
	FormerSpouseElection          *string                   `json:"formerSpouseElection,omitempty"`
	FormerSpouseCustomBaseAmount  *float64                  `json:"formerSpouseCustomBaseAmount,omitempty"`
	FormerSpouseConsent           *string                   `json:"formerSpouseConsent,omitempty"`
	CurrentSpouseConsentForFormer *string                   `json:"currentSpouseConsentForFormer,omitempty"`
	CurrentSpouseWaiverForSelf    *string                   `json:"currentSpouseWaiverForSelf,omitempty"`
	InsurableInterestDetails      *InsurableInterestDetails `json:"insurableInterestDetails,omitempty"`
}

// FormerSpouseSurvivorDetails placeholder for frontend's object type.
// Specific fields to be determined if/when frontend provides them.
type FormerSpouseSurvivorDetails struct {
	CourtOrderNumber    *string  `json:"courtOrderNumber,omitempty"`
	DateOfCourtOrder    *string  `json:"dateOfCourtOrder,omitempty"`    // YYYY-MM-DD
	PortionAwarded      *string  `json:"portionAwarded,omitempty"`      // 'full' | 'partial' | 'none' | 'pro-rata'
	CustomAmountAwarded *float64 `json:"customAmountAwarded,omitempty"` // If portion is 'partial' and a specific amount
}
