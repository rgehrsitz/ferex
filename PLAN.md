# Ferex Modernization Plan

## 1. Research & Requirements (in progress)
- [x] Review high-level project brief and legacy documentation (ferex_spec, UI notes, tax rules, etc.).
- [ ] Extract must-have features for MVP: pension modeling (FERS/CSRS), annuity supplement, Social Security, TSP projections, taxation, scenario comparison.
- [ ] Identify shared data model requirements for CLI/TUI and Wails GUI.

## 2. Architecture & Tech Stack
- [ ] Define clean Go package layout: `core` domain models/services, `storage` (if needed), `ui/cli`, `ui/desktop`.
- [ ] Establish configuration & dependency injection approach (simple constructors + interfaces).
- [ ] Specify data serialization format for saved scenarios (.ferex) and ensure backward compatibility where reasonable.

## 3. Core Domain Implementation
- [ ] Implement data structures for user profiles, service history, income streams, tax settings.
- [ ] Port/refactor calculation engines: pension, SRS, Social Security estimator, TSP growth/withdrawals, taxation.
- [ ] Build Monte Carlo framework (pluggable random source, scenario inputs).
- [ ] Implement validation and error reporting utilities.

## 4. Shared Application Services
- [ ] Scenario manager (load/save, comparison, versioning metadata).
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
*Last updated: 2025-02-14*
