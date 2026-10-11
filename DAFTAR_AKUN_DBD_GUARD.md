# Daftar Akun dan Hak Akses Sistem DBD Guard (Kota Depok)

Dokumen ini memuat seluruh data akun, peran pengguna (*Role-Based Access Control* - RBAC), kredensial login, cakupan wilayah, serta rute dashboard yang tersedia di dalam sistem **DBD Guard**.

---

## 1. Status Layanan Sistem Saat Ini

| Komponen | Status | Port / URL | Keterangan |
|---|---|---|---|
| **Database** | 🟢 Berjalan (Docker) | `localhost:5433` | PostgreSQL 15 + PostGIS (`dbd_guard`) |
| **Backend API** | 🟢 Berjalan (Golang Fiber) | `http://localhost:8080` | Endpoint Base: `/api/v1`, Health: `/health` |
| **Frontend Web** | 🟢 Berjalan (Next.js 16) | `http://localhost:3000` | Portal Publik & Dashboard Petugas |

> **Halaman Login Petugas & Admin**: [http://localhost:3000/login](http://localhost:3000/login)  
> **Password Default Seluruh Akun**: `password123`

---

## 2. Akun Demo Utama (7 Peran Sesuai PRD)

Berikut adalah akun cepat yang disiapkan khusus untuk mempermudah proses simulasi, demonstrasi, dan pengujian seluruh alur penanganan DBD Kota Depok:

| No | Peran (Role) | Nama Pengguna | Email Login | Password | Scope Wilayah | Rute Halaman |
|---|---|---|---|---|---|---|
| 1 | **Admin Developer** (`ROLE_ADMIN_DEVELOPER`) | Super Admin Developer | `dev@dbdguard.id` *(alt: `admin@dbdguard.id`)* | `password123` | Global (Seluruh Sistem) | [`/developer`](http://localhost:3000/developer) |
| 2 | **Dinas Kesehatan (Dinkes)** (`ROLE_DINKES`) | Dinas Kesehatan Kota Depok | `admin@depok.go.id` | `password123` | Global Kota Depok (11 Kecamatan) | [`/dinkes`](http://localhost:3000/dinkes) |
| 3 | **Kepala Puskesmas** (`ROLE_KEPALA_PUSKESMAS`) | Kepala Puskesmas Beji | `kepala.puskesmas@depok.go.id` | `password123` | Puskesmas Beji 1 | [`/puskesmas`](http://localhost:3000/puskesmas) |
| 4 | **Petugas Puskesmas (Surveilans/PE)** (`ROLE_PETUGAS_PUSKESMAS`) | Petugas PE Puskesmas Beji | `petugas.puskesmas@depok.go.id` | `password123` | Puskesmas Beji 1 | [`/puskesmas`](http://localhost:3000/puskesmas) |
| 5 | **Tim Fogging** (`ROLE_TIM_FOGGING`) | Tim Fogging Puskesmas Beji | `tim.fogging@depok.go.id` | `password123` | Puskesmas Beji 1 | [`/fogging`](http://localhost:3000/fogging) |
| 6 | **Kader Wilayah RW** (`ROLE_KADER`) | Kader Wilayah Kelurahan Beji | `kader@depok.go.id` | `password123` | Kelurahan Beji (Wilayah Kerja) | [`/kader`](http://localhost:3000/kader) |
| 7 | **Warga (Pelapor)** (`ROLE_WARGA`) | Warga Pelapor Depok | `warga@gmail.com` *(Atau Lapor Anonim)* | `password123` | Publik Bebas | [`/`](http://localhost:3000) & [`/report`](http://localhost:3000/report) |

---

## 3. Rincian Hak Akses & Tanggung Jawab per Peran (RBAC)

### 1. Admin Developer (`ROLE_ADMIN_DEVELOPER`)
* **Login URL**: [http://localhost:3000/login](http://localhost:3000/login)
* **Akses**: `["all"]`
* **Fitur Utama**:
  * Manajemen akun dinas kesehatan dan monitoring sistem.
  * Audit log keamanan dan pengaturan parameter operasional platform.
  * Pemantauan status kesehatan server dan dependensi.

### 2. Dinas Kesehatan Kota Depok (`ROLE_DINKES`)
* **Login URL**: [http://localhost:3000/login](http://localhost:3000/login)
* **Akses**: `["read_all", "approve_stock"]`
* **Fitur Utama**:
  * GIS Heatmap sebaran kasus DBD dan jentik se-Kota Depok (11 Kecamatan, 63 Kelurahan).
  * Early Warning System (EWS) dan monitoring ambang batas KLB (Kejadian Luar Biasa).
  * Manajemen inventaris logistik terpusat (Gudang Utama Dinkes) dan persetujuan transfer stok ke Puskesmas.
  * Pemantauan SLA verifikasi kader (3x24 jam) dan respon PE puskesmas.

### 3. Kepala Puskesmas (`ROLE_KEPALA_PUSKESMAS`)
* **Login URL**: [http://localhost:3000/login](http://localhost:3000/login)
* **Akses**: `["verify", "pe", "fogging", "inventory"]` (Tingkat Manajerial)
* **Fitur Utama**:
  * Meninjau hasil Penyelidikan Epidemiologi (PE) dari Petugas Puskesmas.
  * Memberikan persetujuan (*approval*) atau penolakan jadwal *fogging fokus*.
  * Mengajukan permintaan stok logistik (insektisida, solar, abate) ke Dinkes.
  * Monitoring ketersediaan stok obat di gudang puskesmas setempat.

### 4. Petugas Puskesmas (`ROLE_PETUGAS_PUSKESMAS`)
* **Login URL**: [http://localhost:3000/login](http://localhost:3000/login)
* **Akses**: `["verify", "pe", "inventory"]` (Tingkat Operasional)
* **Fitur Utama**:
  * Menerima daftar laporan warga yang telah diverifikasi oleh Kader (atau dieskalasi otomatis oleh SLA).
  * Menginput hasil Penyelidikan Epidemiologi (PE) di lapangan:
    * Keberadaan jentik dalam radius 100 meter (Ya/Tidak).
    * Jumlah penderita demam di sekitar lokasi.
    * Foto bukti pemeriksaan jentik/lingkungan.
  * Melihat rekomendasi otomatis dari Decision Support System (DSS) sistem: *Fogging Fokus*, *Edukasi & PSN*, atau *Tertunda*.

### 5. Tim Fogging Lapangan (`ROLE_TIM_FOGGING`)
* **Login URL**: [http://localhost:3000/login](http://localhost:3000/login)
* **Akses**: `["fogging"]`
* **Fitur Utama**:
  * Tampilan ramah mobile (PWA view) untuk petugas di lapangan.
  * Melihat jadwal dan lokasi tugas fogging hari ini (Aturan bisnis: maksimal 1 penugasan per tim per hari).
  * Mengubah status penugasan secara real-time: **Menuju Lokasi (OTW)** → **Mulai Pelaksanaan** → **Selesai**.
  * Mengunggah foto bukti pelaksanaan fogging dan mencatat pemakaian logistik (misal: Malathion 2 Liter, Solar 10 Liter). Stok gudang otomatis terpotong saat penugasan ditandai selesai.

### 6. Kader Wilayah RW (`ROLE_KADER`)
* **Login URL**: [http://localhost:3000/login](http://localhost:3000/login)
* **Akses**: `["verify"]`
* **Fitur Utama**:
  * Notifikasi instan saat ada warga di lingkungan RW-nya yang mengajukan laporan.
  * Melakukan verifikasi faktual lapangan dalam tenggat waktu **SLA 3x24 Jam** (Valid / Ditolak).
  * Fitur deteksi dan penggabungan laporan ganda (*duplicate clustering*) dalam radius 100 meter.
  * Menerima alarm Early Warning System (EWS) jika ada potensi bahaya di wilayah sekitar.

### 7. Warga / Pelapor (`ROLE_WARGA` & Publik)
* **Portal Publik**: [http://localhost:3000](http://localhost:3000)
* **Formulir Laporan**: [http://localhost:3000/report](http://localhost:3000/report)
* **Tracking Status**: [http://localhost:3000/track](http://localhost:3000/track)
* **Akses**: Bebas (Tidak wajib login / tersedia opsi Lapor Anonim)
* **Fitur Utama**:
  * Melaporkan dugaan kasus DBD atau genangan jentik nyamuk lengkap dengan foto dan titik peta.
  * Memperoleh nomor resi pelaporan (*Reference ID*) unik (format: `DBD-DPK-YYYYXXXX`).
  * Melacak progres penanganan laporan dari verifikasi kader hingga fogging selesai.
  * Mengajukan banding jika laporan ditolak secara keliru oleh kader.
  * Memberikan ulasan dan rating bintang pasca-fogging selesai dilaksanakan.

---

## 4. Daftar Akun Lengkap Se-Kota Depok (194 Akun Aktif)

Basis data telah diisi secara komprehensif dengan **194 akun aktif** untuk seluruh fasilitas kesehatan dan wilayah administratif Kota Depok (Password default seluruh akun: `password123`):

### A. Tabel Akun 38 Puskesmas Se-Kota Depok
Setiap Puskesmas memiliki akun Kepala Puskesmas, Petugas PE, dan Tim Fogging.

| No | Nama Puskesmas | Kecamatan | Email Kepala Puskesmas | Email Petugas PE | Email Tim Fogging |
|---|---|---|---|---|---|
| 1 | Puskesmas Beji | Beji | `kepala.pusk1@depok.go.id` | `petugas.pusk1@depok.go.id` | `fogging.pusk1@depok.go.id` |
| 2 | Puskesmas Cipayung | Cipayung | `kepala.pusk2@depok.go.id` | `petugas.pusk2@depok.go.id` | `fogging.pusk2@depok.go.id` |
| 3 | Puskesmas Bojongsari | Bojongsari | `kepala.pusk3@depok.go.id` | `petugas.pusk3@depok.go.id` | `fogging.pusk3@depok.go.id` |
| 4 | Puskesmas Cilodong | Cilodong | `kepala.pusk4@depok.go.id` | `petugas.pusk4@depok.go.id` | `fogging.pusk4@depok.go.id` |
| 5 | Puskesmas Cimanggis | Cimanggis | `kepala.pusk5@depok.go.id` | `petugas.pusk5@depok.go.id` | `fogging.pusk5@depok.go.id` |
| 6 | Puskesmas Cinere | Cinere | `kepala.pusk6@depok.go.id` | `petugas.pusk6@depok.go.id` | `fogging.pusk6@depok.go.id` |
| 7 | Puskesmas Limo | Limo | `kepala.pusk7@depok.go.id` | `petugas.pusk7@depok.go.id` | `fogging.pusk7@depok.go.id` |
| 8 | Puskesmas Pancoran Mas | Pancoran Mas | `kepala.pusk8@depok.go.id` | `petugas.pusk8@depok.go.id` | `fogging.pusk8@depok.go.id` |
| 9 | Puskesmas Sawangan | Sawangan | `kepala.pusk9@depok.go.id` | `petugas.pusk9@depok.go.id` | `fogging.pusk9@depok.go.id` |
| 10 | Puskesmas Sukmajaya | Sukmajaya | `kepala.pusk10@depok.go.id` | `petugas.pusk10@depok.go.id` | `fogging.pusk10@depok.go.id` |
| 11 | Puskesmas Tapos | Tapos | `kepala.pusk11@depok.go.id` | `petugas.pusk11@depok.go.id` | `fogging.pusk11@depok.go.id` |
| 12 | Puskesmas Tanah Baru | Beji | `kepala.pusk12@depok.go.id` | `petugas.pusk12@depok.go.id` | `fogging.pusk12@depok.go.id` |
| 13 | Puskesmas Kemiri Muka | Beji | `kepala.pusk13@depok.go.id` | `petugas.pusk13@depok.go.id` | `fogging.pusk13@depok.go.id` |
| 14 | Puskesmas Duren Seribu | Bojongsari | `kepala.pusk14@depok.go.id` | `petugas.pusk14@depok.go.id` | `fogging.pusk14@depok.go.id` |
| 15 | Puskesmas Pondok Petir | Bojongsari | `kepala.pusk15@depok.go.id` | `petugas.pusk15@depok.go.id` | `fogging.pusk15@depok.go.id` |
| 16 | Puskesmas Kalimulya | Cilodong | `kepala.pusk16@depok.go.id` | `petugas.pusk16@depok.go.id` | `fogging.pusk16@depok.go.id` |
| 17 | Puskesmas Jatimulya | Cilodong | `kepala.pusk17@depok.go.id` | `petugas.pusk17@depok.go.id` | `fogging.pusk17@depok.go.id` |
| 18 | Puskesmas Harjamukti | Cimanggis | `kepala.pusk18@depok.go.id` | `petugas.pusk18@depok.go.id` | `fogging.pusk18@depok.go.id` |
| 19 | Puskesmas Pasir Gunung Selatan | Cimanggis | `kepala.pusk19@depok.go.id` | `petugas.pusk19@depok.go.id` | `fogging.pusk19@depok.go.id` |
| 20 | Puskesmas Tugu | Cimanggis | `kepala.pusk20@depok.go.id` | `petugas.pusk20@depok.go.id` | `fogging.pusk20@depok.go.id` |
| 21 | Puskesmas Gandul | Cinere | `kepala.pusk21@depok.go.id` | `petugas.pusk21@depok.go.id` | `fogging.pusk21@depok.go.id` |
| 22 | Puskesmas Pangkalan Jati | Cinere | `kepala.pusk22@depok.go.id` | `petugas.pusk22@depok.go.id` | `fogging.pusk22@depok.go.id` |
| 23 | Puskesmas Ratujaya | Cipayung | `kepala.pusk23@depok.go.id` | `petugas.pusk23@depok.go.id` | `fogging.pusk23@depok.go.id` |
| 24 | Puskesmas Bojong Pondok Terong | Cipayung | `kepala.pusk24@depok.go.id` | `petugas.pusk24@depok.go.id` | `fogging.pusk24@depok.go.id` |
| 25 | Puskesmas Krukut | Limo | `kepala.pusk25@depok.go.id` | `petugas.pusk25@depok.go.id` | `fogging.pusk25@depok.go.id` |
| 26 | Puskesmas Meruyung | Limo | `kepala.pusk26@depok.go.id` | `petugas.pusk26@depok.go.id` | `fogging.pusk26@depok.go.id` |
| 27 | Puskesmas Depok Jaya | Pancoran Mas | `kepala.pusk27@depok.go.id` | `petugas.pusk27@depok.go.id` | `fogging.pusk27@depok.go.id` |
| 28 | Puskesmas Rangkapan Jaya | Pancoran Mas | `kepala.pusk28@depok.go.id` | `petugas.pusk28@depok.go.id` | `fogging.pusk28@depok.go.id` |
| 29 | Puskesmas Mampang | Pancoran Mas | `kepala.pusk29@depok.go.id` | `petugas.pusk29@depok.go.id` | `fogging.pusk29@depok.go.id` |
| 30 | Puskesmas Pasir Putih | Sawangan | `kepala.pusk30@depok.go.id` | `petugas.pusk30@depok.go.id` | `fogging.pusk30@depok.go.id` |
| 31 | Puskesmas Pengasinan | Sawangan | `kepala.pusk31@depok.go.id` | `petugas.pusk31@depok.go.id` | `fogging.pusk31@depok.go.id` |
| 32 | Puskesmas Cinangka | Sawangan | `kepala.pusk32@depok.go.id` | `petugas.pusk32@depok.go.id` | `fogging.pusk32@depok.go.id` |
| 33 | Puskesmas Abadijaya | Sukmajaya | `kepala.pusk33@depok.go.id` | `petugas.pusk33@depok.go.id` | `fogging.pusk33@depok.go.id` |
| 34 | Puskesmas Bakti Jaya | Sukmajaya | `kepala.pusk34@depok.go.id` | `petugas.pusk34@depok.go.id` | `fogging.pusk34@depok.go.id` |
| 35 | Puskesmas Sukamaju Baru | Tapos | `kepala.pusk35@depok.go.id` | `petugas.pusk35@depok.go.id` | `fogging.pusk35@depok.go.id` |
| 36 | Puskesmas Jatijajar | Tapos | `kepala.pusk36@depok.go.id` | `petugas.pusk36@depok.go.id` | `fogging.pusk36@depok.go.id` |
| 37 | Puskesmas Sukatani | Tapos | `kepala.pusk37@depok.go.id` | `petugas.pusk37@depok.go.id` | `fogging.pusk37@depok.go.id` |
| 38 | Puskesmas Cilangkap | Tapos | `kepala.pusk38@depok.go.id` | `petugas.pusk38@depok.go.id` | `fogging.pusk38@depok.go.id` |

> *Catatan:* Puskesmas Beji (No. 1) juga dapat diakses menggunakan email demo cepat: `kepala.puskesmas@depok.go.id`, `petugas.puskesmas@depok.go.id`, dan `tim.fogging@depok.go.id`.

---

### B. Tabel Akun 63 Kader Kelurahan (Tingkat RW Se-Kota Depok)
Setiap Kelurahan memiliki akun Kader Wilayah yang terhubung ke Puskesmas pengampunya:

| No | Kelurahan | Kecamatan | Puskesmas Pengampu | Email Login Kader |
|---|---|---|---|---|
| 1 | Kelurahan Beji | Beji | Puskesmas Beji | `kader.kel1@dbdguard.id` |
| 2 | Kelurahan Beji Timur | Beji | Puskesmas Beji | `kader.kel2@dbdguard.id` |
| 3 | Kelurahan Kemiri Muka | Beji | Puskesmas Kemiri Muka | `kader.kel3@dbdguard.id` |
| 4 | Kelurahan Kukusan | Beji | Puskesmas Beji | `kader.kel4@dbdguard.id` |
| 5 | Kelurahan Pondok Cina | Beji | Puskesmas Beji | `kader.kel5@dbdguard.id` |
| 6 | Kelurahan Tanah Baru | Beji | Puskesmas Tanah Baru | `kader.kel6@dbdguard.id` |
| 7 | Kelurahan Bojongsari | Bojongsari | Puskesmas Bojongsari | `kader.kel7@dbdguard.id` |
| 8 | Kelurahan Bojongsari Baru | Bojongsari | Puskesmas Bojongsari | `kader.kel8@dbdguard.id` |
| 9 | Kelurahan Curug | Bojongsari | Puskesmas Bojongsari | `kader.kel9@dbdguard.id` |
| 10 | Kelurahan Duren Mekar | Bojongsari | Puskesmas Bojongsari | `kader.kel10@dbdguard.id` |
| 11 | Kelurahan Duren Seribu | Bojongsari | Puskesmas Duren Seribu | `kader.kel11@dbdguard.id` |
| 12 | Kelurahan Pondok Petir | Bojongsari | Puskesmas Pondok Petir | `kader.kel12@dbdguard.id` |
| 13 | Kelurahan Serua | Bojongsari | Puskesmas Bojongsari | `kader.kel13@dbdguard.id` |
| 14 | Kelurahan Cilodong | Cilodong | Puskesmas Cilodong | `kader.kel14@dbdguard.id` |
| 15 | Kelurahan Jatimulya | Cilodong | Puskesmas Jatimulya | `kader.kel15@dbdguard.id` |
| 16 | Kelurahan Kalibaru | Cilodong | Puskesmas Cilodong | `kader.kel16@dbdguard.id` |
| 17 | Kelurahan Kalimulya | Cilodong | Puskesmas Kalimulya | `kader.kel17@dbdguard.id` |
| 18 | Kelurahan Sukamaju | Cilodong | Puskesmas Cilodong | `kader.kel18@dbdguard.id` |
| 19 | Kelurahan Cisalak Pasar | Cimanggis | Puskesmas Cimanggis | `kader.kel19@dbdguard.id` |
| 20 | Kelurahan Curug Cimanggis | Cimanggis | Puskesmas Cimanggis | `kader.kel20@dbdguard.id` |
| 21 | Kelurahan Harjamukti | Cimanggis | Puskesmas Harjamukti | `kader.kel21@dbdguard.id` |
| 22 | Kelurahan Mekarsari | Cimanggis | Puskesmas Cimanggis | `kader.kel22@dbdguard.id` |
| 23 | Kelurahan Pasir Gunung Selatan | Cimanggis | Puskesmas Pasir Gunung Selatan | `kader.kel23@dbdguard.id` |
| 24 | Kelurahan Tugu | Cimanggis | Puskesmas Tugu | `kader.kel24@dbdguard.id` |
| 25 | Kelurahan Cinere | Cinere | Puskesmas Cinere | `kader.kel25@dbdguard.id` |
| 26 | Kelurahan Gandul | Cinere | Puskesmas Gandul | `kader.kel26@dbdguard.id` |
| 27 | Kelurahan Pangkalan Jati | Cinere | Puskesmas Pangkalan Jati | `kader.kel27@dbdguard.id` |
| 28 | Kelurahan Pangkalan Jati Baru | Cinere | Puskesmas Cinere | `kader.kel28@dbdguard.id` |
| 29 | Kelurahan Bojong Pondok Terong | Cipayung | Puskesmas Bojong Pondok Terong | `kader.kel29@dbdguard.id` |
| 30 | Kelurahan Cipayung | Cipayung | Puskesmas Cipayung | `kader.kel30@dbdguard.id` |
| 31 | Kelurahan Cipayung Jaya | Cipayung | Puskesmas Cipayung | `kader.kel31@dbdguard.id` |
| 32 | Kelurahan Pondok Jaya | Cipayung | Puskesmas Cipayung | `kader.kel32@dbdguard.id` |
| 33 | Kelurahan Ratujaya | Cipayung | Puskesmas Ratujaya | `kader.kel33@dbdguard.id` |
| 34 | Kelurahan Grogol | Limo | Puskesmas Limo | `kader.kel34@dbdguard.id` |
| 35 | Kelurahan Krukut | Limo | Puskesmas Krukut | `kader.kel35@dbdguard.id` |
| 36 | Kelurahan Limo | Limo | Puskesmas Limo | `kader.kel36@dbdguard.id` |
| 37 | Kelurahan Meruyung | Limo | Puskesmas Meruyung | `kader.kel37@dbdguard.id` |
| 38 | Kelurahan Depok | Pancoran Mas | Puskesmas Pancoran Mas | `kader.kel38@dbdguard.id` |
| 39 | Kelurahan Depok Jaya | Pancoran Mas | Puskesmas Depok Jaya | `kader.kel39@dbdguard.id` |
| 40 | Kelurahan Mampang | Pancoran Mas | Puskesmas Mampang | `kader.kel40@dbdguard.id` |
| 41 | Kelurahan Pancoran Mas | Pancoran Mas | Puskesmas Pancoran Mas | `kader.kel41@dbdguard.id` |
| 42 | Kelurahan Rangkapan Jaya | Pancoran Mas | Puskesmas Rangkapan Jaya | `kader.kel42@dbdguard.id` |
| 43 | Kelurahan Rangkapan Jaya Baru | Pancoran Mas | Puskesmas Pancoran Mas | `kader.kel43@dbdguard.id` |
| 44 | Kelurahan Bedahan | Sawangan | Puskesmas Sawangan | `kader.kel44@dbdguard.id` |
| 45 | Kelurahan Cinangka | Sawangan | Puskesmas Cinangka | `kader.kel45@dbdguard.id` |
| 46 | Kelurahan Kedaung | Sawangan | Puskesmas Sawangan | `kader.kel46@dbdguard.id` |
| 47 | Kelurahan Pasir Putih | Sawangan | Puskesmas Pasir Putih | `kader.kel47@dbdguard.id` |
| 48 | Kelurahan Pengasinan | Sawangan | Puskesmas Pengasinan | `kader.kel48@dbdguard.id` |
| 49 | Kelurahan Sawangan | Sawangan | Puskesmas Sawangan | `kader.kel49@dbdguard.id` |
| 50 | Kelurahan Sawangan Baru | Sawangan | Puskesmas Sawangan | `kader.kel50@dbdguard.id` |
| 51 | Kelurahan Abadijaya | Sukmajaya | Puskesmas Abadijaya | `kader.kel51@dbdguard.id` |
| 52 | Kelurahan Bakti Jaya | Sukmajaya | Puskesmas Bakti Jaya | `kader.kel52@dbdguard.id` |
| 53 | Kelurahan Cisalak | Sukmajaya | Puskesmas Sukmajaya | `kader.kel53@dbdguard.id` |
| 54 | Kelurahan Mekar Jaya | Sukmajaya | Puskesmas Sukmajaya | `kader.kel54@dbdguard.id` |
| 55 | Kelurahan Sukmajaya | Sukmajaya | Puskesmas Sukmajaya | `kader.kel55@dbdguard.id` |
| 56 | Kelurahan Tirtajaya | Sukmajaya | Puskesmas Sukmajaya | `kader.kel56@dbdguard.id` |
| 57 | Kelurahan Cilangkap | Tapos | Puskesmas Cilangkap | `kader.kel57@dbdguard.id` |
| 58 | Kelurahan Cimpaeun | Tapos | Puskesmas Tapos | `kader.kel58@dbdguard.id` |
| 59 | Kelurahan Jatijajar | Tapos | Puskesmas Jatijajar | `kader.kel59@dbdguard.id` |
| 60 | Kelurahan Leuwinanggung | Tapos | Puskesmas Tapos | `kader.kel60@dbdguard.id` |
| 61 | Kelurahan Sukamaju Baru | Tapos | Puskesmas Sukamaju Baru | `kader.kel61@dbdguard.id` |
| 62 | Kelurahan Sukatani | Tapos | Puskesmas Sukatani | `kader.kel62@dbdguard.id` |
| 63 | Kelurahan Tapos | Tapos | Puskesmas Tapos | `kader.kel63@dbdguard.id` |

> *Catatan:* Kelurahan Beji (No. 1) juga dapat diakses menggunakan email demo cepat: `kader@depok.go.id`.

---

### C. Akun Sampel Warga Terdaftar (10 Akun)
Digunakan untuk simulasi pelaporan warga non-anonim dan pemberian rating ulasan pasca-fogging:

| No | Nama Pengguna | Email Login | Password | Peran |
|---|---|---|---|---|
| 1 | Warga 1 (Kota Depok) | `warga1@gmail.com` | `password123` | `ROLE_WARGA` |
| 2 | Warga 2 (Kota Depok) | `warga2@gmail.com` | `password123` | `ROLE_WARGA` |
| 3 | Warga 3 (Kota Depok) | `warga3@gmail.com` | `password123` | `ROLE_WARGA` |
| 4 | Warga 4 (Kota Depok) | `warga4@gmail.com` | `password123` | `ROLE_WARGA` |
| 5 | Warga 5 (Kota Depok) | `warga5@gmail.com` | `password123` | `ROLE_WARGA` |
| 6 | Warga 6 (Kota Depok) | `warga6@gmail.com` | `password123` | `ROLE_WARGA` |
| 7 | Warga 7 (Kota Depok) | `warga7@gmail.com` | `password123` | `ROLE_WARGA` |
| 8 | Warga 8 (Kota Depok) | `warga8@gmail.com` | `password123` | `ROLE_WARGA` |
| 9 | Warga 9 (Kota Depok) | `warga9@gmail.com` | `password123` | `ROLE_WARGA` |
| 10 | Warga 10 (Kota Depok) | `warga10@gmail.com` | `password123` | `ROLE_WARGA` |
> *Catatan:* Akun demo warga utama adalah `warga@gmail.com` (Bpk. Hendra Wijaya).

---

## 5. Alur Pengujian Cepat (*End-to-End Simulation*)

Untuk menguji alur kerja lengkap sistem dari hulu ke hilir:

1. **Langkah 1 (Warga)**: Buka [http://localhost:3000/report](http://localhost:3000/report), isi formulir laporan warga baru, lalu catat nomor tiketnya (misal: `DBD-DPK-20261001`).
2. **Langkah 2 (Kader)**: Login di [http://localhost:3000/login](http://localhost:3000/login) dengan `kader@depok.go.id` / `password123`. Buka dashboard verifikasi, lalu klik **Setujui / Verifikasi**.
3. **Langkah 3 (Petugas Puskesmas)**: Login dengan `petugas.puskesmas@depok.go.id` / `password123`. Lakukan input hasil PE (Jentik = "Ya", Demam = 3). Sistem akan menghasilkan rekomendasi *Fogging Fokus*.
4. **Langkah 4 (Kepala Puskesmas)**: Login dengan `kepala.puskesmas@depok.go.id` / `password123`. Klik menu persetujuan fogging dan setujui penugasan.
5. **Langkah 5 (Tim Fogging)**: Login dengan `tim.fogging@depok.go.id` / `password123`. Buka penugasan di [/fogging](http://localhost:3000/fogging), klik **Mulai** lalu **Selesaikan Tugas** dengan memasukkan bukti pemakaian insektisida.
6. **Langkah 6 (Dinkes)**: Login dengan `admin@depok.go.id` / `password123`. Buka [/dinkes](http://localhost:3000/dinkes) untuk memantau pembaruan peta GIS, status stok, dan log audit secara terpusat.
