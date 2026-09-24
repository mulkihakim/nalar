import type { StaffLogItem } from '../types'

interface LogsTableProps {
  logs: StaffLogItem[]
  isLoading?: boolean
}

export function LogsTable({ logs, isLoading }: LogsTableProps) {
  if (isLoading) {
    return <div className="py-12 text-center text-slate-400">Memuat log pengerjaan...</div>
  }

  if (logs.length === 0) {
    return (
      <div className="py-12 text-center text-slate-400 text-sm">
        Belum ada log percobaan pengerjaan pada ujian ini.
      </div>
    )
  }

  return (
    <div className="overflow-x-auto rounded-xl border border-slate-200">
      <table className="w-full text-left border-collapse text-sm">
        <thead>
          <tr className="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
            <th className="py-3.5 px-4">Waktu</th>
            <th className="py-3.5 px-4">Nama Siswa</th>
            <th className="py-3.5 px-4 text-center">Percobaan #</th>
            <th className="py-3.5 px-4 text-center">Slot</th>
            <th className="py-3.5 px-4">Opsi yang Di-drop</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100 bg-white">
          {logs.map((log) => {
            const isGround = log.slot === 'ground'
            return (
              <tr key={log.id} className="hover:bg-slate-50/70 transition-colors">
                <td className="py-3 px-4 text-xs text-slate-400 whitespace-nowrap">
                  {new Date(log.created_at).toLocaleString('id-ID', {
                    day: 'numeric',
                    month: 'short',
                    hour: '2-digit',
                    minute: '2-digit',
                    second: '2-digit',
                  })}
                </td>
                <td className="py-3 px-4 font-semibold text-slate-900 whitespace-nowrap">
                  {log.student_name}
                </td>
                <td className="py-3 px-4 text-center font-medium text-slate-600 whitespace-nowrap">
                  #{log.attempt_no}
                </td>
                <td className="py-3 px-4 text-center whitespace-nowrap">
                  <span
                    className={`text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider ${
                      isGround
                        ? 'bg-blue-50 text-blue-700 border border-blue-200'
                        : 'bg-purple-50 text-purple-700 border border-purple-200'
                    }`}
                  >
                    {log.slot}
                  </span>
                </td>
                <td className="py-3 px-4 text-slate-800 text-xs leading-relaxed max-w-md">
                  {log.option_text || `Opsi #${log.option_id}`}
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}
