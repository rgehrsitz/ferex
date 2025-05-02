<script lang="ts">
  import Header from './components/Header.svelte';
  import ScenarioSidebar from './components/ScenarioSidebar.svelte';
  import ScenarioTabs from './components/ScenarioTabs.svelte';
  import CompareView from './components/CompareView.svelte';
  import type { Scenario, ScenarioData } from './types/scenario.js';
  import { SaveScenarios, LoadScenarios, GetSavedScenarios } from '../wailsjs/go/main/App';

  function selectScenario(id: number): void { selectedScenarioId = id; compareScenarioId = null; }
  function selectCompare(id: number): void { compareScenarioId = id; }
  
  // State variables for file management dialogs
  let showSaveDialog = false;
  let showLoadDialog = false;
  let filename = "my_scenarios";
  let savedFiles: string[] = [];
  let selectedFile = "";
  let statusMessage = "";
  let isLoading = false;

  // Default data for new scenario
  const getDefaultScenarioData = (): ScenarioData => ({
    pension: {
      retirementSystem: 'FERS',
      highThreeSalary: 100000,
      yearsOfService: 30,
      retirementAge: 62,
      unusedSickLeave: 1000,
      militaryService: 0,
      isPartTime: false,
      partTimeProrationFactor: 1.0,
      csrsOffset: false,
      survivorBenefit: 'full'
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
      additionalIncome: 0,
      additionalDeductions: 0,
      stateIncomeTaxRate: 0.05
    },
    cola: {
      assumedInflationRate: 0.025,
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
        },
        {
          id: crypto.randomUUID(),
          name: 'Rental property',
          amount: 1800,
          frequency: 'monthly',
          startAge: 62,
          endAge: 90,
          applyCola: true
        }
      ]
    }
  });

  let scenarios: Scenario[] = [
    { id: 1, name: "Scenario A", data: getDefaultScenarioData() },
    { id: 2, name: "Scenario B", data: getDefaultScenarioData() }
  ];
  let selectedScenarioId: number | null = 1;
  let compareScenarioId: number | null = null;

  function addScenario(): void {
    const id = Math.max(...scenarios.map(s => s.id)) + 1;
    scenarios = [...scenarios, { 
      id, 
      name: `Scenario ${String.fromCharCode(64 + id)}`, 
      data: getDefaultScenarioData() 
    }];
    selectedScenarioId = id;
    compareScenarioId = null;
  }
  function duplicateScenario(id: number): void {
    const orig = scenarios.find(s => s.id === id);
    if (!orig) return;
    const newId = Math.max(...scenarios.map(s => s.id)) + 1;
    
    // Create deep clone to ensure all nested objects are properly copied
    const clonedData = JSON.parse(JSON.stringify(orig.data));
    
    // Make sure otherIncome sources have unique IDs in the duplicate
    if (clonedData.otherIncome?.sources?.length > 0) {
      clonedData.otherIncome.sources = clonedData.otherIncome.sources.map(source => ({
        ...source,
        id: crypto.randomUUID()
      }));
    }
    
    scenarios = [...scenarios, { id: newId, name: orig.name + " Copy", data: clonedData }];
  }
  function deleteScenario(id: number): void {
    scenarios = scenarios.filter(s => s.id !== id);
    if (selectedScenarioId === id) selectedScenarioId = scenarios[0]?.id || null;
    if (compareScenarioId === id) compareScenarioId = null;
  }

  // Save and load functions
  async function handleSave() {
    showSaveDialog = true;
    statusMessage = "";
  }
  
  async function handleLoad() {
    showLoadDialog = true;
    statusMessage = "";
    
    try {
      isLoading = true;
      const result = await GetSavedScenarios();
      savedFiles = result.files;
      if (savedFiles.length > 0) {
        selectedFile = savedFiles[0];
      }
    } catch (error) {
      console.error("Error getting saved scenarios:", error);
      statusMessage = "Error loading saved files list.";
    } finally {
      isLoading = false;
    }
  }
  
  async function saveScenarios() {
    if (!filename) {
      statusMessage = "Please enter a filename.";
      return;
    }
    
    try {
      isLoading = true;
      // Make sure filename has .json extension
      let filenameToSave = filename;
      if (!filenameToSave.endsWith('.json')) {
        filenameToSave += '.json';
      }
      
      const result = await SaveScenarios({
        scenarios: scenarios,
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
  
  async function loadScenarios() {
    if (!selectedFile) {
      statusMessage = "Please select a file to load.";
      return;
    }
    
    try {
      isLoading = true;
      const result = await LoadScenarios({
        filename: selectedFile
      });
      
      if (result.success) {
        scenarios = result.scenarios;
        selectedScenarioId = scenarios[0]?.id || null;
        compareScenarioId = null;
        
        statusMessage = result.message;
        setTimeout(() => {
          showLoadDialog = false;
          statusMessage = "";
        }, 2000);
      } else {
        statusMessage = `Error: ${result.message}`;
      }
    } catch (error) {
      console.error("Error loading scenarios:", error);
      statusMessage = "Error loading scenarios. Please try again.";
    } finally {
      isLoading = false;
    }
  }
  
  let selectedScenario: Scenario | undefined;
  let compareScenario: Scenario | undefined;
  $: selectedScenario = scenarios.find(s => s.id === selectedScenarioId);
  $: compareScenario = compareScenarioId ? scenarios.find(s => s.id === compareScenarioId) : undefined;
</script>

<div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100" id="main-container">
  <Header 
    on:save={handleSave}
    on:load={handleLoad}
  />
  <div class="flex overflow-hidden">
    <ScenarioSidebar
      {scenarios}
      {selectedScenarioId}
      {compareScenarioId}
      on:add={addScenario}
      on:duplicate={e => duplicateScenario(e.detail)}
      on:delete={e => deleteScenario(e.detail)}
      on:select={e => selectScenario(e.detail)}
      on:compare={e => selectCompare(e.detail)}
    />
    <main class="flex-1 p-6 overflow-auto">
      {#if compareScenario}
        <CompareView {selectedScenario} {compareScenario} />
      {:else if selectedScenario}
        <ScenarioTabs bind:scenario={selectedScenario} />
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
            on:click={() => { showSaveDialog = false; statusMessage = ""; }}
            class="px-4 py-2 text-sm font-medium rounded-md bg-gray-100 hover:bg-gray-200 
                   dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200"
            disabled={isLoading}
          >
            Cancel
          </button>
          <button 
            on:click={saveScenarios}
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
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Select File</label>
            <select 
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
            on:click={() => { showLoadDialog = false; statusMessage = ""; }}
            class="px-4 py-2 text-sm font-medium rounded-md bg-gray-100 hover:bg-gray-200 
                   dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200"
            disabled={isLoading}
          >
            Cancel
          </button>
          <button 
            on:click={loadScenarios}
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
