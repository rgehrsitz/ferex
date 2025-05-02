<script lang="ts">
  import type { OtherIncomeData, OtherIncomeSource } from '../types/scenario';
  import { v4 as uuidv4 } from 'uuid';
  
  export let data: OtherIncomeData = {
    sources: []
  };
  
  // Initialize sources array if it doesn't exist
  if (!data.sources) {
    data.sources = [];
  }
  
  // Add default income streams if none exist - but only visually, not to the actual data binding
  // This prevents the UI from showing defaults that aren't actually saved
  let displaySources = [];
  
  // When data.sources changes, update displaySources
  $: {
    // Create a deep copy to avoid direct mutation
    if (!data.sources || data.sources.length === 0) {
      // Ensure we have at least one income source
      displaySources = [
        {
          id: uuidv4(),
          name: 'Part-time work',
          amount: 1200,
          frequency: 'monthly',
          startAge: 62,
          endAge: 70,
          applyCola: true
        }
      ];
      
      // If this is a blank slate, add second default option
      if (!data.sources || data.sources.length === 0) {
        displaySources.push({
          id: uuidv4(),
          name: 'Rental property',
          amount: 1800,
          frequency: 'monthly',
          startAge: 62,
          endAge: 90,
          applyCola: true
        });
      }
      
      // Update the actual data to match display sources
      data.sources = JSON.parse(JSON.stringify(displaySources));
    } else {
      // Make sure we have at least one source
      if (data.sources.length === 0) {
        data.sources = [{
          id: uuidv4(),
          name: 'New Income',
          amount: 1000,
          frequency: 'monthly',
          startAge: 62,
          endAge: 90,
          applyCola: false
        }];
      }
      
      // Update display sources from data
      displaySources = [...data.sources];
    }
  }

  const frequencyOptions = [
    { value: 'monthly', label: 'Monthly' },
    { value: 'annual', label: 'Annual' }
  ];

  function addIncomeStream() {
    const newSource = {
      id: uuidv4(),
      name: 'New Income',
      amount: 1000,
      frequency: 'monthly',
      startAge: 62,
      endAge: 90,
      applyCola: false
    };
    
    // Create a new array with the added source
    displaySources = [...displaySources, newSource];
    
    // Ensure deep copy when updating the original data
    data.sources = JSON.parse(JSON.stringify(displaySources));
    
    // Re-calculate total
    totalAnnualIncome = calculateTotalAnnualIncome();
    
    console.log('Added new income source. Total sources:', displaySources.length);
  }

  function removeIncomeStream(id: string) {
    console.log("Removing income stream with ID:", id);
    
    // Make sure ID is a string
    const idStr = String(id);
    
    // Filter out the matching source
    const updatedSources = displaySources.filter(source => String(source.id) !== idStr);
    
    // If removing would result in empty sources, keep at least one
    if (updatedSources.length === 0) {
      console.log("Cannot remove last income source - at least one must remain");
      return; // Don't update sources if it would leave an empty array
    }
    
    // Update display sources with the filtered array
    displaySources = updatedSources;
    
    // Ensure deep copy when updating the original data
    data.sources = JSON.parse(JSON.stringify(displaySources));
    
    // Re-calculate total
    totalAnnualIncome = calculateTotalAnnualIncome();
  }

  function calculateTotalAnnualIncome() {
    return data.sources.reduce((total, source) => {
      const annualAmount = source.frequency === 'monthly' 
        ? source.amount * 12 
        : source.amount;
      
      return total + annualAmount;
    }, 0);
  }

  $: totalAnnualIncome = calculateTotalAnnualIncome();
</script>

<div>
  <div class="space-y-8">
    {#each displaySources as source, index}
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 border border-gray-200 dark:border-gray-700">
        <div class="flex justify-between items-start mb-4">
          <div class="space-y-1">
            <h3 class="text-lg font-medium text-gray-800 dark:text-gray-200">
              {source.name}
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {source.frequency === 'monthly' ? 'Monthly' : 'Annual'} income
            </p>
          </div>
          <button 
            type="button"
            class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 px-2 py-1 rounded"
            on:click|preventDefault={() => removeIncomeStream(source.id)}
          >
            Remove
          </button>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`name-${source.id}`}>
              Description
            </label>
            <input
              id={`name-${source.id}`}
              type="text"
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={source.name}
              on:change={() => {
                // Update the original data on change
                data.sources = JSON.parse(JSON.stringify(displaySources));
              }}
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`frequency-${source.id}`}>
              Frequency
            </label>
            <select 
              id={`frequency-${source.id}`}
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={source.frequency}
              on:change={() => {
                data.sources = JSON.parse(JSON.stringify(displaySources));
              }}
            >
              {#each frequencyOptions as option}
                <option value={option.value}>{option.label}</option>
              {/each}
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`amount-${source.id}`}>
              Amount
            </label>
            <div class="relative rounded-md shadow-sm">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
              </div>
              <input
                id={`amount-${source.id}`}
                type="number"
                min="0"
                step="100"
                class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                bind:value={source.amount}
              on:change={() => {
                data.sources = JSON.parse(JSON.stringify(displaySources));
              }}
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`startAge-${source.id}`}>
              Start Age
            </label>
            <input
              id={`startAge-${source.id}`}
              type="number"
              min="0"
              max="120"
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={source.startAge}
              on:change={() => {
                data.sources = JSON.parse(JSON.stringify(displaySources));
              }}
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`endAge-${source.id}`}>
              End Age
            </label>
            <input
              id={`endAge-${source.id}`}
              type="number"
              min={source.startAge || 0}
              max="120"
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={source.endAge}
              on:change={() => {
                data.sources = JSON.parse(JSON.stringify(displaySources));
              }}
            />
          </div>
        </div>

        <div class="flex items-center gap-4">
          <div class="flex items-center">
            <input
              id={`cola-${source.id}`}
              type="checkbox"
              class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
              bind:checked={source.applyCola}
              on:change={() => {
                data.sources = JSON.parse(JSON.stringify(displaySources));
              }}
            />
            <label for={`cola-${source.id}`} class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
              Apply COLA
            </label>
          </div>
        </div>
      </div>
    {/each}

    <div class="flex justify-center">
      <button 
        class="px-4 py-2 bg-primary-600 hover:bg-primary-700 dark:bg-primary-700 dark:hover:bg-primary-600 text-white rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
        on:click={addIncomeStream}
      >
        + Add Income Stream
      </button>
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Other Income Summary</h3>
    <div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Total Annual Additional Income</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">${totalAnnualIncome.toLocaleString(undefined, { maximumFractionDigits: 0 })}</div>
      </div>
    </div>
  </div>
</div>