# Pilot Preparation & Demonstration Package — DBD Guard
**Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi (Pemerintah Kota Depok)**

- **Tanggal Dokumen:** 01 Oktober 2026
- **Versi Dokumen:** 2.0 (Post-Claim Verification & Comprehensive Baseline)
- **Status Proyek Terkini:** **CONDITIONALLY PILOT READY — LIMITED PILOT ONLY** & **DEMO READY — PENDING DRY-RUN**
- **Ruang Lingkup Pilot Terbatas:** **Puskesmas Beji** (Kelurahan Beji) & **Puskesmas Pancoran Mas** (Kelurahan Pancoran Mas)
- **Prinsip Rekayasa Perangkat Lunak:** **ZERO CODE MODIFICATION** *(Tidak ada perubahan kode, penambahan fitur, atau modifikasi aturan bisnis)*
- **Kepatuhan Bukti & Penelusuran (Evidence & Traceability):** Seluruh status, klaim, data akun, alur, dan ambang batas merujuk langsung ke dokumen resmi proyek (*Source of Truth*).

---

## 1. Executive Summary

Dokumen **Pilot Preparation & Demonstration Package** ini disusun sebagai panduan operasional terpadu yang memadukan kebutuhan:
1. **Executive Demonstration & Stakeholder Presentation:** Menyajikan pembuktian keandalan sistem (*golden flow* dan *negative security flow*) di hadapan Walikota, Kepala Dinas Kesehatan, Kepala Diskominfo, dan Kepala Puskesmas se-Kota Depok.
2. **Limited Pilot Preparation & Operational Briefing:** Menyiapkan seluruh instrumen pelaksanaan uji coba terbatas (*field pilot*) bagi tim operasional di Puskesmas Beji dan Puskesmas Pancoran Mas.
3. **Pilot Monitoring, SOP, & Incident Handling:** Mengawal jalannya masa percontohan selama 30 hari kalender dengan parameter pemantauan risiko terukur (*Key Risk Indicators*), SOP lapangan, dan protokol penanganan gangguan (*Business Continuity Plan*).

### Status Kesiapan Resmi:
Sesuai ketetapan dokumen [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Section 16.1 & 16.2]:
* **DEMO READY — PENDING DRY-RUN:** Seluruh alur kerja standar maupun alur pengecualian dan proteksi keamanan telah terverifikasi secara fungsional backend pada build aktif tanpa error kompilasi [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 3.1 & Sec. 13], namun memerlukan pelaksanaan simulasi uji coba kering (*dry-run*) antarmuka web dan mobile sebelum sesi presentasi eksekutif.
* **CONDITIONALLY PILOT READY — LIMITED PILOT ONLY:** Siap digelar secara terbatas di 2 Puskesmas percontohan (Puskesmas Beji & Pancoran Mas) dengan syarat pemenuhan prasyarat operasional (*preconditions*), penyediaan jalur komunikasi siaga, dan penerbitan Surat Perintah Tugas (SPT) dari Kepala Dinas Kesehatan Kota Depok [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2.1].
* **NOT CITYWIDE PRODUCTION (GATED):** Peluncuran skala penuh se-Kota Depok (38 Puskesmas, 63 Kelurahan) tetap terkunci (*gated*) hingga 3 agenda keputusan tata kelola birokrasi (**DEC-01, DEC-02, DEC-03**) disahkan secara resmi oleh pimpinan daerah [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2.4 & BUSINESS_DECISION_GOVERNANCE_DBD_GUARD.md].

### 1.1 Laporan Perubahan Dokumen (Koreksi Pass 2)

Berdasarkan pelaksanaan audit sitasi mendalam (*Citation Audit & Ground Truth Verification*), seluruh ketidaksesuaian rujukan, angka target evaluasi, dan klaim kesiapan telah diperbaiki langsung pada dokumen ini:

| Lokasi | Sebelum (*Klaim Lama*) | Sesudah (*Koreksi Pass 2*) | Alasan Koreksi | Sumber Acuan (*Ground Truth*) |
|:---|:---:|:---:|:---|:---|
| **Header & Bab 1** | `DEMONSTRATION READY` | `DEMO READY — PENDING DRY-RUN` | Belum ada log pengukuran dry-run UI resmi. | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 13] |
| **Bab 1.1, 21, 22** | Daftar 6 syarat beragam versi | 7 Syarat resmi terpadu identik | Menjadikan Final Readiness sebagai acuan tunggal; dry-run dipisahkan ke Prasyarat Demo. | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 12 & 16.2] |
| **Checklist #4 (Akun)** | `[X] GO (VERIFIED)` | `[ ] PARTIAL (DIRENCANAKAN)` | Kredensial demo siap di file akun, namun akun personil riil baru dibuat pasca-SPT. | [`DAFTAR_AKUN_DBD_GUARD.md` → Sec. 2-3] |
| **Checklist #6 (SOP)** | `[X] GO (VERIFIED)` | `[ ] PARTIAL (DIRENCANAKAN)` | Draf teks ada di Bab 12, namun pengesahan formal & cetak buku saku belum dieksekusi. | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 (line 160)] |
| **Checklist #7 (WhatsApp)** | `[X] GO (VERIFIED)` | `[ ] PARTIAL (DIRENCANAKAN)` | Konsep eskalasi ada di Bab 2.3, pembuatan grup WA riil dijadwalkan pasca-SPT. | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 8.2 & Sec. 12 #8] |
| **Checklist #9 (Form BCP)**| `[X] GO (VERIFIED)` | `[ ] PARTIAL (DIRENCANAKAN)` | Standar form rangkap 3 disetujui, namun pencetakan fisik form belum ada di faskes. | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 & Sec. 12 #10] |
| **Checklist #12 (Stakeholder)**| `[X] GO (VERIFIED)` | `[ ] PARTIAL / PENDING` | Pemahaman formal pimpinan baru terjadi saat paparan eksekutif. | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2] |
| **Success Criteria #2 (SLA)** | Angka kaku `>= 85%` | `TBD — Business Owner (Usulan tim teknis: >= 85%)` | Target persentase kepatuhan tidak ada di PRD; PRD hanya menetapkan durasi batas 72 jam. Atribusi SK Kadinkes dihapus. | [`PRD_DBD_Guard_Komprehensif.md` → Bab 5.1 #2] |
| **Success Criteria #6 (Uptime)**| Angka kaku `>= 99.0%` | `TBD — Business Owner (Usulan tim teknis: >= 99.0%)` | Tidak ada dokumen SLA tertulis dari Diskominfo di repositori. Atribusi Diskominfo dihapus. | `TBD — Business Owner` |
| **Success Criteria #7 (Adopsi)**| Angka kaku `>= 80%` | `TBD — Business Owner (Usulan tim teknis: >= 80%)` | Tidak ada angka target numerik adopsi di PRD. Atribusi Bidang P2P dihapus. | [`PRD_DBD_Guard_Komprehensif.md` → Bab 3] |
| **Success Criteria #8 (IKM)** | Angka kaku `>= 80` | `TBD — Business Owner (Usulan tim teknis: survei)` | Standar survei IKM belum tersedia di repositori. Atribusi Setda dihapus. | `TBD — Business Owner` |
| **Jadwal Pra-Pilot** | Istilah "H-2 / H-3" | `TBD — ditetapkan setelah SPT terbit` | Istilah H-2/H-3 tidak ditemukan di repositori; jadwal riil bergantung pada tanggal terbit SPT. | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2.1] |
| **Bab 4.4 (Reset Demo)** | Perintah `-Action restore -DbName dbd_guard_demo` | `.\scripts\backup_restore.ps1 -SourceDB dbd_guard -RestoreDB dbd_guard_restore_test` | Menyesuaikan dengan parameter riil skrip `backup_restore.ps1`. Target DB adalah `dbd_guard_restore_test`. | [`scripts/backup_restore.ps1` lines 4-11] |
| **Bab 4.4 (Nama DB Pilot)** | `dbd_guard_pilot` | `[nama DB pilot — TBD]` | Nama `dbd_guard_pilot` tidak ditemukan pada file konfigurasi repositori (`.env`, `.env.example`). | [`backend/.env.example` line 9 (`DB_NAME=dbd_guard`)] |
| **Bab 4.4 (Durasi Restore)**| Klaim `<10 detik` | `durasi belum diukur` | Tidak ada log pengukuran durasi restorasi di repositori; dicatat sebagai action item uji. | Action Item |
| **Bab 5.1 (Run Sheet)** | 40 Menit (13 baris) | 35 Menit (Sesi Utama 25m + Q&A 10m) | Total run sheet dibatasi $\le 35$ menit; Audit Trail & Pen-Test dipindahkan ke cadangan Q&A. | Standar alokasi waktu presentasi |
| **Bab 21 (Keputusan)** | `GO FOR LIMITED PILOT` | `CONDITIONAL GO` | Wajib selalu menyertakan kata CONDITIONAL dan 7 syarat terpadu. | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.1] |

### 1.2 Hal yang Tidak Dapat Diperbaiki Lewat Dokumen (Memerlukan Tindakan Fisik / Keputusan Stakeholder)

Penyelesaian status kesiapan berikut berada di luar wewenang penyuntingan dokumen dan **wajib dieksekusi secara fisik atau diputuskan oleh pemangku kepentingan**:
1. **Penandatanganan Fisik SPT Pilot:** Penerbitan dan penandatanganan Surat Perintah Tugas resmi oleh Kepala Dinas Kesehatan Kota Depok.
2. **Pengesahan Tiga Agenda Kebijakan Birokrasi:** Keputusan tertulis pimpinan daerah atas lembar ballot **DEC-01** (alur draf stok), **DEC-02** (kapasitas fogging KLB), dan **DEC-03** (jendela kedaluwarsa banding).
3. **Pembuatan Berkas Fisik Dump Seed Demo:** Eksekusi pembuatan file `backups/dbd_guard_demo_seed.sql.gz` dari data sintetis.
4. **Perekaman Video Cadangan Demo:** Pembuatan direktori `docs/demo_artifacts/` dan penyimpanan rekaman video *walkthrough* resolusi tinggi.
5. **Uji Pengukuran Durasi Pemulihan Basis Data:** Eksekusi skrip `backup_restore.ps1` pada lingkungan demo untuk mencatat durasi riil restorasi.
6. **Pelaksanaan Dry-Run Terpandu 2x:** Simulasi alur presentasi menggunakan stopwatch dengan pencatatan log hasil tiap langkah.
7. **Pencetakan Fisik Dokumen Lapangan:** Pencetakan formulir penyelidikan epidemiologi Kemenkes rangkap 3 dan buku saku SOP operasional untuk sanitarian faskes.
8. **Pembentukan Grup Komunikasi Riil:** Pembuatan grup WhatsApp Siaga Pilot dengan memasukkan nomor telepon aktif seluruh pejabat dan staf yang ditugaskan.

---

## 2. Pilot Scope

Ruang lingkup uji coba terbatas dikunci secara ketat untuk mengisolasi risiko operasional dan memaksimalkan fokus pendampingan teknis:

### 2.1 Wilayah Kerja Percontohan
* **Puskesmas Percontohan 1: UPTD Puskesmas Beji**
  - *Kecamatan:* Beji
  - *Kelurahan Binaan Pilot:* Kelurahan Beji (Kode Wilayah Faskes terhubung di `master_kelurahans`).
  - *Karakteristik Wilayah:* Kepadatan penduduk tinggi, wilayah pemukiman heterogen, riwayat endemis sedang.
* **Puskesmas Percontohan 2: UPTD Puskesmas Pancoran Mas**
  - *Kecamatan:* Pancoran Mas
  - *Kelurahan Binaan Pilot:* Kelurahan Pancoran Mas.
  - *Karakteristik Wilayah:* Pemukiman padat perkotaan, bantaran sungai, riwayat kasus aktif DBD awal musim hujan.

### 2.2 Peran Pengguna yang Terlibat (*Participating Roles*)
Sesuai pemetaan peran PRD Bab 3 dan dokumen [`DAFTAR_AKUN_DBD_GUARD.md` → Section 2]:
1. **Warga Masyarakat (Pelapor):** Warga domisili Kelurahan Beji dan Kelurahan Pancoran Mas (akses publik bebas).
2. **Kader Jumantik / RW:** 5-10 Kader terpilih per kelurahan pilot (Peran: `ROLE_KADER`).
3. **Petugas Puskesmas (Surveilans/PE):** 2 Petugas Sanitarian per Puskesmas pilot (Peran: `ROLE_PETUGAS_PUSKESMAS`).
4. **Kepala Puskesmas:** Pimpinan faskes Beji & Pancoran Mas sebagai Kuasa Pengguna Anggaran (Peran: `ROLE_KEPALA_PUSKESMAS`).
5. **Tim Fogging Lapangan:** 1 Regu Fogging beranggotakan 2 operator per faskes (Peran: `ROLE_TIM_FOGGING`).
6. **Dinas Kesehatan (Dinkes Kota Depok):** Tim Surveilans P2P & Pengelola Gudang Farmasi Balai Kota (Peran: `ROLE_DINKES`).
7. **Administrator Sistem:** Tim Pengembang & PIC Diskominfo (Peran: `ROLE_ADMIN_DEVELOPER`).

### 2.3 Struktur Tim Pendukung & Jalur Eskalasi (*Support & Escalation Path*)
* **Jalur Koordinasi Cepat (Level 1 - Frontline Support):**
  - Media: WhatsApp Group Siaga Pilot DBD Guard (PIC Faskes Beji, PIC Faskes Panmas, Helpdesk Teknis).
  - Waktu Respon: $<15$ menit pada jam kerja (08.00 - 16.00 WIB).
* **Jalur Penanganan Medis / Kebijakan (Level 2 - Faskes & Dinkes Escalation):**
  - Penanggung Jawab: Kepala Seksi Surveilans & Imunisasi Dinkes Depok.
  - Kewenangan: Penyelesaian sengketa wilayah, konfirmasi hasil uji lab faskes, diskresi transfer stok antar-gudang.
* **Jalur Pemulihan Sistem (Level 3 - Engineering & Infrastructure Support):**
  - Penanggung Jawab: Lead Backend Developer & Admin Infrastruktur Diskominfo.
  - Kewenangan: Pemulihan database, reset sesi pengguna, penanganan kegagalan koneksi container Docker.

---

## 3. Pilot Preconditions & Traceability Matrix

### 3.1 Pilot Precondition Checklist

| No | Item Prasyarat Operasional | Persyaratan Kesiapan (*Requirement*) | Penanggung Jawab (*Owner*) | Status Kesiapan | Bukti Faktual & Rujukan (*Evidence*) |
|:---:|:---|:---|:---:|:---:|:---|
| 1 | **Pilot Accounts** | 7 peran terdaftar dengan kredensial & yurisdiksi Beji & Panmas. | Tim Pengembang & Faskes | `[ ] PARTIAL (DIRENCANAKAN)` | Kredensial akun uji terdefinisi di [`DAFTAR_AKUN_DBD_GUARD.md` → Sec. 2-3]; pembuatan akun personil riil Puskesmas & Kader menunggu penetapan jadwal (TBD — ditetapkan setelah SPT terbit, Owner: Dinkes & Faskes). |
| 2 | **Location Mapping** | Relasi Kecamatan $\rightarrow$ Kelurahan $\rightarrow$ Puskesmas terverifikasi di DB. | Tim Basis Data | `[X] GO (VERIFIED)` | Query tabel `master_kelurahans` & [`IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md` → Sec. 4 ("Core Logic & Repository")] |
| 3 | **Database Engine** | PostgreSQL 15 + PostGIS aktif, skema 15 tabel operasional utuh. | Diskominfo / Dev | `[X] GO (VERIFIED)` | Docker container port `5433`, migration clean [`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md` → Sec. 4 ("Zero-Regression Full Verification Suite")] |
| 4 | **Backup Automation** | Skrip `backup_restore.ps1` siap terjadwal harian pukul 23:59 WIB. | Tim Infrastruktur | `[ ] PARTIAL (DIRENCANAKAN)` | Skrip fisik ada di [`scripts/backup_restore.ps1`] dan live drill sukses [`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md` → Sec. 4.3]; otomasi cron server produksi menunggu jadwal rilis (TBD — ditetapkan setelah SPT terbit, Owner: Tim Infrastruktur Diskominfo). |
| 5 | **Notification Layer** | Fallback asynchronous logging aktif; template email & FCM siap. | Tim Backend | `[ ] PARTIAL (DIRENCANAKAN)` | Async logger aktif di usecase; integrasi SMTP/FCM resmi menunggu alokasi server pemda [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 5/RSK-13 & Sec. 10/Residual-03] |
| 6 | **SOP Operasional** | Panduan operasional sebelum, saat, dan pasca kerja didistribusikan. | Dinkes & Faskes | `[ ] PARTIAL (DIRENCANAKAN)` | Draf dokumen SOP disusun di Bab 12 dokumen ini & [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 (line 160)]; pencetakan fisik & sosialisasi nakes dijadwalkan (TBD — ditetapkan setelah SPT terbit, Owner: Dinkes & Faskes). |
| 7 | **Support Channel** | Kanal komunikasi instan WhatsApp Siaga Pilot aktif. | Tim Pengembang | `[ ] PARTIAL (DIRENCANAKAN)` | Struktur eskalasi dispesifikasi di Bab 2.3 dokumen ini & [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 8.2 & Sec. 12 #8]; pembentukan grup WhatsApp Siaga Pilot dijadwalkan (TBD — ditetapkan setelah SPT terbit, Owner: Tim Pengembang & Faskes). |
| 8 | **Manual Fallback** | Blangko kertas Form PE Kemenkes & buku register BAST disiapkan. | Sanitarian Faskes| `[ ] PARTIAL (DIRENCANAKAN)` | Standar prosedur BCP terdefinisi di [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 & Sec. 12 #10] serta [`RISK_MANAGEMENT_DBD_GUARD.md` → Sec. 11.2]; pengadaan cetak fisik blangko dijadwalkan (TBD — ditetapkan setelah SPT terbit, Owner: Sanitarian Faskes). |
| 9 | **Training / Briefing**| Sesi pengarahan teknis kader & petugas faskes percontohan. | Dinkes & Dev | `[ ] PARTIAL (DIRENCANAKAN)` | Modul materi siap di Bab 13 & 18; pelaksanaan workshop tatap muka dijadwalkan sebelum pembukaan akses publik (TBD — ditetapkan setelah SPT terbit, Owner: Dinkes & Dev) [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2.3 & Sec. 12 #11]. |
| 10 | **Monitoring KRI** | Dasbor pemantauan SLA, log audit, dan indikator risiko siap dipantau. | Surveilans Dinkes| `[X] GO (VERIFIED)` | Modul dasbor `/dinkes` & audit logs terverifikasi via test suite [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 8.2 & Sec. 9 #6 ("Pilot Readiness Monitoring")]. |

### 3.2 Traceability Matrix Prasyarat Conditional Go

Pemetaan keterlacakan 7 syarat gerbang **CONDITIONAL GO** dari dokumen acuan tunggal [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md`] menuju item checklist operasional dan bukti verifikasi:

```text
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                TRACEABILITY MATRIX PRASYARAT CONDITIONAL GO                                      │
├───────────────────────────────────┬───────────────────────────┬───────────┬──────────────────────────────────────┤
│ Syarat Resmi Conditional Go       │ Item Checklist Terkait    │ Status    │ Bukti Sumber (Ground Truth)          │
├───────────────────────────────────┼───────────────────────────┼───────────┼──────────────────────────────────────┤
│ 1. Otorisasi Resmi SPT Kadinkes   │ Checklist #11 (Bab 21)    │ PENDING   │ FINAL_READINESS: Sec. 16.2.1 & 12 #7 │
│ 2. Sesi Briefing Teknis Lapangan  │ Checklist #12 / Precond #9│ PARTIAL   │ FINAL_READINESS: Sec. 16.2.3 & 12 #11│
│ 3. Akun Personil Riil Lapangan    │ Checklist #4 / Precond #1 │ PARTIAL   │ FINAL_READINESS: Sec. 6 & 12 #5      │
│ 4. Form Fisik Manual BCP Kemenkes │ Checklist #9 / Precond #8 │ PARTIAL   │ FINAL_READINESS: Sec. 6 & 12 #10     │
│ 5. Pengesahan & Cetak SOP Faskes  │ Checklist #6 / Precond #6 │ PARTIAL   │ FINAL_READINESS: Sec. 6 (line 160)   │
│ 6. Aktivasi WhatsApp Siaga Pilot  │ Checklist #7 / Precond #7 │ PARTIAL   │ FINAL_READINESS: Sec. 8.2 & 12 #8    │
│ 7. Karantina Kebijakan (DEC-01..3)│ Tata Kelola Baseline      │ VERIFIED  │ FINAL_READINESS: Sec. 11 & 16.2.2    │
└───────────────────────────────────┴───────────────────────────┴───────────┴──────────────────────────────────────┘
```

---

## 4. Demo Dataset Specification

Sesuai aturan keamanan data pribadi (**UU PDP No. 27/2022**), seluruh data demonstrasi menggunakan **spesifikasi data sintetis (dummy)** yang mencerminkan topografi wilayah Kota Depok:

### 4.1 Master Wilayah & Lokasi Kasus Sintetis
* **Kasus Demo 1 (Normal Positive Cluster - Beji):**
  - *Pelapor:* "Budi Santoso" (Warga Terdaftar) / No. Telp: `0812-9988-7711` (Dummy)
  - *Alamat:* Jl. H. Asmawi No. 14, RT 03 / RW 05, Kelurahan Beji, Kecamatan Beji, Kota Depok.
  - *Koordinat GPS:* Lat `-6.372514`, Long `106.821345` (Valid area Beji).
  - *Gejala Warga:* Demam tinggi mendadak 3 hari, bintik merah pada kulit.
  - *Foto Bukti Warga:* `dummy_rujukan_rs_beji.jpg` (Spesifikasi: Gambar sintetis surat pengantar berstempel fiktif "Klinik Sehat Sejahtera").
* **Kasus Demo 2 (Negative Hoax / Dispute Cluster - Pancoran Mas):**
  - *Pelapor:* "Anonim" (`anonymous_flag = true`) / No. Telp: Masked
  - *Alamat:* Jl. Mampang Indah RT 02 / RW 04, Kelurahan Pancoran Mas.
  - *Koordinat GPS:* Lat `-6.398112`, Long `106.804221` (Valid area Panmas).
  - *Deskripsi:* Laporan genangan air di lahan kosong tanpa bukti penderita demam.
  - *Foto Bukti Warga:* `dummy_genangan_saluran.jpg`.

### 4.2 Spesifikasi Master Logistik Awal Demo
* **Gudang Utama Dinkes (Balai Kota Depok):**
  - *Insektisida (Malathion 96 EC):* $250.0\text{ Liter}$
  - *Larvasida (Abate 1 SG):* $100.0\text{ Kilogram}$
* **Gudang Puskesmas Beji:**
  - *Insektisida (Malathion 96 EC):* $25.0\text{ Liter}$
  - *Solar / Pelarut:* $80.0\text{ Liter}$
* **Gudang Puskesmas Pancoran Mas:**
  - *Insektisida (Malathion 96 EC):* $15.0\text{ Liter}$
  - *Solar / Pelarut:* $50.0\text{ Liter}$

> *Catatan Pengembang:* Seluruh dataset di atas telah terkonfigurasi pada skrip migrasi awal (`db_seed.sql`). Apabila diperlukan reset database instan, gunakan prosedur pemulihan pada Bab 5.

---

### 4.3 Spesifikasi Lingkungan Runtime Demo
* **Aplikasi Web:** `http://localhost:3000` (Next.js 16 App Router)
* **Backend API Gateway:** `http://localhost:8080/api/v1` (Golang Fiber v2)
* **Basis Data:** `localhost:5433` (PostgreSQL 15 Container: `dbd_guard_db`)
* **Peramban yang Direkomendasikan:** Google Chrome / Microsoft Edge (Mode Incognito atau Multi-Profile untuk mendemokan multi-peran secara serentak).

### 4.4 Prosedur Reset Cepat Sebelum Demo (Quick Reset Routine)
Jalankan prosedur pemulihan berikut di direktori root `DBD_Guard` untuk mereset seluruh status database demonstrasi ke titik nol bersih (*clean state*):

> [!CAUTION]
> **ATURAN MUTLAK ISOLASI BASIS DATA (DATABASE ISOLATION):**
> Skrip pemulihan demonstrasi ini **HANYA BOLEH DIJALANKAN PADA DATABASE LOKAL DEMO/TEST (`dbd_guard_restore_test` di port `5433`)**. DILARANG KERAS mengeksekusi restore point demo terhadap database operasional percontohan (`[nama DB pilot — TBD]`) pada server staging/Diskominfo agar tidak menghapus atau menimpa data transaksi riil faskes!
> Parameter `-RestoreDB` pada `scripts/backup_restore.ps1` secara baku mengisolasi target restore ke `dbd_guard_restore_test`, sehingga database sumber operasional tidak akan tertimpa.

*Spesifikasi Teknis Skrip `scripts/backup_restore.ps1`:*
Sesuai deklarasi parameter pada baris 4-11 berkas [`scripts/backup_restore.ps1`]:
* `$ContainerName` (default: `"dbd_guard_db"`)
* `$PgUser` (default: `"postgres"`)
* `$SourceDB` (default: `"dbd_guard"`)
* `$RestoreDB` (default: `"dbd_guard_restore_test"`) — *Target database restore terisolasi bawaan skrip*
* `$BackupDir` (default: `"e:\Projects\DBD_Guard\backend\storage\backups"`)
* `$RetentionDays` (default: `7`)

*Status Keberadaan Berkas Fisik Reset Demo:*
* `scripts/backup_restore.ps1`: **`[X] VERIFIED TERSEDIA (6.960 bytes)`** di repositori lokal [`scripts/backup_restore.ps1`].
* `backups/dbd_guard_demo_seed.sql.gz`: **`[ ] BELUM TERSEDIA FISIK (FOLLOW-UP / ACTION ITEM)`** — Pengembang wajib mengeksekusi dump awal seed data demo ke file terkompresi sebelum sesi presentasi.
* `docs/demo_artifacts/`: **`[ ] BELUM TERSEDIA FISIK (FOLLOW-UP / ACTION ITEM)`** — Direktori penyimpanan video cadangan (*Full HD walkthrough*) dan tangkapan layar eksekutif wajib dibuat sebelum demonstrasi pimpinan.

```powershell
# Direktori Kerja: E:\Projects\DBD_Guard
# 1. Jalankan skrip restore point demonstrasi ke DB lokal terisolasi (BUKAN DB OPERASIONAL PILOT)
#    Sintaks parameter presisi sesuai baris 4-11 scripts/backup_restore.ps1:
powershell -ExecutionPolicy Bypass -File .\scripts\backup_restore.ps1 -SourceDB dbd_guard -RestoreDB dbd_guard_restore_test

# 2. Verifikasi status container database lokal
docker ps --filter "name=dbd_guard_db"

# 3. Jalankan pengujian verifikasi cepat backend
cd backend
go test -run "TestP0_01|TestIMP001_TEST01" ./internal/usecase/...
```
*Expected Output:* Restorasi selesai (durasi belum diukur — dicatat sebagai action item uji coba), konsistensi data 100%, pengujian verifikasi berstatus **PASS**.

### 4.5 Demo Failure & Contingency Plan (Rencana Kontinjensi Kegagalan Demo)
Jika terjadi kendala teknis saat presentasi di depan pimpinan:
1. **Kegagalan Jaringan Lokal / Port Konflik:** Presenter segera mengalihkan tampilan ke rekaman demonstrasi terotentikasi (*Full HD Video Walkthrough*) yang disiapkan pada folder `docs/demo_artifacts/` (*Item tindak lanjut sebelum hari demonstrasi*).
2. **Keterlambatan Render Tile Peta Leaflet:** Tetap lanjutkan presentasi dengan menjelaskan bahwa Leaflet.js telah memetakan koordinat lat/long secara presisi di basis data, dan caching offline proxy Diskominfo akan mengatasi hal ini saat rilis server Balai Kota.

### 4.6 Prasyarat Eksekusi Demonstrasi (Demo Preconditions Checklist)
Sebelum sesi demonstrasi eksekutif di hadapan Walikota dan Kepala Dinas Kesehatan diselenggarakan, daftar prasyarat teknis berikut wajib dipenuhi dan diuji seluruhnya:

| No | Item Prasyarat Demonstrasi | Status | Bukti / Catatan Tindak Lanjut (*Action Item*) |
|:---:|:---|:---:|:---|
| 1 | **Seed Dump Data Sintetis Siap** | `[ ] PENDING` | Berkas `backups/dbd_guard_demo_seed.sql.gz` dibuat dan diverifikasi bebas PII warga riil. |
| 2 | **Uji Restorasi & Pengukuran Durasi** | `[ ] PENDING` | Eksekusi restore ke DB `dbd_guard_restore_test` diuji dan durasinya dicatat dalam detik. |
| 3 | **Dry-Run Penuh 2x dengan Stopwatch** | `[ ] PENDING` | Simulasi end-to-end 2 kali dengan log lolos/gagal per langkah skenario ($\le 25$ menit). |
| 4 | **Video Walkthrough Cadangan Siaga** | `[ ] PENDING` | File video resolusi tinggi tersimpan pada folder `docs/demo_artifacts/`. |
| 5 | **Skrip Uji Penetrasi Negatif Siap** | `[ ] PENDING` | Token JWT pengujian penolakan HTTP 403 (IMP-001, SoD, PII) disiapkan dalam bentuk tombol/skrip cepat. |
| 6 | **Penetapan Moda per Skenario Demo** | `[ ] PENDING` | Keputusan definitif alokasi moda `[LIVE]`, `[REKAMAN CADANGAN]`, atau `[FALLBACK SLIDE]` disahkan presenter. |

---

## 5. Golden Flow Demo (Alur Pelaporan Standar Positif Kasus)

### 5.1 Demo Run Sheet (1-Halaman Panduan Eksekutif)
*Target Durasi Total:* **35 Menit** (25 Menit Demonstrasi Sistem Inti, 10 Menit Tanya Jawab Stakeholder).

| No | Durasi | Segmen Skenario Demo | Peran / Aktor | Layar / Media | Moda Presentasi | Bukti Uji & Kesiapan | Inti Pesan Eksekutif |
|:---:|:---:|:---|:---|:---|:---:|:---|:---|
| 1 | 00 - 02m | **Pembukaan & Urgensi** | Presenter | Slide 1 - 2 | `[LIVE]` | Dokumen Slide Bab 18 | Latar belakang DBD Depok & urgensi satu pintu |
| 2 | 02 - 05m | **Pelaporan Warga (Golden)** | Warga Publik | `/report` & `/track` | `[LIVE]` / `[REKAMAN CADANGAN]` | `TestCitizenReport_Submission` (PASS) | Kemudahan input laporan, bukti foto, resi unik |
| 3 | 05 - 08m | **Klaim & Verifikasi Kader** | Kader RW Beji | `/kader` | `[LIVE]` / `[REKAMAN CADANGAN]` | `TestIMP001_TEST01_SameRegion` (PASS) | Otorisasi yurisdiksi ketat & kepatuhan SLA 72 jam |
| 4 | 08 - 11m | **Penyelidikan PE & DSS** | Sanitarian | `/puskesmas` | `[LIVE]` / `[REKAMAN CADANGAN]` | `TestPESubmission_Success` (PASS) | Survei jentik, HI/ABJ, rekomendasi otomatis DSS |
| 5 | 11 - 13m | **Persetujuan Fogging (SoD)** | Kepala Faskes | `/puskesmas` | `[LIVE]` / `[REKAMAN CADANGAN]` | `TestP0_03_and_04_DSS_and_SoD` (PASS) | Pemisahan wewenang validator $\ne$ investigator |
| 6 | 13 - 16m | **Eksekusi Fogging & Bukti** | Tim Fogging | `/fogging` (Mobile) | `[LIVE]` / `[REKAMAN CADANGAN]` | State machine PASS; **PERLU DRY-RUN (Mobile UI)** | Bukti foto fisik wajib & pemotongan obat otomatis |
| 7 | 16 - 18m | **Monitoring Spasial (GIS)** | Dinkes Depok | `/dinkes` | `[LIVE]` / `[REKAMAN CADANGAN]` | Backend PASS; **PERLU DRY-RUN (Leaflet Tile)** | Peta persebaran kasus aktif & agregasi respon |
| 8 | 18 - 21m | **Exception Flow: Sengketa & Banding**| Warga & Faskes | `/track` & `/puskesmas` | `[LIVE]` / `[REKAMAN CADANGAN]` | `TestP0_06_AppealToPELifecycle` (PASS); **PERLU DRY-RUN (UI Modal & Takeover)** | Hak banding warga berbukti lab & faskes takeover |
| 9 | 21 - 23m | **Logistik: Transfer Stok BMD** | Dinkes & Faskes | `/puskesmas` | `[LIVE]` / `[REKAMAN CADANGAN]` | `TestP0_01_InventoryAccountingIntegrity` (PASS) | Konservasi stok 2-tahap (In-Transit) anti-kebocoran |
| 10 | 23 - 25m | **DSS: Clinical Override** | Kepala Puskesmas | `/puskesmas` | `[LIVE]` / `[REKAMAN CADANGAN]` | `TestP0_03_DSS_Override_WithReason` (PASS) | Kedaulatan dokter atasi kondisi kritis balita |
| 11 | 25 - 35m | **Penutup & Diskusi Pilot** | Presenter & Kadinkes | Slide 10 - 13 | `[LIVE]` | Dokumen Slide Bab 18 | Pembahasan SPT, kesiapan pilot Beji & Panmas |

### 5.1.1 Skenario Cadangan untuk Sesi Tanya Jawab (Deep-Dive Cadangan)
Untuk menjaga total durasi presentasi utama tetap berada pada batas $\le 35$ menit, dua skenario teknis mendalam berikut dialokasikan sebagai **materi demonstrasi cadangan** yang dapat ditampilkan sewaktu-waktu apabila diminta oleh penguji / stakeholder:
1. **Audit Trail Forensik & EWS Alert Spasial:**
   - *Estimasi Waktu:* 3 Menit
   - *Moda Presentasi:* **`[LIVE DB / LOG]`** / **`[FALLBACK SLIDE]`**
   - *Peran & Media:* Administrator IT / Dasbor Dinkes (`/dinkes`) dan kueri PostgreSQL langsung.
   - *Tujuan:* Membuktikan log mutasi *append-only* anti-manipulasi dan mekanisme *early warning* kluster KLB kepada Inspektorat / Diskominfo.
2. **Uji Penetrasi Keamanan & RBAC (Direct API Bypass):**
   - *Estimasi Waktu:* 3 Menit
   - *Moda Presentasi:* **`[LIVE CURL / POSTMAN]`** / **`[FALLBACK SLIDE]`**
   - *Peran & Media:* Penguji Keamanan / Terminal Bash / Postman.
   - *Tujuan:* Menunjukkan respon seketika `HTTP 403 Forbidden` saat penyerang mencoba menembak API bypass lintas yurisdiksi atau melanggar prinsip *Four-Eyes* kepada tim keamanan Diskominfo.

### 5.2 Skenario Pembuktian Alur Utama (*Happy Path*)
Skenario pembuktian alur utama dari warga hingga tindakan selesai:

| Langkah | Aktor & URL | Aksi yang Dilakukan | Hasil yang Diharapkan (Expected) | Bukti Faktual & Kesiapan (*Evidence*) |
|:---:|:---|:---|:---|:---:|
| **Step 1** | **Warga**<br>`/report` | Mengisi form laporan kasus Demam Berdarah di Jl. H. Asmawi RW 05 Beji, melampirkan foto rujukan klinik, submit. | Laporan tersimpan, tiket berstatus `MENUNGGU_PE_KADER`, resi unik terbit (format: `DBD-DPK-2026XXXX`). | `TestCitizenReport_Submission` (PASS) |
| **Step 2** | **Kader**<br>`/kader` | Login sebagai `kader@depok.go.id` (Wilayah Beji). Melihat antrean laporan baru, menekan tombol **"Klaim Tugas"**. | Klaim berhasil; nama investigator terisi; otorisasi yurisdiksi lolos karena wilayah laporan cocok. | `TestIMP001_TEST01_SameRegion_ClaimSuccess` (PASS) |
| **Step 3** | **Kader**<br>`/kader` | Kader mendatangi lokasi, memeriksa keberadaan jentik di penampungan air warga, menekan **"Verifikasi Valid"**. | Status tiket bermutasi menjadi `MENUNGGU_PE_PUSKESMAS`; notifikasi terkirim ke faskes pembina. | `TestReportVerification_Success` (PASS) |
| **Step 4** | **Petugas PE**<br>`/puskesmas` | Login sebagai `petugas.puskesmas@depok.go.id`. Menginput hasil survei PE: Jentik Positif ($HI = 15\%$), Penderita Demam = 3 Orang. | Data PE tersimpan; algoritma DSS mengevaluasi aturan bisnis PRD Bab 5.2: **Rekomendasi = FOGGING_FOKUS**. | `TestPESubmission_Success` & formula DSS (PASS) |
| **Step 5** | **Kepala Faskes**<br>`/puskesmas` | Login sebagai `kepala.puskesmas@depok.go.id`. Membuka tab approval PE, meninjau hasil, menekan **"Setujui & Jadwalkan Fogging"**. | Jadwal fogging terbentuk; status kasus bermutasi ke `DIJADWALKAN_FOGGING`. SoD guard memvalidasi validator $\ne$ investigator. | `TestP0_03_and_04_DSS_and_SoD` (PASS) |
| **Step 6** | **Tim Fogging**<br>`/fogging` | Login sebagai `tim.fogging@depok.go.id` (Tampilan Mobile). Mengubah status menjadi **OTW** $\rightarrow$ **Mulai Pelaksanaan**. | Status penugasan terupdate real-time; jam mulai tercatat otomatis di sistem. | `schedule_handler.go` state machine (PASS); **PERLU DRY-RUN (Mobile UI)** |
| **Step 7** | **Tim Fogging**<br>`/fogging` | Petugas selesai menyemprot, mengunggah foto dokumentasi pengasapan, mencatat obat terpakai (2 L Malathion), tekan **"Selesai"**. | Status berubah `COMPLETED`. Sistem memverifikasi syarat foto bukti fisik terpenuhi. | Invariant `fogging_evidences > 0` (PASS); **PERLU DRY-RUN (Frontend File Upload)** |
| **Step 8** | **Sistem Logistik**<br>`/puskesmas` | Sistem secara otomatis memotong stok gudang faskes secara atomik. | Saldo Malathion Puskesmas Beji berkurang dari $25.0\text{ L} \rightarrow 23.0\text{ L}$. Zero leak! | `TestP0_01_InventoryAccountingIntegrity` (PASS) |
| **Step 9** | **Dinkes**<br>`/dinkes` | Login sebagai `admin@depok.go.id`. Membuka peta GIS Kota Depok. | Muncul marker kasus hijau (Selesai Penanganan) di RW 05 Beji; statistik respon harian bertambah. | `TestCountCasesByArea` (PASS); **PERLU DRY-RUN (Leaflet Tile Rendering)** |

---

## 6. Exception Flow Demo (Sengketa, Banding, & Eskalasi)
*Moda Presentasi:* **`[LIVE]`** (Pilihan Utama) / **`[REKAMAN CADANGAN]`** (Kontinjensi).

Skenario pembuktian penanganan laporan yang ditolak kader namun disengketakan warga:

```text
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                EXCEPTION & ESCALATION WORKFLOW                                   │
│                                                                                                  │
│  [Warga Melapor] ──► [Kader Tolak: REJECTED_HOAX] ──► [Warga Ajukan Banding + Surat Lab RS]      │
│                                                                   │                              │
│  [Kepala Puskesmas Validasi Akhir] ◄── [Puskesmas Takeover & PE Ulang] ◄── [Eskalasi ke Faskes]  │
└──────────────────────────────────────────────────────────────────────────────────────────────────┘
```

| Langkah | Aktor & Layar | Tindakan yang Diperagakan | Hasil Sistem yang Terbukti (*Verified Result*) | Bukti Uji & Kesiapan (*Test Evidence*) |
|:---:|:---|:---|:---|:---:|
| **Ex-1** | **Kader** (`/kader`) | Kader menolak laporan warga di Panmas karena dianggap hoaks / genangan air biasa. | Status tiket bermutasi menjadi `REJECTED_HOAX`. Alur pelaporan tidak dimatikan sepihak. | `TestP0_06_AppealToPELifecycle` (PASS) |
| **Ex-2** | **Warga** (`/track`) | Warga memeriksa tracking resi, melihat status ditolak, menekan tombol **"Ajukan Banding"** dan mengunggah foto bukti hasil lab trombosit RS. | Kuota banding terpakai (`appeal_count = 1`); status tiket bermutasi menjadi `DIESKALASIKAN`. | `TestP0_06_AppealToPELifecycle` Step 1 (PASS); **PERLU DRY-RUN (Frontend Modal Banding)** |
| **Ex-3** | **Puskesmas** (`/puskesmas`)| Petugas Surveilans Puskesmas menerima tiket eskalasi banding, menekan tombol **"Takeover Investigasi"**. | Akun investigator dialihkan dari kader ke petugas faskes (`investigator_id` terupdate). | `TestP0_06_AppealToPELifecycle` Step 3 (PASS); **PERLU DRY-RUN (Frontend Takeover Button)** |
| **Ex-4** | **Petugas Faskes** | Petugas Puskesmas turun langsung ke lapangan, mendapati penderita demam nyata, menginput hasil PE revisi. | Data PE kedua tersimpan valid berpasangan dengan tiket banding warga. | `TestP0_06_AppealToPELifecycle` Step 4 (PASS) |
| **Ex-5** | **Puskesmas** (`/puskesmas`)| Kepala Puskesmas mereview sengketa banding dan mengesahkan status kasus menjadi `TERVERIFIKASI`. | Sengketa warga tuntas terselesaikan secara adil berbasis bukti laboratorium sahih. | `TestP0_06_AppealToPELifecycle` Step 6 (PASS) |

---

## 7. Security Demonstration (Pembuktian Ketahanan Sistem)
*Moda Presentasi:* **`[LIVE API / CURL / POSTMAN]`** (Pilihan Utama) / **`[FALLBACK SLIDE]`** (Kontinjensi).

Presenter dapat mendemonstrasikan secara langsung kegagalan skenario manipulasi berbahaya (*negative testing*) untuk meyakinkan tim audit keamanan Diskominfo:

### Skenario 1: Percobaan Klaim Lintas Wilayah oleh Kader (Bypass Direct API)
* **Aksi Penyerang:** Kader Kelurahan Pancoran Mas mencoba mengirimkan request HTTP POST via cURL/Postman untuk mengklaim laporan milik Kelurahan Beji:
  ```bash
  curl -X POST http://localhost:8080/api/v1/kader/reports/REPORT-BEJI-UUID/claim \
    -H "Authorization: Bearer <TOKEN-KADER-PANCORAN-MAS>"
  ```
* **Respon Sistem:** **`HTTP 403 FORBIDDEN`**
* **Payload Respon:**
  ```json
  {
    "error": "Akses ditolak: Laporan berada di luar yurisdiksi wilayah kerja Anda",
    "message": "akses ditolak: laporan di luar yurisdiksi wilayah kerja"
  }
  ```
* **Bukti Pengujian:** Lolos 100% pada `TestIMP001_HTTP_CrossRegion_Forbidden` dan `TestIMP001_HTTP_DirectAPIBypass_Forbidden` [`IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md` → Sec. 7].

### Skenario 2: Pelanggaran Pemisahan Wewenang (Self-Validation SoD Bypass)
* **Aksi Penyerang:** Petugas Puskesmas yang bertindak sebagai investigator lapangan mencoba menembak API validasi untuk menyetujui hasil PE buatannya sendiri tanpa lewat Kepala Puskesmas:
  ```bash
  curl -X POST http://localhost:8080/api/v1/puskesmas/pe/PE-UUID/validate \
    -H "Authorization: Bearer <TOKEN-PETUGAS-INVESTIGATOR-YANG-SAMA>"
  ```
* **Respon Sistem:** **`HTTP 403 FORBIDDEN`**
* **Penjelasan:** Sistem mendeteksi `peResult.OfficerID == validatorID` dan secara mutlak memblokir aksi untuk menegakkan *Four-Eyes Principle*.
* **Bukti Pengujian:** Lolos pada `TestP0_04_SoD_SameActor_Validation_MustFail` [`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md` → Sec. 3].

### Skenario 3: Penyerobotan Akses Berkas Foto Medis (Unauthorized Media Access)
* **Aksi Penyerang:** Pengguna membuka URL foto rekam medis pelapor secara langsung tanpa menyertakan token autentikasi atau menggunakan token dari wilayah faskes lain.
* **Respon Sistem:** **`HTTP 401 Unauthorized`** atau **`HTTP 403 Forbidden`** (Akses berkas statis ditutup total; media disajikan melalui *Private Authenticated Stream*).
* **Bukti Pengujian:** Lolos pada `TestP0_02_SensitiveFileStorageProtection` [`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md` → Sec. 3].

---

## 8. Inventory Accounting Demonstration
*Moda Presentasi:* **`[LIVE WEB]`** (Pilihan Utama) / **`[REKAMAN CADANGAN]`** (Kontinjensi).

Mendemonstrasikan kepatuhan mutlak sistem terhadap kaidah tata kelola Barang Milik Daerah (**Permendagri No. 19/2016**):

### 8.1 Rumus Konservasi Total Stok Logistik
Sistem menjamin invariant matematis berikut selalu terpenuhi pada setiap detik mutasi:
$$\mathbf{Stok\_Gudang\_Asal + Stok\_Gudang\_Tujuan + Stok\_Dalam\_Perjalanan\ (IN\_TRANSIT) = Total\_Aset\_Riil}$$

```text
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                             ALUR MUTASI DUA TAHAP (STOCK TRANSFER)                               │
│                                                                                                  │
│  [Gudang Dinkes: 250 L]                                              [Puskesmas Beji: 25 L]      │
│            │                                                                    │                │
│            ▼ (Dinkes Approve Permintaan 10 L)                                   │                │
│  [Gudang Dinkes: 240 L] ──► [IN_TRANSIT: 10 L (Kurir Fisik)]                    │                │
│                                           │                                     │                │
│                                           ▼ (Puskesmas Cek Fisik & Tekan Terima)▼                │
│  [Gudang Dinkes: 240 L] ──► [IN_TRANSIT: 0 L] ─────────────────────► [Puskesmas Beji: 35 L]      │
└──────────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Pembuktian Saat Demonstrasi
1. **Kondisi Awal (Before):**
   - Gudang Farmasi Dinkes: $250.0\text{ Liter}$
   - Gudang Puskesmas Beji: $25.0\text{ Liter}$
   - In-Transit: $0.0\text{ Liter}$
   - *Total Sistem:* $\mathbf{275.0\text{ Liter}}$
2. **Persetujuan Mutasi (Transfer In-Transit):**
   - Kepala Dinkes menyetujui pengiriman $10.0\text{ Liter}$ Malathion ke Beji.
   - Saldo Gudang Dinkes terpotong menjadi $240.0\text{ Liter}$.
   - Saldo Puskesmas Beji **belum bertambah** (tetap $25.0\text{ Liter}$).
   - Status tercatat `IN_TRANSIT` sebesar $10.0\text{ Liter}$.
   - *Total Sistem:* $240 + 25 + 10 = \mathbf{275.0\text{ Liter}}$ (Konsisten 100%).
3. **Penerimaan Fisik (Receive & Reconciliation):**
   - Pengelola Barang Puskesmas Beji memeriksa fisik jeriken, mencocokkan nomor seri BAST, lalu menekan tombol **"Konfirmasi Penerimaan Barang"**.
   - Status transfer berubah menjadi `COMPLETED`.
   - Saldo Puskesmas Beji resmi bertambah menjadi $35.0\text{ Liter}$.
   - *Total Sistem:* $240 + 35 + 0 = \mathbf{275.0\text{ Liter}}$ (Akuntansi Sempurna).
*Bukti Pengujian Backend:* `TestP0_01_InventoryAccountingIntegrity` skenario B, D, dan E (Reject & Refund) [`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md` → Sec. 3].

---

## 9. DSS Demonstration (Clinical Override & Human Decision)
*Moda Presentasi:* **`[LIVE WEB]`** (Pilihan Utama) / **`[REKAMAN CADANGAN]`** (Kontinjensi).

Mendemonstrasikan posisi Decision Support System (DSS) sebagai **instrumen pembantu (*Advisory System*)**, bukan pemutus mandiri yang memberangus kedaulatan dokter:

### 9.1 Skenario Rekomendasi vs Diskresi Medis
1. **Input PE Lapangan:**
   - Hasil jentik negatif ($HI = 2\%$, di bawah batas $5\%$).
   - Jumlah penderita demam sekitar $= 1$ orang.
2. **Kalkulasi Algoritma DSS:**
   - Rekomendasi otomatis: **EDUKASI_DAN_PSN** (Tidak disarankan fogging karena tidak memenuhi kriteria epidemiologis).
3. **Keputusan Manusia (Human Decision):**
   - Dokter / Kepala Puskesmas mengetahui bahwa pasien yang terjangkit demam tersebut merupakan anak balita dengan komplikasi kritis (DHF Grade III / Syok) di wilayah padat pemukiman.
   - Kepala Puskesmas menggunakan tombol **"Override Rekomendasi DSS"**.
   - Sistem mewajibkan pengisian kolom teks: *"Alasan Medis Override"*.
   - Input Alasan: *"Pasien balita mengalami DSS/DHF Grade III syok; terdapat riwayat kematian kasus DBD di RT bersebelahan 2 minggu lalu. Fogging fokus darurat diperlukan untuk memutus rantai vektor lokal."*
4. **Pencatatan Audit Trail:**
   Sistem menyimpan audit trail lengkap: `ACTION = OVERRIDE_DSS`, `ACTOR = dr. Kepala Puskesmas Beji`, `REASON = [Teks di atas]`, `TIMESTAMP = [Waktu Riil]`.
*Bukti Pengujian Backend:* `TestP0_03_DSS_Override_WithoutReason_MustFail` dan `TestP0_03_DSS_Override_WithReason_Success` [`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md` → Sec. 3].

---

## 10. EWS / KLB Demonstration
*Moda Presentasi:* **`[LIVE WEB]`** (Pilihan Utama) / **`[REKAMAN CADANGAN]`** (Kontinjensi).

Mendemonstrasikan algoritma surveilans Early Warning System (EWS) dalam mendeteksi ancaman Kejadian Luar Biasa (KLB):

### 10.1 Logika Perhitungan Paralel Wilayah
Sistem DBD Guard tidak menggunakan kueri tunggal yang sempit, melainkan mengevaluasi kasus secara paralel pada dua tingkat agregasi wilayah:
1. **Tingkat Kelurahan:** Memantau akumulasi kasus terkonfirmasi dalam 14 hari terakhir. Ambang batas siaga: $\ge 3$ kasus per RW atau $\ge 5$ kasus per Kelurahan.
2. **Tingkat Kecamatan:** Memantau tren kenaikan kasus mingguan dibandingkan rata-rata historis (kenaikan $>100\%$).

### 10.2 Skenario Demo EWS
* Saat kasus di Kelurahan Beji mencapai ambang batas, dasbor Dinas Kesehatan (`/dinkes`) secara otomatis menampilkan:
  - Indikator Visual: **Banner Merah Berkedip (*High Alert KLB*)** — **PERLU DRY-RUN (Komponen UI Alert)**.
  - Peta Spasial: Pewarnaan poligon wilayah Kelurahan Beji menjadi **Merah Pekat (High Risk Zone)** — **PERLU DRY-RUN (Leaflet Poligon Layer)**.
  - Rekomendasi Dinkes: Mobilisasi tim supervisi epidemiologi dan penyediaan cadangan insektisida Balai Kota.
* *Bukti Pengujian Backend:* `TestCountCasesByArea_KelurahanAndKecamatanScope` (PASS) [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 5/RSK-06].

---

## 11. Audit Trail Demonstration
*Moda Presentasi:* **`[LIVE DB / LOG]`** (Pilihan Utama) / **`[FALLBACK SLIDE]`** (Kontinjensi).

Menunjukkan akuntabilitas yuridis sistem di mana setiap manipulasi data dapat dipertanggungjawabkan di hadapan aparat pengawas internal pemerintah (Inspektorat / BPK):

```text
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                   CONTOH LOG AUDIT TRAIL SISTEM                                  │
├────────────────────┬──────────────────┬─────────────────┬───────────────────┬────────────────────┤
│ Timestamp          │ User / Role      │ Action          │ Target Record     │ Detail Perubahan   │
├────────────────────┼──────────────────┼─────────────────┼───────────────────┼────────────────────┤
│ 2026-10-01 08:30   │ kader.beji01     │ CLAIM_TASK      │ RPT-202610-0012   │ Status: PENDING->  │
│                    │ (Kader Beji)     │                 │                   │ IN_INVESTIGATION   │
│ 2026-10-01 09:15   │ kader.beji01     │ SUBMIT_PE       │ PE-202610-0008    │ Larva Index: 28%   │
│                    │                  │                 │                   │ Breteau Index: 35  │
│ 2026-10-01 10:00   │ kapus.beji       │ APPROVE_SCHEDULE│ SCHED-202610-0004 │ Action: FOGGING    │
│                    │ (Kepala Faskes)  │                 │                   │ Date: 2026-10-02   │
│ 2026-10-01 10:30   │ logistik.beji    │ DISPENSE_STOCK  │ TX-STK-202610-003 │ Malathion -5L      │
│                    │                  │                 │                   │ Sisa: 25L          │
│ 2026-10-01 11:00   │ dinkes.admin     │ OVERRIDE_DSS    │ RPT-202610-0015   │ Reason: Kluster    │
│                    │                  │                 │                   │ Sekolah Aktif      │
└────────────────────┴──────────────────┴─────────────────┴───────────────────┴────────────────────┘
```

*Bukti Pengujian & Akuntabilitas:*
- Log audit dicatat secara *append-only* (tidak dapat diubah/dihapus oleh pengguna biasa) via trigger dan event model database [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 7 & PRD Bab 8.2].
- Mendukung verifikasi forensik kepatuhan SLA 3x24 jam dan transparansi alokasi insektisida.

---

## 12. Pilot Standard Operating Procedure (SOP)
*Status Kesiapan:* **`PARTIAL (DIRENCANAKAN)`** — *Draft SOP selesai disusun dalam dokumen ini; pencetakan fisik buku saku dan sosialisasi kepada tenaga kesehatan Puskesmas dijadwalkan setelah SPT terbit (TBD — Owner: Dinkes P2P & Tim Teknis) [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 & Sec. 12 #9].*

Panduan alur operasional harian bagi petugas faskes, kader jumantik, dan tim logistik selama masa pilot:

### 12.1 Pra-Operasional (Pukul 07.30 - 08.00 WIB)
1. **Pengecekan Kesiapan Sistem:** Administrator Diskominfo dan Koordinator Puskesmas memastikan peladen backend dan antarmuka web aktif (*HTTP 200 OK* pada endpoint `/health`).
2. **Sinkronisasi Tugas Kader:** Kader Jumantik login ke aplikasi, memeriksa daftar penugasan PE baru di wilayah kelurahan masing-masing, dan memastikan konektivitas perangkat seluler.
3. **Verifikasi Stok Awal:** Petugas farmasi/logistik Puskesmas memeriksa fisik insektisida dan memastikan data stok di aplikasi sesuai dengan stok riil di gudang faskes.

### 12.2 Selama Jam Operasi (Pukul 08.00 - 16.00 WIB)
1. **Triase Laporan Warga:** Staf Puskesmas memvalidasi laporan warga yang masuk dari kelurahan wilayah kerjanya.
2. **Pelaksanaan PE Lapangan:** Sanitarian faskes / Kader memeriksa jentik dalam radius 100 meter dari rumah penderita DBD, mendokumentasikan foto, dan menginput form PE digital.
3. **Persetujuan Jadwal:** Kepala Puskesmas mereview hasil PE dan memberikan keputusan tindakan (Fogging Fokus / PSN).
4. **Eksekusi Pengabutan:** Tim fogging menyemprot sesuai radius blok penanganan, mengunggah foto dokumentasi asli, dan mencatat pemakaian bahan bakar/obat.

### 12.3 Prosedur Penanganan Gangguan (Jika Terjadi Error / Kendala)
1. **Tangkapan Layar Error:** Pengguna mengambil tangkapan layar (*screenshot*) pesan error dan mencatat nomor resi laporan yang bermasalah.
2. **Lapor Kanal Bantuan:** Kirimkan bukti error ke WhatsApp Siaga Pilot.
3. **Dilarang Melakukan Bypass Manual:** Pengguna dilarang memanipulasi database secara langsung atau memaksakan klaim di luar prosedur.
4. **Aktivasi Cadangan Manual (BCP):** Jika sistem padam $>2$ jam, catat hasil survei PE menggunakan formulir kertas standar Kemenkes.

### 12.4 Pasca Jam Operasi (Pukul 16.00 - 17.00 WIB)
1. **Rekonsiliasi Harian:** Pengelola logistik mencocokkan jumlah pemakaian obat digital dengan fisik jeriken di gudang ($1 \times 24$ jam).
2. **Review Kepatuhan SLA:** Surveilans memeriksa laporan yang mendekati batas kedaluwarsa 72 jam untuk disiapkan penugasan eskalasi.
3. **Pencadangan Otomatis:** Sistem mengeksekusi automasi backup harian pada pukul 23:59 WIB.

---

## 13. Business Continuity Plan (BCP) & Contingency
*Status Kesiapan:* **`PARTIAL (DIRENCANAKAN)`** — *Protokol kontinjensi terdefinisi; pencetakan fisik formulir darurat Kemenkes dan distribusi ke 2 faskes dijadwalkan setelah SPT terbit (TBD — Owner: Sanitarian Puskesmas Beji & Panmas) [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 8/BCP-01 & Sec. 11, `FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 & Sec. 12 #10].*

Matriks tindakan pemulihan operasional jika terjadi kegagalan infrastruktur selama masa uji coba pilot:

| Kondisi Kegagalan | Tindakan Operasional Darurat (*Contingency Procedure*) | Penanggung Jawab (*Owner*) |
|:---|:---|:---:|
| **Peladen Backend Down / Crash** | 1. Segera alihkan pencatatan ke Formulir Kertas Kemenkes Rangkap 3.<br>2. Tim Teknis melakukan restart container via Docker Compose.<br>3. Input susulan (*back-entry*) dilakukan dalam $1 \times 24$ jam pasca-pulih [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6]. | Tim Teknis & Petugas Faskes |
| **Basis Data Corrupt / Rusak** | 1. Eksekusi skrip pemulihan: `powershell -ExecutionPolicy Bypass -File .\scripts\backup_restore.ps1 -SourceDB dbd_guard -RestoreDB dbd_guard_restore_test`.<br>2. Pulihkan data ke titik cadangan terakhir (RPO $< 24$ jam, durasi pemulihan belum diukur).<br>3. Rekonsiliasi log transaksi hari berjalan.<br>*(Catatan Keamanan: Pengujian restorasi hanya ditujukan ke database pengujian `dbd_guard_restore_test`, DILARANG KERAS menyentuh database operasional pilot `[nama DB pilot — TBD]`).* | Administrator Diskominfo |
| **Layanan Email SMTP / FCM Padam** | 1. Sistem beralih ke mode *asynchronous log capture* (alur transaksi tetap berjalan normal).<br>2. Notifikasi darurat tugas kader dikirimkan manual via WhatsApp Siaga Pilot. | Tim Pengembang & Koordinator Kader |
| **Koneksi Internet Seluler Putus di Lokasi PE** | 1. Petugas PE mencatat hasil survei pada draf offline formulir kertas.<br>2. Input data ke aplikasi dilakukan setibanya petugas di Puskesmas. | Petugas Lapangan PE |
| **Kader Berhalangan Hadir / Sakit** | 1. Jika tugas terbengkalai $> 72$ jam, sistem mengeksekusi *auto-release cron*.<br>2. Kepala Puskesmas melakukan tindakan *takeover* langsung ke petugas faskes. | Kepala Puskesmas |

---

## 14. Incident Handling & Severity Classification

| Tingkat Keparahan | Kriteria Insiden | Contoh Kasus | Waktu Respon Target (SLA) | Prosedur Penanganan & Eskalasi |
|:---:|:---|:---|:---|:---:|
| **Severity 1 (Kritis)** | Seluruh layanan sistem padam total; data tidak dapat diakses; menghentikan operasi pilot kedua faskes. | Basis data PostgreSQL lock total; backend panic crash; kebocoran PII terbuka. | Respon: $<15$ Menit<br>Resolusi: $<2$ Jam | Lapor langsung ke Lead Dev & Kadinkes. Aktifkan BCP manual penuh. Eksekusi live restore drill jika DB corrupt. |
| **Severity 2 (Mayor)** | Alur kerja penting terganggu; transaksi kunci terhenti namun sistem masih dapat diakses sebagian. | Tombol klaim kader error 500; transfer stok gagal memotong saldo; foto bukti tidak terunggah. | Respon: $<30$ Menit<br>Resolusi: $<4$ Jam | Eskalasi ke Tim Backend. Periksa log Docker. Gunakan bypass administratif sementara jika disetujui Kadinkes. |
| **Severity 3 (Minor)** | Gangguan kosmetik antarmuka atau fitur pendukung; tidak memblokir penanganan kasus DBD. | Salah ketik label UI; marker peta agak lambat merender; filter tanggal laporan meleset 1 hari. | Respon: $<2$ Jam<br>Resolusi: $<24$ Jam | Catat pada backlog perbaikan mingguan. Patch dirilis pada siklus maintenance berkala. |

---

## 15. Pilot Success Criteria

Kriteria keberhasilan terukur (*Measurable Success Criteria*) untuk mengevaluasi kelayakan DBD Guard selama 30 hari masa uji coba pilot:

| No | Indikator Keberhasilan | Target Evaluasi Pilot | Dasar Rujukan & Status Angka | Metode Pengukuran |
|:---:|:---|:---|:---|:---|
| 1 | **Pelaporan Masyarakat** | 100% Laporan tersimpan | Integritas transaksi database tanpa data loss [`PRD_DBD_Guard_Komprehensif.md` → Bab 8.2 & `FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 5.1] | Rekonsiliasi nomor tiket resi terbit |
| 2 | **Kepatuhan Waktu Respon (SLA)** | **TBD — Business Owner**<br>*(Usulan tim teknis — bukan ketetapan: $\ge 85\%$ Selesai $< 72$ Jam)* | Batas waktu SLA 72 Jam mengacu pada PRD Bab 4.1; target persentase kelulusan belum ditetapkan oleh Business Owner [`PRD_DBD_Guard_Komprehensif.md` → Bab 4.1] | Kueri durasi tiket s/d laporan PE faskes |
| 3 | **Efektivitas Otorisasi Yurisdiksi** | 0 Kasus klaim liar (100% blokir 403) | Teruji via automated test suite otorisasi yurisdiksi [`IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md` → Sec. 7] | Audit error log HTTP usecase claim |
| 4 | **Akurasi Akuntansi Logistik** | Selisih Fisik vs DB = 0 (*Zero Discrepancy*) | Stok aman $\ge 20$ L [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 11.1], rekonsiliasi harian $1 \times 24$ jam [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6], dan model atomik transfer 2-tahap [`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md` → Sec. 3/P0-01] | Berita Acara Rekonsiliasi Fisik Gudang |
| 5 | **Kepatuhan Dokumentasi Fogging** | 100% Memiliki Foto Riil | Invariant penutupan tiket penugasan fogging [`PRD_DBD_Guard_Komprehensif.md` → Bab 3 & Bab 5.1, `FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 5.5] | Kueri relasi tabel `fogging_evidences` |
| 6 | **Kestabilan Sistem (Uptime)** | **TBD — Business Owner**<br>*(Usulan tim teknis — bukan ketetapan: $\ge 99.0\%$ Jam Kerja)* | Target ketersediaan layanan belum ditetapkan oleh Business Owner [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 7] | Log monitoring server & health check probe |
| 7 | **Adopsi Pengguna Lapangan** | **TBD — Business Owner**<br>*(Usulan tim teknis — bukan ketetapan: $\ge 80\%$ Kader Aktif Login)* | Target adopsi kader lapangan belum ditetapkan oleh Business Owner [`PRD_DBD_Guard_Komprehensif.md` → Bab 3] | Kueri field `last_login` tabel `users` |
| 8 | **Kepuasan Layanan Publik** | **TBD — Business Owner**<br>*(Usulan tim teknis — bukan ketetapan: target Indeks Kepuasan Masyarakat)* | Standar IKM belum ditetapkan oleh Business Owner [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2] | Kuesioner/survei kepuasan pelapor warga |

---

## 16. Pilot Monitoring & Key Risk Indicators (KRI)

Mengacu pada Bab 12 dokumen [`RISK_MANAGEMENT_DBD_GUARD_v2.md`], parameter monitoring berikut wajib dipantau mingguan oleh Tim Surveilans Dinkes:

| Kode KRI | Nama Indikator Pemantauan | Formula Perhitungan | Ambang Peringatan (*Warning*) | Ambang Bahaya (*Critical*) | Frekuensi Pantau |
|:---:|:---|:---|:---:|:---:|:---:|
| **KRI-01** | Pelanggaran Batas SLA Laporan | $rac{	ext{Laporan }>72	ext{ Jam}}{	ext{Total Laporan Masuk}} 	imes 100\%$ | $> 5\%$ / Minggu | $> 15\%$ / Minggu | Mingguan |
| **KRI-02** | Frekuensi Transfer Darurat Stok | Jumlah transaksi `StockTransfer` antar-Puskesmas | $> 3$ kali / Bulan | $> 8$ kali / Bulan | Bulanan |
| **KRI-03** | Rasio Override Rekomendasi DSS | $rac{	ext{Kasus Di-Override}}{	ext{Total Hasil PE}} 	imes 100\%$ | $> 15\%$ / Bulan | $> 35\%$ / Bulan | Bulanan |
| **KRI-04** | Rasio Deteksi Laporan Spam | $rac{	ext{Laporan Ditolak Spam}}{	ext{Total Laporan Masuk}} 	imes 100\%$ | $> 5\%$ / Minggu | $> 12\%$ / Minggu | Mingguan |
| **KRI-05** | Tingkat Inaktivitas Relawan Kader | $rac{	ext{Kader Tidak Aktif }>14	ext{ Hari}}{	ext{Total Kader Terdaftar}} 	imes 100\%$ | $> 20\%$ Kader | $> 40\%$ Kader | Bulanan |
| **KRI-06** | Waktu Tunggu Approval Jadwal | Rata-rata durasi antrean persetujuan Kepala Faskes | $> 12$ Jam | $> 24$ Jam | Harian |
| **KRI-07** | Percobaan Pelanggaran Yurisdiksi | Jumlah penolakan HTTP 403 pada klaim tugas | $> 1\%$ Percobaan | $> 5\%$ Percobaan | Mingguan |

---

## 17. Daily Checklist (Operational Checklist)

Checklist harian yang wajib diisi oleh Koordinator Sistem di Puskesmas percontohan setiap sore hari:

| Komponen Pengawasan | Status Harian | Bukti Pemeriksaan (*Evidence*) | Petugas Pemeriksa (*Owner*) |
|:---|:---:|:---|:---:|
| **1. Backend Service** | [ ] NORMAL / [ ] ERROR | HTTP 200 pada endpoint `/health` | Helpdesk Teknis |
| **2. Basis Data PostgreSQL**| [ ] NORMAL / [ ] ERROR | Koneksi pool stabil; disk space $< 80\%$ | Helpdesk Teknis |
| **3. Automasi Pencadangan** | [ ] NORMAL / [ ] ERROR | File `.sql.gz` terbentuk di folder backup | Tim Infrastruktur |
| **4. Layanan Notifikasi** | [ ] NORMAL / [ ] ERROR | Log pengiriman email/FCM tercatat | Tim Pengembang |
| **5. Akun Pengguna Faskes** | [ ] NORMAL / [ ] ERROR | Petugas & kader dapat login lancar | Sanitarian Faskes |
| **6. Antrean Laporan Warga**| [ ] NORMAL / [ ] ERROR | Tidak ada tiket menggantung tanpa status | Surveilans Puskesmas |
| **7. Progress Penyelidikan PE**| [ ] NORMAL / [ ] ERROR | Seluruh PE tervalidasi dokter faskes | Dokter / Kepala Faskes |
| **8. Realisasi Jadwal Fogging**| [ ] NORMAL / [ ] ERROR | Foto bukti terunggah; obat terpotong | Koordinator Fogging |
| **9. Sinkronisasi Inventaris**| [ ] NORMAL / [ ] ERROR | Saldo digital = sisa jeriken fisik | Pengelola Barang |
| **10. Alert Sistem EWS** | [ ] NORMAL / [ ] ALERT | Cek indikator status KLB kelurahan | Surveilans Dinkes |
| **11. Integritas Audit Logs** | [ ] NORMAL / [ ] ERROR | Seluruh mutasi hari ini terekam | Tim Keamanan IT |

---

## 18. Stakeholder Presentation Flow (13 Slide Panduan Paparan)

Struktur paparan eksekutif untuk meyakinkan Walikota Depok dan Kepala Dinas Kesehatan:

* **Slide 1 — Judul & Masalah:**
  "Tantangan Penanganan DBD Kota Depok: Keterlambatan respon laporan warga, risiko kebocoran logistik insektisida, dan keterbatasan koordinasi lintas faskes."
* **Slide 2 — Solusi Terpadu DBD Guard:**
  "Sistem informasi satu pintu berbasis geospasial dan alur kerja terstandar yang menghubungkan Warga $\rightarrow$ Kader RW $\rightarrow$ Puskesmas $\rightarrow$ Dinas Kesehatan."
* **Slide 3 — Alur Kerja 7 Peran (Business Workflow):**
  Diagram visual transisi status tiket: Pelaporan publik, verifikasi SLA 72 jam, PE lapangan, rekomendasi DSS, jadwal fogging, dan pemotongan stok otomatis.
* **Slide 4 — Arsitektur & Prinsip 100% Open-Source:**
  Teknologi modern berkinerja tinggi (Go Fiber, PostgreSQL PostGIS, Next.js 16) tanpa biaya lisensi pihak ketiga berbayar (bebas Google Maps API / Twilio berbayar).
* **Slide 5 — Perjalanan Penjaminan Mutu (Quality Assurance):**
  Perjalanan audit ketat: *Deep Bug Audit $\rightarrow$ P0 Fundamental Fix $\rightarrow$ Live DR Drill $\rightarrow$ Operational UAT $\rightarrow$ IMP-001 Jurisdiction Fix $\rightarrow$ Regression Verified*.
* **Slide 6 — Keamanan Data & Pemisahan Wewenang (Security & SoD):**
  Perlindungan rekam medis pelapor (UU PDP No. 27/2022) dan penegakan programatik *Four-Eyes Principle* pada validasi hasil penyelidikan.
* **Slide 7 — Kedaulatan Medis pada DSS (Advisory Decision):**
  Sistem pakar memberikan rekomendasi objektif berbasis aturan Kemenkes, namun kedaulatan penanganan darurat tetap berada pada diskresi klinis Kepala Puskesmas (*Clinical Override*).
* **Slide 8 — Integritas Akuntansi Logistik Aset Daerah (BMD):**
  Penerapan mutasi transfer stok dua pihak (*StockTransfer*) secara atomik tanpa saldo minus, menjamin kepatuhan audit BPK/Inspektorat.
* **Slide 9 — Sistem Peringatan Dini (EWS & GIS Heatmap):**
  Visualisasi spasial sebaran jentik dan kasus aktif untuk memetakan kluster rawan wabah di 11 Kecamatan Kota Depok secara real-time.
* **Slide 10 — Lingkup Uji Coba Terbatas (Pilot Project Scope):**
  Fokus pelaksanaan percontohan di Puskesmas Beji dan Puskesmas Pancoran Mas selama 30 hari kalender.
* **Slide 11 — Transparansi Batasan Sistem (Known Limitations):**
  Ketergantungan infrastruktur server publik dan kebutuhan dokumentasi fisik BAST yang berjalan beriringan dengan sistem digital.
* **Slide 12 — Tiga Agenda Kebijakan Birokrasi (Pending Decisions):**
  Penyampaian objektif DEC-01 (alur stok obat), DEC-02 (kapasitas fogging KLB), dan DEC-03 (batas banding) untuk memperoleh pengesahan pimpinan.
* **Slide 13 — Rencana Aksi Lanjutan & Rekomendasi:**
  Permohonan penerbitan Surat Perintah Tugas (SPT) Pilot Project dan jadwal pelaksanaan *kick-off* lapangan.

---

## 19. Stakeholder Questions (FAQ)

Daftar jawaban resmi berbasis bukti atas pertanyaan yang berpotensi diajukan pimpinan daerah:

### Kategori 1: Tata Kelola & Bisnis (Governance & Business)
1. **Siapa pemilik resmi sistem ini?**
   * *Jawaban:* Dinas Kesehatan Kota Depok bertindak sebagai Pemilik Bisnis & Data (*Business & Data Owner*), sedangkan Diskominfo Kota Depok bertindak sebagai Pengelola Infrastruktur (*Infrastructure Host*). [`FINAL_READINESS` → Sec. 7]
2. **Apakah seluruh faskes wajib menggunakan aplikasi ini?**
   * *Jawaban:* Pada tahap pilot, kewajiban hanya berlaku bagi Puskesmas Beji dan Pancoran Mas berdasarkan Surat Perintah Tugas Kadinkes. Kewajiban menyeluruh akan diatur melalui Surat Edaran / Peraturan Walikota sebelum rilis kota penuh. [`RISK_MANAGEMENT_v2` → Sec. 7/RSK-07]
3. **Bagaimana pembagian tanggung jawab hukum jika terjadi kesalahan di lapangan?**
   * *Jawaban:* Jejak audit mencatat secara permanen siapa yang melakukan input data, verifikasi, atau persetujuan. Tanggung jawab operasional berada pada pejabat faskes yang menandatangani tiket secara digital. [`FINAL_READINESS` → Sec. 7]

### Kategori 2: Keamanan & Privasi Data (Security & PDP)
4. **Siapa saja yang dapat melihat data pribadi dan foto rekam medis pelapor?**
   * *Jawaban:* Akses dibatasi secara ketat berdasarkan peran dan wilayah kerja. Berkas foto hanya dapat dibuka oleh petugas faskes dan kader yang memiliki yurisdiksi atas wilayah laporan tersebut melalui token privat JWT. Publik hanya melihat nomor resi tanpa data PII. [`P0_FIX_REPORT` → Sec. 3/P0-02]
5. **Apakah kader dapat mengklaim laporan di luar kelurahannya?**
   * *Jawaban:* Tidak bisa. Melalui kontrol backend IMP-001, sistem secara mutlak menolak klaim tugas lintas wilayah dengan status HTTP 403 Forbidden. [`IMP001_VERIFICATION` → Sec. 6]

### Kategori 3: Logistik & Keuangan Daerah (Inventory & BMD)
6. **Bagaimana memastikan stok insektisida di aplikasi tidak bocor atau selisih dengan fisik gudang?**
   * *Jawaban:* Sistem menggunakan pencatatan ganda atomik. Mutasi antar-faskes wajib melalui status `IN_TRANSIT` yang memerlukan konfirmasi fisik dari faskes penerima. Setiap pemotongan stok wajib dikaitkan dengan tiket penugasan fogging yang memiliki bukti foto fisik. [`P0_FIX_REPORT` → Sec. 3/P0-01]
7. **Bagaimana jika barang yang dikirim dari Dinkes rusak di perjalanan?**
   * *Jawaban:* Faskes penerima memiliki hak menolak mutasi (*Reject Transfer*). Sistem secara otomatis mengembalikan (*refund*) saldo ke gudang asal dan mencatat insiden pada berita acara digital. [`P0_FIX_REPORT` → Sec. 3/P0-01/E]

### Kategori 4: Keputusan Klinis & DSS (Decision Support System)
8. **Apakah sistem dapat memaksa Puskesmas melakukan fogging jika dokter tidak setuju?**
   * *Jawaban:* Tidak. Rekomendasi DSS murni bersifat saran (*Advisory*). Kepala Puskesmas memegang wewenang penuh untuk menolak atau mengubah rekomendasi melalui fitur *Clinical Override* dengan mencantumkan alasan medis tertulis. [`P0_FIX_REPORT` → Sec. 3/P0-03]

### Kategori 5: Kesiapan Pilot & Operasional (Pilot Operations)
9. **Mengapa uji coba tahap pertama hanya dilakukan di Beji dan Pancoran Mas?**
   * *Jawaban:* Puskesmas Beji dan Pancoran Mas mewakili profil demografi perkotaan yang ideal untuk menguji beban kerja aplikasi, kesiapan kader relawan, serta koordinasi faskes secara terfokus sebelum diperluas ke 38 Puskesmas lainnya. [`FINAL_READINESS` → Sec. 8]
10. **Kapan evaluasi hasil pilot dilakukan?**
    * *Jawaban:* Evaluasi berkala dilakukan mingguan oleh Tim Surveilans Dinkes, dan evaluasi komprehensif penentu kelayakan rilis kota penuh dilakukan pada hari ke-30 masa uji coba pilot. [`FINAL_READINESS` → Sec. 16.2]

---

## 20. Known Limitations

Secara transparan dan bertanggung jawab, berikut adalah batasan operasional dan teknis sistem yang diketahui (*Known Limitations*) yang harus dipahami oleh pemangku kepentingan:

1. **Ketergantungan Saluran Komunikasi Publik Gratis:**
   Sesuai batasan mutlak proyek (100% Bebas Biaya Lisensi), sistem tidak menggunakan gateway SMS berbayar ataupun WhatsApp Business API resmi. Notifikasi mengandalkan email SMTP dan push notification Firebase (FCM). [`RISK_MANAGEMENT_v2` → Sec. 5/RSK-13]
2. **Ketergantungan Peladen Peta Publik OpenStreetMap:**
   Peta web memuat layer dari OSM publik. Rekomendasi teknis telah disampaikan ke Diskominfo untuk mengaktifkan local reverse caching proxy di Data Center Balai Kota guna mengantisipasi pembatasan IP saat rilis kota masif. [`RISK_MANAGEMENT_v2` → Sec. 5/RSK-18]
3. **Pembatasan Kaku Jadwal Fogging Reguler:**
   Kapasitas penjadwalan saat ini mengunci batas 1 tugas per hari per regu tanpa fitur *emergency override* (menunggu penetapan formal **DEC-02**). [`BUSINESS_DECISION` → Sec. 5]
4. **Alur Draf Permintaan Stok Otomatis:**
   Rekomendasi stok tertunda saat ini langsung menerbitkan tiket berstatus `PENDING` ke Dinas Kesehatan tanpa menunggu persetujuan internal Kepala Puskesmas (menunggu penetapan formal **DEC-01**). [`BUSINESS_DECISION` → Sec. 4]
5. **Ketiadaan Batas Waktu Temporal pada Hak Banding:**
   Warga dapat mengajukan banding 1 kali kapan pun sepanjang laporan berstatus ditolak hoaks, tanpa batas waktu 48 jam atau 72 jam (menunggu penetapan formal **DEC-03**). [`BUSINESS_DECISION` → Sec. 6]
6. **Kebutuhan Tertib Dokumen Fisik Berdampingan:**
   Pencatatan digital tidak menghapuskan kewajiban hukum penandatanganan fisik dokumen Berita Acara Serah Terima (BAST) barang milik daerah saat serah terima fisik jeriken obat. [`RISK_MANAGEMENT_v2` → Sec. 10/Residual-01]

---

## 21. Pilot Go / No-Go Checklist

Daftar periksa penentu kelayakan gerbang uji coba terbatas (*Go / No-Go Gate*) dengan pemisahan tegas antara kesiapan teknis yang **sudah diverifikasi (VERIFIED)** dan kesiapan operasional administratif yang **masih direncanakan (PARTIAL / PENDING)**:

| No | Parameter Kesiapan Gerbang Uji Coba | Status Kesiapan | Dasar Rujukan & Bukti Verifikasi | Catatan Kepatuhan & Syarat Pemenuhan |
|:---:|:---|:---:|:---|:---|
| 1 | **Technical Prerequisite PASS 100%** | `[X] GO (VERIFIED)` | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 3.1] | Backend Go compile 0 error, Frontend build 0 error, DB migration 0 pending. |
| 2 | **Security & Compliance PASS (RBAC, PII, IMP-001)** | `[X] GO (VERIFIED)` | [`IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md` → Sec. 7] | Celah P0/P1 tertutup; otorisasi yurisdiksi kader teruji via automated test suite. |
| 3 | **Business Workflow PASS (End-to-End Core)** | `[X] GO (VERIFIED)` | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 5] | Golden flow pelaporan warga s/d verifikasi logistik berjalan konsisten. |
| 4 | **Penyediaan Akun Operasional Faskes Pilot** | `[ ] PARTIAL (DIRENCANAKAN)` | [`DAFTAR_AKUN_DBD_GUARD.md` → Sec. 3] | Akun demo/seed tersedia; akun operasional riil personil faskes Beji & Panmas wajib di-provisioning pra-pilot. |
| 5 | **Pemetaan Wilayah Kelurahan Beji & Panmas** | `[X] GO (VERIFIED)` | [`IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md` → Sec. 4] | Master data poligon dan koordinat batas wilayah kerja 2 kelurahan terkunci valid pada migrasi master. |
| 6 | **Distribusi Fisik Dokumen SOP Operasional** | `[ ] PARTIAL (DIRENCANAKAN)` | Bab 12 dokumen ini & [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 & Sec. 12 #9] | Draf SOP telah siap; pencetakan fisik buku saku & sosialisasi nakes: TBD — ditetapkan setelah SPT terbit (Owner: Dinkes P2P & Tim Teknis). |
| 7 | **Saluran Bantuan WhatsApp Siaga Pilot** | `[ ] PARTIAL (DIRENCANAKAN)` | Bab 2.3 & 3.1 dokumen ini, [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 11] & [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 8.2 & Sec. 12 #8] | Struktur WAG dirancang; pembuatan grup fisik & penambahan nomor personil: TBD — ditetapkan setelah SPT terbit (Owner: Dinkes P2P). |
| 8 | **Prosedur Backup & Restore Teruji Mandiri** | `[X] GO (VERIFIED)` | [`scripts/backup_restore.ps1`] & [`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md` → Sec. 4.3] | Skrip operasional tersedia (6.960 bytes); drill restorasi terbukti sukses 100% match tanpa data loss. |
| 9 | **Distribusi Fisik Formulir Cadangan BCP** | `[ ] PARTIAL (DIRENCANAKAN)` | [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 8/BCP-01 & Sec. 11] & [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 & Sec. 12 #10] | Format Form PE Kemenkes disetujui; pencetakan fisik & penempatan di 2 faskes: TBD — ditetapkan setelah SPT terbit (Owner: Sanitarian Beji & Panmas). |
| 10 | **Isolasi Lingkup Terkunci di 2 Faskes Saja** | `[X] GO (VERIFIED)` | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 2] | Terkunci ketat di Puskesmas Beji & Puskesmas Pancoran Mas; 9 puskesmas lain dikecualikan. |
| 11 | **Surat Perintah Tugas (SPT) Resmi Kadinkes** | `[ ] PENDING / TBD (MENUNGGU)` | [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2 & Sec. 12 #7] | Menunggu penandatanganan dan penerbitan nomor surat resmi oleh Kepala Dinas Kesehatan Kota Depok. |
| 12 | **Pemahaman Keterbatasan Sistem oleh Stakeholder** | `[ ] PARTIAL / PENDING (DIRENCANAKAN)` | Bab 20 dokumen ini & [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2] | Batasan sistem didokumentasikan transparan; sosialisasi resmi dijadwalkan saat paparan eksekutif & sesi pengarahan teknis. |

---

### KEPUTUSAN FINAL: CONDITIONAL GO

> **KEPUTUSAN GERBANG PILOT (GO / NO-GO GATE):**
> Sistem DBD Guard dinyatakan berstatus **`CONDITIONAL GO`** untuk Uji Coba Terbatas (*Limited Pilot*) di 2 Faskes (Puskesmas Beji & Puskesmas Pancoran Mas).
> 
> **Sistem DILARANG KERAS dioperasikan secara operasional di lapangan sebelum TUJUH (7) SYARAT PRA-PELUNCURAN berikut dipenuhi dan diverifikasi seluruhnya:**
> 1. **Penerbitan & Penandatanganan SPT Resmi Kadinkes:** Surat Perintah Tugas resmi diterbitkan dan ditandatangani oleh Kepala Dinas Kesehatan Kota Depok sebagai payung hukum penugasan pilot [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2 & Sec. 12 #7].
> 2. **Provisioning & Verifikasi Akun Operasional Personil Riil:** Pembuatan dan aktivasi kredensial pengguna riil untuk staf Puskesmas Beji, Puskesmas Pancoran Mas, dan kader jumantik terkait pada database operasional [`DAFTAR_AKUN_DBD_GUARD.md` → Sec. 3 & `FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6].
> 3. **Pencetakan & Penempatan Fisik Formulir Cadangan BCP:** Formulir penyelidikan epidemiologi manual standar Kemenkes rangkap 3 dan register BAST logistik telah dicetak dan tersedia di meja kerja sanitarian kedua faskes [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 8/BCP-01 & `FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 & Sec. 12 #10].
> 4. **Penggandaan & Penyerahan Fisik Buku Saku SOP Operasional:** Dokumen SOP Operasional (Bab 12) dicetak dalam bentuk buku saku ringkas dan dibagikan kepada seluruh personil faskes yang terlibat [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6 & Sec. 12 #9].
> 5. **Pembentukan & Aktivasi Kanal WhatsApp Group Siaga Pilot:** Pembentukan grup WhatsApp operasional siaga yang beranggotakan perwakilan Dinkes, Kepala Puskesmas, Sanitarian, Pengelola Logistik, dan Tim Teknis Pengembang dengan nomor telepon riil [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 11 & `FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 8.2 & Sec. 12 #8].
> 6. **Penyelenggaraan Sesi Pengarahan Teknis (Briefing) & Simulasi Pengguna:** Penyelenggaraan sesi pengarahan teknis setengah hari dan simulasi terpandu bagi para Kepala Puskesmas, sanitarian, regu fogging, dan perwakilan kader RW terpilih sebelum peluncuran lapangan [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2 #3 & Sec. 12 #11].
> 7. **Pemenuhan Batasan Operasional Khusus (Karantina Kebijakan DEC-01/02/03 & Isolasi 2 Faskes):** Operasi terbatas terkunci hanya di Puskesmas Beji dan Puskesmas Pancoran Mas dalam kondisi stok aman ($\ge 20$ L), status non-KLB, serta penanganan sengketa via tim siaga Dinkes [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2 #2 & Sec. 11.1].

---

## 22. Final Readiness & Discrepancy Note

### 22.1 Status Kesiapan Akhir Resmi:
$$\mathbf{CONDITIONAL\ GO\ —\ CONDITIONALLY\ PILOT\ READY\ (LIMITED\ PILOT\ ONLY)}$$
$$\mathbf{\&}$$
$$\mathbf{DEMO\ READY\ —\ PENDING\ DRY-RUN\ (EXECUTIVE\ DEMO)}$$

### 22.2 Catatan Rekonsiliasi & Resolusi Diskrepansi (*Discrepancy Note*):
1. **Prioritas Dokumen Acuan:**
   Jika terdapat perbedaan interpretasi antara instruksi awal dan kondisi kode/dokumen rujukan, maka ketentuan yang tertuang dalam dokumen [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md`], [`RISK_MANAGEMENT_DBD_GUARD_v2.md`], dan [`BUSINESS_DECISION_GOVERNANCE_DBD_GUARD.md`] berlaku sebagai kebenaran mutlak (*Source of Truth*).
2. **Karantina Isu Kebijakan:**
   Ketiga isu tata kelola birokrasi (**DEC-01, DEC-02, DEC-03**) secara konsisten diposisikan sebagai **`PENDING STAKEHOLDER DECISION`**. Status ini bersifat **NON-BLOCKING** bagi pelaksanaan uji coba terbatas di Puskesmas Beji dan Puskesmas Pancoran Mas karena volume kasus berada pada batas aman reguler, namun menjadi **BLOCKER MUTLAK** bagi peluncuran penuh se-Kota Depok.
3. **Ketiadaan Klaim Absolut:**
   Sistem dinyatakan siap uji coba lapangan secara bertanggung jawab berdasarkan bukti pengujian empiris, tanpa klaim keliru seperti "100% bebas risiko" atau "100% kebal peretasan". Seluruh risiko residual telah dipetakan tindakan mitigasinya.

---

## 23. Action Tracker

Matriks pemantauan tindak lanjut komprehensif untuk seluruh syarat Conditional Go, prasyarat demonstrasi, keputusan tata kelola, dan artefak tertunda:

| No | Item Tindak Lanjut | Kategori | Penanggung Jawab (*Owner*) | Tenggat Waktu | Status | Evidence / Catatan Rujukan |
|:---:|:---|:---:|:---|:---:|:---:|:---|
| **1** | Penerbitan SPT Resmi Pilot Kadinkes | Syarat Conditional Go | Kepala Dinas Kesehatan Kota Depok | TBD — ditetapkan pimpinan | `PENDING` | Lembar disposisi / SK Kadinkes [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2] |
| **2** | Provisioning Akun Riil Nakes & Kader | Syarat Conditional Go | Administrator Dinkes & Koordinator Faskes | TBD — setelah SPT terbit | `PENDING` | Log aktivasi user riil di database [`DAFTAR_AKUN_DBD_GUARD.md` → Sec. 3] |
| **3** | Pencetakan & Distribusi Formulir Fisik BCP | Syarat Conditional Go | Sanitarian Puskesmas Beji & Panmas | TBD — setelah SPT terbit | `PENDING` | Tanda terima fisik formulir rangkap 3 [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 8/BCP-01] |
| **4** | Penggandaan & Penyerahan Buku Saku SOP | Syarat Conditional Go | Tim Teknis & Dinkes P2P | TBD — setelah SPT terbit | `PENDING` | Bukti serah terima buku saku fisik [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 6] |
| **5** | Pembentukan WhatsApp Group Siaga Pilot | Syarat Conditional Go | Tim Surveilans Dinkes P2P | TBD — setelah SPT terbit | `PENDING` | Tautan & daftar anggota nomor riil WAG [`RISK_MANAGEMENT_DBD_GUARD_v2.md` → Sec. 11] |
| **6** | Briefing Teknis & Simulasi Pengguna Lapangan | Syarat Conditional Go | Tim Teknis Pengembang & Dinkes | TBD — sebelum hari-H pilot | `PENDING` | Daftar hadir & notula sesi briefing [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 16.2 #3] |
| **7** | Verifikasi Batasan Non-KLB & Stok $\ge 20$ L | Syarat Conditional Go | Kepala Puskesmas Beji & Panmas | TBD — sebelum hari-H pilot | `PENDING` | Berita acara cek fisik gudang & status non-KLB [`FINAL_PILOT_READINESS_ASSESSMENT_DBD_GUARD.md` → Sec. 11.1] |
| **8** | Pembuatan Seed Dump Sintetis Bebas Data Nyata | Prasyarat Demo | Tim Database / Backend | TBD — pra-demo | `PENDING` | File `dbd_guard_demo_seed.sql.gz` terbentuk & diverifikasi tanpa PII nyata |
| **9** | Uji Durasi Restorasi ke DB Demo | Prasyarat Demo | Tim DevOps / Infrastruktur | TBD — pra-demo | `PENDING` | Log durasi riil eksekusi `backup_restore.ps1` ke `dbd_guard_restore_test` |
| **10** | Dry-Run Penuh 2x dengan Stopwatch | Prasyarat Demo | Lead Presenter & Tim Teknis | TBD — pra-demo | `PENDING` | Log catatan waktu lolos/gagal per langkah $\le 35$ menit |
| **11** | Pembuatan Video Walkthrough Cadangan | Prasyarat Demo | Tim Dokumentasi & Frontend | TBD — pra-demo | `PENDING` | File video MP4 tersimpan di folder `docs/demo_artifacts/` |
| **12** | Validasi Token & Script Negative Test cURL | Prasyarat Demo | Tim QA / Backend | TBD — pra-demo | `PENDING` | Script cURL teruji menghasilkan HTTP 403 / token JWT valid |
| **13** | Keputusan Final Moda per Skenario Demo | Prasyarat Demo | Lead Presenter | TBD — pasca dry-run | `PENDING` | Matriks moda definitif (LIVE / REKAMAN / FALLBACK) |
| **14** | Penetapan Kebijakan DEC-01 (Draf Stok Otomatis) | Tata Kelola (Governance) | Kadinkes & Kepala Puskesmas | TBD — sebelum rilis kota | `PENDING` | Lembar pengesahan formal Opsi A atau Opsi B [`BUSINESS_DECISION_GOVERNANCE_DBD_GUARD.md` → Sec. 4] |
| **15** | Penetapan Kebijakan DEC-02 (Kapasitas Fogging KLB) | Tata Kelola (Governance) | Kadinkes & Tim Anggaran | TBD — sebelum rilis kota | `PENDING` | Lembar pengesahan formal Opsi A atau Opsi B [`BUSINESS_DECISION_GOVERNANCE_DBD_GUARD.md` → Sec. 5] |
| **16** | Penetapan Kebijakan DEC-03 (Batas Waktu Banding) | Tata Kelola (Governance) | Kadinkes & Tim Hukum Pemkot | TBD — sebelum rilis kota | `PENDING` | Lembar pengesahan formal Opsi A, B, atau C [`BUSINESS_DECISION_GOVERNANCE_DBD_GUARD.md` → Sec. 6] |
| **17** | Pembuatan Artefak Fisik Seed Dump & Folder Demo | Artefak / File | Tim Teknis Pengembang | TBD — pra-demo | `PENDING` | Keberadaan file `dbd_guard_demo_seed.sql.gz` & direktori `docs/demo_artifacts/` |

---
*Dokumen Pilot Preparation & Demonstration Package ini disahkan secara digital sebagai paket panduan resmi pelaksanaan demonstrasi eksekutif dan peluncuran uji coba terbatas sistem DBD Guard Kota Depok.*
