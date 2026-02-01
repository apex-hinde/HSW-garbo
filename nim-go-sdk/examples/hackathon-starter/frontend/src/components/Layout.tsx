import React from 'react'
import { Link, NavLink } from 'react-router-dom'
import { NimChat } from '@liminalcash/nim-chat'

export default function Layout({ children }: { children: React.ReactNode }) {
  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'
  const apiUrl = import.meta.env.VITE_API_URL || 'https://api.liminal.cash'

  return (
    <div className="app">
      <header className="topnav">
        <div className="topnav-inner">
          {/* Brand (icon + "liminal payroll and analytics") links to Dashboard */}
          <Link to="/" className="brandlink" aria-label="Go to dashboard">
            <img className="brandicon" src="/liminal-mark.svg" alt="" />
            <span className="brandtext">liminal payroll</span>
          </Link>

          {/* Only two nav items */}
          <nav className="topnav-links">
            <NavLink to="/employees" className={({ isActive }) => `topnav-link ${isActive ? 'active' : ''}`}>
              Employees
            </NavLink>
            <NavLink to="/run" className={({ isActive }) => `topnav-link ${isActive ? 'active' : ''}`}>
              Commands
            </NavLink>
            <NavLink to="/cash-flow" className={({ isActive }) => `topnav-link ${isActive ? 'active' : ''}`}>
              Cash Flow
            </NavLink>
          </nav>
        </div>
      </header>

      <main className="page">
        {children}
      </main>

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
