<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import ApexCharts from 'apexcharts';
  export let result: Record<string, any> = {};
  let chartEl: HTMLDivElement;
  let chart: ApexCharts | null = null; // Svelte 5 idiom: plain let for local state

  function getPercentileSeries(result) {
    if (!result || !result.YearlyBalances) return { series: [], years: 0 };
    const years = result.YearlyBalances[0]?.length || 0;
    const percentiles = [10, 50, 90];
    const percentileSeries = percentiles.map(() => []);
    for (let y = 0; y < years; y++) {
      const balances = result.YearlyBalances.map((sim: number[]) => sim[y]);
      balances.sort((a, b) => a - b);
      percentiles.forEach((p, i) => {
        const idx = Math.floor((p / 100) * (balances.length - 1));
        percentileSeries[i].push(balances[idx]);
      });
    }
    return {
      series: [
        { name: '10th Percentile', data: percentileSeries[0] },
        { name: 'Median', data: percentileSeries[1] },
        { name: '90th Percentile', data: percentileSeries[2] }
      ],
      years
    };
  }

  function renderChart() {
    if (!chartEl || !result || !result.YearlyBalances) return;
    const { series, years } = getPercentileSeries(result);
    if (chart) {
      chart.destroy();
      chart = null;
    }
    chart = new ApexCharts(chartEl, {
      chart: { type: 'line', height: 350, toolbar: { show: false } },
      series,
      xaxis: {
        categories: Array.from({length: years}, (_, i) => `Year ${i+1}`),
        title: { text: 'Year' }
      },
      yaxis: {
        title: { text: 'Balance ($)' },
        min: 0
      },
      stroke: { width: 2 },
      markers: { size: 0 },
      legend: { position: 'top' },
      title: { text: 'Monte Carlo Projected Balances', align: 'center' },
      colors: ['#ef4444', '#2563eb', '#059669'],
      tooltip: { y: { formatter: val => `$${Math.round(val).toLocaleString()}` } }
    });
    chart.render();
  }

  $: if (result && result.YearlyBalances) renderChart(); // Svelte 5 idiom: $: for reactive effects

  // Svelte 5 idiom: ensure runes and correct reactivity
  onDestroy(() => { if (chart) chart.destroy(); });
</script>

{#if result && result.SuccessRate !== undefined}
  <div class="p-4 bg-white rounded shadow my-4 max-w-3xl mx-auto">
    <h2 class="text-lg font-bold mb-2">Monte Carlo Results</h2>
    <div class="mb-2">Success Rate: <span class="font-mono">{(result.SuccessRate * 100).toFixed(1)}%</span></div>
    <div class="mb-2">Median Ending Balance: <span class="font-mono">${result.Percentiles?.[50]?.toLocaleString() || 'N/A'}</span></div>
    <div class="mb-2">10th-90th Percentile Range: <span class="font-mono">${result.Percentiles?.[10]?.toLocaleString() || 'N/A'} - ${result.Percentiles?.[90]?.toLocaleString() || 'N/A'}</span></div>
    <div bind:this={chartEl} class="w-full" style="min-height:350px;"></div>
  </div>
{/if}
