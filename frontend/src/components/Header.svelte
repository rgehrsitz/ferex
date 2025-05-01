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
    // Apply to both document.documentElement AND main-container for maximum compatibility
    if (darkMode) {
      document.documentElement.classList.add('dark');
      document.body.classList.add('dark');
      // Also apply to main container as fallback (may be needed in Wails environment)
      const mainContainer = document.getElementById('main-container');
      if (mainContainer) mainContainer.classList.add('dark');
      console.log('Added dark class', darkMode);
    } else {
      document.documentElement.classList.remove('dark');
      document.body.classList.remove('dark');
      // Also remove from main container
      const mainContainer = document.getElementById('main-container');
      if (mainContainer) mainContainer.classList.remove('dark');
      console.log('Removed dark class', darkMode);
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
