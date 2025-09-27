package main

import (
	"ferex/backend/config"
	"fmt"
	"log"
)

func main() {
	// Load the TSP configuration
	tspConfig, err := config.GetTSPConfig()
	if err != nil {
		log.Fatalf("Error loading TSP configuration: %v", err)
	}

	// Display the contribution limits
	fmt.Println("TSP Contribution Limits:")
	fmt.Printf("  Year: %d\n", tspConfig.ContributionLimits.Year)
	fmt.Printf("  Max Regular Contribution Limit: $%.2f\n", tspConfig.ContributionLimits.MaxRegularContributionLimit)
	fmt.Printf("  Catch-Up Contribution Limit: $%.2f\n", tspConfig.ContributionLimits.CatchUpContributionLimit)
	fmt.Printf("  Agency Automatic Contribution: %.2f%%\n", tspConfig.ContributionLimits.AgencyAutomaticContributionPercentage*100)
	fmt.Printf("  Max Agency Matching: %.2f%%\n", tspConfig.ContributionLimits.MaxAgencyMatchPercentage*100)
	fmt.Printf("  Standard Employee Contribution for Max Match: %.2f%%\n", tspConfig.ContributionLimits.StandardEmployeeContributionForMaxMatch*100)

	// Display L-Fund compositions
	fmt.Println("\nL-Fund Base Compositions:")
	for fundName, composition := range tspConfig.LFundBaseCompositions {
		fmt.Printf("  %s:\n", fundName)
		if composition.G != nil {
			fmt.Printf("    G: %.2f%%\n", *composition.G)
		}
		if composition.F != nil {
			fmt.Printf("    F: %.2f%%\n", *composition.F)
		}
		if composition.C != nil {
			fmt.Printf("    C: %.2f%%\n", *composition.C)
		}
		if composition.S != nil {
			fmt.Printf("    S: %.2f%%\n", *composition.S)
		}
		if composition.I != nil {
			fmt.Printf("    I: %.2f%%\n", *composition.I)
		}
	}

	// Display L-Fund target composition
	fmt.Println("\nL-Fund Target Composition:")
	if tspConfig.LFundTargetComposition.G != nil {
		fmt.Printf("  G: %.2f%%\n", *tspConfig.LFundTargetComposition.G)
	}
	if tspConfig.LFundTargetComposition.F != nil {
		fmt.Printf("  F: %.2f%%\n", *tspConfig.LFundTargetComposition.F)
	}
	if tspConfig.LFundTargetComposition.C != nil {
		fmt.Printf("  C: %.2f%%\n", *tspConfig.LFundTargetComposition.C)
	}
	if tspConfig.LFundTargetComposition.S != nil {
		fmt.Printf("  S: %.2f%%\n", *tspConfig.LFundTargetComposition.S)
	}
	if tspConfig.LFundTargetComposition.I != nil {
		fmt.Printf("  I: %.2f%%\n", *tspConfig.LFundTargetComposition.I)
	}

	// Display RMD tables
	fmt.Println("\nRMD Uniform Lifetime Table:")
	uniformLifetimeTable := tspConfig.GetUniformLifetimeTableAsMap()
	for age := 72; age <= 120; age++ {
		if factor, ok := uniformLifetimeTable[age]; ok {
			fmt.Printf("  Age %d: %.1f\n", age, factor)
		}
	}

	// Display RMD start age rules
	fmt.Println("\nRMD Start Age Rules:")
	for _, rule := range tspConfig.RMDStartAgeRules {
		fmt.Printf("  Birth Years %d-%d: Start Age %d\n", rule.BirthYearStart, rule.BirthYearEnd, rule.StartAge)
	}

	// Test getting RMD start age for a specific birth year
	testBirthYear := 1960
	fmt.Printf("\nRMD Start Age for Birth Year %d: %d\n", testBirthYear, tspConfig.GetRMDStartAge(testBirthYear))
}
