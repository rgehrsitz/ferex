import './style.css';
import AppWrapper from './AppWrapper.svelte';
import { createRoot } from 'svelte';

try {
  const target = document.getElementById('app');
  if (!target) {
    console.error('Cannot find #app element to mount application');
    document.body.innerHTML = '<div style="color: red; padding: 20px;">Error: Cannot find mount point. Please check the console.</div>';
  } else {
    console.log('Mounting application to #app element');
    createRoot(AppWrapper, { target });
  }
} catch (err) {
  console.error('Failed to mount application:', err);
  document.body.innerHTML = '<div style="color: red; padding: 20px;">Error: Failed to initialize application. Please check the console.</div>';
}
