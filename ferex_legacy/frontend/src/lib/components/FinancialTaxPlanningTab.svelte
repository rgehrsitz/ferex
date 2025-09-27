<script lang="ts">
  import type { ScenarioInput } from "../../types";
  import { createEventDispatcher, onMount } from "svelte";
  import RecurringIncomeManager from "./RecurringIncomeManager.svelte";
  import OneTimeIncomeEventManager from "./OneTimeIncomeEventManager.svelte";
  import FederalTaxAssumptions from "./FederalTaxAssumptions.svelte";
  import StateTaxAssumptions from "./StateTaxAssumptions.svelte";
  import OtherAssets from "./OtherAssets.svelte";
  import RetirementSpendingGoals from "./RetirementSpendingGoals.svelte";

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
  const dispatch = createEventDispatcher();

  onMount(() => {
    // console.log('[DEBUG] FinancialTaxPlanningTab mounted. Inputs:', JSON.parse(JSON.stringify(inputs)));
  });

  function handleChange() {
    dispatch("change");
  }
</script>

<div class="space-y-8">
  <FederalTaxAssumptions bind:inputs on:change={handleChange} />
  <StateTaxAssumptions bind:inputs on:change={handleChange} />

  <!-- Other Income & Assets Section -->
  <section class="space-y-4 p-4 border border-gray-200 rounded-md shadow-sm">
    <h3 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-3">
      Other Income & Assets
    </h3>
    <div class="space-y-6">
      <!-- Other Recurring Income Sources -->
      <div class="border-t border-gray-200 pt-4">
        <!-- <h4 class="text-md font-medium text-gray-600 mb-2">Other Recurring Income Sources</h4> -->
        <RecurringIncomeManager
          bind:sources={inputs.otherRecurringIncomeSources}
          on:change={handleChange}
        />
      </div>

      <!-- One-Time Income Events -->
      <div class="border-t border-gray-200 pt-4">
        <!-- <h4 class="text-md font-medium text-gray-600 mb-2">One-Time Income Events</h4> -->
        <OneTimeIncomeEventManager
          bind:events={inputs.oneTimeIncomeEvents}
          on:change={handleChange}
        />
      </div>

      <OtherAssets bind:inputs on:change={handleChange} />
    </div>
  </section>

  <RetirementSpendingGoals bind:inputs on:change={handleChange} />
</div>

<!-- Add more sections as needed, e.g., Liabilities, Legacy Goals, Spouse's Info -->
