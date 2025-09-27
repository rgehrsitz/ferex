<script lang="ts">
  import type { Scenario } from '../../types';
  import { createEventDispatcher } from 'svelte'; 

  // Props definition using Svelte 5 runes
  let { scenarios, activeScenarioId, activeScenarioName } = $props<{
    scenarios: Scenario[];
    activeScenarioId: string | undefined;
    activeScenarioName: string;
  }>();

  const dispatch = createEventDispatcher<{
    selectScenario: string; // payload is scenarioId
    updateScenarioName: string; // payload is newName (activeScenarioId is implicit)
    addScenario: void;
    duplicateScenario: void;
    deleteScenario: void;
  }>();

  let scenarioNameInput = $state(activeScenarioName);

  // Effect to synchronize local scenarioNameInput when activeScenarioName prop changes
  $effect(() => {
    scenarioNameInput = activeScenarioName;
  });

  function handleNameInput() {
    // Dispatch event when the input value changes
    // This allows App.svelte to handle the actual state update.
    dispatch('updateScenarioName', scenarioNameInput);
  }

  // Computed property for disabling delete button
  let isDeleteDisabled = $derived(scenarios.length <= 1 || !activeScenarioId);

</script>

<div class="mb-8 p-6 bg-slate-50 rounded-lg shadow-md">
  <h2 class="text-xl font-semibold text-gray-700 mb-4">Manage Scenarios</h2>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4 items-end">
    <div>
      <label for="scenario-select" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Active Scenario:</label>
      <select
        id="scenario-select"
        value={activeScenarioId ?? ''}
        onchange={(e) => dispatch('selectScenario', (e.target as HTMLSelectElement).value)}
        class="text-sm rounded-lg block w-full p-2.5 bg-white text-gray-900 border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 dark:focus:ring-blue-500 dark:focus:border-blue-500"
      >
        {#each scenarios as scenario (scenario.id)}
          <option value={scenario.id}>{scenario.name}</option>
        {/each}
      </select>
    </div>
    <div>
      <label for="scenario-name" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Scenario Name:</label>
      <input
        type="text"
        id="scenario-name"
        bind:value={scenarioNameInput}
        oninput={handleNameInput} 
        class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
        placeholder="Enter scenario name"
        disabled={!activeScenarioId}
      />
    </div>
  </div>
  <div class="mt-6 flex flex-wrap gap-3 justify-center md:justify-start">
    <button
      onclick={() => dispatch('addScenario')}
      class="px-4 py-2 text-sm font-medium text-white bg-green-600 hover:bg-green-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2"
    >
      Add New
    </button>
    <button
      onclick={() => dispatch('duplicateScenario')}
      disabled={!activeScenarioId}
      class="px-4 py-2 text-sm font-medium text-white bg-blue-500 hover:bg-blue-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50"
    >
      Duplicate Active
    </button>
    <button
      onclick={() => dispatch('deleteScenario')}
      disabled={isDeleteDisabled}
      class="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 disabled:opacity-50"
    >
      Delete Active
    </button>
  </div>
</div>
