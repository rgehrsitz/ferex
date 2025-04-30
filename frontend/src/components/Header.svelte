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
    <span class="font-bold text-lg text-gray-900 dark:text-gray-100">FEREX</span>
    <span class="text-xs text-gray-500 dark:text-gray-400 ml-2">Federal Employee Retirement Explorer</span>
  </div>
  <div class="flex items-center gap-4">
    <!-- Placeholder for user info, notifications, etc. -->
    <button on:click={toggleDarkMode} class="px-3 py-1 rounded bg-gray-200 dark:bg-gray-800 text-gray-800 dark:text-gray-100 border border-gray-300 dark:border-gray-700 hover:bg-gray-300 dark:hover:bg-gray-700 transition">
      {darkMode ? '🌙 Dark' : '☀️ Light'}
    </button>
  </div>
</header>
