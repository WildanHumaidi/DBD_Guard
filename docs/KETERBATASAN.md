# Keterbatasan Sistem & Catatan Transparansi Pengembangan
### Sistem Monitoring dan SPK Penanganan DBD Kota Depok (DBD Guard)
**Penyusun**: Mahasiswa PKL Dinas Komunikasi dan Informatika (Diskominfo) Kota Depok

---

## Ringkasan Eksekutif

Dokumen ini disusun sebagai bentuk **kejujuran dan transparansi teknis** dalam laporan pelaksanaan Praktek Kerja Lapangan (PKL) di Diskominfo Kota Depok. Meskipun sistem **DBD Guard** telah berhasil mengimplementasikan seluruh kebutuhan fungsional inti (arsitektur Go Fiber, Next.js PWA, PostGIS, SPK berbasis aturan, inventaris terpadu, dan RBAC 7 peran), terdapat sejumlah batasan teknis, operasional, dan manajerial yang harus diketahui secara terbuka oleh penguji maupun mentor PKL.

---

## Daftar Keterbatasan Sistem Saat Ini

### 1. Mekanisme Backup Database Masih Berada di Server yang Sama
* **Kondisi Saat Ini**:
  Sistem telah dilengkapi dengan skrip pencadangan basis data otomatis (`scripts/backup_restore.ps1`) yang menghasilkan file dump PostgreSQL (`.dump`) di direktori lokal `backend/storage/backups/`.
* **Keterbatasan**:
  Penyimpanan cadangan ini masih berada pada mesin/disk fisik yang sama dengan server aplikasi dan database utama (*local storage*). Sistem belum mengintegrasikan mekanisme *offsite backup* otomatis (seperti replikasi berkala ke server backup terpisah, NAS Diskominfo, atau *cloud storage* independen).
* **Risiko**:
  Apabila terjadi kerusakan fatal pada perangkat keras (*hard disk failure*) atau bencana fisik di lokasi server, salinan cadangan berisiko hilang bersamaan dengan data utama.
* **Rekomendasi Lanjutan**:
  Mengonfigurasi cronjob atau skrip transfer berkala (misal menggunakan rsync / SCP / S3 API) yang memindahkan file dump ke *cold storage* terisolasi milik Pemerintah Kota Depok setiap pergantian hari.

---

### 2. Pengujian Ulang Pemulihan Bencana (*Restore Drill*) Belum Dilakukan Berkala
* **Kondisi Saat Ini**:
  Prosedur pemulihan (*restore procedure*) telah terdokumentasi dan terbukti dapat dijalankan saat pengujian awal.
* **Keterbatasan**:
  Pelaksanaan simulasi pemulihan bencana (*restore drill*) belum dijadwalkan secara otomatis dan belum diuji coba secara berkala pada lingkungan tiruan (*sandbox environment*).
* **Risiko**:
  Tanpa drill berkala, ada kemungkinan file dump yang dihasilkan di masa depan mengalami korupsi data (*silent corruption*) atau memiliki perbedaan skema migrasi yang tidak disadari sebelum insiden nyata terjadi.
* **Rekomendasi Lanjutan**:
  Membuat pipeline otomasi pengujian bulanan yang secara otomatis memulihkan file dump terbaru ke instance PostgreSQL sementara dan menjalankan *health-check query* untuk memvalidasi integritas data.

---

### 3. Pilot Operasional Belum Memperoleh Persetujuan Resmi dari Dinas Kesehatan
* **Kondisi Saat Ini**:
  Aplikasi ini dirancang dan dikembangkan sebagai proyek PKL di bawah binaan **Dinas Komunikasi dan Informatika (Diskominfo) Kota Depok** untuk menjawab permasalahan koordinasi manual penanganan DBD.
* **Keterbatasan**:
  Hingga saat ini, sistem berada pada status **Prototipe Siap Uji Coba (*Pilot Ready*)**, namun **belum memperoleh persetujuan regulasi formal** atau penetapan Surat Keputusan (SK) operasional dari **Dinas Kesehatan (Dinkes) Kota Depok**.
* **Dampak**:
  Aplikasi belum dapat dideploy secara resmi ke seluruh 38 Puskesmas dan ribuan kader lingkungan sebelum adanya kesepakatan koordinasi lintas Organisasi Perangkat Daerah (OPD) antara Diskominfo dan Dinas Kesehatan.
* **Rekomendasi Lanjutan**:
  Menjadwalkan audiensi resmi dan demonstrasi sistem kepada Kepala Bidang Pencegahan dan Pengendalian Penyakit (P2P) Dinas Kesehatan Kota Depok guna mendapatkan validasi kebijakan, pembobotan SOP klinis, serta penetapan wilayah kelurahan percontohan (*pilot area*).

---

### 4. Ketergantungan pada Server Tile Publik OpenStreetMap (OSM)
* **Kondisi Saat Ini**:
  Peta pemantauan spasial menggunakan pustaka *open-source* Leaflet.js dengan sumber peta (*tile layer*) default dari server publik OpenStreetMap (`https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png`).
* **Keterbatasan**:
  Server tile publik OSM beroperasi dengan prinsip *fair use* sukarela, tanpa jaminan *SLA uptime*, dan melarang penggunaan lalu lintas skala besar (*bulk heavy traffic*). Apabila server publik OSM mengalami lonjakan akses atau internet lokal terputus, gambar latar belakang peta jalanan tidak akan termuat (meskipun titik koordinat kasus dari PostGIS tetap ditampilkan).
* **Mitigasi yang Telah Dibuat**:
  Kode sumber sistem telah disesuaikan agar URL tile peta bersifat fleksibel dan dapat dikonfigurasi melalui *environment variable* (`NEXT_PUBLIC_MAP_TILE_URL`).
* **Rekomendasi Lanjutan**:
  Diskominfo Kota Depok disarankan mendeploy *tile cache* lokal mandiri (misal menggunakan TileServer GL, OpenMapTiles, atau GeoServer) di dalam jaringan intra-pemerintah (*Government Intra-Network* / Intranet) agar peta dapat diakses dengan latensi sangat rendah dan 100% independen tanpa tergantung internet eksternal.

---

### 5. Ketiadaan Model Prediksi Cerdas (AI / Machine Learning)
* **Kondisi Saat Ini**:
  Sistem Pendukung Keputusan (SPK) yang dibangun beroperasi secara murni berbasis aturan deterministik (*Rule-Based Decision Support System*).
* **Keterbatasan**:
  Sistem saat ini **tidak memiliki** model *Machine Learning* untuk meramalkan persebaran wabah di masa depan (*predictive outbreak modeling*) maupun analisis spasial cuaca/curah hujan ekstrem dari BMKG.
* **Alasan Teknis**:
  Keputusan ini diambil secara sadar agar rekomendasi tindakan klinis (*Fogging Fokus* vs *Edukasi PSN*) bersifat transparan, dapat diverifikasi langsung (*explainable*), tidak memicu halusinasi, dan patuh pada pedoman teknis surveilans Kementerian Kesehatan RI.

---

## Kesimpulan

Seluruh keterbatasan di atas tidak menghambat berjalannya fitur utama aplikasi dalam simulasi maupun pengujian saat ini. Namun demikian, poin-poin tersebut menjadi rekomendasi strategis bagi pihak Diskominfo Kota Depok apabila hendak melanjutkan proyek ini menuju tahap rilis produksi dan integrasi resmi tingkat kota.
