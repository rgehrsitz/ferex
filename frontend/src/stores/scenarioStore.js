import { derived, get } from 'svelte/store';
import { appData } from './appDataStore';
import { selectedScenarioId, compareScenarioId } from './uiStore';

// Derived view of scenarios from central appData
export const scenarios = derived(
    appData,
    $appData => $appData.scenarios
);

// Derived stores for easy access to selected and compare scenarios
export const selectedScenario = derived(
    [scenarios, selectedScenarioId],
    ([$scenarios, $selectedScenarioId]) =>
        $scenarios.find(s => s.id === $selectedScenarioId)
);

export const compareScenario = derived(
    [scenarios, compareScenarioId],
    ([$scenarios, $compareScenarioId]) =>
        $scenarios.find(s => s.id === $compareScenarioId)
);

// Helper functions for scenario management using appData
export function addScenario(newScenario) {
    // Ensure the scenario has a valid ID - use the createUniqueId function if needed
    if (!newScenario.id || newScenario.id === 0) {
        appData.update(d => {
            // Generate a unique ID based on existing scenarios
            const id = createUniqueId(d.scenarios);
            newScenario.id = id;
            
            return {
                ...d,
                scenarios: [...d.scenarios, newScenario]
            };
        });
    } else {
        // If an ID is already provided, just add the scenario
        appData.update(d => ({
            ...d,
            scenarios: [...d.scenarios, newScenario]
        }));
    }
    
    return newScenario.id;
}

export function updateScenario(updatedScenario) {
    console.log('scenarioStore.updateScenario called with:', updatedScenario);
    
    if (!updatedScenario || !updatedScenario.id) {
        console.error('Invalid scenario - missing ID:', updatedScenario);
        return;
    }
    
    // Make a deep copy of the updated scenario to prevent reference issues
    const clonedScenario = JSON.parse(JSON.stringify(updatedScenario));
    
    appData.update(d => {
        console.log('Current appData.scenarios before update:', d.scenarios);
        
        // Find the scenario we're updating
        const existingScenario = d.scenarios.find(s => s.id === clonedScenario.id);
        
        if (!existingScenario) {
            console.error(`Cannot update scenario with ID ${clonedScenario.id} - not found in store`);
            return d; // Return unmodified store if scenario not found
        }
        
        // Create a new array with the updated scenario
        const newScenarios = d.scenarios.map(s => 
            s.id === clonedScenario.id ? clonedScenario : s
        );
        
        console.log('Updated appData.scenarios after update:', newScenarios);
        
        return {
            ...d,
            scenarios: newScenarios
        };
    });
}

export function deleteScenario(id) {
    appData.update(d => ({
        ...d,
        scenarios: d.scenarios.filter(s => s.id !== id)
    }));
    selectedScenarioId.update(currentId =>
        currentId === id
            ? (get(appData).scenarios[0]?.id || null)
            : currentId
    );
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
                system: 'FERS',
                high3Salary: 100000,
                yearsOfService: 30,
                ageAtRetirement: 62,
                unusedSickLeaveMonths: 6, // Approximately 1000 hours converted to months
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
    const initial = [
        createDefaultScenario(1, 'Scenario A'),
        createDefaultScenario(2, 'Scenario B')
    ];
    appData.update(d => ({ ...d, scenarios: initial }));
    selectedScenarioId.set(initial[0]?.id || null);
    compareScenarioId.set(null);
    return initial;
}