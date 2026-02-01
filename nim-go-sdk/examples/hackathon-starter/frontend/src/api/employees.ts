import { api } from './client'
import type { Employee } from '../types/payroll'

export async function listEmployees() {
  const res = await api<{ employees: Employee[] }>('/api/employees')
  return res.employees || []
}

export function createEmployee(input: Omit<Employee, 'id'>) {
  return api<Employee>('/api/employees', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export function updateEmployee(id: number, input: Omit<Employee, 'id'>) {
  return api<Employee>(`/api/employees/${id}`, {
    method: 'PUT',
    body: JSON.stringify(input),
  })
}

export function deleteEmployee(id: number) {
  return api<{ success: boolean }>(`/api/employees/${id}`, {
    method: 'DELETE',
  })
}
