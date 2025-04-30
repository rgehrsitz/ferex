<script lang="ts">
  import Header from './components/Header.svelte';
  import ScenarioSidebar from './components/ScenarioSidebar.svelte';
  import ScenarioTabs from './components/ScenarioTabs.svelte';
  import CompareView from './components/CompareView.svelte';
  import type { Scenario, ScenarioData } from './types/scenario.js';

  function selectScenario(id: number): void { selectedScenarioId = id; compareScenarioId = null; }
  function selectCompare(id: number): void { compareScenarioId = id; }

  let scenarios: Scenario[] = [
    { id: 1, name: "Scenario A", data: {} },
    { id: 2, name: "Scenario B", data: {} }
  ];
  let selectedScenarioId: number | null = 1;
  let compareScenarioId: number | null = null;

  function addScenario(): void {
    const id = Math.max(...scenarios.map(s => s.id)) + 1;
    scenarios = [...scenarios, { id, name: `Scenario ${String.fromCharCode(64 + id)}`, data: {} }];
    selectedScenarioId = id;
    compareScenarioId = null;
  }
  function duplicateScenario(id: number): void {
    const orig = scenarios.find(s => s.id === id);
    if (!orig) return;
    const newId = Math.max(...scenarios.map(s => s.id)) + 1;
    scenarios = [...scenarios, { id: newId, name: orig.name + " Copy", data: { ...orig.data } }];
  }
  function deleteScenario(id: number): void {
    scenarios = scenarios.filter(s => s.id !== id);
    if (selectedScenarioId === id) selectedScenarioId = scenarios[0]?.id || null;
    if (compareScenarioId === id) compareScenarioId = null;
  }

  let selectedScenario: Scenario | undefined;
  let compareScenario: Scenario | undefined;
  $: selectedScenario = scenarios.find(s => s.id === selectedScenarioId);
  $: compareScenario = compareScenarioId ? scenarios.find(s => s.id === compareScenarioId) : undefined;
</script>

<div class="min-h-screen bg-gray-100 text-gray-900 dark:bg-gray-900 dark:text-gray-100" id="main-container">
  <Header />
  <div class="flex">
    <ScenarioSidebar
      {scenarios}
      {selectedScenarioId}
      {compareScenarioId}
      on:add={addScenario}
      on:duplicate={e => duplicateScenario(e.detail)}
      on:delete={e => deleteScenario(e.detail)}
      on:select={e => selectScenario(e.detail)}
      on:compare={e => selectCompare(e.detail)}
    />
    <main class="flex-1 p-6">
      {#if compareScenario}
        <CompareView {selectedScenario} {compareScenario} />
      {:else if selectedScenario}
        <ScenarioTabs bind:scenario={selectedScenario} />
      {:else}
        <div>No scenarios defined.</div>
      {/if}
    </main>
  </div>
</div>
