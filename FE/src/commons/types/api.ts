export type TApiResponseData<T> = {
  success: true
  message: string
  data: T
}

export type TApiResponsePagination<T> = {
  success: true
  message: string
  items: T[]
  pagination: { page: number; perPage: number; total: number }
}

export type TQueryParams = {
  page?: number
  perPage?: number
  per_page?: number
  [key: string]: unknown
}

export type TApiResponseError = {
  success: false
  error: { code: string; message: string }
}
