export type Role = 'admin' | 'asesor' | 'siswa'

export interface User {
  id: number
  name: string
  username: string
  role: Role
  is_active: boolean
  created_by?: number
  created_at?: string
  updated_at?: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface MeResponse {
  user: User
}

export interface ApiError {
  error: string
}
