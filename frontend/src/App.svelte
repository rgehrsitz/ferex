<script lang="ts">
  import Header from './components/Header.svelte';
  import ScenarioSidebar from './components/ScenarioSidebar.svelte';
  import ScenarioTabs from './components/ScenarioTabs.svelte';
  import CompareView from './components/CompareView.svelte';
  import type { Scenario, ScenarioData } from './types/scenario.js';
  
  // Import stores
  import { 
    scenarios, 
    selectedScenarioId, 
    compareScenarioId,
    selectedScenario,
    compareScenario,
    createDefaultScenario,
    initializeStore,
    createUniqueId
  } from './stores/scenarioStore';
  
  import { api } from './stores/apiStore';
  import { syncWithScenario } from './stores/userDataStore';
  
  // Initialize the store with default scenarios
  initializeStore();
  
  // Create derivation to sync scenario data with user store
  $: { 
    // This will update the userData store whenever the selected scenario changes
    $syncWithScenario; 
  }
  
  function selectScenario(id: number): void { 
    selectedScenarioId.set(id); 
    compareScenarioId.set(null); 
  }
  
  function selectCompare(id: number): void { 
    compareScenarioId.set(id); 
  }
  
  // State variables for file management dialogs
  let showSaveDialog = false;
  let showLoadDialog = false;
  let filename = "my_scenarios";
  let savedFiles: string[] = [];
  let selectedFile = "";
  let statusMessage = "";
  let isLoading = false;

  function addScenario(): void {
    let scenariosList;
    scenarios.subscribe(list => scenariosList = list)();
    
    const id = createUniqueId(scenariosList);
    const newScenario = createDefaultScenario(id, `Scenario ${String.fromCharCode(64 + id)}`);
    
    // Update the store
    scenarios.update(list => [...list, newScenario]);
    selectedScenarioId.set(id);
    compareScenarioId.set(null);
  }
  
  function duplicateScenario(id: number): void {
    let scenariosList;
    scenarios.subscribe(list => scenariosList = list)();
    
    const orig = scenariosList.find(s => s.id === id);
    if (!orig) return;
    
    const newId = createUniqueId(scenariosList);
    
    // Create deep clone to ensure all nested objects are properly copied
    const clonedData = JSON.parse(JSON.stringify(orig.data));
    
    // Make sure otherIncome sources have unique IDs in the duplicate
    if (clonedData.otherIncome?.sources?.length > 0) {
      clonedData.otherIncome.sources = clonedData.otherIncome.sources.map(source => ({
        ...source,
        id: crypto.randomUUID()
      }));
    }
    
    // Update the store
    scenarios.update(list => [...list, { id: newId, name: orig.name + " Copy", data: clonedData }]);
  }
  
  function deleteScenario(id: number): void {
    let scenariosList;
    scenarios.subscribe(list => scenariosList = list)();
    
    // Update the store
    scenarios.update(list => list.filter(s => s.id !== id));
    
    // Update selected ID if needed
    if ($selectedScenarioId === id) {
      const remainingScenarios = scenariosList.filter(s => s.id !== id);
      selectedScenarioId.set(remainingScenarios[0]?.id || null);
    }
    
    // Update compare ID if needed
    if ($compareScenarioId === id) {
      compareScenarioId.set(null);
    }
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
      const result = await api.getSavedScenarios();
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
      let scenariosList;
      scenarios.subscribe(list => scenariosList = list)();
      
      // Make sure filename has .json extension
      let filenameToSave = filename;
      if (!filenameToSave.endsWith('.json')) {
        filenameToSave += '.json';
      }
      
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
  
  async function loadScenarios() {
    if (!selectedFile) {
      statusMessage = "Please select a file to load.";
      return;
    }
    
    try {
      isLoading = true;
      const result = await api.loadScenarios({
        filename: selectedFile
      });
      
      if (result.success) {
        // Update the store with loaded scenarios
        scenarios.set(result.scenarios);
        
        // Set first scenario as selected
        if (result.scenarios.length > 0) {
          selectedScenarioId.set(result.scenarios[0].id);
        } else {
          selectedScenarioId.set(null);
        }
        
        // Clear compare scenario
        compareScenarioId.set(null);
        
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
  
  // We don't need these anymore as they're derived directly in the store
  // We use $selectedScenario and $compareScenario in the template
</script>

<div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100" id="main-container">
  <Header 
    on:save={handleSave}
    on:load={handleLoad}
  />
  <div class="flex overflow-hidden">
    <ScenarioSidebar
      scenarios={$scenarios}
      selectedScenarioId={$selectedScenarioId}
      compareScenarioId={$compareScenarioId}
      on:add={addScenario}
      on:duplicate={e => duplicateScenario(e.detail)}
      on:delete={e => deleteScenario(e.detail)}
      on:select={e => selectScenario(e.detail)}
      on:compare={e => selectCompare(e.detail)}
    />
    <main class="flex-1 p-6 overflow-auto">
      {#if $compareScenario}
        <CompareView selectedScenario={$selectedScenario} compareScenario={$compareScenario} />
      {:else if $selectedScenario}
        <ScenarioTabs bind:scenario={$selectedScenario} />
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
