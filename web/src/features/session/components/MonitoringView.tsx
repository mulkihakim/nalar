import { Modal } from '@/components/Modal'
import { Button } from '@/components/ui/button'
import {
  HoverCard,
  HoverCardTrigger,
  HoverCardContent,
} from '@/components/ui/hover-card'
import type { MonitoringAnalyticsResponse } from '../types'
import { Users, ArrowRight } from 'lucide-react'

interface MonitoringViewProps {
  isOpen: boolean
  onClose: () => void
  data?: MonitoringAnalyticsResponse
  claimText?: string
  isLoading?: boolean
  isSuccessDialog?: boolean
  onNextArgument?: () => void
  isLastArgument?: boolean
}

export function MonitoringView({
  isOpen,
  onClose,
  data,
  claimText,
  isLoading,
  isSuccessDialog,
  onNextArgument,
  isLastArgument,
}: MonitoringViewProps) {
  const groundOptions = data?.options.filter((o) => o.type === 'ground') || []
  const warrantOptions = data?.options.filter((o) => o.type === 'warrant') || []

  const activeClaim = data?.claim_text || claimText || ''

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Monitoring Pola Pilihan Teman Sejawat"
      maxWidth="4xl"
    >
      <div className="space-y-4 py-1">
        {/* Claim Box */}
        {activeClaim && (
          <div className="bg-slate-900 text-white p-4 rounded-xl shadow-xs border border-slate-800">
            <span className="text-[10px] font-bold tracking-widest text-indigo-400 uppercase block mb-1">
              Claim (Pernyataan / Klaim Utama)
            </span>
            <p className="text-sm md:text-base font-bold leading-relaxed">{activeClaim}</p>
          </div>
        )}

        {isLoading ? (
          <div className="space-y-4 py-3 animate-pulse">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2.5">
                <div className="h-4 bg-slate-200 rounded w-1/3"></div>
                <div className="h-16 bg-slate-50 rounded-xl border border-slate-200"></div>
                <div className="h-16 bg-slate-50 rounded-xl border border-slate-200"></div>
              </div>
              <div className="space-y-2.5">
                <div className="h-4 bg-slate-200 rounded w-1/3"></div>
                <div className="h-16 bg-slate-50 rounded-xl border border-slate-200"></div>
                <div className="h-16 bg-slate-50 rounded-xl border border-slate-200"></div>
              </div>
            </div>
            <p className="text-center text-slate-400 text-xs pt-1">Memuat perbandingan pola pilihan teman sejawat...</p>
          </div>
        ) : !data ? (
          <div className="py-12 text-center text-slate-500 text-sm">Tidak ada data analitik.</div>
        ) : (
          <>
            {!data.has_enough_peers && (
              <div className="p-3 bg-amber-50 border border-amber-200 rounded-xl text-amber-800 text-xs flex items-center gap-2">
                <Users className="w-4 h-4 text-amber-600 shrink-0" />
                <span>
                  Data agregat kelompok disembunyikan demi menjaga privasi (memerlukan minimal{' '}
                  {data.min_peers} siswa; saat ini baru {data.total_peers} siswa).
                </span>
              </div>
            )}

            {/* 2 Kolom Opsi: Ground di Kiri, Warrant di Kanan */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-h-[440px] overflow-y-auto pr-1">
              {/* Kolom Ground */}
              <div className="space-y-2.5">
                <div className="flex items-center justify-between px-1">
                  <h4 className="text-xs font-bold text-slate-800 flex items-center gap-1.5 uppercase tracking-wide">
                    <span className="w-2 h-2 rounded-full bg-blue-500"></span>
                    Opsi Ground
                  </h4>
                  <span className="text-[11px] text-slate-400 font-medium">Percobaan Anda : Rata-rata Teman</span>
                </div>

                <div className="space-y-2">
                  {groundOptions.map((opt) => (
                    <div
                      key={opt.option_id}
                      className="p-3 rounded-xl border border-slate-200 bg-white shadow-xs flex flex-col justify-between gap-1.5"
                    >
                      <div className="flex items-start justify-between gap-3">
                        <span className="text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider bg-blue-50 text-blue-700 border border-blue-200 shrink-0">
                          Ground
                        </span>
                        <HoverCard openDelay={100} closeDelay={150}>
                          <HoverCardTrigger asChild>
                            <button
                              type="button"
                              className="text-xs font-bold text-slate-900 tracking-wider bg-slate-100 hover:bg-slate-200 px-2 py-0.5 rounded-md cursor-help transition-colors"
                            >
                              {opt.x} : {opt.y}
                            </button>
                          </HoverCardTrigger>
                          <HoverCardContent
                            align="end"
                            className="w-60 p-3 bg-white border border-slate-200 shadow-md rounded-xl text-xs space-y-2 z-50 text-slate-700"
                          >
                            <div className="font-semibold text-slate-900 flex items-center gap-1.5 border-b border-slate-100 pb-1.5">
                              <Users className="w-3.5 h-3.5 text-indigo-600" />
                              <span>Pola Pilihan Ground</span>
                            </div>
                            <div className="text-[11px] leading-relaxed space-y-1.5">
                              <div className="flex items-center justify-between">
                                <span className="text-slate-500">Percobaan Anda:</span>
                                <span className="font-bold text-slate-900 bg-slate-100 px-1.5 py-0.5 rounded">
                                  {opt.x} kali
                                </span>
                              </div>
                              <div className="flex items-center justify-between">
                                <span className="text-slate-500">Rata-rata Teman:</span>
                                <span className="font-bold text-slate-900 bg-slate-100 px-1.5 py-0.5 rounded">
                                  {opt.y === '-' ? 'Belum ada data' : `${opt.y} kali`}
                                </span>
                              </div>
                            </div>
                          </HoverCardContent>
                        </HoverCard>
                      </div>
                      <p className="text-xs text-slate-800 leading-relaxed">{opt.text}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Kolom Warrant */}
              <div className="space-y-2.5">
                <div className="flex items-center justify-between px-1">
                  <h4 className="text-xs font-bold text-slate-800 flex items-center gap-1.5 uppercase tracking-wide">
                    <span className="w-2 h-2 rounded-full bg-purple-500"></span>
                    Opsi Warrant
                  </h4>
                  <span className="text-[11px] text-slate-400 font-medium">Percobaan Anda : Rata-rata Teman</span>
                </div>

                <div className="space-y-2">
                  {warrantOptions.map((opt) => (
                    <div
                      key={opt.option_id}
                      className="p-3 rounded-xl border border-slate-200 bg-white shadow-xs flex flex-col justify-between gap-1.5"
                    >
                      <div className="flex items-start justify-between gap-3">
                        <span className="text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider bg-purple-50 text-purple-700 border border-purple-200 shrink-0">
                          Warrant
                        </span>
                        <HoverCard openDelay={100} closeDelay={150}>
                          <HoverCardTrigger asChild>
                            <button
                              type="button"
                              className="text-xs font-bold text-slate-900 tracking-wider bg-slate-100 hover:bg-slate-200 px-2 py-0.5 rounded-md cursor-help transition-colors"
                            >
                              {opt.x} : {opt.y}
                            </button>
                          </HoverCardTrigger>
                          <HoverCardContent
                            align="end"
                            className="w-60 p-3 bg-white border border-slate-200 shadow-md rounded-xl text-xs space-y-2 z-50 text-slate-700"
                          >
                            <div className="font-semibold text-slate-900 flex items-center gap-1.5 border-b border-slate-100 pb-1.5">
                              <Users className="w-3.5 h-3.5 text-indigo-600" />
                              <span>Pola Pilihan Warrant</span>
                            </div>
                            <div className="text-[11px] leading-relaxed space-y-1.5">
                              <div className="flex items-center justify-between">
                                <span className="text-slate-500">Percobaan Anda:</span>
                                <span className="font-bold text-slate-900 bg-slate-100 px-1.5 py-0.5 rounded">
                                  {opt.x} kali
                                </span>
                              </div>
                              <div className="flex items-center justify-between">
                                <span className="text-slate-500">Rata-rata Teman:</span>
                                <span className="font-bold text-slate-900 bg-slate-100 px-1.5 py-0.5 rounded">
                                  {opt.y === '-' ? 'Belum ada data' : `${opt.y} kali`}
                                </span>
                              </div>
                            </div>
                          </HoverCardContent>
                        </HoverCard>
                      </div>
                      <p className="text-xs text-slate-800 leading-relaxed">{opt.text}</p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </>
        )}

        {/* Footer Buttons */}
        <div className="flex items-center justify-end gap-2 pt-3 border-t border-slate-100">
          <Button variant="outline" size="sm" onClick={onClose} className="text-xs">
            Tutup
          </Button>
          {isSuccessDialog && onNextArgument && (
            <Button
              size="sm"
              onClick={() => {
                onClose()
                onNextArgument()
              }}
              className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold"
            >
              {isLastArgument ? 'Selesaikan Ujian' : 'Lanjut ke Argumen Berikutnya'}
              <ArrowRight className="w-3.5 h-3.5 ml-1" />
            </Button>
          )}
        </div>
      </div>
    </Modal>
  )
}
