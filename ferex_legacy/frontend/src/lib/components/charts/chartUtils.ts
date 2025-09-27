import {
    Chart,
    CategoryScale,
    LinearScale,
    BarElement,
    LineElement,
    PointElement,
    ArcElement,
    Title,
    Tooltip,
    Legend,
    Filler,
    LineController,
    BarController,
    DoughnutController,
    RadialLinearScale,
    RadarController // Import RadarController
} from 'chart.js';

// Register Chart.js components
Chart.register(
    CategoryScale,
    LinearScale,
    BarElement,
    LineElement,
    PointElement,
    ArcElement,
    Title,
    Tooltip,
    Legend,
    Filler,
    LineController,
    BarController,
    DoughnutController,
    RadialLinearScale,
    RadarController // Register RadarController
);

export { Chart };

// Chart themes for consistency
export const chartColors = {
    pension: '#3B82F6',      // Blue
    socialSecurity: '#10B981', // Green
    tsp: '#8B5CF6',          // Purple
    federal: '#EF4444',      // Red
    state: '#F59E0B',        // Amber
    total: '#6B7280',        // Gray
    background: 'rgba(59, 130, 246, 0.1)',
    border: 'rgba(59, 130, 246, 1)'
};

// Common chart options
export const defaultChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
        mode: 'index' as const,
        intersect: false
    },
    plugins: {
        legend: {
            display: true,
            position: 'bottom' as const,
            labels: {
                usePointStyle: true,
                padding: 20,
                font: {
                    size: 12
                }
            }
        },        tooltip: {
            mode: 'index' as const,
            intersect: false,
            backgroundColor: 'rgba(0, 0, 0, 0.8)',
            titleColor: 'white',
            bodyColor: 'white',
            borderColor: 'rgba(255, 255, 255, 0.2)',
            borderWidth: 1,
            cornerRadius: 8,
            displayColors: true,
            callbacks: {
                label: function (context: any) {
                    const label = context.dataset.label || '';
                    const value = context.parsed.y;
                    const isMonthly = context.chart.options.plugins?.displayMode === 'monthly';
                    return `${label}: ${formatCurrencyWithPeriod(value, isMonthly)}`;
                }
            }
        }
    },
    elements: {
        point: {
            radius: 4,
            hoverRadius: 8,
            hitRadius: 10
        },
        line: {
            tension: 0.1
        }
    },    scales: {
        y: {
            beginAtZero: true,
            ticks: {
                callback: function (value: any, index: number, values: any[]) {
                    // For Chart.js context, we need to access the displayMode differently
                    // This will be handled by individual chart components
                    return formatCurrency(value);
                }
            }
        }
    }
};

// Currency formatter for charts
export function formatCurrency (value: number, isMonthly: boolean = false): string {
    if (value === 0 || value == null) return '$0';
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
        maximumFractionDigits: 0
    }).format(value);
}

// Convert between monthly and yearly values
export function convertToMonthly(yearlyValue: number): number {
    return yearlyValue / 12;
}

export function convertToYearly(monthlyValue: number): number {
    return monthlyValue * 12;
}

// Format currency with appropriate label
export function formatCurrencyWithPeriod(value: number, isMonthly: boolean = false): string {
    const formattedValue = formatCurrency(isMonthly ? convertToMonthly(value) : value);
    return `${formattedValue}${isMonthly ? '/mo' : '/yr'}`;
}

// Chart creation utility
export function createChartInstance (canvas: HTMLCanvasElement, config: any): Chart {
    return new Chart(canvas, config);
}

// Chart export utilities
export function exportChartAsImage (chart: Chart, filename: string = 'chart'): void {
    const url = chart.toBase64Image();
    const link = document.createElement('a');
    link.download = `${filename}.png`;
    link.href = url;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

export function exportChartAsPDF (chart: Chart, filename: string = 'chart'): void {
    // For PDF export, we'd need to integrate with a library like jsPDF
    // For now, just export as image
    exportChartAsImage(chart, filename);
}

// Enhanced chart options with dark mode support
export function getResponsiveChartOptions (darkMode: boolean = false) {
    const textColor = darkMode ? '#E5E7EB' : '#374151';
    const gridColor = darkMode ? '#374151' : '#E5E7EB';

    return {
        ...defaultChartOptions,
        plugins: {
            ...defaultChartOptions.plugins,
            legend: {
                ...defaultChartOptions.plugins.legend,
                labels: {
                    color: textColor
                }
            }
        },
        scales: {
            ...defaultChartOptions.scales,
            x: {
                ticks: {
                    color: textColor
                },
                grid: {
                    color: gridColor
                }
            },
            y: {
                ...defaultChartOptions.scales.y,
                ticks: {
                    ...defaultChartOptions.scales.y.ticks,
                    color: textColor
                },
                grid: {
                    color: gridColor
                }
            }
        }
    };
}
