<script lang="ts">
  import { onMount } from 'svelte';
  import Chart from 'chart.js/auto';
  import 'chartjs-plugin-annotation';

  // Props
  const { 
    type = 'line', 
    data = { labels: [], datasets: [] }, 
    options = {}, 
    height = '400px' 
  } = $props<{
    type?: string;
    data?: { labels: string[]; datasets: any[] };
    options?: Record<string, any>;
    height?: string;
  }>();

  let canvas: HTMLCanvasElement | null = null;
  let chart: Chart | null = null;

  onMount(() => {
    if (!canvas) return;
    chart = new Chart(canvas, {
      type: type as any, // Chart.js type
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

  // Svelte 5 idiom: use $effect for reactivity
  $effect(() => {
    if (chart && data) {
      chart.data = data;
      chart.update();
    }
  });

  $effect(() => {
    if (chart && options) {
      chart.options = {
        responsive: true,
        maintainAspectRatio: false,
        ...options
      };
      chart.update();
    }
  });
</script>

<div style="width: 100%; height: {height};">
  <canvas bind:this={canvas}></canvas>
</div>