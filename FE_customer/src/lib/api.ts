import axios, { type AxiosError } from 'axios'

const apiHost = import.meta.env.VITE_API_HOST?.trim().replace(/^https?:\/\//, '').replace(/\/$/, '')
const apiBaseURL = import.meta.env.VITE_API_BASE_URL?.trim() || (apiHost ? `https://${apiHost}/api/v1` : '/api/v1')

const buildApiUrl = (path: string) => `${apiBaseURL.replace(/\/$/, '')}${path}`

export const api = axios.create({
  baseURL: apiBaseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor — attach access token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor — handle 401 (token expired)
api.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const originalRequest = error.config as typeof error.config & { _retry?: boolean }

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        try {
          // Refresh token also returns envelope: { success, data: { access_token, ... } }
          const { data } = await axios.post(buildApiUrl('/auth/refresh'), {
            refresh_token: refreshToken,
          })
          const newToken = data.data?.access_token ?? data.access_token
          localStorage.setItem('access_token', newToken)
          originalRequest!.headers!.Authorization = `Bearer ${newToken}`
          return api(originalRequest!)
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          localStorage.removeItem('auth_user')
          window.location.href = '/login'
        }
      } else {
        window.location.href = '/login'
      }
    }

    return Promise.reject(error)
  }
)

// Helpers
export const getErrorMessage = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    return error.response?.data?.error?.message ?? error.message
  }
  return 'Đã xảy ra lỗi, vui lòng thử lại.'
}
