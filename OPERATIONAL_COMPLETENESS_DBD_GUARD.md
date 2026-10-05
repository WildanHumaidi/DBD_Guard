# Operational Completeness & Deployment Readiness — DBD Guard
**Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi (Pemerintah Kota Depok)**

- **Tahap Proyek:** BATCH 3 CLOSURE
- **Fokus Utama:** Operational Completeness & Governance Alignment
- **Status Dokumen:** FINAL DRAFT

---

## 1. Executive Summary

Dokumen **Operational Completeness & Deployment Readiness** ini menetapkan kerangka kerja operasional menyeluruh untuk menjamin sistem DBD Guard dapat dijalankan secara terstruktur, dipantau, dipertanggungjawabkan, dan dipulihkan apabila terjadi kendala teknis maupun prosedural. 

Tahap ini **bukan** modifikasi arsitektur atau penambahan fitur fungsional, melainkan pemetaan antara fungsi sistem yang telah matang (Production-Ready MVP) dengan prosedur birokrasi dan operasional nyata di lapangan.

**Keputusan Tata Kelola (Governance Decision):** 
Selama penyusunan SOP dan kesiapan operasional ini, seluruh kebijakan yang belum diputuskan (DEC-01, DEC-02, DEC-03) diperlakukan menggunakan protokol `GOVERNANCE DECISION REQUIRED`.

---

## 2. Role & PIC (Person In Charge) Matrix

Berikut adalah pemetaan kepemilikan dan tanggung jawab masing-masing entitas dalam siklus hidup DBD Guard sesuai struktur organisasi yang berlaku:

| Peran Sistem | PIC Aktual Lapangan | Tanggung Jawab Utama |
| :--- | :--- | :--- |
| **Warga (Public)** | Penduduk Kota Depok | Melaporkan indikasi kasus DBD dengan foto bukti nyata dan titik lokasi GIS. |
| **Kader** | Kader Posyandu/RW Setempat | Merespons tiket laporan $\le$ 72 jam, mengecek kondisi lapangan (verifikasi awal). |
| **Tim Fogging** | Petugas Lapangan Puskesmas | Menerima jadwal fogging, mengeksekusi di lokasi, mengunggah foto bukti eksekusi. |
| **Petugas Faskes** | Petugas Surveilans / Sanitarian | Melakukan Penyelidikan Epidemiologi (PE) komprehensif, input variabel DSS. |
| **Kepala Faskes** | Kepala UPTD Puskesmas | Menyetujui rekomendasi DSS, melakukan *Clinical Override* (SoD), KPA aset inventaris. |
| **Dinas Kesehatan** | PIC TO BE CONFIRMED | Memantau sebaran kasus (Heatmap), manajemen stok insektisida lintas faskes, EWS. |
| **Administrator** | PIC TO BE CONFIRMED | Manajemen akun pengguna, backup database, penanganan insiden sistem (L3). |

---

## 3. SOP Mapping (SOP-01 s/d SOP-19)

Pemetaan 19 Standar Operasional Prosedur (SOP) dengan fungsionalitas *codebase* DBD Guard:

### A. Alur Pelaporan & Verifikasi Kasus
*   **SOP-01 (Pelaporan Warga):** Warga mengirim laporan via portal publik. Warga dapat melampirkan foto atau bukti relevan sesuai field/form yang tersedia dan ketentuan sistem. Sistem menginisiasi tiket dengan status `MENUNGGU_PE_KADER`.
*   **SOP-02 (Klaim Tiket Kader):** Kader mengklaim tiket sesuai batasan otorisasi yurisdiksi wilayah RW-nya (IMP-001). Status berubah menjadi `MENUNGGU_PE_PUSKESMAS` jika valid.
*   **SOP-03 (Penolakan Laporan Hoax):** Kader atau Faskes berhak menolak laporan yang tidak valid atau palsu, mengubah status tiket menjadi `REJECTED_HOAX`.
*   **SOP-04 (Banding Warga):** Warga berhak memproses *appeal* untuk laporan `REJECTED_HOAX` (Maks 1 kali). **Status: DEC-03 PENDING (Temporal Limit GOVERNANCE DECISION REQUIRED)**. *Pilot operating assumption / proposed handling*: No expiry khusus selama pilot.

### B. Penyelidikan Epidemiologi & Intervensi
*   **SOP-05 (Pelaksanaan PE):** Petugas Sanitarian melaksanakan pengumpulan data PE (jentik, kasus demam $\ge$ 3) dan diinput ke sistem untuk kalkulasi Rekomendasi DSS.
*   **SOP-06 (Validasi DSS & SoD):** Kepala Puskesmas wajib menyetujui jadwal atau memberikan *Clinical Override*. Sistem menolak validasi jika dilakukan oleh pelaku PE yang sama (*Four-Eyes Principle*).
*   **SOP-07 (Eksekusi Fogging):** Tim lapangan wajib mengunggah foto bukti pelaksanaan di lokasi untuk mengubah status menjadi `COMPLETED`.
*   **SOP-08 (Pemotongan Inventaris):** Penyelesaian tugas fogging secara otomatis dan mutlak memotong stok inventaris insektisida secara atomik pada database Puskesmas.

### C. Manajemen Logistik & Kedaruratan
*   **SOP-09 (Permintaan Stok Antar Faskes):** Transaksi *Stock Transfer* wajib melewati status `IN_TRANSIT`. **Status: DEC-01 PENDING (Bypass ke Dinkes GOVERNANCE DECISION REQUIRED)**. *Temporary contingency option pending stakeholder approval*: Penanganan transfer antar faskes diselesaikan manual.
*   **SOP-10 (Monitoring Heatmap & EWS):** Tim Surveilans Dinkes menggunakan dasbor analitik harian untuk meninjau radius kerawanan dan potensi KLB.
*   **SOP-11 (Penetapan Wilayah KLB):** Penetapan KLB dan eskalasi EWS. **Status: DEC-02 PENDING (Override Quota KLB GOVERNANCE DECISION REQUIRED)**. *Emergency fallback option — GOVERNANCE DECISION REQUIRED*: Pengerahan > 1 armada fogging jika ditetapkan darurat (Response Capacity). KLB Detection dipisahkan dari Fogging Response Capacity.
*   **SOP-17 (Rekonsiliasi Inventaris Fisik):** Setiap selisih antara stok fisik dan stok sistem wajib melalui proses rekonsiliasi, investigasi, dokumentasi, dan otorisasi sesuai prosedur organisasi yang berlaku. BAST dapat digunakan apabila diwajibkan oleh SOP/ketentuan organisasi terkait.

### D. Keamanan, Sistem & Manajemen Akun
*   **SOP-12 (Manajemen Akun):** Account management dilakukan oleh authorized administrator sesuai RBAC sistem. Organizational PIC aktual tetap PIC TO BE CONFIRMED sampai ditetapkan stakeholder.
*   **SOP-13 (Pencadangan Database):** *Automated Cron Backup* (`pg_dump`) dieksekusi setiap malam. (Retensi: GOVERNANCE DECISION REQUIRED).
*   **SOP-14 (Pemulihan Bencana):** Prosedur restorasi *live* `backup_restore.ps1` untuk skenario korupsi data.
*   **SOP-15 (Reset Sandi):** Reset kata sandi administratif wajib divalidasi oleh IT Dinkes.
*   **SOP-16 (Oto-Rilis Tugas Kader):** Eksekusi *cron job* SLA SLA-Cron untuk membebaskan tiket kader yang mangkir $> 72$ jam.
*   **SOP-18 (Transisi Manual / BCP):** *Business Continuity Procedure / Proposed Fallback*. Saat peladen lumpuh, operasional beralih ke formulir PE rangkap 3 konvensional; re-input dilakukan pasca sistem hidup. (Detail re-entry, rekonsiliasi, pencegahan duplikat: GAP — OPERATIONAL DECISION REQUIRED).
*   **SOP-19 (Audit Jejak Digital):** Inspeksi `audit_logs` untuk verifikasi perselisihan wewenang (kapan, siapa, alamat IP).

---

## 4. Incident Management & Escalation Protocol

Insiden di lapangan harus dikelola menggunakan eskalasi 3-Tier. **INCIDENT SLA GOVERNANCE DECISION REQUIRED**.

1.  **Level 1 (Operasional Dasar):** 
    *   **PIC:** Operator Puskesmas / Sanitarian
    *   **Kasus:** Lupa sandi kader, revisi data laporan (typo), stok habis dilaporkan manual.
    *   **Tindakan:** Resolve langsung via dashboard UI faskes.
2.  **Level 2 (Kebijakan Khusus & Logistik):**
    *   **PIC:** PIC TO BE CONFIRMED (Dinkes / Kepala Puskesmas)
    *   **Kasus:** Keputusan transfer logistik antar-Puskesmas, sengketa hasil PE antar-RW, pengerahan BKO KLB darurat.
    *   **Tindakan:** Evaluasi dan eksekusi via Dashboard Dinkes atau perintah fisik.
3.  **Level 3 (Teknis & Kegagalan Sistem):**
    *   **PIC:** PIC TO BE CONFIRMED (Tim Pengembang / Diskominfo)
    *   **Kasus:** Database terkorupsi, kegagalan sinkronisasi notifikasi FCM secara luas, *Downtime* server.
    *   **Tindakan:** Menjalankan **SOP-14 (Disaster Recovery)** dan aktivasi **SOP-18 (BCP Manual)**.

---

## 5. Monitoring & Key Risk Indicators (KRI)

### A. Existing / Evidence-Based KRI

| KRI | Source | Current Measurement | Official Threshold | Status |
| --- | ------ | ------------------- | ------------------ | ------ |
| SLA Pelanggaran | PRD (SLA 3x24 Jam) | Tiket gagal direspons kader > 72 jam | 72 Jam | APPROVED |
| Sistem Notifikasi | Technical Requirement | Bounced emails / undelivered push | N/A | APPROVED |

### B. Proposed KRI

| Proposed KRI | Purpose | Data Source | Suggested Review | Approval Status |
| ------------ | ------- | ----------- | ---------------- | --------------- |
| Breach Rate Target | Kualitas respons kader | SLA_OVERDUE | < 10% | THRESHOLD GOVERNANCE DECISION REQUIRED |
| Hoax Rate | Deteksi penyalahgunaan sistem | REJECTED_HOAX | > 30% | THRESHOLD GOVERNANCE DECISION REQUIRED |
| Low Stock Alert | Mencegah kekosongan logistik | Inventaris Puskesmas | < 20 Liter | THRESHOLD GOVERNANCE DECISION REQUIRED |
| Appeal Bottleneck | Mengukur efisiensi banding | Status APPEAL | > 3 Hari | THRESHOLD GOVERNANCE DECISION REQUIRED |

---

## 6. Backup & Disaster Recovery (DR) Procedures

*   **Prosedur Harian (Backup):** Skrip utilitas terjadwal mencadangkan Postgres volume (Technical fact). Data cadangan harus diverifikasi (Operational policy - Frekuensi: GOVERNANCE DECISION REQUIRED). Retensi dan RPO: GOVERNANCE DECISION REQUIRED.
*   **Prosedur Pemulihan (Restoration):** Dalam kasus bencana data, administrator wajib:
    1.  Membekukan trafik aplikasi (Maintenance Mode / matikan container web).
    2.  Mengeksekusi skrip `backup_restore.ps1` dari titik *snapshot* terakhir.
    3.  Menginstruksikan faskes untuk melakukan rekonsiliasi selisih input manual sejak waktu _snapshot_.

---

## 7. Pilot Readiness (Limited Pilot)

**Target Pilot:** Puskesmas Beji & Puskesmas Pancoran Mas (PENDING CONFIRMATION).
**Durasi Uji Coba Terbatas:** PROPOSED PILOT DURATION (GOVERNANCE DECISION REQUIRED).

1.  **Fokus Uji Coba:**
    *   Uji stres *real-world* untuk *Jurisdiction Guard* (IMP-001).
    *   Kesiapan faskes mengadopsi mekanisme *Four-Eyes Principle* (SoD).
    *   Adaptasi warga terhadap mekanisme banding (*Appeal Lifecycle*).
2.  **Mitigasi Isu Tata Kelola (The "Gated" Policies):**
    *   **DEC-01 (Stok):** *Temporary contingency option pending stakeholder approval*: Ditangani via BAST manual sementara (BUKAN official bypass).
    *   **DEC-02 (KLB Quota):** *Emergency fallback option — GOVERNANCE DECISION REQUIRED*: Dinkes *takeover* manual jika KLB.
    *   **DEC-03 (Waktu Banding):** *Pilot operating assumption*: Dibiarkan terbuka (*no expiration*) khusus selama pilot (GOVERNANCE DECISION REQUIRED).
3.  **Dokumentasi Wajib (Go-Live Gate):**
    *   Surat Perintah Tugas (SPT) Uji Coba dari Kepala Dinas Kesehatan Kota Depok (PENDING CONFIRMATION).
    *   Akun rintisan (*seed accounts*) didistribusikan kepada PIC Beji & Pancoran Mas (PENDING CONFIRMATION).

---

## 8. Governance Decisions Required

| ID     | Decision            | Current State | Operational Impact   | Blocking?       |
| ------ | ------------------- | ------------- | -------------------- | --------------- |
| DEC-01 | Stock Transfer Flow | Pending       | Inventory operations | Yes/Conditional |
| DEC-02 | KLB Capacity        | Pending       | Emergency response   | Yes/Conditional |
| DEC-03 | Appeal Expiration   | Pending       | Appeal workflow      | Yes/Conditional |

---

## 9. Final Operational Readiness Matrix

| Area                     | Status | Evidence | Remaining Dependency |
| ------------------------ | ------ | -------- | -------------------- |
| SOP                      | CONDITIONAL | Dokumen SOP-01 - 19 | Pengesahan Stakeholder |
| PIC                      | CONDITIONAL | Pemetaan Sistem | Konfirmasi Organisasi (PIC TO BE CONFIRMED) |
| Incident                 | CONDITIONAL | Incident 3-Tier | INCIDENT SLA GOVERNANCE DECISION REQUIRED |
| Escalation               | CONDITIONAL | Incident 3-Tier | Pengesahan Stakeholder |
| Monitoring               | CONDITIONAL | PRD | Penetapan Threshold KRI |
| Backup                   | CONDITIONAL | Technical Skrip | Retention & RPO Policy |
| DR                       | CONDITIONAL | Technical Skrip | RTO Policy |
| BCP                      | CONDITIONAL | SOP-18 | Detail Proses Manual (Trigger, Re-entry) |
| Inventory Reconciliation | CONDITIONAL | SOP-17 | DEC-01 Approval |
| Fogging Operations       | CONDITIONAL | PRD & SOP | Persetujuan Operasional Faskes |
| KLB Response             | CONDITIONAL | EWS System | DEC-02 Approval |
| User Onboarding          | CONDITIONAL | SOP-12 | Distribusi Akun |
| Account Lifecycle        | CONDITIONAL | SOP-12 | Pengesahan Stakeholder |
| Support                  | CONDITIONAL | Eskalasi | PIC Support Channel |
| Deployment               | CONDITIONAL | Laporan Batch 1 & 2 | Kesiapan Infrastruktur |
| Data Governance          | CONDITIONAL | - | GOVERNANCE / LEGAL REVIEW REQUIRED |
| Audit Review             | PASS | Audit Logs (SOP-19) | - |
| Pilot Readiness          | CONDITIONALLY READY | - | GOVERNANCE DECISION REQUIRED |

---

## 10. Final Verdict

**FINAL STATUS:**
*   **TECHNICAL:** PASS (Berdasarkan Batch 1 dan Batch 2)
*   **OPERATIONAL:** CONDITIONAL
*   **GOVERNANCE:** PENDING DECISION
*   **PILOT:** CONDITIONALLY READY
*   **CITYWIDE:** NOT YET VERIFIED

**VERDICT: CONDITIONALLY OPERATIONALLY READY**
Sistem technical ready dan operational framework tersedia, namun terdapat governance/PIC/configuration dependency terbatas yang membutuhkan penyelesaian sebelum *citywide deployment*. (Kondisional karena Batch 3 governance belum disahkan, sedangkan Technical Batch 1 & 2 telah terselesaikan).

---

## 11. Batch 3 Closure

Batch 3 Operational Completeness & Governance Reconciliation telah diselesaikan berdasarkan evidence yang tersedia. Operational framework, SOP mapping, incident handling, backup/DR procedures, monitoring framework, pilot readiness, dan governance register telah terdokumentasi. Remaining dependencies berupa governance decision, organizational confirmation, dan pilot authorization tidak diambil alih oleh sistem/dokumen dan tetap menunggu keputusan pihak berwenang.
