<script lang="ts">
  import type { ScenarioInput, LWOPPeriod } from '../../types';
  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  function addLWOPPeriod() {
    if (!inputs.lwopPeriods) {
      inputs.lwopPeriods = [];
    }
    const newLWOP: LWOPPeriod = {
      id: crypto.randomUUID(),
      startDate: '',
      endDate: '',
      type: 'PersonalNonMilitary'
    };
    inputs.lwopPeriods = [...inputs.lwopPeriods, newLWOP];
  }

  function removeLWOPPeriod(idToRemove: string) {
    if (inputs.lwopPeriods) {
      inputs.lwopPeriods = inputs.lwopPeriods.filter((p: LWOPPeriod) => p.id !== idToRemove);
    }
  }
</script>

<div class="pt-6 border-t border-gray-300 mt-6">
  <div class="flex justify-between items-center mb-3">
      <h4 class="text-lg font-semibold text-gray-700">Leave Without Pay (LWOP) Periods</h4>
      <button 
        type="button" 
        onclick={addLWOPPeriod} 
        class="px-4 py-2 text-sm font-medium text-white bg-sky-600 hover:bg-sky-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-1 transition-all duration-150 ease-in-out shadow hover:shadow-md"
      >
        + Add LWOP Period
      </button>
  </div>
  <p class="text-xs text-gray-600 mb-4">Detail periods of LWOP. Certain types and durations can affect creditable service or retirement eligibility (e.g., more than 6 months LWOP in a calendar year is generally not creditable, unless for military duty or OWCP).</p>
  
  {#if inputs.lwopPeriods && inputs.lwopPeriods.length > 0}
    <div class="space-y-4">
    {#each inputs.lwopPeriods as lwopPeriod, i (lwopPeriod.id)}
      <div class="p-4 border border-gray-300 rounded-lg bg-white shadow-sm space-y-3 transition-all hover:shadow-md">
        <div class="flex justify-between items-center border-b border-gray-200 pb-2 mb-3">
          <h5 class="text-base font-medium text-sky-700">LWOP Period {i + 1}</h5>
          <button 
            type="button" 
            onclick={() => removeLWOPPeriod(lwopPeriod.id)} 
            class="text-red-600 hover:text-red-800 text-sm font-semibold p-1 rounded-md hover:bg-red-100 transition-colors duration-150 ease-in-out"
          >
            Remove Period
          </button>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-x-6 gap-y-4">
          <div>
            <label for={`lwop-start-date-${lwopPeriod.id}`} class="block text-sm font-medium text-gray-700 mb-1">Start Date:</label>
            <input type="date" id={`lwop-start-date-${lwopPeriod.id}`} bind:value={lwopPeriod.startDate} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm bg-white text-gray-900">
          </div>
          <div>
            <label for={`lwop-end-date-${lwopPeriod.id}`} class="block text-sm font-medium text-gray-700 mb-1">End Date:</label>
            <input type="date" id={`lwop-end-date-${lwopPeriod.id}`} bind:value={lwopPeriod.endDate} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm bg-white text-gray-900">
          </div>
          <div>
            <label for={`lwop-type-${lwopPeriod.id}`} class="block text-sm font-medium text-gray-700 mb-1">LWOP Type:</label>
            <select id={`lwop-type-${lwopPeriod.id}`} bind:value={lwopPeriod.type} class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-white text-gray-900 border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">
              <option value="PersonalNonMilitary">Personal (Non-Military)</option>
              <option value="MilitaryLWOP">Military LWOP (USERRA)</option>
              <option value="OWCP">OWCP (Workers' Comp)</option>
            </select>
          </div>
        </div>
      </div>
    {/each}
    </div>
  {:else}
     <p class="text-sm text-gray-500 italic">No LWOP periods added yet. Click "+ Add LWOP Period" to begin.</p>
  {/if}
</div>
