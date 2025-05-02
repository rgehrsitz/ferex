<script>
  import { onMount } from 'svelte';
  import ChartComponent from './ChartComponent.svelte';
  
  // Props
  export let selectedProjection = null;
  export let compareProjection = null;
  export let selectedScenario = null;
  export let compareScenario = null;
  
  // State
  let activeScenario = 'A'; // Default to scenario A
  let chartData = {
    labels: [],
    datasets: []
  };
  
  // Chart options
  const chartOptions = {
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
  function toggleScenario(scenario) {
    activeScenario = scenario;
    prepareChartData();
  }
  
  // Prepare chart data when projections or active scenario changes
  $: if (selectedProjection && compareProjection && activeScenario) {
    prepareChartData();
  }
  
  function prepareChartData() {
    // Reset data
    chartData = {
      labels: [],
      datasets: []
    };
    
    // Get the relevant projection based on active scenario
    const projection = activeScenario === 'A' ? selectedProjection : compareProjection;
    const scenario = activeScenario === 'A' ? selectedScenario : compareScenario;
    
    if (!projection?.yearlyData || !scenario) {
      return;
    }
    
    // Generate dates for X-axis
    const dates = projection.yearlyData.map(item => new Date(item.year, 0, 1));
    chartData.labels = dates;
    
    // Income source colors
    const colors = {
      pension: 'rgb(34, 197, 94)', // Green
      socialSecurity: 'rgb(59, 130, 246)', // Blue
      tsp: 'rgb(249, 115, 22)', // Orange
      otherIncome: 'rgb(139, 92, 246)', // Purple
      taxBurden: 'rgb(239, 68, 68)' // Red (for negative values)
    };
    
    // Check which income sources are available
    const hasPension = projection.yearlyData.some(year => year.pensionIncome > 0);
    const hasSocialSecurity = projection.yearlyData.some(year => year.socialSecurityIncome > 0);
    const hasTSP = projection.yearlyData.some(year => year.tspWithdrawal > 0);
    const hasOtherIncome = projection.yearlyData.some(year => year.otherIncome > 0);
    const hasTax = projection.yearlyData.some(year => year.taxBurden > 0);
    
    // Create datasets for each income source
    if (hasPension) {
      chartData.datasets.push({
        label: 'FERS/CSRS Pension',
        data: projection.yearlyData.map((item, index) => ({
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
        data: projection.yearlyData.map((item, index) => ({
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
        data: projection.yearlyData.map((item, index) => ({
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
        data: projection.yearlyData.map((item, index) => ({
          x: dates[index],
          y: Math.round(item.otherIncome / 12) // Convert annual to monthly
        })),
        backgroundColor: colors.otherIncome,
        borderColor: colors.otherIncome,
        borderWidth: 1
      });
    }
    
    // Add retirement age marker if available
    if (scenario?.data?.pension?.retirementAge) {
      const retirementYear = new Date().getFullYear() + 
        (scenario.data.pension.retirementAge - new Date().getFullYear() + 1970);
      
      // Add to chart options annotations
      chartOptions.plugins.annotation.annotations.retirement = {
        type: 'line',
        xMin: new Date(retirementYear, 0, 1),
        xMax: new Date(retirementYear, 0, 1),
        borderColor: 'rgba(255, 99, 132, 0.8)',
        borderWidth: 2,
        borderDash: [5, 5],
        label: {
          content: 'Retirement',
          display: true,
          position: 'start'
        }
      };
    } else {
      // Clear retirement annotation if not available
      delete chartOptions.plugins.annotation.annotations.retirement;
    }
    
    // Add Social Security age marker if available
    if (scenario?.data?.socialSecurity?.startAge) {
      const ssStartAge = scenario.data.socialSecurity.startAge;
      const birthYear = scenario.data.socialSecurity.birthYear || 1970;
      const ssYear = birthYear + ssStartAge;
      
      // Add to chart options annotations
      chartOptions.plugins.annotation.annotations.socialSecurity = {
        type: 'line',
        xMin: new Date(ssYear, 0, 1),
        xMax: new Date(ssYear, 0, 1),
        borderColor: 'rgba(59, 130, 246, 0.8)',
        borderWidth: 2,
        borderDash: [5, 5],
        label: {
          content: 'Social Security',
          display: true,
          position: 'start'
        }
      };
    } else {
      // Clear Social Security annotation if not available
      delete chartOptions.plugins.annotation.annotations.socialSecurity;
    }
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