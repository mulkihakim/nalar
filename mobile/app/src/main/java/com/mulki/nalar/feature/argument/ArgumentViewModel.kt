package com.mulki.nalar.feature.argument

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.mulki.nalar.core.network.ApiService
import com.mulki.nalar.core.network.dto.ArgumentItem
import com.mulki.nalar.core.network.dto.ConfirmRequest
import com.mulki.nalar.core.network.dto.ConfirmResult
import com.mulki.nalar.core.network.dto.DropRequest
import com.mulki.nalar.core.network.dto.MonitoringResponse
import com.mulki.nalar.core.network.dto.OptionItem
import com.mulki.nalar.core.network.dto.SessionDetail
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import org.json.JSONObject

data class ArgumentUiState(
    val session: SessionDetail,
    val currentArgumentIndex: Int = 0,
    val selectedGroundOption: OptionItem? = null,
    val selectedWarrantOption: OptionItem? = null,
    val activeBottomSheetType: String? = null, // "ground" | "warrant" | null
    val isConfirming: Boolean = false,
    val confirmResult: ConfirmResult? = null,
    val isShowingMonitoring: Boolean = false,
    val monitoringData: MonitoringResponse? = null,
    val isLoadingMonitoring: Boolean = false,
    val isSessionCompleted: Boolean = false,
    val errorMessage: String? = null
)

class ArgumentViewModel(
    initialSession: SessionDetail,
    private val apiService: ApiService
) : ViewModel() {

    private val _uiState: MutableStateFlow<ArgumentUiState>

    init {
        // Cari argumen yang belum selesai pertama kali
        val firstIncompleteIndex = initialSession.arguments.indexOfFirst { !it.completed }
        val startIndex = if (firstIncompleteIndex >= 0) firstIncompleteIndex else 0

        _uiState = MutableStateFlow(
            ArgumentUiState(
                session = initialSession,
                currentArgumentIndex = startIndex,
                isSessionCompleted = initialSession.status == "completed" || (firstIncompleteIndex == -1 && initialSession.arguments.isNotEmpty())
            )
        )
    }

    val uiState: StateFlow<ArgumentUiState> = _uiState.asStateFlow()

    fun openBottomSheet(type: String) {
        _uiState.value = _uiState.value.copy(activeBottomSheetType = type)
    }

    fun dismissBottomSheet() {
        _uiState.value = _uiState.value.copy(activeBottomSheetType = null)
    }

    fun selectOption(option: OptionItem) {
        val currentState = _uiState.value
        val currentArg = getCurrentArgument() ?: return
        val slotType = currentState.activeBottomSheetType ?: option.type

        if (slotType == "ground") {
            _uiState.value = currentState.copy(
                selectedGroundOption = option,
                activeBottomSheetType = null,
                errorMessage = null
            )
        } else if (slotType == "warrant") {
            _uiState.value = currentState.copy(
                selectedWarrantOption = option,
                activeBottomSheetType = null,
                errorMessage = null
            )
        }

        // Catat drop ke backend (asinkron, F-09)
        viewModelScope.launch {
            try {
                apiService.recordDrop(
                    sessionId = currentState.session.id,
                    argumentId = currentArg.id,
                    request = DropRequest(option_id = option.id, slot = slotType)
                )
            } catch (e: Exception) {
                // Drop error dicatat silent atau log
            }
        }
    }

    fun confirmAnswers() {
        val currentState = _uiState.value
        val ground = currentState.selectedGroundOption
        val warrant = currentState.selectedWarrantOption
        val currentArg = getCurrentArgument() ?: return

        if (ground == null || warrant == null) {
            _uiState.value = currentState.copy(errorMessage = "Pilih kedua slot Ground dan Warrant terlebih dahulu")
            return
        }

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isConfirming = true, errorMessage = null)
            try {
                val response = apiService.confirmArgument(
                    sessionId = currentState.session.id,
                    argumentId = currentArg.id,
                    request = ConfirmRequest(
                        ground_option_id = ground.id,
                        warrant_option_id = warrant.id
                    )
                )

                if (response.isSuccessful && response.body() != null) {
                    val result = response.body()!!

                    // Tandai argumen saat ini selesai di lokal jika benar
                    val updatedArguments = currentState.session.arguments.map {
                        if (it.id == currentArg.id && result.is_correct) {
                            it.copy(completed = true)
                        } else it
                    }

                    _uiState.value = _uiState.value.copy(
                        isConfirming = false,
                        confirmResult = result,
                        session = currentState.session.copy(
                            arguments = updatedArguments,
                            completed_arguments = if (result.is_correct) currentState.session.completed_arguments + 1 else currentState.session.completed_arguments
                        ),
                        isSessionCompleted = result.is_session_completed
                    )
                } else {
                    val errorMsg = try {
                        val errorBody = response.errorBody()?.string()
                        JSONObject(errorBody ?: "").optString("error", "Konfirmasi gagal")
                    } catch (e: Exception) {
                        "Konfirmasi gagal (${response.code()})"
                    }
                    _uiState.value = _uiState.value.copy(
                        isConfirming = false,
                        errorMessage = errorMsg
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isConfirming = false,
                    errorMessage = e.localizedMessage ?: "Tidak dapat terhubung ke server"
                )
            }
        }
    }

    fun handleContinueAfterConfirm() {
        val currentState = _uiState.value
        val result = currentState.confirmResult ?: return

        if (result.is_correct) {
            if (currentState.session.mode == "social") {
                // Di mode analitik sosial, selalu tampilkan monitoring view (termasuk pada soal terakhir)
                loadMonitoringView()
            } else if (result.is_session_completed) {
                _uiState.value = currentState.copy(
                    confirmResult = null,
                    isSessionCompleted = true
                )
            } else {
                advanceToNextArgument()
            }
        }
    }

    private fun loadMonitoringView() {
        val currentArg = getCurrentArgument() ?: return
        val sessionId = _uiState.value.session.id

        _uiState.value = _uiState.value.copy(
            confirmResult = null,
            isShowingMonitoring = true,
            isLoadingMonitoring = true
        )

        viewModelScope.launch {
            try {
                val res = apiService.getMonitoring(sessionId, currentArg.id)
                if (res.isSuccessful && res.body() != null) {
                    _uiState.value = _uiState.value.copy(
                        monitoringData = res.body(),
                        isLoadingMonitoring = false
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingMonitoring = false)
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingMonitoring = false)
            }
        }
    }

    fun continueFromMonitoring() {
        val currentState = _uiState.value
        val isAllCompleted = currentState.session.arguments.all { it.completed }
        val isLastIndex = currentState.currentArgumentIndex >= currentState.session.arguments.size - 1

        if (currentState.isSessionCompleted || isAllCompleted || isLastIndex) {
            _uiState.value = currentState.copy(
                isShowingMonitoring = false,
                isSessionCompleted = true
            )
        } else {
            _uiState.value = currentState.copy(isShowingMonitoring = false)
            advanceToNextArgument()
        }
    }

    private fun advanceToNextArgument() {
        val currentState = _uiState.value
        val nextIndex = currentState.currentArgumentIndex + 1

        if (nextIndex < currentState.session.arguments.size) {
            _uiState.value = currentState.copy(
                currentArgumentIndex = nextIndex,
                selectedGroundOption = null,
                selectedWarrantOption = null,
                confirmResult = null,
                errorMessage = null
            )
        } else {
            _uiState.value = currentState.copy(
                isSessionCompleted = true,
                confirmResult = null
            )
        }
    }

    fun tryAgain() {
        _uiState.value = _uiState.value.copy(confirmResult = null)
    }

    fun getCurrentArgument(): ArgumentItem? {
        val state = _uiState.value
        return state.session.arguments.getOrNull(state.currentArgumentIndex)
    }

    class Factory(
        private val session: SessionDetail,
        private val apiService: ApiService
    ) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            if (modelClass.isAssignableFrom(ArgumentViewModel::class.java)) {
                return ArgumentViewModel(session, apiService) as T
            }
            throw IllegalArgumentException("Unknown ViewModel class")
        }
    }
}
