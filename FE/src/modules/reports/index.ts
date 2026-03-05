import { api } from '@/utils/fetcher'
import { ENDPOINTS } from '@/commons/endpoint'
import type { TApiResponsePagination, TQueryParams } from '@/commons/types/api'

export type TReportStockIn = {
  id: string
  product_id: string
  quantity: number
  status: 'received'
  notes: string
  created_at: string
  updated_at: string
}

export type TReportStockOut = {
  id: string
  product_id: string
  quantity: number
  status: 'shipped'
  notes: string
  created_at: string
  updated_at: string
}

export const reportsApi = {
  stockIn: (params?: TQueryParams) =>
    api.get<TApiResponsePagination<TReportStockIn>>(
      ENDPOINTS.REPORTS_STOCK_IN,
      params as Record<string, unknown>,
    ),
  stockOut: (params?: TQueryParams) =>
    api.get<TApiResponsePagination<TReportStockOut>>(
      ENDPOINTS.REPORTS_STOCK_OUT,
      params as Record<string, unknown>,
    ),
}
