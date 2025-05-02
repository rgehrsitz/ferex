<script lang="ts">
  import { onMount } from 'svelte';
  import { CalculateRetirementProjection } from '../../wailsjs/go/main/App';
  import { main } from '../../wailsjs/go/models';
  import type { Scenario } from '../types/scenario';
  import MonthlyIncomeChart from './charts/SimpleMontlyIncomeChart.svelte';
  import CumulativeIncomeChart from './charts/SimpleCumulativeIncomeChart.svelte';
  import MonthlyIncomeDeltaChart from './charts/SimpleMonthlyDeltaChart.svelte';
  import IncomeSourceChart from './charts/SimpleIncomeSourceChart.svelte';

  export let selectedScenario: Scenario;
  export let compareScenario: Scenario;
  
  let loading = false;
  let error = '';
  let selectedProjection: any = null;
  let compareProjection: any = null;
  
  // Summary statistics
  let selectedNetTotal = 0;
  let compareNetTotal = 0;
  let netIncomeDifference = 0;
  let cumulativeDifference = 0;
  
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
      const scenarioInput2 = createProjectionInput(compareScenario);
      
      // Run projections in parallel
      const [projection1, projection2] = await Promise.all([
        CalculateRetirementProjection(scenarioInput1),
        CalculateRetirementProjection(scenarioInput2)
      ]);
      
      selectedProjection = projection1;
      compareProjection = projection2;
      
      // Calculate summary statistics
      if (projection1 && projection2) {
        selectedNetTotal = projection1.totalNetIncome;
        compareNetTotal = projection2.totalNetIncome;
        netIncomeDifference = compareNetTotal - selectedNetTotal;
        
        // Calculate cumulative difference over projection period
        const maxYears = Math.max(
          projection1.yearlyData?.length || 0,
          projection2.yearlyData?.length || 0
        );
        
        cumulativeDifference = 0;
        for (let i = 0; i < maxYears; i++) {
          const year1Income = projection1.yearlyData[i]?.netIncome || 0;
          const year2Income = projection2.yearlyData[i]?.netIncome || 0;
          cumulativeDifference += (year2Income - year1Income);
        }
      }
      
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
    pension.system = scenario.data.pension.retirementSystem;
    pension.high3Salary = scenario.data.pension.highThreeSalary;
    pension.yearsOfService = scenario.data.pension.yearsOfService;
    pension.ageAtRetirement = scenario.data.pension.retirementAge;
    pension.unusedSickLeaveMonths = Math.round(scenario.data.pension.unusedSickLeave / 174); // Convert hours to months
    pension.survivorBenefitOption = scenario.data.pension.survivorBenefit;
    pension.isPartTime = scenario.data.pension.isPartTime;
    pension.partTimeProrationFactor = scenario.data.pension.partTimeProrationFactor;
    input.pension = pension;
    
    // Social Security data
    if (scenario.data.socialSecurity) {
      const ss = new main.SocialSecurityInput();
      ss.startAge = scenario.data.socialSecurity.startAge;
      ss.estimatedMonthlyBenefit = scenario.data.socialSecurity.estimatedMonthlyBenefit;
      ss.isEligible = scenario.data.socialSecurity.isEligible;
      ss.birthYear = scenario.data.socialSecurity.birthYear;
      ss.birthMonth = scenario.data.socialSecurity.birthMonth;
      input.socialSecurity = ss;
    }
    
    // TSP data
    if (scenario.data.tsp) {
      const tsp = new main.TSPInput();
      tsp.currentBalance = scenario.data.tsp.currentBalance;
      tsp.traditionalBalance = scenario.data.tsp.traditionalBalance;
      tsp.rothBalance = scenario.data.tsp.rothBalance;
      tsp.annualContribution = scenario.data.tsp.annualContribution;
      tsp.expectedReturnRate = scenario.data.tsp.expectedReturnRate;
      tsp.withdrawalStrategy = scenario.data.tsp.withdrawalStrategy;
      tsp.fixedWithdrawalAmount = scenario.data.tsp.fixedWithdrawalAmount || 0;
      tsp.withdrawalPercentage = scenario.data.tsp.withdrawalPercentage || 0;
      tsp.withdrawalStartAge = scenario.data.tsp.withdrawalStartAge;
      input.tsp = tsp;
    }
    
    // Tax data
    if (scenario.data.tax) {
      const tax = new main.TaxInput();
      tax.filingStatus = scenario.data.tax.filingStatus;
      tax.stateOfResidence = scenario.data.tax.stateOfResidence;
      tax.additionalIncome = scenario.data.tax.additionalIncome;
      tax.additionalDeductions = scenario.data.tax.additionalDeductions;
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
    if (scenario.data.otherIncome) {
      const otherIncome = new main.OtherIncomeInput();
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
      input.otherIncome = otherIncome;
    }
    
    // Projection age range
    input.projectionStartAge = scenario.data.pension.retirementAge;
    input.projectionEndAge = 95; // Default to age 95
    
    return input;
  }
  
  // Run projections when component mounts
  onMount(() => {
    runProjections();
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
            {compareScenario.name} vs {selectedScenario.name}
          </div>
        </div>
        
        <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
          <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">{selectedScenario.name} Total</h3>
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
                <td class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200">{yearData.age}</td>
                <td class="py-2 px-4 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200">{yearData.year}</td>
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
            on:click={() => activeCharts.monthlyIncome = !activeCharts.monthlyIncome}>
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
            on:click={() => activeCharts.cumulativeIncome = !activeCharts.cumulativeIncome}>
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
            on:click={() => activeCharts.monthlyDelta = !activeCharts.monthlyDelta}>
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
            on:click={() => activeCharts.incomeSource = !activeCharts.incomeSource}>
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
        on:click={runProjections}>
        Calculate Projections
      </button>
    </div>
  {/if}
</div>
