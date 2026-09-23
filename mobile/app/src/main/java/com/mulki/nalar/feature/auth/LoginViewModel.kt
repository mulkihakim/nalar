package com.mulki.nalar.feature.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.mulki.nalar.core.auth.TokenManager
import com.mulki.nalar.core.network.ApiService
import com.mulki.nalar.core.network.dto.LoginRequest
import com.mulki.nalar.core.network.dto.UserDto
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import org.json.JSONObject

sealed interface LoginUiState {
    object Idle : LoginUiState
    object Loading : LoginUiState
    data class Success(val user: UserDto) : LoginUiState
    data class Error(val message: String) : LoginUiState
}

class LoginViewModel(
    private val apiService: ApiService,
    private val tokenManager: TokenManager
) : ViewModel() {

    private val _uiState = MutableStateFlow<LoginUiState>(LoginUiState.Idle)
    val uiState: StateFlow<LoginUiState> = _uiState.asStateFlow()

    fun login(username: String, password: String) {
        if (username.isBlank() || password.isBlank()) {
            _uiState.value = LoginUiState.Error("Username dan password wajib diisi")
            return
        }

        viewModelScope.launch {
            _uiState.value = LoginUiState.Loading
            try {
                val response = apiService.login(LoginRequest(username.trim(), password))
                if (response.isSuccessful && response.body() != null) {
                    val body = response.body()!!
                    tokenManager.saveAuth(body.token, body.user)
                    _uiState.value = LoginUiState.Success(body.user)
                } else {
                    val errorBody = response.errorBody()?.string()
                    val errorMsg = try {
                        JSONObject(errorBody ?: "").optString("error", "Login gagal")
                    } catch (e: Exception) {
                        "Login gagal (${response.code()})"
                    }
                    _uiState.value = LoginUiState.Error(errorMsg)
                }
            } catch (e: Exception) {
                _uiState.value = LoginUiState.Error(
                    e.localizedMessage ?: "Tidak dapat terhubung ke server. Periksa koneksi Anda."
                )
            }
        }
    }

    fun resetState() {
        _uiState.value = LoginUiState.Idle
    }

    class Factory(
        private val apiService: ApiService,
        private val tokenManager: TokenManager
    ) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            if (modelClass.isAssignableFrom(LoginViewModel::class.java)) {
                return LoginViewModel(apiService, tokenManager) as T
            }
            throw IllegalArgumentException("Unknown ViewModel class")
        }
    }
}
