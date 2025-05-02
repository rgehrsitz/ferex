<script>
  import { onMount } from 'svelte';
  import SimpleChartComponent from './SimpleChartComponent.svelte';
  
  // Props
  export let selectedProjection = null;
  export let compareProjection = null;
  export let selectedScenario = null;
  export let compareScenario = null;
  
  let chartData = {
    labels: [],
    datasets: []
  };
  
  let chartOptions = {
    scales: {
      y: {
        title: {
          display: true,
          text: 'Monthly Income Difference ($)'
        },
        ticks: {
          callback: (value) => {
            return '$' + value.toLocaleString();
          }
        }
      }
    },
    plugins: {
      legend: {
        position: 'top'
      }
    }
  };
  
  // Generate data for the chart
  function prepareChartData() {
    if (!selectedProjection?.yearlyData || !compareProjection?.yearlyData) {
      return;
    }
    
    // Basic labels (years)
    const labels = selectedProjection.yearlyData.map(d => d.year.toString());
    
    // Calculate differences (B - A)
    const deltaDiff = [];
    const bgColors = [];
    
    for (let i = 0; i < Math.min(selectedProjection.yearlyData.length, compareProjection.yearlyData.length); i++) {
      const monthlyA = selectedProjection.yearlyData[i].netIncome / 12;
      const monthlyB = compareProjection.yearlyData[i].netIncome / 12;
      const diff = Math.round(monthlyB - monthlyA);
      deltaDiff.push(diff);
      
      // Color based on positive/negative
      bgColors.push(diff >= 0 ? 'rgba(34, 197, 94, 0.5)' : 'rgba(239, 68, 68, 0.5)');
    }
    
    // Set chart data
    chartData = {
      labels,
      datasets: [
        {
          label: `Monthly Difference (${compareScenario?.name || 'B'} - ${selectedScenario?.name || 'A'})`,
          data: deltaDiff,
          backgroundColor: bgColors,
          borderColor: 'rgba(217, 70, 239, 0.8)', // Purple
          borderWidth: 1
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
  <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">Monthly Income Delta (B - A)</h2>
  
  {#if selectedProjection && compareProjection}
    <SimpleChartComponent type="bar" data={chartData} options={chartOptions} />
    <div class="text-sm text-gray-500 dark:text-gray-400 mt-2">
      <span class="inline-block w-3 h-3 bg-green-500 mr-1"></span> Positive values: {compareScenario?.name || 'Scenario B'} provides more income
      <span class="inline-block w-3 h-3 bg-red-500 ml-3 mr-1"></span> Negative values: {selectedScenario?.name || 'Scenario A'} provides more income
    </div>
  {:else}
    <div class="flex items-center justify-center h-[400px] bg-gray-50 dark:bg-gray-700 rounded">
      <p class="text-gray-500 dark:text-gray-400">No projection data available</p>
    </div>
  {/if}
</div>