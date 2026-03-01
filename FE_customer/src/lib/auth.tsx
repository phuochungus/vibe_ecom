import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'
import { api } from './api'
import type { User, LoginResponse } from '@/types'

interface AuthContextValue {
  user: User | null
  isAuthenticated: boolean
  login: (identifier: string, password: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

function getUserFromStorage(): User | null {
  try {
    const raw = localStorage.getItem('auth_user')
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(getUserFromStorage)

  const login = useCallback(async (identifier: string, password: string) => {
    // Backend always returns envelope: { success, data: LoginResponse, meta, ... }
    const res = await api.post<{ success: boolean; data: LoginResponse }>('/auth/login', { identifier, password })
    const payload = res.data.data

    // Customer frontend: check USER role
    if (payload.user.role !== 'USER') {
      throw new Error('Tài khoản không có quyền truy cập')
    }

    localStorage.setItem('access_token', payload.access_token)
    localStorage.setItem('refresh_token', payload.refresh_token)
    localStorage.setItem('auth_user', JSON.stringify(payload.user))
    setUser(payload.user)
  }, [])

  const logout = useCallback(async () => {
    try {
      await api.post('/auth/logout')
    } catch {
      // ignore
    } finally {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('auth_user')
      localStorage.removeItem('shopping_cart')
      setUser(null)
    }
  }, [])

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
