<script>
  import { onMount, onDestroy } from 'svelte';
  import Chart from 'chart.js/auto';
  import 'chartjs-plugin-annotation';
  
  // Props
  export let type = 'line'; // line, bar, pie, etc.
  export let data = {
    labels: [],
    datasets: []
  };
  export let options = {};
  export let height = '400px';
  
  let canvas;
  let chart;
  
  onMount(() => {
    if (!canvas) return;
    
    // Create chart instance
    chart = new Chart(canvas, {
      type,
      data,
      options: {
        responsive: true,
        maintainAspectRatio: false,
        ...options
      }
    });
    
    return () => {
      if (chart) chart.destroy();
    };
  });
  
  // Update chart when data changes
  $: if (chart && data) {
    chart.data = data;
    chart.update();
  }
  
  // Update chart when options change
  $: if (chart && options) {
    chart.options = {
      responsive: true,
      maintainAspectRatio: false,
      ...options
    };
    chart.update();
  }
</script>

<div style="width: 100%; height: {height};">
  <canvas bind:this={canvas}></canvas>
</div>