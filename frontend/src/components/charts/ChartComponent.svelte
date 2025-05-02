<script>
  import { onMount, onDestroy } from 'svelte';
  import { Chart, registerables } from 'chart.js';
  import annotationPlugin from 'chartjs-plugin-annotation';
  
  // Register all Chart.js components needed
  Chart.register(...registerables);
  Chart.register(annotationPlugin);
  
  // Props
  export let type = 'line'; // line, bar, pie, etc.
  export let data = {
    labels: [],
    datasets: []
  };
  export let options = {};
  export let width = '100%';
  export let height = '400px';
  export let id = `chart-${Math.random().toString(36).substring(2, 9)}`;
  
  let canvas;
  let chart;
  
  onMount(() => {
    // Add a slight delay to ensure the canvas is fully rendered
    setTimeout(() => {
      if (canvas) {
        try {
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
        } catch (err) {
          console.error('Error creating chart:', err);
        }
      }
    }, 100);
  });
  
  // Update chart when data or options change
  $: if (chart && data) {
    try {
      chart.data = data;
      chart.update('active');
    } catch (err) {
      console.error('Error updating chart data:', err);
    }
  }
  
  $: if (chart && options) {
    try {
      chart.options = {
        responsive: true,
        maintainAspectRatio: false,
        ...options
      };
      chart.update('active');
    } catch (err) {
      console.error('Error updating chart options:', err);
    }
  }
  
  onDestroy(() => {
    if (chart) {
      chart.destroy();
    }
  });
</script>

<div style:width={width} style:height={height}>
  <canvas bind:this={canvas} id={id}></canvas>
</div>