import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line,
  CartesianGrid,
  ScatterChart,
  Scatter,
  Legend,
} from 'recharts'
import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { listEmployees } from '../api/employees'

// Adjust this type if your backend uses different field names
type Employee = {
  id?: string
  first_name: string
  last_name: string
  recipient: string
  wage: string | number
  department: string
}

function parseWage(w: Employee['wage']): number {
  if (typeof w === 'number') return w
  // allow "$2,500.00" etc
  const cleaned = String(w).replace(/[^0-9.-]/g, '')
  const n = Number(cleaned)
  return Number.isFinite(n) ? n : NaN
}

function quantile(sorted: number[], q: number): number {
  if (sorted.length === 0) return NaN
  const pos = (sorted.length - 1) * q
  const base = Math.floor(pos)
  const rest = pos - base
  if (sorted[base + 1] === undefined) return sorted[base]
  return sorted[base] + rest * (sorted[base + 1] - sorted[base])
}

export default function Dashboard() {
  const employeesQuery = useQuery<Employee[], Error>({
  queryKey: ['employees-mock'],
  queryFn: async () => {
    const res = await fetch('/mock-employees.json')
    if (!res.ok) throw new Error('Failed to load mock-employees.json')
    return res.json()
  },
})

  const analytics = useMemo(() => {
    const emps = employeesQuery.data ?? []

    const rows = emps.map((e, idx) => {
      const wageNum = parseWage(e.wage)
      return {
        key: e.id ?? `${e.recipient ?? 'emp'}-${idx}`,
        name: `${e.first_name} ${e.last_name}`.trim(),
        recipient: e.recipient?.trim() ?? '',
        department: (e.department?.trim() || 'Unassigned') as string,
        wage: wageNum,
        wageRaw: e.wage,
        wageValid: Number.isFinite(wageNum) && wageNum > 0,
      }
    })

    const valid = rows.filter((r) => r.wageValid)
    const invalid = rows.filter((r) => !r.wageValid)

    // KPIs
    const headcount = rows.length
    const validCount = valid.length
    const totalPayroll = valid.reduce((s, r) => s + r.wage, 0)
    const avgWage = validCount ? totalPayroll / validCount : 0

    // Dept aggregates
    const deptMap = new Map<
      string,
      { department: string; headcount: number; total: number; avg: number; max: number; min: number; wages: number[] }
    >()

    for (const r of valid) {
      const d = r.department
      if (!deptMap.has(d)) {
        deptMap.set(d, { department: d, headcount: 0, total: 0, avg: 0, max: r.wage, min: r.wage, wages: [] })
      }
      const obj = deptMap.get(d)!
      obj.headcount += 1
      obj.total += r.wage
      obj.max = Math.max(obj.max, r.wage)
      obj.min = Math.min(obj.min, r.wage)
      obj.wages.push(r.wage)
    }

    const deptAgg = Array.from(deptMap.values()).map((d) => ({
      ...d,
      avg: d.headcount ? d.total / d.headcount : 0,
    }))

    // 1) Total wage cost by department
    const totalByDept = deptAgg
      .slice()
      .sort((a, b) => b.total - a.total)
      .map((d) => ({ department: d.department, total: +d.total.toFixed(2) }))

    // 2) Average wage by department
    const avgByDept = deptAgg
      .slice()
      .sort((a, b) => b.avg - a.avg)
      .map((d) => ({ department: d.department, avg: +d.avg.toFixed(2) }))

    // 3) Headcount by department
    const headcountByDept = deptAgg
      .slice()
      .sort((a, b) => b.headcount - a.headcount)
      .map((d) => ({ department: d.department, headcount: d.headcount }))

    // 4) Share of payroll by department (pie)
    const payrollShare = deptAgg
      .slice()
      .sort((a, b) => b.total - a.total)
      .map((d) => ({ name: d.department, value: +d.total.toFixed(2) }))

    // 5) Wage stats per dept (min/p25/median/p75/max) — "box plot table"
    const deptBoxStats = deptAgg
      .slice()
      .sort((a, b) => b.total - a.total)
      .map((d) => {
        const s = d.wages.slice().sort((a, b) => a - b)
        return {
          department: d.department,
          min: +quantile(s, 0).toFixed(2),
          p25: +quantile(s, 0.25).toFixed(2),
          median: +quantile(s, 0.5).toFixed(2),
          p75: +quantile(s, 0.75).toFixed(2),
          max: +quantile(s, 1).toFixed(2),
        }
      })

    // 6) Wage histogram
    const wages = valid.map((r) => r.wage).sort((a, b) => a - b)
    const minW = wages[0] ?? 0
    const maxW = wages[wages.length - 1] ?? 0
    const bins = 10
    const width = maxW > minW ? (maxW - minW) / bins : 1
    const hist = Array.from({ length: bins }, (_, i) => ({
      bucket: `${(minW + i * width).toFixed(0)}–${(minW + (i + 1) * width).toFixed(0)}`,
      count: 0,
    }))
    for (const w of wages) {
      const idx = Math.min(bins - 1, Math.max(0, Math.floor((w - minW) / width)))
      hist[idx].count += 1
    }

    // 7) Percentile curve
    const percentile = Array.from({ length: 21 }, (_, i) => {
      const p = i * 5
      const v = wages.length ? quantile(wages, p / 100) : 0
      return { percentile: p, wage: +v.toFixed(2) }
    })

    // 8) Top vs rest (stacked bar)
    const sortedDesc = valid.slice().sort((a, b) => b.wage - a.wage)
    const topN = 10
    const topSum = sortedDesc.slice(0, topN).reduce((s, r) => s + r.wage, 0)
    const restSum = sortedDesc.slice(topN).reduce((s, r) => s + r.wage, 0)
    const topVsRest = [{ label: 'Payroll', top: +topSum.toFixed(2), rest: +restSum.toFixed(2) }]

    // 9) Pareto (top earners + cumulative %)
    const paretoN = Math.min(15, sortedDesc.length)
    let cum = 0
    const pareto = sortedDesc.slice(0, paretoN).map((r, idx) => {
      cum += r.wage
      const cumPct = totalPayroll ? (cum / totalPayroll) * 100 : 0
      return { name: r.name || `Emp ${idx + 1}`, wage: +r.wage.toFixed(2), cumPct: +cumPct.toFixed(2) }
    })

    // 10) Top earners list (bar)
    const topEarners = sortedDesc.slice(0, 10).map((r) => ({ name: r.name, wage: +r.wage.toFixed(2) }))

    // 11) Dept top earner
    const deptTop = deptAgg
      .slice()
      .sort((a, b) => b.max - a.max)
      .map((d) => ({ department: d.department, topWage: +d.max.toFixed(2) }))

    // 12) Wage dot plot by dept (scatter)
    const scatter = valid.map((r) => ({ department: r.department, wage: +r.wage.toFixed(2), name: r.name }))

    // 13) Recipient coverage
    const recipientPresent = rows.filter((r) => r.recipient.length > 0).length
    const recipientMissing = rows.length - recipientPresent
    const recipientCoverage = [
      { name: 'Has recipient', value: recipientPresent },
      { name: 'Missing recipient', value: recipientMissing },
    ]

    // 14) Duplicate recipients
    const recMap = new Map<string, number>()
    for (const r of rows) {
      if (!r.recipient) continue
      recMap.set(r.recipient, (recMap.get(r.recipient) ?? 0) + 1)
    }
    const duplicateRecipients = Array.from(recMap.entries())
      .filter(([, c]) => c > 1)
      .slice(0, 10)
      .map(([recipient, count]) => ({ recipient, count }))

    // 15) Wage validity checks
    const wageValidity = [
      { name: 'Valid wage', value: valid.length },
      { name: 'Invalid / missing wage', value: invalid.length },
    ]

    return {
      headcount,
      totalPayroll: +totalPayroll.toFixed(2),
      avgWage: +avgWage.toFixed(2),

      totalByDept,
      avgByDept,
      headcountByDept,
      payrollShare,
      deptBoxStats,

      hist,
      percentile,
      topVsRest,
      pareto,
      topEarners,
      deptTop,
      scatter,
      recipientCoverage,
      duplicateRecipients,
      wageValidity,
    }
  }, [employeesQuery.data])

  if (employeesQuery.isLoading) return <p>Loading dashboard…</p>
  if (employeesQuery.error) return <p className="error">Failed: {employeesQuery.error.message}</p>

  const pieColors = ['#ff7a00', '#ffb266', '#ffd7b0', '#ff9b33', '#ff8a1a', '#ffa858']

  return (
    <div className="dash">
      {/* KPI Row */}
      <section className="card kpi">
        <div className="kpi-item">
          <div className="kpi-label">Headcount</div>
          <div className="kpi-value">{analytics.headcount}</div>
        </div>
        <div className="kpi-item">
          <div className="kpi-label">Total Payroll</div>
          <div className="kpi-value">${analytics.totalPayroll.toLocaleString()}</div>
        </div>
        <div className="kpi-item">
          <div className="kpi-label">Average Wage</div>
          <div className="kpi-value">${analytics.avgWage.toLocaleString()}</div>
        </div>
      </section>

      {/* Department: totals / averages / headcount */}
      <section className="card">
        <h2 className="section-title">Total wage cost by department</h2>
        <div className="chart">
          <ResponsiveContainer>
            <BarChart data={analytics.totalByDept}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="department" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="total" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </section>

      <section className="card">
        <h2 className="section-title">Share of payroll by department</h2>
        <div className="chart">
          <ResponsiveContainer>
            <PieChart>
              <Tooltip />
              <Pie data={analytics.payrollShare} dataKey="value" nameKey="name" outerRadius={90}>
                {analytics.payrollShare.map((_, i) => (
                  <Cell key={i} fill={pieColors[i % pieColors.length]} />
                ))}
              </Pie>
            </PieChart>
          </ResponsiveContainer>
        </div>
      </section>

      <section className="card">
        <h2 className="section-title">Wage percentile curve</h2>
        <div className="chart">
          <ResponsiveContainer>
            <LineChart data={analytics.percentile}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="percentile" />
              <YAxis />
              <Tooltip />
              <Line type="monotone" dataKey="wage" dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </section>

      <section className="card">
        <h2 className="section-title">Pareto of top earners (wage + cumulative %)</h2>
        <div className="chart">
          <ResponsiveContainer>
            <LineChart data={analytics.pareto}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" tick={{ fontSize: 11 }} interval={0} />
              <YAxis yAxisId="left" />
              <YAxis yAxisId="right" orientation="right" domain={[0, 100]} />
              <Tooltip />
              <Legend />
              <Line yAxisId="left" type="monotone" dataKey="wage" dot={false} name="Wage" />
              <Line yAxisId="right" type="monotone" dataKey="cumPct" dot={false} name="Cumulative %" />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </section>

      {/* Rankings */}
      <section className="card">
        <h2 className="section-title">Top earners</h2>
        <div className="chart">
          <ResponsiveContainer>
            <BarChart data={analytics.topEarners} layout="vertical">
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis type="number" />
              <YAxis type="category" dataKey="name" width={140} />
              <Tooltip />
              <Bar dataKey="wage" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </section>

      {/* Scatter/outliers */}
      <section className="card">
        <h2 className="section-title">Wages by department (outlier view)</h2>
        <div className="chart">
          <ResponsiveContainer>
            <ScatterChart>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="department" type="category" />
              <YAxis dataKey="wage" type="number" />
              <Tooltip cursor={{ strokeDasharray: '3 3' }} />
              <Scatter name="Employees" data={analytics.scatter} />
            </ScatterChart>
          </ResponsiveContainer>
        </div>
      </section>

      {/* “Box plot” substitute table (per dept quantiles) */}
      <section className="card">
        <h2 className="section-title">Department wage spread (min / p25 / median / p75 / max)</h2>
        <div style={{ overflowX: 'auto' }}>
          <table className="table">
            <thead>
              <tr>
                <th>Department</th>
                <th>Min</th>
                <th>P25</th>
                <th>Median</th>
                <th>P75</th>
                <th>Max</th>
              </tr>
            </thead>
            <tbody>
              {analytics.deptBoxStats.map((d) => (
                <tr key={d.department}>
                  <td>{d.department}</td>
                  <td>{d.min}</td>
                  <td>{d.p25}</td>
                  <td>{d.median}</td>
                  <td>{d.p75}</td>
                  <td>{d.max}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  )
}
