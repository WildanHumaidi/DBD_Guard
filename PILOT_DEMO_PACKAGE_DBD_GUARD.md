# PILOT PREPARATION & DEMONSTRATION PACKAGE
**Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi (Pemerintah Kota Depok)**

- **Tanggal Dokumen:** 01 Oktober 2026
- **Status Sistem:** **CONDITIONALLY PILOT READY (LIMITED PILOT ONLY) & DEMONSTRATION READY**
- **Prinsip Operasional:** Zero Code Modification (Berbasis build yang telah terverifikasi)

---

## 1. EXECUTIVE SUMMARY & GOVERNANCE STATE

Sistem DBD Guard Kota Depok telah mencapai tahap kematangan teknis **Technical PASS** pasca-eksekusi *P0 Fundamental Fix* dan *IMP-001 (Kader Jurisdiction Guard)*. Seluruh celah kritis kebocoran stok obat, privasi berkas medis, kelemahan segregasi wewenang (SoD), dan bypass klaim lintas yurisdiksi telah ditutup secara sempurna dengan 100% keberhasilan *test suite*.

**Status Peluncuran:**
* **Kesiapan Demonstrasi (Executive Demo):** **100% READY.** Seluruh alur fungsional dan lapisan keamanan dapat didemonstrasikan di hadapan pimpinan.
* **Kesiapan Uji Coba Lapangan (Limited Pilot):** **CONDITIONALLY READY.** Siap dieksekusi di wilayah **Puskesmas Beji** dan **Puskesmas Pancoran Mas**, di bawah payung Surat Perintah Tugas (SPT) Kepala Dinas Kesehatan.
* **Keputusan Tata Kelola Tertunda (DEC-01, DEC-02, DEC-03):** Dinilai sebagai **Non-Blocking** untuk pilot terbatas, dengan asumsi operasional berjalan pada volume stok aman dan status non-KLB. 

Dokumen ini disusun untuk memandu jalannya uji coba terkendali melalui 4 fase operasional: **Prepare $\rightarrow$ Demonstrate $\rightarrow$ Operate $\rightarrow$ Monitor**.

---

## 2. PREPARE: PRECONDITIONS & ENVIRONMENT SETUP

Persiapan operasional lapangan wajib memastikan 100% prasyarat terpenuhi sebelum sistem dibuka untuk relawan kader dan publik.

### A. Prasyarat Infrastruktur & Data
1. **Lingkungan Server & Database:** Docker Container untuk backend, frontend, dan PostgreSQL aktif 24/7.
2. **Master Data Wilayah:** 63 Kelurahan di 11 Kecamatan Kota Depok telah dipetakan terhadap 38 Puskesmas.
3. **Disaster Recovery:** Skrip cadangan `backup_restore.ps1` aktif menjadwalkan pencadangan otomatis (telah lulus *live drill* pemulihan 100%).
4. **Data Awal Inventaris:** Gudang logistik Puskesmas Beji dan Pancoran Mas tercatat aman ($\ge 20$ Liter) untuk menghindari kebutuhan draf pengajuan obat darurat selama pilot.

### B. Prasyarat Pengguna & Regulasi
1. **Akun Pengguna 7 Peran:** Terdistribusi sesuai berkas `DAFTAR_AKUN_DBD_GUARD.md` (Warga, Kader, Puskesmas [Petugas PE, Kepala, Fogging], Dinkes [Surveilans, Pengelola Barang]).
2. **Legalisasi Pelaksanaan:** Penandatanganan **Nota Dinas / Surat Perintah Tugas (SPT)** dari Kepala Dinas Kesehatan khusus untuk dua puskesmas percontohan (pengganti sementara Perwal/SK).
3. **Komunikasi Darurat:** Pembuatan grup WhatsApp koordinasi siaga (*Level 1 Faskes, Level 2 Dinkes, Level 3 Dev*).

---

## 3. DEMONSTRATE: STAKEHOLDER DEMONSTRATION SCRIPT

Panduan *Executive Demo* ini memvalidasi integritas operasional dan ketangguhan keamanan sistem di hadapan pimpinan daerah.

### Skenario A: End-to-End Golden Flow (Siklus Wabah Normal)
**Tujuan:** Memperlihatkan sinergi pelaporan publik hingga eksekusi fogging.
1. **Warga Lapor (`/report`):** Aktor Warga mengisi formulir publik (input titik GPS Beji, unggah foto jentik). Tiket terbit: `MENUNGGU_PE_KADER`.
2. **Kader Wilayah Klaim (`/kader`):** Kader RW Beji login. Sistem memvalidasi wilayah kerja dan mengizinkan klaim.
3. **Kader Verifikasi Awal:** Kader konfirmasi kasus valid. Tiket mutasi ke `MENUNGGU_PE_PUSKESMAS`.
4. **Penyelidikan Epidemiologi Faskes (`/puskesmas`):** Petugas Sanitarian Beji melakukan input PE (Jentik +, $\ge 3$ Penderita Demam). 
5. **DSS Rekomendasi:** Sistem Pakar memproses kalkulasi $\rightarrow$ Hasil: **FOGGING FOKUS**.
6. **Approval Kepala Faskes:** Kepala Puskesmas login, verifikasi SoD (*Four-Eyes Principle*), dan menerbitkan Jadwal Fogging.
7. **Eksekusi Fogging (`/fogging`):** Tim Lapangan login, tiba di lokasi, **wajib mengunggah foto bukti** penyemprotan, lalu klik Selesai.
8. **Dinkes Monitoring (`/dinkes`):** Dasbor Heatmap real-time menyala merah; stok insektisida di Puskesmas Beji otomatis terpotong sesuai akuntansi sistem.

### Skenario B: Security Flow & Jurisdiction Guard (IMP-001)
**Tujuan:** Mendemonstrasikan lapisan pertahanan sistem dari akses tidak sah.
1. **Uji Klaim Lintas Yurisdiksi:** Kader dari wilayah Pancoran Mas mencoba mengklaim laporan kasus yang berada di Kelurahan Beji (via UI maupun *bypass API/Curl*).
2. **Hasil Demo:** Sistem seketika me-reject dengan status **`403 FORBIDDEN`** ("Akses ditolak: laporan di luar yurisdiksi wilayah kerja").

### Skenario C: Exception Flow (Banding & Segregation of Duties)
**Tujuan:** Memperagakan mekanisme sengketa (*Appeal*) dan proteksi anti-rangkap wewenang.
1. Laporan warga ditolak kader sebagai *HOAX*. Warga mengajukan **BANDING** dengan melampirkan foto diagnosis Dokter RS.
2. Laporan naik eskalasi, di-takeover oleh Puskesmas. Petugas Sanitarian turun tangan melakukan verifikasi ulang.
3. **SoD Guard:** Sanitarian yang sama mencoba bertindak sebagai Kepala Puskesmas untuk menyetujui hasil kerjanya sendiri.
4. **Hasil Demo:** Sistem memblokir aksi dengan **`403 FORBIDDEN`** (*OfficerID == ValidatorID*). Kepala Puskesmas wajib memvalidasi.

---

## 4. OPERATE: TRAINING & OPERATION MANUAL

Panduan operasional di masa uji coba (Pilot Terbatas):

### Ketentuan Operasional Reguler
1. **Manajemen Beban Kerja:** Setiap regu fogging lapangan secara *hardcoded* dibatasi maksimal **1 tugas per hari** demi keselamatan K3.
2. **Notifikasi:** Tidak menggunakan SMS Gateway berbayar. Seluruh notifikasi dikirimkan melalui layanan **Email (SMTP)** dan **Firebase Push Notification (FCM)** yang 100% gratis.
3. **Peta Publik:** Sistem memuat *map tiles* dari OpenStreetMap (OSM). Petugas disarankan menggunakan aplikasi secara wajar untuk menghindari blokir IP (*rate-limiting*) dari pihak OSM.

### Panduan Mengelola Isu Kebijakan yang Tertunda (Quarantine Rules)
Mengingat DEC-01, DEC-02, dan DEC-03 masih **PENDING**, aturan berikut wajib dipatuhi selama Pilot:
- **Terkait Stok Obat:** Jika muncul notifikasi `FOGGING_TERTUNDA` akibat stok kosong, faskes harus berkoordinasi manual dengan Dinkes untuk pengambilan obat darurat (karena alur auto-draft langsung meluncur ke Dinkes).
- **Terkait KLB:** Jika terjadi lebih dari 1 lokasi penyemprotan darurat per hari di wilayah Beji/Pancoran Mas, koordinasi perbantuan regu **BKO Dinkes** dilakukan secara luring (manual via WhatsApp).
- **Terkait Banding Warga:** Pelapor yang laporannya ditolak dipersilakan melakukan banding maksimal 1 kali; faskes tetap melayani tanpa mempermasalahkan jeda waktu keterlambatan.

---

## 5. MONITOR: MONITORING & SUPPORT PROTOCOL

Pengawasan masa pilot diukur menggunakan indikator KRI (*Key Risk Indicators*) dan penanganan eskalasi.

### Indikator Keberhasilan Uji Coba (KRI Monitoring)
Pemantauan dilakukan secara mingguan oleh Kepala Seksi Surveilans Dinkes:
* **KRI-01 (SLA Laporan):** Ambang peringatan > 5% laporan dibiarkan melebihi 72 jam.
* **KRI-04 (Deteksi Spam):** Ambang peringatan > 5% laporan harian diblokir oleh geofencing spam beruntun.
* **KRI-06 (Durasi Approval Faskes):** Kepala Puskesmas wajib menyetujui jadwal fogging dalam rentang waktu < 12 jam.
* **KRI-07 (Penolakan Klaim Wilayah):** Log error HTTP 403 dipantau untuk memastikan IMP-001 sukses mencegah klaim liar.

### Jalur Eskalasi Insiden Darurat Lapangan
1. **Level 1 (Operasional - Ringan):** Koordinator Lapangan Puskesmas menangani kebingungan input warga/kader atau reset kata sandi.
2. **Level 2 (Kebijakan/Stok - Menengah):** Kepala Bidang P2P Dinas Kesehatan menyetujui *StockTransfer* manual, BAST, atau diskresi pergeseran jadwal antar faskes jika terjadi kluster darurat.
3. **Level 3 (Teknis Sistem - Kritis):** Tim Pengembang Web / DevOps turun tangan jika terjadi galat *downtime*, blokir server peta (OSM), limit kuota email tercapai, atau kebutuhan restorasi *Disaster Recovery*.

### Business Continuity Plan (BCP) - Panduan Pemadaman Server
Apabila terjadi gangguan server lebih dari $2 \times 24$ jam di puncak musim hujan:
1. Pihak Puskesmas wajib bertransisi menggunakan formulir investigasi PE cetak rangkap 3 sesuai standar Kemenkes.
2. Seluruh stok logistik yang digunakan dicatat manual dalam buku besar BMD.
3. Rekonsiliasi data digital susulan wajib dilakukan dalam tempo $1 \times 24$ jam pasca-server kembali aktif.
