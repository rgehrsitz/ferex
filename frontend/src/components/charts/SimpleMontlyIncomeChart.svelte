<script lang="ts">
  import { onMount } from 'svelte';
  import SimpleChartComponent from './SimpleChartComponent.svelte';

  // Props
  export let selectedProjection: any = null;
  export let compareProjection: any = null;
  export let selectedScenario: any = null;
  export let compareScenario: any = null;

  let chartData: { labels: string[]; datasets: any[] } = { labels: [], datasets: [] };

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
        position: 'top'
      },
      tooltip: {
        callbacks: {
          title: (context: any) => {
            // Extract age from the label (e.g., "Age 65" → "Age 65")
            return context[0].label;
          },
          label: (context: any) => {
            // Format the value as currency
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
    
    // Convert income data to monthly values
    const dataScenarioA = selectedProjection.yearlyData.map((d: any) => Math.round(d.netIncome / 12));
    const dataScenarioB = compareProjection.yearlyData.map((d: any) => Math.round(d.netIncome / 12));
    
    // Set chart data
    chartData = {
      labels: ageLabels,
      datasets: [
        {
          label: selectedScenario?.name || 'Scenario A',
          data: dataScenarioA,
          borderColor: 'rgb(59, 130, 246)', // Blue
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderWidth: 2
        },
        {
          label: compareScenario?.name || 'Scenario B',
          data: dataScenarioB,
          borderColor: 'rgb(34, 197, 94)', // Green
          backgroundColor: 'rgba(34, 197, 94, 0.1)',
          borderWidth: 2
        }
      ]
    };
  }
  
  // Update data when projections change
  $: if (selectedProjection && compareProjection) {
    prepareChartData();
  }
  
  // Initialize on mount
  onMount(() => {
    if (selectedProjection && compareProjection) {
      prepareChartData();
    }
  });
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
  <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">Monthly Income Comparison</h2>
  
  {#if selectedProjection && compareProjection}
    <SimpleChartComponent type="line" data={chartData} options={chartOptions} />
  {:else}
    <div class="flex items-center justify-center h-[400px] bg-gray-50 dark:bg-gray-700 rounded">
      <p class="text-gray-500 dark:text-gray-400">No projection data available</p>
    </div>
  {/if}
</div>