<script>
  import { onMount } from 'svelte';
  import ChartComponent from './ChartComponent.svelte';
  import annotationPlugin from 'chartjs-plugin-annotation';
  
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
  
  import { format, parseISO } from 'date-fns';

  // Prepare chart data when projections change
  $: if (selectedProjection && compareProjection) {
    prepareChartData();
  }
  
  // Initialize chart on mount to ensure it's properly rendered
  onMount(() => {
    if (selectedProjection && compareProjection) {
      prepareChartData();
    }
  });
  
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
    const dates = selectedProjection.yearlyData.map(item => {
      // Format dates as strings for Chart.js to avoid timezone issues
      return format(new Date(parseInt(item.year, 10), 0, 1), 'yyyy-MM-dd');
    });
    chartData.labels = dates;
    
    // Create datasets for selected scenario
    chartData.datasets.push({
      label: selectedScenario?.name || 'Scenario A',
      data: selectedProjection.yearlyData.map((item, index) => ({
        x: dates[index],
        y: Math.round(item.netIncome / 12) // Convert annual to monthly
      })),
      borderColor: 'rgb(59, 130, 246)', // Blue
      backgroundColor: 'rgba(59, 130, 246, 0.1)',
      borderWidth: 2,
      tension: 0.4,
      pointRadius: 2,
      pointHoverRadius: 5
    });
    
    // Create datasets for compare scenario
    chartData.datasets.push({
      label: compareScenario?.name || 'Scenario B',
      data: compareProjection.yearlyData.map((item, index) => {
        const matchingDate = index < dates.length ? dates[index] : null;
        return {
          x: matchingDate,
          y: Math.round(item.netIncome / 12) // Convert annual to monthly
        };
      }).filter(item => item.x !== null),
      borderColor: 'rgb(34, 197, 94)', // Green
      backgroundColor: 'rgba(34, 197, 94, 0.1)',
      borderWidth: 2,
      tension: 0.4,
      pointRadius: 2,
      pointHoverRadius: 5
    });
    
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
    
    // Add Social Security age markers if available
    if (selectedScenario?.data?.socialSecurity?.startAge) {
      const ssStartAge = selectedScenario.data.socialSecurity.startAge;
      const birthYear = selectedScenario.data.socialSecurity.birthYear || 1970;
      const ssYear = birthYear + ssStartAge;
      
      // Add to chart options annotations
      chartOptions.plugins.annotation.annotations.ssA = {
        type: 'line',
        xMin: new Date(ssYear, 0, 1),
        xMax: new Date(ssYear, 0, 1),
        borderColor: 'rgba(255, 159, 64, 0.8)',
        borderWidth: 2,
        borderDash: [5, 5],
        label: {
          content: 'A Social Security',
          display: true,
          position: 'start'
        }
      };
    }
    
    if (compareScenario?.data?.socialSecurity?.startAge) {
      const ssStartAge = compareScenario.data.socialSecurity.startAge;
      const birthYear = compareScenario.data.socialSecurity.birthYear || 1970;
      const ssYear = birthYear + ssStartAge;
      
      // Add to chart options annotations
      chartOptions.plugins.annotation.annotations.ssB = {
        type: 'line',
        xMin: new Date(ssYear, 0, 1),
        xMax: new Date(ssYear, 0, 1),
        borderColor: 'rgba(153, 102, 255, 0.8)',
        borderWidth: 2,
        borderDash: [5, 5],
        label: {
          content: 'B Social Security',
          display: true,
          position: 'start'
        }
      };
    }
  }
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
  <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">Monthly Income Comparison</h2>
  
  {#if selectedProjection && compareProjection}
    <div class="h-[400px]">
      <ChartComponent type="line" data={chartData} options={chartOptions} />
    </div>
  {:else}
    <div class="flex items-center justify-center h-[400px] bg-gray-50 dark:bg-gray-700 rounded">
      <p class="text-gray-500 dark:text-gray-400">No projection data available</p>
    </div>
  {/if}
</div>