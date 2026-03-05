import { api } from '@/utils/fetcher'
import { ENDPOINTS } from '@/commons/endpoint'
import type { TApiResponsePagination, TApiResponseData, TQueryParams } from '@/commons/types/api'

export type TStockInStatus = 'CREATED' | 'IN_PROGRESS' | 'DONE' | 'CANCELLED'

export type TStockIn = {
  id: string
  product_id: string
  quantity: number
  status: TStockInStatus
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

export type TStockInDetail = {
  id: string
  product_id: string
  product_sku: string
  product_name: string
  quantity: number
  unit_price: number
  status: TStockInStatus
  performed_by: string
  location: string
  notes: string
  created_at: string
  updated_at: string
  timeline: TTransactionTimelineEntry[]
}

export type TCreateStockInPayload = {
  product_id: string
  quantity: number
  notes: string
}

export type TAdvanceStockInPayload = {
  status: 'IN_PROGRESS' | 'DONE'
}

export const stockInApi = {
  list: (params?: TQueryParams) =>
    api.get<TApiResponsePagination<TStockIn>>(ENDPOINTS.STOCK_IN, params as Record<string, unknown>),
  create: (payload: TCreateStockInPayload) =>
    api.post<TApiResponseData<TStockIn>>(ENDPOINTS.STOCK_IN, payload),
  getDetail: (id: string) =>
    api.get<TApiResponseData<TStockInDetail>>(ENDPOINTS.STOCK_IN_DETAIL(id)),
  advance: (id: string, payload: TAdvanceStockInPayload) =>
    api.patch<TApiResponseData<TStockIn>>(ENDPOINTS.STOCK_IN_ADVANCE(id), payload),
  remove: (id: string) =>
    api.delete<TApiResponseData<null>>(ENDPOINTS.STOCK_IN_DELETE(id)),
}
