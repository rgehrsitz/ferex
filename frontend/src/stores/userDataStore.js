import { writable, derived, get } from 'svelte/store';
import { selectedScenario } from './scenarioStore';

// User demographic data store for sharing between components
export const userData = writable({
    birthYear: null,
    birthMonth: null,
    retirementAge: null
});

// Update userData when relevant parts of selected scenario change
export const syncWithScenario = derived(
    selectedScenario,
    ($selectedScenario) => {
        if ($selectedScenario) {
            userData.update(data => ({
                ...data,
                birthYear: $selectedScenario.data.socialSecurity?.birthYear || data.birthYear,
                birthMonth: $selectedScenario.data.socialSecurity?.birthMonth || data.birthMonth,
                retirementAge: $selectedScenario.data.pension?.retirementAge || data.retirementAge
            }));
        }
        return $selectedScenario;
    }
);

// Helper function to update user data directly
export function updateUserData(newData) {
    userData.update(currentData => ({
        ...currentData,
        ...newData
    }));
}

// Helper to get current user data
export function getUserData() {
    return get(userData);
}