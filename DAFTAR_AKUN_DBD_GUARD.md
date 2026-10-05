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

## 4. Daftar Akun Lengkap Se-Kota Depok (Seeded Data)

Selain akun demo utama di atas, basis data telah diisi secara komprehensif dengan akun untuk seluruh fasilitas kesehatan dan wilayah administratif Kota Depok:

### A. Akun 38 Puskesmas (Kepala, Petugas PE, Tim Fogging)
Format penamaan email standar:
* **Kepala Puskesmas**: `kepala.pusk1@depok.go.id` s/d `kepala.pusk38@depok.go.id`
* **Petugas PE**: `petugas.pusk1@depok.go.id` s/d `petugas.pusk38@depok.go.id`
* **Tim Fogging**: `fogging.pusk1@depok.go.id` s/d `fogging.pusk38@depok.go.id`
* **Password**: `password123`

Contoh 5 Puskesmas Pertama:
| Puskesmas | Kecamatan | Email Kepala | Email Petugas PE | Email Tim Fogging |
|---|---|---|---|---|
| Puskesmas Bojongsari 1 | Bojongsari | `kepala.pusk1@depok.go.id` | `petugas.pusk1@depok.go.id` | `fogging.pusk1@depok.go.id` |
| Puskesmas Cilodong 2 | Cilodong | `kepala.pusk2@depok.go.id` | `petugas.pusk2@depok.go.id` | `fogging.pusk2@depok.go.id` |
| Puskesmas Cimanggis 3 | Cimanggis | `kepala.pusk3@depok.go.id` | `petugas.pusk3@depok.go.id` | `fogging.pusk3@depok.go.id` |
| Puskesmas Cinere 4 | Cinere | `kepala.pusk4@depok.go.id` | `petugas.pusk4@depok.go.id` | `fogging.pusk4@depok.go.id` |
| Puskesmas Cipayung 5 | Cipayung | `kepala.pusk5@depok.go.id` | `petugas.pusk5@depok.go.id` | `fogging.pusk5@depok.go.id` |

### B. Akun 63 Kader Kelurahan (Kecamatan Se-Kota Depok)
Format penamaan email standar:
* **Email**: `kader.kel1@dbdguard.id` s/d `kader.kel63@dbdguard.id`
* **Password**: `password123`

### C. Akun Warga Terdaftar (10 Akun Sampel)
* **Email**: `warga1@gmail.com` s/d `warga10@gmail.com`
* **Password**: `password123`

---

## 5. Alur Pengujian Cepat (*End-to-End Simulation*)

Untuk menguji alur kerja lengkap sistem dari hulu ke hilir:

1. **Langkah 1 (Warga)**: Buka [http://localhost:3000/report](http://localhost:3000/report), isi formulir laporan warga baru, lalu catat nomor tiketnya (misal: `DBD-DPK-20261001`).
2. **Langkah 2 (Kader)**: Login di [http://localhost:3000/login](http://localhost:3000/login) dengan `kader@depok.go.id` / `password123`. Buka dashboard verifikasi, lalu klik **Setujui / Verifikasi**.
3. **Langkah 3 (Petugas Puskesmas)**: Login dengan `petugas.puskesmas@depok.go.id` / `password123`. Lakukan input hasil PE (Jentik = "Ya", Demam = 3). Sistem akan menghasilkan rekomendasi *Fogging Fokus*.
4. **Langkah 4 (Kepala Puskesmas)**: Login dengan `kepala.puskesmas@depok.go.id` / `password123`. Klik menu persetujuan fogging dan setujui penugasan.
5. **Langkah 5 (Tim Fogging)**: Login dengan `tim.fogging@depok.go.id` / `password123`. Buka penugasan di [/fogging](http://localhost:3000/fogging), klik **Mulai** lalu **Selesaikan Tugas** dengan memasukkan bukti pemakaian insektisida.
6. **Langkah 6 (Dinkes)**: Login dengan `admin@depok.go.id` / `password123`. Buka [/dinkes](http://localhost:3000/dinkes) untuk memantau pembaruan peta GIS, status stok, dan log audit secara terpusat.
