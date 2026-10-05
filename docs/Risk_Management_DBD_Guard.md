# Dokumen Manajemen Risiko (Risk Management Document)

**Proyek:** DBD Guard — Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi Kota Depok  
**Versi:** 3.0 — Diperbarui 24 September 2026  
**Status:** Fase Pengembangan Aktif (MVP Backend & Frontend)  
**Klasifikasi:** Dokumen Internal — Untuk Keperluan Presentasi Stakeholder

> Dokumen ini merupakan revisi menyeluruh dari versi sebelumnya. Manajemen risiko DBD Guard tidak hanya mencakup kegagalan teknis, tetapi juga risiko bahwa sistem tidak digunakan sesuai tujuan, menghasilkan data yang tidak dapat dipercaya, mengganggu proses kerja lapangan, menciptakan ketergantungan pada individu tertentu, menghasilkan rekomendasi yang salah, menimbulkan kebingungan mengenai akuntabilitas keputusan, atau kehilangan kepercayaan pemangku kepentingan.

---

## 1. Pendahuluan

Dokumen ini disusun untuk mengidentifikasi, mengevaluasi, dan menyusun strategi mitigasi terhadap seluruh potensi risiko yang dapat memengaruhi keberhasilan **DBD Guard** sebagai sistem informasi kesehatan publik Kota Depok.

- **Tujuan Utama:** Memastikan keandalan sistem respons cepat wabah DBD sekaligus menjamin keberlanjutan operasional jangka panjang.
- **Prinsip Dasar:** Risiko operasional dan bisnis harus dinilai berdampingan dengan risiko teknologi, karena sistem mendukung proses pelaporan warga, penyelidikan epidemiologi (PE), validasi lapangan, manajemen inventori, dan pengelolaan fogging secara terintegrasi.

---

## 2. Ringkasan Status Komponen Sistem

| Komponen Sistem | Status | Keterangan |
|:---|:---:|:---|
| Backend (Go/Fiber) | Aktif | Seluruh domain utama (Report, PE, Fogging, Inventory, EWS) terimplementasi. |
| Database (PostGIS) | Aktif | Berjalan melalui Docker Compose. |
| Frontend (Next.js) | Aktif | Enam modul Dashboard telah beroperasi. |
| Otentikasi & Keamanan | Aktif | JWT, RBAC 7 peran, dan Rate Limiter berjalan. |
| Layanan Notifikasi | Parsial | Infrastruktur siap; masih mode simulasi, kredensial produksi belum dikonfigurasi. |
| Penyimpanan Berkas | Belum | Masih penyimpanan lokal; MinIO belum diimplementasikan. |
| Aplikasi Seluler | Belum | Future Improvement sesuai PRD. |
| Load Testing | Belum | Belum pernah dilakukan pengujian beban. |
| CI/CD Pipeline | Belum | Build masih dilakukan secara manual. |

---

## 3. Kategori Risiko

Sistem ini dipantau melalui tujuh area risiko:

1. **Business & Governance (BR)** — Adopsi, tata kelola, anggaran, akuntabilitas, kepercayaan publik.
2. **Operational & Epidemiologi (OP)** — Proses lapangan, alur kerja, SDM, respons wabah.
3. **Technology & Infrastructure (TC)** — Peladen, basis data, penyimpanan, skalabilitas.
4. **Security & Privacy (SC)** — Keamanan siber, privasi data medis, hak warga.
5. **Compliance & Regulatory (CP)** — SOP Kemenkes, regulasi Pemkot, K3 Fogging, UU PDP.
6. **Development & Delivery (DV)** — Testing, deployment, pipeline, kualitas kode.
7. **Integration & External Dependency (IE)** — SatuSehat, penyedia notifikasi, ekosistem nasional.

**Legenda Status:** `Telah Dimitigasi` | `Parsial` | `Belum Dimulai` | `Perlu Review`

---

## 4. Register Risiko

> **Keterangan Kolom:**
> - **Inherent Risk:** Tingkat risiko sebelum ada mitigasi apapun.
> - **Residual Risk:** Tingkat risiko yang tersisa setelah mitigasi diterapkan.
> - **Risk Owner:** Pihak yang bertanggung jawab memastikan mitigasi berjalan.
> - **KRI / Pemicu:** Indikator peringatan dini sebelum risiko menjadi insiden.

---

### A. Business & Governance

*Tantangan strategis terkait tata kelola, keberlanjutan, anggaran, akuntabilitas, dan kepercayaan pemangku kepentingan.*

| ID | Deskripsi Risiko | Dampak | Probabilitas | Inherent Risk | Risk Owner | Status | Strategi Mitigasi | KRI / Pemicu | Residual Risk |
|:---|:---|:---:|:---:|:---:|:---|:---:|:---|:---|:---:|
| **BR-01** | **Adopsi Pengguna Rendah**<br>Warga, Kader, atau Puskesmas tidak menggunakan sistem secara konsisten dan tetap mengandalkan pelaporan manual. | Sangat Tinggi | Tinggi | Sangat Tinggi | Product Owner / Dinkes | Parsial | Pengembangan UI/UX berbasis masukan pengguna lapangan; program sosialisasi dan pelatihan berbasis peran; monitoring tingkat penggunaan aktif. | Persentase laporan masuk via sistem <70%; persentase Kader aktif <60%. | Sedang |
| **BR-02** | **Ketidakjelasan Governance & Ownership**<br>Tidak ada pihak yang secara eksplisit bertanggung jawab atas operasional, keputusan bisnis, keamanan data, dan pemeliharaan sistem pasca-serah terima. | Tinggi | Sedang | Tinggi | Manajemen Dinkes | Belum Dimulai | Penetapan struktur kepemilikan sistem: Product Owner, System Owner, Data Owner, Security Owner, dan Operational Owner. Setiap perubahan aturan bisnis harus memiliki pihak yang menyetujui. | Tidak ada PIC yang jelas untuk isu sistem; keputusan teknis selalu bergantung pada tim developer. | Sedang |
| **BR-03** | **Ketergantungan pada Tim Pengembang Awal**<br>Setelah proyek selesai, hanya tim developer awal yang memahami arsitektur dan cara menjalankan sistem. | Tinggi | Sedang | Tinggi | Technical Owner / Dinkes IT | Belum Dimulai | Penyusunan dokumentasi arsitektur, API, basis data, deployment, dan runbook operasional. Sesi Knowledge Transfer berkala dengan tim IT Dinkes. Serah terima kepemilikan repositori dan kredensial. | Jumlah staf internal yang mampu melakukan deployment/recovery kurang dari 2 orang. | Rendah |
| **BR-04** | **Ketidakcukupan Anggaran Operasional (OPEX)**<br>Tidak ada alokasi dana APBD yang memadai untuk peladen, storage, domain, pemeliharaan, keamanan, dan dukungan teknis rutin. | Sangat Tinggi | Sedang | Sangat Tinggi | Manajemen Dinkes | Belum Dimulai | Penyusunan estimasi OPEX tahunan yang mencakup seluruh komponen (infrastruktur, backup, notifikasi, monitoring, keamanan, maintenance, disaster recovery). Advokasi penganggaran APBD khusus. | Tidak ada pos anggaran pemeliharaan sistem dalam APBD tahun berikutnya. | Sedang |
| **BR-05** | **Resistensi Organisasi terhadap Perubahan**<br>Puskesmas dan Kader tetap menjalankan proses manual meskipun sistem telah tersedia, menciptakan dua alur data paralel yang tidak sinkron. | Tinggi | Tinggi | Sangat Tinggi | Operational Owner / Dinkes | Belum Dimulai | Penerbitan payung hukum (SK Kepala Dinas atau Peraturan Walikota) yang mewajibkan penggunaan sistem. Penetapan SOP resmi dan jalur eskalasi ketika workflow tidak dijalankan. | Jumlah laporan manual yang harus diinput ulang ke sistem lebih dari 10% per bulan. | Sedang |
| **BR-06** | **Kualitas Data Lapangan Rendah**<br>Pengguna memasukkan data yang tidak lengkap, salah, terlambat, atau tidak konsisten (lokasi salah, identitas tidak lengkap, hasil PE tidak diisi). | Sangat Tinggi | Tinggi | Sangat Tinggi | Data Owner / Operational Owner | Parsial | Validasi wajib pada field kritis; master data untuk input lokasi; audit trail; dashboard kualitas data; rekonsiliasi berkala; monitoring anomali data. | Persentase laporan dengan data tidak lengkap lebih dari 5%; jumlah record tanpa pembaruan status lebih dari 15%. | Sedang |
| **BR-07** | **DSS Menghasilkan Rekomendasi yang Tidak Tepat**<br>Sistem merekomendasikan tindakan yang salah akibat data input buruk, business rule yang keliru, atau formula yang tidak sesuai kondisi lapangan. | Sangat Tinggi | Sedang | Tinggi | Business Owner + Technical Owner | Parsial | DSS diposisikan sebagai alat bantu keputusan, bukan pengganti keputusan petugas. Simpan semua input, rule aktif, nilai perhitungan, rekomendasi, dan keputusan final beserta aktor dan timestamp untuk keperluan audit. | Frekuensi override manual terhadap rekomendasi DSS lebih dari 30% dalam periode tertentu. | Sedang |
| **BR-08** | **Ketidakjelasan Akuntabilitas Keputusan**<br>Ketika rekomendasi sistem berbeda dengan keputusan petugas, tidak jelas siapa yang bertanggung jawab dan bagaimana keputusan tersebut dapat diaudit. | Tinggi | Sedang | Tinggi | Business Owner / Dinkes | Parsial | Pemisahan eksplisit antara "System Recommendation" dan "Human Decision" dalam basis data. Setiap keputusan yang berbeda dari rekomendasi sistem wajib menyertakan alasan override dan identitas pengambil keputusan. | Tidak ada rekaman alasan override pada keputusan yang berbeda dari rekomendasi DSS. | Rendah |
| **BR-09** | **Kegagalan SLA Operasional**<br>Laporan, PE, validasi, atau permintaan fogging tidak diproses sesuai target waktu yang ditetapkan. | Tinggi | Sedang | Tinggi | Operational Owner | Telah Dimitigasi | SLA Cron Job berjalan otomatis. Penetapan SLA per proses (acknowledgment laporan, klaim, PE, validasi, approval fogging, stock request) dengan eskalasi otomatis jika terlampaui. | Persentase pelanggaran SLA lebih dari 10% per minggu; median waktu respons naik signifikan. | Rendah |
| **BR-10** | **Ketergantungan pada Key Person**<br>Workflow bergantung pada satu Kader, operator, admin, atau kepala unit tertentu sehingga berhenti ketika orang tersebut tidak tersedia. | Tinggi | Sedang | Tinggi | Operational Owner | Parsial | Implementasi fitur pendelegasian tugas dan mekanisme reassignment. Monitoring aktor yang tidak aktif dalam periode tertentu. Rotasi peran backup per wilayah. | Satu aktor menangani lebih dari 50% laporan aktif di satu wilayah. | Sedang |
| **BR-11** | **Ketimpangan Akses Digital**<br>Sebagian warga atau petugas mengalami hambatan menggunakan sistem karena keterbatasan perangkat, koneksi internet, atau kemampuan literasi digital. | Sedang | Tinggi | Tinggi | Product Owner / Dinkes | Parsial | Antarmuka web yang ringan dan responsif (PWA); dukungan perangkat seluler kelas menengah; kanal fallback untuk pelaporan (Call Center 112); bantuan operator di Puskesmas. | Jumlah kelurahan dengan tingkat pelaporan digital kurang dari 20% dari populasi. | Sedang |
| **BR-12** | **Hilangnya Kepercayaan Pemangku Kepentingan**<br>Beberapa kesalahan sistem yang terlihat oleh pengguna menyebabkan seluruh sistem dianggap tidak dapat dipercaya. | Sangat Tinggi | Rendah | Tinggi | Product Owner / Dinkes | Parsial | Transparansi rekomendasi sistem; audit trail yang dapat diakses; manajemen insiden yang responsif; publikasi status sistem saat terjadi gangguan; dashboard kualitas data untuk supervisor. | Keluhan publik terkait sistem meningkat; penggunaan aktif turun lebih dari 20% dalam satu bulan. | Rendah |
| **BR-13** | **Pergantian Pejabat dan Perubahan Prioritas Politik**<br>Pergantian Kepala Dinas atau Walikota dapat mengubah prioritas anggaran dan dukungan terhadap proyek ini. | Tinggi | Sedang | Tinggi | Manajemen Dinkes | Belum Dimulai | Dokumentasi dampak dan nilai investasi sistem (laporan ROI). Pelibatan aktif DPRD dan pemangku kepentingan lintas jabatan. Pengembangan sistem yang terlembagakan, bukan bergantung pada individu. | Pergantian Kepala Dinas tanpa proses serah terima yang terstruktur. | Sedang |
| **BR-14** | **Potensi Manipulasi atau Penyembunyian Data Publik**<br>Administrator tingkat Dinkes dapat memodifikasi data KLB atau menyembunyikan informasi dari Portal Publik untuk kepentingan tertentu. | Sangat Tinggi | Rendah | Tinggi | Security Owner / Audit | Parsial | Audit log mencatat setiap perubahan data sensitif oleh administrator. Pemisahan peran yang ketat antara pengubah data dan pengawas. Pertimbangan mekanisme pelaporan anomali ke lembaga pengawas internal. | Perubahan massal data oleh satu akun admin tanpa otorisasi; log audit menunjukkan penghapusan data secara tidak wajar. | Rendah |
| **BR-15** | **Scope Creep dan Perubahan Kebutuhan Tidak Terkendali**<br>Permintaan fitur baru terus bertambah tanpa prioritas dan mekanisme persetujuan perubahan yang jelas. | Sedang | Tinggi | Tinggi | Product Owner / Technical Owner | Belum Dimulai | Setiap permintaan perubahan dikategorikan: Critical, Required, Optional, atau Future Improvement. Wajib melalui proses Change Request dengan impact assessment sebelum dieksekusi. | Jumlah fitur baru yang masuk tanpa melalui proses Change Request lebih dari 3 per sprint. | Rendah |

---

### B. Operational & Epidemiologi

*Kendala terkait proses kerja lapangan, alur respons wabah, kualitas penyelidikan epidemiologi, dan kapasitas sumber daya manusia.*

| ID | Deskripsi Risiko | Dampak | Probabilitas | Risk Owner | Status | Strategi Mitigasi | KRI / Pemicu |
|:---|:---|:---:|:---:|:---|:---:|:---|:---|
| **OP-01** | **Kader Tidak Merespons Laporan**<br>Kader tidak mengklaim atau menindaklanjuti laporan warga dalam batas waktu SLA. | Tinggi | Sedang | Operational Owner | Telah Dimitigasi | SLA Cron Job melepas klaim otomatis setelah 24 jam; Force Unassign tersedia untuk supervisor; pelanggaran tercatat di Audit Log. | Persentase klaim yang terlepas otomatis oleh sistem lebih dari 15% per minggu. |
| **OP-02** | **Kehabisan Stok Insektisida saat Fogging Disetujui**<br>Puskesmas tidak memiliki stok cukup ketika jadwal fogging telah disetujui. | Sangat Tinggi | Sedang | Operational Owner | Telah Dimitigasi | DSS menghasilkan status FOGGING_TERTUNDA dan auto-generate StockRequest ke Dinkes; fitur Alokasi Silang memindahkan stok dari Puskesmas lain secara otomatis. | Jumlah permintaan fogging dengan status FOGGING_TERTUNDA akibat stok lebih dari 3 per minggu. |
| **OP-03** | **Laporan Palsu (Hoaks/Spam)**<br>Warga mengirimkan laporan fiktif yang menghabiskan kapasitas respons petugas. | Sedang | Sedang | Operational Owner | Telah Dimitigasi | Clustering duplikat (<100m, 7 hari); Anti-Spam Clustering (<50m, 1 jam, lebih dari 3 laporan) menolak otomatis; aksi REJECTED_HOAX tersedia untuk Kader; opsi banding tersedia. | Persentase laporan yang ditandai hoaks/duplikat lebih dari 10% dari total laporan masuk. |
| **OP-04** | **Tim Fogging Tidak Melaporkan Penyelesaian**<br>Fogging selesai namun status di sistem masih In Progress karena Tim tidak memperbarui. | Sedang | Sedang | Operational Owner | Telah Dimitigasi | Sistem mewajibkan pengisian usedItems saat update status ke COMPLETED; stok dipotong secara atomik; upload bukti tercatat di Audit Log. | Jumlah jadwal fogging dengan status In Progress lebih dari 48 jam setelah tanggal pelaksanaan. |
| **OP-05** | **Task Hoarding oleh Kader**<br>Kader mengklaim banyak laporan namun tidak menyelesaikannya sehingga menghambat respons. | Tinggi | Sedang | Operational Owner | Telah Dimitigasi | Fitur ReturnReportTask; AutoReleaseExpiredTasks setelah 24 jam; ForceUnassignReport tersedia untuk supervisor; seluruh aksi tercatat. | Satu Kader memiliki lebih dari 5 klaim aktif tanpa progress selama lebih dari 12 jam. |
| **OP-06** | **Waktu Tunggu Pelaporan Terlalu Lama**<br>Warga baru melaporkan kasus demam pada hari ke-5 atau ke-7, padahal masa kritis respons DBD adalah 1-3 hari pertama. | Sangat Tinggi | Tinggi | Product Owner / Dinkes | Belum Dimulai | Kampanye edukasi "Lapor Hari Pertama Demam". Notifikasi pengingat kepada warga yang pernah melapor. Simplifikasi alur pelaporan seminimal mungkin. | Median waktu antara onset demam dan pelaporan ke sistem lebih dari 3 hari. |
| **OP-07** | **Data Tidak Representatif akibat Under-Reporting**<br>Hanya sebagian kecil warga yang melapor via sistem, sehingga heatmap dan DSS tidak menggambarkan kondisi wabah yang sebenarnya. | Sangat Tinggi | Tinggi | Data Owner / Dinkes | Belum Dimulai | Integrasi data kasus dari Puskesmas dan rumah sakit sebagai sumber verifikasi silang. Analisis gap antara kasus terkonfirmasi secara medis dengan laporan masuk sistem. | Jumlah kasus DBD terkonfirmasi medis di suatu wilayah jauh melebihi laporan yang masuk via sistem. |
| **OP-08** | **Konflik Data PE dengan Data Rekam Medis**<br>Data demam dari laporan warga tidak diverifikasi silang dengan data diagnosis resmi RS atau Puskesmas. | Tinggi | Sedang | Data Owner / Dinkes | Belum Dimulai | Pemetaan mekanisme rekonsiliasi data antara DBD Guard dan sistem rekam medis yang ada. Fitur impor data kasus terkonfirmasi dari sumber medis resmi. | Ditemukan kasus DBD terkonfirmasi laboratorium yang tidak tercatat dalam sistem DBD Guard. |
| **OP-09** | **Kesalahan Koordinat GPS pada Laporan**<br>Lokasi laporan tidak akurat akibat sinyal GPS lemah (indoor, area padat) sehingga tim fogging mendatangi lokasi yang keliru. | Tinggi | Sedang | Technical Owner | Parsial | Validasi radius input koordinat; peta interaktif untuk konfirmasi lokasi oleh Kader saat PE; mekanisme koreksi koordinat pasca-klaim. | Jumlah pengaduan dari Tim Fogging terkait lokasi tidak ditemukan atau tidak sesuai. |
| **OP-10** | **Kepala Puskesmas sebagai Bottleneck Approval**<br>Seluruh persetujuan fogging harus melalui Kepala Puskesmas; ketidakhadiran satu orang ini menghentikan seluruh alur. | Tinggi | Sedang | Operational Owner / Puskesmas | Belum Dimulai | Implementasi mekanisme pendelegasian wewenang approval kepada Wakil Kepala atau petugas yang ditunjuk. SLA monitoring untuk proses approval. | Jumlah permintaan fogging yang menunggu approval lebih dari 24 jam. |
| **OP-11** | **Beban Kerja Kader Tidak Seimbang**<br>Satu Kader menanggung area dengan volume laporan jauh lebih tinggi dari Kader lain, berpotensi menyebabkan burnout dan abandon task. | Sedang | Sedang | Operational Owner | Belum Dimulai | Dashboard distribusi beban kerja per Kader untuk supervisor. Mekanisme redistribusi tugas oleh Kepala Puskesmas. Monitoring durasi penyelesaian rata-rata per Kader. | Satu Kader menangani lebih dari 40% total laporan aktif di satu wilayah kerja Puskesmas. |
| **OP-12** | **Tidak Ada Mekanisme Joint-Response Lintas Puskesmas**<br>Penyebaran wabah yang melintasi batas wilayah kerja dua atau lebih Puskesmas tidak memiliki alur koordinasi di dalam sistem. | Tinggi | Sedang | Business Owner / Dinkes | Belum Dimulai | Pengembangan fitur notifikasi lintas-Puskesmas untuk kasus yang berada di radius perbatasan wilayah. Alur eskalasi ke Dinkes sebagai koordinator untuk kasus KLB lintas wilayah. | Teridentifikasi kluster wabah yang melibatkan lebih dari satu wilayah kerja Puskesmas tanpa koordinasi terdokumentasi. |
| **OP-13** | **Turnover Kader yang Tinggi**<br>Kader DBD bersifat sukarela dengan pergantian yang sering; sistem tidak memiliki mekanisme onboarding cepat untuk Kader pengganti. | Sedang | Tinggi | Operational Owner / Puskesmas | Belum Dimulai | Panduan penggunaan singkat berbasis peran. Fitur onboarding mandiri. Pendampingan awal oleh Puskesmas. Monitoring Kader yang belum aktif selama periode tertentu. | Persentase akun Kader yang tidak aktif lebih dari 30 hari melebihi 25%. |
| **OP-14** | **Tim Fogging Tidak Familiar dengan Sistem Digital**<br>Anggota Tim Fogging mungkin bukan pengguna teknologi yang terbiasa, sehingga update status real-time tidak dilakukan. | Sedang | Sedang | Operational Owner | Belum Dimulai | Pelatihan penggunaan sistem khusus Tim Fogging. Antarmuka yang disederhanakan untuk aksi Tim Fogging. Monitoring status tugas yang tidak diperbarui. | Persentase tugas fogging yang statusnya diperbarui lebih dari 6 jam setelah waktu selesai melebihi 20%. |

---

### C. Technology & Infrastructure

*Isu yang melibatkan kinerja peladen, basis data, penyimpanan berkas, dan ketahanan arsitektur.*

| ID | Deskripsi Risiko | Dampak | Probabilitas | Risk Owner | Status | Strategi Mitigasi | KRI / Pemicu |
|:---|:---|:---:|:---:|:---|:---:|:---|:---|
| **TC-01** | **Kegagalan Peladen (Server Downtime)**<br>Peladen mati sehingga seluruh layanan pelaporan lumpuh. | Sangat Tinggi | Sedang | Technical Owner | Parsial | Perencanaan Dockerisasi penuh backend dan frontend; health check endpoint; reverse proxy; monitoring uptime. | Uptime sistem di bawah 99% dalam satu bulan. |
| **TC-02** | **Kegagalan Transmisi Notifikasi**<br>Email atau FCM gagal terkirim sehingga Kader tidak mengetahui adanya laporan baru. | Sedang | Tinggi | Technical Owner | Parsial | Konfigurasi kredensial Firebase dan SMTP di lingkungan produksi; implementasi mekanisme retry queue; pencatatan kegagalan pengiriman. | Persentase notifikasi gagal terkirim lebih dari 5% dalam satu hari. |
| **TC-03** | **Inkonsistensi Data Geospasial**<br>Perhitungan radius atau clustering tidak akurat. | Tinggi | Rendah | Technical Owner | Telah Dimitigasi | PostgreSQL + PostGIS dengan query geospasial native; indeks GIST; komponen GIS heatmap sudah terimplementasi. | Laporan duplikat atau kluster yang lolos dari filter radius sistem. |
| **TC-04** | **Sistem Penyimpanan Lokal Tidak Terskala**<br>Foto bukti disimpan di filesystem lokal, rentan hilang dan tidak dapat diskalakan secara horizontal. | Tinggi | Tinggi | Technical Owner | Belum Dimulai | Integrasi MinIO sebagai Object Storage self-hosted; migrasi fungsi saveFile() ke MinIO SDK; konfigurasi sebagai layanan dalam Docker Compose. | Kapasitas penyimpanan lokal mencapai lebih dari 70%. |
| **TC-05** | **Titik Tunggal Kegagalan pada Basis Data**<br>Hanya satu instance PostgreSQL tanpa replikasi atau backup terjadwal. | Sangat Tinggi | Sedang | Technical Owner | Belum Dimulai | Konfigurasi pg_dump terjadwal ke volume terpisah; untuk produksi: pertimbangkan PostgreSQL streaming replication atau managed database. | Tidak ada rekaman backup berhasil dalam 24 jam terakhir. |
| **TC-06** | **Kredensial Default pada Lingkungan Pengembangan**<br>JWT Secret dan password database masih menggunakan nilai default yang belum dirotasi. | Tinggi | Rendah | Technical Owner / Security Owner | Parsial | Guard log.Fatal jika JWT_SECRET kosong di mode produksi sudah ada; deployment checklist yang mewajibkan rotasi semua secret sebelum produksi. | Kredensial default terdeteksi aktif di lingkungan staging atau produksi. |
| **TC-07** | **Performa Sistem saat Beban Puncak KLB**<br>Belum ada load testing; pada musim puncak DBD (Oktober-Maret) laporan bisa masuk ratusan per hari secara bersamaan. | Sangat Tinggi | Sedang | Technical Owner | Belum Dimulai | Pelaksanaan load testing sebelum peluncuran produksi; optimasi query basis data; implementasi connection pooling; monitoring response time. | Response time API meningkat lebih dari 2x lipat dibanding baseline. |
| **TC-08** | **Tidak Ada Rencana Migrasi Data**<br>Jika arsitektur perlu diupgrade, tidak ada prosedur migrasi data historis yang terdokumentasi. | Sedang | Rendah | Technical Owner | Belum Dimulai | Dokumentasi skema basis data versi saat ini; skrip migrasi untuk setiap perubahan skema; version control untuk migrasi database. | Perubahan skema basis data tanpa skrip migrasi yang terdokumentasi. |

---

### D. Security & Privacy

*Potensi ancaman siber, pelanggaran privasi informasi medis warga, dan pemenuhan hak subjek data.*

| ID | Deskripsi Risiko | Dampak | Probabilitas | Risk Owner | Status | Strategi Mitigasi | KRI / Pemicu |
|:---|:---|:---:|:---:|:---|:---:|:---|:---|
| **SC-01** | **Kebocoran Data Medis dan Pribadi Warga**<br>Akses tidak sah terhadap data identitas, lokasi, atau kondisi kesehatan warga. | Sangat Tinggi | Rendah | Security Owner | Telah Dimitigasi | RBAC ketat dengan 7 peran; fitur anonymous_flag; rate_limiter; Portal Publik hanya menampilkan agregasi per kelurahan; seluruh API non-publik dilindungi JWT. | Terdeteksi akses tidak sah atau permintaan API anomali dalam volume tinggi. |
| **SC-02** | **Eksploitasi Celah Keamanan Aplikasi**<br>Serangan SQL Injection, XSS, atau CSRF melalui form input. | Tinggi | Rendah | Technical Owner / Security Owner | Telah Dimitigasi | GORM menggunakan parameterized query secara default; validasi input di setiap handler; React DOM escaping; JWT Bearer token sebagai perlindungan CSRF. | Terdeteksi pola input mencurigakan atau error query yang tidak biasa di log. |
| **SC-03** | **Penyalahgunaan Wewenang Internal**<br>Akun administrator diretas atau disalahgunakan oleh pihak dalam untuk memodifikasi data sensitif. | Tinggi | Rendah | Security Owner | Telah Dimitigasi | Audit Log menyeluruh mencatat semua aksi sensitif beserta actor_id, action, entity, dan timestamp; dashboard Audit Log tersedia di panel Developer. | Perubahan data sensitif dalam volume tinggi oleh satu akun dalam waktu singkat. |
| **SC-04** | **Foto Bukti Medis Terekspos Tanpa Autentikasi**<br>Foto surat RS atau kondisi rumah warga dapat diakses tanpa verifikasi hak akses. | Tinggi | Sedang | Security Owner / Technical Owner | Parsial | Implementasi Pre-signed URL melalui MinIO dengan validasi peran pengguna sebelum URL diterbitkan (bergantung pada penyelesaian TC-04). | Berkas foto dapat diakses langsung melalui URL statis tanpa token autentikasi. |
| **SC-05** | **Data Lokasi Presisi Berpotensi Menimbulkan Stigmatisasi Sosial**<br>Koordinat rumah warga yang melaporkan kasus DBD, jika bocor, dapat menyebabkan stigmatisasi terhadap pemilik rumah. | Tinggi | Rendah | Security Owner / Data Owner | Parsial | Portal Publik hanya menampilkan data agregat per kelurahan/RW; koordinat presisi tidak pernah ditampilkan ke publik; pembatasan akses koordinat hanya untuk peran tertentu. | Data koordinat individual terekspos di endpoint publik atau laporan yang dapat diunduh. |
| **SC-06** | **Foto yang Diunggah Mengandung Informasi Berlebih**<br>Foto yang diunggah warga (surat RS, foto rumah) dapat mengandung data pribadi yang tidak disengaja seperti foto KTP atau rekam medis lengkap. | Sedang | Sedang | Security Owner / Product Owner | Belum Dimulai | Panduan unggahan foto yang jelas untuk pengguna; pertimbangan pemrosesan foto otomatis untuk menghapus metadata EXIF; validasi tipe dan ukuran file. | Ditemukannya foto yang mengandung data pribadi berlebih melalui audit berkas. |
| **SC-07** | **Tidak Ada Mekanisme Penghapusan Data (Right to Deletion)**<br>Warga yang ingin menghapus riwayat laporannya tidak memiliki mekanisme yang tersedia, sebagaimana diatur dalam UU PDP. | Tinggi | Rendah | Security Owner / Legal | Belum Dimulai | Perancangan fitur permintaan penghapusan atau anonimisasi data oleh warga; prosedur manual sementara melalui administrator; dokumentasi kebijakan retensi data. | Diterimanya permintaan penghapusan data dari warga tanpa mekanisme penanganan yang tersedia. |

---

### E. Compliance & Regulatory

*Deviasi terhadap SOP Kemenkes, regulasi Pemerintah Kota Depok, keselamatan kerja, dan undang-undang yang berlaku.*

| ID | Deskripsi Risiko | Dampak | Probabilitas | Risk Owner | Status | Strategi Mitigasi | KRI / Pemicu |
|:---|:---|:---:|:---:|:---|:---:|:---|:---|
| **CP-01** | **Keputusan DSS Menyimpang dari SOP Kemenkes**<br>Sistem merekomendasikan tindakan yang tidak sesuai dengan pedoman penanganan DBD Kemenkes. | Tinggi | Rendah | Business Owner / Technical Owner | Telah Dimitigasi | DSS mengikuti aturan PRD secara eksak; rekomendasi bukan keputusan final; alur approval Kepala Puskesmas tetap wajib. | Override rekomendasi DSS oleh petugas yang tidak disertai alasan terdokumentasi. |
| **CP-02** | **Master Data Wilayah Tidak Sinkron**<br>Pemetaan Kecamatan, Kelurahan, dan Puskesmas tidak konsisten dengan data resmi Pemkot Depok. | Tinggi | Sedang | Data Owner / Dinkes | Parsial | Dropdown (bukan free-text) untuk input lokasi di frontend; perlunya mekanisme sinkronisasi berkala dengan data kependudukan resmi. | Ditemukannya data wilayah yang tidak konsisten antara sistem dengan data resmi Pemkot. |
| **CP-03** | **Konflik Ambang Batas SLA (24 jam vs. 72 jam)**<br>AutoReleaseExpiredTasks berjalan dengan threshold 24 jam, sementara PRD menetapkan SLA 3x24 jam (72 jam). | Sedang | Rendah | Technical Owner / Business Owner | Perlu Review | Konfirmasi apakah 24 jam berlaku untuk pelepasan klaim Kader atau perlu disesuaikan ke 72 jam sesuai PRD. | Kader mengalami pelepasan klaim otomatis sebelum batas waktu SLA yang ditetapkan PRD. |
| **CP-04** | **Kepatuhan terhadap UU Pelindungan Data Pribadi**<br>Sistem mengumpulkan data medis dan lokasi warga tanpa persetujuan eksplisit yang terdokumentasi secara hukum. | Sangat Tinggi | Rendah | Legal / Security Owner | Belum Dimulai | Penyusunan Terms of Service dan Privacy Policy yang mengikat secara hukum. Mekanisme persetujuan eksplisit (informed consent) saat pendaftaran warga. | Tidak adanya dokumen ToS dan Privacy Policy yang dapat diakses pengguna sebelum peluncuran produksi. |
| **CP-05** | **Risiko K3: Tidak Ada Konfirmasi Clearance Area Sebelum Fogging**<br>Sistem mencatat jadwal dan lokasi fogging, tetapi tidak memastikan area telah dikosongkan dari manusia sebelum insektisida disemprotkan. | Tinggi | Sedang | Operational Owner / Puskesmas | Belum Dimulai | Penambahan checklist konfirmasi keamanan sebelum Tim Fogging memulai operasi. Notifikasi kepada warga sekitar minimal 2 jam sebelum fogging dilaksanakan. | Laporan kejadian warga terpapar insektisida saat operasi fogging berlangsung. |
| **CP-06** | **Tidak Ada Mekanisme Ekspor Laporan Periodik**<br>Dinkes perlu menyusun laporan bulanan/triwulanan untuk DPRD dan Kemenkes, namun sistem belum memiliki fitur ekspor ini. | Sedang | Sedang | Product Owner / Dinkes | Belum Dimulai | Pengembangan fitur ekspor laporan terstruktur dalam format standar (PDF, Excel) dengan parameter periode, wilayah, dan jenis laporan yang dapat dikonfigurasi. | Permintaan laporan periodik dari Dinkes yang harus diekstrak secara manual dari basis data. |
| **CP-07** | **Perubahan Regulasi Kemenkes di Masa Mendatang**<br>Penyesuaian SOP penanganan wabah DBD dapat membuat logika DSS yang ada menjadi tidak relevan. | Tinggi | Rendah | Business Owner / Technical Owner | Telah Dimitigasi | Arsitektur modular memungkinkan perubahan aturan DSS tanpa mengubah keseluruhan sistem; konfigurasi parameter DSS tersimpan terpisah dari logika inti. | Dikeluarkannya Permenkes atau SE baru yang mengubah kriteria penanganan DBD. |
| **CP-08** | **Tidak Ada Audit Trail untuk Disposal Insektisida**<br>Sistem mencatat pemotongan stok namun tidak mencatat pembuangan sisa limbah insektisida pascafogging. | Sedang | Sedang | Operational Owner / Dinkes | Belum Dimulai | Penambahan field laporan disposal insektisida sisa; integrasi dengan prosedur K3 pengelolaan limbah B3 sesuai regulasi yang berlaku. | Tidak ada dokumentasi disposal insektisida dalam periode operasional fogging. |

---

### F. Development & Delivery

*Tantangan dalam tahapan pengujian, penyebaran, pipeline, dan kualitas kode.*

| ID | Deskripsi Risiko | Dampak | Probabilitas | Risk Owner | Status | Strategi Mitigasi | KRI / Pemicu |
|:---|:---|:---:|:---:|:---|:---:|:---|:---|
| **DV-01** | **Cakupan Pengujian Otomatis Tidak Memadai**<br>Hanya 2 file test yang ada; banyak use case kritis (report, inventory, verification) belum memiliki unit test. | Tinggi | Tinggi | Technical Owner | Parsial | Prioritas unit test untuk: logika rekomendasi DSS, logika alokasi silang, SLA auto-release, dan validasi anti-overbooking. Target code coverage minimal 70% untuk komponen inti. | Code coverage di bawah 50% untuk komponen use case inti. |
| **DV-02** | **Tidak Ada CI/CD Pipeline dan Staging Environment**<br>Build dilakukan secara manual; tidak ada Dockerfile untuk backend/frontend; tidak ada lingkungan staging. | Tinggi | Tinggi | Technical Owner | Belum Dimulai | Pembuatan Dockerfile untuk backend dan frontend; pelengkapan docker-compose dengan semua layanan; konfigurasi GitHub Actions untuk CI minimal (lint, test, build). | Deployment ke produksi dilakukan tanpa melalui proses pengujian otomatis. |
| **DV-03** | **Kesenjangan Lingkungan Development vs. Production**<br>Backend berjalan sebagai file .exe; notifikasi mock; storage lokal; kredensial default. Gap sangat besar dengan kebutuhan produksi. | Sedang | Tinggi | Technical Owner | Belum Dimulai | Penyusunan deployment checklist dan environment variable matrix untuk tiga lingkungan (dev, staging, prod). | Perbedaan perilaku sistem yang signifikan antara lingkungan development dan staging. |
| **DV-04** | **File Executable Tersimpan di Repositori**<br>File api.exe, server.exe, dan binary lainnya (total sekitar 174 MB) tersimpan di direktori backend. | Rendah | Tinggi | Technical Owner | Belum Dimulai | Penambahan *.exe ke .gitignore; pembersihan binary dari riwayat Git menggunakan git filter-repo. | Ukuran repositori terus membesar setiap kali ada build baru. |

---

### G. Integration & External Dependency

*Risiko terkait integrasi dengan ekosistem nasional dan ketergantungan pada layanan eksternal.*

| ID | Deskripsi Risiko | Dampak | Probabilitas | Risk Owner | Status | Strategi Mitigasi | KRI / Pemicu |
|:---|:---|:---:|:---:|:---|:---:|:---|:---|
| **IE-01** | **Isolasi dari Ekosistem SatuSehat Kemenkes**<br>DBD Guard berdiri sendiri tanpa integrasi dengan platform pelaporan kesehatan nasional SatuSehat, berpotensi menjadi sistem redundan. | Tinggi | Sedang | Business Owner / Technical Owner | Belum Dimulai | Pemetaan kebutuhan teknis untuk pengembangan API yang kompatibel dengan standar SatuSehat. Penyiapan mekanisme ekspor/impor data sementara. Penetapan integration roadmap. | Dikeluarkannya kebijakan Kemenkes yang mewajibkan seluruh pelaporan DBD melalui SatuSehat. |
| **IE-02** | **Ketergantungan pada Penyedia Notifikasi Eksternal**<br>Kegagalan layanan Firebase Cloud Messaging atau penyedia SMTP akan menghentikan seluruh alur notifikasi sistem. | Sedang | Sedang | Technical Owner | Parsial | Implementasi abstraction layer untuk layanan notifikasi; mekanisme retry queue; pencatatan log kegagalan; pertimbangan backup provider untuk kondisi darurat. | Kegagalan pengiriman notifikasi dalam volume tinggi yang berlangsung lebih dari 1 jam. |
| **IE-03** | **Tidak Ada Integrasi dengan Data Rekam Medis**<br>Data kasus DBD dari sistem tidak diverifikasi silang dengan data diagnosis dari SIM RS atau SIM Puskesmas. | Tinggi | Sedang | Data Owner / Dinkes | Belum Dimulai | Pemetaan skema integrasi dengan SIM RS/Puskesmas yang ada di Depok. Mekanisme impor data kasus terkonfirmasi laboratorium sebagai sumber data sekunder. | Perbedaan signifikan antara jumlah kasus DBD terkonfirmasi medis dengan data dalam sistem. |
| **IE-04** | **Sistem Belum Dirancang untuk Ekspansi Multi-Wilayah**<br>Jika kota lain ingin mengadopsi DBD Guard, sistem saat ini tidak memiliki arsitektur multi-tenancy. | Rendah | Rendah | Technical Owner | Belum Dimulai | Dokumentasi arsitektur saat ini dan identifikasi komponen yang perlu direfaktor untuk mendukung multi-tenancy di masa depan. | Diterimanya permintaan adopsi dari kota/kabupaten lain tanpa rencana arsitektur yang siap. |

---

## 5. Matriks Ringkasan Risiko (Heat Map)

| Dampak \ Probabilitas | Probabilitas Rendah | Probabilitas Sedang | Probabilitas Tinggi |
|:---|:---|:---|:---|
| **Sangat Tinggi** | SC-01 (Dimitigasi), CP-04 (Belum) | TC-01 (Parsial), TC-05 (Belum), TC-07 (Belum), BR-04 (Belum) | BR-01 (Parsial), OP-06 (Belum), OP-07 (Belum), BR-06 (Parsial) |
| **Tinggi** | CP-01 (Dimitigasi), CP-07 (Dimitigasi), SC-05 (Parsial), SC-07 (Belum), BR-03 (Belum) | BR-02 (Belum), BR-07 (Parsial), BR-08 (Parsial), BR-09 (Dimitigasi), BR-10 (Parsial), BR-13 (Belum), BR-14 (Parsial), OP-08 (Belum), OP-10 (Belum), OP-12 (Belum), CP-05 (Belum), IE-01 (Belum), IE-03 (Belum) | BR-05 (Belum), BR-11 (Parsial), BR-15 (Belum), TC-04 (Belum), DV-01 (Parsial), DV-02 (Belum), OP-13 (Belum) |
| **Sedang** | SC-02 (Dimitigasi), CP-03 (Review), TC-08 (Belum) | TC-02 (Parsial), OP-04 (Dimitigasi), OP-08 (Belum), OP-09 (Parsial), OP-11 (Belum), SC-04 (Parsial), SC-06 (Belum), CP-02 (Parsial), CP-06 (Belum), CP-08 (Belum), IE-02 (Parsial) | DV-03 (Belum), OP-14 (Belum), BR-12 (Parsial) |
| **Rendah** | SC-02 (Dimitigasi), TC-03 (Dimitigasi), TC-06 (Parsial), IE-04 (Belum) | — | DV-04 (Belum) |

---

## 6. Prioritas Tindakan (Action Items)

| Urgensi | ID | Uraian Tugas | Estimasi | Risk Owner |
|:---:|:---|:---|:---|:---|
| **Kritis (P0)** | TC-05 | Setup pg_dump backup terjadwal ke volume terpisah. | 1 Hari | Technical Owner |
| **Kritis (P0)** | TC-04 | Integrasi MinIO sebagai Object Storage; migrasi saveFile(). | 2-3 Hari | Technical Owner |
| **Kritis (P0)** | TC-07 | Pelaksanaan load testing sebelum peluncuran produksi. | 2 Hari | Technical Owner |
| **Kritis (P0)** | BR-02 | Penetapan struktur Product Owner, Data Owner, dan Operational Owner. | 1 Hari | Manajemen Dinkes |
| **Kritis (P0)** | CP-04 | Penyusunan Terms of Service dan Privacy Policy sesuai UU PDP. | 2 Hari | Legal / Security Owner |
| **Kritis (P0)** | BR-04 | Penyusunan estimasi OPEX tahunan untuk penganggaran APBD. | 2 Hari | Manajemen Dinkes |
| **Tinggi (P1)** | DV-02 | Setup Dockerfile, Staging Environment, dan pipeline CI/CD. | 2-3 Hari | Technical Owner |
| **Tinggi (P1)** | TC-02 | Konfigurasi Firebase FCM dan SMTP di lingkungan produksi. | 1-2 Hari | Technical Owner |
| **Tinggi (P1)** | SC-04 | Implementasi Pre-signed URL MinIO untuk akses foto. | 1 Hari | Technical Owner |
| **Tinggi (P1)** | BR-03 | Penyusunan dokumentasi arsitektur, API, deployment, dan runbook. | 3 Hari | Technical Owner |
| **Tinggi (P1)** | BR-05 | Penyusunan draf SK Kepala Dinas sebagai payung hukum wajib pakai. | 5 Hari | Manajemen Dinkes |
| **Tinggi (P1)** | CP-05 | Penambahan checklist keamanan K3 dalam alur operasional fogging. | 1 Hari | Operational Owner |
| **Tinggi (P1)** | OP-10 | Implementasi mekanisme pendelegasian approval Kepala Puskesmas. | 2 Hari | Technical Owner |
| **Menengah (P2)** | BR-01 | Konsolidasi strategi sosialisasi sistem bersama Dinkes dan Puskesmas. | Menyesuaikan Jadwal | Product Owner |
| **Menengah (P2)** | DV-01 | Penulisan unit test untuk DSS, alokasi silang, dan SLA. | 3-4 Hari | Technical Owner |
| **Menengah (P2)** | CP-06 | Pengembangan fitur ekspor laporan periodik (PDF/Excel). | 2-3 Hari | Technical Owner |
| **Menengah (P2)** | IE-01 | Pemetaan kebutuhan integrasi API dengan SatuSehat Kemenkes. | 2 Hari | Business Owner |
| **Menengah (P2)** | OP-06 | Simplifikasi alur pelaporan warga dan kampanye "Lapor Hari Pertama". | Menyesuaikan Jadwal | Product Owner |
| **Menengah (P2)** | CP-03 | Review threshold SLA (24 jam vs. 72 jam) bersama stakeholder. | 0.5 Hari | Business Owner |

---

## 7. Rencana Kontingensi dan Business Continuity

### 7.1 Prosedur Darurat Teknis

Apabila terjadi kondisi kahar atau kegagalan sistem total:

1. **Pemberitahuan Massal:** Admin Utama mengumumkan status gangguan kepada seluruh Puskesmas melalui saluran komunikasi sekunder (WhatsApp Group).
2. **Pengalihan Operasional (Fallback Manual):** Pelaporan warga dialihkan ke Call Center 112. Tim lapangan beralih ke formulir Penyelidikan Epidemiologi berbasis kertas.
3. **Pemulihan Basis Data:** Tim infrastruktur melaksanakan restore dari berkas backup pg_dump terakhir.
4. **Rollback Sistem:** Apabila rilis baru menyebabkan kegagalan operasional, tim teknis melakukan pembalikan ke versi stabil terakhir.
5. **Rekonsiliasi Data:** Admin Puskesmas melaksanakan batch input data manual secara retrospektif setelah sistem pulih.

### 7.2 Business Continuity Decision

Harus ditetapkan secara eksplisit:

| Elemen | Ketentuan |
|:---|:---|
| **Decision Maker** | Siapa yang berwenang menyatakan sistem dalam status UNAVAILABLE dan mengaktifkan mode manual. |
| **Kriteria Pengaktifan** | Sistem tidak dapat diakses lebih dari 2 jam berturut-turut. |
| **Kanal Pemberitahuan** | Saluran komunikasi resmi yang digunakan saat sistem tidak tersedia. |
| **Batas Waktu Input Ulang** | Maksimal 24 jam setelah sistem pulih untuk seluruh data manual. |
| **Proses Verifikasi** | Prosedur verifikasi keakuratan data setelah batch input manual selesai. |

---

## 8. Pemilik Risiko (Risk Ownership Matrix)

| Peran | Tanggung Jawab Risiko Utama |
|:---|:---|
| **Product Owner / Dinkes** | BR-01, BR-11, BR-12, BR-13, BR-15, OP-06, OP-07, CP-06 |
| **Business Owner** | BR-07, BR-08, CP-01, CP-07, IE-01 |
| **Operational Owner** | BR-09, BR-10, CP-05, CP-08, OP-01 s.d. OP-05, OP-10, OP-11, OP-12, OP-13, OP-14 |
| **Technical Owner** | TC-01 s.d. TC-08, DV-01 s.d. DV-04, IE-02, IE-04, OP-09, SC-02 |
| **Data Owner** | BR-06, OP-07, OP-08, CP-02, IE-03 |
| **Security Owner** | SC-01, SC-03 s.d. SC-07, CP-04, BR-14 |
| **Legal / Compliance** | CP-04, SC-07 |
| **Manajemen Dinkes** | BR-02, BR-03, BR-04, BR-05, BR-13 |

---

*Dokumen ini bersifat dinamis (Living Document) dan wajib ditinjau ulang pada setiap Sprint Review, serta minimal setiap 6 bulan sekali bersama Dinas Kesehatan Kota Depok. Setiap pembaruan dokumen harus disertai tanggal revisi dan ringkasan perubahan.*
