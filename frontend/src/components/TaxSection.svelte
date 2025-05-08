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
    data = null, 
    onUpdate = (data: TaxData) => {},
    scenarioName = '',
    currentAge = 65 // Default if not provided
  } = $props<{
    data?: TaxData | null;
    onUpdate?: (data: TaxData) => void;
    scenarioName?: string;
    currentAge?: number;
  }>();
  
  // Using $state for local data, initialized with defaults or incoming data
  let localData = $state<TaxData>(
    data ? { ...defaultData, ...data } : { ...defaultData }
  );
  
  // Store a previous version of localData as a string for comparison
  let previousLocalDataString = $state("");
  
  // Initialize previousLocalDataString after localData is defined
  $effect(() => {
    if (previousLocalDataString === "") {
      previousLocalDataString = JSON.stringify(localData);
    }
  });
  
  // Track if we're updating from props to prevent loops
  let updatingFromProps = $state(false);

  // Manual function to update local data when props change
  function updateLocalDataFromProps() {
    if (!data) return;
    
    try {
      const dataString = JSON.stringify(data);
      // Only update if different
      if (dataString !== previousLocalDataString) {
        console.log('TaxSection: Props changed, updating local data');
        updatingFromProps = true;
        localData = { ...defaultData, ...data };
        previousLocalDataString = dataString;
      }
    } finally {
      updatingFromProps = false;
    }
  }

  // Update local data when props change - but only on initial props or when they actually change
  let previousPropString = $state('');
  
  $effect(() => {
    if (data) {
      const propString = JSON.stringify(data);
      if (propString !== previousPropString) {
        previousPropString = propString;
        updateLocalDataFromProps();
      }
    }
  });

  // Manual function to update the props from localData
  function updateParentData() {
    if (updatingFromProps) return; // Skip if we're currently updating from props
    
    const currentString = JSON.stringify(localData);
    
    // Only update if the data has actually changed
    if (currentString !== previousLocalDataString) {
      console.log('TaxSection: localData changed, updating parent');
      previousLocalDataString = currentString;
      
      // Create a deep copy to avoid reference issues
      const dataForParent = JSON.parse(currentString);
      
      // Notify parent via callback
      onUpdate(dataForParent);
    }
  }

  // Using direct state values for UI binding
  // This is simpler and avoids issues with derived bindings
  // We'll update parent in the change handlers instead of automatically
  
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
  
  async function calculateTaxes() {
    try {
      loading = true;
      error = '';
      
      // Create input object for backend
      const taxInput = new main.TaxInput();
      taxInput.filingStatus = localData.filingStatus;
      taxInput.stateOfResidence = localData.stateOfResidence;
      taxInput.stateIncomeTaxRate = localData.stateIncomeTaxRate;
      taxInput.totalIncome = sampleAnnualIncome;
      taxInput.pensionIncome = samplePensionIncome;
      taxInput.socialSecurityIncome = sampleSocialSecurityIncome;
      taxInput.tspWithdrawals = sampleTSPWithdrawals;
      taxInput.nonTaxableIncome = 0;
      taxInput.itemizedDeductions = localData.itemizedDeductions;
      taxInput.federalTaxCredits = localData.federalTaxCredits;
      taxInput.stateTaxCredits = localData.stateTaxCredits;
      taxInput.age = currentAge; // Use the age passed from parent (derived from birth year)
      taxInput.spouseAge = localData.spouseAge;
      
      // Call backend API
      calculationResult = await CalculateTaxLiability(taxInput);
      
    } catch (err) {
      console.error("Error calculating tax liability:", err);
      error = 'Failed to calculate tax liability. Please try again.';
    } finally {
      loading = false;
    }
  }

  // Function to handle explicit field changes
  // No need to create a data alias - we need to update all references in the template

  function handleFieldChange() {
    // Update parent with new data first
    updateParentData();
    // Then calculate taxes with updated values
    calculateTaxes();
  }

  // Trigger calculation when component mounts or when prop values change
  $effect(() => {
    if (localData && currentAge !== undefined) {
      calculateTaxes();
    }
  });

  // Create a compatibility object for backward references
  // This is a temporary fix to prevent "data is not defined" errors
  // We don't need this since data is already defined as a prop

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
          bind:value={localData.filingStatus}
          onchange={handleFieldChange}
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
          bind:value={localData.stateOfResidence}
          onchange={handleFieldChange}
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
            bind:value={localData.stateIncomeTaxRate}
            onchange={() => {
              if (typeof localData.stateIncomeTaxRate === 'string') {
                localData.stateIncomeTaxRate = parseFloat(localData.stateIncomeTaxRate);
              }
              handleFieldChange();
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
      
      {#if localData.filingStatus === 'married_joint'}
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
            bind:value={localData.spouseAge}
            onchange={() => {
              if (typeof localData.spouseAge === 'string') {
                localData.spouseAge = parseInt(localData.spouseAge);
              }
              handleFieldChange();
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
            bind:value={localData.itemizedDeductions}
            onchange={() => {
              if (typeof localData.itemizedDeductions === 'string') {
                localData.itemizedDeductions = parseFloat(localData.itemizedDeductions);
              }
              handleFieldChange();
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
            bind:value={localData.federalTaxCredits}
            onchange={() => {
              if (typeof localData.federalTaxCredits === 'string') {
                localData.federalTaxCredits = parseFloat(localData.federalTaxCredits);
              }
              handleFieldChange();
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
            bind:value={localData.stateTaxCredits}
            onchange={() => {
              if (typeof localData.stateTaxCredits === 'string') {
                localData.stateTaxCredits = parseFloat(localData.stateTaxCredits);
              }
              handleFieldChange();
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