<script lang="ts">
  // Svelte V5 idioms: props only, no stores
  import { onMount } from 'svelte';
  // Svelte 5 runes mode: use $props for props
  const { selectedScenario, compareScenario } = $props<{
    selectedScenario: import('../types/scenario.js').Scenario | null;
    compareScenario: import('../types/scenario.js').Scenario | null;
  }>();
  import { CalculateRetirementProjection } from '../../wailsjs/go/main/App.js';
  import { main } from '../../wailsjs/go/models.js';
  import type { Scenario } from '../types/scenario.js';
  import MonthlyIncomeChart from './charts/SimpleMontlyIncomeChart.svelte';
  import CumulativeIncomeChart from './charts/SimpleCumulativeIncomeChart.svelte';
  import MonthlyIncomeDeltaChart from './charts/SimpleMonthlyDeltaChart.svelte';
  import IncomeSourceChart from './charts/SimpleIncomeSourceChart.svelte';

  let loading = $state(false);
  let error = $state('');
  let selectedProjection = $state<any>(null);
  let compareProjection = $state<any>(null);

  // Svelte 5 idiom: use $derived for derived values
  const selectedNetTotal = $derived(selectedProjection?.netIncomeTotal ?? 0);
  const compareNetTotal = $derived(compareProjection?.netIncomeTotal ?? 0);
  const netIncomeDifference = $derived(selectedNetTotal - compareNetTotal);
  const cumulativeDifference = $derived((selectedProjection?.cumulativeIncomeTotal ?? 0) - (compareProjection?.cumulativeIncomeTotal ?? 0));

  // Chart visibility controls
  let activeCharts = {
    monthlyIncome: true,
    cumulativeIncome: true,
    monthlyDelta: true,
    incomeSource: true
  };
  
  async function runProjections() {
    try {
      loading = true;
      error = '';
      
      // Create input for first scenario
      const scenarioInput1 = createProjectionInput(selectedScenario);
      
      // Create input for second scenario
      if (!selectedScenario || !compareScenario) {
  loading = false;
  error = 'Both scenarios must be selected for comparison.';
  return;
}
const scenarioInput2 = createProjectionInput(compareScenario);
      
      // Run projections in parallel
      const [projection1, projection2] = await Promise.all([
        CalculateRetirementProjection(scenarioInput1),
        CalculateRetirementProjection(scenarioInput2)
      ]);
      
      selectedProjection = projection1;
      compareProjection = projection2;
      
      // Calculate summary statistics
      // No assignments to $derived variables! Let Svelte reactivity handle updates.
      // If you need a local cumulative sum for another purpose, use a local variable here.
      // All display and logic using selectedNetTotal, compareNetTotal, netIncomeDifference, or cumulativeDifference will update automatically via $derived.
      
    } catch (err) {
      console.error("Error running projections:", err);
      error = 'Failed to calculate retirement projections. Please try again.';
    } finally {
      loading = false;
    }
  }
  
  // Helper function to create projection input
  function createProjectionInput(scenario: Scenario) {
    // Convert frontend scenario to backend input format
    const input = new main.RetirementScenarioInput();
    
    // Pension data
    const pension = new main.PensionInput();
    pension.system = scenario.data.pension.system ?? '';
    pension.high3Salary = scenario.data.pension.high3Salary ?? 0;
    pension.yearsOfService = scenario.data.pension.yearsOfService ?? 0;
    pension.ageAtRetirement = scenario.data.pension.ageAtRetirement ?? 0;
    pension.unusedSickLeaveMonths = scenario.data.pension.unusedSickLeaveMonths ?? 0;
    pension.survivorBenefitOption = scenario.data.pension.survivorBenefitOption ?? '';
    pension.isPartTime = scenario.data.pension.isPartTime ?? false;
    pension.partTimeProrationFactor = scenario.data.pension.partTimeProrationFactor ?? 1;
    input.pension = pension;
    
    // Social Security data
    if (scenario.data.socialSecurity) {
      const ss = new main.SocialSecurityInput();
      ss.startAge = scenario.data.socialSecurity.startAge || 62;
      ss.estimatedMonthlyBenefit = scenario.data.socialSecurity.estimatedMonthlyBenefit || 0;
      ss.isEligible = scenario.data.socialSecurity.isEligible !== undefined ? scenario.data.socialSecurity.isEligible : true;
      ss.birthYear = scenario.data.socialSecurity.birthYear || 1960;
      ss.birthMonth = scenario.data.socialSecurity.birthMonth || 1;
      
      // Set default values for estimation
      ss.estimatedAnnualSalary = 60000; // Default annual salary
      ss.yearsWorked = 35; // Default years worked for estimation
      
      // Add SSA estimates if available
      if ((scenario.data.socialSecurity.ssaEstimateAt62 ?? 0) > 0) {
        ss.userProvidedEstimate62 = scenario.data.socialSecurity.ssaEstimateAt62 ?? 0;
      }
      if ((scenario.data.socialSecurity.ssaEstimateAtFRA ?? 0) > 0) {
        ss.userProvidedEstimateFRA = scenario.data.socialSecurity.ssaEstimateAtFRA ?? 0;
      }
      if ((scenario.data.socialSecurity.ssaEstimateAt70 ?? 0) > 0) {
        ss.userProvidedEstimate70 = scenario.data.socialSecurity.ssaEstimateAt70 ?? 0;
      }
      
      input.socialSecurity = ss;
    }
    
    // TSP data
    if (scenario.data.tsp) {
      const tsp = new main.TSPInput();
      
      // Calculate total balance
      const totalBalance = (scenario.data.tsp.traditionalBalance || 0) + (scenario.data.tsp.rothBalance || 0);
      tsp.currentBalance = totalBalance;
      
      // Set individual balances
      tsp.traditionalBalance = scenario.data.tsp.traditionalBalance || 0;
      tsp.rothBalance = scenario.data.tsp.rothBalance || 0;
      
      // Calculate annual contribution from percentage (estimate)
      const estimatedSalary = 100000; // Default estimate
      const annualContribution = (
        ((scenario.data.tsp.contributionRate || 0) / 100) * estimatedSalary + 
        ((scenario.data.tsp.contributionRateRoth || 0) / 100) * estimatedSalary
      );
      tsp.annualContribution = annualContribution;
      
      // Expected returns and withdrawals
      tsp.expectedReturnRate = (scenario.data.tsp.expectedReturn || 0) / 100;
      tsp.withdrawalStrategy = scenario.data.tsp.withdrawalStrategy || 'fixed';
      
      // Fixed amount is stored monthly but API wants annual
      if (scenario.data.tsp.withdrawalStrategy === 'fixed') {
        tsp.fixedWithdrawalAmount = (scenario.data.tsp.fixedMonthlyWithdrawal || 0) * 12;
      } else if (scenario.data.tsp.withdrawalStrategy === 'percentage') {
        tsp.withdrawalPercentage = (scenario.data.tsp.withdrawalRate || 0) / 100;
      } else {
        // RMD has no parameter
        tsp.withdrawalPercentage = 0;
      }
      
      // Withdrawal start age
      tsp.withdrawalStartAge = scenario.data.tsp.withdrawalStartAge || 62;
      
      input.tsp = tsp;
    }
    
    // Tax data
    if (scenario.data.tax) {
      const tax = new main.TaxInput();
      tax.filingStatus = scenario.data.tax.filingStatus;
      tax.stateOfResidence = scenario.data.tax.stateOfResidence;

      tax.stateIncomeTaxRate = scenario.data.tax.stateIncomeTaxRate;
      input.tax = tax;
    }
    
    // COLA data
    if (scenario.data.cola) {
      const cola = new main.COLAInput();
      cola.assumedInflationRate = scenario.data.cola.assumedInflationRate;
      cola.applyCOLAToPension = scenario.data.cola.applyCOLAToPension;
      cola.applyColaToSocialSecurity = scenario.data.cola.applyColaToSocialSecurity;
      input.cola = cola;
    }
    
    // Other Income data
    const otherIncome = new main.OtherIncomeInput();
    
    if (scenario.data.otherIncome && scenario.data.otherIncome.sources && scenario.data.otherIncome.sources.length > 0) {
      otherIncome.sources = scenario.data.otherIncome.sources.map(src => {
        const source = new main.OtherIncomeSource();
        source.id = src.id;
        source.name = src.name;
        source.amount = src.amount;
        source.frequency = src.frequency;
        source.startAge = src.startAge;
        source.endAge = src.endAge;
        source.applyCola = src.applyCola;
        return source;
      });
    } else {
      otherIncome.sources = []; // Ensure it's initialized with an empty array
    }
    
    input.otherIncome = otherIncome;
    
    // Calculate current age based on birth year
    const currentYear = new Date().getFullYear();
    const currentAge = scenario.data.socialSecurity?.birthYear 
      ? (currentYear - scenario.data.socialSecurity.birthYear) 
      : 45; // Default to age 45 if birth year not available
      
    // Projection age range - start at current age, not retirement age
    input.projectionStartAge = currentAge;
    input.projectionEndAge = 95; // Default to age 95
    
    console.log(`Setting projection range: ${currentAge} to 95 years old`);
    
    return input;
  }
  
  // Run projections when component mounts
  onMount(() => {
    runProjections(); // Svelte 5 idiom: use onMount for side effects
  });
</script>

<div class="p-4 bg-white dark:bg-gray-800 rounded-lg shadow">
  {#if error}
    <div class="bg-red-100 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-700 dark:text-red-300 p-4 rounded-lg mb-4">
      {error}
    </div>
  {/if}
  
  {#if loading}
    <div class="flex justify-center items-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 dark:border-primary-400"></div>
      <span class="ml-2 text-gray-600 dark:text-gray-300">Calculating retirement projections...</span>
    </div>
  {:else if selectedProjection && compareProjection}
    <!-- Summary Section -->
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">Summary Comparison</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
          <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Lifetime Net Income Difference</h3>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200 mt-2">
            {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(netIncomeDifference)}
          </div>
          <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {compareScenario?.name ?? 'N/A'} vs {selectedScenario?.name ?? 'N/A'}
          </div>
        </div>
        
        <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
          <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">{selectedScenario?.name ?? 'Scenario'} Total</h3>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200 mt-2">
            {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(selectedNetTotal)}
          </div>
        </div>
        
        <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
          <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">{compareScenario.name} Total</h3>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200 mt-2">
            {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(compareNetTotal)}
          </div>
        </div>
      </div>
    </div>
    
    <!-- Yearly Comparison Table -->
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">Year by Year Comparison</h2>
      
      <div class="overflow-x-auto">
        <table class="min-w-full bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-700">
              <th class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Age</th>
              <th class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Year</th>
              <th class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{selectedScenario.name} Income</th>
              <th class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{compareScenario.name} Income</th>
              <th class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Difference</th>
            </tr>
          </thead>
          <tbody>
            {#each selectedProjection.yearlyData as yearData, i}
              {@const compareYear = compareProjection.yearlyData[i] || { netIncome: 0 }}
              {@const difference = compareYear.netIncome - yearData.netIncome}
              <tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
                <td class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-sm font-medium text-gray-800 dark:text-gray-200">{yearData.age}</td>
                <td class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-600 dark:text-gray-400">{yearData.year}</td>
                <td class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                  {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(yearData.netIncome)}
                </td>
                <td class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                  {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(compareYear.netIncome)}
                </td>
                <td class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-sm text-right" 
                  class:text-green-600={difference > 0} 
                  class:text-red-600={difference < 0}
                  class:dark:text-green-400={difference > 0} 
                  class:dark:text-red-400={difference < 0}>
                  {new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0, signDisplay: 'always' }).format(difference)}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
    
    <!-- Chart Visualizations -->
    <div class="mb-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200">Visualizations</h2>
        
        <div class="flex flex-wrap gap-2">
          <button 
            class="px-3 py-1 text-sm rounded border"
            class:bg-primary-600={activeCharts.monthlyIncome} 
            class:text-white={activeCharts.monthlyIncome}
            class:border-primary-600={activeCharts.monthlyIncome}
            class:bg-white={!activeCharts.monthlyIncome} 
            class:dark:bg-gray-800={!activeCharts.monthlyIncome}
            class:text-gray-800={!activeCharts.monthlyIncome} 
            class:dark:text-gray-200={!activeCharts.monthlyIncome}
            class:border-gray-300={!activeCharts.monthlyIncome}
            class:dark:border-gray-600={!activeCharts.monthlyIncome}
            onclick={() => activeCharts.monthlyIncome = !activeCharts.monthlyIncome}>
            Monthly Income
          </button>
          
          <button 
            class="px-3 py-1 text-sm rounded border"
            class:bg-primary-600={activeCharts.cumulativeIncome} 
            class:text-white={activeCharts.cumulativeIncome}
            class:border-primary-600={activeCharts.cumulativeIncome}
            class:bg-white={!activeCharts.cumulativeIncome} 
            class:dark:bg-gray-800={!activeCharts.cumulativeIncome}
            class:text-gray-800={!activeCharts.cumulativeIncome} 
            class:dark:text-gray-200={!activeCharts.cumulativeIncome}
            class:border-gray-300={!activeCharts.cumulativeIncome}
            class:dark:border-gray-600={!activeCharts.cumulativeIncome}
            onclick={() => activeCharts.cumulativeIncome = !activeCharts.cumulativeIncome}>
            Cumulative Income
          </button>
          
          <button 
            class="px-3 py-1 text-sm rounded border"
            class:bg-primary-600={activeCharts.monthlyDelta} 
            class:text-white={activeCharts.monthlyDelta}
            class:border-primary-600={activeCharts.monthlyDelta}
            class:bg-white={!activeCharts.monthlyDelta} 
            class:dark:bg-gray-800={!activeCharts.monthlyDelta}
            class:text-gray-800={!activeCharts.monthlyDelta} 
            class:dark:text-gray-200={!activeCharts.monthlyDelta}
            class:border-gray-300={!activeCharts.monthlyDelta}
            class:dark:border-gray-600={!activeCharts.monthlyDelta}
            onclick={() => activeCharts.monthlyDelta = !activeCharts.monthlyDelta}>
            Income Delta
          </button>
          
          <button 
            class="px-3 py-1 text-sm rounded border"
            class:bg-primary-600={activeCharts.incomeSource} 
            class:text-white={activeCharts.incomeSource}
            class:border-primary-600={activeCharts.incomeSource}
            class:bg-white={!activeCharts.incomeSource} 
            class:dark:bg-gray-800={!activeCharts.incomeSource}
            class:text-gray-800={!activeCharts.incomeSource} 
            class:dark:text-gray-200={!activeCharts.incomeSource}
            class:border-gray-300={!activeCharts.incomeSource}
            class:dark:border-gray-600={!activeCharts.incomeSource}
            onclick={() => activeCharts.incomeSource = !activeCharts.incomeSource}>
            Income Sources
          </button>
        </div>
      </div>
      
      <div class="grid grid-cols-1 gap-6">
        <!-- Income Comparison Chart -->
        {#if activeCharts.monthlyIncome}
          <div>
            <MonthlyIncomeChart
              {selectedProjection}
              {compareProjection}
              {selectedScenario}
              {compareScenario}
            />
          </div>
        {/if}
        
        <!-- Cumulative Income Chart -->
        {#if activeCharts.cumulativeIncome}
          <div>
            <CumulativeIncomeChart
              {selectedProjection}
              {compareProjection}
              {selectedScenario}
              {compareScenario}
            />
          </div>
        {/if}
        
        <!-- Monthly Income Delta Chart -->
        {#if activeCharts.monthlyDelta}
          <div>
            <MonthlyIncomeDeltaChart
              {selectedProjection}
              {compareProjection}
              {selectedScenario}
              {compareScenario}
            />
          </div>
        {/if}
        
        <!-- Income Source Breakdown Chart -->
        {#if activeCharts.incomeSource}
          <div>
            <IncomeSourceChart
              {selectedProjection}
              {compareProjection}
              {selectedScenario}
              {compareScenario}
            />
          </div>
        {/if}
      </div>
    </div>
  {:else}
    <div class="text-center py-8">
      <p class="text-gray-600 dark:text-gray-400">No projection data available. Click the Calculate button to run retirement projections.</p>
      <button 
        class="mt-4 px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 dark:bg-primary-700 dark:hover:bg-primary-600"
        onclick={runProjections}>
        Calculate Projections
      </button>
    </div>
  {/if}
</div>
