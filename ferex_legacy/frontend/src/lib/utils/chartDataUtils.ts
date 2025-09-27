import type { RetirementCalculationResult, ScenarioInput } from "../../types";
import type { models } from "../../../wailsjs/go/models";

export interface MonthlyDataPoint {
  year: number;
  month: number;
  date: Date;
  ageYears: number;
  ageMonths: number;
  pension: number;
  socialSecurity: number;
  tspWithdrawalTraditional: number;
  tspWithdrawalRoth: number;
  totalTspWithdrawal: number;
  totalPreTaxIncome: number;
  federalTax: number;
  stateTax: number;
  totalTax: number;
  netCashFlow: number;
  notes?: string;
}

export interface YearlyDataPoint {
  year: number;
  age: number;
  ageYears: number; // Add for consistency
  pension: number;
  socialSecurity: number;
  tspWithdrawal: number;
  totalTspWithdrawal: number; // Add for consistency
  totalIncome: number;
  totalPreTaxIncome: number; // Add for consistency
  totalTax: number;
  netIncome: number;
  netCashFlow: number; // Add for consistency
}

export interface ChartDataset {
  label: string;
  data: number[];
  backgroundColor?: string;
  borderColor?: string;
  borderWidth?: number;
  fill?: boolean | string;
  stack?: string;
}

/**
 * Extract monthly projections from backend results
 */
export function extractMonthlyProjections(
  results: RetirementCalculationResult
): MonthlyDataPoint[] {
  if (
    !results.detailedMonthlyProjections ||
    results.detailedMonthlyProjections.length === 0
  ) {
    console.warn("No detailed monthly projections found in results");
    return [];
  }

  return results.detailedMonthlyProjections.map((projection) => ({
    year: projection.year,
    month: projection.month,
    date: new Date(projection.year, projection.month - 1, 1),
    ageYears: projection.ageYears,
    ageMonths: projection.ageMonths,
    pension: projection.pensionForMonth || 0,
    socialSecurity: projection.socialSecurityForMonth || 0,
    tspWithdrawalTraditional: projection.tspWithdrawalTraditionalForMonth || 0,
    tspWithdrawalRoth: projection.tspWithdrawalRothForMonth || 0,
    totalTspWithdrawal:
      (projection.tspWithdrawalTraditionalForMonth || 0) +
      (projection.tspWithdrawalRothForMonth || 0),
    totalPreTaxIncome: projection.totalPreTaxIncomeForMonth || 0,
    federalTax: projection.allocatedFederalTaxForMonth || 0,
    stateTax: projection.allocatedStateTaxForMonth || 0,
    totalTax:
      (projection.allocatedFederalTaxForMonth || 0) +
      (projection.allocatedStateTaxForMonth || 0),
    netCashFlow: projection.netCashFlowForMonth || 0,
    notes: projection.notes,
  }));
}

/**
 * Aggregate monthly data to yearly data
 */
export function aggregateToYearly(
  monthlyData: MonthlyDataPoint[]
): YearlyDataPoint[] {
  const yearlyMap = new Map<number, YearlyDataPoint>();

  monthlyData.forEach((month) => {
    const existing = yearlyMap.get(month.year);
    if (existing) {
      existing.pension += month.pension;
      existing.socialSecurity += month.socialSecurity;
      existing.tspWithdrawal += month.totalTspWithdrawal;
      existing.totalIncome += month.totalPreTaxIncome;
      existing.totalTax += month.totalTax;
      existing.netIncome += month.netCashFlow;
    } else {
      yearlyMap.set(month.year, {
        year: month.year,
        age: month.ageYears,
        ageYears: month.ageYears,
        pension: month.pension,
        socialSecurity: month.socialSecurity,
        tspWithdrawal: month.totalTspWithdrawal,
        totalTspWithdrawal: month.totalTspWithdrawal,
        totalIncome: month.totalPreTaxIncome,
        totalPreTaxIncome: month.totalPreTaxIncome,
        totalTax: month.totalTax,
        netIncome: month.netCashFlow,
        netCashFlow: month.netCashFlow,
      });
    }
  });

  return Array.from(yearlyMap.values()).sort((a, b) => a.year - b.year);
}

/**
 * Create chart datasets for income over time
 */
export function createIncomeDatasets(
  data: MonthlyDataPoint[] | YearlyDataPoint[],
  isMonthly: boolean = false
): ChartDataset[] {
  const datasets: ChartDataset[] = [
    {
      label: "Pension Income",
      data: data.map((d) => d.pension),
      backgroundColor: "rgba(59, 130, 246, 0.2)",
      borderColor: "rgb(59, 130, 246)",
      borderWidth: 2,
      fill: true,
      stack: "income",
    },
    {
      label: "Social Security Income",
      data: data.map((d) => d.socialSecurity),
      backgroundColor: "rgba(16, 185, 129, 0.2)",
      borderColor: "rgb(16, 185, 129)",
      borderWidth: 2,
      fill: true,
      stack: "income",
    },
    {
      label: "TSP Withdrawal Income",
      data: data.map((d) =>
        "tspWithdrawal" in d ? d.tspWithdrawal : d.totalTspWithdrawal
      ),
      backgroundColor: "rgba(139, 92, 246, 0.2)",
      borderColor: "rgb(139, 92, 246)",
      borderWidth: 2,
      fill: true,
      stack: "income",
    },
  ];

  return datasets;
}

/**
 * Create chart labels for time axis
 */
export function createTimeLabels(
  data: MonthlyDataPoint[] | YearlyDataPoint[],
  isMonthly: boolean = false
): string[] {
  if (isMonthly && "month" in data[0]) {
    return (data as MonthlyDataPoint[]).map(
      (d) => `${d.year}-${d.month.toString().padStart(2, "0")}`
    );
  } else {
    return data.map((d) => `${d.year}`);
  }
}

/**
 * Calculate scenario differences for comparison charts
 */
export function calculateScenarioDifferences(
  scenarioA: MonthlyDataPoint[],
  scenarioB: MonthlyDataPoint[]
): MonthlyDataPoint[] {
  const differences: MonthlyDataPoint[] = [];

  const minLength = Math.min(scenarioA.length, scenarioB.length);

  for (let i = 0; i < minLength; i++) {
    const a = scenarioA[i];
    const b = scenarioB[i];

    differences.push({
      year: a.year,
      month: a.month,
      date: a.date,
      ageYears: a.ageYears,
      ageMonths: a.ageMonths,
      pension: b.pension - a.pension,
      socialSecurity: b.socialSecurity - a.socialSecurity,
      tspWithdrawalTraditional:
        b.tspWithdrawalTraditional - a.tspWithdrawalTraditional,
      tspWithdrawalRoth: b.tspWithdrawalRoth - a.tspWithdrawalRoth,
      totalTspWithdrawal: b.totalTspWithdrawal - a.totalTspWithdrawal,
      totalPreTaxIncome: b.totalPreTaxIncome - a.totalPreTaxIncome,
      federalTax: b.federalTax - a.federalTax,
      stateTax: b.stateTax - a.stateTax,
      totalTax: b.totalTax - a.totalTax,
      netCashFlow: b.netCashFlow - a.netCashFlow,
    });
  }

  return differences;
}

/**
 * Get retirement events for timeline markers
 */
export function getRetirementEvents(
  inputs: ScenarioInput
): Array<{ date: Date; label: string; type: string }> {
  const events: Array<{ date: Date; label: string; type: string }> = [];

  if (inputs.plannedRetirementDate) {
    events.push({
      date: new Date(inputs.plannedRetirementDate),
      label: "Retirement",
      type: "retirement",
    });
  }

  if (inputs.dateOfBirth) {
    const birthDate = new Date(inputs.dateOfBirth);
    const age62Date = new Date(
      birthDate.getFullYear() + 62,
      birthDate.getMonth(),
      birthDate.getDate()
    );
    events.push({
      date: age62Date,
      label: "Age 62 (SS Eligible)",
      type: "social-security",
    });
  }

  return events.sort((a, b) => a.date.getTime() - b.date.getTime());
}

/**
 * Format currency for chart tooltips and labels
 */
export function formatChartCurrency(
  value: number,
  isMonthly: boolean = false
): string {
  const formatted = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    maximumFractionDigits: 0,
  }).format(Math.abs(value));

  const sign = value < 0 ? "-" : "";
  const period = isMonthly ? "/mo" : "/yr";

  return `${sign}${formatted}${period}`;
}

/**
 * Create cumulative income data for waterfall charts
 */
export function createCumulativeData(
  data: MonthlyDataPoint[]
): MonthlyDataPoint[] {
  let cumulativePension = 0;
  let cumulativeSS = 0;
  let cumulativeTSP = 0;
  let cumulativeIncome = 0;
  let cumulativeNet = 0;

  return data.map((month) => {
    cumulativePension += month.pension;
    cumulativeSS += month.socialSecurity;
    cumulativeTSP += month.totalTspWithdrawal;
    cumulativeIncome += month.totalPreTaxIncome;
    cumulativeNet += month.netCashFlow;

    return {
      ...month,
      pension: cumulativePension,
      socialSecurity: cumulativeSS,
      totalTspWithdrawal: cumulativeTSP,
      totalPreTaxIncome: cumulativeIncome,
      netCashFlow: cumulativeNet,
    };
  });
}
