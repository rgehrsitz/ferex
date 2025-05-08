<script lang="ts">
  import { CalculateCOLAAdjustment } from '../../wailsjs/go/main/App.js';
  import { main } from '../../wailsjs/go/models.js';
  import type { COLAData } from '../types/scenario.js';
  import SectionHeader from './SectionHeader.svelte';
  
  import { onMount } from 'svelte';
  
  // Default data values
  const defaultData = {
    assumedInflationRate: 2.5,
    applyCOLAToPension: true,
    applyColaToSocialSecurity: true
  };
  
  const { 
    data: propData = defaultData, 
    scenarioName, 
    onUpdate = (data: COLAData) => {} 
  } = $props<{
    data?: COLAData;
    scenarioName: string;
    onUpdate?: (data: COLAData) => void;
  }>();
  
  // Using $state for local data, initialized with defaults or incoming data
  let localData = $state<COLAData>(
    propData ? { ...defaultData, ...propData } : { ...defaultData }
  );
  
  // Store a previous version of localData as a string for comparison
  let previousLocalDataString = $state("");
  
  // Initialize previousLocalDataString after localData is defined
  $effect(() => {
    if (previousLocalDataString === "") {
      previousLocalDataString = JSON.stringify(localData);
    }
  });
  
  // Track if we're updating from props to prevent loops
  let updatingFromProps = $state(false);

  // Manual function to update local data when props change
  function updateLocalDataFromProps() {
    if (!propData) return;
    
    try {
      const dataString = JSON.stringify(propData);
      // Only update if different
      if (dataString !== previousLocalDataString) {
        console.log('COLASection: Props changed, updating local data');
        updatingFromProps = true;
        localData = { ...defaultData, ...propData };
        previousLocalDataString = dataString;
      }
    } finally {
      updatingFromProps = false;
    }
  }

  // Update local data when props change - but only on initial props or when they actually change
  let previousPropString = $state('');
  
  $effect(() => {
    if (propData) {
      const propString = JSON.stringify(propData);
      if (propString !== previousPropString) {
        previousPropString = propString;
        updateLocalDataFromProps();
      }
    }
  });

  // Manual function to update the props from localData
  function updateParentData() {
    if (updatingFromProps) return; // Skip if we're currently updating from props
    
    const currentString = JSON.stringify(localData);
    
    // Only update if the data has actually changed
    if (currentString !== previousLocalDataString) {
      console.log('COLASection: localData changed, updating parent');
      previousLocalDataString = currentString;
      
      // Create a deep copy to avoid reference issues
      const dataForParent = JSON.parse(currentString);
      
      // Notify parent via callback
      onUpdate(dataForParent);
    }
  }
  
  // For backward compatibility in the template
  const data = $derived(localData);

  let loading = $state(false);
  let error = $state('');
  let calculationResult = $state({
    baseAmount: 0,
    finalAmount: 0,
    totalGrowthRate: 0,
    effectiveAnnualRate: 0,
    yearlyAdjustments: [],
    notes: ''
  });

  // Projection parameters
  let projectionYears = $state(20);
  let baseAmount = $state(30000); // Sample pension amount
  let monthsInFirstYear = $state(12); // Full first year

  // COLA system to display
  let selectedSystem = $state<'FERS' | 'CSRS' | 'Social Security'>('FERS');

  // Needed for calculations
  const currentYear = new Date().getFullYear();

  const pensionOptions = [
    { value: 'FERS', label: 'FERS COLA', description: 'Follows inflation with caps: 2% if inflation is 2-3%, inflation-1% for inflation >3%' },
    { value: 'CSRS', label: 'CSRS COLA', description: 'Matches full inflation rate with no caps' },
    { value: 'Social Security', label: 'Social Security COLA', description: 'Based on CPI-W, usually similar to general inflation' }
  ];
  
  async function calculateCOLA() {
    try {
      loading = true;
      error = '';
      
      // Create input object for backend
      const colaInput = new main.COLAInput();
      colaInput.assumedInflationRate = localData.assumedInflationRate / 100; // Convert to decimal
      colaInput.applyCOLAToPension = localData.applyCOLAToPension;
      colaInput.applyColaToSocialSecurity = localData.applyColaToSocialSecurity;
      colaInput.baseAmount = baseAmount;
      colaInput.retirementSystem = selectedSystem;
      colaInput.retirementAge = 60; // Default for calculation
      colaInput.isSpecialProvision = false;
      colaInput.startYear = currentYear;
      colaInput.projectionYears = projectionYears;
      colaInput.monthsInFirstYear = monthsInFirstYear;
      
      // Call backend API
      calculationResult = await CalculateCOLAAdjustment(colaInput);
      
      // Calculate example COLAs for display
      calculateEstimatedFirstYearCOLA(localData.assumedInflationRate);
      
      return calculationResult;
    } catch (err) {
      console.error("Error calculating COLA adjustment:", err);
      error = 'Failed to calculate COLA adjustment. Please try again.';
    } finally {
      loading = false;
    }
  }
  
  function calculateEstimatedFirstYearCOLA(inflationRate: number) {
    // Calculate FERS COLA based on inflation (as decimal)
    const inflationRateDecimal = inflationRate / 100;
    let fersCOLA = 0;
    
    if (inflationRateDecimal <= 0.02) {
      fersCOLA = inflationRateDecimal;
    } else if (inflationRateDecimal <= 0.03) {
      fersCOLA = 0.02;
    } else {
      fersCOLA = inflationRateDecimal - 0.01;
    }
    
    // Return FERS COLA, CSRS COLA (which is the full inflation rate),
    // and approximate SS COLA (using inflation as proxy)
    return {
      fers: localData.applyCOLAToPension ? fersCOLA * 100 : inflationRate,
      csrs: inflationRate,
      socialSecurity: inflationRate
    };
  }
  
  // Trigger calculation when relevant inputs change
  $effect(() => {
    if (localData.assumedInflationRate !== undefined || 
        localData.applyCOLAToPension !== undefined || 
        localData.applyColaToSocialSecurity !== undefined ||
        selectedSystem) {
      calculateCOLA();
    }
  });

  const estCOLA = $derived(calculateEstimatedFirstYearCOLA(localData.assumedInflationRate));
</script>

<div>
  <SectionHeader sectionName="COLA Calculator" {scenarioName} />
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="assumedInflationRate">
          Projected Inflation Rate
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="assumedInflationRate"
            type="number"
            min="0"
            max="10"
            step="0.1"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={localData.assumedInflationRate}
            onchange={() => {
              if (typeof localData.assumedInflationRate === 'string') {
                localData.assumedInflationRate = parseFloat(localData.assumedInflationRate);
              }
              updateParentData();
            }}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>

      <div class="flex items-center">
        <input
          id="applyCOLAToPension"
          type="checkbox"
          class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          bind:checked={localData.applyCOLAToPension}
          onchange={updateParentData}
        />
        <label for="applyCOLAToPension" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
          Apply COLA to Pension
        </label>
      </div>

      <div class="flex items-center">
        <input
          id="applyColaToSocialSecurity"
          type="checkbox"
          class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          bind:checked={localData.applyColaToSocialSecurity}
          onchange={updateParentData}
        />
        <label for="applyColaToSocialSecurity" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
          Apply COLA to Social Security
        </label>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="projectionYears">
          Projection Years
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="projectionYears"
            type="number"
            min="5"
            max="40"
            step="1"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={projectionYears}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">years</span>
          </div>
        </div>
      </div>
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="selectedSystem">
          Retirement System
        </label>
        <select 
          id="selectedSystem"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={selectedSystem}
        >
          {#each pensionOptions as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </select>
      </div>

      <div class="bg-blue-50 dark:bg-blue-900 p-3 rounded-md mt-2">
        <h4 class="text-sm font-medium text-blue-800 dark:text-blue-300">COLA Rules</h4>
        <ul class="mt-2 text-xs text-blue-700 dark:text-blue-400 space-y-1 list-disc pl-4">
          <li>FERS: If inflation is 2% or less, COLA = Inflation</li>
          <li>FERS: If inflation is >2% but ≤3%, COLA = 2%</li>
          <li>FERS: If inflation is >3%, COLA = Inflation - 1%</li>
          <li>CSRS: COLA = Full Inflation (no caps)</li>
          <li>SS: COLA = Based on CPI-W (typically close to inflation)</li>
          <li>FERS retirees don't receive COLA until age 62 (except special categories)</li>
        </ul>
      </div>
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">COLA Projections</h3>
    
    {#if error}
      <div class="bg-red-100 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-700 dark:text-red-300 p-4 rounded-lg mb-4">
        {error}
      </div>
    {/if}
    
    {#if loading}
      <div class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 dark:border-primary-400"></div>
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">FERS COLA</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{estCOLA.fers.toFixed(1)}%</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">CSRS COLA</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{estCOLA.csrs.toFixed(1)}%</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Social Security COLA</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{estCOLA.socialSecurity.toFixed(1)}%</div>
        </div>
      </div>
      
      {#if calculationResult && calculationResult.yearlyAdjustments && calculationResult.yearlyAdjustments.length > 0}
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow mb-6">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">{selectedSystem} COLA Projection Summary</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Starting Amount</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                ${calculationResult.baseAmount.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Final Amount (after {projectionYears} years)</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                ${calculationResult.finalAmount.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Total Growth</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                {(calculationResult.totalGrowthRate * 100).toFixed(1)}%
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Effective Annual Rate</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                {(calculationResult.effectiveAnnualRate * 100).toFixed(2)}%
              </div>
            </div>
          </div>
        </div>
        
        <!-- Year-by-Year COLA Projection Table (first few years) -->
        <div class="overflow-x-auto">
          <table class="min-w-full bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
            <thead>
              <tr class="bg-gray-50 dark:bg-gray-700">
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Year</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Inflation</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">COLA %</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Starting Amount</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Adjusted Amount</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Total Growth</th>
              </tr>
            </thead>
            <tbody>
              {#each calculationResult.yearlyAdjustments.slice(0, 10) as yearData}
                <tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200">{yearData.year}</td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    {(yearData.inflationRate * 100).toFixed(1)}%
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    {(yearData.colaPercentage * 100).toFixed(1)}%
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.startingAmount.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.adjustedAmount.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    {(yearData.cumulativeGrowth * 100).toFixed(1)}%
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
      
      {#if calculationResult && calculationResult.notes}
        <div class="mt-4 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Notes</h4>
          <div class="text-sm text-gray-600 dark:text-gray-400">
            {calculationResult.notes}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>