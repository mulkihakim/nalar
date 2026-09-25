package com.mulki.nalar.feature.result

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.mulki.nalar.core.network.ApiService
import com.mulki.nalar.core.network.dto.AnalysisResponse
import com.mulki.nalar.core.network.dto.ExamDto
import com.mulki.nalar.core.network.dto.SessionDetail
import com.mulki.nalar.core.network.dto.SessionHistoryItem
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import org.json.JSONObject

data class ExamHistoryItem(
    val exam: ExamDto,
    val sessions: List<SessionHistoryItem> = emptyList(),
    val isLoadingSessions: Boolean = false,
    val isExpanded: Boolean = false
)

data class ResultUiState(
    val isLoading: Boolean = false,
    val exams: List<ExamHistoryItem> = emptyList(),
    val errorMessage: String? = null,
    val selectedSessionForAnalysis: SessionDetail? = null,
    val analysisData: AnalysisResponse? = null,
    val isLoadingAnalysis: Boolean = false
)

class ResultViewModel(
    private val apiService: ApiService
) : ViewModel() {

    private val _uiState = MutableStateFlow(ResultUiState())
    val uiState: StateFlow<ResultUiState> = _uiState.asStateFlow()

    init {
        loadExams()
    }

    fun loadExams() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, errorMessage = null)
            try {
                val res = apiService.getMyExams()
                if (res.isSuccessful && res.body() != null) {
                    val examItems = res.body()!!.map { ExamHistoryItem(exam = it) }
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        exams = examItems
                    )
                    // Auto-load sessions for the first exam if available
                    if (examItems.isNotEmpty()) {
                        toggleExpandExam(examItems[0].exam.id)
                    }
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        errorMessage = "Gagal memuat ujian"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    errorMessage = e.localizedMessage ?: "Tidak dapat terhubung ke server"
                )
            }
        }
    }

    fun toggleExpandExam(examId: Long) {
        val currentExams = _uiState.value.exams
        val targetIndex = currentExams.indexOfFirst { it.exam.id == examId }
        if (targetIndex < 0) return

        val item = currentExams[targetIndex]
        val newExpanded = !item.isExpanded

        val updatedExams = currentExams.toMutableList()
        updatedExams[targetIndex] = item.copy(isExpanded = newExpanded)
        _uiState.value = _uiState.value.copy(exams = updatedExams)

        if (newExpanded && item.sessions.isEmpty()) {
            loadSessionsForExam(examId)
        }
    }

    fun loadSessionsForExam(examId: Long) {
        val currentExams = _uiState.value.exams
        val targetIndex = currentExams.indexOfFirst { it.exam.id == examId }
        if (targetIndex < 0) return

        val updated = currentExams.toMutableList()
        updated[targetIndex] = updated[targetIndex].copy(isLoadingSessions = true)
        _uiState.value = _uiState.value.copy(exams = updated)

        viewModelScope.launch {
            try {
                val res = apiService.getMyExamSessions(examId)
                val sessions = if (res.isSuccessful && res.body() != null) res.body()!! else emptyList()

                val finalExams = _uiState.value.exams.toMutableList()
                val idx = finalExams.indexOfFirst { it.exam.id == examId }
                if (idx >= 0) {
                    finalExams[idx] = finalExams[idx].copy(
                        sessions = sessions,
                        isLoadingSessions = false
                    )
                    _uiState.value = _uiState.value.copy(exams = finalExams)
                }
            } catch (e: Exception) {
                val finalExams = _uiState.value.exams.toMutableList()
                val idx = finalExams.indexOfFirst { it.exam.id == examId }
                if (idx >= 0) {
                    finalExams[idx] = finalExams[idx].copy(isLoadingSessions = false)
                    _uiState.value = _uiState.value.copy(exams = finalExams)
                }
            }
        }
    }

    fun loadAnalysisForSession(sessionId: Long) {
        _uiState.value = _uiState.value.copy(isLoadingAnalysis = true, analysisData = null)
        viewModelScope.launch {
            try {
                val res = apiService.getAnalysis(sessionId)
                if (res.isSuccessful && res.body() != null) {
                    _uiState.value = _uiState.value.copy(
                        analysisData = res.body(),
                        isLoadingAnalysis = false
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingAnalysis = false)
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingAnalysis = false)
            }
        }
    }

    fun clearAnalysis() {
        _uiState.value = _uiState.value.copy(analysisData = null)
    }

    class Factory(private val apiService: ApiService) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            if (modelClass.isAssignableFrom(ResultViewModel::class.java)) {
                return ResultViewModel(apiService) as T
            }
            throw IllegalArgumentException("Unknown ViewModel class")
        }
    }
}
