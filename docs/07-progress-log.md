# 07. Progress Log — Catatan Implementasi per Fitur

Prasyarat: `06-tasks.md` (Definition of Done), `05-tech-stack.md` (batasan library).

> Diperbarui **oleh agent coding**, **setelah** selesai satu fitur/task — bukan sebelum. Tujuan: sesi agent berikutnya tahu apa yang sudah ada tanpa baca ulang seluruh kode.

## Cara menulis entri baru

Tambahkan **di paling atas** bagian `## Log`, format singkat:

```
### YYYY-MM-DD — <ID Fitur> <nama singkat>
- Status: selesai / sebagian (+sisanya) / diblok (+alasan)
- File: daftar path relatif (ringkas, bukan tiap file kalau jumlahnya besar — cukup folder/pola)
- Keputusan: hal yang TIDAK eksplisit di 01–06-*.md tapi harus diputuskan saat coding + alasannya.
  Kalau ini seharusnya dikonfirmasi ke pemilik proyek, WAJIB ditulis di sini.
- Deviasi: bagian kode yang tidak mengikuti dokumen secara ketat + alasan. "tidak ada" kalau tidak ada.
- DoD: nomor kriteria relevan di 06-tasks.md §2 → ✅ / ❌ / N/A
- Test: apa yang ditest, di mana (path)
- TODO: bagian yang dipotong/ditunda (kalau ada)
```

**3 aturan:**
1. Satu entri per sesi kerja/fitur — jangan digabung.
2. Melanjutkan fitur lama → entri baru (jangan edit entri lama).
3. Dokumen (`01`–`06`) yang bertentangan satu sama lain atau dengan kode → tulis "Perlu klarifikasi", jangan menebak.

---

## Log

<!-- Entri terbaru ditambahkan di bawah baris ini, urutan terbaru dulu -->

### 2026-09-23 — F-01 Login screen, role guard, dan token persistence (Web & Mobile)
- **Status:** selesai
- **File dibuat/diubah:**
  - `web/package.json`
  - `web/package-lock.json`
  - `web/src/types/index.ts`
  - `web/src/lib/api-client.ts`
  - `web/src/lib/query-client.ts`
  - `web/src/features/auth/types.ts`
  - `web/src/features/auth/api.ts`
  - `web/src/features/auth/context.tsx`
  - `web/src/features/auth/components/LoginForm.tsx`
  - `web/src/features/auth/components/ProtectedRoute.tsx`
  - `web/src/components/ui/input.tsx`
  - `web/src/components/ui/label.tsx`
  - `web/src/components/ui/card.tsx`
  - `web/src/routes/__root.tsx`
  - `web/src/routes/login.tsx`
  - `web/src/routes/router.tsx`
  - `web/src/routes/admin/users.tsx`
  - `web/src/routes/admin/classes.tsx`
  - `web/src/routes/student/exams.tsx`
  - `web/src/main.tsx`
  - `mobile/gradle/libs.versions.toml`
  - `mobile/app/build.gradle.kts`
  - `mobile/app/src/main/java/com/mulki/nalar/core/network/dto/AuthDto.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/core/auth/TokenManager.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/core/network/ApiService.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/core/network/NetworkModule.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/feature/auth/LoginViewModel.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/feature/auth/LoginScreen.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/MainActivity.kt`
  - `docs/07-progress-log.md`
  - `backend/internal/middleware/cors.go`
  - `backend/cmd/api/main.go`
  - `backend/cmd/seed/main.go`
  - `backend/.env.example`
  - `web/vite.config.ts`
- **Keputusan implementasi:**
  - CORS & Dev Proxy: Menambahkan middleware CORS di backend Go (`internal/middleware/cors.go`) yang menangani preflight `OPTIONS` (204 No Content) dan header `Access-Control-Allow-*`. Selain itu, `web/vite.config.ts` dikonfigurasi dengan proxy `/api` menuju `http://localhost:8080`, dan `api-client.ts` disesuaikan untuk memakai relative URL `/api/v1` sehingga request di dev server berjalan seamlessly tanpa kendala CORS.
  - Layout Panel & Sidebar: Memindahkan menu navigasi dari header ke Sidebar yang dapat di-toggle (`PanelLeft` icon toggle, transisi halus lebar `w-60` ke `w-0`) sesuai `04-rules-design.md` §B.4b.
  - Responsivitas Mobile (Drawer Overlay): Pada layar mobile (< 768px seperti Mobile S 320px), menu sidebar berubah menjadi Floating Drawer Overlay (`fixed inset-y-0 left-0 z-50`) dengan backdrop gelap (`fixed inset-0 bg-slate-900/40 z-40`), sehingga konten utama TIDAK bergeser/tertekan/menyempit. Mengetuk backdrop atau link menu otomatis menutup drawer.
  - Perapian Header Profil: Menata ulang tampilan profil user dan tombol Keluar di header menjadi horizontal sejajar yang rapi (avatar inisial, nama pengguna, badge role pill `indigo-50`, dan tombol Keluar dengan ikon `LogOut`), mengeliminasi tampilan tumpuk atas-bawah yang sempit. Di mobile (< 640px), teks nama dipotong rapi (`truncate`) dan label tombol disederhanakan tanpa merusak layout.
  - Perapian UI Form Login: Mengeliminasi divider bar `CardFooter` yang memotong tombol "Masuk", memposisikan tombol submit langsung di dalam `CardContent`, menambahkan icon branding `BookOpen`, menyempurnakan rounded container & soft shadow, serta menyembunyikan tombol "Masuk" di header navbar saat user sudah berada di halaman login.
  - Script Seed: Menambahkan `backend/cmd/seed/main.go` untuk mempermudah inisialisasi akun pengujian (`admin`, `asesor1`, `siswa1`) dengan password default `password123`.
  - Web UI: Komponen `input`, `label`, dan `card` di-generate murni menggunakan shadcn CLI (`npx shadcn@latest add input label card -y`) sesuai aturan di `05-tech-stack.md` §2.
  - Web Routing & Guard: Menggunakan TanStack Router dengan code-based route definition. `ProtectedRoute` memeriksa status login dan kecocokan role, dengan fallback 403 jika role tidak sesuai dan redirect ke `/login` jika unauthenticated. Auto-redirect setelah login disesuaikan dengan role (`admin` -> `/admin/users`, `asesor` -> `/admin/classes`, `siswa` -> `/student/exams`).
  - Web State: `AuthContext` memvalidasi token tersimpan saat app startup via `GET /auth/me` dan otomatis membersihkan token jika kedaluwarsa/tidak valid.
  - Mobile Stack: Menggunakan Jetpack DataStore Preferences (`TokenManager`) untuk menyimpan token JWT dan info user, Retrofit + Gson + OkHttp dengan `AuthInterceptor` untuk otomatis menyematkan header `Authorization: Bearer <token>`.
  - Mobile Versi Pustaka: Menyesuaikan versi pustaka Android di `libs.versions.toml` agar sepenuhnya kompatibel dengan `compileSdk 34` dan `agp 8.5.2` (`core-ktx 1.13.1`, `activity-compose 1.9.0`, `lifecycle 2.8.4`).
  - Mobile UX: Seluruh tombol interaktif dan input field di `LoginScreen` menerapkan area sentuh minimal 48dp sesuai persyaratan `04-rules-design.md` §A.1.
  - Mobile Networking Dev: Untuk pengujian di device fisik via kabel data USB saat debugging dari Android Studio, menggunakan ADB Reverse Tunneling (`adb reverse tcp:8080 tcp:8080`) dengan `BASE_URL = http://127.0.0.1:8080/api/v1/` guna mengeliminasi kendala Windows Defender Firewall dan AP Isolation pada router Wi-Fi.
- **Deviasi dari dokumen:** tidak ada.
- **Definition of Done:**
  - [x] Web login screen terhubung ke endpoint `POST /auth/login` via React Query.
  - [x] Web role-based route guard melindungi rute admin/asesor vs siswa.
  - [x] Web logout membersihkan token dan mengarahkan kembali ke `/login`.
  - [x] Mobile login screen terhubung ke `ApiService.login`.
  - [x] Mobile menyimpan JWT dan role di Jetpack DataStore.
  - [x] Mobile tombol/area sentuh memenuhi standar aksesibilitas minimum 48dp.
  - [x] Web build (`npm run build`) dan mobile build (`gradlew compileDebugKotlin`) sukses 100%.
- **Test:**
  - Web: TypeScript check dan Vite build produksi (`npm run build`) lulus tanpa error.
  - Mobile: Android Gradle compilation (`compileDebugKotlin`) lulus tanpa error/warning.
- **TODO lanjutan:**
  - Implementasi CRUD Pengguna (F-02) di web admin dan daftar kelas (F-03).


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
  - `docs/07-progress-log.md`
- **Keputusan implementasi:**
  - Menambahkan library `github.com/golang-jwt/jwt/v5` ke `go.mod` untuk token JWT (HMAC-SHA256) dengan klaim `user_id` dan `role`, serta masa aktif default 24 jam.
  - Rate limiting pada endpoint login diimplementasikan menggunakan in-memory sliding bucket berbasis client IP (`internal/middleware/ratelimit.go`, 5 request/menit, response 429 Too Many Requests) tanpa menambah dependency eksternal.
  - Struct `User` di `model.go` diberi tag `json:"-"` pada `PasswordHash` guna mencegah kebocoran hash password di respons JSON API.
  - Endpoint `GET /auth/me` mengembalikan format `{ "user": { ... } }` yang konsisten dengan respons `POST /auth/login`.
  - Logout bersifat stateless di sisi klien (klien menghapus JWT dari penyimpanan lokal), dengan expiry token di server.
  - Verifikasi migrasi: file `backend/migrations/000001_create_users_table.up.sql` sudah mencakup seluruh kolom di `03-architecture.md` §1 (`id`, `name`, `username`, `password_hash`, `role`, `created_by`, `is_active`), sehingga tidak memerlukan migrasi SQL baru.
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
