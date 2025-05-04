import { writable, derived } from 'svelte/store';
import { selectedScenario } from './scenarioStore';

// Store for calculation results
export const calculationResults = writable({
    // Keyed by scenario id
});

// Helper to store results
export function storeCalculationResult(scenarioId, calculationType, result) {
    calculationResults.update(results => {
        // Initialize scenario results object if it doesn't exist
        if (!results[scenarioId]) {
            results[scenarioId] = {};
        }
        
        // Store the calculation result
        results[scenarioId][calculationType] = {
            ...result,
            timestamp: new Date()
        };
        
        return results;
    });
}

// Helper to get a specific calculation result
export function getCalculationResult(scenarioId, calculationType) {
    return new Promise((resolve) => {
        const unsubscribe = calculationResults.subscribe(results => {
            if (results[scenarioId] && results[scenarioId][calculationType]) {
                resolve(results[scenarioId][calculationType]);
            } else {
                resolve(null);
            }
            unsubscribe();
        });
    });
}

// Derived store for current scenario's results
export const currentResults = derived(
    [calculationResults, selectedScenario],
    ([$calculationResults, $selectedScenario]) => {
        if (!$selectedScenario) return null;
        return $calculationResults[$selectedScenario.id] || null;
    }
);

// Function to check if calculation is still valid or needs refresh
export function isCalculationValid(scenarioId, calculationType, scenarioData, maxAge = 5 * 60 * 1000) {
    return new Promise((resolve) => {
        const unsubscribe = calculationResults.subscribe(results => {
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