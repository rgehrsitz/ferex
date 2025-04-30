<script lang="ts">
  import type { Scenario } from '../types/scenario.js';
  export let scenarios: Scenario[] = [];
  export let selectedScenarioId: number;
  export let compareScenarioId;
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
</script>

<aside class="w-64 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 p-4 border-r border-gray-200 dark:border-gray-700 min-h-screen">
  <h2 class="font-bold mb-2">Scenarios</h2>
  <ul>
    {#each scenarios as scenario}
      <li class="mb-1 flex items-center gap-1">
        <input
          class="flex-1 text-left px-2 py-1 rounded bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 border border-gray-300 dark:border-gray-700 {selectedScenarioId === scenario.id ? 'bg-blue-100 dark:bg-blue-900' : ''}"
          value={scenario.name}
          readonly
        />
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
