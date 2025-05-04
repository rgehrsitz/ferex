import { writable } from 'svelte/store';
import { 
    CalculatePension, 
    CalculateSocialSecurity, 
    CalculateTSPProjection,
    CalculateTaxLiability,
    CalculateCOLAAdjustment,
    CalculateRetirementProjection,
    SaveScenarios,
    LoadScenarios,
    GetSavedScenarios
} from '../../wailsjs/go/main/App';

// Store for tracking API call status
export const apiStatus = writable({
    loading: false,
    error: null,
    lastUpdated: null
});

// Helper to wrap API calls with loading and error handling
export async function apiCall(fn, ...args) {
    apiStatus.update(status => ({ ...status, loading: true, error: null }));
    
    try {
        const result = await fn(...args);
        apiStatus.update(status => ({ 
            ...status, 
            loading: false, 
            lastUpdated: new Date() 
        }));
        return result;
    } catch (error) {
        console.error('API call failed:', error);
        apiStatus.update(status => ({ 
            ...status, 
            loading: false, 
            error: error.toString() 
        }));
        throw error;
    }
}

// API function wrappers with loading status
export const api = {
    calculatePension: (input) => apiCall(CalculatePension, input),
    calculateSocialSecurity: (input) => apiCall(CalculateSocialSecurity, input),
    calculateTSP: (input) => apiCall(CalculateTSPProjection, input),
    calculateTax: (input) => apiCall(CalculateTaxLiability, input),
    calculateCOLA: (input) => apiCall(CalculateCOLAAdjustment, input),
    calculateRetirementProjection: (input) => apiCall(CalculateRetirementProjection, input),
    saveScenarios: (input) => apiCall(SaveScenarios, input),
    loadScenarios: (input) => apiCall(LoadScenarios, input),
    getSavedScenarios: () => apiCall(GetSavedScenarios)
};