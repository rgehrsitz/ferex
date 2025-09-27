import { writable } from 'svelte/store';

// Store to track user's preferred chart display mode
export const displayMode = writable<'monthly' | 'yearly'>('yearly');
