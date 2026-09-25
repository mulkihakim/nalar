package com.mulki.nalar.feature.examlist

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.mulki.nalar.core.network.ApiService
import com.mulki.nalar.core.network.dto.ExamDto
import com.mulki.nalar.core.network.dto.SessionHistoryItem
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import org.json.JSONObject

sealed interface ExamListUiState {
    object Loading : ExamListUiState
    data class Success(val exams: List<ExamDto>) : ExamListUiState
    data class Error(val message: String) : ExamListUiState
}

data class ExamHistoryModalState(
    val exam: ExamDto? = null,
    val isLoading: Boolean = false,
    val sessions: List<SessionHistoryItem> = emptyList(),
    val errorMessage: String? = null
)

class ExamListViewModel(
    private val apiService: ApiService
) : ViewModel() {

    private val _uiState = MutableStateFlow<ExamListUiState>(ExamListUiState.Loading)
    val uiState: StateFlow<ExamListUiState> = _uiState.asStateFlow()

    private val _historyModalState = MutableStateFlow(ExamHistoryModalState())
    val historyModalState: StateFlow<ExamHistoryModalState> = _historyModalState.asStateFlow()

    init {
        loadExams()
    }

    fun loadExams() {
        viewModelScope.launch {
            _uiState.value = ExamListUiState.Loading
            try {
                val response = apiService.getMyExams()
                if (response.isSuccessful && response.body() != null) {
                    _uiState.value = ExamListUiState.Success(response.body()!!)
                } else {
                    val errorMsg = try {
                        val errorBody = response.errorBody()?.string()
                        JSONObject(errorBody ?: "").optString("error", "Gagal memuat ujian")
                    } catch (e: Exception) {
                        "Gagal memuat ujian (${response.code()})"
                    }
                    _uiState.value = ExamListUiState.Error(errorMsg)
                }
            } catch (e: Exception) {
                _uiState.value = ExamListUiState.Error(
                    e.localizedMessage ?: "Tidak dapat terhubung ke server."
                )
            }
        }
    }

    fun openExamHistory(exam: ExamDto) {
        _historyModalState.value = ExamHistoryModalState(
            exam = exam,
            isLoading = true
        )
        viewModelScope.launch {
            try {
                val res = apiService.getMyExamSessions(exam.id)
                if (res.isSuccessful && res.body() != null) {
                    _historyModalState.value = _historyModalState.value.copy(
                        isLoading = false,
                        sessions = res.body()!!
                    )
                } else {
                    _historyModalState.value = _historyModalState.value.copy(
                        isLoading = false,
                        errorMessage = "Gagal memuat riwayat ujian"
                    )
                }
            } catch (e: Exception) {
                _historyModalState.value = _historyModalState.value.copy(
                    isLoading = false,
                    errorMessage = e.localizedMessage ?: "Gagal terhubung ke server"
                )
            }
        }
    }

    fun closeExamHistory() {
        _historyModalState.value = ExamHistoryModalState()
    }

    class Factory(private val apiService: ApiService) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            if (modelClass.isAssignableFrom(ExamListViewModel::class.java)) {
                return ExamListViewModel(apiService) as T
            }
            throw IllegalArgumentException("Unknown ViewModel class")
        }
    }
}
