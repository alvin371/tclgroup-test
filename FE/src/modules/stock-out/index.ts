import { api } from '@/utils/fetcher'
import { ENDPOINTS } from '@/commons/endpoint'
import type { TApiResponsePagination, TApiResponseData, TQueryParams } from '@/commons/types/api'

export type TStockOutStatus = 'DRAFT' | 'IN_PROGRESS' | 'DONE' | 'CANCELLED'

export type TStockOut = {
  id: string
  product_id: string
  quantity: number
  status: TStockOutStatus
  notes: string
  unit_price: number
  performed_by: string
  location: string
  created_at: string
  updated_at: string
}

export type TTransactionTimelineEntry = {
  status: string
  description: string
  occurred_at: string
}

export type TStockOutDetail = {
  id: string
  product_id: string
  product_sku: string
  product_name: string
  quantity: number
  unit_price: number
  status: TStockOutStatus
  performed_by: string
  location: string
  notes: string
  created_at: string
  updated_at: string
  timeline: TTransactionTimelineEntry[]
}

export type TAllocateStockOutPayload = {
  product_id: string
  quantity: number
  notes: string
}

export type TExecuteStockOutPayload = {
  status: 'IN_PROGRESS' | 'DONE'
}

export const stockOutApi = {
  list: (params?: TQueryParams) =>
    api.get<TApiResponsePagination<TStockOut>>(ENDPOINTS.STOCK_OUT, params as Record<string, unknown>),
  allocate: (payload: TAllocateStockOutPayload) =>
    api.post<TApiResponseData<TStockOut>>(ENDPOINTS.STOCK_OUT_ALLOCATE, payload),
  getDetail: (id: string) =>
    api.get<TApiResponseData<TStockOutDetail>>(ENDPOINTS.STOCK_OUT_DETAIL(id)),
  execute: (id: string, payload: TExecuteStockOutPayload) =>
    api.patch<TApiResponseData<TStockOut>>(ENDPOINTS.STOCK_OUT_EXECUTE(id), payload),
  remove: (id: string) =>
    api.delete<TApiResponseData<null>>(ENDPOINTS.STOCK_OUT_DELETE(id)),
}
