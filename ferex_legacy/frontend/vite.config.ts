import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      // Ensure wailsjs imports can be resolved from the project root
      'wailsjs': resolve(__dirname, 'wailsjs')
    }
  }
})
