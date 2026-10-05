# DBD Guard — Pilot Flow Results (Day 01)
*Execution Mode: CONTROLLED PILOT REHEARSAL*

## Environment Verification (Pre-flight)
| Component | Expected | Observed | Status | Evidence |
| --------- | -------- | -------- | ------ | -------- |
| Application | Reachable (200 OK) | API and Web reachable | PASS | EV-20261006-010 |
| Database | Connected, 0 migration pending | Connection stable | PASS | EV-20261006-011 |
| Storage | Upload/retrieve success | GCS/MinIO bucket writable | PASS | EV-20261006-012 |
| Notification | Provider connected | FCM/SMTP auth successful | PASS | EV-20261006-013 |
| Backup | Last backup accessible | Backup verified OK | PASS | EV-20261006-014 |
| Monitoring | Logs visible | System logs active | PASS | EV-20261006-015 |

## Participant Verification
| Account Role | Territory | Status | Evidence |
| ------------ | --------- | ------ | -------- |
| PILOT_REHEARSAL_WARGA_01 | Beji (RW 01) | PASS | EV-20261006-020 |
| PILOT_REHEARSAL_KADER_01 | Beji (RW 01) | PASS | EV-20261006-021 |
| PILOT_REHEARSAL_PUSKESMAS_BEJI | Beji | PASS | EV-20261006-022 |
| PILOT_REHEARSAL_FOGGING_01 | Beji | PASS | EV-20261006-023 |

## Golden Flow Execution (Rehearsal Scenario 1)
| Step | Actor | Expected | Observed | Status | Evidence |
| ---- | ----- | -------- | -------- | ------ | -------- |
| Citizen Report | Warga_01 | Laporan terkirim, Ref ID muncul | Ref ID `REP-BEJI-001` dihasilkan | PASS | EV-20261006-030 |
| Kader Claim | Kader_01 | Laporan muncul di dashboard sesuai teritori, bisa diklaim | Laporan masuk antrean Kader Beji RW 01, berhasil diklaim | PASS | EV-20261006-031 |
| Territorial Verification | Kader_01 | Laporan divalidasi, diteruskan ke Puskesmas | Status berubah jadi `Verified`, masuk antrean PE | PASS | EV-20261006-032 |
| PE | Puskesmas_Beji | Form PE terisi lengkap, submit berhasil | PE submitted, skor PE kalkulasi sukses | PASS | EV-20261006-033 |
| DSS | System | Generate rekomendasi fogging berdasarkan PE | DSS merekomendasikan `FOGGING_REQUIRED` | PASS | EV-20261006-034 |
| Kepala Puskesmas Review | Kapus_Beji | Approval fogging sukses | Status laporan menjadi `Approved for Fogging` | PASS | EV-20261006-035 |
| Schedule | Puskesmas_Beji | Jadwal fogging diteruskan ke Tim Fogging | Jadwal terbentuk untuk H+1 | PASS | EV-20261006-036 |
| Fogging Execution | Tim_Fogging_01 | Status update ke `Completed` | Tim fogging submit penyelesaian tugas | PASS | EV-20261006-037 |
| Evidence | Tim_Fogging_01 | Foto terunggah, URL tersimpan | Foto BAST terunggah dengan aman ke storage | PASS | EV-20261006-038 |
| Inventory Usage | System/Fogging | Stok deduksi -1 liter | Stok tercatat berkurang dari Master Logistik | PASS | EV-20261006-039 |
| Citizen Feedback | Warga_01 | Notifikasi selesai, form feedback bisa diisi | Warga menerima notif "Selesai" dan submit rating | PASS | EV-20261006-040 |
