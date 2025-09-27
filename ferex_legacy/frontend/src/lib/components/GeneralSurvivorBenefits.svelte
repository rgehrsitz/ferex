<script lang="ts">
  import type { ScenarioInput } from '../../types';

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

    // Validation state for real-time feedback
    let validationErrors = $state<Record<string, string>>({});
    
</script>

<div class="pt-6 border-t border-gray-300 mt-6">
  <h4 class="text-md font-semibold text-gray-700 mb-3">Survivor Benefits - General</h4>
  <div class="space-y-4">
    <div class="flex items-start">
      <div class="flex items-center h-5">
        <input 
          id="isMarriedAtRetirement" 
          type="checkbox" 
          bind:checked={inputs.isMarriedAtRetirement} 
          class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded"
        >
      </div>
      <div class="ml-3 text-sm">
        <label for="isMarriedAtRetirement" class="font-medium text-gray-700">Is Married at Retirement?</label>
        <p class="text-xs text-gray-500">This impacts eligibility for spousal survivor benefits.</p>
      </div>
    </div>
    <div class="flex items-start">
      <div class="flex items-center h-5">
        <input 
          id="hasFormerSpouseEntitlement" 
          type="checkbox" 
          bind:checked={inputs.hasFormerSpouseEntitlement} 
          class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded"
        >
      </div>
      <div class="ml-3 text-sm">
        <label for="hasFormerSpouseEntitlement" class="font-medium text-gray-700">Has Former Spouse Court-Ordered Entitlement?</label>
        <p class="text-xs text-gray-500">Check this if a court order mandates survivor benefits for a former spouse.</p>
      </div>
    </div>

    {#if inputs.isMarriedAtRetirement}
    <div class="grid grid-cols-1 gap-y-3 gap-x-4 sm:grid-cols-6">
      <div class="sm:col-span-3">
        <label for="spouseBirthDate" class="block text-sm font-medium text-gray-700">Spouse's Date of Birth</label>

          <input 
            type="date" 
            name="spouseBirthDate" 
            id="spouseBirthDate" 
            bind:value={inputs.spouseBirthDate}
            class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm {validationErrors.dateSequence ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'}"
          >
          {#if validationErrors.dateSequence}
          <p class="mt-1 text-xs text-red-600">{validationErrors.dateSequence}</p>
        {/if}
        <p class="mt-1 text-xs text-gray-500">Required if married at retirement for accurate survivor benefit calculations.</p>
      </div>
    </div>
    {/if}
  </div>
</div>
