<script lang="ts">
  import { CalculatePension } from '../wailsjs/go/main/App.js';
  let system = $state("FERS");
  let high3Salary = $state(0);
  let yearsOfService = $state(0);
  let ageAtRetirement = $state(0);
  let pensionResult = $state<any>(null);
  let error = $state("");

  function validateInputs() {
    if (high3Salary <= 0 || yearsOfService < 0 || ageAtRetirement <= 0) {
      error = "Please enter valid positive values for all fields.";
      return false;
    }
    error = "";
    return true;
  }

  function calculatePension() {
    if (!validateInputs()) return;
    CalculatePension({
      System: system,
      High3Salary: high3Salary,
      YearsOfService: yearsOfService,
      AgeAtRetirement: ageAtRetirement
    }).then(result => {
      pensionResult = result;
    });
  }
</script>

<main class="min-h-screen flex items-center justify-center bg-gray-900">
  <div class="bg-white rounded-xl shadow-lg p-8 w-full max-w-md">
    <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">Federal Pension Calculator</h1>
    <form class="space-y-4" on:submit|preventDefault={calculatePension}>
      <div>
        <label class="block text-gray-700 font-medium mb-1" for="system">System</label>
        <select bind:value={system} id="system" class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500">
          <option value="FERS">FERS</option>
          <option value="CSRS">CSRS</option>
        </select>
      </div>
      <div>
        <label class="block text-gray-700 font-medium mb-1" for="high3Salary">High 3 Salary</label>
        <input autocomplete="off" bind:value={high3Salary} id="high3Salary" type="number" min="0" class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500" />
      </div>
      <div>
        <label class="block text-gray-700 font-medium mb-1" for="yearsOfService">Years of Service</label>
        <input autocomplete="off" bind:value={yearsOfService} id="yearsOfService" type="number" min="0" class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500" />
      </div>
      <div>
        <label class="block text-gray-700 font-medium mb-1" for="ageAtRetirement">Age at Retirement</label>
        <input autocomplete="off" bind:value={ageAtRetirement} id="ageAtRetirement" type="number" min="0" class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500" />
      </div>
      {#if error}
        <div class="text-red-600 text-sm">{error}</div>
      {/if}
      <button type="submit" class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition">Calculate Pension</button>
    </form>
    {#if pensionResult}
      <div class="mt-6 bg-blue-50 rounded p-4 text-gray-800">
        <h2 class="font-semibold text-lg mb-2">Results</h2>
        <div><span class="font-medium">Annual Pension:</span> {pensionResult.annualPension?.toLocaleString(undefined, { style: 'currency', currency: 'USD' })}</div>
        <div><span class="font-medium">Monthly Pension:</span> {pensionResult.monthlyPension?.toLocaleString(undefined, { style: 'currency', currency: 'USD' })}</div>
        {#if pensionResult.notes}
          <div class="mt-2 text-sm text-gray-600 whitespace-pre-line">{pensionResult.notes}</div>
        {/if}
      </div>
    {/if}
  </div>
</main>
