<script lang="ts">
  import type { ScenarioInput } from '../../types';
  import { createEventDispatcher } from 'svelte';

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
  const dispatch = createEventDispatcher();

  function handleChange() {
    dispatch('change');
  }

  // Initialize TSP Withdrawal values if they don't exist
  $effect(() => {
    if (inputs.tspWithdrawalStrategy === undefined) inputs.tspWithdrawalStrategy = 'None';
    if (inputs.tspWithdrawalPercentageValue === undefined) inputs.tspWithdrawalPercentageValue = null;
    if (inputs.tspWithdrawalFixedAmountValue === undefined) inputs.tspWithdrawalFixedAmountValue = null;
  });

  const tspWithdrawalStrategyOptions = [
    { value: 'None', label: 'None' },
    { value: 'FixedAmountYearly', label: 'Fixed Amount Annually' },
    { value: 'FixedAmountMonthly', label: 'Fixed Amount Monthly' },
    { value: 'PercentageOfBalanceYearly', label: 'Percentage of Balance Annually' },
    { value: 'IRSMinimumRequiredDistribution', label: 'IRS Minimum Required Distribution' },
    { value: 'AsNeededToMeetIncomeGoal', label: 'As Needed to Meet Income Goal' },
  ];
</script>

<div class="mb-4">
  <label for="tspWithdrawalStrategy" class="block text-sm font-medium text-gray-700 mb-1">TSP Withdrawal Strategy Post-Retirement:</label>
  <select id="tspWithdrawalStrategy" bind:value={inputs.tspWithdrawalStrategy} onblur={handleChange} onkeydown={(e) => e.key === 'Enter' && handleChange()} class="mt-1 block w-full p-2 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
    {#each tspWithdrawalStrategyOptions as option}
      <option value={option.value}>{option.label}</option>
    {/each}
  </select>
{#if inputs.tspWithdrawalStrategy === 'PercentageOfBalanceYearly'}
  <div class="mb-4 mt-2">
    <label for="tspWithdrawalPercentageValue" class="block text-sm font-medium text-gray-700 mb-1">Withdrawal Percentage (% of Balance Annually):</label>
    <input 
      type="number" 
      id="tspWithdrawalPercentageValue" 
      value={inputs.tspWithdrawalPercentageValue ? (inputs.tspWithdrawalPercentageValue * 100).toString() : ''}
      oninput={(e) => {
        const target = e.target as HTMLInputElement;
        const percentageValue = parseFloat(target.value);
        if (!isNaN(percentageValue)) {
          // Convert percentage to decimal for storage (e.g., 4% becomes 0.04)
          inputs.tspWithdrawalPercentageValue = percentageValue / 100;
        } else {
          inputs.tspWithdrawalPercentageValue = null;
        }
        handleChange();
      }}
      onblur={handleChange} onkeydown={(e) => e.key === 'Enter' && handleChange()}
      class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white text-gray-700" 
      placeholder="e.g., 4" 
      min="0" 
      max="100" 
      step="0.1" 
    />
    <small class="text-gray-500 text-xs">Enter as percentage (e.g., 4 for 4% annually)</small>
  </div>
{/if}
{#if inputs.tspWithdrawalStrategy === 'FixedAmountYearly' || inputs.tspWithdrawalStrategy === 'FixedAmountMonthly'}
  <div class="mb-4 mt-2">
    <label for="tspWithdrawalFixedAmountValue" class="block text-sm font-medium text-gray-700 mb-1">
      Withdrawal Amount ({inputs.tspWithdrawalStrategy === 'FixedAmountYearly' ? 'Annually' : 'Monthly'}):
    </label>
    <input 
      type="number" 
      id="tspWithdrawalFixedAmountValue" 
      bind:value={inputs.tspWithdrawalFixedAmountValue} 
      onblur={handleChange} onkeydown={(e) => e.key === 'Enter' && handleChange()}
      class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white text-gray-700" 
      placeholder={inputs.tspWithdrawalStrategy === 'FixedAmountYearly' ? 'e.g., 20000' : 'e.g., 1500'} 
      min="0" 
      step={inputs.tspWithdrawalStrategy === 'FixedAmountYearly' ? '100' : '50'}
    />
  </div>
{/if}
</div>
