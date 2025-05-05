// Utility to create a default Scenario object
// Types are imported from the TypeScript definitions for compatibility
// If you are using TypeScript, you can convert this file to .ts

export function createDefaultScenario(id, name) {
  return {
    id,
    name,
    data: {
      pension: {
        system: 'FERS',
        high3Salary: 100000,
        yearsOfService: 30,
        ageAtRetirement: 62,
        unusedSickLeaveMonths: 6,
        survivorBenefitOption: 'full',
        isPartTime: false,
        partTimeProrationFactor: 1.0,
        militaryService: 0
      },
      socialSecurity: {
        startAge: 62,
        estimatedMonthlyBenefit: 2000,
        isEligible: true,
        birthYear: 1970,
        birthMonth: 1
      },
      tsp: {
        traditionalBalance: 400000,
        rothBalance: 100000,
        contributionRate: 5,
        catchUpContribution: 0,
        expectedAnnualReturn: 6
      },
      tax: {
        filingStatus: 'single',
        state: 'VA',
        standardDeduction: true,
        itemizedDeductions: 0
      },
      cola: {
        baseYear: 2025,
        annualCOLAPercent: 2.0
      },
      otherIncome: {
        description: '',
        amount: 0,
        startAge: 62,
        endAge: 67
      }
    }
  };
}
