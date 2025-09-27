<script lang="ts">
  import type { ScenarioInput } from "../../types";
  import { createEventDispatcher } from "svelte";

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
  const dispatch = createEventDispatcher();

  function handleChange() {
    dispatch("change");
  }

  // Initialize Other Assets values if they don't exist
  $effect(() => {
    inputs.otherTaxableAccountBalance =
      inputs.otherTaxableAccountBalance === undefined
        ? null
        : inputs.otherTaxableAccountBalance;

    if (
      inputs.otherTaxableAccountAllocation === undefined ||
      inputs.otherTaxableAccountAllocation === null
    ) {
      inputs.otherTaxableAccountAllocation = {
        stocksPercent: 60,
        bondsPercent: 30,
        cashPercent: 10,
      }; // Default values
    } else {
      inputs.otherTaxableAccountAllocation.stocksPercent =
        inputs.otherTaxableAccountAllocation.stocksPercent === undefined
          ? 60
          : inputs.otherTaxableAccountAllocation.stocksPercent;
      inputs.otherTaxableAccountAllocation.bondsPercent =
        inputs.otherTaxableAccountAllocation.bondsPercent === undefined
          ? 30
          : inputs.otherTaxableAccountAllocation.bondsPercent;
      inputs.otherTaxableAccountAllocation.cashPercent =
        inputs.otherTaxableAccountAllocation.cashPercent === undefined
          ? 10
          : inputs.otherTaxableAccountAllocation.cashPercent;
    }

    inputs.otherTaxableAccountExpectedGrowth =
      inputs.otherTaxableAccountExpectedGrowth === undefined
        ? 5.0
        : inputs.otherTaxableAccountExpectedGrowth;
  });
</script>

<div class="border-t border-gray-200 pt-4">
  <h4 class="text-md font-medium text-gray-600 mb-2">
    Taxable Investment Accounts (Non-TSP)
  </h4>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
    <div>
      <label
        for="otherTaxableAccountBalance"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Total Balance</label
      >
      <input
        type="number"
        id="otherTaxableAccountBalance"
        bind:value={inputs.otherTaxableAccountBalance}
        onblur={handleChange}
        onkeydown={(e) => e.key === "Enter" && handleChange()}
        min="0"
        step="any"
        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        placeholder="e.g., 100000"
      />
    </div>
    <div>
      <label
        for="otherTaxableAccountExpectedGrowth"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Expected Annual Return (%)</label
      >
      <input
        type="number"
        id="otherTaxableAccountExpectedGrowth"
        bind:value={inputs.otherTaxableAccountExpectedGrowth}
        onblur={handleChange}
        onkeydown={(e) => e.key === "Enter" && handleChange()}
        min="0"
        step="any"
        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        placeholder="e.g., 5.0"
      />
    </div>
  </div>

  {#if inputs.otherTaxableAccountAllocation}
    <div class="mt-4 pt-2 border-t border-gray-100">
      <p class="text-sm font-medium text-gray-700 mb-1">
        Asset Allocation (%):
      </p>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-x-4 gap-y-2">
        <div>
          <label
            for="taxableStocksPercent"
            class="block text-xs font-medium text-gray-600">Stocks</label
          >
          <input
            type="number"
            id="taxableStocksPercent"
            bind:value={inputs.otherTaxableAccountAllocation.stocksPercent}
            onblur={handleChange}
            onkeydown={(e) => e.key === "Enter" && handleChange()}
            min="0"
            max="100"
            step="1"
            class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm"
            placeholder="e.g., 60"
          />
        </div>
        <div>
          <label
            for="taxableBondsPercent"
            class="block text-xs font-medium text-gray-600">Bonds</label
          >
          <input
            type="number"
            id="taxableBondsPercent"
            bind:value={inputs.otherTaxableAccountAllocation.bondsPercent}
            onblur={handleChange}
            onkeydown={(e) => e.key === "Enter" && handleChange()}
            min="0"
            max="100"
            step="1"
            class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm"
            placeholder="e.g., 30"
          />
        </div>
        <div>
          <label
            for="taxableCashPercent"
            class="block text-xs font-medium text-gray-600">Cash</label
          >
          <input
            type="number"
            id="taxableCashPercent"
            bind:value={inputs.otherTaxableAccountAllocation.cashPercent}
            onblur={handleChange}
            onkeydown={(e) => e.key === "Enter" && handleChange()}
            min="0"
            max="100"
            step="1"
            class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm"
            placeholder="e.g., 10"
          />
        </div>
      </div>
      <p class="text-xs text-gray-500 mt-1">
        Ensure these sum to 100% if applicable.
      </p>
    </div>
  {/if}
</div>
