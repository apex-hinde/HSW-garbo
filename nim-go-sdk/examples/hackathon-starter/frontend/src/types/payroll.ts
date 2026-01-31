export type Employee = {
  id: string
  name: string
  liminalUser: string // e.g. "@alice" or a userId
  salary: string      // keep as string to avoid float issues
  currency: string    // "USD" | "EUR" | "LIL" ...
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
  employeeIds: string[]
}

export type ExecutePayrollRunResponse = {
  success: boolean
  transactionIds?: string[]
  error?: string
}
