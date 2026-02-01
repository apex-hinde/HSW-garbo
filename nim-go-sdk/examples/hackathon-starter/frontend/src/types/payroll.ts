export type Employee = {
  id: number
  firstName: string
  lastName: string
  recipient: string
  wage: number
  department: string
}

export type PayrollRun = {
  id: string
  periodStart: string // ISO
  periodEnd: string   // ISO
  currency: string
  createdAt: string   // ISO
  status: 'draft' | 'processing' | 'paid' | 'failed'
}

export type CreatePayrollRunRequest = {
  periodStart: string
  periodEnd: string
  employeeIds: number[]
}

export type ExecutePayrollRunResponse = {
  success: boolean
  transactionIds?: string[]
  error?: string
}
