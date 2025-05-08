<script lang="ts">
  import { onMount } from 'svelte';
  import SimpleChartComponent from './SimpleChartComponent.svelte';

  // Props
  const { selectedProjection, compareProjection, selectedScenario, compareScenario } = $props<{
    selectedProjection: any;
    compareProjection: any;
    selectedScenario: any;
    compareScenario: any;
  }>();

  // State
  let activeScenario = $state<'A' | 'B'>('A'); // Default to scenario A
  let chartData = $state<{ labels: string[]; datasets: any[] }>({ labels: [], datasets: [] });

  let chartOptions: Record<string, any> = {
    scales: {
      x: {
        stacked: true,
        title: {
          display: true,
          text: 'Age'
        }
      },
      y: {
        stacked: true,
        title: {
          display: true,
          text: 'Monthly Income ($)'
        },
        ticks: {
          callback: (value: number) => {
            return '$' + value.toLocaleString();
          }
        }
      }
    },
    plugins: {
      legend: {
        position: 'bottom'
      },
      tooltip: {
        callbacks: {
          title: (context: any) => {
            return context[0].label;
          },
          label: (context: any) => {
            return `${context.dataset.label}: ${new Intl.NumberFormat('en-US', { 
              style: 'currency', 
              currency: 'USD',
              maximumFractionDigits: 0 
            }).format(context.raw)}`;
          },
          footer: (context: any) => {
            // Calculate total across all datasets for this x value
            const total = context.reduce((sum: number, item: any) => sum + item.raw, 0);
            return `Total: ${new Intl.NumberFormat('en-US', { 
              style: 'currency', 
              currency: 'USD',
              maximumFractionDigits: 0 
            }).format(total)}`;
          }
        }
      }
    }
  };

  // Toggle between scenarios
  function toggleScenario(scenario: 'A' | 'B') {
    activeScenario = scenario;
    prepareChartData();
  }

  // Generate data for the chart
  function prepareChartData(): void {
    // Get the relevant projection based on active scenario
    const projection = activeScenario === 'A' ? selectedProjection : compareProjection;
    if (!projection?.yearlyData) {
      return;
    }
    // Use age as the main label for X-axis
    const ageLabels = projection.yearlyData.map((d: any) => `Age ${d.age}`);
    // Create arrays for different income sources
    const pensionData = projection.yearlyData.map((d: any) => Math.round(d.pensionIncome / 12));
    const ssData = projection.yearlyData.map((d: any) => Math.round(d.socialSecurityIncome / 12));
    const tspData = projection.yearlyData.map((d: any) => Math.round(d.tspWithdrawal / 12));
    const otherData = projection.yearlyData.map((d: any) => Math.round(d.otherIncome / 12));
    // Set chart data
    chartData = {
      labels: ageLabels,
      datasets: [
        {
          label: 'FERS/CSRS Pension',
          data: pensionData,
          backgroundColor: 'rgb(34, 197, 94)', // Green
          borderColor: 'rgb(34, 197, 94)',
          borderWidth: 1
        },
        {
          label: 'Social Security',
          data: ssData,
          backgroundColor: 'rgb(59, 130, 246)', // Blue
          borderColor: 'rgb(59, 130, 246)',
          borderWidth: 1
        },
        {
          label: 'TSP Withdrawals',
          data: tspData,
          backgroundColor: 'rgb(249, 115, 22)', // Orange
          borderColor: 'rgb(249, 115, 22)',
          borderWidth: 1
        },
        {
          label: 'Other Income',
          data: otherData,
          backgroundColor: 'rgb(139, 92, 246)', // Purple
          borderColor: 'rgb(139, 92, 246)',
          borderWidth: 1
        }
      ]
    };
  }
  
  // Update data when projections or active scenario changes (Svelte 5 runes mode)
  $effect(() => {
    if ((activeScenario === 'A' && selectedProjection) || 
        (activeScenario === 'B' && compareProjection)) {
      prepareChartData();
    }
  });
  
  // Initialize on mount
  onMount(() => {
    if ((activeScenario === 'A' && selectedProjection) || 
        (activeScenario === 'B' && compareProjection)) {
      prepareChartData();
    }
  });
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
        onclick={() => toggleScenario('A')}>
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
        onclick={() => toggleScenario('B')}>
        {compareScenario?.name || 'Scenario B'}
      </button>
    </div>
  </div>
  
  {#if (activeScenario === 'A' && selectedProjection) || (activeScenario === 'B' && compareProjection)}
    <SimpleChartComponent type="bar" data={chartData} options={chartOptions} />
  {:else}
    <div class="flex items-center justify-center h-[400px] bg-gray-50 dark:bg-gray-700 rounded">
      <p class="text-gray-500 dark:text-gray-400">No projection data available</p>
    </div>
  {/if}
</div>