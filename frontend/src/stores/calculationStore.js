import { writable, derived } from 'svelte/store';
import { appData } from './appDataStore';
import { selectedScenarioId } from './uiStore';
// Placeholder for calculation logic (to be implemented)
import { runCalculation } from '../calculations.js';

// Central store for all calculated results, keyed by scenario ID
export const calculationData = writable({});

// Automatically recalculate when the selected scenario changes
const selectedScenario = derived(
    [appData, selectedScenarioId],
    ([$appData, $selectedScenarioId]) =>
        $appData.scenarios.find(s => s.id === $selectedScenarioId) || null
);

selectedScenario.subscribe(scenario => {
    if (scenario) {
        const result = runCalculation(scenario);
        calculationData.update(d => ({ ...d, [scenario.id]: result }));
    }
});

// Derived store for current scenario's calculation result
export const currentCalculation = derived(
    [calculationData, selectedScenario],
    ([$calculationData, $selectedScenario]) =>
        $selectedScenario ? $calculationData[$selectedScenario.id] : null
);

// Helper to get a specific calculation result
export function getCalculationResult(scenarioId, calculationType) {
    return new Promise((resolve) => {
        const unsubscribe = calculationData.subscribe(results => {
            if (results[scenarioId] && results[scenarioId][calculationType]) {
                resolve(results[scenarioId][calculationType]);
            } else {
                resolve(null);
            }
            unsubscribe();
        });
    });
}

// Function to check if calculation is still valid or needs refresh
export function isCalculationValid(scenarioId, calculationType, scenarioData, maxAge = 5 * 60 * 1000) {
    return new Promise((resolve) => {
        const unsubscribe = calculationData.subscribe(results => {
            // If no results for this scenario or calculation, it's invalid
            if (!results[scenarioId] || !results[scenarioId][calculationType]) {
                resolve(false);
                unsubscribe();
                return;
            }
            
            const result = results[scenarioId][calculationType];
            
            // Check if the calculation is recent enough
            const now = new Date();
            const timestamp = new Date(result.timestamp);
            if (now - timestamp > maxAge) {
                resolve(false);
                unsubscribe();
                return;
            }
            
            // For more complex validation, you could add logic here to check 
            // if relevant input fields have changed since the last calculation
            
            resolve(true);
            unsubscribe();
        });
    });
}

// Function to explicitly store a calculation result for a scenario and type
export function storeCalculationResult(scenarioId, calculationType, result) {
  calculationData.update(d => {
    const prev = d[scenarioId] || {};
    return {
      ...d,
      [scenarioId]: {
        ...prev,
        [calculationType]: {
          ...result,
          timestamp: new Date().toISOString()
        }
      }
    };
  });
}