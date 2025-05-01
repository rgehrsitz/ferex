<script lang="ts">
  export let data: any = {
    inflationRate: 2.5,
    useHistoricalCOLA: true,
    customCOLARate: 2.0,
    socialSecurityCOLA: 2.3,
    fersCOLA: 2.0,
    csrsCOLA: 2.5,
    applyCOLACap: true,
    applyFERSAge62Rule: true
  };

  const pensionOptions = [
    { name: 'FERS COLA', description: 'Follows inflation with caps: 2% if inflation is 2-3%, inflation-1% for inflation >3%' },
    { name: 'CSRS COLA', description: 'Matches full inflation rate with no caps' },
    { name: 'Social Security COLA', description: 'Based on CPI-W, usually similar to general inflation' }
  ];

  function calculateEstimatedFirstYearCOLA(inflationRate: number) {
    // Calculate FERS COLA based on inflation
    let fersCOLA = 0;
    
    if (inflationRate <= 2) {
      fersCOLA = inflationRate;
    } else if (inflationRate <= 3) {
      fersCOLA = 2;
    } else {
      fersCOLA = inflationRate - 1;
    }
    
    // Return FERS COLA, CSRS COLA (which is the full inflation rate),
    // and approximate SS COLA (using inflation as proxy)
    return {
      fers: data.applyCOLACap ? fersCOLA : inflationRate,
      csrs: inflationRate,
      socialSecurity: inflationRate
    };
  }

  $: estCOLA = calculateEstimatedFirstYearCOLA(data.inflationRate);
</script>

<div>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="inflationRate">
          Projected Inflation Rate
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="inflationRate"
            type="number"
            min="0"
            max="10"
            step="0.1"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.inflationRate}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>

      <div class="flex items-center">
        <input
          id="useHistoricalCOLA"
          type="checkbox"
          class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          bind:checked={data.useHistoricalCOLA}
        />
        <label for="useHistoricalCOLA" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
          Use calculated COLA rates based on projected inflation
        </label>
      </div>

      {#if !data.useHistoricalCOLA}
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="customCOLARate">
            Custom COLA Rate (for all income sources)
          </label>
          <div class="relative rounded-md shadow-sm">
            <input
              id="customCOLARate"
              type="number"
              min="0"
              max="10"
              step="0.1"
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={data.customCOLARate}
            />
            <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
              <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
            </div>
          </div>
        </div>
      {/if}
    </div>

    <div class="space-y-4">
      <div class="flex items-center">
        <input
          id="applyCOLACap"
          type="checkbox"
          class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          bind:checked={data.applyCOLACap}
        />
        <label for="applyCOLACap" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
          Apply FERS COLA cap
        </label>
      </div>
      
      <div class="flex items-center">
        <input
          id="applyFERSAge62Rule"
          type="checkbox"
          class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          bind:checked={data.applyFERSAge62Rule}
        />
        <label for="applyFERSAge62Rule" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
          Apply FERS age 62 COLA rule
        </label>
        <div class="ml-2 text-xs text-gray-500 dark:text-gray-400">
          (Most FERS retirees don't receive COLAs until age 62)
        </div>
      </div>

      <div class="bg-blue-50 dark:bg-blue-900 p-3 rounded-md mt-2">
        <h4 class="text-sm font-medium text-blue-800 dark:text-blue-300">COLA Rules</h4>
        <ul class="mt-2 text-xs text-blue-700 dark:text-blue-400 space-y-1 list-disc pl-4">
          <li>FERS: If inflation is 2% or less, COLA = Inflation</li>
          <li>FERS: If inflation is >2% but ≤3%, COLA = 2%</li>
          <li>FERS: If inflation is >3%, COLA = Inflation - 1%</li>
          <li>CSRS: COLA = Full Inflation (no caps)</li>
          <li>SS: COLA = Based on CPI-W (typically close to inflation)</li>
        </ul>
      </div>
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Estimated First-Year COLA</h3>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">FERS COLA</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{data.useHistoricalCOLA ? estCOLA.fers.toFixed(1) : data.customCOLARate.toFixed(1)}%</div>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">CSRS COLA</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{data.useHistoricalCOLA ? estCOLA.csrs.toFixed(1) : data.customCOLARate.toFixed(1)}%</div>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Social Security COLA</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{data.useHistoricalCOLA ? estCOLA.socialSecurity.toFixed(1) : data.customCOLARate.toFixed(1)}%</div>
      </div>
    </div>
  </div>
</div>