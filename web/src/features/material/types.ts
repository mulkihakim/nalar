import { z } from 'zod'
import type { Material, Argument, Option, OptionType } from '../../types'

export const materialSchema = z.object({
  title: z.string().trim().min(1, 'Judul materi bacaan wajib diisi'),
  content: z.string().trim().min(1, 'Isi materi teks bacaan wajib diisi'),
})

export type CreateMaterialInput = z.infer<typeof materialSchema>
export type UpdateMaterialInput = z.infer<typeof materialSchema>

export const optionInputSchema = z.object({
  type: z.enum(['ground', 'warrant']),
  text: z.string().trim().min(1, 'Teks opsi tidak boleh kosong'),
  is_correct: z.boolean(),
})

export type OptionInput = z.infer<typeof optionInputSchema>

export const argumentSchema = z
  .object({
    claim_text: z.string().trim().min(1, 'Pernyataan klaim wajib diisi'),
    order_no: z.number().int().optional(),
    options: z.array(optionInputSchema),
  })
  .superRefine((data, ctx) => {
    const grounds = data.options.filter((o) => o.type === 'ground')
    const warrants = data.options.filter((o) => o.type === 'warrant')

    if (grounds.length < 3) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Minimal 3 pilihan opsi untuk Ground',
        path: ['options_ground'],
      })
    }
    const groundCorrect = grounds.filter((o) => o.is_correct).length
    if (groundCorrect !== 1) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Pilih tepat 1 jawaban BENAR untuk Ground',
        path: ['options_ground_correct'],
      })
    }

    if (warrants.length < 3) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Minimal 3 pilihan opsi untuk Warrant',
        path: ['options_warrant'],
      })
    }
    const warrantCorrect = warrants.filter((o) => o.is_correct).length
    if (warrantCorrect !== 1) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Pilih tepat 1 jawaban BENAR untuk Warrant',
        path: ['options_warrant_correct'],
      })
    }
  })

export type CreateArgumentInput = z.infer<typeof argumentSchema>

export type { Material, Argument, Option, OptionType }
