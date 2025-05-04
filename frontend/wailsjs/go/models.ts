export namespace main {
	
	export class COLAInput {
	    assumedInflationRate: number;
	    applyCOLAToPension: boolean;
	    applyColaToSocialSecurity: boolean;
	    baseAmount: number;
	    retirementSystem: string;
	    retirementAge: number;
	    isSpecialProvision: boolean;
	    startYear: number;
	    projectionYears: number;
	    monthsInFirstYear: number;
	
	    static createFrom(source: any = {}) {
	        return new COLAInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.assumedInflationRate = source["assumedInflationRate"];
	        this.applyCOLAToPension = source["applyCOLAToPension"];
	        this.applyColaToSocialSecurity = source["applyColaToSocialSecurity"];
	        this.baseAmount = source["baseAmount"];
	        this.retirementSystem = source["retirementSystem"];
	        this.retirementAge = source["retirementAge"];
	        this.isSpecialProvision = source["isSpecialProvision"];
	        this.startYear = source["startYear"];
	        this.projectionYears = source["projectionYears"];
	        this.monthsInFirstYear = source["monthsInFirstYear"];
	    }
	}
	export class COLAYearData {
	    year: number;
	    inflationRate: number;
	    startingAmount: number;
	    colaPercentage: number;
	    adjustedAmount: number;
	    cumulativeGrowth: number;
	
	    static createFrom(source: any = {}) {
	        return new COLAYearData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.inflationRate = source["inflationRate"];
	        this.startingAmount = source["startingAmount"];
	        this.colaPercentage = source["colaPercentage"];
	        this.adjustedAmount = source["adjustedAmount"];
	        this.cumulativeGrowth = source["cumulativeGrowth"];
	    }
	}
	export class COLAResult {
	    baseAmount: number;
	    finalAmount: number;
	    totalGrowthRate: number;
	    effectiveAnnualRate: number;
	    yearlyAdjustments: COLAYearData[];
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new COLAResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.baseAmount = source["baseAmount"];
	        this.finalAmount = source["finalAmount"];
	        this.totalGrowthRate = source["totalGrowthRate"];
	        this.effectiveAnnualRate = source["effectiveAnnualRate"];
	        this.yearlyAdjustments = this.convertValues(source["yearlyAdjustments"], COLAYearData);
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
	
	export class GetSavedScenariosResult {
	    files: string[];
	
	    static createFrom(source: any = {}) {
	        return new GetSavedScenariosResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = source["files"];
	    }
	}
	export class LoadScenarioInput {
	    filename: string;
	
	    static createFrom(source: any = {}) {
	        return new LoadScenarioInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	    }
	}
	export class OtherIncomeSource {
	    id: string;
	    name: string;
	    amount: number;
	    frequency: string;
	    startAge: number;
	    endAge?: number;
	    applyCola: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OtherIncomeSource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.amount = source["amount"];
	        this.frequency = source["frequency"];
	        this.startAge = source["startAge"];
	        this.endAge = source["endAge"];
	        this.applyCola = source["applyCola"];
	    }
	}
	export class OtherIncomeInput {
	    sources: OtherIncomeSource[];
	
	    static createFrom(source: any = {}) {
	        return new OtherIncomeInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sources = this.convertValues(source["sources"], OtherIncomeSource);
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
	export class TaxInput {
	    filingStatus: string;
	    stateOfResidence: string;
	    stateIncomeTaxRate: number;
	    totalIncome: number;
	    pensionIncome: number;
	    socialSecurityIncome: number;
	    tspWithdrawals: number;
	    otherIncome: number;
	    nonTaxableIncome: number;
	    itemizedDeductions: number;
	    federalTaxCredits: number;
	    stateTaxCredits: number;
	    age: number;
	    spouseAge: number;
	
	    static createFrom(source: any = {}) {
	        return new TaxInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filingStatus = source["filingStatus"];
	        this.stateOfResidence = source["stateOfResidence"];
	        this.stateIncomeTaxRate = source["stateIncomeTaxRate"];
	        this.totalIncome = source["totalIncome"];
	        this.pensionIncome = source["pensionIncome"];
	        this.socialSecurityIncome = source["socialSecurityIncome"];
	        this.tspWithdrawals = source["tspWithdrawals"];
	        this.otherIncome = source["otherIncome"];
	        this.nonTaxableIncome = source["nonTaxableIncome"];
	        this.itemizedDeductions = source["itemizedDeductions"];
	        this.federalTaxCredits = source["federalTaxCredits"];
	        this.stateTaxCredits = source["stateTaxCredits"];
	        this.age = source["age"];
	        this.spouseAge = source["spouseAge"];
	    }
	}
	export class TSPInput {
	    currentBalance: number;
	    traditionalBalance: number;
	    rothBalance: number;
	    annualContribution: number;
	    expectedReturnRate: number;
	    withdrawalStrategy: string;
	    fixedWithdrawalAmount: number;
	    withdrawalPercentage: number;
	    withdrawalStartAge: number;
	    birthYear: number;
	    birthMonth: number;
	    retirementAge: number;
	
	    static createFrom(source: any = {}) {
	        return new TSPInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentBalance = source["currentBalance"];
	        this.traditionalBalance = source["traditionalBalance"];
	        this.rothBalance = source["rothBalance"];
	        this.annualContribution = source["annualContribution"];
	        this.expectedReturnRate = source["expectedReturnRate"];
	        this.withdrawalStrategy = source["withdrawalStrategy"];
	        this.fixedWithdrawalAmount = source["fixedWithdrawalAmount"];
	        this.withdrawalPercentage = source["withdrawalPercentage"];
	        this.withdrawalStartAge = source["withdrawalStartAge"];
	        this.birthYear = source["birthYear"];
	        this.birthMonth = source["birthMonth"];
	        this.retirementAge = source["retirementAge"];
	    }
	}
	export class SocialSecurityInput {
	    startAge: number;
	    estimatedMonthlyBenefit: number;
	    isEligible: boolean;
	    birthYear: number;
	    birthMonth: number;
	    currentAge: number;
	    estimatedAnnualSalary: number;
	    yearsWorked: number;
	    userProvidedEstimate62: number;
	    userProvidedEstimateFRA: number;
	    userProvidedEstimate70: number;
	
	    static createFrom(source: any = {}) {
	        return new SocialSecurityInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startAge = source["startAge"];
	        this.estimatedMonthlyBenefit = source["estimatedMonthlyBenefit"];
	        this.isEligible = source["isEligible"];
	        this.birthYear = source["birthYear"];
	        this.birthMonth = source["birthMonth"];
	        this.currentAge = source["currentAge"];
	        this.estimatedAnnualSalary = source["estimatedAnnualSalary"];
	        this.yearsWorked = source["yearsWorked"];
	        this.userProvidedEstimate62 = source["userProvidedEstimate62"];
	        this.userProvidedEstimateFRA = source["userProvidedEstimateFRA"];
	        this.userProvidedEstimate70 = source["userProvidedEstimate70"];
	    }
	}
	export class PensionInput {
	    system: string;
	    high3Salary: number;
	    yearsOfService: number;
	    ageAtRetirement: number;
	    unusedSickLeaveMonths: number;
	    survivorBenefitOption: string;
	    isPartTime: boolean;
	    partTimeProrationFactor: number;
	
	    static createFrom(source: any = {}) {
	        return new PensionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.system = source["system"];
	        this.high3Salary = source["high3Salary"];
	        this.yearsOfService = source["yearsOfService"];
	        this.ageAtRetirement = source["ageAtRetirement"];
	        this.unusedSickLeaveMonths = source["unusedSickLeaveMonths"];
	        this.survivorBenefitOption = source["survivorBenefitOption"];
	        this.isPartTime = source["isPartTime"];
	        this.partTimeProrationFactor = source["partTimeProrationFactor"];
	    }
	}
	export class ScenarioData {
	    pension: PensionInput;
	    socialSecurity: SocialSecurityInput;
	    tsp: TSPInput;
	    tax: TaxInput;
	    cola: COLAInput;
	    otherIncome: OtherIncomeInput;
	
	    static createFrom(source: any = {}) {
	        return new ScenarioData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pension = this.convertValues(source["pension"], PensionInput);
	        this.socialSecurity = this.convertValues(source["socialSecurity"], SocialSecurityInput);
	        this.tsp = this.convertValues(source["tsp"], TSPInput);
	        this.tax = this.convertValues(source["tax"], TaxInput);
	        this.cola = this.convertValues(source["cola"], COLAInput);
	        this.otherIncome = this.convertValues(source["otherIncome"], OtherIncomeInput);
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
	export class Scenario {
	    id: number;
	    name: string;
	    data: ScenarioData;
	
	    static createFrom(source: any = {}) {
	        return new Scenario(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.data = this.convertValues(source["data"], ScenarioData);
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
	export class LoadScenarioResult {
	    success: boolean;
	    message: string;
	    scenarios: Scenario[];
	
	    static createFrom(source: any = {}) {
	        return new LoadScenarioResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.scenarios = this.convertValues(source["scenarios"], Scenario);
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
	
	
	
	export class PensionResult {
	    annualPension: number;
	    monthlyPension: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new PensionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.annualPension = source["annualPension"];
	        this.monthlyPension = source["monthlyPension"];
	        this.notes = source["notes"];
	    }
	}
	export class YearlyProjectionData {
	    age: number;
	    year: number;
	    pensionIncome: number;
	    socialSecurity: number;
	    tspWithdrawal: number;
	    otherIncome: number;
	    totalGrossIncome: number;
	    federalTax: number;
	    stateTax: number;
	    totalTaxes: number;
	    netIncome: number;
	    tspBalance: number;
	
	    static createFrom(source: any = {}) {
	        return new YearlyProjectionData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.age = source["age"];
	        this.year = source["year"];
	        this.pensionIncome = source["pensionIncome"];
	        this.socialSecurity = source["socialSecurity"];
	        this.tspWithdrawal = source["tspWithdrawal"];
	        this.otherIncome = source["otherIncome"];
	        this.totalGrossIncome = source["totalGrossIncome"];
	        this.federalTax = source["federalTax"];
	        this.stateTax = source["stateTax"];
	        this.totalTaxes = source["totalTaxes"];
	        this.netIncome = source["netIncome"];
	        this.tspBalance = source["tspBalance"];
	    }
	}
	export class RetirementProjectionResult {
	    yearlyData: YearlyProjectionData[];
	    totalGrossIncome: number;
	    totalNetIncome: number;
	    totalTaxes: number;
	    maxTSPBalance: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new RetirementProjectionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.yearlyData = this.convertValues(source["yearlyData"], YearlyProjectionData);
	        this.totalGrossIncome = source["totalGrossIncome"];
	        this.totalNetIncome = source["totalNetIncome"];
	        this.totalTaxes = source["totalTaxes"];
	        this.maxTSPBalance = source["maxTSPBalance"];
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
	export class RetirementScenarioInput {
	    pension: PensionInput;
	    socialSecurity: SocialSecurityInput;
	    tsp: TSPInput;
	    tax: TaxInput;
	    cola: COLAInput;
	    otherIncome: OtherIncomeInput;
	    projectionStartAge: number;
	    projectionEndAge: number;
	
	    static createFrom(source: any = {}) {
	        return new RetirementScenarioInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pension = this.convertValues(source["pension"], PensionInput);
	        this.socialSecurity = this.convertValues(source["socialSecurity"], SocialSecurityInput);
	        this.tsp = this.convertValues(source["tsp"], TSPInput);
	        this.tax = this.convertValues(source["tax"], TaxInput);
	        this.cola = this.convertValues(source["cola"], COLAInput);
	        this.otherIncome = this.convertValues(source["otherIncome"], OtherIncomeInput);
	        this.projectionStartAge = source["projectionStartAge"];
	        this.projectionEndAge = source["projectionEndAge"];
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
	export class SaveScenarioInput {
	    scenarios: Scenario[];
	    filename: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveScenarioInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scenarios = this.convertValues(source["scenarios"], Scenario);
	        this.filename = source["filename"];
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
	export class SaveScenarioResult {
	    success: boolean;
	    message: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveScenarioResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.path = source["path"];
	    }
	}
	
	
	
	export class SocialSecurityResult {
	    estimatedMonthlyAt62: number;
	    estimatedMonthlyAtFRA: number;
	    estimatedMonthlyAt70: number;
	    estimatedAnnualAt62: number;
	    estimatedAnnualAtFRA: number;
	    estimatedAnnualAt70: number;
	    claimingAge: number;
	    claimingMonthlyAmount: number;
	    claimingAnnualAmount: number;
	    fullRetirementAge: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new SocialSecurityResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.estimatedMonthlyAt62 = source["estimatedMonthlyAt62"];
	        this.estimatedMonthlyAtFRA = source["estimatedMonthlyAtFRA"];
	        this.estimatedMonthlyAt70 = source["estimatedMonthlyAt70"];
	        this.estimatedAnnualAt62 = source["estimatedAnnualAt62"];
	        this.estimatedAnnualAtFRA = source["estimatedAnnualAtFRA"];
	        this.estimatedAnnualAt70 = source["estimatedAnnualAt70"];
	        this.claimingAge = source["claimingAge"];
	        this.claimingMonthlyAmount = source["claimingMonthlyAmount"];
	        this.claimingAnnualAmount = source["claimingAnnualAmount"];
	        this.fullRetirementAge = source["fullRetirementAge"];
	        this.notes = source["notes"];
	    }
	}
	
	export class TSPYearData {
	    age: number;
	    year: number;
	    startingBalance: number;
	    contributions: number;
	    returns: number;
	    withdrawals: number;
	    endingBalance: number;
	
	    static createFrom(source: any = {}) {
	        return new TSPYearData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.age = source["age"];
	        this.year = source["year"];
	        this.startingBalance = source["startingBalance"];
	        this.contributions = source["contributions"];
	        this.returns = source["returns"];
	        this.withdrawals = source["withdrawals"];
	        this.endingBalance = source["endingBalance"];
	    }
	}
	export class TSPProjectionResult {
	    yearlyData: TSPYearData[];
	    maxBalance: number;
	    finalBalance: number;
	    totalContributions: number;
	    totalWithdrawals: number;
	    totalReturns: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new TSPProjectionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.yearlyData = this.convertValues(source["yearlyData"], TSPYearData);
	        this.maxBalance = source["maxBalance"];
	        this.finalBalance = source["finalBalance"];
	        this.totalContributions = source["totalContributions"];
	        this.totalWithdrawals = source["totalWithdrawals"];
	        this.totalReturns = source["totalReturns"];
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
	
	
	export class TaxResult {
	    federalTax: number;
	    stateTax: number;
	    totalTax: number;
	    effectiveFederalRate: number;
	    effectiveStateRate: number;
	    effectiveTotalRate: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new TaxResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.federalTax = source["federalTax"];
	        this.stateTax = source["stateTax"];
	        this.totalTax = source["totalTax"];
	        this.effectiveFederalRate = source["effectiveFederalRate"];
	        this.effectiveStateRate = source["effectiveStateRate"];
	        this.effectiveTotalRate = source["effectiveTotalRate"];
	        this.notes = source["notes"];
	    }
	}

}

