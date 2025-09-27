<script lang="ts">
  import type { Scenario } from '../../types';
  // Props for scenarios and activeScenarioId, Svelte 5 runes
  // Default scenarios to empty array if not provided
let { scenarios = [], activeScenarioId } = $props<{ scenarios?: Scenario[]; activeScenarioId?: string | undefined }>();
  // Event dispatcher to notify parent when scenario selection changes
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher<{ selectScenario: string }>();

  // Find the active scenario object
  // Use fallback to empty array to avoid errors if scenarios is undefined
let activeScenario = $derived((scenarios ?? []).find((s: Scenario) => s.id === activeScenarioId));

  // Helper to format date
  function formatDate(dateStr: string | undefined): string {
    if (!dateStr) return 'N/A';
    const d = new Date(dateStr);
    return d.toLocaleDateString();
  }

  // Handler for clicking a scenario row
  function selectScenario(id: string) {
    dispatch('selectScenario', id);
  }

  import HomePageKeyMetrics from './HomePageKeyMetrics.svelte';
</script>

<div class="flex flex-col gap-6">
  <h1 class="text-2xl font-bold text-blue-900 dark:text-blue-200">Welcome to FeReX</h1>

  <!-- Active Scenario Summary Card -->
  {#if activeScenario}
    <div class="bg-blue-50 dark:bg-blue-900/30 rounded-lg shadow p-6 flex flex-col md:flex-row gap-6 items-center justify-between border border-blue-100 dark:border-blue-800">
      <div>
        <div class="text-sm text-blue-700 dark:text-blue-200 font-semibold">Active Scenario</div>
        <div class="text-2xl font-bold text-blue-900 dark:text-blue-100">{activeScenario.name}</div>
        <div class="text-xs text-gray-500 dark:text-gray-300 mt-1">Last Modified: {formatDate(activeScenario.inputs?.lastModified)}</div>
      </div>
      <div class="flex-1">
        <HomePageKeyMetrics scenario={activeScenario} />
      </div>
    </div>
  {/if}

  <!-- Recent Scenarios Table -->
  <div class="bg-white dark:bg-gray-800 rounded shadow p-4 mt-4">
    <div class="font-semibold mb-2">Recent Scenarios</div>
    <table class="w-full text-sm">
      <thead>
        <tr class="text-left text-gray-500 dark:text-gray-300">
          <th class="py-2">Name</th>
          <th class="py-2">Date Modified</th>
          <th class="py-2">Key Metric</th>
        </tr>
      </thead>
      <tbody class="text-gray-800 dark:text-gray-100">
        {#each scenarios as scenario}
          <tr
            class="border-t border-gray-200 dark:border-gray-700 hover:bg-blue-50 dark:hover:bg-blue-900/30 cursor-pointer {scenario.id === activeScenarioId ? 'bg-blue-100 dark:bg-blue-900/50 font-semibold' : ''}"
            onclick={() => selectScenario(scenario.id)}
          >
            <td class="py-2">{scenario.name}</td>
            <td class="py-2">{formatDate(scenario.inputs?.lastModified)}</td>
            <td class="py-2">N/A</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
