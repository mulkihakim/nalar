import { z } from 'zod'
import type { Class, User } from '../../types'

export const classSchema = z.object({
  name: z.string().trim().min(1, 'Nama kelas wajib diisi'),
})

export type CreateClassInput = z.infer<typeof classSchema>
export type UpdateClassInput = z.infer<typeof classSchema>

export const addMemberSchema = z.object({
  user_id: z.number().int().positive('Pilih siswa yang valid'),
})

export type AddMemberInput = z.infer<typeof addMemberSchema>

export type { Class, User }
