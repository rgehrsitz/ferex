<script lang="ts">
  import PensionSection from './PensionSection.svelte';
  import SocialSecuritySection from './SocialSecuritySection.svelte';
  import SocialSecurityWrapper from './SocialSecurityWrapper.svelte';
  import TSPSection from './TSPSection.svelte';
  import TaxSection from './TaxSection.svelte';
  import COLASection from './COLASection.svelte';
  import OtherIncomeSection from './OtherIncomeSection.svelte';
  import { onMount, onDestroy } from 'svelte';
  import type { Scenario, ScenarioData } from '../types/scenario.js';
  import Dynamic from './Dynamic.svelte';

  // Use a bindable prop for two-way binding
  let { scenario } = $props<{ scenario: Scenario }, { scenario: Scenario }>();

  // State with Svelte 5 runes
  let tabsByScenario = $state<Record<string | number, string>>({});
  let componentRefs = $state<Record<string, any>>({});

  let statusMessage = $state<{ text: string; type: string } | null>(null);
  let statusTimeout = $state<ReturnType<typeof setTimeout> | null>(null);

  // Derived active tab based on the current scenario
  let activeTab = $derived(tabsByScenario[scenario?.id] || 'Pension');

  // Reset componentRefs whenever scenario.id changes (no dependency loop)
  $effect(() => {
    componentRefs = {};
    console.log('ScenarioTabs: Scenario changed, reset componentRefs');
  });

  // Function to set the tab for the current scenario
  function setTab(newTab: string) {
    console.log(`Setting tab for scenario ${scenario.id} to ${newTab}`);
    tabsByScenario = { ...tabsByScenario, [scenario.id]: newTab };
  }

  
  // Tab definitions
  const tabs = [
    { label: 'Pension', comp: PensionSection, prop: 'pension' },
    { label: 'Social Security', comp: SocialSecuritySection, prop: 'socialSecurity' },
    { label: 'TSP', comp: TSPSection, prop: 'tsp' },
    { label: 'Tax', comp: TaxSection, prop: 'tax' },
    { label: 'COLA', comp: COLASection, prop: 'cola' },
    { label: 'Other Income', comp: OtherIncomeSection, prop: 'otherIncome' }
  ];
  
  // Listen for global calculate events
  function handleGlobalCalculate(event: CustomEvent<any>) {
    console.log('ScenarioTabs received global-calculate event', event.detail);
    
    // Check if this event is for the current scenario
    if (event.detail.scenarioId !== scenario.id) {
      console.log('Event is for a different scenario, ignoring');
      return;
    }
    
    // Show calculation status
    showStatusMessage({ 
      text: 'Calculating all sections...',
      type: 'info'
    });
    
    // Trigger calculations in all tab components
    calculateAllSections();
  }
  
  function calculateAllSections(): void {
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
  function showStatusMessage(message: { text: string; type: string }): void {
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
    console.log(`ScenarioTabs for scenario ${scenario.id} mounted - adding event listeners`);
    document.addEventListener('global-calculate', handleGlobalCalculate as EventListener);
  });
  
  onDestroy(() => {
    console.log(`ScenarioTabs for scenario ${scenario.id} destroyed - removing event listeners`);
    document.removeEventListener('global-calculate', handleGlobalCalculate as EventListener);
    if (statusTimeout) {
      clearTimeout(statusTimeout);
    }
    
    // Clean up any references to this component
    componentRefs = {};
  });

  // Handle section updates with both direct binding and event dispatching
  function handleSectionUpdate<K extends keyof ScenarioData>(
    prop: K, 
    updatedData: ScenarioData[K]
  ): void {
    if (!scenario || !scenario.data) {
      console.warn('Cannot update scenario: scenario or scenario.data is undefined');
      return;
    }
    
    try {
      // Create an updated scenario with the new data
      const updatedScenario = {
        ...scenario,
        data: {
          ...scenario.data,
          [prop]: updatedData
        }
      };
      
      // Update the local scenario via two-way binding
      scenario = updatedScenario;
      
      // Also dispatch an event to notify parent components
      dispatch('scenarioChange', updatedScenario);
      
      console.log(`Updated scenario.data.${String(prop)}`);
    } catch (error) {
      console.error('Error updating scenario:', error);
    }
  }
</script>

<div>
  <!-- Tab navigation -->
  <div class="flex border-b border-gray-200 dark:border-gray-700 mb-4">
    {#each tabs as t}
      <button 
        class="px-4 py-2 font-medium text-sm transition-colors duration-200 relative" 
        class:text-primary-600={activeTab===t.label} 
        class:text-gray-500={activeTab!==t.label}
        class:dark:text-primary-400={activeTab===t.label} 
        class:dark:text-gray-400={activeTab!==t.label} 
        class:hover:text-primary-700={activeTab!==t.label}
        class:dark:hover:text-primary-300={activeTab!==t.label}
        onclick={() => setTab(t.label)}
      >
        {t.label}
        {#if activeTab === t.label}
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

  <!-- Tab content -->
  {#each tabs as t}
  {#if activeTab === t.label}
    <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      {#if t.label === 'Pension' && scenario && scenario.data}
        <PensionSection 
          data={scenario.data.pension} 
          scenarioName={scenario.name || ''} 
          onUpdate={(data: any) => scenario && handleSectionUpdate('pension', data)} 
        />
      {:else if t.comp}
        <!-- Debug logging removed -->
        <Dynamic 
          component={t.comp}
          data={(scenario?.data && (scenario.data as any)[t.prop]) || {}}
          scenarioId={scenario?.id || 0}
          scenarioName={scenario?.name || ''}
          onUpdate={(data: any) => scenario && handleSectionUpdate(t.prop as keyof ScenarioData, data)}
          bind:this={componentRefs[t.prop]}
          currentAge={t.prop === 'tax' && scenario?.data && (scenario.data as any).socialSecurity?.birthYear 
            ? new Date().getFullYear() - (scenario.data as any).socialSecurity.birthYear 
            : undefined}
        />
      {/if}
    </div>
  {/if}
{/each}
</div>