export function createDefaultScenario(id, name) {
    const defaultData = {
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
            contributionRateRoth: 0,
            expectedReturn: 6,
            withdrawalStrategy: 'percentage',
            withdrawalRate: 4,
            fixedMonthlyWithdrawal: 0,
            withdrawalStartAge: 60
        },
        tax: {
            filingStatus: 'single',
            stateOfResidence: 'VA',
            stateIncomeTaxRate: 0.05,
            itemizedDeductions: 0,
            federalTaxCredits: 0,
            stateTaxCredits: 0,
            age: 62,
            spouseAge: 60
        },
        cola: {
            assumedInflationRate: 2.0,
            applyCOLAToPension: true,
            applyColaToSocialSecurity: true
        },
        otherIncome: {
            sources: []
        }
    };
    return {
        id,
        name,
        data: defaultData
    };
}
