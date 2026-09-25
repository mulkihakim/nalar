package com.mulki.nalar.feature.session

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.mulki.nalar.core.network.ApiService
import com.mulki.nalar.core.network.dto.ExamDto
import com.mulki.nalar.core.network.dto.SessionDetail
import com.mulki.nalar.core.network.dto.SessionStatusResponse
import com.mulki.nalar.core.network.dto.StartSessionRequest
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import org.json.JSONObject

data class SessionDialogState(
    val selectedExam: ExamDto? = null,
    val sessionStatus: SessionStatusResponse? = null,
    val isCheckingStatus: Boolean = false,
    val isStarting: Boolean = false,
    val conflictCurrentMode: String? = null,
    val conflictTargetMode: String? = null,
    val errorMessage: String? = null,
    val startedSession: SessionDetail? = null
)

class SessionViewModel(
    private val apiService: ApiService
) : ViewModel() {

    private val _state = MutableStateFlow(SessionDialogState())
    val state: StateFlow<SessionDialogState> = _state.asStateFlow()

    fun openModeSelector(exam: ExamDto) {
        _state.value = SessionDialogState(
            selectedExam = exam,
            isCheckingStatus = true
        )
        viewModelScope.launch {
            try {
                val res = apiService.getSessionStatus(exam.id)
                if (res.isSuccessful && res.body() != null) {
                    _state.value = _state.value.copy(
                        sessionStatus = res.body(),
                        isCheckingStatus = false
                    )
                } else {
                    _state.value = _state.value.copy(isCheckingStatus = false)
                }
            } catch (e: Exception) {
                _state.value = _state.value.copy(isCheckingStatus = false)
            }
        }
    }

    fun selectMode(mode: String, confirmModeChange: Boolean = false) {
        val exam = _state.value.selectedExam ?: return
        viewModelScope.launch {
            _state.value = _state.value.copy(
                isStarting = true,
                errorMessage = null
            )
            try {
                val response = apiService.startSession(
                    examId = exam.id,
                    request = StartSessionRequest(
                        mode = mode,
                        confirm_mode_change = confirmModeChange
                    )
                )

                if (response.isSuccessful && response.body() != null) {
                    _state.value = _state.value.copy(
                        isStarting = false,
                        conflictCurrentMode = null,
                        conflictTargetMode = null,
                        startedSession = response.body()
                    )
                } else if (response.code() == 409) {
                    val errorBody = response.errorBody()?.string()
                    val currentMode = try {
                        JSONObject(errorBody ?: "").optString("current_mode", "standard")
                    } catch (e: Exception) {
                        "standard"
                    }
                    _state.value = _state.value.copy(
                        isStarting = false,
                        conflictCurrentMode = currentMode,
                        conflictTargetMode = mode
                    )
                } else {
                    val errorBody = response.errorBody()?.string()
                    val errorMsg = try {
                        JSONObject(errorBody ?: "").optString("error", "Gagal memulai sesi")
                    } catch (e: Exception) {
                        "Gagal memulai sesi (${response.code()})"
                    }
                    _state.value = _state.value.copy(
                        isStarting = false,
                        errorMessage = errorMsg
                    )
                }
            } catch (e: Exception) {
                _state.value = _state.value.copy(
                    isStarting = false,
                    errorMessage = e.localizedMessage ?: "Gagal terhubung ke server"
                )
            }
        }
    }

    fun confirmModeChange() {
        val targetMode = _state.value.conflictTargetMode ?: return
        selectMode(mode = targetMode, confirmModeChange = true)
    }

    fun cancelModeConflict() {
        _state.value = _state.value.copy(
            conflictCurrentMode = null,
            conflictTargetMode = null
        )
    }

    fun clearStartedSession() {
        _state.value = SessionDialogState()
    }

    fun dismissDialog() {
        _state.value = SessionDialogState()
    }

    class Factory(private val apiService: ApiService) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            if (modelClass.isAssignableFrom(SessionViewModel::class.java)) {
                return SessionViewModel(apiService) as T
            }
            throw IllegalArgumentException("Unknown ViewModel class")
        }
    }
}
