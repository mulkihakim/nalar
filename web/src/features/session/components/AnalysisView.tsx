import type { AnalysisAnalyticsResponse, AnalysisOptionStat } from '../types'
import {
  HoverCard,
  HoverCardTrigger,
  HoverCardContent,
} from '@/components/ui/hover-card'
import { Info, Users, UserCheck } from 'lucide-react'

interface AnalysisViewProps {
  data?: AnalysisAnalyticsResponse
  isLoading?: boolean
  isStaffView?: boolean
  onSelectOption?: (option: AnalysisOptionStat) => void
}

export function AnalysisView({ data, isLoading, isStaffView = false, onSelectOption }: AnalysisViewProps) {
  if (isLoading) {
    return <div className="py-12 text-center text-slate-400">Memuat analisis kelompok...</div>
  }

  if (!data) {
    return <div className="py-12 text-center text-slate-400">Data analisis tidak tersedia.</div>
  }

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="bg-indigo-50/70 border border-indigo-100 p-5 rounded-2xl flex items-start gap-4">
        <Info className="w-6 h-6 text-indigo-600 flex-shrink-0 mt-0.5" />
        <div>
          <h4 className="font-bold text-slate-900 text-sm md:text-base">
            Analisis Kelompok Teman Sejawat
          </h4>
          <p className="text-xs text-slate-600 mt-1">
            Format angka perbandingan: <strong>Total Percobaan Memilih : Jumlah Siswa Memilih</strong>.
            {isStaffView && (
              <span className="block mt-1 text-indigo-700 font-medium">
                💡 Klik pada kartu opsi untuk melihat daftar nama siswa yang memilih opsi tersebut.
              </span>
            )}
          </p>
        </div>
      </div>

      {!data.has_enough_peers && (
        <div className="p-4 bg-amber-50 border border-amber-200 rounded-xl text-amber-800 text-xs flex items-center gap-2">
          <Users className="w-4 h-4 text-amber-600 flex-shrink-0" />
          <span>
            Data agregat kelompok disembunyikan demi menjaga privasi (memerlukan minimal {data.min_peers}{' '}
            siswa telah mengerjakan ujian; saat ini baru {data.total_peers} siswa).
          </span>
        </div>
      )}

      {/* List Argumen & Analisis Opsi */}
      <div className="space-y-6">
        {data.arguments.map((arg, idx) => (
          <div
            key={arg.argument_id}
            className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm space-y-4"
          >
            <div className="border-b border-slate-100 pb-3">
              <span className="text-[11px] font-bold text-indigo-600 uppercase tracking-wider">
                Argumen #{idx + 1}
              </span>
              <h5 className="font-bold text-slate-900 text-sm mt-0.5">{arg.claim_text}</h5>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {arg.options.map((opt) => {
                const isGround = opt.type === 'ground'
                const maxAttempts = Math.max(...arg.options.map((o) => o.total_attempts), 1)
                const percentage = data.has_enough_peers
                  ? Math.round((opt.total_attempts / maxAttempts) * 100)
                  : 0

                return (
                  <div
                    key={opt.option_id}
                    onClick={() => {
                      if (isStaffView && onSelectOption) {
                        onSelectOption(opt)
                      }
                    }}
                    className={`p-3.5 rounded-xl border border-slate-200 bg-slate-50/60 flex flex-col justify-between gap-2 transition-all ${
                      isStaffView
                        ? 'cursor-pointer hover:border-indigo-300 hover:bg-indigo-50/20 hover:shadow-xs'
                        : ''
                    }`}
                  >
                    <div>
                      <div className="flex items-center justify-between mb-1.5">
                        <span
                          className={`text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider ${
                            isGround
                              ? 'bg-blue-50 text-blue-700 border border-blue-200'
                              : 'bg-purple-50 text-purple-700 border border-purple-200'
                          }`}
                        >
                          {isGround ? 'Ground' : 'Warrant'}
                        </span>

                        <div className="flex items-center gap-1.5">
                          {isStaffView && opt.students && opt.students.length > 0 && (
                            <span className="text-[10px] text-slate-500 flex items-center gap-0.5">
                              <UserCheck className="w-3 h-3 text-indigo-600" />
                              {opt.students.length}
                            </span>
                          )}
                          {!isStaffView ? (
                            <HoverCard openDelay={100} closeDelay={150}>
                              <HoverCardTrigger asChild>
                                <button
                                  type="button"
                                  className="text-xs font-bold text-slate-900 tracking-wider bg-slate-100 hover:bg-slate-200 px-2 py-0.5 rounded-md cursor-help transition-colors"
                                >
                                  {opt.ratio}
                                </button>
                              </HoverCardTrigger>
                              <HoverCardContent
                                align="end"
                                className="w-64 p-3 bg-white border border-slate-200 shadow-md rounded-xl text-xs space-y-2 z-50 text-slate-700"
                              >
                                <div className="font-semibold text-slate-900 flex items-center gap-1.5 border-b border-slate-100 pb-1.5">
                                  <Users className="w-3.5 h-3.5 text-indigo-600" />
                                  <span>Statistik Kelompok Siswa</span>
                                </div>
                                <div className="text-[11px] leading-relaxed space-y-1.5">
                                  <div className="flex items-center justify-between">
                                    <span className="text-slate-500">Total Dipilih:</span>
                                    <span className="font-bold text-slate-900 bg-slate-100 px-1.5 py-0.5 rounded">
                                      {opt.total_attempts} kali
                                    </span>
                                  </div>
                                  <div className="flex items-center justify-between">
                                    <span className="text-slate-500">Siswa Unik:</span>
                                    <span className="font-bold text-slate-900 bg-slate-100 px-1.5 py-0.5 rounded">
                                      {opt.unique_students} orang
                                    </span>
                                  </div>
                                </div>
                              </HoverCardContent>
                            </HoverCard>
                          ) : (
                            <span className="text-xs font-bold text-slate-900 tracking-wider bg-slate-100 px-2 py-0.5 rounded-md">
                              {opt.ratio}
                            </span>
                          )}
                        </div>
                      </div>
                      <p className="text-xs text-slate-800 line-clamp-2">{opt.text}</p>
                    </div>

                    {/* Bar Distribusi */}
                    {data.has_enough_peers && (
                      <div className="w-full bg-slate-200 h-1.5 rounded-full overflow-hidden mt-1">
                        <div
                          className={`h-full rounded-full transition-all duration-500 ${
                            isGround ? 'bg-blue-500' : 'bg-purple-500'
                          }`}
                          style={{ width: `${percentage}%` }}
                        />
                      </div>
                    )}
                  </div>
                )
              })}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
