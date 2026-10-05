# Implementation Plan — DBD Guard
**Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi (Pemerintah Kota Depok)**

- **Tanggal Dokumen:** 29 September 2026
- **Status Dokumen:** Master Implementation Blueprint
- **Prinsip Rekayasa Perangkat Lunak:** **Minimal → Surgical → Traceable → Testable → Reversible**
- **Aturan Eksekusi Saat Ini:** **ZERO CODE MODIFICATION** (Dokumen ini berfungsi sebagai cetak biru teknis & bisnis sebelum modifikasi kode dilakukan)
- **Status Kesiapan Akhir:** **PARTIALLY READY** (Group A: Ready for Implementation; Group B: Blocked by Business Decision)

---

## 1. Executive Summary

Proyek **DBD Guard Kota Depok** telah melalui rangkaian audit ketat yang mencakup *Deep Bug Audit*, *P0 Fundamental Fix & Re-Verification*, *Disaster Recovery Drill*, *Business & Operational UAT*, *UAT Reconciliation*, serta *Business Decision & Governance Resolution*. Hasil verifikasi membuktikan bahwa baseline arsitektur teknis sistem berada pada status **Technical PASS**.

Dokumen **Implementation Plan** ini disusun sebagai *engineering blueprint* yang menjembatani temuan validasi dan keputusan tata kelola menuju eksekusi teknis yang terarah. Selaras dengan disiplin tata kelola birokrasi pemerintahan:
1. **GROUP A (Confirmed Implementation):** Terdiri dari cacat implementasi keamanan bisnis **UAT-FIND-001 (Kader Cross-Region Task Claim Bypass - P1)** yang telah terverifikasi secara teknis dan memiliki spesifikasi perbaikan langsung tanpa memerlukan perdebatan kebijakan.
2. **GROUP B (Pending Stakeholder Decisions):** Terdiri dari tiga agenda kebijakan birokrasi (**DECISION-01** alur stok obat, **DECISION-02** kapasitas fogging KLB, dan **DECISION-03** batas waktu kedaluwarsa banding). Seluruh opsi teknis dipetakan dampaknya secara komprehensif, namun **pengerjaan kode diblokir total** hingga pemangku kepentingan (*stakeholders*) menandatangani lembar pengesahan resmi.

Perbaikan teknis berpegang teguh pada aturan **"Don't Fix What Is Not Broken"**: fondasi atomik konservasi stok insektisida, proteksi foto medis via storage key, peran pemisahan wewenang (*SoD*), dan mesin status laporan yang telah berstatus PASS tidak boleh diotak-atik.

---

## 2. Current Baseline

Berikut adalah kondisi faktual seluruh komponen sistem DBD Guard yang menjadi titik tolak (*baseline*) rencana implementasi:

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                             BASELINE SISTEM AKTUAL DBD GUARD                                │
├────────────────────────────┬─────────────────────────────┬──────────────────────────────────┤
│ Komponen / Layer           │ Kondisi Aktual              │ Catatan Integritas & Verifikasi  │
├────────────────────────────┼─────────────────────────────┼──────────────────────────────────┤
│ 1. Backend Core            │ Go 1.22+, Fiber v2, GORM    │ Test suite passing 100% (exit 0) │
│ 2. Database Schema         │ PostgreSQL 15               │ 15 tabel operasional, SoD OK     │
│ 3. P0 Invariants           │ Zero Leaks Verified         │ Akuntansi stok & foto terlindungi│
│ 4. Disaster Recovery       │ Hot Drill Verified          │ Live restore drill match 100%    │
│ 5. Frontend Web            │ Next.js 14 (App Router)     │ Dashboard 5 role operasional     │
│ 6. Open Confirmed Defect   │ 1 Cacat (UAT-FIND-001)      │ Bypass klaim tugas lintas faskes │
│ 7. Open Business Decision  │ 3 Isu Kebijakan             │ Menunggu sign-off pemda          │
└────────────────────────────┴─────────────────────────────┴──────────────────────────────────┘
```

---

## 3. Confirmed Implementation Items

Kelompok **GROUP A** mencakup item yang telah terkonfirmasi cacat (*confirmed defect*), memiliki dampak keamanan bisnis, dan wajib segera diperbaiki tanpa menunggu keputusan birokrasi:

* **Item ID:** `IMP-001`
* **Sumber Temuan:** `UAT-FIND-001` (Hasil Rekonsiliasi UAT)
* **Judul:** Penegakan Otorisasi Teritorial Wilayah Kerja pada Klaim Tugas Investigasi Kader (*Kader Cross-Region Task Claim Authorization*)
* **Klasifikasi:** **P1 BUSINESS / SECURITY BUG**
* **Status Kesiapan:** **READY FOR IMPLEMENTATION**

---

## 4. UAT-FIND-001 Implementation

### 4.1 Deskripsi Masalah & Analisis Faktual
- **Current Behavior:**
  Kader Posyandu/Kelurahan mengklaim laporan melalui endpoint `POST /api/v1/kader/reports/:id/claim`. Pada file `backend/internal/delivery/http/report_handler.go:278`, handler hanya mengekstrak `actorID` via `middleware.GetUserID(c)`. Pada `backend/internal/usecase/report_usecase.go:273`, fungsi `ClaimReport(ctx, id, kaderID)` hanya mengecek:
  1. Status laporan bernilai `MENUNGGU_PE_KADER`;
  2. Kolom `investigator_id IS NULL`.
  Sistem sama sekali tidak memeriksa apakah wilayah laporan cocok dengan wilayah binaan (`location_id`) milik akun kader yang sedang login.
- **Root Cause:**
  Pembatasan wilayah kerja kader baru sebatas filter kueri penampilan daftar laporan di antarmuka web (`GetReports`), tetapi **luput diimplementasikan pada lapisan eksekusi mutasi klaim tugas di usecase (*missing business authorization guard*)**.
- **Expected Behavior:**
  Otorisasi wilayah wajib ditegakkan secara mutlak di lapisan *Backend Usecase*. Sistem wajib menolak klaim apabila laporan berada di luar wilayah Puskesmas/Kelurahan binaan kader dengan mengembalikan respon **`403 FORBIDDEN`**.

### 4.2 Kriteria Penerimaan (*Acceptance Criteria*)

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                           ACCEPTANCE CRITERIA MATRIX (UAT-FIND-001)                         │
├───────┬──────────────────────┬──────────────────────┬───────────────────────────────────────┤
│ Kasus │ Kondisi Aktor        │ Kondisi Laporan      │ Hasil yang Diharapkan (Expected)      │
├───────┼──────────────────────┼──────────────────────┼───────────────────────────────────────┤
│ A     │ Kader A (Wilayah X)  │ Laporan di Wilayah X │ CLAIM SUCCESS (HTTP 200 OK)           │
│ B     │ Kader A (Wilayah X)  │ Laporan di Wilayah Y │ 403 FORBIDDEN ("Akses Ditolak: Di luar│
│       │                      │                      │ yurisdiksi wilayah kerja")            │
│ C     │ Direct API Call      │ Laporan di Wilayah Y │ 403 FORBIDDEN (Bypass UI tertutup)    │
│ D     │ Dua Kader Bersamaan  │ Laporan di Wilayah X │ Tepat SATU klaim sukses (Atomic lock),│
│       │ (Race Condition)     │                      │ kader kedua menerima error konflik    │
│ E     │ Aktor Non-Kader      │ Laporan di Wilayah X │ 403 FORBIDDEN (Role guard middleware) │
└───────┴──────────────────────┴──────────────────────┴───────────────────────────────────────┘
```

### 4.3 Pemanfaatan Komponen Eksisting (Minimal Change Rule)
Sesuai prinsip rekayasa minimal, perbaikan tidak memerlukan penambahan tabel atau repositori baru karena komponen pendukung **sudah tersedia di codebase**:
1. `middleware.GetLocationID(c)` sudah tersedia di `backend/internal/delivery/http/middleware/auth_middleware.go:105`.
2. `reportRepo.VerifyReportLocation(ctx, reportID, locationID)` sudah tersedia di `backend/internal/repository/report_repository.go:351` yang menghubungkan `report_locations.kelurahan` dengan `master_kelurahans.puskesmas_id`.

---

## 5. Pending Stakeholder Decisions

Kelompok **GROUP B** memuat keputusan bisnis yang belum disahkan oleh pemangku kepentingan Kota Depok. Rencana implementasi teknis disiapkan untuk masing-masing cabang opsi, namun implementasinya **DITANGGUHKAN (*GATED*)**:

### 5.1 DEC-01 — Stock Request dari FOGGING_TERTUNDA
* **Status:** **PENDING STAKEHOLDER DECISION**
* **Cabang Kebijakan:**
  * **Option A:** Auto-draft langsung berstatus `PENDING` menuju Dinas Kesehatan (Status Quo PRD Bab 5.2).
  * **Option B:** Auto-draft berstatus `DRAFT_PUSKESMAS` dan wajib disetujui internal oleh Kepala Puskesmas sebelum menjadi `PENDING` (Permendagri 19/2016).

### 5.2 DEC-02 — Emergency Fogging Capacity (Kapasitas Saat KLB)
* **Status:** **PENDING STAKEHOLDER DECISION**
* **Cabang Kebijakan:**
  * **Option A:** Strict 1 task per hari per regu tanpa pengecualian (K3 mutlak).
  * **Option B:** Emergency override hingga 2 tugas per hari saat kelurahan berstatus KLB aktif.
  * **Option C:** Kuota tim faskes tetap 1 tugas/hari, kelebihan tugas didelegasikan ke Regu BKO Dinkes.
  * **Option D:** Hybrid (Lokal maksimal 2 tugas untuk radius $<500$ m, selebihnya didelegasikan ke BKO).

### 5.3 DEC-03 — Batas Kedaluwarsa Hak Banding Warga (*Appeal Expiration*)
* **Status:** **PENDING STAKEHOLDER DECISION**
* **Cabang Kebijakan:**
  * **Option A:** Maksimal 1 kali pengajuan tanpa batas waktu kedaluwarsa (*No Expiry*).
  * **Option B:** Maksimal 1 kali pengajuan dengan batas waktu ketat 48 Jam.
  * **Option C:** Maksimal 1 kali pengajuan dengan batas waktu 72 Jam (Selaras SLA SPM Dinkes Depok).

---

## 6. Decision-Dependent Implementation Options

Analisis komparatif arsitektur dan dampak teknis untuk setiap opsi kebijakan yang mungkin disahkan:

### 6.1 Matriks Teknis DEC-01 (Alur Permintaan Stok)

| Komponen Teknis | OPTION A (Auto-Direct Dinkes) | OPTION B (Internal Gate Puskesmas) |
|---|---|---|
| **Backend Layer** | Tidak ada perubahan; `pe_usecase.go:267` tetap menetapkan `Status: PENDING`. | Ubah `pe_usecase.go` mengisi `Status: DRAFT_PUSKESMAS`; tambah method `SubmitStockRequest(ctx, id, kepalaPuskesmasID)`. |
| **Frontend Web** | Tidak ada perubahan; Dinkes langsung melihat request baru di dashboard. | Tambah tab *"Draf Pengajuan Faskes"* di panel Puskesmas `/puskesmas/inventory-approval` dengan tombol *"Kirim ke Dinkes"*. |
| **State Machine** | `PENDING -> APPROVED / REJECTED`. | `DRAFT_PUSKESMAS -> PENDING -> APPROVED / REJECTED`. |
| **Database Schema** | Skema saat ini tetap (kolom status `stock_requests`). | Tambah nilai `'DRAFT_PUSKESMAS'` pada constraint check enum status `stock_requests`. |
| **Authorization** | Dilakukan oleh Petugas PE saat verifikasi kasus. | Otorisasi transisi draf dikunci ketat untuk `RoleKepalaPuskesmas`. |
| **Audit Logging** | Event: `STOCK_REQUEST_AUTO_GENERATED`. | Event: `STOCK_REQUEST_DRAFT_CREATED` & `STOCK_REQUEST_SUBMITTED_BY_KEPALA_PUSKESMAS`. |
| **Test Impact** | Suite test inventaris eksisting tetap valid. | Tambah test scenario otorisasi pengajuan draf faskes di `inventory_usecase_test.go`. |

### 6.2 Matriks Teknis DEC-02 (Kapasitas Fogging KLB)

| Parameter Desain | OPTION A (Strict) | OPTION B (KLB Override) | OPTION C (Dinkes BKO) | OPTION D (Hybrid B+C) |
|---|---|---|---|---|
| **Workflow Impact** | Tolak jadwal jika tim sudah ada tugas pada tanggal yang sama. | Cek status KLB; jika aktif, izinkan jadwal ke-2 untuk tim yang sama. | Jika tim faskes penuh, rute penugasan dialihkan ke pool Tim BKO Dinkes. | Evaluasi jarak/radius; izinkan tugas ke-2 lokal atau eskalasi ke BKO Dinkes. |
| **Authorization** | `RolePetugasPuskesmas`. | `RolePetugasPuskesmas` + verifikasi status KLB aktif. | `RoleDinkes` mengalokasikan armada cadangan kota. | Berjenjang: Faskes (tugas ke-2) lalu Dinkes (tugas ke-3). |
| **Data Model** | Tidak ada perubahan tabel. | Tambah kolom `is_emergency_override` & `override_reason` di `fogging_schedules`. | Tambah relasi `origin_agency` / flag `is_bko` pada entitas `fogging_teams`. | Gabungan kolom override pada jadwal dan penanda armada BKO Dinkes. |
| **Schedule Guard** | `count >= 1 -> Error`. | `if !isKLB && count >= 1 -> Error; if isKLB && count >= 2 -> Error`. | `count >= 1 -> Alihkan ke Tim BKO`. | Validasi bertingkat kuota lokal dan armada bantuan. |
| **Notifikasi** | Normal. | Notifikasi siaga K3 ke koordinator fogging faskes. | Notifikasi perbantuan tugas ke posko BKO Dinkes Balai Kota. | Notifikasi hybrid berjenjang. |
| **Audit Trail** | Standar jadwal. | `FOGGING_SCHEDULE_KLB_OVERRIDE` (Wajib catat alasan). | `FOGGING_TASK_DELEGATED_TO_DINKES_BKO`. | Catat event override lokal dan delegasi dinas. |
| **Test Suite** | Test eksisting PASS. | Tambah test case penjadwalan ganda saat KLB. | Tambah test penugasan armada BKO lintas faskes. | Tambah test suite eskalasi darurat hibrida. |
| **Dependensi Ops** | Nol. | SOP uang lembur & rotasi operator K3 lokal. | SOP kesiapan operasional tim mobil pickup Dinkes. | Koordinasi terpadu faskes & dinas kota. |

### 6.3 Matriks Teknis DEC-03 (Batas Waktu Banding)

| Parameter Desain | OPTION A (No Expiry) | OPTION B (48 Jam) | OPTION C (72 Jam / SPM) |
|---|---|---|---|
| **Database Impact** | Tidak ada perubahan kolom. | Tidak ada perubahan skema (memanfaatkan kolom `updated_at`). | Tidak ada perubahan skema (memanfaatkan kolom `updated_at`). |
| **Repository Query** | Kueri atomic hanya mengecek `appeal_count < 1`. | `WHERE report_id = ? AND appeal_count < 1 AND updated_at >= NOW() - INTERVAL '48 hours'`. | `WHERE report_id = ? AND appeal_count < 1 AND updated_at >= NOW() - INTERVAL '72 hours'`. |
| **Usecase Validation** | Hanya validasi kuota 1 kali. | Validasi `time.Since(report.UpdatedAt) <= 48*time.Hour`. | Validasi `time.Since(report.UpdatedAt) <= 72*time.Hour`. |
| **Frontend UX** | Tombol banding aktif selama kuota banding belum terpakai. | Tombol dinonaktifkan jika $>48$ jam, tampil label *"Batas Banding (48 Jam) Kedaluwarsa"*. | Tombol dinonaktifkan jika $>72$ jam, tampil countdown waktu sisa banding. |
| **Notifikasi** | Normal saat warga submit banding. | Reminder SMS/FCM pada jam ke-36 sebelum hak banding hangus. | Reminder FCM pada jam ke-48 sebelum hak banding hangus. |
| **Migration** | No database migration required. | No database migration required. | No database migration required. |
| **Automated Test** | Test eksisting PASS. | Tambah unit test banding valid $<48$j dan ditolak $>48$j. | Tambah unit test banding valid $<72$j dan ditolak $>72$j. |

---

## 7. Dependency Map

Diagram alur dependensi pengerjaan teknis:

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                           DBD GUARD IMPLEMENTATION DEPENDENCY MAP                           │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

 [TRACK 1: CONFIRMED BUG FIX (GROUP A)]
 
  UAT-FIND-001 Confirmed (P1)
         │
         ▼
  Inspect Existing Code (handler, usecase, repo) ──[Sudah Selesai: Zero New Abstraction]
         │
         ▼
  Surgical Code Fix:
  1. report_handler.go (Ekstrak location_id dari JWT)
  2. report_usecase.go (Validasi VerifyReportLocation)
         │
         ▼
  Automated Test Matrix Suite (5 Scenarios)
         │
         ▼
  Regression Testing (P0 Invariants, SLA Cron, Kader Normal Claim)
         │
         ▼
  Release Gate 1 (Technical PASS) ──► READY FOR GROUP A MERGE


 [TRACK 2: GOVERNANCE & BUSINESS DECISIONS (GROUP B)]
 
  Stakeholder Policy Deliberation (Pemkot & Dinkes Depok)
         │
         ├───────────────────────────────┬───────────────────────────────┐
         ▼                               ▼                               ▼
    DEC-01 Decided                  DEC-02 Decided                  DEC-03 Decided
  (Option A vs B)                (Option A/B/C/D)                (Option A/B/C)
         │                               │                               │
         ▼                               ▼                               ▼
  BR-INV-04 Finalized             BR-FOG-05 Finalized             BR-REP-08 Finalized
         │                               │                               │
         ▼                               ▼                               ▼
  Implementation Option           Implementation Option           Implementation Option
         │                               │                               │
         └───────────────────────────────┼───────────────────────────────┘
                                         │
                                         ▼
                             Release Gate 2 & 3 (Business & Governance PASS)
                                         │
                                         ▼
                             Release Gate 4 (Pilot Rollout Depok)
```

**Aturan Dependensi:** Track 1 (Group A) bersifat independen murni dan tidak memiliki dependensi terhadap Track 2 (Group B).

---

## 8. Implementation Master Table

| ID | Item Pekerjaan | Tipe | Prioritas | Sumber Temuan | Status Kesiapan | Perubahan Backend | Perubahan Frontend | Perubahan DB | Pengujian | Dependensi |
|---|---|---|:---:|---|:---:|---|---|---|---|---|
| **IMP-001** | Penegakan Otorisasi Teritorial Wilayah Kerja Klaim Kader | Security / Bug Fix | **P1** | `UAT-FIND-001` | **CONFIRMED (Ready)** | `report_handler.go` & `report_usecase.go` | Tidak ada | Tidak ada | Unit, Concurrency, & Integration | Tidak ada |
| **IMP-002A**| Auto-Draft Permintaan Stok Langsung ke Dinkes | Business Rule | P2 | `DEC-01 (Option A)` | **PENDING DECISION** | Nol (Status Quo) | Nol | Tidak ada | Regresi inventaris | Sign-off DEC-01 |
| **IMP-002B**| Alur Verifikasi Draf Stok oleh Kepala Puskesmas | Governance | P2 | `DEC-01 (Option B)` | **PENDING DECISION** | `pe_usecase.go` & `inventory_usecase.go` | Tab draf faskes di panel Puskesmas | Check enum status `stock_requests` | Unit test otorisasi faskes | Sign-off DEC-01 |
| **IMP-003A**| Batas Kaku Penjadwalan 1 Tugas Fogging Per Hari | Operational | P2 | `DEC-02 (Option A)` | **PENDING DECISION** | Nol (Status Quo) | Nol | Tidak ada | Regresi penjadwalan | Sign-off DEC-02 |
| **IMP-003B**| Emergency Override Penjadwalan saat Wilayah KLB | Operational | P2 | `DEC-02 (Option B)` | **PENDING DECISION** | `fogging_schedule_usecase.go` | Form input alasan override KLB | Kolom override di `fogging_schedules` | Multi-tugas KLB test | Sign-off DEC-02 |
| **IMP-003C**| Penugasan Tim Cadangan BKO Dinkes | Governance | P2 | `DEC-02 (Option C)` | **PENDING DECISION** | `fogging_schedule_usecase.go` | Dropdown armada BKO kota | Flag tim BKO di `fogging_teams` | Delegasi armada dinas test | Sign-off DEC-02 |
| **IMP-003D**| Mekanisme Kontingensi Hibrida Fogging KLB | Operational | P2 | `DEC-02 (Option D)` | **PENDING DECISION** | Mesin orkestrasi jadwal bertingkat | Form eskalasi berjenjang | Kolom override & BKO | Integrasi skenario darurat | Sign-off DEC-02 |
| **IMP-004A**| Hak Banding Warga 1 Kali Tanpa Batas Waktu | Business Rule | P3 | `DEC-03 (Option A)` | **PENDING DECISION** | Nol (Status Quo) | Nol | Tidak ada | Regresi banding | Sign-off DEC-03 |
| **IMP-004B**| Batas Waktu Kedaluwarsa Banding 48 Jam | Business Rule | P3 | `DEC-03 (Option B)` | **PENDING DECISION** | `report_repository.go` filter interval | Label expired di portal warga | Tidak ada | Test banding 48 jam | Sign-off DEC-03 |
| **IMP-004C**| Batas Waktu Kedaluwarsa Banding 72 Jam (SLA SPM) | Business Rule | P3 | `DEC-03 (Option C)` | **PENDING DECISION** | `report_repository.go` filter interval | Countdown sisa jam di portal warga | Tidak ada | Test banding 72 jam | Sign-off DEC-03 |

---

## 9. Detailed Implementation Plan

Berikut adalah rencana teknis mendalam untuk item **IMP-001**:

### IMP-001 — Kader Cross-Region Task Claim Authorization

#### Source
Hasil rekonsiliasi audit operasional UAT (`UAT-FIND-001`), diklasifikasikan sebagai **CONFIRMED P1 BUSINESS / SECURITY BUG**.

#### Problem
Kader di Puskesmas Beji dapat mengklaim tugas investigasi laporan yang berada di Kelurahan Pancoran Mas jika menembak endpoint API `POST /api/v1/kader/reports/:id/claim` secara langsung dengan menyertakan JWT token kader.

#### Expected Behavior
Sistem menolak mutasi klaim dan mengembalikan HTTP 403 Forbidden apabila `puskesmas_id` dari kelurahan laporan tidak sesuai dengan `location_id` kader yang tersimpan pada token JWT.

#### Current Behavior
Usecase `ClaimReport` hanya memeriksa status laporan dan ketiadaan investigator lain, sehingga klaim lintas wilayah berhasil diproses.

#### Root Cause
Usecase `ClaimReport` tidak menerima parameter identitas lokasi aktor (`userLocationID`) dan tidak memanggil fungsi validasi wilayah yang sudah tersedia di repositori.

#### Proposed Change
1. Modifikasi kontrak interface `ReportUsecase`:
   ```go
   ClaimReport(ctx context.Context, id uuid.UUID, kaderID uuid.UUID, locationID *uuid.UUID) error
   ```
2. Pada `backend/internal/delivery/http/report_handler.go:ClaimReport`:
   - Ekstrak `locationID, _ := middleware.GetLocationID(c)`
   - Teruskan `locationID` ke usecase.
   - Jika usecase mengembalikan error `domain.ErrUnauthorizedJurisdiction` atau pesan otorisasi wilayah, kembalikan respon HTTP `fiber.StatusForbidden` (403).
3. Pada `backend/internal/usecase/report_usecase.go:ClaimReport`:
   - Tambahkan validasi:
     ```go
     if locationID != nil {
         allowed, err := u.reportRepo.VerifyReportLocation(ctx, id, *locationID)
         if err != nil {
             return err
         }
         if !allowed {
             return domain.ErrUnauthorizedJurisdiction
         }
     }
     ```
4. Tambahkan definisi error di `backend/internal/domain/report.go`:
   ```go
   var ErrUnauthorizedJurisdiction = errors.New("akses ditolak: laporan di luar yurisdiksi wilayah kerja")
   ```

#### Affected Components
- **Backend:** `ReportHandler`, `ReportUsecase`, `domain.Report`.
- **Frontend:** Tidak ada perubahan kode (UI sudah memfilter laporan via `GetReports` dan siap menerima error jika terjadi penolakan).
- **Database:** Tidak ada perubahan skema (memanfaatkan data relasi pada tabel `master_kelurahans`).
- **Infrastructure:** Tidak ada dampak.

#### Files yang Dimodifikasi
1. `backend/internal/domain/report.go`
2. `backend/internal/delivery/http/report_handler.go`
3. `backend/internal/usecase/report_usecase.go`
4. `backend/internal/usecase/bug08_claim_race_test.go` (Update signature pemanggilan)
5. `backend/internal/usecase/report_usecase_test.go` (Tambah suite pengujian baru)

#### API Impact
- Endpoint: `POST /api/v1/kader/reports/:id/claim`
- Tambahan Response Code: `403 Forbidden` jika klaim melanggar yurisdiksi. Format respon:
  ```json
  {
    "error": "Akses ditolak: Laporan berada di luar yurisdiksi wilayah kerja Anda",
    "message": "akses ditolak: laporan di luar yurisdiksi wilayah kerja"
  }
  ```

#### Database Impact
**NOL (Zero DB Impact):** Kueri verifikasi menggunakan indeks yang telah ada pada tabel `report_locations` dan `master_kelurahans`.

#### Business Rule Impact
Menegakkan Aturan Bisnis Bab 3 PRD DBD Guard mengenai batas yurisdiksi teritorial faskes.

#### Security Impact
Menutup kerentanan *Broken Object Level Authorization (BOLA / IDOR)* pada alur penugasan faskes.

#### Regression Risk
Sangat rendah (*Surgical Fix*). Hanya menyentuh alur otorisasi fungsi klaim kader tanpa memodifikasi alur pembatalan, pelacakan, ataupun verifikasi.

#### Test Plan
- Unit test mock usecase untuk skenario wilayah sama dan wilayah berbeda.
- Integration test live database dengan multi-Puskesmas.
- Concurrency test (uji balap klaim kader).

#### Acceptance Criteria
1. Kader A (Beji) klaim Laporan A (Beji) $\rightarrow$ **200 OK**.
2. Kader A (Beji) klaim Laporan B (Pancoran Mas) $\rightarrow$ **403 Forbidden**.
3. Direct curl tanpa header wilayah $\rightarrow$ **403 Forbidden**.
4. Dua kader satu kelurahan klaim bersamaan $\rightarrow$ **1 Sukses, 1 Gagal**.

#### Rollback Consideration
Perubahan sepenuhnya *stateless* pada lapisan aplikasi. Jika terjadi anomali, kode dapat di-revert secara instan melalui git tanpa memerlukan rollback basis data.

---

## 10. Test Strategy

Pengujian disusun berlapis untuk menjamin keandalan mutlak:

### 10.1 Unit Test (Business Logic Isolation)
- Memvalidasi logika `ClaimReport` terhadap kembalian `VerifyReportLocation`.
- Menguji skenario `locationID == nil` (misal user super-admin atau migrasi lama).

### 10.2 Integration Test (Handler to DB Verification)
- Menjalankan request HTTP Fiber mock dengan konteks JWT nyata yang berisi klaim `location_id`.
- Memverifikasi join multi-tabel (`citizen_reports -> report_locations -> master_kelurahans`).

### 10.3 E2E / Business Scenario
- Skenario alur hidup lengkap: Warga lapor di Kelurahan Beji $\rightarrow$ Kader Beji login $\rightarrow$ Klaim tugas $\rightarrow$ Sukses $\rightarrow$ Kader Sukmajaya login $\rightarrow$ Coba klaim laporan Beji $\rightarrow$ Ditolak 403.

### 10.4 UAT-FIND-001 Test Matrix

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                             TEST EXECUTION MATRIX (UAT-FIND-001)                            │
├─────────────────────────┬───────────────────────────────────┬─────────────┬─────────────────┤
│ Skenario Uji            │ Kondisi Masukan                   │ Ekspektasi  │ Status Baseline │
├─────────────────────────┼───────────────────────────────────┼─────────────┼─────────────────┤
│ 1. Same Region Claim    │ Kader ID & Laporan di Faskes Sama │ Allowed 200 │ TBD (To Test)   │
│ 2. Cross Region Claim   │ Kader ID Faskes A, Laporan FaskesB│ Denied 403  │ TBD (To Test)   │
│ 3. Direct API Call      │ Curl POST langsung ke endpoint ID │ Denied 403  │ TBD (To Test)   │
│ 4. Unauthorized Role    │ User Role Petugas / Warga         │ Denied 403  │ TBD (To Test)   │
│ 5. Concurrent Claim     │ 2 Goroutine klaim laporan sama    │ 1 Win, 1 Err│ TBD (To Test)   │
└─────────────────────────┴───────────────────────────────────┴─────────────┴─────────────────┘
```

---

## 11. Regression Strategy

Untuk menjamin bahwa implementasi `IMP-001` tidak merusak fungsi yang sudah PASS, guardrail regresi berikut wajib dijalankan:

1. **Kader Normal Claim Workflow:** Memastikan kader tetap dapat mengklaim laporan di wilayah binaannya secara lancar tanpa hambatan baru.
2. **SLA Escalation Cron:** Memastikan cron eskalasi 24 jam (`backend/internal/delivery/cron/sla_cron.go`) tetap berfungsi memutasikan laporan tak bertuan menjadi `ESCALATED_TO_PUSKESMAS`.
3. **Puskesmas Takeover:** Memastikan Petugas Puskesmas tetap dapat mengambil alih laporan yang tereskalasi tanpa terhambat validasi klaim kader.
4. **Alur Penyelidikan Epidemiologi (PE):** Memastikan proses pengisian hasil PE dan foto bukti tidak terganggu.
5. **Appeal-to-PE State Machine:** Memastikan laporan hasil banding yang kembali ke status `MENUNGGU_PE_PUSKESMAS` tidak terdampak oleh guard klaim kader.
6. **Notification System:** Memastikan event notifikasi FCM ke kader lokal tetap terkirim secara tepat sasaran.

---

## 12. Migration Plan

### Group A (IMP-001)
> **NO DATABASE MIGRATION REQUIRED.**
Perbaikan otorisasi teritorial memanfaatkan skema dan relasi basis data yang telah ada (`master_kelurahans.puskesmas_id` dan `report_locations.kelurahan`). Tidak ada penambahan tabel, kolom, indeks, ataupun mutasi data historis.

### Group B (Tergantung Pilihan Stakeholder)
- Jika **DEC-01 (Option B)** dipilih:
  Memerlukan penambahan nilai `'DRAFT_PUSKESMAS'` pada constraint check enum status tabel `stock_requests`.
- Jika **DEC-02 (Option B atau D)** dipilih:
  Memerlukan migrasi penambahan kolom `is_emergency_override BOOLEAN DEFAULT FALSE` dan `override_reason TEXT` pada tabel `fogging_schedules`.
- Jika **DEC-03** dipilih:
  **No database migration required** (memanfaatkan kolom `updated_at` yang telah ada pada tabel `citizen_reports`).

---

## 13. Rollout Plan

Peluncuran implementasi dibagi ke dalam fase berjenjang (*phased deployment*):

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                   JADWAL FASE ROLLOUT                                       │
├─────────┬─────────────────────────────┬─────────────────────────────────────────────────────┤
│ Fase    │ Nama Fase                   │ Cakupan Pekerjaan & Gate Check                      │
├─────────┼─────────────────────────────┼─────────────────────────────────────────────────────┤
│ Phase A │ Fix UAT-FIND-001            │ Surgical fix pada handler & usecase backend         │
│ Phase B │ Automated Testing           │ Eksekusi Unit Test & Matrix Test UAT-FIND-001       │
│ Phase C │ Regression Suite            │ Eksekusi penuh `go test ./...` & E2E sanity check   │
│ Phase D │ Stakeholder Finalization    │ Penandatanganan ballot sheet kebijakan (DEC 1, 2, 3)│
│ Phase E │ Implement Chosen Policies   │ Eksekusi teknis kebijakan birokrasi terpilih        │
│ Phase F │ Final Operational UAT       │ Verifikasi terpadu 5 role pemangku kepentingan      │
│ Phase G │ Production Pilot            │ Pilot faskes percontohan (Puskesmas Beji & Pancoran)│
└─────────┴─────────────────────────────┴─────────────────────────────────────────────────────┘
```

---

## 14. Release Gates

Kriteria kelulusan formal sebelum sistem dinyatakan dapat beroperasi penuh di Kota Depok:

### Gate 1 — Technical Gate
- [ ] Build backend & frontend sukses tanpa warning kritis (*Zero Build Errors*).
- [ ] Automated unit & integration tests **PASS 100%**.
- [ ] Seluruh skenario pada *UAT-FIND-001 Test Matrix* berstatus **PASS**.
- [ ] Regression suite seluruh P0 invariants berstatus **PASS**.

### Gate 2 — Business Gate
- [ ] Alur operasional 5 peran utama terverifikasi tanpa *deadlock*.
- [ ] Konfirmasi aturan bisnis teritorial kader telah efektif di lapangan.
- [ ] Tidak ada anomali kebocoran data medis pelapor warga.

### Gate 3 — Governance Gate
- [ ] **DECISION-01** telah disahkan (Option A atau Option B) oleh Kadinkes Depok.
- [ ] **DECISION-02** telah disahkan (Option A, B, C, atau D) oleh Bidang P2P.
- [ ] **DECISION-03** telah disahkan (Option A, B, atau C) oleh tim hukum faskes.

### Gate 4 — Pilot Gate
- [ ] Dokumen SOP fisik dan petunjuk teknis telah didistribusikan ke Puskesmas percontohan.
- [ ] Penanggung jawab operasional faskes dan kontak PIC tanggap darurat telah ditetapkan.
- [ ] Prosedur pemulihan bencana (*Disaster Recovery*) dan rollback telah siap siaga.

---

## 15. Risk Traceability

Matriks penelusuran hubungan antara risiko audit, temuan UAT, keputusan bisnis, implementasi, dan pengujian:

| Identifikasi Risiko | Temuan UAT | Keputusan Tata Kelola | Item Rencana Implementasi | Metode Pengujian | Status Akhir |
|---|---|---|---|---|:---:|
| **RSK-SEC-04** (Bypass Otorisasi Wilayah) | `UAT-FIND-001` | Bug Teknis (Must Fix) | `IMP-001` (Jurisdiction Guard) | Unit Test + Direct API Call | **READY FOR IMPLEMENTATION** |
| **RSK-GOV-01** (Pelanggaran KPA Aset BMD) | `UAT-FIND-002` | `DECISION-01` | `IMP-002A` / `IMP-002B` | Usecase Authorization Test | **PENDING STAKEHOLDER DECISION** |
| **RSK-OPS-05** (Kelumpuhan Tanggap KLB) | `UAT-FIND-003` | `DECISION-02` | `IMP-003A/B/C/D` | Schedule Multi-Task Test | **PENDING STAKEHOLDER DECISION** |
| **RSK-LEG-02** (Ketidakpastian Hak Banding) | `UAT-FIND-004` | `DECISION-03` | `IMP-004A/B/C` | Repository Interval Filter Test | **PENDING STAKEHOLDER DECISION** |

---

## 16. Open Items

Daftar item yang masih berstatus terbuka dan memerlukan tindak lanjut:

| ID Item | Deskripsi Item Terbuka | Pihak Penanggung Jawab | Target Penyelesaian | Tindakan yang Dibutuhkan |
|---|---|---|---|---|
| **OPEN-01** | Pengesahan Lembar Ballot DEC-01 (Stok Obat) | Kadinkes & Seksi Farmasi | Menjelang Phase D | Memilih Option A atau Option B |
| **OPEN-02** | Pengesahan Lembar Ballot DEC-02 (Fogging KLB) | Kabid P2P & Koordinator Fogging | Menjelang Phase D | Memilih Option A, B, C, atau D |
| **OPEN-03** | Pengesahan Lembar Ballot DEC-03 (Batas Banding)| Bidang Surveilans & Hukum | Menjelang Phase D | Memilih Option A, B, atau C |

---

## 17. Final Implementation Readiness

### **STATUS KESIAPAN KESELURUHAN: PARTIALLY READY**

Rincian kesiapan per kelompok kerja:
1. **GROUP A (Confirmed Items — IMP-001):**
   **`READY FOR IMPLEMENTATION`**
   Seluruh spesifikasi teknis, rancangan perubahan minimal (*surgical fix*), lokasi file, kriteria penerimaan, dan skenario pengujian telah lengkap 100%. Tim pengembang dapat langsung mengeksekusi perbaikan kode untuk `IMP-001` segera setelah instruksi eksekusi diberikan.
2. **GROUP B (Pending Governance Decisions — DEC 01, 02, 03):**
   **`BLOCKED BY BUSINESS DECISION`**
   Implementasi kode pada kelompok ini **terkunci total** demi menjaga integritas kepatuhan hukum dan birokrasi Pemerintah Kota Depok, hingga lembar pengesahan ditandatangani secara resmi oleh pejabat berwenang.
