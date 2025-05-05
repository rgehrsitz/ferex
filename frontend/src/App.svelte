<script lang="ts">
  import Header from './components/Header.svelte';
  import ScenarioSidebar from './components/ScenarioSidebar.svelte';
  import ScenarioTabs from './components/ScenarioTabs.svelte';
  import CompareView from './components/CompareView.svelte';
  import type { Scenario, ScenarioData } from './types/scenario.js';

  // Import consolidated stores
  import {
    scenarios,
    selectedScenario,
    compareScenario,
    addScenario,
    deleteScenario as deleteScenarioStore,
    createDefaultScenario,
    initializeStore,
    updateScenario as updateScenarioStore
  } from './stores/scenarioStore.js';
  import { selectedScenarioId, compareScenarioId } from './stores/uiStore.js';
  import { api } from './stores/apiStore.js';
  import { appData } from './stores/appDataStore.js';
  import { updateUserProfile, getUserProfile } from './stores/userDataStore.js';

  // Debug: log scenarios store whenever it changes
  $: console.log('App scenarios store updated', $scenarios);

  // Initialize the store with default scenarios
  initializeStore();

  function addScenarioHandler(): void {
    const id = addScenario(createDefaultScenario(undefined, undefined));
    selectedScenarioId.set(id);
    compareScenarioId.set(null);
  }

  function duplicateScenarioHandler(id: number): void {
    let orig;
    scenarios.subscribe(list => orig = list.find(s => s.id === id))();
    if (!orig) return;
    const clone = createDefaultScenario(undefined, `${orig.name} Copy`);
    clone.data = JSON.parse(JSON.stringify(orig.data));
    const newId = addScenario(clone);
    selectedScenarioId.set(newId);
  }

  function deleteScenarioHandler(id: number): void {
    deleteScenarioStore(id);
  }

  async function saveScenarios() {
    if (!filename) {
      statusMessage = "Please enter a filename.";
      return;
    }

    try {
      isLoading = true;
      
      // Get a snapshot of the current scenarios from the store
      let scenariosList;
      scenarios.subscribe(list => {
        // Make a deep copy to avoid reference issues
        scenariosList = JSON.parse(JSON.stringify(list));
      })();
      
      console.log('Saving scenarios payload', JSON.stringify(scenariosList, null, 2));
      
      // Make sure filename has .json extension
      let filenameToSave = filename;
      if (!filenameToSave.endsWith('.json')) {
        filenameToSave += '.json';
      }
      
      // Add additional logs for debugging
      console.log('Scenarios store content before saving:', $scenarios);
      console.log('Selected scenario data before saving:', $selectedScenario?.data);
      console.log('Pension data of selected scenario:', $selectedScenario?.data?.pension);
      
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
      const list = payload.scenarios ?? payload;
      if (!Array.isArray(list)) {
        console.warn('Unexpected payload.scenarios', payload);
        statusMessage = 'Error: invalid data';
        return;
      }
      // update central appData store with loaded scenarios
      // Using a deep clone to ensure references are broken
      const deepClonedList = JSON.parse(JSON.stringify(list));
      console.log('Loading scenarios with deep cloned data:', deepClonedList);
      
      appData.update(d => ({ ...d, scenarios: deepClonedList }));
      
      // Force selection to trigger reactivity
      selectedScenarioId.set(null);
      compareScenarioId.set(null);
      
      // Then set the first scenario as selected (if any)
      if (list.length > 0) {
        setTimeout(() => {
          selectedScenarioId.set(list[0].id);
        }, 0);
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
    statusMessage = "";
  }

  // Show the Load dialog and fetch saved file list
  async function handleLoad(): Promise<void> {
    showLoadDialog = true;
    statusMessage = "";
    isLoading = true;
    try {
      const response = await api.getSavedScenarios();
      console.log('getSavedScenarios response', response);
      // Wails bindings often wrap data in a `result` field
      const payload = response?.result ?? response;
      if (Array.isArray(payload)) {
        savedFiles = payload;
      } else if (Array.isArray(payload.files)) {
        savedFiles = payload.files;
      } else {
        console.warn('Unexpected payload for saved files', payload);
        savedFiles = [];
      }
    } catch (err) {
      console.error('Error fetching saved files:', err);
      statusMessage = 'Error fetching saved files.';
    } finally {
      isLoading = false;
    }
  }

  // State variables for file management dialogs
  let showSaveDialog = false;
  let showLoadDialog = false;
  let filename = "my_scenarios";
  let savedFiles: string[] = [];
  let selectedFile = "";
  let statusMessage = "";
  let isLoading = false;

  function selectScenario(id: number): void { 
    selectedScenarioId.set(id); 
    compareScenarioId.set(null); 
  }

  function selectCompare(id: number): void { 
    compareScenarioId.set(id); 
  }

  // Handle scenario updates from child sections
  function handleUpdateScenario(event): void {
    console.log('App.handleUpdateScenario received', event.detail);
    
    // Make sure we're getting valid pension data
    if (event.detail?.data?.pension) {
      console.log('Pension data being saved:', event.detail.data.pension);
    }
    
    // Ensure we're passing a complete scenario object and not just a reference
    const updatedScenario = JSON.parse(JSON.stringify(event.detail));
    console.log('App updating scenario store with:', updatedScenario);
    updateScenarioStore(updatedScenario);
  }
  
  // Handle global calculation request
  function handleGlobalCalculate(): void {
    console.log('App.handleGlobalCalculate triggered');
    
    if (!$selectedScenario) {
      console.warn('No scenario selected, cannot calculate');
      return;
    }
    
    // Send an event to ScenarioTabs to trigger calculation in all sections
    // We'll create a custom event for this that will be captured by ScenarioTabs
    const calculationEvent = new CustomEvent('global-calculate', {
      detail: { scenarioId: $selectedScenario.id }
    });
    
    // Dispatch the event to the document so it bubbles up
    document.dispatchEvent(calculationEvent);
    
    console.log('Global calculation event dispatched');
  }
</script>

<div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100" id="main-container">
  <Header 
    on:save={handleSave}
    on:load={handleLoad}
    on:calculate={handleGlobalCalculate}
  />
  <div class="flex overflow-hidden">
    <ScenarioSidebar
      scenarios={$scenarios}
      selectedScenarioId={$selectedScenarioId}
      compareScenarioId={$compareScenarioId}
      on:add={addScenarioHandler}
      on:duplicate={e => duplicateScenarioHandler(e.detail)}
      on:delete={e => deleteScenarioHandler(e.detail)}
      on:select={e => selectScenario(e.detail)}
      on:compare={e => selectCompare(e.detail)}
    />
    <main class="flex-1 p-6 overflow-auto">
      {#if $compareScenario}
        <CompareView selectedScenario={$selectedScenario} compareScenario={$compareScenario} />
      {:else if $selectedScenario}
        <ScenarioTabs scenario={$selectedScenario} on:update-scenario={handleUpdateScenario} />
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
