# 05. Tech Stack — Keputusan Terkunci (Jangan Diganti Tanpa Konfirmasi)

Prasyarat: baca `01-product.md` dan `03-architecture.md`.

> Dokumen ini mengunci library dan struktur folder yang **sudah terinstal** di proyek. Agent coding **tidak boleh** mengganti/menambah library inti (router, ORM, UI kit, dst) dengan alternatif lain tanpa bertanya ke pemilik proyek dulu — cukup pakai yang tercantum di sini.

## 1. Backend (Go)

| Kebutuhan | Library | Catatan |
|---|---|---|
| Router HTTP | `go-chi/chi/v5` | Idiomatic, kompatibel `net/http` standar — bukan framework "berat" seperti Gin/Fiber |
| ORM | `gorm.io/gorm` + `gorm.io/driver/postgres` | Dipakai untuk query CRUD, **bukan** untuk `AutoMigrate` (lihat §1.2) |
| Load `.env` | `joho/godotenv` | Dev lokal saja |
| Hash password | `golang.org/x/crypto/bcrypt` | Sesuai `04-rules-design.md` §A.1 |
| Auth token (JWT) | `github.com/golang-jwt/jwt/v5` | Sudah ditambahkan ke `go.mod` (lihat `07-progress-log.md` entri F-01) |

### 1.1 Struktur folder (wajib diikuti, sesuai yang sudah ada)

```
backend/
  cmd/api/main.go          # entrypoint: wiring db, router, semua modul domain
  internal/
    db/db.go               # koneksi gorm + runner migrasi
    user/                  # 1 folder = 1 domain/tabel utama
      handler.go           # layer HTTP: decode request, panggil service, encode response
      service.go           # business logic (validasi domain, aturan dari 02-flows.md/04-rules-design.md)
      repository.go        # satu-satunya layer yang boleh panggil gorm/db langsung
      model.go             # struct gorm untuk tabel
  migrations/
    000001_create_users_table.up.sql
    000001_create_users_table.down.sql
```

**Setiap domain baru** (`class`, `material`, `exam`, `session`, `attempt_log`, dst — lihat ERD di `03-architecture.md` §1) dibuat dengan pola folder yang **identik** dengan `internal/user`: `handler.go`, `service.go`, `repository.go`, `model.go`.

Aturan layer (tidak boleh dilanggar, ini yang membuat kode testable):
- `handler.go` **hanya** boleh memanggil `service.go`, tidak boleh menyentuh gorm.
- `service.go` **hanya** boleh memanggil `repository.go` domainnya sendiri (atau repository domain lain lewat interface bila perlu lintas domain), tidak boleh menyusun query gorm sendiri.
- `repository.go` satu-satunya tempat query gorm/SQL untuk domain itu.

### 1.2 Migrasi database

Constraint di `03-architecture.md` §2 (unique index, FK, check constraint) **harus** ditulis sebagai SQL murni di `migrations/*.up.sql` / `*.down.sql`, penomoran berurutan (`000001_`, `000002_`, ...). **Jangan** pakai `db.AutoMigrate()` gorm untuk membuat/mengubah tabel produksi — itu tidak bisa menegakkan partial unique index (butuh untuk kriteria penerimaan #14 di `06-tasks.md`) atau check constraint dengan rapi. `gorm` tetap dipakai untuk query, bukan untuk schema management.

### 1.3 Testing

Unit test ditaruh di folder domain yang sama, mengikuti konvensi Go (`service_test.go`, `repository_test.go`). Untuk repository yang butuh DB nyata (constraint unique/FK), gunakan Postgres nyata di test (bukan SQLite) karena kita bergantung pada fitur Postgres (partial unique index) — detail setup test DB didiskusikan saat mulai hari 1 (`06-tasks.md`), catat keputusannya di `07-progress-log.md`.

## 2. Web (React)

| Kebutuhan | Library | Catatan |
|---|---|---|
| UI Kit | **shadcn/ui + Tailwind CSS** | **Menggantikan Chakra UI.** Alasan: komponen di-*copy* ke repo (bukan dependency tertutup), gampang dikustomisasi, dan jadi standar de-facto di ekosistem React+Tailwind saat ini |
| Routing | `@tanstack/react-router` | Type-safe, sudah terpasang |
| Data fetching | `@tanstack/react-query` | **Wajib** dipakai untuk semua panggilan API ke backend (GET/POST dsb) — jangan `fetch` manual di `useEffect` |
| Form + validasi | `react-hook-form` + `zod` (resolver) | Skema zod idealnya mencerminkan validasi server (misal 4 ground + 4 warrant di `03-architecture.md`) supaya pesan error konsisten |

> **Aksi yang perlu dilakukan:** hapus `@chakra-ui/react` dari `package.json`, install `tailwindcss` + jalankan `npx shadcn@latest init`. Import komponen (`Button`, `Input`, dst) lewat shadcn generator, bukan lewat npm package biasa.

### 2.1 Struktur folder (wajib diikuti — belum ada di proyek, buat baru)

Prinsip sama dengan backend: **1 domain = 1 folder**, supaya agent yang baru baca kode langsung tahu di mana mencari logic exam/session/dst tanpa loncat-loncat.

```
web/src/
  routes/                      # tanstack-router: HANYA wiring halaman + layout, tidak ada logic fetch/form di sini
    __root.tsx
    login.tsx
    admin/
      users.tsx
      classes.tsx
      materials.tsx
      exams.tsx
      results.tsx
    student/                   # F-14 — siswa lewat web (drag-and-drop, §1.4a di 02-flows.md)
      exams.tsx
      exam.$examId.tsx

  features/                    # domain logic, mirror nama domain di 03-architecture.md
    auth/
      api.ts                   # semua useQuery/useMutation react-query untuk domain ini
      types.ts
    user/
      api.ts
      components/
      types.ts
    class/
    material/
    exam/
    session/                   # drop, confirm, hasil, analitik sosial

  components/
    ui/                        # HANYA hasil `npx shadcn add <component>` — jangan tulis manual
    layout/                    # Navbar, Sidebar, dst

  lib/
    api-client.ts              # wrapper fetch/axios: base URL + inject token dari auth
    query-client.ts            # instance QueryClient

  types/
    index.ts                   # tipe yang dipakai lintas domain (User, Role, dst)
```

Aturan layer (setara handler/service/repository di backend):
- `routes/*.tsx` **hanya** memanggil hook dari `features/<domain>/api.ts` dan merender `components/`, tidak boleh menulis `useQuery`/`fetch` langsung.
- `features/<domain>/api.ts` **satu-satunya** tempat panggilan API untuk domain itu — komponen tidak boleh `fetch` manual (sesuai keputusan wajib pakai `@tanstack/react-query` di atas).
- `components/ui/` murni hasil generator shadcn; kalau butuh varian custom, bikin komponen baru di `components/` yang *membungkus* `components/ui/`, jangan edit file hasil generate.

## 3. Mobile (Kotlin/Android)

> Stack sudah dikonfirmasi dan berjalan (lihat `07-progress-log.md` entri F-01): Jetpack Compose, Retrofit + Gson + OkHttp, Coroutines/StateFlow, Jetpack DataStore untuk token.

| Kebutuhan | Library | Alasan singkat |
|---|---|---|
| UI | Jetpack Compose | Standar modern Android, deklaratif, lebih sedikit boilerplate dibanding XML/View |
| Networking | Retrofit + Gson + OkHttp | Dipakai dengan `AuthInterceptor` untuk menyematkan header `Authorization: Bearer <token>` otomatis |
| Async/state | Coroutines + `StateFlow`, `ViewModel` | Idiomatic Android modern, cocok untuk pola "ketuk slot → bottom sheet" di `02-flows.md` §1.4b |
| Simpan token | Jetpack `DataStore` Preferences | Pengganti `SharedPreferences`, async dan type-safe |

### 3.1 Struktur folder (wajib diikuti)

`com.mulki.nalar` adalah application ID aktual di proyek ini (lihat `AndroidManifest`/`build.gradle.kts`).

```
app/src/main/java/com/mulki/nalar/
  MainActivity.kt

  core/
    network/
      ApiService.kt           # interface Retrofit — endpoint sesuai `03-architecture.md` §3
      dto/                     # data class request/response, 1 file per domain (ExamDto.kt, SessionDto.kt, dst)
      NetworkModule.kt         # provide Retrofit + OkHttp (base URL, auth interceptor)
    auth/
      TokenManager.kt          # baca/simpan JWT via DataStore

  feature/                    # 1 folder = 1 domain/fitur, mirror ID fitur di 01-product.md §B
    auth/
      LoginScreen.kt
      LoginViewModel.kt
    examlist/                 # F-06
      ExamListScreen.kt
      ExamListViewModel.kt
    session/                  # F-07, F-10 — pilih mode, dialog konfirmasi ganti mode
      ModeSelectDialog.kt
      SessionViewModel.kt
    argument/                 # F-08 — layar tap-to-select (lihat 02-flows.md §1.4b)
      ArgumentScreen.kt
      ArgumentViewModel.kt
      OptionBottomSheet.kt
    result/                   # F-12
      ResultScreen.kt
      ResultViewModel.kt
    analytics/                # monitoring & analysis analitik sosial

  ui/
    theme/                     # Color.kt, Type.kt, Theme.kt (Compose)
    components/                # composable generik (dipakai lintas fitur)
```

Aturan layer (versi ringan, karena scope 2–3 hari — **boleh dinaikkan ke Repository terpisah kalau nanti butuh unit test tanpa jaringan**, catat di `07-progress-log.md` kalau berubah):
- `*Screen.kt` (Compose) **hanya** merender state dari `*ViewModel.kt`, tidak boleh panggil `ApiService` langsung.
- `*ViewModel.kt` memanggil `ApiService` dan mengekspos `StateFlow` ke Screen — 1 ViewModel per fitur di atas, jangan digabung.
- `dto/` harus mengikuti bentuk JSON di `03-architecture.md` §3 apa adanya (nama field sama, tidak diterjemahkan), supaya serialisasi tidak meleset.

## 4. Aturan untuk agent coding

- Sebelum menulis kode untuk fitur apa pun, cek dulu file ini untuk memastikan library/struktur folder yang dipakai konsisten dengan yang sudah terinstal.
- Kalau kebutuhan baru muncul dan **tidak ada** di tabel di atas (misal butuh library CSV untuk F-16, atau library UUID), pilih **satu** kandidat paling standar, sebutkan alasannya, dan catat di `07-progress-log.md` sebagai "keputusan implementasi" — jangan diam-diam tambah banyak dependency sekaligus.
