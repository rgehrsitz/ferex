<script lang="ts">
  export let data: any = {
    filingStatus: 'joint',
    federalWithholding: 22,
    stateWithholding: 7,
    stateOfResidence: 'MD',
    useItemizedDeductions: false,
    estimatedItemizedDeductions: 12000,
    age65OrOlder: false,
    spouseAge65OrOlder: false
  };

  const filingStatusOptions = [
    { value: 'single', label: 'Single' },
    { value: 'joint', label: 'Married Filing Jointly' },
    { value: 'separate', label: 'Married Filing Separately' },
    { value: 'hoh', label: 'Head of Household' }
  ];

  const stateOptions = [
    { value: 'AL', label: 'Alabama' },
    { value: 'AK', label: 'Alaska' },
    { value: 'AZ', label: 'Arizona' },
    { value: 'AR', label: 'Arkansas' },
    { value: 'CA', label: 'California' },
    { value: 'CO', label: 'Colorado' },
    { value: 'CT', label: 'Connecticut' },
    { value: 'DE', label: 'Delaware' },
    { value: 'FL', label: 'Florida' },
    { value: 'GA', label: 'Georgia' },
    { value: 'HI', label: 'Hawaii' },
    { value: 'ID', label: 'Idaho' },
    { value: 'IL', label: 'Illinois' },
    { value: 'IN', label: 'Indiana' },
    { value: 'IA', label: 'Iowa' },
    { value: 'KS', label: 'Kansas' },
    { value: 'KY', label: 'Kentucky' },
    { value: 'LA', label: 'Louisiana' },
    { value: 'ME', label: 'Maine' },
    { value: 'MD', label: 'Maryland' },
    { value: 'MA', label: 'Massachusetts' },
    { value: 'MI', label: 'Michigan' },
    { value: 'MN', label: 'Minnesota' },
    { value: 'MS', label: 'Mississippi' },
    { value: 'MO', label: 'Missouri' },
    { value: 'MT', label: 'Montana' },
    { value: 'NE', label: 'Nebraska' },
    { value: 'NV', label: 'Nevada' },
    { value: 'NH', label: 'New Hampshire' },
    { value: 'NJ', label: 'New Jersey' },
    { value: 'NM', label: 'New Mexico' },
    { value: 'NY', label: 'New York' },
    { value: 'NC', label: 'North Carolina' },
    { value: 'ND', label: 'North Dakota' },
    { value: 'OH', label: 'Ohio' },
    { value: 'OK', label: 'Oklahoma' },
    { value: 'OR', label: 'Oregon' },
    { value: 'PA', label: 'Pennsylvania' },
    { value: 'RI', label: 'Rhode Island' },
    { value: 'SC', label: 'South Carolina' },
    { value: 'SD', label: 'South Dakota' },
    { value: 'TN', label: 'Tennessee' },
    { value: 'TX', label: 'Texas' },
    { value: 'UT', label: 'Utah' },
    { value: 'VT', label: 'Vermont' },
    { value: 'VA', label: 'Virginia' },
    { value: 'WA', label: 'Washington' },
    { value: 'WV', label: 'West Virginia' },
    { value: 'WI', label: 'Wisconsin' },
    { value: 'WY', label: 'Wyoming' },
    { value: 'DC', label: 'District of Columbia' }
  ];

  function calculateEffectiveTaxRate(federalRate: number, stateRate: number) {
    // Simple approximation of combined effective rate
    return federalRate + stateRate;
  }

  $: effectiveTaxRate = calculateEffectiveTaxRate(data.federalWithholding, data.stateWithholding);
</script>

<div>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="filingStatus">
          Federal Filing Status
        </label>
        <select 
          id="filingStatus"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.filingStatus}
        >
          {#each filingStatusOptions as status}
            <option value={status.value}>{status.label}</option>
          {/each}
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="federalWithholding">
          Federal Tax Rate
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="federalWithholding"
            type="number"
            min="0"
            max="50"
            step="0.1"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.federalWithholding}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="stateOfResidence">
          State of Residence
        </label>
        <select 
          id="stateOfResidence"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.stateOfResidence}
        >
          {#each stateOptions as state}
            <option value={state.value}>{state.label}</option>
          {/each}
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="stateWithholding">
          State Tax Rate
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="stateWithholding"
            type="number"
            min="0"
            max="15"
            step="0.1"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.stateWithholding}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>
    </div>

    <div class="space-y-4">
      <div class="flex items-center">
        <input
          id="useItemizedDeductions"
          type="checkbox"
          class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          bind:checked={data.useItemizedDeductions}
        />
        <label for="useItemizedDeductions" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
          Use Itemized Deductions
        </label>
      </div>

      {#if data.useItemizedDeductions}
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="estimatedItemizedDeductions">
            Estimated Itemized Deductions
          </label>
          <div class="relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
            </div>
            <input
              id="estimatedItemizedDeductions"
              type="number"
              min="0"
              step="500"
              class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={data.estimatedItemizedDeductions}
            />
          </div>
        </div>
      {/if}

      <div class="flex items-center">
        <input
          id="age65OrOlder"
          type="checkbox"
          class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          bind:checked={data.age65OrOlder}
        />
        <label for="age65OrOlder" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
          Age 65 or Older (for additional standard deduction)
        </label>
      </div>

      {#if data.filingStatus === 'joint'}
        <div class="flex items-center">
          <input
            id="spouseAge65OrOlder"
            type="checkbox"
            class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
            bind:checked={data.spouseAge65OrOlder}
          />
          <label for="spouseAge65OrOlder" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
            Spouse Age 65 or Older
          </label>
        </div>
      {/if}
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Tax Summary</h3>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Federal Rate</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{data.federalWithholding}%</div>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">State Rate</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{data.stateWithholding}%</div>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Combined Effective Rate</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{effectiveTaxRate.toFixed(1)}%</div>
      </div>
    </div>
  </div>
</div>