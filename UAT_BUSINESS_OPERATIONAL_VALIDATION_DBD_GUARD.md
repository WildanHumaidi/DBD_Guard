# Business & Operational Validation (UAT) — DBD Guard
**Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi (Pemerintah Kota Depok)**

- **Tanggal Pelaksanaan:** 29 September 2026
- **Lead UAT Assessor:** Antigravity AI (Pair Programming System Evaluator)
- **Status Rilis Teknis Sebelumnya:** `PASS` (P0 Fundamental Fix & Deep Bug Audit Re-Verification)
- **Status Akhir UAT Bisnis & Operasional:** **CONDITIONALLY BUSINESS READY**

---

## 1. Executive Summary

Setelah tahapan audit kode mendalam (*Deep Bug Audit* #1 dan #2), perbaikan fundamental level P0 (*P0 Fundamental Fix*), pengujian bertingkat (*Unit*, *Integration*, *End-to-End*, dan *Disaster Recovery Restoration Drill*), sistem **DBD Guard** telah mencapai status kelulusan teknis (**Technical PASS**). 

Tahap **Business & Operational Validation / User Acceptance Testing (UAT)** ini dilakukan secara komprehensif untuk menjawab pertanyaan fundamental:
> *"Apakah sistem yang secara teknis sudah PASS benar-benar menjalankan proses bisnis DBD Guard sesuai requirement, business rules, peran kewilayahan, alur operasional lapangan, dan dinamika nyata stakeholder kesehatan di Kota Depok?"*

### Hasil Validasi Utama:
1. **Integritas Alur End-to-End (FLOW-01):** Alur utama dari pelaporan warga (Kecamatan $\rightarrow$ Kelurahan $\rightarrow$ RW), penugasan dan investigasi lapangan oleh Kader Jumantik, evaluasi Penyelidikan Epidemiologi (PE), sistem penasihat rekomendasi (*Decision Support System* / DSS), persetujuan Kepala Puskesmas, eksekusi Tim Fogging, hingga pemotongan stok otomatis di gudang farmasi Puskesmas telah **berhasil dijalankan secara koheren dan konsisten**.
2. **Kedaulatan Keputusan Manusia (*Human Authority*):** DSS terbukti bertindak murni sebagai alat bantu rekomendasi (*advisory*), bukan pengganti wewenang klinis/manajerial (*decision replacement*). Hak *override* Kepala Puskesmas dan Dinkes berfungsi dengan pencatatan audit alasan wajib (*mandatory override reason*).
3. **Pemisahan Tugas (*Segregation of Duties* / *Four-Eyes Principle*):** Investigasi PE dan validasi hasil PE terbukti tidak dapat dieksekusi oleh individu yang sama, menjamin objektivitas intervensi kesehatan lingkungan.
4. **Akuntansi Logistik (*Inventory Zero-Leakage*):** Perpindahan insektisida dan larvasida dari Gudang Utama Dinkes ke Puskesmas melalui status `IN_TRANSIT` mempertahankan prinsip konservasi stok absolut ($\text{Stok Asal} + \text{Stok Tujuan} + \text{In-Transit} = \text{Konstan}$). Pembatalan transfer mengembalikan stok secara akurat (*refund reversal*).
5. **Siklus Banding Warga (*Appeal Unblocking*):** Laporan yang ditolak sebagai Hoax berhasil diajukan banding, diverifikasi ulang oleh Puskesmas, diambil alih (*takeover*), dan dilanjutkan ke PE tanpa terjadi *deadlock state*.
6. **Kesenjangan Operasional & Tata Kelola (*Findings*):** Ditemukan sejumlah catatan operasional penting (P1 & P2), antara lain:
   - Validasi yurisdiksi batas wilayah kerja (*location boundary check*) pada endpoint klaim tugas Kader (`ClaimReport`) belum dikunci ketat di level usecase API (baru difilter pada tampilan antarmuka UI).
   - Pengajuan stok darurat (*auto-draft stock request*) dari rekomendasi `FOGGING_TERTUNDA` langsung terbit dengan status `PENDING` ke Dinkes tanpa melalui penandatanganan/konfirmasi formal Kepala Puskesmas.
   - Kebijakan kapasitas fogging *anti-overbooking* (maksimal 1 tugas per hari per tim) berpotensi menjadi *bottleneck* operasional saat terjadi lonjakan kasus (*outbreak surge*) atau penetapan status KLB.

---

## 2. Test Scope

Pengujian UAT mencakup seluruh domain fungsional dan operasional yang diatur di dalam PRD dan Regulasi Kota Depok:

```text
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           CAKUPAN UAT DBD GUARD                                 │
├─────────────────────┬───────────────────────────────┬───────────────────────────┤
│ Stakeholder / Role  │ Domain Bisnis                 │ Alur Validasi             │
├─────────────────────┼───────────────────────────────┼───────────────────────────┤
│ 1. Warga Pelapor    │ Pelaporan Publik & Privasi    │ Input, Geotag, Anonim,    │
│                     │                               │ Lacak Ref, Banding, Cancel│
│ 2. Kader Jumantik   │ Verifikasi & PE Tingkat RW    │ Klaim, Tolak Hoax, Input  │
│                     │                               │ PE, Edukasi PSN, SLA 72 Jam│
│ 3. Puskesmas        │ Manajemen Medis & Logistik    │ Validasi PE, Override DSS,│
│                     │                               │ Takeover, Approval Jadwal │
│ 4. Tim Fogging      │ Eksekusi Lapangan & Bukti     │ Terima Tugas, Update OTW/ │
│                     │                               │ Selesai, Pemakaian Stok   │
│ 5. Dinas Kesehatan  │ Pengawasan Kota & EWS/KLB     │ Agregasi GIS, Buffer Stok,│
│                     │                               │ Ambang Batas KLB, Audit   │
└─────────────────────┴───────────────────────────────┴───────────────────────────┘
```

---

## 3. Source of Truth

Seluruh pengujian, ekspektasi, dan klasifikasi hasil UAT didasarkan secara ketat pada dokumen referensi resmi proyek:
1. **PRD DBD Guard Komprehensif** (`PRD_DBD_Guard_Komprehensif.md`)
2. **Aturan Bisnis & Domain Kota Depok** (`.agents/rules/business_rules.md`)
3. **Arsitektur & Batasan Kode** (`.agents/rules/structure.md`, `code_style.md`)
4. **Dokumen Rekayasa Kebutuhan Terpadu v1.1** (`docs/DBD_Guard_Dokumen_Terpadu_Rekayasa_Kebutuhan_v1_1.docx`)
5. **Laporan Audit Manajemen Risiko** (`RISK_MANAGEMENT_DBD_GUARD.md`)
6. **Laporan P0 Fix & Re-Verifikasi Audit** (`P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md`)
7. **Skema Database & Relasi Kontrak** (PostgreSQL + PostGIS)

---

## 4. Test Environment

- **Sistem Operasi:** Windows 11 Pro / PowerShell 7
- **Database Engine:** PostgreSQL 16 dengan Ekstensi PostGIS (Spasial)
- **Backend Runtime:** Go 1.22+ (Fiber Framework, Modular Monolith)
- **Frontend Web:** Next.js 16.3.4 (React 19, Tailwind CSS, Mode PWA)
- **Peta Spasial:** Leaflet.js + OpenStreetMap (OSM Tile Service)
- **Penyimpanan Berkas:** Local Encrypted File Store (`./storage/photos`, RBAC-gated)
- **Layanan Notifikasi:** Firebase Cloud Messaging (FCM) & SMTP Mockup Service

---

## 5. Role-Based UAT

### 5.1 Warga (Pelapor)

#### Scenario W-01: Submit Citizen Report (Identified)
- **Deskripsi:** Warga melaporkan dugaan kasus demam berdarah dengan mengisi data diri lengkap, alamat, RW/RT, kelurahan, gejala, foto bukti surat diagnosis/lingkungan, serta titik lokasi GPS.
- **Expected Business Behavior:** 
  1. Data tersimpan pada tabel `reporters`, `citizen_reports`, `report_locations`, dan `report_photos` secara atomik.
  2. Diterbitkan *Reference Number* unik format `REP-YYYYMMDD-XXXXXXXX`.
  3. Status awal laporan adalah `MENUNGGU_PE_KADER`.
  4. Laporan yang berjarak <100 meter dari laporan aktif dalam kurun 7 hari otomatis ditandai `is_possible_duplicate = true`.
  5. Pengiriman laporan berulang dalam kurun 1 jam pada radius <50m ($\ge 3$ laporan) otomatis ditandai `REJECTED` (*Anti-Spam*).
  6. Notifikasi diteruskan ke Kader di RW/Kelurahan terkait.
- **Actual System Behavior:** Berhasil 100%. Data tersimpan utuh dalam transaksi, No. Referensi terbit, status `MENUNGGU_PE_KADER`, flag duplikasi spasial aktif, dan spam clustering menolak laporan liar secara otomatis.
- **Status:** `PASS`
- **Evidence:** Uji `TestAnonymousReport_DirectAPISubmission_EndToEndInvariants` dan fungsi `SubmitCitizenReport` pada `report_usecase.go`.

#### Scenario W-02: Anonymous Report (Privasi PII Terjamin)
- **Deskripsi:** Warga memilih opsi *Lapor Anonim* (`anonymous_flag = true`) untuk menjaga kerahasiaan identitas pasien di lingkungan tempat tinggal.
- **Expected Business Behavior:**
  1. Nama pemohon di database dipaksa menjadi `"Anonymous"`.
  2. Nomor kontak/telepon dikosongkan total (`""`).
  3. Tidak ada SMS, Email, maupun push notification yang dikirimkan ke pelapor.
  4. Pelapor tetap dapat melacak progres laporan secara mandiri menggunakan Nomor Referensi.
  5. Pengingat evaluasi pasca-fogging (hari ke-14) otomatis dilewati (*skip*) agar tidak membocorkan identitas.
- **Actual System Behavior:** Berhasil 100%. Sanitasi PII berjalan di layer usecase sebelum persistensi data. No. Referensi tetap aktif dan dapat dilacak di `/track`.
- **Status:** `PASS`
- **Evidence:** Uji `TestAnonymousReport_EndToEndPIIProtection` pada `bug10_anonymous_reporting_test.go`.

#### Scenario W-03: Track Report & Isolasi Akses Data
- **Deskripsi:** Warga memeriksa status tindak lanjut laporannya melalui portal `/track`.
- **Expected Business Behavior:**
  1. Pencarian via Nomor Referensi menampilkan status riil dan histori tahapan penanganan.
  2. Pencarian alternatif via Nomor Handphone hanya mengembalikan ringkasan umum (No. Ref, Tanggal, Status, Kategori, Kecamatan).
  3. Data sensitif (nama pasien, foto surat rekam medis, alamat detail tingkat nomor rumah) tidak diekspos kepada publik.
- **Actual System Behavior:** Berhasil. Endpoint `/reports/phone/:phone` mengembalikan DTO `ReportSummaryResponse` yang mengecualikan alamat detail dan nama pasien. Berkas foto dilindungi otentikasi JWT.
- **Status:** `PASS`
- **Evidence:** `TrackReportsByPhone` pada `report_usecase.go` dan `auth_middleware.go`.

#### Scenario W-04: Appeal Rejection (Banding Warga)
- **Deskripsi:** Warga yang laporannya ditolak sebagai Hoax mengajukan banding dengan melampirkan keterangan tambahan bahwa pasien terbukti positif di rumah sakit.
- **Expected Business Behavior:**
  1. Tombol banding hanya muncul jika status laporan `REJECTED` atau `REJECTED_HOAX`.
  2. Banding hanya dapat diajukan maksimal 1 kali (`appeal_count < 1`).
  3. Status berubah menjadi `APPEALED` dan dialihkan ke antrean telaah Puskesmas.
  4. Pengajuan banding kedua otomatis ditolak sistem.
- **Actual System Behavior:** Berhasil. Kuota banding terkunci secara atomik di database (`SubmitAppealAtomic`). Upaya banding berulang mengembalikan HTTP 400.
- **Status:** `PASS`
- **Evidence:** Uji `TestBUGREG001_FullLifecycle_AppealToPEValidation` dan `SubmitAppealAtomic` di `report_repository.go`.

#### Scenario W-05: Cancel Report
- **Deskripsi:** Warga membatalkan laporannya sendiri karena salah memasukkan informasi alamat.
- **Expected Business Behavior:**
  1. Pembatalan hanya diizinkan jika laporan masih berstatus `MENUNGGU_PE_KADER` dan belum diklaim oleh Kader (`investigator_id == nil`).
  2. Apabila Kader telah mengklaim tugas lapangan, pembatalan mandiri oleh warga dilarang untuk mencegah hilangnya akuntabilitas waktu operasional petugas.
- **Actual System Behavior:** Berhasil. API menolak pembatalan dengan pesan informatif jika tugas sudah diambil oleh Kader.
- **Status:** `PASS`
- **Evidence:** `CancelReport` pada `report_usecase.go` baris 290–309.

---

### 5.2 Kader (Verifikator Lapangan & Jumantik)

#### Scenario K-01: Claim Report (Ambil Tugas Penyelidikan)
- **Deskripsi:** Kader membuka antarmuka verifikasi RW dan mengklaim tugas penyelidikan lapangan.
- **Expected Business Behavior:**
  1. Kader hanya melihat daftar laporan di yurisdiksi wilayah tugasnya.
  2. Dua kader yang mengklik tombol klaim pada milidetik yang sama tidak boleh menyebabkan tugas ganda (*race condition*).
  3. Satu kader berhasil menjadi pemilik tugas, kader lainnya menerima notifikasi bahwa laporan telah diambil.
- **Actual System Behavior:**
  - Penanganan konkurensi: `PASS`. Operasi database atomik `WHERE report_id = ? AND investigator_id IS NULL` mencegah *double-claim*.
  - Batasan Yurisdiksi Wilayah: `PARTIAL`. Filter wilayah berjalan dengan baik pada kueri `GetReports` di UI, namun endpoint `/api/v1/kader/reports/:id/claim` tidak memvalidasi ulang kecocokan `user.LocationID` dengan `report.LocationID`.
- **Status:** `PARTIAL`
- **Evidence:** Uji `TestBUGNEW008_ClaimReport_AtomicRaceCondition` lulus; Analisis usecase pada `ClaimReport` (`report_usecase.go:273`).

#### Scenario K-02: Verifikasi & Penolakan Laporan Hoax
- **Deskripsi:** Kader mendatangi lokasi dan mendapati bahwa alamat fiktif atau warga sekitar tidak pernah ada yang terjangkit demam.
- **Expected Business Behavior:**
  1. Kader dapat menandai laporan sebagai Hoax dengan menyertakan alasan penolakan wajib.
  2. Status laporan bermutasi menjadi `REJECTED_HOAX`.
  3. Kader dilarang menggunakan endpoint ini untuk langsung "menyetujui" laporan tanpa melakukan input formulir PE lengkap.
- **Actual System Behavior:** Berhasil 100%. Sistem menolak request jika `decision != "REJECTED"`, memaksa kader mengisi PE formal untuk meneruskan laporan.
- **Status:** `PASS`
- **Evidence:** `verification_usecase.go` baris 56–63.

#### Scenario K-03: Penyelidikan Epidemiologi (PE Lapangan)
- **Deskripsi:** Kader melakukan pemeriksaan jentik nyamuk dan pendataan penderita demam dalam radius 100 meter.
- **Expected Business Behavior:**
  1. Pemeriksaan wajib menyertakan edukasi 3M Plus / PSN (`psn_edukasi_done = true`).
  2. Hanya kader yang mengklaim laporan (`investigator_id == kader.ID`) yang diizinkan menginput hasil PE.
  3. Sistem mengalkulasi rekomendasi awal DSS dan memutasi status laporan menjadi `MENUNGGU_VALIDASI_PUSKESMAS`.
  4. Tercatat riwayat audit `SUBMIT_PE`.
- **Actual System Behavior:** Berhasil 100%. Laporan berpindah status ke antrean Puskesmas, validasi kepemilikan bertindak tegas, dan riwayat audit tersimpan.
- **Status:** `PASS`
- **Evidence:** `pe_usecase.go` baris 41–129.

#### Scenario K-04: Eskalasi Otomatis Batas Waktu SLA (72 Jam)
- **Deskripsi:** Laporan warga diabaikan oleh Kader/RT selama lebih dari 3x24 jam (72 jam) tanpa ada klaim atau verifikasi.
- **Expected Business Behavior:**
  1. Layanan cron/auto-release mendeteksi laporan yang melewati batas waktu 72 jam.
  2. Status laporan otomatis bermutasi menjadi `ESKALASI_PUSKESMAS`.
  3. Petugas yang ditugaskan sebelumnya (jika ada yang mandek) dicabut (`investigator_id = nil`).
  4. Laporan masuk ke radar intervensi darurat Puskesmas dengan label *"RT Tidak Responsif"*.
  5. Audit log mencatat peristiwa `SLA_ESCALATION`.
- **Actual System Behavior:** Berhasil 100%. Simulasi batas waktu 72 jam mengeksekusi eskalasi secara tepat dan idempoten (tidak memicu eskalasi berulang pada pemanggilan berikutnya).
- **Status:** `PASS`
- **Evidence:** Uji `TestAutoReleaseExpiredTasks_72hSLA` pada `sla_usecase_test.go`.

---

### 5.3 Petugas & Kepala Puskesmas

#### Scenario P-01: Takeover Laporan Tereskalasi
- **Deskripsi:** Petugas Puskesmas mengambil alih laporan yang tereskalasi akibat kader pasif atau laporan hasil banding warga.
- **Expected Business Behavior:**
  1. Petugas Puskesmas memiliki wewenang untuk mengeksekusi *takeover* pada laporan berstatus `MENUNGGU_PE_KADER`, `ESKALASI_PUSKESMAS`, atau `APPEALED`.
  2. Sistem mencatat ID pemeriksa sebelumnya, ID pemeriksa baru, dan alasan pengambilalihan pada *audit trail*.
  3. Laporan yang sudah selesai (`COMPLETED`) dilarang keras untuk diambil alih.
- **Actual System Behavior:** Berhasil 100%. Petugas Puskesmas berhasil mengklaim laporan tereskalasi dan menyerahkan hasil PE secara legal.
- **Status:** `PASS`
- **Evidence:** Uji `TestTakeoverWorkflow` pada `takeover_usecase_test.go`.

#### Scenario P-02: Validasi PE & Prinsip Pemisahan Tugas (Four-Eyes Principle)
- **Deskripsi:** Puskesmas meninjau hasil investigasi lapangan kader sebelum menentukan intervensi lingkungan.
- **Expected Business Behavior:**
  1. Petugas yang melakukan PE dilarang memvalidasi pekerjaannya sendiri (*Segregation of Duties*).
  2. Validasi harus dilakukan oleh pejabat yang berbeda (Kepala Puskesmas / Dokter Penanggung Jawab).
- **Actual System Behavior:** Berhasil 100%. Jika `peResult.OfficerID == validatorID`, usecase menolak eksekusi dengan pesan tegas `pelanggaran pemisahan tugas: pemeriksa PE dilarang memvalidasi laporannya sendiri (Four-Eyes Principle)`.
- **Status:** `PASS`
- **Evidence:** Uji `TestP0_03_and_04_DSS_and_SoD` pada `p0_verification_test.go`.

#### Scenario P-03: Rekomendasi DSS Terbimbing
- **Deskripsi:** Sistem mengevaluasi variabel epidemiologi lapangan dan stok gudang farmasi Puskesmas.
- **Expected Business Behavior:**
  1. Jika Jentik = "Ya" atau Penderita Demam $\ge 3$, dan Stok Insektisida Puskesmas $\ge 10$ Liter $\rightarrow$ Rekomendasi: `FOGGING_FOKUS`.
  2. Jika Jentik = "Ya" atau Penderita Demam $\ge 3$, dan Stok Insektisida $< 10$ Liter $\rightarrow$ Rekomendasi: `FOGGING_TERTUNDA` (Sistem menginisiasi draf permintaan stok ke Dinkes).
  3. Jika Jentik = "Tidak" dan Penderita Demam $< 3$ $\rightarrow$ Rekomendasi: `EDUKASI_PSN` (Kasus ditutup tanpa fogging, status laporan: `CLOSED_NO_FOGGING`).
- **Actual System Behavior:** Berhasil 100%. Perhitungan rekomendasi berjalan deterministik sesuai matriks PRD.
- **Status:** `PASS`
- **Evidence:** `pe_usecase.go` baris 69–85 dan 178–190.

#### Scenario P-04: Human Advisory Override (Pengesampingan Rekomendasi DSS)
- **Deskripsi:** Kepala Puskesmas menilai kondisi lapangan di pemukiman padat memerlukan fogging fokus meskipun jentik tidak ditemukan pada saat kunjungan siang hari.
- **Expected Business Behavior:**
  1. Kepala Puskesmas dapat mengesampingkan output DSS (misal dari `EDUKASI_PSN` menjadi `FOGGING_FOKUS`).
  2. Pengesampingan mewajibkan pengisian alasan tertulis (*override reason*).
  3. Sistem mencatat audit log khusus `OVERRIDE_DSS_RECOMMENDATION` yang merekam rekomendasi awal sistem, keputusan akhir manusia, dan alasan justifikasi klinis.
  4. Pengesampingan tanpa alasan wajib ditolak oleh API.
- **Actual System Behavior:** Berhasil 100%. Validasi input menolak override tanpa alasan. Audit log mencatat rekaman komparatif dengan lengkap.
- **Status:** `PASS`
- **Evidence:** Uji `TestP0_03_DSS_Override_WithoutReason_MustFail` dan `TestP0_03_DSS_Override_WithReason_Success`.

#### Scenario P-05: Persetujuan Jadwal Fogging (*Approval*)
- **Deskripsi:** Kepala Puskesmas menyetujui pengajuan fogging (`domain.FoggingRequest`) agar dapat dialokasikan ke jadwal tim pelaksana.
- **Expected Business Behavior:**
  1. Hanya Kepala Puskesmas yang berwenang menyetujui pengajuan fogging.
  2. Penjadwalan hanya dapat dilakukan setelah persetujuan berstatus `APPROVED`.
  3. Waktu pelaksanaan fogging wajib minimal 24 jam setelah pengajuan disetujui (Business Rule BR5: persiapan warga menutup makanan dan mengamankan hewan peliharaan).
- **Actual System Behavior:** Berhasil 100%. Endpoint `/fogging-schedules` menolak pembuatan jadwal sebelum 24 jam dari waktu persetujuan.
- **Status:** `PASS`
- **Evidence:** `fogging_schedule_usecase.go` baris 76–78.

---

### 5.4 Tim Fogging (Regu Lapangan)

#### Scenario F-01: Penerimaan Surat Penugasan Kerja
- **Deskripsi:** Regu fogging menerima notifikasi tugas kerja harian beserta detail koordinat lokasi, RW/RT sasaran, dan volume estimasi bahan.
- **Expected Business Behavior:**
  1. Jadwal tampil di dashboard Tim Fogging `/fogging/schedules`.
  2. Tim dapat memantau status persiapan logistik.
  3. Penugasan tidak melebihi kapasitas kerja regu (1 tugas per regu per hari).
- **Actual System Behavior:** Berhasil 100%. Antarmuka web PWA responsif menyajikan penugasan dengan jelas.
- **Status:** `PASS`
- **Evidence:** Halaman `frontend-web/src/app/fogging/schedules/page.tsx`.

#### Scenario F-02: Eksekusi Fogging & Pemotongan Logistik Atomik
- **Deskripsi:** Regu fogging tiba di lokasi, menyemprot area sasaran, mendokumentasikan foto bukti, dan menyelesaikan tugas.
- **Expected Business Behavior:**
  1. Tim mengupdate status dari `SCHEDULED` $\rightarrow$ `ON_THE_WAY` $\rightarrow$ `IN_PROGRESS` $\rightarrow$ `COMPLETED`.
  2. Saat memilih status `COMPLETED`, tim **wajib** mengisi rincian barang yang digunakan (misal: Malathion/Insektisida 2 Liter) dan mengunggah foto bukti penyemprotan.
  3. Stok di gudang farmasi Puskesmas terkait otomatis terpotong secara atomik pada saat tombol "Selesai" ditekan.
  4. Laporan warga bertransisi menjadi `COMPLETED`.
- **Actual System Behavior:** Berhasil 100%. Pemotongan stok dilakukan di dalam transaksi database tunggal bersamaan dengan pembuatan `InventoryUsageLog`.
- **Status:** `PASS`
- **Evidence:** `fogging_schedule_usecase.go` baris 172–210 dan uji `TestStockDeductionAtomic`.

#### Scenario F-03: Penyelesaian Tugas Tidak Sah (Tanpa Logistik / Tanpa Bukti)
- **Deskripsi:** Tim fogging mencoba mengeklik "Selesai" tanpa mengunggah foto bukti lapangan atau tanpa mengisi volume bahan yang terpakai.
- **Expected Business Behavior:**
  1. Backend menolak penyelesaian tugas dengan pesan `penggunaan inventory wajib diisi saat menyelesaikan jadwal fogging`.
  2. Frontend mengunci tombol submit dan menampilkan pesan validasi bahwa bukti foto dan pemakaian logistik adalah syarat mutlak penyelesaian.
- **Actual System Behavior:** Berhasil 100%. Validasi ganda (frontend & backend) menggagalkan upaya manipulasi penyelesaian tugas fiktif.
- **Status:** `PASS`
- **Evidence:** `fogging_schedule_usecase.go:149` dan `frontend-web/src/app/fogging/schedules/page.tsx:80-90`.

---

### 5.5 Dinas Kesehatan (Dinkes Kota Depok)

#### Scenario D-01: Pengawasan Spasial & Pemantauan Agregat
- **Deskripsi:** Pengelola program DBD Dinkes memantau penyebaran kasus se-Kota Depok melalui dashboard GIS Heatmap dan Skor Risiko Kewilayahan.
- **Expected Business Behavior:**
  1. Dinkes memiliki wewenang pengawasan global tingkat kota (11 Kecamatan, 63 Kelurahan).
  2. Data kasus disajikan secara agregat untuk melindungi privasi warga.
  3. Dinkes dapat melihat indikator performa SLA tiap Puskesmas dan Kader.
- **Actual System Behavior:** Berhasil 100%. API `/gis/heatmap` dan `/analytics/sla` mengembalikan metrik performa dan pemetaan spasial akurat.
- **Status:** `PASS`
- **Evidence:** `gis_handler.go` dan `analytics_handler.go`.

---

## 6. End-to-End Business Scenario (FLOW-01)

Pengujian komprehensif alur hidup penuh (*Full Lifecycle*) 1 kasus DBD dari warga melapor hingga evaluasi pasca-tindakan:

```text
  [Warga] Lapor Kasus di Beji (Gejala DBD, Anak Demam)
     │
     ▼
  [Sistem] Generate Ref No: REP-20260929-BEJI01 ──► [Kader] Notifikasi Masuk
     │                                                    │
     │                                                    ▼
     │                                            [Kader] Klaim Tugas Lapangan
     │                                                    │
     │                                                    ▼
     │                                            [Kader] Investigasi PE di Lokasi
     │                                            (Jentik: Ya, Demam: 3 Org, 3M: Ya)
     │                                                    │
     │                                                    ▼
     │                                            [Sistem] Hitung Rekomendasi:
     │                                            "FOGGING_FOKUS" (Stok Cukup)
     │                                                    │
     │                                                    ▼
  [Puskesmas] Validasi PE ─────────────────────────► Status: WAITING_FOGGING_APPROVAL
     │
     ▼
  [Kepala Puskesmas] Setujui Fogging ──────────────► Status: APPROVED
     │
     ▼
  [Puskesmas/Dinkes] Jadwalkan Tim (+48 Jam) ──────► Status: FOGGING_SCHEDULED
     │
     ▼
  [Tim Fogging] OTW ──► Eksekusi di Lokasi ────────► Status: FOGGING_IN_PROGRESS
     │
     ▼
  [Tim Fogging] Unggah Foto & Input Pemakaian 2L ──► Status: COMPLETED
     │                                                    │
     ▼                                                    ▼
  [Sistem] Potong Stok Gudang Puskesmas (-2L)      [Warga] Terima Notifikasi Selesai
     │                                                    │
     ▼                                                    ▼
  [Sistem] Audit Trail & Heatmap Dinkes Terupdate  [Warga] Kirim Rating Bintang 5
```

### Hasil Uji FLOW-01:
- Seluruh 11 tahapan transisi status berjalan mulus tanpa kegagalan integritas basis data.
- Tidak terjadi *orphan state* atau transaksi gantung.
- Warga berhasil memverifikasi status akhir dan memberikan feedback kepuasan pelayanan.

---

## 7. Inventory Workflow Validation

Pengujian tata kelola rantai pasok logistik insektisida dan larvasida se-Kota Depok:

```text
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    ALUR PERPINDAHAN INVENTARIS LOGISTIK                         │
│                                                                                 │
│   [Gudang Dinkes]                   [In-Transit]              [Gudang Puskesmas]│
│     (Stok: 100)                     (Stok: 0)                   (Stok: 5)       │
│          │                               │                             │        │
│          ▼ [Approve Request]             │                             │        │
│     (Stok: 80) ────────────────────► (Stok: 20)                        │        │
│                                          │                             │        │
│                                          ▼ [Puskesmas Terima]          │        │
│                                      (Stok: 0) ───────────────────► (Stok: 25)  │
│                                                                                 │
│   TOTAL SISTEM: 80 + 0 + 25 = 105 (KONSERVASI STOK TERJAGA 100%)                │
└─────────────────────────────────────────────────────────────────────────────────┘
```

| ID Skenario | Operasi Logistik | Kondisi Awal | Tindakan | Hasil Aktual | Status | Bukti Uji |
|---|---|---|---|---|---|---|
| **I-01** | Permintaan Stok Puskesmas | Stok Puskesmas 2L (Kritis) | Puskesmas ajukan 10L ke Dinkes | Tercatat `StockRequest` status `PENDING` di Dinkes | `PASS` | `inventory_usecase_test.go` |
| **I-02** | Persetujuan & Transfer Dinkes | Dinkes Pusat: 100L | Dinkes setujui permintaan 20L | Dinkes berkurang ke 80L; Terbit `StockTransfer` 20L `IN_TRANSIT` | `PASS` | `TestStockTransfer_AccountingConservation` |
| **I-03** | Penerimaan Barang Puskesmas | Transfer 20L `IN_TRANSIT` | Puskesmas klik "Terima Barang" | Transfer menjadi `COMPLETED`; Gudang Puskesmas bertambah +20L | `PASS` | `TestStockTransfer_Completion` |
| **I-04** | Penolakan Transfer (Barang Rusak) | Transfer 20L `IN_TRANSIT` | Puskesmas klik "Tolak Barang" | Transfer `REJECTED`; Stok Dinkes Pusat di-refund kembali +20L | `PASS` | `TestStockTransfer_RejectionRefund` |
| **I-05** | Operasi Ilegal / Duplikasi | Transfer sudah `COMPLETED` | Mencoba klik terima untuk kedua kalinya | Ditolak: `hanya transfer dengan status IN_TRANSIT yang dapat diproses` | `PASS` | `inventory_usecase.go:262` |

---

## 8. Appeal → PE Validation (Siklus Banding)

Pengujian khusus untuk memastikan bahwa penolakan laporan warga tidak menjadi *dead-end state*:

```text
  Laporan Ditolak (REJECTED_HOAX) 
       │
       ▼ (Warga Banding: 1x Kuota)
  Status: APPEALED
       │
       ▼ (Puskesmas Teliti Berkas Lab RS)
  Puskesmas Menyetujui Banding
       │
       ▼ (Invariant State Machine)
  Status: ESKALASI_PUSKESMAS
       │
       ▼ (Petugas Puskesmas Ambil Alih Tugas)
  Takeover Report (Investigator = Petugas Puskesmas)
       │
       ▼ (Penyelidikan Lapangan oleh Puskesmas)
  Submit Hasil PE (Jentik: Positif, Demam: 3 Orang)
       │
       ▼ (Validasi Medis oleh Kepala Puskesmas - Four Eyes Principle)
  Status: WAITING_FOGGING_APPROVAL
```

- **Temuan UAT:** Seluruh rantai transisi banding berjalan tanpa hambatan (*unblocked*), membuktikan resolusi penuh atas anomali P0-06 sebelumnya.
- **Status:** `PASS`
- **Bukti:** `backend/internal/usecase/bug_reg001_appeal_pe_lifecycle_test.go` (100% Passed).

---

## 9. DSS & Human Decision Validation

| Parameter Pengujian | Logika Sistem | Pengujian Lapangan | Hasil Evaluasi | Status |
|---|---|---|---|---|
| **Kriteria Fogging Fokus** | Larva = true ATAU Demam $\ge 3$, Stok $\ge 10$L | Jentik positif di bak mandi, 3 tetangga demam tinggi | Sistem merekomendasikan `FOGGING_FOKUS` | `PASS` |
| **Kriteria Fogging Tertunda** | Kriteria terpenuhi, namun Stok $< 10$L | Jentik positif, 4 orang demam, stok gudang 4L | Sistem merekomendasikan `FOGGING_TERTUNDA` & menerbitkan draf stok | `PASS` |
| **Kriteria Edukasi PSN** | Larva = false DAN Demam $< 3$ | Tidak ada jentik, 0 penderita demam lain | Sistem merekomendasikan `EDUKASI_PSN` | `PASS` |
| **Mandatory Override Reason** | Override wajib mencantumkan alasan | Mencoba override dari PSN ke Fogging tanpa alasan | Ditolak oleh usecase dengan HTTP 400 | `PASS` |
| **Four-Eyes Principle (SoD)** | Pemeriksa PE $\ne$ Validator | User Sanitarian A memeriksa PE dan memvalidasi sendiri | Ditolak sistem: Pelanggaran pemisahan tugas | `PASS` |

---

## 10. EWS / KLB Validation

Pengujian sistem alarm epidemiologi dini pada tingkat kewilayahan:

### 10.1 Scenario EWS-01: Ambang Batas Kelurahan
- **Skenario:** Kelurahan Sukatani memiliki konfigurasi ambang batas minimal 5 kasus dalam 7 hari. Dimasukkan 5 kasus terverifikasi.
- **Hasil:** Sistem otomatis mendeteksi lonjakan kasus dan menerbitkan status `KlbAlert` aktif untuk Kelurahan Sukatani. Notifikasi peringatan dini terkirim ke dashboard Dinkes.
- **Status:** `PASS`

### 10.2 Scenario EWS-02: Agregasi Tingkat Kecamatan
- **Skenario:** Konfigurasi KLB diubah ke skala `KECAMATAN` dengan batas minimal 6 kasus. Kasus tersebar di beberapa kelurahan (Sukatani: 3 kasus, Sukamaju Baru: 2 kasus, Tapos: 2 kasus $\rightarrow$ Total: 7 kasus di Kecamatan Tapos).
- **Hasil:** Kueri agregasi spasial berhasil menjumlahkan kasus lintas-kelurahan yang berada di bawah naungan Kecamatan Tapos dan memicu peringatan KLB Kecamatan secara akurat.
- **Status:** `PASS`
- **Bukti:** Uji `TestCheckAndTriggerKLB_KecamatanScope` pada `bug08_area_aggregation_test.go`.

### 10.3 Scenario EWS-03: Kasus Non-Fogging Tetap Masuk Surveilans
- **Skenario:** Kasus DBD yang hasil PE-nya merekomendasikan `EDUKASI_PSN` (tidak difogging) diuji apakah tetap dihitung ke dalam ambang batas KLB.
- **Hasil:** Berhasil. Kueri `CountCasesByArea` menghitung seluruh laporan yang berstatus `VERIFIED`, `WAITING_FOGGING_APPROVAL`, `FOGGING_SCHEDULED`, `FOGGING_IN_PROGRESS`, dan `COMPLETED`. Kasus PSN tetap berkontribusi terhadap perhitungan epidemiologi kota.
- **Status:** `PASS`

---

## 11. Negative Business Testing

Pengujian terhadap potensi penyalahgunaan sistem (*abuse*), manipulasi bisnis, dan celah otorisasi:

```text
┌─────────────────────────────────────────────────────────────────────────────────┐
│                      HASIL NEGATIVE BUSINESS TESTING                            │
├─────┬───────────────────────────────────────┬────────────────────┬──────────────┤
│ No  │ Skenario Pengujian Negatif            │ Reaksi Sistem      │ Kesimpulan   │
├─────┼───────────────────────────────────────┼────────────────────┼──────────────┤
│ 1   │ Petugas memvalidasi laporan PE-nya    │ Ditolak UseCase    │ Aman (PASS)  │
│     │ sendiri (Bypass Four-Eyes Principle)  │ (SoD Enforced)     │              │
│ 2   │ Klaim ganda laporan oleh 2 Kader      │ Atomic Lock Aktif  │ Aman (PASS)  │
│     │ secara bersamaan (Race Condition)     │ Hanya 1 pemenang   │              │
│ 3   │ Penyelesaian fogging fiktif tanpa     │ Ditolak Validasi   │ Aman (PASS)  │
│     │ foto bukti atau tanpa pemakaian bahan │ Mandatory Fields   │              │
│ 4   │ Pengambilan alih laporan yang sudah   │ Ditolak State      │ Aman (PASS)  │
│     │ selesai (Takeover on COMPLETED)       │ Machine Guard      │              │
│ 5   │ Pembatalan laporan oleh warga setelah │ Ditolak: Tugas     │ Aman (PASS)  │
│     │ diambil oleh Kader                    │ sedang berjalan    │              │
│ 6   │ Unduh berkas medis/foto laporan warga │ Ditolak HTTP 401   │ Aman (PASS)  │
│     │ oleh pengguna publik/anonim           │ (RBAC Protected)   │              │
│ 7   │ Penerimaan transfer logistik dua kali │ Ditolak: Status    │ Aman (PASS)  │
│     │ (Double Receipt of In-Transit Stock)  │ bukan IN_TRANSIT   │              │
│ 8   │ Pemotongan stok melebihi sisa fisik   │ Ditolak: Saldo     │ Aman (PASS)  │
│     │ yang tersedia di gudang Puskesmas     │ tidak mencukupi    │              │
│ 9   │ Klaim tugas lintas kelurahan via API  │ Diizinkan di API   │ GAP DITEMUKAN│
│     │ langsung (tanpa melalui filter UI)    │ (Kurang guard API) │ (PARTIAL)    │
└─────┴───────────────────────────────────────┴────────────────────┴──────────────┘
```

---

## 12. Business Continuity Validation

Pengujian keandalan operasional saat terjadi gangguan pada infrastruktur pendukung:

### BC-01: Backend Service Temporarily Unavailable
- **Pengamatan:** PWA Frontend menampilkan pesan koneksi terputus yang ramah pengguna (*graceful offline notification*) tanpa menampilkan *blank page* atau *unhandled exception*.
- **Status:** `PASS`

### BC-02: Kegagalan Layanan Notifikasi (FCM / SMTP Down)
- **Pengamatan:** Dalam situasi server SMTP Pemkot Depok atau koneksi FCM mengalami *timeout*, transaksi inti (misal: pengiriman laporan warga atau persetujuan fogging) **tetap berhasil tersimpan** berkat implementasi *defensive nil checks* dan pemrosesan notifikasi *asynchronous non-blocking*.
- **Status:** `PASS`

### BC-03: Pasca-Pemulihan Bencana (*Post-Disaster Recovery Operation*)
- **Pengamatan:** Setelah database dipulihkan dari cadangan fisik melalui `scripts/backup_restore.ps1`, seluruh workflow dapat dilanjutkan secara normal tanpa anomali relasi kunci asing (*foreign keys*) maupun inkonsistensi saldo stok.
- **Status:** `PASS`

### BC-04: Ketidakhadiran Aktor Utama (*Personnel Absence*)
- **Pengamatan:** Apabila Kader Jumantik di suatu RW berhalangan atau tidak aktif dalam 72 jam, sistem memiliki mekanisme *delegated fallback*:
  1. Auto-eskalasi menarik tugas kembali ke Puskesmas.
  2. Petugas Puskesmas dapat melakukan *takeover* langsung.
  3. Dinkes memegang wewenang pengawasan global untuk melakukan intervensi jika Puskesmas kekurangan personel.
- **Status:** `PASS`

---

## 13. Data Meaning Validation

Validasi kesesuaian antara representasi data teknis dengan makna operasional riil di lapangan:

1. **`Status = 'COMPLETED'` pada Citizen Report:**
   - *Makna Bisnis Riil:* Area telah disemprot secara fisik oleh regu fogging, volume insektisida telah dipotong dari pembukuan farmasi, bukti foto telah diunggah, dan warga berhak memberikan rating efektivitas penanganan.
   - *Hasil Evaluasi:* Sesuai 100%. Status `COMPLETED` tidak pernah diberikan jika pemotongan stok belum berhasil dieksekusi.
2. **`Quantity` vs `ReservedQuantity` pada Inventaris:**
   - *Makna Bisnis Riil:* Ketersediaan fisik bahan di gudang Puskesmas dikurangi alokasi yang sudah dipesan untuk jadwal fogging aktif. Mencegah Kepala Puskesmas menyetujui jadwal fogging baru saat stok riil sudah terikat dengan penugasan lain.
   - *Hasil Evaluasi:* Sesuai 100%. Algoritma rekomendasi DSS selalu mengevaluasi `(item.Quantity - item.ReservedQuantity) >= 10.0`.
3. **`Status = 'ESKALASI_PUSKESMAS'`:**
   - *Makna Bisnis Riil:* Kegagalan respon kader tingkat RW dalam batas waktu toleransi 72 jam, atau persetujuan banding atas penolakan keliru. Laporan memerlukan tindakan langsung dari tenaga sanitarian profesional Puskesmas.
   - *Hasil Evaluasi:* Sesuai 100%.
4. **`KlbAlert.IsActive = true`:**
   - *Makna Bisnis Riil:* Ambang batas kasus positif dalam rentang 7 hari terpenuhi di wilayah tersebut dan memerlukan mobilisasi sumber daya pencegahan ekstra dari Dinas Kesehatan.
   - *Hasil Evaluasi:* Sesuai 100%.

---

## 14. UAT Result Matrix

| ID | Skenario Pengujian | Stakeholder | Perilaku yang Diharapkan | Perilaku Aktual Sistem | Status | Bukti / Catatan |
|---|---|---|---|---|---|---|
| **UAT-W01** | Pengiriman Laporan Warga Lengkap | Warga | Data tersimpan, No. Ref terbit, duplikasi terdeteksi | Berhasil tersimpan, ref unik `REP-...`, radius check aktif | `PASS` | `SubmitCitizenReport` (`report_usecase.go:54`) |
| **UAT-W02** | Perlindungan Laporan Anonim | Warga | Identitas pelapor tidak dicatat, tanpa notif PII | Nama `"Anonymous"`, telepon kosong, tidak ada notif | `PASS` | `TestAnonymousReport_EndToEndPIIProtection` |
| **UAT-W03** | Isolasi Data Pelacakan Mandiri | Warga | Lacak via No. Ref aman; via No. HP disamarkan | DTO `ReportSummaryResponse` memproteksi PII & foto medis | `PASS` | `TrackReportsByPhone` (`report_usecase.go:165`) |
| **UAT-W04** | Pengajuan Banding Penolakan | Warga | Laporan `REJECTED_HOAX` dapat dibanding 1x | Status `APPEALED`, kuota banding terkunci atomik | `PASS` | `SubmitAppealAtomic` (`report_repository.go:95`) |
| **UAT-W05** | Pembatalan Laporan Mandiri | Warga | Batal berhasil jika belum diklaim kader | Berhasil dibatalkan; ditolak jika sudah diambil kader | `PASS` | `CancelReport` (`report_usecase.go:290`) |
| **UAT-K01** | Klaim Tugas Investigasi Lapangan | Kader | Satu tugas hanya boleh diambil satu kader | Race condition tertangani; atomic lock aktif | `PASS` | `TestBUGNEW008_ClaimReport_AtomicRaceCondition` |
| **UAT-K02** | Validasi Yurisdiksi Wilayah Klaim | Kader | Kader hanya bisa klaim laporan di kelurahannya | Filter UI aktif; namun API belum memvalidasi LocationID | `PARTIAL` | `ClaimReport` (`report_usecase.go:273`) |
| **UAT-K03** | Penolakan Laporan Hoax | Kader | Laporan fiktif ditolak dengan alasan wajib | Status `REJECTED_HOAX`, alasan tersimpan di audit | `PASS` | `VerifyReport` (`verification_usecase.go:56`) |
| **UAT-K04** | Input Hasil PE Lapangan | Kader | Wajib edukasi 3M Plus, hitung jentik & demam | Tersimpan di `pe_results`, status `MENUNGGU_VALIDASI` | `PASS` | `SubmitPEResult` (`pe_usecase.go:41`) |
| **UAT-K05** | Eskalasi SLA 72 Jam | Kader | Laporan >72 jam otomatis ditarik ke Puskesmas | Status `ESKALASI_PUSKESMAS`, investigator di-unassign | `PASS` | `TestAutoReleaseExpiredTasks_72hSLA` |
| **UAT-P01** | Takeover Laporan Tereskalasi | Puskesmas | Puskesmas mengambil alih penanganan PE | Pengambilalihan sukses, audit merekam mutasi | `PASS` | `TestTakeoverWorkflow` |
| **UAT-P02** | Penegakan Four-Eyes Principle | Puskesmas | Pemeriksa PE dilarang memvalidasi sendiri | Ditolak usecase jika OfficerID == ValidatorID | `PASS` | `ValidatePEResult` (`pe_usecase.go:146`) |
| **UAT-P03** | Rekomendasi Penasihat DSS | Puskesmas | Rekomendasi deterministik berbasis bukti & stok | Sesuai matriks PRD (Fogging/Tertunda/PSN) | `PASS` | `pe_usecase.go:69-85` |
| **UAT-P04** | Override Rekomendasi DSS | Puskesmas | Override diizinkan dengan alasan wajib & audit | Berhasil, audit mencatat rekomendasi vs keputusan | `PASS` | `TestP0_03_DSS_Override_WithReason_Success` |
| **UAT-P05** | Persetujuan Jadwal Fogging | Puskesmas | Kepala Puskesmas menyetujui jadwal $\ge 24$ jam | Jadwal tervalidasi, ditolak jika $< 24$ jam dari approve | `PASS` | `fogging_schedule_usecase.go:76` |
| **UAT-F01** | Pembatasan Kapasitas Harian Tim | Tim Fogging | 1 tim maksimal 1 penugasan fogging per hari | Terkunci di usecase `CountSchedulesByTeamAndDate` | `PASS` | `fogging_schedule_usecase.go:90` |
| **UAT-F02** | Eksekusi & Pemotongan Stok | Tim Fogging | Status `COMPLETED` otomatis potong stok gudang | Pemotongan atomik, reservasi dilepas, status selesai | `PASS` | `fogging_schedule_usecase.go:193` |
| **UAT-F03** | Kewajiban Bukti Foto Eksekusi | Tim Fogging | Tidak bisa selesai tanpa foto & rincian bahan | Ditolak backend & form terkunci di frontend | `PASS` | `fogging/schedules/page.tsx:80` |
| **UAT-I01** | Permintaan Stok Puskesmas | Farmasi | Pengajuan stok tercatat di dashboard Dinkes | Tercatat `StockRequest` status `PENDING` | `PASS` | `inventory_usecase.go:113` |
| **UAT-I02** | Konservasi Stok Logistik Transfer | Farmasi | Total stok sistem konstan selama pengiriman | Sumber berkurang, transit bertambah, tujuan tetap | `PASS` | `TestStockTransfer_AccountingConservation` |
| **UAT-I03** | Penerimaan & Penyelesaian Transfer | Farmasi | Barang diterima menambah stok gudang tujuan | Status `COMPLETED`, stok tujuan bertambah fisik | `PASS` | `inventory_usecase.go:281` |
| **UAT-I04** | Pembatalan / Retur Transfer | Farmasi | Barang ditolak dikembalikan ke stok pengirim | Status `REJECTED`, refund penuh ke stok asal | `PASS` | `TestStockTransfer_RejectionRefund` |
| **UAT-D01** | Ambang Batas KLB Kelurahan | Dinkes | Alarm otomatis aktif jika kasus $\ge$ threshold | `KlbAlert` aktif otomatis, notifikasi terbit | `PASS` | `CheckAndTriggerKLB` (`ews_klb_usecase.go:56`) |
| **UAT-D02** | Agregasi Ambang Batas Kecamatan | Dinkes | Akumulasi kasus lintas kelurahan dalam 1 camat | Teragregasi via PostGIS, alarm kecamatan aktif | `PASS` | `TestCheckAndTriggerKLB_KecamatanScope` |
| **UAT-D03** | Pemantauan Spasial Heatmap | Dinkes | Pemetaan sebaran kasus terpusat se-Depok | Endpoint `/gis/heatmap` menyajikan agregasi aman | `PASS` | `gis_handler.go:20` |

---

## 15. Business & Operational Findings

### Finding UAT-FIND-001 (P1 — High Impact): Ketiadaan Validasi Yurisdiksi Wilayah pada Endpoint Klaim Kader
- **Skenario:** Kader dari Kelurahan A (misal: Beji) menembak API `/api/v1/kader/reports/:id/claim` secara langsung dengan menyertakan UUID laporan dari Kelurahan B (misal: Pancoran Mas).
- **Perilaku yang Diharapkan:** Sistem menolak klaim tugas dengan pesan `akses ditolak: laporan berada di luar wilayah yurisdiksi kerja Anda`.
- **Perilaku Aktual:** Sistem menerima klaim tersebut karena `ClaimReport` di `report_usecase.go` hanya memeriksa apakah status laporan `MENUNGGU_PE_KADER` dan `InvestigatorID == nil`, tanpa membandingkan `user.LocationID` dengan `report.Location.PuskesmasID`.
- **First Divergence:** `backend/internal/usecase/report_usecase.go:273` (`ClaimReport`).
- **Root Cause:** Pengecekan wilayah baru diimplementasikan pada kueri pembacaan antarmuka (`GetReports`), tetapi belum ditanamkan pada logika mutasi klaim tugas.
- **Dampak Bisnis:** Kader nakal atau keliru dapat mengklaim tugas di luar wilayah kerjanya, berpotensi memicu konflik operasional antar-kader atau laporan terbengkalai karena jarak tempuh lapangan yang tidak realistis.
- **Affected Role / Workflow:** Kader Jumantik / Alur Verifikasi & Klaim Tugas.
- **Klasifikasi:** `Business Logic Bug` / `Governance Gap`.
- **Rekomendasi:** Tambahkan verifikasi `userLocationID` pada `ClaimReport` di usecase sebelum memanggil repository atomik.
- **Needs Code Change?:** `YES` (Pasca-UAT).

---

### Finding UAT-FIND-002 (P2 — Medium Impact): Auto-Draft Stock Request Langsung Masuk Antrean Dinkes Tanpa Verifikasi Kepala Puskesmas
- **Skenario:** Validasi PE menghasilkan rekomendasi `FOGGING_TERTUNDA` akibat stok insektisida Puskesmas $< 10$ Liter.
- **Perilaku yang Diharapkan:** Sistem membuat draf permintaan stok (`StockRequest`) dengan status `DRAFT` yang memerlukan peninjauan dan tanda tangan/persetujuan Kepala Puskesmas sebelum dikirimkan ke Dinas Kesehatan.
- **Perilaku Aktual:** Sistem langsung membuat `StockRequest` dengan status `PENDING` (`pe_usecase.go:267`), sehingga pengajuan logistik langsung muncul di dashboard Kepala Dinas Kesehatan tanpa persetujuan Kepala Puskesmas.
- **First Divergence:** `backend/internal/usecase/pe_usecase.go:267`.
- **Root Cause:** Alur dibuat otomatis penuh (*automated shortcut*) untuk efisiensi, namun menabrak hierarki birokrasi pengadaan obat farmasi pemerintah.
- **Dampak Bisnis:** Kepala Puskesmas kehilangan kontrol anggaran dan visibilitas atas permintaan obat yang diajukan atas nama faskesnya ke Dinas Kesehatan.
- **Affected Role / Workflow:** Kepala Puskesmas & Pengelola Logistik Dinkes / Alur Permintaan Stok.
- **Klasifikasi:** `Governance Gap` / `Workflow Gap`.
- **Rekomendasi:** Ubah status default draf menjadi `DRAFT_PUSKESMAS` dan sediakan tombol *"Kirim ke Dinkes"* di dashboard Kepala Puskesmas.
- **Needs Code Change?:** `BUSINESS DECISION` (Tergantung apakah Dinkes menghendaki otomatisasi penuh atau persetujuan berjenjang).

---

### Finding UAT-FIND-003 (P2 — Medium Impact): Kebijakan Kapasitas 1 Tugas/Hari/Tim Membatasi Penanganan Saat Wabah / KLB
- **Skenario:** Terjadi lonjakan kasus DBD (outbreak) di 3 RW berdekatan pada satu kelurahan dalam minggu yang sama, sementara Puskesmas hanya memiliki 1 regu fogging.
- **Perilaku yang Diharapkan:** Sistem mengizinkan penjadwalan darurat untuk menyemprot area bersebelahan pada hari yang sama dengan persetujuan khusus Kepala Puskesmas atau Dinkes.
- **Perilaku Aktual:** Usecase `fogging_schedule_usecase.go:94` mengunci mati kapasitas: `if count >= 1 { return errors.New("tim fogging sudah memiliki jadwal pada hari tersebut (maksimal 1 tugas per hari)") }`.
- **First Divergence:** `backend/internal/usecase/fogging_schedule_usecase.go:94`.
- **Root Cause:** Penerapan kaku aturan bisnis PRD 5.2 tanpa klausul pengecualian kondisi darurat (*Emergency Override*).
- **Dampak Bisnis:** Penanganan kasus kedua dan ketiga terpaksa ditunda 24 hingga 48 jam kemudian, meningkatkan risiko penularan lokal nyamuk *Aedes aegypti*.
- **Affected Role / Workflow:** Tim Fogging, Kepala Puskesmas / Alur Penjadwalan Eksekusi.
- **Klasifikasi:** `Operational Gap` / `Business Decision Required`.
- **Rekomendasi:** Tambahkan flag *KLB Emergency Multi-Tasking* yang dapat diaktifkan oleh Kepala Puskesmas atau Dinkes untuk mengizinkan maksimal 2–3 tugas per hari khusus saat status KLB aktif di area tersebut.
- **Needs Code Change?:** `BUSINESS DECISION`.

---

### Finding UAT-FIND-004 (P2 — Medium Impact): Batas Waktu Banding Warga (48 Jam) Belum Ditegakkan Secara Temporal di Database
- **Skenario:** Laporan warga ditolak sebagai hoax pada 3 minggu yang lalu. Warga membuka tautan pelacakan hari ini dan mengklik tombol ajukan banding.
- **Perilaku yang Diharapkan:** Banding ditolak karena telah melewati batas kedaluwarsa 48 jam sejak penolakan diterbitkan (sesuai SOP & PRD).
- **Perilaku Aktual:** Banding berhasil diterima (`APPEALED`) karena kueri SQL `SubmitAppealAtomic` hanya memeriksa `appeal_count < 1` dan status penolakan, tanpa memeriksa selisih waktu `updated_at`.
- **First Divergence:** `backend/internal/repository/report_repository.go:97`.
- **Root Cause:** Pengecekan temporal 48 jam tertulis dalam dokumen SOP tetapi terlewat dalam klausul `WHERE` kueri update SQL.
- **Dampak Bisnis:** Kasus-kasus lama yang sudah kedaluwarsa dapat dihidupkan kembali sewaktu-waktu oleh warga, mengacaukan antrean investigasi aktif Puskesmas.
- **Affected Role / Workflow:** Warga & Puskesmas / Alur Banding.
- **Klasifikasi:** `Requirement Gap` / `Business Logic Bug`.
- **Rekomendasi:** Tambahkan klausul `AND updated_at >= NOW() - INTERVAL '48 hours'` pada fungsi `SubmitAppealAtomic`.
- **Needs Code Change?:** `YES` (Pasca-UAT).

---

### Finding UAT-FIND-005 (P3 — Low Impact): Ketergantungan Pelacakan Manual untuk Evaluasi Warga Non-PWA
- **Skenario:** Warga melapor tanpa akun melalui web publik dan tidak menginstal PWA. Hari ke-14 pasca-fogging, cron pengingat feedback berjalan.
- **Perilaku yang Diharapkan:** Warga menerima pengingat untuk mengisi ulasan efektivitas fogging.
- **Perilaku Aktual:** Cron menandai `feedback_requested = true`, namun karena gateway SMS/WhatsApp dilarang oleh batasan biaya open-source (PRD Bab 4), warga hanya bisa mengisi jika secara sadar berinisiatif membuka kembali tautan `/track`.
- **First Divergence:** `backend/internal/usecase/report_usecase.go:439`.
- **Root Cause:** Batasan PRD 100% Gratis/Open-Source yang meniadakan SMS/WA Gateway berbayar.
- **Dampak Bisnis:** Tingkat pengisian evaluasi efektivitas fogging (*citizen feedback rate*) akan sangat rendah pada populasi warga umum.
- **Affected Role / Workflow:** Warga / Alur Evaluasi Pasca-Fogging.
- **Klasifikasi:** `Operational Gap` / `Business Constraint`.
- **Rekomendasi:** Optimalkan pengiriman email gratis (SMTP Pemkot Depok) bagi warga yang menyertakan email non-anonim, serta sediakan QR code pada stiker fisik yang ditempelkan tim fogging di rumah warga saat penyemprotan.
- **Needs Code Change?:** `NO` (Mitigasi operasional SOP lapangan).

---

## 16. Risk & Control Implications

Evaluasi terhadap dampak temuan UAT terhadap matriks manajemen risiko:

| Kode Risiko | Judul Risiko | Level Inheren | Kontrol Saat Ini | Status Kontrol Pasca-UAT | Efektivitas Kontrol |
|---|---|---|---|---|---|
| **RSK-OP-01** | Kebocoran Wilayah Klaim Kader | Sedang | Filter kueri `GetReports` di UI | Celah di API langsung (`UAT-FIND-001`) | **Perlu Peningkatan (Need API Guard)** |
| **RSK-OP-02** | Pengadaan Stok Bypassing Manajerial | Sedang | Otomatisasi `FOGGING_TERTUNDA` | Draf langsung masuk antrean Dinkes (`UAT-FIND-002`) | **Perlu Keputusan Bisnis Dinkes** |
| **RSK-OP-03** | Penumpukan Kasus Terhambat Kapasitas | Tinggi | Batasan 1 tugas/hari | Kaku saat kondisi KLB (`UAT-FIND-003`) | **Perlu Emergency Exception Policy** |
| **RSK-OP-04** | Kebocoran Dokumen Sensitif Warga | Kritis | Middleware Protected RBAC & PostGIS | Terbukti kokoh di seluruh endpoint foto | **Efektif (100% Mitigated)** |
| **RSK-OP-05** | Penyimpangan Akuntansi Logistik | Kritis | Single Accounting & In-Transit Lock | Konservasi stok terbukti absolut | **Efektif (100% Mitigated)** |
| **RSK-OP-06** | Keputusan Klinis Digantikan Robot | Kritis | Advisory DSS & Mandatory Override | Kepala Puskesmas memegang wewenang mutlak | **Efektif (100% Mitigated)** |

---

## 17. Open Business Decisions (Keputusan Bisnis yang Dibutuhkan)

Sebelum rilis operasional penuh ke 38 Puskesmas se-Kota Depok, pihak pimpinan Dinas Kesehatan dan Tim Manajemen DBD Guard perlu mengonfirmasi 3 butir keputusan:

1. **Kebijakan Draf Permintaan Stok Otomatis:**
   - *Opsi A (Status Quo):* Sistem langsung mengirimkan permintaan stok ke Dinkes saat rekomendasi `FOGGING_TERTUNDA` terbit untuk mempercepat pengiriman logistik.
   - *Opsi B (Birokratis Formal - Direkomendasikan):* Draf permintaan stok tersimpan di Puskesmas dan wajib disetujui secara eksplisit oleh Kepala Puskesmas sebelum masuk ke radar Dinas Kesehatan.
2. **Kapasitas Regu Fogging saat Kondisi Darurat (KLB):**
   - *Opsi A:* Mempertahankan batas kaku 1 tugas/hari per tim untuk menjaga keselamatan personel dari paparan asap insektisida berlebih.
   - *Opsi B (Direkomendasikan):* Mengizinkan *Emergency Surge* hingga 2–3 tugas per hari jika status KLB aktif di wilayah tersebut, didukung suplai BKO dari Tim Reaksi Cepat Dinkes Kota Depok.
3. **Masa Kedaluwarsa Banding Laporan Hoax:**
   - Menetapkan secara baku di dalam regulasi apakah hak banding warga hangus setelah **48 Jam** (2 hari) atau diperpanjang hingga **72 Jam** (3 hari) sejak diterbitkannya penolakan oleh Kader.

---

## 18. Final UAT Verdict

Berdasarkan metodologi pengujian berjenjang, verifikasi peran stakeholder, pengujian skenario lintas-peran (*End-to-End*), pengujian integritas akuntansi inventaris, serta validasi operasional lapangan:

### **VERDICT: CONDITIONALLY BUSINESS READY**

### Justifikasi Putusan:
1. **Kesiapan Inti Bisnis (*Core Business Readiness*):** Seluruh alur kerja kritis kesehatan masyarakat (Pelaporan $\rightarrow$ Verifikasi $\rightarrow$ PE $\rightarrow$ Rekomendasi Penasihat DSS $\rightarrow$ Persetujuan Kepala Puskesmas $\rightarrow$ Eksekusi Fogging Lapangan $\rightarrow$ Pemotongan Stok Farmasi $\rightarrow$ Deteksi KLB) telah teruji secara komprehensif dan **beroperasi secara valid dan akuntabel**.
2. **Tidak Ada Bug P0 Terbuka:** Tidak ada anomali fatal, *deadlock state*, kebocoran data privasi warga, atau ketidakseimbangan akuntansi inventaris yang tersisa pada sistem.
3. **Syarat Menuju Full Operational Release:**
   - Menerapkan perbaikan minor (P1) pada penegakan batas wilayah di endpoint klaim kader (`UAT-FIND-001`).
   - Menyepakati 3 butir keputusan bisnis terbuka bersama Kepala Dinas Kesehatan Kota Depok terkait draf stok otomatis, kapasitas darurat KLB, dan jendela waktu banding warga.

Sistem DBD Guard Kota Depok secara substansi telah **siap untuk memasuki tahap uji coba operasional lapangan terbatas (*Pilot Run / Soft Launch*)** di wilayah kerja percontohan sebelum diimplementasikan secara menyeluruh di 11 Kecamatan Kota Depok.
