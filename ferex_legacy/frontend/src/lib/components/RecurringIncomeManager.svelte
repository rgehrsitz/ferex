<script lang="ts">
  import type { OtherRecurringIncomeSource } from '../../types';
  import { createEventDispatcher } from 'svelte';

  let { sources = $bindable([]) } = $props<{ sources?: OtherRecurringIncomeSource[] }>();


  const dispatch = createEventDispatcher();

  // Explicitly define the shape for the form state, aligning with OtherRecurringIncomeSource
  type NewSourceForm = {
    name: string;
    amount: number | null;
    frequency: 'Monthly' | 'Annually';
    startDate: string; // YYYY-MM-DD
    endDate: string;   // YYYY-MM-DD, can be empty string for optional
    isInflationAdjusted: boolean;
  };

  let newSource: NewSourceForm = $state({
    name: '',
    amount: null,
    frequency: 'Annually',
    startDate: '',
    endDate: '',
    isInflationAdjusted: false
  });

  function generateId() {
    return Math.random().toString(36).substring(2, 15);
  }

  function addSource() {
    // Basic validation - enhance as needed
    if (!newSource.name || typeof newSource.amount !== 'number' || newSource.amount <= 0 || !newSource.startDate) {
      alert('Please provide a valid name, positive amount, and start date (YYYY-MM-DD) for the income source.');
      return;
    }
    if (newSource.endDate && newSource.startDate && newSource.endDate < newSource.startDate) {
        alert('End date cannot be before start date.');
        return;
    }

    const completeNewSource: OtherRecurringIncomeSource = {
        id: generateId(),
        name: newSource.name,
        amount: newSource.amount, // Known to be number
        frequency: newSource.frequency,
        startDate: newSource.startDate,
        endDate: newSource.endDate || null, // Store as null if empty string
        isInflationAdjusted: newSource.isInflationAdjusted
    };

    sources = [...sources, completeNewSource];
    newSource = { name: '', amount: null, frequency: 'Annually', startDate: '', endDate: '', isInflationAdjusted: false }; // Reset form
    dispatch('change');
  }

  function removeSource(index: number) {
    sources = sources.filter((source: OtherRecurringIncomeSource, i: number) => i !== index);
    dispatch('change');
  }

  function handleChange() {
    dispatch('change'); // Dispatch change when individual fields of existing items are edited
  }
</script>

<div class="space-y-4">
  <h5 class="text-sm font-medium text-gray-700">Manage Recurring Income Sources</h5>
  
  <!-- Form to add new source -->
  <div class="p-3 border border-gray-200 rounded-md bg-gray-50 space-y-3">
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
      <div>
        <label for="newSourceName" class="block text-xs font-medium text-gray-600">Name</label>
        <input type="text" id="newSourceName" bind:value={newSource.name} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm" placeholder="e.g., Rental Property 1" />
      </div>
      <div>
        <label for="newSourceAmount" class="block text-xs font-medium text-gray-600">Amount ($)</label>
        <input type="number" id="newSourceAmount" bind:value={newSource.amount} min="0" step="any" class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm" placeholder="e.g., 12000" />
      </div>
      <div>
        <label for="newSourceFrequency" class="block text-xs font-medium text-gray-600">Frequency</label>
        <select id="newSourceFrequency" bind:value={newSource.frequency} class="mt-1 block w-full p-2 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md shadow-sm sm:text-sm">
          <option value="Annually">Annually</option>
          <option value="Monthly">Monthly</option>
        </select>
      </div>
      <div>
        <label for="newSourceStartDate" class="block text-xs font-medium text-gray-600">Start Date</label>
        <input type="date" id="newSourceStartDate" bind:value={newSource.startDate} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm" />
      </div>
      <div>
        <label for="newSourceEndDate" class="block text-xs font-medium text-gray-600">End Date (Optional)</label>
        <input type="date" id="newSourceEndDate" bind:value={newSource.endDate} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm" />
      </div>
      <div class="flex items-end">
        <label class="flex items-center text-xs font-medium text-gray-600">
          <input type="checkbox" bind:checked={newSource.isInflationAdjusted} class="mr-2 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500" />
          Inflation Adjusted?
        </label>
      </div>
    </div>
    <button onclick={addSource} class="px-3 py-1.5 bg-blue-500 text-white text-sm font-medium rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
      Add Income Source
    </button>
  </div>

  <!-- List of existing sources -->
  {#if sources && sources.length > 0}
    <div class="space-y-2 pt-3">
      <h6 class="text-xs font-medium text-gray-500">Existing Sources:</h6>
      <ul class="divide-y divide-gray-200">
        {#each sources as source, index (source.id)} <!-- Use unique ID as key -->
          <li class="py-3 px-2 my-1 border border-gray-100 rounded-md hover:bg-gray-50">
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-x-4 gap-y-3 items-center">
              <div>
                <label for={`sourceName-${source.id}`} class="text-xs text-gray-500">Name</label>
                <input type="text" id={`sourceName-${source.id}`} bind:value={source.name} onchange={handleChange} class="text-sm p-1.5 border border-gray-300 rounded-md w-full"/>
              </div>
              <div>
                <label for={`sourceAmount-${source.id}`} class="text-xs text-gray-500">Amount ($)</label>
                <input type="number" id={`sourceAmount-${source.id}`} bind:value={source.amount} onchange={handleChange} min="0" step="any" class="text-sm p-1.5 border border-gray-300 rounded-md w-full"/>
              </div>
              <div>
                <label for={`sourceFrequency-${source.id}`} class="text-xs text-gray-500">Frequency</label>
                <select id={`sourceFrequency-${source.id}`} bind:value={source.frequency} onchange={handleChange} class="text-sm p-1.5 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md w-full">
                  <option value="Annually">Annually</option>
                  <option value="Monthly">Monthly</option>
                </select>
              </div>
              <div>
                <label for={`sourceStartDate-${source.id}`} class="text-xs text-gray-500">Start Date</label>
                <input type="date" id={`sourceStartDate-${source.id}`} bind:value={source.startDate} onchange={handleChange} class="text-sm p-1.5 border border-gray-300 rounded-md w-full"/>
              </div>
              <div>
                <label for={`sourceEndDate-${source.id}`} class="text-xs text-gray-500">End Date</label>
                <input type="date" id={`sourceEndDate-${source.id}`} bind:value={source.endDate} onchange={handleChange} class="text-sm p-1.5 border border-gray-300 rounded-md w-full"/>
              </div>
              <div class="flex items-center pt-3">
                <input type="checkbox" bind:checked={source.isInflationAdjusted} onchange={handleChange} class="mr-2 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500" id={`inflation-${source.id}`} />
                <label for={`inflation-${source.id}`} class="text-xs text-gray-600">Inflation Adj.</label>
              </div>
            </div>
            <div class="mt-2 flex justify-end">
                <button onclick={() => removeSource(index)} class="text-red-500 hover:text-red-700 text-xs font-medium py-1 px-2 rounded hover:bg-red-50">
                  Remove Source
                </button>
            </div>
          </li>
        {/each}
      </ul>
    </div>
  {:else}
    <p class="text-xs text-gray-500 pt-2">No recurring income sources added yet.</p>
  {/if}
</div>
