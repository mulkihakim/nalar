import type { StaffResultSummary } from '../types'
import { Badge } from '@/components/ui/badge'
import { CheckCircle2, Clock } from 'lucide-react'

interface ResultsTableProps {
  results: StaffResultSummary[]
  isLoading?: boolean
}

export function ResultsTable({ results, isLoading }: ResultsTableProps) {
  if (isLoading) {
    return <div className="py-12 text-center text-slate-400">Memuat data hasil siswa...</div>
  }

  if (results.length === 0) {
    return (
      <div className="py-12 text-center text-slate-400 text-sm">
        Belum ada siswa yang mengerjakan ujian ini.
      </div>
    )
  }

  return (
    <div className="overflow-x-auto rounded-xl border border-slate-200">
      <table className="w-full text-left border-collapse text-sm">
        <thead>
          <tr className="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
            <th className="py-3.5 px-4">Nama Siswa</th>
            <th className="py-3.5 px-4">Username</th>
            <th className="py-3.5 px-4 text-center">Percobaan Ke</th>
            <th className="py-3.5 px-4 text-center">Mode</th>
            <th className="py-3.5 px-4 text-center">Argumen Selesai</th>
            <th className="py-3.5 px-4 text-center">Total Percobaan</th>
            <th className="py-3.5 px-4 text-center">Status</th>
            <th className="py-3.5 px-4 text-right">Waktu Mulai</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100 bg-white">
          {results.map((r) => (
            <tr key={r.id} className="hover:bg-slate-50/70 transition-colors">
              <td className="py-3 px-4 font-semibold text-slate-900">{r.student_name}</td>
              <td className="py-3 px-4 text-slate-500">{r.username}</td>
              <td className="py-3 px-4 text-center font-medium">#{r.attempt_no}</td>
              <td className="py-3 px-4 text-center">
                <Badge variant="outline" className="capitalize text-[10px]">
                  {r.mode}
                </Badge>
              </td>
              <td className="py-3 px-4 text-center text-slate-700">
                {r.completed_arguments} / {r.total_arguments}
              </td>
              <td className="py-3 px-4 text-center text-slate-700 font-medium">
                {r.total_attempts} kali
              </td>
              <td className="py-3 px-4 text-center">
                {r.status === 'completed' ? (
                  <Badge variant="outline" className="bg-emerald-50 text-emerald-700 border-emerald-200 text-[10px]">
                    <CheckCircle2 className="w-3 h-3 mr-0.5" /> Selesai
                  </Badge>
                ) : (
                  <Badge variant="outline" className="bg-amber-50 text-amber-700 border-amber-200 text-[10px]">
                    <Clock className="w-3 h-3 mr-0.5" /> Berjalan
                  </Badge>
                )}
              </td>
              <td className="py-3 px-4 text-right text-xs text-slate-400">
                {new Date(r.started_at).toLocaleString('id-ID', {
                  day: 'numeric',
                  month: 'short',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
