import { z } from 'zod'
import type { User, Role } from '../../types'

export const createUserSchema = z.object({
  name: z.string().trim().min(1, 'Nama lengkap wajib diisi'),
  username: z
    .string()
    .trim()
    .min(3, 'Username minimal 3 karakter')
    .regex(/^[a-zA-Z0-9._-]+$/, 'Username hanya boleh mengandung huruf, angka, titik, minus, atau underscore'),
  password: z.string().min(6, 'Password minimal 6 karakter'),
  role: z.enum(['admin', 'asesor', 'siswa']),
})

export type CreateUserInput = z.infer<typeof createUserSchema>

export const updateUserSchema = z.object({
  name: z.string().trim().min(1, 'Nama lengkap wajib diisi'),
  password: z
    .string()
    .transform((val) => val.trim())
    .refine((val) => val === '' || val.length >= 6, {
      message: 'Password baru minimal 6 karakter jika ingin diubah',
    })
    .optional(),
  is_active: z.boolean(),
})

export type UpdateUserInput = z.infer<typeof updateUserSchema>

export type { User, Role }
