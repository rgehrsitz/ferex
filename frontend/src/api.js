// api.js - Direct API integrations for Svelte V5
import * as AppModule from '../wailsjs/go/main/App.js';

// Add error handling wrapper to each function
const wrapWithErrorHandling = (fn, name) => {
  if (!fn) {
    console.warn(`API function ${name} not found, using mock`);
    return async () => ({ success: false, error: `API function ${name} not available` });
  }
  
  return async (...args) => {
    try {
      console.log(`Calling ${name} with args:`, args);
      const result = await fn(...args);
      console.log(`${name} result:`, result);
      return result;
    } catch (err) {
      console.error(`Error in ${name}:`, err);
      return { success: false, error: err.toString() };
    }
  };
};

// Export the API with all available functions
export const api = {
  calculatePension: wrapWithErrorHandling(AppModule.CalculatePension, 'calculatePension'),
  calculateSocialSecurity: wrapWithErrorHandling(AppModule.CalculateSocialSecurity, 'calculateSocialSecurity'),
  calculateTSP: wrapWithErrorHandling(AppModule.CalculateTSPProjection, 'calculateTSP'),
  calculateTax: wrapWithErrorHandling(AppModule.CalculateTaxLiability, 'calculateTax'),
  calculateCOLA: wrapWithErrorHandling(AppModule.CalculateCOLAAdjustment, 'calculateCOLA'),
  calculateRetirementProjection: wrapWithErrorHandling(AppModule.CalculateRetirementProjection, 'calculateRetirementProjection'),
  saveScenarios: wrapWithErrorHandling(AppModule.SaveScenarios, 'saveScenarios'),
  loadScenarios: wrapWithErrorHandling(AppModule.LoadScenarios, 'loadScenarios'),
  getSavedScenarios: wrapWithErrorHandling(AppModule.GetSavedScenarios, 'getSavedScenarios')
};

// Log for debugging
console.log('API functions initialized:', Object.keys(api));