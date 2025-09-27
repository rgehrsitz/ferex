<script lang="ts">
  import type { CalculationAuditTrail, AuditStep } from "../types/auditTrail";

  interface Props {
    auditTrail: CalculationAuditTrail | null;
    showDetailed?: boolean;
  }

  let { auditTrail, showDetailed = false }: Props = $props();
  let expandedSteps = $state<Set<number>>(new Set());
  let showFormulas = $state(false);
  let showInputs = $state(false);

  function toggleStep(stepNumber: number) {
    if (expandedSteps.has(stepNumber)) {
      expandedSteps.delete(stepNumber);
    } else {
      expandedSteps.add(stepNumber);
    }
    expandedSteps = new Set(expandedSteps); // Trigger reactivity
  }

  function toggleAllSteps() {
    if (expandedSteps.size === auditTrail?.steps.length) {
      expandedSteps.clear();
    } else {
      expandedSteps = new Set(
        auditTrail?.steps.map((s: AuditStep) => s.stepNumber) || []
      );
    }
  }

  function formatCurrency(value: number): string {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(value);
  }

  function formatNumber(value: number): string {
    if (Number.isInteger(value)) {
      return value.toString();
    }
    return value.toFixed(4).replace(/\.?0+$/, "");
  }

  function getSeverityColor(stepName: string): string {
    if (
      stepName.toLowerCase().includes("reduction") ||
      stepName.toLowerCase().includes("cost")
    ) {
      return "text-red-600 dark:text-red-400";
    }
    if (
      stepName.toLowerCase().includes("supplement") ||
      stepName.toLowerCase().includes("credit")
    ) {
      return "text-green-600 dark:text-green-400";
    }
    return "text-blue-600 dark:text-blue-400";
  }
</script>

{#if auditTrail}
  <div
    class="calculation-audit-trail bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700"
  >
    <!-- Header -->
    <div class="p-4 border-b border-gray-200 dark:border-gray-700">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            Calculation Audit Trail
          </h3>
          <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
            {auditTrail.calculationType}
          </p>
          <p class="text-xs text-gray-500 dark:text-gray-500 mt-1">
            {auditTrail.inputSummary}
          </p>
        </div>
        <div class="flex items-center space-x-2">
          <button
            onclick={() => toggleAllSteps()}
            class="px-3 py-1 text-xs font-medium text-blue-600 dark:text-blue-400
                 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-md transition-colors"
          >
            {expandedSteps.size === auditTrail.steps.length
              ? "Collapse All"
              : "Expand All"}
          </button>
        </div>
      </div>

      <!-- Controls -->
      {#if showDetailed}
        <div class="flex items-center space-x-4 mt-3">
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={showFormulas}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700 dark:text-gray-300"
              >Show Formulas</span
            >
          </label>
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={showInputs}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700 dark:text-gray-300"
              >Show Input Values</span
            >
          </label>
        </div>
      {/if}
    </div>

    <!-- Final Result Summary -->
    <div
      class="p-4 bg-blue-50 dark:bg-blue-900/20 border-b border-gray-200 dark:border-gray-700"
    >
      <div class="flex items-center justify-between">
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300"
          >Final Monthly Pension:</span
        >
        <span class="text-lg font-bold text-blue-600 dark:text-blue-400">
          {formatCurrency(auditTrail.finalResult)}
        </span>
      </div>
    </div>

    <!-- Calculation Steps -->
    <div class="divide-y divide-gray-200 dark:divide-gray-700">
      {#each auditTrail.steps as step}
        <div class="audit-step">
          <!-- Step Header -->
          <button
            onclick={() => toggleStep(step.stepNumber)}
            class="w-full p-4 text-left hover:bg-gray-50 dark:hover:bg-gray-700/50
               transition-colors focus:outline-none focus:bg-gray-50 dark:focus:bg-gray-700/50"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <div
                  class="flex-shrink-0 w-8 h-8 bg-blue-100 dark:bg-blue-900/40
                        rounded-full flex items-center justify-center"
                >
                  <span
                    class="text-sm font-medium text-blue-600 dark:text-blue-400"
                  >
                    {step.stepNumber}
                  </span>
                </div>
                <div>
                  <h4
                    class="text-sm font-medium {getSeverityColor(
                      step.stepName
                    )}"
                  >
                    {step.stepName}
                  </h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    {step.description}
                  </p>
                </div>
              </div>
              <div class="flex items-center space-x-3">
                <span class="text-sm font-medium text-gray-900 dark:text-white">
                  {typeof step.result === "number"
                    ? formatCurrency(step.result)
                    : step.result}
                </span>
                <svg
                  class="w-4 h-4 text-gray-400 transform transition-transform duration-200
                     {expandedSteps.has(step.stepNumber) ? 'rotate-180' : ''}"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M19 9l-7 7-7-7"
                  />
                </svg>
              </div>
            </div>
          </button>

          <!-- Expanded Step Details -->
          {#if expandedSteps.has(step.stepNumber)}
            <div class="px-4 pb-4 ml-11 space-y-3">
              <!-- Formula -->
              {#if showFormulas && step.formula}
                <div class="bg-gray-50 dark:bg-gray-700/50 rounded-md p-3">
                  <h5
                    class="text-xs font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    Formula:
                  </h5>
                  <code
                    class="text-xs text-gray-800 dark:text-gray-200 font-mono"
                  >
                    {step.formula}
                  </code>
                </div>
              {/if}

              <!-- Inputs -->
              {#if showInputs && step.inputs && Object.keys(step.inputs).length > 0}
                <div class="bg-gray-50 dark:bg-gray-700/50 rounded-md p-3">
                  <h5
                    class="text-xs font-medium text-gray-700 dark:text-gray-300 mb-2"
                  >
                    Input Values:
                  </h5>
                  <div class="grid grid-cols-2 gap-2">
                    {#each Object.entries(step.inputs) as [key, value]}
                      <div class="flex justify-between">
                        <span class="text-xs text-gray-600 dark:text-gray-400"
                          >{key}:</span
                        >
                        <span
                          class="text-xs text-gray-800 dark:text-gray-200 font-mono"
                        >
                          {typeof value === "number"
                            ? formatNumber(value)
                            : value}
                        </span>
                      </div>
                    {/each}
                  </div>
                </div>
              {/if}

              <!-- Calculation Details -->
              {#if step.calculation}
                <div class="bg-blue-50 dark:bg-blue-900/20 rounded-md p-3">
                  <h5
                    class="text-xs font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    Calculation:
                  </h5>
                  <code
                    class="text-xs text-gray-800 dark:text-gray-200 font-mono whitespace-pre-wrap"
                  >
                    {step.calculation}
                  </code>
                </div>
              {/if}

              <!-- Notes -->
              {#if step.notes}
                <div class="bg-amber-50 dark:bg-amber-900/20 rounded-md p-3">
                  <h5
                    class="text-xs font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    Notes:
                  </h5>
                  <p class="text-xs text-gray-600 dark:text-gray-400">
                    {step.notes}
                  </p>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Warnings -->
    {#if auditTrail.warnings && auditTrail.warnings.length > 0}
      <div
        class="p-4 bg-amber-50 dark:bg-amber-900/20 border-t border-gray-200 dark:border-gray-700"
      >
        <h4 class="text-sm font-medium text-amber-800 dark:text-amber-300 mb-2">
          ⚠️ Important Considerations
        </h4>
        <ul class="space-y-1">
          {#each auditTrail.warnings as warning}
            <li class="text-xs text-amber-700 dark:text-amber-400">
              • {warning}
            </li>
          {/each}
        </ul>
      </div>
    {/if}

    <!-- OPM References -->
    {#if auditTrail.ompReferences && auditTrail.ompReferences.length > 0}
      <div
        class="p-4 bg-gray-50 dark:bg-gray-700/50 border-t border-gray-200 dark:border-gray-700"
      >
        <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          📚 References
        </h4>
        <ul class="space-y-1">
          {#each auditTrail.ompReferences as reference}
            <li class="text-xs text-gray-600 dark:text-gray-400">
              • {reference}
            </li>
          {/each}
        </ul>
      </div>
    {/if}
  </div>
{:else}
  <div class="p-8 text-center text-gray-500 dark:text-gray-400">
    <p class="text-sm">No audit trail available for this calculation.</p>
    <p class="text-xs mt-1">
      Run calculation with detailed audit trail enabled to see step-by-step
      breakdown.
    </p>
  </div>
{/if}

<style>
  .calculation-audit-trail {
    max-height: 80vh;
    overflow-y: auto;
  }

  .audit-step:last-child {
    border-bottom: none;
  }
  /* Custom scrollbar for dark mode */
  .calculation-audit-trail::-webkit-scrollbar {
    width: 6px;
  }

  .calculation-audit-trail::-webkit-scrollbar-track {
    background: #f3f4f6;
  }

  :global(.dark) .calculation-audit-trail::-webkit-scrollbar-track {
    background: #374151;
  }

  .calculation-audit-trail::-webkit-scrollbar-thumb {
    background: #d1d5db;
    border-radius: 9999px;
  }

  :global(.dark) .calculation-audit-trail::-webkit-scrollbar-thumb {
    background: #4b5563;
  }

  .calculation-audit-trail::-webkit-scrollbar-thumb:hover {
    background: #9ca3af;
  }

  :global(.dark) .calculation-audit-trail::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
  }
</style>
