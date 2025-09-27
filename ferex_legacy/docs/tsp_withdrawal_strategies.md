# TSP Withdrawal Strategies

This document explains the various TSP withdrawal strategies implemented in the FEREX application.

## Available Withdrawal Strategies

The application supports the following withdrawal strategies:

1. **FixedAmountYearly**: Withdraws a fixed dollar amount each year
2. **FixedAmountMonthly**: Withdraws a fixed dollar amount each month (calculated as annual amount)
3. **PercentageOfBalanceYearly**: Withdraws a percentage of the current balance each year
4. **IRSMinimumRequiredDistribution**: Withdraws the IRS Required Minimum Distribution (RMD) amount
5. **InflationAdjustedFixedAmount**: Withdraws a fixed amount that increases with inflation each year

## Withdrawal Order Options

For each withdrawal strategy, you can specify the order in which funds are withdrawn:

1. **ProRata**: Withdraws proportionally from Traditional and Roth accounts based on their relative balances
2. **TraditionalFirst**: Withdraws from Traditional TSP first, then Roth TSP when Traditional is depleted
3. **RothFirst**: Withdraws from Roth TSP first, then Traditional TSP when Roth is depleted

## Strategy Details

### FixedAmountYearly

This strategy withdraws a fixed dollar amount each year until funds are depleted or the projection period ends.

- **Input Parameter**: `WithdrawalFixedAmountValue` - The annual amount to withdraw
- **Considerations**: This strategy does not adjust for inflation, which means the purchasing power of the withdrawal amount decreases over time

### FixedAmountMonthly

Similar to FixedAmountYearly, but the input is a monthly amount that gets converted to an annual amount for the calculation.

- **Input Parameter**: `WithdrawalFixedAmountValue` - The monthly amount to withdraw
- **Calculation**: Annual withdrawal = Monthly amount × 12

### PercentageOfBalanceYearly

Withdraws a percentage of the current balance each year. This strategy automatically adjusts the withdrawal amount as the balance changes.

- **Input Parameter**: `WithdrawalPercentageValue` - The percentage of the balance to withdraw annually (e.g., 4%)
- **Calculation**: Annual withdrawal = Current balance × (Percentage / 100)
- **Considerations**: This strategy can help extend the life of the portfolio since withdrawals decrease as the balance decreases

### IRSMinimumRequiredDistribution

Withdraws the IRS Required Minimum Distribution (RMD) amount, which is calculated based on the account balance and the participant's age.

- **Calculation**: RMD = Account balance ÷ Distribution period factor
- **Applicability**: Only applies when the participant reaches the RMD age (currently 73 for those born between 1951-1959, and 75 for those born in 1960 or later)
- **Notes**: The distribution period factors come from the IRS Uniform Lifetime Table and are configured in the application

### InflationAdjustedFixedAmount

Withdraws a fixed amount that increases with inflation each year, helping to maintain purchasing power over time.

- **Input Parameters**: 
  - `WithdrawalFixedAmountValue` - The initial annual amount to withdraw
  - `ExpectedAnnualInflationRate` - The expected annual inflation rate
- **Calculation**: Annual withdrawal = Initial amount × (1 + Inflation rate)^year_number

## Implementation Notes

- All withdrawal strategies ensure that the withdrawal amount never exceeds the available balance
- For the ProRata withdrawal order, if one account is depleted, the remaining needed amount is taken from the other account
- The application tracks and reports the amount withdrawn from each account type (Traditional and Roth)
- RMD calculations follow the latest IRS guidelines and use the appropriate distribution period factors based on the participant's birth year
