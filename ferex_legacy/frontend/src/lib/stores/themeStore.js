// Theme state management using Svelte stores
import { writable } from 'svelte/store';

// Initialize the theme store with the current theme
function createThemeStore () {
  // Initially check for saved theme or system preference
  const getInitialTheme = () => {
    // Check for server-side rendering
    if (typeof window === 'undefined') return true;

    const savedTheme = localStorage.getItem('theme');

    if (savedTheme === 'dark') return true;
    if (savedTheme === 'light') return false;

    // If no stored preference, use system preference or default to dark
    return window.matchMedia('(prefers-color-scheme: dark)').matches || true;
  };

  const isDark = writable(getInitialTheme());

  // Subscribe to the store and update the DOM whenever it changes
  const { subscribe, set } = isDark;

  return {
    subscribe,
    toggle: () => {
      isDark.update(currentValue => {
        const newValue = !currentValue;

        // Update the DOM
        if (newValue) {
          document.documentElement.classList.add('dark');
          localStorage.setItem('theme', 'dark');
        } else {
          document.documentElement.classList.remove('dark');
          localStorage.setItem('theme', 'light');
        }

        // Debug
        // console.log('Theme toggled to:', newValue ? 'dark' : 'light');
        // console.log('HTML classes:', document.documentElement.className);

        return newValue;
      });
    },
    // Explicit setters
    setDark: () => {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
      set(true);
    },
    setLight: () => {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
      set(false);
    },
    // Initialize the DOM based on the current state
    initialize: () => {
      const isDarkMode = getInitialTheme();
      if (isDarkMode) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
      set(isDarkMode);
    }
  };
}

export const theme = createThemeStore();
