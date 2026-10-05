import Link from 'next/link';
import { Shield, Phone, Mail, MapPin } from 'lucide-react';

export default function Footer() {
  return (
    <footer className="bg-surface border-t border-border pt-16 pb-8 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-12 mb-12">
          {/* Brand Info */}
          <div className="md:col-span-2">
            <div className="flex items-center gap-3 mb-6">
              <div className="bg-primary p-2 rounded-xl">
                <Shield className="h-6 w-6 text-white" />
              </div>
              <span className="font-bold text-2xl text-text tracking-tight">DBD Guard</span>
            </div>
            <p className="text-text-muted leading-relaxed max-w-sm">
              Sistem informasi terpadu pelaporan, edukasi, dan pencegahan Demam Berdarah Dengue (DBD) untuk masyarakat Kota Depok.
            </p>
          </div>

          {/* Quick Links */}
          <div>
            <h4 className="font-bold text-lg text-text mb-6">Tautan Cepat</h4>
            <ul className="space-y-4 text-sm font-medium text-text-muted">
              <li><Link href="/" className="hover:text-primary transition-colors">Beranda</Link></li>
              <li><Link href="/#layanan" className="hover:text-primary transition-colors">Layanan Kami</Link></li>
              <li><Link href="/#alur" className="hover:text-primary transition-colors">Alur Kerja</Link></li>
              <li><Link href="/#edukasi" className="hover:text-primary transition-colors">Edukasi 3M Plus</Link></li>
              <li><Link href="/#faq" className="hover:text-primary transition-colors">FAQ</Link></li>
              <li><Link href="/track" className="hover:text-primary transition-colors">Lacak Laporan</Link></li>
              <li><Link href="/#publik" className="hover:text-primary transition-colors">Data Publik & Rating</Link></li>
            </ul>
          </div>

          {/* Contact / Help */}
          <div>
            <h4 className="font-bold text-lg text-text mb-6">Hubungi Kami</h4>
            <ul className="space-y-4 text-sm font-medium text-text-muted">
              <li className="flex items-start gap-3">
                <MapPin className="w-5 h-5 text-primary shrink-0" />
                <span>Dinas Kesehatan Kota Depok<br/>Jl. Margonda Raya No. 54</span>
              </li>
              <li className="flex items-center gap-3">
                <Phone className="w-5 h-5 text-primary shrink-0" />
                <span>(021) 119 - Darurat Medis</span>
              </li>
              <li className="flex items-center gap-3">
                <Mail className="w-5 h-5 text-primary shrink-0" />
                <span>dbdguard@depok.go.id</span>
              </li>
            </ul>
          </div>
        </div>

        <div className="border-t border-border pt-8 flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-text-muted font-medium">
          <p>&copy; {new Date().getFullYear()} Pemerintah Kota Depok. Hak Cipta Dilindungi.</p>
          <div className="flex gap-6 flex-wrap justify-center">
            <Link href="/kebijakan-privasi" className="hover:text-primary transition-colors">Kebijakan Privasi</Link>
            <Link href="/syarat-ketentuan" className="hover:text-primary transition-colors">Syarat & Ketentuan</Link>
            <Link href="/login" className="hover:text-primary transition-colors font-bold border-l border-border pl-6">Login Pegawai</Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
