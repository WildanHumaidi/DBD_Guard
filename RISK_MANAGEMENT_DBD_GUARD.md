# Dokumen Manajemen Risiko (Risk Management Document) — DBD Guard

**Proyek:** DBD Guard (Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi Kota Depok)  
**Versi Dokumen:** 4.0 — Audit Menyeluruh Bisnis, Operasional, dan Tata Kelola  
**Tanggal Audit:** 28 September 2026  
**Status Audit:** Selesai (Final Comprehensive Assessment)  
**Klasifikasi:** Internal — Laporan Manajemen & Pengambilan Keputusan Pemangku Kepentingan  

---

## 1. Executive Summary

Audit Manajemen Risiko ini dilakukan secara menyeluruh terhadap seluruh ekosistem **DBD Guard Kota Depok**, mencakup regulasi/PRD, arsitektur teknis, implementasi kode backend (Golang/Fiber) dan frontend (Next.js/React), skema basis data geospasial (PostgreSQL/PostGIS), alur operasional lapangan, hingga dinamika birokrasi Pemerintah Kota Depok.

Fokus utama audit ini bukanlah semata-mata mencari *software bug* teknis, melainkan menjawab pertanyaan mendasar:
> *"Apakah sistem ini memiliki kelemahan bisnis, celah proses, atau kondisi yang dapat menyebabkan tujuan operasional kesehatan masyarakat gagal meskipun aplikasinya secara teknis terlihat berjalan normal?"*

### Temuan Kunci Audit (Key Findings)

1. **Celah Fatal Akuntansi Inventaris ("Alokasi Silang Otomatis"):**
   Sistem mengimplementasikan fitur pemotongan stok otomatis dari Puskesmas donor (`donor warehouse`) apabila Puskesmas pelaksana kekurangan insektisida. Secara digital stok terpotong, namun **tanpa ada serah terima fisik (BAST), tanpa konfirmasi Puskesmas donor, dan tanpa alur pengiriman nyata**. Hal ini menciptakan "stok hantu" (*phantom inventory*), melanggar regulasi pengelolaan Barang Milik Daerah (BMD/Permendagri No. 19/2016), dan berisiko menjadi temuan pidana maladministrasi oleh Inspektorat Kota Depok atau BPK.

2. **Moral Hazard & Absensi Verifikasi Fisik pada Penyelidikan Epidemiologi (PE):**
   Sistem sangat bergantung pada kejujuran input Kader/Petugas PE. Tidak ada kontrol validasi geofencing GPS, metadata EXIF foto, atau uji silang data demam dengan rekam medis faskes. Akibatnya, timbul dua potensi distorsi ekstrem:
   - *Under-reporting sengaja:* Petugas menginput jentik nihil dan demam 0 agar tidak repot mengurus fogging.
   - *Over-reporting pesanan:* Petugas menginput jentik positif karena desakan sosial warga/tokoh lokal ("minta difogging"), memboroskan insektisida dan mempercepat resistensi insektisida vektor nyamuk *Aedes aegypti*.

3. **Pemberangusan Wewenang Klinis Kepala Puskesmas oleh Logika DSS:**
   Jika hasil PE menghasilkan kombinasi jentik negatif dan demam < 3, sistem secara otomatis mengunci status laporan menjadi `CLOSED_NO_FOGGING` tanpa menerbitkan `FoggingRequest`. Kepala Puskesmas kehilangan diskresi medis untuk memerintahkan fogging selektif pada area yang secara epidemiologis berisiko tinggi. Keputusan medis krusial telah diambil alih secara mutlak oleh algoritma kode.

4. **Kebocoran Informasi Medis & Pelanggaran UU PDP pada Berkas Unggahan:**
   Direktori penyimpanan berkas bukti (`/storage`) dipublikasikan langsung melalui `app.Static("/storage", "./storage")` tanpa lapisan middleware autentikasi maupun otorisasi peran. Foto surat keterangan rawat inap rumah sakit, kartu identitas, dan foto hunian warga dapat diakses bebas oleh siapa pun di internet yang memiliki tautan atau melakukan enumerasi berkas.

5. **Bottleneck Kapasitas & Aturan Anti-Overbooking yang Kaku:**
   Batasan PRD bahwa 1 Tim Fogging hanya boleh mengeksekusi maksimal 1 tugas per hari memicu kelumpuhan operasional saat terjadi lonjakan kasus (outbreak). Jika suatu kecamatan hanya memiliki 1 tim operasional, laporan berstatus darurat terpaksa mengantre berhari-hari, melanggar batas waktu kritis penularan DBD (kurun 1-3 hari pertama).

6. **Blindspot Surveilans pada Early Warning System (EWS):**
   Algoritma penghitungan ambang batas Kejadian Luar Biasa (KLB) hanya menghitung laporan dengan status tindakan fogging (`VERIFIED`, `WAITING_FOGGING_APPROVAL`, `FOGGING_SCHEDULED`, `FOGGING_IN_PROGRESS`, `COMPLETED`). Laporan PE yang ditangani dengan edukasi/PSN (`CLOSED_NO_FOGGING`) diabaikan 100% dari kalkulasi KLB. Sistem buta terhadap kluster wabah nyata yang ditangani secara non-kimiawi.

7. **Risiko Keberlanjutan Organisasi & Hambatan Adopsi Lapangan:**
   Kader RT/RW bersifat relawan dengan tingkat pergantian (*turnover*) sangat tinggi, tanpa insentif anggaran operasional pulsa/kuota. Sistem belum memiliki modul manajemen akun Kader mandiri. Di sisi lain, belum ada payung hukum formal (Peraturan Walikota / SK Kepala Dinas) yang mewajibkan Puskesmas menggunakan DBD Guard, sehingga resistensi kembali ke cara manual (WhatsApp Group) sangat tinggi.

---

## 2. Ruang Lingkup (Scope)

### 2.1 Ruang Lingkup Sistem yang Diaudit
Audit mencakup seluruh domain fungsional DBD Guard sesuai dokumen PRD Komprehensif:
1. **Modul Pelaporan Warga:** Pelaporan publik (reguler & anonim), pelacakan nomor referensi, pembatalan, dan banding (*appeal*).
2. **Modul Verifikasi & Penyelidikan Epidemiologi (PE):** Verifikasi hoaks oleh Kader RW, pelimpahan/eskalasi 72 jam, pengambilalihan (*takeover*) oleh Puskesmas, dan input indikator epidemiologi.
3. **Modul Decision Support System (DSS) & Persetujuan:** Logika formula otomatisasi rekomendasi (Fogging Fokus / Tertunda / Edukasi PSN) dan persetujuan Kepala Puskesmas.
4. **Modul Penjadwalan & Eksekusi Lapangan:** Penugasan tim fogging BKO/UPT, aturan anti-overbooking, pencatatan material terpakai, dan unggah foto bukti.
5. **Modul Manajemen Logistik & Inventaris:** Stok multi-gudang (Dinkes & 38 Puskesmas), permintaan stok (*StockRequest*), transfer stok (*StockTransfer*), dan alokasi silang (*cross-allocation*).
6. **Modul Surveilans Kewilayahan, GIS, & EWS/KLB:** Heatmap geospasial berbasis PostGIS, deteksi radius kluster, dan alarm otomatis KLB.
7. **Modul Tata Kelola Pengguna & Audit Trail:** Manajemen peran RBAC, pembatasan wilayah kerja (*location scoping*), dan perekaman log jejak audit.

### 2.2 Batasan di Luar Lingkup (Non-Goals)
Sesuai batasan PRD Bagian 2, audit ini **tidak** mengevaluasi fungsionalitas berikut sebagai defisiensi:
- Modul Rekam Medis Elektronik (RME/EMR) mendalam untuk diagnosis klinis pasien rumah sakit.
- Modul kepegawaian, penggajian, dan presensi dinas/Puskesmas.
- Penanganan penyakit menular tropis di luar Demam Berdarah Dengue.
- Integrasi perangkat keras IoT/drone penyemprot otomatis.

### 2.3 Entitas Pemangku Kepentingan Terdampak
- **Dinas Kesehatan Kota Depok:** Bidang Pencegahan dan Pengendalian Penyakit (P2P), Seksi Surveilans dan Imunisasi, Pengelola Gudang Farmasi/Logistik Dinkes.
- **Puskesmas (38 Unit Se-Kota Depok):** Kepala Puskesmas, Dokter Fungsional, Petugas Surveilans/Sanitarian, Pengelola Barang Puskesmas.
- **Kader Kesehatan Wilayah:** Kader Jumantik, Kader Posyandu, dan Pengurus RW/RT di 63 Kelurahan.
- **Tim Pelaksana Fogging:** Regu BKO Dinkes dan Tim Fogging internal UPT Puskesmas.
- **Masyarakat Kota Depok:** Pelapor kasus, pasien terduga DBD, dan warga di sekitar radius paparan 100 meter.

---

## 3. Metode Penilaian Risiko (Risk Assessment Method)

Metode evaluasi risiko menggunakan matriks kualitatif dan kuantitatif standar manajemen risiko sektor publik (berkiblat pada ISO 31000 dan kerangka manajemen risiko SPBE Kementerian PAN-RB).

### 3.1 Skala Kemungkinan Terjadi (Likelihood)

| Nilai | Kategori | Frekuensi Kualitatif | Probabilitas Operasional |
|:---:|:---|:---|:---|
| **1** | Sangat Rendah | Sangat jarang terjadi; kondisi luar biasa atau anomali langka. | < 5% per tahun |
| **2** | Rendah | Kemungkinan kecil terjadi; pernah terjadi di proyek sejenis. | 5% - 20% per tahun |
| **3** | Sedang | Dapat terjadi beberapa kali dalam setahun pada kondisi normal. | 21% - 50% per tahun |
| **4** | Tinggi | Sangat mungkin terjadi secara rutin dalam siklus bulanan. | 51% - 80% per tahun |
| **5** | Sangat Tinggi | Hampir pasti terjadi secara berulang dalam operasional harian/mingguan. | > 80% per tahun |

### 3.2 Skala Dampak Keparahan (Impact)

| Nilai | Kategori | Dampak Operasional & Layanan | Dampak Finansial & Hukum | Dampak Kesehatan & Reputasi |
|:---:|:---|:---|:---|:---|
| **1** | Sangat Rendah | Gangguan operasional minor (< 1 jam); tidak ada data hilang. | Kerugian finansial tidak material (< Rp 1 juta). | Keluhan perorangan; tidak ada eksposur media. |
| **2** | Rendah | Gangguan layanan lokal terisolasi (< 6 jam); koreksi manual cepat. | Pemborosan logistik minor (< Rp 10 juta); teguran internal. | Ketidaknyamanan segelintir warga; tidak berdampak klinis. |
| **3** | Sedang | Gangguan alur kerja tingkat Puskesmas (< 24 jam); SLA terlampaui. | Kerugian logistik sedang (Rp 10 - 50 juta); ketidaksesuaian SOP. | Keterlambatan intervensi wabah lokal; eksposur media sosial terbatas. |
| **4** | Tinggi | Kelumpuhan respons wabah pada beberapa kelurahan; data rusak/hilang. | Temuan audit Inspektorat (Rp 50 - 250 juta); risiko tuntutan hukum. | Kluster DBD meluas tanpa intervensi; penurunan kepercayaan publik kota. |
| **5** | Sangat Tinggi | Kelumpuhan total sistem kota Depok; kegagalan identifikasi KLB. | Temuan BPK/APH; pelanggaran UU PDP sanksi denda pidana/perdata. | Korban jiwa akibat keterlambatan respons; krisis reputasi Pemerintah Kota. |

### 3.3 Matriks Tingkat Risiko (Risk Matrix)

Perkalian antara **Likelihood (L)** dan **Impact (I)** menghasilkan **Nilai Risiko (Score 1 - 25)**:
- **Rendah (Low / 1 - 5):** Risiko dapat diterima; dipantau melalui prosedur rutin.
- **Sedang (Medium / 6 - 11):** Memerlukan tindakan mitigasi terencana dan pengawasan berkala.
- **Tinggi (High / 12 - 19):** Memerlukan intervensi manajemen segera dan alokasi sumber daya khusus.
- **Kritis / Ekstrem (Critical / 20 - 25):** Ancaman fatal; operasional tidak boleh dilanjutkan tanpa perbaikan mendasar.

```
Dampak (Impact)
  5 |  5 (M)  | 10 (M)  | 15 (H)  | 20 (C)  | 25 (C)
  4 |  4 (L)  |  8 (M)  | 12 (H)  | 16 (H)  | 20 (C)
  3 |  3 (L)  |  6 (M)  |  9 (M)  | 12 (H)  | 15 (H)
  2 |  2 (L)  |  4 (L)  |  6 (M)  |  8 (M)  | 10 (M)
  1 |  1 (L)  |  2 (L)  |  3 (L)  |  4 (L)  |  5 (M)
    +-------------------------------------------------
        1         2         3         4         5
                     Kemungkinan (Likelihood)
```

### 3.4 Kategori Bukti (Evidence Classification)
Setiap penilaian dalam dokumen ini disandarkan pada bukti objektif dengan label:
- **Verified:** Terverifikasi langsung melalui inspeksi kode sumber dan skema basis data aktif.
- **Static Analysis:** Ditemukan melalui penelusuran alur logika statis dan relasi antarmuka/kontrak API.
- **Test Verified:** Dibuktikan melalui eksekusi unit test, skenario integrasi, atau pengujian automasi.
- **Runtime Verified:** Dibuktikan saat aplikasi berjalan di lingkungan server/lokal.
- **Needs Verification:** Berpotensi terjadi secara empiris namun memerlukan konfirmasi data riil instansi.

---

## 4. Kategori Risiko (Risk Categories)

Untuk memudahkan pemetaan dan penunjukan *Risk Owner*, risiko diklasifikasikan ke dalam 7 area strategis:

1. **BGR (Business, Governance & Strategic Risk):** Tata kelola kelembagaan, kepemilikan sistem, payung hukum, komitmen anggaran daerah, resistensi birokrasi, dan akuntabilitas keputusan.
2. **OER (Operational & Epidemiological Field Risk):** Hambatan eksekusi lapangan, koordinasi lintas wilayah, beban kerja kader, alur kerja darurat, dan keselamatan kerja (K3).
3. **FAR (Fraud, Abuse & Moral Hazard Risk):** Penyalahgunaan kewenangan, manipulasi data lapangan, rekayasa kebutuhan fogging, kebocoran/penggelapan logistik, dan laporan palsu terorganisir.
4. **DAR (Data, Decision Support & Algorithm Risk):** Bias algoritma DSS, kesalahan formula otomatisasi, hilangnya diskresi klinis dokter, kualitas data input warga, dan distorsi statistik KLB.
5. **TIR (Technology, Architecture & Infrastructure Risk):** Titik kegagalan tunggal (SPOF), skalabilitas server, keandalan basis data spatial, kuota pihak ketiga, dan performa beban puncak.
6. **SCR (Security, Privacy & Compliance Risk):** Kebocoran data pribadi (PII), akses tanpa hak pada rekam medis, pelanggaran UU No. 27/2022 (UU PDP), kepatuhan regulasi Kemenkes, dan audit log integrity.
7. **BCR (Business Continuity & Disaster Recovery Risk):** Kelangsungan layanan saat sistem lumpuh (*blackout*), prosedur fallback manual berbasis kertas, dan rekonsiliasi data pasca-bencana.

---

## 5. Register Risiko (Risk Register)

Tabel berikut menyajikan register risiko komprehensif DBD Guard yang telah dievaluasi berdasarkan kondisi kode sumber, arsitektur, dan konteks operasional riil Kota Depok:

| ID | Kategori | Risiko Spesifik | Akar Penyebab (Cause) | Dampak Bisnis (Business Impact) | Kontrol Eksisting | Celah Kontrol (Control Gap) | L | I | Inherent Risk | Residual Risk | Perlakuan Risiko (Treatment) | Pemilik Risiko (Risk Owner) | Status Bukti |
|:---|:---:|:---|:---|:---|:---|:---|:---:|:---:|:---:|:---:|:---|:---|:---:|
| **RSK-01** | FAR | **Manipulasi Stok Melalui Alokasi Silang Otomatis** | Logika `inventory_usecase.go` otomatis memotong stok Puskesmas donor tanpa approval atau transfer fisik riil saat gudang lokal minus. | Laporan mutasi barang fiktif; selisih stok fisik vs digital; temuan audit BPK/Inspektorat atas pelanggaran Permendagri 19/2016 (BMD). | Audit log mencatat mutasi stok donor. | Tidak ada verifikasi fisik, tidak ada Berita Acara Serah Terima (BAST), tidak ada konfirmasi penerimaan oleh donor. | 4 | 5 | **20 (C)** | **16 (H)** | Mitigate: Hapus pemotongan lintas gudang otomatis; wajibkan alur formal *StockTransfer* dengan otorisasi dua pihak. | Pengelola Barang Dinkes / Puskesmas | Verified |
| **RSK-02** | DAR | **Pemberangusan Diskresi Medis oleh Algoritma DSS** | Hasil PE dengan jentik=false dan demam=0 langsung mengubah status laporan menjadi `CLOSED_NO_FOGGING` tanpa membuat `FoggingRequest`. | Pasien demam berulang di kluster non-jentik tidak tertangani; Kepala Puskesmas tidak bisa menyetujui fogging meskipun ada indikasi klinis kuat. | Nilai variabel PE tersimpan di tabel `pe_results`. | Sistem tidak menyediakan opsi override klinis bagi Kepala Puskesmas untuk merekomendasikan fogging manual. | 4 | 4 | **16 (H)** | **12 (H)** | Mitigate: Posisikan DSS murni sebagai rekomendasi; seluruh hasil PE wajib masuk antrean tinjauan dokter/Kepala Puskesmas dengan alasan override terdokumentasi. | Kepala Puskesmas / Seksi Surveilans Dinkes | Verified |
| **RSK-03** | SCR | **Ekspos Terbuka Foto Rekam Medis & PII Pelapor** | Endpoint `app.Static("/storage", "./storage")` menyajikan direktori unggahan foto tanpa middleware token autentikasi/otorisasi. | Foto diagnosis rumah sakit, KTP, dan hunian pribadi bocor ke publik; sanksi administratif dan pidana denda UU PDP No. 27/2022. | Penamaan acak berbasis UUID pada berkas unggahan. | Tidak ada autentikasi token; tidak ada pre-signed URL; tidak ada validasi hak akses pemohon berkas. | 5 | 5 | **25 (C)** | **15 (H)** | Mitigate: Tutup direktori statis publik; implementasikan handler pengunduhan terproteksi dengan validasi RBAC atau MinIO Pre-signed URL berbatas waktu. | Technical Owner / Data Protection Officer | Verified |
| **RSK-04** | FAR | **Moral Hazard Input Data PE Lapangan** | Input variabel jentik dan jumlah demam bersifat murni deklaratif teks tanpa validasi geolokasi atau silang rekam medis. | Terjadinya fogging pesanan (korupsi logistik) atau penolakan fogging fiktif untuk menghindari beban kerja lapangan. | Catatan audit trail menyimpan identitas investigator. | Tidak ada validasi GPS geofencing saat submit; tidak ada verifikasi silang laboratorium darah faskes. | 4 | 4 | **16 (H)** | **12 (H)** | Mitigate: Wajibkan geotagging foto jentik; terapkan audit sampling PE oleh supervisor Puskesmas; integrasikan verifikasi lab RSUD. | Seksi P2P Dinkes / Sanitarian Puskesmas | Static Analysis |
| **RSK-05** | OER | **Kelumpuhan Respons Wabah Akibat Anti-Overbooking Kaku** | Sistem menolak pembuatan jadwal jika tim fogging telah memiliki 1 jadwal pada tanggal yang sama (`CountSchedulesByTeamAndDate >= 1`). | Penanganan kluster KLB tertunda berhari-hari; percepatan transmisi gigitan nyamuk di wilayah padat penduduk. | Pengecekan kuota harian di backend sebelum jadwal disimpan. | Tidak ada mekanisme pembagian shift darurat, perbanyakan sub-tim, atau status kedaruratan KLB untuk melipatgandakan kapasitas. | 4 | 5 | **20 (C)** | **16 (H)** | Mitigate: Sediakan toggle "Mode Kedaruratan KLB" yang mengizinkan multi-jadwal berbasis persetujuan Kepala Dinas; rekrut mitra pihak ketiga terverifikasi. | Bidang P2P Dinkes / Kepala Puskesmas | Verified |
| **RSK-06** | DAR | **Blindspot Early Warning System (EWS) KLB** | Fungsi `CountCasesByArea` hanya menghitung laporan dengan status intervensi fogging, mengecualikan status `CLOSED_NO_FOGGING`. | Kluster kasus demam tinggi yang tidak difogging (misal hanya abate/PSN) tidak memicu alarm KLB; Dinkes terlambat menetapkan status darurat. | Penghitungan agregasi otomatis di repository basis data. | Kriteria SQL mengecualikan kasus terkonfirmasi non-fogging; cakupan wilayah KLB (*RegionScope*) terkunci kaku antara Kelurahan ATAU Kecamatan. | 4 | 5 | **20 (C)** | **15 (H)** | Mitigate: Ubah query EWS agar menghitung seluruh laporan terkonfirmasi terlepas dari metode intervensi; aktifkan evaluasi paralel Kelurahan dan Kecamatan. | Epidemiolog Kesehatan Dinkes | Verified |
| **RSK-07** | BGR | **Resistensi Pengguna Lapangan & Pengabaian Sistem** | Beban kerja administratif ganda tanpa mandat hukum resmi (SK Walikota) dan minimnya literasi digital di tingkat relawan RW. | Kader dan Puskesmas kembali mengandalkan WhatsApp; sistem menjadi "kuburan data" yang tidak diperbarui secara real-time. | Antarmuka web PWA responsif. | Tidak adanya sanksi administratif bagi unit kerja yang lambat; tidak ada anggaran kompensasi kuota data/pulsa pelaporan relawan. | 4 | 4 | **16 (H)** | **12 (H)** | Mitigate: Terbitkan Perwal/SK Kadinkes penggunaan wajib sistem; sediakan insentif operasional pulsa pelaporan; buat dashboard pemeringkatan kepatuhan Puskesmas. | Kepala Dinas Kesehatan Kota Depok | Needs Verification |
| **RSK-08** | BGR | **Pelanggaran Segregation of Duties pada Alur Takeover** | Handler `ValidatePEResult` tidak memvalidasi apakah validator berbeda dengan pemeriksa lapangan (`peResult.OfficerID != validatorID`). | Satu orang petugas Puskesmas dapat mengambil alih laporan, menginput PE rekayasa, lalu memvalidasi hasilnya sendiri tanpa supervisi. | Pengecekan peran pengguna melalui middleware JWT. | Tidak ada aturan pemisahan tugas (Four-Eyes Principle) di tingkat usecase bisnis. | 3 | 4 | **12 (H)** | **8 (M)** | Mitigate: Tambahkan validasi eksplisit bahwa ID pembuat PE dilarang keras menjadi validator laporan yang sama. | Koordinator Tata Kelola IT Dinkes | Verified |
| **RSK-09** | TIR | **Titik Kegagalan Tunggal Basis Data & Ketiadaan DR** | Seluruh operasional bergantung pada satu kontainer PostgreSQL tanpa replikasi streaming dan tanpa backup otomatis terjadwal. | Kehilangan data seluruh riwayat kasus wabah dan mutasi logistik jika media penyimpanan peladen rusak (*disk crash*). | PostGIS extension aktif; query spatial berjalan normal. | Belum ada cron script `pg_dump` ke penyimpanan sekunder; tidak ada skenario Failover / Disaster Recovery plan tertulis. | 3 | 5 | **15 (H)** | **10 (M)** | Mitigate: Pasang cronjob pencadangan otomatis harian terenkripsi ke off-site storage; siapkan prosedur pemulihan basis data teruji berkala. | Administrator Infrastruktur IT | Verified |
| **RSK-10** | SCR | **Pelanggaran Privasi melalui Endpoint Pelacakan Kontak** | Endpoint publik `GET /api/v1/reports/phone/:phone` menyajikan daftar seluruh laporan hanya berdasarkan parameter nomor telepon tanpa verifikasi OTP. | Pihak ketiga dapat mengintai riwayat kesehatan dan alamat tetangga hanya dengan menebak/memasukkan nomor telepon warga. | Rate limiter publik membatasi volume permintaan IP. | Tidak ada verifikasi kepemilikan nomor (OTP/autentikasi); data lokasi terekspos tanpa persetujuan subjek data. | 4 | 4 | **16 (H)** | **8 (M)** | Mitigate: Hapus endpoint pencarian nomor bebas; wajibkan login warga atau verifikasi kode OTP SMS/WhatsApp sebelum menampilkan daftar laporan. | Technical Owner / Tim Pengembang | Verified |
| **RSK-11** | FAR | **Penyelesaian Fogging Fiktif Tanpa Verifikasi Foto** | Method `UpdateExecutionStatus` mengubah status jadwal menjadi `COMPLETED` hanya dengan memvalidasi keberadaan `usedItems`, tanpa mengecek unggahan bukti foto. | Tim Fogging dapat melaporkan tugas telah selesai dan menghabiskan stok obat tanpa benar-benar melakukan penyemprotan di lapangan. | Tabel `fogging_evidences` tersedia untuk menyimpan foto bukti. | Pengecekan keberadaan rekaman foto bukti di basis data tidak dijadikan prasyarat penutupan status jadwal. | 4 | 4 | **16 (H)** | **8 (M)** | Mitigate: Wajibkan query verifikasi minimal 1 foto bukti terunggah sebelum status jadwal fogging dapat diubah menjadi `COMPLETED`. | Seksi Pengendalian Vektor Dinkes | Verified |
| **RSK-12** | OER | **Monopoli Penugasan oleh Kader (Task Hoarding)** | Endpoint klaim laporan (`ClaimReport`) tidak membatasi kuota maksimum laporan aktif yang dapat dipegang oleh seorang Kader. | Satu relawan mengklaim puluhan laporan sekaligus tetapi tidak sempat memeriksa, menyebabkan penanganan terlantar hingga batas SLA 72 jam. | Auto-release cron melepaskan tugas setelah 72 jam; Force Unassign tersedia bagi admin. | Tidak ada pembatasan kuota konkurensi klaim (misal maksimal 3 laporan aktif per kader). | 4 | 3 | **12 (H)** | **6 (M)** | Mitigate: Batasi klaim aktif maksimal 3 tugas per akun Kader; tolak klaim baru jika tugas sebelumnya belum diselesaikan atau dikembalikan. | Operational Owner / Puskesmas | Verified |
| **RSK-13** | TIR | **Kegagalan Notifikasi Kritis Akibat Limitasi Provider Gratis** | Pengiriman notifikasi bergantung pada SMTP server publik/Gmail gratis dan Firebase Cloud Messaging tier gratis tanpa sistem antrean failover. | Petugas Puskesmas dan warga tidak menerima peringatan dini KLB atau pemberitahuan jadwal fogging saat kuota harian provider terlampaui. | Integrasi kode FCM dan Gomail tersedia di backend. | Kredensial produksi belum aktif (masih mode mock); tidak ada antrean pesan (Message Queue) bertoleransi kegagalan (*retry mechanism*). | 4 | 4 | **16 (H)** | **8 (M)** | Mitigate: Gunakan SMTP resmi Pemkot Depok (`@depok.go.id`); pasang message broker (Redis/RabbitMQ) dengan strategi exponential backoff retry. | Technical Owner / Diskominfo Depok | Verified |
| **RSK-14** | FAR | **Penolakan Palsu pada Kluster Ledakan Kasus Nyata** | Logika anti-spam otomatis menolak (`REJECTED`) setiap laporan baru jika terdapat >= 3 laporan dalam radius 50m dan rentang 1 jam. | Laporan keluarga yang serumah atau tetangga yang bersamaan terjangkit DBD langsung ditolak sistem tanpa verifikasi manusia, memicu kematian preventable. | Pencatatan audit trail alasan `SPAM_CLUSTERING_DETECTED`. | Sistem berasumsi multi-laporan berdekatan pasti spam, mengabaikan karakteristik penularan DBD kluster keluarga yang sangat cepat. | 3 | 5 | **15 (H)** | **10 (M)** | Mitigate: Ubah status laporan berulang menjadi `POTENTIAL_CLUSTER_VERIFY` (antrean prioritas investigasi), bukan otomatis ditolak mentah-mentah. | Epidemiolog Dinkes / Tim Teknis | Verified |
| **RSK-15** | BGR | **Ketiadaan Modul Siklus Hidup Akun Relawan Lapangan** | Panel Developer hanya menyediakan antarmuka pembuatan akun Dinkes dan Puskesmas, tanpa ada modul pembuatan dan rotasi akun Kader RW. | Petugas terpaksa berbagi akun massal atau bergantung pada seeder manual tim developer; tidak ada mekanisme pencabutan hak akses kader yang berhenti. | Skema otentikasi JWT dan tabel `users` memiliki atribut `active_status`. | Tidak ada endpoint dan UI manajemen pengguna untuk peran Kader dan Tim Fogging di dashboard Puskesmas. | 4 | 4 | **16 (H)** | **8 (M)** | Mitigate: Bangun modul administrasi relawan wilayah di dashboard Petugas Puskesmas untuk mendaftarkan dan menonaktifkan Kader di wilayah kerjanya. | Tim Pengembang Frontend/Backend | Verified |
| **RSK-16** | BCR | **Ketiadaan Prosedur Transisi Operasional Manual** | Tidak ada SOP atau instrumen formal berbasis kertas jika peladen mati total lebih dari 2x24 jam saat musim puncak wabah (Desember - Maret). | Penyelidikan epidemiologi dan fogging mandek total; pencatatan stok gudang berantakan; data historis manual tidak pernah diinput ulang. | Layanan darurat kota Call Center 112 tersedia di Pemkot Depok. | Belum ada template formulir PE fisik darurat, buku log pemakaian obat fisik, dan protokol rekonsiliasi data pasca-pemulihan sistem. | 3 | 4 | **12 (H)** | **8 (M)** | Mitigate: Terbitkan Standard Operating Procedure (SOP) Kedaruratan Manual; siapkan formulir cetak standar di setiap Puskesmas; tetapkan protokol input susulan 24 jam pasca-pulih. | Manajemen Dinkes / Bagian Organisasi | Needs Verification |
| **RSK-17** | SCR | **Ketiadaan Mekanisme Penghapusan Data (Right to Erasure)** | Basis data tidak mendukung alur anonimisasi permanen atau penghapusan data identitas warga yang menarik kembali laporannya. | Pelanggaran hak subjek data Pasal 43 UU No. 27/2022 tentang Pelindungan Data Pribadi (UU PDP). | Tersedia fitur pelaporan anonim sejak awal pengisian. | Warga pelapor reguler tidak memiliki opsi mandiri maupun administratif untuk menghapus data kontak dan alamatnya pasca-penanganan selesai. | 3 | 4 | **12 (H)** | **6 (M)** | Mitigate: Sediakan fitur anonimisasi otomatis untuk data pribadi pelapor (nama & kontak) setelah masa retensi surveilans (misal 30 hari pasca-laporan ditutup). | Legal Officer / Tim Teknis | Static Analysis |
| **RSK-18** | TIR | **Pembatasan Lisensi Layanan Peta OpenStreetMap Publik** | Frontend mengandalkan tile OpenStreetMap gratis (`tile.openstreetmap.org`) tanpa peladen tile mandiri atau caching layer. | Tile peta gagal dimuat (layar peta abu-abu/blank) karena diblokir oleh OSM Policy akibat lonjakan permintaan warga saat puncak wabah DBD. | Integrasi library Leaflet.js open-source. | Tidak ada reverse-proxy caching, tidak ada fallback server tile, tidak ada peladen tile lokal mandiri. | 3 | 4 | **12 (H)** | **6 (M)** | Mitigate: Bangun reverse-proxy tile cache mandiri di lingkungan Pemkot Depok atau gunakan penyedia tile open-source dengan perjanjian SLA yang jelas. | Technical Owner / Diskominfo | Static Analysis |

---

## 6. Business Weakness & Business Loophole

Bagian ini membedah secara khusus kelemahan rancangan proses bisnis, insentif yang salah arah (*perverse incentives*), serta celah operasional yang memungkinkan manipulasi atau kegagalan sistem.

| ID | Area Bisnis | Kelemahan / Celah Proses (Loophole) | Mengapa Ini Berbahaya (Why It Matters) | Potensi Penyalahgunaan / Skenario Kegagalan | Dampak Bisnis Nyata | Rekomendasi Solusi Bisnis |
|:---|:---|:---|:---|:---|:---|:---|
| **BW-01** | Pengelolaan Inventaris | **Pengurangan Stok Otomatis Antar-Puskesmas ("Alokasi Silang Hantu")** | Mengizinkan satu entitas memotong aset fisik entitas lain tanpa dokumen hukum serah terima adalah pelanggaran tata kelola keuangan negara. | Tim Fogging di Puskesmas Beji kehabisan obat, sistem diam-diam memotong stok milik Puskesmas Cimanggis. Petugas Cimanggis heran stok digitalnya berkurang 10 liter padahal barangnya masih di gudang atau sebaliknya. | Korupsi logistik tidak terlacak; selisih pembukuan neraca aset daerah; temuan pidana maladministrasi oleh Badan Pemeriksa Keuangan (BPK). | Hapus fungsi auto-deduct donor warehouse. Ganti dengan alur *Emergency Stock Transfer* yang mewajibkan konfirmasi kirim dari donor dan verifikasi terima dari pemohon. |
| **BW-02** | Verifikasi & Penyelidikan (PE) | **Ketiadaan Verifikasi Kehadiran Fisik Lapangan pada PE** | Kader/Petugas dapat mengisi form PE dari mana saja tanpa pernah menginjakkan kaki di lokasi warga yang demam. | Kader yang enggan repot segera mencentang "Jentik: Tidak" dan "Demam: 0" dari rumah. Laporan ditutup otomatis sebagai edukasi, padahal di lokasi terdapat sarang nyamuk masif. | Kematian warga akibat keterlambatan fogging; kluster wabah meledak di RW bersangkutan; hilangnya legitimasi program kesehatan kota. | Wajibkan validasi koordinat GPS saat tombol "Submit PE" ditekan (toleransi radius 150m dari lokasi laporan); sertakan foto jentik ber-watermark waktu dan koordinat. |
| **BW-03** | Algoritma Rekomendasi (DSS) | **Sistem Menutup Kasus Secara Sepihak Tanpa Tinjauan Klinis Dokter** | Kode program (`pe_usecase.go:217`) langsung memberi vonis `CLOSED_NO_FOGGING` jika jentik nihil, memotong wewenang dokter fungsional. | Pasien balita dengan trombosit anjlok dan syok DBD di pemukiman padat tidak bisa difogging karena kader tidak menemukan jentik saat inspeksi singkat 10 menit. | Malpraktik institusional; Kepala Puskesmas tidak berdaya mengubah rekomendasi sistem; tanggung jawab medis menjadi kabur antara dokter dan pembuat software. | Seluruh hasil PE wajib berstatus rekomendasi awal; status akhir dan disposisi fogging mutlak berada di tangan Kepala Puskesmas dengan alasan medis tercatat. |
| **BW-04** | Validasi dan Pengawasan PE | **Bypass Segregation of Duties pada Pengambilalihan Tugas** | Petugas Puskesmas yang melakukan *Takeover* dari Kader yang lalai dapat merangkap sebagai pihak yang memvalidasi hasil investigasinya sendiri. | Petugas Puskesmas A mengambil alih laporan RW 02, mengisi hasil PE fiktif, lalu meng-approve PE tersebut di dashboard tanpa ada orang kedua yang memeriksa (*no four-eyes check*). | Ketiadaan akuntabilitas internal; celah kolusi antara oknum Puskesmas dan vendor jasa chemical fogging liar; manipulasi indikator kinerja surveilans. | Terapkan aturan ketat: Akun yang bertindak sebagai investigator lapangan dilarang secara sistemik menjadi validator pada laporan yang sama (`OfficerID != ValidatorID`). |
| **BW-05** | Pelaporan Warga Publik | **Pembatalan Laporan Tanpa Bukti Otorisasi Pemilik Laporan** | Siapa pun yang mengetahui nomor referensi (`REP-YYYYMMDD-XXXX`) dapat membatalkan laporan yang belum diklaim melalui API publik. | Tetangga atau pihak yang tidak ingin lingkungannya disorot kasus DBD dapat menebak nomor referensi atau mencegat URL tracking untuk membatalkan laporan warga lain. | Laporan kasus terhapus/batal tanpa sepengetahuan pelapor asli; penanganan wabah terbengkalai. | Wajibkan kode verifikasi (PIN rahasia 4-digit atau OTP) yang diberikan kepada warga saat pertama kali melapor untuk aksi pembatalan laporan. |
| **BW-06** | Eksekusi Jadwal Fogging | **Penutupan Jadwal Fogging Tanpa Kewajiban Unggah Bukti Lapangan** | Tim Fogging dapat menekan status `COMPLETED` dan memotong stok gudang hanya dengan mengisi angka pemakaian bahan tanpa melampirkan foto. | Petugas nakal melaporkan telah menyemprot 10 liter malathion padahal obat dijual ke pihak swasta atau tidak disemprotkan; laporan ditutup selesai di sistem. | Penggelapan insektisida daerah; pencemaran lingkungan; warga tidak mendapatkan proteksi fogging yang dilaporkan selesai. | Sistem wajib memblokir tombol status `COMPLETED` jika jumlah rekaman di tabel `fogging_evidences` untuk jadwal tersebut masih nol. |
| **BW-07** | Deteksi Spam Otomatis | **Pemberangusan Kasus Ledakan Wabah Keluarga oleh Filter Anti-Spam** | Aturan menolak otomatis laporan jika ada >= 3 laporan dalam radius 50m dalam 1 jam secara keliru menganggap wabah serentak sebagai serangan spam. | Dalam satu kluster keluarga besar di perkampungan padat, 4 anggota keluarga demam mendadak dan melapor secara berurutan; laporan ke-3 dan ke-4 ditolak otomatis sebagai hoaks. | Kegagalan mendeteksi ledakan kasus fokal; warga kehilangan kepercayaan karena laporannya divonis spam oleh mesin. | Ubah perlakuan spam otomatis: jangan ditolak, melainkan ditandai sebagai `FLAGGED_CLUSTER_SUSPECT` dan diprioritaskan untuk verifikasi telepon segera oleh kader/Puskesmas. |
| **BW-08** | Pelacakan Nomor Telepon | **Pencarian Riwayat Laporan Tanpa Otentikasi Kepemilikan** | Endpoint publik `/reports/phone/:phone` mengembalikan seluruh riwayat laporan medis dan alamat presisi hanya berbekal input nomor ponsel. | Oknum pihak ketiga memasukkan nomor ponsel target warga untuk mengetahui apakah rumahnya sedang terjangkit DBD, alamat lengkap, dan nomor kontak keluarga. | Stigmatisasi sosial terhadap keluarga pasien DBD; pelanggaran kerahasiaan data pribadi; potensi pemerasan atau penipuan berkedok petugas kesehatan. | Nonaktifkan endpoint publik berbasis nomor telepon terbuka; integrasikan fitur pelacakan multi-laporan hanya melalui tautan verifikasi OTP SMS/WhatsApp. |
| **BW-09** | Manajemen Logistik Lapangan | **Stok Terpotong Permanen Meskipun Penggunaan Logistik Ditolak** | Pemotongan stok gudang riil terjadi saat eksekusi selesai, namun modul persetujuan logistik sekunder yang menolak pemakaian bahan tidak mengembalikan (*refund*) stok. | Puskesmas menolak laporan penggunaan bahan tim fogging karena dianggap tidak wajar/boros. Status penggunaan menjadi `REJECTED`, tetapi saldo fisik gudang di database tetap berkurang. | Selisih permanen pada pembukuan inventaris gudang farmasi Puskesmas; ketidaksesuaian laporan pertanggungjawaban barang tahunan. | Sinkronkan alur logistik: Jika penolakan pemakaian obat diverifikasi benar adanya penyelewengan, terapkan entri jurnal koreksi (*StockAdjustment*) yang memulihkan saldo buku. |
| **BW-10** | Administrasi Relawan | **Absensi Manajemen Siklus Akun Relawan Lapangan (Kader)** | Tidak ada antarmuka bagi Puskesmas untuk menambah, mengubah nomor HP, atau menonaktifkan Kader Jumantik di wilayah kerjanya. | Kader yang pindah domisili atau mengundurkan diri tetap memegang akun aktif; kader baru tidak bisa mendaftar mandiri dan terpaksa menggunakan kredensial lama. | Kebocoran data warga ke mantan relawan; kebingungan identitas petugas pemeriksa lapangan; lumpuhnya pelaporan saat terjadi suksesi kepengurusan RT/RW. | Bangun antarmuka terdesentralisasi di portal Puskesmas untuk mengelola direktori relawan per RW, lengkap dengan mekanisme penonaktifan akun instan. |

---

## 7. Risiko Operasional Lapangan (Operational Risk)

### 7.1 Keterbatasan Literasi Digital & Perangkat Relawan RW
Sebagian besar Kader Jumantik di Kota Depok adalah ibu-ibu PKK atau tokoh warga senior yang memiliki keterbatasan perangkat telepon pintar (*smartphone* dengan spesifikasi rendah, keterbatasan memori internal, atau versi peramban usang) serta keterbatasan kuota internet mandiri. Mengandalkan antarmuka web PWA berbasis browser Next.js tanpa kompensasi biaya data dan tanpa pelatihan tatap muka intensif berisiko memicu penolakan massal. Petugas lapangan akan memilih mencatat di buku tulis dan melapor via grup WhatsApp informal.

### 7.2 Masalah Geospasial pada Pemukiman Padat (Urban Canyon / Gang Sempit)
Kota Depok memiliki karakteristik pemukiman sangat padat di beberapa kecamatan (seperti Pancoran Mas, Sukmajaya, dan Beji) dengan deretan gang sempit. Akurasi GPS pada perangkat seluler warga di area indoor atau beratap seng padat kerap meleset antara 50 hingga 150 meter.
- **Dampak Operasional:** Tim fogging mendatangi gang atau RT yang keliru, memicu konflik sosial antar-warga ("Rumah yang sakit di RT 02, kenapa yang difogging RT 04?").
- **Mitigasi:** Wajibkan Kader mengonfirmasi titik patokan fisik (*landmark*) dan nomor rumah manual pada saat verifikasi lapangan, bukan semata-mata mengandalkan pin koordinat digital.

### 7.3 Risiko Keselamatan & Kesehatan Kerja (K3) Operasional Fogging
Proses penyemprotan insektisida termal (*thermal fogging*) menggunakan senyawa kimia berbahaya (misal organofosfat/malathion atau sintetis piretroid) yang dilarutkan dalam solar/minyak.
- **Celah Sistem:** Sistem DBD Guard mencatat jadwal dan penyelesaian fogging, tetapi **tidak memiliki instrumen keselamatan K3**:
  1. Tidak ada checklist konfirmasi pengosongan area (evakuasi balita, lansia, penderita asma, dan hewan peliharaan) 2 jam sebelum operasi.
  2. Tidak ada pencatatan kelayakan Alat Pelindung Diri (APD: masker respirator, kacamata goggle) bagi operator mesin fogging.
  3. Tidak ada prosedur penanganan limbah sisa insektisida dan pencucian mesin pasca-operasi.
- **Dampak:** Keracunan inhalasi akut pada warga sekitar atau operator fogging, yang berpotensi memicu tuntutan hukum pidana kelalaian kepada Dinas Kesehatan.

### 7.4 Hambatan Koordinasi Wabah Lintas Batas Puskesmas
Vektor nyamuk *Aedes aegypti* memiliki daya terbang efektif 100-200 meter tanpa memedulikan batas administratif wilayah kerja Puskesmas. Apabila suatu kluster kasus berada persis di perbatasan dua kelurahan yang dinaungi oleh Puskesmas berbeda (misalnya perbatasan Kelurahan Pondok Cina dan Kelurahan Kemiri Muka):
- Sistem saat ini mengisolasi akses laporan hanya kepada Puskesmas yang terpetakan pada kelurahan bersangkutan (`mk.puskesmas_id = *locationID`).
- Puskesmas tetangga tidak dapat melihat adanya kasus demam berjarak 50 meter di seberang jalannya.
- Tidak ada fitur *Joint Fogging Coordination*, mengakibatkan penanganan parsial di mana nyamuk hanya berpindah ke pemukiman di sisi jalan sebelah tanpa diberantas tuntas.

---

## 8. Risiko Data & Keputusan (Data & Decision Risk)

### 8.1 Bias Representasi Data akibat Under-Reporting
Heatmap risiko geospasial pada portal publik dan dashboard Dinkes dibangun murni dari data laporan yang masuk via sistem (`citizen_reports`). 
- Pada wilayah dengan masyarakat berpendidikan tinggi dan melek digital, laporan warga akan masuk deras sehingga wilayah tampak berwarna merah pekat (*High Risk*).
- Sebaliknya, pada wilayah kantong pemukiman padat kumuh dengan literasi digital minim, tidak ada warga yang melapor via web sehingga wilayah tampak hijau aman (*Low Risk*), padahal di lapangan terjadi penularan wabah DBD hebat yang tidak terdeteksi.
- **Dampak:** Alokasi logistik insektisida dan jadwal fogging pemerintah kota terdistorsi ke wilayah yang melek teknologi, sementara wilayah rentan yang sesungguhnya terabaikan.

### 8.2 Inkompatibilitas dengan Data Rekam Medis (Ketiadaan Sinkronisasi Faskes)
DBD Guard beroperasi dalam ruang isolasi data yang terpisah dari ekosistem Rekam Medis Elektronik (RME) Puskesmas (SIMPus) maupun Sistem Informasi Manajemen Rumah Sakit (SIMRS) di RSUD Kota Depok dan RS Swasta se-Kota Depok.
- Kasus pasien yang sudah dinyatakan positif uji laboratorium darah (NS1 positif atau trombositopenia parah) di IGD RSUD tidak otomatis mengalir ke DBD Guard.
- Sistem hanya menunggu warga atau keluarganya secara sukarela membuka website dan melapor mandiri.
- Kesenjangan waktu (*lagging*) pelaporan warga rata-rata berkisar 4-7 hari setelah onset demam, melewati periode emas intervensi transmisi awal nyamuk.

### 8.3 Distorsi Logika Agregasi Kejadian Luar Biasa (KLB)
Tabel `klb_configurations` menyimpan parameter kaku `region_scope` ("KELURAHAN" atau "KECAMATAN").
- Pada implementasi kode usecase (`ews_klb_usecase.go:67`):
  ```go
  if areaType != config.RegionScope {
      return nil
  }
  ```
- Jika Dinkes mengonfigurasi skop pemantauan pada level "KELURAHAN", maka akumulasi kasus yang tersebar di 6 kelurahan dalam 1 kecamatan yang sama tidak akan pernah memicu alarm KLB tingkat Kecamatan, meskipun total kasus sekecamatan sudah melebihi 100 orang.
- Sebaliknya, jika dikonfigurasi "KECAMATAN", lonjakan dahsyat terisolasi di 1 RW padat tidak akan memicu alarm Kelurahan sebelum seluruh kecamatan meledak.
- Selain itu, query penghitungan kasus hanya menjaring status bertindakan fogging, sehingga ribuan kasus yang ditangani secara suportif/PSN diabaikan dari indikator epidemiologi kota.

---

## 9. Risiko Teknologi & Infrastruktur (Technology & Infrastructure Risk)

### 9.1 Titik Kegagalan Tunggal (Single Point of Failure / SPOF)
Arsitektur saat ini dijalankan sebagai satu kesatuan kontainer/proses tanpa mekanisme redundansi (*failover*):
- **Basis Data:** PostgreSQL/PostGIS beroperasi pada satu instans penyimpanan tunggal. Tidak ada *read-replica*, tidak ada *standby node*, dan tidak ada skrip pencadangan berkala (`pg_dump` cron) otomatis ke off-site storage. Jika media disk peladen mengalami *bad sector* atau korupsi data, seluruh sistem lumpuh permanen.
- **Backend Service:** Fiber backend berjalan sebagai proses monolitik tunggal tanpa orkestrasi load balancing di belakang reverse proxy ganda.

### 9.2 Risiko Penyimpanan Berkas Lokal (Local File Storage Scalability)
Foto bukti laporan warga dan bukti fogging disimpan langsung pada sistem berkas lokal peladen (`storage/photos/...`):
- Berkas citra resolusi tinggi dari kamera telepon pintar modern berukuran antara 3 MB hingga 8 MB per foto.
- Dengan estimasi 1.000 laporan aktif per bulan disertai masing-masing 2 foto dan 1 foto bukti fogging, konsumsi ruang disk mencapai 15-25 GB per kuartal.
- Tanpa proses kompresi otomatis di sisi backend dan tanpa migrasi ke Object Storage mandiri (seperti MinIO terdistribusi), media penyimpanan peladen akan mengalami kehabisan kapasitas (*disk exhaustion*) yang berakibat pada kegagalan penulisan basis data (*database crash*).

### 9.3 Limitasi Ketergantungan Eksternal Pihak Ketiga (Third-Party Quotas)
Meskipun proyek menetapkan arsitektur "100% Gratis / Open-Source", sistem mengandalkan gerbang komunikasi gratis yang memiliki batasan ketat:
- **Google SMTP / Email Publik:** Memiliki batas kuota ketat (misal 500 email per hari untuk akun Gmail gratis). Saat terjadi KLB di mana ribuan notifikasi dikirimkan serentak, layanan akan langsung diblokir sementara (*rate-limit suspended*), menyebabkan antrean notifikasi macet total.
- **OpenStreetMap Public Tile Server:** Kebijakan penggunaan OSM (*Tile Usage Policy*) secara eksplisit melarang penggunaan aplikasi sektor publik skala kota dengan volume tinggi tanpa caching proxy lokal. Peladen OSM dapat menjatuhkan sanksi pemblokiran IP address peladen Pemkot Depok, menyebabkan antarmuka peta menjadi kosong (*blank map*).

### 9.4 Ketiadaan Pipeline Otomasi Pengujian & Deployment (CI/CD)
Kode sumber belum dilengkapi dengan lingkungan pementasan (*staging environment*), pipeline automasi pengujian (GitHub Actions / GitLab CI), dan kontainerisasi penuh untuk frontend dan backend. Setiap proses pembaruan kode dilakukan manual di peladen produksi, memperbesar risiko *human error* saat kompilasi biner atau migrasi skema basis data di tengah jam sibuk pelayanan.

---

## 10. Risiko Keamanan, Privasi & Kepatuhan (Security & Compliance Risk)

### 10.1 Pelanggaran Undang-Undang Pelindungan Data Pribadi (UU No. 27/2022)
DBD Guard memproses **Data Pribadi Spesifik** (data kesehatan, data riwayat penyakit anak, dan data biometrik/foto kondisi hunian warga) yang diatur secara ketat dalam UU PDP:
1. **Absensi Persetujuan Eksplisit (Informed Consent):** Form pelaporan warga (`/report`) tidak menyertakan kotak centang persetujuan pemrosesan data pribadi dan syarat ketentuan layanan (*Terms of Service & Privacy Policy*) yang mengikat secara hukum.
2. **Ketiadaan Kebijakan Retensi & Penghapusan Data (Right to Erasure):** Tidak ada alur bagi warga untuk memohon penghapusan nomor kontak dan nama mereka pasca-wabah selesai. Data tersimpan permanen di basis data tanpa enkripsi tingkat kolom (*column-level encryption at rest*).
3. **Sanksi Hukum:** Pelanggaran kebocoran data kesehatan pribadi mengancam Pengendali Data Publik (Dinas Kesehatan Kota Depok) dengan sanksi administratif berupa denda hingga miliaran rupiah, ganti rugi perdata, dan tuntutan pidana terhadap oknum yang membocorkan.

### 10.2 Paparan Endpoint Tidak Terotentikasi pada Berkas & PII
- Direktori statis `/storage` yang diekspos melalui `app.Static` memungkinkan pengunduhan foto surat rujukan rumah sakit yang memuat NIK, nama lengkap pasien, tanggal lahir, dan rekam diagnosis klinis tanpa perlu login token JWT.
- Endpoint `/api/v1/reports/phone/:phone` memungkinkan pihak mana pun melakukan pengikisan data (*scraping*) riwayat kesehatan masyarakat Depok hanya dengan mengiterasi basis nomor seluler lokal (misal format `0812xxxxxxxx`).

### 10.3 Pelanggaran Regulasi Pengelolaan Aset Daerah (Permendagri No. 19/2016)
Penyelenggaraan alokasi silang otomatis di mana stok insektisida milik satu Puskesmas dialihkan secara virtual ke Puskesmas lain melanggar prinsip dasar Pengelolaan Barang Milik Daerah:
- Setiap perpindahan barang antar-Pengguna Barang wajib disertai Berita Acara Serah Terima (BAST) fisik yang ditandatangani oleh kedua Kepala Puskesmas bersangkutan dan diketahui oleh Kepala Dinas Kesehatan.
- Mutasi virtual instan di basis data tanpa dokumen penyerta akan diklasifikasikan sebagai *laporan pertanggungjawaban palsu* saat dilakukan audit fisik oleh Inspektorat Kota Depok.

---

## 11. Risiko Kelangsungan Bisnis & Rencana Kontingensi (Business Continuity Risk)

### 11.1 Prosedur Darurat Operasional (System Downtime SOP)
Apabila peladen utama DBD Guard mengalami kelumpuhan teknis (*down*), proses pencegahan wabah DBD Kota Depok **tidak boleh terhenti**. Ditetapkan protokol kontingensi berikut:

```
[SISTEM DOWN > 2 JAM]
         |
         v
[Deklarasi Darurat oleh Kepala Dinas Kesehatan]
         |
         +--> 1. Kanal Warga: Dialihkan ke Call Center Kota Depok 112 / SMS Gateway
         |
         +--> 2. Alur Lapangan: Penggunaan Formulir Kertas Standar PE Kemenkes
         |
         +--> 3. Alur Logistik: Pengeluaran Stok berbasis Buku Register Manual & BAST Fisik
         |
         +--> 4. Komunikasi Cepat: Koordinasi Satgas DBD melalui WhatsApp Group Khusus
         |
         v
[SISTEM PULIH / RESTORE SELESAI]
         |
         v
[Protokol Rekonsiliasi Data Susulan Maksimal 24 Jam]
```

### 11.2 Protokol Pengambilan Keputusan BCP (Business Continuity Decision Matrix)

| Elemen Kedaruratan | Ketentuan Tetap | Penanggung Jawab |
|:---|:---|:---|
| **Kriteria Aktivasi BCP** | Sistem tidak dapat diakses (5xx error / server mati) melebihi 120 menit berturut-turut pada jam kerja (07.00 - 17.00 WIB). | Kepala Seksi Surveilans Dinkes |
| **Otoritas Deklarator** | Pejabat berwenang yang berhak menyatakan sistem "UNAVAILABLE" dan menginstruksikan peralihan ke SOP Kertas Manual. | Kepala Bidang P2P Dinkes |
| **Kanal Informasi Darurat** | Pengumuman resmi melalui saluran darurat radio komando pemkot, surat kawat dinas, dan WhatsApp Group "Koordinator Puskesmas Se-Depok". | Diskominfo / Humas Dinkes |
| **Instrumen Pengganti** | Formulir Penyelidikan Epidemiologi Standar Kemenkes (Form PE-DBD/2020) rangkap 3 berkarbon. | Sanitarian Puskesmas |
| **Batas Waktu Entri Susulan** | Maksimal 24 jam kalender setelah sistem dinyatakan pulih untuk menginput seluruh data manual ke aplikasi. | Operator Data Puskesmas |
| **Verifikasi Integritas** | Validasi silang antara jumlah stok insektisida riil di gudang fisik dengan saldo akhir database pasca-input susulan. | Pengelola Barang Puskesmas |

---

## 12. Analisis Kontrol Internal (Control Analysis)

Evaluasi menyeluruh terhadap efektivitas kontrol internal yang saat ini tertanam di dalam arsitektur DBD Guard:

| Jenis Kontrol | Mekanisme Kontrol yang Ada | Evaluasi Efektivitas | Status Bukti | Celah Kontrol yang Ditemukan |
|:---|:---|:---:|:---:|:---|
| **Preventif** | Rate Limiter pada endpoint publik pelaporan warga (`PublicRateLimiter`). | Efektif | Test Verified | Hanya membatasi serangan banjir traffic IP; tidak mencegah spammer berganti IP atau input manual berulang. |
| **Preventif** | Pengecekan radius duplikat laporan (< 100m, 7 hari). | Parsial | Verified | Hanya menandai `is_possible_duplicate = true`; tidak memblokir atau menggabungkan tiket laporan secara otomatis. |
| **Preventif** | Validasi input koordinat geografis wilayah Indonesia di handler API. | Efektif | Verified | Koordinat (0,0) atau luar negeri ditolak; namun akurasi titik rumah spesifik di gang sempit tetap tidak tervalidasi. |
| **Preventif** | Isolasi data wilayah Puskesmas (`LocationID scoping`). | Parsial | Verified | Puskesmas hanya melihat kelurahan binaan; namun membutakan koordinasi kluster wabah yang berada di perbatasan wilayah. |
| **Detektif** | Pencatatan rekam jejak audit (`AuditLog`) pada transaksi krusial. | Efektif | Verified | Aksi PE, Approval Fogging, dan Mutasi Stok tercatat; namun belum diproteksi terhadap risiko penghapusan oleh Super Admin. |
| **Detektif** | SLA Cron Job pemantau laporan kedaluwarsa 72 jam (`sla_cron.go`). | Efektif | Test Verified | Berhasil mendeteksi laporan macet; namun eskalasi hanya mengubah status tanpa memicu penugasan satgas darurat aktif. |
| **Korektif** | Fitur pelepasan paksa penugasan kader (`ForceUnassignReport`). | Parsial | Verified | Supervisor Puskesmas dapat membebaskan tiket laporan; namun aksi ini bersifat reaktif dan bergantung pada inisiatif manusia. |
| **Korektif** | Alur Banding Warga (`SubmitAppeal`) untuk laporan yang ditolak hoaks. | Efektif | Test Verified | Memberikan hak jawab bagi warga; dibatasi 1 kali; diteruskan ke tingkat Puskesmas untuk evaluasi ulang. |
| **Otomatis** | Perhitungan formula rekomendasi tindakan DSS. | Cacat Logika | Verified | Mengunci status laporan tanpa fogging secara sepihak, menghilangkan diskresi pertimbangan klinis dokter. |
| **Manual** | Tinjauan dan persetujuan pengajuan fogging oleh Kepala Puskesmas. | Parsial | Verified | Memberikan kontrol manusia; namun menjadi titik kemacetan (*bottleneck*) fatal jika Kepala Puskesmas berhalangan hadir. |
| **SoD** | Pembagian wewenang antara Kader (PE) dan Kepala Puskesmas (Approval). | Cacat Logika | Verified | Dilanggar pada alur *Takeover*: satu orang petugas Puskesmas dapat merangkap sebagai investigator dan validator sekaligus. |

---

## 13. Rencana Perlakuan Risiko (Risk Treatment Plan)

Rencana aksi konkret mitigasi risiko disusun berdasarkan skala prioritas penanganan untuk memandu tim pengembang dan pimpinan Dinas Kesehatan Kota Depok:

### 13.1 Prioritas Kritis (P0 — Wajib Diselesaikan Segera / Blocker Rilis Produksi)

1. **P0-1: Restrukturisasi Total Modul Logistik "Alokasi Silang" (Terkait RSK-01 / BW-01)**
   - *Tindakan:* Hapus mekanisme pemotongan stok otomatis dari donor warehouse di `inventory_usecase.go`. Ganti dengan penerbitan tiket `StockRequest` resmi ke Dinas Kesehatan. Pemindahan stok antar-Puskesmas wajib melalui modul `StockTransfer` dua arah dengan Berita Acara Serah Terima digital ber-OTP.
   - *Penanggung Jawab:* Technical Lead & Pengelola Barang Dinkes.
   - *Target Waktu:* 3 Hari Kerja.
   - *Indikator Keberhasilan:* Saldo stok gudang tidak pernah berkurang tanpa bukti transaksi persetujuan dari kepala gudang pemilik stok.

2. **P0-2: Isolasi Keamanan Direktori Berkas Medis (Terkait RSK-03 / BW-08)**
   - *Tindakan:* Hapus `app.Static("/storage", "./storage")` dari router publik. Buat endpoint terproteksi `/api/v1/storage/photos/:id` yang memvalidasi header token JWT dan hak akses peran pemohon (hanya peran Dinkes, Puskesmas wilayah terkait, dan pelapor pemilik berkas). Implementasikan generator tautan sementara bertanda tangan (*Pre-signed URL*) jika menggunakan MinIO.
   - *Penanggung Jawab:* Backend Developer & Security Specialist.
   - *Target Waktu:* 2 Hari Kerja.
   - *Indikator Keberhasilan:* Berkas foto diagnosis rumah sakit menghasilkan kode status `401 Unauthorized` jika diakses tanpa token kredensial.

3. **P0-3: Restorasi Diskresi Klinis Dokter pada Formula DSS (Terkait RSK-02 / BW-03)**
   - *Tindakan:* Hapus penutupan status otomatis `CLOSED_NO_FOGGING` di usecase PE. Setiap hasil PE wajib diteruskan ke antrean validasi dokter/Kepala Puskesmas. Dokter diberikan hak *Override Recommendation* dengan kewajiban mengisi alasan pertimbangan epidemiologis yang tercatat di audit log.
   - *Penanggung Jawab:* Tim Analis Bisnis & Dokter Surveilans Dinkes.
   - *Target Waktu:* 2 Hari Kerja.
   - *Indikator Keberhasilan:* Kepala Puskesmas dapat menyetujui fogging selektif pada hasil PE non-jentik disertai catatan pertimbangan klinis.

4. **P0-4: Penegakan Pemisahan Tugas (Segregation of Duties) pada Takeover (Terkait RSK-08 / BW-04)**
   - *Tindakan:* Tambahkan validasi programatik pada `ValidatePEResult`:
     ```go
     if peResult.OfficerID == validatorID {
         return errors.New("pelanggaran pemisahan tugas: pemeriksa lapangan dilarang memvalidasi laporannya sendiri")
     }
     ```
   - *Penanggung Jawab:* Backend Developer.
   - *Target Waktu:* 1 Hari Kerja.
   - *Indikator Keberhasilan:* Upaya validasi mandiri ditolak sistem dengan pesan error 403 Forbidden.

5. **P0-5: Konfigurasi Cadangan Basis Data Otomatis & Disaster Recovery (Terkait RSK-09)**
   - *Tindakan:* Pasang cron script pencadangan berkala `pg_dump` setiap 6 jam ke media penyimpanan sekunder terpisah dari kontainer Docker. Siapkan dan uji coba skrip restorasi berkala (*mock recovery test*).
   - *Penanggung Jawab:* DevOps / Tim Infrastruktur Diskominfo.
   - *Target Waktu:* 1 Hari Kerja.
   - *Indikator Keberhasilan:* File dump terarsip rapi harian dan berhasil dipulihkan pada server cadangan dalam waktu < 30 menit.

---

### 13.2 Prioritas Tinggi (P1 — Selesai Sebelum Pilot Project Lapangan)

1. **P1-1: Validasi Kewajiban Unggah Foto Bukti Fogging (Terkait RSK-11 / BW-06)**
   - *Tindakan:* Tambahkan pengecekan integritas di `UpdateExecutionStatus`: jika status diubah menjadi `COMPLETED`, lakukan query ke tabel `fogging_evidences`. Tolak permintaan jika jumlah bukti foto = 0.
   - *Penanggung Jawab:* Backend Developer.
   - *Target Waktu:* 1 Hari Kerja.

2. **P1-2: Pembangunan Antarmuka Administrasi Relawan Kader RW (Terkait RSK-15 / BW-10)**
   - *Tindakan:* Bangun modul "Manajemen Kader Wilayah" di dashboard Puskesmas. Petugas Puskesmas dapat membuatkan akun Kader baru per RW, mereset password, dan menonaktifkan kader yang demisioner secara instan.
   - *Penanggung Jawab:* Frontend & Backend Developer.
   - *Target Waktu:* 3 Hari Kerja.

3. **P1-3: Pembenahan Kueri Agregasi EWS dan Threshold Wilayah (Terkait RSK-06)**
   - *Tindakan:* Perbaiki kueri `CountCasesByArea` agar menjaring seluruh laporan terkonfirmasi (`VERIFIED`, `CLOSED_NO_FOGGING`, `COMPLETED`). Ubah konfigurasi EWS agar memeriksa ambang batas Kelurahan dan Kecamatan secara paralel simultan.
   - *Penanggung Jawab:* Backend Developer & Epidemiolog.
   - *Target Waktu:* 2 Hari Kerja.

4. **P1-4: Pembuatan Modul Geofencing Presensi Pemeriksaan PE (Terkait RSK-04 / BW-02)**
   - *Tindakan:* Manfaatkan Geolocation API browser untuk mencocokkan koordinat Kader saat klik "Submit PE" dengan koordinat rumah pelapor. Berikan peringatan jika jarak > 200 meter.
   - *Penanggung Jawab:* Frontend Developer.
   - *Target Waktu:* 2 Hari Kerja.

5. **P1-5: Penyusunan Regulasi Daerah (Payung Hukum Wajib Pakai) (Terkait RSK-07)**
   - *Tindakan:* Dinkes menyusun draf Surat Keputusan (SK) Kepala Dinas Kesehatan dan advokasi Peraturan Walikota (Perwal) tentang Penyelenggaraan Sistem Pelaporan DBD Digital Terpadu Kota Depok.
   - *Penanggung Jawab:* Bagian Hukum & Manajemen Dinkes.
   - *Target Waktu:* 14 Hari Kerja.

---

### 13.3 Prioritas Sedang (P2 — Pemeliharaan & Penyempurnaan Berkelanjutan)

1. **P2-1: Penerapan Kuota Maksimum Klaim Kader (Anti-Task Hoarding) (Terkait RSK-12)**
   - *Tindakan:* Batasi seorang akun Kader hanya dapat memiliki maksimal 3 tugas verifikasi/PE aktif secara bersamaan.
   - *Target Waktu:* 1 Hari Kerja.

2. **P2-2: Kebijakan Retensi & Pemenuhan Hak Hapus Data Pribadi Warga (Terkait RSK-17)**
   - *Tindakan:* Bangun cronjob anonimisasi otomatis: hapus nomor kontak dan nama pelapor dari tabel `reporters` 30 hari setelah kasus berstatus `COMPLETED` atau `CLOSED_NO_FOGGING`.
   - *Target Waktu:* 2 Hari Kerja.

3. **P2-3: Caching Reverse-Proxy Peladen Tile Peta Mandiri (Terkait RSK-18)**
   - *Tindakan:* Pasang reverse-proxy Nginx dengan disk caching untuk tile OpenStreetMap di peladen lokal Pemkot Depok guna mencegah pemblokiran IP oleh OSM foundation.
   - *Target Waktu:* 2 Hari Kerja.

4. **P2-4: Mekanisme Delegasi Otorisasi Kepala Puskesmas (Terkait Bottleneck Approval)**
   - *Tindakan:* Fitur pelimpahan wewenang sementara (*Acting Approval / Plt*) kepada dokter fungsional yang ditunjuk saat Kepala Puskesmas dinas luar/cuti.
   - *Target Waktu:* 2 Hari Kerja.

---

## 14. Indikator Risiko Utama & Pemantauan (Key Risk Indicators / KRI)

Tabel indikator peringatan dini (*early warning metrics*) yang wajib dipantau mingguan oleh Pengelola Sistem dan pimpinan Dinas Kesehatan Kota Depok:

| Indikator (KRI) | Definisi Formula Pengukuran | Batas Waspada (Warning) | Batas Bahaya (Critical) | Frekuensi Pantau | Tindakan Jika Melewati Batas Bahaya |
|:---|:---|:---:|:---:|:---:|:---|
| **KRI-01: Tingkat Pelanggaran SLA Laporan** | $\frac{\text{Laporan melewati SLA 72 Jam}}{\text{Total Laporan Masuk}} \times 100\%$ | > 5% per minggu | > 15% per minggu | Mingguan | Rapat koordinasi eskalasi Puskesmas; aktivasi Satgas Jemput Bola Dinkes. |
| **KRI-02: Frekuensi Alokasi Silang Logistik** | Jumlah insiden alokasi silang otomatis antar-Puskesmas dalam 1 bulan. | > 2 kali sebulan | > 5 kali sebulan | Bulanan | Audit investigasi fisik gudang farmasi; evaluasi kuota distribusi logistik tahunan. |
| **KRI-03: Rasio Override Rekomendasi DSS** | $\frac{\text{Approval Fogging Berbeda dari DSS}}{\text{Total Approval Fogging}} \times 100\%$ | > 15% per bulan | > 30% per bulan | Bulanan | Peninjauan ulang formula bobot DSS oleh tim epidemiolog bersama Kemenkes. |
| **KRI-04: Penolakan Laporan oleh Filter Spam** | $\frac{\text{Laporan Ditolak Otomatis Sistem (Spam)}}{\text{Total Laporan Masuk}} \times 100\%$ | > 5% per minggu | > 12% per minggu | Mingguan | Verifikasi manual sampel laporan spam; evaluasi ambang batas radius waktu filter. |
| **KRI-05: Tingkat Inaktivitas Akun Relawan** | $\frac{\text{Akun Kader Tidak Login > 30 Hari}}{\text{Total Akun Kader Terdaftar}} \times 100\%$ | > 20% | > 40% | Bulanan | Rekonsiliasi data kader bersama lurah/RW; pembukaan pendaftaran relawan baru. |
| **KRI-06: Antrean Persetujuan Kepala Puskesmas** | Rata-rata durasi antrean jadwal menunggu approval Kepala Puskesmas. | > 12 jam | > 24 jam | Harian | Kirim notifikasi darurat SMS/WA ke Kepala Puskesmas; alihkan ke Plt/Wakil. |
| **KRI-07: Konsumsi Kapasitas Penyimpanan Foto** | Persentase keterisian ruang disk direktori berkas unggahan bukti. | > 70% kapasitas | > 85% kapasitas | Mingguan | Pembersihan arsip lama terkompresi; penambahan volume storage server. |
| **KRI-08: Selisih Pembukuan Logistik Fisik vs Web** | $\frac{\text{Selisih Fisik Stok vs Saldo Sistem}}{\text{Saldo Sistem}} \times 100\%$ | > 2% selisih | > 5% selisih | Bulanan | Lakukan *Stock Opname* total; laporkan ke Inspektorat jika ada dugaan penggelapan. |

---

## 15. Keterlacakan Terhadap Hasil Deep Bug Audit (Traceability to Deep Bug Audit)

Tabel pemetaan ini membuktikan bahwa temuan teknis pada dua audit mendalam sebelumnya (*Deep Bug Audit Phase 1* dan *Bug Review Report*) bukanlah sekadar cacat kode isolatif, melainkan representasi dari kelemahan sistemik tata kelola dan proses bisnis:

| ID Bug Asal | Deskripsi Cacat Teknis (Code Level) | Celah Kontrol yang Ditunjukkan | Transformasi Menjadi Risiko Bisnis / Operasional | ID Risiko Terkait |
|:---|:---|:---|:---|:---:|
| **BUG-01 (DBA-1)** | Endpoint Appeal menolak status `REJECTED_HOAX` dan UI mengabaikan batas 1x banding. | Kontrol hak banding warga tidak konsisten dengan status penolakan kader. | Warga yang difitnah laporannya hoaks kehilangan hak pemulihan penanganan; kluster nyata terabaikan; timbul sengketa warga vs kader. | **RSK-04** |
| **BUG-02 (DBA-1)** | Jadwal fogging crash saat array `usedItems` kosong; duplikasi entri audit logistik. | Ketiadaan validasi rantai proses logistik lapangan vs penutupan tugas. | Tim fogging tidak bisa melaporkan tugas selesai di aplikasi; tim memilih jalur offline; pembukuan obat Puskesmas mandek. | **RSK-11** |
| **BUG-03 (DBA-1)** | Puskesmas menerima transfer stok otomatis gagal karena payload boolean `is_accepted`. | Kontrol penerimaan barang logistik cacat kontrak data frontend-backend. | Stok obat kiriman Dinkes ditolak sistem dan kembali secara fiktif ke Dinkes padahal fisik obat sudah di Puskesmas; krisis logistik. | **RSK-01** |
| **BUG-05 (DBA-1)** | Otorisasi usecase PE menolak eksekusi `Takeover` resmi oleh petugas Puskesmas (403). | Alur delegasi dan pelimpahan wewenang terputus pada lapisan bisnis. | Laporan warga yang ditelantarkan Kader RT tidak dapat diselamatkan oleh dokter Puskesmas; kegagalan SLA penanganan wabah. | **RSK-08** |
| **BUG-08 (DBA-1)** | Pemicuan KLB Kecamatan membypass data dan di-hardcode hanya untuk Kelurahan. | Ketiadaan sinkronisasi agregasi bertingkat pada arsitektur surveilans. | Kegagalan peringatan dini KLB tingkat kota; Dinkes tidak menyadari wabah telah lintas batas kelurahan; keterlambatan deklarasi bencana. | **RSK-06** |
| **BUG-09 (DBA-1)** | Persetujuan permintaan stok Puskesmas menambah saldo virtual tanpa memotong gudang Dinkes. | Pelanggaran prinsip *Single-Entry/Double-Entry Accounting* logistik. | Terjadinya pencetakan stok fiktif dari ketiadaan (*inflationary inventory*); pembukuan neraca aset farmasi Dinkes korup total. | **RSK-01** |
| **BUG-10 (DBA-1)** | Mode lapor anonim di frontend tetap mewajibkan nama dan nomor telepon pelapor. | Kontrol privasi dan perlindungan PII warga bocor sejak antarmuka terdepan. | Warga enggan melapor karena takut identitasnya dihakimi warga lain; pelanggaran prinsip kepatuhan hukum privasi UU PDP. | **RSK-10** |
| **LOGIC-01 (BRR)** | Usecase mengizinkan pembuatan multi-jadwal fogging pada 1 permintaan yang sama. | Ketiadaan constraint integritas relasi 1:1 antara permintaan dan jadwal. | Pemborosan jadwal penyemprotan ganda pada satu lokasi; pembagian logistik membengkak dua kali lipat untuk kasus yang sama. | **RSK-05** |
| **LOGIC-04 (BRR)** | Tidak ada pengecekan duplikasi hasil PE pada satu laporan warga. | Kontrol idempotensi pada entri formulir penyelidikan tidak tersedia. | Dua kader berbeda menginput hasil PE bertolak belakang pada laporan yang sama; memicu anomali ganda pada antrean approval. | **RSK-04** |
| **LOGIC-06 (BRR)** | Usecase otentikasi tidak memeriksa status `ActiveStatus = false` akun saat login. | Kegagalan penegakan sanksi penonaktifan kredensial pengguna. | Kader atau staf Puskesmas yang telah dipecat/berhenti tetap dapat masuk ke sistem dan mengakses data sensitif kesehatan warga. | **RSK-15** |
| **SEC-01 (BRR)** | Handler menerima ID aktor (`actor_id`) langsung dari payload JSON request body. | Kegagalan autentikasi berbasis konteks klaim token JWT. | Pengguna dapat memalsukan identitas pengambil keputusan (*identity spoofing*) dengan mengirimkan UUID pejabat lain pada request API. | **RSK-08** |
| **CONC-01 (BRR)** | Pengecekan saldo stok inventaris dilakukan di luar transaksi atomik database. | Kerentanan persaingan akses konkuren (*race condition*) pada saldo barang. | Dua jadwal fogging bersamaan menarik sisa stok terakhir, menghasilkan saldo gudang minus secara fisik dan digital. | **RSK-01** |

---

## 16. Isu Terbuka & Keputusan Strategis yang Diperlukan (Open Issues / Decisions Required)

Berikut adalah daftar keputusan strategis non-teknis yang memerlukan kesepakatan tingkat tinggi antara Tim Pengembang, Kepala Dinas Kesehatan Kota Depok, dan Bagian Hukum Pemkot Depok:

1. **Kebijakan Hukum Alokasi Silang Logistik:**
   - *Pilihan A:* Menghapus total alokasi silang otomatis. Jika Puskesmas kehabisan stok, jadwal ditunda hingga stok dikirim resmi dari Gudang Farmasi Dinkes.
   - *Pilihan B:* Mempertahankan pinjam-meminjam antar-Puskesmas terdekat, namun mewajibkan fitur penandatanganan elektronik Berita Acara Serah Terima (BAST) fisik dan persetujuan pejabat pengelola barang kedua Puskesmas sebelum sistem memotong stok.
   - *Rekomendasi Audit:* **Pilihan B** untuk fleksibilitas lapangan, namun dengan kontrol ketat dokumen BAST digital.

2. **Hierarki Penentuan Fogging (Algoritma vs Klinis):**
   - *Isu:* Apakah rekomendasi DSS bersifat *Binding* (mengikat) atau *Advisory* (hanya saran)?
   - *Rekomendasi Audit:* Posisikan DSS murni sebagai *Advisory*. Kepala Puskesmas wajib memegang hak veto akhir (Override) untuk menyetujui fogging pada laporan dengan PE negatif jika terdapat indikasi epidemiologis lain, dengan kewajiban mencatatkan alasan klinis formal.

3. **Status Hukum Kerahasiaan Data Kesehatan Publik:**
   - *Isu:* Bagaimana menyelaraskan transparansi Heatmap Publik dengan kepatuhan UU PDP Pasal 43?
   - *Rekomendasi Audit:* Batasi agregasi visualisasi publik minimal pada level Kelurahan (tidak boleh menampilkan koordinat presisi rumah tangga atau nomor RW secara telanjang). Hapus seluruh fitur pelacakan nomor telepon publik terbuka.

4. **Alokasi Anggaran Operasional (OPEX) Keberlanjutan Sistem:**
   - *Isu:* Penggunaan server gratis (SMTP, FCM, Docker lokal) tidak memadai untuk produksi skala kota.
   - *Rekomendasi Audit:* Dinkes Depok wajib mengajukan pos anggaran pemeliharaan sistem informasi dalam APBD Perubahan, mencakup: Cloud Hosting bersertifikasi ISO 27001, kuota gateway notifikasi resmi, dan insentif paket data relawan Jumantik.

5. **Penetapan Payung Hukum Wajib Pakai:**
   - *Isu:* Puskesmas dan Kader rentan mengabaikan sistem jika tidak ada instruksi kedinasan.
   - *Rekomendasi Audit:* Penerbitan Surat Edaran / Instruksi Walikota Depok tentang Kewajiban Pencatatan Laporan dan Intervensi Fogging DBD secara Terpadu melalui Aplikasi DBD Guard sebelum peluncuran resmi.

---

## 17. Penilaian Akhir & Kesimpulan (Final Assessment)

Berdasarkan seluruh temuan dalam Risk Management Audit ini, ditarik kesimpulan komprehensif untuk menjawab 16 pertanyaan fundamental proyek:

1. **Apa risiko terbesar pada bisnis DBD Guard?**  
   Risiko terbesar adalah **kegagalan adopsi dan hilangnya kepercayaan institusional** akibat data lapangan yang dapat dimanipulasi, pembukuan inventaris yang melanggar hukum keuangan daerah, serta sistem yang berpotensi memvonis salah kondisi wabah di lapangan.

2. **Apa kelemahan proses bisnisnya?**  
   Proses bisnis terlalu berasumsi bahwa seluruh aktor (warga, relawan, petugas) memiliki integritas sempurna, akses internet stabil, dan disiplin tinggi, tanpa menyiapkan kontrol pencegahan kecurangan, kontrol geofencing, atau alur fallback manual jika sistem macet.

3. **Apa business loophole yang dapat dimanfaatkan user?**  
   Celah "Alokasi Silang Otomatis" untuk mengaburkan saldo gudang, celah penutupan fogging tanpa foto bukti, dan celah pengambilalihan tugas PE yang dapat divalidasi sendiri oleh satu aktor yang sama.

4. **Apa proses yang dapat dimanipulasi?**  
   Input hasil PE (jentik dan penderita demam), pencatatan angka liter pemakaian insektisida, dan pembatalan laporan warga yang belum diverifikasi.

5. **Apa proses yang terlalu bergantung pada kejujuran user?**  
   Penyelidikan Epidemiologi (PE) oleh Kader. Sistem mempercayai begitu saja apakah kader benar-benar memeriksa bak mandi dalam radius 100 meter atau hanya mengisi formulir dari rumah.

6. **Apa proses yang tidak memiliki accountability yang jelas?**  
   Penolakan pemakaian inventaris oleh Puskesmas yang tidak memulihkan saldo fisik gudang, serta pembatalan laporan publik yang dapat dieksekusi tanpa verifikasi identitas pemilik tiket.

7. **Apa risiko jika data salah tetapi sistem tetap terlihat normal?**  
   Terjadinya wabah tersembunyi (*silent outbreak*). Sistem menampilkan zona hijau aman karena laporan ditolak atau ditandai non-fogging, sementara warga di lapangan mengalami komplikasi dengue berat tanpa pernah terdeteksi oleh Dinas Kesehatan.

8. **Apa risiko jika stakeholder tidak menggunakan sistem?**  
   Sistem menjadi investasi mubazir (*dead capital*). Penanganan DBD di Kota Depok kembali ke metode kuno berbasis surat manual dan chat WhatsApp yang lambat, parsial, dan tidak terukur.

9. **Apa risiko sustainability dan biaya operasional?**  
   Sistem terancam mati pasca-serah terima pengembang karena ketiadaan alokasi dana pemeliharaan server, kuota notifikasi komersial, perpanjangan domain, dan ketergantungan pada tim pengembang awal.

10. **Apa risiko jika orang tertentu tidak tersedia?**  
    Kemacetan persetujuan (*approval bottleneck*) pada Kepala Puskesmas yang dapat menghentikan seluruh operasi fogging jika pejabat bersangkutan sedang dinas luar atau sakit, tanpa ada mekanisme pendelegasian wewenang otomatis.

11. **Apa risiko jika sistem down?**  
    Kelumpuhan total layanan respons wabah karena belum tersedianya SOP darurat formulir fisik Kemenkes terpadu dan ketiadaan jalur komando cadangan yang terstandardisasi.

12. **Apa risiko dari dependency pihak eksternal?**  
    Pemblokiran layanan oleh OpenStreetMap akibat pelanggaran kuota akses gratis, serta kegagalan transmisi notifikasi email/FCM gratis saat terjadi lonjakan kasus berskala masif.

13. **Apa kelemahan yang ditunjukkan oleh hasil dua Deep Bug Audit?**  
    Kedua audit sebelumnya membuktikan bahwa sistem memiliki kerapuhan mendasar pada **konsistensi transaksi basis data, kontrak data frontend-backend, tata kelola akun pengguna, dan isolasi privasi data medis**.

14. **Apa control gap yang paling penting?**  
    Ketiadaan kontrol otentikasi pada direktori berkas medis (`/storage`), ketiadaan kontrol pemisahan tugas (*Segregation of Duties*) pada takeover, dan ketiadaan verifikasi fisik atas pemotongan stok logistik.

15. **Apa risiko yang belum terlihat hanya dari sisi teknis?**  
    Risiko yuridis formal: ancaman sanksi pidana dan denda UU Perlindungan Data Pribadi terhadap Pemkot Depok, temuan maladministrasi barang milik daerah oleh BPK, serta potensi konflik sosial horizontal antar-warga akibat salah sasaran lokasi fogging.

### Rekomendasi Status Kelayakan Sistem:
Sistem DBD Guard **BELUM LAYAK (CONDITIONAL HOLD)** untuk diluncurkan ke tingkat operasional penuh Pemerintah Kota Depok sebelum **5 Tindakan Prioritas Kritis (P0-1 s/d P0-5)** diselesaikan dan diuji coba secara ketat bersama perwakilan Dinas Kesehatan dan perwakilan Puskesmas percontohan.

---
*Dokumen Audit Manajemen Risiko ini disusun secara independen, formal, dan objektif sebagai landasan pengambilan keputusan strategis dan penyempurnaan sistem DBD Guard Kota Depok.*
