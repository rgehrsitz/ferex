<script lang="ts">
  import type { ScenarioInput } from "../../types";
  import { createEventDispatcher } from "svelte";

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
  const dispatch = createEventDispatcher();

  function handleChange() {
    // console.log('[DEBUG] StateTaxAssumptions change:', {
    //   stateOfResidenceForTax: inputs.stateOfResidenceForTax,
    //   retirementEffectiveStateTaxRate: inputs.retirementEffectiveStateTaxRate
    // });
    dispatch("change");
  }

  // US States and territories list
  const stateOptions = [
    { value: "", label: "Select State/Territory" },
    { value: "Alabama", label: "Alabama (AL)" },
    { value: "Alaska", label: "Alaska (AK)" },
    { value: "Arizona", label: "Arizona (AZ)" },
    { value: "Arkansas", label: "Arkansas (AR)" },
    { value: "California", label: "California (CA)" },
    { value: "Colorado", label: "Colorado (CO)" },
    { value: "Connecticut", label: "Connecticut (CT)" },
    { value: "Delaware", label: "Delaware (DE)" },
    { value: "District of Columbia", label: "District of Columbia (DC)" },
    { value: "Florida", label: "Florida (FL)" },
    { value: "Georgia", label: "Georgia (GA)" },
    { value: "Hawaii", label: "Hawaii (HI)" },
    { value: "Idaho", label: "Idaho (ID)" },
    { value: "Illinois", label: "Illinois (IL)" },
    { value: "Indiana", label: "Indiana (IN)" },
    { value: "Iowa", label: "Iowa (IA)" },
    { value: "Kansas", label: "Kansas (KS)" },
    { value: "Kentucky", label: "Kentucky (KY)" },
    { value: "Louisiana", label: "Louisiana (LA)" },
    { value: "Maine", label: "Maine (ME)" },
    { value: "Maryland", label: "Maryland (MD)" },
    { value: "Massachusetts", label: "Massachusetts (MA)" },
    { value: "Michigan", label: "Michigan (MI)" },
    { value: "Minnesota", label: "Minnesota (MN)" },
    { value: "Mississippi", label: "Mississippi (MS)" },
    { value: "Missouri", label: "Missouri (MO)" },
    { value: "Montana", label: "Montana (MT)" },
    { value: "Nebraska", label: "Nebraska (NE)" },
    { value: "Nevada", label: "Nevada (NV)" },
    { value: "New Hampshire", label: "New Hampshire (NH)" },
    { value: "New Jersey", label: "New Jersey (NJ)" },
    { value: "New Mexico", label: "New Mexico (NM)" },
    { value: "New York", label: "New York (NY)" },
    { value: "North Carolina", label: "North Carolina (NC)" },
    { value: "North Dakota", label: "North Dakota (ND)" },
    { value: "Ohio", label: "Ohio (OH)" },
    { value: "Oklahoma", label: "Oklahoma (OK)" },
    { value: "Oregon", label: "Oregon (OR)" },
    { value: "Pennsylvania", label: "Pennsylvania (PA)" },
    { value: "Rhode Island", label: "Rhode Island (RI)" },
    { value: "South Carolina", label: "South Carolina (SC)" },
    { value: "South Dakota", label: "South Dakota (SD)" },
    { value: "Tennessee", label: "Tennessee (TN)" },
    { value: "Texas", label: "Texas (TX)" },
    { value: "Utah", label: "Utah (UT)" },
    { value: "Vermont", label: "Vermont (VT)" },
    { value: "Virginia", label: "Virginia (VA)" },
    { value: "Washington", label: "Washington (WA)" },
    { value: "West Virginia", label: "West Virginia (WV)" },
    { value: "Wisconsin", label: "Wisconsin (WI)" },
    { value: "Wyoming", label: "Wyoming (WY)" },
  ];

  // Initialize State Tax Assumptions values if they don't exist
  $effect(() => {
    inputs.stateOfResidenceForTax =
      inputs.stateOfResidenceForTax === undefined
        ? ""
        : inputs.stateOfResidenceForTax;
    inputs.retirementEffectiveStateTaxRate =
      inputs.retirementEffectiveStateTaxRate === undefined
        ? null
        : inputs.retirementEffectiveStateTaxRate;
  });

  // Update effective tax rate when state changes (if using automatic calculation)
  function handleStateChange() {
    // This could be enhanced to automatically set the effective rate based on the state's configuration
    // For now, we'll just dispatch the change
    dispatch("change");
  }
</script>

<section class="space-y-4 p-4 border border-gray-200 rounded-md shadow-sm">
  <h3 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-3">
    State Tax Assumptions
  </h3>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
    <div>
      <label
        for="stateOfResidenceForTax"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Retirement State of Residence</label
      >
      <select
        id="stateOfResidenceForTax"
        bind:value={inputs.stateOfResidenceForTax}
        onchange={handleStateChange}
        class="mt-1 block w-full p-2 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
      >
        {#each stateOptions as option}
          <option value={option.value}>{option.label}</option>
        {/each}
      </select>
      <p class="mt-1 text-xs text-gray-500">
        State tax rules will be applied automatically based on federal pension,
        TSP, and Social Security taxation by state
      </p>
    </div>

    <div>
      <label
        for="retirementEffectiveStateTaxRate"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Override Effective State Tax Rate (%)</label
      >
      <input
        type="number"
        id="retirementEffectiveStateTaxRate"
        bind:value={inputs.retirementEffectiveStateTaxRate}
        onchange={handleChange}
        min="0"
        max="100"
        step="any"
        class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        placeholder="Leave blank to use automatic calculation"
      />
      <p class="mt-1 text-xs text-gray-500">
        Optional: Override the automatic state tax calculation with a custom
        effective rate
      </p>
    </div>
  </div>
</section>
