import { useQuery } from '@tanstack/react-query'
import { api } from '../api/client'
import { listPayrollRuns } from '../api/payroll'
import type { PayrollRun } from '../types/payroll'

type BalanceResp = {
  balances: { currency: string; amount: string; usdValue: string }[]
  totalUsd: string
}

export default function Dashboard() {
  const balance = useQuery<BalanceResp, Error>({
    queryKey: ['balance'],
    queryFn: () => api<BalanceResp>('/api/balance'),
  })

  const runs = useQuery<PayrollRun[], Error>({
    queryKey: ['payrollRuns'],
    queryFn: listPayrollRuns,
  })

  return (
    <div className="grid">
      <section className="card">
        <h2>Wallet Balance</h2>
        {balance.isLoading && <p>Loading...</p>}
        {balance.error && <p className="error">Failed: {balance.error.message}</p>}
        {balance.data && (
          <>
            <div className="big">${balance.data.totalUsd}</div>
            <ul>
              {balance.data.balances.map((b: BalanceResp['balances'][number]) => (
                <li key={b.currency}>
                  <strong>{b.currency}</strong>: {b.amount} (≈ ${b.usdValue})
                </li>
              ))}
            </ul>
          </>
        )}
      </section>

      <section className="card">
        <h2>Recent Payroll Runs</h2>
        {runs.isLoading && <p>Loading...</p>}
        {runs.error && <p className="error">Failed: {runs.error.message}</p>}
        {runs.data && runs.data.length === 0 && <p>No runs yet.</p>}
        {runs.data && runs.data.length > 0 && (
          <ul>
            {runs.data.slice(0, 5).map((r: PayrollRun) => (
              <li key={r.id}>
                <strong>{r.status}</strong> — {new Date(r.createdAt).toLocaleString()} — {r.currency}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}
