import { api } from './client'
import type { Employee } from '../types/payroll'

export function listEmployees() {
  return api<Employee[]>('/api/employees')
}

export function createEmployee(input: Omit<Employee, 'id'>) {
  return api<Employee>('/api/employees', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export function updateEmployee(id: string, input: Omit<Employee, 'id'>) {
  return api<Employee>(`/api/employees/${id}`, {
    method: 'PUT',
    body: JSON.stringify(input),
  })
}

export function deleteEmployee(id: string) {
  return api<{ success: boolean }>(`/api/employees/${id}`, {
    method: 'DELETE',
  })
}
