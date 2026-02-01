import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { Employee } from '../types/payroll'
import { listEmployees, createEmployee, deleteEmployee } from '../api/employees'

const empty: Omit<Employee, 'id'> = { firstName: '', lastName: '', recipient: '', wage: 0, department: '' }

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

  const delMut = useMutation<{ success: boolean }, Error, number>({
    mutationFn: deleteEmployee,
    onSuccess: async () => qc.invalidateQueries({ queryKey: ['employees'] }),
  })

  const canSubmit = useMemo(() => {
    return (
      form.firstName.trim().length > 0 &&
      form.lastName.trim().length > 0 &&
      form.recipient.trim().length > 0 &&
      form.wage > 0 &&
      form.department.trim().length > 0
    )
  }, [form])

  return (
    <div className="stack">
      <section className="card">
        <h2>Add Employee</h2>

        <div className="form-row">
          <input
            placeholder="First Name"
            value={form.firstName}
            onChange={(e) => setForm((f) => ({ ...f, firstName: e.target.value }))}
          />
          <input
            placeholder="Last Name"
            value={form.lastName}
            onChange={(e) => setForm((f) => ({ ...f, lastName: e.target.value }))}
          />
          <input
            placeholder="Recipient (e.g. @alice)"
            value={form.recipient}
            onChange={(e) => setForm((f) => ({ ...f, recipient: e.target.value }))}
          />
          <input
            placeholder="Wage"
            type="number"
            value={form.wage || ''}
            onChange={(e) => setForm((f) => ({ ...f, wage: parseFloat(e.target.value) || 0 }))}
          />
          <input
            placeholder="Department"
            value={form.department}
            onChange={(e) => setForm((f) => ({ ...f, department: e.target.value }))}
          />

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
                <th>Recipient</th>
                <th>Wage</th>
                <th>Department</th>
                <th />
              </tr>
            </thead>
            <tbody>
              {employees.data.map((emp: Employee) => (
                <tr key={emp.id}>
                  <td>{emp.firstName} {emp.lastName}</td>
                  <td>{emp.recipient}</td>
                  <td>{emp.wage}</td>
                  <td>{emp.department}</td>
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


      </section>
    </div>
  )
}
