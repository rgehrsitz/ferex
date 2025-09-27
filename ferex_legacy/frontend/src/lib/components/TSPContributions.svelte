<script lang="ts">
  import type { ScenarioInput } from '../../types';

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  const tspContributeUntilOptions = [
    { value: 'Retirement', label: 'Retirement' },
    { value: 'SpecificAge', label: 'Specific Age' },
  ];
</script>

<div class="mb-4">
  <label for="tspBalanceTraditional" class="block text-sm font-medium text-gray-700 mb-1">Current Traditional TSP Balance:</label>
  <input type="number" id="tspBalanceTraditional" bind:value={inputs.tspBalanceTraditional} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white text-gray-700" placeholder="e.g., 500000" step="1000" />
</div>

<div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4 mb-4">
  <div>
    <label for="tspAnnualContributionPreRetirement" class="block text-sm font-medium text-gray-700 mb-1">Annual TSP Contribution:</label>
    <input type="number" id="tspAnnualContributionPreRetirement" bind:value={inputs.tspAnnualContributionPreRetirement} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white text-gray-700" placeholder="Amount" />
  </div>
  <div>
    <label for="tspContributeUntil" class="block text-sm font-medium text-gray-700 mb-1">Contribute Until:</label>
    <select id="tspContributeUntil" bind:value={inputs.tspContributeUntil} class="mt-1 block w-full p-2 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
      <option value={null}>Select Option</option>
      {#each tspContributeUntilOptions as option}
        <option value={option.value}>{option.label}</option>
      {/each}
    </select>
  </div>
</div>

<div class="mb-4">
  <label for="tspExpectedAnnualGrowthRatePreRetirement" class="block text-sm font-medium text-gray-700 mb-1">Expected Annual Growth Rate (Pre-Retirement %):</label>
  <input type="number" id="tspExpectedAnnualGrowthRatePreRetirement" bind:value={inputs.tspExpectedAnnualGrowthRatePreRetirement} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white text-gray-700" placeholder="e.g., 5.0" step="0.1" />
</div>
