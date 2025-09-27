<script lang="ts">
  let open = $state(false);
  let tab = $state('calculator');
  function togglePanel() { open = !open; }
</script>

<div class={`fixed top-0 right-0 h-full z-40 transition-transform duration-300 ${open ? 'translate-x-0' : 'translate-x-full'}`}
     style="width: 380px;">
  <div class="flex items-center justify-between h-14 px-4 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700">
    <span class="font-bold text-lg text-blue-900 dark:text-blue-200">Quick Panel</span>
    <button class="ml-auto text-gray-400 hover:text-blue-500" onclick={togglePanel} aria-label="Close panel">
      <svg width="24" height="24" fill="none" stroke="currentColor"><path d="M6 18L18 6M6 6l12 12"/></svg>
    </button>
  </div>
  <div class="flex flex-col h-full bg-white dark:bg-gray-900">
    <div class="flex gap-2 px-4 py-2 border-b border-gray-200 dark:border-gray-700">
      <button class={`px-3 py-1 rounded ${tab === 'calculator' ? 'bg-blue-100 dark:bg-blue-800' : ''}`} onclick={() => tab = 'calculator'}>Calculator</button>
      <button class={`px-3 py-1 rounded ${tab === 'recent' ? 'bg-blue-100 dark:bg-blue-800' : ''}`} onclick={() => tab = 'recent'}>Recent Scenarios</button>
    </div>
    <div class="flex-1 overflow-y-auto p-4">
      {#if tab === 'calculator'}
        <div>Calculator widget goes here</div>
      {:else if tab === 'recent'}
        <div>Recent scenarios list goes here</div>
      {/if}
    </div>
  </div>
</div>

<!-- Toggle button (can be placed anywhere in the app) -->
<button class="fixed top-1/2 right-0 z-50 bg-blue-600 text-white px-2 py-1 rounded-l shadow" onclick={togglePanel} aria-label="Open quick panel">
  <svg width="24" height="24" fill="none" stroke="currentColor"><path d="M4 12h16M16 8l4 4-4 4"/></svg>
</button>
