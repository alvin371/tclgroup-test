import { api } from '@/utils/fetcher'
import { ENDPOINTS } from '@/commons/endpoint'
import type { TApiResponsePagination, TApiResponseData, TQueryParams } from '@/commons/types/api'

export type TInventoryItem = {
  id: string
  product_id: string
  product_name: string
  product_sku: string
  physical_stock: number
  reserved: number
  available_stock: number
  updated_at: string
}

export type TAdjustInventoryPayload = {
  new_qty: number
  notes: string
}

export const inventoryApi = {
  list: (params?: TQueryParams) =>
    api.get<TApiResponsePagination<TInventoryItem>>(ENDPOINTS.INVENTORY, params as Record<string, unknown>),
  adjust: (productId: string, payload: TAdjustInventoryPayload) =>
    api.patch<TApiResponseData<TInventoryItem>>(ENDPOINTS.INVENTORY_ADJUST(productId), payload),
}
