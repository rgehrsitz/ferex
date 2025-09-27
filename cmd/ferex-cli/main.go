package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ferexapp/ferex/internal/app"
	"github.com/ferexapp/ferex/internal/core"
	filestore "github.com/ferexapp/ferex/internal/storage/file"
)

var (
	storeDir string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ferex",
		Short: "Ferex CLI for managing retirement scenarios",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if storeDir == "" {
				return errors.New("--store directory is required")
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVar(&storeDir, "store", "scenarios", "directory to store .ferex files")

	rootCmd.AddCommand(newProfileTemplateCmd())
	rootCmd.AddCommand(newScenarioCreateCmd())
	rootCmd.AddCommand(newScenarioListCmd())
	rootCmd.AddCommand(newScenarioShowCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newProfileTemplateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "profile-template",
		Short: "Print a starter UserProfile JSON document",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile := core.UserProfile{
				ID:               "PROFILE-ID-REPLACE",
				DisplayName:      "Sample User",
				BirthDate:        time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				HireDate:         time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
				RetirementSystem: core.RetirementSystemFERS,
				ServiceHistory: []core.ServicePeriod{
					{
						StartDate:    time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
						EndDate:      time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
						Type:         core.ServiceTypeFERSCIVIL,
						HoursPerWeek: 40,
						DepositPaid:  true,
					},
				},
				Compensation: []core.CompensationEntry{
					{
						StartDate: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
						PayType:   core.PayTypeBaseSalary,
						Amount:    125000,
					},
				},
				SocialSecurity: core.SocialSecurityRecord{
					Age62Estimate:           2200,
					FullRetirementAgeMonths: 804,
					LastUpdated:             time.Now().UTC(),
				},
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(profile)
		},
	}
}

func newScenarioCreateCmd() *cobra.Command {
	var (
		profilePath    string
		scenarioName   string
		retirementDate string
		supplement     bool
	)
	cmd := &cobra.Command{
		Use:   "scenario-create",
		Short: "Create a scenario from a profile JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if profilePath == "" {
				return errors.New("--profile is required")
			}
			if scenarioName == "" {
				return errors.New("--name is required")
			}
			if retirementDate == "" {
				return errors.New("--retirement-date is required")
			}
			data, err := os.ReadFile(profilePath)
			if err != nil {
				return err
			}
			var profile core.UserProfile
			if err := json.Unmarshal(data, &profile); err != nil {
				return fmt.Errorf("parse profile: %w", err)
			}
			date, err := time.Parse(time.DateOnly, retirementDate)
			if err != nil {
				return fmt.Errorf("parse retirement date: %w", err)
			}
			store, err := filestore.NewScenarioStore(storeDir)
			if err != nil {
				return err
			}
			manager := app.NewScenarioManager(store, nil)
			scenario, err := manager.CreateScenario(context.Background(), profile, scenarioName, func(s *core.Scenario) {
				s.RetirementDate = date
				s.IncludeSupplement = supplement
				if s.TSPPlan.Withdrawal.Strategy == "" {
					s.TSPPlan.Withdrawal.Strategy = core.WithdrawalStrategyFixedAmount
					s.TSPPlan.Withdrawal.Frequency = core.PaymentFrequencyMonthly
					s.TSPPlan.Withdrawal.FixedAmount = 2000
					s.TSPPlan.Withdrawal.Source = core.WithdrawalSourceProportional
				}
				if len(s.TaxSettings.FederalBrackets) == 0 {
					s.TaxSettings = defaultTaxSettings()
				}
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created scenario %s (%s)\n", scenario.ID, scenario.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Saved to %s\n", filepath.Join(storeDir, scenario.ID+".ferex"))
			return nil
		},
	}
	cmd.Flags().StringVar(&profilePath, "profile", "", "path to profile JSON file")
	cmd.Flags().StringVar(&scenarioName, "name", "", "scenario name")
	cmd.Flags().StringVar(&retirementDate, "retirement-date", "", "retirement date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&supplement, "supplement", true, "include FERS supplement")
	return cmd
}

func newScenarioListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scenario-list",
		Short: "List saved scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := filestore.NewScenarioStore(storeDir)
			if err != nil {
				return err
			}
			manager := app.NewScenarioManager(store, nil)
			items, err := manager.ListScenarios(context.Background())
			if err != nil {
				return err
			}
			if len(items) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No scenarios found.")
				return nil
			}
			for _, item := range items {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\n", item.ID, item.Name, item.RetirementDate.Format(time.DateOnly), item.UpdatedAt.Format(time.RFC3339))
			}
			return nil
		},
	}
}

func newScenarioShowCmd() *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "scenario-show",
		Short: "Show a scenario JSON payload",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return errors.New("--id is required")
			}
			store, err := filestore.NewScenarioStore(storeDir)
			if err != nil {
				return err
			}
			manager := app.NewScenarioManager(store, nil)
			scenario, err := manager.GetScenario(context.Background(), id)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(scenario)
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "scenario id")
	return cmd
}

func defaultTaxSettings() core.TaxSettings {
	brackets := []core.TaxBracket{
		{Threshold: 0, Rate: 0.1},
		{Threshold: 11000, Rate: 0.12},
		{Threshold: 44725, Rate: 0.22},
		{Threshold: 95375, Rate: 0.24},
		{Threshold: 182100, Rate: 0.32},
		{Threshold: 231250, Rate: 0.35},
		{Threshold: 578125, Rate: 0.37},
	}
	return core.TaxSettings{
		FilingStatus:      core.TaxFilingSingle,
		StandardDeduction: 14600,
		FederalBrackets:   brackets,
		TaxFreeBasis:      0,
		StateName:         "",
		StateRate:         0,
	}
}

func init() {
	cobra.OnInitialize(func() {
		if storeDir != "" {
			if err := os.MkdirAll(storeDir, 0o755); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to prepare store directory: %v\n", err)
				os.Exit(1)
			}
		}
	})
}
