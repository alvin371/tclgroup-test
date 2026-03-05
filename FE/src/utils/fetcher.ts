import type { TApiResponseError } from '@/commons/types/api'

type TApi = {
  baseUrl: string
}

type TQuery = Record<string, unknown> | string

type TRequestOptions = {
  params?: TQuery
} & RequestInit

type TGetOptions = TRequestOptions | Record<string, unknown>

type TDownloadOptions = {
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  payload?: unknown
  params?: TQuery
} & RequestInit

const REQUEST_INIT_KEYS = new Set([
  'body',
  'cache',
  'credentials',
  'duplex',
  'headers',
  'integrity',
  'keepalive',
  'method',
  'mode',
  'priority',
  'redirect',
  'referrer',
  'referrerPolicy',
  'signal',
  'window',
])

const ABSOLUTE_URL_PATTERN = /^[a-zA-Z][a-zA-Z\d+.-]*:/

const isObject = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null

const looksLikeRequestInit = (value: Record<string, unknown>) =>
  Object.keys(value).some((key) => REQUEST_INIT_KEYS.has(key))

const toInputString = (input: RequestInfo | URL) => {
  if (typeof input === 'string') return input
  if (input instanceof URL) return input.toString()
  if (typeof Request !== 'undefined' && input instanceof Request) return input.url
  return input.toString()
}

const joinUrl = (baseUrl: string, input: string) => {
  if (ABSOLUTE_URL_PATTERN.test(input)) return input

  const normalizedBase = baseUrl.replace(/\/+$/, '')
  const normalizedInput = input.replace(/^\/+/, '')

  if (!normalizedBase) return `/${normalizedInput}`
  return `${normalizedBase}/${normalizedInput}`
}

const toSearchParams = (params: TQuery) => {
  if (typeof params === 'string') return params.replace(/^\?/, '')

  const search = new URLSearchParams()

  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null) continue

    if (Array.isArray(value)) {
      for (const item of value) {
        if (item === undefined || item === null) continue
        search.append(key, String(item))
      }
      continue
    }

    search.append(key, String(value))
  }

  return search.toString()
}

const withQuery = (url: string, params?: TQuery) => {
  if (!params) return url

  const query = toSearchParams(params)
  if (!query) return url

  return `${url}${url.includes('?') ? '&' : '?'}${query}`
}

const parseResponseBody = async (response: Response) => {
  if (response.status === 204) return null

  const contentType = response.headers.get('content-type') ?? ''

  if (contentType.includes('application/json')) {
    try {
      return await response.json()
    } catch {
      return null
    }
  }

  try {
    const text = await response.text()
    return text || null
  } catch {
    return null
  }
}

const toApiError = (data: unknown, response: Response): TApiResponseError => {
  if (
    isObject(data) &&
    data.success === false &&
    isObject(data.error) &&
    typeof data.error.code === 'string' &&
    typeof data.error.message === 'string'
  ) {
    return data as TApiResponseError
  }

  const fallbackMessage =
    (isObject(data) && typeof data.error_message === 'string' && data.error_message) ||
    (isObject(data) && typeof data.message === 'string' && data.message) ||
    response.statusText ||
    'Unknown error'

  return {
    success: false,
    error: {
      code: String(response.status || 'UNKNOWN_ERROR'),
      message: fallbackMessage,
    },
  }
}

const normalizeGetOptions = (options?: TGetOptions): TRequestOptions | undefined => {
  if (!options) return undefined
  if (!isObject(options)) return undefined

  if ('params' in options) return options as TRequestOptions
  if (looksLikeRequestInit(options)) return options as TRequestOptions

  return { params: options }
}

const splitOptions = (options?: TRequestOptions) => {
  if (!options) return { params: undefined as TQuery | undefined, requestOptions: {} as RequestInit }

  const { params, ...requestOptions } = options
  return { params, requestOptions }
}

const parseJsonOrThrow = async <Resp>(response: Response, allowFailure = false) => {
  const data = await parseResponseBody(response)

  if (!response.ok && !allowFailure) {
    throw toApiError(data, response)
  }

  return data as Resp
}

export const Api = ({ baseUrl }: TApi) => {
  return {
    get: async <Resp>(
      input: RequestInfo | URL,
      options?: TGetOptions,
      isCatch = true,
    ) => {
      const normalizedOptions = normalizeGetOptions(options)
      const { params, requestOptions } = splitOptions(normalizedOptions)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: 'GET',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            'Content-Type': 'application/json',
            ...(requestOptions.headers || {}),
          },
        },
      )

      return parseJsonOrThrow<Resp>(response, !isCatch)
    },

    post: async <Resp>(
      input: RequestInfo | URL,
      payload?: unknown,
      options?: TRequestOptions,
    ) => {
      const { params, requestOptions } = splitOptions(options)
      const isFormData = payload instanceof FormData

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: 'POST',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
            ...(requestOptions.headers || {}),
          },
          body: isFormData ? payload : payload === undefined ? undefined : JSON.stringify(payload),
        },
      )

      return parseJsonOrThrow<Resp>(response)
    },

    authPost: async <Resp>(
      input: RequestInfo | URL,
      payload?: unknown,
      options?: TRequestOptions,
    ) => {
      const { params, requestOptions } = splitOptions(options)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: 'POST',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            'Content-Type': 'application/json',
            ...(requestOptions.headers || {}),
          },
          body: payload === undefined ? undefined : JSON.stringify(payload),
        },
      )

      const preloginToken = response.headers.get('Prelogin-Token')
      const resJson = await parseJsonOrThrow<Record<string, unknown>>(response)

      if (!isObject(resJson)) return resJson as Resp

      return {
        ...resJson,
        data: {
          ...(isObject(resJson.data) ? resJson.data : {}),
          preloginToken,
        },
      } as Resp
    },

    put: async <Resp>(
      input: RequestInfo | URL,
      payload?: unknown,
      options?: TRequestOptions,
    ) => {
      const { params, requestOptions } = splitOptions(options)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: 'PUT',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            'Content-Type': 'application/json',
            ...(requestOptions.headers || {}),
          },
          body: payload === undefined ? undefined : JSON.stringify(payload),
        },
      )

      return parseJsonOrThrow<Resp>(response)
    },

    patch: async <Resp>(
      input: RequestInfo | URL,
      payload?: unknown,
      options?: TRequestOptions,
    ) => {
      const { params, requestOptions } = splitOptions(options)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: 'PATCH',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            'Content-Type': 'application/json',
            ...(requestOptions.headers || {}),
          },
          body: payload === undefined ? undefined : JSON.stringify(payload),
        },
      )

      return parseJsonOrThrow<Resp>(response)
    },

    delete: async <Resp>(
      input: RequestInfo | URL,
      payload?: unknown,
      options?: TRequestOptions,
    ) => {
      const { params, requestOptions } = splitOptions(options)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: 'DELETE',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            'Content-Type': 'application/json',
            ...(requestOptions.headers || {}),
          },
          body: payload === undefined ? undefined : JSON.stringify(payload),
        },
      )

      return parseJsonOrThrow<Resp>(response)
    },

    sendFile: async <Resp>(
      input: RequestInfo | URL,
      payload: BodyInit,
      options?: {
        params?: TQuery
        method?: 'POST' | 'PUT' | 'PATCH'
      } & RequestInit,
    ) => {
      const { params, requestOptions } = splitOptions(options)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: options?.method ?? 'POST',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            ...(requestOptions.headers || {}),
          },
          body: payload,
        },
      )

      return parseJsonOrThrow<Resp>(response)
    },

    downloadFile: async (
      input: RequestInfo | URL,
      fileName?: string,
      options?: TRequestOptions,
    ) => {
      const { params, requestOptions } = splitOptions(options)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: 'GET',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            ...(requestOptions.headers || {}),
          },
        },
      )

      if (!response.ok) {
        const errorData = await parseResponseBody(response)
        throw toApiError(errorData, response)
      }

      if (!fileName) return response

      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const anchor = document.createElement('a')

      anchor.href = url
      anchor.download = fileName
      document.body.appendChild(anchor)
      anchor.click()
      document.body.removeChild(anchor)
      window.URL.revokeObjectURL(url)
    },

    downloadFileOi: async (
      input: RequestInfo | URL,
      fileName?: string,
      options?: TRequestOptions,
    ) => {
      const { params, requestOptions } = splitOptions(options)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method: 'GET',
          cache: 'no-store',
          ...requestOptions,
          headers: {
            ...(requestOptions.headers || {}),
          },
        },
      )

      if (!response.ok) {
        const errorData = await parseResponseBody(response)
        throw toApiError(errorData, response)
      }

      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const downloadLink = document.createElement('a')

      downloadLink.href = url
      downloadLink.download =
        fileName || `Report-${new Date().toLocaleDateString('id-ID')}.xlsx`

      document.body.appendChild(downloadLink)
      downloadLink.click()
      document.body.removeChild(downloadLink)
      window.URL.revokeObjectURL(url)
    },

    download: async (
      input: RequestInfo | URL,
      fileName?: string,
      options?: TDownloadOptions,
    ) => {
      const { params, payload, method = 'POST', ...requestOptions } =
        options || ({} as TDownloadOptions)

      const response = await fetch(
        withQuery(joinUrl(baseUrl, toInputString(input)), params),
        {
          method,
          cache: 'no-store',
          ...requestOptions,
          headers: {
            'Content-Type': 'application/json',
            ...(requestOptions.headers || {}),
          },
          body: payload === undefined ? undefined : JSON.stringify(payload),
        },
      )

      if (!response.ok) {
        const errorData = await parseResponseBody(response)
        throw toApiError(errorData, response)
      }

      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const downloadLink = document.createElement('a')

      downloadLink.href = url
      downloadLink.download =
        fileName || `Report-${new Date().toLocaleDateString('id-ID')}.xlsx`

      document.body.appendChild(downloadLink)
      downloadLink.click()
      document.body.removeChild(downloadLink)
      window.URL.revokeObjectURL(url)
    },
  }
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

export const api = Api({ baseUrl: API_BASE_URL })
