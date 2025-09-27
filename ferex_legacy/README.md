# FeReX - Federal Retirement Scenario Explorer <img src="ferex.svg" height="75px">

[![License](https://img.shields.io/badge/License-MIT-blue)](#license)

Ferex is an experimental desktop application designed to help U.S. Federal employees model their retirement options.  It combines a Go backend for heavy calculation with a modern Svelte 5 frontend powered by the [Wails](https://wails.io/) framework.

The project includes logic for FERS and CSRS pension calculations, Thrift Savings Plan (TSP) growth and withdrawal strategies, Social Security modelling and Monte Carlo simulations.  Scenario data is stored in simple JSON files with the `.ferex` extension.

## Features

- **Comprehensive scenario model** – each `ScenarioVariant` stores all inputs required for a retirement scenario and is fully self contained【F:backend/scenario/file_format.go†L17-L24】.
- **FERS annuity calculations** using user provided service dates and salary information【F:backend/models/fers.go†L3-L15】.
- **Monte Carlo simulation engine** for projecting investment balances and probability of depletion over time【F:backend/calculation/monte_carlo.go†L10-L12】.
- **Configurable financial parameters** such as TSP contribution limits, L‑Fund compositions and RMD rules loaded from `config/ferex_config.json`【F:config/ferex_config.json†L1-L20】【F:config/ferex_config.json†L140-L157】.
- **Multiple TSP withdrawal strategies** including fixed amount, percentage of balance and IRS minimum distributions【F:docs/tsp_withdrawal_strategies.md†L5-L22】.

## Tech Stack

- **Go 1.23** for backend logic and tests
- **Svelte 5** with TypeScript and Tailwind CSS for the frontend
- **Wails v2.10** to package everything into a desktop application

## Getting Started

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 16+](https://nodejs.org/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Development

Run the application with live reload:

```bash
wails dev
```

### Building

Create a distributable build:

```bash
wails build
```

### Testing

Run all backend unit tests with:

```bash
go test ./...
```

## Project Structure

- `frontend/` – Svelte UI and assets
- `backend/` – calculation packages and scenario file handling
- `config/` – JSON files for financial parameters and tax tables
- `.ferex` files – example scenarios that can be loaded by the app

## Configuration

`config/ferex_config.json` contains adjustable values like TSP contribution limits and tax brackets.  The configuration loader is implemented in `backend/config` and is used throughout calculations.

## License

Ferex is released under the [MIT License](LICENSE).

## Acknowledgements

This project builds on the excellent work of the Wails, Svelte and Tailwind communities.
