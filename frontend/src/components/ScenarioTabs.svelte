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
      {#key tab}
        <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
          <svelte:component
            this={t.comp}
            data={scenario.data[t.prop]}
            scenarioId={scenario.id}
            scenarioName={scenario.name}
          />
        </div>
      {/key}
    {/if}
  {/each}
</div>
