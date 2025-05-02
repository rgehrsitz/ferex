export interface PensionData {
  retirementSystem: 'FERS' | 'CSRS' | 'CSRS_OFFSET';
  highThreeSalary: number;
  yearsOfService: number;
  retirementAge: number;
  unusedSickLeave: number; // in hours
  militaryService: number; // in years
  isPartTime: boolean;
  partTimeProrationFactor: number; // 0.1 to 1.0
  csrsOffset: boolean; // Only relevant for CSRS_OFFSET
  survivorBenefit: 'none' | 'partial' | 'full';
}

export interface SocialSecurityData {
  startAge: number;
  estimatedMonthlyBenefit: number;
  isEligible: boolean;
  birthYear: number;
  birthMonth: number;
}

export interface TSPData {
  currentBalance: number;
  traditionalBalance: number; // Pre-tax
  rothBalance: number; // After-tax
  annualContribution: number;
  expectedReturnRate: number; // As decimal (e.g., 0.06 for 6%)
  withdrawalStrategy: 'fixed' | 'rmd' | 'percentage';
  fixedWithdrawalAmount?: number;
  withdrawalPercentage?: number;
  withdrawalStartAge: number;
}

export interface TaxData {
  filingStatus: 'single' | 'married_joint' | 'married_separate' | 'head_of_household';
  stateOfResidence: string;
  additionalIncome: number;
  additionalDeductions: number;
  stateIncomeTaxRate: number; // As decimal (e.g., 0.05 for 5%)
}

export interface COLAData {
  assumedInflationRate: number; // As decimal (e.g., 0.025 for 2.5%)
  applyCOLAToPension: boolean;
  applyColaToSocialSecurity: boolean;
}

export interface OtherIncomeData {
  sources: OtherIncomeSource[];
}

export interface OtherIncomeSource {
  id: string;
  name: string;
  amount: number;
  frequency: 'monthly' | 'annual';
  startAge: number;
  endAge?: number;
  applyCola: boolean;
}

export interface ScenarioData {
  pension: PensionData;
  socialSecurity?: SocialSecurityData;
  tsp?: TSPData;
  tax?: TaxData;
  cola?: COLAData;
  otherIncome?: OtherIncomeData;
}

export interface Scenario {
  id: number;
  name: string;
  data: ScenarioData;
}
