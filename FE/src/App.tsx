import { BrowserRouter, Routes, Route, Navigate } from 'react-router'
import { Providers } from '@/providers'
import AppLayout from '@/app/layout'
import InventoryPage from '@/app/inventory/page'
import StockInPage from '@/app/stock-in/page'
import StockOutPage from '@/app/stock-out/page'
import ReportsPage from '@/app/reports/page'
import TransactionDetailPage from '@/app/transaction-detail/page'
import { ROUTES } from '@/commons/route'

export default function App() {
  return (
    <Providers>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<AppLayout />}>
            <Route index element={<Navigate to={ROUTES.INVENTORY} replace />} />
            <Route path={ROUTES.INVENTORY} element={<InventoryPage />} />
            <Route path={ROUTES.STOCK_IN} element={<StockInPage />} />
            <Route path="/stock-in/:id" element={<TransactionDetailPage type="stock-in" />} />
            <Route path={ROUTES.STOCK_OUT} element={<StockOutPage />} />
            <Route path="/stock-out/:id" element={<TransactionDetailPage type="stock-out" />} />
            <Route path={ROUTES.REPORTS} element={<ReportsPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </Providers>
  )
}
