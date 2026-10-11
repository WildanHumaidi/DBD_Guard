# Aturan Bisnis & Domain (Business Constraints): DBD Guard Kota Depok

Saat menulis kode, melakukan *refactoring*, atau mendesain API, Anda **WAJIB** menaati aturan bisnis berikut yang ditarik dari PRD. Pelanggaran terhadap aturan ini dianggap *fatal error*:

1.  **SLA Verifikasi RT/RW & Eskalasi**: Laporan warga memiliki batas waktu **3x24 jam** untuk diverifikasi oleh RT/RW. Jika lewat, sistem otomatis mengeskalasi laporan tersebut ke Puskesmas terkait dengan status khusus `"Eskalasi (RT Tidak Responsif)"`.
2.  **Mekanisme Banding Laporan**: Setiap penolakan laporan wajib menyertakan alasan (`reason`). Warga kemudian diberikan opsi banding yang akan meneruskan laporan langsung ke Puskesmas dengan status `"Banding Warga"`.
3.  **Clustering Laporan Duplikat**: Laporan baru yang berada dalam radius **< 100 meter** dari laporan berstatus aktif (dalam 7 hari terakhir) ditandai sebagai `"Kemungkinan Duplikat"` dan dilaporkan ke RT/RW untuk opsi verifikasi silang / gabung.
4.  **Kontradiksi Anonimitas vs Notifikasi**:
    *   Jika `anonymous_flag = true`: Sistem **TIDAK menyimpan** email/kontak, dan **TIDAK mengirimkan** notifikasi apapun (FCM/Email). Warga murni melacak lewat No. Ref/QR.
    *   Jika `anonymous_flag = false`: Notifikasi FCM dan Email otomatis berjalan di setiap perubahan status.
5.  **Master Data Wilayah Pasti**: Hirarki kaku 11 Kecamatan -> 63 Kelurahan -> 38 Puskesmas. Pemetaan Kelurahan ke Puskesmas adalah tepat 1:1. Semua input wilayah harus menggunakan *Dropdown*, DILARANG keras menyediakan *free-text* lokasi untuk konsistensi routing.
6.  **Proteksi Data Privasi Warga**: API / Halaman Publik HANYA menampilkan *agregasi* (contoh: "3 Laporan Genangan di Kel. Sukatani"), BUKAN data telanjang (koordinat presisi tingkat rumah tangga, atau nama asli).
7.  **Siklus Fogging & Inventaris**: 
    *   **Keputusan Fogging BUKAN Otomatis**: Sistem hanya mengeluarkan Rekomendasi (Fogging Fokus / Tertunda / Edukasi). Penjadwalan mutlak ditekan oleh Kepala Puskesmas.
    *   **Anti-Overbooking**: 1 Tim Fogging maksimal hanya memiliki 1 jadwal aktif/tugas per hari.
    *   **Pemotongan Stok Lapangan**: Tim Fogging WAJIB mengunggah foto bukti dan meng-input angka pemakaian stok riil (liter/unit) saat klik "Selesai", yang secara *atomic* memotong stok di Gudang Puskesmas terkait.
8.  **Audit Trail**: Setiap perubahan status laporan, eskalasi, persetujuan jadwal, dan pemotongan stok HUKUMNYA WAJIB dimasukkan ke tabel `AUDIT_LOGS` atau tabel histori terkait. Jangan *update* basis data krusial secara diam-diam.
