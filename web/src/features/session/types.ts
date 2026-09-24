export type OptionType = 'ground' | 'warrant'

export type SessionMode = 'standard' | 'help' | 'social'
export type SessionStatus = 'in_progress' | 'completed'

export interface SessionStatusResponse {
  has_active_session: boolean
  session_id?: number
  mode?: SessionMode
}

export interface OptionItem {
  id: number
  argument_id: number
  type: OptionType
  text: string
  is_correct?: boolean
}

export interface ArgumentItem {
  id: number
  claim_text: string
  order_no: number
  completed: boolean
  options: OptionItem[]
}

export interface SessionDetail {
  id: number
  exam_id: number
  exam_title: string
  material_title: string
  material_content: string
  student_id: number
  attempt_no: number
  mode: SessionMode
  status: SessionStatus
  current_argument_id?: number
  arguments: ArgumentItem[]
  completed_arguments: number
  total_arguments: number
}

export interface StartSessionPayload {
  mode: SessionMode
  confirm_mode_change?: boolean
}

export interface DropPayload {
  option_id: number
  slot: OptionType
}

export interface ConfirmPayload {
  ground_option_id: number
  warrant_option_id: number
}

export interface ConfirmResult {
  is_correct: boolean
  completed: boolean
  is_session_completed: boolean
  message: string
  ground_correct?: boolean
  warrant_correct?: boolean
}

export interface MonitoringOptionStat {
  option_id: number
  type: OptionType
  text: string
  x: number
  y: string
}

export interface MonitoringAnalyticsResponse {
  argument_id: number
  claim_text?: string
  has_enough_peers: boolean
  min_peers: number
  total_peers: number
  options: MonitoringOptionStat[]
}

export interface StudentChooserInfo {
  student_id: number
  student_name: string
  username: string
  attempts: number
  last_attempt: string
}

export interface AnalysisOptionStat {
  option_id: number
  type: OptionType
  text: string
  total_attempts: number
  unique_students: number
  ratio: string
  students?: StudentChooserInfo[]
}

export interface AnalysisArgumentStat {
  argument_id: number
  claim_text: string
  order_no: number
  options: AnalysisOptionStat[]
}

export interface AnalysisAnalyticsResponse {
  exam_id: number
  has_enough_peers: boolean
  min_peers: number
  total_peers: number
  arguments: AnalysisArgumentStat[]
}

export interface ArgumentHistoryDetail {
  argument_id: number
  claim_text: string
  attempts: number
  is_clean: boolean
}

export interface SessionHistoryItem {
  id: number
  attempt_no: number
  mode: SessionMode
  status: SessionStatus
  started_at: string
  completed_at?: string
  total_attempts: number
  completed_arguments: number
  total_arguments: number
  argument_details: ArgumentHistoryDetail[]
}

// Staff Types
export interface StaffResultSummary {
  id: number
  student_id: number
  student_name: string
  username: string
  attempt_no: number
  mode: SessionMode
  status: SessionStatus
  started_at: string
  completed_at?: string
  total_attempts: number
  completed_arguments: number
  total_arguments: number
}

export interface StaffLogItem {
  id: number
  session_id: number
  attempt_no: number
  student_id: number
  student_name: string
  argument_id: number
  option_id: number
  option_text: string
  slot: OptionType
  created_at: string
}
