<script lang="ts">
import type { ChartOptions } from 'chart.js';
export let selectedProjection: any;
export let compareProjection: any;
export let selectedScenario: any;
export let compareScenario: any;

import { onMount } from 'svelte';
import ChartComponent from './ChartComponent.svelte';

let chartData: { labels: Date[]; datasets: any[] } = {
  labels: [],
  datasets: []
};

let chartOptions: ChartOptions = {
  scales: {
    x: {
      type: 'time',
      time: {
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
        text: 'Cumulative Income ($)'
      },
      grid: {
        color: 'rgba(200, 200, 200, 0.2)'
      },
      ticks: {
        callback: function(tickValue: string | number, index: number, ticks: any[]): string {
          const value = typeof tickValue === 'number' ? tickValue : parseFloat(tickValue);
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
  const dates = selectedProjection.yearlyData.map((item: any) => new Date(item.year, 0, 1));
  chartData.labels = dates;

  // Calculate cumulative values for Scenario A
  let cumulativeA = 0;
  const cumulativeDataA = selectedProjection.yearlyData.map((item: any, index: number) => {
    cumulativeA += item.netIncome;
    return {
      x: dates[index],
      y: cumulativeA
    };
  });

    // Calculate cumulative values for Scenario B
  let cumulativeB = 0;
  const cumulativeDataB = compareProjection.yearlyData.map((item: any, index: number) => {
    cumulativeB += item.netIncome;
    const matchingDate = index < dates.length ? dates[index] : null;
    return {
      x: matchingDate,
      y: cumulativeB
    };
  }).filter((item: { x: Date | null; y: number }) => item.x !== null);

  // Create datasets
  chartData.datasets.push({
    label: selectedScenario?.name || 'Scenario A',
    data: cumulativeDataA,
    borderColor: 'rgb(59, 130, 246)', // Blue
    backgroundColor: 'rgba(59, 130, 246, 0.1)',
    borderWidth: 2,
    tension: 0.4,
    pointRadius: 2,
    pointHoverRadius: 5,
    fill: true
  });

  chartData.datasets.push({
    label: compareScenario?.name || 'Scenario B',
    data: cumulativeDataB,
    borderColor: 'rgb(34, 197, 94)', // Green
    backgroundColor: 'rgba(34, 197, 94, 0.1)',
    borderWidth: 2,
    tension: 0.4,
    pointRadius: 2,
    pointHoverRadius: 5,
    fill: true
  });

  // Add retirement age markers if available
  if (selectedScenario?.data?.pension?.retirementAge) {
    const retirementYear = new Date().getFullYear() + 
      (selectedScenario.data.pension.retirementAge - new Date().getFullYear() + 1970);
    // Add to chart options annotations (with undefined checks)
    if (chartOptions.plugins && chartOptions.plugins.annotation && chartOptions.plugins.annotation.annotations) {
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
  }

  if (compareScenario?.data?.pension?.retirementAge) {
    const retirementYear = new Date().getFullYear() + 
      (compareScenario.data.pension.retirementAge - new Date().getFullYear() + 1970);
    // Add to chart options annotations (with undefined checks)
    if (chartOptions.plugins && chartOptions.plugins.annotation && chartOptions.plugins.annotation.annotations) {
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
}
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
  <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-4">Cumulative Income Comparison</h2>
  
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