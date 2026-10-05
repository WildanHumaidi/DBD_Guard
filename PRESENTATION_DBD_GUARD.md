# DBD Guard: Sistem Pelaporan & Penanganan DBD Terintegrasi Kota Depok

Status Kesiapan: Conditionally Pilot Ready (Uji Coba Terbatas) & Demo Ready, Menunggu Simulasi Kering | 1 Oktober 2026

---

## 0. Ringkasan 1 Halaman

Intinya: DBD Guard adalah platform terpadu untuk percepatan respon penanganan DBD Kota Depok dan pengendalian logistik insektisida yang siap diuji coba terbatas di dua puskesmas.

Kota Depok menghadapi tantangan keterlambatan respon laporan warga, risiko selisih stok insektisida, dan pencatatan manual yang belum terintegrasi. DBD Guard hadir sebagai solusi satu pintu berbasis geospasial yang menghubungkan Warga, Kader RW, Puskesmas, dan Dinas Kesehatan tanpa biaya lisensi perangkat lunak.

Saat ini sistem berstatus Conditionally Pilot Ready untuk uji coba terbatas di Puskesmas Beji dan Puskesmas Pancoran Mas selama 30 hari. Seluruh fungsi alur utama dan proteksi keamanan telah disiapkan, dengan status demo siap dipaparkan menunggu simulasi kering (dry-run).

Dukungan yang dimohonkan dari pimpinan daerah hari ini meliputi penerbitan Surat Perintah Tugas (SPT) pelaksanaan pilot dan arahan kebijakan tata kelola sistem.

---

## 1. Masalah yang Diselesaikan

Intinya: DBD Guard mengubah proses manual yang terpisah-pisah menjadi alur kerja digital yang terpadu, transparan, dan dapat dipantau setiap saat.

Sebelum adanya sistem ini, koordinasi penanganan demam berdarah menghadapi berbagai kendala administratif di lapangan. Tabel berikut merangkum perbedaan operasional sebelum dan sesudah implementasi sistem:

| Aspek Operasional | Kondisi Saat Ini (Sebelum Sistem) | Dengan DBD Guard (Setelah Sistem) |
|:---|:---|:---|
| **Saluran Pelaporan** | Tersebar di berbagai kanal obrolan dan pesan pribadi tanpa nomor tanda terima resmi. | Satu portal publik terpusat dengan nomor resi pelacakan langsung untuk warga. |
| **Penyelidikan Lapangan** | Formulir kertas manual rentan tercecer dan pencatatan jentik dilakukan manual. | Formulir digital terstandar dengan dokumentasi foto serta pencatatan temuan jentik dan warga demam. |
| **Keputusan Tindakan** | Diskusi informal tanpa dokumentasi tertulis mengenai alasan tindakan. | Rekomendasi otomatis sistem dengan rekam diskresi klinis dokter. |
| **Pengelolaan Obat** | Pencatatan buku gudang manual dengan risiko selisih stok fisik dan digital. | Mutasi transfer stok dua pihak dengan pemotongan otomatis berbasis foto bukti kerja. |
| **Akuntabilitas & Audit** | Sulit membuktikan riwayat perubahan jika terjadi sengketa atau audit pemeriksa. | Jejak audit permanen mencatat waktu, pelaku, dan rincian perubahan; tidak dapat diubah atau dihapus oleh pengguna biasa. |

> **Untuk Presenter:**
> * Tunjukkan: Tabel perbandingan di atas kepada audiens.
> * Sampaikan: "Sistem ini bukan sekadar aplikasi lapor, melainkan sistem kendali operasional dan perlindungan aset daerah."
> * Waktu: 2 menit.

---

## 2. Cara Kerja Sistem

Intinya: Sistem menghubungkan tujuh peran pengguna dalam satu alur kerja tertutup dari pelaporan awal hingga pemotongan stok obat.

Alur penanganan bergerak secara berurutan sesuai kewenangan tugas masing-masing pihak:

```text
[Warga Melapor] -> [Kader Cek Lokasi] -> [Petugas Survei PE]
                                                |
[Stok Terpotong] <- [Tim Fogging Semprot] <- [Kepala Faskes Setujui]
        |
[Dashboard Peta & EWS Dinkes Terupdate]
```

Tujuh peran pengguna memiliki tanggung jawab yang terbagi secara tegas sesuai spesifikasi kebutuhan produk:

| Peran Pengguna | Tanggung Jawab Utama |
|:---|:---|
| **Warga (Pelapor)** | Mengirimkan laporan dugaan sarang nyamuk dengan foto bukti dan titik koordinat rumah, serta memantau status laporan. |
| **Kader Wilayah (Jumantik)** | Memverifikasi keabsahan laporan warga di wilayahnya. |
| **Petugas Puskesmas (Sanitarian)** | Melaksanakan Penyelidikan Epidemiologi (PE) dan mencatat parameter klinis lingkungan di lapangan. |
| **Kepala Puskesmas** | Meninjau rekomendasi, menyetujui jadwal fogging, dan mengajukan permintaan stok ke Dinkes. |
| **Tim Fogging** | Melaksanakan penyemprotan di blok sasaran dan mengunggah foto bukti kerja serta volume obat terpakai. |
| **Dinas Kesehatan (Dinkes)** | Memantau peta sebaran kasus seluruh kota, mengelola inventaris pusat insektisida, dan memonitor peringatan dini. |
| **Admin Developer** | Mengelola akun pengguna, memelihara keandalan teknis sistem, dan mengawasi jejak audit sistem. |

> **Untuk Presenter:**
> * Tunjukkan: Bagan alur penanganan dan daftar peran.
> * Sampaikan: "Setiap peran memiliki batas kewenangan yang terkunci rapi; tidak ada peran yang bisa melompati prosedur."
> * Waktu: 3 menit.

---

## 3. Yang Akan Kita Tunjukkan Hari Ini (Alur Demo)

Intinya: Demonstrasi hari ini memperlihatkan alur penanganan lengkap dari laporan warga sampai pembaruan peta, serta pembuktian kontrol keamanan sistem.

Total alokasi waktu presentasi dan demonstrasi adalah 35 menit dengan rincian segmen sebagai berikut:

| Waktu | Skenario yang Ditunjukkan | Peran Terlibat | Pesan Kunci untuk Audiens | Moda Demo |
|:---:|:---|:---|:---|:---:|
| **Menit 01-05** | Pembukaan & Latar Belakang Masalah | Presenter | Urgensi tata kelola DBD Kota Depok yang akuntabel. | SLIDE |
| **Menit 06-11** | Pelaporan Warga & Verifikasi Kader RW | Warga, Kader | Laporan warga diterima kader kelurahan yang berwenang. | SLIDE / tangkapan layar |
| **Menit 12-17** | Penyelidikan Faskes & Rekomendasi Sistem | Sanitarian | Sistem mengevaluasi kriteria tindakan secara objektif. | SLIDE / tangkapan layar |
| **Menit 18-22** | Persetujuan Tindakan & Eksekusi Fogging | Kepala Faskes, Fogger | Keputusan di tangan dokter; foto bukti wajib diunggah. | SLIDE / tangkapan layar |
| **Menit 23-25** | Pemotongan Stok Obat & Peta Pantau | Petugas, Dinkes | Logistik terpotong akurat dan peta wilayah terbarui. | SLIDE / tangkapan layar |
| **Menit 26-28** | Uji Ketahanan & Penolakan Pelanggaran | Sistem Keamanan | Bukti sistem menolak klaim tugas di luar wilayah kerja. | SLIDE / tangkapan layar |
| **Menit 29-35** | Diskusi, Skenario Cadangan & Tanya Jawab | Seluruh Stakeholder | Pembahasan langkah pilot dan peragaan skenario cadangan. | SLIDE |

Skenario cadangan mencakup pengambilalihan penugasan darurat (takeover) dan diskresi klinis dokter (DSS override) yang disiapkan sebagai peragaan pada sesi tanya jawab. Waktu yang tertera pada blok 'Untuk Presenter' adalah panduan bicara di dalam sesi, bukan waktu tambahan.

> **Untuk Presenter:**
> * Tunjukkan: Jadwal alur demonstrasi agar audiens memahami batasan waktu.
> * Sampaikan: "Peragaan menggunakan slide dan tangkapan layar karena simulasi kering tim masih dijadwalkan sebelum uji coba langsung."
> * Waktu: 1 menit.

---

## 4. Kontrol & Keamanan

Intinya: DBD Guard dibangun dengan prinsip pencegahan dini sehingga penyimpangan wewenang dan pelanggaran wilayah kerja otomatis ditolak oleh sistem.

Sistem tidak hanya diuji pada kondisi berjalan lancar, melainkan telah melalui pengujian ketat terhadap skenario penyimpangan:

1. **Penolakan Klaim Lintas Wilayah:**
   Kader kelurahan tertentu tidak dapat mengklaim laporan dari kelurahan lain. Ketika dicoba, sistem menampilkan penolakan akses demi menjaga ketertiban wilayah kerja masing-masing.

2. **Pencegahan Verifikasi Diri Sendiri (Prinsip Empat Mata):**
   Petugas yang turun melakukan Penyelidikan Epidemiologi (PE) dilarang keras memvalidasi hasil PE-nya sendiri untuk mencegah manipulasi data lapangan.

3. **Perlindungan Data Pelapor dan Berkas Lapangan:**
   Foto kondisi rumah pelapor dilindungi jalur khusus dan hanya dapat dibuka oleh petugas berwenang yang bertugas. Publik umum hanya dapat melihat nomor tiket tanpa identitas nama maupun nomor telepon pelapor.

> **Untuk Presenter:**
> * Tunjukkan: Layar pesan penolakan akses saat pengguna mencoba melanggar batas wilayah.
> * Sampaikan: "Keamanan data warga dan dokumentasi lapangan terlindungi melalui pembatasan akses berbasis peran."
> * Waktu: 3 menit.

---

## 5. Keputusan Manusia vs Rekomendasi Sistem (DSS)

Intinya: Sistem pendukung keputusan memberikan saran objektif berbasis pedoman teknis, namun wewenang akhir tindakan sepenuhnya berada pada diskresi klinis Kepala Puskesmas.

Algoritma sistem mengevaluasi indikator lapangan dan ketersediaan logistik untuk menghasilkan rekomendasi terstandar:

| Kondisi Lapangan | Rekomendasi Otomatis Sistem | Tindakan Lanjutan |
|:---|:---|:---|
| Jentik = Ya ATAU demam >= 3, stok tersedia | **Fogging Fokus** | Kepala Puskesmas menerbitkan persetujuan jadwal penyemprotan. |
| Jentik = Ya ATAU demam >= 3, stok kurang | **Tertunda (Draf Permintaan Stok)** | Sistem membuat draf permintaan stok; alur persetujuannya menunggu keputusan DEC-01. |
| Jentik = Tidak DAN demam = 0 | **Edukasi & PSN** | Edukasi warga dan PSN (3M Plus). |
| Kombinasi lain: mengikuti aturan bisnis lengkap | *Merujuk pada Dokumen Aturan Bisnis Lengkap.* | - |

Apabila dokter atau Kepala Puskesmas menilai kondisi lapangan membutuhkan tindakan berbeda dari saran sistem, pimpinan faskes berhak mengambil keputusan pengalihan (override). Seluruh pengalihan wajib menyertakan alasan medis tertulis dan tercatat permanen pada sistem.

> **Untuk Presenter:**
> * Tunjukkan: Rekomendasi pada layar faskes dan formulir pengisian alasan pengalihan medis.
> * Sampaikan: "Aplikasi ini memandu objektivitas, tetapi tidak mengambil alih kedaulatan dokter pimpinan faskes."
> * Waktu: 3 menit.

---

## 6. Integritas Stok

Intinya: Setiap liter insektisida dipantau dengan pencatatan ganda sehingga perpindahan barang antar-puskesmas tidak dapat disalahgunakan.

Sistem menerapkan prinsip bahwa total stok pengirim, stok penerima, dan barang dalam perjalanan selalu seimbang:

```text
Saat Transit:
[Dinkes: 80L] -> [Transit: 20L] -> [Puskesmas: 10L] (Total 110L)

Setelah Diterima:
[Dinkes: 80L] -> [Transit: 0L] -> [Puskesmas: 30L] (Total 110L)
```

Contoh simulasi mutasi obat tercatat sebagai berikut:

| Tahap Mutasi | Stok Gudang Asal | Status Pengiriman | Stok Faskes Penerima | Keterangan Akuntansi |
|:---|:---:|:---:|:---:|:---|
| **1. Kondisi Awal** | 100 Liter | Belum Dikirim | 10 Liter | Saldo tercatat sebelum transaksi pengiriman. |
| **2. Pengiriman Barang** | 80 Liter | 20 Liter Dalam Perjalanan | 10 Liter | Stok asal berkurang, masuk status dalam transit. |
| **3. Konfirmasi Penerimaan** | 80 Liter | Selesai Diterima | 30 Liter | Fisik diverifikasi, stok penerima bertambah sah. |
| **4. Eksekusi Fogging** | 80 Liter | - | 25 Liter | Pemotongan 5 Liter wajib menyertakan foto bukti. |

> **Untuk Presenter:**
> * Tunjukkan: Menu inventaris dan histori transaksi transfer stok.
> * Sampaikan: "Pencatatan ini menjaga faskes dari temuan selisih barang oleh pemeriksa keuangan daerah."
> * Waktu: 3 menit.

---

## 7. Peringatan Dini (EWS/KLB)

Intinya: Sistem secara otomatis memantau ambang batas kasus per wilayah kelurahan untuk memberikan sinyal siaga dini kepada dinas kesehatan.

Indikator sistem memberikan penanda visual otomatis saat terjadi konsentrasi kasus:

* **Sinyal Waspada Kelurahan:** Timbul ketika akumulasi laporan positif di suatu kelurahan melonjak melampaui ambang batas yang ditetapkan Dinkes.
* **Pewarnaan Peta Spasial:** Wilayah kelurahan yang mengalami lonjakan kasus berubah warna menjadi zona perhatian khusus pada peta interaktif.
* **Bukan Ramalan Statistik:** Sistem murni menampilkan fakta kenaikan kasus riil yang dilaporkan, tanpa membuat klaim prediksi buatan yang belum terbukti.

> **Untuk Presenter:**
> * Tunjukkan: Tampilan peta sebaran wilayah Kota Depok dengan kelurahan berstatus siaga.
> * Sampaikan: "Dinas Kesehatan dapat langsung memprioritaskan alokasi bantuan ke kelurahan yang menyala merah."
> * Waktu: 2 menit.

---

## 8. Jejak Audit

Intinya: Seluruh tindakan operasional tercatat dalam buku catatan digital yang tidak dapat diubah maupun dihapus oleh pengguna biasa.

Setiap mutasi status dan persetujuan penanganan terdokumentasi lengkap dalam log sistem (*contoh ilustrasi*):

| Waktu Kejadian | Nama Pengguna | Peran | Tindakan | Rincian Catatan |
|:---:|:---|:---|:---|:---|
| 08:30 WIB | Kader Kelurahan Beji | Kader Wilayah (Jumantik) | Menolak Laporan | Laporan dinilai hoaks oleh kader. |
| 09:00 WIB | Sanitarian Puskesmas Beji | Petugas Puskesmas (Sanitarian) | Takeover Penugasan | Puskesmas mengambil alih laporan hasil banding yang dieskalasi. |
| 09:30 WIB | Sanitarian Puskesmas Beji | Petugas Puskesmas (Sanitarian) | Input PE | Mencatat temuan lapangan: jentik positif, demam 3 orang. |
| 10:00 WIB | Kepala Puskesmas Beji | Kepala Puskesmas | Diskresi Klinis (DSS Override) | Rekomendasi fogging dibatalkan dengan alasan medis. |
| 10:30 WIB | Petugas Dinkes | Dinas Kesehatan (Dinkes) | Transfer Stok | Dinkes menyetujui transfer stok ke Puskesmas. |

Catatan ini menjadi jejak audit kepatuhan batas waktu respon penanganan laporan dan transparansi penggunaan bahan kimia daerah.

> **Untuk Presenter:**
> * Tunjukkan: Tabel log audit pada dashboard administrator.
> * Sampaikan: "Bila sewaktu-waktu ada aduan warga atau pemeriksaan inspektorat, kronologi penanganan terbaca detik demi detik."
> * Waktu: 2 menit.

---

## 9. Perjalanan Kualitas

Intinya: DBD Guard telah melewati rangkaian pengujian fungsional dan verifikasi risiko sebelum dinyatakan siap untuk demonstrasi eksekutif.

Pengembangan sistem melalui tahapan evaluasi kualitas berjenjang:

```text
[Audit Celah] -> [Manajemen Risiko] -> [Perbaikan Masalah Kritis]
                                                |
[Kesiapan Pilot] <- [Uji Regresi] <- [Pengetatan Batas Wilayah]
```

Rincian pencapaian tiap fase penjaminan mutu:

* **Audit Celah Awal:** Mengidentifikasi potensi titik lemah penanganan laporan dan tata kelola akun.
* **Manajemen Risiko:** Menyusun mitigasi terhadap risiko server padam, koneksi terputus, dan keterbatasan anggaran.
* **Perbaikan Masalah Kritis:** Mengunci transaksi logistik agar saldo tidak bisa minus dan menutup akses berkas foto pelapor.
* **Uji Simulasi Pemulihan Bencana:** Uji pemulihan data (restore drill) dilaksanakan pada tahap perbaikan dan datanya cocok; salinan cadangan masih di server yang sama. Drill ulang di lingkungan pilot dijadwalkan (TBD).
* **Uji Validasi Operasional:** Memastikan alur pelaporan warga hingga persetujuan faskes berjalan konsisten.
* **Pengetatan Batas Wilayah:** Memastikan penolakan akses otomatis saat ada upaya klaim lintas batas kelurahan.
* **Uji Regresi:** Memastikan perbaikan celah keamanan tidak merusak fungsi fitur yang telah berjalan stabil.
* **Penilaian Kesiapan Akhir:** Menetapkan postur kesiapan bersyarat untuk pelaksanaan uji coba terbatas di dua faskes.

> **Untuk Presenter:**
> * Tunjukkan: Bagan perjalanan pengujian.
> * Sampaikan: "Aplikasi ini telah diuji secara teknis berkali-kali untuk memastikan stabilitasnya di hadapan bapak dan ibu."
> * Waktu: 2 menit.

---

## 10. Rencana Pilot

Intinya: Uji coba terbatas dilaksanakan di dua puskesmas selama 30 hari kalender untuk memvalidasi kesiapan operasional lapangan sebelum perluasan kota.

Rencana pelaksanaan pilot dirancang terkendali:

* **Wilayah Uji Coba:** Puskesmas Beji dan Puskesmas Pancoran Mas (mencakup Kelurahan Beji dan Kelurahan Pancoran Mas). Faskes lain belum diaktifkan selama masa pilot.
* **Durasi Pelaksanaan:** 30 hari kalender sejak kick-off (tanggal ditetapkan setelah SPT terbit).
* **Personil Lapangan:** Sanitarian faskes, dokter kepala faskes, perwakilan kader RW setempat, tim logistik faskes, dan tim fogging.
* **Indikator Evaluasi:**
  1. Keutuhan data pelaporan masyarakat yang masuk.
  2. Kepatuhan waktu respon penanganan laporan di bawah 72 jam (target persentase ditetapkan bersama Dinkes).
  3. Ketepatan otorisasi wilayah kerja petugas lapangan.
  4. Kesesuaian fisik stok gudang dengan saldo tercatat digital.
  5. Kelengkapan foto bukti pelaksanaan pada setiap penugasan.
  6. Kestabilan server dan tingkat kenyamanan kader dalam menggunakan aplikasi.
* **Prosedur Cadangan Operasional:** Jika terjadi pemadaman listrik atau server terganggu berkepanjangan (ambang batas durasi belum ditetapkan), petugas mencatat hasil survei pada formulir kertas standar Kementerian Kesehatan rangkap tiga, lalu diinput susulan setelah layanan pulih.
* **Usulan tim teknis:** bila stok faskes < 20 L atau terjadi KLB, pilot dijeda dan dieskalasi ke Dinkes; pengambil keputusan: TBD - Business Owner.

> **Untuk Presenter:**
> * Tunjukkan: Ruang lingkup pilot dan mekanisme pelaporan berkala.
> * Sampaikan: "Uji coba ini bukan rilis terburu-buru, melainkan langkah terukur untuk melatih personil lapangan."
> * Waktu: 3 menit.

---

## 11. Syarat Sebelum Pilot Dimulai

Intinya: Uji coba lapangan baru akan dibuka setelah tujuh persyaratan administratif dan operasional dipenuhi secara utuh.

Daftar prasyarat resmi penentu peluncuran pilot:

| No | Persyaratan Operasional & Administratif | Penanggung Jawab | Status Saat Ini |
|:---:|:---|:---|:---:|
| 1 | Penerbitan Surat Perintah Tugas (SPT) Resmi Kadinkes | Kepala Dinas Kesehatan Kota Depok | 🔴 Belum Terbit |
| 2 | Pembuatan Akun Pengguna Riil Staf Faskes dan Kader | TBD | 🔴 Akun pilot belum dibuat; akun demo tidak dipakai di pilot |
| 3 | Pencetakan & Penempatan Fisik Formulir Cadangan | TBD | 🔴 Belum Dicetak |
| 4 | Penggandaan & Penyerahan Buku Saku Panduan Operasional | TBD | 🟡 Bersyarat (Draf Buku Siap) |
| 5 | Aktivasi Kanal WhatsApp Group Siaga Pendampingan | TBD | 🟡 Bersyarat (Struktur Siap) |
| 6 | Penyelenggaraan Pelatihan Singkat & Simulasi Pengguna | TBD | 🔴 Belum Terlaksana |
| 7 | Verifikasi Batasan Wilayah Non-KLB & Stok Aman (≥ 20 L) | TBD | 🟡 Bersyarat (Menunggu Cek Fisik) |

> **Untuk Presenter:**
> * Tunjukkan: Tabel status prasyarat dengan jujur.
> * Sampaikan: "Kami berkomitmen tidak akan meluncurkan pilot sebelum ketujuh syarat ini dipenuhi bersama bapak dan ibu."
> * Waktu: 2 menit.

---

## 12. Keterbatasan yang Kami Sampaikan Terbuka

Intinya: Kami menyampaikan batasan sistem secara jujur agar para pemangku kepentingan memahami mitigasi yang disiapkan bersama.

Keterbatasan operasional dan teknis yang perlu dipahami:

| Batasan yang Diketahui | Konsekuensi Lapangan | Mitigasi / Rencana Mitigasi |
|:---|:---|:---|
| **Saluran Notifikasi Bebas Biaya** | Sistem mengandalkan protokol email (SMTP) dan push notification ponsel (FCM); tidak memakai WhatsApp resmi berbayar. | Rencana: grup WhatsApp koordinasi lapangan (belum aktif; syarat pilot no. 5). |
| **Penyimpanan Foto di Disk Lokal** | Foto kondisi rumah dan bukti kerja disimpan di disk lokal peladen (host), belum terintegrasi kluster object storage mandiri. | Prosedur pencadangan berkas foto belum ada (TBD). |
| **Layanan Peta Publik Gratis** | Tampilan peta mengandalkan server peta terbuka (OpenStreetMap) yang rentan melambat bila diakses serentak. | Rekomendasi: Diskominfo memasang cache peta lokal (belum dilaksanakan). |
| **Kewajiban Berkas Fisik Serah Terima** | Sistem mencatat mutasi digital, namun regulasi penatausahaan aset daerah mewajibkan berkas fisik tertulis. | Petugas logistik tetap menandatangani dokumen Berita Acara Serah Terima (BAST) fisik. |
| **Cadangan belum off-site** | File backup tersimpan lokal di host | Penyimpanan off-site belum tersedia (TBD). |

> **Untuk Presenter:**
> * Tunjukkan: Tabel keterbatasan sistem secara transparan.
> * Sampaikan: "Transparansi adalah kunci agar tidak ada kejutan di lapangan saat pilot berjalan."
> * Waktu: 3 menit.

---

## 13. Keputusan yang Masih Menunggu

Intinya: Tiga agenda tata kelola birokrasi menunggu arahan pimpinan untuk peluncuran skala penuh, namun tidak menghalangi uji coba terbatas.

Tiga keputusan tata kelola yang dimohonkan arahannya:

1. **Alur Permintaan Tambahan Stok Insektisida:**
   * *Pertanyaan:* Apakah draf permintaan stok dari sistem harus ditandatangani Kepala Puskesmas terlebih dahulu sebelum masuk ke Dinas Kesehatan, atau boleh otomatis terbit? Saat ini draf langsung masuk antrean Dinkes berstatus PENDING tanpa persetujuan Kepala Puskesmas.
   * *Pemilik Keputusan:* Kepala Dinas Kesehatan Kota Depok & Kepala UPTD Puskesmas se-Kota Depok.
   * *Dampak ke Pilot:* Tidak menghambat pilot terbatas selama stok kedua faskes terverifikasi >= 20 liter (syarat pilot no. 7, saat ini menunggu cek fisik).

2. **Fleksibilitas Kapasitas Tugas Fogging saat Kejadian Luar Biasa:**
   * *Pertanyaan:* Apakah sistem boleh membuka batas 1 tugas per hari per regu jika terjadi lonjakan kasus darurat di suatu kelurahan?
   * *Pemilik Keputusan:* Kepala Dinas Kesehatan Kota Depok & Kepala Bidang P2P.
   * *Dampak ke Pilot:* Tidak menghambat pilot terbatas selama kedua faskes berstatus non-KLB (syarat pilot no. 7).

3. **Batas Waktu Pengajuan Sanggahan oleh Warga:**
   * *Pertanyaan:* Berapa batas hari maksimal bagi warga untuk mengajukan sanggahan (misalnya 48 jam atau 72 jam) setelah laporannya ditolak?
   * *Pemilik Keputusan:* Kepala Dinas Kesehatan Kota Depok & Tim Bidang P2P.
   * *Dampak ke Pilot:* Tidak menghambat pilot karena banding saat ini dibatasi 1 kali per laporan.

> **Untuk Presenter:**
> * Tunjukkan: Tiga poin keputusan kebijakan.
> * Sampaikan: "Kami tidak mendahului wewenang pimpinan; keputusan ini kami serahkan sepenuhnya kepada arahan bapak dan ibu."
> * Waktu: 3 menit.

---

## 14. Langkah Selanjutnya

Intinya: Rangkaian tahapan tindak lanjut yang terstruktur menuju pelaksanaan uji coba terbatas di lapangan.

Jadwal pelaksanaan tahapan tindak lanjut dirancang berurutan:

```text
[Simulasi Kering Tim (Dry-Run)] -> Jadwal: Belum terlaksana (prasyarat sebelum demo langsung)
       ↓ 
[Paparan Eksekutif Hari Ini]    -> Moda: Slide & tangkapan layar (35 menit)
       ↓ 
[Penerbitan SPT Resmi Kadinkes] -> Jadwal: Menunggu disposisi pimpinan dinas
       ↓ 
[Pelatihan Petugas & Kader]     -> Jadwal: Ditetapkan setelah SPT terbit
       ↓ 
[Kick-Off Pilot]                -> Jadwal: Puskesmas Beji & Pancoran Mas (durasi 30 hari kalender sejak kick-off (tanggal ditetapkan setelah SPT terbit))
       ↓ 
[Evaluasi & Laporan Akhir]      -> Jadwal: Hari ke-30 pasca peluncuran pilot
```

Tahapan simulasi kering tim (dry-run) ditempatkan sebelum sesi demonstrasi langsung agar seluruh pemateri menguasai alur sistem. Segera setelah Surat Perintah Tugas ditandatangani, tim akan menyelenggarakan pelatihan teknis dan mendistribusikan buku saku panduan operasional.

> **Untuk Presenter:**
> * Tunjukkan: Bagan langkah lanjutan kepada para pimpinan.
> * Sampaikan: "Segera setelah SPT ditandatangani, tim siap menyelenggarakan pelatihan dan membagikan buku panduan."
> * Waktu: 2 menit.

---

## 15. Tanya Jawab Antisipasi

Intinya: Rangkuman jawaban lugas berbasis dokumen atas pertanyaan yang berpotensi diajukan oleh para pimpinan daerah.

### Kategori Tata Kelola & Bisnis
1. **Siapa pemilik resmi dari sistem DBD Guard ini?**
   Dinas Kesehatan Kota Depok bertindak sebagai pemilik bisnis dan penanggung jawab data kesehatan, sedangkan Dinas Komunikasi dan Informatika terlibat dalam penyediaan server dan layanan pendukung. Dasar hukum formal berupa Peraturan Walikota mengenai kepemilikan definitif saat ini belum ditetapkan (TBD).
2. **Apakah puskesmas lain di luar Beji dan Pancoran Mas wajib ikut menggunakan aplikasi ini?**
   Belum. Kewajiban saat ini hanya berlaku bagi Puskesmas Beji dan Puskesmas Pancoran Mas selama 30 hari kalender sejak kick-off (tanggal ditetapkan setelah SPT terbit). Regulasi kewajiban menyeluruh bagi seluruh faskes belum ditetapkan (TBD) dan akan dirumuskan setelah evaluasi pilot.

### Kategori Keamanan & Perlindungan Data
3. **Bagaimana menjamin foto rumah warga dan rekam medik pelapor tidak bocor ke publik?**
   Sistem ini bukan rekam medis elektronik. Foto dan lampiran medis hanya dapat diakses petugas Puskesmas dan Dinkes. Publik hanya melihat nomor resi tanpa identitas pelapor.
4. **Apakah kader kelurahan lain bisa mengklaim laporan wilayah tetangganya?**
   Tidak bisa. Sistem secara otomatis menolak akses setiap upaya klaim tugas di luar kelurahan terdaftar pengguna demi menjaga integritas wilayah kerja.

### Kategori Keputusan Medis (DSS)
5. **Apakah rekomendasi fogging dari aplikasi bersifat mengikat keputusan dokter?**
   Tidak. Rekomendasi sistem murni bersifat saran pendukung berbasis pedoman teknis. Kepala Puskesmas memiliki wewenang penuh untuk mengubah keputusan dengan mencantumkan alasan medis tertulis.

### Kategori Logistik & Keuangan Daerah
6. **Bagaimana memastikan obat semprot di puskesmas tidak hilang atau selisih pembukuan?**
   Setiap pengurangan stok wajib dikaitkan dengan tiket penugasan fogging yang memiliki foto bukti di lapangan. Perpindahan barang antar-faskes wajib dikonfirmasi fisik oleh petugas penerima sebelum saldo berpindah.

### Kategori Operasional Pilot
7. **Mengapa hanya Puskesmas Beji dan Puskesmas Pancoran Mas yang dipilih untuk tahap pertama?**
   Belum ditetapkan; alasan pemilihan kedua puskesmas akan dijelaskan bersama Dinkes. Penetapan resmi wilayah pilot menunggu SPT Kepala Dinas Kesehatan (saat ini belum terbit).
8. **Kapan hasil evaluasi pilot akan dilaporkan kepada Walikota?**
   Evaluasi berkala dipantau oleh tim surveilans dinas kesehatan, dan evaluasi komprehensif dijadwalkan pada hari ke-30 masa uji coba. Tanggal pasti pelaporan akhir kepada Walikota belum ditetapkan (TBD) karena bergantung pada penetapan tanggal kick-off setelah penerbitan SPT.

---

## Lampiran

### Lampiran A: Glosarium Istilah

1. **Penyelidikan Epidemiologi (PE):** Kegiatan pencarian penderita demam dan pemeriksaan jentik nyamuk di rumah penderita serta radius 100 meter di sekitarnya.
2. **Kader Jumantik:** Relawan warga terlatih yang bertugas memantau keberadaan jentik nyamuk di lingkungan rukun warga.
3. **Sanitarian:** Tenaga kesehatan lingkungan di puskesmas yang bertanggung jawab atas survei penyehatan dan pencegahan penyakit lingkungan.
4. **Sistem Pendukung Keputusan (DSS):** Perangkat lunak yang mengolah data survei jentik dan kasus demam untuk memberikan saran tindakan penanganan yang terstandar.
5. **Diskresi Klinis (Clinical Override):** Hak dokter atau pimpinan faskes untuk mengubah rekomendasi sistem berdasarkan pertimbangan medis lapangan.
6. **Prinsip Empat Mata (Four-Eyes Principle):** Aturan keamanan di mana petugas Penyelidikan Epidemiologi tidak boleh memvalidasi hasil PE-nya sendiri.
7. **Sistem Peringatan Dini (EWS):** Fitur pemantauan otomatis yang mendeteksi kenaikan kasus melebihi batas wajar pada tingkat kelurahan.
8. **Fogging Fokus:** Tindakan pengasapan insektisida pada radius tertentu di sekitar lokasi penemuan kasus positif demam berdarah.
9. **Pemberantasan Sarang Nyamuk (PSN):** Gerakan 3M Plus untuk menguras, menutup, dan mendaur ulang tempat penampungan air.
10. **Rencana Kontinjensi (BCP):** Prosedur darurat menggunakan formulir kertas manual saat jaringan internet atau peladen padam.
11. **Berita Acara Serah Terima (BAST):** Dokumen fisik tanda terima resmi perpindahan barang milik daerah antar-pengelola logistik.
12. **Fogging Tertunda:** Rekomendasi sistem saat kriteria fogging terpenuhi tetapi stok insektisida faskes tidak mencukupi, yang secara otomatis memicu draf permohonan stok ke dinas.
13. **Uji Simulasi Kering (Dry-Run):** Latihan peragaan alur sistem oleh tim pemateri dengan stopwatch sebelum sesi presentasi di hadapan pimpinan.
14. **Uji Regresi:** Pengujian ulang menyeluruh untuk memastikan penambahan perbaikan sistem tidak menimbulkan kesalahan pada fitur lain.
15. **Surat Perintah Tugas (SPT):** Dokumen surat resmi dari pimpinan dinas kesehatan sebagai dasar legalitas penugasan uji coba di lapangan.

---

### Lampiran B: Cheat Sheet Presenter

Panduan cepat satu halaman bagi pemateri selama sesi demonstrasi berlangsung:

| Segmen Acara | Alokasi Waktu | Kalimat Pembuka yang Disarankan | Kalimat Penutup yang Disarankan | Tindakan Darurat Jika Kendala Teknis |
|:---|:---:|:---|:---|:---|
| **1. Pembukaan & Latar Belakang** | 5 Menit (Menit 01-05) | "Selamat pagi bapak dan ibu pimpinan, hari ini kami hadir membawa solusi terpadu penanganan DBD Kota Depok." | "...karenanya, DBD Guard hadir menjembatani warga hingga dinas kesehatan." | Lanjutkan slide presentasi; tidak ada ketergantungan sistem. |
| **2. Alur Pelaporan Warga & Kader** | 6 Menit (Menit 06-11) | "Mari kita perhatikan alur saat seorang warga melapor dan diterima kader kelurahannya." | "Laporan tervalidasi di lapangan tanpa menunggu disposisi berhari-hari." | Tampilkan slide alur laporan dan tangkapan layar verifikasi kader. |
| **3. Survei Faskes & Rekomendasi** | 6 Menit (Menit 12-17) | "Setelah pengecekan awal, sanitarian puskesmas mencatat temuan survei ke dalam sistem." | "Sistem langsung memunculkan saran objektif berbasis aturan bisnis sistem." | Tampilkan slide formulir survei dan tabel logika rekomendasi. |
| **4. Persetujuan & Eksekusi Fogging** | 5 Menit (Menit 18-22) | "Kini Kepala Puskesmas meninjau hasil survei untuk memutuskan jadwal penyemprotan." | "Tindakan fogging disetujui, dan foto bukti kerja wajib diunggah." | Tampilkan slide persetujuan faskes dan bukti foto lapangan. |
| **5. Pemotongan Stok & Peta Pantau** | 3 Menit (Menit 23-25) | "Begitu penyemprotan tuntas, stok obat terpotong otomatis dan peta wilayah diperbarui." | "Logistik terpantau akurat tanpa risiko selisih pembukuan gudang." | Tampilkan slide rekonsiliasi saldo dan pembaruan peta wilayah. |
| **6. Uji Ketahanan & Batas Wilayah** | 3 Menit (Menit 26-28) | "Sistem ini dirancang tangguh terhadap upaya penyalahgunaan wewenang dan wilayah." | "Terbukti bahwa batas yurisdiksi dan data warga terkunci secara aman." | Tampilkan slide bukti penolakan klaim tugas lintas wilayah. |
| **7. Diskusi, Cadangan & Tanya Jawab**| 7 Menit (Menit 29-35) | "Kami persilakan bapak dan ibu memberikan arahan serta pertanyaan kebijakan." | "Kami memohon arahan pimpinan untuk penerbitan surat perintah tugas pilot." | Peragakan skenario cadangan (takeover tugas dan DSS override). |

*Pedoman Kontinjensi Demo:*
* Mengingat berkas rekaman video mandiri dan log simulasi kering (dry-run) belum tersedia, seluruh segmen peragaan mengandalkan slide dan tangkapan layar.
* Skenario cadangan mencakup pengambilalihan penugasan (takeover) dan diskresi klinis dokter (DSS override) disiapkan khusus untuk menjawab pendalaman materi pada sesi tanya jawab. Catatan: Waktu yang tertera pada panduan 'Untuk Presenter' adalah panduan kecepatan bicara, bukan tambahan waktu.

---

### Lampiran C: Tabel Pemetaan Klaim ke Dokumen Sumber

Seluruh fakta, status kesiapan, kontrol keamanan, dan data bisnis dalam dokumen presentasi ini bersumber dari dokumen resmi proyek DBD Guard:

| Topik & Klaim Presentasi | Fakta yang Dipaparkan | Dokumen Sumber Resmi | Bagian / Identifikasi Sumber |
|:---|:---|:---|:---|
| **Status Kesiapan Pilot** | Conditionally Pilot Ready (Uji Coba Terbatas di Beji & Panmas) | Laporan Kesiapan Akhir Pilot | Bagian 1.1, Bagian 15, dan Bagian 16 |
| **Status Demonstrasi** | Demo Ready, menunggu simulasi kering (dry-run) | Panduan Persiapan & Demo Pilot | Bab 1.1, Bab 4.6, dan Bab 22 |
| **Batas Respon Laporan (SLA)** | Batas waktu respon laporan warga maksimal 72 jam | Dokumen Spesifikasi Kebutuhan Produk | PRD Bab 5.1 (Verifikasi Kader SLA 3x24 jam) |
| **Kontrol Batas Wilayah** | Penolakan otomatis klaim penugasan lintas kelurahan | Laporan Verifikasi Batas Wilayah | IMP-001 Bagian 4, 6, dan 7 |
| **Prinsip Empat Mata (SoD)** | Larangan validasi mandiri oleh petugas survei yang sama | Dokumen Spesifikasi Kebutuhan Produk | IMP-001 Bagian 8 (regresi TestP0_03_and_04) |
| **Logika Rekomendasi DSS** | Aturan jentik, demam, dan stok menghasilkan fogging/tertunda/PSN | Dokumen Spesifikasi Kebutuhan Produk | PRD Bab 5.2 dan Aturan Bisnis |
| **Kedaulatan Keputusan Medis** | Rekomendasi sistem bersifat saran; dokter memegang hak pengalihan | Laporan Audit Masalah Kritis | Laporan P0 Bagian 3 (P0-03) |
| **Akuntansi Mutasi Obat** | Saldo pengirim, penerima, dan transit seimbang tanpa minus | Laporan Audit Masalah Kritis | Laporan P0 Bagian 3 (P0-01) |
| **Kewajiban Foto Bukti** | Pemotongan stok fogging wajib menyertakan foto asli lapangan | Dokumen Spesifikasi Kebutuhan Produk | PRD Bab 3, 5.1, dan Aturan Bisnis |
| **Penyimpanan Foto di Disk Lokal**| Berkas foto disimpan pada disk lokal peladen (host) | Laporan Kesiapan Akhir Pilot | Bagian 14 |
| **Layanan Peta Publik** | Tampilan peta mengandalkan OpenStreetMap | Laporan Kesiapan Akhir Pilot | Bagian 14 |
| **Saluran Notifikasi Bebas Biaya** | Notifikasi mengandalkan email SMTP dan FCM tanpa WhatsApp API | Dokumen Manajemen Risiko v2.0 | Bab 5 (RSK-13) |
| **Ambang Stok Aman Pilot** | Batas minimal ketersediaan obat faskes pilot ≥ 20 Liter | Laporan Kesiapan Akhir Pilot | Bagian 11.1 dan Bagian 16.2 |
| **Pencadangan & Pemulihan** | Uji pemulihan data (restore drill) dilaksanakan pada tahap perbaikan dan datanya cocok; salinan cadangan masih di server yang sama. Drill ulang di lingkungan pilot dijadwalkan (TBD). | Laporan P0 Bagian 3.4 dan 5.3; Laporan Kesiapan Akhir Pilot Bagian 3.3 dan 9 |
| **Rencana Darurat Manual (BCP)**| Formulir kertas Kemenkes rangkap 3 jika peladen padam | Dokumen Manajemen Risiko v2.0 | Bab 8 (BCP-01) dan Bab 11 |
| **Tujuh Syarat Pilot** | Daftar 7 prasyarat sebelum operasional pilot dimulai | Panduan Persiapan & Demo Pilot | Bab 3.2, Bab 21, dan Bab 23 |
| **Empat Batasan Terbuka** | Transparansi ketergantungan notifikasi, peta, disk lokal, dan BAST | Laporan Kesiapan Akhir Pilot | Bagian 14 |
| **Tiga Keputusan Birokrasi** | DEC-01 (Stok), DEC-02 (KLB), dan DEC-03 (Banding) | Dokumen Tata Kelola Keputusan Bisnis | Bagian 4, 5, dan 6 |

---
*Dokumen Presentasi Eksekutif DBD Guard ini disusun sebagai materi panduan paparan kepada pimpinan daerah dan pemangku kepentingan Kota Depok (status pengesahan dokumen formal: Belum ditetapkan / TBD).*
