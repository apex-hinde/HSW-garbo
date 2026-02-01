import { useMemo, useState } from 'react'

type Category =
  | 'Banking & Wallet'
  | 'Savings & Investment'
  | 'Employee Directory'
  | 'Payroll Management'
  | 'Key Features'

type Command = {
  id: string
  category: Category
  name: string
  description: string
  requiresConfirmation?: boolean
  notes?: string
}

const COMMANDS: Command[] = [
  // 💰 Banking & Wallet Tools
  {
    id: 'get_balance',
    category: 'Banking & Wallet',
    name: 'get_balance',
    description: 'Check your wallet balance (can filter by currency).',
  },
  {
    id: 'get_transactions',
    category: 'Banking & Wallet',
    name: 'get_transactions',
    description: 'View recent transaction history.',
  },
  {
    id: 'get_profile',
    category: 'Banking & Wallet',
    name: 'get_profile',
    description: 'Get your profile information.',
  },
  {
    id: 'search_users',
    category: 'Banking & Wallet',
    name: 'search_users',
    description: 'Find users by display tag or name.',
  },
  {
    id: 'send_money',
    category: 'Banking & Wallet',
    name: 'send_money',
    description: 'Send money to another user.',
    requiresConfirmation: true,
    notes: 'Requires confirmation before execution.',
  },

  // 🏦 Savings & Investment Tools
  {
    id: 'get_savings_balance',
    category: 'Savings & Investment',
    name: 'get_savings_balance',
    description: 'Check savings positions and current APY.',
  },
  {
    id: 'get_vault_rates',
    category: 'Savings & Investment',
    name: 'get_vault_rates',
    description: 'View current APY rates for savings vaults.',
  },
  {
    id: 'deposit_savings',
    category: 'Savings & Investment',
    name: 'deposit_savings',
    description: 'Move money into savings.',
    requiresConfirmation: true,
    notes: 'Requires confirmation before execution.',
  },
  {
    id: 'withdraw_savings',
    category: 'Savings & Investment',
    name: 'withdraw_savings',
    description: 'Take money out of savings.',
    requiresConfirmation: true,
    notes: 'Requires confirmation before execution.',
  },

  // 👥 Employee Directory Tools
  {
    id: 'create_employee',
    category: 'Employee Directory',
    name: 'create_employee',
    description: 'Add new employee to directory.',
  },
  {
    id: 'get_employee',
    category: 'Employee Directory',
    name: 'get_employee',
    description: 'Get employee details by ID.',
  },
  {
    id: 'list_employees',
    category: 'Employee Directory',
    name: 'list_employees',
    description: 'Show all employees.',
  },
  {
    id: 'update_employee',
    category: 'Employee Directory',
    name: 'update_employee',
    description: 'Modify employee information.',
  },
  {
    id: 'delete_employee',
    category: 'Employee Directory',
    name: 'delete_employee',
    description: 'Remove employee from directory.',
  },
  {
    id: 'list_employees_by_department',
    category: 'Employee Directory',
    name: 'list_employees_by_department',
    description: 'Filter employees by department.',
  },
  {
    id: 'count_employees',
    category: 'Employee Directory',
    name: 'count_employees',
    description: 'Get total employee count.',
  },

  // 💼 Payroll Management Tools
  {
    id: 'payroll_check',
    category: 'Payroll Management',
    name: 'payroll_check',
    description: 'Check if payroll is completed.',
  },
  {
    id: 'fulfill_remaining_payroll',
    category: 'Payroll Management',
    name: 'fulfill_remaining_payroll',
    description: 'Process payroll for all employees.',
  },
]

// Key Features shown as “cards” too
const FEATURES: { title: string; body: string }[] = [
  { title: 'Confirmation required', body: 'All money movements require your confirmation before execution.' },
  { title: 'Multi-currency', body: 'Works with USD, EUR, and LIL.' },
  { title: 'Employee management', body: 'Includes wages, departments, and recipient handles.' },
  { title: 'Payroll bulk payments', body: 'Payroll tools handle bulk payments to your entire team.' },
]

function copyToClipboard(text: string) {
  if (navigator.clipboard?.writeText) return navigator.clipboard.writeText(text)
  const ta = document.createElement('textarea')
  ta.value = text
  document.body.appendChild(ta)
  ta.select()
  document.execCommand('copy')
  document.body.removeChild(ta)
  return Promise.resolve()
}

export default function RunPayroll() {
  const [query, setQuery] = useState('')
  const [category, setCategory] = useState<'All' | Category>('All')
  const [copied, setCopied] = useState<string | null>(null)

  const categories = useMemo(() => {
    const set = new Set<Category>()
    for (const c of COMMANDS) set.add(c.category)
    return ['All', ...Array.from(set)] as const
  }, [])

  const grouped = useMemo(() => {
    const q = query.trim().toLowerCase()
    const filtered = COMMANDS.filter((c) => {
      const matchesCategory = category === 'All' ? true : c.category === category
      const matchesQuery =
        !q ||
        c.name.toLowerCase().includes(q) ||
        c.description.toLowerCase().includes(q) ||
        (c.notes ?? '').toLowerCase().includes(q)
      return matchesCategory && matchesQuery
    })

    const map = new Map<Category, Command[]>()
    for (const cmd of filtered) {
      if (!map.has(cmd.category)) map.set(cmd.category, [])
      map.get(cmd.category)!.push(cmd)
    }

    // keep category order consistent
    const order: Category[] = [
      'Banking & Wallet',
      'Savings & Investment',
      'Employee Directory',
      'Payroll Management',
    ]

    return order
      .filter((cat) => map.has(cat))
      .map((cat) => ({ category: cat, commands: map.get(cat)! }))
  }, [query, category])

  const totalShown = grouped.reduce((s, g) => s + g.commands.length, 0)

  return (
    <div className="stack">
      <section className="card">
        <h2 className="section-title">Commands</h2>
        <p className="muted" style={{ marginTop: 6 }}>
          Quick reference for available tools. Use search + filters to find what you need.
        </p>

        <div style={{ display: 'grid', gap: 12, marginTop: 14 }}>
          <input
            placeholder="Search (e.g. payroll, savings, send_money)"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />

          <div className="pills" aria-label="Command categories">
            {categories.map((c) => (
              <button
                key={c}
                className={`pill ${category === c ? 'active' : ''}`}
                onClick={() => setCategory(c)}
                type="button"
              >
                {c}
              </button>
            ))}
          </div>

          <div className="muted" style={{ fontSize: 13 }}>
            Showing <strong>{totalShown}</strong> command{totalShown === 1 ? '' : 's'}
          </div>
        </div>
      </section>

      <section className="card">
        <h2 className="section-title">Key Features</h2>
        <div style={{ display: 'grid', gap: 12, gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))' }}>
          {FEATURES.map((f) => (
            <div
              key={f.title}
              style={{
                padding: 14,
                borderRadius: 16,
                border: '1px solid rgba(18,18,18,0.10)',
                background: 'rgba(18,18,18,0.02)',
              }}
            >
              <div style={{ fontWeight: 900 }}>{f.title}</div>
              <div className="muted" style={{ marginTop: 6 }}>
                {f.body}
              </div>
            </div>
          ))}
        </div>
      </section>

      {grouped.map((group) => (
        <section key={group.category} className="card">
          <h2 className="section-title">{group.category}</h2>

          <div style={{ display: 'grid', gap: 12 }}>
            {group.commands.map((cmd) => (
              <div
                key={cmd.id}
                style={{
                  padding: 14,
                  borderRadius: 16,
                  border: '1px solid rgba(18,18,18,0.10)',
                  background: 'rgba(18,18,18,0.02)',
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', gap: 12 }}>
                  <div>
                    <div style={{ display: 'flex', gap: 10, alignItems: 'center', flexWrap: 'wrap' }}>
                      <code style={{ fontWeight: 900, fontSize: 14 }}>{cmd.name}</code>
                      {cmd.requiresConfirmation && (
                        <span
                          style={{
                            fontSize: 12,
                            fontWeight: 800,
                            padding: '4px 10px',
                            borderRadius: 999,
                            border: '1px solid rgba(255,122,0,0.35)',
                            background: 'rgba(255,122,0,0.10)',
                          }}
                        >
                          Confirmation required
                        </span>
                      )}
                    </div>
                    <div className="muted" style={{ marginTop: 6 }}>
                      {cmd.description}
                    </div>
                    {cmd.notes && (
                      <div className="muted" style={{ marginTop: 6, fontSize: 13 }}>
                        {cmd.notes}
                      </div>
                    )}
                  </div>

                  <button
                    className="btn"
                    type="button"
                    onClick={async () => {
                      await copyToClipboard(cmd.name)
                      setCopied(cmd.id)
                      setTimeout(() => setCopied(null), 900)
                    }}
                    aria-label={`Copy ${cmd.name}`}
                  >
                    {copied === cmd.id ? 'Copied' : 'Copy'}
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>
      ))}
    </div>
  )
}
