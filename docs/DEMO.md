# Panduan Demonstrasi Sistem DBD Guard (Kota Depok)
### Proyek PKL Dinas Komunikasi dan Informatika (Diskominfo) Kota Depok
**Judul**: *"Rancang Bangun Dashboard Monitoring dan Sistem Pendukung Keputusan Penanganan Demam Berdarah Dengue (DBD) Berbasis Web pada Dinas Komunikasi dan Informatika Kota Depok"*

---

> [!NOTE]
> Seluruh data yang digunakan dalam demo ini berlabel **DATA SIMULASI** (100% fiktif), tanpa menggunakan maupun meniru data riil warga guna mematuhi prinsip perlindungan privasi data pribadi (UU PDP).

---

## 1. Persiapan Singkat Sebelum Demo

Pastikan seluruh layanan sistem berjalan di mesin lokal:

| Layanan | Perintah / Status | Port / URL |
|---|---|---|
| **Database** (PostgreSQL 15 + PostGIS) | `docker compose up -d` | `localhost:5433` |
| **Backend API** (Go Fiber) | `cd backend && go run ./cmd/api/main.go` | `http://localhost:8080` |
| **Frontend Web** (Next.js 16) | `cd frontend-web && npm run dev` | `http://localhost:3000` |

### Reset Data Demo Cepat (Idempoten & Bersih)
Sebelum memulai presentasi di hadapan mentor, jalankan perintah satu baris berikut untuk mereset seluruh data ke kondisi awal yang konsisten:

**Windows (PowerShell):**
```powershell
.\scripts\seed_demo.ps1
```

**Linux / macOS:**
```bash
./scripts/seed_demo.sh
```
*Waktu eksekusi hanya butuh ~3 detik.*

---

## 2. Tabel Akun Demo (7 Peran RBAC)

Semua akun menggunakan kata sandi (*default password*): **`password123`**  
URL Halaman Login Petugas: [http://localhost:3000/login](http://localhost:3000/login)

| Peran (*Role*) | Nama Akun Demo | Email Login | Cakupan Wilayah | Rute Akses |
|---|---|---|---|---|
| **Dinas Kesehatan (Dinkes)** | Dinas Kesehatan Kota Depok | `admin@depok.go.id` | Seluruh Kota Depok (11 Kecamatan) | [`/dinkes`](http://localhost:3000/dinkes) |
| **Kepala Puskesmas** | Dr. H. Kepala Puskesmas Beji | `kepala.puskesmas@depok.go.id` | Puskesmas Beji | [`/puskesmas`](http://localhost:3000/puskesmas) |
| **Petugas Puskesmas (PE)** | Petugas PE Puskesmas Beji | `petugas.puskesmas@depok.go.id` | Puskesmas Beji | [`/puskesmas/pe-validation`](http://localhost:3000/puskesmas/pe-validation) |
| **Tim Fogging** | Tim Fogging Reaksi Cepat Beji | `tim.fogging@depok.go.id` | Puskesmas Beji | [`/fogging`](http://localhost:3000/fogging) |
| **Kader Wilayah RW** | Ibu Kader RW 04 Kelurahan Beji | `kader@depok.go.id` | RW 04, Kelurahan Beji | [`/kader`](http://localhost:3000/kader) |
| **Warga (Pelapor)** | Bpk. Hendra Wijaya | `warga@gmail.com` *(atau Anonim)* | Publik Bebas | [`/report`](http://localhost:3000/report) |
| **Admin Developer** | Super Admin Developer | `dev@dbdguard.id` | Global Sistem | [`/developer`](http://localhost:3000/developer) |

---

## 3. Rundown Presentasi Demo (7–10 Menit)

Alur presentasi dirancang runtut mengikuti siklus hidup penanganan kasus DBD:

```
[Warga Melapor] ➔ [Verifikasi Kader] ➔ [PE & SPK Puskesmas] ➔ [Approval Kepala Pusk] ➔ [Eksekusi Fogging & Potong Stok] ➔ [Monitoring Dinkes]
```

### Menit 00:00 – 01:30 | 1. Portal Publik Warga & Pelaporan
* **Akses**: Buka [http://localhost:3000](http://localhost:3000) dan [http://localhost:3000/report](http://localhost:3000/report).
* **Poin Demonstrasi**:
  1. Tunjukkan kemudahan formulir pelaporan berjenjang (*Dropdown kaku*: Kecamatan $\rightarrow$ Kelurahan $\rightarrow$ RW, tanpa free-text lokasi liar).
  2. Jelaskan fitur **Perlindungan Privasi Warga**: Opsi *Lapor Anonim* yang tidak menyimpan kontak/email di database.
  3. Tunjukkan halaman pelacakan status di [`/track`](http://localhost:3000/track) menggunakan Nomor Referensi (contoh demo: `DBD-DPK-2026-DEMO01`).

### Menit 01:30 – 03:00 | 2. Dashboard Kader RW & Verifikasi Wilayah
* **Akses**: Login sebagai `kader@depok.go.id` di [http://localhost:3000/login](http://localhost:3000/login) $\rightarrow$ diarahkan ke [`/kader`](http://localhost:3000/kader).
* **Poin Demonstrasi**:
  1. Kader menerima antrean laporan di wilayah kerjanya dengan tenggat **SLA 3×24 Jam**.
  2. **Skenario Ditolak Lintas Wilayah**: Tunjukkan laporan `DBD-DPK-2026-REJECT01` yang ditolak kader karena lokasi berada di luar batas administratif RW 04 / Kelurahan Beji, lengkap dengan alasan penolakan dan jejak audit.

### Menit 03:00 – 05:30 | 3. Sistem Pendukung Keputusan (SPK) Validasi PE Puskesmas
* **Akses**: Login sebagai `petugas.puskesmas@depok.go.id` $\rightarrow$ buka [`/puskesmas/pe-validation`](http://localhost:3000/puskesmas/pe-validation).
* **Poin Demonstrasi**:
  1. Klik tombol **Validasi** pada kasus `DBD-DPK-2026-SPK01`.
  2. Sorot panel transparan **"Mengapa rekomendasi ini?"**:
     - **Input 1: Jentik** = Positif / Ditemukan dalam radius 100m.
     - **Input 2: Penderita Demam** = 4 orang (Ambang batas $\ge 3$).
     - **Input 3: Stok Insektisida** = 23 Liter (Ambang batas $\ge 10$ Liter).
     - **Aturan Terpicu**: Menyoroti **Aturan 1 (Fogging Fokus)** dengan penjelasan berbasis pedoman teknis operasional Kemenkes, tanpa klaim AI *blackbox*.
  3. Tunjukkan fitur **Advisory Human Override**: Kepala Puskesmas tetap memiliki hak klinis melakukan *override* dengan kewajiban mengisi alasan resmi yang terekam pada audit trail (*Four-Eyes Principle*).

### Menit 05:30 – 07:00 | 4. Eksekusi Fogging & Pemotongan Stok Otomatis
* **Akses**: Login sebagai `tim.fogging@depok.go.id` $\rightarrow$ buka [`/fogging`](http://localhost:3000/fogging).
* **Poin Demonstrasi**:
  1. Tampilan ramah perangkat bergerak (*mobile-first PWA*).
  2. Penugasan 1 tugas per tim per hari (*Anti-Overbooking*).
  3. Alur status: Menuju Lokasi (OTW) $\rightarrow$ Mulai $\rightarrow$ Selesai.
  4. Saat klik **Selesai**, tim wajib mengunggah bukti foto dan mencatat volume logistik terpakai (2 Liter Insektisida). Stok Gudang Puskesmas Beji langsung terpotong secara *atomic* dari 25L menjadi 23L.

### Menit 07:00 – 09:00 | 5. Dashboard Monitoring Dinas Kesehatan (Dinkes)
* **Akses**: Login sebagai `admin@depok.go.id` $\rightarrow$ buka [`/dinkes`](http://localhost:3000/dinkes).
* **Poin Demonstrasi**:
  1. **Badge Data Simulasi**: Memastikan kepatuhan data demo sintetis.
  2. **5 Kartu Ringkasan**: Laporan Masuk, Menunggu Verifikasi, Dalam Penyelidikan, Fogging Selesai, dan Kepatuhan SLA 3x24 Jam.
  3. **Filter Interaktif**: Filter Kecamatan (11 Kecamatan Kota Depok) dan Filter Periode (7 Hari, 30 Hari, Semua).
  4. **Peta Sebaran & Status Siaga (EWS/KLB)**: Tunjukkan notifikasi peringatan KLB aktif pada **Kelurahan Beji** (6 Kasus Aktif) yang melebihi ambang batas risiko.
  5. **Panel Stok Insektisida Sentral**: Monitoring stok Gudang Utama Dinkes vs 38 Puskesmas, serta deteksi otomatis puskesmas berstatus stok menipis (< 10L, contoh: Puskesmas Cipayung).
  6. **Ringkasan Log Audit**: Tunjukkan jejak transparansi rekaman aksi pengguna terbaru.

### Menit 09:00 – 10:00 | 6. Penutup & Nilai Tambah bagi Diskominfo
* Sistem dibangun **100% Bebas Biaya Pihak Ketiga** (tanpa Google Maps API / Twilio berbayar).
* Menggunakan OpenStreetMap dengan konfigurasi Tile Layer yang fleksibel untuk dialihkan ke server tile lokal Diskominfo.
* Arsitektur modular yang stabil dan siap dievaluasi.

---

## 4. Rencana Cadangan (Contingency Plan)

| Potensi Kendala | Kemungkinan Penyebab | Tindakan Mitigasi Cepat Saat Presentasi |
|---|---|---|
| **Peta OSM Lambat / Internet Offline** | Jaringan tempat presentasi tidak stabil | Seluruh backend, kalkulasi PostGIS, marker, dan antarmuka berjalan 100% offline di `localhost`. Jika tile OSM publik lambat, peta tetap merender titik kasus. Sistem juga mendukung environment variable `NEXT_PUBLIC_MAP_TILE_URL` untuk diarahkan ke tile server lokal jika tersedia. |
| **Koneksi Database Terputus** | Docker desktop tertidur atau port konflik | Buka terminal dan jalankan `docker compose up -d` lalu periksa koneksi dengan `.\scripts\seed_demo.ps1`. |
| **Salah Klik / Data Demo Berantakan** | Human error saat simulasi | Buka terminal, jalankan `.\scripts\seed_demo.ps1`. Dalam 3 detik, seluruh database kembali bersih dan siap demo ulang. |
| **Pertanyaan Mentor tentang AI / Prediksi** | Mentor menanyakan model machine learning | Jelaskan secara tegas dan jujur bahwa sistem DBD Guard ini berfokus pada **Sistem Pendukung Keputusan Berbasis Aturan (Rule-Based DSS)** yang mengacu pada parameter pasti (Jentik, Demam $\ge 3$, Stok Insektisida $\ge 10$ L) sesuai pedoman surveilans Kemenkes, sehingga keputusannya dapat dipertanggungjawabkan (*explainable*) dan tidak menimbulkan risiko *hallucination*. |
