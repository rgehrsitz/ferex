<script lang="ts">
import type { ChartOptions, ChartData } from 'chart.js';
import { onMount } from 'svelte';
import ChartComponent from './ChartComponent.svelte';

// Props
export let selectedProjection: any = null;
export let compareProjection: any = null;
export let selectedScenario: any = null;
export let compareScenario: any = null;

// State
let activeScenario: string = 'A'; // Default to scenario A
let chartData: ChartData = {
  labels: [],
  datasets: []
};

$: if (selectedProjection && compareProjection && selectedScenario && compareScenario) {
  prepareChartData();
}

// Chart options
const chartOptions: ChartOptions = {
  scales: {
    x: {
      type: 'time',
      stacked: true,
      time: {
        unit: 'year',
        displayFormats: {
          year: 'yyyy'
        }
      },
      title: {
        display: true,
        text: 'Date'
      },
      grid: {
        color: 'rgba(200, 200, 200, 0.2)'
      },
      ticks: {
        callback: (value: number | string) => {
          return '$' + value.toLocaleString();
        }
      }
    },
    y: {
      stacked: true,
      title: {
        display: true,
        text: 'Monthly Income ($)'
      },
      grid: {
        color: 'rgba(200, 200, 200, 0.2)'
      },
      ticks: {
        callback: (value: number | string) => {
          return '$' + value.toLocaleString();
        }
      }
    }
  },
  plugins: {
    tooltip: {
      callbacks: {
        label: function(context: any) {
          let label = context.dataset.label || '';
          if (label) {
            label += ': ';
          }
          if (context.parsed.y !== null) {
            label += new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(context.parsed.y);
          }
          return label;
        }
      }
    },
    legend: {
      position: 'bottom',
      labels: {
        usePointStyle: true,
        padding: 20
      }
    },
    annotation: {
      annotations: {}
    }
  },
  responsive: true,
  maintainAspectRatio: false
};

// Toggle between scenarios
function toggleScenario(scenario: string) {
  activeScenario = scenario;
  prepareChartData();
}

// Prepare chart data when projections or active scenario changes
function prepareChartData() {
  // Reset data
  chartData = {
    labels: [],
    datasets: []
  };

  // Get the relevant projection based on active scenario
  const projection = activeScenario === 'A' ? selectedProjection : compareProjection;
  const scenario = activeScenario === 'A' ? selectedScenario : compareScenario;
  const dates = projection.yearlyData.map((year: Record<string, number>, index: number) => new Date(year.year, 0, 1));
  chartData.labels = dates;

  // Income source colors
  const colors = {
    pension: 'rgb(34, 197, 94)', // Green
    socialSecurity: 'rgb(59, 130, 246)', // Blue
    tsp: 'rgb(249, 115, 22)', // Orange
    otherIncome: 'rgb(139, 92, 246)', // Purple
    taxBurden: 'rgb(239, 68, 68)' // Red (for negative values)
  };

  if (!projection?.yearlyData || !scenario) {
    return;
  }

  // Check which income sources are available
  const hasPension = projection.yearlyData.some((year: Record<string, number>) => year.pensionIncome > 0);
  const hasSocialSecurity = projection.yearlyData.some((year: Record<string, number>) => year.socialSecurityIncome > 0);
  const hasTSP = projection.yearlyData.some((year: Record<string, number>) => year.tspWithdrawal > 0);
  const hasOtherIncome = projection.yearlyData.some((year: Record<string, number>) => year.otherIncome > 0);
  const hasTax = projection.yearlyData.some((year: Record<string, number>) => year.taxBurden > 0);

  // Create datasets for each income source
  if (hasPension) {
    chartData.datasets.push({
      label: 'FERS/CSRS Pension',
      data: projection.yearlyData.map((item: Record<string, number>, index: number) => ({
        x: dates[index],
        y: Math.round(item.pensionIncome / 12) // Convert annual to monthly
      })),
      backgroundColor: colors.pension,
      borderColor: colors.pension,
      borderWidth: 1
    });
  }
  
  if (hasSocialSecurity) {
    chartData.datasets.push({
      label: 'Social Security',
      data: projection.yearlyData.map((item: Record<string, number>, index: number) => ({
        x: dates[index],
        y: Math.round(item.socialSecurityIncome / 12) // Convert annual to monthly
      })),
      backgroundColor: colors.socialSecurity,
      borderColor: colors.socialSecurity,
      borderWidth: 1
    });
  }
  
  if (hasTSP) {
    chartData.datasets.push({
      label: 'TSP Withdrawals',
      data: projection.yearlyData.map((item: Record<string, number>, index: number) => ({
        x: dates[index],
        y: Math.round(item.tspWithdrawal / 12) // Convert annual to monthly
      })),
      backgroundColor: colors.tsp,
      borderColor: colors.tsp,
      borderWidth: 1
    });
  }
  
  if (hasOtherIncome) {
    chartData.datasets.push({
      label: 'Other Income',
      data: projection.yearlyData.map((item: Record<string, number>, index: number) => ({
        x: dates[index],
        y: Math.round(item.otherIncome / 12) // Convert annual to monthly
      })),
      backgroundColor: colors.otherIncome,
      borderColor: colors.otherIncome,
      borderWidth: 1
    });
  }
  // (Retirement and Social Security annotation logic continues here...)
}

</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
  <div class="flex justify-between items-center mb-4">
    <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200">Income Source Breakdown</h2>
    
    <div class="flex space-x-2">
      <button 
        class="px-3 py-1 text-sm rounded" 
        class:bg-primary-600={activeScenario === 'A'} 
        class:text-white={activeScenario === 'A'}
        class:bg-gray-100={activeScenario !== 'A'} 
        class:dark:bg-gray-700={activeScenario !== 'A'}
        class:text-gray-800={activeScenario !== 'A'} 
        class:dark:text-gray-200={activeScenario !== 'A'}
        on:click={() => toggleScenario('A')}>
        {selectedScenario?.name || 'Scenario A'}
      </button>
      
      <button 
        class="px-3 py-1 text-sm rounded" 
        class:bg-primary-600={activeScenario === 'B'} 
        class:text-white={activeScenario === 'B'}
        class:bg-gray-100={activeScenario !== 'B'} 
        class:dark:bg-gray-700={activeScenario !== 'B'}
        class:text-gray-800={activeScenario !== 'B'} 
        class:dark:text-gray-200={activeScenario !== 'B'}
        on:click={() => toggleScenario('B')}>
        {compareScenario?.name || 'Scenario B'}
      </button>
    </div>
  </div>
  
  {#if (activeScenario === 'A' && selectedProjection) || (activeScenario === 'B' && compareProjection)}
    <div class="h-[400px]">
      <ChartComponent type="bar" data={chartData} options={chartOptions} />
    </div>
  {:else}
    <div class="flex items-center justify-center h-[400px] bg-gray-50 dark:bg-gray-700 rounded">
      <p class="text-gray-500 dark:text-gray-400">No projection data available</p>
    </div>
  {/if}
</div>