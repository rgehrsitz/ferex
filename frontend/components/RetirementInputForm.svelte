<script lang="ts">
import { createEventDispatcher } from 'svelte';
import { writable } from 'svelte/store';

// Simple state for MVP, can expand as needed
const input = writable({
    system: 'FERS', // Add system selector
    high3Salary: 90000,
    yearsOfService: 30,
    ageAtRetirement: 62,
    initialBalance: 500000,
    annualWithdrawal: 40000,
    expectedReturn: 0.06,
    returnStdDev: 0.12,
    inflationMean: 0.025,
    inflationStdDev: 0.01,
    numSimulations: 1000,
    years: 30
});

const dispatch = createEventDispatcher();

function submitForm() {
    input.subscribe((val) => {
        dispatch('submit', val);
    })();
}
</script>

<form class="bg-green-50 rounded-lg shadow-2xl w-[700px] mx-auto p-12 space-y-10" on:submit|preventDefault={submitForm}>
    <h2 class="text-3xl font-extrabold text-center text-green-800 mb-12 tracking-tight">Retirement Input</h2>
    <div class="flex flex-col space-y-2 mb-6">
      <label class="text-lg font-bold text-slate-800" for="retirement-system">Retirement System</label>
      <select class="block w-full rounded border border-slate-400 px-4 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="retirement-system" bind:value={$input.system}>
        <option value="FERS">FERS</option>
        <option value="CSRS">CSRS</option>
        <option value="CSRS Offset">CSRS Offset</option>
      </select>
    </div>
    <div class="grid grid-cols-2 gap-8">
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-200" for="high3-salary">High-3 Salary</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="high3-salary" type="number" bind:value={$input.high3Salary} min="0" />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="years-of-service">Years of Service</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="years-of-service" type="number" bind:value={$input.yearsOfService} min="0" />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="age-at-retirement">Age at Retirement</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="age-at-retirement" type="number" bind:value={$input.ageAtRetirement} min="0" />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="initial-tsp-balance">Initial TSP Balance</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="initial-tsp-balance" type="number" bind:value={$input.initialBalance} min="0" />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="annual-withdrawal">Annual Withdrawal</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="annual-withdrawal" type="number" bind:value={$input.annualWithdrawal} min="0" />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="expected-return">Expected Return (%)</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="expected-return" type="number" step="0.01" bind:value={$input.expectedReturn} />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="return-std-dev">Return Std Dev (%)</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="return-std-dev" type="number" step="0.01" bind:value={$input.returnStdDev} />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="inflation-mean">Inflation Mean (%)</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="inflation-mean" type="number" step="0.01" bind:value={$input.inflationMean} />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="inflation-std-dev">Inflation Std Dev (%)</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="inflation-std-dev" type="number" step="0.01" bind:value={$input.inflationStdDev} />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="num-simulations">Monte Carlo Simulations</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="num-simulations" type="number" bind:value={$input.numSimulations} min="1" />
        </div>
        <div class="flex flex-col space-y-1">
            <label class="text-base font-bold text-slate-800" for="years-to-simulate">Years to Simulate</label>
            <input class="block w-full rounded border border-slate-300 px-3 py-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500" id="years-to-simulate" type="number" bind:value={$input.years} min="1" />
        </div>
    </div>
    <button class="w-full mt-12 py-3 px-4 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-bold shadow-xl transition-colors duration-150 text-xl" type="submit">Calculate</button>
</form>

