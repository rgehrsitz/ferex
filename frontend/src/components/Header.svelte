<script lang="ts">
  import { onMount } from 'svelte';
  let darkMode = false;

  // On mount, check localStorage or system preference
  onMount(() => {
    if (localStorage.getItem('darkMode')) {
      darkMode = localStorage.getItem('darkMode') === 'true';
    } else {
      darkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    updateHtmlClass();
  });

  function toggleDarkMode() {
    darkMode = !darkMode;
    localStorage.setItem('darkMode', darkMode.toString());
    updateHtmlClass();
  }

  function updateHtmlClass() {
    if (darkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
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
