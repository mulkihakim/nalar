package com.mulki.nalar.feature.argument

import androidx.activity.compose.BackHandler
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.automirrored.filled.MenuBook
import androidx.compose.material.icons.filled.ArrowDropDown
import androidx.compose.material.icons.filled.Check
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
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
import com.mulki.nalar.core.network.dto.OptionItem
import com.mulki.nalar.core.network.dto.SessionDetail
import com.mulki.nalar.core.network.dto.formatModeLabel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ArgumentScreen(
    viewModel: ArgumentViewModel,
    onSessionFinished: (SessionDetail) -> Unit,
    onExit: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()
    var showExitDialog by remember { mutableStateOf(false) }
    var showMaterialDialog by remember { mutableStateOf(false) }

    // Intersep tombol back hardware
    BackHandler {
        showExitDialog = true
    }

    // Jika sesi selesai dan dialog/monitoring sudah ditutup, picu navigasi ke ringkasan
    LaunchedEffect(uiState.isSessionCompleted, uiState.isShowingMonitoring, uiState.confirmResult) {
        if (uiState.isSessionCompleted && !uiState.isShowingMonitoring && uiState.confirmResult == null) {
            onSessionFinished(uiState.session)
        }
    }

    if (uiState.isShowingMonitoring) {
        val isLast = uiState.currentArgumentIndex >= uiState.session.arguments.size - 1
        MonitoringView(
            data = uiState.monitoringData,
            isLoading = uiState.isLoadingMonitoring,
            isLastArgument = isLast,
            onContinue = { viewModel.continueFromMonitoring() }
        )
        return
    }

    val currentArg = viewModel.getCurrentArgument()
    val totalArgs = uiState.session.arguments.size
    val currentPosition = (uiState.currentArgumentIndex + 1).coerceAtMost(totalArgs)

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Column {
                        Text(
                            text = "Argumen $currentPosition dari $totalArgs",
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Bold
                        )
                        Text(
                            text = uiState.session.exam_title,
                            style = MaterialTheme.typography.labelSmall,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    }
                },
                navigationIcon = {
                    IconButton(onClick = { showExitDialog = true }) {
                        Icon(imageVector = Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Keluar")
                    }
                },
                actions = {
                    IconButton(onClick = { showMaterialDialog = true }) {
                        Icon(imageVector = Icons.AutoMirrored.Filled.MenuBook, contentDescription = "Baca Materi")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.surface,
                    titleContentColor = MaterialTheme.colorScheme.onSurface
                )
            )
        },
        containerColor = MaterialTheme.colorScheme.background
    ) { innerPadding ->
        if (currentArg == null) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(innerPadding),
                contentAlignment = Alignment.Center
            ) {
                CircularProgressIndicator(color = MaterialTheme.colorScheme.primary)
            }
        } else {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(innerPadding)
                    .padding(16.dp)
                    .verticalScroll(rememberScrollState())
            ) {
                // Mode indicator badge
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .clip(RoundedCornerShape(8.dp))
                            .background(
                                if (uiState.session.mode == "help") {
                                    MaterialTheme.colorScheme.tertiaryContainer
                                } else {
                                    MaterialTheme.colorScheme.primaryContainer
                                }
                            )
                            .padding(horizontal = 10.dp, vertical = 4.dp)
                    ) {
                        Text(
                            text = "Mode: ${formatModeLabel(uiState.session.mode)}",
                            style = MaterialTheme.typography.labelSmall,
                            fontWeight = FontWeight.SemiBold,
                            color = if (uiState.session.mode == "help") {
                                MaterialTheme.colorScheme.onTertiaryContainer
                            } else {
                                MaterialTheme.colorScheme.onPrimaryContainer
                            }
                        )
                    }

                    Text(
                        text = "Percobaan #${uiState.session.attempt_no}",
                        style = MaterialTheme.typography.labelSmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }

                Spacer(modifier = Modifier.height(12.dp))

                // Claim Card (Bagian Atas sesuai 02-flows.md §1.4b)
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(16.dp),
                    colors = CardDefaults.cardColors(
                        containerColor = MaterialTheme.colorScheme.surface
                    ),
                    elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
                ) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text(
                            text = "Pernyataan / Claim:",
                            style = MaterialTheme.typography.labelMedium,
                            fontWeight = FontWeight.SemiBold,
                            color = MaterialTheme.colorScheme.primary
                        )
                        Spacer(modifier = Modifier.height(6.dp))
                        Text(
                            text = currentArg.claim_text,
                            style = MaterialTheme.typography.bodyLarge,
                            fontWeight = FontWeight.Bold,
                            color = MaterialTheme.colorScheme.onSurface,
                            lineHeight = 24.sp
                        )
                    }
                }

                Spacer(modifier = Modifier.height(16.dp))

                Text(
                    text = "Susun Argumen Toulmin:",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.SemiBold,
                    color = MaterialTheme.colorScheme.onSurface
                )
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = "Ketuk slot di bawah ini untuk memilih opsi yang tepat.",
                    style = MaterialTheme.typography.labelSmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )

                Spacer(modifier = Modifier.height(12.dp))

                // Slot Ground Card (Tersusun vertikal sesuai 02-flows.md §1.4b)
                SlotCard(
                    label = "Slot Ground (Data / Fakta Pendukung)",
                    selectedOption = uiState.selectedGroundOption,
                    placeholder = "Pilih Ground",
                    onClick = { viewModel.openBottomSheet("ground") }
                )

                Spacer(modifier = Modifier.height(12.dp))

                // Slot Warrant Card (Tersusun vertikal di bawah Ground)
                SlotCard(
                    label = "Slot Warrant (Prinsip / Penjamin Keterkaitan)",
                    selectedOption = uiState.selectedWarrantOption,
                    placeholder = "Pilih Warrant",
                    onClick = { viewModel.openBottomSheet("warrant") }
                )

                if (!uiState.errorMessage.isNullOrBlank()) {
                    Spacer(modifier = Modifier.height(12.dp))
                    Surface(
                        color = MaterialTheme.colorScheme.errorContainer,
                        shape = RoundedCornerShape(8.dp),
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Text(
                            text = uiState.errorMessage!!,
                            color = MaterialTheme.colorScheme.onErrorContainer,
                            style = MaterialTheme.typography.labelMedium,
                            modifier = Modifier.padding(12.dp)
                        )
                    }
                }

                Spacer(modifier = Modifier.height(24.dp))

                // Tombol Confirm (Aktif jika kedua slot terisi, 04-rules-design.md §B.5)
                val isBothSlotsFilled = uiState.selectedGroundOption != null && uiState.selectedWarrantOption != null

                Button(
                    onClick = { viewModel.confirmAnswers() },
                    enabled = isBothSlotsFilled && !uiState.isConfirming,
                    modifier = Modifier
                        .fillMaxWidth()
                        .heightIn(min = 48.dp), // Area sentuh minimal 48dp
                    shape = RoundedCornerShape(12.dp)
                ) {
                    if (uiState.isConfirming) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(20.dp),
                            color = MaterialTheme.colorScheme.onPrimary,
                            strokeWidth = 2.dp
                        )
                    } else {
                        Icon(imageVector = Icons.Default.Check, contentDescription = null)
                        Spacer(modifier = Modifier.width(8.dp))
                        Text("Konfirmasi Jawaban", fontSize = 16.sp)
                    }
                }
            }
        }
    }

    // Modal Bottom Sheet untuk Pemilihan Opsi Ground / Warrant
    if (uiState.activeBottomSheetType != null && currentArg != null) {
        val type = uiState.activeBottomSheetType!!
        val optionsForSlot = currentArg.options.filter { it.type == type }
        val currentSelectedId = if (type == "ground") {
            uiState.selectedGroundOption?.id
        } else {
            uiState.selectedWarrantOption?.id
        }

        OptionBottomSheet(
            title = if (type == "ground") "Pilih Opsi Ground" else "Pilih Opsi Warrant",
            options = optionsForSlot,
            selectedOptionId = currentSelectedId,
            onSelectOption = { option -> viewModel.selectOption(option) },
            onDismiss = { viewModel.dismissBottomSheet() }
        )
    }

    // Modal Feedback Hasil Confirm
    if (uiState.confirmResult != null) {
        ConfirmResultDialog(
            result = uiState.confirmResult!!,
            mode = uiState.session.mode,
            onContinue = { viewModel.handleContinueAfterConfirm() },
            onTryAgain = { viewModel.tryAgain() }
        )
    }

    // Dialog Baca Materi Teks
    if (showMaterialDialog) {
        AlertDialog(
            onDismissRequest = { showMaterialDialog = false },
            containerColor = MaterialTheme.colorScheme.surface,
            title = {
                Text(
                    text = uiState.session.material_title,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
            },
            text = {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .heightIn(max = 350.dp)
                        .verticalScroll(rememberScrollState())
                ) {
                    Text(
                        text = uiState.session.material_content.ifBlank { "Tidak ada isi materi." },
                        style = MaterialTheme.typography.bodyMedium,
                        lineHeight = 22.sp
                    )
                }
            },
            confirmButton = {
                Button(onClick = { showMaterialDialog = false }) {
                    Text("Tutup")
                }
            }
        )
    }

    // Dialog Konfirmasi Keluar dari Pengerjaan
    if (showExitDialog) {
        AlertDialog(
            onDismissRequest = { showExitDialog = false },
            containerColor = MaterialTheme.colorScheme.surface,
            title = {
                Text("Keluar dari Latihan?", fontWeight = FontWeight.Bold)
            },
            text = {
                Text("Progres pengerjaan Anda tetap tersimpan di server. Anda dapat melanjutkannya nanti kapan saja.")
            },
            confirmButton = {
                Button(
                    onClick = {
                        showExitDialog = false
                        onExit()
                    },
                    modifier = Modifier.heightIn(min = 48.dp)
                ) {
                    Text("Ya, Keluar")
                }
            },
            dismissButton = {
                OutlinedButton(
                    onClick = { showExitDialog = false },
                    modifier = Modifier.heightIn(min = 48.dp)
                ) {
                    Text("Batal")
                }
            }
        )
    }
}

@Composable
private fun SlotCard(
    label: String,
    selectedOption: OptionItem?,
    placeholder: String,
    onClick: () -> Unit
) {
    val isFilled = selectedOption != null

    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = if (isFilled) {
                MaterialTheme.colorScheme.primaryContainer.copy(alpha = 0.3f)
            } else {
                MaterialTheme.colorScheme.surface
            }
        ),
        // Sesuai 04-rules-design.md §B.5:
        // Slot kosong: border slate-200. Slot terisi: border primary (indigo)
        border = BorderStroke(
            width = if (isFilled) 1.5.dp else 1.dp,
            color = if (isFilled) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.outline
        )
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(14.dp)
        ) {
            Text(
                text = label,
                style = MaterialTheme.typography.labelSmall,
                fontWeight = FontWeight.SemiBold,
                color = if (isFilled) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.onSurfaceVariant
            )

            Spacer(modifier = Modifier.height(8.dp))

            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                if (isFilled) {
                    Text(
                        text = selectedOption!!.text,
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurface,
                        modifier = Modifier.weight(1f)
                    )
                } else {
                    Text(
                        text = placeholder,
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        modifier = Modifier.weight(1f)
                    )
                }

                Spacer(modifier = Modifier.width(8.dp))

                Icon(
                    imageVector = Icons.Default.ArrowDropDown,
                    contentDescription = "Buka Pilihan",
                    tint = if (isFilled) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
        }
    }
}
