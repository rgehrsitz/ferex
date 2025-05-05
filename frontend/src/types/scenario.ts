export interface PensionData {
  // Fields that match the backend struct names
  system?: string;
  high3Salary?: number;
  yearsOfService: number;
  ageAtRetirement?: number;
  unusedSickLeaveMonths?: number;
  isPartTime: boolean;
  partTimeProrationFactor: number; // 0.1 to 1.0
  survivorBenefitOption?: string;
  
  // Additional frontend fields
  militaryService?: number; // in years
  csrsOffset?: boolean;
}

export interface SocialSecurityData {
  startAge: number;
  estimatedMonthlyBenefit: number;
  isEligible: boolean;
  birthYear: number;
  birthMonth: number;
  ssaEstimateAt62?: number;
  ssaEstimateAtFRA?: number;
  ssaEstimateAt70?: number;
}

export interface TSPData {
  traditionalBalance: number; // Pre-tax
  rothBalance: number; // After-tax
  contributionRate: number; // Traditional TSP contribution percentage 
  contributionRateRoth: number; // Roth TSP contribution percentage
  expectedReturn: number; // As percentage (e.g., 6 for 6%)
  withdrawalStrategy: 'fixed' | 'rmd' | 'percentage';
  withdrawalRate: number; // Percentage of balance for percentage strategy
  fixedMonthlyWithdrawal: number; // Monthly withdrawal amount for fixed strategy
  withdrawalStartAge: number;
}

export interface TaxData {
  filingStatus: 'single' | 'married_joint' | 'married_separate' | 'head_of_household';
  stateOfResidence: string;
  stateIncomeTaxRate: number; // As decimal (e.g., 0.05 for 5%)
  itemizedDeductions: number;
  federalTaxCredits: number;
  stateTaxCredits: number;
  age: number;
  spouseAge: number;
}

export interface COLAData {
  assumedInflationRate: number; // As percentage (e.g., 2.5 for 2.5%)
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
