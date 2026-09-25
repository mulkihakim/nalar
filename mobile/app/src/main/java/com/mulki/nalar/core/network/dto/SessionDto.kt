package com.mulki.nalar.core.network.dto

data class SessionStatusResponse(
    val has_active_session: Boolean,
    val session_id: Long? = null,
    val mode: String? = null
)

data class StartSessionRequest(
    val mode: String,
    val confirm_mode_change: Boolean = false
)

data class OptionItem(
    val id: Long,
    val argument_id: Long,
    val type: String, // "ground" | "warrant"
    val text: String,
    val is_correct: Boolean? = null
)

data class ArgumentItem(
    val id: Long,
    val claim_text: String,
    val order_no: Int,
    val completed: Boolean,
    val options: List<OptionItem>
)

data class SessionDetail(
    val id: Long,
    val exam_id: Long,
    val exam_title: String,
    val material_title: String,
    val material_content: String,
    val student_id: Long,
    val attempt_no: Int,
    val mode: String, // "standard" | "help" | "social"
    val status: String, // "in_progress" | "completed"
    val current_argument_id: Long? = null,
    val arguments: List<ArgumentItem>,
    val completed_arguments: Int,
    val total_arguments: numberTypeAlias
)

private typealias numberTypeAlias = Int

data class DropRequest(
    val option_id: Long,
    val slot: String // "ground" | "warrant"
)

data class ConfirmRequest(
    val ground_option_id: Long,
    val warrant_option_id: Long
)

data class ConfirmResult(
    val is_correct: Boolean,
    val completed: Boolean,
    val is_session_completed: Boolean,
    val message: String,
    val ground_correct: Boolean? = null,
    val warrant_correct: Boolean? = null
)

data class MonitoringOptionStat(
    val option_id: Long,
    val type: String,
    val text: String,
    val x: Int,
    val y: String
)

data class MonitoringResponse(
    val argument_id: Long,
    val claim_text: String? = null,
    val has_enough_peers: Boolean,
    val min_peers: Int,
    val total_peers: Int,
    val options: List<MonitoringOptionStat>
)

data class AnalysisOptionStat(
    val option_id: Long,
    val type: String,
    val text: String,
    val total_attempts: Int,
    val unique_students: Int,
    val ratio: String
)

data class AnalysisArgumentStat(
    val argument_id: Long,
    val claim_text: String,
    val order_no: Int,
    val options: List<AnalysisOptionStat>
)

data class AnalysisResponse(
    val exam_id: Long,
    val has_enough_peers: Boolean,
    val min_peers: Int,
    val total_peers: Int,
    val arguments: List<AnalysisArgumentStat>
)

data class ArgumentHistoryDetail(
    val argument_id: Long,
    val claim_text: String,
    val attempts: Int,
    val is_clean: Boolean
)

data class SessionHistoryItem(
    val id: Long,
    val attempt_no: Int,
    val mode: String,
    val status: String,
    val started_at: String,
    val completed_at: String? = null,
    val total_attempts: Int,
    val completed_arguments: Int,
    val total_arguments: Int,
    val argument_details: List<ArgumentHistoryDetail>? = emptyList()
)

fun formatModeLabel(mode: String?): String {
    return when (mode) {
        "standard" -> "Standar"
        "help" -> "Bantuan"
        "social" -> "Analitik Sosial"
        else -> mode ?: "Standar"
    }
}
