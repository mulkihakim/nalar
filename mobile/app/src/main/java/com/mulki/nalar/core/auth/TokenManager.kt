package com.mulki.nalar.core.auth

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.longPreferencesKey
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import com.mulki.nalar.core.network.dto.UserDto
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "auth_prefs")

class TokenManager(private val context: Context) {
    companion object {
        private val TOKEN_KEY = stringPreferencesKey("jwt_token")
        private val USER_ID_KEY = longPreferencesKey("user_id")
        private val USER_NAME_KEY = stringPreferencesKey("user_name")
        private val USER_USERNAME_KEY = stringPreferencesKey("user_username")
        private val USER_ROLE_KEY = stringPreferencesKey("user_role")
    }

    val token: Flow<String?> = context.dataStore.data.map { preferences ->
        preferences[TOKEN_KEY]
    }

    val role: Flow<String?> = context.dataStore.data.map { preferences ->
        preferences[USER_ROLE_KEY]
    }

    val userName: Flow<String?> = context.dataStore.data.map { preferences ->
        preferences[USER_NAME_KEY]
    }

    suspend fun saveAuth(token: String, user: UserDto) {
        context.dataStore.edit { preferences ->
            preferences[TOKEN_KEY] = token
            preferences[USER_ID_KEY] = user.id
            preferences[USER_NAME_KEY] = user.name
            preferences[USER_USERNAME_KEY] = user.username
            preferences[USER_ROLE_KEY] = user.role
        }
    }

    suspend fun clearAuth() {
        context.dataStore.edit { preferences ->
            preferences.clear()
        }
    }
}
