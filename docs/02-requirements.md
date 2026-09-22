# 02. Requirements — Fitur, Prioritas, Mode Belajar

Prasyarat: baca `01-overview.md` untuk konteks istilah domain.

## 1. Fitur dan prioritas

Prioritas: **Must** (wajib untuk demo), **Should** (penting tapi bisa dipotong jika waktu mepet), **Could** (nice-to-have). Urutan pemotongan jika waktu mepet ada di `06-tasks.md`.

| ID | Fitur | Prioritas | Platform |
|---|---|---|---|
| F-01 | Login/logout, role-based access | Must | Semua |
| F-02 | Admin kelola asesor; admin/asesor buat siswa | Must | Web |
| F-03 | CRUD kelas dan anggota kelas | Must | Web |
| F-04 | CRUD materi dan argumen (claim + 4 ground + 4 warrant, validasi) | Must | Web |
| F-05 | CRUD ujian: pilih materi, aktif/nonaktif, atur akses kelas/siswa, atur `arguments_per_session` | Must | Web |
| F-06 | Siswa melihat daftar ujian yang boleh diakses | Must | Mobile |
| F-07 | Siswa memilih mode, lalu server membuat sesi baru (subset argumen acak) atau melanjutkan sesi berjalan | Must | Mobile + backend |
| F-08 | Pengerjaan: baca materi, susun argumen dengan tap-to-select (ketuk slot → pilih opsi di bottom sheet), Confirm, lanjut setelah benar | Must | Mobile |
| F-09 | Pencatatan log setiap drop | Must | Backend |
| F-10 | Sesi berjalan bisa dilanjutkan setelah keluar/gangguan jaringan, termasuk konfirmasi jika mode belajar diubah | Must | Mobile + backend |
| F-11 | Mode standar, bantuan, dan analitik sosial | Must | Mobile + backend |
| F-12 | Hasil sendiri untuk siswa; hasil semua siswa untuk admin/asesor | Must | Mobile + web |
| F-13 | Log percobaan dengan identitas untuk admin/asesor | Must | Web |
| F-14 | Siswa mengerjakan lewat web (reuse API yang sama) | Should | Web |
| F-15 | Urutan opsi diacak, stabil per sesi | Should | Backend |
| F-16 | Ekspor hasil/log ke CSV | Should | Web |
| F-17 | Impor siswa lewat CSV, mode gelap | Could | Web |

> Catatan F-08 vs F-14: layar pengerjaan di mobile pakai tap-to-select (lihat `03-flows.md` §4b), sedangkan versi web (F-14) tetap pakai drag-and-drop (§4a) karena "mengikuti desain lama" yang sudah baik di layar besar. Endpoint backend sama untuk keduanya.

## 2. Mode belajar

| Mode | Perilaku |
|---|---|
| **Standar** | Siswa tidak diberi tahu mana opsi yang benar atau salah, kecuali hasil Confirm (benar/belum tepat). |
| **Bantuan** | Siswa mengetahui bagian mana yang salah dan benar (lihat asumsi A3 di `05-rules.md`). |
| **Analitik sosial** | Perilaku seperti standar, ditambah tampilan perbandingan pilihan kelompok (detail di `03-flows.md` §8). |

Siswa memilih salah satu dari ketiganya setiap kali menekan sebuah ujian. Untuk sesi baru, mode itu langsung berlaku. Untuk sesi yang sedang berjalan, mode lama tetap dipakai kecuali siswa memilih mode lain dan mengonfirmasi perubahannya (lihat `03-flows.md` §1, langkah 2).
