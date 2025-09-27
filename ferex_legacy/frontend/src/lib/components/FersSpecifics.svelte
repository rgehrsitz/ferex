<script lang="ts">
  import type { ScenarioInput, InsurableInterestDetails } from "../../types";

  let { inputs = $bindable(), showSRSInput = false } = $props<{
    inputs: ScenarioInput;
    showSRSInput?: boolean;
  }>();

  // Reactive effect to set default for estimatedSSBenefitAt62ForSRS
  $effect(() => {
    if (showSRSInput && inputs.estimatedSSBenefitAt62ForSRS === null) {
      // console.log('FRP DEBUG (FersSpecifics): Defaulting estimatedSSBenefitAt62ForSRS to 0');
      inputs.estimatedSSBenefitAt62ForSRS = 0;
    }
  });

  function handleFersSurvivorSpouseElectionChange() {
    if (inputs.survivorBenefitFERS) {
      if (inputs.survivorBenefitFERS.spouseElection === "InsurableInterest") {
        if (!inputs.survivorBenefitFERS.insurableInterestDetails) {
          inputs.survivorBenefitFERS.insurableInterestDetails = {
            beneficiaryDOB: "",
            relationship: "",
          };
        }
      } else {
        inputs.survivorBenefitFERS.insurableInterestDetails = null;
      }
    }
  }
</script>

<div class="pt-4 border-t border-gray-200 mt-6">
  <h4 class="text-md font-semibold text-gray-700 mb-4">
    FERS Specific Details
  </h4>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
    <div>
      <label
        for="fersCoverageType"
        class="block text-sm font-medium text-gray-700"
        >FERS Coverage Type:</label
      >
      <select
        id="fersCoverageType"
        bind:value={inputs.fersCoverageType}
        class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
      >
        <option value={null}>-- Select FERS Coverage --</option>
        <option value="FERS">FERS (Hired before 2013)</option>
        <option value="FERS_RAE">FERS-RAE (Hired 2013)</option>
        <option value="FERS_FRAE">FERS-FRAE (Hired 2014+)</option>
      </select>
    </div>

    <div class="col-span-1 md:col-span-2 pt-2">
      <div class="flex items-start">
        <div class="flex items-center h-5">
          <input
            id="isSubjectToSRSEarningsTest"
            type="checkbox"
            bind:checked={inputs.isSubjectToSRSEarningsTest}
            class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
            disabled={!showSRSInput}
          />
        </div>
        <div class="ml-3 text-sm">
          <label
            for="isSubjectToSRSEarningsTest"
            class="font-medium text-gray-700"
            >Subject to FERS Annuity Supplement Earnings Test?</label
          >
          <p
            class="text-xs text-gray-500"
            class:text-gray-400={!showSRSInput}
            class:italic={!showSRSInput}
          >
            Only applicable if retiring before Minimum Retirement Age (MRA) and
            receiving the FERS Supplement. The earnings test does not apply once
            MRA is reached. Test is waived if retiring MRA+10/20/30.
            {showSRSInput ? "" : "(SRS not applicable or age >= 62)"}
          </p>
        </div>
      </div>
    </div>

    {#if showSRSInput}
      <div class="col-span-1 md:col-span-2">
        <label
          for="estimatedSSBenefitAt62ForSRS"
          class="block text-sm font-medium text-gray-700"
          >Estimated Social Security at Age 62 (for FERS Supplement/SRS
          calculation):</label
        >
        <input
          type="number"
          id="estimatedSSBenefitAt62ForSRS"
          bind:value={inputs.estimatedSSBenefitAt62ForSRS}
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          placeholder="e.g., 1500"
        />
        <p class="mt-1 text-xs text-gray-500">
          This value is needed for the FERS Annuity Supplement calculation
          because your planned retirement age is under 62. Your general Social
          Security claiming details are on the 'Social Security & TSP' tab.
        </p>
      </div>
    {/if}

    <div class="flex items-start col-span-1 md:col-span-2 mt-2">
      <div class="flex items-center h-5">
        <input
          id="didSwitchFromCSRS"
          type="checkbox"
          bind:checked={inputs.didSwitchFromCSRS}
          class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
        />
      </div>
      <div class="ml-3 text-sm">
        <label for="didSwitchFromCSRS" class="font-medium text-gray-700"
          >Did you switch from CSRS to FERS?</label
        >
      </div>
    </div>

    {#if inputs.didSwitchFromCSRS}
      <div class="col-span-1 md:col-span-2">
        <label
          for="switchedToFERSDate"
          class="block text-sm font-medium text-gray-700"
          >Date Switched to FERS:</label
        >
        <input
          type="date"
          id="switchedToFERSDate"
          bind:value={inputs.switchedToFERSDate}
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        />
      </div>
    {/if}

    <div class="col-span-1 md:col-span-2">
      <label
        for="survivorBenefitElectionFERS"
        class="block text-sm font-medium text-gray-700"
        >FERS Survivor Benefit (Spouse):</label
      >
      <select
        id="survivorBenefitElectionFERS"
        bind:value={inputs.survivorBenefitFERS.spouseElection}
        onchange={handleFersSurvivorSpouseElectionChange}
        class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
      >
        <option value="None">None</option>
        <option value="Full50Percent">Full Spouse (50%)</option>
        <option value="Partial25Percent">Partial Spouse (25%)</option>
        <option value="InsurableInterest">Insurable Interest</option>
      </select>
      {#if inputs.survivorBenefitFERS && inputs.survivorBenefitFERS.spouseElection === "InsurableInterest" && inputs.survivorBenefitFERS.insurableInterestDetails}
        <div class="mt-3 pl-4 border-l-2 border-sky-500">
          <h5 class="text-sm font-medium text-sky-700 mb-2">
            Insurable Interest Details:
          </h5>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
            <div>
              <label
                for="fersInsurableInterestDOB"
                class="block text-xs font-medium text-gray-600"
                >Beneficiary's Date of Birth:</label
              >
              <input
                type="date"
                id="fersInsurableInterestDOB"
                bind:value={
                  inputs.survivorBenefitFERS.insurableInterestDetails
                    .beneficiaryDOB
                }
                class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-sky-500 focus:border-sky-500 sm:text-sm"
              />
            </div>
            <div>
              <label
                for="fersInsurableInterestRelationship"
                class="block text-xs font-medium text-gray-600"
                >Relationship to Beneficiary:</label
              >
              <input
                type="text"
                id="fersInsurableInterestRelationship"
                bind:value={
                  inputs.survivorBenefitFERS.insurableInterestDetails
                    .relationship
                }
                class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-sky-500 focus:border-sky-500 sm:text-sm"
                placeholder="e.g., Sibling, Child"
              />
            </div>
          </div>
        </div>
      {/if}
    </div>

    {#if inputs.hasFormerSpouseEntitlement}
      <div
        class="col-span-1 md:col-span-2 mt-4 pt-4 border-t border-dashed border-gray-400"
      >
        <h5 class="text-sm font-semibold text-gray-600 mb-2">
          Former Spouse Survivor Benefit (FERS)
        </h5>
        <div>
          <label
            for="survivorBenefitFERSFormerSpouse"
            class="block text-sm font-medium text-gray-700"
            >Election for Former Spouse:</label
          >
          <select
            id="survivorBenefitFERSFormerSpouse"
            bind:value={inputs.survivorBenefitFERS.formerSpouseElection}
            class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
            disabled={!inputs.survivorBenefitFERS}
          >
            <option value="None">None</option>
            <option value="Full50Percent">Full (50% of your annuity)</option>
            <option value="Partial25Percent"
              >Partial (25% of your annuity)</option
            >
          </select>
          {#if inputs.survivorBenefitFERS && inputs.survivorBenefitFERS.formerSpouseElection && inputs.survivorBenefitFERS.formerSpouseElection !== "None"}
            <div class="mt-3 text-xs text-gray-500">
              <p>
                Note: Ensure this election complies with the court order.
                Current spouse consent may be required if total spouse/former
                spouse benefits exceed maximums.
              </p>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>
