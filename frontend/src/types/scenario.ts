export interface ScenarioData {
  pension?: any;
  socialSecurity?: any;
  tsp?: any;
  tax?: any;
  cola?: any;
  otherIncome?: any;
}
export interface Scenario {
  id: number;
  name: string;
  data: ScenarioData;
}
