<script lang="ts">
  import { onMount } from "svelte";
  import { theme } from "../../stores/themeStore";
  import type { Scenario } from "../../../types";
  import { createEventDispatcher } from "svelte";

  // Scenario props
  let { scenarios = [], activeScenarioId } = $props<{
    scenarios?: Scenario[];
    activeScenarioId?: string;
  }>();

  // Svelte 5 state for renaming logic
  let isRenaming = $state(false);
  let renameValue = $state("");

  // Find active scenario name
  let scenarioName = $derived(() => {
    const found = scenarios.find((s: Scenario) => s.id === activeScenarioId);
    return found ? found.name : "No Scenario";
  });

  function startRename() {
    if (!activeScenarioId) return;
    const found = scenarios.find((s: Scenario) => s.id === activeScenarioId);
    if (found) {
      renameValue = found.name;
      isRenaming = true;
    }
  }

  function cancelRename() {
    isRenaming = false;
    renameValue = "";
  }

  function commitRename() {
    if (!activeScenarioId) return;
    const trimmed = renameValue.trim();
    if (
      trimmed &&
      scenarios.find((s: Scenario) => s.id === activeScenarioId)?.name !==
        trimmed
    ) {
      dispatch("renameScenario", { id: activeScenarioId, name: trimmed });
    }
    isRenaming = false;
    renameValue = "";
  }

  function handleRenameKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      commitRename();
    } else if (e.key === "Escape") {
      cancelRename();
    }
  }

  let showNotifications = $state(false);
  let notificationCount = $state(2); // placeholder
  let showUserMenu = $state(false);
  let isDark = $state(true);

  const dispatch = createEventDispatcher<{
    openFile: void;
    saveFile: void;
    selectScenario: string;
    addScenario: void;
    duplicateScenario: void;
    deleteScenario: void;
    renameScenario: { id: string; name: string };
  }>();

  onMount(() => {
    const unsubscribe = theme.subscribe((value) => {
      isDark = value;
    });
    return unsubscribe;
  });

  function toggleTheme() {
    theme.toggle();
  }
  function toggleNotifications() {
    showNotifications = !showNotifications;
  }
  function toggleUserMenu() {
    showUserMenu = !showUserMenu;
  }
  function handleOpenFile() {
    dispatch("openFile");
  }
  function handleSaveFile() {
    dispatch("saveFile");
  }
  function handleScenarioSwitch(e: Event) {
    const id = (e.target as HTMLSelectElement).value;
    dispatch("selectScenario", id);
  }
  function handleAdd() {
    dispatch("addScenario");
  }
  function handleDuplicate() {
    dispatch("duplicateScenario");
  }
  function handleDelete() {
    dispatch("deleteScenario");
  }
</script>

<header
  class="h-14 flex items-center justify-between px-4 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700"
>
  <div class="flex items-center gap-4">
    <span class="font-bold text-blue-900 dark:text-blue-200 text-lg">FeReX</span
    >
    <div class="relative flex items-center gap-2">
      <label for="scenario-switcher" class="sr-only">Active Scenario</label>
      <select
        id="scenario-switcher"
        class="rounded border-gray-300 px-2 py-1 bg-white text-gray-900 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:ring-blue-500 focus:border-blue-500"
        aria-label="Active Scenario"
        onchange={handleScenarioSwitch}
        bind:value={activeScenarioId}
        disabled={isRenaming}
      >
        {#each scenarios as s: { id: string; name: string } (s.id)}
          <option value={s.id}>{s.name}</option>
        {/each}
      </select>
      {#if isRenaming}
        <input
          class="rounded border-gray-300 px-2 py-1 ml-1 w-36 focus:ring-blue-500 focus:border-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
          bind:value={renameValue}
          onblur={commitRename}
          onkeydown={handleRenameKeydown}
          aria-label="Rename scenario"
          maxlength={40}
        />
        <button
          class="ml-1 px-2 py-1 rounded bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200 text-xs font-medium"
          title="Cancel rename"
          onclick={cancelRename}
        >
          ✕
        </button>
      {:else}
        <button
          class="ml-1 px-2 py-1 rounded bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-600 text-xs font-medium"
          title="Rename Scenario"
          onclick={startRename}
          aria-label="Rename scenario"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="inline h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15.232 5.232l3.536 3.536M9 13h3l8-8a2.828 2.828 0 00-4-4l-8 8v3z"
            /></svg
          >
        </button>
      {/if}
      <button
        class="ml-1 px-2 py-1 rounded bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 hover:bg-green-200 dark:hover:bg-green-800 text-xs font-medium"
        title="Add Scenario"
        onclick={handleAdd}
      >
        +
      </button>
      <button
        class="px-2 py-1 rounded bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 hover:bg-blue-200 dark:hover:bg-blue-800 text-xs font-medium disabled:opacity-50"
        title="Duplicate Scenario"
        onclick={handleDuplicate}
        disabled={!activeScenarioId}
      >
        ⧉
      </button>
      <button
        class="px-2 py-1 rounded bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 hover:bg-red-200 dark:hover:bg-red-800 text-xs font-medium disabled:opacity-50"
        title="Delete Scenario"
        onclick={handleDelete}
        disabled={!activeScenarioId || scenarios.length === 1}
      >
        🗑
      </button>
    </div>
  </div>
  <div class="flex items-center gap-4">
    <!-- Open Button -->
    <button
      type="button"
      class="px-3 py-1 rounded bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 hover:bg-blue-200 dark:hover:bg-blue-800 font-medium transition-colors"
      onclick={handleOpenFile}
      aria-label="Open scenario file"
      title="Open... (Ctrl+O)"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="inline h-5 w-5 mr-1"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 4v16m8-8H4"
        /></svg
      >
      Open
    </button>
    <!-- Save Button -->
    <button
      type="button"
      class="px-3 py-1 rounded bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 hover:bg-green-200 dark:hover:bg-green-800 font-medium transition-colors"
      onclick={handleSaveFile}
      aria-label="Save scenario file"
      title="Save... (Ctrl+S)"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="inline h-5 w-5 mr-1"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M5 13l4 4L19 7"
        /></svg
      >
      Save
    </button>
    <!-- Theme Toggle Button -->
    <button
      type="button"
      class="rounded-full p-2 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
      title={isDark ? "Switch to light mode" : "Switch to dark mode"}
      onclick={toggleTheme}
      aria-label="Toggle theme"
    >
      {#if isDark}
        <!-- Sun icon -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 text-yellow-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          ><circle cx="12" cy="12" r="5" fill="currentColor" /><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 3v1m0 16v1m8.66-13.66l-.71.71M5.05 18.36l-.71.71m13.66 8.66l-.71-.71M5.05 5.64l-.71-.71M21 12h1M3 12H2m16.24 7.76l-.71-.71M7.76 7.76l-.71-.71"
          /></svg
        >
      {:else}
        <!-- Moon icon -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 text-gray-700 dark:text-gray-200"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 12.79A9 9 0 1111.21 3a7 7 0 109.79 9.79z"
          /></svg
        >
      {/if}
    </button>
    <button
      class="relative"
      onclick={toggleNotifications}
      aria-label="Notifications"
    >
      <svg width="24" height="24" fill="none" stroke="currentColor"
        ><circle cx="12" cy="12" r="10" /><path d="M12 8v4l3 3" /></svg
      >
      {#if notificationCount > 0}
        <span
          class="absolute -top-1 -right-1 bg-red-500 text-white rounded-full text-xs w-5 h-5 flex items-center justify-center"
          >{notificationCount}</span
        >
      {/if}
    </button>
    {#if showNotifications}
      <div
        class="absolute right-16 top-14 w-72 bg-white dark:bg-gray-800 shadow-lg rounded p-4 z-50"
      >
        <div class="font-semibold mb-2">Notifications</div>
        <ul class="space-y-2">
          <li class="text-sm">Monte Carlo run complete</li>
          <li class="text-sm">Scenario B saved</li>
        </ul>
        <button class="mt-2 text-xs text-blue-600 hover:underline"
          >Mark all as read</button
        >
      </div>
    {/if}
    <button
      class="ml-2 flex items-center"
      onclick={toggleUserMenu}
      aria-label="User menu"
    >
      <svg width="24" height="24" fill="none" stroke="currentColor"
        ><circle cx="12" cy="8" r="4" /><path
          d="M4 20c0-4 8-4 8-4s8 0 8 4"
        /></svg
      >
    </button>
    {#if showUserMenu}
      <div
        class="absolute right-4 top-14 w-48 bg-white dark:bg-gray-800 shadow-lg rounded p-2 z-50"
      >
        <a
          class="block px-3 py-2 hover:bg-blue-100 dark:hover:bg-blue-800 rounded"
          href="/">Profile</a
        >
        <a
          class="block px-3 py-2 hover:bg-blue-100 dark:hover:bg-blue-800 rounded"
          href="/">Docs</a
        >
        <a
          class="block px-3 py-2 hover:bg-blue-100 dark:hover:bg-blue-800 rounded"
          href="/">Keyboard Shortcuts</a
        >
        <a
          class="block px-3 py-2 hover:bg-blue-100 dark:hover:bg-blue-800 rounded"
          href="/">About</a
        >
        <a
          class="block px-3 py-2 hover:bg-blue-100 dark:hover:bg-blue-800 rounded"
          href="/">Logout</a
        >
      </div>
    {/if}
  </div>
</header>
