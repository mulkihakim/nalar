const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

export const TOKEN_STORAGE_KEY = 'nalar_auth_token'

export function getStoredToken(): string | null {
  return localStorage.getItem(TOKEN_STORAGE_KEY)
}

export function setStoredToken(token: string): void {
  localStorage.setItem(TOKEN_STORAGE_KEY, token)
}

export function removeStoredToken(): void {
  localStorage.removeItem(TOKEN_STORAGE_KEY)
}

interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>
}

export async function apiClient<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { params, headers, ...restOptions } = options

  let url = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint : `/${endpoint}`}`

  if (params) {
    const searchParams = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        searchParams.append(key, String(value))
      }
    })
    const queryString = searchParams.toString()
    if (queryString) {
      url += (url.includes('?') ? '&' : '?') + queryString
    }
  }

  const token = getStoredToken()
  const defaultHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (token) {
    defaultHeaders['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(url, {
    ...restOptions,
    headers: {
      ...defaultHeaders,
      ...headers,
    },
  })

  if (!response.ok) {
    let errorMessage = `HTTP error ${response.status}`
    try {
      const errorData = await response.json()
      if (errorData?.error) {
        errorMessage = errorData.error
      }
    } catch {
      // response is not json
    }

    if (response.status === 401) {
      removeStoredToken()
    }

    throw new Error(errorMessage)
  }

  return response.json() as Promise<T>
}
