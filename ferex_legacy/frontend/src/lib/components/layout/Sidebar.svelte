<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  const { activeView = "home" } = $props();
  const dispatch = createEventDispatcher();

  // Sidebar navigation items with SVG icons
  const navItems = [
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/></svg>', 
      label: 'Home', 
      route: 'home' 
    },
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"/></svg>', 
      label: 'Scenario Editor', 
      route: 'scenario' 
    },
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z"/></svg>', 
      label: 'Results Summary', 
      route: 'results' 
    },
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"/></svg>', 
      label: 'Income Projections', 
      route: 'income-projections' 
    },
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z"/><path fill-rule="evenodd" d="M4 5a2 2 0 012-2v1a1 1 0 102 0V3h4v1a1 1 0 102 0V3a2 2 0 012 2v6a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm8 8a1 1 0 100-2 1 1 0 000 2zm-3-1a1 1 0 11-2 0 1 1 0 012 0zm-1-4a1 1 0 100-2 1 1 0 000 2zm3-1a1 1 0 11-2 0 1 1 0 012 0z" clip-rule="evenodd"/></svg>', 
      label: 'Scenario Comparison', 
      route: 'scenario-comparison' 
    },
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/></svg>', 
      label: 'Risk Analysis', 
      route: 'risk-analysis' 
    },
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>', 
      label: 'Household Planning', 
      route: 'household' 
    },
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clip-rule="evenodd"/></svg>', 
      label: 'Reports & Exports', 
      route: 'reports' 
    },
    { 
      icon: '<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd"/></svg>', 
      label: 'Settings', 
      route: 'settings' 
    }
  ];
  let collapsed = $state(false);

  function toggleSidebar() {
    collapsed = !collapsed;
  }
  function handleNav(route: string) {
    dispatch('changeView', { view: route });
  }
</script>

<aside class={`h-full bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-700 flex flex-col transition-all duration-300 ${collapsed ? 'w-16' : 'w-56'}`}>
  <div class="flex items-center justify-between h-14 px-4 border-b border-gray-200 dark:border-gray-700">
    <span class="font-bold text-lg text-blue-900 dark:text-blue-200">{!collapsed ? 'FeReX' : 'F'}</span>
    <button class="ml-auto text-gray-400 hover:text-blue-500" onclick={toggleSidebar} aria-label="Toggle sidebar">
      {#if collapsed}
        <svg width="24" height="24" fill="none" stroke="currentColor"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
      {:else}
        <svg width="24" height="24" fill="none" stroke="currentColor"><path d="M19 12H5m7 7l-7-7 7-7"/></svg>
      {/if}
    </button>
  </div>
  <nav class="flex-1 py-4 space-y-1">
    {#each navItems as item}
      <button type="button"
        class={`w-full flex items-center px-4 py-2 rounded transition-colors ${collapsed ? 'justify-center' : ''} ${activeView === item.route ? 'bg-blue-600 text-white dark:bg-blue-400 dark:text-gray-900' : 'text-gray-700 dark:text-gray-200 hover:bg-blue-100 dark:hover:bg-blue-800'}`}
        title={item.label}
        onclick={() => handleNav(item.route)}>
        <span class="flex items-center justify-center">{@html item.icon}</span>
        {#if !collapsed}
          <span class="ml-3">{item.label}</span>
        {/if}
      </button>
    {/each}
  </nav>
</aside>
