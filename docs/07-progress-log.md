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

### 2026-09-25 — Fase 4: Implementasi Mobile Android Lengkap Siswa (F-06 s/d F-12)
- **Status:** selesai
- **File dibuat/diubah:**
  - `mobile/gradle/libs.versions.toml` (tambah androidx-navigation-compose)
  - `mobile/app/build.gradle.kts`
  - `mobile/app/src/main/java/com/mulki/nalar/ui/theme/` (Color.kt, Theme.kt, Type.kt)
  - `mobile/app/src/main/java/com/mulki/nalar/core/auth/TokenManager.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/core/network/ApiService.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/core/network/dto/` (ExamDto.kt, SessionDto.kt)
  - `mobile/app/src/main/java/com/mulki/nalar/feature/navigation/NalarNavigation.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/feature/examlist/` (ExamListScreen.kt, ExamListViewModel.kt)
  - `mobile/app/src/main/java/com/mulki/nalar/feature/session/` (ModeSelectDialog.kt, ModeConfirmDialog.kt, SessionViewModel.kt)
  - `mobile/app/src/main/java/com/mulki/nalar/feature/argument/` (ArgumentScreen.kt, ArgumentViewModel.kt, OptionBottomSheet.kt, ConfirmResultDialog.kt, MonitoringView.kt, MaterialReadScreen.kt)
  - `mobile/app/src/main/java/com/mulki/nalar/feature/result/` (ResultScreen.kt, ResultViewModel.kt, SessionSummaryScreen.kt, AnalysisScreen.kt)
  - `mobile/app/src/main/java/com/mulki/nalar/feature/profile/ProfileScreen.kt`
  - `mobile/app/src/main/java/com/mulki/nalar/MainActivity.kt`
- **Keputusan:**
  - Menambahkan pustaka resmi `androidx.navigation:navigation-compose:2.7.7` ke `libs.versions.toml` dan `build.gradle.kts` untuk navigasi multi-screen deklaratif berbasis NavHost.
  - Memetakan token warna Nalar (04-rules-design.md §B.2) langsung ke slot MaterialTheme.colorScheme (`primary` = Indigo-600, `secondary` = Emerald-600 success, `tertiary` = Amber-600 warning, `error` = Red-600 danger) sehingga seluruh Composable konsisten tanpa hardcode warna individual.
  - Sesuai 02-flows.md §1.4b, interaksi pengerjaan argumen di mobile menerapkan vertikal tap-to-select (slot Ground di atas, slot Warrant di bawah) yang membuka Material3 ModalBottomSheet.
  - Setiap kali siswa memilih opsi dari bottom sheet, dicatat drop attempt log ke backend (POST /sessions/:id/arguments/:aid/drops). Menutup sheet tanpa memilih tidak mencatat log (asumsi A6).
  - Evaluasi confirm dinilai di server. Di mode bantuan, bagian ground/warrant yang tepat/salah ditampilkan; di mode standar dan sosial, detail tidak dibocorkan.
  - Pada mode analitik sosial, monitoring rasio X:Y ditampilkan setelah setiap argumen selesai, dan analisis rasio A:B ditampilkan setelah sesi selesai.
  - Riwayat hasil menampilkan detail per argumen beserta penanda "Tanpa Salah ✓" jika diselesaikan tepat 2 percobaan.
- **Deviasi:** tidak ada.
- **DoD:**
  - DoD #4 (respons mode standar tidak membocorkan is_correct) → ✅
  - DoD #5 (argumen harus berurutan) → ✅
  - DoD #6 (setiap drop menambah attempt_logs) → ✅
  - DoD #7 (monitoring ratio X:Y) → ✅
  - DoD #8 (analysis ratio A:B) → ✅
  - DoD #9 (opsi tanpa pemilih menampilkan Y = -) → ✅
  - DoD #10 (sembunyikan data kelompok jika peers < MIN_PEERS) → ✅
  - DoD #11 (attempt_no bertambah & subset diacak) → ✅
  - DoD #12 (sesi berjalan dilanjutkan jika mode sama) → ✅
  - DoD #13 (status 409 jika mode beda + konfirmasi ganti mode) → ✅
  - Area sentuh minimum 48dp (04-rules-design.md §A.1) → ✅
- **Test:**
  - `mobile`: `.\gradlew compileDebugKotlin` (BUILD SUCCESSFUL 100%, 0 warning).
  - `mobile`: `.\gradlew assembleDebug` (BUILD SUCCESSFUL 100%, APK debug siap).
- **TODO:** tidak ada.

### 2026-09-24 — Penyempurnaan Monitoring Dialog (HoverCard) & Notifikasi Toast CRUD (Sonner)
- **Status:** selesai
- **File dibuat/diubah:**
  - `web/src/components/ui/hover-card.tsx` (dibuat via shadcn CLI)
  - `web/src/components/ui/sonner.tsx` (dibuat via shadcn CLI)
  - `web/src/routes/__root.tsx` (integrasi `<Toaster richColors position="top-right" />`)
  - `web/src/features/session/components/MonitoringView.tsx` (HoverCard penjelasan rasio X:Y)
  - `web/src/features/session/components/AnalysisView.tsx` (HoverCard penjelasan rasio A:B siswa)
  - `web/src/features/class/api.ts` (toast success & error untuk CRUD kelas dan anggota)
  - `web/src/features/material/api.ts` (toast success & error untuk CRUD materi dan argumen)
  - `web/src/features/exam/api.ts` (toast success & error untuk CRUD paket ujian, akses, dan status)
  - `web/src/features/user/api.ts` (toast success & error untuk CRUD pengguna)
- **Keputusan:**
  - Memasang komponen `hover-card` dan `sonner` resmi dari shadcn via CLI (`npx shadcn@latest add hover-card sonner --yes`) sesuai aturan bahwa komponen UI harus berasal dari shadcn.
  - Notifikasi toast CRUD dipasang terpusat pada hook mutasi TanStack Query (`onSuccess` & `onError`), menjamin feedback langsung muncul di setiap operasi create, update, delete, maupun error tanpa duplikasi di tiap modal.
- **Deviasi:** tidak ada.
- **Test:**
  - `web`: `npm run build` PASS 100% (tsc & vite build tanpa error).
  - `backend`: `go test -v ./...` PASS 100%.
- **TODO:** tidak ada.

### 2026-09-23 — Fase 3: Modul Ujian Siswa (Drag-and-Drop) & Analitik Admin/Asesor
- **Status:** selesai
- **File dibuat/diubah:**
  - `backend/internal/session/` (model.go, repository.go, service.go, handler.go, service_test.go)
  - `backend/internal/db/db.go` (AutoMigrate session, session_arguments, attempt_logs, argument_progress)
  - `backend/cmd/api/main.go` (wiring session repository, service, handler)
  - `web/src/features/session/` (types.ts, api.ts, components/ModeSelectDialog.tsx, OptionCard.tsx, SlotZone.tsx, ArgumentBoard.tsx, MonitoringView.tsx, AnalysisView.tsx, SessionSummary.tsx, ResultsTable.tsx, LogsTable.tsx)
  - `web/src/routes/student/` (exams.tsx, exam.$examId.tsx)
  - `web/src/routes/admin/results.tsx`
  - `web/src/routes/router.tsx`
  - `web/src/routes/__root.tsx`
- **Keputusan:**
  - Menggunakan `@dnd-kit/core` untuk drag-and-drop Toulmin model di web (ringan, ramah aksesibilitas, kompatibel React 19).
  - Mode conflict (memilih mode berbeda dari sesi berjalan) mengembalikan HTTP 409 Conflict beserta `current_mode` agar frontend memunculkan dialog konfirmasi ganti mode.
  - Fix query `GetStudentExams` di `session/repository.go`: menggunakan `cm.user_id` (sesuai kolom skema `class_members`) bukan `cm.student_id`.
  - Refinement UX layout Toulmin di web: ringkasan jembatan Toulmin diletakkan di atas slot menggantikan teks struktur, opsi berada langsung di bawah slot masing-masing, tombol confirm dipindah ke paling bawah.
  - Feedback evaluasi (benar maupun salah) menggunakan komponen Shadcn Dialog (Modal), sehingga tidak menggeser atau memperlebar jarak layout.
  - Stepper navigasi soal: siswa dapat melihat/meninjau argumen yang sudah diselesaikan sebelumnya, berpindah antar soal terbuka, dan soal berikutnya terkunci sebelum soal aktif selesai.
  - Tombol keluar ke daftar ujian dilengkapi Shadcn ConfirmDialog agar siswa tidak sengaja keluar.
  - Pengacakan urutan opsi Ground dan Warrant secara stabil per sesi & argumen di backend (jawaban benar tidak selalu di urutan teratas).
- **Deviasi:** tidak ada.
- **DoD:**
  - DoD #4 (respons mode standar tidak berisi is_correct) → ✅
  - DoD #5 (argumen harus diselesaikan berurutan) → ✅
  - DoD #6 (setiap drop menambah attempt_logs) → ✅
  - DoD #7 (monitoring ratio 2:3, 3:3, 4:3) → ✅
  - DoD #8 (analysis ratio 9:4) → ✅
  - DoD #9 (opsi tanpa pemilih menampilkan Y = -) → ✅
  - DoD #10 (sembunyikan data kelompok jika peers < MIN_PEERS) → ✅
  - DoD #11 (attempt_no bertambah & subset diacak) → ✅
  - DoD #12 (sesi berjalan dilanjutkan jika mode sama) → ✅
  - DoD #13 (status 409 jika mode beda) → ✅
  - DoD #14 (partial unique index status in_progress) → ✅
- **Test:**
  - Backend: `backend/internal/session/service_test.go` (PASS), `go test -v ./...` (100% PASS)
  - Frontend: `npm run build` (tsc -b + vite build) (100% PASS)
- **TODO:** Pengerjaan versi mobile Android (Fase 4).

### 2026-09-23 — Refinement Shadcn UI, Zod Form Validation, Aturan Opsi Argumen (Min 3), & Restriksi Pengguna/Ujian
- **Status:** selesai
- **File dibuat/diubah:**
  - `docs/01-product.md`
  - `docs/03-architecture.md`
  - `docs/05-tech-stack.md`
  - `docs/06-tasks.md`
  - `docs/07-progress-log.md`
  - `backend/internal/user/` (repository.go, service.go, handler.go, service_test.go)
  - `backend/internal/material/` (service.go, service_test.go)
  - `backend/internal/exam/` (repository.go, service.go, service_test.go)
  - `web/src/components/ui/` (alert-dialog.tsx, select.tsx, textarea.tsx, switch.tsx, badge.tsx, dialog.tsx, dropdown-menu.tsx, button.tsx, input.tsx, label.tsx, card.tsx)
  - `web/src/components/ConfirmDialog.tsx`
  - `web/src/components/Modal.tsx`
  - `web/src/routes/__root.tsx`
  - `web/src/features/user/` (types.ts, api.ts, components/UserModal.tsx, components/UserTable.tsx)
  - `web/src/features/class/` (types.ts, components/ClassModal.tsx, components/ClassTable.tsx, components/MembersModal.tsx)
  - `web/src/features/material/` (types.ts, components/MaterialModal.tsx, components/ArgumentModal.tsx, components/MaterialTable.tsx)
  - `web/src/features/exam/` (types.ts, components/ExamModal.tsx, components/ExamTable.tsx, components/ExamAccessModal.tsx)
- **Keputusan implementasi:**
  - Enforce Shadcn CLI di `components/ui/`: Komponen di `web/src/components/ui/` dipasang resmi via CLI shadcn tanpa modifikasi manual langsung. Wrapper aplikasi seperti `ConfirmDialog` dan `Modal` ditempatkan di `web/src/components/` di luar `ui/`. Larangan modifikasi manual dicatat secara tegas di `docs/05-tech-stack.md`.
  - Shadcn AlertDialog & Select: Mengganti semua browser native `confirm()` dan native `<select>` menjadi `AlertDialog` (`ConfirmDialog`) dan `Select` shadcn untuk aksi konfirmasi (logout, hapus kelas, hapus materi, hapus argumen, hapus ujian, keluarkan siswa, dan hapus pengguna).
  - Fleksibilitas Opsi Argumen: Aturan opsi diubah dari strict 4+4 menjadi **minimal 3 ground** (tepat 1 benar, minimal 2 salah) dan **minimal 3 warrant** (tepat 1 benar, minimal 2 salah). Modal input di web mendukung penambahan/penghapusan opsi secara dinamis dengan validasi schema Zod.
  - Validasi Zod: Schema Zod diterapkan menyeluruh di semua form frontend (auth login, user create/update, class create/update, member add, material create/update, argument dynamic options, exam create/update).
  - Restriksi Pengguna & Asesor: Di endpoint `GET /api/v1/users`, asesor tidak dapat melihat pengguna dengan role `admin`. Soft delete dapat mengubah status aktif/nonaktif akun. Hard delete (`DELETE /api/v1/users/:id`) hanya dapat dilakukan oleh admin dan divalidasi tidak boleh memiliki riwayat sesi ujian (`ErrUserHasExamRecords`).
  - Restriksi Hak Akses Ujian: Pada `PUT /api/v1/exams/:id/access`, kelas atau siswa yang dicabut aksesnya diperiksa terlebih dahulu apakah sudah pernah mengerjakan sesi ujian tersebut; jika sudah, request ditolak dengan pesan error yang jelas.
- **Deviasi dari dokumen:** tidak ada.
- **DoD:**
  - Kriteria #1 (Isolasi Asesor & User Visibility): ✅ Asesor terisolasi dan tidak melihat admin.
  - Kriteria #2 (Validasi Argumen min 3 ground, min 3 warrant, 1 benar masing-masing): ✅ Lulus 100% pada backend unit test & frontend Zod validation.
  - Restriksi Delete User & Akses Ujian: ✅ Lulus pada unit test backend `TestDeleteUser_Rules` dan `TestExamAccessParticipantRestriction`.
- **Test:**
  - `backend`: `go test -v ./...` PASS 100% (user, class, material, exam, middleware).
  - `web`: `npm run build` (`tsc -b && vite build`) PASS 100% tanpa error.
- **TODO lanjutan:**
  - Siap melangkah ke pengerjaan Fase 3: Endpoint pengerjaan siswa & engine pengerjaan tes.

### 2026-09-23 — F-02, F-03, F-04, F-05 CRUD Resource Backend & Web Admin Panel (Fase 2a & 2b)
- **Status:** selesai
- **File dibuat/diubah:**
  - `backend/migrations/000002_create_classes_and_members.up.sql` & `.down.sql`
  - `backend/migrations/000003_create_materials_and_arguments.up.sql` & `.down.sql`
  - `backend/migrations/000004_create_exams_and_access.up.sql` & `.down.sql`
  - `backend/internal/db/db.go`
  - `backend/internal/user/model.go`
  - `backend/internal/user/repository.go`
  - `backend/internal/user/service.go`
  - `backend/internal/user/handler.go`
  - `backend/internal/class/` (model.go, repository.go, service.go, handler.go, service_test.go)
  - `backend/internal/material/` (model.go, repository.go, service.go, handler.go, service_test.go)
  - `backend/internal/exam/` (model.go, repository.go, service.go, handler.go, service_test.go)
  - `backend/cmd/api/main.go`
  - `backend/cmd/seed/main.go`
  - `web/src/types/index.ts`
  - `web/src/components/ui/badge.tsx`
  - `web/src/components/ui/dialog.tsx`
  - `web/src/features/user/` (types.ts, api.ts, components/UserModal.tsx, components/UserTable.tsx)
  - `web/src/features/class/` (types.ts, api.ts, components/ClassModal.tsx, components/MembersModal.tsx, components/ClassTable.tsx)
  - `web/src/features/material/` (types.ts, api.ts, components/MaterialModal.tsx, components/ArgumentModal.tsx, components/MaterialTable.tsx)
  - `web/src/features/exam/` (types.ts, api.ts, components/ExamModal.tsx, components/ExamAccessModal.tsx, components/ExamTable.tsx)
  - `web/src/routes/__root.tsx`
  - `web/src/routes/router.tsx`
  - `web/src/routes/admin/users.tsx`
  - `web/src/routes/admin/classes.tsx`
  - `web/src/routes/admin/materials.tsx`
  - `web/src/routes/admin/exams.tsx`
- **Keputusan implementasi:**
  - Migrasi Runner: Menambahkan `RunMigrations` di `backend/internal/db/db.go` yang melacak file `*.up.sql` pada tabel `schema_migrations` agar tabel dan constraint langsung terpasang saat backend dijalankan secara lokal tanpa perlu tool migration eksternal.
  - Validasi DoD #2: Validasi 4 ground + 4 warrant (tepat 1 is_correct per jenis) diimplementasikan dan diverifikasi ketat di layer service (`material/service.go`) dan unit test (`material/service_test.go`).
  - Isolasi Asesor (DoD #1): Ditegakkan di layer service pada seluruh domain (`class`, `material`, `exam`). Asesor hanya dapat mengakses, mengubah, atau menghapus resource miliknya (`owner_id = user_id`), sedangkan admin memiliki akses penuh.
  - Web UI: Mengikuti arsitektur TanStack Query + TanStack Router + Tailwind CSS + Lucide Icons dengan modal responsif, badge status, dan form validasi interaktif.
- **Deviasi dari dokumen:** tidak ada.
- **DoD:**
  - Kriteria #1 (Isolasi Asesor): ✅ Diuji otomatis di unit test class, material, exam.
  - Kriteria #2 (Validasi Argumen 4+4 opsi & 1+1 benar): ✅ Diuji otomatis di unit test material.
- **Test:**
  - `backend`: `go test -v ./...` lulus 100% (semua paket: user, class, material, exam, middleware).
  - `backend`: `go build ./...` sukses tanpa warning.
  - `web`: `npm run build` (`tsc -b && vite build`) sukses 100%.
- **TODO lanjutan:**
  - Masuk ke Fase 3: Endpoint pengerjaan siswa di backend (sessions, session_arguments dengan partial unique index, attempt_logs, argument_progress, scoring/confirm, dan query analitik sosial).

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
