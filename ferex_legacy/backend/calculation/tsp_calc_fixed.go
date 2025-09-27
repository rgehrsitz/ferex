package calculation

// This is a proposed fix for the TSP withdrawal calculation order of operations issue
// The current implementation withdraws money BEFORE calculating growth, which means
// a 2% withdrawal rate requires more than 2% growth to maintain balance.

// CURRENT PROBLEMATIC ORDER (in the existing code):
// 1. Calculate targetAnnualWithdrawal = detail.BeginningBalanceTotal * input.WithdrawalPercentageValue
// 2. Subtract withdrawal: tradBalanceForWithdrawal -= withdrawnTrad
// 3. Calculate growth on reduced balance: tradGrowth := tradBalanceForWithdrawal * postRetirementReturn
// 4. Add growth: tradBalanceForWithdrawal += tradGrowth
// 5. Subtract expense ratio: tradBalanceForWithdrawal *= (1.0 - input.ExpenseRatio)

// PROPOSED CORRECTED ORDER:
// 1. Calculate targetAnnualWithdrawal = detail.BeginningBalanceTotal * input.WithdrawalPercentageValue
// 2. Calculate growth on FULL beginning balance: tradGrowth := tradBalanceForWithdrawal * postRetirementReturn
// 3. Add growth: tradBalanceForWithdrawal += tradGrowth
// 4. Subtract withdrawal: tradBalanceForWithdrawal -= withdrawnTrad
// 5. Subtract expense ratio: tradBalanceForWithdrawal *= (1.0 - input.ExpenseRatio)

// ALTERNATIVE APPROACH (more realistic):
// Calculate growth for the portion of the year before withdrawal, then withdraw mid-year
// This would require more complex timing assumptions but would be more accurate

// The key insight is that with the current order:
// - 2% withdrawal + 2% growth = net loss due to order of operations
// - You need approximately 2.04% growth to break even with 2% withdrawals
// - This explains why TSP withdrawals appear to decline over time even with reasonable return assumptions

// RECOMMENDATION:
// Replace the withdrawal simulation loop in CalculateTSP() with the corrected order of operations
// This will make the TSP projections much more realistic and align with user expectations
