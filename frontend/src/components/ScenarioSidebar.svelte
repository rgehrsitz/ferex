<script lang="ts">
  import type { Scenario } from '../types/scenario.js';
  export let scenarios: Scenario[] = [];
  export let selectedScenarioId: number;
  export let compareScenarioId;
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
  
  // State for editing scenario names
  let editingScenarioId: number | null = null;
  let editName: string = '';

  // Svelte 5 idiom: use let for local state, $: for derived values if needed

  // Custom action to focus input element when it's created
  // This is a better a11y practice than using autofocus attribute
  function focusInput(node: HTMLElement) {
    // Focus the element after a short delay to ensure it's in the DOM
    setTimeout(() => {
      if (node && typeof node.focus === 'function') {
        node.focus();
      }
    }, 10);
    
    return {}; // Return empty object for action API
  }
  
  // Begin editing a scenario name
  function startEditing(scenario: Scenario) {
    editingScenarioId = scenario.id;
    editName = scenario.name;
  }
  
  // Save edited name
  function saveEdit() {
    if (editingScenarioId !== null && editName.trim() !== '') {
      dispatch('rename', { id: editingScenarioId, name: editName.trim() });
      editingScenarioId = null;
    }
  }
  
  // Cancel editing without saving
  function cancelEdit() {
    editingScenarioId = null;
  }
  
  // Handle keyboard events for the edit input
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      saveEdit();
    } else if (event.key === 'Escape') {
      cancelEdit();
    }
  }
</script>

<aside class="w-64 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 p-4 border-r border-gray-200 dark:border-gray-700 min-h-screen overflow-y-auto">
  <h2 class="font-bold mb-4 text-gray-800 dark:text-gray-200">Scenarios</h2>
  <ul class="space-y-2">
    {#each scenarios as scenario}
      <li class="flex flex-col gap-2">
        <div class="flex items-center">
          {#if editingScenarioId === scenario.id}
            <!-- Editing Mode -->
            <div class="flex-1 flex items-center px-1 py-1 border border-primary-500 dark:border-primary-400 rounded-md bg-white dark:bg-gray-800">
              <input 
                type="text" 
                bind:value={editName} 
                class="flex-1 px-2 py-1 text-sm bg-transparent border-none focus:outline-none text-gray-800 dark:text-gray-200"
                on:keydown={handleKeydown}
                use:focusInput
              />
              <button 
                class="text-xs px-1 text-green-600 dark:text-green-400 hover:text-green-700 dark:hover:text-green-300"
                title="Save" 
                on:click={saveEdit}
              >
                ✓
              </button>
              <button 
                class="text-xs px-1 text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300"
                title="Cancel" 
                on:click={cancelEdit}
              >
                ✕
              </button>
            </div>
          {:else}
            <!-- Display Mode -->
            <button 
              class="flex-1 text-left px-3 py-2 rounded-md text-sm font-medium
              {selectedScenarioId === scenario.id 
                ? 'bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200' 
                : 'hover:bg-gray-100 text-gray-800 dark:hover:bg-gray-700 dark:text-gray-200'} 
              border border-gray-200 dark:border-gray-700 transition-colors"
              on:click={() => dispatch('select', scenario.id)}
            >
              {scenario.name}
            </button>
          {/if}
        </div>
        <div class="flex items-center gap-1 px-1">
          <button 
            class="text-xs px-2 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400" 
            title="Rename" 
            on:click={() => startEditing(scenario)}
          >
            <span class="text-sm">✏️</span>
          </button>
          <button 
            class="text-xs px-2 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400" 
            title="Duplicate" 
            on:click={() => dispatch('duplicate', scenario.id)}
          >
            <span class="text-sm">⎘</span>
          </button>
          <button 
            class="text-xs px-2 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400" 
            title="Delete" 
            on:click={() => dispatch('delete', scenario.id)}
          >
            <span class="text-sm">🗑️</span>
          </button>
          {#if compareScenarioId !== scenario.id}
            <button 
              class="text-xs px-2 py-1 rounded bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 ml-auto"
              title="Compare" 
              on:click={() => dispatch('compare', scenario.id)}
            >
              Compare
            </button>
          {/if}
        </div>
      </li>
    {/each}
  </ul>
  <button 
    class="mt-4 w-full px-3 py-2 rounded-md text-sm font-medium bg-primary-600 hover:bg-primary-700 dark:bg-primary-700 dark:hover:bg-primary-600 text-white transition-colors" 
    on:click={() => dispatch('add')}
  >
    + Add Scenario
  </button>
</aside>
