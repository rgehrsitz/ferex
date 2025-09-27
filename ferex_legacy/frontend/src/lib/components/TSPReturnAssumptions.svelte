<script lang="ts">
  import type { ScenarioInput, TSPReturnAssumptions } from "../../types";

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  $effect(() => {
    if (inputs.userReturnAssumptionsTSP == null) {
      inputs.userReturnAssumptionsTSP = {
        G: null,
        F: null,
        C: null,
        S: null,
        I: null,
      };
    }
  });

  function handlePercentageChange(
    fund: keyof TSPReturnAssumptions,
    event: Event
  ) {
    const target = event.target as HTMLInputElement;
    const percentageValue = parseFloat(target.value);

    if (!isNaN(percentageValue) && inputs.userReturnAssumptionsTSP) {
      inputs.userReturnAssumptionsTSP[fund] = percentageValue;
    } else if (inputs.userReturnAssumptionsTSP) {
      inputs.userReturnAssumptionsTSP[fund] = null;
    }
  }

  function getDisplayValue(fund: keyof TSPReturnAssumptions): string {
    if (
      !inputs.userReturnAssumptionsTSP ||
      inputs.userReturnAssumptionsTSP[fund] == null
    ) {
      return "";
    }
    return inputs.userReturnAssumptionsTSP[fund]!.toString();
  }
</script>

<div class="mb-6">
  <h4 class="text-md font-medium text-gray-800 mb-3">
    TSP Fund Return Assumptions (Annual %)
  </h4>

  <div
    class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3 mb-4"
  >
    <p class="text-sm text-blue-700 dark:text-blue-300">
      <strong>Optional:</strong> Set expected annual return rates for each TSP fund.
      If left blank, the system will use a default 5% return rate.
    </p>
  </div>

  <div class="grid grid-cols-2 md:grid-cols-5 gap-3">
    <div>
      <label
        for="tspReturnGFund"
        class="block text-sm font-medium text-gray-700 mb-1"
      >
        G Fund (%)
      </label>
      <input
        type="number"
        id="tspReturnGFund"
        bind:value={inputs.userReturnAssumptionsTSP.G}
        min="0"
        max="20"
        step="0.1"
        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
        placeholder="e.g., 2.5"
      />
      <small class="text-gray-500 text-xs">Conservative bonds</small>
    </div>

    <div>
      <label
        for="tspReturnFFund"
        class="block text-sm font-medium text-gray-700 mb-1"
      >
        F Fund (%)
      </label>
      <input
        type="number"
        id="tspReturnFFund"
        bind:value={inputs.userReturnAssumptionsTSP.F}
        min="0"
        max="20"
        step="0.1"
        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
        placeholder="e.g., 3.5"
      />
      <small class="text-gray-500 text-xs">Bond index</small>
    </div>

    <div>
      <label
        for="tspReturnCFund"
        class="block text-sm font-medium text-gray-700 mb-1"
      >
        C Fund (%)
      </label>
      <input
        type="number"
        id="tspReturnCFund"
        bind:value={inputs.userReturnAssumptionsTSP.C}
        min="0"
        max="20"
        step="0.1"
        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
        placeholder="e.g., 7.0"
      />
      <small class="text-gray-500 text-xs">S&P 500</small>
    </div>

    <div>
      <label
        for="tspReturnSFund"
        class="block text-sm font-medium text-gray-700 mb-1"
      >
        S Fund (%)
      </label>
      <input
        type="number"
        id="tspReturnSFund"
        bind:value={inputs.userReturnAssumptionsTSP.S}
        min="0"
        max="20"
        step="0.1"
        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
        placeholder="e.g., 8.0"
      />
      <small class="text-gray-500 text-xs">Small cap stocks</small>
    </div>

    <div>
      <label
        for="tspReturnIFund"
        class="block text-sm font-medium text-gray-700 mb-1"
      >
        I Fund (%)
      </label>
      <input
        type="number"
        id="tspReturnIFund"
        bind:value={inputs.userReturnAssumptionsTSP.I}
        min="0"
        max="20"
        step="0.1"
        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
        placeholder="e.g., 6.5"
      />
      <small class="text-gray-500 text-xs">International stocks</small>
    </div>
  </div>

  <div class="mt-3">
    <small class="text-gray-600">
      <strong>Tip:</strong> Historical long-term averages: G Fund ~2-3%, F Fund ~3-4%,
      C Fund ~7-10%, S Fund ~8-11%, I Fund ~6-8%. Adjust based on your expectations.
    </small>
  </div>
</div>
