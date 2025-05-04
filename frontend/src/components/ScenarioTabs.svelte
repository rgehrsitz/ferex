<script lang="ts">
  import PensionSection from './PensionSection.svelte';
  import SocialSecuritySection from './SocialSecuritySection.svelte';
  import TSPSection from './TSPSection.svelte';
  import TaxSection from './TaxSection.svelte';
  import COLASection from './COLASection.svelte';
  import OtherIncomeSection from './OtherIncomeSection.svelte';
  export let scenario;
  let tab = 'Pension';
  const tabs = [
    { label: 'Pension', comp: PensionSection, prop: 'pension' },
    { label: 'Social Security', comp: SocialSecuritySection, prop: 'socialSecurity' },
    { label: 'TSP', comp: TSPSection, prop: 'tsp' },
    { label: 'Tax', comp: TaxSection, prop: 'tax' },
    { label: 'COLA', comp: COLASection, prop: 'cola' },
    { label: 'Other Income', comp: OtherIncomeSection, prop: 'otherIncome' }
  ];
</script>

<div>
  <div class="flex border-b border-gray-200 dark:border-gray-700 mb-4">
    {#each tabs as t}
      <button 
        class="px-4 py-2 font-medium text-sm transition-colors duration-200 relative" 
        class:text-primary-600={tab===t.label} 
        class:text-gray-500={tab!==t.label}
        class:dark:text-primary-400={tab===t.label} 
        class:dark:text-gray-400={tab!==t.label} 
        class:hover:text-primary-700={tab!==t.label}
        class:dark:hover:text-primary-300={tab!==t.label}
        on:click={() => tab=t.label}
      >
        {t.label}
        {#if tab === t.label}
          <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary-600 dark:bg-primary-400"></div>
        {/if}
      </button>
    {/each}
  </div>
  {#each tabs as t}
    {#if tab === t.label}
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <!-- Initialize the data object if it doesn't exist to prevent errors -->
        {#if !scenario.data[t.prop]}
          <!-- This ensures we don't get undefined errors when binding data -->
          {#if t.prop === 'socialSecurity'}
            {scenario.data.socialSecurity = {
              startAge: 62,
              estimatedMonthlyBenefit: 2000,
              isEligible: true,
              birthYear: 1970,
              birthMonth: 1
            }}
          {:else if t.prop === 'tsp'}
            {scenario.data.tsp = {
              traditionalBalance: 400000,
              rothBalance: 100000,
              contributionRate: 5,
              contributionRateRoth: 5,
              expectedReturn: 6,
              withdrawalStrategy: 'fixed',
              fixedMonthlyWithdrawal: 2000,
              withdrawalRate: 4,
              withdrawalStartAge: 62
            }}
          {:else if t.prop === 'tax'}
            {scenario.data.tax = {
              filingStatus: 'married_joint',
              stateOfResidence: 'VA',
              stateIncomeTaxRate: 0.05,
              itemizedDeductions: 0,
              federalTaxCredits: 0,
              stateTaxCredits: 0,
              age: 62,
              spouseAge: 62
            }}
          {:else if t.prop === 'cola'}
            {scenario.data.cola = {
              assumedInflationRate: 0.025,
              applyCOLAToPension: true,
              applyColaToSocialSecurity: true
            }}
          {:else if t.prop === 'otherIncome'}
            {scenario.data.otherIncome = {
              sources: []
            }}
          {/if}
        {/if}
        
        <svelte:component 
          this={t.comp} 
          bind:data={scenario.data[t.prop]} 
          scenarioId={scenario.id}
          scenarioName={scenario.name} 
        />
      </div>
    {/if}
  {/each}
</div>
