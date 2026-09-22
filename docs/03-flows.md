# 03. Alur dan Aturan Pengerjaan

Prasyarat: `01-overview.md` (istilah domain), `02-requirements.md` (fitur terkait: F-06 s/d F-12, F-14).
Terkait: `04-architecture.md` untuk kontrak endpoint yang dipanggil di setiap langkah.

## 1. Alur mulai/lanjut sesi dan pengerjaan argumen

1. Siswa membuka daftar ujian. Hanya ujian **aktif** yang siswa akses lewat kelasnya atau penugasan individu yang muncul.
2. Siswa menekan sebuah ujian. Server memeriksa apakah siswa punya **sesi berstatus berjalan** pada ujian itu:
   - **Tidak ada** (belum pernah, atau semua sesi sebelumnya sudah selesai): siswa memilih mode (`02-requirements.md` §2), lalu server membuat sesi baru, mengambil maksimal 3 argumen secara acak dari materi (`arguments_per_session`), dan menyimpannya sebagai session argument dengan urutan tetap.
   - **Ada sesi berjalan**: siswa tetap diminta memilih mode.
     - Jika mode yang dipilih **sama** dengan mode sesi berjalan: langsung lanjutkan ke argumen pertama yang belum selesai. Riwayat argumen yang sudah dijawab benar pada sesi itu tidak hilang, dan subset argumennya tidak diacak ulang.
     - Jika mode yang dipilih **berbeda**: tampilkan konfirmasi "mode belajar akan diubah dari X ke Y, lanjutkan sesi ini?". Jika dikonfirmasi, `sessions.mode` diperbarui dan siswa melanjutkan ke argumen yang belum selesai (subset argumen tetap sama, progres tetap sama). Jika dibatalkan, siswa tetap di halaman daftar ujian.
3. Siswa membaca materi, lalu masuk ke argumen pertama yang belum selesai pada sesi itu.
4. Layar argumen, dengan indikator "Argumen X dari N" mengacu ke posisi dalam session argument (N maksimal 3). Tata letak dan interaksi **berbeda per platform**, tapi keduanya memanggil endpoint drop/confirm yang sama (lihat `04-architecture.md` §API):

   **4a. Web (mengikuti desain lama, drag and drop)**
   - Claim di bagian atas.
   - Slot **Ground** di kiri dan slot **Warrant** di kanan. Satu panah dari ground ke claim, dan satu panah dari warrant yang menunjuk ke panah tersebut (warrant menjembatani ground dan claim).
   - Daftar opsi ground dan warrant di bawah slot masing-masing, masing-masing 4 opsi.
   - Tombol **Confirm**, aktif jika kedua slot terisi.

   **4b. Mobile (tap-to-select, bukan drag and drop)**
   - Claim di bagian atas.
   - Slot **Ground** dan slot **Warrant** ditampilkan sebagai dua kartu kosong tersusun vertikal (ground lalu warrant), dengan placeholder "Pilih Ground" / "Pilih Warrant" saat kosong. Panah/diagram penghubung ground-warrant-claim pada desain lama **tidak dipakai** di mobile karena ruang layar terbatas.
   - Daftar 4 opsi ground dan 4 opsi warrant **tidak ditampilkan langsung di layar**; opsi baru muncul saat slot terkait ditekan.
   - Tombol **Confirm**, aktif jika kedua slot terisi.
5. Siswa mengisi slot, mekanismenya berbeda per platform:
   - **Web:** siswa men-*drop* opsi ke slot. Opsi ground hanya bisa ke slot ground, dan opsi warrant hanya ke slot warrant (divalidasi di klien dan server).
   - **Mobile:** siswa menekan slot yang masih kosong. Menekan slot Ground membuka lembar pilihan (bottom sheet) berisi 4 opsi ground argumen tersebut; menekan slot Warrant membuka lembar berisi 4 opsi warrant. Siswa memilih satu opsi pada lembar tersebut, lembar tertutup, dan opsi itu mengisi slot. Slot yang sudah terisi tetap bisa ditekan ulang untuk membuka lembar pilihan dan mengganti isinya.
   - Di kedua platform, setiap pengisian/penggantian slot dicatat sebagai satu **attempt log** (payload dan endpoint sama; istilah "drop" pada data/API tetap dipakai sebagai konsep, terlepas dari interaksi UI-nya berupa seret atau ketuk-pilih). Mengganti isi slot dengan opsi lain dicatat sebagai drop baru.
6. Tekan **Confirm**. Server menilai:
   - Ground dan warrant benar: argumen ditandai selesai (argument progress), argumen berikutnya dalam session argument terbuka.
   - Salah satu atau keduanya salah: siswa diberi tahu belum tepat (detail sesuai mode) dan boleh mencoba lagi tanpa batas.
7. Setelah argumen terakhir dalam sesi itu selesai, sesi berstatus **selesai** dan siswa melihat ringkasan hasil (§3 di bawah). Siswa bisa memulai sesi baru (percobaan baru) kapan pun, dengan subset argumen yang baru diacak.
8. Kebenaran jawaban **selalu dinilai di server**. Field `is_correct` tidak dikirim ke klien kecuali sesuai aturan mode.
9. Timer pada desain lama dihapus.
10. Jika koneksi putus atau aplikasi ditutup di tengah sesi, status sesi tetap **berjalan** dan mengikuti aturan poin 2 saat diakses kembali.

## 2. Log dan analitik sosial

### 2.1 Definisi dasar
- **Percobaan (attempt):** satu kali opsi di-drop ke slot. Satu baris di `attempt_logs`.
- **Kelompok pembanding:** semua siswa yang memiliki minimal satu sesi pada ujian yang sama.
- Perhitungan **X** (percobaan siswa sendiri) dan agregat kelompok menggabungkan attempt log dari **seluruh sesi** siswa pada ujian itu, bukan hanya sesi yang sedang berjalan (lihat asumsi A9 di `05-rules.md`).
- Perhitungan dilakukan di backend per ujian, per argumen, per opsi.

### 2.2 Tampilan Monitoring (setiap argumen selesai, mode analitik sosial)
Untuk setiap opsi (4 ground dan 4 warrant), tampilkan **X : Y**.

- `X` = jumlah percobaan siswa itu sendiri pada opsi tersebut.
- `Y` = rata-rata percobaan per siswa yang memilih opsi tersebut.

```text
Y = round(total_percobaan_semua_siswa_pada_opsi / jumlah_siswa_unik_yang_memilih_opsi)
```

Contoh: opsi dipilih 3 siswa dengan percobaan 2, 3, dan 4. Total 9, siswa unik 3, sehingga Y = 3. Siswa yang memilih 2 kali melihat **2:3**.

`Y` bukan rata-rata seluruh siswa di kelas, karena penyebutnya hanya siswa unik yang pernah memilih opsi tersebut. Pembulatan mengikuti sistem lama (setengah dibulatkan menjauhi nol, setara `math.Round` di Go).

### 2.3 Tampilan Analysis (setelah semua argumen selesai, mode analitik sosial)
Untuk setiap opsi pada setiap argumen, tampilkan **A : B**.

- `A` = total percobaan seluruh kelompok pada opsi tersebut.
- `B` = jumlah siswa unik yang pernah memilih opsi tersebut.

Contoh: **9:4** berarti opsi dipilih total 9 kali oleh 4 siswa berbeda.

| Tampilan | Angka pertama | Angka kedua |
|---|---|---|
| Monitoring | Percobaan siswa sendiri (X) | Rata-rata percobaan siswa pemilih (Y) |
| Analysis | Total percobaan kelompok (A) | Jumlah siswa unik pemilih (B) |

### 2.4 Aturan tambahan
- Jika tidak ada siswa yang memilih opsi, `Y` ditampilkan `-` (hindari pembagian nol).
- **Anonimitas:** data kelompok hanya tampil ke siswa jika jumlah siswa dengan sesi pada ujian itu >= `MIN_PEERS` (default 3, dapat dikonfigurasi). Identitas siswa lain tidak pernah dikirim ke klien siswa.
- Tampilan admin/asesor memakai data yang sama, tetapi menampilkan **nama siswa** dan rincian percobaan per siswa.
- Sebaiknya UI juga menampilkan bar distribusi selain angka rasio.

## 3. Hasil

Tidak ada nilai numerik. Karena satu siswa boleh punya banyak sesi (percobaan) pada satu ujian, hasil sesi berisi:
- Nomor percobaan, status (berjalan/selesai), mode yang dipakai, dan jumlah argumen selesai dari total pada sesi itu.
- Per argumen (dalam sesi itu): total percobaan dan penanda **tanpa salah** (selesai dengan tepat 2 percobaan, yaitu 1 ground dan 1 warrant).
- Total percobaan seluruh sesi itu.

Siswa melihat **riwayat semua sesinya** pada suatu ujian (daftar percobaan dan hasil masing-masing). Admin/asesor melihat tabel seluruh siswa per ujian, dengan seluruh sesi tiap siswa, dan detailnya dapat dibuka per sesi (opsi mana dipilih, urutan, jumlah, waktu pencatatan).
