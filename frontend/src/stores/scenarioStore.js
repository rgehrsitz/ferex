import { writable, derived } from 'svelte/store';

// Main scenario store
export const scenarios = writable([]);
export const selectedScenarioId = writable(null);
export const compareScenarioId = writable(null);

// Derived stores for easy access to selected and compare scenarios
export const selectedScenario = derived(
    [scenarios, selectedScenarioId],
    ([$scenarios, $selectedScenarioId]) => {
        return $scenarios.find(s => s.id === $selectedScenarioId);
    }
);

export const compareScenario = derived(
    [scenarios, compareScenarioId],
    ([$scenarios, $compareScenarioId]) => {
        return $scenarios.find(s => s.id === $compareScenarioId);
    }
);

// Helper functions for scenario management
export function addScenario(newScenario) {
    scenarios.update(list => [...list, newScenario]);
    return newScenario.id;
}

export function updateScenario(updatedScenario) {
    scenarios.update(list => 
        list.map(scenario => 
            scenario.id === updatedScenario.id ? updatedScenario : scenario
        )
    );
}

export function deleteScenario(id) {
    scenarios.update(list => list.filter(scenario => scenario.id !== id));
    
    // Update selected ID if the deleted scenario was selected
    selectedScenarioId.update(currentId => {
        if (currentId === id) {
            // Get first available scenario, or null if none left
            return list[0]?.id || null;
        }
        return currentId;
    });
    
    // Update compare ID if the deleted scenario was being compared
    compareScenarioId.update(currentId => 
        currentId === id ? null : currentId
    );
}

// Utility function to create unique IDs
export function createUniqueId(scenarios) {
    if (scenarios.length === 0) return 1;
    return Math.max(...scenarios.map(s => s.id)) + 1;
}

// Function to create a default scenario
export function createDefaultScenario(id, name) {
    return {
        id,
        name,
        data: {
            pension: {
                retirementSystem: 'FERS',
                highThreeSalary: 100000,
                yearsOfService: 30,
                retirementAge: 62,
                unusedSickLeave: 1000,
                militaryService: 0,
                isPartTime: false,
                partTimeProrationFactor: 1.0,
                csrsOffset: false,
                survivorBenefit: 'full'
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
                contributionRateRoth: 5,
                expectedReturn: 6,
                withdrawalStrategy: 'fixed',
                fixedMonthlyWithdrawal: 2000,
                withdrawalRate: 4,
                withdrawalStartAge: 62
            },
            tax: {
                filingStatus: 'married_joint',
                stateOfResidence: 'VA',
                additionalIncome: 0,
                additionalDeductions: 0,
                stateIncomeTaxRate: 0.05
            },
            cola: {
                assumedInflationRate: 0.025,
                applyCOLAToPension: true,
                applyColaToSocialSecurity: true
            },
            otherIncome: {
                sources: [
                    {
                        id: crypto.randomUUID(),
                        name: 'Part-time work',
                        amount: 1200,
                        frequency: 'monthly',
                        startAge: 62,
                        endAge: 70,
                        applyCola: true
                    },
                    {
                        id: crypto.randomUUID(),
                        name: 'Rental property',
                        amount: 1800,
                        frequency: 'monthly',
                        startAge: 62,
                        endAge: 90,
                        applyCola: true
                    }
                ]
            }
        }
    };
}

// Load initial scenarios
export function initializeStore() {
    const initialScenarios = [
        createDefaultScenario(1, "Scenario A"),
        createDefaultScenario(2, "Scenario B")
    ];
    
    scenarios.set(initialScenarios);
    selectedScenarioId.set(1);
    compareScenarioId.set(null);
    
    return initialScenarios;
}