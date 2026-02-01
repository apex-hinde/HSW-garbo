import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Employees from './pages/Employees'
import RunPayroll from './pages/RunPayroll'
import CashFlow from './pages/CashFlow'

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/employees" element={<Employees />} />
        <Route path="/run" element={<RunPayroll />} />
        <Route path="/cash-flow" element={<CashFlow />} />
      </Routes>
    </Layout>
  )
}
