import type { ScenarioInput } from '../../types';
import { CalculateTSP } from 'wailsjs/go/main/App';

/**
 * Service for handling TSP-related API calls
 */
export class TSPService {
  /**
   * Calculates TSP projections using the provided scenario inputs
   * The backend CalculateTSP function expects a scenario.ScenarioVariant
   * and handles the mapping to TSPCalculationInput internally
   *
   * @param inputs The scenario inputs containing TSP data
   * @returns Promise with calculation results
   */
  static async calculateTSP (inputs: ScenarioInput) {
    try {
      // console.log('TSPService: Sending TSP calculation request with ScenarioInput...');

      // The backend CalculateTSP function expects a scenario.ScenarioVariant
      // and will handle the mapping to TSPCalculationInput internally
      const result = await CalculateTSP(inputs as any);

      // console.log('TSPService: TSP calculation completed successfully');
      return result;
    } catch (error) {
      // console.error('TSPService: Error calculating TSP:', error);
      throw error;
    }
  }
}
