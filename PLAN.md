# Ferex Modernization Plan

## 1. Research & Requirements (in progress)
- [x] Review high-level project brief and legacy documentation (ferex_spec, UI notes, tax rules, etc.).
- [x] Extract must-have features for MVP: pension modeling (FERS/CSRS), annuity supplement, Social Security, TSP projections, taxation, scenario comparison.
  - Core calculations: high-3 pension engine for FERS & CSRS (including CSRS component transfers), annuity reduction rules, SRS eligibility/earnings test, SSA benefit estimator with claiming options, TSP growth & withdrawal strategies (installments, lump sum, RMDs, Roth vs. Traditional sourcing), Monte Carlo success metrics, consolidated tax engine (federal brackets, SS provisional income, Simplified Method basis recovery).
  - Scenario tooling: profile master data feeding scenarios, configurable retirement dates/survivor options, scenario A/B comparison dashboards, household combined view, notification/events feed, reporting exports (PDF/CSV) for summaries and yearly cash flow.
  - User experience: desktop-first Wails UI with sidebar navigation, scenario selector, quick calculator tray, inline validation/tooltips, responsive layouts with analysis charts (income trends, fan charts, tornado sensitivity).
- [x] Identify shared data model requirements for CLI/TUI and Wails GUI.
  - Entities: `UserProfile` (demographics, baseline service history, SSA inputs, retirement system), `ServicePeriod` (type, hours, deposit status), `CompensationHistory` (salary/high-3 inputs), `Scenario` (links profile + overrides, assumptions, retirement options), `IncomeStream` (pension/TSP/other), `TSPPlan` (balances, allocation, withdrawal policies), `TaxSettings` (filing status, deductions, state rules), `SimulationConfig` (Monte Carlo params), `ReportConfig`.
  - Shared managers: scenario repository (load/save `.ferex` JSON), calculation orchestrator, validation & diagnostics, notification log, export formatter.

## 2. Architecture & Tech Stack
- [x] Define clean Go package layout: `internal/core` domain models, `internal/app` services, `internal/storage` adapters, `cmd/ferex-cli` entry point.
- [x] Establish configuration & dependency injection approach (ScenarioManager constructors with pluggable stores/clocks; CLI wiring via flags).
- [x] Specify data serialization format for saved scenarios (`.ferex` JSON written by file-backed ScenarioStore with future backward-compatibility hooks).

## 3. Core Domain Implementation
- [x] Implement data structures for user profiles, service history, income streams, tax settings.
- [ ] Port/refactor calculation engines: pension, SRS, Social Security estimator, TSP growth/withdrawals, taxation.
- [ ] Build Monte Carlo framework (pluggable random source, scenario inputs).
- [x] Implement validation and error reporting utilities.

## 4. Shared Application Services
- [x] Scenario manager (load/save, comparison, versioning metadata).
- [ ] Reporting output (tabular summaries, yearly cashflow, charts data).
- [ ] Logging and diagnostics for supportability.

## 5. Interfaces
- [ ] Command-line interface (Cobra or Bubble Tea) for scenario management and simulation runs.
- [ ] Desktop GUI using Wails v2 (likely Svelte/React-lite or Vanilla JS) that consumes shared backend services via bindings.
- [ ] Ensure both interfaces share identical core logic with minimal duplication.

## 6. Testing & QA
- [ ] Establish unit and integration test suites for core calculations.
- [ ] Create golden sample scenarios (e.g., `MyFerexScenario`, `Rob`) to validate results against legacy expectations.
- [ ] Configure CI instructions (GitHub Actions suggestion) for future automation.

## 7. Documentation & Developer Experience
- [ ] Update README with build/run instructions for CLI and GUI.
- [ ] Provide API docs (Go doc comments) and architecture overview diagrams.
- [ ] Document contribution guidelines and release checklist.

## 8. Future Enhancements (Backlog)
- [ ] Additional benefit modules (FEHB, FEGLI premiums, state-specific tax tables).
- [ ] Integration with external data sources (SSA earnings records import).
- [ ] Advanced analytics (stress tests, custom dashboards).

---
*Last updated: 2025-02-15*
