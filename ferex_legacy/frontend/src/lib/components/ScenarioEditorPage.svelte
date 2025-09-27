<script lang="ts">
  import AnnuityServiceDetailsTab from './AnnuityServiceDetailsTab.svelte';
  import SocialSecurityTSPTab from './SocialSecurityTSPTab.svelte';
  import InsuranceTab from './InsuranceTab.svelte';
  import FinancialTaxPlanningTab from './FinancialTaxPlanningTab.svelte';
  import type { Scenario, ScenarioInput } from '../../types';

  // TabName type for stricter prop typing
  type TabName = 'Annuity & Service Details' | 'Social Security & TSP' | 'Insurance' | 'Financial & Tax Planning';

  // Props for all scenario and tab state/handlers
  let { scenarios = $bindable(), activeScenarioId, setActiveScenarioId, activeTab, setActiveTab, addNewScenario, deleteActiveScenario, duplicateActiveScenario, TABS } = $props<{
    scenarios: Scenario[];
    activeScenarioId: string | undefined;
    setActiveScenarioId: (id: string) => void;
    activeTab: TabName;
    setActiveTab: (tab: TabName) => void;
    addNewScenario: () => void;
    deleteActiveScenario: () => void;
    duplicateActiveScenario: () => void;
    TABS: readonly TabName[];
  }>();

  // Derive active scenario from the bindable scenarios array
  let activeScenario = $derived(scenarios.find((s: Scenario) => s.id === activeScenarioId));

  function handleTabChange(tab: TabName) {
    setActiveTab(tab);
  }

  function handleAdd() {
    addNewScenario();
  }
  function handleDelete() {
    deleteActiveScenario();
  }
  function handleDuplicate() {
    duplicateActiveScenario();
  }

  // Scenario switcher handler
  function handleScenarioSwitch(event: Event) {
    const id = (event.target as HTMLSelectElement).value;
    if (id && scenarios.some((s: Scenario) => s.id === id)) {
      setActiveScenarioId(id);
    }
  }
</script>

<div class="flex flex-col gap-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold text-blue-900 dark:text-blue-200">Scenario Editor</h1>
    
    <!-- Calculation Status Indicator -->
    {#if activeScenario?.calculating}
      <div class="flex items-center gap-2 text-blue-600 dark:text-blue-400">
        <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="text-sm font-medium">Calculating...</span>
      </div>
    {:else if activeScenario?.results}
      <div class="flex items-center gap-2 text-green-600 dark:text-green-400">
        <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
        </svg>
        <span class="text-sm font-medium">Calculated</span>
      </div>
    {:else if activeScenario?.error}
      <div class="flex items-center gap-2 text-red-600 dark:text-red-400">
        <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
        </svg>
        <span class="text-sm font-medium">Error</span>
      </div>
    {/if}
  </div>
  
  <div class="bg-white dark:bg-gray-800 rounded shadow p-6">
    <div class="border-b border-gray-200 dark:border-gray-700 mb-4">
      <nav class="flex gap-4">
        {#each TABS as tab}
          <button
            class="py-2 px-3 font-semibold border-b-2 focus:outline-none {activeTab === tab ? 'text-blue-700 dark:text-blue-200 border-blue-600 dark:border-blue-400' : 'text-gray-500 dark:text-gray-400 border-transparent'}"
            onclick={() => handleTabChange(tab as TabName)}
          >
            {tab}
          </button>
        {/each}
      </nav>
    </div>
    <div>
      {#if activeTab === 'Annuity & Service Details' && activeScenario}
        <AnnuityServiceDetailsTab bind:inputs={activeScenario.inputs} />
      {:else if activeTab === 'Social Security & TSP' && activeScenario}
        <SocialSecurityTSPTab bind:inputs={activeScenario.inputs} />
      {:else if activeTab === 'Insurance' && activeScenario}
        <InsuranceTab bind:inputs={activeScenario.inputs} />
      {:else if activeTab === 'Financial & Tax Planning' && activeScenario}
        <FinancialTaxPlanningTab bind:inputs={activeScenario.inputs} />
      {/if}
    </div>
  </div>
</div>
