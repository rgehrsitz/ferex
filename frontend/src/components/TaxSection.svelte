<script lang="ts">
  import { CalculateTaxLiability } from '../../wailsjs/go/main/App';
  import { main } from '../../wailsjs/go/models';
  import type { TaxData } from '../types/scenario';
  
  export let data: TaxData;
  
  let loading = false;
  let error = '';
  let calculationResult: any = {
    federalTax: 0,
    stateTax: 0,
    totalTax: 0,
    effectiveFederalRate: 0,
    effectiveStateRate: 0,
    effectiveTotalRate: 0,
    notes: ''
  };
  
  // Sample projected income (would be derived from other sections)
  let sampleAnnualIncome = 100000;
  let sampleSocialSecurityIncome = 30000;
  let samplePensionIncome = 50000;
  let sampleTSPWithdrawals = 20000;

  const filingStatusOptions = [
    { value: 'single', label: 'Single' },
    { value: 'married_joint', label: 'Married Filing Jointly' },
    { value: 'married_separate', label: 'Married Filing Separately' },
    { value: 'head_of_household', label: 'Head of Household' }
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
  
  async function calculateTaxes() {
    try {
      loading = true;
      error = '';
      
      // Create input object for backend
      const taxInput = new main.TaxInput();
      taxInput.filingStatus = data.filingStatus;
      taxInput.stateOfResidence = data.stateOfResidence;
      taxInput.stateIncomeTaxRate = data.stateIncomeTaxRate;
      taxInput.totalIncome = sampleAnnualIncome;
      taxInput.pensionIncome = samplePensionIncome;
      taxInput.socialSecurityIncome = sampleSocialSecurityIncome;
      taxInput.tspWithdrawals = sampleTSPWithdrawals;
      taxInput.nonTaxableIncome = 0;
      taxInput.itemizedDeductions = data.itemizedDeductions;
      taxInput.federalTaxCredits = data.federalTaxCredits;
      taxInput.stateTaxCredits = data.stateTaxCredits;
      taxInput.age = data.age;
      taxInput.spouseAge = data.spouseAge;
      
      // Call backend API
      calculationResult = await CalculateTaxLiability(taxInput);
      
    } catch (err) {
      console.error("Error calculating tax liability:", err);
      error = 'Failed to calculate tax liability. Please try again.';
    } finally {
      loading = false;
    }
  }

  // Trigger calculation when relevant inputs change
  $: {
    if (data.filingStatus !== undefined || 
        data.stateOfResidence || 
        data.stateIncomeTaxRate !== undefined ||
        data.itemizedDeductions !== undefined ||
        data.age !== undefined ||
        data.spouseAge !== undefined) {
      calculateTaxes();
    }
  }

  $: effectiveFederalRate = calculationResult.effectiveFederalRate * 100 || 0;
  $: effectiveStateRate = calculationResult.effectiveStateRate * 100 || 0;
  $: effectiveTotalRate = calculationResult.effectiveTotalRate * 100 || 0;
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
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="stateIncomeTaxRate">
          State Tax Rate
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="stateIncomeTaxRate"
            type="number"
            min="0"
            max="15"
            step="0.1"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.stateIncomeTaxRate}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="age">
          Your Age
        </label>
        <input
          id="age"
          type="number"
          min="55"
          max="100"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.age}
        />
      </div>
      
      {#if data.filingStatus === 'married_joint'}
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="spouseAge">
            Spouse Age
          </label>
          <input
            id="spouseAge"
            type="number"
            min="55"
            max="100"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.spouseAge}
          />
        </div>
      {/if}
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="itemizedDeductions">
          Itemized Deductions
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="itemizedDeductions"
            type="number"
            min="0"
            step="500"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.itemizedDeductions}
          />
        </div>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          Enter 0 to use standard deduction
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="federalTaxCredits">
          Federal Tax Credits
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="federalTaxCredits"
            type="number"
            min="0"
            step="100"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.federalTaxCredits}
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="stateTaxCredits">
          State Tax Credits
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="stateTaxCredits"
            type="number"
            min="0"
            step="100"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.stateTaxCredits}
          />
        </div>
      </div>
      
      <div class="bg-blue-50 dark:bg-blue-900 p-4 rounded-lg mt-4">
        <h4 class="text-sm font-medium text-blue-800 dark:text-blue-300 mb-2">Sample Income (for tax calculation)</h4>
        <ul class="text-xs text-blue-700 dark:text-blue-400">
          <li>Annual Income: ${sampleAnnualIncome.toLocaleString()}</li>
          <li>Pension: ${samplePensionIncome.toLocaleString()}</li>
          <li>Social Security: ${sampleSocialSecurityIncome.toLocaleString()}</li>
          <li>TSP Withdrawals: ${sampleTSPWithdrawals.toLocaleString()}</li>
        </ul>
        <p class="text-xs text-blue-700 dark:text-blue-400 mt-2">
          This sample income is used to calculate your approximate tax rates. Your actual retirement income may differ.
        </p>
      </div>
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Tax Summary</h3>
    
    {#if error}
      <div class="bg-red-100 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-700 dark:text-red-300 p-4 rounded-lg mb-4">
        {error}
      </div>
    {/if}
    
    {#if loading}
      <div class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 dark:border-primary-400"></div>
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
            <div class="text-sm text-gray-500 dark:text-gray-400">Federal Tax</div>
            <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
              ${calculationResult.federalTax.toLocaleString(undefined, { maximumFractionDigits: 0 })}
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
            <div class="text-sm text-gray-500 dark:text-gray-400">State Tax</div>
            <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
              ${calculationResult.stateTax.toLocaleString(undefined, { maximumFractionDigits: 0 })}
            </div>
          </div>
        </div>
        
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Total Tax</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${calculationResult.totalTax.toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Federal Effective Rate</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{effectiveFederalRate.toFixed(1)}%</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">State Effective Rate</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{effectiveStateRate.toFixed(1)}%</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Combined Effective Rate</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{effectiveTotalRate.toFixed(1)}%</div>
        </div>
      </div>
      
      {#if calculationResult.notes}
        <div class="mt-4 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Notes</h4>
          <div class="text-sm text-gray-600 dark:text-gray-400">
            {calculationResult.notes}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>