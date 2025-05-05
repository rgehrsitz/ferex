<script lang="ts">
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
          text: 'Monthly Income Difference ($)'
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
            const value: number = context.raw;
            const formattedValue = new Intl.NumberFormat('en-US', { 
              style: 'currency', 
              currency: 'USD',
              signDisplay: 'always',
              maximumFractionDigits: 0 
            }).format(value);
            
            return `Difference: ${formattedValue}`;
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
    // Calculate differences (B - A)
    const deltaDiff: number[] = [];
    const bgColors: string[] = [];
    for (let i = 0; i < Math.min(selectedProjection.yearlyData.length, compareProjection.yearlyData.length); i++) {
      const monthlyA = selectedProjection.yearlyData[i].netIncome / 12;
      const monthlyB = compareProjection.yearlyData[i].netIncome / 12;
      const diff = Math.round(monthlyB - monthlyA);
      deltaDiff.push(diff);
      // Color based on positive/negative
      bgColors.push(diff >= 0 ? 'rgba(34, 197, 94, 0.5)' : 'rgba(239, 68, 68, 0.5)');
    }
    chartData = {
      labels: ageLabels,
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

  // Prepare chart data whenever projections or scenarios change
  $: prepareChartData();
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