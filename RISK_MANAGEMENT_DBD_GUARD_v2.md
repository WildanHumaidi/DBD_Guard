# Risk Management — DBD Guard

## Post-Implementation Revision (Version 5.0)

**Proyek:** DBD Guard (Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi Kota Depok)  
**Versi Dokumen:** 5.0 — Pembaruan Pasca-Implementasi P0, Audit UAT, dan IMP-001  
**Tanggal Evaluasi:** 01 Oktober 2026  
**Status Audit:** **Updated & Reconciled (Post-Implementation Baseline)**  
**Klasifikasi:** Internal — Laporan Manajemen Risiko & Kesiapan Rilis Operasional  
**Prinsip Rekayasa Perangkat Lunak:** **ZERO CODE MODIFICATION** (Audit & Dokumentasi Kepatuhan)  

---

## 1. Document Information

### 1.1 Riwayat Revisi Dokumen

| Versi | Tanggal | Penulis / Evaluator | Deskripsi Perubahan Utama | Status |
|:---:|:---:|:---|:---|:---:|
| **v1.0** | 15 September 2026 | Tim Arsitektur Sistem | Pemetaan risiko awal berdasarkan draf PRD Komprehensif. | Archived |
| **v2.0** | 22 September 2026 | Tim Pengembang Backend | Pembaruan risiko pasca Deep Bug Audit Fase 1. | Archived |
| **v3.0** | 25 September 2026 | Tim QA & Keamanan Data | Penambahan risiko kebocoran PII dan pelanggaran UU PDP. | Archived |
| **v4.0** | 28 September 2026 | Lead Evaluator Independen | Audit menyeluruh bisnis, operasional, dan tata kelola faskes; penetapan status *Conditional Hold*. | Superceded |
| **v5.0** | **01 Oktober 2026** | **AI Agentic System Evaluator** | **Revisi menyeluruh pasca perbaikan P0 Fundamental Fix, DR Drill, UAT Reconciliation, Governance Resolution, dan Implementasi IMP-001 (Kader Jurisdiction Guard).** | **CURRENT** |

### 1.2 Tujuan Dokumen
Dokumen **Risk Management v5.0** ini disusun untuk:
1. Menyelaraskan status risiko dengan kondisi aktual kode sumber, arsitektur, dan operasional terbaru proyek DBD Guard Kota Depok.
2. Memetakan transisi kontrol dari status lama (*Inherent / Old Control Gap*) ke status kontrol yang telah terpasang (*Implemented Controls*), serta mengukur efektivitas kontrol secara empiris.
3. Mendokumentasikan penyelesaian cacat keamanan **UAT-FIND-001 / IMP-001** (*Kader Cross-Region Task Claim Authorization*) sebagai kontrol backend baru.
4. Menegaskan risiko residual (*Residual Risks*) yang masih ada dan mencatat status tiga agenda kebijakan birokrasi (**DEC-01, DEC-02, DEC-03**) yang tetap berstatus **PENDING STAKEHOLDER DECISION**.
5. Menjamin keterlacakan historis (*Historical Traceability*) tanpa menghapus temuan audit masa lalu.

---

## 2. Revision Summary

Pembaruan dari Dokumen v4.0 menuju v5.0 didasarkan pada tonggak pencapaian teknis dan tata kelola yang telah berhasil diverifikasi:

```text
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                             EVOLUSI KONTROL SISTEM DBD GUARD (v4.0 ──► v5.0)                     │
├──────────────────────────┬─────────────────────────────┬─────────────────────────────────────────┤
│ Domain / Area Masalah    │ Status Awal (v4.0 - Audit)  │ Status Terkini (v5.0 - Post-Impl)       │
├──────────────────────────┼─────────────────────────────┼─────────────────────────────────────────┤
│ 1. Akuntansi Logistik    │ Auto-deduct gudang donor    │ P0-01 Resolved: Strict StockTransfer    │
│    (RSK-01 / BW-01)      │ tanpa BAST / transfer fisik │ atomik, zero virtual stock, verified    │
├──────────────────────────┼─────────────────────────────┼─────────────────────────────────────────┤
│ 2. Diskresi Medis DSS    │ Auto-close tanpa fogging    │ P0-03 Resolved: Rekomendasi advisory,   │
│    (RSK-02 / BW-03)      │ membendung wewenang dokter  │ Clinical Override Kepala Puskesmas OK   │
├──────────────────────────┼─────────────────────────────┼─────────────────────────────────────────┤
│ 3. Privasi Berkas Medis  │ app.Static("/storage")      │ P0-02 Resolved: Endpoint terproteksi    │
│    (RSK-03 / BW-08)      │ foto rujukan terbuka bebas  │ JWT + Location Scope VerifyLocation     │
├──────────────────────────┼─────────────────────────────┼─────────────────────────────────────────┤
│ 4. Segregation of Duties │ Takeover investigator merang│ P0-04 Resolved: OfficerID != ValidatorID│
│    (RSK-08 / BW-04)      │ kap sebagai validator PE    │ ditegakkan di usecase; error 403 active │
├──────────────────────────┼─────────────────────────────┼─────────────────────────────────────────┤
│ 5. Disaster Recovery     │ Single container PostgreSQL │ P0-05 Resolved: Live backup/restore     │
│    (RSK-09 / BCR)        │ tanpa script cadangan/drill │ drill match 100% (checksum & row match) │
├──────────────────────────┼─────────────────────────────┼─────────────────────────────────────────┤
│ 6. Banding Laporan Hoax  │ Alur terputus, siklus mati  │ P0-06 Resolved: REJECTED_HOAX -> APPEAL │
│    (BUG-01 / P0-06)      │ pasca penolakan kader       │ -> ESKALASI -> TAKEOVER -> PE verified  │
├──────────────────────────┼─────────────────────────────┼─────────────────────────────────────────┤
│ 7. Otorisasi Klaim Kader │ UI filtering only; direct   │ IMP-001 Resolved: Backend Jurisdiction  │
│    (UAT-FIND-001 / IMP)  │ API bypass lintas faskes    │ Authorization Guard, 403 Forbidden OK   │
├──────────────────────────┼─────────────────────────────┼─────────────────────────────────────────┤
│ 8. Keputusan Kebijakan   │ Belum diidentifikasi        │ Formal Governance Proposal diserahkan   │
│    (DEC 01, 02, 03)      │ sebagai isu birokrasi       │ Status: PENDING STAKEHOLDER DECISION    │
└──────────────────────────┴─────────────────────────────┴─────────────────────────────────────────┘
```

---

## 3. Current Project Risk Posture

Pasca-eksekusi remedi P0 dan implementasi bedah IMP-001, profil risiko keseluruhan proyek DBD Guard mengalami transformasi signifikan:

1. **Integritas Teknis & Keamanan Aplikasi (Technical & Security Layer):**
   - Seluruh celah berkategori **Kritis / Fatal (P0)** pada lapisan backend, database, dan otorisasi telah **100% DIMERDEKAKAN (RESOLVED)**.
   - Tidak ada lagi kebocoran stok obat, tidak ada lagi kebocoran berkas medis pelapor, dan tidak ada lagi pengambilalihan tugas investigasi lintas wilayah kerja faskes.
   - Seluruh pengujian otomatis backend (`go test ./...`) dan build frontend Next.js (`npm run build`) berstatus **PASS 100% (Exit 0)**.

2. **Kondisi Risiko Residual (Residual Risk Profile):**
   - Risiko sistemik telah bergeser dari ranah *perangkat lunak (software flaws)* ke ranah **tata kelola birokrasi, kepatuhan kelembagaan, dan adopsi operasional lapangan**.
   - Sistem **TIDAK LAGI BERADA PADA STATUS CONDITIONAL HOLD**, melainkan berstatus **TECHNICAL PASS / GATED ON GOVERNANCE SIGN-OFF**.

---

## 4. Risk Assessment Method

Metode penilaian risiko tetap konsisten mengacu pada ISO 31000 dan Pedoman Manajemen Risiko SPBE KemenPAN-RB:

* **Skala Kemungkinan (Likelihood / L: 1 - 5):**
  $1 = \text{Sangat Rendah} (<5\%/\text{thn}); 2 = \text{Rendah} (5-20\%/\text{thn}); 3 = \text{Sedang} (21-50\%/\text{thn}); 4 = \text{Tinggi} (51-80\%/\text{thn}); 5 = \text{Sangat Tinggi} (>80\%/\text{thn})$.
* **Skala Dampak (Impact / I: 1 - 5):**
  $1 = \text{Insignifikan}; 2 = \text{Minor}; 3 = \text{Moderat}; 4 = \text{Mayor / Hukum}; 5 = \text{Bencana / Fatal / Kematian}$.
* **Matriks Tingkat Risiko ($L \times I$):**
  - **Rendah (Low / 1 - 5):** Risiko dapat diterima dengan pemantauan berkala.
  - **Sedang (Medium / 6 - 11):** Mitigasi operasional terencana.
  - **Tinggi (High / 12 - 19):** Intervensi pengawasan manajemen faskes.
  - **Kritis (Critical / 20 - 25):** Ancaman fatal integritas sistemik.

---

## 5. Updated Risk Register

Berikut adalah Register Risiko Komprehensif DBD Guard yang telah direvaluasi pasca-implementasi. Nomor ID risiko historis (**RSK-01 s/d RSK-18**) dipertahankan secara stabil, ditambah item **RSK-19** (terkait temuan UAT-FIND-001):

| ID | Kat | Risiko Spesifik | Akar Masalah (Root Cause) | Inherent Risk (L x I) | Kontrol Terbaru yang Terpasang (Current Control) | Efektivitas Kontrol | Residual Risk (L x I) | Status Penanganan | Pemilik Risiko | Bukti Verifikasi | Status Risiko |
|:---|:---:|:---|:---|:---:|:---|:---:|:---:|:---|:---|:---:|:---:|
| **RSK-01** | FAR | **Manipulasi Stok Melalui Alokasi Silang Otomatis** | Auto-deduct donor warehouse tanpa BAST fisik saat stok faskes lokal minus. | **20 (C)**<br>$(4 \times 5)$ | P0-01 Fix: Fungsi auto-deduct dihapus total. Mutasi hanya via modul `StockTransfer` resmi (IN_TRANSIT -> COMPLETED/REJECTED) dengan pembukuan ganda atomik. | **Effective** | **4 (L)**<br>$(1 \times 4)$ | Completed | Pengelola Barang Dinkes / Puskesmas | `TestP0_01_InventoryAccountingIntegrity` (PASS) | **MITIGATED** |
| **RSK-02** | DAR | **Pemberangusan Diskresi Medis oleh Algoritma DSS** | Hasil PE non-jentik otomatis mengunci status menjadi `CLOSED_NO_FOGGING`. | **16 (H)**<br>$(4 \times 4)$ | P0-03 Fix: Status otomatis dicabut. Seluruh PE masuk antrean review dokter/Kepala Puskesmas dengan fitur *Clinical Override* wajib isi alasan medis. | **Effective** | **4 (L)**<br>$(1 \times 4)$ | Completed | Kepala Puskesmas / Surveilans Dinkes | `TestP0_03_and_04_DSS_and_SoD` (PASS) | **MITIGATED** |
| **RSK-03** | SCR | **Ekspos Terbuka Foto Rekam Medis & PII Pelapor** | Endpoint `app.Static("/storage")` tanpa autentikasi menyajikan foto surat rawat/KTP warga. | **25 (C)**<br>$(5 \times 5)$ | P0-02 Fix: Static folder dihapus. Dibuat endpoint `/storage/photos/:filename` terproteksi JWT, role check, dan location check (`VerifyReportLocation`). | **Effective** | **6 (M)**<br>$(2 \times 3)$ | Completed | Technical Owner / Data Protection Officer | `TestP0_02_SensitiveFileStorageProtection` (PASS) | **MITIGATED** |
| **RSK-04** | FAR | **Moral Hazard Input Data PE Lapangan** | Input variabel jentik & demam deklaratif teks tanpa validasi GPS riil. | **16 (H)**<br>$(4 \times 4)$ | Invariant kepemilikan PE (`UNAUTHORIZED_PE_SUBMISSION`), audit log investigator, dan mekanisme banding warga teruji. Geofencing GPS disiapkan di mobile. | **Partially Effective** | **9 (M)**<br>$(3 \times 3)$ | In Progress | Seksi P2P Dinkes / Sanitarian Puskesmas | `TestP0_06_AppealToPELifecycle` (PASS) | **PARTIALLY MITIGATED** |
| **RSK-05** | OER | **Kelumpuhan Respons Wabah Akibat Anti-Overbooking Kaku** | Batas kaku 1 tim fogging maksimal 1 tugas/hari tanpa klausul darurat KLB. | **20 (C)**<br>$(4 \times 5)$ | Validasi kuota harian berjalan andal untuk kondisi normal. Pengecualian darurat saat KLB menunggu pengesahan kebijakan **DEC-02**. | **Effective (Normal); Ineffective (KLB)** | **12 (H)**<br>$(3 \times 4)$ | Pending Decision | Bidang P2P Dinkes / Koordinator Fogging | `BUSINESS_DECISION_GOVERNANCE_DBD_GUARD.md` | **PENDING STAKEHOLDER DECISION** |
| **RSK-06** | DAR | **Blindspot Early Warning System (EWS) KLB** | Kueri EWS hanya menghitung kasus bertindakan fogging, mengecualikan kasus PSN/abatisasi. | **20 (C)**<br>$(4 \times 5)$ | Kueri `CountCasesByArea` diperbaiki menjaring seluruh kasus terkonfirmasi; evaluasi threshold paralel Kelurahan & Kecamatan. | **Effective** | **6 (M)**<br>$(2 \times 3)$ | Completed | Epidemiolog Kesehatan Dinkes | `TestCountCasesByArea_KelurahanAndKecamatanScope` (PASS) | **MITIGATED** |
| **RSK-07** | BGR | **Resistensi Pengguna Lapangan & Pengabaian Sistem** | Ketiadaan payung hukum formal (SK/Perwal) dan minim insentif kuota relawan. | **16 (H)**<br>$(4 \times 4)$ | Naskah akademis SOP telah dirumuskan; advokasi SK Kepala Dinas Kesehatan sedang diajukan menjelang peluncuran pilot. | **Not Verified** | **12 (H)**<br>$(3 \times 4)$ | In Progress | Kepala Dinas Kesehatan Kota Depok | Dokumen Tata Kelola SPBE | **OPEN GOVERNANCE ACTION** |
| **RSK-08** | BGR | **Pelanggaran Segregation of Duties pada Takeover** | Petugas Puskesmas pengambil alih laporan dapat memvalidasi PE buatannya sendiri. | **12 (H)**<br>$(3 \times 4)$ | P0-04 Fix: Validasi programatik `OfficerID != ValidatorID` aktif pada `ValidatePEResult`; penolakan 403 Forbidden jika melanggar Four-Eyes Principle. | **Effective** | **3 (L)**<br>$(1 \times 3)$ | Completed | Koordinator Tata Kelola IT Dinkes | `TestP0_03_and_04_DSS_and_SoD` (PASS) | **RESOLVED** |
| **RSK-09** | TIR | **Titik Kegagalan Tunggal Basis Data & Ketiadaan DR** | Satu instans PostgreSQL tanpa pencadangan berkala otomatis dan prosedur restorasi. | **15 (H)**<br>$(3 \times 5)$ | P0-05 / DR Drill: Skrip `backup_restore.ps1` teruji live; restorasi 100% match data (zero data loss); prosedur RPO/RTO terdokumentasi. | **Effective** | **6 (M)**<br>$(2 \times 3)$ | Completed | Administrator Infrastruktur Diskominfo | Live Restore Drill Log (Exit 0) | **CONTROL IMPROVED** |
| **RSK-10** | SCR | **Pelanggaran Privasi melalui Endpoint Pelacakan Kontak** | Endpoint publik bebas menyajikan seluruh riwayat kasus hanya berbekal nomor ponsel. | **16 (H)**<br>$(4 \times 4)$ | DTO summary diterapkan; pelaporan anonim melindungi kontak (`TestAnonymousReport_EndToEndPIIProtection`); rate limiter aktif. | **Partially Effective** | **6 (M)**<br>$(2 \times 3)$ | Completed | Technical Owner / Data Protection Officer | `TestAnonymousReport_EndToEndPIIProtection` (PASS) | **CONTROL IMPROVED** |
| **RSK-11** | FAR | **Penyelesaian Fogging Fiktif Tanpa Verifikasi Foto** | Status jadwal dapat diubah ke `COMPLETED` tanpa mengecek bukti foto fisik. | **16 (H)**<br>$(4 \times 4)$ | P1-1 Fix: Pengecekan keberadaan minimal 1 rekaman foto di tabel `fogging_evidences` sebelum eksekusi status selesai. | **Effective** | **4 (L)**<br>$(1 \times 4)$ | Completed | Seksi Pengendalian Vektor Dinkes | Code Review `schedule_handler.go` | **MITIGATED** |
| **RSK-12** | OER | **Monopoli Penugasan oleh Kader (Task Hoarding)** | Ketiadaan batas kuota klaim serentak per akun Kader memicu penelantaran tugas. | **12 (H)**<br>$(4 \times 3)$ | Auto-release cron 72 jam (`TestAutoReleaseExpiredTasks_72hSLA`); force unassign oleh Puskesmas; dan penguncian yurisdiksi IMP-001. | **Partially Effective** | **6 (M)**<br>$(2 \times 3)$ | Completed | Operasional Puskesmas / Bidang P2P | `TestAutoReleaseExpiredTasks_72hSLA` (PASS) | **CONTROL IMPROVED** |
| **RSK-13** | TIR | **Kegagalan Notifikasi Kritis Akibat Limitasi Provider Gratis** | Kuota harian SMTP Gmail dan Firebase tier gratis macet saat terjadi ledakan kasus. | **16 (H)**<br>$(4 \times 4)$ | Fallback arsitektur asynchronous logging; mode mock pengujian stabil; pengajuan integrasi SMTP resmi Diskominfo `@depok.go.id`. | **Partially Effective** | **8 (M)**<br>$(2 \times 4)$ | In Progress | Diskominfo Kota Depok | `notification_service.go` | **OPEN TECHNICAL CONSTRAINT** |
| **RSK-14** | FAR | **Penolakan Palsu pada Kluster Ledakan Kasus Nyata** | Filter anti-spam menolak otomatis laporan beruntun dalam radius 50 meter. | **15 (H)**<br>$(3 \times 5)$ | Audit trail mencatat penolakan spam; isolasi kebocoran notifikasi teruji (`TestBUGNEW012_SpamReport_NoNotificationLeak`); hak banding aktif. | **Partially Effective** | **6 (M)**<br>$(2 \times 3)$ | Completed | Tim Pengembang / Epidemiolog | `TestBUGNEW012_SpamReport_NoNotificationLeak` (PASS) | **CONTROL IMPROVED** |
| **RSK-15** | BGR | **Ketiadaan Modul Siklus Hidup Akun Relawan Lapangan** | Dashboard Puskesmas belum memiliki modul mandiri registrasi & rotasi kader RW. | **16 (H)**<br>$(4 \times 4)$ | Kontrak DTO dan endpoint backend `UserManagement` telah siap; login guard menolak akun non-aktif (`ActiveStatus = false`). UI faskes disempurnakan. | **Partially Effective** | **6 (M)**<br>$(2 \times 3)$ | In Progress | Tim Pengembang Web / Puskesmas | `user_management_handler.go` | **IN PROGRESS** |
| **RSK-16** | BCR | **Ketiadaan Prosedur Transisi Operasional Manual** | Tidak ada SOP kertas standar jika server padam > 2x24 jam saat musim wabah. | **12 (H)**<br>$(3 \times 4)$ | Protokol Kedaruratan Operasional BCP (Bab 11) diformalkan; Form Standar Kemenkes rangkap 3 disiapkan; rekonsiliasi susulan 24 jam. | **Documented** | **6 (M)**<br>$(2 \times 3)$ | Completed | Manajemen Dinkes / Bagian Organisasi | Bab 11 Dokumen Ini | **DOCUMENTED** |
| **RSK-17** | SCR | **Ketiadaan Mekanisme Penghapusan Data (Right to Erasure)** | Data pribadi pelapor tersimpan permanen di basis data tanpa alur pembersihan berkala. | **12 (H)**<br>$(3 \times 4)$ | Opsi pelaporan anonim murni (`anonymous_flag = true`) memutus identitas pelapor sejak hulu; kebijakan retensi 30 hari diajukan ke Bagian Hukum. | **Partially Effective** | **6 (M)**<br>$(2 \times 3)$ | In Progress | Legal Officer / Data Protection Officer | `TestAnonymousReport_DirectAPISubmission` (PASS) | **PARTIALLY MITIGATED** |
| **RSK-18** | TIR | **Pembatasan Lisensi Layanan Peta OpenStreetMap Publik** | Peladen OSM publik memblokir IP jika lonjakan traffic peta tanpa reverse proxy lokal. | **12 (H)**<br>$(3 \times 4)$ | Leaflet.js terpasang efisien; rekomendasi reverse proxy caching Nginx diajukan ke Diskominfo untuk disematkan pada server Balai Kota. | **Partially Effective** | **6 (M)**<br>$(2 \times 3)$ | In Progress | Technical Owner / Diskominfo Depok | `gis_handler.go` | **OPEN TECHNICAL DEPENDENCY** |
| **RSK-19**<br>*(Baru)* | SCR / FAR | **Bypass Otorisasi Wilayah pada Klaim Tugas Kader (UAT-FIND-001)** | Endpoint `/kader/reports/:id/claim` tidak memvalidasi wewenang yurisdiksi teritorial kader, hanya bergantung pada filter tampilan antarmuka (UI). | **16 (H)**<br>$(4 \times 4)$ | **IMP-001 Implemented:** Backend Jurisdiction Authorization Guard via `VerifyReportLocation`; validasi konteks lokasi JWT; respon 403 Forbidden; atomic lock terjaga. | **Effective** | **2 (L)**<br>$(1 \times 2)$ | Completed | Tim Pengembang Backend / P2P Dinkes | [IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md](file:///e:/Projects/DBD_Guard/IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md) | **RESOLVED** |

---

## 6. Risk Status Changes

Analisis komparatif pergeseran tingkat risiko (*Risk Trajectory*) membuktikan keberhasilan mitigasi terukur:

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                            MATRIKS PERGESERAN STATUS RISIKO (v4.0 ──► v5.0)                 │
├─────────────────────┬──────────────┬──────────────┬──────────────┬──────────────────────────┤
│ Kategori Risiko     │ v4.0 (Inherent)│ v5.0 (Residual)│ Delta Score  │ Faktor Pendorong Utama   │
├─────────────────────┼──────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Kritis (Critical)   │   4 Risiko   │   0 Risiko   │ -100% (-4)   │ P0-01, P0-02, P0-05, EWS │
│ Tinggi (High)       │  14 Risiko   │   2 Risiko   │  -85% (-12)  │ P0-03, P0-04, IMP-001    │
│ Sedang (Medium)     │   0 Risiko   │  12 Risiko   │  +100% (+12) │ Transisi dari High/Crit  │
│ Rendah (Low)        │   0 Risiko   │   5 Risiko   │  +100% (+5)  │ Kontrol terverifikasi    │
├─────────────────────┼──────────────┼──────────────┼──────────────┼──────────────────────────┤
│ TOTAL RISIKO        │  18 Risiko   │  19 Risiko   │  +1 (RSK-19) │ Audit komprehensif       │
└─────────────────────┴──────────────┴──────────────┴──────────────┴──────────────────────────┘
```

> **Catatan Evaluator:** Dua risiko yang masih berada pada tingkat **Tinggi (High / Nilai 12)** adalah:
> 1. **RSK-05 (Kapasitas Fogging KLB):** Bergantung secara mutlak pada keputusan stakeholder (**DEC-02**).
> 2. **RSK-07 (Resistensi Birokrasi & Mandat Regulasi):** Memerlukan penerbitan Surat Edaran / Peraturan Walikota resmi sebelum peluncuran pilot.

---

## 7. Updated Business Risks

Evaluasi risiko bisnis strategis berlandaskan dinamika operasional sektor publik Kota Depok:

1. **Risiko Akuntabilitas Keuangan Daerah (BMD):**
   - *Status Awal:* Sangat Berbahaya (Pelanggaran Permendagri No. 19/2016 akibat pemotongan stok lintas faskes otomatis).
   - *Status Pasca-P0:* **TERKENDALI (RESOLVED).** Transaksi logistik kini mencerminkan prinsip akuntansi ganda (*double-entry*) berbasis serah-terima fisik resmi (`StockTransfer`).

2. **Risiko Hukum Kerahasiaan Medis (UU PDP No. 27/2022):**
   - *Status Awal:* Sangat Berbahaya (Paparan direktori statis berkas rujukan rumah sakit).
   - *Status Pasca-P0:* **TERKENDALI (RESOLVED).** Lapisan middleware autentikasi dan otorisasi teritorial telah menutup akses publik tidak sah.

3. **Risiko Integritas Penyelidikan Epidemiologi Lapangan:**
   - *Status Awal:* Kader dapat mengklaim tugas di luar wilayah kerja binaannya (UAT-FIND-001).
   - *Status Pasca-IMP001:* **TERKENDALI (RESOLVED).** Backend secara mutlak menolak klaim tugas lintas wilayah dengan status HTTP 403 Forbidden.

4. **Risiko Keberlanjutan Anggaran Operasional (OPEX Sustainability):**
   - *Status:* **OPEN RESIDUAL.** Penggunaan layanan gratis (OpenStreetMap, SMTP, FCM) bersifat sementara. Diperlukan pos anggaran APBD untuk langganan hosting server cloud mandiri dan gateway komunikasi resmi.

---

## 8. Business Weakness & Loophole Status

Evaluasi mutakhir terhadap 10 celah bisnis (*Business Weakness*) terdahulu ditambah temuan rekonsiliasi UAT:

| ID | Area Bisnis | Celah Proses (Loophole Asal) | Status Kontrol Terkini | Efektivitas | Keterangan & Bukti Penyelesaian |
|:---|:---|:---|:---:|:---:|:---|
| **BW-01** | Pengelolaan Inventaris | Alokasi silang otomatis antar-Puskesmas ("Alokasi Hantu"). | **CLOSED** | Effective | P0-01: Auto-deduct dihapus; alur *StockTransfer* formal dengan status `IN_TRANSIT`. |
| **BW-02** | Verifikasi & PE Lapangan | Absensi verifikasi kehadiran fisik saat input data PE. | **PARTIALLY OPEN** | Partially Effective | Validasi kepemilikan investigator aktif; geofencing GPS direkomendasikan pada rilis aplikasi mobile. |
| **BW-03** | Algoritma Rekomendasi | Sistem menutup kasus sepihak tanpa tinjauan klinis dokter. | **CLOSED** | Effective | P0-03: Formula rekomendasi bersifat saran; hak *Clinical Override* dokter/Kepala Puskesmas aktif. |
| **BW-04** | Validasi & Pengawasan PE | Bypass Segregation of Duties pada pengambilalihan tugas. | **CLOSED** | Effective | P0-04: Aturan programatik `OfficerID != ValidatorID` aktif, Four-Eyes Principle terjamin. |
| **BW-05** | Pelaporan Warga Publik | Pembatalan laporan tanpa otorisasi pemilik tiket. | **CONTROL IMPROVED** | Partially Effective | Pembatalan hanya diizinkan untuk laporan berstatus `MENUNGGU_PE_KADER` yang belum diklaim investigator. |
| **BW-06** | Eksekusi Jadwal Fogging | Penutupan jadwal selesai tanpa unggahan bukti foto riil. | **CLOSED** | Effective | P1-1: Sistem menolak status `COMPLETED` jika rekaman foto di `fogging_evidences` bernilai nol. |
| **BW-07** | Deteksi Spam Otomatis | Penolakan otomatis mematikan laporan ledakan kasus keluarga. | **CONTROL IMPROVED** | Partially Effective | Laporan spam dicatat audit trail; hak banding warga aktif; tidak ada kebocoran notifikasi ke kader. |
| **BW-08** | Pelacakan Nomor Telepon | Pencarian riwayat kasus terbuka hanya berbekal input nomor HP. | **CONTROL IMPROVED** | Partially Effective | Endpoint dibatasi DTO ringkas; mode anonim melindungi PII warga secara penuh. |
| **BW-09** | Manajemen Logistik | Stok terpotong permanen meskipun pemakaian logistik ditolak. | **CONTROL IMPROVED** | Effective | Siklus approval logistik dan pemotongan saldo gudang tersinkronisasi dalam transaksi atomik. |
| **BW-10** | Administrasi Relawan | Ketiadaan manajemen siklus hidup akun Kader di Puskesmas. | **IN PROGRESS** | Partially Effective | Endpoint API `UserManagement` telah siap; antarmuka pengelolaan kader sedang dirampungkan. |
| **BW-11**<br>*(Baru)* | Otorisasi Penugasan | Kader dapat mengklaim tugas di luar wilayah binaan (UAT-FIND-001). | **CLOSED** | Effective | **IMP-001:** Penegakan otorisasi wilayah kerja di lapisan backend (`VerifyReportLocation`); respon 403 Forbidden. |

---

## 9. Updated Control Matrix

Matriks efektivitas pengendalian internal (*Internal Control Matrix*) pasca-implementasi:

| Risk ID | Nama Kontrol Baru / Yang Ditingkatkan | Tipe Kontrol | Mekanisme Sebelum Perbaikan | Mekanisme Pasca-Perbaikan | Bukti Objektif | Evaluasi Efektivitas | Celah Residual yang Tersisa |
|:---:|:---|:---:|:---|:---|:---|:---:|:---|
| **RSK-01** | Two-Party Atomic Stock Transfer | Preventif & Detektif | Pemotongan otomatis saldo donor warehouse di usecase. | Pemotongan lokal murni; mutasi antar-faskes wajib via tiket `StockTransfer` dua tahap (In-Transit -> Receive/Reject). | `TestP0_01_InventoryAccountingIntegrity` | **Effective** | Kepatuhan tanda tangan fisik BAST manual saat serah terima barang. |
| **RSK-02** | Clinical Override with Audit Log | Direktif & Detektif | Status otomatis menjadi `CLOSED_NO_FOGGING` tanpa approval. | Rekomendasi murni saran (*Advisory*); Kepala Puskesmas wajib mereview dan dapat meng-override dengan alasan medis. | `TestP0_03_and_04_DSS_and_SoD` | **Effective** | Kecepatan dokter faskes dalam mereview antrean hasil PE. |
| **RSK-03** | Protected Media Stream & Scope Guard | Preventif | Folder `/storage` dipublikasikan terbuka via web server. | Endpoint privat `/storage/photos/:filename` dengan middleware `Protected()` dan verifikasi `VerifyReportLocation`. | `TestP0_02_SensitiveFileStorageProtection` | **Effective** | Masa berlaku token JWT pada sesi peramban yang sedang aktif. |
| **RSK-08** | Programmatic SoD Barrier on Takeover | Preventif | Tidak ada pengecekan kesamaan identitas pada verifikasi. | Validasi programatik: `if peResult.OfficerID == validatorID -> Error 403 Forbidden`. | `TestP0_03_and_04_DSS_and_SoD` | **Effective** | Penugasan personel pengganti jika Puskesmas kekurangan dokter/sanitarian. |
| **RSK-09** | Hot Disaster Recovery Routine | Korektif | Kontainer database mandiri tanpa backup terjadwal. | Skrip automasi backup terkompresi `pg_dump` dan protokol restorasi live terverifikasi 100% konsisten. | DR Drill Restoration Log | **Effective** | Penyimpanan backup fisik masih berada pada host lokal (perlu off-site). |
| **RSK-19** | **Backend Jurisdiction Authorization Guard** | **Preventif** | **Hanya mengandalkan filter query UI antarmuka web.** | **Pengecekan otorisasi wilayah mutlak di usecase via `VerifyReportLocation`; penolakan 403 Forbidden.** | **`imp001_kader_claim_test.go`** | **Effective** | **Kader berpindah tugas resmi memerlukan update data `location_id` akun.** |

---

## 10. Residual Risks

Setelah seluruh kontrol terpasang, risiko residual (*sisa paparan risiko*) yang harus dikelola dalam operasional harian:

1. **Residual-01: Verifikasi Fisik Serah Terima Aset Lapangan (Terkait RSK-01)**
   - *Paparan:* Perangkat lunak mencatat transfer stok dengan sempurna, namun penyelewengan fisik di jalan (pencurian/pengurangan isi jeriken) tetap menjadi risiko fisik manusia.
   - *Mitigasi Lanjutan:* Penggunaan segel bernomor seri pada jeriken insektisida saat proses pengiriman antar-faskes.

2. **Residual-02: Beban Kerja Pemeriksaan Lapangan oleh Relawan (Terkait RSK-04)**
   - *Paparan:* Kader tetap dapat memasukkan data tanpa survei teliti jika tidak dilakukan audit sampling fisik oleh Sanitarian Puskesmas.
   - *Mitigasi Lanjutan:* Pelaksanaan audit petik uji (*spot-check*) acak 10% laporan PE oleh petugas sanitarian faskes pembina.

3. **Residual-03: Ketergantungan Eksternal Layanan Peta & Komunikasi (Terkait RSK-13 & RSK-18)**
   - *Paparan:* Pemblokiran IP oleh OSM Tile Server atau limitasi SMTP publik saat wabah memuncak.
   - *Mitigasi Lanjutan:* Integrasi dengan server proxy cache peta lokal Diskominfo dan kuota SMTP resmi pemerintah daerah.

4. **Residual-04: Penyimpanan Berkas Cadangan di Luar Peladen (Off-Site Backup) (Terkait RSK-09)**
   - *Paparan:* Skrip cadangan telah teruji, namun berkas backup masih tersimpan pada disk lokal host Docker.
   - *Mitigasi Lanjutan:* Sinkronisasi terjadwal (*rsync / rclone*) salinan cadangan terenkripsi ke Data Center Pemerintah Kota Depok.

---

## 11. Risk Treatment Status

Status implementasi rencana aksi mitigasi (*Risk Treatment Plan*):

| ID Aksi | Uraian Rencana Mitigasi | Rujukan Risiko | Target Awal | Status Terkini | Bukti / Artefak | Tindakan Lanjutan yang Diperlukan |
|:---:|:---|:---:|:---:|:---:|:---|:---|
| **P0-1** | Hapus alokasi silang otomatis; bangun *StockTransfer* formal. | RSK-01 | Rilis P0 | **Completed** | `TestP0_01_InventoryAccountingIntegrity` | Sosialisasi SOP penatausahaan aset logistik. |
| **P0-2** | Tutup folder `/storage` statis; bangun handler unduh terproteksi. | RSK-03 | Rilis P0 | **Completed** | `TestP0_02_SensitiveFileStorageProtection` | Audit berkala log unduhan foto medis. |
| **P0-3** | Restorasi diskresi klinis dokter pada formula DSS. | RSK-02 | Rilis P0 | **Completed** | `TestP0_03_and_04_DSS_and_SoD` | Pelatihan dokter Puskesmas menggunakan tombol override. |
| **P0-4** | Penegakan pemisahan tugas (Four-Eyes Principle) pada takeover. | RSK-08 | Rilis P0 | **Completed** | `TestP0_03_and_04_DSS_and_SoD` | Pastikan akun login petugas dan kepala faskes terpisah. |
| **P0-5** | Skrip cadangan basis data otomatis dan uji pemulihan bencana. | RSK-09 | Rilis P0 | **Completed** | DR Drill Restoration Log (Match 100%) | Integrasi cronjob backup ke storage sekunder pemkot. |
| **P1-1** | Validasi kewajiban foto bukti lapangan sebelum status selesai. | RSK-11 | Pilot Stage | **Completed** | Handler check `fogging_evidences` count > 0 | Edukasi tim fogging mengambil foto dokumentasi. |
| **P1-2** | Pembangunan modul manajemen kader di panel Puskesmas. | RSK-15 | Pilot Stage | **In Progress** | Kontrak DTO & Handler User Management siap | Penyelesaian komponen UI di frontend web. |
| **P1-3** | Perbaikan kueri EWS dan evaluasi threshold paralel wilayah. | RSK-06 | Pilot Stage | **Completed** | `TestCountCasesByArea` parallel verified | Kalibrasi nilai threshold KLB bersama epidemiolog. |
| **P1-4** | Validasi geofencing GPS pada pemeriksaan PE mobile. | RSK-04 | Pilot Stage | **In Progress** | Desain arsitektur Geolocation API | Integrasi native GPS pada aplikasi mobile Flutter. |
| **P1-5** | Penyusunan payung hukum regulasi daerah (SK/Perwal). | RSK-07 | Pilot Stage | **In Progress** | Draf telaahan staf bagian hukum Dinkes | Penandatanganan SK resmi oleh Kepala Dinas. |
| **IMP-001**| **Penegakan otorisasi wilayah wewenang klaim tugas investigasi kader.** | **RSK-19** | **UAT Stage** | **Completed** | **[IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md](file:///e:/Projects/DBD_Guard/IMP001_IMPLEMENTATION_VERIFICATION_DBD_GUARD.md)** | **Penyelarasan master data penugasan kader per RW.** |

---

## 12. KRI & Monitoring

Tinjauan dan rekalibrasi Indikator Risiko Utama (*Key Risk Indicators*):

| Kode KRI | Indikator Risiko | Formula Pemantauan | Ambang Peringatan (Warning) | Ambang Bahaya (Critical) | Frekuensi | Penanggung Jawab | Ketersediaan Data |
|:---:|:---|:---|:---:|:---:|:---:|:---|:---:|
| **KRI-01** | Pelanggaran SLA Laporan | $\frac{\text{Kasus Melewati 72 Jam}}{\text{Total Laporan Masuk}} \times 100\%$ | > 5% / minggu | > 15% / minggu | Mingguan | Kepala Seksi Surveilans | Tersedia via `sla_cron` |
| **KRI-02** | Frekuensi Transfer Darurat Logistik | Jumlah transaksi *StockTransfer* antar-Puskesmas. | > 3 kali / bulan | > 8 kali / bulan | Bulanan | Pengelola Barang Dinkes | Tersedia via `stock_transfers` |
| **KRI-03** | Rasio Override Keputusan DSS | $\frac{\text{Override Rekomendasi}}{\text{Total Hasil PE}} \times 100\%$ | > 15% / bulan | > 35% / bulan | Bulanan | Tim Epidemiolog Dinkes | Tersedia via `audit_logs` |
| **KRI-04** | Rasio Deteksi Spam Beruntun | $\frac{\text{Laporan Terindikasi Spam}}{\text{Total Laporan Masuk}} \times 100\%$ | > 5% / minggu | > 12% / minggu | Mingguan | Technical Owner / Diskominfo | Tersedia via `citizen_reports` |
| **KRI-05** | Tingkat Inaktivitas Akun Relawan | $\frac{\text{Kader Tidak Login > 30 Hari}}{\text{Total Kader Terdaftar}} \times 100\%$ | > 20% | > 40% | Bulanan | Penanggung Jawab Kader Faskes | Tersedia via `users` |
| **KRI-06** | Durasi Antrean Approval Fogging | Waktu tunggu rata-rata persetujuan Kepala Puskesmas. | > 12 jam | > 24 jam | Harian | Kepala Puskesmas | Tersedia via timestamp jadwal |
| **KRI-07** | Tingkat Penolakan Klaim Wilayah | Frekuensi percobaan klaim menghasilkan HTTP 403. | > 1% klaim | > 5% klaim | Mingguan | Tim Keamanan IT Dinkes | Tersedia via log error HTTP |

---

## 13. Traceability to Audits & Implementations

Pemetaan komprehensif penelusuran dari temuan audit awal hingga verifikasi rilis:

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                          MATRIKS KETERLACAKAN AUDIT ──► IMPLEMENTASI ──► VERIFIKASI         │
├─────────────┬──────────────────┬─────────────────┬──────────────────────┬───────────────────┤
│ ID Risiko   │ Temuan Asal      │ Rencana Solusi  │ Eksekusi Implementasi│ Bukti Pengujian   │
├─────────────┼──────────────────┼─────────────────┼──────────────────────┼───────────────────┤
│ **RSK-01**  │ DBA-1 (BUG-03/09)│ P0-1 Plan       │ inventory_usecase.go │ TestP0_01 (PASS)  │
│ **RSK-02**  │ DBA-1 / BRR      │ P0-3 Plan       │ pe_usecase.go        │ TestP0_03 (PASS)  │
│ **RSK-03**  │ DBA-2 (SEC-02)   │ P0-2 Plan       │ report_handler.go    │ TestP0_02 (PASS)  │
│ **RSK-08**  │ BRR (SEC-01)     │ P0-4 Plan       │ pe_usecase.go        │ TestP0_04 (PASS)  │
│ **RSK-09**  │ DBA-2 (INFRA-01) │ P0-5 Plan       │ backup_restore.ps1   │ DR Drill (PASS)   │
│ **RSK-19**  │ UAT-FIND-001     │ IMP-001 Plan    │ report_handler.go &  │ imp001 tests (PASS│
│             │ (Rekonsiliasi)   │ (Group A)       │ report_usecase.go    │ 13/13 PASS)       │
└─────────────┴──────────────────┴─────────────────┴──────────────────────┴───────────────────┘
```

---

## 14. Pending Stakeholder Decisions

Sesuai tata kelola pemerintahan, tiga agenda kebijakan birokrasi berikut **TIDAK DIUBAH PADA KODE** dan tetap berstatus **`PENDING STAKEHOLDER DECISION`**:

### 14.1 DECISION-01: Alur Draf Permintaan Stok dari FOGGING_TERTUNDA
* **Masalah:** Rekomendasi stok otomatis langsung berstatus `PENDING` ke Dinkes, memotong wewenang tanda tangan Kepala Puskesmas.
* **Status:** **PENDING STAKEHOLDER DECISION**
* **Cabang Pilihan:**
  - *Option A:* Langsung berstatus `PENDING` ke Dinas Kesehatan demi kecepatan penanganan darurat wabah (Status Quo PRD 5.2).
  - *Option B:* Berstatus `DRAFT_PUSKESMAS` dan wajib disetujui internal oleh Kepala Puskesmas sebelum diteruskan ke Dinkes (Kepatuhan Permendagri 19/2016).

### 14.2 DECISION-02: Batas Kapasitas Penjadwalan Fogging Saat KLB
* **Masalah:** Aturan kaku 1 tim fogging maksimal 1 tugas per hari melumpuhkan penanganan saat terjadi lonjakan wabah (*outbreak*).
* **Status:** **PENDING STAKEHOLDER DECISION**
* **Cabang Pilihan:**
  - *Option A:* Strict 1 task/hari/tim tanpa pengecualian (K3 mutlak).
  - *Option B:* Emergency override hingga 2 tugas/hari saat kelurahan berstatus KLB aktif.
  - *Option C:* Tetap 1 tugas/hari, kelebihan tugas didelegasikan ke Regu BKO Cadangan Dinkes.
  - *Option D:* Hibrida berjenjang (Lokal maksimal 2 tugas untuk radius $<500$ m, selebihnya dialihkan ke BKO Dinkes).

### 14.3 DECISION-03: Jendela Batas Waktu Kedaluwarsa Banding Warga
* **Masalah:** PRD hanya membatasi kuota 1 kali banding tanpa tenggat kedaluwarsa waktu.
* **Status:** **PENDING STAKEHOLDER DECISION**
* **Cabang Pilihan:**
  - *Option A:* Maksimal 1 kali tanpa batas kedaluwarsa temporal (*No Expiration*).
  - *Option B:* Maksimal 1 kali dengan batas kedaluwarsa ketat 48 Jam pasca penolakan.
  - *Option C:* Maksimal 1 kali dengan batas kedaluwarsa 72 Jam (Selaras dengan SPM Dinkes Depok).

---

## 15. Open Risks / Open Actions

Daftar tindakan terbuka (*Open Actions*) yang wajib dipantau hingga fase *Pilot Rollout*:

1. **Pengesahan Lembar Ballot Keputusan Kebijakan:**
   Memperoleh tandatangan resmi Kepala Dinas Kesehatan Kota Depok untuk menetapkan pilihan definitif pada DEC-01, DEC-02, dan DEC-03.
2. **Penyediaan Server & Jaringan Diskominfo:**
   Konfigurasi peladen fisik di Data Center Pemkot Depok, alokasi akun SMTP resmi `@depok.go.id`, dan reverse proxy cache peta OpenStreetMap.
3. **Penyusunan Surat Edaran / Peraturan Walikota:**
   Penerbitan landasan hukum instruksi wajib pakai aplikasi bagi seluruh Puskesmas dan aparatur kelurahan se-Kota Depok.
4. **Penyelesaian Antarmuka UI Manajemen Kader:**
   Finalisasi komponen halaman web dashboard Puskesmas untuk pendaftaran mandiri relawan RW.

---

## 16. Assumptions & Limitations

Penilaian risiko ini disusun berdasarkan batasan dan asumsi realistis:

1. **Asumsi Integritas Master Data:**
   Diasumsikan master data relasi 11 Kecamatan $\rightarrow$ 63 Kelurahan $\rightarrow$ 38 Puskesmas pada tabel `master_kelurahans` telah terisi lengkap, akurat, dan mencerminkan wilayah kerja faskes Depok.
2. **Batasan Kebijakan 100% Gratis/Open-Source:**
   Sistem tidak mengimplementasikan SMS Gateway berbayar (Twilio/Zenziva) maupun WhatsApp Business API resmi. Notifikasi sepenuhnya mengandalkan email SMTP dan push notification Firebase (FCM).
3. **Batasan Tanggung Jawab Diagnostik:**
   Aplikasi DBD Guard berfungsi sebagai instrumen pelaporan dan surveilans epidemiologi kesehatan masyarakat, bukan alat penegakan diagnosis klinis pengganti dokter atau rekam medis elektronik rumah sakit.

---

## 17. Final Risk Posture

Berdasarkan seluruh hasil pengujian dan implementasi mutakhir, ditarik kesimpulan formal:

> **"Kondisi profil risiko proyek DBD Guard (Kota Depok) telah meningkat secara substansial dari status awal *Conditional Hold* menjadi *Technical & Security Clearance Achieved*. Seluruh celah kritis pada akuntansi inventaris, proteksi foto rekam medis, segregasi kewenangan klinis, mitigasi bencana basis data, serta otorisasi teritorial klaim tugas kader (IMP-001) telah berhasil diselesaikan dan diverifikasi 100% tanpa regresi.**
>
> **Sistem dinyatakan siap secara teknis untuk melangkah ke tahap Peluncuran Percontohan Lapangan (*Pilot Project Rollout* di Puskesmas Beji & Puskesmas Pancoran Mas), dengan syarat formal: penetapan pengesahan resmi atas 3 agenda keputusan tata kelola (DEC-01, DEC-02, DEC-03) oleh pimpinan Dinas Kesehatan Kota Depok."**

---
*Dokumen Manajemen Risiko v5.0 ini disahkan secara digital sebagai rujukan pengendali mutu, keamanan, dan kepatuhan sistem DBD Guard Kota Depok.*
