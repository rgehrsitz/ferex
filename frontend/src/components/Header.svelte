<script lang="ts">
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();
  let darkMode: boolean = false;
  let isCalculating: boolean = false;

  // On mount, check localStorage or system preference
  onMount(() => {
    console.log('[DarkMode] onMount');
    if (localStorage.getItem('darkMode')) {
      darkMode = localStorage.getItem('darkMode') === 'true';
    } else {
      darkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    // Force dark mode initially - remove this line once theme toggle works
    darkMode = true;
    updateHtmlClass();
    setTimeout(() => {
      updateHtmlClass();
      console.log('[DarkMode] Forced update after mount', darkMode);
    }, 100);
  });

  function toggleDarkMode() {
    darkMode = !darkMode;
    console.log('Toggle dark mode:', darkMode);
    localStorage.setItem('darkMode', darkMode.toString());
    updateHtmlClass();
  }

  function updateHtmlClass() {
    try {
      console.log('[DarkMode] Updating HTML class with darkMode =', darkMode);
      
      // Apply to document.documentElement
      if (darkMode) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
      
      // Apply to document.body
      if (darkMode) {
        document.body.classList.add('dark');
      } else {
        document.body.classList.remove('dark');
      }
      
      // Apply to main container
      const mainContainer = document.getElementById('main-container');
      if (mainContainer) {
        if (darkMode) {
          mainContainer.classList.add('dark');
        } else {
          mainContainer.classList.remove('dark');
        }
      } else {
        console.warn('[DarkMode] Main container not found');
      }
      
      // Debug current class state
      console.log('[DarkMode] documentElement classList:', document.documentElement.classList);
      console.log('[DarkMode] body classList:', document.body.classList);
      if (mainContainer) {
        console.log('[DarkMode] mainContainer classList:', mainContainer.classList);
      }
      
    } catch (error) {
      console.error('[DarkMode] Error updating HTML class:', error);
    }
  }
  
  function handleSaveClick() {
    dispatch('save');
  }
  
  function handleLoadClick() {
    dispatch('load');
  }
  
  function handleCalculateClick() {
    isCalculating = true;
    dispatch('calculate');
    // Reset the calculating state after a short delay
    setTimeout(() => {
      isCalculating = false;
    }, 1000);
  }
</script>

<header class="flex items-center justify-between px-6 py-3 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 shadow-sm">
  <div class="flex items-center gap-2">
    <span class="font-bold text-lg text-primary-600 dark:text-primary-400">FEREX</span>
    <span class="text-xs text-gray-600 dark:text-gray-300 ml-2">Federal Employee Retirement Explorer</span>
  </div>
  <div class="flex items-center gap-3">
    <!-- Save and Load buttons -->
    <button 
      on:click={handleSaveClick} 
      class="px-3 py-1 rounded-md text-sm font-medium 
      bg-green-100 hover:bg-green-200 dark:bg-green-800 dark:hover:bg-green-700 
      text-green-800 dark:text-green-200 
      border border-green-300 dark:border-green-600 
      transition-colors duration-200 flex items-center gap-1"
    >
      <span>💾</span>
      <span>Save</span>
    </button>
    
    <button 
      on:click={handleLoadClick} 
      class="px-3 py-1 rounded-md text-sm font-medium 
      bg-blue-100 hover:bg-blue-200 dark:bg-blue-800 dark:hover:bg-blue-700 
      text-blue-800 dark:text-blue-200 
      border border-blue-300 dark:border-blue-600 
      transition-colors duration-200 flex items-center gap-1"
    >
      <span>📂</span>
      <span>Load</span>
    </button>
    
    <!-- Calculate Button -->
    <button 
      on:click={handleCalculateClick} 
      class="px-3 py-1 rounded-md text-sm font-medium 
      bg-purple-100 hover:bg-purple-200 dark:bg-purple-800 dark:hover:bg-purple-700 
      text-purple-800 dark:text-purple-200 
      border border-purple-300 dark:border-purple-600 
      transition-colors duration-200 flex items-center gap-1"
    >
      <span>📊</span>
      <span>{isCalculating ? 'Calculating...' : 'Calculate All'}</span>
    </button>
    
    <!-- Dark Mode Toggle -->
    <button 
      on:click={toggleDarkMode} 
      class="px-3 py-1 rounded-md text-sm font-medium 
      bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 
      text-gray-800 dark:text-gray-200 
      border border-gray-300 dark:border-gray-600 
      transition-colors duration-200 flex items-center gap-1"
    >
      <span>{darkMode ? '🌙' : '☀️'}</span>
      <span>{darkMode ? 'Dark' : 'Light'}</span>
    </button>
  </div>
</header>