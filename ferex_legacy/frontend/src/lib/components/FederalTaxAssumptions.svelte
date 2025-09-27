<script lang="ts">
  import type { ScenarioInput } from "../../types";
  import { createEventDispatcher } from "svelte";

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
  const dispatch = createEventDispatcher();

  function handleChange() {
    // console.log("[DEBUG] FederalTaxAssumptions change:", {
    //   federalTaxFilingStatus: inputs.federalTaxFilingStatus,
    //   federalTaxNumberOfDependents: inputs.federalTaxNumberOfDependents,
    //   stateOfResidenceForTax: inputs.stateOfResidenceForTax,
    // });
    dispatch("change");
  }

  const federalFilingStatusOptions = [
    { value: "Single", label: "Single" },
    { value: "MarriedFilingJointly", label: "Married Filing Jointly" },
    { value: "MarriedFilingSeparately", label: "Married Filing Separately" },
    { value: "HeadOfHousehold", label: "Head of Household" },
    { value: "QualifyingWidow(er)", label: "Qualifying Widow(er)" },
  ];

  const federalTaxLawAssumptionOptions = [
    { value: "CurrentLaw", label: "Current Law (as of today)" },
    { value: "ExtendsTCJA", label: "TCJA Provisions Extended Permanently" },
    {
      value: "Custom",
      label: "Custom (details not modeled, implies user adjustment elsewhere)",
    },
  ];

  // Initialize Federal Tax Assumptions values if they don't exist
  $effect(() => {
    inputs.federalTaxFilingStatus =
      inputs.federalTaxFilingStatus === undefined
        ? "Single"
        : inputs.federalTaxFilingStatus;
    inputs.federalTaxNumberOfDependents =
      inputs.federalTaxNumberOfDependents === undefined
        ? 0
        : inputs.federalTaxNumberOfDependents;
    inputs.federalTaxItemizedDeductions =
      inputs.federalTaxItemizedDeductions === undefined
        ? null
        : inputs.federalTaxItemizedDeductions;
    inputs.federalOtherTaxableIncomeAnnual =
      inputs.federalOtherTaxableIncomeAnnual === undefined
        ? null
        : inputs.federalOtherTaxableIncomeAnnual;
    inputs.federalTaxCreditsAnnual =
      inputs.federalTaxCreditsAnnual === undefined
        ? null
        : inputs.federalTaxCreditsAnnual;
    inputs.federalTaxLawAssumption =
      inputs.federalTaxLawAssumption === undefined
        ? "CurrentLaw"
        : inputs.federalTaxLawAssumption;
  });
</script>

<section class="space-y-4 p-4 border border-gray-200 rounded-md shadow-sm">
  <h3 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-3">
    Federal Tax Assumptions
  </h3>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
    <div>
      <label
        for="federalTaxFilingStatus"
        class="block text-sm font-medium text-gray-700">Filing Status</label
      >
      <select
        id="federalTaxFilingStatus"
        bind:value={inputs.federalTaxFilingStatus}
        onchange={handleChange}
        class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
      >
        {#each federalFilingStatusOptions as option}
          <option value={option.value}>{option.label}</option>
        {/each}
      </select>
    </div>

    <div>
      <label
        for="federalTaxNumberOfDependents"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Number of Dependents</label
      >
      <input
        type="number"
        id="federalTaxNumberOfDependents"
        bind:value={inputs.federalTaxNumberOfDependents}
        onchange={handleChange}
        min="0"
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm"
        placeholder="0"
      />
    </div>

    <div>
      <label
        for="federalTaxItemizedDeductions"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Itemized Deductions Amount (Annual)</label
      >
      <input
        type="number"
        id="federalTaxItemizedDeductions"
        bind:value={inputs.federalTaxItemizedDeductions}
        onchange={handleChange}
        min="0"
        step="any"
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm"
        placeholder="e.g., 15000"
      />
    </div>

    <div>
      <label
        for="federalOtherTaxableIncomeAnnual"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Other Taxable Income (Annual)</label
      >
      <input
        type="number"
        id="federalOtherTaxableIncomeAnnual"
        bind:value={inputs.federalOtherTaxableIncomeAnnual}
        onchange={handleChange}
        step="any"
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm"
        placeholder="e.g., 5000"
      />
    </div>

    <div>
      <label
        for="federalTaxCreditsAnnual"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Tax Credits (Annual)</label
      >
      <input
        type="number"
        id="federalTaxCreditsAnnual"
        bind:value={inputs.federalTaxCreditsAnnual}
        onchange={handleChange}
        min="0"
        step="any"
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sms"
        placeholder="e.g., 2000"
      />
    </div>

    <div>
      <label
        for="federalTaxLawAssumption"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Tax Law Assumption</label
      >
      <select
        id="federalTaxLawAssumption"
        bind:value={inputs.federalTaxLawAssumption}
        onchange={handleChange}
        class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
      >
        {#each federalTaxLawAssumptionOptions as option}
          <option value={option.value}>{option.label}</option>
        {/each}
      </select>
    </div>
  </div>
</section>
