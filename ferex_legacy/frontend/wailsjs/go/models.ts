export namespace models {
	
	export class AuditStep {
	    stepNumber: number;
	    stepName: string;
	    description: string;
	    formula: string;
	    inputs: Record<string, any>;
	    calculation: string;
	    result: number;
	    notes?: string;
	
	    static createFrom(source: any = {}) {
	        return new AuditStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stepNumber = source["stepNumber"];
	        this.stepName = source["stepName"];
	        this.description = source["description"];
	        this.formula = source["formula"];
	        this.inputs = source["inputs"];
	        this.calculation = source["calculation"];
	        this.result = source["result"];
	        this.notes = source["notes"];
	    }
	}
	export class COLACalculationInput {
	    InitialAmount: number;
	    COLARate: number;
	    Years: number;
	    COLAPolicy: string;
	    StartYear: number;
	
	    static createFrom(source: any = {}) {
	        return new COLACalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.InitialAmount = source["InitialAmount"];
	        this.COLARate = source["COLARate"];
	        this.Years = source["Years"];
	        this.COLAPolicy = source["COLAPolicy"];
	        this.StartYear = source["StartYear"];
	    }
	}
	export class COLACalculationResult {
	    ProjectedAmounts: number[];
	    FinalAmount: number;
	    CumulativeCOLA: number;
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new COLACalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ProjectedAmounts = source["ProjectedAmounts"];
	        this.FinalAmount = source["FinalAmount"];
	        this.CumulativeCOLA = source["CumulativeCOLA"];
	        this.Notes = source["Notes"];
	    }
	}
	export class CSRSCalculationInput {
	    High3Salary: number;
	    UnusedSickLeaveHours: number;
	    SurvivorBenefitType: string;
	    IsPartTime: boolean;
	    PartTimeProrationFactor: number;
	    EmployeeContributions: number;
	    IsCSRSOffset: boolean;
	    YearsOfOffsetService: number;
	    SSAt62WithOffset: number;
	    SSAt62WithoutOffset: number;
	    DateOfBirth: string;
	    ServiceComputationDate: string;
	    PlannedRetirementDate: string;
	
	    static createFrom(source: any = {}) {
	        return new CSRSCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.High3Salary = source["High3Salary"];
	        this.UnusedSickLeaveHours = source["UnusedSickLeaveHours"];
	        this.SurvivorBenefitType = source["SurvivorBenefitType"];
	        this.IsPartTime = source["IsPartTime"];
	        this.PartTimeProrationFactor = source["PartTimeProrationFactor"];
	        this.EmployeeContributions = source["EmployeeContributions"];
	        this.IsCSRSOffset = source["IsCSRSOffset"];
	        this.YearsOfOffsetService = source["YearsOfOffsetService"];
	        this.SSAt62WithOffset = source["SSAt62WithOffset"];
	        this.SSAt62WithoutOffset = source["SSAt62WithoutOffset"];
	        this.DateOfBirth = source["DateOfBirth"];
	        this.ServiceComputationDate = source["ServiceComputationDate"];
	        this.PlannedRetirementDate = source["PlannedRetirementDate"];
	    }
	}
	export class CSRSCalculationResult {
	    IRSSimplifiedMethodExclusion: number;
	    monthlyGrossAnnuity: number;
	    sickLeaveServiceCredit: number;
	    totalServiceYears: number;
	    isProrated: boolean;
	    monthlyGrossAnnuityBeforeProration: number;
	    monthlyOffsetReduction: number;
	    monthlyNetAnnuity: number;
	    monthlySurvivorBenefitReduction: number;
	    monthlyFinalAnnuity: number;
	    retirementType: string;
	    prorationFactor: number;
	    notes: string;
	    grossAnnuity: number;
	    grossAnnuityBeforeProration: number;
	    offsetReduction: number;
	    netAnnuity: number;
	    survivorBenefitReduction: number;
	    finalAnnuity: number;
	
	    static createFrom(source: any = {}) {
	        return new CSRSCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IRSSimplifiedMethodExclusion = source["IRSSimplifiedMethodExclusion"];
	        this.monthlyGrossAnnuity = source["monthlyGrossAnnuity"];
	        this.sickLeaveServiceCredit = source["sickLeaveServiceCredit"];
	        this.totalServiceYears = source["totalServiceYears"];
	        this.isProrated = source["isProrated"];
	        this.monthlyGrossAnnuityBeforeProration = source["monthlyGrossAnnuityBeforeProration"];
	        this.monthlyOffsetReduction = source["monthlyOffsetReduction"];
	        this.monthlyNetAnnuity = source["monthlyNetAnnuity"];
	        this.monthlySurvivorBenefitReduction = source["monthlySurvivorBenefitReduction"];
	        this.monthlyFinalAnnuity = source["monthlyFinalAnnuity"];
	        this.retirementType = source["retirementType"];
	        this.prorationFactor = source["prorationFactor"];
	        this.notes = source["notes"];
	        this.grossAnnuity = source["grossAnnuity"];
	        this.grossAnnuityBeforeProration = source["grossAnnuityBeforeProration"];
	        this.offsetReduction = source["offsetReduction"];
	        this.netAnnuity = source["netAnnuity"];
	        this.survivorBenefitReduction = source["survivorBenefitReduction"];
	        this.finalAnnuity = source["finalAnnuity"];
	    }
	}
	export class CalculationAuditTrail {
	    calculationType: string;
	    inputSummary: string;
	    steps: AuditStep[];
	    finalResult: number;
	    warnings?: string[];
	    ompReferences?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CalculationAuditTrail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.calculationType = source["calculationType"];
	        this.inputSummary = source["inputSummary"];
	        this.steps = this.convertValues(source["steps"], AuditStep);
	        this.finalResult = source["finalResult"];
	        this.warnings = source["warnings"];
	        this.ompReferences = source["ompReferences"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ServicePeriod {
	    id: string;
	    serviceCategory: string;
	    civilianServiceType?: string;
	    militaryServiceType?: string;
	    depositRedepositPaymentStatus: string;
	    systemDuringService: string;
	    startDate: string;
	    endDate: string;
	    isPartTime: boolean;
	    hoursPerWeekIfPartTime?: number;
	    notes?: string;
	
	    static createFrom(source: any = {}) {
	        return new ServicePeriod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.serviceCategory = source["serviceCategory"];
	        this.civilianServiceType = source["civilianServiceType"];
	        this.militaryServiceType = source["militaryServiceType"];
	        this.depositRedepositPaymentStatus = source["depositRedepositPaymentStatus"];
	        this.systemDuringService = source["systemDuringService"];
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
	        this.isPartTime = source["isPartTime"];
	        this.hoursPerWeekIfPartTime = source["hoursPerWeekIfPartTime"];
	        this.notes = source["notes"];
	    }
	}
	export class FERSCalculationInput {
	    High3Salary: number;
	    UnusedSickLeaveHours: number;
	    EmployeeContributions: number;
	    ServicePeriods: ServicePeriod[];
	    SurvivorBenefitElection: string;
	    expectedSSBenefitAt62: number;
	    DateOfBirth: string;
	    ServiceComputationDate: string;
	    PlannedRetirementDate: string;
	    SwitchedToFERSDate: string;
	
	    static createFrom(source: any = {}) {
	        return new FERSCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.High3Salary = source["High3Salary"];
	        this.UnusedSickLeaveHours = source["UnusedSickLeaveHours"];
	        this.EmployeeContributions = source["EmployeeContributions"];
	        this.ServicePeriods = this.convertValues(source["ServicePeriods"], ServicePeriod);
	        this.SurvivorBenefitElection = source["SurvivorBenefitElection"];
	        this.expectedSSBenefitAt62 = source["expectedSSBenefitAt62"];
	        this.DateOfBirth = source["DateOfBirth"];
	        this.ServiceComputationDate = source["ServiceComputationDate"];
	        this.PlannedRetirementDate = source["PlannedRetirementDate"];
	        this.SwitchedToFERSDate = source["SwitchedToFERSDate"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FERSCalculationResult {
	    IRSSimplifiedMethodExclusion: number;
	    monthlyPension: number;
	    annualPension: number;
	    monthlyEarlyRetirementReduction: number;
	    SickLeaveServiceCredit: number;
	    ProrationApplied: boolean;
	    monthlyProratedPension: number;
	    monthlySurvivorBenefitReduction: number;
	    totalServiceYears: number;
	    monthlyBasicAnnuity: number;
	    isEligibleForSupplement: boolean;
	    fersSupplement: number;
	    retirementType: string;
	    prorationFactor: number;
	    notes: string;
	    auditTrail?: CalculationAuditTrail;
	    basicAnnuity: number;
	    EarlyRetirementReduction: number;
	    ProratedPension: number;
	    SurvivorBenefitReduction: number;
	
	    static createFrom(source: any = {}) {
	        return new FERSCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IRSSimplifiedMethodExclusion = source["IRSSimplifiedMethodExclusion"];
	        this.monthlyPension = source["monthlyPension"];
	        this.annualPension = source["annualPension"];
	        this.monthlyEarlyRetirementReduction = source["monthlyEarlyRetirementReduction"];
	        this.SickLeaveServiceCredit = source["SickLeaveServiceCredit"];
	        this.ProrationApplied = source["ProrationApplied"];
	        this.monthlyProratedPension = source["monthlyProratedPension"];
	        this.monthlySurvivorBenefitReduction = source["monthlySurvivorBenefitReduction"];
	        this.totalServiceYears = source["totalServiceYears"];
	        this.monthlyBasicAnnuity = source["monthlyBasicAnnuity"];
	        this.isEligibleForSupplement = source["isEligibleForSupplement"];
	        this.fersSupplement = source["fersSupplement"];
	        this.retirementType = source["retirementType"];
	        this.prorationFactor = source["prorationFactor"];
	        this.notes = source["notes"];
	        this.auditTrail = this.convertValues(source["auditTrail"], CalculationAuditTrail);
	        this.basicAnnuity = source["basicAnnuity"];
	        this.EarlyRetirementReduction = source["EarlyRetirementReduction"];
	        this.ProratedPension = source["ProratedPension"];
	        this.SurvivorBenefitReduction = source["SurvivorBenefitReduction"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FormerSpouseSurvivorDetails {
	    courtOrderNumber?: string;
	    dateOfCourtOrder?: string;
	    portionAwarded?: string;
	    customAmountAwarded?: number;
	
	    static createFrom(source: any = {}) {
	        return new FormerSpouseSurvivorDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.courtOrderNumber = source["courtOrderNumber"];
	        this.dateOfCourtOrder = source["dateOfCourtOrder"];
	        this.portionAwarded = source["portionAwarded"];
	        this.customAmountAwarded = source["customAmountAwarded"];
	    }
	}
	export class HealthPremiumCalculationInput {
	    FEHBPremium: number;
	    MedicarePremium: number;
	    IncludeFEHB: boolean;
	    IncludeMedicare: boolean;
	    COLARate: number;
	    YearsToProject: number;
	    OtherHealthPremium: number;
	
	    static createFrom(source: any = {}) {
	        return new HealthPremiumCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FEHBPremium = source["FEHBPremium"];
	        this.MedicarePremium = source["MedicarePremium"];
	        this.IncludeFEHB = source["IncludeFEHB"];
	        this.IncludeMedicare = source["IncludeMedicare"];
	        this.COLARate = source["COLARate"];
	        this.YearsToProject = source["YearsToProject"];
	        this.OtherHealthPremium = source["OtherHealthPremium"];
	    }
	}
	export class HealthPremiumCalculationResult {
	    ProjectedPremiums: number[];
	    TotalPremiums: number;
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new HealthPremiumCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ProjectedPremiums = source["ProjectedPremiums"];
	        this.TotalPremiums = source["TotalPremiums"];
	        this.Notes = source["Notes"];
	    }
	}
	export class InsurableInterestDetails {
	    relationshipToEmployee: string;
	    dateOfBirth: string;
	
	    static createFrom(source: any = {}) {
	        return new InsurableInterestDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.relationshipToEmployee = source["relationshipToEmployee"];
	        this.dateOfBirth = source["dateOfBirth"];
	    }
	}
	export class LWOPPeriod {
	    id: string;
	    startDate: string;
	    endDate: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new LWOPPeriod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
	        this.type = source["type"];
	    }
	}
	export class MonteCarloInput {
	    NumSimulations: number;
	    Years: number;
	    InitialBalance: number;
	    AnnualWithdrawal: number;
	    ExpectedReturn: number;
	    ReturnStdDev: number;
	    InflationMean: number;
	    InflationStdDev: number;
	    Seed: number;
	
	    static createFrom(source: any = {}) {
	        return new MonteCarloInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.NumSimulations = source["NumSimulations"];
	        this.Years = source["Years"];
	        this.InitialBalance = source["InitialBalance"];
	        this.AnnualWithdrawal = source["AnnualWithdrawal"];
	        this.ExpectedReturn = source["ExpectedReturn"];
	        this.ReturnStdDev = source["ReturnStdDev"];
	        this.InflationMean = source["InflationMean"];
	        this.InflationStdDev = source["InflationStdDev"];
	        this.Seed = source["Seed"];
	    }
	}
	export class MonteCarloResult {
	    SuccessRate: number;
	    Percentiles: Record<number, number>;
	    YearlyBalances: number[][];
	    DepletionProbabilities: number[];
	
	    static createFrom(source: any = {}) {
	        return new MonteCarloResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SuccessRate = source["SuccessRate"];
	        this.Percentiles = source["Percentiles"];
	        this.YearlyBalances = source["YearlyBalances"];
	        this.DepletionProbabilities = source["DepletionProbabilities"];
	    }
	}
	export class MonthlyFinancialProjection {
	    year: number;
	    month: number;
	    ageYears: number;
	    ageMonths: number;
	    pensionForMonth: number;
	    socialSecurityForMonth: number;
	    tspWithdrawalTraditionalForMonth: number;
	    tspWithdrawalRothForMonth: number;
	    healthPremiumsForMonth?: number;
	    totalPreTaxIncomeForMonth: number;
	    allocatedFederalTaxForMonth: number;
	    allocatedStateTaxForMonth?: number;
	    netCashFlowForMonth: number;
	    notes?: string;
	
	    static createFrom(source: any = {}) {
	        return new MonthlyFinancialProjection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.month = source["month"];
	        this.ageYears = source["ageYears"];
	        this.ageMonths = source["ageMonths"];
	        this.pensionForMonth = source["pensionForMonth"];
	        this.socialSecurityForMonth = source["socialSecurityForMonth"];
	        this.tspWithdrawalTraditionalForMonth = source["tspWithdrawalTraditionalForMonth"];
	        this.tspWithdrawalRothForMonth = source["tspWithdrawalRothForMonth"];
	        this.healthPremiumsForMonth = source["healthPremiumsForMonth"];
	        this.totalPreTaxIncomeForMonth = source["totalPreTaxIncomeForMonth"];
	        this.allocatedFederalTaxForMonth = source["allocatedFederalTaxForMonth"];
	        this.allocatedStateTaxForMonth = source["allocatedStateTaxForMonth"];
	        this.netCashFlowForMonth = source["netCashFlowForMonth"];
	        this.notes = source["notes"];
	    }
	}
	export class OneTimeExpenseEvent {
	    id: string;
	    name: string;
	    amount?: number;
	    date: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new OneTimeExpenseEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.amount = source["amount"];
	        this.date = source["date"];
	        this.type = source["type"];
	    }
	}
	export class OneTimeIncomeEvent {
	    id: string;
	    name: string;
	    amount?: number;
	    date: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new OneTimeIncomeEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.amount = source["amount"];
	        this.date = source["date"];
	        this.type = source["type"];
	    }
	}
	export class OtherRecurringIncomeSource {
	    id: string;
	    name: string;
	    amount?: number;
	    frequency: string;
	    startDate: string;
	    endDate?: string;
	    isInflationAdjusted?: boolean;
	    annualIncreaseRate?: number;
	    subjectToFICA?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OtherRecurringIncomeSource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.amount = source["amount"];
	        this.frequency = source["frequency"];
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
	        this.isInflationAdjusted = source["isInflationAdjusted"];
	        this.annualIncreaseRate = source["annualIncreaseRate"];
	        this.subjectToFICA = source["subjectToFICA"];
	    }
	}
	export class SurvivorBenefitCalculationInput {
	    PensionType: string;
	    InitialAnnuity: number;
	    SurvivorElection: string;
	    SpouseAge: number;
	    RetireeAgeAtDeath: number;
	    COLARate: number;
	    YearsToProject: number;
	    IncludeSSSurvivor: boolean;
	    SSSurvivorAmount: number;
	    IncludeTSP: boolean;
	    TSPBalanceAtDeath: number;
	    OtherSurvivorIncome: number;
	
	    static createFrom(source: any = {}) {
	        return new SurvivorBenefitCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PensionType = source["PensionType"];
	        this.InitialAnnuity = source["InitialAnnuity"];
	        this.SurvivorElection = source["SurvivorElection"];
	        this.SpouseAge = source["SpouseAge"];
	        this.RetireeAgeAtDeath = source["RetireeAgeAtDeath"];
	        this.COLARate = source["COLARate"];
	        this.YearsToProject = source["YearsToProject"];
	        this.IncludeSSSurvivor = source["IncludeSSSurvivor"];
	        this.SSSurvivorAmount = source["SSSurvivorAmount"];
	        this.IncludeTSP = source["IncludeTSP"];
	        this.TSPBalanceAtDeath = source["TSPBalanceAtDeath"];
	        this.OtherSurvivorIncome = source["OtherSurvivorIncome"];
	    }
	}
	export class SocialSecurityCalculationInput {
	    BirthYear: number;
	    CurrentAge: number;
	    EarningsHistory: number[];
	    EstimatedAnnualSalary: number;
	    YearsWorked: number;
	    UserProvidedEstimate62: number;
	    UserProvidedEstimateFRA: number;
	    UserProvidedEstimate70: number;
	    ClaimAge: number;
	
	    static createFrom(source: any = {}) {
	        return new SocialSecurityCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BirthYear = source["BirthYear"];
	        this.CurrentAge = source["CurrentAge"];
	        this.EarningsHistory = source["EarningsHistory"];
	        this.EstimatedAnnualSalary = source["EstimatedAnnualSalary"];
	        this.YearsWorked = source["YearsWorked"];
	        this.UserProvidedEstimate62 = source["UserProvidedEstimate62"];
	        this.UserProvidedEstimateFRA = source["UserProvidedEstimateFRA"];
	        this.UserProvidedEstimate70 = source["UserProvidedEstimate70"];
	        this.ClaimAge = source["ClaimAge"];
	    }
	}
	export class TaxCalculationInput {
	    IRSSimplifiedMethodExclusion: number;
	    FilingStatus: string;
	    TaxYear: number;
	    Age: number;
	    NumberOfDependents: number;
	    IsBlind: boolean;
	    GrossPension: number;
	    TaxablePension: number;
	    TSPWithdrawal: number;
	    TSPRothWithdrawal: number;
	    SocialSecurity: number;
	    OtherTaxableIncome: number;
	    StateOfResidence: string;
	    StateTaxableIncome: number;
	    Deductions: number;
	    TaxCredits: number;
	    TaxLawAssumption: string;
	
	    static createFrom(source: any = {}) {
	        return new TaxCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IRSSimplifiedMethodExclusion = source["IRSSimplifiedMethodExclusion"];
	        this.FilingStatus = source["FilingStatus"];
	        this.TaxYear = source["TaxYear"];
	        this.Age = source["Age"];
	        this.NumberOfDependents = source["NumberOfDependents"];
	        this.IsBlind = source["IsBlind"];
	        this.GrossPension = source["GrossPension"];
	        this.TaxablePension = source["TaxablePension"];
	        this.TSPWithdrawal = source["TSPWithdrawal"];
	        this.TSPRothWithdrawal = source["TSPRothWithdrawal"];
	        this.SocialSecurity = source["SocialSecurity"];
	        this.OtherTaxableIncome = source["OtherTaxableIncome"];
	        this.StateOfResidence = source["StateOfResidence"];
	        this.StateTaxableIncome = source["StateTaxableIncome"];
	        this.Deductions = source["Deductions"];
	        this.TaxCredits = source["TaxCredits"];
	        this.TaxLawAssumption = source["TaxLawAssumption"];
	    }
	}
	export class TSPReturnAssumptions {
	    g?: number;
	    f?: number;
	    c?: number;
	    s?: number;
	    i?: number;
	    overall?: number;
	    useOverallForPre?: boolean;
	    useOverallForPost?: boolean;
	    volatilityG?: number;
	    volatilityF?: number;
	    volatilityC?: number;
	    volatilityS?: number;
	    volatilityI?: number;
	
	    static createFrom(source: any = {}) {
	        return new TSPReturnAssumptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.g = source["g"];
	        this.f = source["f"];
	        this.c = source["c"];
	        this.s = source["s"];
	        this.i = source["i"];
	        this.overall = source["overall"];
	        this.useOverallForPre = source["useOverallForPre"];
	        this.useOverallForPost = source["useOverallForPost"];
	        this.volatilityG = source["volatilityG"];
	        this.volatilityF = source["volatilityF"];
	        this.volatilityC = source["volatilityC"];
	        this.volatilityS = source["volatilityS"];
	        this.volatilityI = source["volatilityI"];
	    }
	}
	export class TSPFundAllocationPercentages {
	    g?: number;
	    f?: number;
	    c?: number;
	    s?: number;
	    i?: number;
	    lFundName?: string;
	
	    static createFrom(source: any = {}) {
	        return new TSPFundAllocationPercentages(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.g = source["g"];
	        this.f = source["f"];
	        this.c = source["c"];
	        this.s = source["s"];
	        this.i = source["i"];
	        this.lFundName = source["lFundName"];
	    }
	}
	export class TSPCalculationInput {
	    baseSalaryForContributions: number;
	    currentAgeYears: number;
	    birthYear: number;
	    retirementAgeYears: number;
	    currentTraditionalBalance: number;
	    currentRothBalance: number;
	    currentLoanBalance?: number;
	    currentAllocation: TSPFundAllocationPercentages;
	    employeeContributionAmount?: number;
	    employeeContributionPercentage?: number;
	    isContributionPercentage: boolean;
	    contributeUntilRetirement: boolean;
	    contributionStopAge?: number;
	    catchUpContributionsEligible: boolean;
	    traditionalContributionAllocationPct: number;
	    rothContributionAllocationPct: number;
	    contributionFundAllocation: TSPFundAllocationPercentages;
	    userReturnAssumptions: TSPReturnAssumptions;
	    expenseRatio?: number;
	    expectedAnnualInflationRate: number;
	    futureAllocationStrategy?: string;
	    postRetirementAllocation?: TSPFundAllocationPercentages;
	    withdrawalStrategy: string;
	    withdrawalFixedAmountValue?: number;
	    withdrawalPercentageValue?: number;
	    withdrawalStartDate: string;
	    withdrawalStartAge?: number;
	    withdrawalOrder?: string;
	    yearsToProjectWithdrawals: number;
	
	    static createFrom(source: any = {}) {
	        return new TSPCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.baseSalaryForContributions = source["baseSalaryForContributions"];
	        this.currentAgeYears = source["currentAgeYears"];
	        this.birthYear = source["birthYear"];
	        this.retirementAgeYears = source["retirementAgeYears"];
	        this.currentTraditionalBalance = source["currentTraditionalBalance"];
	        this.currentRothBalance = source["currentRothBalance"];
	        this.currentLoanBalance = source["currentLoanBalance"];
	        this.currentAllocation = this.convertValues(source["currentAllocation"], TSPFundAllocationPercentages);
	        this.employeeContributionAmount = source["employeeContributionAmount"];
	        this.employeeContributionPercentage = source["employeeContributionPercentage"];
	        this.isContributionPercentage = source["isContributionPercentage"];
	        this.contributeUntilRetirement = source["contributeUntilRetirement"];
	        this.contributionStopAge = source["contributionStopAge"];
	        this.catchUpContributionsEligible = source["catchUpContributionsEligible"];
	        this.traditionalContributionAllocationPct = source["traditionalContributionAllocationPct"];
	        this.rothContributionAllocationPct = source["rothContributionAllocationPct"];
	        this.contributionFundAllocation = this.convertValues(source["contributionFundAllocation"], TSPFundAllocationPercentages);
	        this.userReturnAssumptions = this.convertValues(source["userReturnAssumptions"], TSPReturnAssumptions);
	        this.expenseRatio = source["expenseRatio"];
	        this.expectedAnnualInflationRate = source["expectedAnnualInflationRate"];
	        this.futureAllocationStrategy = source["futureAllocationStrategy"];
	        this.postRetirementAllocation = this.convertValues(source["postRetirementAllocation"], TSPFundAllocationPercentages);
	        this.withdrawalStrategy = source["withdrawalStrategy"];
	        this.withdrawalFixedAmountValue = source["withdrawalFixedAmountValue"];
	        this.withdrawalPercentageValue = source["withdrawalPercentageValue"];
	        this.withdrawalStartDate = source["withdrawalStartDate"];
	        this.withdrawalStartAge = source["withdrawalStartAge"];
	        this.withdrawalOrder = source["withdrawalOrder"];
	        this.yearsToProjectWithdrawals = source["yearsToProjectWithdrawals"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SRSCalculationInput {
	    EstimatedSocialSecurityAt62: number;
	    YearsOfFERSService: number;
	    RetirementAge: number;
	    MRA: number;
	    IsImmediateUnreducedAnnuity: boolean;
	    ProjectedEarnedIncome: number;
	    RetirementYear: number;
	
	    static createFrom(source: any = {}) {
	        return new SRSCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EstimatedSocialSecurityAt62 = source["EstimatedSocialSecurityAt62"];
	        this.YearsOfFERSService = source["YearsOfFERSService"];
	        this.RetirementAge = source["RetirementAge"];
	        this.MRA = source["MRA"];
	        this.IsImmediateUnreducedAnnuity = source["IsImmediateUnreducedAnnuity"];
	        this.ProjectedEarnedIncome = source["ProjectedEarnedIncome"];
	        this.RetirementYear = source["RetirementYear"];
	    }
	}
	export class RetirementCalculationInput {
	    CalculationSystem: string;
	    FERSInput: FERSCalculationInput;
	    CSRSInput: CSRSCalculationInput;
	    SRSInput: SRSCalculationInput;
	    TSPInput: TSPCalculationInput;
	    TaxInput: TaxCalculationInput;
	    SocialSecurityInput: SocialSecurityCalculationInput;
	    COLAInput: COLACalculationInput;
	    SurvivorInput: SurvivorBenefitCalculationInput;
	    HealthInput: HealthPremiumCalculationInput;
	    MonteCarloInput: MonteCarloInput;
	
	    static createFrom(source: any = {}) {
	        return new RetirementCalculationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CalculationSystem = source["CalculationSystem"];
	        this.FERSInput = this.convertValues(source["FERSInput"], FERSCalculationInput);
	        this.CSRSInput = this.convertValues(source["CSRSInput"], CSRSCalculationInput);
	        this.SRSInput = this.convertValues(source["SRSInput"], SRSCalculationInput);
	        this.TSPInput = this.convertValues(source["TSPInput"], TSPCalculationInput);
	        this.TaxInput = this.convertValues(source["TaxInput"], TaxCalculationInput);
	        this.SocialSecurityInput = this.convertValues(source["SocialSecurityInput"], SocialSecurityCalculationInput);
	        this.COLAInput = this.convertValues(source["COLAInput"], COLACalculationInput);
	        this.SurvivorInput = this.convertValues(source["SurvivorInput"], SurvivorBenefitCalculationInput);
	        this.HealthInput = this.convertValues(source["HealthInput"], HealthPremiumCalculationInput);
	        this.MonteCarloInput = this.convertValues(source["MonteCarloInput"], MonteCarloInput);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SurvivorBenefitCalculationResult {
	    InitialSurvivorAnnuity: number;
	    ProjectedAnnuities: number[];
	    TotalSurvivorIncome: number;
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new SurvivorBenefitCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.InitialSurvivorAnnuity = source["InitialSurvivorAnnuity"];
	        this.ProjectedAnnuities = source["ProjectedAnnuities"];
	        this.TotalSurvivorIncome = source["TotalSurvivorIncome"];
	        this.Notes = source["Notes"];
	    }
	}
	export class SocialSecurityCalculationResult {
	    EstimatedAt62: number;
	    EstimatedAtFRA: number;
	    EstimatedAt70: number;
	    ClaimingAge: number;
	    ClaimingAmount: number;
	    Notes: string;
	    projectedAnnualBenefitTimeline?: number[];
	
	    static createFrom(source: any = {}) {
	        return new SocialSecurityCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EstimatedAt62 = source["EstimatedAt62"];
	        this.EstimatedAtFRA = source["EstimatedAtFRA"];
	        this.EstimatedAt70 = source["EstimatedAt70"];
	        this.ClaimingAge = source["ClaimingAge"];
	        this.ClaimingAmount = source["ClaimingAmount"];
	        this.Notes = source["Notes"];
	        this.projectedAnnualBenefitTimeline = source["projectedAnnualBenefitTimeline"];
	    }
	}
	export class TaxCalculationResult {
	    GrossIncome: number;
	    TaxableIncome: number;
	    FederalTaxOwed: number;
	    StateTaxOwed: number;
	    NetAfterTaxIncome: number;
	    EffectiveTaxRate: number;
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new TaxCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.GrossIncome = source["GrossIncome"];
	        this.TaxableIncome = source["TaxableIncome"];
	        this.FederalTaxOwed = source["FederalTaxOwed"];
	        this.StateTaxOwed = source["StateTaxOwed"];
	        this.NetAfterTaxIncome = source["NetAfterTaxIncome"];
	        this.EffectiveTaxRate = source["EffectiveTaxRate"];
	        this.Notes = source["Notes"];
	    }
	}
	export class YearlyTSPWithdrawalDetail {
	    year: number;
	    age: number;
	    beginningBalanceTotal: number;
	    beginningBalanceTrad: number;
	    beginningBalanceRoth: number;
	    growthAmount: number;
	    totalWithdrawn: number;
	    traditionalWithdrawn: number;
	    rothWithdrawn: number;
	    endingBalanceTotal: number;
	    endingBalanceTrad: number;
	    endingBalanceRoth: number;
	    rmdAmount?: number;
	    notes?: string;
	
	    static createFrom(source: any = {}) {
	        return new YearlyTSPWithdrawalDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.age = source["age"];
	        this.beginningBalanceTotal = source["beginningBalanceTotal"];
	        this.beginningBalanceTrad = source["beginningBalanceTrad"];
	        this.beginningBalanceRoth = source["beginningBalanceRoth"];
	        this.growthAmount = source["growthAmount"];
	        this.totalWithdrawn = source["totalWithdrawn"];
	        this.traditionalWithdrawn = source["traditionalWithdrawn"];
	        this.rothWithdrawn = source["rothWithdrawn"];
	        this.endingBalanceTotal = source["endingBalanceTotal"];
	        this.endingBalanceTrad = source["endingBalanceTrad"];
	        this.endingBalanceRoth = source["endingBalanceRoth"];
	        this.rmdAmount = source["rmdAmount"];
	        this.notes = source["notes"];
	    }
	}
	export class TSPCalculationResult {
	    projectedTraditionalBalanceAtRetirement: number;
	    projectedRothBalanceAtRetirement: number;
	    totalProjectedBalanceAtRetirement: number;
	    withdrawalSchedule?: YearlyTSPWithdrawalDetail[];
	    notes?: string;
	
	    static createFrom(source: any = {}) {
	        return new TSPCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectedTraditionalBalanceAtRetirement = source["projectedTraditionalBalanceAtRetirement"];
	        this.projectedRothBalanceAtRetirement = source["projectedRothBalanceAtRetirement"];
	        this.totalProjectedBalanceAtRetirement = source["totalProjectedBalanceAtRetirement"];
	        this.withdrawalSchedule = this.convertValues(source["withdrawalSchedule"], YearlyTSPWithdrawalDetail);
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SRSCalculationResult {
	    AnnualSRSAmount: number;
	    MonthlySRSAmount: number;
	    EarningsTestReduction: number;
	    IsEligible: boolean;
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new SRSCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AnnualSRSAmount = source["AnnualSRSAmount"];
	        this.MonthlySRSAmount = source["MonthlySRSAmount"];
	        this.EarningsTestReduction = source["EarningsTestReduction"];
	        this.IsEligible = source["IsEligible"];
	        this.Notes = source["Notes"];
	    }
	}
	export class RetirementCalculationResult {
	    FERSResult: FERSCalculationResult;
	    CSRSResult: CSRSCalculationResult;
	    SRSResult: SRSCalculationResult;
	    TSPResult: TSPCalculationResult;
	    TaxResult: TaxCalculationResult;
	    SocialSecurityResult: SocialSecurityCalculationResult;
	    COLAResult: COLACalculationResult;
	    SurvivorResult: SurvivorBenefitCalculationResult;
	    HealthResult: HealthPremiumCalculationResult;
	    MonteCarloResult: MonteCarloResult;
	    NetAfterTaxIncome: number;
	    EffectiveTaxRate: number;
	    TotalRetirementIncome: number;
	    Notes: string;
	    detailedMonthlyProjections?: MonthlyFinancialProjection[];
	
	    static createFrom(source: any = {}) {
	        return new RetirementCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FERSResult = this.convertValues(source["FERSResult"], FERSCalculationResult);
	        this.CSRSResult = this.convertValues(source["CSRSResult"], CSRSCalculationResult);
	        this.SRSResult = this.convertValues(source["SRSResult"], SRSCalculationResult);
	        this.TSPResult = this.convertValues(source["TSPResult"], TSPCalculationResult);
	        this.TaxResult = this.convertValues(source["TaxResult"], TaxCalculationResult);
	        this.SocialSecurityResult = this.convertValues(source["SocialSecurityResult"], SocialSecurityCalculationResult);
	        this.COLAResult = this.convertValues(source["COLAResult"], COLACalculationResult);
	        this.SurvivorResult = this.convertValues(source["SurvivorResult"], SurvivorBenefitCalculationResult);
	        this.HealthResult = this.convertValues(source["HealthResult"], HealthPremiumCalculationResult);
	        this.MonteCarloResult = this.convertValues(source["MonteCarloResult"], MonteCarloResult);
	        this.NetAfterTaxIncome = source["NetAfterTaxIncome"];
	        this.EffectiveTaxRate = source["EffectiveTaxRate"];
	        this.TotalRetirementIncome = source["TotalRetirementIncome"];
	        this.Notes = source["Notes"];
	        this.detailedMonthlyProjections = this.convertValues(source["detailedMonthlyProjections"], MonthlyFinancialProjection);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class SocialSecurityCreditedEarning {
	    year: number;
	    earnings: number;
	
	    static createFrom(source: any = {}) {
	        return new SocialSecurityCreditedEarning(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.earnings = source["earnings"];
	    }
	}
	export class SurvivorBenefitCSRSInput {
	    election?: string;
	    customBaseAmountForPartial?: number;
	    formerSpouseElection?: string;
	    formerSpouseCustomBaseAmount?: number;
	    formerSpouseConsent?: string;
	    currentSpouseConsentForFormer?: string;
	    currentSpouseWaiverForSelf?: string;
	    insurableInterestDetails?: InsurableInterestDetails;
	
	    static createFrom(source: any = {}) {
	        return new SurvivorBenefitCSRSInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.election = source["election"];
	        this.customBaseAmountForPartial = source["customBaseAmountForPartial"];
	        this.formerSpouseElection = source["formerSpouseElection"];
	        this.formerSpouseCustomBaseAmount = source["formerSpouseCustomBaseAmount"];
	        this.formerSpouseConsent = source["formerSpouseConsent"];
	        this.currentSpouseConsentForFormer = source["currentSpouseConsentForFormer"];
	        this.currentSpouseWaiverForSelf = source["currentSpouseWaiverForSelf"];
	        this.insurableInterestDetails = this.convertValues(source["insurableInterestDetails"], InsurableInterestDetails);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class SurvivorBenefitFERSInput {
	    spouseElection?: string;
	    formerSpouseElection?: string;
	    formerSpouseConsent?: string;
	    currentSpouseConsentForFormer?: string;
	    currentSpouseWaiverForSelf?: string;
	    insurableInterestDetails?: InsurableInterestDetails;
	
	    static createFrom(source: any = {}) {
	        return new SurvivorBenefitFERSInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.spouseElection = source["spouseElection"];
	        this.formerSpouseElection = source["formerSpouseElection"];
	        this.formerSpouseConsent = source["formerSpouseConsent"];
	        this.currentSpouseConsentForFormer = source["currentSpouseConsentForFormer"];
	        this.currentSpouseWaiverForSelf = source["currentSpouseWaiverForSelf"];
	        this.insurableInterestDetails = this.convertValues(source["insurableInterestDetails"], InsurableInterestDetails);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class TSPFundAllocation {
	    G?: number;
	    F?: number;
	    C?: number;
	    S?: number;
	    I?: number;
	    LFundName?: string;
	
	    static createFrom(source: any = {}) {
	        return new TSPFundAllocation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.G = source["G"];
	        this.F = source["F"];
	        this.C = source["C"];
	        this.S = source["S"];
	        this.I = source["I"];
	        this.LFundName = source["LFundName"];
	    }
	}
	
	export class TSPIndividualReturnAssumptions {
	    G?: number;
	    F?: number;
	    C?: number;
	    S?: number;
	    I?: number;
	
	    static createFrom(source: any = {}) {
	        return new TSPIndividualReturnAssumptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.G = source["G"];
	        this.F = source["F"];
	        this.C = source["C"];
	        this.S = source["S"];
	        this.I = source["I"];
	    }
	}
	
	export class TSPVolatilityRates {
	    gStdDev?: number;
	    fStdDev?: number;
	    cStdDev?: number;
	    sStdDev?: number;
	    iStdDev?: number;
	
	    static createFrom(source: any = {}) {
	        return new TSPVolatilityRates(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gStdDev = source["gStdDev"];
	        this.fStdDev = source["fStdDev"];
	        this.cStdDev = source["cStdDev"];
	        this.sStdDev = source["sStdDev"];
	        this.iStdDev = source["iStdDev"];
	    }
	}
	
	
	export class TaxableAccountAssetAllocation {
	    stocksPercent?: number;
	    bondsPercent?: number;
	    cashPercent?: number;
	
	    static createFrom(source: any = {}) {
	        return new TaxableAccountAssetAllocation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stocksPercent = source["stocksPercent"];
	        this.bondsPercent = source["bondsPercent"];
	        this.cashPercent = source["cashPercent"];
	    }
	}

}

export namespace scenario {
	
	export class ScenarioVariant {
	    variantId: string;
	    variantName: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    lastModified: any;
	    notes?: string;
	    calculationSystem: string;
	    high3Salary?: number;
	    serviceComputationDate: string;
	    dateOfBirth: string;
	    plannedRetirementDate: string;
	    unusedSickLeaveHoursAtRetirement?: number;
	    employeeContributions?: number;
	    servicePeriods?: models.ServicePeriod[];
	    lwopPeriods?: models.LWOPPeriod[];
	    isVeraDsRetirement?: boolean;
	    mraPlus10PostponeAnnuityStartDate?: string;
	    isDeferredRetirement?: boolean;
	    deferredRetirementAnnuityStartDate?: string;
	    isDisabilityRetirement?: boolean;
	    isMarriedAtRetirement?: boolean;
	    hasFormerSpouseEntitlement?: boolean;
	    ageOfInsurableInterestBeneficiary?: number;
	    formerSpouseSurvivorBenefitDetails?: models.FormerSpouseSurvivorDetails;
	    fersCoverageType?: string;
	    estimatedSSBenefitAt62ForSRS?: number;
	    didSwitchFromCSRS?: boolean;
	    switchedToFERSDate?: string;
	    survivorBenefitFERS?: models.SurvivorBenefitFERSInput;
	    isCSRSOffset?: boolean;
	    yearsOfOffsetService?: number;
	    ssAt62WithOffset?: number;
	    ssAt62WithoutOffset?: number;
	    survivorBenefitCSRS?: models.SurvivorBenefitCSRSInput;
	    userProvidedSSBenefitAmount1?: number;
	    userProvidedSSBenefitClaimingAge1?: number;
	    userProvidedSSBenefitAmount2?: number;
	    userProvidedSSBenefitClaimingAge2?: number;
	    ssBenefitSpousalOption?: string;
	    ssBenefitSpousalAmount?: number;
	    ssBenefitSurvivorOption?: string;
	    ssBenefitSurvivorAmount?: number;
	    userSSClaimingAge?: number;
	    userAssumedSSCOLA?: number;
	    calculateHistoricalWEPGPO?: boolean;
	    socialSecurityEstimateSource?: string;
	    socialSecurityEstimate?: number;
	    socialSecurityFRA?: number;
	    socialSecurityBenefitStartDate?: string;
	    socialSecurityCreditedEarnings?: models.SocialSecurityCreditedEarning[];
	    socialSecurityBenefitTable?: models.SocialSecurityCreditedEarning[];
	    spouseBirthDate?: string;
	    spouseSSClaimingAge?: number;
	    spouseSocialSecurityEstimate?: number;
	    isSpouseSubjectToWEP?: boolean;
	    spouseWEPNonCoveredPensionAmount?: number;
	    isSpouseSubjectToGPO?: boolean;
	    spouseGPONonCoveredSpousalBenefitAmount?: number;
	    tspBalanceTraditional?: number;
	    tspBalanceRoth?: number;
	    tspLoanBalance?: number;
	    tspAnnualContributionPreRetirement?: number;
	    tspContributionPercentagePreRetirement?: number;
	    tspContributeUntil?: string;
	    tspContributionStopAge?: number;
	    tspCatchUpContributionsEligible?: boolean;
	    tspContributionAllocationTraditionalPercent?: number;
	    tspContributionAllocationRothPercent?: number;
	    tspContributionAllocationToFunds?: models.TSPFundAllocation;
	    tspCurrentAllocationToFunds?: models.TSPFundAllocation;
	    tspPostRetirementAllocation?: models.TSPFundAllocation;
	    tspFutureAllocationStrategy?: string;
	    tspFutureAllocationTarget?: models.TSPFundAllocation;
	    tspLifecycleFundTargetDate?: string;
	    userReturnAssumptionsTSP?: models.TSPIndividualReturnAssumptions;
	    tspExpenseRatio?: number;
	    tspExpectedAnnualGrowthRatePreRetirement?: number;
	    tspWithdrawalStrategy?: string;
	    tspWithdrawalFixedAmountValue?: number;
	    tspWithdrawalPercentageValue?: number;
	    tspWithdrawalStartAge?: number;
	    tspWithdrawalStartDate?: string;
	    tspWithdrawalOrder?: string;
	    applyRMDsToTSP?: boolean;
	    tspExpectedAnnualGrowthRatePostRetirement?: number;
	    tspVolatilityAssumptions?: models.TSPVolatilityRates;
	    tspAnnuityType?: string;
	    tspAnnuitySpouseAge?: number;
	    tspAnnuityFeatures?: string[];
	    tspAnnuityAmountForPurchase?: number;
	    tspAnnuityPurchaseAge?: number;
	    fegliBasicCoverage?: boolean;
	    fegliOptionAAmount?: number;
	    fegliOptionBCoverageMultiples?: number;
	    fegliOptionCCoverageMultiples?: number;
	    fegliPost65Reduction?: string;
	    fegliBasicRetiredCoverage?: string;
	    fegliOptionARetiredCoverage?: string;
	    fegliOptionBRetiredCoverage?: string;
	    fegliOptionCRetiredCoverage?: string;
	    fehbPlanName?: string;
	    fehbPlanCode?: string;
	    fehbCoverageType?: string;
	    fehbPremiumPreRetirementMonthly?: number;
	    fehbPremiumPostRetirementMonthly?: number;
	    fehbContinueInRetirement?: boolean;
	    fehbBiweeklyDeduction?: number;
	    fehbContinuedInRetirement?: boolean;
	    fehbProjectedPremiumIncreaseRate?: number;
	    userAssumedFehbPremiumGrowthRate?: number;
	    fltcipEnrolled?: boolean;
	    fltcipCurrentStatus?: string;
	    fltcipDailyBenefitAmount?: number;
	    fltcipBenefitPeriod?: string;
	    fltcipInflationProtection?: string;
	    fltcipWaitingPeriod?: string;
	    fltcipPremiumMonthly?: number;
	    fltcipContinueInRetirement?: boolean;
	    fedvipDentalPlanName?: string;
	    fedvipDentalPremium?: number;
	    fedvipVisionPlanName?: string;
	    fedvipVisionPremium?: number;
	    otherLifeInsuranceCoverageAmount?: number;
	    otherLifeInsurancePremiumAnnual?: number;
	    disabilityInsuranceCoverageAmount?: number;
	    disabilityInsurancePremiumAnnual?: number;
	    longTermCareInsurancePolicyDetails?: string;
	    federalTaxFilingStatus?: string;
	    federalTaxNumberOfDependents?: number;
	    federalTaxExtraWithholding?: number;
	    federalTaxAdjustmentsToIncome?: number;
	    federalTaxItemizedDeductions?: number;
	    federalOtherTaxableIncomeAnnual?: number;
	    federalTaxCreditsAnnual?: number;
	    federalTaxLawAssumption?: string;
	    userAssumedFederalTaxRateChange?: string;
	    federalTaxUserAssumedRateChangeAmount?: number;
	    federalTaxUserAssumedRateChangeYear?: number;
	    stateOfResidenceForTax?: string;
	    retirementEffectiveStateTaxRate?: number;
	    stateTaxFilingStatus?: string;
	    stateTaxNumberOfDependents?: number;
	    stateTaxExtraWithholding?: number;
	    stateTaxAdjustmentsToIncome?: number;
	    stateTaxItemizedDeductions?: number;
	    userAssumedStateTaxRateChange?: string;
	    stateTaxUserAssumedRateChangeAmount?: number;
	    stateTaxUserAssumedRateChangeYear?: number;
	    taxTreatmentOfMilitaryRetirementPay?: string;
	    otherTaxableAccountBalance?: number;
	    otherTaxableAccountAllocation?: models.TaxableAccountAssetAllocation;
	    otherTaxableAccountAnnualContribution?: number;
	    otherTaxableAccountExpectedGrowth?: number;
	    otherTaxAdvantagedAccountBalance?: number;
	    otherTaxAdvantagedAccountExpectedGrowth?: number;
	    otherTaxAdvantagedAccountAnnualContribution?: number;
	    otherRecurringIncomeSources?: models.OtherRecurringIncomeSource[];
	    otherRecurringMonthlyExpensesPreRetirement?: number;
	    otherRecurringMonthlyExpensesPostRetirement?: number;
	    oneTimeIncomeEvents?: models.OneTimeIncomeEvent[];
	    oneTimeExpenseEvents?: models.OneTimeExpenseEvent[];
	    primaryHomeValue?: number;
	    primaryHomeMortgageBalance?: number;
	    primaryHomeMortgageInterestRate?: number;
	    primaryHomeMortgageRemainingTermYears?: number;
	    primaryHomePropertyTaxRate?: number;
	    primaryHomeInsuranceRate?: number;
	    primaryHomeHoaMonthly?: number;
	    sellPrimaryHomeInRetirement?: boolean;
	    sellPrimaryHomeAge?: number;
	    downsizeOrMoveInRetirement?: boolean;
	    downsizeOrMoveAge?: number;
	    estimatedNetProceedsFromSale?: number;
	    costOfNewHomeOrRental?: number;
	    otherMajorAssets?: string;
	    otherMajorLiabilities?: string;
	    generalInflationRate?: number;
	    userAssumedGeneralInflationRate?: number;
	    userAssumedPensionCOLA_FERS?: number;
	    userAssumedPensionCOLA_CSRS?: number;
	    colaAnnuity?: number;
	    colaSocialSecurity?: number;
	    wageInflationRate?: number;
	    userLifeExpectancy?: number;
	    spouseDateOfBirth?: string;
	    spouseLifeExpectancy?: number;
	    monteCarloEnabled?: boolean;
	    monteCarloNumSimulations?: number;
	    monteCarloConfidenceLevel?: number;
	    desiredAnnualSpendingInRetirement?: number;
	    isSpendingGoalInflationAdjusted?: boolean;
	    healthcareCostInflationRate?: number;
	    userLongTermCareNeed?: string;
	    spouseLongTermCareNeed?: string;
	    longTermCareCostStartAge?: number;
	
	    static createFrom(source: any = {}) {
	        return new ScenarioVariant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.variantId = source["variantId"];
	        this.variantName = source["variantName"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.notes = source["notes"];
	        this.calculationSystem = source["calculationSystem"];
	        this.high3Salary = source["high3Salary"];
	        this.serviceComputationDate = source["serviceComputationDate"];
	        this.dateOfBirth = source["dateOfBirth"];
	        this.plannedRetirementDate = source["plannedRetirementDate"];
	        this.unusedSickLeaveHoursAtRetirement = source["unusedSickLeaveHoursAtRetirement"];
	        this.employeeContributions = source["employeeContributions"];
	        this.servicePeriods = this.convertValues(source["servicePeriods"], models.ServicePeriod);
	        this.lwopPeriods = this.convertValues(source["lwopPeriods"], models.LWOPPeriod);
	        this.isVeraDsRetirement = source["isVeraDsRetirement"];
	        this.mraPlus10PostponeAnnuityStartDate = source["mraPlus10PostponeAnnuityStartDate"];
	        this.isDeferredRetirement = source["isDeferredRetirement"];
	        this.deferredRetirementAnnuityStartDate = source["deferredRetirementAnnuityStartDate"];
	        this.isDisabilityRetirement = source["isDisabilityRetirement"];
	        this.isMarriedAtRetirement = source["isMarriedAtRetirement"];
	        this.hasFormerSpouseEntitlement = source["hasFormerSpouseEntitlement"];
	        this.ageOfInsurableInterestBeneficiary = source["ageOfInsurableInterestBeneficiary"];
	        this.formerSpouseSurvivorBenefitDetails = this.convertValues(source["formerSpouseSurvivorBenefitDetails"], models.FormerSpouseSurvivorDetails);
	        this.fersCoverageType = source["fersCoverageType"];
	        this.estimatedSSBenefitAt62ForSRS = source["estimatedSSBenefitAt62ForSRS"];
	        this.didSwitchFromCSRS = source["didSwitchFromCSRS"];
	        this.switchedToFERSDate = source["switchedToFERSDate"];
	        this.survivorBenefitFERS = this.convertValues(source["survivorBenefitFERS"], models.SurvivorBenefitFERSInput);
	        this.isCSRSOffset = source["isCSRSOffset"];
	        this.yearsOfOffsetService = source["yearsOfOffsetService"];
	        this.ssAt62WithOffset = source["ssAt62WithOffset"];
	        this.ssAt62WithoutOffset = source["ssAt62WithoutOffset"];
	        this.survivorBenefitCSRS = this.convertValues(source["survivorBenefitCSRS"], models.SurvivorBenefitCSRSInput);
	        this.userProvidedSSBenefitAmount1 = source["userProvidedSSBenefitAmount1"];
	        this.userProvidedSSBenefitClaimingAge1 = source["userProvidedSSBenefitClaimingAge1"];
	        this.userProvidedSSBenefitAmount2 = source["userProvidedSSBenefitAmount2"];
	        this.userProvidedSSBenefitClaimingAge2 = source["userProvidedSSBenefitClaimingAge2"];
	        this.ssBenefitSpousalOption = source["ssBenefitSpousalOption"];
	        this.ssBenefitSpousalAmount = source["ssBenefitSpousalAmount"];
	        this.ssBenefitSurvivorOption = source["ssBenefitSurvivorOption"];
	        this.ssBenefitSurvivorAmount = source["ssBenefitSurvivorAmount"];
	        this.userSSClaimingAge = source["userSSClaimingAge"];
	        this.userAssumedSSCOLA = source["userAssumedSSCOLA"];
	        this.calculateHistoricalWEPGPO = source["calculateHistoricalWEPGPO"];
	        this.socialSecurityEstimateSource = source["socialSecurityEstimateSource"];
	        this.socialSecurityEstimate = source["socialSecurityEstimate"];
	        this.socialSecurityFRA = source["socialSecurityFRA"];
	        this.socialSecurityBenefitStartDate = source["socialSecurityBenefitStartDate"];
	        this.socialSecurityCreditedEarnings = this.convertValues(source["socialSecurityCreditedEarnings"], models.SocialSecurityCreditedEarning);
	        this.socialSecurityBenefitTable = this.convertValues(source["socialSecurityBenefitTable"], models.SocialSecurityCreditedEarning);
	        this.spouseBirthDate = source["spouseBirthDate"];
	        this.spouseSSClaimingAge = source["spouseSSClaimingAge"];
	        this.spouseSocialSecurityEstimate = source["spouseSocialSecurityEstimate"];
	        this.isSpouseSubjectToWEP = source["isSpouseSubjectToWEP"];
	        this.spouseWEPNonCoveredPensionAmount = source["spouseWEPNonCoveredPensionAmount"];
	        this.isSpouseSubjectToGPO = source["isSpouseSubjectToGPO"];
	        this.spouseGPONonCoveredSpousalBenefitAmount = source["spouseGPONonCoveredSpousalBenefitAmount"];
	        this.tspBalanceTraditional = source["tspBalanceTraditional"];
	        this.tspBalanceRoth = source["tspBalanceRoth"];
	        this.tspLoanBalance = source["tspLoanBalance"];
	        this.tspAnnualContributionPreRetirement = source["tspAnnualContributionPreRetirement"];
	        this.tspContributionPercentagePreRetirement = source["tspContributionPercentagePreRetirement"];
	        this.tspContributeUntil = source["tspContributeUntil"];
	        this.tspContributionStopAge = source["tspContributionStopAge"];
	        this.tspCatchUpContributionsEligible = source["tspCatchUpContributionsEligible"];
	        this.tspContributionAllocationTraditionalPercent = source["tspContributionAllocationTraditionalPercent"];
	        this.tspContributionAllocationRothPercent = source["tspContributionAllocationRothPercent"];
	        this.tspContributionAllocationToFunds = this.convertValues(source["tspContributionAllocationToFunds"], models.TSPFundAllocation);
	        this.tspCurrentAllocationToFunds = this.convertValues(source["tspCurrentAllocationToFunds"], models.TSPFundAllocation);
	        this.tspPostRetirementAllocation = this.convertValues(source["tspPostRetirementAllocation"], models.TSPFundAllocation);
	        this.tspFutureAllocationStrategy = source["tspFutureAllocationStrategy"];
	        this.tspFutureAllocationTarget = this.convertValues(source["tspFutureAllocationTarget"], models.TSPFundAllocation);
	        this.tspLifecycleFundTargetDate = source["tspLifecycleFundTargetDate"];
	        this.userReturnAssumptionsTSP = this.convertValues(source["userReturnAssumptionsTSP"], models.TSPIndividualReturnAssumptions);
	        this.tspExpenseRatio = source["tspExpenseRatio"];
	        this.tspExpectedAnnualGrowthRatePreRetirement = source["tspExpectedAnnualGrowthRatePreRetirement"];
	        this.tspWithdrawalStrategy = source["tspWithdrawalStrategy"];
	        this.tspWithdrawalFixedAmountValue = source["tspWithdrawalFixedAmountValue"];
	        this.tspWithdrawalPercentageValue = source["tspWithdrawalPercentageValue"];
	        this.tspWithdrawalStartAge = source["tspWithdrawalStartAge"];
	        this.tspWithdrawalStartDate = source["tspWithdrawalStartDate"];
	        this.tspWithdrawalOrder = source["tspWithdrawalOrder"];
	        this.applyRMDsToTSP = source["applyRMDsToTSP"];
	        this.tspExpectedAnnualGrowthRatePostRetirement = source["tspExpectedAnnualGrowthRatePostRetirement"];
	        this.tspVolatilityAssumptions = this.convertValues(source["tspVolatilityAssumptions"], models.TSPVolatilityRates);
	        this.tspAnnuityType = source["tspAnnuityType"];
	        this.tspAnnuitySpouseAge = source["tspAnnuitySpouseAge"];
	        this.tspAnnuityFeatures = source["tspAnnuityFeatures"];
	        this.tspAnnuityAmountForPurchase = source["tspAnnuityAmountForPurchase"];
	        this.tspAnnuityPurchaseAge = source["tspAnnuityPurchaseAge"];
	        this.fegliBasicCoverage = source["fegliBasicCoverage"];
	        this.fegliOptionAAmount = source["fegliOptionAAmount"];
	        this.fegliOptionBCoverageMultiples = source["fegliOptionBCoverageMultiples"];
	        this.fegliOptionCCoverageMultiples = source["fegliOptionCCoverageMultiples"];
	        this.fegliPost65Reduction = source["fegliPost65Reduction"];
	        this.fegliBasicRetiredCoverage = source["fegliBasicRetiredCoverage"];
	        this.fegliOptionARetiredCoverage = source["fegliOptionARetiredCoverage"];
	        this.fegliOptionBRetiredCoverage = source["fegliOptionBRetiredCoverage"];
	        this.fegliOptionCRetiredCoverage = source["fegliOptionCRetiredCoverage"];
	        this.fehbPlanName = source["fehbPlanName"];
	        this.fehbPlanCode = source["fehbPlanCode"];
	        this.fehbCoverageType = source["fehbCoverageType"];
	        this.fehbPremiumPreRetirementMonthly = source["fehbPremiumPreRetirementMonthly"];
	        this.fehbPremiumPostRetirementMonthly = source["fehbPremiumPostRetirementMonthly"];
	        this.fehbContinueInRetirement = source["fehbContinueInRetirement"];
	        this.fehbBiweeklyDeduction = source["fehbBiweeklyDeduction"];
	        this.fehbContinuedInRetirement = source["fehbContinuedInRetirement"];
	        this.fehbProjectedPremiumIncreaseRate = source["fehbProjectedPremiumIncreaseRate"];
	        this.userAssumedFehbPremiumGrowthRate = source["userAssumedFehbPremiumGrowthRate"];
	        this.fltcipEnrolled = source["fltcipEnrolled"];
	        this.fltcipCurrentStatus = source["fltcipCurrentStatus"];
	        this.fltcipDailyBenefitAmount = source["fltcipDailyBenefitAmount"];
	        this.fltcipBenefitPeriod = source["fltcipBenefitPeriod"];
	        this.fltcipInflationProtection = source["fltcipInflationProtection"];
	        this.fltcipWaitingPeriod = source["fltcipWaitingPeriod"];
	        this.fltcipPremiumMonthly = source["fltcipPremiumMonthly"];
	        this.fltcipContinueInRetirement = source["fltcipContinueInRetirement"];
	        this.fedvipDentalPlanName = source["fedvipDentalPlanName"];
	        this.fedvipDentalPremium = source["fedvipDentalPremium"];
	        this.fedvipVisionPlanName = source["fedvipVisionPlanName"];
	        this.fedvipVisionPremium = source["fedvipVisionPremium"];
	        this.otherLifeInsuranceCoverageAmount = source["otherLifeInsuranceCoverageAmount"];
	        this.otherLifeInsurancePremiumAnnual = source["otherLifeInsurancePremiumAnnual"];
	        this.disabilityInsuranceCoverageAmount = source["disabilityInsuranceCoverageAmount"];
	        this.disabilityInsurancePremiumAnnual = source["disabilityInsurancePremiumAnnual"];
	        this.longTermCareInsurancePolicyDetails = source["longTermCareInsurancePolicyDetails"];
	        this.federalTaxFilingStatus = source["federalTaxFilingStatus"];
	        this.federalTaxNumberOfDependents = source["federalTaxNumberOfDependents"];
	        this.federalTaxExtraWithholding = source["federalTaxExtraWithholding"];
	        this.federalTaxAdjustmentsToIncome = source["federalTaxAdjustmentsToIncome"];
	        this.federalTaxItemizedDeductions = source["federalTaxItemizedDeductions"];
	        this.federalOtherTaxableIncomeAnnual = source["federalOtherTaxableIncomeAnnual"];
	        this.federalTaxCreditsAnnual = source["federalTaxCreditsAnnual"];
	        this.federalTaxLawAssumption = source["federalTaxLawAssumption"];
	        this.userAssumedFederalTaxRateChange = source["userAssumedFederalTaxRateChange"];
	        this.federalTaxUserAssumedRateChangeAmount = source["federalTaxUserAssumedRateChangeAmount"];
	        this.federalTaxUserAssumedRateChangeYear = source["federalTaxUserAssumedRateChangeYear"];
	        this.stateOfResidenceForTax = source["stateOfResidenceForTax"];
	        this.retirementEffectiveStateTaxRate = source["retirementEffectiveStateTaxRate"];
	        this.stateTaxFilingStatus = source["stateTaxFilingStatus"];
	        this.stateTaxNumberOfDependents = source["stateTaxNumberOfDependents"];
	        this.stateTaxExtraWithholding = source["stateTaxExtraWithholding"];
	        this.stateTaxAdjustmentsToIncome = source["stateTaxAdjustmentsToIncome"];
	        this.stateTaxItemizedDeductions = source["stateTaxItemizedDeductions"];
	        this.userAssumedStateTaxRateChange = source["userAssumedStateTaxRateChange"];
	        this.stateTaxUserAssumedRateChangeAmount = source["stateTaxUserAssumedRateChangeAmount"];
	        this.stateTaxUserAssumedRateChangeYear = source["stateTaxUserAssumedRateChangeYear"];
	        this.taxTreatmentOfMilitaryRetirementPay = source["taxTreatmentOfMilitaryRetirementPay"];
	        this.otherTaxableAccountBalance = source["otherTaxableAccountBalance"];
	        this.otherTaxableAccountAllocation = this.convertValues(source["otherTaxableAccountAllocation"], models.TaxableAccountAssetAllocation);
	        this.otherTaxableAccountAnnualContribution = source["otherTaxableAccountAnnualContribution"];
	        this.otherTaxableAccountExpectedGrowth = source["otherTaxableAccountExpectedGrowth"];
	        this.otherTaxAdvantagedAccountBalance = source["otherTaxAdvantagedAccountBalance"];
	        this.otherTaxAdvantagedAccountExpectedGrowth = source["otherTaxAdvantagedAccountExpectedGrowth"];
	        this.otherTaxAdvantagedAccountAnnualContribution = source["otherTaxAdvantagedAccountAnnualContribution"];
	        this.otherRecurringIncomeSources = this.convertValues(source["otherRecurringIncomeSources"], models.OtherRecurringIncomeSource);
	        this.otherRecurringMonthlyExpensesPreRetirement = source["otherRecurringMonthlyExpensesPreRetirement"];
	        this.otherRecurringMonthlyExpensesPostRetirement = source["otherRecurringMonthlyExpensesPostRetirement"];
	        this.oneTimeIncomeEvents = this.convertValues(source["oneTimeIncomeEvents"], models.OneTimeIncomeEvent);
	        this.oneTimeExpenseEvents = this.convertValues(source["oneTimeExpenseEvents"], models.OneTimeExpenseEvent);
	        this.primaryHomeValue = source["primaryHomeValue"];
	        this.primaryHomeMortgageBalance = source["primaryHomeMortgageBalance"];
	        this.primaryHomeMortgageInterestRate = source["primaryHomeMortgageInterestRate"];
	        this.primaryHomeMortgageRemainingTermYears = source["primaryHomeMortgageRemainingTermYears"];
	        this.primaryHomePropertyTaxRate = source["primaryHomePropertyTaxRate"];
	        this.primaryHomeInsuranceRate = source["primaryHomeInsuranceRate"];
	        this.primaryHomeHoaMonthly = source["primaryHomeHoaMonthly"];
	        this.sellPrimaryHomeInRetirement = source["sellPrimaryHomeInRetirement"];
	        this.sellPrimaryHomeAge = source["sellPrimaryHomeAge"];
	        this.downsizeOrMoveInRetirement = source["downsizeOrMoveInRetirement"];
	        this.downsizeOrMoveAge = source["downsizeOrMoveAge"];
	        this.estimatedNetProceedsFromSale = source["estimatedNetProceedsFromSale"];
	        this.costOfNewHomeOrRental = source["costOfNewHomeOrRental"];
	        this.otherMajorAssets = source["otherMajorAssets"];
	        this.otherMajorLiabilities = source["otherMajorLiabilities"];
	        this.generalInflationRate = source["generalInflationRate"];
	        this.userAssumedGeneralInflationRate = source["userAssumedGeneralInflationRate"];
	        this.userAssumedPensionCOLA_FERS = source["userAssumedPensionCOLA_FERS"];
	        this.userAssumedPensionCOLA_CSRS = source["userAssumedPensionCOLA_CSRS"];
	        this.colaAnnuity = source["colaAnnuity"];
	        this.colaSocialSecurity = source["colaSocialSecurity"];
	        this.wageInflationRate = source["wageInflationRate"];
	        this.userLifeExpectancy = source["userLifeExpectancy"];
	        this.spouseDateOfBirth = source["spouseDateOfBirth"];
	        this.spouseLifeExpectancy = source["spouseLifeExpectancy"];
	        this.monteCarloEnabled = source["monteCarloEnabled"];
	        this.monteCarloNumSimulations = source["monteCarloNumSimulations"];
	        this.monteCarloConfidenceLevel = source["monteCarloConfidenceLevel"];
	        this.desiredAnnualSpendingInRetirement = source["desiredAnnualSpendingInRetirement"];
	        this.isSpendingGoalInflationAdjusted = source["isSpendingGoalInflationAdjusted"];
	        this.healthcareCostInflationRate = source["healthcareCostInflationRate"];
	        this.userLongTermCareNeed = source["userLongTermCareNeed"];
	        this.spouseLongTermCareNeed = source["spouseLongTermCareNeed"];
	        this.longTermCareCostStartAge = source["longTermCareCostStartAge"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UserProfile {
	    employeeName?: string;
	    birthDate: string;
	    mraYears: number;
	    mraMonths: number;
	
	    static createFrom(source: any = {}) {
	        return new UserProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.employeeName = source["employeeName"];
	        this.birthDate = source["birthDate"];
	        this.mraYears = source["mraYears"];
	        this.mraMonths = source["mraMonths"];
	    }
	}
	export class FerexFile {
	    fileFormatVersion: string;
	    lastOpenedAppVersion?: string;
	    userProfile: UserProfile;
	    variants: ScenarioVariant[];
	    lastViewedVariantID?: string;
	
	    static createFrom(source: any = {}) {
	        return new FerexFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileFormatVersion = source["fileFormatVersion"];
	        this.lastOpenedAppVersion = source["lastOpenedAppVersion"];
	        this.userProfile = this.convertValues(source["userProfile"], UserProfile);
	        this.variants = this.convertValues(source["variants"], ScenarioVariant);
	        this.lastViewedVariantID = source["lastViewedVariantID"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

