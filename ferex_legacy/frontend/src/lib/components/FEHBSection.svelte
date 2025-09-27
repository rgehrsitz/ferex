<script lang="ts">
  import type { ScenarioInput } from "../../types";
  import { createEventDispatcher } from "svelte";

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  // Initialize FEHB values if they don't exist
  $effect(() => {
    if (inputs.fehbBiweeklyDeduction === undefined)
      inputs.fehbBiweeklyDeduction = null;
    if (inputs.fehbProjectedPremiumIncreaseRate === undefined)
      inputs.fehbProjectedPremiumIncreaseRate = 5.0;
    if (inputs.fehbContinuedInRetirement === undefined)
      inputs.fehbContinuedInRetirement = true;
  });

  const dispatch = createEventDispatcher();

  function handleFehbChange() {
    dispatch("change");
  }

  // Calculate estimated annual cost for user
  let estimatedAnnualCost = $derived(
    inputs.fehbBiweeklyDeduction ? inputs.fehbBiweeklyDeduction * 26 : 0,
  );
</script>

<div class="space-y-6">
  <!-- Current FEHB Coverage -->
  <div class="bg-gray-50 p-4 rounded-md">
    <h4 class="text-md font-medium text-gray-700 mb-3">
      Current FEHB Coverage
    </h4>
    <p class="text-sm text-gray-600 mb-4">
      Enter information from your Leave and Earnings Statement (LES)
    </p>

    <div class="max-w-md">
      <!-- Biweekly Deduction -->
      <div>
        <label
          for="fehbBiweeklyDeduction"
          class="block text-sm font-medium text-gray-700 mb-1"
        >
          Biweekly Deduction (from LES)
        </label>
        <div class="mt-1 relative rounded-md shadow-sm">
          <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
          >
            <span class="text-gray-500 sm:text-sm">$</span>
          </div>
          <input
            type="number"
            id="fehbBiweeklyDeduction"
            bind:value={inputs.fehbBiweeklyDeduction}
            onchange={handleFehbChange}
            min="0"
            step="0.01"
            class="pl-8 pr-3 py-2 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            placeholder="185.50"
          />
        </div>
        <p class="mt-1 text-xs text-gray-500">
          Look for "FEHB" on your LES deductions
        </p>
      </div>
    </div>

    {#if estimatedAnnualCost > 0}
      <div class="mt-3 p-3 bg-blue-50 rounded-md">
        <p class="text-sm text-blue-800">
          <strong>Estimated Annual Cost:</strong>
          ${estimatedAnnualCost.toLocaleString()}
          <span class="text-blue-600"
            >(${inputs.fehbBiweeklyDeduction} × 26 pay periods)</span
          >
        </p>
      </div>
    {/if}
  </div>

  <!-- Retirement Planning -->
  <div class="space-y-4">
    <h4 class="text-md font-medium text-gray-700">Retirement Planning</h4>

    <!-- Continue into Retirement -->
    <div class="inline-flex items-center w-auto">
      <div class="inline-flex items-center w-auto">
        <input
          id="fehbContinuedInRetirement"
          type="checkbox"
          bind:checked={inputs.fehbContinuedInRetirement}
          onchange={handleFehbChange}
          class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
        />
      </div>
      <div class="ml-3 text-sm">
        <label
          for="fehbContinuedInRetirement"
          class="inline-flex w-auto items-center text-sm font-medium text-gray-700"
        >
          Continue FEHB into Retirement
        </label>
        <p class="text-gray-500">
          {inputs.fehbContinuedInRetirement
            ? "FEHB will be continued into retirement. The government will continue to pay its share of the premium."
            : "You will need to secure other health insurance in retirement."}
        </p>
      </div>
    </div>

    <!-- Premium Increase Rate -->
    <div class="w-48">
      <label
        for="fehbProjectedPremiumIncreaseRate"
        class="block text-sm font-medium text-gray-700 mb-1"
      >
        Projected Annual Premium Increase (%)
      </label>
      <div class="mt-1 relative rounded-md shadow-sm">
        <input
          type="number"
          id="fehbProjectedPremiumIncreaseRate"
          bind:value={inputs.fehbProjectedPremiumIncreaseRate}
          onchange={handleFehbChange}
          min="0"
          max="20"
          step="0.1"
          class="pr-8 py-2 px-2 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
        <div
          class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none"
        >
          <span class="text-gray-500 sm:text-sm">%</span>
        </div>
      </div>
      <p class="mt-1 text-xs text-gray-500">
        Historical average is around 5-7%
      </p>
    </div>
  </div>

  {#if inputs.fehbContinuedInRetirement}
    <div class="bg-blue-50 p-4 rounded-md">
      <h4 class="text-sm font-medium text-blue-800">
        FEHB Retirement Eligibility Requirements
      </h4>
      <ul class="mt-2 text-sm text-blue-700 space-y-1 list-disc list-inside">
        <li>
          Must be enrolled in FEHB for the 5 years immediately before
          retirement, or
        </li>
        <li>
          Since your first opportunity to enroll (if less than 5 years of
          service)
        </li>
        <li>Must be eligible for an immediate annuity</li>
      </ul>
      <p class="mt-2 text-xs text-blue-600">
        <strong>Note:</strong> The government continues to pay its share (typically
        70-75%) of the premium in retirement.
      </p>
    </div>
  {/if}
</div>
