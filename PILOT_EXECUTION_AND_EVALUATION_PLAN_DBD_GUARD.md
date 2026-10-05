# Pilot Execution & Evaluation Plan — DBD Guard

## 1. Pilot Objective
Menentukan bagaimana DBD Guard akan dijalankan, diamati, dicatat, dan dievaluasi selama limited pilot sehingga hasil pilot dapat digunakan sebagai dasar keputusan untuk melanjutkan, memperbaiki, atau memperluas penggunaan sistem. Pilot ini merupakan validasi nyata terhadap penggunaan sistem, proses kerja, adopsi pengguna, operasional, reliability, dan governance readiness dalam scope terbatas.

## 2. Pilot Scope
* Lokasi Pilot: TARGET PILOT — PENDING CONFIRMATION (Puskesmas Beji & Puskesmas Pancoran Mas - *menunggu konfirmasi stakeholder*)
* Durasi: PENDING CONFIRMATION
* Jumlah Pengguna: PENDING CONFIRMATION
* Jumlah Kelurahan/RW: PENDING CONFIRMATION
* Jumlah Laporan Target: PENDING CONFIRMATION
* Tanggal Mulai: DATE PENDING CONFIRMATION
* Tanggal Selesai: DATE PENDING CONFIRMATION
* PIC: PIC TO BE CONFIRMED

## 3. Entry Gate
Status: [PASS / CONDITIONAL / PENDING / BLOCKED]

### A. Governance
* [ ] Stakeholder authorization: PENDING
* [ ] SPT / approval bila diperlukan: PENDING
* [ ] DEC-01 status: GOVERNANCE DECISION REQUIRED
* [ ] DEC-02 status: GOVERNANCE DECISION REQUIRED
* [ ] DEC-03 status: GOVERNANCE DECISION REQUIRED
* [ ] PIC confirmation: PENDING

### B. Technical
* [ ] Application deployment: PENDING
* [ ] Database: PENDING
* [ ] Storage: PENDING
* [ ] Notification: PENDING
* [ ] Backup: PENDING
* [ ] Monitoring: PENDING
* [ ] Domain/network: PENDING
* [ ] Authentication: PENDING
* [ ] Accounts: PENDING

### C. Operational
* [ ] SOP tersedia: PENDING
* [ ] Escalation path: PENDING
* [ ] Support channel: PENDING
* [ ] Incident process: PENDING
* [ ] Manual fallback: PENDING
* [ ] Inventory reconciliation procedure: PENDING

### D. User Readiness
* [ ] Accounts created: PENDING
* [ ] Roles assigned: PENDING
* [ ] Territory assigned: PENDING
* [ ] Onboarding: PENDING
* [ ] Training: PENDING
* [ ] Credentials distributed: PENDING

### E. Data Readiness
* [ ] Master location: PENDING
* [ ] Puskesmas: PENDING
* [ ] Kader: PENDING
* [ ] Inventory: PENDING
* [ ] Fogging team: PENDING
* [ ] Test/demo data cleanup: PENDING

*Catatan: Pilot tidak boleh dinyatakan STARTED jika mandatory entry gate masih BLOCKED.*

## 4. Participant Matrix
| Role              | Expected Users | Confirmed Users | Training | Account | Territory | Status |
| ----------------- | -------------: | --------------: | -------- | ------- | --------- | ------ |
| Warga             | PENDING CONF.  | PENDING CONF.   | PENDING  | PENDING | N/A       | PENDING|
| Kader             | PENDING CONF.  | PENDING CONF.   | PENDING  | PENDING | PENDING   | PENDING|
| Petugas Puskesmas | PENDING CONF.  | PENDING CONF.   | PENDING  | PENDING | PENDING   | PENDING|
| Kepala Puskesmas  | PENDING CONF.  | PENDING CONF.   | PENDING  | PENDING | PENDING   | PENDING|
| Tim Fogging       | PENDING CONF.  | PENDING CONF.   | PENDING  | PENDING | PENDING   | PENDING|
| Dinkes            | PENDING CONF.  | PENDING CONF.   | PENDING  | PENDING | PENDING   | PENDING|
| Administrator     | PENDING CONF.  | PENDING CONF.   | PENDING  | PENDING | PENDING   | PENDING|

## 5. Golden Flow
Flow utama yang diamati selama pilot:
`Warga ↓ Citizen Report ↓ Reference ID ↓ Kader Claim ↓ Territorial Verification ↓ PE ↓ DSS ↓ Kepala Puskesmas Review ↓ Fogging Decision ↓ Schedule ↓ Fogging Execution ↓ Evidence ↓ Inventory Usage ↓ Citizen Feedback`

Untuk setiap tahap akan diobservasi:
* **Actor**: Siapa yang melakukan aksi.
* **Expected behavior**: Perilaku sistem yang diharapkan.
* **System behavior**: Perilaku sistem aktual (Observed).
* **Evidence**: Bukti transaksi (screenshot, log, db record).
* **Success criteria**: Kondisi berhasilnya aksi.
* **Failure handling**: Penanganan jika gagal.

## 6. Exception Scenarios
Skenario negatif yang wajib diuji/diamati (Trigger ↓ Expected ↓ Observed ↓ Evidence ↓ Incident? ↓ Resolution):

### Security
* Cross-region Kader claim
* Unauthorized PE validation
* Unauthorized fogging approval
* Unauthorized sensitive document access

### Business
* Rejected report
* Appeal
* Takeover
* Duplicate report
* Duplicate merge
* Insufficient stock
* Fogging postponed

### Operational
* Notification failure
* Storage failure
* API failure
* Database temporary failure
* User credential issue
* Cron failure

### EWS/KLB
* Threshold not reached
* Threshold reached
* Alert generated
* Alert recipient scope

## 7. Day-1 Runbook
### Before Start
* Verify environment, services, database, storage, notification, accounts, backup, monitoring.

### Start of Day
* Announce pilot start
* Confirm PIC
* Confirm support channel
* Verify pilot scope

### During Day
* Monitor system health, user activity
* Record incidents
* Collect user feedback
* Record deviations from SOP

### End of Day
* Daily reconciliation
* Incident review
* Unresolved issue review
* Data integrity check
* Backup check
* Pilot diary update

## 8. Daily Monitoring
| Metric               | Result | Target/Reference | Status | Evidence |
| -------------------- | -----: | ---------------- | ------ | -------- |
| Reports submitted    |        | PENDING CONF.    |        |          |
| Reports claimed      |        | PENDING CONF.    |        |          |
| Verification pending |        | PENDING CONF.    |        |          |
| SLA overdue          |        | PENDING CONF.    |        |          |
| PE submitted         |        | PENDING CONF.    |        |          |
| DSS generated        |        | PENDING CONF.    |        |          |
| Fogging approved     |        | PENDING CONF.    |        |          |
| Fogging completed    |        | PENDING CONF.    |        |          |
| Feedback received    |        | PENDING CONF.    |        |          |
| Failed notifications |        | PENDING CONF.    |        |          |
| Incidents            |        | PENDING CONF.    |        |          |
| Inventory variance   |        | PENDING CONF.    |        |          |

## 9. User Adoption
Pilot harus mengevaluasi penggunaan sistem tanpa ketergantungan berlebihan pada developer via observasi, wawancara terstruktur, survei, tiket support, dan incident record.
* **Kader**: Mudah memahami verification, mudah melakukan claim, memahami territorial restriction, memahami EWS.
* **Puskesmas**: Mudah melakukan PE, memahami DSS, memahami human approval, memahami inventory.
* **Kepala Puskesmas**: Memahami recommendation, memahami override, memahami accountability.
* **Fogging Team**: Mudah melihat assignment, mudah update status, mudah upload evidence, memahami usage input.
* **Warga**: Mudah membuat report, memahami tracking, memahami feedback, memahami anonymous reporting.

## 10. Incident Management
Log insiden dicatat terpisah. Kategori insiden: Technical, Business, Security, Data, Operational, Usability, Governance. Klasifikasi insiden: `SYSTEM BUG`, `USER ERROR`, `PROCESS GAP`, `TRAINING GAP`, `CONFIGURATION ISSUE`, `GOVERNANCE ISSUE`, `EXTERNAL FAILURE`.

## 11. Change Control
Setiap perubahan di tengah pilot harus melewati flow: `Issue ↓ Impact assessment ↓ Decision ↓ Change ↓ Test ↓ Deployment ↓ Verification`. Critical fix wajib menyertakan evidence regression.
| Change ID | Reason | Risk | Approval | Test | Deployment | Result |
| --------- | ------ | ---- | -------- | ---- | ---------- | ------ |

## 12. Data Integrity Monitoring
Selama pilot, monitor: Orphan records, duplicate processing, incorrect status transition, inventory mismatch, duplicate inventory deduction, incorrect ownership, incorrect territory, missing audit log, failed notification references, storage reference mismatch.

## 13. Inventory Reconciliation
(System Stock vs Physical Stock)
* Item, Location, System quantity, Physical quantity, Difference, Cause, Evidence, Authorization, Correction.
* Frekuensi rekonsiliasi: OPERATIONAL DECISION REQUIRED.

## 14. Security Monitoring
Monitor (Status: Observed / Normal baseline / Needs Investigation):
* Failed authentication, 403 spikes, 429 spikes, suspicious report volume, unauthorized access attempts, sensitive document access, anonymous PII exposure attempts, unusual inventory activity, abnormal account behavior.

## 15. EWS/KLB Validation
`Report aggregation ↓ Threshold evaluation ↓ EWS/KLB event ↓ Notification / recipient resolution ↓ Human response ↓ Operational action`
Jika threshold tidak tercapai selama pilot, statusnya adalah `NOT OBSERVED DURING PILOT` (bukan FAILED). Threshold: THRESHOLD GOVERNANCE DECISION REQUIRED.

## 16. Notification Monitoring
Track: Notification intended, recipient, provider, send attempt, retry, success/failure, final status (Created, Sent, Delivered). Jika provider tidak menyediakan delivery confirmation, tidak boleh mengklaim delivered.

## 17. User Feedback
Mekanisme pengumpulan feedback dari tiap role (Warga, Kader, dsb.) dengan pertanyaan minimal:
1. Apakah workflow mudah dipahami?
2. Apakah sistem mengurangi pekerjaan manual?
3. Apakah informasi yang diperlukan tersedia?
4. Bagian mana yang paling menghambat?
5. Apakah ada ketergantungan pada manual/WhatsApp?
6. Apa kesalahan yang paling sering terjadi?
7. Apakah sistem dipercaya untuk workflow operasional?

## 18. End Gate
Checklist sebelum pilot ditutup (Technical, Business, Operational, Security, Governance).

## 19. Evaluation Scorecard
Skor akhir: PASS / CONDITIONAL / NEEDS IMPROVEMENT / FAIL.
| Dimension         | Evidence | Status | Finding |
| ----------------- | -------- | ------ | ------- |
| Technical         |          |        |         |
| Business Workflow |          |        |         |
| Security          |          |        |         |
| Data Integrity    |          |        |         |
| Usability         |          |        |         |
| Operations        |          |        |         |
| Governance        |          |        |         |
| Adoption          |          |        |         |

## 20. Exit Decision Criteria
Hasil akhir: CONTINUE / CONTINUE WITH IMPROVEMENTS / EXTEND PILOT / PAUSE / STOP.

## 21. Wider Deployment Assessment
Evaluasi kesiapan citywide (Technical, operational, governance, infrastructure capacity, user adoption, incident history, data quality, support capacity, inventory readiness).
Status: NOT RECOMMENDED / CONDITIONALLY RECOMMENDED / RECOMMENDED / NOT YET ASSESSABLE.

## 22. Evidence Collection
Setiap hasil wajib didukung bukti yang dapat diaudit (screenshot, API response, dll) dengan klasifikasi: RUNTIME VERIFIED, TEST VERIFIED, DOCUMENTED, USER OBSERVED, NOT OBSERVED.

## 23. Remaining Governance Dependencies
DEC-01, DEC-02, DEC-03: GOVERNANCE DECISION REQUIRED.

## 24. Pilot Conclusion
Menyimpulkan apakah pengguna dapat menjalankan workflow, sistem stabil, operasional sesuai SOP, masalah nyata teridentifikasi, dan kecukupan evidence untuk wider deployment.
