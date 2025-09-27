import './style.css';
import './app.css'; // Custom app styles including dark mode
import App from './App.svelte';
import { mount } from 'svelte';
import { theme } from './lib/stores/themeStore';

// Initialize theme (adds/removes dark class based on saved preference)
theme.initialize();

mount(App, { target: document.getElementById('app')! });
