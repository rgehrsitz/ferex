<script lang="ts">
  import type { ScenarioInput } from '../../types';
  import { createEventDispatcher } from 'svelte';

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
  
  // Initialize FEGLI values if they don't exist
  $effect(() => {
    if (inputs.fegliBasicCoverage === undefined) inputs.fegliBasicCoverage = true;
    if (inputs.fegliOptionAAmount === undefined) inputs.fegliOptionAAmount = 10000;
    if (inputs.fegliOptionBCoverageMultiples === undefined) inputs.fegliOptionBCoverageMultiples = 1;
    if (inputs.fegliOptionCCoverageMultiples === undefined) inputs.fegliOptionCCoverageMultiples = 0;
    if (inputs.fegliPost65Reduction === undefined) inputs.fegliPost65Reduction = 'FullReduction';
    if (inputs.fegliBasicHeldIntoRetirement === undefined) inputs.fegliBasicHeldIntoRetirement = true;
    if (inputs.fegliBasicPostRetirementReductionChoice === undefined) inputs.fegliBasicPostRetirementReductionChoice = '75Reduce';
  })

  const dispatch = createEventDispatcher();

  function handleFegliChange() {
    dispatch('change');
  }
</script>

<div class="space-y-4">
  <!-- Basic Insurance -->
  <div class="space-y-2">
    <div class="inline-flex items-center w-auto">
      <input
        type="checkbox"
        id="fegliBasicCoverage"
        bind:checked={inputs.fegliBasicCoverage}
        onchange={handleFegliChange}
        class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
      />
      <label for="fegliBasicCoverage" class="ml-2 inline-flex w-auto items-center text-sm font-medium text-gray-700">
        Basic Insurance (1x salary, rounded up to next $1,000 + $2,000)
      </label>
    </div>

    {#if inputs.fegliBasicCoverage}
      <div class="ml-6 space-y-2">
        <div class="inline-flex items-center w-auto">
          <input
            type="radio"
            id="fegliFullReduction"
            name="fegliPost65Reduction"
            value="FullReduction"
            bind:group={inputs.fegliPost65Reduction}
            onchange={handleFegliChange}
            class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
          />
          <label for="fegliFullReduction" class="ml-2 inline-flex w-auto items-center text-sm text-gray-700">
            Full Reduction at age 65 (coverage reduces 2% per month until 25% of original amount remains)
          </label>
        </div>
        <div class="inline-flex items-center w-auto">
          <input
            type="radio"
            id="fegliNoReduction"
            name="fegliPost65Reduction"
            value="NoReduction"
            bind:group={inputs.fegliPost65Reduction}
            onchange={handleFegliChange}
            class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
          />
          <label for="fegliNoReduction" class="ml-2 inline-flex w-auto items-center text-sm text-gray-700">
            No Reduction (full premium required)
          </label>
        </div>
        <div class="inline-flex items-center w-auto">
          <input
            type="radio"
            id="fegliPartialReduction75"
            name="fegliPost65Reduction"
            value="PartialReduction75"
            bind:group={inputs.fegliPost65Reduction}
            onchange={handleFegliChange}
            class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
          />
          <label for="fegliPartialReduction75" class="ml-2 inline-flex w-auto items-center text-sm text-gray-700">
            Partial Reduction at age 65 (coverage reduces 1% per month until 50% of original amount remains)
          </label>
        </div>

        <div class="mt-2">
          <div class="inline-flex items-center w-auto">
            <input
              type="checkbox"
              id="fegliBasicHeldIntoRetirement"
              bind:checked={inputs.fegliBasicHeldIntoRetirement}
              onchange={handleFegliChange}
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label for="fegliBasicHeldIntoRetirement" class="ml-2 inline-flex w-auto items-center text-sm font-medium text-gray-700">
              Continue Basic Insurance into Retirement
            </label>
          </div>

          {#if inputs.fegliBasicHeldIntoRetirement}
            <div class="ml-6 mt-2 space-y-2">
              <p class="text-sm text-gray-500">Post-Retirement Reduction:</p>
              <div class="space-y-1">
                <div class="inline-flex items-center w-auto">
                  <input
                    type="radio"
                    id="reduce75"
                    name="fegliBasicPostRetirementReduction"
                    value="75Reduce"
                    bind:group={inputs.fegliBasicPostRetirementReductionChoice}
                    onchange={handleFegliChange}
                    class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
                  />
                  <label for="reduce75" class="ml-2 inline-flex w-auto items-center text-sm text-gray-700">
                    75% Reduction (no cost, coverage reduces 2% per month)
                  </label>
                </div>
                <div class="inline-flex items-center w-auto">
                  <input
                    type="radio"
                    id="reduce50"
                    name="fegliBasicPostRetirementReduction"
                    value="50Reduce"
                    bind:group={inputs.fegliBasicPostRetirementReductionChoice}
                    onchange={handleFegliChange}
                    class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
                  />
                  <label for="reduce50" class="ml-2 inline-flex w-auto items-center text-sm text-gray-700">
                    50% Reduction (pay full premium, coverage reduces 1% per month)
                  </label>
                </div>
                <div class="inline-flex items-center w-auto">
                  <input
                    type="radio"
                    id="noReduce"
                    name="fegliBasicPostRetirementReduction"
                    value="NoReduce"
                    bind:group={inputs.fegliBasicPostRetirementReductionChoice}
                    onchange={handleFegliChange}
                    class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
                  />
                  <label for="noReduce" class="ml-2 inline-flex w-auto items-center text-sm text-gray-700">
                    No Reduction (pay full premium + 2% of face value)
                  </label>
                </div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>

  <!-- Option A - Standard Insurance -->
  <div class="space-y-2">
    <div class="inline-flex items-center w-auto">
      <input
        type="checkbox"
        id="fegliOptionA"
        checked={!!inputs.fegliOptionAAmount && inputs.fegliOptionAAmount > 0}
        onchange={(e) => {
          const target = e.target as HTMLInputElement;
          inputs.fegliOptionAAmount = target.checked ? 10000 : 0;
          handleFegliChange();
        }}
        class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
      />
      <label for="fegliOptionA" class="ml-2 inline-flex w-auto items-center text-sm font-medium text-gray-700">
        Option A - Standard Insurance
      </label>
    </div>

    {#if inputs.fegliOptionAAmount && inputs.fegliOptionAAmount > 0}
      <div class="ml-6">
        <div class="inline-flex items-center w-auto">
          <span class="text-sm text-gray-500 mr-2">Coverage Amount:</span>
          <div class="relative rounded-md shadow-sm w-32">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <span class="text-gray-500 sm:text-sm">$</span>
            </div>
            <input
              type="number"
              bind:value={inputs.fegliOptionAAmount}
              onchange={handleFegliChange}
              min="0"
              step="1000"
              class="pl-8 pr-3 py-2 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              placeholder="10,000"
            />
          </div>
        </div>
      </div>
    {/if}
  </div>

  <!-- Option B - Additional Insurance -->
  <div class="space-y-2">
    <div class="inline-flex items-center w-auto">
      <input
        type="checkbox"
        id="fegliOptionB"
        checked={!!inputs.fegliOptionBCoverageMultiples && inputs.fegliOptionBCoverageMultiples > 0}
        onchange={(e) => {
          const target = e.target as HTMLInputElement;
          inputs.fegliOptionBCoverageMultiples = target.checked ? 1 : 0;
          handleFegliChange();
        }}
        class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
      />
      <label for="fegliOptionB" class="ml-2 inline-flex w-auto items-center text-sm font-medium text-gray-700">
        Option B - Additional Insurance (multiples of salary)
      </label>
    </div>

    {#if inputs.fegliOptionBCoverageMultiples && inputs.fegliOptionBCoverageMultiples > 0}
      <div class="ml-6">
        <div class="inline-flex items-center w-auto">
          <span class="text-sm text-gray-500 mr-2">Number of multiples:</span>
          <select
            bind:value={inputs.fegliOptionBCoverageMultiples}
            onchange={handleFegliChange}
            class="mt-1 block w-20 pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
          >
            {#each [1, 2, 3, 4, 5] as multiple}
              <option value={multiple}>{multiple}</option>
            {/each}
          </select>
        </div>
      </div>
    {/if}
  </div>

  <!-- Option C - Family Insurance -->
  <div class="space-y-2">
    <div class="inline-flex items-center w-auto">
      <input
        type="checkbox"
        id="fegliOptionC"
        checked={!!inputs.fegliOptionCCoverageMultiples && inputs.fegliOptionCCoverageMultiples > 0}
        onchange={(e) => {
          const target = e.target as HTMLInputElement;
          inputs.fegliOptionCCoverageMultiples = target.checked ? 1 : 0;
          handleFegliChange();
        }}
        class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
      />
      <label for="fegliOptionC" class="ml-2 inline-flex w-auto items-center text-sm font-medium text-gray-700">
        Option C - Family Insurance (coverage for spouse/children)
      </label>
    </div>

    {#if inputs.fegliOptionCCoverageMultiples && inputs.fegliOptionCCoverageMultiples > 0}
      <div class="ml-6">
        <div class="inline-flex items-center w-auto">
          <span class="text-sm text-gray-500 mr-2">Number of multiples:</span>
          <select
            bind:value={inputs.fegliOptionCCoverageMultiples}
            onchange={handleFegliChange}
            class="mt-1 block w-20 pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
          >
            {#each [1, 2, 3, 4, 5] as multiple}
              <option value={multiple}>{multiple}</option>
            {/each}
          </select>
        </div>
        <p class="mt-1 text-xs text-gray-500">
          Note: Each multiple provides $5,000 for spouse and $2,500 per child
        </p>
      </div>
    {/if}
  </div>
</div>
