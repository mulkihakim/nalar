package com.mulki.nalar.core.network

import com.mulki.nalar.core.auth.TokenManager
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit

object NetworkModule {
    // HP terhubung via USB ke Android Studio dengan 'adb reverse tcp:8080 tcp:8080'
    // Gunakan 127.0.0.1 agar langsung tembus ke backend di laptop tanpa terhalang firewall Wi-Fi
    private const val BASE_URL = "http://127.0.0.1:8080/api/v1/"

    fun provideApiService(tokenManager: TokenManager): ApiService {
        val authInterceptor = Interceptor { chain ->
            val original = chain.request()
            val token = runBlocking { tokenManager.token.firstOrNull() }
            val requestBuilder = original.newBuilder()
            if (!token.isNullOrBlank()) {
                requestBuilder.addHeader("Authorization", "Bearer $token")
            }
            chain.proceed(requestBuilder.build())
        }

        val loggingInterceptor = HttpLoggingInterceptor().apply {
            level = HttpLoggingInterceptor.Level.BODY
        }

        val okHttpClient = OkHttpClient.Builder()
            .addInterceptor(authInterceptor)
            .addInterceptor(loggingInterceptor)
            .connectTimeout(15, TimeUnit.SECONDS)
            .readTimeout(15, TimeUnit.SECONDS)
            .build()

        val retrofit = Retrofit.Builder()
            .baseUrl(BASE_URL)
            .client(okHttpClient)
            .addConverterFactory(GsonConverterFactory.create())
            .build()

        return retrofit.create(ApiService::class.java)
    }
}
