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

// Class Domain
export interface Class {
  id: number
  name: string
  owner_id: number
  owner?: User
  members?: User[]
  created_at?: string
  updated_at?: string
}

// Material Domain
export type OptionType = 'ground' | 'warrant'

export interface Option {
  id: number
  argument_id: number
  type: OptionType
  text: string
  is_correct: boolean
  created_at?: string
  updated_at?: string
}

export interface Argument {
  id: number
  material_id: number
  claim_text: string
  order_no: number
  options?: Option[]
  created_at?: string
  updated_at?: string
}

export interface Material {
  id: number
  title: string
  content: string
  owner_id: number
  owner?: User
  arguments?: Argument[]
  created_at?: string
  updated_at?: string
}

// Exam Domain
export interface Exam {
  id: number
  title: string
  material_id: number
  material?: Material
  owner_id: number
  owner?: User
  is_active: boolean
  arguments_per_session: number
  classes?: Class[]
  students?: User[]
  created_at?: string
  updated_at?: string
}
