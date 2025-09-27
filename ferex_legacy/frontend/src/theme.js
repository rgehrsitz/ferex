// Theme management for the entire application

// Function to initialize theme based on localStorage or system preference
export function initializeTheme () {
  const savedTheme = localStorage.getItem('theme');

  // Check if we have a saved preference or should use system preference
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark');
    return true; // Return the current state (dark=true)
  } else if (savedTheme === 'light') {
    document.documentElement.classList.remove('dark');
    return false; // Return the current state (dark=false)
  } else {
    // Default to dark mode for this app
    document.documentElement.classList.add('dark');
    return true;
  }
}

// Function to toggle theme
/**
 * @param {boolean} isDark
 */
export function toggleTheme (isDark) {
  if (isDark) {
    document.documentElement.classList.add('dark');
    localStorage.setItem('theme', 'dark');
  } else {
    document.documentElement.classList.remove('dark');
    localStorage.setItem('theme', 'light');
  }

  // Ensure Tailwind is aware of the change
  document.documentElement.style.colorScheme = isDark ? 'dark' : 'light';

  // Help diagnostics
  // console.log('Theme toggled to:', isDark ? 'dark' : 'light');
  // console.log('HTML classes:', document.documentElement.className);
}
