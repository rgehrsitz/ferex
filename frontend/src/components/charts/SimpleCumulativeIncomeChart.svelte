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

  let chartData = $state<{ labels: string[]; datasets: any[] }>({ labels: [], datasets: [] });

  let chartOptions: Record<string, any> = {
    scales: {
      x: {
        title: {
          display: true,
          text: 'Age'
        }
      },
      y: {
        title: {
          display: true,
          text: 'Cumulative Income ($)'
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
        position: 'top'
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
          }
        }
      }
    }
  };

  // Generate data for the chart
  function prepareChartData(): void {
    if (!selectedProjection?.yearlyData || !compareProjection?.yearlyData) {
      return;
    }
    
    // Use age as the main label for X-axis
    const ageLabels = selectedProjection.yearlyData.map((d: any) => `Age ${d.age}`);
    
    // Calculate cumulative values
    let cumulativeA = 0;
    const dataScenarioA = selectedProjection.yearlyData.map((d: any) => {
      cumulativeA += d.netIncome;
      return cumulativeA;
    });
    
    let cumulativeB = 0;
    const dataScenarioB = compareProjection.yearlyData.map((d: any) => {
      cumulativeB += d.netIncome;
      return cumulativeB;
    });
    
    // Set chart data
    chartData = {
      labels: ageLabels,
      datasets: [
        {
          label: selectedScenario?.name || 'Scenario A',
          data: dataScenarioA,
          borderColor: 'rgb(59, 130, 246)', // Blue
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderWidth: 2,
          fill: true
        },
        {
          label: compareScenario?.name || 'Scenario B',
          data: dataScenarioB,
          borderColor: 'rgb(34, 197, 94)', // Green
          backgroundColor: 'rgba(34, 197, 94, 0.1)',
          borderWidth: 2,
          fill: true
        }
      ]
    };
  }
  
  // Update data when projections change (Svelte 5 runes mode)
  $effect(() => {
    if (selectedProjection && compareProjection) {
      prepareChartData();
    }
  });
  
  // Initialize on mount
  onMount(() => {
    if (selectedProjection && compareProjection) {
      prepareChartData();
    }
  });
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
  <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">Cumulative Income Comparison</h2>
  
  {#if selectedProjection && compareProjection}
    <SimpleChartComponent type="line" data={chartData} options={chartOptions} />
  {:else}
    <div class="flex items-center justify-center h-[400px] bg-gray-50 dark:bg-gray-700 rounded">
      <p class="text-gray-500 dark:text-gray-400">No projection data available</p>
    </div>
  {/if}
</div>