<script lang="ts">
  import type { ScenarioInput, ServicePeriod } from '../../types';
  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  function addServicePeriod() {
    if (!inputs.servicePeriods) {
      inputs.servicePeriods = [];
    }
    const newPeriod: ServicePeriod = {
      id: crypto.randomUUID(),
      startDate: '',
      endDate: '',
      serviceCategory: 'Civilian', // Default category
      civilianServiceType: 'RegularDeductionFERS', // Default for Civilian
      militaryServiceType: null,
      depositRedepositPaymentStatus: 'NotApplicable',
      systemDuringService: inputs.calculationSystem || 'None',
      isPartTime: false,
      hoursPerWeekIfPartTime: null,
      notes: null
    };
    inputs.servicePeriods = [...inputs.servicePeriods, newPeriod];
  }

  function removeServicePeriod(idToRemove: string) {
    if (inputs.servicePeriods) {
      inputs.servicePeriods = inputs.servicePeriods.filter((p: ServicePeriod) => p.id !== idToRemove);
    }
  }
</script>

<div class="pt-6 border-t border-gray-300 mt-6">
    <h4 class="text-md font-semibold text-gray-700 mb-3">Creditable Service Periods</h4>
    <p class="text-xs text-gray-600 mb-3">Detail all distinct periods of federal civilian and military service. These will be used by the backend to calculate total creditable service.</p>
    {#if inputs.servicePeriods && inputs.servicePeriods.length > 0}
      {#each inputs.servicePeriods as period, i (period.id)}
        <div class="p-3 mb-3 border border-gray-200 rounded-lg bg-white shadow-sm space-y-2 transition-all hover:shadow-md">
          <div class="flex justify-between items-center border-b border-gray-100 pb-1.5 mb-2">
            <h5 class="text-sm font-medium text-gray-600">Service Period {i + 1}</h5>
            <button 
              type="button" 
              onclick={() => removeServicePeriod(period.id)} 
              class="text-red-500 hover:text-red-700 text-xs font-semibold p-1 rounded-md hover:bg-red-50 transition-colors duration-150 ease-in-out"
            >
              Remove
            </button>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-x-4 gap-y-3 text-sm">
            <div>
              <label for={`period-start-date-${period.id}`} class="block text-xs font-medium text-gray-600">Start Date:</label>
              <input type="date" id={`period-start-date-${period.id}`} bind:value={period.startDate} class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-xs">
            </div>
            <div>
              <label for={`period-end-date-${period.id}`} class="block text-xs font-medium text-gray-600">End Date:</label>
              <input type="date" id={`period-end-date-${period.id}`} bind:value={period.endDate} class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-xs">
            </div>
            <div>
              <label for={`period-service-category-${period.id}`} class="block text-xs font-medium text-gray-600">Service Category:</label>
              <select id={`period-service-category-${period.id}`} bind:value={period.serviceCategory} class="mt-1 block w-full pl-2 pr-8 py-1.5 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">
                <option value="Civilian">Civilian</option>
                <option value="Military">Military</option>
              </select>
            </div>

            {#if period.serviceCategory === 'Civilian'}
              <div>
                <label for={`period-civilian-service-type-${period.id}`} class="block text-xs font-medium text-gray-600">Civilian Service Type:</label>
                <select id={`period-civilian-service-type-${period.id}`} bind:value={period.civilianServiceType} class="mt-1 block w-full pl-2 pr-8 py-1.5 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">
                  <option value="RegularDeductionFERS">Regular Deduction (FERS)</option>
                  <option value="NonDeductionPre10_82CSRS">Non-Deduction Pre-10/82 (CSRS)</option>
                  <option value="RefundedServiceRedepositOwedFERS">Refunded Service, Redeposit Owed (FERS)</option>
                  <option value="TemporaryServiceDepositPaidFERS">Temporary Service, Deposit Paid (FERS)</option>
                  <!-- Add other civilian types as needed -->
                </select>
              </div>
            {/if}

            {#if period.serviceCategory === 'Military'}
              <div>
                <label for={`period-military-service-type-${period.id}`} class="block text-xs font-medium text-gray-600">Military Service Type:</label>
                <select id={`period-military-service-type-${period.id}`} bind:value={period.militaryServiceType} class="mt-1 block w-full pl-2 pr-8 py-1.5 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">
                  <option value="Post1956MilitaryServiceDepositPaid">Post-1956 Military Service, Deposit Paid</option>
                  <option value="Post1956MilitaryServiceDepositOwed">Post-1956 Military Service, Deposit Owed</option>
                  <!-- Add other military types as needed -->
                </select>
              </div>
            {/if}

            <div>
              <label for={`period-deposit-status-${period.id}`} class="block text-xs font-medium text-gray-600">Deposit/Redeposit Status:</label>
              <select id={`period-deposit-status-${period.id}`} bind:value={period.depositRedepositPaymentStatus} class="mt-1 block w-full pl-2 pr-8 py-1.5 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">
                <option value="NotApplicable">Not Applicable</option>
                <option value="PaidInFull">Paid In Full</option>
                <option value="OwedOrPartiallyPaid">Owed or Partially Paid</option>
                <option value="AwaitingDetermination">Awaiting Determination</option>
              </select>
            </div>

            <div class="flex items-center col-span-1 md:col-span-2 lg:col-span-1">
              <input type="checkbox" id={`period-is-part-time-${period.id}`} bind:checked={period.isPartTime} class="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500">
              <label for={`period-is-part-time-${period.id}`} class="ml-2 block text-xs font-medium text-gray-600">Is Part-Time?</label>
            </div>

            {#if period.isPartTime}
              <div>
                <label for={`period-hours-per-week-${period.id}`} class="block text-xs font-medium text-gray-600">Hours/Week (if Part-Time):</label>
                <input type="number" id={`period-hours-per-week-${period.id}`} bind:value={period.hoursPerWeekIfPartTime} class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-xs" placeholder="e.g. 20">
              </div>
            {/if}
             <div>
              <label for={`period-system-${period.id}`} class="block text-xs font-medium text-gray-600">System During Service:</label>
              <select id={`period-system-${period.id}`} bind:value={period.systemDuringService} class="mt-1 block w-full pl-2 pr-8 py-1.5 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">
                <option value="None">None (e.g., Military, Non-Gov)</option>
                <option value="FERS">FERS</option>
                <option value="CSRS">CSRS</option>
              </select>
            </div>
            <div class="col-span-1 md:col-span-2 lg:col-span-4">
              <label for={`period-notes-${period.id}`} class="block text-xs font-medium text-gray-600">Notes (Optional):</label>
              <input type="text" id={`period-notes-${period.id}`} bind:value={period.notes} class="mt-1 block w-full px-2 py-1.5 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-xs" placeholder="e.g., Details about this service period">
            </div>
          </div>
        </div>
      {/each}
    {:else}
      <p class="text-sm text-gray-500 italic">No service periods added yet.</p>
    {/if}
    <button 
      type="button" 
      onclick={addServicePeriod} 
      class="mt-4 px-3 py-1.5 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
    >
      Add Service Period
    </button>
</div>
