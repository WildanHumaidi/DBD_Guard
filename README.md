# DBD Guard: Sistem Monitoring, Pelaporan, & SPK Penanganan Demam Berdarah Dengue (DBD) Terintegrasi Kota Depok

Sistem informasi manajemen kesehatan terpadu berbasis web dan spasial (GIS) untuk penanganan demam berdarah di wilayah **Kota Depok** (11 Kecamatan, 63 Kelurahan, dan 38 Puskesmas). Mengintegrasikan alur pelaporan warga, verifikasi kader, penyelidikan epidemiologi (PE), sistem pendukung keputusan (SPK) berbasis aturan, penjadwalan fogging, hingga pemantauan stok logistik terpusat.

---

## 1. Ringkasan Eksekutif & Latar Belakang Masalah

Sebelum sistem DBD Guard diimplementasikan, alur penanganan kasus DBD di tingkat daerah menghadapi beberapa tantangan operasional:
1. **Koordinasi Manual & Parsial**: Warga melapor melalui grup WhatsApp atau surat pengantar RT/RW yang rawan tercecer dan lambat direspons.
2. **Ketiadaan Pemantauan SLA Ketat**: Tenggat respon verifikasi awal lapangan (SLA 3x24 jam) tidak terpantau secara transparan, seringkali menyebabkan penanganan terlambat saat jentik telah menyebar.
3. **Laporan Tumpang Tindih (Duplikasi Spasial)**: Banyak warga dalam radius 100 meter yang sama melaporkan satu kasus demam secara terpisah, memicu dobel penugasan petugas surveilans.
4. **Keputusan Fogging Subjektif**: Fogging seringkali diminta secara reaktif tanpa didukung bukti epidemiologi (keberadaan jentik dan kasus demam penyerta).
5. **Kebocoran & Diskrepansi Stok**: Penggunaan insektisida, bubuk abate, dan solar tidak terhubung langsung dengan penyelesaian tindakan di lapangan.

**Solusi DBD Guard**:
DBD Guard menghubungkan seluruh pemangku kepentingan (warga, kader RW, puskesmas, tim fogging, dan dinas kesehatan) ke dalam satu platform digital terpadu dengan aturan bisnis ketat, otomasi sistem pendukung keputusan, dan audit log lengkap.

---

## 2. Arsitektur Teknologi & Prinsip Desain

Sistem DBD Guard dibangun secara mandiri dengan **100% teknologi bebas biaya / open-source**, tanpa dependensi pada API pihak ketiga berbayar (seperti Google Maps API, Twilio, atau Fonnte WA Gateway):

| Lapisan Sistem | Teknologi / Pustaka | Peran & Alasan Pemilihan |
|---|---|---|
| **Backend API** | Go (Golang) 1.22+ & Fiber v2 | Kinerja tinggi, penggunaan memori rendah, konkurensi native (Goroutines) untuk proses notifikasi asinkron dan audit trail. |
| **Database** | PostgreSQL 15 + PostGIS | Penyimpanan data relasional dan kalkulasi geometri spasial native (`ST_DWithin`, `ST_DistanceSphere`) untuk deteksi duplikasi radius 100m dan pemetaan titik kasus. |
| **Frontend Web** | Next.js 16 (App Router, Turbopack) & React 19 | Server-side rendering (SSR), Progressive Web App (PWA) readiness, performa cepat dengan caching modern. |
| **Desain & UI** | Tailwind CSS v4, Lucide Icons | Antarmuka bersih, responsif lintas perangkat (Desktop & Mobile), kontras tinggi, ramah bagi petugas lapangan. |
| **Pemetaan (GIS)** | Leaflet.js & OpenStreetMap (OSM) | Peta interaktif bebas lisensi komersial dengan URL tile yang dapat dikonfigurasi (`NEXT_PUBLIC_MAP_TILE_URL`). |
| **Notifikasi** | Firebase Cloud Messaging (FCM) & SMTP Email | Push notification web/mobile dan pengiriman email konfirmasi status tanpa biaya lisensi bulanan. |
| **Object Storage** | MinIO / Local File Storage | Penyimpanan foto bukti laporan dan dokumentasi pelaksanaan fogging berstandar S3 API. |

---

## 3. Matriks Hak Akses & Peran Pengguna (RBAC 7 Peran)

Sistem menerapkan kontrol akses berbasis peran (Role-Based Access Control) dengan yurisdiksi wilayah kerja yang ketat:

| No | Peran (Role) | Kredensial Demo | Lingkup Wilayah | Rute Utama | Tanggung Jawab Operasional |
|---|---|---|---|---|---|
| 1 | **Warga (Pelapor)** | Publik Bebas / `warga@gmail.com` | Publik / Mandiri | `/`, `/report`, `/track` | Melaporkan dugaan DBD/jentik dengan foto & koordinat, melacak status via nomor referensi unik. |
| 2 | **Kader Lingkungan RW** | `kader@depok.go.id` | Kelurahan Kerja (misal: Beji) | `/kader`, `/kader/verification` | Menerima notifikasi warga, klaim tugas, verifikasi lapangan dalam batas SLA 3x24 jam, penggabungan laporan duplikat. |
| 3 | **Petugas Puskesmas (Surveilans/PE)** | `petugas.puskesmas@depok.go.id` | Puskesmas (misal: Beji) | `/puskesmas`, `/puskesmas/pe-validation` | Penyelidikan Epidemiologi (PE) radius 100m, input keberadaan jentik dan penderita demam, telaah rekomendasi SPK. |
| 4 | **Kepala Puskesmas** | `kepala.puskesmas@depok.go.id` | Puskesmas (misal: Beji) | `/puskesmas`, `/puskesmas/fogging-requests` | Memberikan persetujuan manajerial jadwal fogging, override rekomendasi SPK, pengajuan transfer stok obat. |
| 5 | **Tim Fogging Lapangan** | `tim.fogging@depok.go.id` | Puskesmas Wilayah | `/fogging`, `/fogging/schedules` | Tampilan mobile (PWA), pelaporan eksekusi fogging (OTW -> Mulai -> Selesai), unggah bukti foto & pemotongan stok riil. |
| 6 | **Dinas Kesehatan (Dinkes)** | `admin@depok.go.id` | Kota Depok (11 Kecamatan) | `/dinkes`, `/dinkes/map`, `/dinkes/klb`, `/dinkes/sla` | GIS Heatmap sebaran kasus se-Depok, pemantauan Early Warning System (EWS) KLB, audit SLA, gudang logistik utama. |
| 7 | **Admin Developer** | `dev@dbdguard.id` | Global Sistem | `/developer`, `/developer/status`, `/developer/audit-logs` | Pemantauan kesehatan server/database, pemeriksaan jejak audit keamanan, pemeliharaan master data & RBAC. |

> **Password Default Akun Demo**: `password123`

---

## 4. Alur Bisnis End-to-End (Siklus Penanganan Kasus)

```
[1. Warga Melapor]
       │  (Foto, Koordinat GPS, RT/RW, Alamat)
       ▼
[2. Deteksi Duplikat Spasial (100m)]
       │  (Jika duplikat -> ditandai untuk opsi gabung)
       ▼
[3. Notifikasi & Verifikasi Kader RW]
       │  (Tenggat Waktu: SLA 3x24 Jam)
       ├─── (Lewat 3x24 Jam) ──> [Eskalasi Otomatis ke Puskesmas]
       ▼
[4. Penyelidikan Epidemiologi (PE)]
       │  (Pemeriksaan Jentik & Kasus Demam Sekitar Radius 100 Meter)
       ▼
[5. Sistem Pendukung Keputusan (SPK/DSS)]
       │  (Evaluasi Aturan: Fogging Fokus vs Edukasi PSN vs Tertunda)
       ▼
[6. Persetujuan Kepala Puskesmas & Penjadwalan]
       │  (Maksimal 1 tugas per hari per tim fogging)
       ▼
[7. Eksekusi Lapangan & Pemotongan Stok Otomatis]
       │  (Unggah Foto Bukti & Pemakaian Insektisida/Solar/Abate)
       ▼
[8. Pelacakan Warga & Pemantauan Agregat Dinkes]
          (Status selesai terupdate, peta GIS & ambang KLB terhitung)
```

---

## 5. Fitur Unggulan Sistem

### A. Sistem Pendukung Keputusan Berbasis Aturan (Rule-Based DSS v1.0)
Sistem tidak menyerahkan keputusan penindakan hanya berdasarkan insting pribadi. Modul SPK mengevaluasi kondisi lapangan dengan pohon aturan transparan:
* **Aturan 1 (Fogging Fokus)**: Ditemukan jentik dalam radius 100m ATAU penderita demam lain $\ge 3$ orang, DAN stok insektisida di gudang puskesmas $\ge 10$ Liter.
* **Aturan 2 (Fogging Tertunda)**: Indikasi medis terpenuhi, namun stok insektisida $< 10$ Liter. Sistem memicu peringatan transfer stok segera.
* **Aturan 3 (Edukasi & PSN 3M Plus)**: Tidak ditemukan jentik dan penderita demam $< 3$ orang. Tindakan cukup melalui larvasidasi abate dan kerja bakti warga.
* **Advisory Override**: Kepala Puskesmas memiliki wewenang manajerial untuk melakukan *override* atas rekomendasi sistem dengan kewajiban mengisi justifikasi medis tertulis pada audit log.

### B. Early Warning System (EWS) & Ambang Batas KLB
* Memantau akumulasi kasus per kelurahan secara berkala.
* Status **WASPADA / KELURAHAN SIAGA** diaktifkan otomatis apabila terjadi lonjakan kasus $\ge 2\times$ lipat dibandingkan periode dasar, atau telah melampaui batas ambang bahaya ($\ge 5$ kasus aktif dalam kurun waktu 14 hari).

### C. Penggabungan Laporan Duplikat Spasial
* Memanfaatkan query PostGIS `ST_DWithin(geom, NEW_GEOM, 100)`.
* Apabila warga dalam satu blok perumahan melaporkan kasus yang sama, kader diberikan rekomendasi klastering untuk menggabungkan laporan tersebut guna mencegah inefisiensi surveilans.

### D. Manajemen Inventaris Terintegrasi & Anti-Kebocoran
* Alur perpindahan logistik: Gudang Utama Dinkes $\rightarrow$ Gudang Puskesmas $\rightarrow$ Penggunaan Lapangan.
* Stok bahan kimia (Malathion/Cypermethrin, Solar, Bubuk Abate) terpotong secara transaksional (*atomic database update*) tepat saat tim fogging menyelesaikan tugas di aplikasi dan mengunggah foto bukti.

---

## 6. Panduan Instalasi & Menjalankan Sistem Lokal

### Prasyarat:
- **Go**: Versi 1.22 atau lebih baru
- **Node.js**: Versi 20.x atau 22.x
- **PostgreSQL**: Versi 15+ dengan ekstensi **PostGIS** aktif (Port default lokal: `5432` atau `5433`)
- **Docker** (Opsional, jika menjalankan database via container)

### Langkah 1: Klon Repositori & Persiapan Direktori
```bash
git clone https://github.com/dbd-guard/dbd-guard.git
cd DBD_Guard
```

### Langkah 2: Menjalankan Database PostgreSQL + PostGIS
Jika menggunakan Docker Compose:
```bash
docker compose up -d
```
Atau pastikan PostgreSQL lokal aktif dan buat database:
```sql
CREATE DATABASE dbd_guard;
\c dbd_guard
CREATE EXTENSION IF NOT EXISTS postgis;
```

### Langkah 3: Konfigurasi Environment Backend
Salin file konfigurasi pada folder `backend`:
```bash
cd backend
cp .env.example .env
```
Sesuaikan konfigurasi koneksi basis data:
```env
APP_ENV=development
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=dbd_guard
JWT_SECRET=supersecretjwtkey_dbd_guard_2026
```

### Langkah 4: Menjalankan Database Seeder Idempoten
Sebelum menjalankan API, isi master data dan skenario demo:
```bash
go run cmd/seeder/main.go
```
*Hasil: Seluruh 38 Puskesmas, 63 Kelurahan, 7 akun pengguna demo, dan data simulasi kasus akan terisi dalam kurun waktu ~3 detik.*

### Langkah 5: Menjalankan Backend API Server
```bash
go run cmd/api/main.go
```
Backend berjalan pada: `http://localhost:8080` (Health Check: `http://localhost:8080/health`).

### Langkah 6: Konfigurasi & Menjalankan Frontend Web
Buka terminal baru:
```bash
cd frontend-web
npm install
npm run dev
```
Frontend berjalan pada: `http://localhost:3000`.

---

## 7. Panduan Menjalankan Seluruh Pengujian (Test Suites)

DBD Guard dilengkapi rangkaian pengujian berlapis untuk memastikan stabilitas kode dan integritas logika bisnis:

### A. Pengujian Unit & Integrasi Backend
Menjalankan seluruh test suite repository dan usecase:
```bash
cd backend
go test -v -count=1 ./...
```
Cakupan pengujian meliputi:
- Uji isolasi yurisdiksi wilayah (`report_repository_test.go`)
- Uji pencegahan race condition klaim tugas kader (`bug08_claim_race_test.go`)
- Uji pemotongan stok fogging (`bug02_fogging_stock_test.go`)
- Uji matriks aturan rekomendasi SPK (`bug13_dss_decision_matrix_test.go`)
- Uji eskalasi SLA 3x24 jam (`sla_usecase_test.go`)
- Uji agregasi ambang batas KLB (`bug08_area_aggregation_test.go`)

### B. Pengujian Simulasi Lifecycle E2E Backend
Menjalankan simulasi skenario hulu-hilir (Warga $\rightarrow$ Kader $\rightarrow$ Puskesmas $\rightarrow$ Fogging $\rightarrow$ KLB):
```bash
cd backend
go run cmd/e2e/main.go
```

### C. Verifikasi Kompilasi & Tipe Data Frontend
Menjalankan validasi TypeScript dan pembuatan build produksi Next.js:
```bash
cd frontend-web
npm run build
```

---

## 8. Bukti Hasil Pengujian End-to-End (Desktop & Mobile)

Rangkaian pengujian antarmuka end-to-end telah dijalankan pada 2 viewport (Desktop 1280x800 dan Mobile 390x844). Sebanyak **62 bukti screenshot (31 desktop + 31 mobile)** telah didokumentasikan di direktori berikut:

- **Desktop Viewport**: `screenshots/e2e_evidence/desktop/`
- **Mobile Viewport**: `screenshots/e2e_evidence/mobile/`

### Indeks Bukti Pengujian Utama:
1. `01_home_desktop.png` / `01_home_mobile.png` — Beranda Utama Portal DBD Guard Kota Depok
2. `02_form_lapor_warga.png` — Formulir Pengaduan Warga (Gejala, Peta, RT/RW)
3. `03_lacak_laporan_warga.png` — Halaman Tracking Status Laporan Real-Time
4. `04_portal_transparansi_publik.png` — Portal Transparansi Statistik Kasus Publik
5. `05_halaman_login.png` — Antarmuka Login Petugas & Admin Berbasis RBAC
6. `06_kader_dashboard.png` — Dashboard Verifikator Lapangan Kader Lingkungan
7. `07_kader_verifikasi_laporan.png` — Antrean Laporan Warga & Klaim Yurisdiksi
8. `08_puskesmas_dashboard.png` — Dashboard Surveilans Puskesmas Beji
9. `09_puskesmas_validasi_pe_spk.png` — Antrean Validasi Hasil PE & Rekomendasi SPK
10. `10_puskesmas_fogging_requests.png` — Persetujuan Tindakan & Penjadwalan Fogging
11. `11_puskesmas_ews_klb.png` — Monitoring Status Kelurahan Siaga (EWS KLB)
12. `12_puskesmas_inventory_approval.png` — Persetujuan Permintaan Logistik Puskesmas
13. `13_tim_fogging_dashboard.png` — Dashboard Lapangan Tim Eksekutor Fogging
14. `14_tim_fogging_jadwal_pelaksanaan.png` — Penugasan Harian Tim Fogging (PWA View)
15. `15_dinkes_dashboard.png` — Dashboard Eksekutif Dinas Kesehatan Kota Depok
16. `16_dinkes_peta_gis.png` — GIS Heatmap Sebaran Kasus 11 Kecamatan Kota Depok
17. `17_dinkes_monitoring_klb.png` — Pemantauan Ambang Batas & Deklarasi KLB Kota
18. `18_dinkes_pemantauan_sla.png` — Audit Kepatuhan SLA 3x24 Jam Kader se-Kota Depok
19. `19_dinkes_gudang_inventaris.png` — Manajemen Gudang Utama & Transfer Logistik Dinkes
20. `20_developer_dashboard.png` — Dashboard Pemeliharaan Super Admin Developer
21. `21_developer_system_status.png` — Monitoring Kesehatan Layanan & Latensi API
22. `22_developer_audit_logs.png` — Jejak Audit Keamanan & Riwayat Perubahan Data
23. `23_developer_user_management.png` — Manajemen Akun & Otoritas 7 Peran Pengguna
24. `24_live_report_form_filled.png` — Bukti Pengisian Data Formulir Warga Secara Live
25. `25_live_report_submission_success.png` — Dialog Berhasil & Terbitnya Nomor Resi Resmi
26. `26_live_track_laporan_timeline.png` — Pelacakan Live Status Kasus Baru pada Timeline
27. `27_kader_laporan_baru_masuk.png` — Verifikasi Laporan Baru Masuk ke Inbox Kader
28. `28_puskesmas_spk_modal_detail.png` — Modal Keputusan SPK Rule-Based & Advisory Override
29. `29_tim_fogging_penugasan_lapangan.png` — Detail Penugasan Lapangan Tim Fogging
30. `30_tim_fogging_modal_update_status.png` — Modal Pemotongan Stok Bahan & Unggah Bukti Pelaksanaan

---

## 9. Struktur Direktori Proyek

```
DBD_Guard/
├── backend/                             # Layanan API & Logika Bisnis (Golang)
│   ├── cmd/
│   │   ├── api/main.go                  # Entry point HTTP server Fiber
│   │   ├── e2e/main.go                  # Skrip simulasi pengujian E2E backend
│   │   └── seeder/main.go               # Seeder data master & skenario demo
│   ├── internal/
│   │   ├── config/                      # Manajemen environment, DB, migrasi
│   │   ├── delivery/http/               # Handler HTTP & routing Fiber
│   │   ├── domain/                      # Entity model, DTO, interface kontrak
│   │   ├── infrastructure/              # Layanan notifikasi FCM & MinIO storage
│   │   ├── repository/                  # Akses basis data PostgreSQL + PostGIS
│   │   └── usecase/                     # Logika bisnis, SLA cron, SPK, EWS
│   └── go.mod
│
├── frontend-web/                        # Antarmuka Pengguna Web & PWA (Next.js 16)
│   ├── src/
│   │   ├── app/                         # App router (Rute publik & tiap peran)
│   │   │   ├── developer/               # Rute Super Admin Developer
│   │   │   ├── dinkes/                  # Rute Dinas Kesehatan Kota Depok
│   │   │   ├── fogging/                 # Rute Tim Fogging Lapangan
│   │   │   ├── kader/                   # Rute Kader Lingkungan RW
│   │   │   ├── puskesmas/               # Rute Petugas & Kepala Puskesmas
│   │   │   ├── publik/                  # Portal transparansi data agregat
│   │   │   ├── report/                  # Formulir pelaporan kasus warga
│   │   │   ├── track/                   # Pelacakan status laporan real-time
│   │   │   └── login/                   # Halaman autentikasi login terpadu
│   │   ├── components/                  # Komponen UI modular (Peta, Navbar, Modal)
│   │   └── lib/                         # Konfigurasi API axios, cookie, Leaflet
│   └── package.json
│
├── prd/                                 # Dokumen Spesifikasi Kebutuhan Produk (PRD)
│   └── PRD_DBD_Guard_Komprehensif.md
│
├── docs/                                # Panduan Teknis & Dokumen Laporan PKL
│   ├── DEMO.md                          # Panduan alur presentasi demo PKL
│   ├── KETERBATASAN.md                  # Laporan batasan teknis & transparansi
│   └── Risk_Management_DBD_Guard.md     # Analisis manajemen risiko sistem
│
├── screenshots/                         # Bukti visual hasil pengujian otomatis
│   └── e2e_evidence/
│       ├── desktop/                     # 31 Screenshot resolusi Desktop (1280x800)
│       └── mobile/                      # 31 Screenshot resolusi Mobile (390x844)
│
├── scripts/                             # Skrip automasi & seeding lokal
└── README.md                            # Dokumentasi induk proyek DBD Guard
```

---

## 10. Catatan Kepatuhan Regulasi & Perlindungan Data Pribadi

1. **Prinsip Kerahasiaan Data (UU PDP)**: Seluruh data publik yang ditampilkan pada `/publik` dan peta sebaran hanya berupa agregat kelurahan/kecamatan, tanpa membeberkan koordinat presisi rumah tangga atau identitas pasien.
2. **Opsi Anonimitas Pelapor**: Warga yang memilih opsi laporan anonim tidak diwajibkan menginput kontak/email dan dapat melacak penanganan murni melalui nomor referensi acak sistem.
3. **Pemberitahuan Data Simulasi**: Sesuai dengan tujuan prototipe evaluasi PKL, seluruh data yang termuat dalam seeder berlabel **DATA SIMULASI** dan bersifat fiktif untuk melindungi privasi warga Kota Depok.
