# Product Requirements Document (PRD) Komprehensif: DBD Guard (Kota Depok)

## 1. Product Overview & Problem Statement
**DBD Guard** adalah platform digital terintegrasi untuk pelaporan kasus Demam Berdarah Dengue (DBD) warga, verifikasi kewilayahan oleh Kader, Penyelidikan Epidemiologi (PE) oleh Puskesmas, penilaian risiko wilayah, serta penjadwalan *fogging* dan pengelolaan inventaris yang transparan.

Sistem ini dirancang secara spesifik untuk lingkup **Pemerintah Kota Depok**, bertujuan untuk menggantikan alur kerja manual (seperti surat menyurat dan pesan teks tidak terstruktur) menjadi alur digital yang sistematis, transparan, dan berbasis *Risk Scoring*. Hal ini mempercepat tindakan preventif dan kuratif tanpa mengebiri wewenang klinis dari pihak Kepala Puskesmas setempat.

---

## 2. Goals & Non-Goals (Scope)

### Goals (Tujuan Sistem)
1. **Digitalisasi Pelaporan**: Mempermudah warga dalam melaporkan kasus DBD atau jentik secara real-time.
2. **Standardisasi Alur**: Memberikan pedoman aksi (*Decision Support System*) berbasis bukti lapangan.
3. **Sentralisasi Inventaris**: Memantau ketersediaan stok *fogging* se-Kota Depok untuk meminimalkan keterlambatan aksi akibat kehabisan logistik.
4. **Early Warning System (EWS) & KLB**: Memberikan alarm otomatis ketika kepadatan kasus di suatu wilayah mencapai ambang batas Kejadian Luar Biasa (KLB).
5. **Akurasi Data Kewilayahan**: Pemetaan dan notifikasi tepat sasaran menggunakan struktur data administratif Kecamatan → Kelurahan → RW.

### Non-Goals (Out of Scope / Di Luar Cakupan)
Agar fokus tetap pada penanganan terintegrasi, hal-hal berikut **TIDAK** ditangani oleh sistem:
1. **Sistem Rekam Medis (EMR)**: Sistem ini bukan aplikasi rekam medis elektronik. Data pasien dari rumah sakit di luar diagnosis awal untuk keperluan surveilans tidak disimpan secara detail.
2. **Manajemen Kepagawaian Dinkes/Puskesmas**: HRD, penggajian, dan manajemen presensi *bukan* bagian dari sistem.
3. **Penanganan Penyakit Menular Lainnya**: Saat ini difokuskan murni untuk DBD.
4. **Hardware/IoT Integration**: Tidak menangani perangkat otomatis seperti drone penyemprot nyamuk.

---

## 3. Target Users & Roles (RBAC)

Sistem beroperasi dengan model Role-Based Access Control (RBAC).

| Role | Tanggung Jawab Utama | Akses Sistem |
|---|---|---|
| **Warga (Pelapor)** | Melaporkan kasus, memantau status laporan (*tracking*), dan memberikan rating evaluasi pasca-fogging. Dapat melapor anonim. | Portal Publik (Web) |
| **Kader Wilayah** | Bertugas memverifikasi laporan tingkat **RW** (valid/tidak valid), serta menerima EWS dari wilayah sekitar. Menggantikan peran RT/RW secara terpadu. | Dashboard Kader |
| **Petugas Puskesmas** | Menerima laporan terverifikasi, melakukan input hasil Penyelidikan Epidemiologi (PE) sederhana, dan mengelola stok gudang tingkat Puskesmas. | Dashboard Puskesmas |
| **Kepala Puskesmas** | Memberikan persetujuan (*approval*) jadwal *fogging* berdasarkan rekomendasi sistem, dan mengajukan permintaan stok obat ke Dinas Kesehatan. | Dashboard Puskesmas (Manager) |
| **Tim Fogging** | Menerima penugasan eksekusi, meng-update status lapangan (OTW, Mulai, Selesai), serta **wajib** mengunggah foto bukti dan melaporkan volume logistik yang terpakai. | Web/PWA App (Mobile View) |
| **Dinas Kesehatan (Dinkes)** | Memantau seluruh sistem se-Kota Depok via GIS Heatmap, menyetujui distribusi stok logistik, mengatur *threshold* KLB, dan mengontrol performa Kader/Puskesmas. | Dashboard Dinkes |
| **Admin Developer** | Mengelola kredensial akun, pengaturan tampilan portal publik (banner, info), serta memantau kesehatan server. | Dashboard Developer |

---

## 4. Technology Stack & Architecture

Sistem dituntut berbiaya operasional sangat rendah (100% Gratis/Open-Source untuk pihak ketiga).
Pendekatan Arsitektur: **Modular Monolith**.

### Tech Stack
*   **Backend & API**: **Golang (Go)** dengan framework Fiber. Digunakan karena ringan, cepat, dan efisien secara *memory*.
*   **Database**: **PostgreSQL + PostGIS**. Untuk mendukung logika *spatial*/radius secara *native*.
*   **Frontend Web**: **Next.js (React)** dengan styling **Tailwind CSS**. Menggunakan mode PWA (*Progressive Web App*) agar dapat digunakan secara responsif di *mobile* tanpa butuh rilis native app pada tahap MVP.
*   **Maps & GIS**: **Leaflet.js + OpenStreetMap (OSM)**. Bebas lisensi, tanpa *API Key* Google Maps.
*   **Notifikasi Push**: **Firebase Cloud Messaging (FCM)**.
*   **Notifikasi Email**: **Gomail/SMTP** (gratis menggunakan SMTP server pemerintah kota/Gmail).
*   **File Storage**: **MinIO** (S3-compatible, *self-hosted*).

---

## 5. Functional Requirements & Workflow

### 5.1. Alur Laporan Utama
1. **Lapor & Track**: Warga melapor dengan mengisi formulir berjenjang (Kecamatan → Kelurahan → RW). Sistem memberikan *Reference ID* untuk pelacakan.
2. **Verifikasi Kader (SLA 3x24 Jam)**: Kader menerima notifikasi, meninjau, dan memverifikasi laporan. Jika >3x24 jam diabaikan, laporan **Otomatis Dieskalasi** ke Puskesmas. Jika ditolak, pelapor mendapat opsi **Banding**.
3. **Penyelidikan Epidemiologi (PE)**: Petugas Puskesmas menginput dua variabel kunci:
   - Ditemukan jentik di radius 100m? (Ya / Tidak)
   - Jumlah orang demam di sekitar radius? (Angka)
4. **Decision Support System (DSS)**: Sistem memberi rekomendasi berdasarkan input PE dan Ketersediaan Stok Puskesmas (Lihat bagian 5.2).
5. **Eksekusi & Sinkronisasi Inventaris**: Jika Kepala Puskesmas *approve* fogging, Tim Fogging menuju lokasi. Saat eksekusi "Selesai", Tim wajib memasukkan jumlah logistik (misal: Obat A 2 Liter). Stok Puskesmas otomatis terpotong saat itu juga.
6. **Evaluasi Warga**: Warga pelapor memberikan *feedback* bintang setelah area di-*fogging*.

### 5.2. Aturan Bisnis (Business Rules)
*   **Anti-Overbooking Jadwal**: Satu entitas Tim Fogging hanya bisa dijadwalkan **1 tugas per hari**.
*   **Rekomendasi Otomatis (DSS)**:
    *   Jika Jentik = "Ya" atau Warga Demam ≥ 3, dan Stok **Tersedia** → Rekomendasi: **Fogging Fokus** 🔴.
    *   Jika Jentik = "Ya" atau Warga Demam ≥ 3, dan Stok **Kurang** → Rekomendasi: **Tertunda (Auto-Draft Request Stok)** 🟡.
    *   Jika Jentik = "Tidak" dan Warga Demam = 0 → Rekomendasi: **Edukasi & PSN** 🟢.
*   **Early Warning System (EWS)**: Pemberitahuan "Bahaya" di *dashboard* Dinkes dan peringatan beruntun ke Kader sekitar RW terdampak jika parameter Kejadian Luar Biasa (KLB) terlewati.
*   **Clustering Duplikat**: Jika ada >1 laporan dalam radius <100 meter pada kurun 7 hari, akan dilabeli sebagai **Potensi Duplikat**, Kader bisa menggabungkannya menjadi 1 kasus kluster.

---

## 6. Data, Privacy, & Security Requirements

### Perlindungan Privasi (PII - Personally Identifiable Information)
1. **Opsi Anonim**: Warga dapat memilih opsi "Lapor Anonim" di mana email/nomor kontak tidak disave di *database*, tidak ada *push notification*, dan resi hanya bisa diakses lewat manual copy-paste (QR).
2. **Data Masking Publik**: Pada Dashboard/Heatmap Publik, detail pelapor dan alamat persis tidak boleh diekspos (misalnya hanya ditampilkan di radius kasar Kelurahan/RW). Hanya Dinkes & Puskesmas terkait yang boleh melihat *geotagging* akurat.
3. **Data Medis**: Foto medis (misal: surat RS) yang dilampirkan hanya bisa diakses oleh level Puskesmas dan Dinkes.

### Autentikasi
*   Seluruh portal *dashboard* (kecuali publik) mewajibkan autentikasi JWT *(JSON Web Tokens)* atau *Secure HttpOnly Cookies*.
*   Log Audit merekam aksi perubahan krusial (seperti modifikasi batas KLB, *approval* fogging, alokasi inventaris).

---

## 7. Non-Functional Requirements

### User Experience (UX) Standard
1. **Empty States**: Jika belum ada data/laporan (misal tabel kosong), *frontend* wajib menampilkan ilustrasi "*Belum ada data*" bukan membiarkan komponen *blank*.
2. **Loading States**: Saat menunggu respons API, tombol wajib dalam state *disabled/spinning* atau menampilkan tampilan *skeleton loading* untuk mencegah *double-submit*.
3. **Error Handling**: Semua *error* (baik 4xx maupun 5xx) wajib ditangkap dan dikomunikasikan via *Toast/Alert Notification* yang beradab (bukan JSON *raw error* ke layar).

### Performa & Skalabilitas
*   Endpoint pelaporan dan *read* publik (GIS Heatmap) dilindungi dengan *rate limiter* (*throttle*) untuk mencegah *spam* dan *DDoS*.
*   Kueri jarak radius geografis (100 meter) menggunakan indeks *GIST spatial* dari PostGIS untuk performa <200ms.

---

## 8. Future Improvements (Tahap Selanjutnya)

Komponen-komponen ini tidak termasuk dalam fase rilis awal (MVP) namun menjadi target evolusi sistem:
1. **Native Mobile App (Flutter)**: Pembungkusan PWA menjadi aplikasi terinstal di App Store dan Google Play Store.
2. **Integrasi Eksternal (API RSUD)**: Sinkronisasi data real-time dengan RSUD Kota Depok terkait pasien *rawat inap* DBD (sehingga laporan tak hanya mengandalkan inisiatif warga).
3. **Prediksi Wabah Berbasis AI/Cuaca**: Mengawinkan data *historical* laporan dengan curah hujan Kota Depok untuk memprediksi kluster wabah bulan berikutnya.
