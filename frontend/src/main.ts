import './style.css';
import App from './App.svelte';
import { mount } from 'svelte';

try {
  const target = document.querySelector('#app');
  if (!target) throw new Error('No #app element found');
  mount(App, { target });
} catch (err) {
  console.error('Failed to mount application:', err);
  document.body.innerHTML = '<div style="color: red; padding: 20px;">Error: Failed to initialize application. Please check the console.</div>';
}
