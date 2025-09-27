<script lang="ts">
  import type { ScenarioInput } from "../../types";

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  const LFundOptions = [
    "L Income",
    "L 2025",
    "L 2030",
    "L 2035",
    "L 2040",
    "L 2045",
    "L 2050",
    "L 2055",
    "L 2060",
    "L 2065",
  ];

  // Initialize allocations if they don't exist
  $effect(() => {
    if (inputs.tspCurrentAllocationToFunds == null) {
      inputs.tspCurrentAllocationToFunds = {
        G: 0,
        F: 0,
        C: 0,
        S: 0,
        I: 0,
        LFundName: null,
      };
    }
    if (inputs.tspPostRetirementAllocation == null) {
      inputs.tspPostRetirementAllocation = {
        G: 0,
        F: 0,
        C: 0,
        S: 0,
        I: 0,
        LFundName: null,
      };
    }
  });

  const lFundCompositions = {
    "L Income": { G: 71.22, F: 6.51, C: 11.72, S: 2.81, I: 7.74 },
    "L 2025": { G: 69.32, F: 6.64, C: 12.71, S: 3.17, I: 8.16 },
    "L 2030": { G: 42.0, F: 5.0, C: 27.0, S: 9.0, I: 17.0 },
    "L 2035": { G: 34.0, F: 5.0, C: 31.0, S: 11.0, I: 19.0 },
    "L 2040": { G: 26.0, F: 5.0, C: 35.0, S: 13.0, I: 21.0 },
    "L 2045": { G: 18.0, F: 5.0, C: 39.0, S: 15.0, I: 23.0 },
    "L 2050": { G: 10.0, F: 5.0, C: 43.0, S: 17.0, I: 25.0 },
    "L 2055": { G: 6.0, F: 4.0, C: 47.0, S: 19.0, I: 24.0 },
    "L 2060": { G: 6.0, F: 4.0, C: 47.0, S: 19.0, I: 24.0 },
    "L 2065": { G: 6.0, F: 4.0, C: 47.0, S: 19.0, I: 24.0 },
  };

  function handleLFundChange(lFundName: string | null) {
    if (inputs.tspCurrentAllocationToFunds) {
      const newInputs = { ...inputs };
      const newAllocation = {
        ...newInputs.tspCurrentAllocationToFunds,
        LFundName: lFundName,
      };

      if (
        lFundName &&
        lFundCompositions[lFundName as keyof typeof lFundCompositions]
      ) {
        const composition =
          lFundCompositions[lFundName as keyof typeof lFundCompositions];
        newAllocation.G = composition.G;
        newAllocation.F = composition.F;
        newAllocation.C = composition.C;
        newAllocation.S = composition.S;
        newAllocation.I = composition.I;
      }

      newInputs.tspCurrentAllocationToFunds = newAllocation;
      inputs = newInputs;
    }
  }
</script>

{#if inputs.tspCurrentAllocationToFunds}
  <div class="mb-4">
    <p class="block text-sm font-medium text-gray-700 mb-2">
      Current TSP Allocation (%):
    </p>
    <div class="grid grid-cols-3 md:grid-cols-5 gap-2">
      <div>
        <label
          for="tspGFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">G Fund:</label
        >
        <input
          type="number"
          id="tspGFund"
          bind:value={inputs.tspCurrentAllocationToFunds.G}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
      <div>
        <label
          for="tspFFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">F Fund:</label
        >
        <input
          type="number"
          id="tspFFund"
          bind:value={inputs.tspCurrentAllocationToFunds.F}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
      <div>
        <label
          for="tspCFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">C Fund:</label
        >
        <input
          type="number"
          id="tspCFund"
          bind:value={inputs.tspCurrentAllocationToFunds.C}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
      <div>
        <label
          for="tspSFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">S Fund:</label
        >
        <input
          type="number"
          id="tspSFund"
          bind:value={inputs.tspCurrentAllocationToFunds.S}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
      <div>
        <label
          for="tspIFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">I Fund:</label
        >
        <input
          type="number"
          id="tspIFund"
          bind:value={inputs.tspCurrentAllocationToFunds.I}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
    </div>
    <small class="text-gray-500"
      >Ensure G, F, C, S, I percentages sum to 100 if L-Fund is not selected.</small
    >
  </div>
  <div class="mb-4">
    <label for="tspLFund" class="block text-sm font-medium text-gray-700 mb-1"
      >L Fund (Overrides G,F,C,S,I if selected):</label
    >
    <select
      id="tspLFund"
      value={inputs.tspCurrentAllocationToFunds.LFundName}
      onchange={(e) => {
        const target = e.target as HTMLSelectElement;
        const value = target.value === "null" ? null : target.value;
        handleLFundChange(value);
      }}
      class="mt-1 block w-full p-2 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
    >
      <option value={null}>Not Selected</option>
      {#each LFundOptions as lFund}
        <option value={lFund}>{lFund}</option>
      {/each}
    </select>
  </div>
{/if}

<!-- Post-Retirement Allocation Section -->
<div class="mt-6 mb-4">
  <h4 class="text-md font-medium text-gray-800 mb-2">
    Post-Retirement TSP Allocation
  </h4>

  <label
    class="flex items-start gap-2 text-sm text-gray-700 cursor-pointer mb-3"
  >
    <input
      id="usePostRetirementAllocation"
      type="checkbox"
      class="mt-1 h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
      checked={inputs.tspFutureAllocationStrategy ===
        "UsePostRetirementAllocation"}
      onchange={() => {
        inputs.tspFutureAllocationStrategy =
          inputs.tspFutureAllocationStrategy === "UsePostRetirementAllocation"
            ? "MaintainCurrent"
            : "UsePostRetirementAllocation";
      }}
    />
    <span>Use different allocation strategy after retirement</span>
  </label>

  <small class="text-gray-500 mb-3"
    >Define a different asset allocation for your TSP funds after retirement.</small
  >

  {#if inputs.tspFutureAllocationStrategy === "UsePostRetirementAllocation"}
    <div class="grid grid-cols-3 md:grid-cols-5 gap-2">
      <div>
        <label
          for="tspPostRetirementGFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">G Fund:</label
        >
        <input
          type="number"
          id="tspPostRetirementGFund"
          bind:value={inputs.tspPostRetirementAllocation.G}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
      <div>
        <label
          for="tspPostRetirementFFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">F Fund:</label
        >
        <input
          type="number"
          id="tspPostRetirementFFund"
          bind:value={inputs.tspPostRetirementAllocation.F}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
      <div>
        <label
          for="tspPostRetirementCFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">C Fund:</label
        >
        <input
          type="number"
          id="tspPostRetirementCFund"
          bind:value={inputs.tspPostRetirementAllocation.C}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
      <div>
        <label
          for="tspPostRetirementSFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">S Fund:</label
        >
        <input
          type="number"
          id="tspPostRetirementSFund"
          bind:value={inputs.tspPostRetirementAllocation.S}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
      <div>
        <label
          for="tspPostRetirementIFund"
          class="block text-xs font-medium text-gray-600 mb-0.5">I Fund:</label
        >
        <input
          type="number"
          id="tspPostRetirementIFund"
          bind:value={inputs.tspPostRetirementAllocation.I}
          min="0"
          max="100"
          class="mt-1 block w-full p-1.5 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm bg-white text-gray-700"
          placeholder="%"
        />
      </div>
    </div>
  {/if}
</div>
