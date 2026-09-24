import { useState } from 'react'
import type { SessionDetail, AnalysisAnalyticsResponse } from '../types'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { AnalysisView } from './AnalysisView'
import { CheckCircle2, RotateCcw, ListChecks, Award, Users } from 'lucide-react'

interface SessionSummaryProps {
  session: SessionDetail
  analysisData?: AnalysisAnalyticsResponse
  isLoadingAnalysis?: boolean
  onRestart: () => void
  onBackToList: () => void
  onReviewQuestions?: () => void
}

export function SessionSummary({
  session,
  analysisData,
  isLoadingAnalysis,
  onRestart,
  onBackToList,
  onReviewQuestions,
}: SessionSummaryProps) {
  const [activeTab, setActiveTab] = useState<'summary' | 'social'>('summary')

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header Banner */}
      <div className="bg-emerald-600 text-white p-6 md:p-8 rounded-3xl shadow-sm flex flex-col md:flex-row items-center justify-between gap-6">
        <div className="flex items-center gap-4 text-center md:text-left">
          <div className="w-14 h-14 bg-white/10 rounded-2xl flex items-center justify-center flex-shrink-0">
            <Award className="w-8 h-8 text-emerald-100" />
          </div>
          <div>
            <span className="text-xs font-bold uppercase tracking-wider text-emerald-200">
              Percobaan #{session.attempt_no} Selesai
            </span>
            <h2 className="text-xl md:text-2xl font-bold mt-0.5">Latihan Berhasil Diselesaikan!</h2>
            <p className="text-xs md:text-sm text-emerald-100 mt-1">
              Anda telah menyusun seluruh argumen Toulmin dengan tepat pada sesi ini.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          {onReviewQuestions && (
            <Button
              variant="outline"
              onClick={onReviewQuestions}
              className="border-white/30 text-white bg-white/10 hover:bg-white/20"
            >
              Tinjau Soal
            </Button>
          )}
          <Button
            variant="outline"
            onClick={onBackToList}
            className="border-white/30 text-white bg-white/10 hover:bg-white/20"
          >
            Daftar Ujian
          </Button>
          <Button
            onClick={onRestart}
            className="bg-white text-emerald-800 hover:bg-emerald-50 font-semibold"
          >
            <RotateCcw className="w-4 h-4 mr-2" />
            Latihan Lagi
          </Button>
        </div>
      </div>

      {/* Mode Sosial Tab Switcher */}
      {session.mode === 'social' && (
        <div className="flex border-b border-slate-200">
          <button
            type="button"
            onClick={() => setActiveTab('summary')}
            className={`pb-3 px-4 text-sm font-semibold flex items-center gap-2 border-b-2 transition-all ${
              activeTab === 'summary'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700'
            }`}
          >
            <ListChecks className="w-4 h-4" />
            Ringkasan Sesi Anda
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('social')}
            className={`pb-3 px-4 text-sm font-semibold flex items-center gap-2 border-b-2 transition-all ${
              activeTab === 'social'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700'
            }`}
          >
            <Users className="w-4 h-4" />
            Analitik Sosial Kelompok
          </button>
        </div>
      )}

      {/* Tab Konten: Summary */}
      {activeTab === 'summary' && (
        <div className="space-y-6">
          {/* Stat Cards */}
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            <div className="p-5 rounded-2xl bg-white border border-slate-200 shadow-sm">
              <span className="text-xs font-semibold text-slate-400 uppercase">Argumen Selesai</span>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                {session.completed_arguments} / {session.total_arguments}
              </p>
            </div>
            <div className="p-5 rounded-2xl bg-white border border-slate-200 shadow-sm">
              <span className="text-xs font-semibold text-slate-400 uppercase">Mode Belajar</span>
              <p className="text-xl font-bold text-slate-900 mt-1 capitalize">{session.mode}</p>
            </div>
            <div className="p-5 rounded-2xl bg-white border border-slate-200 shadow-sm col-span-2 md:col-span-1">
              <span className="text-xs font-semibold text-slate-400 uppercase">Status</span>
              <div className="mt-1">
                <Badge variant="outline" className="bg-emerald-50 text-emerald-700 border-emerald-200">
                  <CheckCircle2 className="w-3.5 h-3.5 mr-1" /> Selesai
                </Badge>
              </div>
            </div>
          </div>

          {/* List Argumen */}
          <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6">
            <h4 className="font-bold text-slate-900 mb-4 text-base">Rincian Argumen yang Dipelajari</h4>
            <div className="divide-y divide-slate-100">
              {session.arguments.map((arg, idx) => (
                <div key={arg.id} className="py-3.5 flex items-center justify-between gap-4">
                  <div className="flex items-start gap-3">
                    <span className="w-6 h-6 rounded-full bg-slate-100 text-slate-600 text-xs font-bold flex items-center justify-center flex-shrink-0 mt-0.5">
                      {idx + 1}
                    </span>
                    <div>
                      <p className="text-sm font-medium text-slate-900 leading-snug">
                        {arg.claim_text}
                      </p>
                    </div>
                  </div>
                  <Badge variant="outline" className="bg-emerald-50 text-emerald-700 border-emerald-200 text-xs flex-shrink-0">
                    Selesai
                  </Badge>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Tab Konten: Social Analytics */}
      {activeTab === 'social' && (
        <AnalysisView data={analysisData} isLoading={isLoadingAnalysis} />
      )}
    </div>
  )
}
