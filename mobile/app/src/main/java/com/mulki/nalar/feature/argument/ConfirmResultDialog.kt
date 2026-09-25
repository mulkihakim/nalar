package com.mulki.nalar.feature.argument

import androidx.compose.foundation.background
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
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Close
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.mulki.nalar.core.network.dto.ConfirmResult

@Composable
fun ConfirmResultDialog(
    result: ConfirmResult,
    mode: String,
    onContinue: () -> Unit,
    onTryAgain: () -> Unit
) {
    AlertDialog(
        onDismissRequest = {
            if (result.is_correct) onContinue() else onTryAgain()
        },
        containerColor = MaterialTheme.colorScheme.surface,
        title = {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                if (result.is_correct) {
                    Icon(
                        imageVector = Icons.Default.CheckCircle,
                        contentDescription = null,
                        tint = MaterialTheme.colorScheme.secondary,
                        modifier = Modifier.size(28.dp)
                    )
                    Text(
                        text = "Jawaban Tepat!",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.onSurface
                    )
                } else {
                    Icon(
                        imageVector = Icons.Default.Close,
                        contentDescription = null,
                        tint = MaterialTheme.colorScheme.error,
                        modifier = Modifier.size(28.dp)
                    )
                    Text(
                        text = "Belum Tepat",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.onSurface
                    )
                }
            }
        },
        text = {
            Column(
                verticalArrangement = Arrangement.spacedBy(10.dp),
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(
                    text = result.message,
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurface
                )

                // Mode bantuan: tampilkan mana yang benar dan mana yang salah (04-rules-design.md §A.2 A3)
                if (!result.is_correct && mode == "help") {
                    Column(
                        verticalArrangement = Arrangement.spacedBy(6.dp),
                        modifier = Modifier
                            .fillMaxWidth()
                            .clip(RoundedCornerShape(8.dp))
                            .background(MaterialTheme.colorScheme.surfaceVariant)
                            .padding(12.dp)
                    ) {
                        Text(
                            text = "Petunjuk Mode Bantuan:",
                            style = MaterialTheme.typography.labelSmall,
                            fontWeight = FontWeight.Bold,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )

                        // Ground evaluation
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text(
                                text = "Slot Ground: ",
                                style = MaterialTheme.typography.labelMedium,
                                color = MaterialTheme.colorScheme.onSurface
                            )
                            if (result.ground_correct == true) {
                                Text(
                                    text = "Sudah Tepat ✓",
                                    style = MaterialTheme.typography.labelMedium,
                                    fontWeight = FontWeight.Bold,
                                    color = MaterialTheme.colorScheme.secondary
                                )
                            } else {
                                Text(
                                    text = "Belum Tepat ✗",
                                    style = MaterialTheme.typography.labelMedium,
                                    fontWeight = FontWeight.Bold,
                                    color = MaterialTheme.colorScheme.error
                                )
                            }
                        }

                        // Warrant evaluation
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text(
                                text = "Slot Warrant: ",
                                style = MaterialTheme.typography.labelMedium,
                                color = MaterialTheme.colorScheme.onSurface
                            )
                            if (result.warrant_correct == true) {
                                Text(
                                    text = "Sudah Tepat ✓",
                                    style = MaterialTheme.typography.labelMedium,
                                    fontWeight = FontWeight.Bold,
                                    color = MaterialTheme.colorScheme.secondary
                                )
                            } else {
                                Text(
                                    text = "Belum Tepat ✗",
                                    style = MaterialTheme.typography.labelMedium,
                                    fontWeight = FontWeight.Bold,
                                    color = MaterialTheme.colorScheme.error
                                )
                            }
                        }
                    }
                }
            }
        },
        confirmButton = {
            if (result.is_correct) {
                Button(
                    onClick = onContinue,
                    modifier = Modifier.heightIn(min = 48.dp)
                ) {
                    val buttonLabel = when {
                        mode == "social" -> "Lihat Monitoring Teman"
                        result.is_session_completed -> "Lihat Ringkasan Hasil"
                        else -> "Lanjut ke Argumen Berikutnya"
                    }
                    Text(buttonLabel, fontSize = 14.sp)
                }
            } else {
                Button(
                    onClick = onTryAgain,
                    modifier = Modifier.heightIn(min = 48.dp)
                ) {
                    Text("Coba Lagi", fontSize = 14.sp)
                }
            }
        }
    )
}
