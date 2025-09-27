<script lang="ts">
  import type {
    ScenarioInput,
    TSPFundAllocation,
    TSPReturnAssumptions,
  } from "../../types";
  import SocialSecuritySection from "./SocialSecuritySection.svelte";
  import TSPContributions from "./TSPContributions.svelte";
  import TSPAllocations from "./TSPAllocations.svelte";
  import TSPReturnAssumptionsComponent from "./TSPReturnAssumptions.svelte";
  import TSPWithdrawals from "./TSPWithdrawals.svelte";

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
</script>

<div class="space-y-8">
  <SocialSecuritySection bind:inputs />

  <!-- TSP Section -->
  <section class="space-y-4 p-4 border border-gray-200 rounded-md shadow-sm">
    <h3 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-3">
      Thrift Savings Plan (TSP)
    </h3>

    <TSPContributions bind:inputs />
    <TSPAllocations bind:inputs />
    <TSPReturnAssumptionsComponent bind:inputs />
    <TSPWithdrawals bind:inputs />
  </section>
</div>

<style>
  /* Removed @apply rules as Tailwind classes are now applied directly in the HTML */
</style>
