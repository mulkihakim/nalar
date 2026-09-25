package com.mulki.nalar.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

// Mapping palet Nalar ke MaterialTheme.colorScheme (04-rules-design.md §B.2 & §B.6)
// primary = Indigo (utama), secondary = Emerald (success), tertiary = Amber (warning), error = Red (danger)
private val NalarLightColorScheme = lightColorScheme(
    primary = Primary,
    onPrimary = Color.White,
    primaryContainer = PrimaryContainer,
    onPrimaryContainer = PrimaryHover,

    secondary = Success,
    onSecondary = Color.White,
    secondaryContainer = SuccessContainer,
    onSecondaryContainer = OnSuccessContainer,

    tertiary = Warning,
    onTertiary = Color.White,
    tertiaryContainer = WarningContainer,
    onTertiaryContainer = OnWarningContainer,

    error = Danger,
    onError = Color.White,
    errorContainer = DangerContainer,
    onErrorContainer = OnDangerContainer,

    background = Background,
    onBackground = TextPrimary,

    surface = Surface,
    onSurface = TextPrimary,
    surfaceVariant = SurfaceVariant,
    onSurfaceVariant = TextSecondary,

    outline = Border,
    outlineVariant = BorderStrong
)

@Composable
fun NalarTheme(
    content: @Composable () -> Unit
) {
    MaterialTheme(
        colorScheme = NalarLightColorScheme,
        typography = Typography,
        content = content
    )
}