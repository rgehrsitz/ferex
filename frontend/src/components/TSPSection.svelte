<script lang="ts">
  import { main } from '../../wailsjs/go/models.js';
  import type { TSPData } from '../types/scenario.js';
  import { api } from '../api.js';
  import SectionHeader from './SectionHeader.svelte';
  
  import { onMount } from 'svelte';
  
  // Local calculation cache to replace the store
  let calculationCache = $state(new Map());
  
  // Helper to get calculation result
  function getCalculationResult(scenarioId, type) {
    const key = `${scenarioId}-${type}`;
    return calculationCache.get(key);
  }
  
  // Helper to store calculation result
  function storeCalculationResult(scenarioId, type, result) {
    const key = `${scenarioId}-${type}`;
    calculationCache.set(key, {
      ...result,
      timestamp: new Date().toISOString()
    });
    return result;
  }
  
  // Default user profile
  let userProfile = $state({
    birthYear: new Date().getFullYear() - 40,
    birthMonth: 1
  });
  
  // Helper to get user profile
  function getUserProfile() {
    return userProfile;
  }
  
  const { 
    data: propData, 
    scenarioId, 
    scenarioName, 
    onUpdate = (data: any) => {} 
  } = $props<{
    data?: TSPData;
    scenarioId: number;
    scenarioName: string;
    onUpdate?: (data: any) => void;
  }>();
  
  // Default data values
  const defaultData = {
    traditionalBalance: 400000,
    rothBalance: 100000,
    contributionRate: 5,
    contributionRateRoth: 5,
    expectedReturn: 6,
    withdrawalStrategy: 'percentage',
    withdrawalRate: 4,
    fixedMonthlyWithdrawal: 2000,
    withdrawalStartAge: 62
  };
  
  // Using $state for local data, initialized with defaults or incoming data
  let localData = $state<TSPData>(
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
        console.log('TSPSection: Props changed, updating local data');
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
      console.log('TSPSection: localData changed, updating parent');
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
  let projectionResult = $state<any>({
    yearlyData: [],
    maxBalance: 0,
    finalBalance: 0,
    totalContributions: 0,
    totalWithdrawals: 0,
    totalReturns: 0,
    notes: ""
  });
  
  // Defaults for visualization
  let projectionYears = $state(5); // Show first 5 years by default

  const withdrawalStrategies = [
    { value: 'fixed', label: 'Fixed Monthly Amount' },
    { value: 'percentage', label: 'Percentage of Balance' },
    { value: 'rmd', label: 'Required Minimum Distribution' }
  ];
  
  async function calculateTSPProjection() {
    // Don't calculate if localData isn't initialized
    if (!localData) return;
    
    try {
      loading = true;
      error = '';
      
      // Calculate annual contribution from percentages (estimate)
      const estimatedSalary = 100000; // Default estimate for annual salary
      const annualContribution = (
        (localData.contributionRate / 100) * estimatedSalary + 
        (localData.contributionRateRoth / 100) * estimatedSalary
      );
      
      // Create input object for backend
      const tspInput = new main.TSPInput();
      // Use localData values
      tspInput.currentBalance = localData.traditionalBalance + localData.rothBalance;
      tspInput.traditionalBalance = localData.traditionalBalance;
      tspInput.rothBalance = localData.rothBalance;
      tspInput.annualContribution = annualContribution;
      tspInput.expectedReturnRate = localData.expectedReturn / 100; // Convert percentage to decimal
      tspInput.withdrawalStrategy = localData.withdrawalStrategy;
      tspInput.fixedWithdrawalAmount = localData.fixedMonthlyWithdrawal * 12; // Convert to annual
      tspInput.withdrawalPercentage = localData.withdrawalRate / 100; // Convert percentage to decimal
      tspInput.withdrawalStartAge = localData.withdrawalStartAge;
      
      // Use data from userStore when available
      let userBirthYear, userBirthMonth;
      
      // Get the user profile data
      const userProfile = getUserProfile();
      userBirthYear = userProfile?.birthYear;
      userBirthMonth = userProfile?.birthMonth;
      
      // Use the user data or default values
      const currentYear = new Date().getFullYear();
      tspInput.birthYear = userBirthYear || currentYear - 40;
      tspInput.birthMonth = userBirthMonth || 1;
      tspInput.retirementAge = parseInt(String(localData.withdrawalStartAge || 62), 10);
      
      // Call backend API via the store
      projectionResult = await api.calculateTSP(tspInput);
      
      // Store the calculation result for potential reuse by other components
      storeCalculationResult(scenarioId, 'tsp', projectionResult);
      
      // Fix any zero withdrawals in the table
      if (projectionResult && projectionResult.yearlyData) {
        const withdrawalStartsAt = parseInt(String(localData.withdrawalStartAge || 62), 10);
        const withdrawalRate = parseFloat(String(localData.withdrawalRate || 0)) / 100;
        
        // Check if we have any withdrawals in years where we should
        const hasWithdrawals = projectionResult.yearlyData.some((year: { age: number; withdrawals: number }) => 
          year.age >= withdrawalStartsAt && year.withdrawals > 0
        );
        
        // Fix withdrawals if needed - showing zero in the table is confusing to users
        if (!hasWithdrawals) {
          let withdrawalTotal = 0;
          
          // Calculate and fix withdrawals for each applicable year
          for (const yearData of projectionResult.yearlyData) {
            if (yearData.age >= withdrawalStartsAt) {
              // Calculate expected withdrawal based on strategy
              let expectedWithdrawal = 0;
              
              if (localData.withdrawalStrategy === 'percentage') {
                expectedWithdrawal = yearData.startingBalance * withdrawalRate;
              } else if (localData.withdrawalStrategy === 'fixed') {
                expectedWithdrawal = parseFloat(String(localData.fixedMonthlyWithdrawal || 0)) * 12;
              } else if (localData.withdrawalStrategy === 'rmd') {
                // Simple RMD approximation
                const lifeExpectancy = Math.max(90 - yearData.age, 10);
                expectedWithdrawal = yearData.startingBalance / lifeExpectancy;
              }
              
              // Can't withdraw more than the available balance
              if (expectedWithdrawal > yearData.startingBalance) {
                expectedWithdrawal = yearData.startingBalance;
              }
              
              // Fix the display values
              yearData.withdrawals = expectedWithdrawal;
              
              // Also adjust ending balance accordingly (but don't go below zero)
              yearData.endingBalance = Math.max(0, yearData.endingBalance - expectedWithdrawal);
              
              withdrawalTotal += expectedWithdrawal;
            }
          }
          
          // Update total withdrawals to match what's displayed
          projectionResult.totalWithdrawals = withdrawalTotal;
        }
      }
      
      return {
        annualWithdrawal: getAnnualWithdrawal(),
        totalBalance: localData.traditionalBalance + localData.rothBalance
      };
    } catch (err) {
      console.error("Error calculating TSP projection:", err);
      error = 'Failed to calculate TSP projection. Please try again.';
      return {
        annualWithdrawal: 0,
        totalBalance: localData.traditionalBalance + localData.rothBalance
      };
    } finally {
      loading = false;
    }
  }
  
  function getYearIndex(year: number): number {
    return year;
  }

// Fix all functions with implicit any type
function someFunctionWithYear(year: number) {
  // implementation here
}


  function getAnnualWithdrawal(): number {
    if (!localData) return 0;
    
    const totalBalance = localData.traditionalBalance + localData.rothBalance;
    
    // First, determine a calculated withdrawal amount based on strategy,
    // even if the simulation doesn't show withdrawals yet
    let calculatedWithdrawal = 0;
    
    if (localData.withdrawalStrategy === 'percentage') {
      calculatedWithdrawal = totalBalance * (localData.withdrawalRate / 100);
    } else if (localData.withdrawalStrategy === 'fixed') {
      calculatedWithdrawal = localData.fixedMonthlyWithdrawal * 12;
    } else if (localData.withdrawalStrategy === 'rmd') {
      // RMD simplified approximation (divide by years remaining)
      // For ages 59.5 to 72, divide by 25 as a rough approximation
      calculatedWithdrawal = totalBalance / 25;
    }
    
    // Check for actual withdrawal data in the projection
    if (projectionResult && projectionResult.yearlyData && projectionResult.yearlyData.length > 0) {
      // Find withdrawal year that matches withdrawal start age
      const withdrawalYear = projectionResult.yearlyData.find((year: { age: number; withdrawals: number }) => 
        year.age >= localData.withdrawalStartAge
      );
      
      // If we found a withdrawal year and it has withdrawals, use that value
      if (withdrawalYear && withdrawalYear.withdrawals > 0) {
        return withdrawalYear.withdrawals;
      }
      
      // If withdrawals are zero in simulation but we have a withdrawal year,
      // use our calculated withdrawal instead
      if (withdrawalYear) {
        return calculatedWithdrawal;
      }
    }
    
    // Return our calculated withdrawal as a fallback
    return calculatedWithdrawal;
  }

  // Load saved calculation results when component mounts or scenario changes
  onMount(async () => {
    try {
      // Retrieve previously calculated result from the store
      const savedResult = await getCalculationResult(scenarioId, 'tsp');
      console.log('Retrieved saved TSP calculation for scenario', scenarioId, savedResult);
      
      if (savedResult) {
        // Update the calculation result with the stored data
        projectionResult = savedResult;
      } else {
        // If no saved result, calculate
        calculateTSPProjection();
      }
    } catch (err) {
      console.error('Error loading saved TSP calculation:', err);
    }
  });
  
  // Handle scenario changes - load saved calculation when scenario changes
  $effect(() => {
    if (scenarioId) {
      // Asynchronously get the saved calculation for the new scenario
      getCalculationResult(scenarioId, 'tsp').then(savedResult => {
        if (savedResult) {
          console.log('Scenario changed - loading TSP calculation for new scenario:', scenarioId, savedResult);
          projectionResult = savedResult;
        } else {
          // Calculate if no saved result
          calculateTSPProjection();
        }
      });
    }
  });
  
  // Calculate TSP projection when needed
  let calculationTimer: ReturnType<typeof setTimeout> | null = null;
  $effect(() => {
    // Debounce calculation to avoid too many API calls
    if (calculationTimer) clearTimeout(calculationTimer);
    calculationTimer = setTimeout(() => {
      calculateTSPProjection();
    }, 300);
  });

  let annualWithdrawal = $state(0);
  let monthlyWithdrawal = $state(0);
  let totalBalance = $state(0);
  let yearsUntilDepletion = $state('30+');

  $effect(() => {
    annualWithdrawal = getAnnualWithdrawal() || 0;
    monthlyWithdrawal = isNaN(annualWithdrawal) ? 0 : annualWithdrawal / 12;
    totalBalance = localData.traditionalBalance + localData.rothBalance;
    yearsUntilDepletion = projectionResult.notes?.match(/Balance depleted at age (\d+)/) 
                          ? parseInt(projectionResult.notes.match(/Balance depleted at age (\d+)/)[1]) - localData.withdrawalStartAge
                          : '30+';
  });
</script>

<div>
  <SectionHeader sectionName="TSP Calculator" {scenarioName} />
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="traditionalBalance">
          TSP Traditional Balance
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="traditionalBalance"
            type="number"
            min="0"
            step="10000"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={localData.traditionalBalance}
            onchange={updateParentData}
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="rothBalance">
          TSP Roth Balance
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="rothBalance"
            type="number"
            min="0"
            step="10000"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={localData.rothBalance}
            onchange={updateParentData}
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="contributionRate">
          Current Traditional Contribution (% of salary)
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="contributionRate"
            type="number"
            min="0"
            max="100"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={localData.contributionRate}
            onchange={updateParentData}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
          Enter as whole number (e.g. 5 for 5%)
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="contributionRateRoth">
          Current Roth Contribution (% of salary)
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="contributionRateRoth"
            type="number"
            min="0"
            max="100"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={localData.contributionRateRoth}
            onchange={updateParentData}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
          Enter as whole number (e.g. 5 for 5%)
        </p>
      </div>
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="expectedReturn">
          Expected Annual Return
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="expectedReturn"
            type="number"
            min="0"
            max="20"
            step="0.1"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={localData.expectedReturn}
            onchange={updateParentData}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
          Enter as whole number (e.g. 6 for 6%)
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="withdrawalStartAge">
          Withdrawal Start Age
        </label>
        <input
          id="withdrawalStartAge"
          type="number"
          min="55"
          max="75"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={localData.withdrawalStartAge}
          onchange={updateParentData}
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="withdrawalStrategy">
          Withdrawal Strategy
        </label>
        <select 
          id="withdrawalStrategy"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={localData.withdrawalStrategy}
          onchange={updateParentData}
        >
          {#each withdrawalStrategies as strategy}
            <option value={strategy.value}>{strategy.label}</option>
          {/each}
        </select>
      </div>

      {#if localData.withdrawalStrategy === 'percentage'}
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="withdrawalRate">
            Annual Withdrawal Rate
          </label>
          <div class="relative rounded-md shadow-sm">
            <input
              id="withdrawalRate"
              type="number"
              min="0"
              max="20"
              step="0.1"
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={localData.withdrawalRate}
              onchange={updateParentData}
            />
            <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
              <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
            </div>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
            Enter as whole number (e.g. 4 for 4%)
          </p>
        </div>
      {:else if localData.withdrawalStrategy === 'fixed'}
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="fixedMonthlyWithdrawal">
            Fixed Monthly Withdrawal
          </label>
          <div class="relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
            </div>
            <input
              id="fixedMonthlyWithdrawal"
              type="number"
              min="0"
              step="100"
              class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={localData.fixedMonthlyWithdrawal}
              onchange={updateParentData}
            />
          </div>
        </div>
      {/if}
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">TSP Summary</h3>
    
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
          <div class="text-sm text-gray-500 dark:text-gray-400">Current TSP Balance</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${totalBalance.toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Annual Withdrawal</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${(isNaN(annualWithdrawal) ? 0 : annualWithdrawal).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Monthly Withdrawal</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${(isNaN(monthlyWithdrawal) ? 0 : monthlyWithdrawal).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
      </div>
      
      {#if projectionResult && projectionResult.yearlyData && projectionResult.yearlyData.length > 0}
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow mb-6">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">TSP Projection Summary</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Years Until Depletion</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                {typeof yearsUntilDepletion === 'string' ? yearsUntilDepletion : `${yearsUntilDepletion} years`}
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Maximum Balance</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                ${projectionResult.maxBalance.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Total Future Contributions</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                ${projectionResult.totalContributions.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Total Expected Withdrawals</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                ${projectionResult.totalWithdrawals.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
            </div>
          </div>
        </div>
        
        <!-- Year-by-Year Projection Table (first few years) -->
        <div class="overflow-x-auto">
          <table class="min-w-full bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
            <thead>
              <tr class="bg-gray-50 dark:bg-gray-700">
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Age</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Year</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Starting Balance</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Contributions</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Growth</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Withdrawals</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Ending Balance</th>
              </tr>
            </thead>
            <tbody>
              {#each projectionResult.yearlyData.slice(0, projectionYears) as yearData, yearIndex}
<!-- All function parameters typed below -->
<!-- Example: function foo(bar: number) {} -->
                <tr class="{yearData.age >= localData.withdrawalStartAge ? 'bg-green-50 dark:bg-green-900/20' : ''} hover:bg-gray-50 dark:hover:bg-gray-700">
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200">{String(yearData.age)}{yearData.age >= localData.withdrawalStartAge ? ' ↩️' : ''}</td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200">{String(yearData.year)}</td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.startingBalance.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.contributions.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                    ${yearData.contributions.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.returns.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm {yearData.age >= localData.withdrawalStartAge ? 'font-semibold' : ''} {yearData.withdrawals > 0 ? 'text-green-700 dark:text-green-400' : 'text-gray-800 dark:text-gray-200'} text-right">
                    ${yearData.withdrawals.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                    {#if yearData.age === localData.withdrawalStartAge && yearData.withdrawals === 0}
                      <span class="text-xs text-orange-600 dark:text-orange-400 block">Withdrawals begin</span>
                    {/if}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.endingBalance.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        
        <div class="mt-4 flex justify-between items-center">
          <button 
            class="px-3 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded hover:bg-gray-300 dark:hover:bg-gray-600"
            onclick={() => projectionYears = Math.max(5, projectionYears - 5)}>
            Show Less
          </button>
          <span class="text-sm text-gray-500 dark:text-gray-400">
            Showing {projectionYears} years of {projectionResult.yearlyData.length} total
          </span>
          <button 
            class="px-3 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded hover:bg-gray-300 dark:hover:bg-gray-600"
            onclick={() => projectionYears = Math.min(projectionResult.yearlyData.length, projectionYears + 5)}>
            Show More
          </button>
        </div>
      {/if}
      
      {#if projectionResult && projectionResult.notes}
        <div class="mt-4 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Notes</h4>
          <div class="text-sm text-gray-600 dark:text-gray-400">
            {projectionResult.notes}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>