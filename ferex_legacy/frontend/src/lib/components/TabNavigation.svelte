<script lang="ts">
  // No specific imports needed for Svelte 5 runes for props and events here
  // if types are simple. For complex event payloads, createEventDispatcher might be used.

  // Props definition using Svelte 5 runes
  let { TABS, activeTab } = $props<{
    TABS: readonly string[]; // Using 'readonly string[]' for better type safety if TABS is constant
    activeTab: string;
  }>();

  // Svelte 5 simplifies event dispatching for basic cases, but for typed events:
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher<{
    selectTab: string; // payload is tabName
  }>();

</script>

<div class="mb-6 border-b border-gray-300">
  <nav class="-mb-px flex space-x-2 md:space-x-4 overflow-x-auto" aria-label="Tabs">
    {#each TABS as tabName (tabName)}
      <button
        onclick={() => dispatch('selectTab', tabName)}
        class="whitespace-nowrap py-3 px-2 md:px-4 border-b-2 font-medium text-sm transition-colors duration-150 ease-in-out
               {activeTab === tabName 
                 ? 'border-blue-600 text-blue-700' 
                 : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
        aria-current={activeTab === tabName ? 'page' : undefined}
      >
        {tabName}
      </button>
    {/each}
  </nav>
</div>
