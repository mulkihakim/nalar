# 04. Rules & Design — Non-Fungsional dan Sistem Desain

> File ini memuat aturan non-fungsional, asumsi produk, dan panduan visual / sistem desain. Ini rujukan default: kalau ada keputusan teknis/visual yang tidak eksplisit disebut di file lain, cek dulu di sini sebelum agent coding berasumsi.

---

## Bagian A — Rules (Non-Fungsional & Asumsi Produk)

### A.1 Persyaratan non-fungsional

- **Keamanan:** password di-hash (bcrypt/argon2), token dengan claim role, otorisasi diperiksa di server untuk setiap endpoint (termasuk kepemilikan data asesor), rate limiting pada login.
- **Privasi:** data seminimal mungkin. Tampilkan pemberitahuan bahwa pilihan siswa dicatat, dan identitas hanya terlihat oleh admin/asesor.
- **Keandalan:** drop dan confirm tidak boleh hilang saat koneksi putus sesaat (retry/antrian sederhana di klien).
- **UX mobile:** interaksi utama adalah **ketuk slot → pilih opsi di bottom sheet** (bukan drag and drop; lihat `02-flows.md` §1.4b), untuk menghindari kerumitan implementasi drag-and-drop custom di Android serta keterbatasan lebar layar untuk teks ground/warrant yang panjang. Area sentuh minimal 48dp.
- **Kualitas:** dokumentasi API (OpenAPI), README dengan diagram arsitektur dan screenshot, seed data yang bisa dijalankan dengan satu perintah.
- **Tanpa Docker/deploy:** semua bisa dijalankan lokal dengan setup minimal.

### A.2 Asumsi (bisa diubah — konfirmasi ke pemilik proyek kalau agent ragu)

- **A1.** Kelompok pembanding = semua siswa dengan sesi pada ujian yang sama (bukan per kelas), karena satu ujian bisa diberikan ke beberapa kelas atau siswa individu.
- **A2.** Semua 4 opsi per jenis ditampilkan (4 ground dan 4 warrant). Jumlah tampil bisa dijadikan pengaturan di masa mendatang bila diperlukan.
- **A3.** Mode bantuan: setelah Confirm, siswa melihat bagian mana (ground/warrant) yang benar dan yang salah. Alternatif: opsi salah ditandai sejak awal.
- **A4.** Mode dipilih siswa setiap kali mengakses ujian. Untuk sesi berjalan, mode lama tetap dipakai kecuali siswa memilih mode lain dan mengonfirmasinya secara eksplisit (lihat `02-flows.md` §1 dan `03-architecture.md` §3).
- **A5.** Asesor hanya melihat kelas, materi, ujian, dan siswa yang berada di kelasnya.
- **A6.** Melepas opsi kembali ke daftar (web) atau menutup bottom sheet tanpa memilih (mobile) tidak dicatat. Hanya opsi yang benar-benar mengisi slot yang dicatat sebagai drop.
- **A7.** Tidak ada nilai numerik. "Hasil" berarti kemajuan penyelesaian dan jumlah percobaan, per sesi (percobaan).
- **A8.** Materi harus punya argumen lebih dari `arguments_per_session` agar pengacakan bermakna; jika materi punya argumen ≤ `arguments_per_session`, semua argumen dipakai. Pemilihan argumen tiap sesi baru diacak independen (boleh terulang dari sesi sebelumnya, tidak dijamin argumen berbeda).
- **A9.** Perhitungan X/Y/A/B pada analitik sosial (`02-flows.md` §2) menggabungkan attempt log dari **semua sesi** siswa pada ujian itu (bukan hanya sesi aktif), karena sifatnya latihan berulang. Kalau maksudnya per sesi saja, ini perlu diubah — beri tahu pemilik proyek.

---

## Bagian B — Design System (Palet, Tipografi, Layout)

Prasyarat: `05-tech-stack.md` (library terkunci: shadcn/ui + Tailwind untuk web, Jetpack Compose untuk mobile).

### B.1 Prinsip

- Satu aksen warna brand dipakai konsisten, sisanya netral (abu-slate) — hindari banyak warna cerah bersaing sekaligus.
- Web dan mobile berbagi palet dan filosofi yang sama, tapi komponen native masing-masing platform (jangan meniru mentah komponen iOS/Material di web atau sebaliknya).
- Karena ini "latihan belajar" (`01-product.md` §A.1), nuansa visual sebaiknya tenang dan fokus-baca (bukan gamifikasi/ramai), konsisten dengan non-tujuan produk (tidak ada skor/gamifikasi).

### B.2 Palet warna (token)

| Token | Hex | Kelas Tailwind terdekat | Penggunaan |
|---|---|---|---|
| `primary` | `#4F46E5` | `indigo-600` | Logo, tombol utama (Confirm, Masuk, Simpan), item nav aktif |
| `primary-hover` | `#4338CA` | `indigo-700` | Hover/active state tombol primary |
| `background` | `#F8FAFC` | `slate-50` | Latar halaman (di luar card/sidebar) |
| `surface` | `#FFFFFF` | `white` | Card, sidebar, header, bottom sheet |
| `border` | `#E2E8F0` | `slate-200` | Border card, input, divider |
| `text-primary` | `#0F172A` | `slate-900` | Judul, teks utama |
| `text-secondary` | `#64748B` | `slate-500` | Deskripsi, label, teks pembantu |
| `success` | `#059669` | `emerald-600` | Argumen selesai, ujian aktif, "tanpa salah" |
| `danger` | `#DC2626` | `red-600` | Error, ujian nonaktif, validasi gagal |
| `warning` | `#D97706` | `amber-600` | Indikator mode Bantuan, peringatan ganti mode |

**Cara pakai di web (shadcn/Tailwind):** cukup pakai kelas Tailwind di atas langsung (`bg-indigo-600`, `text-slate-500`, dst) — tidak perlu bikin custom CSS variable dulu, ini paling sederhana untuk pemula. Kalau nanti F-17 (mode gelap) dikerjakan, baru pindahkan token ini ke CSS variable shadcn (`--primary` dkk di `index.css`) supaya bisa di-swap per tema.

**Cara pakai di mobile (Compose):** definisikan sekali di `ui/theme/Color.kt`, lalu pasang ke `MaterialTheme.colorScheme` (map `primary` → `#4F46E5`, `error` → `#DC2626`, dst) agar semua komponen Material3 otomatis konsisten tanpa hardcode warna di tiap screen.

### B.3 Tipografi

- **Font:** Inter (sans-serif) untuk web maupun mobile. Alasan: gratis, sangat mudah dibaca di ukuran kecil (penting untuk teks *ground*/*warrant* yang panjang), dan tersedia langsung lewat Google Fonts (web) atau bisa dibundel sebagai font resource (Android) tanpa lisensi tambahan.
- **Skala** (dipakai konsisten di kedua platform):

| Peran | Ukuran | Weight |
|---|---|---|
| Judul halaman (mis. "Masuk ke Nalar", "Daftar Ujian") | 24px | Bold (700) |
| Judul card/section | 18px | Semibold (600) |
| Body / claim / ground / warrant | 16px | Regular (400) |
| Label input, caption, metadata | 13px | Medium (500) |

Hindari all-caps untuk label (lihat catatan di skill `frontend-design`) — pakai sentence case biasa, misal "Username" bukan "USERNAME".

### B.4 Layout Web

Web punya dua kebutuhan berbeda (lihat peran di `01-product.md` §A.3), jadi dua pola layout:

#### B.4a. Halaman auth (login) — standalone, tanpa chrome aplikasi

```
┌─────────────────────────────┐
│                             │
│         [Logo Nalar]        │
│      ┌───────────────┐      │
│      │   Card login   │      │
│      │  (centered)    │      │
│      └───────────────┘      │
│                             │
└─────────────────────────────┘
```
Background `bg-slate-50` penuh, card `bg-white` di tengah (`items-center justify-center`, `max-w-sm`), tanpa header/sidebar.

#### B.4b. Panel Admin/Asesor — header + sidebar

Dipakai untuk F-02 s/d F-05, F-13, F-16 (semua fitur `Platform: Web` yang menyasar admin/asesor).

```
┌───────────────────────────────────────────┐
│ Header: [Logo]              [Nama · Role · Logout] │
├──────────┬──────────────────────────────────┤
│ Sidebar  │  Konten halaman aktif             │
│ Dashboard│  (max-width, padding konsisten)   │
│ Kelas    │                                    │
│ Materi   │                                    │
│ Ujian    │                                    │
│ Siswa*   │                                    │
│ Hasil&Log│                                    │
└──────────┴──────────────────────────────────┘
```
- Header fixed, tinggi 64px, `border-b border-slate-200`, latar putih.
- Sidebar fixed, lebar 240px, item menu aktif diberi `bg-indigo-50 text-indigo-600` (bukan warna solid penuh — supaya tidak terlalu ramai).
- Menu "Kelola Asesor" (item bertanda \*) hanya tampil untuk role Admin, sesuai tabel hak akses `01-product.md` §A.3.
- Konten utama: `max-w-6xl mx-auto px-6 py-8`, supaya tabel/form tidak melebar penuh di monitor lebar.
- Sidebar collapsible ke ikon-saja di lebar sempit (opsional, boleh dilewati untuk MVP demo).

#### B.4c. Siswa via web (F-14, Should)

Karena scope siswa di web sempit (cuma: lihat ujian → baca materi → susun argumen drag-drop `02-flows.md` §1.4a → lihat hasil), **tidak pakai sidebar**. Cukup header sederhana (logo kiri, nama siswa + logout kanan), konten di bawahnya. Layar argumen mengikuti tata letak web (drag-drop) sesuai `02-flows.md` §1.4a — jangan disamakan dengan panel admin.

### B.5 Layout Mobile (Android)

Scope siswa di mobile juga sempit (F-06 s/d F-12): daftar ujian, pengerjaan, hasil. Pola yang disarankan:

- **Bottom navigation** (Material3 `NavigationBar`) dengan 2–3 tujuan: **Ujian** (daftar + entry point pengerjaan), **Hasil** (riwayat sesi, F-12), **Profil** (info akun + logout). Ini pola standar Android untuk app dengan sedikit top-level destination — lebih pas daripada drawer/hamburger untuk kasus sesempit ini.
- **Top app bar** di setiap screen: judul screen, tombol back saat masuk ke sub-halaman (mis. dari daftar ujian ke layar argumen).
- **Daftar ujian:** list `Card` (judul ujian, badge status jika relevan, tap untuk mulai/lanjut sesi).
- **Layar argumen (F-08):** sesuai `02-flows.md` §1.4b — dua slot card kosong tersusun vertikal. Tambahan visual:
  - Slot kosong: border `slate-200`, placeholder teks abu (`text-secondary`).
  - Slot terisi: border `primary` (indigo) — **hanya menandai "terisi", bukan "benar/salah"**, supaya tidak melanggar aturan penerimaan #4 di `06-tasks.md` (mode standar tidak boleh membocorkan `is_correct`).
  - Tombol Confirm: `primary`, nonaktif (abu, disabled) sampai kedua slot terisi.
- **Bottom sheet pilihan opsi:** Material3 `ModalBottomSheet`, list 4 opsi dengan area sentuh minimal 48dp (sesuai §A.1 di atas — UX mobile).
- **Indikator progres:** teks kecil "Argumen X dari N" di top app bar atau tepat di bawahnya, bukan progress bar besar (produk ini bukan gamifikasi, sesuai non-tujuan di `01-product.md` §A.2).

### B.6 Penerapan & maintenance

- Web: jangan tulis hex manual berulang di banyak file — pakai kelas Tailwind semantik dari tabel §B.2 secara konsisten.
- Mobile: satu sumber kebenaran di `ui/theme/Color.kt` + `Theme.kt`, semua screen ambil warna dari `MaterialTheme.colorScheme`, jangan hardcode `Color(0xFF...)` di composable individual.
- Kalau ada kebutuhan token baru (misal warna khusus untuk badge "aktif/nonaktif" ujian) yang tidak ada di tabel §B.2, tambahkan ke tabel ini dulu sebelum dipakai di kode, supaya tidak drift antar komponen.
