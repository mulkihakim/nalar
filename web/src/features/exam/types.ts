import { z } from 'zod'
import type { Exam, Material, Class, User } from '../../types'

export const examSchema = z.object({
  title: z.string().trim().min(1, 'Judul paket ujian wajib diisi'),
  material_id: z.number().int().positive('Pilih materi bacaan untuk ujian'),
  arguments_per_session: z
    .number()
    .int('Jumlah argumen harus berupa angka bulat')
    .min(1, 'Jumlah argumen per sesi minimal 1')
    .max(50, 'Jumlah argumen per sesi maksimal 50'),
  is_active: z.boolean(),
})

export type CreateExamInput = z.infer<typeof examSchema> & {
  class_ids?: number[]
  student_ids?: number[]
}

export type UpdateExamInput = Partial<z.infer<typeof examSchema>>

export const examAccessSchema = z.object({
  class_ids: z.array(z.number().int().positive()),
  student_ids: z.array(z.number().int().positive()),
})

export type ExamAccessInput = z.infer<typeof examAccessSchema>

export type { Exam, Material, Class, User }
