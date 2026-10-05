# IMP-001 Implementation & Verification Report

**Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi (Kota Depok)**  
**Item:** IMP-001 / UAT-FIND-001 — Kader Cross-Region Task Claim Authorization  
**Tanggal Eksekusi:** 01 Oktober 2026  
**Status Akhir:** **RESOLVED**  

---

## 1. Summary

Implementasi **IMP-001** berhasil diselesaikan sesuai spesifikasi dokumen [IMPLEMENTATION_PLAN_DBD_GUARD.md](file:///e:/Projects/DBD_Guard/IMPLEMENTATION_PLAN_DBD_GUARD.md) dan panduan tata kelola bisnis DBD Guard Kota Depok. 

Perbaikan ini menegakkan batas otorisasi wilayah teritorial (*territorial jurisdiction boundary*) secara mutlak di lapisan *Backend Core* (`report_handler.go` & `report_usecase.go`). Sistem kini memverifikasi relasi teritorial antara wilayah laporan (`report_locations.kelurahan`) dan wilayah wewenang kader (`location_id` pada JWT context) menggunakan repositori `VerifyReportLocation`. 

Dengan perbaikan ini, celah keamanan bisnis **UAT-FIND-001** (Kader dapat mengklaim tugas di luar wilayah Puskesmas/Kelurahan binaan melalui panggilan API langsung) telah ditutup secara permanen dengan respon **HTTP 403 Forbidden**, tanpa mengubah skema basis data, tanpa merusak atomisitas klaim konkuren, dan tanpa menyebabkan regresi pada alur kerja lain.

---

## 2. Original Finding

* **ID Temuan:** `UAT-FIND-001` / `RSK-SEC-04`
* **Klasifikasi:** **P1 BUSINESS / SECURITY DEFECT**
* **Deskripsi:**
  Endpoint `POST /api/v1/kader/reports/:id/claim` memvalidasi status laporan (`MENUNGGU_PE_KADER`) dan memastikan `investigator_id` masih kosong, namun **tidak memverifikasi apakah laporan yang diklaim berada dalam wilayah kerja Puskesmas/Kelurahan yang menjadi wewenang Kader**.
  Filter wilayah sebelumnya hanya diterapkan di antarmuka frontend (`GetReports`), sehingga seorang Kader dapat memotong batasan wilayah dengan mengirimkan request HTTP langsung (via cURL, Postman, atau skrip) untuk mengklaim tugas di Puskesmas lain.

---

## 3. Root Cause

1. **Missing Business Authorization Guard:**
   Pada fungsi `ClaimReport(ctx, id, kaderID)` di usecase `report_usecase.go`, parameter `locationID` tidak diikutsertakan dalam kontrak fungsi, dan tidak ada pemanggilan ke helper otorisasi lokasi yang telah tersedia di repositori.
2. **Untrusted Frontend Boundary Assumption:**
   Otorisasi wilayah hanya bergantung pada pemfilteran tampilan daftar laporan pada frontend web, yang bukan merupakan *security boundary*.
3. **Repository Typo:**
   Kueri verifikasi wilayah `VerifyReportLocation` di `report_repository.go` memiliki ketidaksesuaian nama kolom (`mk.kelurahan_id` alih-alih `mk.id`), yang langsung diperbaiki secara bedah (*surgical*) sehingga kueri integrasi PostgreSQL berjalan 100% valid.

---

## 4. Implemented Change

Perbaikan dilakukan mengikuti prinsip **Minimal → Surgical → Traceable → Testable → Reversible**:

1. **Domain Layer ([report.go](file:///e:/Projects/DBD_Guard/backend/internal/domain/report.go)):**
   Menambahkan definisi error domain standar:
   - `domain.ErrUnauthorizedJurisdiction`: `"akses ditolak: laporan di luar yurisdiksi wilayah kerja"`
   - `domain.ErrMissingLocationContext`: `"akses ditolak: akun tidak memiliki wilayah kerja yang valid"`

2. **Usecase Layer ([report_usecase.go](file:///e:/Projects/DBD_Guard/backend/internal/usecase/report_usecase.go)):**
   - Mengubah tanda tangan method interface `ReportUsecase`:
     ```go
     ClaimReport(ctx context.Context, id uuid.UUID, kaderID uuid.UUID, locationID *uuid.UUID) error
     ```
   - Menambahkan guard validasi:
     1. Menolak jika `locationID == nil` (`ErrMissingLocationContext`).
     2. Memverifikasi yurisdiksi via `u.reportRepo.VerifyReportLocation(ctx, id, *locationID)`. Jika tidak cocok, tolak dengan `ErrUnauthorizedJurisdiction`.
     3. Mempertahankan invariant status `MENUNGGU_PE_KADER` dan `investigator_id IS NULL`.
     4. Menjalankan mutasi atomik `ClaimReportAtomic(ctx, id, kaderID)`.
     5. Mencatat audit trail `CLAIM_REPORT` hanya jika klaim atomik sukses.

3. **HTTP Delivery Layer ([report_handler.go](file:///e:/Projects/DBD_Guard/backend/internal/delivery/http/report_handler.go)):**
   - Mengekstrak `locationID, _ := middleware.GetLocationID(c)` dari konteks JWT tepercaya (bukan dari request body).
   - Meneruskan `locationID` ke usecase `ClaimReport`.
   - Memetakan `domain.ErrUnauthorizedJurisdiction` dan `domain.ErrMissingLocationContext` ke respon **HTTP 403 Forbidden**:
     ```json
     {
       "error": "Akses ditolak: Laporan berada di luar yurisdiksi wilayah kerja Anda",
       "message": "akses ditolak: laporan di luar yurisdiksi wilayah kerja"
     }
     ```
   - Memetakan error konflik status/klaim ke **HTTP 409 Conflict**.

4. **Repository Layer ([report_repository.go](file:///e:/Projects/DBD_Guard/backend/internal/repository/report_repository.go)):**
   - Memperbaiki klausa kolom pada `VerifyReportLocation` dari `(mk.puskesmas_id = ? OR mk.kelurahan_id = ?)` menjadi `(mk.puskesmas_id = ? OR mk.id = ?)`, menyelaraskan dengan skema entitas `MasterKelurahan`.

5. **Caller Updates:**
   - Memperbarui `cmd/e2e/main.go` dan `bug08_claim_race_test.go` untuk menyertakan `locationID`.

---

## 5. Files Changed

| File | Tipe Perubahan | Deskripsi Perubahan |
|---|---|---|
| `backend/internal/domain/report.go` | Modified | Menambahkan `ErrUnauthorizedJurisdiction` dan `ErrMissingLocationContext` |
| `backend/internal/usecase/report_usecase.go` | Modified | Menambahkan parameter `locationID`, pengecekan teritorial, dan audit trail |
| `backend/internal/delivery/http/report_handler.go` | Modified | Ekstraksi `GetLocationID(c)` dari JWT dan pemetaan respon HTTP 403 |
| `backend/internal/repository/report_repository.go` | Modified (Surgical) | Koreksi nama kolom `mk.kelurahan_id` -> `mk.id` pada `VerifyReportLocation` |
| `backend/internal/usecase/bug08_claim_race_test.go` | Modified | Penyesuaian pemanggilan `ClaimReport` dengan `locationID` dan mock helper |
| `backend/cmd/e2e/main.go` | Modified | Penyesuaian pemanggilan `ClaimReport` dengan `kader.LocationID` |
| `backend/internal/repository/report_repository_test.go` | Modified | Menambahkan pengujian integrasi database `TestVerifyReportLocation` |
| `backend/internal/usecase/imp001_kader_claim_test.go` | Added (New) | Suite unit test untuk seluruh skenario bisnis IMP-001 |
| `backend/internal/delivery/http/imp001_kader_claim_handler_test.go` | Added (New) | Suite HTTP integration test untuk seluruh skenario API endpoint |

---

## 6. API Behavior

### Endpoint
`POST /api/v1/kader/reports/:id/claim`

### Headers
* `Authorization: Bearer <JWT_TOKEN>` (Wajib, memuat `user_id`, `role: "Kader"`, dan `location_id`)

### Response Matrix

#### 1. Same Region Claim (Success)
* **Status:** `200 OK`
* **Body:**
  ```json
  {
    "message": "Berhasil mengklaim tugas laporan"
  }
  ```

#### 2. Cross Region Claim (Forbidden)
* **Status:** `403 Forbidden`
* **Body:**
  ```json
  {
    "error": "Akses ditolak: Laporan berada di luar yurisdiksi wilayah kerja Anda",
    "message": "akses ditolak: laporan di luar yurisdiksi wilayah kerja"
  }
  ```

#### 3. Missing Location Context (Forbidden)
* **Status:** `403 Forbidden`
* **Body:**
  ```json
  {
    "error": "Akses ditolak: Laporan berada di luar yurisdiksi wilayah kerja Anda",
    "message": "akses ditolak: akun tidak memiliki wilayah kerja yang valid"
  }
  ```

#### 4. Unauthorized Role (Forbidden by RBAC)
* **Status:** `403 Forbidden`
* **Body:**
  ```json
  {
    "error": "Access denied: insufficient permissions"
  }
  ```

#### 5. Unauthenticated Request
* **Status:** `401 Unauthorized`
* **Body:**
  ```json
  {
    "error": "Missing authorization header"
  }
  ```

---

## 7. Test Matrix

Seluruh skenario pengujian yang diwajibkan telah diimplementasikan dan diverifikasi:

| Test ID | Skenario | Kondisi Masukan | Ekspektasi | Hasil Aktual | Status |
|:---:|---|---|---|---|:---:|
| **TEST-01** | Same Region Claim | Kader Location = X, Report Location = X | HTTP 200, `investigator_id` terisi, audit trail tercatat | Status 200 OK, `investigator_id` terisi Kader, 1 audit log | **PASS** |
| **TEST-02** | Cross Region Claim | Kader Location = X, Report Location = Y | HTTP 403, `investigator_id` tidak berubah, tanpa audit trail | Status 403 Forbidden, `investigator_id` tetap nil, 0 audit log | **PASS** |
| **TEST-03** | Direct API Bypass | Curl/Postman langsung ke endpoint ID luar wilayah | HTTP 403 Forbidden | Status 403 Forbidden, mutasi dibatalkan | **PASS** |
| **TEST-04** | Unauthorized Role | Akun Petugas Puskesmas / Warga panggil claim | HTTP 403 Forbidden (RBAC) | Status 403 Forbidden, middleware menolak | **PASS** |
| **TEST-05** | Concurrent Claim | 2 Goroutine Kader wilayah sama klaim bersamaan | Tepat 1 menang (HTTP 200), 1 kalah (Konflik) | Tepat 1 Sukses, 1 Gagal, tidak ada korupsi data | **PASS** |
| **TEST-06** | Missing Location Context | Token Kader tanpa klaim `location_id` | HTTP 403 Forbidden secara aman | Status 403 Forbidden, tidak ada fallback permissif | **PASS** |
| **TEST-07** | Invalid Report Status | Laporan berstatus `COMPLETED` / `APPEALED` | Error ditolak | Klaim ditolak, tidak ada perubahan investigator | **PASS** |
| **TEST-08** | Already Claimed Report | Laporan sudah memiliki `investigator_id` | Error ditolak | Klaim ditolak, investigator awal tidak tertimpa | **PASS** |
| **TEST-09** | Live PostgreSQL Join | Verifikasi integrasi PostGIS & relasi kelurahan | `VerifyReportLocation` match Puskesmas & Kelurahan | Kueri SQL berjalan akurat (True pada faskes sama, False pada faskes beda) | **PASS** |

---

## 8. Regression Test Result

Setelah implementasi IMP-001, seluruh suite pengujian otomatis dijalankan untuk memverifikasi tidak adanya regresi:

### Backend Test Suite (`go test -v ./...`)
* **Package `internal/config`:** **PASS** (0.15s)
* **Package `internal/delivery/http`:** **PASS** (0.20s — 7/7 HTTP scenarios passed)
* **Package `internal/repository`:** **PASS** (0.28s — Real DB integration verified)
* **Package `internal/usecase`:** **PASS** (0.33s — 6/6 Usecase scenarios passed)
* **P0 Invariants Status:**
  - `TestP0_01_InventoryAccountingIntegrity`: **PASS** (Zero leaks, strict conservation)
  - `TestP0_02_SensitiveFileStorageProtection`: **PASS** (Photo access permission scope)
  - `TestP0_03_and_04_DSS_and_SoD`: **PASS** (SoD barrier intact)
  - `TestP0_06_AppealToPELifecycle`: **PASS** (Full lifecycle REJECTED_HOAX -> APPEALED -> ESKALASI -> TAKEOVER -> SUBMIT_PE -> VALIDATE intact)
  - `TestBUGNEW008_ClaimReport_AtomicRaceCondition`: **PASS** (Atomic claim preserved)
  - `TestBUGREG001_FullLifecycle_AppealToPEValidation`: **PASS**

### Frontend Web Build (`npm run build`)
* **Next.js Version:** 16.3.4 (Turbopack)
* **Route Generation:** 32/32 routes compiled successfully
* **TypeScript Compilation:** Finished in 4.7s with **Zero Type Errors**
* **Exit Code:** `0`

---

## 9. Database / Schema Impact

* **Dampak Skema:** **NOL (Zero Database Migration Required).**
* Tidak ada penambahan tabel, kolom, index, ataupun modifikasi DDL.
* Kueri otorisasi memanfaatkan indeks relasi yang telah ada pada tabel `report_locations` dan `master_kelurahans` (`puskesmas_id` dan `id`).

---

## 10. Security Impact

* **Penutupan Kerentanan BOLA / IDOR:**
  Menutup celah *Broken Object Level Authorization* (OWASP API1:2023) pada mutasi penugasan kader. Kader tidak lagi dapat memanipulasi parameter URL `:id` untuk mengambil alih tugas investigasi di luar wilayah binaan.
* **Defense in Depth:**
  Keamanan tidak lagi bergantung pada lapisan presentasi (frontend filtering), melainkan ditegakkan pada lapisan otorisasi domain aplikasi (*Backend Usecase Layer*).

---

## 11. Business Impact

* **Kepatuhan Teritorial Faskes:**
  Menjamin kepatuhan penuh terhadap Peraturan Daerah dan PRD Bab 3 terkait batas wilayah kerja operasional Puskesmas dan Kelurahan di Kota Depok.
* **Integritas Penyelidikan Epidemiologi (PE):**
  Hanya kader yang secara sah bertugas di wilayah pelaporan yang dapat menjadi `investigator_id` dan memasukkan data observasi jentik, fever count, dan edukasi 3M Plus.

---

## 12. Final Status

### **RESOLVED**

Perbaikan **IMP-001** telah selesai 100%, seluruh kriteria penerimaan terpenuhi, seluruh uji unit, konkurensi, integrasi HTTP, dan regresi berstatus **PASS**, serta tidak ada dependensi atau dampak pada item Group B (DEC-01, DEC-02, DEC-03).
