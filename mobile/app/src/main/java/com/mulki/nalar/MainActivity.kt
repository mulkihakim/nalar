package com.mulki.nalar

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.viewModels
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.mulki.nalar.core.auth.TokenManager
import com.mulki.nalar.core.network.ApiService
import com.mulki.nalar.core.network.NetworkModule
import com.mulki.nalar.feature.auth.LoginScreen
import com.mulki.nalar.feature.auth.LoginViewModel
import com.mulki.nalar.feature.navigation.NalarApp
import com.mulki.nalar.ui.theme.NalarTheme
import kotlinx.coroutines.launch

class MainActivity : ComponentActivity() {

    private lateinit var tokenManager: TokenManager
    private lateinit var apiService: ApiService

    private val loginViewModel: LoginViewModel by viewModels {
        LoginViewModel.Factory(apiService, tokenManager)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        tokenManager = TokenManager(applicationContext)
        apiService = NetworkModule.provideApiService(tokenManager)

        enableEdgeToEdge()
        setContent {
            NalarTheme {
                val token by tokenManager.token.collectAsState(initial = null)
                val userName by tokenManager.userName.collectAsState(initial = null)
                val username by tokenManager.username.collectAsState(initial = null)
                val userRole by tokenManager.role.collectAsState(initial = null)
                val scope = rememberCoroutineScope()

                if (token.isNullOrBlank()) {
                    Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                        LoginScreen(
                            viewModel = loginViewModel,
                            onLoginSuccess = { _ ->
                                // Token otomatis tersimpan di DataStore
                            },
                            modifier = Modifier.padding(innerPadding)
                        )
                    }
                } else if (!userRole.isNullOrBlank() && userRole != "siswa") {
                    Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                        NonStudentNoticeScreen(
                            name = userName ?: "Pengguna",
                            role = userRole ?: "",
                            onLogout = {
                                scope.launch {
                                    tokenManager.clearAuth()
                                    loginViewModel.resetState()
                                }
                            },
                            modifier = Modifier.padding(innerPadding)
                        )
                    }
                } else {
                    NalarApp(
                        apiService = apiService,
                        userName = userName ?: "Siswa",
                        username = username ?: "siswa",
                        userRole = userRole ?: "siswa",
                        onLogout = {
                            scope.launch {
                                tokenManager.clearAuth()
                                loginViewModel.resetState()
                            }
                        }
                    )
                }
            }
        }
    }
}

@Composable
private fun NonStudentNoticeScreen(
    name: String,
    role: String,
    onLogout: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Card(
            modifier = Modifier.fillMaxWidth(),
            colors = CardDefaults.cardColors(
                containerColor = MaterialTheme.colorScheme.surface
            ),
            elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
        ) {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(24.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                Text(
                    text = "Halo, $name",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
                Text(
                    text = "Akun Anda memiliki peran ${role.uppercase()}. Aplikasi mobile Nalar dikhususkan untuk Siswa mengerjakan latihan. Silakan gunakan Nalar versi Web untuk mengakses panel pengelolaan Admin/Asesor.",
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
                Spacer(modifier = Modifier.height(8.dp))
                Button(
                    onClick = onLogout,
                    modifier = Modifier
                        .fillMaxWidth()
                        .heightIn(min = 48.dp)
                ) {
                    Text("Keluar & Ganti Akun Siswa")
                }
            }
        }
    }
}