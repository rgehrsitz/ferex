<script>
  import App from './App.svelte';
  import { onMount } from 'svelte';

  let hasError = $state(false);
  let errorMessage = $state('');
  let mounted = $state(false);

  onMount(() => {
    console.log('AppWrapper mounted successfully');
    mounted = true;
  });

  function handleError(error) {
    console.error('App error caught:', error);
    hasError = true;
    errorMessage = error?.message || 'An unknown error occurred';
  }

  // Setup window error handler
  $effect(() => {
    const errorHandler = (event) => {
      console.error('Global error caught:', event.error);
      hasError = true;
      errorMessage = event.error?.message || event.message || 'An unknown error occurred';
      event.preventDefault();
    };

    window.addEventListener('error', errorHandler);
    window.addEventListener('unhandledrejection', (event) => {
      errorHandler({ error: event.reason, message: 'Unhandled Promise Rejection' });
    });

    return () => {
      window.removeEventListener('error', errorHandler);
      window.removeEventListener('unhandledrejection', errorHandler);
    };
  });
</script>

<div class="app-wrapper min-h-screen w-full flex flex-col bg-white dark:bg-gray-800">
  <div class="p-4 bg-blue-100 dark:bg-blue-900 mb-2">
    <h1 class="text-lg font-bold">Ferex Debug Mode</h1>
    <p>AppWrapper mounted: {mounted ? 'Yes' : 'No'}</p>
  </div>

  {#if hasError}
    <div class="p-4 bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200 rounded-md m-4">
      <h2 class="text-xl font-bold mb-2">Error</h2>
      <p>{errorMessage}</p>
      <pre class="mt-4 p-2 bg-gray-100 dark:bg-gray-800 overflow-auto text-sm">{errorMessage}</pre>
    </div>
  {:else}
    <div class="mx-auto w-full">
      <try>
        <App />
        <catch let:error>
          {handleError(error)}
          <p>Error rendering App component</p>
        </catch>
      </try>
    </div>
  {/if}
</div>