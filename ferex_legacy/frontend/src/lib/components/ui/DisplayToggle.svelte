<script lang="ts">
  import { displayMode } from '../../stores/displayPreferences';
  
  export let className = '';
</script>

<div class="flex items-center {className}">
  <span class="text-sm text-gray-600 dark:text-gray-400 pr-2">Display as:</span>
  <div class="relative inline-flex items-center">
    <button 
      class="px-3 py-1 text-sm font-medium rounded-l-md transition-colors {$displayMode === 'yearly' ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-800 dark:bg-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'}" 
      on:click={() => $displayMode = 'yearly'}
    >
      Yearly
    </button>
    <button 
      class="px-3 py-1 text-sm font-medium rounded-r-md transition-colors {$displayMode === 'monthly' ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-800 dark:bg-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'}" 
      on:click={() => $displayMode = 'monthly'}
    >
      Monthly
    </button>
  </div>
</div>
