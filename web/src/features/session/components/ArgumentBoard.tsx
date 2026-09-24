import { useState, useEffect } from 'react'
import {
  DndContext,
  type DragEndEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import type { ArgumentItem, OptionItem, ConfirmResult, SessionMode } from '../types'
import { OptionCard } from './OptionCard'
import { SlotZone } from './SlotZone'
import { Button } from '@/components/ui/button'
import { Modal } from '@/components/Modal'
import { ArrowRight, CheckCircle2, AlertTriangle, Users, Sparkles } from 'lucide-react'

interface ArgumentBoardProps {
  argument: ArgumentItem
  mode: SessionMode
  onDropOption: (option: OptionItem, slot: 'ground' | 'warrant') => void
  onConfirm: (groundOptionId: number, warrantOptionId: number) => Promise<ConfirmResult | null>
  onNextArgument: () => void
  isLastArgument: boolean
  isSubmitting?: boolean
  showMonitoring?: boolean
  onShowMonitoring?: (isSuccess?: boolean) => void
  isReadOnly?: boolean
  isAllCompleted?: boolean
  onFinishExam?: () => void
}

export function ArgumentBoard({
  argument,
  mode,
  onDropOption,
  onConfirm,
  onNextArgument,
  isLastArgument,
  isSubmitting,
  showMonitoring,
  onShowMonitoring,
  isReadOnly = false,
  isAllCompleted = false,
  onFinishExam,
}: ArgumentBoardProps) {
  const [chosenGround, setChosenGround] = useState<OptionItem | null>(null)
  const [chosenWarrant, setChosenWarrant] = useState<OptionItem | null>(null)
  const [confirmResult, setConfirmResult] = useState<ConfirmResult | null>(null)
  const [isResultModalOpen, setIsResultModalOpen] = useState(false)
  const [isDone, setIsDone] = useState<boolean>(argument.completed)
  const [autoCloseSeconds, setAutoCloseSeconds] = useState<number | null>(null)

  useEffect(() => {
    setIsDone(argument.completed)
    if (argument.completed) {
      const correctGround = argument.options.find((o) => o.type === 'ground' && o.is_correct)
      const correctWarrant = argument.options.find((o) => o.type === 'warrant' && o.is_correct)
      setChosenGround(correctGround || null)
      setChosenWarrant(correctWarrant || null)
    } else {
      setChosenGround(null)
      setChosenWarrant(null)
    }
    setConfirmResult(null)
    setIsResultModalOpen(false)
  }, [argument.id, argument.completed, argument.options])

  useEffect(() => {
    if (isResultModalOpen && confirmResult && !confirmResult.is_correct && mode !== 'help') {
      setAutoCloseSeconds(3)
      const timer = setInterval(() => {
        setAutoCloseSeconds((prev) => {
          if (prev === null || prev <= 1) {
            clearInterval(timer)
            setIsResultModalOpen(false)
            return null
          }
          return prev - 1
        })
      }, 1000)
      return () => clearInterval(timer)
    } else {
      setAutoCloseSeconds(null)
    }
  }, [isResultModalOpen, confirmResult, mode])

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 5,
      },
    })
  )

  const groundOptions = argument.options.filter((o) => o.type === 'ground')
  const warrantOptions = argument.options.filter((o) => o.type === 'warrant')

  const handleDragEnd = (event: DragEndEvent) => {
    if (isReadOnly || isDone) return
    const { active, over } = event
    if (!over) return

    const option = active.data.current?.option as OptionItem | undefined
    const slot = over.data.current?.slot as 'ground' | 'warrant' | undefined

    if (!option || !slot) return

    // Validasi opsi hanya bisa masuk ke slot sesuai tipenya
    if (option.type !== slot) {
      return
    }

    if (slot === 'ground') {
      setChosenGround(option)
      onDropOption(option, 'ground')
    } else if (slot === 'warrant') {
      setChosenWarrant(option)
      onDropOption(option, 'warrant')
    }

    if (confirmResult && !confirmResult.is_correct) {
      setConfirmResult(null)
    }
  }

  const handleSelectDirectly = (option: OptionItem) => {
    if (isReadOnly || isDone) return
    if (option.type === 'ground') {
      setChosenGround(option)
      onDropOption(option, 'ground')
    } else {
      setChosenWarrant(option)
      onDropOption(option, 'warrant')
    }
    if (confirmResult && !confirmResult.is_correct) {
      setConfirmResult(null)
    }
  }

  const handleConfirmClick = async () => {
    if (!chosenGround || !chosenWarrant || isReadOnly || isDone) return

    const res = await onConfirm(chosenGround.id, chosenWarrant.id)
    if (res) {
      setConfirmResult(res)
      if (res.completed) {
        setIsDone(true)
      }
      if (res.is_correct && mode === 'social' && onShowMonitoring) {
        onShowMonitoring(true)
      } else {
        setIsResultModalOpen(true)
      }
    }
  }

  // Tentukan status validasi per slot untuk mode bantuan
  let groundStatus: 'correct' | 'incorrect' | null = null
  let warrantStatus: 'correct' | 'incorrect' | null = null

  if (confirmResult) {
    if (confirmResult.is_correct) {
      groundStatus = 'correct'
      warrantStatus = 'correct'
    } else if (mode === 'help') {
      if (confirmResult.ground_correct !== undefined) {
        groundStatus = confirmResult.ground_correct ? 'correct' : 'incorrect'
      }
      if (confirmResult.warrant_correct !== undefined) {
        warrantStatus = confirmResult.warrant_correct ? 'correct' : 'incorrect'
      }
    }
  } else if (isDone || argument.completed) {
    groundStatus = 'correct'
    warrantStatus = 'correct'
  }

  return (
    <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
      <div className="space-y-4">
        {/* Claim Header (Toulmin Model Top) */}
        <div className="bg-slate-900 text-white p-5 rounded-2xl shadow-sm border border-slate-800">
          <div className="flex items-center justify-between mb-1.5">
            <span className="text-xs font-bold tracking-widest text-indigo-400 uppercase">
              Claim (Pernyataan / Klaim Utama)
            </span>
            {isDone && (
              <span className="text-xs text-emerald-400 font-semibold flex items-center gap-1">
                <CheckCircle2 className="w-3.5 h-3.5" /> Sudah Dijawab Benar
              </span>
            )}
          </div>
          <h3 className="text-base md:text-lg font-bold leading-relaxed">{argument.claim_text}</h3>
        </div>

        {/* Toulmin Bridge Label (Menggantikan Teks Struktur Argumen Toulmin) */}
        <div className="flex items-center justify-center py-1 text-slate-500 gap-2 sm:gap-4 text-xs font-medium bg-slate-100/80 rounded-xl px-4 border border-slate-200/60">
          <span className="font-semibold text-blue-600">Ground (Fakta)</span>
          <ArrowRight className="w-3.5 h-3.5 text-slate-400 shrink-0" />
          <span className="bg-purple-100/80 text-purple-700 px-2.5 py-0.5 rounded-full font-semibold border border-purple-200 text-[11px]">
            Dijembatani oleh Warrant
          </span>
          <ArrowRight className="w-3.5 h-3.5 text-slate-400 shrink-0" />
          <span className="font-semibold text-slate-900">Claim (Klaim)</span>
        </div>

        {/* 2 Slot Bersanding (Ground & Warrant) */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Slot Ground (Kiri) */}
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-xs">
            <SlotZone
              slotType="ground"
              chosenOption={chosenGround || undefined}
              onClear={() => {
                if (!isDone && !isReadOnly) setChosenGround(null)
              }}
              disabled={isDone || isReadOnly}
              validationStatus={groundStatus}
            />
          </div>

          {/* Slot Warrant (Kanan) */}
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-xs">
            <SlotZone
              slotType="warrant"
              chosenOption={chosenWarrant || undefined}
              onClear={() => {
                if (!isDone && !isReadOnly) setChosenWarrant(null)
              }}
              disabled={isDone || isReadOnly}
              validationStatus={warrantStatus}
            />
          </div>
        </div>

        {/* Banner Review jika Soal ini Sudah Selesai */}
        {(isDone || argument.completed) && (
          <div className="p-4 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-900 flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-2.5 text-xs sm:text-sm">
              <CheckCircle2 className="w-5 h-5 text-emerald-600 shrink-0" />
              <span>
                Argumen ini telah diselesaikan dengan benar pada sesi ini.
              </span>
            </div>
            <div className="flex items-center gap-2 shrink-0">
              {mode === 'social' && showMonitoring && onShowMonitoring && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onShowMonitoring(false)}
                  className="border-emerald-300 text-emerald-800 bg-white hover:bg-emerald-50 text-xs"
                >
                  <Users className="w-3.5 h-3.5 mr-1" />
                  Lihat Pola Teman
                </Button>
              )}
              {isAllCompleted && onFinishExam && (
                <Button
                  size="sm"
                  onClick={onFinishExam}
                  className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold"
                >
                  Lihat Hasil Ujian
                  <ArrowRight className="w-3.5 h-3.5 ml-1" />
                </Button>
              )}
            </div>
          </div>
        )}

        {/* Opsi Pilihan Ground & Warrant Tepat di Bawah Slot Masing-Masing */}
        {!isDone && !isReadOnly && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-1">
            {/* Opsi Ground */}
            <div className="space-y-2.5">
              <div className="flex items-center justify-between px-1">
                <h4 className="text-xs font-bold text-slate-800 flex items-center gap-1.5 uppercase tracking-wide">
                  <span className="w-2 h-2 rounded-full bg-blue-500"></span>
                  Pilihan Ground (Data / Bukti)
                </h4>
                <span className="text-[11px] text-slate-400">Seret ke slot atau klik</span>
              </div>
              <div className="space-y-2">
                {groundOptions.map((opt) => (
                  <div key={opt.id} onClick={() => handleSelectDirectly(opt)}>
                    <OptionCard
                      option={opt}
                      isChosen={chosenGround?.id === opt.id}
                      disabled={isDone}
                    />
                  </div>
                ))}
              </div>
            </div>

            {/* Opsi Warrant */}
            <div className="space-y-2.5">
              <div className="flex items-center justify-between px-1">
                <h4 className="text-xs font-bold text-slate-800 flex items-center gap-1.5 uppercase tracking-wide">
                  <span className="w-2 h-2 rounded-full bg-purple-500"></span>
                  Pilihan Warrant (Prinsip / Alasan)
                </h4>
                <span className="text-[11px] text-slate-400">Seret ke slot atau klik</span>
              </div>
              <div className="space-y-2">
                {warrantOptions.map((opt) => (
                  <div key={opt.id} onClick={() => handleSelectDirectly(opt)}>
                    <OptionCard
                      option={opt}
                      isChosen={chosenWarrant?.id === opt.id}
                      disabled={isDone}
                    />
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Tombol Confirm Jawaban di Paling Bawah */}
        {!isDone && !isReadOnly && (
          <div className="pt-4 flex justify-end">
            <Button
              size="lg"
              disabled={!chosenGround || !chosenWarrant || isSubmitting}
              onClick={handleConfirmClick}
              className="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-8 shadow-sm text-sm"
            >
              {isSubmitting ? 'Memeriksa...' : 'Confirm Jawaban'}
            </Button>
          </div>
        )}

        {/* Jika semua argumen sudah selesai, tampilkan tombol ke Halaman Hasil */}
        {(isDone || isReadOnly || argument.completed) && isAllCompleted && onFinishExam && (
          <div className="pt-4 flex justify-end">
            <Button
              size="lg"
              onClick={onFinishExam}
              className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold px-8 shadow-sm text-sm"
            >
              Lihat Hasil Ujian
              <ArrowRight className="w-4 h-4 ml-2" />
            </Button>
          </div>
        )}

        {/* Dialog Evaluasi Jawaban Shadcn Modal */}
        <Modal
          isOpen={isResultModalOpen}
          onClose={() => setIsResultModalOpen(false)}
          title={confirmResult?.is_correct ? 'Jawaban Benar!' : 'Susunan Belum Tepat'}
          maxWidth="md"
        >
          <div className="py-2 space-y-4">
            {confirmResult?.is_correct ? (
              <div className="text-center py-2 space-y-3">
                <div className="w-12 h-12 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto">
                  <CheckCircle2 className="w-7 h-7" />
                </div>
                <div>
                  <h4 className="font-bold text-slate-900 text-base">Tepat Sekali!</h4>
                  <p className="text-xs text-slate-500 mt-1 leading-relaxed">
                    {confirmResult.message || 'Susunan argumen Anda sudah benar.'}
                  </p>
                </div>

                <div className="pt-4 flex flex-col gap-2">
                  {mode === 'social' && showMonitoring && onShowMonitoring && (
                    <Button
                      variant="outline"
                      onClick={() => {
                        setIsResultModalOpen(false)
                        onShowMonitoring(false)
                      }}
                      className="border-indigo-200 text-indigo-700 hover:bg-indigo-50 text-xs"
                    >
                      <Users className="w-4 h-4 mr-1.5" />
                      Lihat Pola Pilihan Teman (X : Y)
                    </Button>
                  )}
                  <Button
                    onClick={() => {
                      setIsResultModalOpen(false)
                      onNextArgument()
                    }}
                    className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold"
                  >
                    {isLastArgument ? 'Selesaikan Sesi Latihan' : 'Lanjut ke Argumen Berikutnya'}
                    <ArrowRight className="w-4 h-4 ml-1.5" />
                  </Button>
                </div>
              </div>
            ) : (
              <div className="py-2 space-y-3">
                <div className="flex items-start gap-3 p-3.5 bg-red-50 text-red-900 rounded-xl border border-red-200">
                  <AlertTriangle className="w-5 h-5 text-red-600 shrink-0 mt-0.5" />
                  <div className="text-xs space-y-1">
                    <p className="font-semibold text-sm">
                      {confirmResult?.message || 'Susunan argumen belum tepat.'}
                    </p>
                    <p className="text-red-700">
                      Silakan coba lagi dengan meninjau kembali kaitan antara fakta (Ground), alasan pengait (Warrant), dan pernyataan (Claim).
                    </p>
                  </div>
                </div>

                {mode === 'help' && (
                  <div className="p-3 bg-amber-50 rounded-xl border border-amber-200 text-xs text-amber-900 space-y-2">
                    <div className="flex items-center gap-1.5 font-semibold text-amber-800">
                      <Sparkles className="w-4 h-4 text-amber-600" />
                      Petunjuk Evaluasi Mode Bantuan:
                    </div>
                    <div className="grid grid-cols-2 gap-2 text-center">
                      <div className={`p-2 rounded-lg border font-medium ${groundStatus === 'correct' ? 'bg-emerald-100/70 border-emerald-300 text-emerald-800' : 'bg-red-100/70 border-red-300 text-red-800'}`}>
                        Ground: {groundStatus === 'correct' ? 'Benar ✓' : 'Belum Tepat ✗'}
                      </div>
                      <div className={`p-2 rounded-lg border font-medium ${warrantStatus === 'correct' ? 'bg-emerald-100/70 border-emerald-300 text-emerald-800' : 'bg-red-100/70 border-red-300 text-red-800'}`}>
                        Warrant: {warrantStatus === 'correct' ? 'Benar ✓' : 'Belum Tepat ✗'}
                      </div>
                    </div>
                    <p className="text-[11px] text-amber-700">
                      Perhatikan penanda pada masing-masing slot untuk mengetahui bagian mana yang perlu diperbaiki.
                    </p>
                  </div>
                )}

                {mode !== 'help' && autoCloseSeconds !== null && (
                  <p className="text-[11px] text-slate-500 text-center">
                    Pemberitahuan ini akan tertutup otomatis dalam <span className="font-semibold text-slate-700">{autoCloseSeconds} detik</span>...
                  </p>
                )}

                <div className="pt-2 flex justify-end">
                  <Button
                    onClick={() => setIsResultModalOpen(false)}
                    className="bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-semibold px-5"
                  >
                    Coba Lagi
                  </Button>
                </div>
              </div>
            )}
          </div>
        </Modal>
      </div>
    </DndContext>
  )
}
