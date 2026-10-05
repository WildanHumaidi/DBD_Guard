import Header from '@/components/layout/Header';
import Footer from '@/components/layout/Footer';
import { ShieldCheck } from 'lucide-react';

export default function KebijakanPrivasi() {
  return (
    <main className="flex min-h-screen flex-col bg-bg font-sans pt-24">
      <Header />
      
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 w-full">
        <div className="bg-surface rounded-3xl p-8 sm:p-12 shadow-xl shadow-slate-200/50 dark:shadow-none border border-border">
          <div className="flex items-center gap-4 mb-8">
            <div className="w-14 h-14 bg-primary/10 text-primary rounded-2xl flex items-center justify-center shrink-0">
              <ShieldCheck className="w-7 h-7" />
            </div>
            <div>
              <h1 className="text-3xl font-extrabold text-text">Kebijakan Privasi</h1>
              <p className="text-text-muted mt-1">Terakhir diperbarui: {new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}</p>
            </div>
          </div>

          <div className="space-y-8 text-text-muted leading-relaxed">
            <section>
              <h2 className="text-xl font-bold text-text mb-3">1. Pendahuluan</h2>
              <p>
                Pemerintah Kota Depok (selanjutnya disebut "Kami") berkomitmen untuk melindungi privasi dan keamanan data pribadi Anda. Kebijakan Privasi ini menjelaskan bagaimana kami mengumpulkan, menggunakan, menyimpan, dan melindungi informasi pribadi yang Anda berikan melalui platform DBD Guard.
              </p>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">2. Informasi yang Kami Kumpulkan</h2>
              <p>Saat Anda menggunakan layanan DBD Guard, khususnya saat membuat laporan (kecuali Anda memilih fitur pelaporan anonim), kami mungkin mengumpulkan informasi berikut:</p>
              <ul className="list-disc pl-5 mt-2 space-y-1">
                <li>Nama lengkap</li>
                <li>Nomor telepon dan/atau alamat email</li>
                <li>Lokasi (koordinat GPS) saat melakukan pelaporan</li>
                <li>Foto atau media pendukung terkait temuan jentik atau kasus DBD</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">3. Penggunaan Informasi</h2>
              <p>Informasi yang dikumpulkan digunakan secara eksklusif untuk tujuan:</p>
              <ul className="list-disc pl-5 mt-2 space-y-1">
                <li>Verifikasi dan tindak lanjut laporan oleh Kader dan Puskesmas setempat.</li>
                <li>Pemetaan (Geographic Information System) untuk area rawan DBD di Kota Depok.</li>
                <li>Pengiriman notifikasi status pelaporan dan informasi kesehatan.</li>
                <li>Keperluan analitik statistik yang telah dianominasi untuk kepentingan Dinas Kesehatan.</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">4. Keamanan dan Perlindungan Data</h2>
              <p>
                Kami menerapkan standar keamanan teknis dan operasional yang ketat, sesuai dengan Undang-Undang Perlindungan Data Pribadi (UU PDP) dan regulasi e-Government. Data Anda disimpan dalam server terenkripsi yang berlokasi di pusat data nasional dan hanya dapat diakses oleh petugas berwenang (Kader, Tenaga Medis Puskesmas, dan Dinas Kesehatan).
              </p>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">5. Data Publik (Open Data)</h2>
              <p>
                Pada laman Data Publik, informasi yang ditampilkan adalah murni agregasi statistik (contoh: jumlah kasus di tingkat Kelurahan). Kami tidak akan pernah mempublikasikan data identitas individu (Nama, Kontak, dan Koordinat persis tingkat rumah tangga) ke ruang publik.
              </p>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">6. Hak Pengguna</h2>
              <p>
                Anda berhak meminta penghapusan kontak pribadi Anda dari database kami setelah kasus atau laporan yang Anda buat dinyatakan selesai (Status: Selesai). Anda juga dapat memilih opsi <strong>"Lapor Secara Anonim"</strong> yang secara otomatis mencegah sistem kami menyimpan kontak pribadi apa pun.
              </p>
            </section>
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
