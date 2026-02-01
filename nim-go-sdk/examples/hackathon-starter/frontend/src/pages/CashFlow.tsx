import {
    Line,
    BarChart,
    Bar,
    AreaChart,
    Area,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer,
    ComposedChart,
    ReferenceLine,
} from 'recharts'
import { useState, useEffect } from 'react'

interface HistoricalDataPoint {
    day: number
    date: string
    actual_amount: number
    predicted_amount: number
    residual: number
    transaction_count: number
}

interface PredictionDataPoint {
    day: number
    date: string
    predicted_amount: number
}

interface ModelStats {
    equation: string
    weight: number
    bias: number
    r_squared: number
    mse: number
    rmse: number
}

interface Insights {
    trend: string
    total_days: number
    total_amount: number
    total_transactions: number
    avg_amount_per_day: number
    avg_transactions_per_day: number
    min_daily_amount: number
    max_daily_amount: number
    date_range: {
        start: string
        end: string
    }
}

interface CashFlowData {
    model: ModelStats
    insights: Insights
    predictions: PredictionDataPoint[]
    historical_data: HistoricalDataPoint[]
}

export default function CashFlow() {
    const [cashFlowData, setCashFlowData] = useState<CashFlowData | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [predictionDays, setPredictionDays] = useState(30)

    // Mock data for demonstration - in production, this would come from your backend
    useEffect(() => {
        // Simulate API call
        const fetchCashFlowData = async () => {
            try {
                setLoading(true)
                // TODO: Replace with actual API call to your backend
                // const response = await fetch('/api/cash-flow-analysis')
                // const data = await response.json()

                // For now, using mock data
                const mockData: CashFlowData = generateMockData(predictionDays)
                setCashFlowData(mockData)
                setError(null)
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load cash flow data')
            } finally {
                setLoading(false)
            }
        }

        fetchCashFlowData()
    }, [predictionDays])

    if (loading) return <p>Loading cash flow analysis…</p>
    if (error) return <p className="error">Error: {error}</p>
    if (!cashFlowData) return <p>No data available</p>

    const { model, insights, predictions, historical_data } = cashFlowData

    // Combine historical and predictions for continuous view
    const combinedData = [
        ...historical_data.map(d => ({
            day: d.day,
            date: d.date,
            actual: d.actual_amount,
            predicted: d.predicted_amount,
            type: 'historical' as const,
        })),
        ...predictions.map(d => ({
            day: d.day,
            date: d.date,
            actual: null,
            predicted: d.predicted_amount,
            type: 'forecast' as const,
        })),
    ]

    // Trend indicator styling
    const trendColor = insights.trend === 'increasing' ? '#10b981' :
        insights.trend === 'decreasing' ? '#ef4444' : '#6b7280'
    const trendIcon = insights.trend === 'increasing' ? '↗' :
        insights.trend === 'decreasing' ? '↘' : '→'

    return (
        <div className="cash-flow">
            {/* Header Section */}
            <section className="card">
                <h1 className="page-title">Cash Flow Analysis</h1>
                <p className="page-subtitle">
                    Advanced OLS regression analysis of transaction patterns with predictive forecasting
                </p>
            </section>

            {/* KPI Cards */}
            <section className="card kpi">
                <div className="kpi-item">
                    <div className="kpi-label">Trend</div>
                    <div className="kpi-value" style={{ color: trendColor }}>
                        {trendIcon} {insights.trend}
                    </div>
                </div>
                <div className="kpi-item">
                    <div className="kpi-label">Total Flow</div>
                    <div className="kpi-value">${insights.total_amount.toLocaleString()}</div>
                </div>
                <div className="kpi-item">
                    <div className="kpi-label">Avg per Day</div>
                    <div className="kpi-value">${insights.avg_amount_per_day.toFixed(2)}</div>
                </div>
                <div className="kpi-item">
                    <div className="kpi-label">R² Score</div>
                    <div className="kpi-value">{model.r_squared.toFixed(4)}</div>
                </div>
            </section>

            {/* Historical vs Predicted - Main Chart */}
            <section className="card">
                <h2 className="section-title">Historical Data vs Model Predictions</h2>
                <div className="chart">
                    <ResponsiveContainer>
                        <ComposedChart data={historical_data}>
                            <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                            <XAxis
                                dataKey="date"
                                tick={{ fontSize: 11 }}
                                angle={-45}
                                textAnchor="end"
                                height={80}
                            />
                            <YAxis
                                label={{ value: 'Amount ($)', angle: -90, position: 'insideLeft' }}
                            />
                            <Tooltip
                                contentStyle={{
                                    backgroundColor: '#fff',
                                    border: '1px solid #e5e7eb',
                                    borderRadius: '8px',
                                    padding: '12px'
                                }}
                            />
                            <Legend />
                            <Area
                                type="monotone"
                                dataKey="actual_amount"
                                fill="#3b82f6"
                                fillOpacity={0.1}
                                stroke="none"
                                name="Actual Range"
                            />
                            <Line
                                type="monotone"
                                dataKey="actual_amount"
                                stroke="#3b82f6"
                                strokeWidth={2}
                                dot={{ r: 3 }}
                                name="Actual Amount"
                            />
                            <Line
                                type="monotone"
                                dataKey="predicted_amount"
                                stroke="#f59e0b"
                                strokeWidth={2}
                                strokeDasharray="5 5"
                                dot={false}
                                name="Model Prediction"
                            />
                        </ComposedChart>
                    </ResponsiveContainer>
                </div>
            </section>

            {/* Combined Historical + Forecast */}
            <section className="card">
                <h2 className="section-title">
                    Cash Flow Forecast ({predictionDays} days ahead)
                </h2>
                <div className="forecast-controls">
                    <label htmlFor="prediction-days">Forecast Period:</label>
                    <select
                        id="prediction-days"
                        value={predictionDays}
                        onChange={(e) => setPredictionDays(Number(e.target.value))}
                        className="select-input"
                    >
                        <option value={7}>7 days</option>
                        <option value={14}>14 days</option>
                        <option value={30}>30 days</option>
                        <option value={60}>60 days</option>
                        <option value={90}>90 days</option>
                    </select>
                    <span className="forecast-separator">or</span>
                    <label htmlFor="custom-days">Custom:</label>
                    <input
                        id="custom-days"
                        type="number"
                        min="1"
                        max="365"
                        value={predictionDays}
                        onChange={(e) => {
                            const value = Number(e.target.value)
                            if (value >= 1 && value <= 365) {
                                setPredictionDays(value)
                            }
                        }}
                        className="number-input"
                        placeholder="Days"
                    />
                    <span className="forecast-unit">days</span>
                </div>
                <div className="chart">
                    <ResponsiveContainer>
                        <AreaChart data={combinedData}>
                            <defs>
                                <linearGradient id="colorActual" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.8} />
                                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0.1} />
                                </linearGradient>
                                <linearGradient id="colorForecast" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.6} />
                                    <stop offset="95%" stopColor="#10b981" stopOpacity={0.05} />
                                </linearGradient>
                            </defs>
                            <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                            <XAxis
                                dataKey="date"
                                tick={{ fontSize: 11 }}
                                interval={Math.floor(combinedData.length / 10)}
                            />
                            <YAxis
                                label={{ value: 'Amount ($)', angle: -90, position: 'insideLeft' }}
                            />
                            <Tooltip
                                contentStyle={{
                                    backgroundColor: '#fff',
                                    border: '1px solid #e5e7eb',
                                    borderRadius: '8px'
                                }}
                            />
                            <Legend />
                            <ReferenceLine
                                x={historical_data[historical_data.length - 1]?.date}
                                stroke="#6b7280"
                                strokeDasharray="3 3"
                                label="Today"
                            />
                            <Area
                                type="monotone"
                                dataKey="actual"
                                stroke="#3b82f6"
                                strokeWidth={2}
                                fill="url(#colorActual)"
                                name="Historical"
                            />
                            <Area
                                type="monotone"
                                dataKey="predicted"
                                stroke="#10b981"
                                strokeWidth={2}
                                strokeDasharray="5 5"
                                fill="url(#colorForecast)"
                                name="Forecast"
                            />
                        </AreaChart>
                    </ResponsiveContainer>
                </div>
            </section>

            {/* Transaction Volume */}
            <section className="card">
                <h2 className="section-title">Daily Transaction Volume</h2>
                <div className="chart">
                    <ResponsiveContainer>
                        <BarChart data={historical_data}>
                            <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                            <XAxis
                                dataKey="date"
                                tick={{ fontSize: 11 }}
                                interval={Math.floor(historical_data.length / 8)}
                            />
                            <YAxis
                                label={{ value: 'Transaction Count', angle: -90, position: 'insideLeft' }}
                            />
                            <Tooltip />
                            <Bar
                                dataKey="transaction_count"
                                fill="#06b6d4"
                                radius={[4, 4, 0, 0]}
                            />
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </section>

            {/* Insights Summary */}
            <section className="card">
                <h2 className="section-title">Statistical Insights</h2>
                <div className="insights-grid">
                    <div className="insight-card">
                        <div className="insight-label">Analysis Period</div>
                        <div className="insight-value">{insights.total_days} days</div>
                        <div className="insight-detail">
                            {insights.date_range.start} to {insights.date_range.end}
                        </div>
                    </div>
                    <div className="insight-card">
                        <div className="insight-label">Total Transactions</div>
                        <div className="insight-value">{insights.total_transactions}</div>
                        <div className="insight-detail">
                            Avg {insights.avg_transactions_per_day.toFixed(1)} per day
                        </div>
                    </div>
                    <div className="insight-card">
                        <div className="insight-label">Daily Range</div>
                        <div className="insight-value">
                            ${insights.min_daily_amount.toFixed(2)} - ${insights.max_daily_amount.toFixed(2)}
                        </div>
                        <div className="insight-detail">
                            Spread: ${(insights.max_daily_amount - insights.min_daily_amount).toFixed(2)}
                        </div>
                    </div>
                    <div className="insight-card">
                        <div className="insight-label">Model Accuracy</div>
                        <div className="insight-value">
                            {(model.r_squared * 100).toFixed(1)}%
                        </div>
                        <div className="insight-detail">
                            R² coefficient of determination
                        </div>
                    </div>
                </div>
            </section>

            {/* Future Predictions Table */}
            <section className="card">
                <h2 className="section-title">Detailed Forecast</h2>
                <div style={{ overflowX: 'auto' }}>
                    <table className="table">
                        <thead>
                            <tr>
                                <th>Date</th>
                                <th>Day #</th>
                                <th>Predicted Amount</th>
                            </tr>
                        </thead>
                        <tbody>
                            {predictions.slice(0, 10).map((pred) => (
                                <tr key={pred.day}>
                                    <td>{pred.date}</td>
                                    <td>{pred.day}</td>
                                    <td>${pred.predicted_amount.toFixed(2)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
                {predictions.length > 10 && (
                    <p className="table-note">Showing first 10 of {predictions.length} predictions</p>
                )}
            </section>
        </div>
    )
}

// Mock data generator for demonstration
function generateMockData(predictionDays: number): CashFlowData {
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 60) // 60 days of historical data

    const historical_data: HistoricalDataPoint[] = []
    const baseAmount = 1000
    const trend = 5 // increasing trend

    for (let i = 0; i < 60; i++) {
        const date = new Date(startDate)
        date.setDate(date.getDate() + i)

        const predicted = baseAmount + trend * (i + 1) + Math.random() * 50
        const actual = predicted + (Math.random() - 0.5) * 200

        historical_data.push({
            day: i + 1,
            date: date.toISOString().split('T')[0],
            actual_amount: Math.round(actual * 100) / 100,
            predicted_amount: Math.round(predicted * 100) / 100,
            residual: Math.round((actual - predicted) * 100) / 100,
            transaction_count: Math.floor(Math.random() * 20) + 5,
        })
    }

    const predictions: PredictionDataPoint[] = []
    const lastDay = historical_data[historical_data.length - 1].day

    for (let i = 0; i < predictionDays; i++) {
        const date = new Date(startDate)
        date.setDate(date.getDate() + 60 + i)

        const predicted = baseAmount + trend * (lastDay + i + 1) + Math.random() * 50

        predictions.push({
            day: lastDay + i + 1,
            date: date.toISOString().split('T')[0],
            predicted_amount: Math.round(predicted * 100) / 100,
        })
    }

    const totalAmount = historical_data.reduce((sum, d) => sum + d.actual_amount, 0)
    const totalTransactions = historical_data.reduce((sum, d) => sum + d.transaction_count, 0)

    return {
        model: {
            equation: `y = ${baseAmount.toFixed(2)} + ${trend.toFixed(2)}x`,
            weight: trend,
            bias: baseAmount,
            r_squared: 0.8523,
            mse: 2341.56,
            rmse: 48.39,
        },
        insights: {
            trend: 'increasing',
            total_days: 60,
            total_amount: Math.round(totalAmount * 100) / 100,
            total_transactions: totalTransactions,
            avg_amount_per_day: Math.round((totalAmount / 60) * 100) / 100,
            avg_transactions_per_day: Math.round((totalTransactions / 60) * 100) / 100,
            min_daily_amount: Math.min(...historical_data.map(d => d.actual_amount)),
            max_daily_amount: Math.max(...historical_data.map(d => d.actual_amount)),
            date_range: {
                start: historical_data[0].date,
                end: historical_data[historical_data.length - 1].date,
            },
        },
        predictions,
        historical_data,
    }
}
