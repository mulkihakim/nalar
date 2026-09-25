package com.mulki.nalar.feature.navigation

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.Assignment
import androidx.compose.material.icons.automirrored.outlined.Assignment
import androidx.compose.material.icons.filled.Assessment
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.outlined.Assessment
import androidx.compose.material.icons.outlined.Person
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import com.mulki.nalar.core.network.ApiService
import com.mulki.nalar.core.network.dto.SessionDetail
import com.mulki.nalar.feature.argument.ArgumentScreen
import com.mulki.nalar.feature.argument.ArgumentViewModel
import com.mulki.nalar.feature.argument.MaterialReadScreen
import com.mulki.nalar.feature.examlist.ExamListScreen
import com.mulki.nalar.feature.examlist.ExamListViewModel
import com.mulki.nalar.feature.profile.ProfileScreen
import com.mulki.nalar.feature.result.AnalysisScreen
import com.mulki.nalar.feature.result.ResultScreen
import com.mulki.nalar.feature.result.ResultViewModel
import com.mulki.nalar.feature.result.SessionSummaryScreen
import com.mulki.nalar.feature.session.ModeConfirmDialog
import com.mulki.nalar.feature.session.ModeSelectDialog
import com.mulki.nalar.feature.session.SessionViewModel

sealed class BottomNavItem(
    val route: String,
    val title: String,
    val selectedIcon: ImageVector,
    val unselectedIcon: ImageVector
) {
    object Exams : BottomNavItem("exams", "Ujian", Icons.AutoMirrored.Filled.Assignment, Icons.AutoMirrored.Outlined.Assignment)
    object Results : BottomNavItem("results", "Hasil", Icons.Filled.Assessment, Icons.Outlined.Assessment)
    object Profile : BottomNavItem("profile", "Profil", Icons.Filled.Person, Icons.Outlined.Person)
}

@Composable
fun NalarApp(
    apiService: ApiService,
    userName: String,
    username: String,
    userRole: String,
    onLogout: () -> Unit
) {
    val navController = rememberNavController()
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentRoute = navBackStackEntry?.destination?.route

    // ViewModel Singletons untuk tab
    val examListViewModel: ExamListViewModel = viewModel(factory = ExamListViewModel.Factory(apiService))
    val sessionViewModel: SessionViewModel = viewModel(factory = SessionViewModel.Factory(apiService))
    val resultViewModel: ResultViewModel = viewModel(factory = ResultViewModel.Factory(apiService))

    val sessionDialogState by sessionViewModel.state.collectAsState()

    // State untuk pengerjaan aktif
    var activeSessionDetail by remember { mutableStateOf<SessionDetail?>(null) }

    // Pantau ketika sesi berhasil dibuat / dilanjutkan dari dialog
    LaunchedEffect(sessionDialogState.startedSession) {
        val started = sessionDialogState.startedSession
        if (started != null) {
            activeSessionDetail = started
            sessionViewModel.clearStartedSession()
            // Tampilkan layar baca materi terlebih dahulu sesuai alur 02-flows.md §1 langkah 3
            navController.navigate("material_read")
        }
    }

    val bottomNavItems = listOf(
        BottomNavItem.Exams,
        BottomNavItem.Results,
        BottomNavItem.Profile
    )

    // Tampilkan BottomBar hanya di top-level tabs
    val showBottomBar = currentRoute in listOf(
        BottomNavItem.Exams.route,
        BottomNavItem.Results.route,
        BottomNavItem.Profile.route
    )

    Scaffold(
        bottomBar = {
            if (showBottomBar) {
                NavigationBar(
                    containerColor = MaterialTheme.colorScheme.surface,
                    contentColor = MaterialTheme.colorScheme.onSurface
                ) {
                    bottomNavItems.forEach { item ->
                        val isSelected = currentRoute == item.route
                        NavigationBarItem(
                            selected = isSelected,
                            onClick = {
                                examListViewModel.closeExamHistory()
                                resultViewModel.clearAnalysis()
                                if (currentRoute != item.route) {
                                    if (item.route == BottomNavItem.Results.route) {
                                        resultViewModel.loadExams()
                                    } else if (item.route == BottomNavItem.Exams.route) {
                                        examListViewModel.loadExams()
                                    }
                                    navController.navigate(item.route) {
                                        popUpTo(BottomNavItem.Exams.route) {
                                            inclusive = false
                                        }
                                        launchSingleTop = true
                                    }
                                }
                            },
                            icon = {
                                Icon(
                                    imageVector = if (isSelected) item.selectedIcon else item.unselectedIcon,
                                    contentDescription = item.title
                                )
                            },
                            label = { Text(item.title, style = MaterialTheme.typography.labelSmall) },
                            colors = NavigationBarItemDefaults.colors(
                                selectedIconColor = MaterialTheme.colorScheme.primary,
                                selectedTextColor = MaterialTheme.colorScheme.primary,
                                unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                                unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant,
                                indicatorColor = MaterialTheme.colorScheme.primaryContainer
                            )
                        )
                    }
                }
            }
        }
    ) { innerPadding ->
        NavHost(
            navController = navController,
            startDestination = BottomNavItem.Exams.route,
            modifier = Modifier
                .fillMaxSize()
                .padding(innerPadding)
        ) {
            // Tab 1: Daftar Ujian
            composable(BottomNavItem.Exams.route) {
                ExamListScreen(
                    viewModel = examListViewModel,
                    onSelectExam = { exam -> sessionViewModel.openModeSelector(exam) }
                )
            }

            // Tab 2: Hasil & Riwayat Sesi
            composable(BottomNavItem.Results.route) {
                ResultScreen(
                    viewModel = resultViewModel,
                    onResumeSession = { exam, _ ->
                        sessionViewModel.openModeSelector(exam)
                    }
                )
            }

            // Tab 3: Profil
            composable(BottomNavItem.Profile.route) {
                ProfileScreen(
                    name = userName,
                    username = username,
                    role = userRole,
                    onLogout = onLogout
                )
            }

            // Sub-screen: Baca Materi
            composable("material_read") {
                val session = activeSessionDetail
                if (session != null) {
                    MaterialReadScreen(
                        examTitle = session.exam_title,
                        materialTitle = session.material_title,
                        materialContent = session.material_content,
                        onStartExam = {
                            navController.navigate("argument_screen") {
                                popUpTo("material_read") { inclusive = true }
                            }
                        },
                        onBack = {
                            navController.popBackStack()
                        }
                    )
                }
            }

            // Sub-screen: Layar Pengerjaan Argumen Tap-to-Select
            composable("argument_screen") {
                val session = activeSessionDetail
                if (session != null) {
                    val argumentViewModel: ArgumentViewModel = viewModel(
                        key = "arg_vm_${session.id}_${session.mode}",
                        factory = ArgumentViewModel.Factory(session, apiService)
                    )

                    ArgumentScreen(
                        viewModel = argumentViewModel,
                        onSessionFinished = { finishedSession ->
                            activeSessionDetail = finishedSession
                            navController.navigate("session_summary") {
                                popUpTo("argument_screen") { inclusive = true }
                            }
                        },
                        onExit = {
                            navController.navigate(BottomNavItem.Exams.route) {
                                popUpTo(BottomNavItem.Exams.route) { inclusive = true }
                            }
                        }
                    )
                }
            }

            // Sub-screen: Ringkasan Selesai Sesi
            composable("session_summary") {
                val session = activeSessionDetail
                if (session != null) {
                    SessionSummaryScreen(
                        session = session,
                        onViewAnalysis = {
                            resultViewModel.loadAnalysisForSession(session.id)
                            navController.navigate("session_analysis")
                        },
                        onFinish = {
                            examListViewModel.loadExams()
                            resultViewModel.loadExams()
                            navController.navigate(BottomNavItem.Exams.route) {
                                popUpTo(BottomNavItem.Exams.route) { inclusive = true }
                            }
                        }
                    )
                }
            }

            // Sub-screen: Analisis Sosial A:B
            composable("session_analysis") {
                val resultState by resultViewModel.uiState.collectAsState()
                AnalysisScreen(
                    data = resultState.analysisData,
                    isLoading = resultState.isLoadingAnalysis,
                    onBack = {
                        navController.popBackStack()
                    }
                )
            }
        }

        // Dialog Pemilihan Mode (F-07)
        if (sessionDialogState.selectedExam != null && sessionDialogState.conflictCurrentMode == null) {
            ModeSelectDialog(
                exam = sessionDialogState.selectedExam!!,
                sessionStatus = sessionDialogState.sessionStatus,
                isCheckingStatus = sessionDialogState.isCheckingStatus,
                isStarting = sessionDialogState.isStarting,
                errorMessage = sessionDialogState.errorMessage,
                onSelectMode = { mode -> sessionViewModel.selectMode(mode) },
                onDismiss = { sessionViewModel.dismissDialog() }
            )
        }

        // Dialog Konfirmasi Perubahan Mode (F-10 ketika 409 Conflict)
        if (sessionDialogState.conflictCurrentMode != null && sessionDialogState.conflictTargetMode != null) {
            ModeConfirmDialog(
                currentMode = sessionDialogState.conflictCurrentMode!!,
                targetMode = sessionDialogState.conflictTargetMode!!,
                isLoading = sessionDialogState.isStarting,
                onConfirm = { sessionViewModel.confirmModeChange() },
                onCancel = { sessionViewModel.cancelModeConflict() }
            )
        }
    }
}
