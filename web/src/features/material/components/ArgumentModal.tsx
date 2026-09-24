import { useState, useEffect } from 'react'
import { CheckCircle2, AlertCircle, Plus, Trash2 } from 'lucide-react'
import { Modal } from '../../../components/Modal'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { Label } from '../../../components/ui/label'
import type { Argument, OptionInput } from '../types'
import { argumentSchema } from '../types'
import { useCreateArgument, useUpdateArgument } from '../api'

interface ArgumentModalProps {
  isOpen: boolean
  onClose: () => void
  materialId: number
  argumentToEdit?: Argument | null
}

const createInitialGrounds = (): OptionInput[] => [
  { type: 'ground', text: '', is_correct: true },
  { type: 'ground', text: '', is_correct: false },
  { type: 'ground', text: '', is_correct: false },
]

const createInitialWarrants = (): OptionInput[] => [
  { type: 'warrant', text: '', is_correct: true },
  { type: 'warrant', text: '', is_correct: false },
  { type: 'warrant', text: '', is_correct: false },
]

export function ArgumentModal({
  isOpen,
  onClose,
  materialId,
  argumentToEdit,
}: ArgumentModalProps) {
  const isEditing = !!argumentToEdit
  const [claimText, setClaimText] = useState('')
  const [groundOptions, setGroundOptions] = useState<OptionInput[]>(createInitialGrounds)
  const [warrantOptions, setWarrantOptions] = useState<OptionInput[]>(createInitialWarrants)
  const [errorMsg, setErrorMsg] = useState('')

  const createMutation = useCreateArgument()
  const updateMutation = useUpdateArgument()

  useEffect(() => {
    if (argumentToEdit && argumentToEdit.options) {
      setClaimText(argumentToEdit.claim_text)
      const grounds = argumentToEdit.options.filter((o) => o.type === 'ground')
      const warrants = argumentToEdit.options.filter((o) => o.type === 'warrant')

      const loadedGrounds: OptionInput[] = grounds.map((g) => ({
        type: 'ground',
        text: g.text,
        is_correct: g.is_correct,
      }))
      while (loadedGrounds.length < 3) {
        loadedGrounds.push({ type: 'ground', text: '', is_correct: false })
      }
      if (!loadedGrounds.some((g) => g.is_correct) && loadedGrounds.length > 0) {
        loadedGrounds[0].is_correct = true
      }

      const loadedWarrants: OptionInput[] = warrants.map((w) => ({
        type: 'warrant',
        text: w.text,
        is_correct: w.is_correct,
      }))
      while (loadedWarrants.length < 3) {
        loadedWarrants.push({ type: 'warrant', text: '', is_correct: false })
      }
      if (!loadedWarrants.some((w) => w.is_correct) && loadedWarrants.length > 0) {
        loadedWarrants[0].is_correct = true
      }

      setGroundOptions(loadedGrounds)
      setWarrantOptions(loadedWarrants)
    } else {
      setClaimText('')
      setGroundOptions(createInitialGrounds())
      setWarrantOptions(createInitialWarrants())
    }
    setErrorMsg('')
  }, [argumentToEdit, isOpen])

  const setCorrectGround = (index: number) => {
    setGroundOptions((prev) =>
      prev.map((opt, i) => ({ ...opt, is_correct: i === index }))
    )
  }

  const setCorrectWarrant = (index: number) => {
    setWarrantOptions((prev) =>
      prev.map((opt, i) => ({ ...opt, is_correct: i === index }))
    )
  }

  const updateGroundText = (index: number, text: string) => {
    setGroundOptions((prev) =>
      prev.map((opt, i) => (i === index ? { ...opt, text } : opt))
    )
  }

  const updateWarrantText = (index: number, text: string) => {
    setWarrantOptions((prev) =>
      prev.map((opt, i) => (i === index ? { ...opt, text } : opt))
    )
  }

  const addGroundOption = () => {
    setGroundOptions((prev) => [
      ...prev,
      { type: 'ground', text: '', is_correct: false },
    ])
  }

  const removeGroundOption = (index: number) => {
    if (groundOptions.length <= 3) return
    const wasCorrect = groundOptions[index].is_correct
    const next = groundOptions.filter((_, i) => i !== index)
    if (wasCorrect && next.length > 0) {
      next[0].is_correct = true
    }
    setGroundOptions(next)
  }

  const addWarrantOption = () => {
    setWarrantOptions((prev) => [
      ...prev,
      { type: 'warrant', text: '', is_correct: false },
    ])
  }

  const removeWarrantOption = (index: number) => {
    if (warrantOptions.length <= 3) return
    const wasCorrect = warrantOptions[index].is_correct
    const next = warrantOptions.filter((_, i) => i !== index)
    if (wasCorrect && next.length > 0) {
      next[0].is_correct = true
    }
    setWarrantOptions(next)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMsg('')

    const allOptions = [...groundOptions, ...warrantOptions]
    const validation = argumentSchema.safeParse({
      claim_text: claimText,
      options: allOptions,
    })

    if (!validation.success) {
      const issue = validation.error.issues[0]
      setErrorMsg(issue ? issue.message : 'Data argumen tidak valid')
      return
    }

    try {
      if (isEditing && argumentToEdit) {
        await updateMutation.mutateAsync({
          materialId,
          argumentId: argumentToEdit.id,
          data: {
            claim_text: validation.data.claim_text,
            options: validation.data.options,
          },
        })
      } else {
        await createMutation.mutateAsync({
          materialId,
          data: {
            claim_text: validation.data.claim_text,
            options: validation.data.options,
          },
        })
      }
      onClose()
    } catch (err: any) {
      setErrorMsg(err.message || 'Gagal menyimpan argumen')
    }
  }

  const isPending = createMutation.isPending || updateMutation.isPending

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Ubah Argumen & Opsi' : 'Tambah Argumen Baru'}
      description="Model Toulmin: Claim + Opsi Ground (min. 3, 1 Benar) + Opsi Warrant (min. 3, 1 Benar)"
      maxWidth="4xl"
    >
      <form onSubmit={handleSubmit} className="space-y-6">
        {errorMsg && (
          <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg flex items-center gap-2">
            <AlertCircle className="w-4 h-4 shrink-0" />
            <span>{errorMsg}</span>
          </div>
        )}

        {/* Claim Input */}
        <div className="space-y-1.5 bg-slate-50 p-4 rounded-xl border border-slate-200">
          <Label htmlFor="claim-text" className="font-bold text-slate-800">
            Pernyataan Claim (Klaim Argumen)
          </Label>
          <Input
            id="claim-text"
            placeholder="misal: Verifikasi fakta mengurangi penyebaran hoaks di internet."
            value={claimText}
            onChange={(e) => setClaimText(e.target.value)}
            disabled={isPending}
            className="bg-white"
            autoFocus
          />
        </div>

        {/* Ground and Warrant Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Ground Options */}
          <div className="space-y-3 bg-amber-50/40 p-4 rounded-xl border border-amber-200/60 flex flex-col justify-between">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <div>
                  <span className="text-sm font-bold text-amber-900">
                    Opsi Ground ({groundOptions.length} opsi)
                  </span>
                  <p className="text-[11px] text-amber-700">Data/fakta pendukung, pilih 1 benar (min. 3)</p>
                </div>
              </div>

              <div className="space-y-2.5 max-h-[380px] overflow-y-auto pr-1">
                {groundOptions.map((opt, idx) => (
                  <div
                    key={idx}
                    className={`p-2.5 rounded-lg border transition-all ${
                      opt.is_correct
                        ? 'bg-emerald-50/80 border-emerald-300 ring-1 ring-emerald-300'
                        : 'bg-white border-slate-200'
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1.5">
                      <span className="text-xs font-semibold text-slate-600">
                        Ground #{idx + 1}
                      </span>
                      <div className="flex items-center gap-1.5">
                        <button
                          type="button"
                          onClick={() => setCorrectGround(idx)}
                          className={`text-xs font-semibold px-2 py-0.5 rounded-full flex items-center gap-1 cursor-pointer transition-colors ${
                            opt.is_correct
                              ? 'bg-emerald-600 text-white shadow-xs'
                              : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                          }`}
                        >
                          <CheckCircle2 className="w-3 h-3" />
                          <span>{opt.is_correct ? 'Kunci Benar' : 'Jadikan Benar'}</span>
                        </button>
                        {groundOptions.length > 3 && (
                          <button
                            type="button"
                            onClick={() => removeGroundOption(idx)}
                            className="p-1 text-slate-400 hover:text-red-600 rounded-md transition-colors cursor-pointer"
                            title="Hapus opsi ini"
                          >
                            <Trash2 className="w-3.5 h-3.5" />
                          </button>
                        )}
                      </div>
                    </div>
                    <Input
                      placeholder={`Teks opsi ground ${idx + 1}...`}
                      value={opt.text}
                      onChange={(e) => updateGroundText(idx, e.target.value)}
                      disabled={isPending}
                      className="bg-white text-xs"
                    />
                  </div>
                ))}
              </div>
            </div>

            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={addGroundOption}
              className="w-full mt-2 border-dashed border-amber-300 text-amber-800 hover:bg-amber-100/50 cursor-pointer"
            >
              <Plus className="w-3.5 h-3.5 mr-1" />
              <span>Tambah Opsi Ground</span>
            </Button>
          </div>

          {/* Warrant Options */}
          <div className="space-y-3 bg-indigo-50/40 p-4 rounded-xl border border-indigo-200/60 flex flex-col justify-between">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <div>
                  <span className="text-sm font-bold text-indigo-900">
                    Opsi Warrant ({warrantOptions.length} opsi)
                  </span>
                  <p className="text-[11px] text-indigo-700">Penalaran/penjamin klaim, pilih 1 benar (min. 3)</p>
                </div>
              </div>

              <div className="space-y-2.5 max-h-[380px] overflow-y-auto pr-1">
                {warrantOptions.map((opt, idx) => (
                  <div
                    key={idx}
                    className={`p-2.5 rounded-lg border transition-all ${
                      opt.is_correct
                        ? 'bg-emerald-50/80 border-emerald-300 ring-1 ring-emerald-300'
                        : 'bg-white border-slate-200'
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1.5">
                      <span className="text-xs font-semibold text-slate-600">
                        Warrant #{idx + 1}
                      </span>
                      <div className="flex items-center gap-1.5">
                        <button
                          type="button"
                          onClick={() => setCorrectWarrant(idx)}
                          className={`text-xs font-semibold px-2 py-0.5 rounded-full flex items-center gap-1 cursor-pointer transition-colors ${
                            opt.is_correct
                              ? 'bg-emerald-600 text-white shadow-xs'
                              : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                          }`}
                        >
                          <CheckCircle2 className="w-3 h-3" />
                          <span>{opt.is_correct ? 'Kunci Benar' : 'Jadikan Benar'}</span>
                        </button>
                        {warrantOptions.length > 3 && (
                          <button
                            type="button"
                            onClick={() => removeWarrantOption(idx)}
                            className="p-1 text-slate-400 hover:text-red-600 rounded-md transition-colors cursor-pointer"
                            title="Hapus opsi ini"
                          >
                            <Trash2 className="w-3.5 h-3.5" />
                          </button>
                        )}
                      </div>
                    </div>
                    <Input
                      placeholder={`Teks opsi warrant ${idx + 1}...`}
                      value={opt.text}
                      onChange={(e) => updateWarrantText(idx, e.target.value)}
                      disabled={isPending}
                      className="bg-white text-xs"
                    />
                  </div>
                ))}
              </div>
            </div>

            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={addWarrantOption}
              className="w-full mt-2 border-dashed border-indigo-300 text-indigo-800 hover:bg-indigo-100/50 cursor-pointer"
            >
              <Plus className="w-3.5 h-3.5 mr-1" />
              <span>Tambah Opsi Warrant</span>
            </Button>
          </div>
        </div>

        <div className="flex items-center justify-end space-x-2 pt-4 border-t border-slate-100">
          <Button type="button" variant="outline" onClick={onClose} disabled={isPending}>
            Batal
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? 'Menyimpan...' : isEditing ? 'Simpan Argumen' : 'Tambahkan Argumen'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
