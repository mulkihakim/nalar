# 01. Overview — Nalar

## 1. Ringkasan

Nalar adalah platform **latihan argumentasi berbasis model Toulmin**. Asesor menyusun materi bacaan beserta argumen (claim). Siswa menyusun argumen dengan memilih *ground* dan *warrant* yang tepat untuk sebuah *claim* (drag and drop di web, tap-to-select di mobile — detail di `03-flows.md`). Sistem mencatat setiap percobaan memilih opsi, lalu menampilkan hasil dan perbandingan pilihan antarsiswa secara anonim (analitik sosial).

Ini **latihan belajar, bukan ujian formal**. Tidak ada nilai numerik atau timer. Siswa boleh mengerjakan ujian yang sama berkali-kali sebagai latihan, dan setiap percobaan mengambil subset argumen secara acak (maksimal 3) dari materi. Siswa lanjut ke argumen berikutnya setelah menemukan jawaban yang benar. Kata "ujian" tetap dipakai sebagai istilah untuk satu paket latihan yang diberikan ke siswa.

## 2. Tujuan dan non-tujuan

**Tujuan**
- Siswa dapat menyelesaikan satu sesi latihan penuh dari HP (Android).
- Admin dan asesor dapat melihat siapa memilih opsi apa, dan berapa kali.
- Repo portofolio yang menunjukkan web + mobile + backend, dengan test untuk logika inti.

**Non-tujuan (versi ini)**
Tipe ujian, konsep eksperimen (kontrol/perlakuan, pre/post), adaptive learning, kategori peserta, skor perilaku (steps/waktu), timer, nilai numerik/ranking antar percobaan, gamifikasi, chatbot/AI, deteksi ekspresi wajah, deploy, Docker, mode offline, multi-bahasa.

> Kalau agent coding diminta menambahkan salah satu hal di atas, itu di luar scope PRD ini — konfirmasi dulu ke pemilik proyek sebelum mengerjakan.

## 3. Peran dan hak akses

| Kemampuan | Admin | Asesor | Siswa |
|---|---|---|---|
| Kelola akun asesor | Ya | - | - |
| Buat akun siswa | Ya | Ya | - |
| Ubah/nonaktifkan akun siswa | Semua | Yang ia buat | - |
| Kelola kelas dan anggota kelas | Semua | Milik sendiri | - |
| Kelola materi dan argumen | Semua | Milik sendiri | - |
| Kelola ujian (materi, status, akses) | Semua | Milik sendiri | - |
| Lihat hasil dan log dengan identitas | Semua | Ujian miliknya | - |
| Lihat daftar ujian yang boleh diakses | - | - | Ya |
| Mengerjakan ujian | - | - | Ya |
| Lihat hasil sendiri dan analitik sosial anonim | - | - | Ya |

Satu siswa boleh menjadi anggota banyak kelas.

## 4. Konsep domain (glosarium)

Istilah ini dipakai konsisten di seluruh dokumentasi dan di kode (nama tabel, endpoint, variabel) — pakai istilah yang sama, jangan diterjemahkan ulang secara bebas.

- **Material:** bacaan/topik. Memiliki beberapa Argument yang berurutan.
- **Argument:** satu *claim* + **4 opsi ground** + **4 opsi warrant**. Tepat 1 ground benar dan 1 warrant benar.
- **Class:** kelompok siswa yang dimiliki satu asesor/admin.
- **Exam:** menghubungkan satu Material dengan status aktif/nonaktif, daftar akses (kelas dan/atau siswa individu), dan `arguments_per_session` (default 3, jumlah maksimum argumen yang diambil secara acak dari materi untuk satu sesi).
- **Session:** satu **percobaan** siswa pada satu ujian. Siswa boleh membuat banyak sesi (belajar berkali-kali). Setiap sesi menyimpan subset argumen (maksimal 3, diacak) yang tetap sepanjang sesi itu, dan mode belajar yang dipilih. Hanya boleh ada **satu sesi berstatus berjalan** per (siswa, ujian) pada satu waktu; sesi baru hanya dibuat setelah sesi sebelumnya berstatus selesai.
- **Session argument:** daftar argumen yang terpilih untuk satu sesi, beserta urutannya. Ditentukan sekali saat sesi dibuat dan tidak berubah meski sesi dilanjutkan nanti.
- **Attempt log:** satu catatan setiap siswa men-*drop* (mengisi slot dengan) sebuah opsi ke slot jawaban. Istilah "drop" tetap dipakai secara konseptual di data/API meskipun di mobile interaksinya berupa tap-to-select, bukan drag-and-drop — lihat `03-flows.md`.
- **Argument progress:** penanda argumen (dalam ruang lingkup sesi tertentu) yang sudah dijawab benar.

Lihat `04-architecture.md` untuk bagaimana istilah-istilah ini dipetakan ke tabel database.
