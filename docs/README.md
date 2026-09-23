# Dokumentasi Proyek Nalar — Panduan untuk Agent Coding

PRD asli (`PRD.md`) dipecah jadi beberapa file fokus di folder ini, supaya konteks yang perlu dibaca untuk satu tugas tidak terlalu panjang dan tidak bercampur dengan hal yang tidak relevan.

## Stack

- **Backend:** Go
- **Web:** React
- **Mobile:** Kotlin (Android)

## Peta dokumen

| File | Isi | Baca saat... |
|---|---|---|
| [`01-overview.md`](./01-overview.md) | Ringkasan produk, tujuan/non-tujuan, peran & hak akses, glosarium istilah domain (Material, Argument, Session, Attempt log, dst) | Konteks dasar — sebelum mengerjakan fitur apa pun |
| [`02-requirements.md`](./02-requirements.md) | Daftar fitur (F-01 s/d F-17) dengan prioritas & platform, serta definisi 3 mode belajar | Menentukan scope satu fitur, atau cek fitur ini prioritas Must/Should/Could |
| [`03-flows.md`](./03-flows.md) | Alur langkah-demi-langkah pengerjaan siswa (mulai/lanjut sesi, layar argumen versi Web vs Mobile, drop/confirm), aturan analitik sosial, struktur hasil | Implementasi endpoint pengerjaan atau layar F-06 s/d F-12 |
| [`04-architecture.md`](./04-architecture.md) | Model data (ERD + skema tabel), constraint database, ringkasan kontrak API | Membuat migrasi DB, model/struct Go, atau memanggil API dari React/Kotlin |
| [`05-rules.md`](./05-rules.md) | Aturan non-fungsional (keamanan, privasi, keandalan, UX mobile), asumsi produk (A1–A9) | Selalu jadi rujukan default — dipakai saat ada keputusan teknis yang tidak eksplisit di file lain |
| [`06-tasks.md`](./06-tasks.md) | Rencana kerja 2–3 hari per platform, urutan pemotongan fitur jika waktu mepet, kriteria penerimaan (Definition of Done) | Merencanakan urutan kerja, atau memverifikasi sebuah fitur benar-benar selesai |
| [`07-tech-stack.md`](./07-tech-stack.md) | Library/framework yang sudah terinstal (Go, React, Kotlin) & struktur folder wajib | Setup dependency baru, atau ragu library apa yang harus dipakai |
| [`08-progress-log.md`](./08-progress-log.md) | Riwayat implementasi tiap fitur, keputusan yang diambil saat coding, status Definition of Done | **Setelah** selesai satu fitur (agent wajib update), atau sebelum lanjut fitur yang berkaitan |

## Cara memakai saat memberi instruksi ke agent

Contoh prompt yang efektif, karena langsung menunjuk dokumen yang relevan:

> "Kerjakan F-08 (layar argumen mobile) sesuai `03-flows.md` bagian 4 (khusus 4b) dan endpoint drop/confirm di `04-architecture.md`. Ikuti aturan non-fungsional di `05-rules.md`."

> "Buatkan migrasi tabel `sessions` dan `session_arguments` sesuai `04-architecture.md`, dengan constraint yang disebutkan di sana."

**Aturan umum:** `01-overview.md` dan `05-rules.md` adalah konteks dasar yang sebaiknya selalu ikut dibaca (relatif pendek). Tambahkan 1–2 file lain sesuai tugas spesifik — hindari menyuruh agent membaca seluruh dokumentasi sekaligus untuk tugas kecil, karena itu justru mengembalikan masalah "terlalu panjang" yang ingin dihindari dengan pemecahan ini.

## Kalau dokumen ini perlu diubah

Jika ada perubahan keputusan produk (seperti perubahan interaksi mobile dari drag-and-drop ke tap-to-select), edit file yang relevan secara langsung (dalam kasus itu: `03-flows.md`), bukan `PRD.md` — anggap `PRD.md` sebagai draf awal yang sudah "dibongkar" ke sini dan tidak lagi menjadi sumber kebenaran tunggal.
