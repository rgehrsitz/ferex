export namespace main {
	
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

}

