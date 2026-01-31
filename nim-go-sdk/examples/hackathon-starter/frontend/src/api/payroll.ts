import { api } from './client'
import type {
  PayrollRun,
  CreatePayrollRunRequest,
  ExecutePayrollRunResponse,
} from '../types/payroll'

export function listPayrollRuns() {
  return api<PayrollRun[]>('/api/payroll/runs')
}

export function createPayrollRun(input: CreatePayrollRunRequest) {
  return api<PayrollRun>('/api/payroll/runs', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export function executePayrollRun(id: string) {
  return api<ExecutePayrollRunResponse>(`/api/payroll/runs/${id}/execute`, {
    method: 'POST',
  })
}
