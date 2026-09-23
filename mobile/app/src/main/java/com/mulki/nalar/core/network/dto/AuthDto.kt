package com.mulki.nalar.core.network.dto

data class LoginRequest(
    val username: String,
    val password: String
)

data class UserDto(
    val id: Long,
    val name: String,
    val username: String,
    val role: String,
    val is_active: Boolean
)

data class LoginResponse(
    val token: String,
    val user: UserDto
)

data class MeResponse(
    val user: UserDto
)
