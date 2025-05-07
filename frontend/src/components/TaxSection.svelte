<script lang="ts">
  import { CalculateTaxLiability } from '../../wailsjs/go/main/App.js';
  import { main } from '../../wailsjs/go/models.js';
  import type { TaxData } from '../types/scenario.js';
  
  import SectionHeader from './SectionHeader.svelte';
  
  // Default data values
  const defaultData = {
    filingStatus: 'single',
    stateOfResidence: 'MD',
    stateIncomeTaxRate: 5,
    itemizedDeductions: 0,
    federalTaxCredits: 0,
    stateTaxCredits: 0,
    spouseAge: 65,
    age: 0
  };

  const { 
    data: propData, 
    onUpdate = (data: any) => {},
    scenarioName = '',
    currentAge = 65 // Default if not provided
  } = $props<{
    data?: TaxData;
    onUpdate?: (data: any) => void;
    scenarioName?: string;
    currentAge?: number;
  }>();
  
  // Create a local copy of data that we can modify
  const data = propData ? { ...propData } : { ...defaultData };
  
  // Set defaults for any missing fields
  data.filingStatus = data.filingStatus || 'single';
  data.stateOfResidence = data.stateOfResidence || 'MD';
  data.stateIncomeTaxRate = data.stateIncomeTaxRate ?? 5;
  data.itemizedDeductions = data.itemizedDeductions ?? 0;
  data.federalTaxCredits = data.federalTaxCredits ?? 0;
  data.stateTaxCredits = data.stateTaxCredits ?? 0;
  // Note: Age is now derived from Social Security data in the parent component
  data.spouseAge = data.spouseAge ?? 65;
  
  // Create local variables for UI binding with $state
  let filingStatus = $state(data.filingStatus);
  let stateOfResidence = $state(data.stateOfResidence);
  let stateIncomeTaxRate = $state(data.stateIncomeTaxRate);
  let itemizedDeductions = $state(data.itemizedDeductions);
  let federalTaxCredits = $state(data.federalTaxCredits);
  let stateTaxCredits = $state(data.stateTaxCredits);
  let spouseAge = $state(data.spouseAge);

  // Use $effect for reactivity
  $effect(() => {
    data.filingStatus = filingStatus;
    data.stateOfResidence = stateOfResidence;
    data.stateIncomeTaxRate = stateIncomeTaxRate;
    data.itemizedDeductions = itemizedDeductions;
    data.federalTaxCredits = federalTaxCredits;
    data.stateTaxCredits = stateTaxCredits;
    data.spouseAge = spouseAge;
    onUpdate({ ...data });
  });

  // currentAge is now received via $props

  
  let loading = $state(false);
  let error = $state('');
  let calculationResult = $state({
    federalTax: 0,
    stateTax: 0,
    totalTax: 0,
    effectiveFederalRate: 0,
    effectiveStateRate: 0,
    effectiveTotalRate: 0,
    notes: ''
  });
  
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
  
  // Update the data object whenever UI variables change
  $effect(() => {
    data.filingStatus = filingStatus;
    data.stateOfResidence = stateOfResidence;
    data.stateIncomeTaxRate = stateIncomeTaxRate;
    data.itemizedDeductions = itemizedDeductions;
    data.federalTaxCredits = federalTaxCredits;
    data.stateTaxCredits = stateTaxCredits;
    // Don't update age in data object - we'll use currentAge passed in from parent
    data.spouseAge = spouseAge;
    
    // Notify parent component of changes
    onUpdate(data);
  });
  
  async function calculateTaxes() {
    try {
      loading = true;
      error = '';
      
      // Create input object for backend
      const taxInput = new main.TaxInput();
      taxInput.filingStatus = filingStatus;
      taxInput.stateOfResidence = stateOfResidence;
      taxInput.stateIncomeTaxRate = stateIncomeTaxRate;
      taxInput.totalIncome = sampleAnnualIncome;
      taxInput.pensionIncome = samplePensionIncome;
      taxInput.socialSecurityIncome = sampleSocialSecurityIncome;
      taxInput.tspWithdrawals = sampleTSPWithdrawals;
      taxInput.nonTaxableIncome = 0;
      taxInput.itemizedDeductions = itemizedDeductions;
      taxInput.federalTaxCredits = federalTaxCredits;
      taxInput.stateTaxCredits = stateTaxCredits;
      taxInput.age = currentAge; // Use the age passed from parent (derived from birth year)
      taxInput.spouseAge = spouseAge;
      
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
  $effect(() => {
    if (filingStatus !== undefined || 
        stateOfResidence || 
        stateIncomeTaxRate !== undefined ||
        itemizedDeductions !== undefined ||
        currentAge !== undefined ||
        spouseAge !== undefined) {
      calculateTaxes();
    }
  });

  const effectiveFederalRate = $derived(calculationResult.effectiveFederalRate * 100 || 0);
  const effectiveStateRate = $derived(calculationResult.effectiveStateRate * 100 || 0);
  const effectiveTotalRate = $derived(calculationResult.effectiveTotalRate * 100 || 0);
</script>

<div>
  <SectionHeader sectionName="Tax Calculator" {scenarioName} />
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="filingStatus">
          Federal Filing Status
        </label>
        <select 
          id="filingStatus"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={filingStatus}
          onchange={() => onUpdate(data)}
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
          bind:value={stateOfResidence}
          onchange={() => onUpdate(data)}
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
            bind:value={stateIncomeTaxRate}
            onchange={() => {
              if (stateIncomeTaxRate) {
                stateIncomeTaxRate = parseFloat(stateIncomeTaxRate.toString());
              }
              onUpdate(data);
            }}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="currentAge">
          Your Age
        </label>
        <div id="currentAge" class="w-full px-3 py-2 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-700 rounded-md text-gray-800 dark:text-gray-200">
          {currentAge} <span class="text-xs text-gray-500 dark:text-gray-400">(from birth year in Social Security tab)</span>
        </div>
      </div>
      
      {#if filingStatus === 'married_joint'}
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
            bind:value={spouseAge}
            onchange={() => {
              if (spouseAge) {
                spouseAge = parseInt(spouseAge.toString());
              }
              onUpdate(data);
            }}
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
            bind:value={itemizedDeductions}
            onchange={() => {
              if (itemizedDeductions) {
                itemizedDeductions = parseFloat(itemizedDeductions.toString());
              }
              onUpdate(data);
            }}
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
            bind:value={federalTaxCredits}
            onchange={() => {
              if (federalTaxCredits) {
                federalTaxCredits = parseFloat(federalTaxCredits.toString());
              }
              onUpdate(data);
            }}
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
            bind:value={stateTaxCredits}
            onchange={() => {
              if (stateTaxCredits) {
                stateTaxCredits = parseFloat(stateTaxCredits.toString());
              }
              onUpdate(data);
            }}
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