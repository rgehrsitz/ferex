<script lang="ts">
  import type { Scenario } from '../types/scenario.js';
  export let scenarios: Scenario[] = [];
  export let selectedScenarioId: number;
  export let compareScenarioId;
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
</script>

<aside class="w-64 bg-white text-gray-900 p-4 border-r min-h-screen">
  <h2 class="font-bold mb-2">Scenarios</h2>
  <ul>
    {#each scenarios as scenario}
      <li class="mb-1 flex items-center gap-1">
        <button
          class="flex-1 text-left px-2 py-1 rounded {selectedScenarioId === scenario.id ? 'bg-blue-100' : ''}"
          on:click={() => dispatch('select', scenario.id)}>
          {scenario.name}
        </button>
        <button title="Duplicate" on:click={() => dispatch('duplicate', scenario.id)}>⎘</button>
        <button title="Delete" on:click={() => dispatch('delete', scenario.id)}>🗑️</button>
        {#if selectedScenarioId !== scenario.id}
          <button title="Compare" on:click={() => dispatch('compare', scenario.id)}>Compare</button>
        {/if}
      </li>
    {/each}
  </ul>
  <button class="mt-2 w-full bg-blue-600 text-white rounded p-2" on:click={() => dispatch('add')}>+ Add Scenario</button>
</aside>
