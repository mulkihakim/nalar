# 06. Tasks — Rencana Kerja & Kriteria Penerimaan

## 1. Rencana 2-3 hari

| Hari | Fokus | Hasil |
|---|---|---|
| 1 | Backend | Auth + role, CRUD kelas/materi/ujian, endpoint pengerjaan (drop, confirm, progress), query analitik, seed data, unit test |
| 2 | Web admin/asesor | Login, kelola user/kelas/materi/ujian, tabel hasil, log dengan identitas, analitik sosial versi staf |
| 3 | Mobile siswa + polish | Login, daftar ujian, pilih mode, layar tap-to-select (lihat `02-flows.md` §1.4b), hasil, analitik sosial; README dan screenshot |

**Jika waktu mepet, potong dalam urutan ini:** F-17, F-16, F-15, F-14, tampilan detail per siswa di F-12 (lihat daftar fitur di `01-product.md` §B). Fitur F-01 sampai F-13 adalah inti demo.

## 2. Kriteria penerimaan (contoh — jadikan test/acceptance check)

1. Asesor tidak bisa melihat/mengubah kelas, materi, atau ujian milik asesor lain (diuji otomatis).
2. Membuat argumen dengan jumlah opsi kurang dari 3 ground atau 3 warrant (minimal 2 salah 1 benar), atau jumlah jawaban benar bukan 1+1, ditolak dengan pesan jelas.
3. Siswa yang tidak terdaftar pada ujian, atau ujian nonaktif, mendapat 403/404.
4. Respons klien siswa pada mode standar tidak berisi `is_correct`.
5. Siswa tidak dapat membuka argumen ke-N sebelum argumen ke-(N-1) selesai **dalam sesi yang sama**.
6. Setiap drop menambah tepat satu baris `attempt_logs`.
7. Untuk data 2, 3, 4 percobaan: monitoring menampilkan `2:3`, `3:3`, `4:3` (test unit fungsi rasio).
8. Untuk total 9 percobaan oleh 4 siswa unik: analysis menampilkan `9:4`.
9. Opsi yang belum dipilih siapa pun menampilkan `Y = -` tanpa error.
10. Data kelompok tidak tampil ke siswa jika jumlah siswa pada ujian < `MIN_PEERS`.
11. Siswa yang menyelesaikan sesi lalu menekan ujian yang sama mendapat sesi baru dengan `attempt_no` bertambah, dan subset argumen dipilih ulang secara acak.
12. Siswa dengan sesi `in_progress` yang memilih mode sama langsung dilanjutkan ke argumen belum selesai, tanpa dialog konfirmasi, dan tanpa argumen yang sudah selesai diulang.
13. Siswa dengan sesi `in_progress` yang memilih mode berbeda dan belum mengonfirmasi mendapat status `409` beserta mode sesi berjalan; setelah konfirmasi, `sessions.mode` berubah dan progres tetap.
14. Tidak mungkin ada dua sesi berstatus `in_progress` untuk (siswa, ujian) yang sama pada waktu bersamaan (diuji dengan constraint/index unik parsial atau setara).

> Saat meminta agent menyelesaikan satu fitur, minta juga agent mencocokkan hasilnya terhadap poin-poin relevan di daftar ini sebagai Definition of Done.
