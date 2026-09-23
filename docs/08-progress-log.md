# 08. Progress Log — Catatan Implementasi per Fitur

Prasyarat: `06-tasks.md` (Definition of Done), `07-tech-stack.md` (batasan library).

> File ini diperbarui **oleh agent coding**, **setelah** selesai mengerjakan satu fitur/task — bukan sebelum, bukan sebagai rencana. Tujuannya: sesi agent berikutnya (atau Anda) tahu apa yang sudah ada, keputusan yang diambil saat coding, dan sejauh mana Definition of Done terpenuhi — tanpa harus membaca ulang seluruh kode.

## Instruksi untuk agent (baca sebelum menutup task)

Setelah menyelesaikan satu fitur (misal F-08, atau "endpoint drop/confirm"), tambahkan entri baru **di paling atas** bagian `## Log` di bawah, dengan format ini:

```
### YYYY-MM-DD — <ID Fitur bila ada> <nama fitur singkat>
- **Status:** selesai / sebagian (sebutkan sisanya) / diblok (sebutkan alasannya)
- **File dibuat/diubah:** daftar path relatif
- **Keputusan implementasi:** hal yang TIDAK dijelaskan eksplisit di 02–07-*.md tapi harus
  diputuskan saat coding, plus alasannya. Kalau agent "menebak" sesuatu yang seharusnya
  dikonfirmasi ke pemilik proyek, WAJIB ditulis di sini, jangan diam-diam diputuskan.
- **Deviasi dari dokumen:** bagian kode yang TIDAK mengikuti PRD/rules/architecture secara
  ketat, dengan alasan. Kalau tidak ada, tulis "tidak ada".
- **Definition of Done:** cocokkan ke nomor kriteria penerimaan relevan di `06-tasks.md` §2,
  tandai per poin: ✅ terpenuhi / ❌ belum / N/A tidak relevan untuk fitur ini
- **Test:** ringkas apa yang ditest dan di mana (path file test)
- **TODO lanjutan:** kalau ada bagian yang dipotong/ditunda
```

Aturan tambahan:
- Satu entri per sesi kerja/fitur, jangan digabung banyak fitur dalam satu entri.
- Kalau memperbaiki/melanjutkan fitur yang sudah ada entrinya, buat entri baru (jangan edit entri lama), supaya riwayat keputusan tetap terlihat.
- Kalau menemukan bagian dokumen (`01`–`07`) yang bertentangan satu sama lain atau dengan kode yang sudah ada, tulis di entri sebagai "Perlu klarifikasi" alih-alih menebak.

## Log

<!-- Entri terbaru ditambahkan di bawah baris ini, urutan terbaru dulu -->

### 2026-09-23 — F-01 Login/logout, role-based access di backend Go
- **Status:** selesai
- **File dibuat/diubah:**
  - `backend/go.mod`
  - `backend/go.sum`
  - `backend/cmd/api/main.go`
  - `backend/internal/db/db.go`
  - `backend/internal/user/model.go`
  - `backend/internal/user/service.go`
  - `backend/internal/user/handler.go`
  - `backend/internal/user/service_test.go`
  - `backend/internal/middleware/auth.go`
  - `backend/internal/middleware/ratelimit.go`
  - `backend/internal/middleware/auth_test.go`
  - `docs/08-progress-log.md`
- **Keputusan implementasi:**
  - Menambahkan library `github.com/golang-jwt/jwt/v5` ke `go.mod` untuk token JWT (HMAC-SHA256) dengan klaim `user_id` dan `role`, serta masa aktif default 24 jam.
  - Rate limiting pada endpoint login diimplementasikan menggunakan in-memory sliding bucket berbasis client IP (`internal/middleware/ratelimit.go`, 5 request/menit, response 429 Too Many Requests) tanpa menambah dependency eksternal.
  - Struct `User` di `model.go` diberi tag `json:"-"` pada `PasswordHash` guna mencegah kebocoran hash password di respons JSON API.
  - Endpoint `GET /auth/me` mengembalikan format `{ "user": { ... } }` yang konsisten dengan respons `POST /auth/login`.
  - Logout bersifat stateless di sisi klien (klien menghapus JWT dari penyimpanan lokal), dengan expiry token di server.
  - Verifikasi migrasi: file `backend/migrations/000001_create_users_table.up.sql` sudah mencakup seluruh kolom di `04-architecture.md` §1 (`id`, `name`, `username`, `password_hash`, `role`, `created_by`, `is_active`), sehingga tidak memerlukan migrasi SQL baru.
  - Database pool di `internal/db/db.go` dikonfigurasi: MaxIdleConns=10, MaxOpenConns=50, ConnMaxLifetime=1 jam.
- **Deviasi dari dokumen:** tidak ada.
- **Definition of Done:**
  - Kriteria penerimaan 1–14 di `06-tasks.md` §2 berfokus pada domain ujian/analitik/sesi/isolasi asesor (N/A untuk F-01 langsung).
  - Catatan DoD Tambahan untuk Auth (F-01) yang perlu dicatat secara eksplisit:
    - [x] Login berhasil mengembalikan token JWT yang valid dan memuat claim `user_id` serta `role`.
    - [x] Login dengan password salah atau username tidak terdaftar ditolak dengan pesan error dan status 401.
    - [x] Login dengan akun non-aktif (`is_active = false`) ditolak dengan status 403.
    - [x] Endpoint `GET /auth/me` terproteksi JWT dan mengembalikan identitas pengguna yang sedang login.
    - [x] Middleware `RequireRole` menolak akses pengguna yang tidak memiliki role yang diizinkan (403 Forbidden).
    - [x] Rate limiting membatasi percobaan login beruntun per IP (429 Too Many Requests).
- **Test:**
  - `backend/internal/user/service_test.go`: test login sukses, password salah, user tidak ditemukan, user tidak aktif, GetMe sukses/not found/inactive, verifikasi claim token JWT.
  - `backend/internal/middleware/auth_test.go`: test validasi Bearer token, header authorization hilang/salah, proteksi RequireRole, rate limiting per IP.
  - Hasil: semua unit test pass (`go test -v ./...`).
- **TODO lanjutan:**
  - Integrasi UI login dan simpan token di web frontend (React) dan mobile (Android) pada task berikutnya.

