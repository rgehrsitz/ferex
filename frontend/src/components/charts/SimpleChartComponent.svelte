<script lang="ts">
  import { onMount } from 'svelte';
  import Chart from 'chart.js/auto';
  import 'chartjs-plugin-annotation';

  // Props
  export let type: string = 'line'; // line, bar, pie, etc.
  export let data: { labels: string[]; datasets: any[] } = { labels: [], datasets: [] };
  export let options: Record<string, any> = {};
  export let height: string = '400px';

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

  // Svelte 5 idiom: $: for reactivity
  $: if (chart && data) {
    chart.data = data;
    chart.update();
  }

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