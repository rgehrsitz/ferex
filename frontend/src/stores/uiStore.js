import { writable } from 'svelte/store';

// UI state store for ephemeral UI flags and selections
export const selectedScenarioId = writable(null);
export const compareScenarioId  = writable(null);

// Additional UI state (e.g., selected tab, modals)
export const uiState = writable({
  selectedTab: 'inputs',
  modalOpen: false
});
