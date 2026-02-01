# HSW-garbo — Liminal Payroll and Analytics Manager

A hackathon project built on the **Liminal / Nim Go SDK hackathon starter**.  
This app helps companies manage an employee directory and run payroll using Liminal’s wallet, savings, and payroll tools (with **confirmation required** for any money movement).

---

## What you get

- **Go backend** exposing APIs + tool execution
- **Web frontend** (Vite + React + TypeScript) for payroll ops + analytics dashboard
- **Employee directory** with wages, departments, and recipient handles
- **Payroll automation helpers** for bulk payments
- **Multi-currency support**: USD, EUR, LIL

---

## Quick start

### 1) Backend (Go)

From the hackathon starter folder:

```bash
cd examples/hackathon-starter
cp .env.example .env
# Add required keys/secrets to .env
go run .
```

Backend typically runs on:

- `http://localhost:8080`

### 2) Frontend (Vite + React + TS)

In another terminal:

```bash
cd examples/hackathon-starter/frontend
npm install
npm run dev
```

Frontend typically runs on:

- `http://localhost:5173`

---

## Commands & tools

### Banking & Wallet Tools

- `get_balance` — Check your wallet balance (can filter by currency)
- `get_transactions` — View recent transaction history
- `get_profile` — Get your profile information
- `search_users` — Find users by display tag or name
- `send_money` — Send money to another user (**requires confirmation**)

### Savings & Investment Tools

- `get_savings_balance` — Check savings positions and current APY
- `get_vault_rates` — View current APY rates for savings vaults
- `deposit_savings` — Move money into savings (**requires confirmation**)
- `withdraw_savings` — Take money out of savings (**requires confirmation**)

### Employee Directory Tools

- `create_employee` — Add new employee to directory
- `get_employee` — Get employee details by ID
- `list_employees` — Show all employees
- `update_employee` — Modify employee information
- `delete_employee` — Remove employee from directory
- `list_employees_by_department` — Filter employees by department
- `count_employees` — Get total employee count

### Payroll Management Tools

- `payroll_check` — Check if payroll is completed
- `fulfill_remaining_payroll` — Process payroll for all employees

---

## Key features

- **All money movements require confirmation** before execution
- **Multi-currency**: USD, EUR, LIL
- Employee management includes wages, departments, and recipient handles
- Payroll tools handle **bulk payments** to the entire team

---

## Project structure

```txt
examples/hackathon-starter/
  frontend/                # Vite + React + TS UI
  internal/
    api/                   # HTTP handlers (REST)
    storage/               # DB + employee storage layer
  main.go                  # server entrypoint
```

---

## Common troubleshooting

### `/api/*` returns the frontend page instead of JSON
This usually means the frontend is not reaching the Go server.

Check:
- the backend is running (e.g. on `:8080`)
- the frontend is calling the correct base URL or using a dev proxy
- try opening backend directly: `http://localhost:8080/api/employees`

---

## Team

Built for a hackathon using the Liminal / Nim Go SDK starter template.
