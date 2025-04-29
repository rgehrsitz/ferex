package main

import (
	"context"
	"fmt"
	"ferex/backend/models"
	"ferex/backend/calculation"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// PensionInput is a minimal struct for frontend pension calculation requests
// (Add more fields as needed for future expansion)
type PensionInput struct {
	System              string  `json:"system"`
	High3Salary         float64 `json:"high3Salary"`
	YearsOfService      float64 `json:"yearsOfService"`
	AgeAtRetirement     int     `json:"ageAtRetirement"`
	UnusedSickLeaveMonths int   `json:"unusedSickLeaveMonths"`
}

// PensionResult is a minimal struct for frontend display
// (Add more fields as needed for future expansion)
type PensionResult struct {
	AnnualPension  float64 `json:"annualPension"`
	MonthlyPension float64 `json:"monthlyPension"`
	Notes          string  `json:"notes"`
}

// CalculatePension computes a FERS or CSRS pension based on user input
//export
func (a *App) CalculatePension(input PensionInput) PensionResult {
	if input.System == "FERS" {
		fersInput := models.FERSCalculationInput{
			High3Salary:           input.High3Salary,
			YearsOfService:        input.YearsOfService,
			AgeAtRetirement:       input.AgeAtRetirement,
			UnusedSickLeaveMonths: input.UnusedSickLeaveMonths,
		}
		fersResult := calculation.CalculateFERSPension(fersInput)
		return PensionResult{
			AnnualPension:  fersResult.AnnualPension,
			MonthlyPension: fersResult.MonthlyPension,
			Notes:          fersResult.Notes,
		}
	} else if input.System == "CSRS" || input.System == "CSRS Offset" {
		csrsInput := models.CSRSCalculationInput{
			High3Salary:           input.High3Salary,
			YearsOfService:        input.YearsOfService,
			UnusedSickLeaveMonths: input.UnusedSickLeaveMonths,
			IsCSRSOffset:          input.System == "CSRS Offset",
			AgeAtRetirement:       input.AgeAtRetirement,
		}
		csrsResult := calculation.CalculateCSRS(csrsInput)
		return PensionResult{
			AnnualPension:  csrsResult.AnnualPension,
			MonthlyPension: csrsResult.MonthlyPension,
			Notes:          csrsResult.Notes,
		}
	}
	return PensionResult{Notes: "Unknown retirement system."}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
