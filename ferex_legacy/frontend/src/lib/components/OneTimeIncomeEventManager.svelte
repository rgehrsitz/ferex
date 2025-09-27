<script lang="ts">
  import type { OneTimeIncomeEvent } from '../../types';
  import { createEventDispatcher } from 'svelte';

  type NewEventForm = {
    name: string;
    amount: number | null;
    date: string; // YYYY-MM-DD
    type: 'Inheritance' | 'SaleOfAsset' | 'Bonus' | 'Other';
  };

  let { events = $bindable([]) } = $props<{ events?: OneTimeIncomeEvent[] }>();
  const dispatch = createEventDispatcher();

  let newEvent: NewEventForm = $state({
    name: '',
    amount: null,
    date: '',
    type: 'Other'
  });

  const eventTypes: ReadonlyArray<NewEventForm['type']> = ['Inheritance', 'SaleOfAsset', 'Bonus', 'Other'] as const;

  function generateId() {
    return Math.random().toString(36).substring(2, 15);
  }

  function addEvent() {
    if (!newEvent.name || typeof newEvent.amount !== 'number' || newEvent.amount <= 0 || !newEvent.date) {
      alert('Please provide a valid name, positive amount, and date (YYYY-MM-DD) for the event.');
      return;
    }


    const completeNewEvent: OneTimeIncomeEvent = {
        id: generateId(),
        name: newEvent.name,
        amount: newEvent.amount, // Known to be number
        date: newEvent.date,
        type: newEvent.type
    };

    events = [...events, completeNewEvent];
    newEvent = { name: '', amount: null, date: '', type: 'Other' }; // Reset form
    dispatch('change');
  }

  function removeEvent(index: number) {
    events = events.filter((_: OneTimeIncomeEvent, i: number) => i !== index);
    dispatch('change');
  }

  function handleChange() {
    dispatch('change');
  }
</script>

<div class="space-y-4">
  <h5 class="text-sm font-medium text-gray-700">Manage One-Time Income Events</h5>
  
  <div class="p-3 border border-gray-200 rounded-md bg-gray-50 space-y-3">
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3 items-end">
      <div>
        <label for="newEventName" class="block text-xs font-medium text-gray-600">Event Name</label>
        <input type="text" id="newEventName" bind:value={newEvent.name} onblur={addEvent} onkeydown={(e) => e.key === 'Enter' && addEvent()} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm" placeholder="e.g., Inheritance" />
      </div>
      <div>
        <label for="newEventDate" class="block text-xs font-medium text-gray-600">Date</label>
        <input type="date" id="newEventDate" bind:value={newEvent.date} onblur={addEvent} onkeydown={(e) => e.key === 'Enter' && addEvent()} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm" />
      </div>
      <div>
        <label for="newEventAmount" class="block text-xs font-medium text-gray-600">Amount ($)</label>
        <input type="number" id="newEventAmount" bind:value={newEvent.amount} step="any" onblur={addEvent} onkeydown={(e) => e.key === 'Enter' && addEvent()} class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm sm:text-sm" placeholder="e.g., 50000" />
      </div>
       <div>
        <label for="newEventType" class="block text-xs font-medium text-gray-600">Type</label>
        <select id="newEventType" bind:value={newEvent.type} onblur={addEvent} onkeydown={(e) => e.key === 'Enter' && addEvent()} class="mt-1 block w-full p-2 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm">
          {#each eventTypes as typeOption}
            <option value={typeOption}>{typeOption}</option>
          {/each}
        </select>
      </div>
    </div>
    <button onclick={addEvent} class="px-3 py-1.5 bg-blue-500 text-white text-sm font-medium rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
      Add One-Time Event
    </button>
  </div>

  {#if events.length > 0}
    <div class="space-y-2 pt-3">
      <h6 class="text-xs font-medium text-gray-500">Existing Events:</h6>
      <ul class="divide-y divide-gray-200">
        {#each events as event, index (event.id)}
          <li class="py-3 px-2 my-1 border border-gray-100 rounded-md hover:bg-gray-50">
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-4 gap-x-4 gap-y-3 items-center">
              <div>
                <label for={`eventName-${event.id}`} class="text-xs text-gray-500">Event Name</label>
                <input type="text" id={`eventName-${event.id}`} bind:value={event.name} onchange={handleChange} class="text-sm p-1.5 border border-gray-300 rounded-md w-full"/>
              </div>
              <div>
                <label for={`eventAmount-${event.id}`} class="text-xs text-gray-500">Amount ($)</label>
                <input type="number" id={`eventAmount-${event.id}`} bind:value={event.amount} onchange={handleChange} min="0" step="any" class="text-sm p-1.5 border border-gray-300 rounded-md w-full"/>
              </div>
              <div>
                <label for={`eventDate-${event.id}`} class="text-xs text-gray-500">Date</label>
                <input type="date" id={`eventDate-${event.id}`} bind:value={event.date} onchange={handleChange} class="text-sm p-1.5 border border-gray-300 rounded-md w-full"/>
              </div>
              <div>
                <label for={`eventType-${event.id}`} class="text-xs text-gray-500">Type</label>
                <select id={`eventType-${event.id}`} bind:value={event.type} onchange={handleChange} class="text-sm p-1.5 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md w-full">
                  {#each eventTypes as typeOption}
                    <option value={typeOption}>{typeOption}</option>
                  {/each}
                </select>
              </div>
            </div>
            <div class="mt-2 flex justify-end">
                <button onclick={() => removeEvent(index)} class="text-red-500 hover:text-red-700 text-xs font-medium py-1 px-2 rounded hover:bg-red-50">
                  Remove Event
                </button>
            </div>
          </li>
        {/each}
      </ul>
    </div>
  {:else}
    <p class="text-xs text-gray-500 pt-2">No one-time income events added yet.</p>
  {/if}
</div>
