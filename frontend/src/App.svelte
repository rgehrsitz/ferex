<script lang="ts">
  import Header from './components/Header.svelte';
  import ScenarioSidebar from './components/ScenarioSidebar.svelte';
  import ScenarioTabs from './components/ScenarioTabs.svelte';
  import CompareView from './components/CompareView.svelte';
  import type { Scenario, ScenarioData } from './types/scenario.js';

  function selectScenario(id: number): void { selectedScenarioId = id; compareScenarioId = null; }
  function selectCompare(id: number): void { compareScenarioId = id; }

  // Default data for new scenario
  const getDefaultScenarioData = (): ScenarioData => ({
    pension: {
      retirementSystem: 'FERS',
      highThreeSalary: 100000,
      yearsOfService: 30,
      retirementAge: 62,
      unusedSickLeave: 1000,
      militaryService: 0,
      isPartTime: false,
      partTimeProrationFactor: 1.0,
      csrsOffset: false,
      survivorBenefit: 'full'
    },
    socialSecurity: {
      startAge: 62,
      estimatedMonthlyBenefit: 2000,
      isEligible: true,
      birthYear: 1970,
      birthMonth: 1
    },
    tsp: {
      currentBalance: 500000,
      traditionalBalance: 400000,
      rothBalance: 100000,
      annualContribution: 20500,
      expectedReturnRate: 0.06,
      withdrawalStrategy: 'fixed',
      fixedWithdrawalAmount: 2000,
      withdrawalPercentage: 0.04,
      withdrawalStartAge: 62
    },
    tax: {
      filingStatus: 'married_joint',
      stateOfResidence: 'VA',
      additionalIncome: 0,
      additionalDeductions: 0,
      stateIncomeTaxRate: 0.05
    },
    cola: {
      assumedInflationRate: 0.025,
      applyCOLAToPension: true,
      applyColaToSocialSecurity: true
    },
    otherIncome: {
      sources: []
    }
  });

  let scenarios: Scenario[] = [
    { id: 1, name: "Scenario A", data: getDefaultScenarioData() },
    { id: 2, name: "Scenario B", data: getDefaultScenarioData() }
  ];
  let selectedScenarioId: number | null = 1;
  let compareScenarioId: number | null = null;

  function addScenario(): void {
    const id = Math.max(...scenarios.map(s => s.id)) + 1;
    scenarios = [...scenarios, { 
      id, 
      name: `Scenario ${String.fromCharCode(64 + id)}`, 
      data: getDefaultScenarioData() 
    }];
    selectedScenarioId = id;
    compareScenarioId = null;
  }
  function duplicateScenario(id: number): void {
    const orig = scenarios.find(s => s.id === id);
    if (!orig) return;
    const newId = Math.max(...scenarios.map(s => s.id)) + 1;
    scenarios = [...scenarios, { id: newId, name: orig.name + " Copy", data: { ...orig.data } }];
  }
  function deleteScenario(id: number): void {
    scenarios = scenarios.filter(s => s.id !== id);
    if (selectedScenarioId === id) selectedScenarioId = scenarios[0]?.id || null;
    if (compareScenarioId === id) compareScenarioId = null;
  }

  let selectedScenario: Scenario | undefined;
  let compareScenario: Scenario | undefined;
  $: selectedScenario = scenarios.find(s => s.id === selectedScenarioId);
  $: compareScenario = compareScenarioId ? scenarios.find(s => s.id === compareScenarioId) : undefined;
</script>

<div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100" id="main-container">
  <Header />
  <div class="flex overflow-hidden">
    <ScenarioSidebar
      {scenarios}
      {selectedScenarioId}
      {compareScenarioId}
      on:add={addScenario}
      on:duplicate={e => duplicateScenario(e.detail)}
      on:delete={e => deleteScenario(e.detail)}
      on:select={e => selectScenario(e.detail)}
      on:compare={e => selectCompare(e.detail)}
    />
    <main class="flex-1 p-6 overflow-auto">
      {#if compareScenario}
        <CompareView {selectedScenario} {compareScenario} />
      {:else if selectedScenario}
        <ScenarioTabs bind:scenario={selectedScenario} />
      {:else}
        <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow text-center">
          <p class="text-gray-700 dark:text-gray-300">No scenarios defined.</p>
        </div>
      {/if}
    </main>
  </div>
</div>
