package com.mulki.nalar.ui.theme

import androidx.compose.ui.graphics.Color

// Nalar Design System — 04-rules-design.md §B.2
val Primary = Color(0xFF4F46E5)             // indigo-600 — tombol utama, logo, item nav aktif
val PrimaryHover = Color(0xFF4338CA)        // indigo-700 — hover/active state
val PrimaryContainer = Color(0xFFEEF2FF)    // indigo-50  — latar item aktif / pill

val Background = Color(0xFFF8FAFC)          // slate-50   — latar halaman
val Surface = Color(0xFFFFFFFF)             // white      — card, header, bottom sheet
val SurfaceVariant = Color(0xFFF1F5F9)      // slate-100  — latar sekunder / kartu netral
val Border = Color(0xFFE2E8F0)              // slate-200  — border card, input, divider
val BorderStrong = Color(0xFFCBD5E1)        // slate-300  — border sekunder

val TextPrimary = Color(0xFF0F172A)         // slate-900  — judul, teks utama
val TextSecondary = Color(0xFF64748B)       // slate-500  — deskripsi, label, teks pembantu

val Success = Color(0xFF059669)             // emerald-600 — argumen selesai, ujian aktif, tanpa salah
val SuccessContainer = Color(0xFFECFDF5)    // emerald-50  — latar badge sukses
val OnSuccessContainer = Color(0xFF065F46)  // emerald-800 — teks pada latar sukses

val Danger = Color(0xFFDC2626)              // red-600     — error, ujian nonaktif, validasi gagal
val DangerContainer = Color(0xFFFEF2F2)     // red-50      — latar pesan error
val OnDangerContainer = Color(0xFF991B1B)   // red-800     — teks pesan error

val Warning = Color(0xFFD97706)             // amber-600   — indikator mode bantuan, peringatan ganti mode
val WarningContainer = Color(0xFFFFFBEB)    // amber-50    — latar badge peringatan
val OnWarningContainer = Color(0xFF92400E)  // amber-800   — teks pada latar peringatan