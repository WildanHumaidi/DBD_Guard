import Header from '@/components/layout/Header';
import Footer from '@/components/layout/Footer';
import { FileText } from 'lucide-react';

export default function SyaratKetentuan() {
  return (
    <main className="flex min-h-screen flex-col bg-bg font-sans pt-24">
      <Header />
      
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 w-full">
        <div className="bg-surface rounded-3xl p-8 sm:p-12 shadow-xl shadow-slate-200/50 dark:shadow-none border border-border">
          <div className="flex items-center gap-4 mb-8">
            <div className="w-14 h-14 bg-primary/10 text-primary rounded-2xl flex items-center justify-center shrink-0">
              <FileText className="w-7 h-7" />
            </div>
            <div>
              <h1 className="text-3xl font-extrabold text-text">Syarat & Ketentuan Layanan</h1>
              <p className="text-text-muted mt-1">Berlaku Efektif: {new Date().toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })}</p>
            </div>
          </div>

          <div className="space-y-8 text-text-muted leading-relaxed">
            <section>
              <h2 className="text-xl font-bold text-text mb-3">1. Penerimaan Syarat & Ketentuan</h2>
              <p>
                Dengan mengakses dan menggunakan portal DBD Guard Kota Depok, Anda menyatakan setuju dan tunduk pada seluruh Syarat dan Ketentuan ini, serta mematuhi semua peraturan perundang-undangan yang berlaku di Republik Indonesia.
              </p>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">2. Kewajiban Pengguna</h2>
              <p>Dalam menggunakan layanan pelaporan di sistem DBD Guard, Anda wajib:</p>
              <ul className="list-disc pl-5 mt-2 space-y-1">
                <li>Memberikan data dan keterangan yang jujur, akurat, dan sesuai dengan fakta di lapangan.</li>
                <li>Tidak mengirimkan laporan palsu, spam, ujaran kebencian, atau hal-hal yang melanggar norma sosial.</li>
                <li>Bertanggung jawab secara moral atas pelaporan yang Anda buat demi kemaslahatan lingkungan dan warga.</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">3. Kewenangan Tindak Lanjut</h2>
              <p>
                DBD Guard adalah sarana penunjang bagi Pemerintah Kota Depok. Semua keputusan tindak lanjut operasional lapangan, termasuk namun tidak terbatas pada <strong>Penyelidikan Epidemiologi (PE)</strong> dan penetapan <strong>Fogging Fokus</strong>, merupakan hak prerogatif dan keputusan mutlak dari pihak medis (Kepala Puskesmas) berdasarkan standar operasional dinas kesehatan.
              </p>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">4. Modifikasi Layanan</h2>
              <p>
                Pemerintah Kota Depok berhak penuh untuk mengubah, memperbarui, menambah, atau menghentikan sebagian maupun seluruh layanan pada DBD Guard kapan saja tanpa pemberitahuan sebelumnya, untuk tujuan pemeliharaan (maintenance) atau peningkatan layanan (upgrade).
              </p>
            </section>

            <section>
              <h2 className="text-xl font-bold text-text mb-3">5. Pembebasan Tanggung Jawab</h2>
              <p>
                Kami terus berupaya memastikan sistem ini berjalan lancar. Namun, Pemerintah Kota Depok tidak menjamin bahwa sistem akan selalu terbebas dari gangguan, cacat, atau jeda *server*. Segala kerugian yang timbul secara tidak langsung dari keterlambatan penanganan karena faktor teknis bukan merupakan tanggung jawab hukum Pemerintah Kota.
              </p>
            </section>
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
