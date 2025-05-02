<script>
  import { onMount } from 'svelte';
  import ChartComponent from './ChartComponent.svelte';
  
  // Props
  export let selectedProjection = null;
  export let compareProjection = null;
  export let selectedScenario = null;
  export let compareScenario = null;
  
  let chartData = {
    labels: [],
    datasets: []
  };
  
  // Chart options
  const chartOptions = {
    scales: {
      x: {
        type: 'time',
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
        }
      },
      y: {
        title: {
          display: true,
          text: 'Monthly Income Difference ($)'
        },
        grid: {
          color: 'rgba(200, 200, 200, 0.2)'
        },
        ticks: {
          callback: (value) => {
            return '$' + value.toLocaleString();
          }
        }
      }
    },
    plugins: {
      tooltip: {
        callbacks: {
          label: function(context) {
            let label = context.dataset.label || '';
            if (label) {
              label += ': ';
            }
            if (context.parsed.y !== null) {
              label += new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', signDisplay: 'always' }).format(context.parsed.y);
            }
            return label;
          }
        }
      },
      legend: {
        position: 'top',
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
  
  // Prepare chart data when projections change
  $: if (selectedProjection && compareProjection) {
    prepareChartData();
  }
  
  function prepareChartData() {
    // Reset data
    chartData = {
      labels: [],
      datasets: []
    };
    
    if (!selectedProjection?.yearlyData || !compareProjection?.yearlyData) {
      return;
    }
    
    // Generate dates for X-axis
    const dates = selectedProjection.yearlyData.map(item => new Date(item.year, 0, 1));
    chartData.labels = dates;
    
    // Calculate monthly differences (B - A)
    const deltaData = [];
    
    for (let i = 0; i < dates.length; i++) {
      const scenarioAIncome = i < selectedProjection.yearlyData.length ? selectedProjection.yearlyData[i].netIncome / 12 : 0;
      const scenarioBIncome = i < compareProjection.yearlyData.length ? compareProjection.yearlyData[i].netIncome / 12 : 0;
      
      deltaData.push({
        x: dates[i],
        y: Math.round(scenarioBIncome - scenarioAIncome) // B - A
      });
    }
    
    // Create dataset for the difference
    chartData.datasets.push({
      label: 'Monthly Difference (B - A)',
      data: deltaData,
      borderColor: 'rgb(217, 70, 239)', // Purple
      backgroundColor: function(context) {
        const value = context.raw.y;
        return value >= 0 ? 'rgba(34, 197, 94, 0.2)' : 'rgba(239, 68, 68, 0.2)';
      },
      borderWidth: 2,
      tension: 0.1,
      pointRadius: 2,
      pointHoverRadius: 5,
      fill: true
    });
    
    // Add a zero line
    chartOptions.plugins.annotation.annotations.zeroLine = {
      type: 'line',
      yMin: 0,
      yMax: 0,
      borderColor: 'rgba(100, 100, 100, 0.5)',
      borderWidth: 1,
      borderDash: [5, 5]
    };
    
    // Add retirement age markers if available
    if (selectedScenario?.data?.pension?.retirementAge) {
      const retirementYear = new Date().getFullYear() + 
        (selectedScenario.data.pension.retirementAge - new Date().getFullYear() + 1970);
      
      // Add to chart options annotations
      chartOptions.plugins.annotation.annotations.retirementA = {
        type: 'line',
        xMin: new Date(retirementYear, 0, 1),
        xMax: new Date(retirementYear, 0, 1),
        borderColor: 'rgba(255, 99, 132, 0.8)',
        borderWidth: 2,
        borderDash: [5, 5],
        label: {
          content: 'A Retirement',
          display: true,
          position: 'start'
        }
      };
    }
    
    if (compareScenario?.data?.pension?.retirementAge) {
      const retirementYear = new Date().getFullYear() + 
        (compareScenario.data.pension.retirementAge - new Date().getFullYear() + 1970);
      
      // Add to chart options annotations
      chartOptions.plugins.annotation.annotations.retirementB = {
        type: 'line',
        xMin: new Date(retirementYear, 0, 1),
        xMax: new Date(retirementYear, 0, 1),
        borderColor: 'rgba(75, 192, 192, 0.8)',
        borderWidth: 2,
        borderDash: [5, 5],
        label: {
          content: 'B Retirement',
          display: true,
          position: 'start'
        }
      };
    }
  }
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
  <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">Monthly Income Delta (B - A)</h2>
  
  {#if selectedProjection && compareProjection}
    <div class="h-[400px]">
      <ChartComponent type="bar" data={chartData} options={chartOptions} />
    </div>
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