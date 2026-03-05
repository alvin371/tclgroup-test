export const ROUTES = {
  INVENTORY: '/inventory',
  STOCK_IN: '/stock-in',
  STOCK_OUT: '/stock-out',
  REPORTS: '/reports',
} as const

export type TRoutePath = (typeof ROUTES)[keyof typeof ROUTES]

export const toStockInDetail = (id: string) => `/stock-in/${id}`
export const toStockOutDetail = (id: string) => `/stock-out/${id}`
