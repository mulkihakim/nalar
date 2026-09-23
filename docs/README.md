# Dokumentasi Proyek Nalar — Panduan untuk Agent Coding

Dokumentasi di folder ini disusun menjadi beberapa file fokus, supaya konteks yang perlu dibaca untuk satu tugas tidak terlalu panjang dan tidak bercampur dengan hal yang tidak relevan. Setiap file memiliki tanggung jawab spesifik sebagai konteks pengerjaan bagi coding agent.

## Stack

- **Backend:** Go
- **Web:** React
- **Mobile:** Kotlin (Android)

## Peta dokumen

| File | Isi | Baca saat... |
|---|---|---|
| [`01-product.md`](./01-product.md) | Ringkasan produk, tujuan/non-tujuan, peran & hak akses, glosarium domain (Bagian A) **+** daftar fitur F-01–F-17 dengan prioritas & platform, definisi 3 mode belajar (Bagian B) | Konteks dasar sebelum mengerjakan fitur apa pun; menentukan scope satu fitur |
| [`02-flows.md`](./02-flows.md) | Alur langkah-demi-langkah pengerjaan siswa (mulai/lanjut sesi, layar argumen versi Web vs Mobile, drop/confirm), aturan analitik sosial, struktur hasil | Implementasi endpoint pengerjaan atau layar F-06 s/d F-12 |
| [`03-architecture.md`](./03-architecture.md) | Model data (ERD + skema tabel), constraint database, ringkasan kontrak API | Membuat migrasi DB, model/struct Go, atau memanggil API dari React/Kotlin |
| [`04-rules-design.md`](./04-rules-design.md) | Aturan non-fungsional (keamanan, privasi, keandalan, UX mobile) & asumsi produk A1–A9 (Bagian A) **+** palet warna, tipografi, layout web & mobile (Bagian B) | Selalu jadi rujukan default untuk keputusan teknis yang tidak eksplisit di file lain; wajib dibaca untuk task apa pun yang menyentuh UI |
| [`05-tech-stack.md`](./05-tech-stack.md) | Library/framework yang sudah terinstal (Go, React, Kotlin) & struktur folder wajib per platform | Setup dependency baru, atau ragu library apa yang harus dipakai |
| [`06-tasks.md`](./06-tasks.md) | Rencana kerja 2–3 hari per platform, urutan pemotongan fitur jika waktu mepet, kriteria penerimaan (Definition of Done) | Merencanakan urutan kerja, atau memverifikasi sebuah fitur benar-benar selesai |
| [`07-progress-log.md`](./07-progress-log.md) | Riwayat implementasi tiap fitur, keputusan yang diambil saat coding, status Definition of Done | **Setelah** selesai satu fitur (agent wajib update), atau sebelum lanjut fitur yang berkaitan |

## Cara memakai saat memberi instruksi ke agent

Contoh prompt yang efektif, karena langsung menunjuk dokumen yang relevan:

> "Kerjakan F-08 (layar argumen mobile) sesuai `02-flows.md` bagian 1 (khusus 1.4b) dan endpoint drop/confirm di `03-architecture.md`. Ikuti aturan non-fungsional di `04-rules-design.md` Bagian A."

> "Buatkan migrasi tabel `sessions` dan `session_arguments` sesuai `03-architecture.md`, dengan constraint yang disebutkan di sana."

**Aturan umum:** `01-product.md` dan `04-rules-design.md` adalah konteks dasar yang sebaiknya selalu ikut dibaca (relatif pendek). Tambahkan 1–2 file lain sesuai tugas spesifik — hindari menyuruh agent membaca seluruh dokumentasi sekaligus untuk tugas kecil.

## Kalau dokumen ini perlu diubah

Jika ada perubahan keputusan produk (seperti perubahan interaksi mobile dari drag-and-drop ke tap-to-select), edit file yang relevan secara langsung (dalam kasus itu: `02-flows.md`). Folder `docs/` ini adalah satu-satunya sumber kebenaran (single source of truth) untuk spesifikasi proyek.

Kalau menambah entri di `07-progress-log.md`, ikuti format ringkas di bagian atas file itu — jangan menyalin seluruh detail implementasi kalau cukup ditulis sebagai poin singkat.
