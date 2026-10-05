# DBD Guard — Governance & Pilot Activation Package

## 1. Executive Summary
Dokumen ini merupakan kerangka persetujuan akhir (*final approval*) yang diwajibkan sebelum sistem DBD Guard dapat memasuki tahap uji coba operasional nyata (*Actual Pilot*). Saat ini, sistem telah dikonfirmasi siap secara **Teknis** (Infrastruktur & NFR) dan **Operasional** (Draft SOP & Kerangka Kerja). Namun, tahap eksekusi saat ini berstatus **PILOT BLOCKED** karena ketiadaan otorisasi administratif. 

Dokumen ini memuat seluruh keputusan krusial yang memerlukan persetujuan dan pengesahan *Stakeholder* agar status pilot dapat berubah menjadi **READY TO START**.

---

## 2. Outstanding Governance Decisions (The "DECs")
Kebijakan berikut akan berdampak langsung pada cara sistem beroperasi di lapangan. Mohon konfirmasi atas pendekatan yang akan diterapkan selama masa *Pilot*.

### DEC-01: Alur Pemindahan Stok Logistik (Stock Transfer Flow)
* **Konteks:** Sistem membutuhkan jejak audit (*audit trail*) ketat untuk setiap pemakaian logistik fogging. Jika Puskesmas A kehabisan stok, apakah dibenarkan melakukan peminjaman ke Dinkes atau Puskesmas B di luar sistem (misal via WhatsApp) yang nantinya direkonsiliasi belakangan?
* **Opsi A (Strict):** Sistem menolak alokasi fogging jika stok di aplikasi habis. *Manual transfer* dilarang keras tanpa pembaruan data sistem terlebih dahulu.
* **Opsi B (Contingency):** Pemindahan darurat diizinkan di luar sistem selama masa *pilot*, dengan kewajiban rekonsiliasi selisih stok (BAST) paling lambat H+1.
* **Keputusan Stakeholder:** [ ] Opsi A / [ ] Opsi B / [ ] Lainnya: ....................

### DEC-02: Batas Kapasitas KLB (KLB Capacity & Fallback)
* **Konteks:** Jika agregasi kasus DBD pada satu wilayah melebihi ambang batas (*threshold* EWS), apakah sistem diizinkan merekomendasikan pengerahan lebih dari 1 (satu) armada fogging secara simultan ke satu titik?
* **Opsi A:** Terbatas pada SOP normal (Maksimal 1 armada per titik, sisa antrean dijadwalkan ulang).
* **Opsi B:** Mengaktifkan mode Darurat KLB (Sistem dapat mengalokasikan armada lintas Puskesmas jika tersedia).
* **Keputusan Stakeholder:** [ ] Opsi A / [ ] Opsi B / [ ] Lainnya: ....................

### DEC-03: Masa Kedaluwarsa Banding (Appeal Expiration)
* **Konteks:** Jika laporan warga ditolak oleh Kader, warga dapat melakukan banding. Berapa lama laporan banding ini akan menggantung di sistem Puskesmas jika tidak ada yang merespons?
* **Opsi A (No Expiry):** Laporan akan terus muncul di *dashboard* Puskesmas sampai ditindaklanjuti secara manual (khusus selama *pilot*).
* **Opsi B (Auto-Expire):** Laporan banding otomatis berstatus "Selesai/Tidak Dapat Diproses" setelah batas waktu (misal: 30 hari).
* **Keputusan Stakeholder:** [ ] Opsi A / [ ] Opsi B / [ ] Lainnya: ....................

---

## 3. Pilot Scope & Constraints Authorization
Mohon otorisasi untuk parameter pembatasan uji coba lapangan:

* **Lokasi Target Pilot:** 
  * [ ] Puskesmas Beji
  * [ ] Puskesmas Pancoran Mas
  * [ ] Lainnya: ....................
* **Durasi Pilot yang Disetujui:** 
  * [ ] 14 Hari Kerja
  * [ ] 30 Hari Kerja
  * [ ] Lainnya: ....................
* **Target Partisipan Warga:** [Tertutup / Terbuka untuk publik di wilayah target]

---

## 4. PIC & Participant Assignments
Sistem membutuhkan individu nyata yang ditugaskan (*assigned*) ke dalam peran-peran berikut beserta surat perintah tugas (SPT) jika diwajibkan oleh birokrasi:

* **Koordinator Tim Pilot (PIC Utama):** ..........................................
* **Kepala Puskesmas (Approver):** ..........................................
* **Petugas Surveilans / Sanitarian Puskesmas (PE & DSS):** ..........................................
* **Koordinator Tim Fogging (Eksekutor & Inventory):** ..........................................
* **Kader Aktif Wilayah Pilot (Verifikator Teritori):** ..........................................

---

## 5. Administrative & Account Readiness Checklist
Sebelum hari pertama (Day-1) Pilot dimulai, *Stakeholder* dan Tim Teknis harus memastikan:

* [ ] **Mapping Teritori:** Data RW/RT untuk Kader wilayah target telah divalidasi.
* [ ] **Pembersihan Data Demo:** Seluruh data simulasi (laporan fiktif, transaksi fiktif) telah dibersihkan dari *database* *production*.
* [ ] **Distribusi Kredensial:** Akun akses telah dicetak/diserahkan kepada seluruh PIC secara aman.
* [ ] **Training/Onboarding:** PIC yang ditunjuk telah memahami alur *Golden Flow* dan penggunaan *support channel* (WhatsApp Group/Helpdesk).
* [ ] **SOP Pengesahan Akhir:** Draf SOP 1-19 diakui sebagai pedoman operasional valid khusus selama masa uji coba ini.

---

## 6. Go-Live Authorization Sign-off
Dengan ditandatanganinya/disetujuinya dokumen ini, maka kendala administratif (*Governance Blockers*) dinyatakan telah diselesaikan. Mode eksekusi akan diubah dari **PILOT BLOCKED** menjadi **ACTUAL PILOT / READY TO START**. 

Setiap insiden, deviasi, maupun hasil laporan di lapangan akan dievaluasi menggunakan *Pilot Execution & Evaluation Plan* untuk menentukan apakah sistem layak dilanjutkan ke tingkat distribusi tingkat kota (*Citywide Deployment*).

**Disetujui Oleh (Authorized Stakeholder):**

Nama  : ..........................................
Jabatan: ..........................................
Tanggal: ..........................................
Tanda Tangan / Persetujuan Digital: 

