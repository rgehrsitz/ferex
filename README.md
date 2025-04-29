# FEREX - Federal Employee Retirement Explorer

![License](https://img.shields.io/badge/license-MIT-blue.svg)

## Overview

FEREX (Federal Employee Retirement Explorer) is a comprehensive desktop application designed to help federal employees calculate and plan their retirement benefits. The application provides accurate calculations for different federal retirement systems, including FERS and CSRS, along with additional retirement planning tools.

## Features

- **Pension Calculators**

  - FERS pension calculation with age and service-based multipliers
  - CSRS and CSRS Offset pension calculations
  - Sick leave credit incorporation
  - Part-time service proration
  - Survivor benefit options

- **Additional Retirement Planning Tools**

  - Social Security benefit estimation
  - Thrift Savings Plan (TSP) projections
  - Health benefit cost calculations
  - Tax impact analysis
  - Survivor benefit calculations

- **Advanced Analysis**
  - Monte Carlo simulations for retirement planning
  - Retirement income projections
  - Cost of living adjustment (COLA) modeling

## Screenshots

_Coming soon_

## Installation

### Prerequisites

- Go 1.23 or later
- Node.js and npm

### Building from Source

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/ferex.git
   cd ferex
   ```

2. Install frontend dependencies:

   ```
   cd frontend
   npm install
   cd ..
   ```

3. Build the application:

   ```
   wails build
   ```

4. The compiled application will be available in the `build/bin` directory.

## Development

### Setting Up the Development Environment

1. Install Wails CLI:

   ```
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

2. Run the application in development mode:
   ```
   wails dev
   ```

This will start a development server with hot reloading for the frontend. You can access the Go backend methods through the browser at http://localhost:34115.

### Project Structure

- `/backend` - Go backend code

  - `/calculation` - Calculation logic for different retirement systems
  - `/models` - Data structures and models
  - `/tests` - Unit tests for the backend

- `/frontend` - Svelte/TypeScript frontend
  - `/src` - Main application code
  - `/components` - Reusable UI components
  - `/wailsjs` - Auto-generated bindings for Go functions

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

Run the test suite with:

```
go test ./backend/tests/...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Wails](https://wails.io/) - The framework used to build this desktop application
- [Svelte](https://svelte.dev/) - The frontend framework
- [Tailwind CSS](https://tailwindcss.com/) - For UI styling

## Contact

Project Link: [https://github.com/yourusername/ferex](https://github.com/yourusername/ferex)
