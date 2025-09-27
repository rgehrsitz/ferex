# FeReX UI Design & Finalized Decisions

This document consolidates the complete Ferex UI vision, including all refinements and agreed‑upon design decisions drawn from our discussions.

---

## 1. Navigation & Layout

### 1.1 Sidebar

* **Collapsible** between full view (labels + icons) and mini‑mode (icons + tooltips).
* **Single main view** visible at a time. Sidebar items:

  1. Home
  2. My Profile
  3. Scenario Editor
  4. Analysis (sub-tabs: Comparison Dashboard, Risk & Monte Carlo)
  5. Household View
  6. Reports & Exports
  7. Settings
* **Quick‑Slide Panel** (from the right edge, available in all views):

  * Quick Calculator widget
  * Recent Scenarios list

### 1.2 Top Bar

* **Left:** FeReX logo/title
* **Center:** Global Scenario Selector dropdown (“Scenario: <name> ▼”)
* **Right:** Notifications (bell icon with badge count; dropdown with timestamped links); Help/User Menu

#### Help/User Menu Contents

1. Profile (opens My Profile)
2. Documentation & Tutorials
3. Keyboard Shortcuts
4. About & Version
5. Log out

---

## 2. Home Page

* **Key Metrics Cards:** concise stats (e.g., Current Net Income; Projected TSP Balance; Estimated Pension).
* **Mini Sparkline Chart:** net income trend over next 10 years.
* **Recent Scenarios Table:** columns for Name, Date Modified, Key Metric (e.g., Annual Income), Actions (Open, Duplicate, Delete).
* **Collapsible Right‑Tray:** houses Quick Calculator and quick-access Recent Scenarios.

---

## 3. Master Data & Profile

* **My Profile** page serves as the **single source of truth** for foundational data:

  * Personal Info (Name, DOB, Hire Date)
  * Base Service History (full-/part-time, military)
  * Baseline High‑3 & Salary History
  * SSA Estimate Inputs
  * Default Pension Rules (CSRS vs. FERS, proration)
* Data here **pre‑populates** new scenarios; scenario editor can override specific fields.

---

## 4. Scenario Editor

### 4.1 Single‑Scenario Workflow

* **Full‑width tabbed interface** (tabs: Inputs, Assumptions, Results Preview).
* **Sections** under Inputs & Assumptions:

  1. Retirement Date & Options (MRA rules, survivor election)
  2. TSP Strategy (Traditional vs. Roth, withdrawal schedule)
  3. Taxes & COLA Assumptions
  4. Optional Overrides (service deposits, custom inflation volatility)
* **Inline Validation** and **Info Tooltips** for all fields.
* **Controls:**

  * Save as Scenario A
  * Duplicate as Scenario B (activates A/B builder flow)

### 4.2 Optional A/B Builder Flow

* Triggered after duplicating Scenario A.
* **Two‑pane layout**:

  * **Left (A):** Locked or view‑only master fields
  * **Right (B):** Editable variant inputs (only overrides)
* **Responsive stacking:** becomes tabbed view below a width threshold.

---

## 5. Analysis Section

* **One page** (“Analysis”) with **sub-tabs**:

  1. Comparison Dashboard
  2. Risk & Monte Carlo
* Shared header with scenario selector and date filters.

### 5.1 Comparison Dashboard

* **Summary Table:** Avg Annual Net Income, Lifetime Net Income, TSP Depletion Age, Effective Tax Rate (columns: A, B, Δ).
* **Level 1 Charts (visible by default):**

  * Side‑by‑side line chart (annual net income trends)
  * Delta summary bar chart
* **Level 2 Charts (expandable tabs):**

  * Stacked Area (income by source)
  * Year‑by‑Year Waterfall Δ

### 5.2 Risk & Monte Carlo

* **Controls:** runs, seed, horizon, basic mean/volatility sliders.
* **Overview Metrics:** Success % at ages (90, 95, 100), key percentiles.
* **Advanced Tabs:**

  * Fan Chart (percentile bands)
  * Histogram at selected age
  * Tornado Sensitivity Diagram

---

## 6. Household View

* **Member toggle:** Scenario A, Scenario B, Combined.
* **Combined Stacked Area Chart** of household net income.
* **Survivor Transition** visualization.
* **Aggregated Metrics Table** (household income, joint TSP balance).

---

## 7. Reports & Exports

* **Quick Export:** Predefined templates (Summary PDF; Detailed CSV).
* **Report Builder:** Drag‑and‑drop charts/tables into a custom layout with live preview.
* **Export Options per Chart/Table:** PNG (transparent), CSV, PDF.
* **Global “Export All”** for entire report.

---

## 8. Notifications Behavior

* **Persistent** across sessions (stored in local storage); badge count reflects unread alerts.
* Alerts auto‑purge after 30 days or user‑cleared.
* Notifications link directly to context (e.g., Monte Carlo results).

---

## 9. Responsiveness & Platform Focus

* **Desktop‑first**, optimized for typical telework monitors (1280px+ width).
* Responsive behaviors:

  * Sidebar collapses to icons + tooltips.
  * Cards stack vertically.
  * Advanced panels accordioned or hidden behind tabs.
* **Mobile** support deprioritized for v1.

---

## 10. Theme & Styling

* **Light/Dark Mode toggle.**
* **Primary palette:** Navy blue (#1F3A93).
* **Accent palette:** Teal/Green (#2ECC71).
* **Neutrals:** Grays for backgrounds/borders.
* **Charts:** Use color‑blind–friendly schemes (e.g., blue/orange).
* **Components:** 2xl rounded corners on buttons, soft shadows on cards, consistent padding, grid layouts, xl headings, base body text.

---

*This document encapsulates our finalized UI design for FeReX, balancing powerful analytical features with a clear, accessible user experience.*
