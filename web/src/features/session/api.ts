import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient, getStoredToken, removeStoredToken } from '@/lib/api-client'
import type { Exam } from '@/types'
import type {
  SessionDetail,
  SessionStatusResponse,
  StartSessionPayload,
  DropPayload,
  ConfirmPayload,
  ConfirmResult,
  MonitoringAnalyticsResponse,
  AnalysisAnalyticsResponse,
  SessionHistoryItem,
  StaffResultSummary,
  StaffLogItem,
  SessionMode,
} from './types'

export class ModeConflictApiError extends Error {
  currentMode?: SessionMode
  constructor(message: string, currentMode?: SessionMode) {
    super(message)
    this.name = 'ModeConflictApiError'
    this.currentMode = currentMode
  }
}

// 1. Siswa - Daftar ujian
export function useMyExams() {
  return useQuery({
    queryKey: ['my-exams'],
    queryFn: () => apiClient<Exam[]>('/my/exams'),
  })
}

// 2. Siswa - Cek status sesi berjalan
export function useSessionStatus(examId: number) {
  return useQuery({
    queryKey: ['exam-session-status', examId],
    queryFn: () => apiClient<SessionStatusResponse>(`/exams/${examId}/session-status`),
    enabled: !!examId,
  })
}

// 3. Siswa - Mulai / lanjut sesi (handle 409 conflict dengan custom fetch agar dapat current_mode)
export function useStartSession() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ examId, payload }: { examId: number; payload: StartSessionPayload }) => {
      const token = getStoredToken()
      const res = await fetch(`/api/v1/exams/${examId}/start`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify(payload),
      })

      if (!res.ok) {
        let errData: any = {}
        try {
          errData = await res.json()
        } catch {}
        if (res.status === 409) {
          throw new ModeConflictApiError(
            errData.error || 'Mode belajar berbeda dengan sesi berjalan',
            errData.current_mode
          )
        }
        if (res.status === 401) {
          removeStoredToken()
        }
        throw new Error(errData.error || `HTTP error ${res.status}`)
      }

      return (await res.json()) as SessionDetail
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['exam-session-status', data.exam_id] })
      queryClient.invalidateQueries({ queryKey: ['session-details', data.id] })
      queryClient.invalidateQueries({ queryKey: ['my-exam-sessions', data.exam_id] })
    },
  })
}

// 4. Siswa - Detail sesi aktif
export function useSessionDetails(sessionId: number) {
  return useQuery({
    queryKey: ['session-details', sessionId],
    queryFn: () => apiClient<SessionDetail>(`/sessions/${sessionId}`),
    enabled: !!sessionId,
  })
}

// 5. Siswa - Catat drop opsi
export function useRecordDrop() {
  return useMutation({
    mutationFn: ({
      sessionId,
      argumentId,
      payload,
    }: {
      sessionId: number
      argumentId: number
      payload: DropPayload
    }) =>
      apiClient<{ message: string }>(`/sessions/${sessionId}/arguments/${argumentId}/drops`, {
        method: 'POST',
        body: JSON.stringify(payload),
      }),
  })
}

// 6. Siswa - Konfirmasi jawaban argumen
export function useConfirmArgument() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({
      sessionId,
      argumentId,
      payload,
    }: {
      sessionId: number
      argumentId: number
      payload: ConfirmPayload
    }) =>
      apiClient<ConfirmResult>(`/sessions/${sessionId}/arguments/${argumentId}/confirm`, {
        method: 'POST',
        body: JSON.stringify(payload),
      }),
    onSuccess: (_, vars) => {
      queryClient.invalidateQueries({
        queryKey: ['monitoring-analytics', vars.sessionId, vars.argumentId],
      })
      queryClient.invalidateQueries({ queryKey: ['analysis-analytics', vars.sessionId] })
      queryClient.invalidateQueries({ queryKey: ['session-details', vars.sessionId] })
      queryClient.invalidateQueries({ queryKey: ['exam-session-status'] })
      queryClient.invalidateQueries({ queryKey: ['my-exams'] })
    },
  })
}

// 7. Siswa - Monitoring analitik per argumen
export function useMonitoringAnalytics(sessionId: number, argumentId?: number) {
  return useQuery({
    queryKey: ['monitoring-analytics', sessionId, argumentId],
    queryFn: () =>
      apiClient<MonitoringAnalyticsResponse>(`/sessions/${sessionId}/monitoring/${argumentId}`),
    enabled: !!sessionId && !!argumentId,
  })
}

// 8. Siswa - Analysis analitik keseluruhan
export function useAnalysisAnalytics(sessionId: number, enabled = true) {
  return useQuery({
    queryKey: ['analysis-analytics', sessionId],
    queryFn: () => apiClient<AnalysisAnalyticsResponse>(`/sessions/${sessionId}/analysis`),
    enabled: !!sessionId && enabled,
  })
}

// 9. Siswa - Riwayat seluruh sesi pada ujian ini
export function useMyExamSessions(examId: number) {
  return useQuery({
    queryKey: ['my-exam-sessions', examId],
    queryFn: () => apiClient<SessionHistoryItem[]>(`/my/exams/${examId}/sessions`),
    enabled: !!examId,
  })
}

// 10. Staff - Hasil semua siswa pada suatu ujian
export function useExamResults(examId: number) {
  return useQuery({
    queryKey: ['staff-exam-results', examId],
    queryFn: () => apiClient<StaffResultSummary[]>(`/exams/${examId}/results`),
    enabled: !!examId,
  })
}

// 11. Staff - Log detail pengerjaan
export function useExamLogs(examId: number) {
  return useQuery({
    queryKey: ['staff-exam-logs', examId],
    queryFn: () => apiClient<StaffLogItem[]>(`/exams/${examId}/logs`),
    enabled: !!examId,
  })
}

// 12. Staff - Analitik sosial dari sudut pandang staf
export function useStaffAnalysis(examId: number) {
  return useQuery({
    queryKey: ['staff-exam-analysis', examId],
    queryFn: () => apiClient<AnalysisAnalyticsResponse>(`/exams/${examId}/analysis`),
    enabled: !!examId,
  })
}
