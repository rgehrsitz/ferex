<script lang="ts">
  import PensionSection from './PensionSection.svelte';
  import SocialSecuritySection from './SocialSecuritySection.svelte';
  import TSPSection from './TSPSection.svelte';
  import TaxSection from './TaxSection.svelte';
  import COLASection from './COLASection.svelte';
  import OtherIncomeSection from './OtherIncomeSection.svelte';
  export let scenario;
  let tab = 'Pension';
  const tabs = [
    { label: 'Pension', comp: PensionSection },
    { label: 'Social Security', comp: SocialSecuritySection },
    { label: 'TSP', comp: TSPSection },
    { label: 'Tax', comp: TaxSection },
    { label: 'COLA', comp: COLASection },
    { label: 'Other Income', comp: OtherIncomeSection }
  ];
</script>

<div>
  <div class="flex border-b mb-4">
    {#each tabs as t}
      <button class="px-4 py-2" class:selected={tab===t.label} on:click={() => tab=t.label}>{t.label}</button>
    {/each}
  </div>
  {#each tabs as t}
    {#if tab === t.label}
      <svelte:component this={t.comp} bind:data={scenario.data[t.label.replace(/ /g,'').toLowerCase()]} />
    {/if}
  {/each}
</div>
