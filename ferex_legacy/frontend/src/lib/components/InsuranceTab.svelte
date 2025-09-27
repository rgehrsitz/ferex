<script lang="ts">
  import type { ScenarioInput } from '../../types';
  import FEGLISection from './FEGLISection.svelte';
  import FEHBSection from './FEHBSection.svelte';
  
  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
</script>

<div class="space-y-8">
  <!-- FEGLI Section -->
  <section class="space-y-4 p-4 border border-gray-200 rounded-md shadow-sm">
    <h3 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-3">
      Federal Employees' Group Life Insurance (FEGLI)
    </h3>
    <FEGLISection bind:inputs={inputs} />
  </section>

  <!-- FEHB Section -->
  <section class="space-y-4 p-4 border border-gray-200 rounded-md shadow-sm">
    <h3 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-3">
      Federal Employees Health Benefits (FEHB)
    </h3>
    <FEHBSection bind:inputs={inputs} />
  </section>
</div>
