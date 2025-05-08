<script module lang="ts">
  // Import API functions directly
  import { api } from './api.js';
</script>

<script lang="ts">
  import Header from './components/Header.svelte';
  import ScenarioSidebar from './components/ScenarioSidebar.svelte';
  import ScenarioTabs from './components/ScenarioTabs.svelte';
  import CompareView from './components/CompareView.svelte';
  import type { Scenario, ScenarioData } from './types/scenario.js';
  import { createDefaultScenario } from './utils/createDefaultScenario.js';

  // State with Svelte 5 runes
  let scenarios = $state([createDefaultScenario(1, 'Scenario 1'), createDefaultScenario(2, 'Scenario 2')]);
  let selectedScenarioId = $state<number | null>(1);
  let compareScenarioId = $state<number | null>(null);
  let filename = $state('');
  let statusMessage = $state('');
  let isLoading = $state(false);
  let showSaveDialog = $state(false);
  let showLoadDialog = $state(false);
  let savedFiles = $state<string[]>([]);
  let selectedFile = $state('');

  // Debug: log scenarios array whenever it changes
  $effect(() => {
    // Use $state.snapshot to avoid Svelte proxy warnings
    console.log('App scenarios updated', $state.snapshot(scenarios));
  });

  // Compute index of selected scenario for two-way binding
  let selectedScenarioIndex = $derived(
    scenarios.findIndex((s: Scenario) => s.id === selectedScenarioId)
  );

  // Use a simple reference to the selected scenario that's not reactive itself
  // This avoids creating reactivity cycles while still allowing the component to work
  let storedSelectedScenario: Scenario | null = null;
  
  // Instead of having automatic reactive bindings, we'll use a controlled approach
  // with explicit functions for copying between scenario states

  // Get the current selected scenario from the scenarios array based on the ID
  function getSelectedScenarioFromArray(): Scenario | null {
    if (selectedScenarioIndex < 0) {
      return null;
    }
    return scenarios[selectedScenarioIndex] || null;
  }
  
  // Manual function to update the scenarios array with changes from a component
  function updateScenario(updatedScenario: Scenario): void {
    if (!updatedScenario || !updatedScenario.id) return;
    
    // Find the index of the scenario to update
    const index = scenarios.findIndex(s => s.id === updatedScenario.id);
    if (index < 0) return;
    
    // Create a new array and update the specific scenario
    const newScenarios = [...scenarios];
    newScenarios[index] = { ...updatedScenario };
    
    // Update the scenarios array with the new value
    scenarios = newScenarios;
    
    // Also update our stored selected scenario if it's the one that was updated
    if (selectedScenarioId === updatedScenario.id) {
      storedSelectedScenario = { ...updatedScenario };
    }
    
    console.log(`Updated scenario ${updatedScenario.id} in array`);
  }
  
  // Make sure a scenario has all required data properties
  function ensureScenarioData(scenario: Scenario): Scenario {
    if (!scenario) return scenario;
    
    // Initialize data if it doesn't exist
    if (!scenario.data) {
      scenario.data = {};
    }
    
    // Initialize all required sections
    if (!scenario.data.pension) {
      scenario.data.pension = {
        system: 'FERS',
        high3Salary: 100000,
        yearsOfService: 30,
        ageAtRetirement: 62,
        unusedSickLeaveMonths: 6,
        survivorBenefitOption: 'full',
        isPartTime: false,
        partTimeProrationFactor: 1.0,
        militaryService: 0
      };
    }
    
    if (!scenario.data.socialSecurity) {
      scenario.data.socialSecurity = {
        startAge: 62,
        estimatedMonthlyBenefit: 2000,
        isEligible: true,
        birthYear: 1970,
        birthMonth: 1
      };
    }
    
    if (!scenario.data.tsp) {
      scenario.data.tsp = {
        traditionalBalance: 400000,
        rothBalance: 100000,
        contributionRate: 5,
        contributionRateRoth: 5,
        expectedReturn: 6,
        withdrawalStrategy: 'fixed',
        fixedMonthlyWithdrawal: 2000,
        withdrawalRate: 4,
        withdrawalStartAge: 62
      };
    }
    
    if (!scenario.data.tax) {
      scenario.data.tax = {
        filingStatus: 'married_joint',
        stateOfResidence: 'VA',
        stateIncomeTaxRate: 0.05,
        itemizedDeductions: 0,
        federalTaxCredits: 0,
        stateTaxCredits: 0,
        age: 62,
        spouseAge: 62
      };
    }
    
    if (!scenario.data.cola) {
      scenario.data.cola = {
        assumedInflationRate: 2.5,
        applyCOLAToPension: true,
        applyColaToSocialSecurity: true
      };
    }
    
    if (!scenario.data.otherIncome) {
      scenario.data.otherIncome = {
        sources: []
      };
    }
    
    return scenario;
  }

  // This computed property is what we pass to components
  // It gets the current scenario in a way that won't create an infinite loop
  let selectedScenario = $derived({
    // This is a getter that will be called when components read this property
    get() {
      // If we don't have a stored scenario, or we need to refresh it,
      // get it from the array
      if (!storedSelectedScenario || storedSelectedScenario.id !== selectedScenarioId) {
        const currentScenario = getSelectedScenarioFromArray();
        if (currentScenario) {
          // Store a clone to avoid reference issues
          const scenarioClone = JSON.parse(JSON.stringify(currentScenario));
          // Make sure it has all required data
          storedSelectedScenario = ensureScenarioData(scenarioClone);
        } else {
          storedSelectedScenario = null;
        }
      }
      return storedSelectedScenario;
    },
    // This is a setter that will be called when components update this property
    set(value: Scenario | null) {
      if (!value) {
        storedSelectedScenario = null;
        return;
      }
      // Store the new value with all required data
      storedSelectedScenario = ensureScenarioData({ ...value });
      // Update the scenarios array with the changes
      updateScenario(storedSelectedScenario);
    }
  });

  function addScenarioHandler(): void {
    if (!scenarios || scenarios.length === 0) {
      const newScenario = createDefaultScenario(1, `Scenario 1`);
      scenarios = [newScenario];
      selectedScenarioId = 1;
      compareScenarioId = null;
      return;
    }
    const maxId = Math.max(...scenarios.map((s: Scenario) => s.id), 0);
    const newId = maxId + 1;
    const newScenario = createDefaultScenario(newId, `Scenario ${newId}`);
    scenarios = [...scenarios, newScenario];
    selectedScenarioId = newId;
    compareScenarioId = null;
  }

  function duplicateScenarioHandler(id: number): void {
    const orig = scenarios.find((s: Scenario) => s.id === id);
    if (!orig) return;
    const maxId = Math.max(...scenarios.map((s: Scenario) => s.id), 0);
    const newId = maxId + 1;
    const newScenario = { ...orig, id: newId };
    scenarios = [...scenarios, newScenario];
    selectedScenarioId = newId;
    compareScenarioId = null;
  }

  function deleteScenarioHandler(id: number): void {
    scenarios = scenarios.filter((s: Scenario) => s.id !== id);
    if (selectedScenarioId === id) selectedScenarioId = scenarios.length ? scenarios[0].id : null;
    if (compareScenarioId === id) compareScenarioId = null;
  }

  function renameScenarioHandler({ id, name }: { id: number; name: string }): void {
    scenarios = scenarios.map((s: Scenario) => s.id === id ? { ...s, name } : s);
  }

  function selectScenario(id: number): void {
    if (selectedScenarioId === id) {
      selectedScenarioId = null;
      setTimeout(() => {
        selectedScenarioId = id;
      }, 0);
    } else {
      selectedScenarioId = id;
    }
    compareScenarioId = null;
  }

  function selectCompare(id: number): void {
    compareScenarioId = id;
  }

  // Function no longer needed as we're using two-way binding with bind:scenario

  async function saveScenarios(): Promise<void> {
    if (!filename) {
      statusMessage = "Please enter a filename.";
      return;
    }

    try {
      isLoading = true;
      // Make a deep copy to avoid reference issues
      const scenariosList = JSON.parse(JSON.stringify(scenarios));
      console.log('Saving scenarios payload', JSON.stringify(scenariosList, null, 2));
      // Make sure filename has .json extension
      let filenameToSave = filename;
      if (!filenameToSave.endsWith('.json')) {
        filenameToSave += '.json';
      }
      // Add additional logs for debugging
      const selectedScenario = scenarios.find(s => s.id === selectedScenarioId);
      console.log('Scenarios content before saving:', scenariosList);
      console.log('Selected scenario data before saving:', selectedScenario?.data);
      console.log('Pension data of selected scenario:', selectedScenario?.data?.pension);
      const result = await api.saveScenarios({
        scenarios: scenariosList,
        filename: filenameToSave
      });
      if (result.success) {
        statusMessage = result.message;
        setTimeout(() => {
          showSaveDialog = false;
          statusMessage = "";
        }, 2000);
      } else {
        statusMessage = `Error: ${result.message}`;
      }
    } catch (error) {
      console.error("Error saving scenarios:", error);
      statusMessage = "Error saving scenarios. Please try again.";
    } finally {
      isLoading = false;
    }
  }

  async function loadScenarios(): Promise<void> {
    if (!selectedFile) {
      statusMessage = "Please select a file to load.";
      return;
    }

    showLoadDialog = true;
    statusMessage = "";
    isLoading = true;
    console.log('Loading scenarios from', selectedFile);
    try {
      const response = await api.loadScenarios({ filename: selectedFile });
      console.log('loadScenarios response', response);
      // handle Wails `result` wrapper
      const payload = response?.result ?? response;
      if (payload.success === false) {
        statusMessage = `Error: ${payload.message}`;
        return;
      }
      // get scenarios array
      const list: Scenario[] = payload.scenarios ?? payload;
      if (!Array.isArray(list)) {
        console.warn('Unexpected payload.scenarios', payload);
        statusMessage = 'Error: invalid data';
        return;
      }
      // update central appData store with loaded scenarios
      // Using a deep clone to ensure references are broken
      const deepClonedList = JSON.parse(JSON.stringify(list));
      
      // Validate and fix scenario IDs to ensure uniqueness
      const idMap = new Map();
      deepClonedList.forEach((scenario: Scenario) => {
        // If ID is missing, zero, or already used, assign a new one
        if (!scenario.id || scenario.id === 0 || idMap.has(scenario.id)) {
          // Find an unused ID
          let newId = 1;
          while (idMap.has(newId)) {
            newId++;
          }
          console.log(`Fixing duplicate/invalid scenario ID. Old: ${scenario.id}, New: ${newId}`);
          scenario.id = newId;
        }
        // Mark this ID as used
        idMap.set(scenario.id, true);
      });
      
      console.log('Loading scenarios with validated IDs:', deepClonedList);
      
      // Validate and ensure all required data sections exist in each scenario
      if (deepClonedList.length > 0) {
        // Check each scenario for possible issues and ensure all required data sections exist
        deepClonedList.forEach((scenario: Scenario) => {
          if (!scenario.data) {
            console.warn(`Scenario ${scenario.id} (${scenario.name}) has no data property`);
            scenario.data = {
  pension: {
    system: 'FERS',
    high3Salary: 100000,
    yearsOfService: 30,
    ageAtRetirement: 62,
    unusedSickLeaveMonths: 6,
    survivorBenefitOption: 'full',
    isPartTime: false,
    partTimeProrationFactor: 1.0,
    militaryService: 0
  },
  socialSecurity: {
    startAge: 62,
    estimatedMonthlyBenefit: 2000,
    isEligible: true,
    birthYear: 1970,
    birthMonth: 1
  },
  tsp: {
    traditionalBalance: 400000,
    rothBalance: 100000,
    contributionRate: 5,
    contributionRateRoth: 5,
    expectedReturn: 6,
    withdrawalStrategy: 'fixed',
    fixedMonthlyWithdrawal: 2000,
    withdrawalRate: 4,
    withdrawalStartAge: 62
  },
  tax: {
    filingStatus: 'married_joint',
    stateOfResidence: 'VA',
    stateIncomeTaxRate: 0.05,
    itemizedDeductions: 0,
    federalTaxCredits: 0,
    stateTaxCredits: 0,
    age: 62,
    spouseAge: 62
  },
  cola: {
    assumedInflationRate: 2.5,
    applyCOLAToPension: true,
    applyColaToSocialSecurity: true
  },
  otherIncome: {
    sources: [
      {
        id: crypto.randomUUID(),
        name: 'Part-time work',
        amount: 1200,
        frequency: 'monthly',
        startAge: 62,
        endAge: 70,
        applyCola: true
      }
    ]
  }
};
          }
          
          // Ensure all data sections exist with defaults
          if (!scenario.data.pension) {
            console.warn(`Scenario ${scenario.id} (${scenario.name}) is missing pension data - adding defaults`);
            scenario.data.pension = {
              system: 'FERS',
              high3Salary: 100000,
              yearsOfService: 30,
              ageAtRetirement: 62,
              unusedSickLeaveMonths: 6,
              survivorBenefitOption: 'full',
              isPartTime: false,
              partTimeProrationFactor: 1.0,
              militaryService: 0
            };
          }
          
          if (!scenario.data.socialSecurity) {
            console.warn(`Scenario ${scenario.id} (${scenario.name}) is missing socialSecurity data - adding defaults`);
            scenario.data.socialSecurity = {
              startAge: 62,
              estimatedMonthlyBenefit: 2000,
              isEligible: true,
              birthYear: 1970,
              birthMonth: 1
            };
          }
          
          if (!scenario.data.tsp) {
            console.warn(`Scenario ${scenario.id} (${scenario.name}) is missing tsp data - adding defaults`);
            scenario.data.tsp = {
              traditionalBalance: 400000,
              rothBalance: 100000,
              contributionRate: 5,
              contributionRateRoth: 5,
              expectedReturn: 6,
              withdrawalStrategy: 'fixed',
              fixedMonthlyWithdrawal: 2000,
              withdrawalRate: 4,
              withdrawalStartAge: 62
            };
          }
          
          if (!scenario.data.tax) {
            console.warn(`Scenario ${scenario.id} (${scenario.name}) is missing tax data - adding defaults`);
            scenario.data.tax = {
              filingStatus: 'married_joint',
              stateOfResidence: 'VA',
              stateIncomeTaxRate: 0.05,
              itemizedDeductions: 0,
              federalTaxCredits: 0,
              stateTaxCredits: 0,
              age: 62,
              spouseAge: 62
            };
          }
          
          if (!scenario.data.cola) {
            console.warn(`Scenario ${scenario.id} (${scenario.name}) is missing cola data - adding defaults`);
            scenario.data.cola = {
              assumedInflationRate: 2.5,
              applyCOLAToPension: true,
              applyColaToSocialSecurity: true
            };
          }
          
          if (!scenario.data.otherIncome) {
            console.warn(`Scenario ${scenario.id} (${scenario.name}) is missing otherIncome data - adding defaults`);
            scenario.data.otherIncome = {
              sources: [
                {
                  id: crypto.randomUUID(),
                  name: 'Part-time work',
                  amount: 1200,
                  frequency: 'monthly',
                  startAge: 62,
                  endAge: 70,
                  applyCola: true
                }
              ]
            };
          }
        });
      }
      
      // Legacy store logic removed. State is managed via local variables only.
      // If needed, force selection to trigger reactivity:
      selectedScenarioId = null;
      compareScenarioId = null;
      
      // Update scenarios with the loaded data
      scenarios = deepClonedList;
      
      // Then set the first scenario as selected (if any) with a longer delay
      // to ensure the store is updated with the new scenarios first
      if (deepClonedList.length > 0) {
        console.log('Setting first scenario after load to:', deepClonedList[0].id);
        setTimeout(() => {
          selectedScenarioId = deepClonedList[0].id;
        }, 100);
      }
      statusMessage = payload.message || 'Loaded successfully.';
      setTimeout(() => {
        showLoadDialog = false;
        statusMessage = "";
      }, 2000);
    } catch (error) {
      console.error('Error loading scenarios:', error);
      statusMessage = 'Error loading scenarios. Please try again.';
    } finally {
      isLoading = false;
    }
  }

  // Show the Save dialog
  function handleSave(): void {
    showSaveDialog = true;
  }

  // Show the Load dialog
  function handleLoad(): void {
    showLoadDialog = true;
  }

  // Handle global calculation request
  function handleGlobalCalculate(): void {
    console.log('App.handleGlobalCalculate triggered');
    
    if (!selectedScenarioId) {
      console.warn('No scenario selected, cannot calculate');
      return;
    }
    
    // Send an event to ScenarioTabs to trigger calculation in all sections
    // We'll create a custom event for this that will be captured by ScenarioTabs
    const calculationEvent = new CustomEvent('global-calculate', {
      detail: { scenarioId: selectedScenarioId }
    });
    
    // Dispatch the event to the document so it bubbles up
    document.dispatchEvent(calculationEvent);
    
    console.log('Global calculation event dispatched');
  }
</script>

<div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100" id="main-container">
  <Header 
    onSave={handleSave}
    onLoad={handleLoad}
    onCalculate={handleGlobalCalculate}
  />
  <div class="flex overflow-hidden">
    <ScenarioSidebar
      scenarios={scenarios}
      selectedScenarioId={selectedScenarioId}
      compareScenarioId={compareScenarioId}
      onRename={renameScenarioHandler}
      onSelect={selectScenario}
      onDuplicate={duplicateScenarioHandler}
      onDelete={deleteScenarioHandler}
      onAdd={addScenarioHandler}
      onCompare={selectCompare}
    />
    <main class="flex-1 p-6 overflow-auto">
      {#if compareScenarioId}
        <CompareView
          selectedScenario={scenarios.find((s: Scenario) => s.id === selectedScenarioId) ?? null}
          compareScenario={scenarios.find((s: Scenario) => s.id === compareScenarioId) ?? null}
        />
      {:else if selectedScenarioId}
        {#if selectedScenarioIndex >= 0 && selectedScenario}
          <ScenarioTabs
            bind:scenario={selectedScenario}
            on:scenarioChange={e => updateScenario(e.detail)}
          />
        {/if}
      {:else}
        <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow text-center">
          <p class="text-gray-700 dark:text-gray-300">No scenarios defined.</p>
        </div>
      {/if}
    </main>
  </div>
  
  <!-- Save Dialog -->
  {#if showSaveDialog}
    <div class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 w-96">
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Save Scenarios</h3>
        <div class="mb-4">
          <label for="filename" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Filename</label>
          <input 
            id="filename"
            bind:value={filename}
            type="text"
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md 
                   bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            placeholder="my_scenarios"
          />
        </div>
        
        {#if statusMessage}
          <div class="mb-4 p-2 rounded-md {statusMessage.includes('Error') ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'}">
            {statusMessage}
          </div>
        {/if}
        
        <div class="flex justify-end gap-3">
          <button 
            onclick={() => { showSaveDialog = false; statusMessage = ""; }}
            class="px-4 py-2 text-sm font-medium rounded-md bg-gray-100 hover:bg-gray-200 
                   dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200"
            disabled={isLoading}
          >
            Cancel
          </button>
          <button 
            onclick={saveScenarios}
            class="px-4 py-2 text-sm font-medium rounded-md bg-primary-600 hover:bg-primary-700 
                   dark:bg-primary-700 dark:hover:bg-primary-600 text-white"
            disabled={isLoading}
          >
            {isLoading ? 'Saving...' : 'Save'}
          </button>
        </div>
      </div>
    </div>
  {/if}
  
  <!-- Load Dialog -->
  {#if showLoadDialog}
    <div class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 w-96">
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">Load Scenarios</h3>
        
        {#if isLoading && savedFiles.length === 0}
          <div class="mb-4 text-center">
            <p class="text-gray-700 dark:text-gray-300">Loading saved files...</p>
          </div>
        {:else if savedFiles.length === 0}
          <div class="mb-4 text-center">
            <p class="text-gray-700 dark:text-gray-300">No saved files found.</p>
          </div>
        {:else}
          <div class="mb-4">
            <label for="loadFile" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Select File</label>
            <select 
              id="loadFile"
              bind:value={selectedFile}
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md 
                     bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            >
              {#each savedFiles as file}
                <option value={file}>{file}</option>
              {/each}
            </select>
          </div>
        {/if}
        
        {#if statusMessage}
          <div class="mb-4 p-2 rounded-md {statusMessage.includes('Error') ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'}">
            {statusMessage}
          </div>
        {/if}
        
        <div class="flex justify-end gap-3">
          <button 
            onclick={() => { showLoadDialog = false; statusMessage = ""; }}
            class="px-4 py-2 text-sm font-medium rounded-md bg-gray-100 hover:bg-gray-200 
                   dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200"
            disabled={isLoading}
          >
            Cancel
          </button>
          <button 
            onclick={loadScenarios}
            class="px-4 py-2 text-sm font-medium rounded-md bg-primary-600 hover:bg-primary-700 
                   dark:bg-primary-700 dark:hover:bg-primary-600 text-white"
            disabled={isLoading || savedFiles.length === 0}
          >
            {isLoading ? 'Loading...' : 'Load'}
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>
