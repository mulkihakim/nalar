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
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.mulki.nalar.core.auth.TokenManager
import com.mulki.nalar.core.network.NetworkModule
import com.mulki.nalar.feature.auth.LoginScreen
import com.mulki.nalar.feature.auth.LoginViewModel
import com.mulki.nalar.ui.theme.NalarTheme
import kotlinx.coroutines.launch

class MainActivity : ComponentActivity() {

    private lateinit var tokenManager: TokenManager

    private val loginViewModel: LoginViewModel by viewModels {
        val apiService = NetworkModule.provideApiService(tokenManager)
        LoginViewModel.Factory(apiService, tokenManager)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        tokenManager = TokenManager(applicationContext)

        enableEdgeToEdge()
        setContent {
            NalarTheme {
                val token by tokenManager.token.collectAsState(initial = null)
                val userName by tokenManager.userName.collectAsState(initial = null)
                val userRole by tokenManager.role.collectAsState(initial = null)
                val scope = rememberCoroutineScope()

                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    if (token.isNullOrBlank()) {
                        LoginScreen(
                            viewModel = loginViewModel,
                            onLoginSuccess = { _ ->
                                // Auth token tersimpan di DataStore
                            },
                            modifier = Modifier.padding(innerPadding)
                        )
                    } else {
                        AuthenticatedHome(
                            name = userName ?: "Pengguna",
                            role = userRole ?: "siswa",
                            onLogout = {
                                scope.launch {
                                    tokenManager.clearAuth()
                                }
                            },
                            modifier = Modifier.padding(innerPadding)
                        )
                    }
                }
            }
        }
    }
}

@Composable
fun AuthenticatedHome(
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
                containerColor = MaterialTheme.colorScheme.surfaceVariant
            )
        ) {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(24.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                Text(
                    text = "Selamat Datang!",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
                Text(
                    text = name,
                    fontSize = 18.sp,
                    fontWeight = FontWeight.SemiBold
                )
                Text(
                    text = "Peran: ${role.uppercase()}",
                    fontSize = 14.sp,
                    color = MaterialTheme.colorScheme.primary
                )

                Spacer(modifier = Modifier.height(16.dp))

                Button(
                    onClick = onLogout,
                    modifier = Modifier
                        .fillMaxWidth()
                        .heightIn(min = 48.dp) // Area sentuh minimal 48dp (05-rules.md §1)
                ) {
                    Text("Keluar (Logout)")
                }
            }
        }
    }
}