<script lang="ts">
  export let selectedScenario;
  export let compareScenario;
</script>
<div class="flex gap-8">
  <div class="flex-1 bg-gray-100 p-4 rounded">
    <h3 class="font-bold mb-2">{selectedScenario.name}</h3>
    <p>[Inputs/results for this scenario]</p>
  </div>
  <div class="flex-1 bg-gray-100 p-4 rounded">
    <h3 class="font-bold mb-2">{compareScenario.name}</h3>
    <p>[Inputs/results for this scenario]</p>
  </div>
</div>
<div class="mt-4 bg-white dark:bg-gray-800 p-4 rounded shadow">
  <h4 class="font-bold mb-2">Comparison Visualizations</h4>
  <p>[Charts/graphs will go here]</p>
</div>
