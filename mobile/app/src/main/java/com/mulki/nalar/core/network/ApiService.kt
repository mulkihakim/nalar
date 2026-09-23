package com.mulki.nalar.core.network

import com.mulki.nalar.core.network.dto.LoginRequest
import com.mulki.nalar.core.network.dto.LoginResponse
import com.mulki.nalar.core.network.dto.MeResponse
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST

interface ApiService {
    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<LoginResponse>

    @GET("auth/me")
    suspend fun getMe(): Response<MeResponse>
}
