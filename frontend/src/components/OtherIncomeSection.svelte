<script lang="ts">
  export let data: any = {
    incomeStreams: [
      {
        id: 1,
        name: 'Part-time work',
        amount: 1200,
        frequency: 'monthly',
        startAge: 62,
        endAge: 70,
        taxable: true,
        cola: true
      },
      {
        id: 2,
        name: 'Rental property',
        amount: 1800,
        frequency: 'monthly',
        startAge: 62,
        endAge: 90,
        taxable: true,
        cola: true
      }
    ]
  };

  const frequencyOptions = [
    { value: 'monthly', label: 'Monthly' },
    { value: 'annually', label: 'Annually' },
    { value: 'one-time', label: 'One-Time' }
  ];

  function addIncomeStream() {
    const newId = data.incomeStreams.length > 0 
      ? Math.max(...data.incomeStreams.map(stream => stream.id)) + 1 
      : 1;
    
    data.incomeStreams = [
      ...data.incomeStreams,
      {
        id: newId,
        name: 'New Income',
        amount: 1000,
        frequency: 'monthly',
        startAge: 62,
        endAge: 90,
        taxable: true,
        cola: false
      }
    ];
  }

  function removeIncomeStream(id: number) {
    data.incomeStreams = data.incomeStreams.filter(stream => stream.id !== id);
  }

  function calculateTotalAnnualIncome() {
    return data.incomeStreams.reduce((total, stream) => {
      const annualAmount = stream.frequency === 'monthly' 
        ? stream.amount * 12 
        : stream.frequency === 'annually' 
          ? stream.amount 
          : 0; // One-time amounts are not included in annual total
      
      return total + annualAmount;
    }, 0);
  }

  $: totalAnnualIncome = calculateTotalAnnualIncome();
</script>

<div>
  <div class="space-y-8">
    {#each data.incomeStreams as income, index}
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 border border-gray-200 dark:border-gray-700">
        <div class="flex justify-between items-start mb-4">
          <div class="space-y-1">
            <h3 class="text-lg font-medium text-gray-800 dark:text-gray-200">
              {income.name}
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {income.frequency === 'monthly' ? 'Monthly' : income.frequency === 'annually' ? 'Annual' : 'One-Time'} income
            </p>
          </div>
          <button 
            class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
            on:click={() => removeIncomeStream(income.id)}
          >
            Remove
          </button>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`name-${income.id}`}>
              Description
            </label>
            <input
              id={`name-${income.id}`}
              type="text"
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={income.name}
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`frequency-${income.id}`}>
              Frequency
            </label>
            <select 
              id={`frequency-${income.id}`}
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={income.frequency}
            >
              {#each frequencyOptions as option}
                <option value={option.value}>{option.label}</option>
              {/each}
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`amount-${income.id}`}>
              Amount
            </label>
            <div class="relative rounded-md shadow-sm">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
              </div>
              <input
                id={`amount-${income.id}`}
                type="number"
                min="0"
                step="100"
                class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                bind:value={income.amount}
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`startAge-${income.id}`}>
              Start Age
            </label>
            <input
              id={`startAge-${income.id}`}
              type="number"
              min="0"
              max="120"
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={income.startAge}
            />
          </div>

          {#if income.frequency !== 'one-time'}
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for={`endAge-${income.id}`}>
                End Age
              </label>
              <input
                id={`endAge-${income.id}`}
                type="number"
                min={income.startAge || 0}
                max="120"
                class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                bind:value={income.endAge}
              />
            </div>
          {/if}
        </div>

        <div class="flex items-center gap-4">
          <div class="flex items-center">
            <input
              id={`taxable-${income.id}`}
              type="checkbox"
              class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
              bind:checked={income.taxable}
            />
            <label for={`taxable-${income.id}`} class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
              Taxable
            </label>
          </div>

          {#if income.frequency !== 'one-time'}
            <div class="flex items-center">
              <input
                id={`cola-${income.id}`}
                type="checkbox"
                class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
                bind:checked={income.cola}
              />
              <label for={`cola-${income.id}`} class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
                Apply COLA
              </label>
            </div>
          {/if}
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
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">${totalAnnualIncome.toLocaleString()}</div>
        <div class="text-xs text-gray-500 dark:text-gray-400 mt-1">
          (One-time income not included in total)
        </div>
      </div>
    </div>
  </div>
</div>