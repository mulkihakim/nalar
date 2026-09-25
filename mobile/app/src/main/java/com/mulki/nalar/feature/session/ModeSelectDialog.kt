package com.mulki.nalar.feature.session

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.RadioButton
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.mulki.nalar.core.network.dto.ExamDto
import com.mulki.nalar.core.network.dto.SessionStatusResponse
import com.mulki.nalar.core.network.dto.formatModeLabel

private data class ModeOption(
    val key: String,
    val title: String,
    val description: String
)

private val MODE_OPTIONS = listOf(
    ModeOption(
        key = "standard",
        title = "Mode Standar",
        description = "Hanya menampilkan apakah kombinasi jawaban sudah tepat atau belum."
    ),
    ModeOption(
        key = "help",
        title = "Mode Bantuan",
        description = "Menampilkan petunjuk bagian Ground atau Warrant mana yang sudah tepat/salah."
    ),
    ModeOption(
        key = "social",
        title = "Mode Analitik Sosial",
        description = "Dilengkapi perbandingan pilihan opsi secara anonim dengan rekan lainnya."
    )
)

@Composable
fun ModeSelectDialog(
    exam: ExamDto,
    sessionStatus: SessionStatusResponse?,
    isCheckingStatus: Boolean,
    isStarting: Boolean,
    errorMessage: String?,
    onSelectMode: (String) -> Unit,
    onDismiss: () -> Unit
) {
    var selectedMode by remember { mutableStateOf("standard") }

    LaunchedEffect(sessionStatus) {
        if (sessionStatus?.has_active_session == true && !sessionStatus.mode.isNullOrBlank()) {
            selectedMode = sessionStatus.mode
        }
    }

    AlertDialog(
        onDismissRequest = { if (!isStarting) onDismiss() },
        containerColor = MaterialTheme.colorScheme.surface,
        title = {
            Column {
                Text(
                    text = "Pilih Mode Belajar",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = exam.title,
                    style = MaterialTheme.typography.labelMedium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
        },
        text = {
            Column(
                verticalArrangement = Arrangement.spacedBy(10.dp),
                modifier = Modifier.fillMaxWidth()
            ) {
                if (isCheckingStatus) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(8.dp),
                        modifier = Modifier.padding(vertical = 4.dp)
                    ) {
                        CircularProgressIndicator(modifier = Modifier.size(16.dp), strokeWidth = 2.dp)
                        Text(
                            text = "Memeriksa sesi berjalan...",
                            style = MaterialTheme.typography.labelMedium,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    }
                } else if (sessionStatus?.has_active_session == true) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clip(RoundedCornerShape(8.dp))
                            .background(MaterialTheme.colorScheme.tertiaryContainer)
                            .padding(10.dp)
                    ) {
                        Text(
                            text = "Anda memiliki sesi yang sedang berjalan dengan Mode ${formatModeLabel(sessionStatus.mode)}.",
                            style = MaterialTheme.typography.labelMedium,
                            color = MaterialTheme.colorScheme.onTertiaryContainer
                        )
                    }
                }

                if (!errorMessage.isNullOrBlank()) {
                    Surface(
                        color = MaterialTheme.colorScheme.errorContainer,
                        shape = RoundedCornerShape(8.dp),
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Text(
                            text = errorMessage,
                            color = MaterialTheme.colorScheme.onErrorContainer,
                            style = MaterialTheme.typography.labelMedium,
                            modifier = Modifier.padding(10.dp)
                        )
                    }
                }

                MODE_OPTIONS.forEach { option ->
                    val isSelected = selectedMode == option.key
                    Card(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clickable(enabled = !isStarting) { selectedMode = option.key },
                        shape = RoundedCornerShape(12.dp),
                        colors = CardDefaults.cardColors(
                            containerColor = if (isSelected) {
                                MaterialTheme.colorScheme.primaryContainer.copy(alpha = 0.5f)
                            } else {
                                MaterialTheme.colorScheme.surface
                            }
                        ),
                        border = BorderStroke(
                            width = if (isSelected) 1.5.dp else 1.dp,
                            color = if (isSelected) {
                                MaterialTheme.colorScheme.primary
                            } else {
                                MaterialTheme.colorScheme.outline
                            }
                        )
                    ) {
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(12.dp),
                            verticalAlignment = Alignment.Top
                        ) {
                            RadioButton(
                                selected = isSelected,
                                onClick = { if (!isStarting) selectedMode = option.key },
                                modifier = Modifier.size(20.dp)
                            )
                            Spacer(modifier = Modifier.width(10.dp))
                            Column {
                                Text(
                                    text = option.title,
                                    style = MaterialTheme.typography.bodyMedium,
                                    fontWeight = FontWeight.SemiBold,
                                    color = MaterialTheme.colorScheme.onSurface
                                )
                                Spacer(modifier = Modifier.height(2.dp))
                                Text(
                                    text = option.description,
                                    style = MaterialTheme.typography.labelSmall,
                                    color = MaterialTheme.colorScheme.onSurfaceVariant
                                )
                            }
                        }
                    }
                }

                // Pemberitahuan privasi sesuai 04-rules-design.md §A.1
                Text(
                    text = "Catatan: Setiap pilihan opsi dicatat untuk analisis belajar. Identitas Anda bersifat anonim bagi siswa lain dan hanya terlihat oleh pengajar.",
                    style = MaterialTheme.typography.labelSmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
        },
        confirmButton = {
            Button(
                onClick = { onSelectMode(selectedMode) },
                enabled = !isStarting && !isCheckingStatus,
                modifier = Modifier.heightIn(min = 48.dp)
            ) {
                if (isStarting) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(18.dp),
                        color = MaterialTheme.colorScheme.onPrimary,
                        strokeWidth = 2.dp
                    )
                } else {
                    val label = if (sessionStatus?.has_active_session == true) "Lanjutkan" else "Mulai"
                    Text(label, fontSize = 14.sp)
                }
            }
        },
        dismissButton = {
            OutlinedButton(
                onClick = onDismiss,
                enabled = !isStarting,
                modifier = Modifier.heightIn(min = 48.dp)
            ) {
                Text("Batal", fontSize = 14.sp)
            }
        }
    )
}
