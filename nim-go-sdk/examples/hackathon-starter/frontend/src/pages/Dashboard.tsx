import { useMemo, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../api/client'

type BalanceResp = {
  balances: { currency: string; amount: string; usdValue: string }[]
  totalUsd: string
}

export default function Dashboard() {
  const [currency, setCurrency] = useState<'USD' | 'EUR' | 'LIL'>('USD')

  const balance = useQuery<BalanceResp, Error>({
    queryKey: ['balance'],
    queryFn: () => api<BalanceResp>('/api/balance'),
  })

  const selected = useMemo(() => {
    const row = balance.data?.balances.find((b) => b.currency === currency)
    return row?.amount ?? '0.00'
  }, [balance.data, currency])

  return (
    <div className="stack">
      <section className="card center">
        <div className="muted" style={{ fontWeight: 700, marginBottom: 8 }}>Your balance</div>

        <div className="h1">
          {currency === 'USD' ? `$${selected}` : `${selected} ${currency}`}
        </div>

        <div style={{ marginTop: 14 }}>
          <div className="pills" role="tablist" aria-label="Currency">
            {(['USD', 'EUR', 'LIL'] as const).map((c) => (
              <button
                key={c}
                className={`pill ${currency === c ? 'active' : ''}`}
                onClick={() => setCurrency(c)}
              >
                {c}
              </button>
            ))}
          </div>
        </div>

        <div className="actions" style={{ marginTop: 18 }}>
          <button className="btn primary" onClick={() => alert('Deposit flow (hook to backend later)')}>
            DEPOSIT
          </button>
          <button className="btn" onClick={() => alert('Withdraw flow (hook to backend later)')}>
            WITHDRAW
          </button>
        </div>
      </section>

      <section className="card">
        <h2 className="section-title">Transactions</h2>
        <p className="muted">Your transaction history will appear here</p>
      </section>
    </div>
  )
}
