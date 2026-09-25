package com.mulki.nalar.core.network.dto

data class MaterialSummaryDto(
    val id: Long,
    val title: String,
    val content: String? = null
)

data class ExamDto(
    val id: Long,
    val title: String,
    val material_id: Long,
    val material: MaterialSummaryDto?,
    val is_active: Boolean,
    val arguments_per_session: Int
)
