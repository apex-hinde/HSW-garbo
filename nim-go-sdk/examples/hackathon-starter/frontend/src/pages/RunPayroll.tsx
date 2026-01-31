import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { listEmployees } from '../api/employees'
import { createPayrollRun, executePayrollRun } from '../api/payroll'
import type { Employee, PayrollRun, CreatePayrollRunRequest, ExecutePayrollRunResponse } from '../types/payroll'

export default function RunPayroll() {
  const qc = useQueryClient()

  const employees = useQuery<Employee[], Error>({
    queryKey: ['employees'],
    queryFn: listEmployees,
  })

  const [selected, setSelected] = useState<Record<string, boolean>>({})
  const [periodStart, setPeriodStart] = useState('')
  const [periodEnd, setPeriodEnd] = useState('')

  const selectedIds = useMemo(
    () => Object.entries(selected).filter(([, v]) => v).map(([k]) => k),
    [selected]
  )

  const createRun = useMutation<PayrollRun, Error, CreatePayrollRunRequest>({
    mutationFn: createPayrollRun,
  })

  const execRun = useMutation<ExecutePayrollRunResponse, Error, string>({
    mutationFn: executePayrollRun,
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['payrollRuns'] })
    },
  })

  const canCreate = periodStart && periodEnd && selectedIds.length > 0

  return (
    <div className="stack">
      <section className="card">
        <h2>1) Select pay period</h2>
        <div className="form-row">
          <input type="date" value={periodStart} onChange={(e) => setPeriodStart(e.target.value)} />
          <input type="date" value={periodEnd} onChange={(e) => setPeriodEnd(e.target.value)} />
        </div>
      </section>

      <section className="card">
        <h2>2) Select employees</h2>
        {employees.isLoading && <p>Loading...</p>}
        {employees.error && <p className="error">{employees.error.message}</p>}

        {employees.data && employees.data.length > 0 && (
          <div className="stack">
            {employees.data.map((emp: Employee) => (
              <label key={emp.id} className="check">
                <input
                  type="checkbox"
                  checked={!!selected[emp.id]}
                  onChange={(ev) => setSelected((s) => ({ ...s, [emp.id]: ev.target.checked }))}
                />
                <span>
                  <strong>{emp.name}</strong> — {emp.salary} {emp.currency} — {emp.liminalUser}
                </span>
              </label>
            ))}
          </div>
        )}
      </section>

      <section className="card">
        <h2>3) Create & execute</h2>

        <button
          disabled={!canCreate || createRun.isPending}
          onClick={() =>
            createRun.mutate({
              periodStart: new Date(periodStart).toISOString(),
              periodEnd: new Date(periodEnd).toISOString(),
              employeeIds: selectedIds,
            })
          }
        >
          {createRun.isPending ? 'Creating...' : 'Create payroll run'}
        </button>

        {createRun.data && (
          <div className="stack" style={{ marginTop: 12 }}>
            <div>
              Created run: <code>{createRun.data.id}</code> — status: <strong>{createRun.data.status}</strong>
            </div>
            <button disabled={execRun.isPending} onClick={() => execRun.mutate(createRun.data.id)}>
              {execRun.isPending ? 'Paying...' : 'Execute payroll'}
            </button>
          </div>
        )}

        {createRun.error && <p className="error">{createRun.error.message}</p>}
        {execRun.error && <p className="error">{execRun.error.message}</p>}

        {execRun.data && !execRun.data.success && <p className="error">{execRun.data.error}</p>}
        {execRun.data && execRun.data.success && <p>✅ Payroll executed.</p>}
      </section>
    </div>
  )
}
