# Business Decision & Governance Resolution — DBD Guard
**Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi (Pemerintah Kota Depok)**

- **Tanggal Dokumen:** 29 September 2026
- **Status Dokumen:** Formal Resolution & Decision Proposal
- **Otoritas Pemutus Kebijakan:** Kepala Dinas Kesehatan Kota Depok, Kepala Bidang Pencegahan dan Pengendalian Penyakit (P2P), dan Kepala UPTD Puskesmas se-Kota Depok
- **Status Akhir Dokumen:** **PENDING STAKEHOLDER DECISION**
- **Prinsip Rekayasa Perangkat Lunak:** **ZERO CODE MODIFICATION** (Tidak ada perubahan kode sumber, skema basis data, atau API pada fase ini)

---

## 1. Executive Summary

Setelah merampungkan seluruh tahapan audit teknis dan operasional—meliputi *Technical Testing*, *P0 Fundamental Fix & Re-Verification*, *Disaster Recovery Restoration Drill*, *Business & Operational UAT*, serta *UAT Reconciliation*—seluruh aspek stabilitas teknis arsitektur perangkat lunak **DBD Guard** telah terbukti berstatus **Technical PASS**. Sistem terbukti mampu menjaga akuntansi stok bahan kimia secara atomik tanpa kebocoran, melindungi kerahasiaan data medis pelapor, memisahkan kewenangan (*Segregation of Duties*), serta mengeksekusi mesin status pelaporan secara andal.

Kendati demikian, hasil rekonsiliasi UAT mengonfirmasi adanya wilayah abu-abu (*grey areas*) dalam alur kerja birokrasi dan kebijakan operasional lapangan yang **secara etika rekayasa tidak boleh diputuskan sepihak oleh tim developer**. Hal ini menyangkut:
1. **Tata Kelola Logistik Barang Milik Daerah (BMD):** Apakah draf permintaan insektisida hasil rekomendasi Penyelidikan Epidemiologi (PE) dapat langsung memotong meja pimpinan faskes menuju Dinas Kesehatan demi kecepatan penanganan, atau wajib melalui persetujuan internal Kepala Puskesmas?
2. **Kapasitas Penjadwalan Penanggulangan KLB:** Apakah batas mutlak 1 tugas per hari per regu fogging harus dipertahankan secara kaku demi keselamatan kerja (K3) operator, ataukah diperlukan klausul pengecualian darurat (*emergency override*) saat terjadi lonjakan wabah?
3. **Jendela Waktu Hak Banding Warga:** Apakah hak banding warga atas penolakan laporan Hoax dibatasi kuota 1 kali tanpa batas waktu, ataukah memerlukan batas kedaluwarsa waktu (48 jam / 72 jam) yang memiliki dasar epidemiologis yang jelas?

Dokumen ini disusun untuk menyajikan opsi-opsi kebijakan secara objektif, menyajikan bukti hukum dan teknis (*Source of Truth*), membedah dampak multi-dimensi secara komprehensif, serta merumuskan matriks tata kelola dan estimasi dampak implementasi. Dokumen ini juga mencatat temuan cacat keamanan kontrol wilayah **UAT-FIND-001** sebagai prioritas perbaikan pada *Implementation Plan* berikutnya.

---

## 2. Current Project Status

Berikut adalah matriks kesiapan dan status proyek DBD Guard per 29 September 2026:

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                          STATUS KESIAPAN KESELURUHAN PROYEK DBD GUARD                       │
├────────────────────────────┬─────────────────────────────┬──────────────────────────────────┤
│ Milestone / Fase Audit     │ Status Kesiapan             │ Catatan Ringkas Hasil Verifikasi │
├────────────────────────────┼─────────────────────────────┼──────────────────────────────────┤
│ 1. P0 Fundamental Fix      │ PASS (Resolved & Verified)  │ Konservasi stok atomik, proteksi │
│                            │                             │ foto medis, SoD faskes, DSS fix  │
│ 2. Technical Release Gate  │ PASS                        │ Unit, Integration, & E2E Test OK │
│ 3. Disaster Recovery Drill │ PASS (Verified)             │ Live restore drill match 100%    │
│ 4. Operational UAT         │ COMPLETED (Verified)        │ Validasi 5 aktor peran utama     │
│ 5. UAT Reconciliation      │ VERIFIED                    │ 1 Bug, 2 Tata Kelola, 1 Gap Req  │
│ 6. Governance Resolution   │ IN PROGRESS (Dokumen Ini)   │ Menunggu keputusan stakeholder   │
│ 7. Implementation Plan     │ PENDING STAKEHOLDER DECISION│ Disusun pasca penetapan opsi     │
│ 8. Production Rollout      │ GATED                       │ Menunggu rilis resolusi & P1 fix │
└────────────────────────────┴─────────────────────────────┴──────────────────────────────────┘
```

Kondisi teknis telah stabil, dan tim pengembang **tidak akan melakukan perubahan kode apa pun** sampai pemangku kepentingan (*stakeholders*) menetapkan pilihan kebijakan resmi.

---

## 3. Decision Context

Pembangunan arsitektur DBD Guard beroperasi di bawah payung hukum dan regulasi kesehatan masyarakat Pemerintah Kota Depok:

1. **Permendagri No. 19 Tahun 2016 tentang Pedoman Pengelolaan Barang Milik Daerah (BMD):**
   - Menetapkan bahwa setiap pengeluaran dan mutasi aset/logistik milik pemerintah daerah wajib melalui persetujuan Kuasa Pengguna Anggaran (KPA) / Kepala UPTD.
   - Mengatur bahwa setiap mutasi antar-unit (Dinkes ke Puskesmas) wajib disertai dokumen Berita Acara Serah Terima (BAST).
2. **Pedoman Pengendalian Vektor Kemenkes RI:**
   - Mengatur tata laksana penyemprotan insektisida termal (*thermal fogging*).
   - Menetapkan batasan paparan toksikologi insektisida bagi operator mesin fogging (standar Keselamatan dan Kesehatan Kerja / K3).
   - Menentukan radius blok penanganan (radius 100-200 meter dari titik kasus konfirmasi).
3. **Standar Pelayanan Minimal (SPM) Kesehatan Kota Depok:**
   - Menetapkan batas toleransi waktu respon laporan dugaan penyakit menular (SLA 3x24 jam atau 72 jam) sejak laporan pertama kali terverifikasi.

---

## 4. Decision 01 — Stock Request dari FOGGING_TERTUNDA

### 4.1 Kondisi & Masalah Operasional
Ketika proses verifikasi Penyelidikan Epidemiologi (PE) oleh tim faskes menyimpulkan perlunya fogging (karena jentik positif $\ge 5\%$ atau terdapat $\ge 3$ penderita demam tanpa sebab jelas), namun stok insektisida di Puskesmas setempat $< 10$ Liter, sistem pakar menghasilkan kesimpulan:
$$\text{Rekomendasi DSS: } \mathbf{FOGGING\_TERTUNDA}$$
Secara otomatis, sistem membuat draf entitas `domain.StockRequest` sebesar 10.0 Liter dengan status default `PENDING`.

### 4.2 Isu Tata Kelola (Governance Conflict)
- **Konflik Alur Otorisasi:** Pada implementasi saat ini (`pe_usecase.go:267`), pembuatan `StockRequest` otomatis langsung berstatus `PENDING` sehingga langsung masuk ke dalam antrean approval Kepala Dinas Kesehatan di Balai Kota.
- **Konflik Regulasi:** Endpoint manual pengajuan stok (`POST /api/v1/inventory/requests`) dikunci ketat menggunakan otorisasi peran `RequireRole(domain.RoleKepalaPuskesmas)`. Terjadi ketidaksinkronan tata kelola: alur otomatis memotong meja Kepala Puskesmas, sedangkan alur manual mewajibkan tandatangan Kepala Puskesmas.

### 4.3 Opsi Kebijakan

#### OPTION A: Auto-Draft Langsung Terbit ke Dinkes (Status Quo)
- **Deskripsi:** Rekomendasi `FOGGING_TERTUNDA` langsung menerbitkan `StockRequest` berstatus `PENDING` ke antrean persetujuan Dinkes tanpa melibatkan Kepala Puskesmas.
- **Tujuan Bisnis:** Memangkas birokrasi internal faskes demi mengejar kecepatan distribusi obat untuk menghentikan transmisi DBD.

#### OPTION B: Two-Step Approval Gate (Persetujuan Internal Kepala Puskesmas Dulu)
- **Deskripsi:** Rekomendasi `FOGGING_TERTUNDA` menerbitkan `StockRequest` dengan status `DRAFT_PUSKESMAS`. Kepala Puskesmas wajib mereview dan menekan tombol persetujuan internal sebelum statusnya bermutasi menjadi `PENDING` dan diteruskan ke Dinkes.
- **Tujuan Bisnis:** Menjaga akuntabilitas Barang Milik Daerah (BMD) di mana Kuasa Pengguna Anggaran Puskesmas memegang kendali penuh atas kebutuhan logistik faskesnya.

### 4.4 Evidence Review
Basis bukti hukum dan rancangan teknis untuk Decision 01:

| Evidence | Source | What It Says | Relevance | Klasifikasi Bukti |
|---|---|---|---|---|
| Auto-draft requirement | `PRD_DBD_Guard_Komprehensif.md:76` | *"Tertunda (Auto-Draft Request Stok)"* | Sistem wajib membuat draf stok saat rekomendasi tertunda. | **Explicit Requirement** |
| Otoritas Kepala Puskesmas | `PRD_DBD_Guard_Komprehensif.md:37` | *"Kepala Puskesmas: ...mengajukan permintaan stok obat ke Dinas Kesehatan."* | Wewenang pengajuan stok berada pada pimpinan faskes. | **Explicit Requirement** |
| Wewenang API Manual | `router.go:117` | `RequireRole(domain.RoleKepalaPuskesmas)` | Endpoint manual mewajibkan peran Kepala Puskesmas. | **Existing Design** |
| Status Awal Saat Ini | `pe_usecase.go:267` | `Status: domain.StockRequestStatusPending` | Kode saat ini langsung membuat status PENDING ke Dinkes. | **Existing Design** |
| Risiko Akuntabilitas BMD | `UAT_RECONCILIATION_DBD_GUARD.md:65` | *UAT-FIND-002: Bypassing Kepala Puskesmas approval* | Pejabat faskes tidak mengetahui adanya permintaan obat keluar. | **Audit Finding / Risk** |

### 4.5 Impact Analysis across 9 Dimensions

| Dimensi Analisis | OPTION A (Auto-Direct ke Dinkes) | OPTION B (Internal Gate Kepala Puskesmas) |
|---|---|---|
| **1. Business Impact** | Respon logistik sangat cepat; eliminasi vektor tidak terhambat administrasi lokal. | Menjamin kepatuhan administrasi pengelolaan aset daerah; tertib anggaran faskes. |
| **2. Operational Impact** | Petugas PE tidak perlu melapor manual ke pimpinan; Dinkes langsung memproses alokasi. | Kepala Puskesmas memiliki beban operasional harian untuk memeriksa tab approval draf stok. |
| **3. Security / Governance** | Melemahkan kontrol internal faskes (*Segregation of Duties* terabaikan). | Pengendalian internal dan tata kelola barang milik daerah (BMD) terjamin 100%. |
| **4. Data Impact** | Kolom `requester_id` diisi oleh akun Petugas PE yang memvalidasi kasus. | Kolom `requester_id` dialihkan atau diafirmasi oleh `kepala_puskesmas_id`. |
| **5. Financial / OPEX** | Berisiko terjadi over-request jika faskes sebenarnya memiliki stok penyangga tak tercatat. | Beban anggaran terkendali; faskes dapat mengecek fisik gudang sebelum meminta obat. |
| **6. User Impact** | Petugas Puskesmas merasa efisien; Dinkes mungkin menerima request tanpa verifikasi pimpinan faskes. | Kepala Puskesmas merasa dilibatkan dan memegang kedaulatan manajerial instansinya. |
| **7. Implementation Impact** | **Rendah (Nol):** Kode saat ini sudah berjalan demikian di backend (`pe_usecase.go`). | **Sedang:** Tambah status `DRAFT_PUSKESMAS`, endpoint submit, dan tab UI di Puskesmas. |
| **8. Risk Impact** | Risiko audit kepatuhan oleh Inspektorat Daerah/BPK terkait otorisasi pengadaan barang. | Risiko penundaan logistik (*bottleneck*) jika Kepala Puskesmas sedang dinas luar kota. |
| **9. Accountability Impact** | Petugas PE faskes yang menandatangani verifikasi secara *de facto* menjadi pemohon barang. | Kepala Puskesmas bertanggung jawab mutlak secara yuridis sebagai Kuasa Pengguna Anggaran. |

### 4.6 Governance & Authority Summary
- **Otoritas Pemutus:** Kepala Dinas Kesehatan Kota Depok & Kepala UPTD Puskesmas se-Kota Depok.
- **Pelaksana Sistem:** Petugas PE (pemicu draf) dan Kepala Puskesmas (penyetuju draf).
- **Status Keputusan:** **PENDING STAKEHOLDER DECISION**

---

## 5. Decision 02 — Emergency Fogging Capacity

### 5.1 Kondisi & Masalah Operasional
PRD Bab 5.2 dan Dokumen Aturan Bisnis #7 menetapkan aturan baku anti-overbooking:
> *"Satu entitas Tim Fogging hanya bisa dijadwalkan **1 tugas per hari**."*

Aturan ini bertujuan melindungi operator dari paparan uap insektisida berlebih dan kelelahan fisik akibat menggendong mesin fogging seberat $\approx 15-20\text{ kg}$. Namun, saat terjadi KLB (*outbreak*) di mana $\ge 2$ RW dalam satu kelurahan membutuhkan tindakan segera, Puskesmas yang hanya memiliki 1 regu fogging mengalami kelumpuhan kapasitas (*resource bottleneck*). Kasus di RW berikutnya terpaksa ditunda 24 hingga 48 jam, melanggar standar batas penularan epidemiologis.

### 5.2 Opsi Kebijakan

#### OPTION A: Strict Rule Without Exception (Status Quo)
- **Deskripsi:** Aturan 1 tugas per hari per regu berlaku kaku tanpa klausul pengecualian dalam kondisi apa pun.
- **Tujuan Bisnis:** Mengutamakan keselamatan dan kesehatan kerja (K3) operator teknis.

#### OPTION B: Emergency KLB Multi-Tasking Override
- **Deskripsi:** Sistem menyediakan fitur pembukaan kuota (*emergency override*) yang mengizinkan regu yang sama mengeksekusi maksimal 2 tugas per hari, khusus apabila wilayah kelurahan berstatus `KlbAlert.IsActive = true` dan disetujui pimpinan faskes/Dinkes.
- **Tujuan Bisnis:** Memaksimalkan kapasitas armada lokal untuk pemutusan rantai wabah seketika.

#### OPTION C: External Team / BKO Mobilization
- **Deskripsi:** Aturan 1 tugas/regu/hari tetap kaku, tetapi kelebihan beban tugas otomatis dialihkan ke Regu Bawah Kendali Operasi (BKO) Dinas Kesehatan Kota Depok.
- **Tujuan Bisnis:** Memperluas kapasitas faskes tanpa melanggar standar keselamatan kerja petugas faskes.

#### OPTION D: Hybrid Contingency (Kombinasi B + C)
- **Deskripsi:** Regu lokal diizinkan mengeksekusi maksimal 2 tugas/hari untuk lokasi yang berdampingan (radius $< 500\text{ m}$), dan apabila kebutuhan melebihi 2 tugas, sistem secara otomatis menerbitkan permintaan bantuan BKO Dinkes.
- **Tujuan Bisnis:** Menghadirkan resiliensi tanggap darurat komprehensif menghadapi bencana wabah.

### 5.3 Evidence Review
Basis bukti hukum dan operasional untuk Decision 02:

| Evidence | Source | What It Says | Relevance | Klasifikasi Bukti |
|---|---|---|---|---|
| Kuota Harian Fogging | `PRD_DBD_Guard_Komprehensif.md:73` | *"Satu entitas Tim Fogging hanya bisa dijadwalkan 1 tugas per hari."* | Batas kapasitas dasar sistem. | **Business Rule** |
| Implementasi Penjadwalan | `fogging_schedule_usecase.go:94` | `count >= 1 -> ErrTeamAlreadyScheduled` | Hardcode pembatasan jadwal harian di usecase. | **Existing Design** |
| Risiko Kelumpuhan Respon | `RISK_MANAGEMENT_DBD_GUARD.md:151` | *RSK-05: Kelumpuhan Respons Wabah Akibat Anti-Overbooking Kaku* | Analisis risiko keterlambatan pemutusan transmisi KLB. | **Audit Finding / Risk** |
| Standar Kemenkes K3 | Pedoman Pengendalian Vektor Kemenkes | Batas jam kerja operator fogging termal maksimal 4-5 jam kerja fisik | Perlindungan keselamatan fisik petugas lapangan. | **Explicit Requirement** |

### 5.4 Impact Analysis across 9 Dimensions

| Dimensi Analisis | OPTION A (Strict) | OPTION B (KLB Override) | OPTION C (BKO Dinkes) | OPTION D (Hybrid B+C) |
|---|---|---|---|---|
| **1. Business Impact** | K3 terjamin; mutu pengabutan terjaga. | Respon wabah cepat di tingkat lokal. | Kapasitas kota optimal tanpa overwork lokal. | Resiliensi KLB tertinggi bagi Pemkot Depok. |
| **2. Operational Impact** | Penundaan tindakan pada kluster kedua/ketiga. | Regu lokal bekerja dobel shift pada hari yang sama. | Koordinasi pergeseran armada Dinkes ke faskes. | Alur eskalasi terstruktur dari lokal ke dinas. |
| **3. Security / Governance** | Kepatuhan aturan reguler 100%. | Perlu hak akses khusus (*override role*). | Perlu nota dinas perbantuan antar-tingkat. | Tata kelola berjenjang berbasis level darurat. |
| **4. Data Impact** | Skema database tetap. | Penambahan flag `is_emergency_override` di jadwal. | Penambahan entitas / penanda `is_bko_team`. | Penambahan field override dan relasi BKO. |
| **5. Financial / OPEX** | Tidak ada biaya ekstra. | Alokasi anggaran uang lembur/insentif operator. | Biaya mobilisasi armada pickup & BBM Dinkes. | Kombinasi insentif lokal dan OPEX armada dinas. |
| **6. User Impact** | Operator merasa aman; warga RW lain kecewa. | Operator rentan lelah fisik; warga terlindungi. | Operator faskes terbantu oleh kedatangan tim dinas. | Seluruh pemangku kepentingan terlindungi optimal. |
| **7. Implementation Impact** | **Nol:** Tetap pada kode saat ini. | **Sedang:** Modifikasi guard pengecekan kuota. | **Sedang:** Logika delegasi tim lintas instansi. | **Tinggi:** Mesin orkestrasi eskalasi bertingkat. |
| **8. Risk Impact** | Risiko ledakan kasus baru DBD yang fatal. | Risiko kecelakaan kerja akibat kelelahan operator. | Risiko keterlambatan kedatangan regu BKO dinas. | Kompleksitas manajemen operasional lapangan. |
| **9. Accountability Impact** | Tanggung jawab kaku pada batasan sistem. | Kepala Faskes bertanggung jawab atas K3 regu. | Kadinkes bertanggung jawab atas kesiapan BKO. | Berbagi tanggung jawab (*shared accountability*). |

### 5.5 Governance & Authority Summary
- **Otoritas Pemutus:** Kepala Dinas Kesehatan Kota Depok & Kepala Bidang P2P.
- **Pelaksana Sistem:** Petugas Penjadwalan Fogging & Koordinator Lapangan.
- **Status Keputusan:** **PENDING STAKEHOLDER DECISION**

---

## 6. Decision 03 — Appeal Expiration

### 6.1 Kondisi & Kesenjangan Requirement
Pada alur pelaporan DBD oleh warga, apabila laporan dinyatakan ditolak sebagai *Hoax* oleh Kader setelah pengecekan fisik, warga diberikan hak hukum untuk mengajukan banding (*appeal*).
- **Ketentuan yang Sah di PRD:** Kuota pengajuan banding dibatasi **maksimal 1 kali** per laporan (`appeal_count < 1`).
- **Kesenjangan Requirement (*Requirement Gap*):** Pada dokumen audit sebelumnya diasumsikan batas waktu pengajuan banding adalah **48 jam**. Namun, setelah dilakukan penelusuran menyeluruh pada PRD dan dokumen aturan bisnis, **tidak ditemukan klausul yang menetapkan angka 48 jam**. Ini merupakan asumsi penguji terdahulu dan belum disahkan secara hukum oleh pemangku kepentingan.

### 6.2 Opsi Kebijakan

#### OPTION A: Quota-Only Limitation (Maksimal 1 Kali Tanpa Batas Waktu)
- **Deskripsi:** Hak banding berlaku maksimal 1 kali seumur hidup laporan tanpa ada pembatasan batas jam/hari (*no expiration*).
- **Tujuan Bisnis:** Ramah warga (*citizen-centric*), memberi waktu warga melampirkan bukti surat rawat inap rumah sakit yang seringkali baru keluar beberapa hari pasca penolakan kader.

#### OPTION B: Strict 48-Hour Expiration
- **Deskripsi:** Tombol banding kedaluwarsa tepat 48 jam setelah laporan berstatus `REJECTED_HOAX`.
- **Tujuan Bisnis:** Menutup siklus pelaporan secara cepat dan menjaga kebersihan antrean kerja penyelidikan epidemiologi faskes.

#### OPTION C: 72-Hour Epidemiological Window Expiration
- **Deskripsi:** Hak banding berlaku selama 72 jam (3 hari penuh), diselaraskan dengan batas toleransi waktu respon Standar Pelayanan Minimal (SLA 3x24 jam) penanganan kasus demam berdarah di Kota Depok.
- **Tujuan Bisnis:** Memberikan kompromi ideal antara waktu pengurusan surat diagnosis medis rumah sakit dan jendela masa inkubasi/penularan nyamuk di lapangan.

### 6.3 Evidence Review
Basis bukti hukum dan teknis untuk Decision 03:

| Evidence | Source | What It Says | Relevance | Klasifikasi Bukti |
|---|---|---|---|---|
| Hak Banding Warga | `PRD_DBD_Guard_Komprehensif.md:64` | *"Jika ditolak, pelapor mendapat opsi Banding."* | Menjamin adanya hak banding bagi warga. | **Explicit Requirement** |
| Alur Eskalasi Banding | `business_rules.md:6` | *"Warga kemudian diberikan opsi banding yang akan meneruskan laporan ke Puskesmas..."* | Menetapkan laporan banding langsung naik ke faskes. | **Business Rule** |
| Batasan Kuota 1 Kali | `report_repository.go:97` | `AND appeal_count < 1` | Pembatasan kuota banding maksimal 1 kali. | **Existing Design** |
| Angka 48 Jam | `P0_FIX_AND_DEEP_BUG_AUDIT_REPORT.md:262` | *"...dapat diajukan banding maksimal 1 kali dalam 48 jam"* | Batas 48 jam adalah asumsi audit, bukan requirement. | **Audit Assumption** |
| Siklus Toleransi Depok | SPM Kesehatan Depok | SLA penanganan tersangka demam berdarah adalah 3x24 jam (72 jam) | Basis rasionalitas waktu epidemiologis. | **Explicit Requirement** |

### 6.4 Impact Analysis across 9 Dimensions

| Dimensi Analisis | OPTION A (No Expiry) | OPTION B (48 Jam) | OPTION C (72 Jam) |
|---|---|---|---|
| **1. Business Impact** | Memberikan keadilan maksimal bagi warga pelapor. | Menegakkan kedisiplinan pelaporan cepat. | Selaras dengan ritme penanganan medis dan SPM 72 jam. |
| **2. Operational Impact** | Kasus lama dapat aktif kembali ke antrean faskes. | Antrean kerja Puskesmas cepat bersih dan final. | Beban faskes terukur dalam rentang 3 hari kerja. |
| **3. Security / Governance** | Risiko manipulasi laporan lama yang sudah tidak relevan. | Menutup celah reaktivasi kasus lama secara ketat. | Kontrol integritas waktu terjaga secara rasional. |
| **4. Data Impact** | Kueri database hanya mengecek `appeal_count < 1`. | Kueri mengecek `appeal_count < 1 AND updated_at >= NOW() - 48h`. | Kueri mengecek `appeal_count < 1 AND updated_at >= NOW() - 72h`. |
| **5. Financial / OPEX** | Tidak ada dampak finansial langsung. | Tidak ada dampak finansial langsung. | Tidak ada dampak finansial langsung. |
| **6. User Impact** | Sangat nyaman bagi warga yang sibuk mengurus pasien di RS. | Berisiko memicu komplain warga jika hasil lab RS terbit lambat. | Warga memiliki waktu yang memadai untuk konfirmasi RS. |
| **7. Implementation Impact** | **Nol:** Sistem saat ini sudah berjalan demikian. | **Rendah:** Penambahan filter waktu pada kueri repositori. | **Rendah:** Penambahan filter waktu pada kueri repositori. |
| **8. Risk Impact** | Risiko penyelidikan lapangan sia-sia untuk kasus 1 bulan lalu. | Risiko hilangnya kepercayaan warga terhadap layanan faskes. | Risiko terkendali dalam ambang batas masa inkubasi nyamuk. |
| **9. Accountability Impact** | Tanggung jawab verifikasi ulang tidak memiliki batas kadaluwarsa. | Batas tanggung jawab faskes berakhir pasti pada jam ke-48. | Batas tanggung jawab faskes berakhir pasti pada jam ke-72. |

### 6.5 Governance & Authority Summary
- **Otoritas Pemutus:** Kepala Dinas Kesehatan Kota Depok & Tim Bidang P2P.
- **Pelaksana Sistem:** Warga Masyarakat Pelapor & Petugas Verifikasi Banding di Faskes.
- **Status Keputusan:** **PENDING STAKEHOLDER DECISION**

---

## 7. Decision Matrix

Berikut adalah matriks komparatif komprehensif bagi para pengambil keputusan:

| Decision | Option | Business Benefit | Risk | Operational Impact | Implementation Impact | Evidence | Stakeholder Decision |
|---|---|---|---|---|---|---|:---:|
| **DECISION-01** (Stock Request Alur) | **Option A (Auto-Direct)** | Distribusi insektisida instan ke faskes | Pelanggaran tata kelola KPA faskes | Permintaan obat langsung ke Dinkes | Rendah (Status Quo) | PRD Bab 5.2 | **TBD — Stakeholder Decision** |
| | **Option B (Two-Step Approval)** | Kepatuhan hukum BMD & akuntabilitas KPA | Penundaan alokasi jika pimpinan faskes lambat | Kepala Puskesmas wajib menyetujui draf | Sedang (Status baru & tab UI) | PRD Bab 3 & Permendagri 19/2016 | **TBD — Stakeholder Decision** |
| **DECISION-02** (Kapasitas Fogging KLB) | **Option A (Strict 1 Task)** | Menjamin perlindungan K3 operator teknis | Penumpukan kasus & penularan meluas saat KLB | Jadwal kaku maksimal 1 tugas per hari | Nol (Status Quo) | PRD Bab 5.2 & Aturan Bisnis #7 | **TBD — Stakeholder Decision** |
| | **Option B (KLB Override)** | Mengatasi kluster darurat di wilayah faskes | Risiko kelelahan kerja operator lapangan | Regu lokal bertugas hingga 2 sesi/hari | Sedang (Guard validasi jadwal) | Standar Tanggap Darurat KLB | **TBD — Stakeholder Decision** |
| | **Option C (Dinkes BKO)** | Mutu dan stamina operator faskes terjaga | Ketergantungan pada kesiapan armada Dinkes | Penugasan perbantuan regu dinas | Sedang (Fitur regu eksternal) | Struktur Organisasi P2P Dinkes | **TBD — Stakeholder Decision** |
| | **Option D (Hybrid B+C)** | Resiliensi penanggulangan wabah tertinggi | Beban koordinasi multi-instansi | Eskalasi lokal 2 tugas lalu bantuan dinas | Tinggi (Mesin eskalasi bertingkat) | Best Practice Manajemen Bencana | **TBD — Stakeholder Decision** |
| **DECISION-03** (Jendela Banding) | **Option A (No Expiry)** | Memudahkan warga melengkapi bukti medis | Laporan lama dapat diaktifkan kembali | Kasus lampau berpotensi masuk antrean | Nol (Status Quo) | PRD Bab 5.1 & Aturan Bisnis #6 | **TBD — Stakeholder Decision** |
| | **Option B (48 Jam)** | Menutup siklus kasus secara tegas | Warga kehilangan hak jika surat RS terlambat | Antrean faskes bersih setelah 48 jam | Rendah (Filter interval kueri) | P0 Audit Report (Asumsi) | **TBD — Stakeholder Decision** |
| | **Option C (72 Jam)** | Harmonisasi dengan SPM 72 jam & ritme RS | Laporan aktif bertahan selama 3 hari kerja | Jendela klarifikasi teratur dan wajar | Rendah (Filter interval kueri) | SPM Dinkes Kota Depok | **TBD — Stakeholder Decision** |

---

## 8. Governance & Accountability

Berdasarkan prinsip tata kelola sistem informasi publik, berikut adalah jawaban formal atas sembilan dimensi pertanyaan tata kelola untuk masing-masing domain keputusan:

```text
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                             MATRIKS TATA KELOLA & AKUNTABILITAS SISTEM                           │
├────────────────────────────────┬──────────────────────────┬───────────────────┬──────────────────┤
│ Dimensi Pertanyaan             │ DECISION-01 (Stok Obat)  │ DECISION-02 (KLB) │ DECISION-03 (Ban)│
├────────────────────────────────┼──────────────────────────┼───────────────────┼──────────────────┤
│ 1. Siapa yang berhak memutus?  │ Kepala Dinas Kesehatan   │ Kepala Dinas & P2P│ Kepala Dinas & P2P│
│ 2. Siapa yang menjalankan?     │ Kepala Faskes & Petugas  │ Koordinator Fogging Pelapor Warga   │
│ 3. Siapa yang bertanggung jawab│ Kuasa Pengguna Anggaran  │ Penanggung Jawab K3 Pelapor Kasus   │
│ 4. Siapa yang harus diaudit?   │ Pengelola Gudang Farmasi │ Tim Pelaksana Fog  Inspektorat Faskes│
│ 5. Siapa pemegang hak override?│ Kepala Dinas Kesehatan   │ Kadinkes Depok    │ Kepala Puskesmas │
│ 6. Apakah wajib dicatat audit? │ WAJIB (Event Log Mutasi) │ WAJIB (Event KLB) │ WAJIB (Event Ban)│
│ 7. Apakah memerlukan SOP fisik?│ WAJIB (Format BAST BMD)  │ WAJIB (SOP K3 KLB)│ WAJIB (SOP Band) │
│ 8. Membutuhkan perubahan PRD?  │ Ya (Jika Opsi B dipilih) │ Ya (Jika B/C/D)   │ Ya (Jika B atau C│
│ 9. Membutuhkan perubahan Rule? │ Ya (Aturan Bisnis #7)    │ Ya (Aturan #7)    │ Ya (Aturan #2)   │
└────────────────────────────────┴──────────────────────────┴───────────────────┴──────────────────┘
```

---

## 9. Required Business Rule Changes

Setelah pemangku kepentingan menetapkan pilihan opsi, spesifikasi perubahan aturan bisnis akan ditransformasikan secara formal sebagai berikut:

### 9.1 Transformasi DECISION-01 (Jika Opsi B Ditetapkan)
- **Business Rule ID:** `BR-INV-04`
- **Rumusan Aturan:** *"Setiap rekomendasi FOGGING_TERTUNDA otomatis menerbitkan draf permintaan stok dengan status `DRAFT_PUSKESMAS`. Draf ini WAJIB diverifikasi dan disetujui oleh Kepala Puskesmas setempat sebelum diteruskan ke Dinas Kesehatan sebagai status `PENDING`."*
- **Role / Authority:** `ROLE_KEPALA_PUSKESMAS` (hanya pimpinan faskes yang dapat memutasikan status ke `PENDING`).
- **Acceptance Criteria:**
  1. Status default entitas `StockRequest` dari PE usecase adalah `DRAFT_PUSKESMAS`.
  2. Endpoint `POST /api/v1/puskesmas/stock-requests/:id/submit` hanya dapat diakses oleh peran Kepala Puskesmas.
  3. Dashboard Dinas Kesehatan hanya menampilkan permintaan stok yang telah berstatus `PENDING`.
- **Backend Validation:** Validasi state transition `DRAFT_PUSKESMAS -> PENDING`.
- **Database Requirement:** Penambahan nilai status `'DRAFT_PUSKESMAS'` pada constraint kolom status `stock_requests`.
- **Frontend Behavior:** Tampilan tab antrean *"Draf Pengajuan Faskes"* pada panel pimpinan Puskesmas.
- **Audit Trail:** Mencatat event `SUBMIT_STOCK_REQUEST_BY_KEPALA_PUSKESMAS`.
- **Test Scenario:** `TestSubmitStockRequest_Unauthorized_Forbidden` & `TestSubmitStockRequest_Success`.

### 9.2 Transformasi DECISION-02 (Jika Opsi B atau D Ditetapkan)
- **Business Rule ID:** `BR-FOG-05`
- **Rumusan Aturan:** *"Batasan kuota 1 tugas per hari per regu fogging dapat ditingkatkan menjadi maksimal 2 tugas per hari secara bersyarat apabila kelurahan target berstatus `KlbAlert.IsActive = true` dan disetujui secara tertulis oleh pimpinan operasional."*
- **Role / Authority:** `ROLE_KEPALA_PUSKESMAS` atau `ROLE_DINAS_KESEHATAN`.
- **Acceptance Criteria:**
  1. Pengecekan kuota di `fogging_schedule_usecase.go` memverifikasi status KLB aktif pada kelurahan bersangkutan.
  2. Jika KLB aktif, batas toleransi harian ditingkatkan menjadi 2 tugas per hari.
  3. Form pembuatan jadwal mewajibkan pengisian alasan darurat (*emergency notes*).
- **Backend Validation:** Validasi guard kuota kondisional: `if !isKLB && count >= 1 -> error; if isKLB && count >= 2 -> error`.
- **Database Requirement:** Kolom `is_emergency_override BOOLEAN DEFAULT FALSE` dan `override_reason TEXT` pada tabel `fogging_schedules`.
- **Frontend Behavior:** Form penjadwalan menampilkan banner siaga KLB warna kuning/merah saat regu dijadwalkan untuk tugas kedua.
- **Audit Trail:** Mencatat event `EMERGENCY_FOGGING_SCHEDULE_OVERRIDE`.
- **Test Scenario:** `TestCreateSchedule_KLBOverride_Success` & `TestCreateSchedule_NonKLB_MaxOneTask`.

### 9.3 Transformasi DECISION-03 (Jika Opsi C Ditetapkan)
- **Business Rule ID:** `BR-REP-08`
- **Rumusan Aturan:** *"Warga pelapor berhak mengajukan banding atas laporan yang dinyatakan ditolak sebagai Hoax maksimal 1 kali dalam kurun waktu 72 Jam sejak status penolakan diterbitkan oleh sistem."*
- **Role / Authority:** `ROLE_WARGA` (Pemilik laporan).
- **Acceptance Criteria:**
  1. Kueri `SubmitAppealAtomic` memvalidasi bahwa `updated_at >= NOW() - INTERVAL '72 hours'`.
  2. Apabila pengajuan dilakukan setelah melewati 72 jam, sistem mengembalikan error: `Batas waktu pengajuan banding (72 jam) telah kedaluwarsa`.
  3. Portal pelacakan laporan menonaktifkan tombol banding setelah masa 72 jam terlampaui.
- **Backend Validation:** Validasi interval waktu pada lapisan repository/usecase.
- **Database Requirement:** Tidak memerlukan perubahan kolom; memanfaatkan kolom `updated_at` pada tabel `reports`.
- **Frontend Behavior:** Komponen pelacak menampilkan countdown jam yang tersisa untuk mengajukan banding.
- **Audit Trail:** Mencatat event `REPORT_APPEAL_SUBMITTED` dengan timestamp.
- **Test Scenario:** `TestSubmitAppeal_Within72Hours_Success` & `TestSubmitAppeal_Past72Hours_ExpiredError`.

---

## 10. Implementation Impact

Tabel estimasi perubahan perangkat lunak yang disiapkan untuk *Implementation Plan* mendatang (tanpa melakukan modifikasi kode pada fase ini):

| Decision | Potential Backend Change | Potential Frontend Change | DB Change | Test Impact |
|---|---|---|---|---|
| **DECISION-01** (Opsi B) | Tambah status `DRAFT_PUSKESMAS` di `pe_usecase.go`; tambah method `SubmitStockRequest` di `inventory_usecase.go` | Tambah tab *"Draf Faskes"* di panel Puskesmas `/inventory`; tombol *"Kirim ke Dinkes"* | Tambah nilai enum/constraint `'DRAFT_PUSKESMAS'` pada tabel `stock_requests` | Tambah unit test otorisasi draf di `inventory_usecase_test.go` |
| **DECISION-02** (Opsi B/D) | Modifikasi validasi kuota di `fogging_schedule_usecase.go` dengan parameter `klb_alert` | Tambah toggle *"Pengecualian Darurat KLB"* & input alasan di form jadwal fogging | Tambah kolom `is_emergency_override` dan `override_reason` di `fogging_schedules` | Tambah test scenario multi-tugas saat KLB di `fogging_schedule_usecase_test.go` |
| **DECISION-03** (Opsi B/C) | Modifikasi kueri atomic `SubmitAppealAtomic` di `report_repository.go` dengan filter batas interval waktu | Update logika rendering tombol banding di portal `/track` (disable + label expired) | Tidak ada perubahan skema (memanfaatkan kolom `updated_at`) | Tambah test case appeal kadaluwarsa di `report_repository_test.go` |

---

## 11. UAT-FIND-001

Sebagai hasil penemuan faktual pada tahapan rekonsiliasi UAT, ditemukan satu cacat implementasi keamanan bisnis murni yang tidak memerlukan perdebatan tata kelola, melainkan **wajib diperbaiki (*Must Fix*)**:

### **UAT-FIND-001 — Kader Cross-Region Claim Bypass**
- **Klasifikasi Status:** **CONFIRMED P1 BUSINESS / SECURITY BUG**
- **Current Behavior (Perilaku Saat Ini):**
  Endpoint klaim tugas investigasi kader (`POST /api/v1/kader/reports/:id/claim`) hanya memeriksa apakah status laporan adalah `MENUNGGU_PE_KADER` dan `investigator_id IS NULL`. Endpoint **tidak memeriksa kesesuaian wilayah** antara `userLocationID` milik kader yang sedang login dengan wilayah administrasi laporan.
- **Expected Behavior (Perilaku yang Diharapkan):**
  Kader hanya diizinkan mengklaim laporan yang berada di bawah yurisdiksi faskes/kelurahan binaannya. Jika seorang kader dari Puskesmas Beji mencoba mengklaim laporan di wilayah Puskesmas Pancoran Mas via API, sistem wajib menolak dan mengembalikan respon error `HTTP 403 Forbidden` (*"Akses ditolak: Laporan berada di luar yurisdiksi wilayah kerja Anda"*).
- **Root Cause (Akar Masalah):**
  Filter pembatasan wilayah kerja kader baru diimplementasikan pada lapisan kueri penampilan antarmuka daftar laporan (`report_repository.go:GetReports`), tetapi **terlewat pada usecase eksekusi mutasi klaim** (`report_usecase.go:ClaimReport`).
- **Affected Endpoint:** `POST /api/v1/kader/reports/:id/claim`
- **Affected Layer:**
  1. Lapisan Handler: `internal/delivery/http/report_handler.go` (ekstraksi claim `user_location_id` dari JWT token).
  2. Lapisan Usecase: `internal/usecase/report_usecase.go` (validasi silang antara `userLocationID` dengan kelurahan laporan via master data kelurahan).
- **Security & Business Impact:**
  Terjadinya perebutan tugas antar-kader, ketidakefisienan verifikasi lapangan, pemborosan biaya operasional kader posyandu, dan pelanggaran prinsip *Least Privilege* serta pembatasan kewenangan teritorial faskes.
- **Proposed Fix (Rancangan Perbaikan):**
  Pada `ClaimReport(ctx, reportID, userID, userLocationID)`:
  1. Ambil relasi wilayah faskes dari laporan bersangkutan.
  2. Bandingkan `puskesmas_id` laporan dengan `userLocationID` milik kader.
  3. Tolak transaksi dengan mengembalikan `domain.ErrUnauthorizedJurisdiction` jika terjadi ketidakcocokan wilayah.
- **Verification Criteria (Kriteria Verifikasi Pengujian):**
  - Unit test `TestClaimReport_CrossJurisdiction_ShouldFail`: Kader wilayah A gagal mengklaim laporan wilayah B dengan status HTTP 403.
  - Unit test `TestClaimReport_SameJurisdiction_ShouldPass`: Kader wilayah A sukses mengklaim laporan wilayah A.

---

## 12. Decisions Required from Stakeholders

Pihak pimpinan Dinas Kesehatan Kota Depok dimohon untuk menetapkan pilihan resmi pada lembar keputusan berikut:

```text
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                   LEMBAR PENGESAHAN KEPUTUSAN BISNIS                             │
│                                           PROYEK DBD GUARD                                       │
├────┬────────────────────────────────────────────────────────┬────────────────────────────────────┤
│ No │ Agenda Keputusan Bisnis                                │ Pilihan Opsi Kebijakan             │
├────┼────────────────────────────────────────────────────────┼────────────────────────────────────┤
│ 1  │ Alur Draf Permintaan Stok saat FOGGING_TERTUNDA        │ [  ] OPTION A (Auto-Direct Dinkes) │
│    │                                                        │ [  ] OPTION B (Approval Gate Faskes│
├────┼────────────────────────────────────────────────────────┼────────────────────────────────────┤
│ 2  │ Kapasitas Penjadwalan Regu Fogging saat Kondisi KLB    │ [  ] OPTION A (Strict 1 Task/Hari) │
│    │                                                        │ [  ] OPTION B (KLB Multi-Task)     │
│    │                                                        │ [  ] OPTION C (Dinkes BKO Support) │
│    │                                                        │ [  ] OPTION D (Hybrid B + C)       │
├────┼────────────────────────────────────────────────────────┼────────────────────────────────────┤
│ 3  │ Batas Waktu Kedaluwarsa Hak Banding Laporan Warga      │ [  ] OPTION A (No Expiration / 1x) │
│    │                                                        │ [  ] OPTION B (Strict 48 Jam)      │
│    │                                                        │ [  ] OPTION C (72 Jam / SLA SPM)   │
└────┴────────────────────────────────────────────────────────┴────────────────────────────────────┘

Ditetapkan di : Kota Depok
Pada tanggal  : ___________________ 2026

Pejabat Pengesah Kebijakan:
Kepala Dinas Kesehatan Kota Depok / Kepala Bidang Pencegahan & Pengendalian Penyakit (P2P)




(_____________________________________________________)
NIP.
```

---

## 13. Final Decision Status

### **STATUS KEPUTUSAN KESELURUHAN: PENDING STAKEHOLDER DECISION**

Rangkuman status per item keputusan:
1. **DECISION-01 (Stock Request Flow):** `PENDING STAKEHOLDER DECISION`
2. **DECISION-02 (Emergency Fogging Capacity):** `PENDING STAKEHOLDER DECISION`
3. **DECISION-03 (Appeal Expiration Window):** `PENDING STAKEHOLDER DECISION`
4. **UAT-FIND-001 (Kader Task Claim Jurisdiction):** `CONFIRMED P1 BUSINESS/SECURITY BUG` (Siap masuk *Implementation Plan*)

Seluruh analisis komparatif, dasar hukum, serta dampak implementasi teknis telah disajikan secara lengkap, objektif, dan transparan. Tidak ada kode sumber yang dimodifikasi pada fase ini. Tim pengembang siap menyusun *Implementation Plan* teknis segera setelah lembar pengesahan ditandatangani oleh pemangku kepentingan.
