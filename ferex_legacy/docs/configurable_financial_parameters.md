# Configurable Financial Parameters

## Overview

The FEREX application now supports configurable financial parameters through external configuration files. This feature allows for easy updates to financial values that change periodically (such as TSP contribution limits, L-Fund compositions, and RMD tables) without requiring code changes.

## Configuration Files

### Location

Configuration files are stored in the `config/` directory at the root of the application:

```markdown
ferex/
├── config/
│   └── tsp_config.json
```

### Format

Configuration files use JSON format for easy reading and editing. Each configuration file contains structured data specific to a particular set of financial parameters.

## TSP Configuration

### File: `config/tsp_config.json`

This file contains all configurable parameters related to TSP calculations, including contribution limits, L-Fund compositions, and RMD tables.

### Structure

```json
{
  "contributionLimits": {
    "maxAgencyMatchPercentage": 0.04,
    "agencyAutomaticContributionPercentage": 0.01,
    "standardEmployeeContributionForMaxMatch": 0.05,
    "catchUpContributionLimit": 7500,
    "maxRegularContributionLimit": 23000,
    "year": 2024
  },
  "lFundBaseCompositions": {
    "LIncome": { "g": 71.22, "f": 6.51, "c": 11.72, "s": 2.81, "i": 7.74 },
    "L2025": { "g": 69.32, "f": 6.64, "c": 12.71, "s": 3.17, "i": 8.16 },
    "L2030": { "g": 42.00, "f": 5.00, "c": 27.00, "s": 9.00, "i": 17.00 },
    // Additional L-Funds...
  },
  "lFundTargetComposition": {
    "g": 71.22, "f": 6.51, "c": 11.72, "s": 2.81, "i": 7.74
  },
  "uniformLifetimeTable": {
    "72": 27.4, "73": 26.5, "74": 25.5,
    // Additional ages...
  },
  "rmdStartAgeRules": [
    { "birthYearStart": 1900, "birthYearEnd": 1950, "startAge": 72 },
    { "birthYearStart": 1951, "birthYearEnd": 1959, "startAge": 73 },
    { "birthYearStart": 1960, "birthYearEnd": 2100, "startAge": 75 }
  ]
}
```

### Parameters

#### Contribution Limits

| Parameter | Description | Example Value |
|-----------|-------------|---------------|
| `maxAgencyMatchPercentage` | Maximum percentage of salary that agencies will match | 0.04 (4%) |
| `agencyAutomaticContributionPercentage` | Automatic agency contribution percentage | 0.01 (1%) |
| `standardEmployeeContributionForMaxMatch` | Employee contribution percentage required for maximum agency matching | 0.05 (5%) |
| `catchUpContributionLimit` | Annual catch-up contribution limit for employees age 50+ | 7500 |
| `maxRegularContributionLimit` | Annual regular contribution limit | 23000 |
| `year` | Year these limits apply to | 2024 |

#### L-Fund Compositions

The `lFundBaseCompositions` object contains the base allocation percentages for each L-Fund. Each fund has allocations across the G, F, C, S, and I funds, represented as percentages.

The `lFundTargetComposition` represents the target allocation that L-Funds move toward as they approach their target date (typically the LIncome fund allocation).

#### RMD Tables

The `uniformLifetimeTable` object maps ages to distribution factors used for Required Minimum Distribution (RMD) calculations.

The `rmdStartAgeRules` array defines the RMD start age based on birth year ranges, reflecting the SECURE Act and SECURE 2.0 Act changes.

## Updating Configuration Files

### When to Update

Configuration files should be updated when:

- Annual contribution limits change (typically announced by the IRS/TSP in October-November for the following year)
- L-Fund compositions are adjusted by the TSP
- RMD rules or tables change due to legislation or IRS updates

### How to Update

1. Open the appropriate configuration file in a text editor
2. Modify the values as needed
3. Save the file
4. Restart the application to apply the changes

### Validation

The application performs basic validation on configuration values. If invalid values are detected, the application will log an error and fall back to default values.

## For Developers

### Accessing Configuration Values

The configuration system is implemented in the `config` package. To access TSP configuration values in your code:

```go
import "ferex/backend/config"

// Get the TSP configuration
tspConfig, err := config.GetTSPConfig()
if err != nil {
    // Handle error
}

// Access values
regularLimit := tspConfig.ContributionLimits.MaxRegularContributionLimit
```

### Adding New Configuration Parameters

To add new configurable parameters:

1. Update the appropriate struct in `backend/config/tsp_config.go`
2. Add the new parameter to the JSON configuration file
3. Update the fallback values in the code
4. Use the new parameter in your calculations

### Testing Configuration Changes

A test program is available to verify configuration loading:

```bash
go run cmd/test_tsp_config/main.go
```

This will display all loaded configuration values, which can be used to verify that your changes were applied correctly.

## Troubleshooting

### Common Issues

1. **Configuration file not found**
   - Ensure the file exists in the `config/` directory
   - Check file permissions

2. **Invalid JSON format**
   - Verify the JSON syntax is correct
   - Use a JSON validator if needed

3. **Unexpected values**
   - Check that numeric values use the correct format (e.g., percentages as decimals: 0.05 for 5%)
   - Ensure all required fields are present

### Logs

If the application fails to load a configuration file, it will log an error message and fall back to default values. Check the application logs for details about configuration-related issues.

## Future Enhancements

- Admin UI for updating configuration values
- Support for multiple configuration profiles
- Automatic validation of configuration values against expected ranges
- Additional configurable parameters for other financial calculations
