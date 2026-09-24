package main

import (
	"flag"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mulkihakim/nalar/backend/internal/class"
	"github.com/mulkihakim/nalar/backend/internal/db"
	"github.com/mulkihakim/nalar/backend/internal/exam"
	"github.com/mulkihakim/nalar/backend/internal/material"
	"github.com/mulkihakim/nalar/backend/internal/session"
	"github.com/mulkihakim/nalar/backend/internal/user"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fresh := flag.Bool("fresh", false, "Drop all tables and re-migrate from models before seeding (like migrate:fresh --seed)")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	dsn := db.GetDSN()
	if dsn == "" {
		log.Fatal("database configuration (DATABASE_URL or DB_*) must be set")
	}

	gormDB, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// 1. Migrasi database
	if *fresh {
		log.Println("🔄 Resetting database schema (fresh)...")
		if err := db.ResetAndMigrate(gormDB); err != nil {
			log.Fatalf("failed to reset and migrate database: %v", err)
		}
		log.Println("✅ Database schema recreated successfully via AutoMigrate!")
	} else {
		if err := db.AutoMigrate(gormDB); err != nil {
			log.Printf("warning on auto-migrate: %v", err)
		} else {
			log.Println("✅ Database auto-migrated successfully from models")
		}
	}

	userRepo := user.NewRepository(gormDB)
	classRepo := class.NewRepository(gormDB)
	materialRepo := material.NewRepository(gormDB)
	examRepo := exam.NewRepository(gormDB)
	sessionRepo := session.NewRepository(gormDB)

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	// -------------------------------------------------------------
	// 2. SEED USERS (1 Admin, 2 Asesor, 5 Siswa)
	// -------------------------------------------------------------
	seedUsers := []user.User{
		{
			Name:         "Super Admin",
			Username:     "admin",
			PasswordHash: string(hash),
			Role:         "admin",
			IsActive:     true,
		},
		{
			Name:         "Budi Asesor, M.Pd.",
			Username:     "asesor1",
			PasswordHash: string(hash),
			Role:         "asesor",
			IsActive:     true,
		},
		{
			Name:         "Dr. Dewi Lestari",
			Username:     "asesor2",
			PasswordHash: string(hash),
			Role:         "asesor",
			IsActive:     true,
		},
		{
			Name:         "Siti Rahmawati",
			Username:     "siswa1",
			PasswordHash: string(hash),
			Role:         "siswa",
			IsActive:     true,
		},
		{
			Name:         "Ahmad Fauzi",
			Username:     "siswa2",
			PasswordHash: string(hash),
			Role:         "siswa",
			IsActive:     true,
		},
		{
			Name:         "Rizky Pratama",
			Username:     "siswa3",
			PasswordHash: string(hash),
			Role:         "siswa",
			IsActive:     true,
		},
		{
			Name:         "Nadia Putri",
			Username:     "siswa4",
			PasswordHash: string(hash),
			Role:         "siswa",
			IsActive:     true,
		},
		{
			Name:         "Fajar Nugraha",
			Username:     "siswa5",
			PasswordHash: string(hash),
			Role:         "siswa",
			IsActive:     true,
		},
	}

	createdUsers := make(map[string]*user.User)
	for _, u := range seedUsers {
		existing, err := userRepo.FindByUsername(u.Username)
		if err != nil {
			log.Printf("error checking user %s: %v", u.Username, err)
			continue
		}
		if existing != nil {
			createdUsers[u.Username] = existing
			continue
		}

		userToCreate := u
		if err := userRepo.Create(&userToCreate); err != nil {
			log.Printf("failed to create user %s: %v", u.Username, err)
		} else {
			createdUsers[u.Username] = &userToCreate
			log.Printf("👤 Created user: %s (role: %s, password: password123)", u.Username, u.Role)
		}
	}

	asesor1 := createdUsers["asesor1"]
	asesor2 := createdUsers["asesor2"]
	siswa1 := createdUsers["siswa1"]
	siswa2 := createdUsers["siswa2"]
	siswa3 := createdUsers["siswa3"]
	siswa4 := createdUsers["siswa4"]
	siswa5 := createdUsers["siswa5"]

	if asesor1 == nil || asesor2 == nil {
		log.Println("asesor users not ready, skipping downstream seeding")
		return
	}

	// -------------------------------------------------------------
	// 3. SEED CLASSES (3 Varied Classes with different members)
	// -------------------------------------------------------------
	classesConfig := []struct {
		Name    string
		Owner   *user.User
		Members []*user.User
	}{
		{
			Name:    "Kelas Argumentasi X-A",
			Owner:   asesor1,
			Members: []*user.User{siswa1, siswa2, siswa3},
		},
		{
			Name:    "Kelas Literasi & Logika XI-IPA",
			Owner:   asesor1,
			Members: []*user.User{siswa1, siswa3, siswa4, siswa5},
		},
		{
			Name:    "Kelas Filsafat Kritis XII-IPS",
			Owner:   asesor2,
			Members: []*user.User{siswa2, siswa4, siswa5},
		},
	}

	createdClasses := make(map[string]*class.Class)
	for _, cfg := range classesConfig {
		var c *class.Class
		var found class.Class
		if err := gormDB.Where("name = ? AND owner_id = ?", cfg.Name, cfg.Owner.ID).First(&found).Error; err == nil {
			c = &found
		} else {
			c = &class.Class{
				Name:    cfg.Name,
				OwnerID: cfg.Owner.ID,
			}
			if err := classRepo.Create(c); err != nil {
				log.Printf("failed to create class %s: %v", cfg.Name, err)
				continue
			}
			log.Printf("🏫 Created class: %s (Owner: %s)", c.Name, cfg.Owner.Name)
		}

		// Sync members
		for _, m := range cfg.Members {
			if m != nil {
				_ = classRepo.AddMember(c.ID, m.ID)
			}
		}
		createdClasses[cfg.Name] = c
	}

	// -------------------------------------------------------------
	// 4. SEED MATERIALS & TOULMIN ARGUMENTS
	// -------------------------------------------------------------
	type SeedOption struct {
		Type      string
		Text      string
		IsCorrect bool
	}
	type SeedArgument struct {
		ClaimText string
		OrderNo   int
		Options   []SeedOption
	}
	type SeedMaterial struct {
		Title     string
		Content   string
		Owner     *user.User
		Arguments []SeedArgument
	}

	materialsConfig := []SeedMaterial{
		// Materi 1: Literasi Digital & Berpikir Kritis
		{
			Title: "Pentingnya Berpikir Kritis dalam Era Informasi",
			Content: `Di era informasi digital, arus informasi mengalir sangat deras dari berbagai kanal jejaring sosial dan portal berita daring. 
Masyarakat dihadapkan pada tantangan besar berupa misinformasi, disinformasi, dan polarisasi opini publik.
Kemampuan mengevaluasi bukti empiris (ground) dan menghubungkannya dengan alasan prinsip (warrant) 
merupakan fondasi utama literasi kritis model Toulmin dalam membedakan fakta obyektif dari klaim manipulatif.`,
			Owner: asesor1,
			Arguments: []SeedArgument{
				{
					ClaimText: "Verifikasi fakta dan analisis sumber independen secara signifikan meminimalkan penyebaran hoaks di media sosial.",
					OrderNo:   1,
					Options: []SeedOption{
						{Type: "ground", Text: "Data penelitian Masyarakat Telematika (Mastel) menunjukkan 78% pengguna internet yang melakukan verifikasi silang pada kanal terpercaya berhasil menghindari kepanikan akibat berita palsu.", IsCorrect: true},
						{Type: "ground", Text: "Banyak pengguna media sosial lebih sering membaca judul dan cuplikan berita tanpa membuka tautan artikel lengkap.", IsCorrect: false},
						{Type: "ground", Text: "Algoritma platform digital dirancang untuk memprioritaskan konten dengan tingkat interaksi dan emosi pengguna yang tinggi.", IsCorrect: false},
						{Type: "ground", Text: "Jumlah pengguna telepon pintar di Indonesia meningkat secara pesat dalam lima tahun terakhir di wilayah perkotaan.", IsCorrect: false},
						{Type: "warrant", Text: "Menguji keabsahan informasi terhadap sumber kredibel memvalidasi fakta empiris sebelum klaim diterima atau disebarluaskan.", IsCorrect: true},
						{Type: "warrant", Text: "Informasi yang viral di internet dapat diperbarui sewaktu-waktu oleh komunitas daring tanpa kurasi redaksi.", IsCorrect: false},
						{Type: "warrant", Text: "Membaca artikel panjang dari media cetak konvensional membutuhkan daya konsentrasi yang lebih tinggi.", IsCorrect: false},
						{Type: "warrant", Text: "Koneksi internet berkecepatan tinggi mempercepat proses unduhan dokumen dan tayangan multimedia.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Memahami bias konfirmasi pribadi memperkuat objektivitas seseorang dalam mengevaluasi argumen yang berseberangan.",
					OrderNo:   2,
					Options: []SeedOption{
						{Type: "ground", Text: "Eksperimen psikologi kognitif membuktikan individu yang dilatih mengenali kecenderungan preferensi pribadinya 65% lebih akurat menilai keabsahan data lawan bicara.", IsCorrect: true},
						{Type: "ground", Text: "Kebanyakan orang secara alami merasa tidak nyaman ketika pendapat dasarnya disanggah di ruang publik.", IsCorrect: false},
						{Type: "ground", Text: "Kelompok diskusi tertutup cenderung memperkuat keyakinan yang sudah ada di antara para anggotanya.", IsCorrect: false},
						{Type: "ground", Text: "Buku-buku pengantar logika formal memuat berbagai daftar kesalahan berpikir (logical fallacies) yang umum dijumpai.", IsCorrect: false},
						{Type: "warrant", Text: "Kesadaran atas kecenderungan mencari informasi pendukung semata memungkinkan evaluasi bukti secara netral dan proporsional.", IsCorrect: true},
						{Type: "warrant", Text: "Perdebatan yang berlangsung lama di forum daring sering kali bergeser menjadi serangan personal antar peserta.", IsCorrect: false},
						{Type: "warrant", Text: "Setiap pandangan subjektif memiliki latar belakang kultural dan pengalaman masa lalu yang berbeda-beda.", IsCorrect: false},
						{Type: "warrant", Text: "Menghindari perdebatan yang memanas merupakan salah satu cara menjaga kerukunan dalam pergaulan sosial.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Integrasi kurikulum literasi media digital di sekolah meningkatkan resiliensi siswa terhadap narasi polarisasi ekstrem.",
					OrderNo:   3,
					Options: []SeedOption{
						{Type: "ground", Text: "Studi komparatif di 40 sekolah menunjukkan siswa yang menempuh modul analisis wacana kritis memiliki skor diskriminasi propaganda 42% lebih tinggi.", IsCorrect: true},
						{Type: "ground", Text: "Sebagian besar siswa sekolah menengah menghabiskan lebih dari 4 jam sehari menggunakan gawai pintar.", IsCorrect: false},
						{Type: "ground", Text: "Guru-guru mata pelajaran bahasa dan sosial kerap menyisipkan isu terkini dalam diskusi pembelajaran di kelas.", IsCorrect: false},
						{Type: "ground", Text: "Banyak fasilitas laboratorium komputer di sekolah kini telah tersambung dengan jaringan serat optik nasional.", IsCorrect: false},
						{Type: "warrant", Text: "Pendidikan literasi melatih dekonstruksi retorika manipulatif sehingga siswa mampu mengenali motif di balik pembingkaian narasi ekstrem.", IsCorrect: true},
						{Type: "warrant", Text: "Penggunaan gawai interaktif dalam proses belajar dapat meningkatkan antusiasme siswa dalam menyelesaikan tugas.", IsCorrect: false},
						{Type: "warrant", Text: "Diskusi kelompok di ruang kelas menciptakan interaksi sosial yang mempererat rasa persahabatan antar pelajar.", IsCorrect: false},
						{Type: "warrant", Text: "Perangkat lunak pembatas konten (content filter) dapat dipasang di server sekolah untuk memblokir situs terlarang.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Pemeriksaan rekam jejak dan kredensial ahli merupakan prasyarat esensial sebelum mempercayai opini otoritatif di ruang publik.",
					OrderNo:   4,
					Options: []SeedOption{
						{Type: "ground", Text: "Analisis kepatuhan publik saat krisis kesehatan mencatat kesalahan tafsir berkurang 70% ketika rujukan ilmiah bersumber dari epidemiolog bereputasi internasional.", IsCorrect: true},
						{Type: "ground", Text: "Influencer kesehatan di media sosial memiliki jumlah pengikut dan jangkauan tayangan yang jauh lebih besar daripada akademisi.", IsCorrect: false},
						{Type: "ground", Text: "Masyarakat awam kesulitan membaca laporan jurnal ilmiah yang sarat dengan terminologi teknis kedokteran.", IsCorrect: false},
						{Type: "ground", Text: "Konferensi pers resmi pemerintah disiarkan secara serentak di stasiun televisi dan saluran streaming internet.", IsCorrect: false},
						{Type: "warrant", Text: "Kompetensi metodologis dan pengakuan sejawat (peer review) menjamin keabsahan epistemologis dari klaim keilmuan yang dipublikasikan.", IsCorrect: true},
						{Type: "warrant", Text: "Penyampaian pesan dengan bahasa santai dan visual menarik lebih mudah dicerna oleh masyarakat umum.", IsCorrect: false},
						{Type: "warrant", Text: "Pemerintah memiliki regulasi ketat mengenai penyiaran informasi darurat kepada khalayak ramai.", IsCorrect: false},
						{Type: "warrant", Text: "Banyak portal berita daring mengandalkan pendapatan dari tayangan iklan berbayar per klik (adsense).", IsCorrect: false},
					},
				},
				{
					ClaimText: "Kebiasaan memisahkan antara fakta obyektif dan interpretasi emotif mencegah kepanikan massal saat situasi darurat bencana.",
					OrderNo:   5,
					Options: []SeedOption{
						{Type: "ground", Text: "Laporan penanggulangan bencana mendokumentasikan kawasan dengan pusat informasi berbasis data lapangan memiliki tingkat kepanikan warga 55% lebih rendah.", IsCorrect: true},
						{Type: "ground", Text: "Foto dan rekaman video amatir yang dramatis menyebar sangat cepat di grup percakapan instan warga saat bencana terjadi.", IsCorrect: false},
						{Type: "ground", Text: "Bantuan logistik darurat memerlukan koordinasi lintas lembaga agar dapat disalurkan secara merata ke posko penampungan.", IsCorrect: false},
						{Type: "ground", Text: "Warga yang terdampak bencana membutuhkan dukungan psikososial untuk memulihkan trauma pasca kejadian.", IsCorrect: false},
						{Type: "warrant", Text: "Fakta terukur menyediakan dasar tindakan rasional, sedangkan penilaian emotif yang tidak terverifikasi memicu respons kepanikan yang destruktif.", IsCorrect: true},
						{Type: "warrant", Text: "Kecepatan distribusi sembako sangat menentukan tingkat kepuasan warga terhadap kinerja aparat tanggap darurat.", IsCorrect: false},
						{Type: "warrant", Text: "Dukungan moral dari relawan memberikan ketenangan batin bagi keluarga yang kehilangan tempat tinggal.", IsCorrect: false},
						{Type: "warrant", Text: "Alat komunikasi satelit tetap dapat berfungsi optimal meskipun infrastruktur menara BTS seluler mengalami kerusakan.", IsCorrect: false},
					},
				},
			},
		},

		// Materi 2: Transisi Energi Terbarukan
		{
			Title: "Transisi Energi Terbarukan Menuju Net-Zero Emission",
			Content: `Pemanasan global menuntut transformasi menyeluruh dari pembangkit berbahan bakar fosil menuju sumber energi baru terbarukan.
Meskipun menjanjikan pengurangan emisi gas rumah kaca secara drastis, transisi ini menghadapi tantangan teknis intermitensi suplai cuaca,
keandalan jaringan listrik tegangan tinggi, serta kalkulasi keekonomian penerapan instrumen disinsentif fiskal seperti pajak karbon.`,
			Owner: asesor1,
			Arguments: []SeedArgument{
				{
					ClaimText: "Integrasi sistem penyimpanan baterai (BESS) skala besar menjaga kestabilan frekuensi jaringan listrik bertenaga surya dan angin.",
					OrderNo:   1,
					Options: []SeedOption{
						{Type: "ground", Text: "Uji operasional jaringan listrik di Australia Selatan mencatat sistem baterai merespons deviasi frekuensi dalam waktu kurang dari 200 milidetik, mencegah pemadaman massal.", IsCorrect: true},
						{Type: "ground", Text: "Pembangkit listrik tenaga uap batu bara masih menyumbang porsi energi dasar terbesar di banyak negara berkembang.", IsCorrect: false},
						{Type: "ground", Text: "Biaya produksi panel fotovoltaik surya terus mengalami penurunan tajam selama satu dekade terakhir.", IsCorrect: false},
						{Type: "ground", Text: "Ladang turbin angin lepas pantai membutuhkan area instalasi yang luas serta survei hidrografi mendalam.", IsCorrect: false},
						{Type: "warrant", Text: "Kemampuan respons cepat injeksi daya baterai mengkompensasi fluktuasi intermitensi cuaca secara real-time sehingga keandalan transmisi tetap terjaga.", IsCorrect: true},
						{Type: "warrant", Text: "Ketersediaan cadangan batu bara nasional yang melimpah memberikan rasa aman bagi pasokan listrik industri berat.", IsCorrect: false},
						{Type: "warrant", Text: "Penurunan harga modul pembangkit mempercepat waktu pengembalian modal investasi bagi pengembang swasta.", IsCorrect: false},
						{Type: "warrant", Text: "Wilayah perairan laut terbuka memiliki potensi hembusan angin yang jauh lebih konstan dibanding daratan.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Penerapan pajak karbon yang progresif mempercepat adopsi teknologi bersih di sektor industri manufaktur.",
					OrderNo:   2,
					Options: []SeedOption{
						{Type: "ground", Text: "Laporan OECD menunjukkan negara-negara yang menetapkan harga karbon di atas $50 per ton berhasil menurunkan emisi industri sebesar 23% dalam kurun waktu empat tahun.", IsCorrect: true},
						{Type: "ground", Text: "Industri semen dan baja menyerap tenaga kerja dalam jumlah sangat besar di berbagai kawasan industri terpadu.", IsCorrect: false},
						{Type: "ground", Text: "Beberapa asosiasi pengusaha mengkhawatirkan kenaikan biaya produksi dapat menurunkan daya saing ekspor produk lokal.", IsCorrect: false},
						{Type: "ground", Text: "Program tanggung jawab sosial perusahaan (CSR) sering kali mencakup kegiatan penanaman pohon penghijauan.", IsCorrect: false},
						{Type: "warrant", Text: "Beban biaya emisi secara langsung mengubah struktur kalkulasi finansial, menjadikan investasi teknologi rendah emisi lebih ekonomis dibandingkan membayar denda polusi.", IsCorrect: true},
						{Type: "warrant", Text: "Penyerapan tenaga kerja lokal merupakan prioritas utama pemerintah daerah dalam menjaga stabilitas ekonomi wilayah.", IsCorrect: false},
						{Type: "warrant", Text: "Daya saing harga barang di pasar global sangat dipengaruhi oleh fluktuasi kurs mata uang dan tarif bea masuk antarnegara.", IsCorrect: false},
						{Type: "warrant", Text: "Kegiatan penanaman pohon mangrove memberikan dampak positif bagi keanekaragaman hayati di kawasan pesisir.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Elektrifikasi armada transportasi umum perkotaan mengurangi konsentrasi partikulat berbahaya PM2.5 secara terukur.",
					OrderNo:   3,
					Options: []SeedOption{
						{Type: "ground", Text: "Pengukuran kualitas udara di Shenzhen membuktikan peralihan 100% bus dan taksi listrik berkontribusi pada penurunan konsentrasi PM2.5 tahunan sebesar 31%.", IsCorrect: true},
						{Type: "ground", Text: "Stasiun pengisian daya kendaraan listrik umum (SPKLU) memerlukan investasi infrastruktur kelistrikan tegangan menengah.", IsCorrect: false},
						{Type: "ground", Text: "Waktu tunggu pengisian baterai kendaraan komersial berkisar antara 30 hingga 60 menit dengan fasilitas fast charging.", IsCorrect: false},
						{Type: "ground", Text: "Kemacetan lalu lintas pada jam sibuk di kota metropolitan menambah waktu tempuh komuter secara signifikan.", IsCorrect: false},
						{Type: "warrant", Text: "Penghapusan emisi gas buang knalpot secara langsung mengeliminasi pembentukan polutan aerosol sekunder beracun di koridor jalan raya yang padat.", IsCorrect: true},
						{Type: "warrant", Text: "Penyediaan trafo daya tambahan menjamin pasokan listrik stasiun pengisian kendaraan tetap mencukupi.", IsCorrect: false},
						{Type: "warrant", Text: "Fasilitas ruang tunggu yang nyaman di SPKLU memungkinkan pengemudi beristirahat selama pengisian daya.", IsCorrect: false},
						{Type: "warrant", Text: "Pelebaran jalan arteri dan jalur layang sering dipilih pemda untuk mengurai antrean kendaraan.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Diversifikasi sumber energi baru terbarukan memitigasi kerentanan ketahanan energi nasional dari krisis geopolitik global.",
					OrderNo:   4,
					Options: []SeedOption{
						{Type: "ground", Text: "Negara-negara dengan bauran energi terbarukan lokal di atas 40% mengalami inflasi sektor energi 60% lebih rendah saat krisis pasokan minyak mentah dunia.", IsCorrect: true},
						{Type: "ground", Text: "Harga minyak mentah dunia ditentukan oleh dinamika kuota produksi negara anggota kartel produsen energi internasional.", IsCorrect: false},
						{Type: "ground", Text: "Banyak negara masih mengalokasikan anggaran subsidi bahan bakar minyak (BBM) untuk menjaga daya beli masyarakat miskin.", IsCorrect: false},
						{Type: "ground", Text: "Pipa transmisi gas bumi lintas negara memerlukan perjanjian bilateral keamanan dan jaminan operasional yang kompleks.", IsCorrect: false},
						{Type: "warrant", Text: "Pemanfaatan sumber daya alam domestik yang tak habis memutus rantai ketergantungan pada fluktuasi harga impor bahan bakar fosil.", IsCorrect: true},
						{Type: "warrant", Text: "Keputusan kenaikan harga bahan bakar di tingkat konsumen sering kali memicu penolakan dan aksi unjuk rasa massa.", IsCorrect: false},
						{Type: "warrant", Text: "Anggaran subsidi energi yang besar berisiko memperlebar defisit fiskal belanja tahunan pemerintah pusat.", IsCorrect: false},
						{Type: "warrant", Text: "Pengawasan perbatasan maritim dan darat membutuhkan armada patroli yang selalu siaga menghadapi potensi sabotase.", IsCorrect: false},
					},
				},
			},
		},

		// Materi 3: Etika Kecerdasan Buatan di Pendidikan
		{
			Title: "Etika dan Regulasi Kecerdasan Buatan (AI) di Pendidikan",
			Content: `Integrasi model kecerdasan buatan generatif (LLM) dalam dunia pendidikan membuka peluang besar sekaligus tantangan etika yang kompleks.
Batasan antara pemanfaatan AI sebagai alat bantu pemantik pemikiran kritis dengan tindakan plagiarisme kognitif menjadi semakin kabur.
Diperlukan regulasi yang jelas terkait transparansi, mitigasi bias algoritma, dan perlindungan integritas akademik generasi pembelajar. Pemanasan global menuntut transformasi menyeluruh dari pembangkit berbahan bakar fosil menuju sumber energi baru terbarukan.
Meskipun menjanjikan pengurangan emisi gas rumah kaca secara drastis, transisi ini menghadapi tantangan teknis intermitensi suplai cuaca,
keandalan jaringan listrik tegangan tinggi, serta kalkulasi keekonomian penerapan instrumen disinsentif fiskal seperti pajak karbon. Pemanasan global menuntut transformasi menyeluruh dari pembangkit berbahan bakar fosil menuju sumber energi baru terbarukan.
Meskipun menjanjikan pengurangan emisi gas rumah kaca secara drastis, transisi ini menghadapi tantangan teknis intermitensi suplai cuaca,
keandalan jaringan listrik tegangan tinggi, serta kalkulasi keekonomian penerapan instrumen disinsentif fiskal seperti pajak karbon.`,
			Owner: asesor2,
			Arguments: []SeedArgument{
				{
					ClaimText: "Transparansi pengungkapan penggunaan AI generatif dalam karya tulis akademik melindungi integritas ilmiah dan orisinalitas pemikiran.",
					OrderNo:   1,
					Options: []SeedOption{
						{Type: "ground", Text: "Survei komite etik universitas mendapati 85% institusi yang menerapkan pedoman sitasi AI formal berhasil menurunkan kasus sengketa atribusi dan penulisan sembunyi-sembunyi.", IsCorrect: true},
						{Type: "ground", Text: "Aplikasi generator teks berbasis LLM dapat menyusun kerangka esai dalam hitungan beberapa detik saja.", IsCorrect: false},
						{Type: "ground", Text: "Alat pendeteksi AI otomatis sering menghasilkan kekeliruan pembacaan positif palsu pada tulisan berbahasa non-Inggris.", IsCorrect: false},
						{Type: "ground", Text: "Banyak mahasiswa menggunakan teknologi digital untuk merangkum literatur jurnal yang berbahasa asing.", IsCorrect: false},
						{Type: "warrant", Text: "Kejelasan dokumentasi metodologis membedakan kontribusi analitis orisinal penulis dari sintesis mesin sehingga kebenaran ilmiah tetap dapat diverifikasi.", IsCorrect: true},
						{Type: "warrant", Text: "Kemudahan akses internet nirkabel di kampus mempercepat penyelesaian tugas akhir bagi mahasiswa tingkat akhir.", IsCorrect: false},
						{Type: "warrant", Text: "Perkembangan algoritma model bahasa besar terus disempurnakan dengan penambahan triliunan data teks.", IsCorrect: false},
						{Type: "warrant", Text: "Penguasaan bahasa asing merupakan keunggulan kompetitif yang dibutuhkan lulusan perguruan tinggi di era globalisasi.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Ketergantungan penuh pada luaran AI tanpa evaluasi skeptis melemahkan daya penalaran reflektif dan keterampilan berpikir mandiri siswa.",
					OrderNo:   2,
					Options: []SeedOption{
						{Type: "ground", Text: "Penelitian psikologi kognitif menemukan siswa yang hanya menyalin luaran AI tanpa merevisi memiliki skor analisis kesalahan logika 38% lebih rendah pada ujian tertulis mandiri.", IsCorrect: true},
						{Type: "ground", Text: "Kecerdasan buatan dapat memproses jutaan halaman buku referensi dan menyajikannya dalam format tanya-jawab interaktif.", IsCorrect: false},
						{Type: "ground", Text: "Sistem pembelajaran daring modern kini telah mengintegrasikan fitur asisten virtual untuk membantu kendala teknis siswa.", IsCorrect: false},
						{Type: "ground", Text: "Beban kurikulum yang padat membuat siswa mencari cara tercepat untuk mengumpulkan tugas tepat waktu.", IsCorrect: false},
						{Type: "warrant", Text: "Proses kognitif internal yang aktif dalam menguji dan meragukan konsep sangat diperlukan untuk membentuk pemahaman konsep yang mendalam dan kokoh.", IsCorrect: true},
						{Type: "warrant", Text: "Asisten virtual chatbot yang responsif meringankan beban kerja staf administrasi akademik di luar jam kerja reguler.", IsCorrect: false},
						{Type: "warrant", Text: "Kapasitas pemrosesan komputasi awan memungkinkan respons instan terhadap ribuan kueri siswa secara bersamaan.", IsCorrect: false},
						{Type: "warrant", Text: "Ketepatan waktu pengumpulan tugas merupakan salah satu komponen penilaian kedisiplinan dalam buku pedoman akademik.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Audit independen terhadap dataset pelatihan AI pendidikan penting untuk mencegah pelestarian bias sosio-kultural dan stereotipe diskriminatif.",
					OrderNo:   3,
					Options: []SeedOption{
						{Type: "ground", Text: "Audit algoritmik independen mengungkap model evaluasi esai yang tidak diaudit memberi bobot nilai 18% lebih rendah pada dialek bahasa daerah tertentu meski argumennya valid secara substansi.", IsCorrect: true},
						{Type: "ground", Text: "Perusahaan pengembang AI bersaing ketat untuk merilis model dengan kemampuan komputasi dan jendela konteks terpanjang.", IsCorrect: false},
						{Type: "ground", Text: "Biaya pelatihan model kecerdasan buatan fondasional memerlukan investasi jutaan dolar untuk infrastruktur cip grafis khusus.", IsCorrect: false},
						{Type: "ground", Text: "Banyak sekolah swasta telah mengadopsi perangkat lunak berbasis kecerdasan buatan untuk mengotomatisasi absensi siswa.", IsCorrect: false},
						{Type: "warrant", Text: "Dataset historis yang tidak dikurasi kritis merefleksikan prasangka masa lalu, sehingga verifikasi eksternal mutlak diperlukan guna menjamin keadilan evaluasi bagi seluruh siswa.", IsCorrect: true},
						{Type: "warrant", Text: "Persaingan inovasi teknologi antara perusahaan rintisan mendorong lahirnya beragam alternatif produk perangkat lunak di pasaran.", IsCorrect: false},
						{Type: "warrant", Text: "Peningkatan efisiensi komputasi perangkat keras silikon dapat menghemat konsumsi energi pada pusat data.", IsCorrect: false},
						{Type: "warrant", Text: "Penggunaan kamera pengenal wajah untuk absensi harian dapat mengurangi kecurangan penitipan presensi antar teman sekelas.", IsCorrect: false},
					},
				},
				{
					ClaimText: "Penggunaan AI sebagai mitra dialog dialektis (Socratic partner) justru memperkuat kemampuan perumusan sanggahan (rebuttal) siswa.",
					OrderNo:   4,
					Options: []SeedOption{
						{Type: "ground", Text: "Uji coba terkontrol di kelas debat menunjukkan siswa yang berlatih dengan AI berorientasi pertanyaan klarifikasi menghasilkan sanggahan yang 45% lebih komprehensif.", IsCorrect: true},
						{Type: "ground", Text: "Berbicara di depan umum memerlukan rasa percaya diri dan teknik vokal intonasi yang baik.", IsCorrect: false},
						{Type: "ground", Text: "Modul pembelajaran debat kompetitif biasanya membatasi durasi penyampaian argumen setiap pembicara maksimal tujuh menit.", IsCorrect: false},
						{Type: "ground", Text: "Penggunaan mikrofon nirkabel membantu juri mendengar artikulasi kata pembicara dengan lebih jelas di dalam aula besar.", IsCorrect: false},
						{Type: "warrant", Text: "Umpan balik yang memancing pertimbangan bukti tandingan melatih peserta didik memperkuat tautan antara fakta dasar (ground) dan prinsip pembenaran (warrant).", IsCorrect: true},
						{Type: "warrant", Text: "Latihan pernapasan diafragma sebelum tampil di atas panggung membantu mengurangi kegugupan peserta lomba.", IsCorrect: false},
						{Type: "warrant", Text: "Manajemen waktu yang disiplin dalam kompetisi debat mencegah pemotongan poin oleh dewan juri.", IsCorrect: false},
						{Type: "warrant", Text: "Tata akustik ruangan yang tertata rapi meminimalkan gema suara yang dapat mengganggu konsentrasi penonton.", IsCorrect: false},
					},
				},
			},
		},
	}

	createdMaterials := make(map[string]*material.Material)
	for _, mCfg := range materialsConfig {
		var mat *material.Material
		var found material.Material
		if err := gormDB.Where("title = ? AND owner_id = ?", mCfg.Title, mCfg.Owner.ID).First(&found).Error; err == nil {
			mat = &found
		} else {
			mat = &material.Material{
				Title:   mCfg.Title,
				Content: mCfg.Content,
				OwnerID: mCfg.Owner.ID,
			}
			if err := materialRepo.CreateMaterial(mat); err != nil {
				log.Printf("failed to create material %s: %v", mCfg.Title, err)
				continue
			}
			log.Printf("📚 Created material: %s (Owner: %s)", mat.Title, mCfg.Owner.Name)

			for _, aCfg := range mCfg.Arguments {
				arg := &material.Argument{
					MaterialID: mat.ID,
					ClaimText:  aCfg.ClaimText,
					OrderNo:    aCfg.OrderNo,
				}
				var opts []material.Option
				for _, o := range aCfg.Options {
					opts = append(opts, material.Option{
						Type:      o.Type,
						Text:      o.Text,
						IsCorrect: o.IsCorrect,
					})
				}
				_ = materialRepo.CreateArgument(arg, opts)
			}
		}
		createdMaterials[mCfg.Title] = mat
	}

	mat1 := createdMaterials["Pentingnya Berpikir Kritis dalam Era Informasi"]
	mat2 := createdMaterials["Transisi Energi Terbarukan Menuju Net-Zero Emission"]
	mat3 := createdMaterials["Etika dan Regulasi Kecerdasan Buatan (AI) di Pendidikan"]

	cls1 := createdClasses["Kelas Argumentasi X-A"]
	cls2 := createdClasses["Kelas Literasi & Logika XI-IPA"]
	cls3 := createdClasses["Kelas Filsafat Kritis XII-IPS"]

	// -------------------------------------------------------------
	// 5. SEED EXAMS (4 Paket Ujian: > 2 Paket)
	// -------------------------------------------------------------
	type SeedExam struct {
		Title               string
		Material            *material.Material
		Owner               *user.User
		ArgumentsPerSession int
		Classes             []*class.Class
		DirectStudents      []*user.User
	}

	examsConfig := []SeedExam{
		{
			Title:               "Latihan Argumentasi Toulmin: Literasi Digital",
			Material:            mat1,
			Owner:               asesor1,
			ArgumentsPerSession: 3,
			Classes:             []*class.Class{cls1, cls2},
		},
		{
			Title:               "Uji Pemahaman: Transisi Energi Bersih",
			Material:            mat2,
			Owner:               asesor1,
			ArgumentsPerSession: 3,
			Classes:             []*class.Class{cls2},
		},
		{
			Title:               "Evaluasi Kritis: Etika Kecerdasan Buatan (AI)",
			Material:            mat3,
			Owner:               asesor2,
			ArgumentsPerSession: 3,
			Classes:             []*class.Class{cls3},
		},
		{
			Title:               "Remedial Khusus: Literasi Digital (Mandiri)",
			Material:            mat1,
			Owner:               asesor1,
			ArgumentsPerSession: 2,
			DirectStudents:      []*user.User{siswa5},
		},
	}

	createdExams := make(map[string]*exam.Exam)
	for _, eCfg := range examsConfig {
		if eCfg.Material == nil || eCfg.Owner == nil {
			continue
		}
		var ex *exam.Exam
		var found exam.Exam
		if err := gormDB.Where("title = ? AND owner_id = ?", eCfg.Title, eCfg.Owner.ID).First(&found).Error; err == nil {
			ex = &found
		} else {
			ex = &exam.Exam{
				Title:               eCfg.Title,
				MaterialID:          eCfg.Material.ID,
				OwnerID:             eCfg.Owner.ID,
				IsActive:            true,
				ArgumentsPerSession: eCfg.ArgumentsPerSession,
			}
			var classIDs []uint
			for _, c := range eCfg.Classes {
				if c != nil {
					classIDs = append(classIDs, c.ID)
				}
			}
			var studentIDs []uint
			for _, s := range eCfg.DirectStudents {
				if s != nil {
					studentIDs = append(studentIDs, s.ID)
				}
			}
			if err := examRepo.Create(ex, classIDs, studentIDs); err != nil {
				log.Printf("failed to create exam %s: %v", eCfg.Title, err)
				continue
			}
			log.Printf("📝 Created exam: %s", ex.Title)
		}
		createdExams[eCfg.Title] = ex
	}

	// -------------------------------------------------------------
	// 6. SEED SESSIONS & ATTEMPTS (For Exam 1, to activate Peer Analytics: > 3 Peers)
	// -------------------------------------------------------------
	exam1 := createdExams["Latihan Argumentasi Toulmin: Literasi Digital"]
	if exam1 != nil && mat1 != nil {
		// Dapatkan argument dan options dari mat1
		args, err := sessionRepo.GetArgumentsByMaterialID(mat1.ID)
		if err == nil && len(args) >= 3 {
			// Siswa 1 s/d Siswa 4 (4 unique peers >= MinPeersDefault 3)
			type StudentAttemptPlan struct {
				Student   *user.User
				Mode      session.SessionMode
				Status    session.SessionStatus
				Completed bool
				// Drops per argumen: map index argumen -> list pilihan yang di-drop
				Drops [][]struct {
					Slot      string
					IsCorrect bool
					DistIndex int
				}
			}

			plan := []StudentAttemptPlan{
				{
					Student:   siswa1,
					Mode:      session.ModeStandard,
					Status:    session.StatusCompleted,
					Completed: true,
					Drops: [][]struct {
						Slot      string
						IsCorrect bool
						DistIndex int
					}{
						// Arg 1: Drop ground salah (dist 0), lalu ground benar, lalu warrant benar
						{
							{Slot: "ground", IsCorrect: false, DistIndex: 0},
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
						// Arg 2: Langsung ground benar & warrant benar
						{
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
						// Arg 3: Drop warrant salah (dist 1), lalu ground benar, lalu warrant benar
						{
							{Slot: "warrant", IsCorrect: false, DistIndex: 1},
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
					},
				},
				{
					Student:   siswa2,
					Mode:      session.ModeSocial,
					Status:    session.StatusCompleted,
					Completed: true,
					Drops: [][]struct {
						Slot      string
						IsCorrect bool
						DistIndex int
					}{
						// Arg 1: Drop ground salah (dist 1), lalu ground benar & warrant benar
						{
							{Slot: "ground", IsCorrect: false, DistIndex: 1},
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
						// Arg 2: Drop ground salah (dist 0), drop warrant salah (dist 0), lalu yang benar
						{
							{Slot: "ground", IsCorrect: false, DistIndex: 0},
							{Slot: "warrant", IsCorrect: false, DistIndex: 0},
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
						// Arg 3: Langsung benar
						{
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
					},
				},
				{
					Student:   siswa3,
					Mode:      session.ModeHelp,
					Status:    session.StatusCompleted,
					Completed: true,
					Drops: [][]struct {
						Slot      string
						IsCorrect bool
						DistIndex int
					}{
						// Arg 1: Drop warrant salah (dist 2), lalu ground & warrant benar
						{
							{Slot: "warrant", IsCorrect: false, DistIndex: 2},
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
						// Arg 2: Langsung benar
						{
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
						// Arg 3: Drop ground salah (dist 2), lalu ground & warrant benar
						{
							{Slot: "ground", IsCorrect: false, DistIndex: 2},
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
					},
				},
				{
					Student:   siswa4,
					Mode:      session.ModeStandard,
					Status:    session.StatusInProgress,
					Completed: false,
					Drops: [][]struct {
						Slot      string
						IsCorrect bool
						DistIndex int
					}{
						// Hanya menyelesaikan Arg 1
						{
							{Slot: "ground", IsCorrect: true},
							{Slot: "warrant", IsCorrect: true},
						},
					},
				},
			}

			startTime := time.Now().Add(-2 * time.Hour)
			for _, p := range plan {
				if p.Student == nil {
					continue
				}

				// Cek apakah sudah ada session untuk student ini di exam1
				var existingSess session.Session
				if err := gormDB.Where("exam_id = ? AND student_id = ?", exam1.ID, p.Student.ID).First(&existingSess).Error; err == nil {
					continue
				}

				sess := &session.Session{
					ExamID:    exam1.ID,
					StudentID: p.Student.ID,
					AttemptNo: 1,
					Mode:      p.Mode,
					Status:    p.Status,
					StartedAt: startTime,
				}
				if p.Completed {
					compTime := startTime.Add(15 * time.Minute)
					sess.CompletedAt = &compTime
				}

				argIDs := []uint{args[0].ID, args[1].ID, args[2].ID}
				if err := sessionRepo.CreateSession(sess, argIDs); err != nil {
					log.Printf("failed to create seed session for student %s: %v", p.Student.Username, err)
					continue
				}

				// Log drops & progress
				for argIdx, argDrops := range p.Drops {
					targetArg := args[argIdx]

					// Filter opsi ground & warrant
					var correctGround, correctWarrant *material.Option
					var groundDistractors, warrantDistractors []*material.Option
					for i := range targetArg.Options {
						opt := &targetArg.Options[i]
						if opt.Type == "ground" {
							if opt.IsCorrect {
								correctGround = opt
							} else {
								groundDistractors = append(groundDistractors, opt)
							}
						} else if opt.Type == "warrant" {
							if opt.IsCorrect {
								correctWarrant = opt
							} else {
								warrantDistractors = append(warrantDistractors, opt)
							}
						}
					}

					for _, drop := range argDrops {
						var targetOpt *material.Option
						if drop.Slot == "ground" {
							if drop.IsCorrect {
								targetOpt = correctGround
							} else if len(groundDistractors) > drop.DistIndex {
								targetOpt = groundDistractors[drop.DistIndex]
							}
						} else {
							if drop.IsCorrect {
								targetOpt = correctWarrant
							} else if len(warrantDistractors) > drop.DistIndex {
								targetOpt = warrantDistractors[drop.DistIndex]
							}
						}

						if targetOpt != nil {
							_ = sessionRepo.LogAttempt(&session.AttemptLog{
								SessionID:  sess.ID,
								ArgumentID: targetArg.ID,
								OptionID:   targetOpt.ID,
								Slot:       drop.Slot,
								CreatedAt:  startTime.Add(time.Duration(argIdx*3) * time.Minute),
							})
						}
					}

					// Jika argumen selesai, record progress
					_ = sessionRepo.MarkArgumentComplete(sess.ID, targetArg.ID)
				}

				log.Printf("📊 Seeded session & attempt logs for: %s (Status: %s, Mode: %s)", p.Student.Name, p.Status, p.Mode)
			}
			log.Println("✨ Analytics data seeded successfully for Exam 1 (4 unique peers, peer monitoring active)!")
		}
	}

	log.Println("==================================================")
	log.Println("✅ ALL SEEDING COMPLETED SUCCESSFULLY!")
	log.Println("Default user password: password123")
	log.Println("Users created: admin, asesor1, asesor2, siswa1, siswa2, siswa3, siswa4, siswa5")
	log.Println("Classes created: Kelas X-A, Kelas XI-IPA, Kelas XII-IPS")
	log.Println("Exams created: 4 paket ujian (> 2)")
	log.Println("Peer analytics ready on Exam 1: 4 unique peers (MinPeers = 3 threshold met)")
	log.Println("==================================================")
}
