import type { ScenarioInput } from '../types';

// Validates the core inputs for a retirement scenario.
// Returns an array of error messages. An empty array indicates no errors.
export function validateScenarioInputs (inputs: ScenarioInput): string[] {
  const errorMessages: string[] = [];

  // Core Requirements
  if (!inputs.calculationSystem) {
    errorMessages.push('Retirement system (FERS or CSRS) must be selected.');
  }
  if (!inputs.dateOfBirth) {
    errorMessages.push('Your date of birth is required.');
  }
  if (!inputs.plannedRetirementDate) {
    errorMessages.push('Your planned retirement date is required.');
  }
  if (inputs.high3Salary === null || inputs.high3Salary <= 0) {
    errorMessages.push('High-3 average salary must be a positive number.');
  }
  if (!inputs.serviceComputationDate) {
    errorMessages.push('Your Service Computation Date (SCD) for leave is required.');
  }

  // Date Sanity Checks
  if (inputs.dateOfBirth && inputs.serviceComputationDate) {
    if (new Date(inputs.dateOfBirth) >= new Date(inputs.serviceComputationDate)) {
      errorMessages.push('Date of Birth must be before Service Computation Date.');
    }
  }
  if (inputs.serviceComputationDate && inputs.plannedRetirementDate) {
    if (new Date(inputs.serviceComputationDate) >= new Date(inputs.plannedRetirementDate)) {
      errorMessages.push('Service Computation Date must be before Planned Retirement Date.');
    }
  }

  // Optional Numeric Field Validations (ensure they are not negative if provided)
  if (inputs.unusedSickLeaveHoursAtRetirement !== null && inputs.unusedSickLeaveHoursAtRetirement < 0) {
    errorMessages.push('Unused sick leave hours cannot be negative.');
  }
  if (inputs.employeeContributions !== null && inputs.employeeContributions < 0) {
    errorMessages.push('Employee retirement contributions (e.g., FERS, CSRS) cannot be negative.');
  }

  // Marital Status & Spouse
  if (inputs.isMarriedAtRetirement && !inputs.spouseBirthDate) {
    errorMessages.push('Spouse\'s date of birth is required if married at retirement.');
  }

  // System Specific Validations
  if (inputs.calculationSystem === 'FERS') {
    if ((inputs.estimatedSSBenefitAt62ForSRS ?? -1) < 0) {
      errorMessages.push('Estimated Social Security benefit at 62 (for SRS calculation) cannot be negative. Enter 0 if not applicable or unknown.');
    }
    if (inputs.didSwitchFromCSRS && !inputs.switchedToFERSDate) {
      errorMessages.push('Date Switched to FERS is required if you indicated a switch from CSRS.');
    }
    if (inputs.survivorBenefitFERS) {
      if (inputs.survivorBenefitFERS.spouseElection === 'InsurableInterest' &&
        (!inputs.survivorBenefitFERS.insurableInterestDetails || !inputs.survivorBenefitFERS.insurableInterestDetails.dateOfBirth)) {
        errorMessages.push('Insurable interest details (including DOB) are required for FERS insurable interest election.');
      }
    }
  } else if (inputs.calculationSystem === 'CSRS') {
    if (inputs.isCSRSOffset) {
      if (inputs.yearsOfOffsetService === null || inputs.yearsOfOffsetService < 0) {
        errorMessages.push('Years of CSRS Offset service cannot be negative.');
      }
      if ((inputs.ssAt62WithOffset ?? -1) < 0) {
        errorMessages.push('Social Security at 62 (with Offset) cannot be negative if CSRS Offset applies. Enter 0 if not applicable.');
      }
      if ((inputs.ssAt62WithoutOffset ?? -1) < 0) {
        errorMessages.push('Social Security at 62 (without Offset) cannot be negative if CSRS Offset applies. Enter 0 if not applicable.');
      }
    }
    if (inputs.survivorBenefitCSRS) {
      if (inputs.survivorBenefitCSRS.election === 'InsurableInterest' &&
        (!inputs.survivorBenefitCSRS.insurableInterestDetails || !inputs.survivorBenefitCSRS.insurableInterestDetails.dateOfBirth)) {
        errorMessages.push('Insurable interest details (including DOB) are required for CSRS insurable interest election.');
      }
      if (inputs.survivorBenefitCSRS.election === 'PartialCustomBase' &&
        ((inputs.survivorBenefitCSRS.customBaseAmountForPartial ?? -1) <= 0)) {
        errorMessages.push('A positive custom base amount is required for CSRS partial survivor benefit.');
      }
    }
  }

  // Service Periods Validation
  if (inputs.servicePeriods) {
    for (const sp of inputs.servicePeriods) {
      if (!sp.startDate || !sp.endDate) {
        errorMessages.push(`Service period entry '${sp.id || 'ID missing'}' is missing start or end date.`);
        continue;
      }
      if (new Date(sp.startDate) >= new Date(sp.endDate)) {
        errorMessages.push(`Service period '${sp.id || 'ID missing'}': Start date must be before end date.`);
      }
      // Corrected: Use sp.civilianServiceType and sp.isPartTime for part-time check
      if (sp.serviceCategory === 'Civilian' && sp.isPartTime &&
        (sp.hoursPerWeekIfPartTime === null || sp.hoursPerWeekIfPartTime <= 0 || sp.hoursPerWeekIfPartTime > 80)) {
        errorMessages.push(`Service period '${sp.id || 'ID missing'}' (Part-Time Civilian) requires valid hours per week (1-80).`);
      }
      // Corrected: Use sp.depositRedepositPaymentStatus for deposit check
      // Assuming 'OwedOrPartiallyPaid' or 'AwaitingDetermination' means deposit is not fully paid or status is pending.
      if ((sp.civilianServiceType === 'CivilianTemporary' || sp.serviceCategory === 'Military') &&
        (sp.depositRedepositPaymentStatus === 'OwedOrPartiallyPaid' || sp.depositRedepositPaymentStatus === 'AwaitingDetermination')) {
        errorMessages.push(`Service period '${sp.id || 'ID missing'}' (${sp.civilianServiceType || sp.militaryServiceType}) requires Deposit Paid status to be 'PaidInFull' or 'NotApplicable'. Current: ${sp.depositRedepositPaymentStatus}.`);
      }
    }
  }

  // LWOP Periods Validation
  if (inputs.lwopPeriods) {
    for (const lwop of inputs.lwopPeriods) {
      if (!lwop.startDate || !lwop.endDate) {
        errorMessages.push(`LWOP period entry '${lwop.id || 'ID missing'}' is missing start or end date.`);
        continue;
      }
      if (new Date(lwop.startDate) >= new Date(lwop.endDate)) {
        errorMessages.push(`LWOP period '${lwop.id || 'ID missing'}': Start date must be before end date.`);
      }
    }
  }

  // Placeholder for future validations (TSP, Insurance, etc.)

  return errorMessages;
}
