package com.mulki.nalar.core.network

import com.mulki.nalar.core.network.dto.AnalysisResponse
import com.mulki.nalar.core.network.dto.ConfirmRequest
import com.mulki.nalar.core.network.dto.ConfirmResult
import com.mulki.nalar.core.network.dto.DropRequest
import com.mulki.nalar.core.network.dto.ExamDto
import com.mulki.nalar.core.network.dto.LoginRequest
import com.mulki.nalar.core.network.dto.LoginResponse
import com.mulki.nalar.core.network.dto.MeResponse
import com.mulki.nalar.core.network.dto.MonitoringResponse
import com.mulki.nalar.core.network.dto.SessionDetail
import com.mulki.nalar.core.network.dto.SessionHistoryItem
import com.mulki.nalar.core.network.dto.SessionStatusResponse
import com.mulki.nalar.core.network.dto.StartSessionRequest
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path

interface ApiService {
    // Auth (F-01)
    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<LoginResponse>

    @GET("auth/me")
    suspend fun getMe(): Response<MeResponse>

    // Daftar Ujian Siswa (F-06)
    @GET("my/exams")
    suspend fun getMyExams(): Response<List<ExamDto>>

    // Status Sesi & Mulai/Lanjut Sesi (F-07, F-10)
    @GET("exams/{examId}/session-status")
    suspend fun getSessionStatus(@Path("examId") examId: Long): Response<SessionStatusResponse>

    @POST("exams/{examId}/start")
    suspend fun startSession(
        @Path("examId") examId: Long,
        @Body request: StartSessionRequest
    ): Response<SessionDetail>

    @GET("sessions/{sessionId}")
    suspend fun getSessionDetail(@Path("sessionId") sessionId: Long): Response<SessionDetail>

    // Pengerjaan Argumen: Drop & Confirm (F-08, F-09, F-11)
    @POST("sessions/{sessionId}/arguments/{argumentId}/drops")
    suspend fun recordDrop(
        @Path("sessionId") sessionId: Long,
        @Path("argumentId") argumentId: Long,
        @Body request: DropRequest
    ): Response<Map<String, String>>

    @POST("sessions/{sessionId}/arguments/{argumentId}/confirm")
    suspend fun confirmArgument(
        @Path("sessionId") sessionId: Long,
        @Path("argumentId") argumentId: Long,
        @Body request: ConfirmRequest
    ): Response<ConfirmResult>

    // Analitik Sosial: Monitoring & Analysis (F-11)
    @GET("sessions/{sessionId}/monitoring/{argumentId}")
    suspend fun getMonitoring(
        @Path("sessionId") sessionId: Long,
        @Path("argumentId") argumentId: Long
    ): Response<MonitoringResponse>

    @GET("sessions/{sessionId}/analysis")
    suspend fun getAnalysis(@Path("sessionId") sessionId: Long): Response<AnalysisResponse>

    // Riwayat Sesi Siswa (F-12)
    @GET("my/exams/{examId}/sessions")
    suspend fun getMyExamSessions(@Path("examId") examId: Long): Response<List<SessionHistoryItem>>
}
