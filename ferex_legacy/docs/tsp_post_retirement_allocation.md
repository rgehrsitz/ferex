# TSP Post-Retirement Allocation Strategy

## Overview

The TSP Post-Retirement Allocation Strategy feature allows users to define a different asset allocation for their TSP funds after retirement. This is a common financial planning practice as many retirees prefer to de-risk their portfolios during retirement to reduce volatility and preserve capital.

## Why This Feature Matters

Many financial advisors recommend adjusting your investment allocation when transitioning from the accumulation phase (pre-retirement) to the distribution phase (post-retirement). Typically, this involves:

- Reducing exposure to higher-risk funds (S, I, and sometimes C funds)
- Increasing allocation to more stable funds (G and F funds)
- Potentially switching to more conservative L Funds (e.g., from L2050 to L Income)

This feature gives users the flexibility to model these changes and see how different post-retirement allocations might affect their retirement income and portfolio longevity.

## How to Use

1. Set your pre-retirement allocation using the `CurrentAllocation` field
2. Define your post-retirement allocation using the new `PostRetirementAllocation` field
3. Set `FutureAllocationStrategy` to "UsePostRetirementAllocation"

## Implementation Details

### Model Changes

The following fields have been added to the `TSPCalculationInput` struct:

```go
PostRetirementAllocation TSPFundAllocationPercentages `json:"postRetirementAllocation,omitempty"` // Allocation to use after retirement
```

And the `FutureAllocationStrategy` field now accepts a new value:

```go
FutureAllocationStrategy string `json:"futureAllocationStrategy,omitempty"` // "MaintainCurrent", "MatchContributionAllocation", "UsePostRetirementAllocation"
```

### Available Strategies

1. **MaintainCurrent**: Uses the current allocation throughout pre-retirement and post-retirement
2. **MatchContributionAllocation**: Uses the contribution allocation throughout pre-retirement and post-retirement
3. **UsePostRetirementAllocation**: Uses the current allocation during pre-retirement, then switches to the post-retirement allocation after retirement

### Growth Calculation

The system will automatically switch to using the post-retirement allocation for growth calculations once the user reaches retirement. This transition happens seamlessly in the withdrawal phase calculations.

## Example

```json
{
  "currentAllocation": {
    "c": 40,
    "s": 30,
    "i": 30
  },
  "postRetirementAllocation": {
    "g": 30,
    "f": 20,
    "c": 30,
    "s": 10,
    "i": 10
  },
  "futureAllocationStrategy": "UsePostRetirementAllocation"
}
```

In this example, the user has an aggressive pre-retirement allocation (0% G/F funds) but plans to switch to a more conservative allocation after retirement (50% G/F funds).

## Notes

- If `UsePostRetirementAllocation` is selected but no post-retirement allocation is provided, the system will default to using the current allocation
- When using an L Fund in the post-retirement allocation, the system will still apply the appropriate glide path adjustments based on the projection year
- The withdrawal detail records will include a note when the post-retirement allocation is being used for growth calculations
