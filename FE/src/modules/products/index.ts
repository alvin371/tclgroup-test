import { api } from '@/utils/fetcher'
import { ENDPOINTS } from '@/commons/endpoint'
import type { TApiResponsePagination, TQueryParams } from '@/commons/types/api'

export type TProduct = {
  id: string
  name: string
  sku: string
  customer_id: string
  created_at: string
  updated_at: string
}

export const productsApi = {
  list: (params?: TQueryParams) =>
    api.get<TApiResponsePagination<TProduct>>(ENDPOINTS.PRODUCTS, params as Record<string, unknown>),
}
