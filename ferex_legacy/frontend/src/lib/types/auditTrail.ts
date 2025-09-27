// Types for calculation audit trails and step-by-step calculation breakdown

export interface AuditStep {
    stepNumber: number;
    stepName: string;
    description: string;
    formula: string;
    inputs: Record<string, any>;
    calculation: string;
    result: number;
    notes?: string;
}

export interface CalculationAuditTrail {
    calculationType: string; // "FERS", "CSRS", "TSP", etc.
    inputSummary: string;
    steps: AuditStep[];
    finalResult: number;
    warnings?: string[];
    ompReferences?: string[]; // Links to official OPM documentation
}

// Extended FERS result type with audit trail
export interface FERSCalculationResultWithAudit {
    irsSimplifiedMethodExclusion: number;
    monthlyPension: number;
    annualPension: number;
    monthlyEarlyRetirementReduction: number;
    sickLeaveServiceCredit: number;
    prorationApplied: boolean;
    monthlyProratedPension: number;
    monthlySurvivorBenefitReduction: number;
    totalServiceYears: number;
    monthlyBasicAnnuity: number;
    isEligibleForSupplement: boolean;
    fersSupplement?: number;
    retirementType: string;
    prorationFactor?: number;
    notes?: string;
    auditTrail?: CalculationAuditTrail;
    basicAnnuity: number;
    earlyRetirementReduction: number;
    proratedPension: number;
    survivorBenefitReduction: number;
}
