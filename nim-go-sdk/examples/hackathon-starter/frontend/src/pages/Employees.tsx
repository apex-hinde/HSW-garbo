import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { Employee } from '../types/payroll'
import { listEmployees, createEmployee, deleteEmployee } from '../api/employees'

const empty: Omit<Employee, 'id'> = { name: '', liminalUser: '', salary: '', currency: 'USD' }

export default function Employees() {
  const qc = useQueryClient()
  const [form, setForm] = useState(empty)

  const employees = useQuery<Employee[], Error>({
    queryKey: ['employees'],
    queryFn: listEmployees,
  })

  const createMut = useMutation<Employee, Error, Omit<Employee, 'id'>>({
    mutationFn: createEmployee,
    onSuccess: async () => {
      setForm(empty)
      await qc.invalidateQueries({ queryKey: ['employees'] })
    },
  })

  const delMut = useMutation<{ success: boolean }, Error, string>({
    mutationFn: deleteEmployee,
    onSuccess: async () => qc.invalidateQueries({ queryKey: ['employees'] }),
  })

  const canSubmit = useMemo(() => {
    return (
      form.name.trim().length > 0 &&
      form.liminalUser.trim().length > 0 &&
      form.salary.trim().length > 0 &&
      form.currency.trim().length > 0
    )
  }, [form])

  return (
    <div className="stack">
      <section className="card">
        <h2>Add Employee</h2>

        <div className="form-row">
          <input
            placeholder="Name"
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
          />
          <input
            placeholder="Liminal user (e.g. @alice)"
            value={form.liminalUser}
            onChange={(e) => setForm((f) => ({ ...f, liminalUser: e.target.value }))}
          />
          <input
            placeholder="Salary (e.g. 2500.00)"
            value={form.salary}
            onChange={(e) => setForm((f) => ({ ...f, salary: e.target.value }))}
          />
          <select
            value={form.currency}
            onChange={(e) => setForm((f) => ({ ...f, currency: e.target.value }))}
          >
            <option value="USD">USD</option>
            <option value="EUR">EUR</option>
            <option value="LIL">LIL</option>
          </select>

          <button disabled={!canSubmit || createMut.isPending} onClick={() => createMut.mutate(form)}>
            {createMut.isPending ? 'Adding...' : 'Add'}
          </button>
        </div>

        {createMut.error && <p className="error">{createMut.error.message}</p>}
      </section>

      <section className="card">
        <h2>Employees</h2>

        {employees.isLoading && <p>Loading...</p>}
        {employees.error && <p className="error">{employees.error.message}</p>}
        {employees.data && employees.data.length === 0 && <p>No employees yet.</p>}

        {employees.data && employees.data.length > 0 && (
          <table className="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Liminal</th>
                <th>Salary</th>
                <th>Currency</th>
                <th />
              </tr>
            </thead>
            <tbody>
              {employees.data.map((emp: Employee) => (
                <tr key={emp.id}>
                  <td>{emp.name}</td>
                  <td>{emp.liminalUser}</td>
                  <td>{emp.salary}</td>
                  <td>{emp.currency}</td>
                  <td style={{ textAlign: 'right' }}>
                    <button className="danger" disabled={delMut.isPending} onClick={() => delMut.mutate(emp.id)}>
                      Remove
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}

        <p className="hint">
          Next step: add Edit → call <code>PUT /api/employees/:id</code>
        </p>
      </section>
    </div>
  )
}
