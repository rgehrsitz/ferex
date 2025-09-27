// scenarioState.ts - Central Svelte V5 idiomatic scenario state management
import type { Scenario, ScenarioInput } from '../../types';
import { createNewScenario, addScenario as addScenarioService, deleteScenario as deleteScenarioService, duplicateScenario as duplicateScenarioService } from '../services/scenarioService';

// --- Scenario State ---
const _initialDefaultScenario = createNewScenario('Default Scenario');
export let scenarios = $state<Scenario[]>([_initialDefaultScenario]);
export let activeScenarioId = $state<string | undefined>(_initialDefaultScenario.id);

export let activeScenario = $derived(scenarios.find((s: Scenario) => s.id === activeScenarioId));
export let activeInputs = $derived(activeScenario?.inputs);

// --- Tab State ---
export const TABS = ['Annuity & Service Details', 'Social Security & TSP', 'Insurance', 'Financial & Tax Planning'] as const;
export type TabName = typeof TABS[number];
export let activeTab = $state<TabName>(TABS[0]);
export function setActiveTab (tab: TabName) { activeTab = tab; }

// --- Scenario Management Functions ---
export function addNewScenario () {
  const result = addScenarioService(scenarios);
  scenarios = result.updatedScenarios;
  activeScenarioId = result.newScenarioId;
}
export function deleteActiveScenario () {
  if (!activeScenarioId) return;
  const currentActiveIdBeforeDeletion = activeScenarioId;
  const { updatedScenarios } = deleteScenarioService(scenarios, activeScenarioId);
  scenarios = updatedScenarios;
  if (scenarios.length === 1) {
    activeScenarioId = scenarios[0]?.id;
  } else if (scenarios.every((s: Scenario) => s.id !== currentActiveIdBeforeDeletion)) {
    activeScenarioId = scenarios[0]?.id;
  } else {
    activeScenarioId = currentActiveIdBeforeDeletion;
  }
}
export function duplicateActiveScenario () {
  const result = duplicateScenarioService(scenarios, activeScenarioId);
  if (result) {
    scenarios = result.updatedScenarios;
    activeScenarioId = result.newScenarioId;
  }
}
