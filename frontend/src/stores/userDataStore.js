import { get } from 'svelte/store';
import { appData } from './appDataStore';

// Central user profile helpers using appDataStore
export function updateUserProfile(patch) {
    appData.update(d => ({
        ...d,
        userProfile: {
            ...d.userProfile,
            ...patch
        }
    }));
}

export function getUserProfile() {
    return get(appData).userProfile;
}