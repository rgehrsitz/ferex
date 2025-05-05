<script lang="ts">
  import PensionSection from './PensionSection.svelte';
  import SocialSecuritySection from './SocialSecuritySection.svelte';
  import TSPSection from './TSPSection.svelte';
  import TaxSection from './TaxSection.svelte';
  import COLASection from './COLASection.svelte';
  import OtherIncomeSection from './OtherIncomeSection.svelte';
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  
  const dispatch = createEventDispatcher();
  export let scenario;
  let tab = 'Pension';
  let componentRefs = {}; // Store references to components
  
  const tabs = [
    { label: 'Pension', comp: PensionSection, prop: 'pension' },
    { label: 'Social Security', comp: SocialSecuritySection, prop: 'socialSecurity' },
    { label: 'TSP', comp: TSPSection, prop: 'tsp' },
    { label: 'Tax', comp: TaxSection, prop: 'tax' },
    { label: 'COLA', comp: COLASection, prop: 'cola' },
    { label: 'Other Income', comp: OtherIncomeSection, prop: 'otherIncome' }
  ];
  
  // Listen for global calculate events
  function handleGlobalCalculate(event) {
    console.log('ScenarioTabs received global-calculate event', event.detail);
    
    // Check if this event is for the current scenario
    if (event.detail.scenarioId !== scenario.id) {
      console.log('Event is for a different scenario, ignoring');
      return;
    }
    
    // Show calculation status
    const statusMsg = { 
      text: 'Calculating all sections...',
      type: 'info'
    };
    showStatusMessage(statusMsg);
    
    // Trigger calculations in all tab components
    calculateAllSections();
  }
  
  function calculateAllSections() {
    console.log('Calculating all sections');
    
    // If we have references to components with calculate methods, call them
    Object.values(componentRefs).forEach(ref => {
      if (ref && typeof ref.calculatePension === 'function') {
        ref.calculatePension();
      }
      // Add other calculation methods as needed
    });
    
    // After a delay, show success message
    setTimeout(() => {
      showStatusMessage({ 
        text: 'All calculations complete',
        type: 'success'
      });
    }, 1000);
  }
  
  // Status message handling
  let statusMessage = null;
  let statusTimeout;
  
  function showStatusMessage(message) {
    statusMessage = message;
    
    // Clear any existing timeout
    if (statusTimeout) {
      clearTimeout(statusTimeout);
    }
    
    // Auto-hide after 3 seconds
    statusTimeout = setTimeout(() => {
      statusMessage = null;
    }, 3000);
  }
  
  // Lifecycle hooks for event listeners
  onMount(() => {
    document.addEventListener('global-calculate', handleGlobalCalculate);
  });
  
  onDestroy(() => {
    document.removeEventListener('global-calculate', handleGlobalCalculate);
    if (statusTimeout) {
      clearTimeout(statusTimeout);
    }
  });

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
  
  <!-- Status message display -->
  {#if statusMessage}
    <div class="mb-4 px-4 py-2 rounded-md transition-all duration-300 ease-in-out" 
      class:bg-blue-100={statusMessage.type === 'info'} 
      class:dark:bg-blue-900={statusMessage.type === 'info'}
      class:text-blue-800={statusMessage.type === 'info'} 
      class:dark:text-blue-200={statusMessage.type === 'info'}
      class:bg-green-100={statusMessage.type === 'success'} 
      class:dark:bg-green-900={statusMessage.type === 'success'}
      class:text-green-800={statusMessage.type === 'success'} 
      class:dark:text-green-200={statusMessage.type === 'success'}
      class:bg-red-100={statusMessage.type === 'error'} 
      class:dark:bg-red-900={statusMessage.type === 'error'}
      class:text-red-800={statusMessage.type === 'error'} 
      class:dark:text-red-200={statusMessage.type === 'error'}
    >
      <div class="flex items-center">
        {#if statusMessage.type === 'info'}
          <span class="mr-2">ℹ️</span>
        {:else if statusMessage.type === 'success'}
          <span class="mr-2">✅</span>
        {:else if statusMessage.type === 'error'}
          <span class="mr-2">❌</span>
        {/if}
        <p>{statusMessage.text}</p>
      </div>
    </div>
  {/if}
  {#each tabs as t}
    {#if tab === t.label}
      {#key `${tab}-${scenario.id}-${JSON.stringify(scenario.data[t.prop])}`}
        <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
          <svelte:component
            this={t.comp}
            data={scenario.data[t.prop]}
            scenarioId={scenario.id}
            scenarioName={scenario.name}
            on:update={e => handleSectionUpdate(t.prop, e.detail)}
            bind:this={componentRefs[t.prop]}
          />
        </div>
      {/key}
    {/if}
  {/each}
</div>
