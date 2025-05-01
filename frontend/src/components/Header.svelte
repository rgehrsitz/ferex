<script lang="ts">
  import { onMount } from 'svelte';
  let darkMode = false;

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
</script>

<header class="flex items-center justify-between px-6 py-3 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 shadow-sm">
  <div class="flex items-center gap-2">
    <span class="font-bold text-lg text-primary-600 dark:text-primary-400">FEREX</span>
    <span class="text-xs text-gray-600 dark:text-gray-300 ml-2">Federal Employee Retirement Explorer</span>
  </div>
  <div class="flex items-center gap-4">
    <!-- Placeholder for user info, notifications, etc. -->
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