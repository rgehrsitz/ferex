<script lang="ts">
  import type { ScenarioInput } from "../../types";
  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  // Validation state for real-time feedback
  let validationErrors = $state<Record<string, string>>({});

  // Validation functions
  function validateHigh3Salary(
    value: number | null | undefined,
  ): string | null {
    if (!value || value <= 0) return "High-3 salary must be greater than 0";
    if (value > 500000) return "High-3 salary seems unusually high";
    return null;
  }

  function validateDateSequence(): string | null {
    if (
      !inputs.dateOfBirth ||
      !inputs.serviceComputationDate ||
      !inputs.plannedRetirementDate
    )
      return null;

    const dob = new Date(inputs.dateOfBirth);
    const scd = new Date(inputs.serviceComputationDate);
    const prd = new Date(inputs.plannedRetirementDate);

    if (scd <= dob) return "Service date must be after birth date";
    if (prd <= scd) return "Retirement date must be after service date";
    if (prd <= dob) return "Retirement date must be after birth date";

    // Check minimum retirement age
    const ageAtRetirement = prd.getFullYear() - dob.getFullYear();
    if (ageAtRetirement < 55)
      return "Federal employees typically cannot retire before age 55";

    return null;
  }

  function validateSickLeave(value: number | null | undefined): string | null {
    if (value && value < 0) return "Sick leave hours cannot be negative";
    if (value && value > 3000) return "Sick leave hours seem unusually high";
    return null;
  }

  function calculateRetirementAge(
    birthDateStr: string | null | undefined,
    referenceDateStr: string | null | undefined,
  ): number | null {
    if (!birthDateStr || !referenceDateStr) {
      return null;
    }
    try {
      const birthDateParts = birthDateStr.split("-").map(Number);
      const referenceDateParts = referenceDateStr.split("-").map(Number);

      if (
        birthDateParts.length !== 3 ||
        referenceDateParts.length !== 3 ||
        birthDateParts.some(isNaN) ||
        referenceDateParts.some(isNaN)
      )
        return null;

      const birthDate = new Date(
        Date.UTC(birthDateParts[0], birthDateParts[1] - 1, birthDateParts[2]),
      );
      const referenceDate = new Date(
        Date.UTC(
          referenceDateParts[0],
          referenceDateParts[1] - 1,
          referenceDateParts[2],
        ),
      );

      if (isNaN(birthDate.getTime()) || isNaN(referenceDate.getTime())) {
        return null;
      }

      let age = referenceDate.getUTCFullYear() - birthDate.getUTCFullYear();
      const monthDiff = referenceDate.getUTCMonth() - birthDate.getUTCMonth();
      if (
        monthDiff < 0 ||
        (monthDiff === 0 && referenceDate.getUTCDate() < birthDate.getUTCDate())
      ) {
        age--;
      }
      return age;
    } catch (e) {
      console.error("Error calculating retirement age:", e);
      return null;
    }
  }

  let retirementAgeAtPlannedDate = $derived(
    calculateRetirementAge(inputs.dateOfBirth, inputs.plannedRetirementDate),
  );

  // Real-time validation using $effect
  $effect(() => {
    const errors: Record<string, string> = {};

    const salaryError = validateHigh3Salary(inputs.high3Salary);
    if (salaryError) errors.high3Salary = salaryError;

    const dateError = validateDateSequence();
    if (dateError) errors.dateSequence = dateError;

    const sickLeaveError = validateSickLeave(
      inputs.unusedSickLeaveHoursAtRetirement,
    );
    if (sickLeaveError) errors.sickLeave = sickLeaveError;

    validationErrors = errors;
  });
</script>

<div class="p-4 border border-gray-200 rounded-md bg-gray-50 space-y-6">
  <h3
    class="text-xl font-semibold text-gray-800 mb-4 py-2 border-b border-gray-300"
  >
    Core Annuity & Service
  </h3>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
    <div>
      <label
        for="calculationSystem"
        class="block text-sm font-medium text-gray-700"
        >Retirement System:</label
      >
      <select
        id="calculationSystem"
        bind:value={inputs.calculationSystem}
        class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
      >
        <option value="">-- Select System --</option>
        <option value="FERS">FERS</option>
        <option value="CSRS">CSRS</option>
      </select>
    </div>

    <div>
      <label for="high3Salary" class="block text-sm font-medium text-gray-700"
        >High-3 Average Salary ($):</label
      >
      <input
        type="number"
        id="high3Salary"
        bind:value={inputs.high3Salary}
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm {validationErrors.high3Salary
          ? 'border-red-500 focus:ring-red-500 focus:border-red-500'
          : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'}"
        placeholder="e.g., 125000"
      />
      {#if validationErrors.high3Salary}
        <p class="mt-1 text-xs text-red-600">{validationErrors.high3Salary}</p>
      {/if}
    </div>

    <div>
      <label
        for="serviceComputationDate"
        class="block text-sm font-medium text-gray-700"
        >Overall Service Comp. Date (SCD):</label
      >
      <input
        type="date"
        id="serviceComputationDate"
        bind:value={inputs.serviceComputationDate}
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm {validationErrors.dateSequence
          ? 'border-red-500 focus:ring-red-500 focus:border-red-500'
          : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'}"
        title="This is your overall SCD. Individual periods below contribute to this."
      />
      <p class="mt-1 text-xs text-gray-500">
        This is your overall SCD. Individual periods below help detail its
        calculation.
      </p>
    </div>

    <div>
      <label for="dateOfBirth" class="block text-sm font-medium text-gray-700"
        >Date of Birth:</label
      >
      <input
        type="date"
        id="dateOfBirth"
        bind:value={inputs.dateOfBirth}
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm {validationErrors.dateSequence
          ? 'border-red-500 focus:ring-red-500 focus:border-red-500'
          : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'}"
      />
    </div>

    <div>
      <label
        for="plannedRetirementDate"
        class="block text-sm font-medium text-gray-700"
        >Planned Retirement Date:</label
      >
      <input
        type="date"
        id="plannedRetirementDate"
        bind:value={inputs.plannedRetirementDate}
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm {validationErrors.dateSequence
          ? 'border-red-500 focus:ring-red-500 focus:border-red-500'
          : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'}"
      />
      {#if validationErrors.dateSequence}
        <p class="mt-1 text-xs text-red-600">{validationErrors.dateSequence}</p>
      {/if}
    </div>

    <div>
      <label
        for="unusedSickLeaveHours"
        class="block text-sm font-medium text-gray-700"
        >Unused Sick Leave Hours (at retirement):</label
      >
      <input
        type="number"
        id="unusedSickLeaveHours"
        bind:value={inputs.unusedSickLeaveHoursAtRetirement}
        class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm {validationErrors.sickLeave
          ? 'border-red-500 focus:ring-red-500 focus:border-red-500'
          : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'}"
        placeholder="e.g., 520"
      />
      {#if validationErrors.sickLeave}
        <p class="mt-1 text-xs text-red-600">{validationErrors.sickLeave}</p>
      {/if}
    </div>

    <div>
      <label
        for="employeeContributions"
        class="block text-sm font-medium text-gray-700"
        >Total FERS/CSRS Contributions ($) (Optional):</label
      >
      <input
        type="number"
        id="employeeContributions"
        bind:value={inputs.employeeContributions}
        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        placeholder="e.g., 60000"
        title="Total actual dollars paid into the FERS or CSRS retirement fund. Used for tax calculation of annuity."
      />
      <div class="mt-1 text-xs text-gray-500 space-y-1">
        <p>
          Helps determine tax-free portion of annuity. If unknown, leave blank.
        </p>
        <p>
          <strong>To find this:</strong> Check your annual SF-50s, or contact HR/RPPS
          for your total retirement contributions.
        </p>
        <p>
          <strong>Rough estimate:</strong> For FERS: ~0.8% of salary per year; For
          CSRS: ~7% of salary per year.
        </p>
      </div>
    </div>

    <div>
      <label for="retirementAge" class="block text-sm font-medium text-gray-700"
        >Retirement Age (calculated):</label
      >
      <div
        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-600 sm:text-sm"
      >
        {retirementAgeAtPlannedDate !== null
          ? retirementAgeAtPlannedDate
          : "Enter birth date and retirement date"}
      </div>
      <p class="mt-1 text-xs text-gray-500">Age at planned retirement date</p>
    </div>
  </div>
</div>
