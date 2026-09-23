import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import { getStoredToken, removeStoredToken, setStoredToken } from '../../lib/api-client'
import type { User } from '../../types'
import { fetchMe } from './api'
import type { AuthContextType } from './types'

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(getStoredToken())
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState<boolean>(true)

  useEffect(() => {
    async function initAuth() {
      const storedToken = getStoredToken()
      if (!storedToken) {
        setIsLoading(false)
        return
      }

      try {
        const res = await fetchMe()
        setUser(res.user)
        setToken(storedToken)
      } catch {
        removeStoredToken()
        setToken(null)
        setUser(null)
      } finally {
        setIsLoading(false)
      }
    }

    initAuth()
  }, [])

  const login = (newToken: string, newUser: User) => {
    setStoredToken(newToken)
    setToken(newToken)
    setUser(newUser)
  }

  const logout = () => {
    removeStoredToken()
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated: !!token && !!user,
        isLoading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
