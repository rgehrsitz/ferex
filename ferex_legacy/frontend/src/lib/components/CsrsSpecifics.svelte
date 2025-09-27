<script lang="ts">
  import type { ScenarioInput, InsurableInterestDetails } from '../../types';

  let { 
    inputs = $bindable(), 
    showCSRSOffsetOptions = false 
  } = $props<{ 
    inputs: ScenarioInput, 
    showCSRSOffsetOptions?: boolean 
  }>();

  function handleCsrsSurvivorElectionChange() {
    if (inputs.survivorBenefitCSRS) {
      if (inputs.survivorBenefitCSRS.election === 'PartialCustomBase') {
        if (inputs.survivorBenefitCSRS.customBaseAmountForPartial === null || inputs.survivorBenefitCSRS.customBaseAmountForPartial === undefined) {
          inputs.survivorBenefitCSRS.customBaseAmountForPartial = 0;
        }
      } else {
        inputs.survivorBenefitCSRS.customBaseAmountForPartial = null;
      }
      if (inputs.survivorBenefitCSRS.election === 'InsurableInterest') {
        if (!inputs.survivorBenefitCSRS.insurableInterestDetails) {
          inputs.survivorBenefitCSRS.insurableInterestDetails = { beneficiaryDOB: '', relationship: '' };
        }
      } else {
        // Only nullify if not 'PartialCustomBase' AND not 'InsurableInterest' to avoid clearing one when switching between them
        if (inputs.survivorBenefitCSRS.election !== 'PartialCustomBase') {
            inputs.survivorBenefitCSRS.insurableInterestDetails = null;
        }
      }
    }
  }
</script>

<div class="pt-6 border-t border-gray-300 mt-6">
  <h4 class="text-md font-semibold text-gray-700 mb-4">CSRS Specific Details</h4>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
    <div class="flex items-start col-span-1 md:col-span-2">
        <div class="flex items-center h-5">
            <input id="isCSRSOffset" type="checkbox" bind:checked={inputs.isCSRSOffset} class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded">
        </div>
        <div class="ml-3 text-sm">
            <label for="isCSRSOffset" class="font-medium text-gray-700">Is CSRS Offset?</label>
        </div>
    </div>
    {#if showCSRSOffsetOptions}
      <div>
        <label for="yearsOfOffsetService" class="block text-sm font-medium text-gray-700">Years of CSRS Offset Service:</label>
        <input type="number" id="yearsOfOffsetService" bind:value={inputs.yearsOfOffsetService} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="e.g., 10">
      </div>
      <div>
        <label for="ssAt62WithOffset" class="block text-sm font-medium text-gray-700">SS at 62 (with Offset rules):</label>
        <input type="number" id="ssAt62WithOffset" bind:value={inputs.ssAt62WithOffset} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="e.g., 1200">
      </div>
      <div>
        <label for="ssAt62WithoutOffset" class="block text-sm font-medium text-gray-700">SS at 62 (standard/no Offset):</label>
        <input type="number" id="ssAt62WithoutOffset" bind:value={inputs.ssAt62WithoutOffset} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="e.g., 1500">
      </div>
    {/if}
    <div>
      <label for="survivorBenefitElectionCSRS" class="block text-sm font-medium text-gray-700">CSRS Survivor Benefit (Spouse):</label>
      <select id="survivorBenefitElectionCSRS" bind:value={inputs.survivorBenefitCSRS.election} onchange={handleCsrsSurvivorElectionChange} class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">
        <option value="None">None</option>
        <option value="Full55PercentMax">Full Spouse (Max 55%)</option>
        <option value="PartialCustomBase">Partial Spouse (Custom Base)</option>
        <option value="InsurableInterest">Insurable Interest</option>
      </select>
        {#if inputs.survivorBenefitCSRS && inputs.survivorBenefitCSRS.election === 'InsurableInterest' && inputs.survivorBenefitCSRS.insurableInterestDetails}
          <div class="mt-3 pl-4 border-l-2 border-teal-500">
            <h5 class="text-sm font-medium text-teal-700 mb-2">Insurable Interest Details:</h5>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
              <div>
                <label for="csrsInsurableInterestDOB" class="block text-xs font-medium text-gray-600">Beneficiary's Date of Birth:</label>
                <input type="date" id="csrsInsurableInterestDOB" bind:value={inputs.survivorBenefitCSRS.insurableInterestDetails.beneficiaryDOB} class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm">
              </div>
              <div>
                <label for="csrsInsurableInterestRelationship" class="block text-xs font-medium text-gray-600">Relationship to Beneficiary:</label>
                <input type="text" id="csrsInsurableInterestRelationship" bind:value={inputs.survivorBenefitCSRS.insurableInterestDetails.relationship} class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm" placeholder="e.g., Sibling, Child">
              </div>
            </div>
          </div>
        {/if}
    </div>
    {#if inputs.survivorBenefitCSRS && inputs.survivorBenefitCSRS.election === 'PartialCustomBase'}
      <div>
        <label for="csrSurvivorBenefitBaseAmount" class="block text-sm font-medium text-gray-700">CSRS Survivor Custom Annuity Base ($):</label>
        <input type="number" id="csrSurvivorBenefitBaseAmount" bind:value={inputs.survivorBenefitCSRS.customBaseAmountForPartial} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="e.g., 20000">
      </div>
    {/if}
    {#if inputs.hasFormerSpouseEntitlement}
      <div class="col-span-1 md:col-span-2 mt-4 pt-4 border-t border-dashed border-gray-400">
        <h5 class="text-sm font-semibold text-gray-600 mb-2">Former Spouse Survivor Benefit (CSRS)</h5>
        <div>
          <label for="survivorBenefitCSRSFormerSpouse" class="block text-sm font-medium text-gray-700">Election for Former Spouse:</label>
          <select 
            id="survivorBenefitCSRSFormerSpouse" 
            bind:value={inputs.survivorBenefitCSRS.formerSpouseElection} 
            class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
            disabled={!inputs.survivorBenefitCSRS}
          >
            <option value="None">None</option>
            <option value="Full55PercentMax">Full (Max 55% of base, subject to court order)</option>
            <option value="PartialCustomBase">Partial (Custom Base, subject to court order)</option>
          </select>
          {#if inputs.survivorBenefitCSRS && inputs.survivorBenefitCSRS.formerSpouseElection === 'PartialCustomBase'}
            <div class="mt-3">
              <label for="csrFormerSpouseSurvivorBenefitBaseAmount" class="block text-sm font-medium text-gray-700">Former Spouse Custom Annuity Base ($):</label>
              <input 
                type="number" 
                id="csrFormerSpouseSurvivorBenefitBaseAmount" 
                bind:value={inputs.survivorBenefitCSRS.formerSpouseCustomBaseAmount} 
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" 
                placeholder="e.g., 15000"
                disabled={!inputs.survivorBenefitCSRS}
              >
            </div>
          {/if}
          {#if inputs.survivorBenefitCSRS && inputs.survivorBenefitCSRS.formerSpouseElection && inputs.survivorBenefitCSRS.formerSpouseElection !== 'None'}
            <div class="mt-3 text-xs text-gray-500">
              <p>Note: Ensure this election complies with the court order. Current spouse consent may be required if total spouse/former spouse benefits exceed maximums.</p>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>
