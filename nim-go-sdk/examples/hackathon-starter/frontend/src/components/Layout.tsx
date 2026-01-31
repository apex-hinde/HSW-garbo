import React from 'react'
import { NavLink } from 'react-router-dom'
import { NimChat } from '@liminalcash/nim-chat'

export default function Layout({ children }: { children: React.ReactNode }) {
  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'
  const apiUrl = import.meta.env.VITE_API_URL || 'https://api.liminal.cash'

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">Liminal Payroll</div>

        <nav className="nav">
          <NavLink to="/" end className={({ isActive }) => (isActive ? 'active' : '')}>
            Dashboard
          </NavLink>
          <NavLink to="/employees" className={({ isActive }) => (isActive ? 'active' : '')}>
            Employees
          </NavLink>
          <NavLink to="/run" className={({ isActive }) => (isActive ? 'active' : '')}>
            Run Payroll
          </NavLink>
        </nav>
      </aside>

      <div className="content">
        <header className="topbar">
          <div className="topbar-title">Company Finance Manager</div>
          <div className="topbar-subtitle">Manage employees, salaries, and payouts via Liminal</div>
        </header>

        <div className="page">{children}</div>
      </div>

      <NimChat
        wsUrl={wsUrl}
        apiUrl={apiUrl}
        title="Nim"
        position="bottom-right"
        defaultOpen={false}
      />
    </div>
  )
}
