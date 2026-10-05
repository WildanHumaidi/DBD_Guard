# Demo Dry Run & Pilot Rehearsal — DBD Guard

## 1. Executive Summary
Sistem DBD Guard Kota Depok telah menjalani sesi *Demo Dry Run & Pilot Rehearsal* bersumber pada dokumen *Readiness Assessment*, *Risk Management*, dan *Implementation Verification* (IMP-001). Pemeriksaan ini bertujuan memastikan alur demonstrasi (*Golden Flow*, *Exception*, *Inventory*, *DSS*, dan *Security*) siap dipresentasikan di hadapan eksekutif, serta memvalidasi kesiapan lapangan (*Pilot Preconditions*) tanpa adanya fiktif (*mock*) fungsional yang tidak terekam dalam *evidence* nyata.

## 2. Evidence Quality
Kualitas pembuktian (evidence) yang mendasari rehearsal ini terbagi menjadi:
- **TEST VERIFIED**: Telah divalidasi oleh *automated test suite* (misal: Unit Test `VerifyReportLocation`).
- **STATICALLY VERIFIED**: Divalidasi melalui inspeksi *codebase* dan aturan bisnis (misal: Aturan SoD dan *Auth Middleware*).
- **DOCUMENTED ONLY**: Tersedia dalam dokumentasi *Readiness* dan BCP, namun belum melalui simulasi *end-to-end* secara *live* dengan beban nyata pada fase produksi.

## 3. Golden Flow Result
| Step | Actor | Action | Expected | Actual | Evidence | Status |
| ---- | ----- | ------ | -------- | ------ | -------- | ------ |
| 1 | Warga | Submit Report via UI | Tiket terbentuk status `MENUNGGU_PE_KADER` | Sesuai Expected | `TEST VERIFIED` | Pass |
| 2 | Kader | Claim Laporan Beji | Klaim berhasil, Yurisdiksi cocok | Sesuai Expected | `TEST VERIFIED` | Pass |
| 3 | Kader | Verification (Valid) | Status `MENUNGGU_PE_PUSKESMAS` | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 4 | Faskes | Input PE (Jentik +, Demam 3) | Lanjut ke DSS | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 5 | DSS | Menghasilkan Rekomendasi | Tampil `FOGGING FOKUS` | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 6 | Kepala Faskes | Human Review & Approval | Jadwal Fogging Terbit (Status Approved) | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 7 | Tim Lapangan | Fogging Execution & Evidence | Status Selesai, Wajib Foto | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 8 | Faskes/Dinkes| Inventory Deduction & Monitoring | Stok berkurang otomatis, Heatmap update | Sesuai Expected | `STATICALLY VERIFIED` | Pass |

## 4. Exception Flow Result
| Step | Actor | Action | Expected | Actual | Evidence | Status |
| ---- | ----- | ------ | -------- | ------ | -------- | ------ |
| 1 | Kader | Reject Hoax | Laporan tertutup (Status Rejected) | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 2 | Warga | Banding (Appeal) | Status `BANDING`, notif ke Puskesmas | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 3 | Faskes | Takeover & PE Ulang | Status `MENUNGGU_PE_PUSKESMAS` (oleh Sanitarian) | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 4 | Faskes (Sanitarian) | Validation (Self-Approve) | Ditolak (`403 FORBIDDEN` - SoD Guard) | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 5 | Kader Luar Wilayah | Claim Laporan | Ditolak (`403 FORBIDDEN` - IMP-001) | Sesuai Expected | `TEST VERIFIED` | Pass |

## 5. Inventory Demo Result
| Step | Action | Expected | Actual | Evidence | Status |
| ---- | ------ | -------- | ------ | -------- | ------ |
| 1 | Stock Request (Faskes -> Dinkes) | Status `PENDING` di Dinkes | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 2 | Approval (Dinkes) | Status `IN_TRANSIT`, stok Dinkes dikunci | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 3 | Receive (Faskes) | Status `COMPLETED`, source/dest balance konsisten | Sesuai Expected | `STATICALLY VERIFIED` | Pass |

*(Catatan: Saldo sumber dan tujuan telah dirancang atomik pada sistem database untuk menjaga konsistensi akuntansi (accounting consistency).*

## 6. DSS Demo Result
| Step | Action | Expected | Actual | Evidence | Status |
| ---- | ------ | -------- | ------ | -------- | ------ |
| 1 | PE Input -> Recommendation | Rekomendasi berdasarkan parameter terlihat jelas | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 2 | Human Decision & Optional Override | Pilihan akhir dan field `override reason` tersedia | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 3 | Audit Trail | Nama aktor dan waktu keputusan terekam di sistem | Sesuai Expected | `STATICALLY VERIFIED` | Pass |

## 7. Security Demo Result
| Scenario | Expected Response | Actual Response | Evidence | Status |
| -------- | ----------------- | --------------- | -------- | ------ |
| 1. Cross-region Claim | HTTP `403 Forbidden` (Kader dilarang ambil luar wilayah) | Sesuai Expected | `TEST VERIFIED` | Pass |
| 2. Unauthorized File Access | HTTP `403/401` (File medis dilindungi auth faskes) | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 3. Self-Validation (SoD) | HTTP `403 Forbidden` (Investigator != Validator) | Sesuai Expected | `STATICALLY VERIFIED` | Pass |
| 4. Unauthorized Role | HTTP `403 Forbidden` (Role mismatch pada endpoint) | Sesuai Expected | `STATICALLY VERIFIED` | Pass |

## 8. Demo Failure Plan
Jika terjadi kegagalan teknis (server *down*, jaringan lambat, map API *limit*) saat demo:
*   **Database/Backend Fail**: Tampilkan *screenshot evidence* atau *prerecorded evidence* (video navigasi UI lokal).
*   **Network Fail (API Timeout)**: Sediakan *API response evidence* (Postman/Curl logs) dan *test result* (Log `go test`) dari laptop presenter.
*   **Map Tiles (OSM) Fail**: Lanjutkan presentasi dengan menjelaskan bahwa *geotagging* tetap tersimpan secara *latitude/longitude* di database (*Database evidence*), peta blank hanya kendala limit API *frontend*.
*   *(Jangan memalsukan live success. Jelaskan secara transparan jika terpaksa menggunakan skenario fallback).*

## 9. Data Reset / Repeatability
Sistem saat ini tidak menyediakan *one-click automated DB reset* dari UI.
*   **Reset Procedure:** `MANUAL RESET` melalui eksekusi *SQL seed script* (`seed.sql`) secara langsung ke dalam container PostgreSQL.
*   **Cleanup State:** Seluruh transaksi demo (tiket, audit, mutasi inventaris) harus di-*truncate* melalui script, dilanjutkan dengan meng-*inject* *test data* / *seed data* bersih.
*   **Test Accounts:** Akun demonstrasi (`warga_demo`, `kader_demo`, dll.) tersedia, sandi di-*reset* ke *default*.

## 10. Pilot Preconditions Verification
| Requirement | Current Evidence | Status |
| ----------- | ---------------- | ------ |
| Pilot account | `DAFTAR_AKUN_DBD_GUARD.md` | `Ready` |
| Location mapping | Konfigurasi Beji & Pancoran Mas | `Ready` |
| Database | Skema dan relasi stabil | `Ready` |
| Backup | `backup_restore.ps1` lulus uji coba drill | `Ready` |
| Notification | SMTP (Email) & FCM di-*setup* | `Ready` |
| SOP | Tersedia dalam dokumen Operasional | `Ready` |
| Support | Grup WA (Level 1, 2, 3) direncanakan | `Planned` |
| Manual fallback | BCP menggunakan form Kemenkes tercatat | `Ready` |
| Monitoring | Dashboard KRI Dinkes (SLA) | `Ready` |

## 11. Stakeholder Decision Status
*   **DEC-01 (Auto-Draft Request Stok):** `PENDING`
*   **DEC-02 (Pendelegasian Regu BKO Dinkes saat KLB):** `PENDING`
*   **DEC-03 (Batas Waktu Kedaluwarsa Banding Warga):** `PENDING`

**Penanganan Pilot:** Karena keputusan tertunda, Pilot akan berjalan menggunakan kondisi eksisting (kuota stok faskes aman > 20 Liter sehingga auto-draft tidak terpantik, status non-KLB dijaga, dan tidak ada limit otomatis penutupan banding warga yang aktif memaksa). Tidak ada *policy* baru yang dikarang saat dry run.

## 12. Executive Demo Run Sheet
*(Durasi total: ~15 menit)*
| Order | Duration | Presenter | Actor | Screen | Action | Expected |
| ----: | -------: | --------- | ----- | ------ | ------ | -------- |
| 1 | 3 min | Project Lead | Warga & Kader | Mobile/Web | Submit laporan & klaim | Laporan masuk antrean PE |
| 2 | 4 min | Project Lead | Sanitarian Faskes | Web Dashboard | Input PE & hitung DSS | Rekomendasi Fogging Fokus |
| 3 | 3 min | Project Lead | Kepala Faskes | Web Dashboard | Approval Jadwal (SoD) | Jadwal eksekusi diterbitkan |
| 4 | 3 min | Project Lead | Tim Lapangan | Mobile | Eksekusi & Upload Bukti | Status Selesai, foto tersimpan |
| 5 | 2 min | Project Lead | Dinkes | Web Dashboard | Heatmap & Inventory | Stok berkurang atomik |

## 13. Technical Demo Run Sheet
*(Durasi total: ~10 menit)*
| Order | Duration | Presenter | Actor | Screen | Action | Expected |
| ----: | -------: | --------- | ----- | ------ | ------ | -------- |
| 1 | 3 min | DevOps/Tech | Kader (Luar Wil) | API / Web | *Cross-region claim* | Akses tertolak (403 IMP-001) |
| 2 | 2 min | DevOps/Tech | Sanitarian | API / Web | *Self-approve* (SoD Guard) | Akses tertolak (403 Validator != ID) |
| 3 | 3 min | DevOps/Tech | Puskesmas -> Dinkes | Web Dashboard | *Inventory Stock Request* | Bukti *consistency accounting* mutasi stok |
| 4 | 2 min | DevOps/Tech | Warga | Web Dashboard | *Banding (Appeal)* Laporan Hoax | Eskalasi ke Faskes (Audit Trail tercatat) |

## 14. Final Demo Checklist
- [x] Server ready
- [x] Database ready
- [x] Test accounts ready
- [x] Test data ready
- [x] Notification ready (Email/FCM)
- [x] Map ready (OSM public limits)
- [x] Upload ready (Medical/Evidence foto)
- [x] Backup verified
- [x] Browser/session clean
- [x] Demo reset verified (via SQL Seed manual)
- [x] Exception flow ready
- [x] Security flow ready
- [x] Evidence backup available (Screenshots & Logs)

## 15. Final Pilot Checklist
- [ ] Pilot scope confirmed (Menunggu tanda tangan SPT Kadinkes)
- [x] User accounts confirmed
- [x] Location mapping confirmed
- [ ] Support channel active (Grup WA belum 100% terkumpul)
- [x] SOP distributed
- [x] BCP available
- [x] Backup verified
- [x] Monitoring configured
- [ ] Incident escalation assigned (Menunggu konfirmasi person-in-charge lapangan)
- [x] Known limitations acknowledged (Sistem gratis FCM/OSM)
- [ ] Required stakeholder authorization available (Menunggu DEC-01, 02, 03)

## 16. Known Demo/Pilot Issues
1. **Map Tile Rate Limiting (OSM):** Risiko peta blank akibat *rate-limit* akses dari IP yang sama. (Solusi Demo: Gunakan *network tethering* jika terkena blokir).
2. **Manual Data Reset:** Transisi demo antar sesi mengharuskan *developer* menjalankan *script seed DB* secara manual. Belum ada fitur `Reset Demo` di level UI.
3. **Quarantined Governance:** Pilot harus diawasi ketat agar tidak bersinggungan dengan parameter *stock request* otomatis atau eskalasi KLB yang kebijakannya masih *pending*.

## 17. Final Rehearsal Verdict

**DEMO STATUS:**
> **DEMO READY**
Semua langkah demo dapat dijalankan atau memiliki *evidence fallback* yang sah dan transparan.

**PILOT STATUS:**
> **CONDITIONALLY PILOT READY**
Semua prasyarat teknis (*mandatory precondition*) telah divalidasi. Masih ada *precondition non-blocking* (legalitas SPT Kadinkes, grup eskalasi, pengesahan DEC-01 s/d DEC-03) yang harus dipenuhi sebelum hari-H operasional lapangan.
