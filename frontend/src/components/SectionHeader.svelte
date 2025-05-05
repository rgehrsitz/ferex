<script lang="ts">
  export let sectionName: string;
  export let scenarioName: string;
</script>

<div class="mb-6 flex justify-between items-center">
  <h2 class="text-xl font-medium text-gray-800 dark:text-gray-200">{sectionName}</h2>
  <div class="bg-primary-50 dark:bg-primary-900 px-3 py-1 rounded-full text-sm text-primary-700 dark:text-primary-300 flex items-center">
    <span class="mr-1">📋</span>
    <span>Editing: <span class="font-semibold">{scenarioName}</span></span>
  </div>
</div>