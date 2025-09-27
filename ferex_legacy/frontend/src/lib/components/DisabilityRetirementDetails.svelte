<script lang="ts">
  import type { ScenarioInput } from '../../types';

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
</script>

<!-- Disability Retirement Section -->
<div class="pt-4 border-t border-gray-200 mt-6">
  <h4 class="text-md font-semibold text-gray-700 mb-3">Special Retirement Provisions</h4>
  <div class="space-y-4">
      <div class="flex items-start">
          <div class="flex items-center h-5">
              <input id="isDisabilityRetirement" name="isDisabilityRetirement" type="checkbox" bind:checked={inputs.isDisabilityRetirement} class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded">
          </div>
          <div class="ml-3 text-sm">
              <label for="isDisabilityRetirement" class="font-medium text-gray-700">Is this a Disability Retirement?</label>
          </div>
      </div>

  </div>
</div>
