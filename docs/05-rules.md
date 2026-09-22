# 05. Rules — Non-Fungsional & Asumsi Produk

File ini adalah rujukan default: kalau ada keputusan teknis yang tidak eksplisit disebut di file lain, cek dulu di sini sebelum agent coding menebak sendiri.

## 1. Persyaratan non-fungsional

- **Keamanan:** password di-hash (bcrypt/argon2), token dengan claim role, otorisasi diperiksa di server untuk setiap endpoint (termasuk kepemilikan data asesor), rate limiting pada login.
- **Privasi:** data seminimal mungkin. Tampilkan pemberitahuan bahwa pilihan siswa dicatat, dan identitas hanya terlihat oleh admin/asesor.
- **Keandalan:** drop dan confirm tidak boleh hilang saat koneksi putus sesaat (retry/antrian sederhana di klien).
- **UX mobile:** interaksi utama adalah **ketuk slot → pilih opsi di bottom sheet** (bukan drag and drop; lihat `03-flows.md` §1.4b), untuk menghindari kerumitan implementasi drag-and-drop custom di Android serta keterbatasan lebar layar untuk teks ground/warrant yang panjang. Area sentuh minimal 48dp.
- **Kualitas:** dokumentasi API (OpenAPI), README dengan diagram arsitektur dan screenshot, seed data yang bisa dijalankan dengan satu perintah.
- **Tanpa Docker/deploy:** semua bisa dijalankan lokal dengan setup minimal.

## 2. Asumsi (bisa diubah — konfirmasi ke pemilik proyek kalau agent ragu)

- **A1.** Kelompok pembanding = semua siswa dengan sesi pada ujian yang sama (bukan per kelas), karena satu ujian bisa diberikan ke beberapa kelas atau siswa individu.
- **A2.** Semua 4 opsi per jenis ditampilkan (sistem lama menampilkan 3 secara default). Jumlah tampil bisa dijadikan pengaturan nanti.
- **A3.** Mode bantuan: setelah Confirm, siswa melihat bagian mana (ground/warrant) yang benar dan yang salah. Alternatif: opsi salah ditandai sejak awal.
- **A4.** Mode dipilih siswa setiap kali mengakses ujian. Untuk sesi berjalan, mode lama tetap dipakai kecuali siswa memilih mode lain dan mengonfirmasinya secara eksplisit (lihat `03-flows.md` §1 dan `04-architecture.md` §3).
- **A5.** Asesor hanya melihat kelas, materi, ujian, dan siswa yang berada di kelasnya.
- **A6.** Melepas opsi kembali ke daftar (web) atau menutup bottom sheet tanpa memilih (mobile) tidak dicatat. Hanya opsi yang benar-benar mengisi slot yang dicatat sebagai drop.
- **A7.** Tidak ada nilai numerik. "Hasil" berarti kemajuan penyelesaian dan jumlah percobaan, per sesi (percobaan).
- **A8.** Materi harus punya argumen lebih dari `arguments_per_session` agar pengacakan bermakna; jika materi punya argumen ≤ `arguments_per_session`, semua argumen dipakai. Pemilihan argumen tiap sesi baru diacak independen (boleh terulang dari sesi sebelumnya, tidak dijamin argumen berbeda).
- **A9.** Perhitungan X/Y/A/B pada analitik sosial (`03-flows.md` §2) menggabungkan attempt log dari **semua sesi** siswa pada ujian itu (bukan hanya sesi aktif), karena sifatnya latihan berulang. Kalau maksudnya per sesi saja, ini perlu diubah — beri tahu pemilik proyek.
