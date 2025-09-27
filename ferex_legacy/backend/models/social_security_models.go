package models

// SocialSecurityCreditedEarning corresponds to an element in TypeScript SocialSecurityCreditedEarnings array
type SocialSecurityCreditedEarning struct {
	Year     int     `json:"year"`
	Earnings float64 `json:"earnings"`
}
