# Palet warna DBD Guard

Aplikasi ini berdiri di antara dua hal yang bertentangan: warga harus merasa cukup tenang untuk membuka aplikasinya, tapi cukup waspada untuk benar-benar melapor. Seluruh keputusan warna di dokumen ini berangkat dari ketegangan itu. Setiap warna punya dua nilai — mode terang dan mode gelap.

---

## Tiga prinsip

**Warna peringatan tidak boleh jadi warna merek.** Tombol utama memakai teal, bukan merah atau oranye. Kalau tombol yang paling sering ditekan berwarna merah, warga kehilangan kemampuan membedakan mana tombol dan mana bahaya.

**Merah disimpan, tidak dibelanjakan.** Merah hanya muncul di satu tempat: penanda KLB. Nilai sebuah peringatan datang dari kelangkaannya. Begitu merah dipakai untuk tombol hapus, notifikasi, dan badge, ia berhenti berarti "bahaya".

**Warna selalu ditemani bentuk.** Sekitar satu dari dua belas laki-laki mengalami buta warna merah-hijau, dan pasangan Aman–Rawan justru yang paling sering tertukar. Setiap status wajib punya ikon dan teks, bukan warna saja.

**Dua mode, satu makna.** Nilai heksa boleh berbeda antara mode terang dan gelap, tapi arti warnanya tidak pernah berubah. Hijau tetap berarti selesai di kedua mode; yang disesuaikan hanya keterbacaannya.

---

## Warna merek

Teal adalah titik temu biru dan hijau. DBD tidak lahir di rumah sakit, ia lahir di air yang tergenang di lingkungan warga. Warna merek ini menunjuk ke sana, bukan ke ruang gawat darurat.

### Teal utama — `#0F766E` / `#2DD4BF`

Dipakai di tombol "Lapor sekarang", tautan, ikon aktif, dan garis progres.

Warna ini harus terasa seperti petugas Puskesmas yang menjelaskan dengan tenang, bukan seperti sirene. Teal cukup klinis untuk dipercaya, tapi tidak membawa rasa darurat yang membuat orang menunda membuka aplikasi. Di mode gelap ia dinaikkan ke versi terang, karena teal pekat di atas latar hitam hampir tidak terbaca dan tombol paling penting di aplikasi ini tidak boleh pernah tenggelam.

### Teal pekat — `#134E4A` / `#0F3B38`

Dipakai di header, bilah navigasi, dan area nama instansi.

Ini warna pihak yang bertanggung jawab. Header adalah tempat nama Kelurahan atau Dinas Kesehatan muncul, dan bagian itu perlu terasa lebih berat daripada isi halaman — laporan warga hanya akan ditulis kalau ada yang jelas berdiri di ujung sana untuk menerimanya. Di mode gelap ia justru dibuat lebih gelap lagi supaya header tetap membumi dan tidak melayang di atas konten.

### Teal muda — `#E6F4F1` / `#123936`

Dipakai di latar kartu edukasi, tips 3M, dan kotak informasi.

Ruang untuk membaca tanpa tekanan. Konten edukasi tentang jentik dan 3M perlu terasa berbeda dari data laporan — bukan sesuatu yang harus ditindaklanjuti hari ini, melainkan sesuatu yang dibaca sambil menunggu. Warna ini menandai bagian aplikasi yang tidak menuntut apa pun dari warga.

---

## Netral

Bagian terbesar dari layar sebenarnya tidak berwarna. Netral yang dipakai condong ke biru agar duduk berdampingan dengan teal tanpa terlihat kusam.

### Teks utama — `#0F172A` / `#E8EDF4`

Dipakai di judul, nama pelapor, alamat, dan angka kasus.

Sengaja bukan hitam murni dan bukan putih murni. Keduanya menghasilkan kontras yang terlalu keras dan membuat mata cepat lelah pada teks panjang — padahal alamat dan detail lokasi justru bagian yang paling harus dibaca teliti oleh petugas. Nada kebiruannya membuat teks terasa satu keluarga dengan warna merek.

### Teks sekunder — `#475569` / `#93A3B8`

Dipakai di waktu laporan, jarak, dan keterangan pendukung.

Batas bawah keterbacaan, bukan sekadar abu-abu yang enak dipandang. Aplikasi ini dibuka sambil berdiri di gang pada siang bolong dengan layar yang silau, jadi abu-abu yang lebih terang dari ini tidak boleh dipakai untuk informasi apa pun yang perlu dibaca.

### Garis pemisah — `#E2E8F0` / `#1F2A3C`

Dipakai di tepi kartu, pemisah daftar, dan batas kolom.

Satu kartu laporan adalah satu rumah dan satu keluarga. Garis di sini menjaga agar laporan tidak melebur menjadi daftar panjang yang kehilangan wujud — tipis saja, cukup untuk menandai batas, tidak sampai membuat layar terlihat seperti tabel administrasi.

### Latar halaman — `#F8FAFC` / `#0B1220`

Dipakai sebagai latar belakang seluruh layar.

Mode gelapnya biru malam, bukan hitam pekat. Petugas fogging bekerja subuh dan malam, dan hitam murni membuat teks putih tampak berpendar serta memaksa pupil bekerja ekstra setiap kali berpindah dari layar ke gang yang gelap. Latar yang sedikit terangkat jauh lebih ramah untuk pemakaian lama di lapangan.

### Permukaan kartu — `#FFFFFF` / `#121A2A`

Dipakai di kartu laporan, panel jadwal, dan kotak dialog.

Sedikit lebih terang dari latar, supaya isi yang bisa ditindaklanjuti selalu terlihat naik ke permukaan. Di mode gelap, perbedaan terang antar lapisan menggantikan peran bayangan — bayangan hampir tidak berfungsi di atas latar gelap.

---

## Status laporan

Lima warna ini menceritakan satu perjalanan: dari ketikan warga di ponselnya sampai asap fogging benar-benar keluar. Urutannya berarti, jadi warnanya juga menghangat lalu mereda.

### Tahap satu — Menunggu verifikasi — `#64748B` / `#94A3B8`

Laporan sudah masuk, kader belum meninjau.

Sengaja belum berwarna. Pada tahap ini belum ada yang bisa dipastikan, dan memberi warna terlalu dini akan membuat warga buru-buru panik atau buru-buru lega. Abu-abu netral adalah cara jujur mengatakan "sudah kami terima, belum kami periksa".

### Tahap dua — Terverifikasi — `#2563EB` / `#60A5FA`

Sudah ditinjau kader atau petugas, tercatat resmi.

Biru adalah warna pencatatan. Momen ini penting secara psikologis: laporan warga berubah dari pesan menjadi data yang diakui. Biru dipilih karena ia terasa administratif tanpa terasa mengancam — birokrasi yang bekerja, bukan yang mendiamkan.

### Tahap tiga — Dijadwalkan fogging — `#B45309` / `#FBBF24`

Tanggal sudah ditetapkan, tim belum datang.

Titik paling hangat dalam seluruh alur, karena ini fase yang paling menggantung: keputusan sudah diambil tapi belum terjadi apa-apa. Amber adalah warna lampu kuning — bersiap, bukan berhenti. Ini juga satu-satunya status yang menuntut warga melakukan sesuatu, yaitu menutup makanan dan mengamankan hewan peliharaan sebelum tim tiba.

### Tahap empat — Fogging selesai — `#15803D` / `#4ADE80`

Tim sudah datang, laporan ditutup.

Hijau daun, bukan hijau neon ala notifikasi pembayaran berhasil. Selesainya fogging bukan perayaan, hanya pemulihan — dan nadanya harus mengingatkan bahwa jentik bisa kembali dalam hitungan minggu kalau genangan tidak diurus. Warna yang terlalu meriah di sini justru berbahaya karena membuat orang merasa urusannya sudah tuntas.

### Bukan kasus — `#57534E` / `#A8A29E`

Laporan ganda, salah alamat, atau tidak terbukti.

Warna paling penting dalam daftar ini, dan alasannya tidak terlihat. Status ini sengaja tidak merah. Warga yang laporannya ternyata tidak terbukti tidak boleh merasa dihukum — kalau sekali saja mereka merasa disalahkan, mereka tidak akan melapor lagi, dan laporan berikutnya yang hilang mungkin justru yang sungguhan. Abu-abu kecokelatan ini menutup laporan tanpa menuduh siapa pun.

---

## Zona risiko

Empat warna ini hanya boleh muncul di peta dan badge wilayah. Isinya menaik dari aman ke KLB, jadi jarak antarwarnanya harus terasa seperti tangga, bukan seperti pilihan.

### Aman — `#16A34A` / `#22C55E`

Tidak ada jentik dan tidak ada kasus.

Hijau di sini berarti "belum ada temuan", bukan "sudah bebas selamanya". Karena itu nadanya sedang saja — hijau yang terlalu cerah akan membuat satu RT merasa kebal dan berhenti memeriksa bak mandi.

### Waspada — `#EAB308` / `#FACC15`

Ditemukan jentik, belum ada kasus manusia.

Inilah tahap di mana pencegahan masih jauh lebih murah daripada fogging, dan warna kuning bertugas menarik mata ke sana lebih dulu. Kuning dipilih, bukan oranye, supaya masih ada satu anak tangga tersisa sebelum keadaan benar-benar memburuk.

### Rawan — `#EA580C` / `#FB923C`

Ada kasus terkonfirmasi di wilayah tersebut.

Batas antara pencegahan dan penanganan. Oranye sengaja dibuat cukup jauh dari kuning agar perpindahan satu wilayah dari Waspada ke Rawan terlihat jelas di peta walaupun dilihat sekilas dari jauh.

### KLB — `#B91C1C` / `#EF4444`

Kejadian Luar Biasa yang ditetapkan Dinas Kesehatan.

Satu-satunya tempat merah diizinkan muncul di seluruh aplikasi. Tidak ada tombol merah, tidak ada notifikasi merah, tidak ada badge merah di mana pun — sehingga ketika warga akhirnya melihat warna ini, mereka tahu tanpa perlu membaca bahwa keadaannya berbeda. Kekuatan warna ini seluruhnya berasal dari seberapa jarang ia dipakai.

---

## Catatan penerapan

**Warna zona di peta pakai isi transparan.** Isi poligon di sekitar 35 persen opasitas dengan garis tepi warna penuh. Kalau solid, nama gang dan jalan di bawahnya hilang, padahal justru itu yang dicari petugas.

**Amber terang tidak boleh jadi warna teks.** Di mode terang, gunakan `#B45309` untuk teks dan `#F59E0B` hanya sebagai isian. Kuning dan amber cerah di atas putih gagal memenuhi rasio kontras minimum.

**Status dan zona boleh berbagi keluarga warna hangat.** Keduanya tidak pernah muncul di permukaan yang sama: status selalu berbentuk chip di dalam kartu, zona selalu berupa bidang di atas peta. Jangan pernah menampilkan keduanya berdampingan tanpa label.

**Ikuti pengaturan sistem, sediakan sakelar manual.** Default-nya mengikuti mode perangkat, tapi tetap sediakan pilihan manual di menu. Petugas lapangan sering perlu memaksa mode terang di bawah matahari langsung.

---

## Daftar lengkap

| Warna | Token | Terang | Gelap | Dipakai untuk |
|---|---|---|---|---|
| Teal utama | `primary` | `#0F766E` | `#2DD4BF` | Tombol utama, tautan, ikon aktif |
| Teal pekat | `primary-deep` | `#134E4A` | `#0F3B38` | Header dan navigasi |
| Teal muda | `primary-tint` | `#E6F4F1` | `#123936` | Latar kartu edukasi |
| Teks utama | `text` | `#0F172A` | `#E8EDF4` | Judul dan isi utama |
| Teks sekunder | `text-muted` | `#475569` | `#93A3B8` | Waktu, jarak, keterangan |
| Garis | `border` | `#E2E8F0` | `#1F2A3C` | Tepi kartu dan pemisah |
| Latar | `bg` | `#F8FAFC` | `#0B1220` | Latar seluruh layar |
| Permukaan | `surface` | `#FFFFFF` | `#121A2A` | Kartu dan panel |
| Menunggu | `status-wait` | `#64748B` | `#94A3B8` | Belum diverifikasi |
| Terverifikasi | `status-verified` | `#2563EB` | `#60A5FA` | Sudah ditinjau petugas |
| Dijadwalkan | `status-scheduled` | `#B45309` | `#FBBF24` | Jadwal fogging ditetapkan |
| Selesai | `status-done` | `#15803D` | `#4ADE80` | Fogging sudah dilakukan |
| Bukan kasus | `status-void` | `#57534E` | `#A8A29E` | Ditutup tanpa tindakan |
| Zona aman | `risk-low` | `#16A34A` | `#22C55E` | Peta dan badge wilayah |
| Zona waspada | `risk-medium` | `#EAB308` | `#FACC15` | Peta dan badge wilayah |
| Zona rawan | `risk-high` | `#EA580C` | `#FB923C` | Peta dan badge wilayah |
| Zona KLB | `risk-outbreak` | `#B91C1C` | `#EF4444` | Hanya untuk status KLB |

---

## Siap tempel: CSS custom properties

```css
:root {
  --primary:          #0F766E;
  --primary-deep:     #134E4A;
  --primary-tint:     #E6F4F1;
  --text:             #0F172A;
  --text-muted:       #475569;
  --border:           #E2E8F0;
  --bg:               #F8FAFC;
  --surface:          #FFFFFF;
  --status-wait:      #64748B;
  --status-verified:  #2563EB;
  --status-scheduled: #B45309;
  --status-done:      #15803D;
  --status-void:      #57534E;
  --risk-low:         #16A34A;
  --risk-medium:      #EAB308;
  --risk-high:        #EA580C;
  --risk-outbreak:    #B91C1C;
}

@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) {
    --primary:          #2DD4BF;
    --primary-deep:     #0F3B38;
    --primary-tint:     #123936;
    --text:             #E8EDF4;
    --text-muted:       #93A3B8;
    --border:           #1F2A3C;
    --bg:               #0B1220;
    --surface:          #121A2A;
    --status-wait:      #94A3B8;
    --status-verified:  #60A5FA;
    --status-scheduled: #FBBF24;
    --status-done:      #4ADE80;
    --status-void:      #A8A29E;
    --risk-low:         #22C55E;
    --risk-medium:      #FACC15;
    --risk-high:        #FB923C;
    --risk-outbreak:    #EF4444;
  }
}
```
