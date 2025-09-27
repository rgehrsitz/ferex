<script lang="ts">
  import type { ScenarioInput, ServicePeriod, LWOPPeriod, InsurableInterestDetails } from '../../types';
  import CoreAnnuityDetails from './CoreAnnuityDetails.svelte';
  import ServiceHistory from './ServiceHistory.svelte';
  import LwopHistory from './LwopHistory.svelte';
  import FersSpecifics from './FersSpecifics.svelte';
  import CsrsSpecifics from './CsrsSpecifics.svelte';
  import DisabilityRetirementDetails from './DisabilityRetirementDetails.svelte';
  import GeneralSurvivorBenefits from './GeneralSurvivorBenefits.svelte';

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  // Initialize arrays if they are undefined
  if (!inputs.servicePeriods) {
    inputs.servicePeriods = [];
  }
  if (!inputs.lwopPeriods) {
    inputs.lwopPeriods = [];
  }

  // Auto-creation logic moved to calculationService.ts to ensure it runs every calculation

  let showFERSSpecific = $derived(inputs.calculationSystem === 'FERS');
  let showCSRSSpecific = $derived(inputs.calculationSystem === 'CSRS');
  let showCSRSOffsetOptions = $derived(inputs.isCSRSOffset && inputs.calculationSystem === 'CSRS');

  // Calculate retirement age for SRS logic
  function calculateAge(birthDateStr: string | null | undefined, referenceDateStr: string | null | undefined): number | null {
    if (!birthDateStr || !referenceDateStr) {
      return null;
    }
    try {
      const birthDateParts = birthDateStr.split('-').map(Number);
      const referenceDateParts = referenceDateStr.split('-').map(Number);

      if (birthDateParts.length !== 3 || referenceDateParts.length !== 3 || birthDateParts.some(isNaN) || referenceDateParts.some(isNaN)) return null;

      const birthDate = new Date(Date.UTC(birthDateParts[0], birthDateParts[1] - 1, birthDateParts[2]));
      const referenceDate = new Date(Date.UTC(referenceDateParts[0], referenceDateParts[1] - 1, referenceDateParts[2]));

      if (isNaN(birthDate.getTime()) || isNaN(referenceDate.getTime())) {
        return null;
      }

      let age = referenceDate.getUTCFullYear() - birthDate.getUTCFullYear();
      const monthDiff = referenceDate.getUTCMonth() - birthDate.getUTCMonth();
      if (monthDiff < 0 || (monthDiff === 0 && referenceDate.getUTCDate() < birthDate.getUTCDate())) {
        age--;
      }
      return age;
    } catch (e) {
      console.error("Error calculating age:", e);
      return null;
    }
  }

  let retirementAgeAtPlannedDate = $derived(
    calculateAge(inputs.dateOfBirth, inputs.plannedRetirementDate)
  );

  let showSRSInput = $derived(
    inputs.calculationSystem === 'FERS' &&
    inputs.dateOfBirth && 
    inputs.plannedRetirementDate && 
    retirementAgeAtPlannedDate !== null &&
    retirementAgeAtPlannedDate < 62
  );

  // Reactive effect to set default for estimatedSSBenefitAt62ForSRS
  // function handleCsrsSurvivorElectionChange() { // Moved to CsrsSpecifics.svelte
    if (inputs.survivorBenefitCSRS) {
      if (inputs.survivorBenefitCSRS.election === 'PartialCustomBase') {
        // Ensure customBaseAmountForPartial is initialized if not already
        if (inputs.survivorBenefitCSRS.customBaseAmountForPartial === null || inputs.survivorBenefitCSRS.customBaseAmountForPartial === undefined) {
          inputs.survivorBenefitCSRS.customBaseAmountForPartial = 0;
        }
      } else {
        inputs.survivorBenefitCSRS.customBaseAmountForPartial = null;
      }
      // Placeholder for insurable interest, similar to FERS, if we add that option
      // if (inputs.survivorBenefitCSRS.election === 'InsurableInterest') { // Moved to CsrsSpecifics.svelte
      //   if (!inputs.survivorBenefitCSRS.insurableInterestDetails) {
      //     inputs.survivorBenefitCSRS.insurableInterestDetails = { beneficiaryDOB: '', relationship: '' };
      //   }
      // } else {
      //   // inputs.survivorBenefitCSRS.insurableInterestDetails = null; // Only if not 'PartialCustomBase'
      // }
    }
  // }
</script>
<div class="flex flex-col space-y-4">
  <CoreAnnuityDetails bind:inputs />
  <DisabilityRetirementDetails bind:inputs />
  <ServiceHistory bind:inputs />
  <LwopHistory bind:inputs />
  <GeneralSurvivorBenefits bind:inputs />
  {#if showFERSSpecific}
    <FersSpecifics bind:inputs {showSRSInput} />
  {/if}
  {#if showCSRSSpecific}
    <CsrsSpecifics bind:inputs {showCSRSOffsetOptions} />
  {/if}
</div>
