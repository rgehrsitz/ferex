<script lang="ts">
  import PensionSection from './PensionSection.svelte';
  import SocialSecuritySection from './SocialSecuritySection.svelte';
  import TSPSection from './TSPSection.svelte';
  import TaxSection from './TaxSection.svelte';
  import COLASection from './COLASection.svelte';
  import OtherIncomeSection from './OtherIncomeSection.svelte';
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
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

  // Emit full scenario update on child section changes
  function handleSectionUpdate(prop, updatedData) {
    console.log('ScenarioTabs.handleSectionUpdate received', prop, updatedData);
    console.log('Before update, scenario.data[prop] is', scenario.data[prop]);
    
    // Create a deep clone of the scenario to prevent reference issues
    const newScenario = JSON.parse(JSON.stringify({
      ...scenario,
      data: { ...scenario.data, [prop]: updatedData }
    }));
    
    console.log('ScenarioTabs dispatching update-scenario', newScenario);
    console.log('New scenario will have pension data:', newScenario.data.pension);
    dispatch('update-scenario', newScenario);
  }
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
            on:update={e => handleSectionUpdate(t.prop, e.detail)}
          />
        </div>
      {/key}
    {/if}
  {/each}
</div>
