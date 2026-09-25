package com.mulki.nalar.feature.session

import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.size
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.mulki.nalar.core.network.dto.formatModeLabel

@Composable
fun ModeConfirmDialog(
    currentMode: String,
    targetMode: String,
    isLoading: Boolean,
    onConfirm: () -> Unit,
    onCancel: () -> Unit
) {
    AlertDialog(
        onDismissRequest = { if (!isLoading) onCancel() },
        containerColor = MaterialTheme.colorScheme.surface,
        title = {
            Text(
                text = "Konfirmasi Perubahan Mode",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Bold
            )
        },
        text = {
            Text(
                text = "Mode belajar akan diubah dari ${formatModeLabel(currentMode)} ke ${formatModeLabel(targetMode)}, lanjutkan sesi ini?",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurface
            )
        },
        confirmButton = {
            Button(
                onClick = onConfirm,
                enabled = !isLoading,
                modifier = Modifier.heightIn(min = 48.dp)
            ) {
                if (isLoading) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(18.dp),
                        color = MaterialTheme.colorScheme.onPrimary,
                        strokeWidth = 2.dp
                    )
                } else {
                    Text("Ya, Lanjutkan", fontSize = 14.sp)
                }
            }
        },
        dismissButton = {
            OutlinedButton(
                onClick = onCancel,
                enabled = !isLoading,
                modifier = Modifier.heightIn(min = 48.dp)
            ) {
                Text("Batal", fontSize = 14.sp)
            }
        }
    )
}
