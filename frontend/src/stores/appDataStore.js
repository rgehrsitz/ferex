import { writable } from 'svelte/store';

// Central store for all user-editable/persisted data
export const appData = writable({
  scenarios: [],          // array of scenario objects
  userProfile: {          // demographic & global inputs
    birthYear: null,
    birthMonth: null,
    retirementAge: null
  },
  preferences: {}         // theme, units, etc.
});
