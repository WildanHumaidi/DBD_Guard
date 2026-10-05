"use client";

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { ArrowLeft, Send, UploadCloud, CheckCircle, MapPin } from 'lucide-react';
import api from '@/lib/api';

// Dummy data for Depok for MVP
const DEPOK_DATA: Record<string, string[]> = {
  "Beji": ["Beji", "Beji Timur", "Kemiri Muka", "Kukusan", "Pondok Cina", "Tanah Baru"],
  "Pancoran Mas": ["Depok", "Depok Jaya", "Mampang", "Pancoran Mas", "Rangkapan Jaya", "Rangkapan Jaya Baru"],
  "Cimanggis": ["Cisalak Pasar", "Curug", "Harjamukti", "Mekarsari", "Pasir Gunung Selatan", "Tugu"],
  "Sukmajaya": ["Abadijaya", "Bakti Jaya", "Cisalak", "Mekarjaya", "Sukmajaya", "Tirtajaya"],
  "Cilodong": ["Cilodong", "Jatimulya", "Kalibaru", "Kalimulya", "Sukamaju"],
  "Limo": ["Grogol", "Krukut", "Limo", "Meruyung"],
  "Cinere": ["Cinere", "Gandul", "Pangkalan Jati", "Pangkalan Jati Baru"],
  "Cipayung": ["Bojong Pondok Terong", "Cipayung", "Cipayung Jaya", "Pondok Jaya", "Ratujaya"],
  "Sawangan": ["Bedahan", "Cinangka", "Kedaung", "Pasir Putih", "Pengasinan", "Sawangan Baru", "Sawangan Lama"],
  "Bojongsari": ["Bojongsari Baru", "Bojongsari Lama", "Curug", "Duren Mekar", "Duren Seribu", "Pondok Petir", "Serua"],
  "Tapos": ["Cilangkap", "Cimpaeun", "Jatijajar", "Leuwinanggung", "Sukamaju Baru", "Sukatani", "Tapos"]
};

// Coordinate mapping for each Kecamatan in Depok
const KECAMATAN_COORDS: Record<string, {lat: number, lng: number}> = {
  "Beji": { lat: -6.3711, lng: 106.8227 },
  "Pancoran Mas": { lat: -6.4020, lng: 106.8126 },
  "Cimanggis": { lat: -6.3804, lng: 106.8653 },
  "Sukmajaya": { lat: -6.3986, lng: 106.8368 },
  "Cilodong": { lat: -6.4293, lng: 106.8407 },
  "Limo": { lat: -6.3770, lng: 106.7797 },
  "Cinere": { lat: -6.3475, lng: 106.7828 },
  "Cipayung": { lat: -6.4278, lng: 106.7925 },
  "Sawangan": { lat: -6.4116, lng: 106.7589 },
  "Bojongsari": { lat: -6.4011, lng: 106.7410 },
  "Tapos": { lat: -6.4237, lng: 106.8778 }
};

export default function ReportPage() {
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [referenceNumber, setReferenceNumber] = useState('');
  
  const [formData, setFormData] = useState({
    reporterName: '',
    reporterPhone: '',
    patientName: '',
    patientAge: '',
    patientAddress: '',
    rt: '',
    rw: '',
    symptoms: '',
    kecamatan: '',
    kelurahan: '',
  });

  const [location, setLocation] = useState<{lat: number, lng: number} | null>(null);
  const [locationError, setLocationError] = useState('');
  const [files, setFiles] = useState<File[]>([]);
  const [agreeTerms, setAgreeTerms] = useState(false);
  const [isAnonymous, setIsAnonymous] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFiles(Array.from(e.target.files));
    }
  };

  useEffect(() => {
    // Try to get user's location automatically
    if ("geolocation" in navigator) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setLocation({
            lat: position.coords.latitude,
            lng: position.coords.longitude
          });
        },
        (error) => {
          console.warn("Geolocation not available or denied:", error.message);
          setLocationError("Gagal mendapatkan lokasi otomatis. Laporan akan menggunakan lokasi kelurahan secara default.");
          // Default to center of Depok
          setLocation({ lat: -6.4025, lng: 106.7942 });
        }
      );
    } else {
      setLocation({ lat: -6.4025, lng: 106.7942 });
    }
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleKecamatanChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedKecamatan = e.target.value;
    setFormData({ 
      ...formData, 
      kecamatan: selectedKecamatan,
      kelurahan: '' // Reset kelurahan when kecamatan changes
    });
    
    // Update location to the center of the selected Kecamatan
    if (KECAMATAN_COORDS[selectedKecamatan]) {
      setLocation(KECAMATAN_COORDS[selectedKecamatan]);
      setLocationError(''); // Clear error if we have a valid fallback location
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.kecamatan || !formData.kelurahan) {
      alert("Harap pilih Kecamatan dan Kelurahan.");
      return;
    }

    setLoading(true);
    
    try {
      const formDataToSend = new FormData();
      formDataToSend.append('name', isAnonymous ? 'Anonymous' : (formData.reporterName || 'Anonymous'));
      formDataToSend.append('contact', isAnonymous ? '' : formData.reporterPhone);
      formDataToSend.append('anonymous_flag', isAnonymous ? 'true' : 'false');
      
      const description = `Pasien: ${formData.patientName || '-'} (Umur: ${formData.patientAge || '-'})\nGejala: ${formData.symptoms}`;
      formDataToSend.append('description', description);
      
      formDataToSend.append('category', 'KASUS_DBD');
      
      if (location) {
        formDataToSend.append('latitude', location.lat.toString());
        formDataToSend.append('longitude', location.lng.toString());
      } else {
        formDataToSend.append('latitude', '-6.4025'); 
        formDataToSend.append('longitude', '106.7942');
      }
      
      formDataToSend.append('address', formData.patientAddress);
      formDataToSend.append('rt', formData.rt);
      formDataToSend.append('rw', formData.rw);
      formDataToSend.append('kelurahan', formData.kelurahan);
      formDataToSend.append('kecamatan', formData.kecamatan);
      
      // Append photos if any
      files.forEach((file) => {
        formDataToSend.append('photos', file);
      });

      const response = await api.post('/reports', formDataToSend, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      
      if (response.data && response.data.data) {
        setReferenceNumber(response.data.data.reference_no);
        setSuccess(true);
      }
    } catch (error: any) {
      console.error('Failed to submit report', error);
      alert(error.response?.data?.error || 'Gagal mengirim laporan. Silakan coba lagi.');
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <main className="min-h-screen bg-bg flex flex-col items-center justify-center p-4 relative overflow-hidden">
        <div className="absolute top-0 right-0 w-96 h-96 bg-primary/10 rounded-full blur-3xl animate-blob pointer-events-none"></div>
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-teal-400/10 rounded-full blur-3xl animate-blob animation-delay-400 pointer-events-none"></div>

        <div className="max-w-md w-full glass-panel p-10 rounded-3xl text-center relative z-10 opacity-0 animate-fade-in-up">
          <div className="w-20 h-20 bg-emerald-500/20 rounded-full flex items-center justify-center mx-auto mb-6 shadow-[0_0_15px_rgba(16,185,129,0.3)]">
            <CheckCircle className="h-10 w-10 text-emerald-500" />
          </div>
          <h2 className="text-2xl font-bold text-text mb-2">Laporan Berhasil!</h2>
          <p className="text-text-muted mb-6">
            Terima kasih atas partisipasi Anda. Laporan Anda telah kami terima dan akan segera diverifikasi oleh petugas.
          </p>
          <div className="bg-surface/50 backdrop-blur-sm border border-border p-5 rounded-2xl mb-8">
            <p className="text-sm text-text-muted mb-1">Nomor Referensi Anda:</p>
            <p className="text-3xl font-mono font-black tracking-wider text-primary glow-text">{referenceNumber}</p>
          </div>
          <div className="flex flex-col gap-3">
            <Link 
              href={`/track?ref=${referenceNumber}`}
              className="w-full py-4 bg-gradient-to-r from-primary to-teal-500 hover:from-primary-deep hover:to-primary text-white rounded-xl font-bold transition-all inline-flex items-center justify-center shadow-[0_0_15px_rgba(45,212,191,0.4)] hover:shadow-[0_0_20px_rgba(45,212,191,0.6)] hover:-translate-y-1"
            >
              Lacak Status Laporan
            </Link>
            <Link 
              href="/"
              className="w-full py-4 bg-surface/50 text-text-muted border border-border rounded-xl font-bold hover:bg-surface hover:text-text transition-all inline-block"
            >
              Kembali ke Beranda
            </Link>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-bg pb-20 overflow-hidden">
      <div className="relative bg-gradient-to-br from-primary-deep via-primary to-teal-800 pt-10 pb-40 px-4 sm:px-6">
        <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-white/10 rounded-full blur-[100px] -translate-y-1/2 translate-x-1/4 animate-blob"></div>
        <div className="absolute bottom-0 left-0 w-[400px] h-[400px] bg-teal-300/20 rounded-full blur-[80px] translate-y-1/2 -translate-x-1/4 animate-blob animation-delay-400"></div>
        
        <div className="max-w-3xl mx-auto relative z-10 opacity-0 animate-fade-in-up">
          <Link href="/" className="inline-flex items-center px-4 py-2 rounded-full bg-white/10 hover:bg-white/20 backdrop-blur-md border border-white/20 text-white mb-8 transition-all">
            <ArrowLeft className="h-4 w-4 mr-2" />
            <span className="text-sm font-semibold">Kembali</span>
          </Link>
          <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-4 tracking-tight">Buat Laporan Baru</h1>
          <p className="text-teal-50 text-lg leading-relaxed max-w-2xl">Bantu kami mempercepat penanganan kasus DBD di lingkungan Anda (Kota Depok) dengan melaporkan kejadian atau temuan jentik.</p>
        </div>
      </div>

      <div className="max-w-3xl mx-auto px-4 sm:px-6 -mt-24 relative z-20">
        <div className="glass-panel rounded-3xl p-6 sm:p-10 opacity-0 animate-fade-in-up animation-delay-200">
          <form onSubmit={handleSubmit} className="space-y-8">
            
            {/* Informasi Pelapor */}
            <div className="space-y-6">
              <h3 className="text-xl font-bold text-text border-b border-border pb-3 flex items-center gap-2">
                <span className="bg-primary/10 text-primary w-8 h-8 rounded-lg flex items-center justify-center text-sm">1</span>
                Informasi Pelapor
              </h3>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-text">
                    Nama Pelapor {!isAnonymous && <span className="text-red-500">*</span>}
                  </label>
                  <input
                    required={!isAnonymous}
                    disabled={isAnonymous}
                    name="reporterName"
                    value={isAnonymous ? 'Anonymous' : formData.reporterName}
                    onChange={handleChange}
                    type="text"
                    className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text disabled:opacity-60"
                    placeholder={isAnonymous ? 'Disembunyikan (Anonim)' : 'Nama Anda'}
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-text">
                    Nomor WhatsApp {!isAnonymous && <span className="text-red-500">*</span>}
                  </label>
                  <input
                    required={!isAnonymous}
                    disabled={isAnonymous}
                    name="reporterPhone"
                    value={isAnonymous ? '' : formData.reporterPhone}
                    onChange={handleChange}
                    type="tel"
                    className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text disabled:opacity-60"
                    placeholder={isAnonymous ? 'Disembunyikan (Anonim)' : '081234567890'}
                  />
                </div>
              </div>

              {/* Checkbox Anonim (Rule 5) */}
              <div className="flex items-center gap-3 p-3.5 bg-surface/30 rounded-xl border border-border">
                <input
                  id="anonymous_flag"
                  type="checkbox"
                  checked={isAnonymous}
                  onChange={(e) => setIsAnonymous(e.target.checked)}
                  className="w-4 h-4 rounded text-primary focus:ring-primary border-border cursor-pointer"
                />
                <label htmlFor="anonymous_flag" className="text-sm font-medium text-text cursor-pointer select-none">
                  Laporkan secara Anonim (Nama dan kontak Anda tidak akan disimpan ataupun dibagikan)
                </label>
              </div>
            </div>

            {/* Informasi Pasien/Lokasi */}
            <div className="space-y-6">
              <h3 className="text-xl font-bold text-text border-b border-border pb-3 flex items-center gap-2 mt-8">
                <span className="bg-primary/10 text-primary w-8 h-8 rounded-lg flex items-center justify-center text-sm">2</span>
                Informasi Kejadian / Pasien
              </h3>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-text">Nama Pasien (Jika ada)</label>
                  <input name="patientName" value={formData.patientName} onChange={handleChange} type="text" className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text" placeholder="Kosongkan jika hanya genangan/jentik" />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-text">Umur Pasien (Tahun)</label>
                  <input name="patientAge" value={formData.patientAge} onChange={handleChange} type="number" min="0" className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text" placeholder="Misal: 25" />
                </div>
              </div>
              
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-text">Kecamatan (Depok) <span className="text-red-500">*</span></label>
                  <select required name="kecamatan" value={formData.kecamatan} onChange={handleKecamatanChange} className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text appearance-none cursor-pointer">
                    <option value="" disabled>Pilih Kecamatan</option>
                    {Object.keys(DEPOK_DATA).sort().map(kec => (
                      <option key={kec} value={kec}>{kec}</option>
                    ))}
                  </select>
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-text">Kelurahan <span className="text-red-500">*</span></label>
                  <select required name="kelurahan" value={formData.kelurahan} onChange={handleChange} disabled={!formData.kecamatan} className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text disabled:opacity-50 appearance-none cursor-pointer">
                    <option value="" disabled>Pilih Kelurahan</option>
                    {formData.kecamatan && DEPOK_DATA[formData.kecamatan]?.map(kel => (
                      <option key={kel} value={kel}>{kel}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-semibold text-text">Alamat Lengkap Kejadian <span className="text-red-500">*</span></label>
                <textarea required name="patientAddress" value={formData.patientAddress} onChange={handleChange} rows={2} className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text resize-none" placeholder="Jl. Mawar No. 123"></textarea>
              </div>
              
              <div className="grid grid-cols-2 gap-6">
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-text">RT <span className="text-red-500">*</span></label>
                  <input required name="rt" value={formData.rt} onChange={handleChange} type="text" className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text" placeholder="001" />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-semibold text-text">RW <span className="text-red-500">*</span></label>
                  <input required name="rw" value={formData.rw} onChange={handleChange} type="text" className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text" placeholder="005" />
                </div>
              </div>

              {locationError && (
                <div className="bg-amber-500/10 border border-amber-500/20 text-amber-600 text-sm p-4 rounded-xl flex items-start gap-3">
                  <MapPin className="h-5 w-5 shrink-0 mt-0.5" />
                  <p>{locationError}</p>
                </div>
              )}

              <div className="space-y-2">
                <label className="text-sm font-semibold text-text">Gejala yang dialami (Bila ada)</label>
                <textarea name="symptoms" value={formData.symptoms} onChange={handleChange} rows={2} className="w-full px-4 py-3.5 bg-surface/50 backdrop-blur-sm rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all font-medium text-text resize-none" placeholder="Demam tinggi, muncul bintik merah..."></textarea>
              </div>
            </div>

            {/* Media Upload */}
            <div className="space-y-6">
              <h3 className="text-xl font-bold text-text border-b border-border pb-3 flex items-center gap-2 mt-8">
                <span className="bg-primary/10 text-primary w-8 h-8 rounded-lg flex items-center justify-center text-sm">3</span>
                Bukti Foto
              </h3>
              <div className="relative border-2 border-dashed border-primary/30 rounded-2xl p-10 flex flex-col items-center justify-center text-center bg-surface/30 hover:bg-surface/60 transition-colors cursor-pointer overflow-hidden group">
                <input 
                  type="file" 
                  multiple 
                  accept="image/*" 
                  onChange={handleFileChange}
                  className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
                />
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mb-4 group-hover:scale-110 group-hover:bg-primary/20 transition-all">
                  <UploadCloud className="h-8 w-8 text-primary" />
                </div>
                <p className="text-base font-semibold text-text mb-1">
                  {files.length > 0 ? <span className="text-primary glow-text">{files.length} foto terpilih</span> : 'Klik untuk mengunggah foto'}
                </p>
                <p className="text-sm text-text-muted">Maks. 5MB (JPG, PNG)</p>
              </div>
            </div>

            {/* Persetujuan */}
            <div className="flex items-start gap-4 bg-surface/50 backdrop-blur-sm p-6 rounded-2xl border border-border mt-8">
              <input 
                type="checkbox" 
                id="agreeTerms" 
                checked={agreeTerms}
                onChange={(e) => setAgreeTerms(e.target.checked)}
                className="mt-1 w-5 h-5 rounded border-border text-primary focus:ring-primary focus:ring-offset-bg cursor-pointer shrink-0 transition-all"
                required
              />
              <label htmlFor="agreeTerms" className="text-sm text-text-muted leading-relaxed cursor-pointer select-none">
                Saya menyatakan bahwa data yang saya berikan adalah benar, dan saya menyetujui <Link href="/kebijakan-privasi" className="text-primary font-bold hover:underline" target="_blank">Kebijakan Privasi</Link> serta <Link href="/syarat-ketentuan" className="text-primary font-bold hover:underline" target="_blank">Syarat & Ketentuan</Link> yang berlaku di layanan DBD Guard Kota Depok.
              </label>
            </div>

            <div className="pt-6">
              <button 
                type="submit" 
                disabled={loading || !agreeTerms}
                className="w-full flex items-center justify-center gap-2 bg-gradient-to-r from-primary to-teal-500 hover:from-primary-deep hover:to-primary text-white py-5 rounded-xl font-extrabold text-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-[0_0_20px_rgba(45,212,191,0.4)] hover:shadow-[0_0_25px_rgba(45,212,191,0.6)] hover:-translate-y-1 group"
              >
                {loading ? (
                  <span className="flex items-center gap-2"><div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div> Mengirim...</span>
                ) : (
                  <>Kirim Laporan Sekarang <Send className="h-5 w-5 group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" /></>
                )}
              </button>
            </div>

          </form>
        </div>
      </div>
    </main>
  );
}

