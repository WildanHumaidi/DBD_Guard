"use client";

import { useState, Suspense } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { Search, ArrowLeft, Clock, CheckCircle2, ShieldAlert, Truck, Smartphone, Key, Star, AlertCircle, MessageSquare, Gavel, Shield } from 'lucide-react';
import api from '@/lib/api';

function TrackContent() {
  const searchParams = useSearchParams();
  const initialRef = searchParams.get('ref') || '';
  
  const [searchType, setSearchType] = useState<'ref' | 'phone'>('ref');
  const [refNumber, setRefNumber] = useState(initialRef);
  const [phoneNumber, setPhoneNumber] = useState('');
  const [loading, setLoading] = useState(false);
  const [report, setReport] = useState<any>(null);
  const [phoneReports, setPhoneReports] = useState<any[]>([]);
  const [error, setError] = useState('');

  // Appeal State
  const [showAppealForm, setShowAppealForm] = useState(false);
  const [appealReason, setAppealReason] = useState('');
  const [appealLoading, setAppealLoading] = useState(false);
  const [appealMessage, setAppealMessage] = useState('');
  const [appealError, setAppealError] = useState('');

  // Feedback State
  const [showFeedbackForm, setShowFeedbackForm] = useState(false);
  const [feedbackRating, setFeedbackRating] = useState(0);
  const [feedbackHover, setFeedbackHover] = useState(0);
  const [feedbackNotes, setFeedbackNotes] = useState('');
  const [feedbackLoading, setFeedbackLoading] = useState(false);
  const [feedbackMessage, setFeedbackMessage] = useState('');
  const [feedbackError, setFeedbackError] = useState('');

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (searchType === 'ref' && !refNumber.trim()) return;
    if (searchType === 'phone' && !phoneNumber.trim()) return;
    
    setLoading(true);
    setError('');
    setReport(null);
    setPhoneReports([]);
    resetAppealFeedback();
    
    try {
      if (searchType === 'ref') {
        const response = await api.get(`/reports/track/${refNumber.trim()}`);
        if (response.data && response.data.data) {
          setReport(response.data.data);
        } else {
          setError('Laporan tidak ditemukan.');
        }
      } else {
        const response = await api.get(`/reports/phone/${phoneNumber.trim()}`);
        if (response.data && response.data.data && response.data.data.length > 0) {
          setPhoneReports(response.data.data);
        } else {
          setError('Tidak ada laporan yang ditemukan untuk nomor handphone tersebut.');
        }
      }
    } catch (err: any) {
      if (err.response?.status === 404) {
        setError('Data tidak ditemukan.');
      } else {
        setError('Terjadi kesalahan sistem. Silakan coba lagi nanti.');
      }
    } finally {
      setLoading(false);
    }
  };

  const resetAppealFeedback = () => {
    setShowAppealForm(false);
    setAppealReason('');
    setAppealMessage('');
    setAppealError('');
    setShowFeedbackForm(false);
    setFeedbackRating(0);
    setFeedbackNotes('');
    setFeedbackMessage('');
    setFeedbackError('');
  };

  // Submit Appeal
  const handleSubmitAppeal = async () => {
    if (!appealReason.trim()) {
      setAppealError('Alasan banding wajib diisi.');
      return;
    }
    setAppealLoading(true);
    setAppealError('');
    try {
      const res = await api.post(`/reports/${report.report_id}/appeal`, {
        reason: appealReason.trim(),
      });
      setAppealMessage(res.data?.message || 'Banding berhasil diajukan!');
      setShowAppealForm(false);
      // Update local report status to reflect appeal
      setReport({ ...report, status: 'APPEALED' });
    } catch (err: any) {
      setAppealError(
        err.response?.data?.message || err.response?.data?.error || 'Gagal mengajukan banding. Coba lagi nanti.'
      );
    } finally {
      setAppealLoading(false);
    }
  };

  // Submit Feedback
  const handleSubmitFeedback = async () => {
    if (feedbackRating < 1 || feedbackRating > 5) {
      setFeedbackError('Silakan pilih rating 1-5 terlebih dahulu.');
      return;
    }
    setFeedbackLoading(true);
    setFeedbackError('');
    try {
      const res = await api.post(`/reports/${report.reference_no}/feedback`, {
        rating: feedbackRating,
        notes: feedbackNotes.trim(),
      });
      setFeedbackMessage(res.data?.message || 'Terima kasih atas penilaian Anda!');
      setShowFeedbackForm(false);
      // Update local report to reflect given feedback
      setReport({ ...report, feedback_rating: feedbackRating, feedback_notes: feedbackNotes });
    } catch (err: any) {
      setFeedbackError(
        err.response?.data?.message || err.response?.data?.error || 'Gagal mengirim penilaian. Coba lagi nanti.'
      );
    } finally {
      setFeedbackLoading(false);
    }
  };

  const [cancelLoading, setCancelLoading] = useState(false);
  const [cancelError, setCancelError] = useState('');

  const handleCancelReport = async () => {
    if (!confirm('Apakah Anda yakin ingin membatalkan laporan ini?')) return;
    setCancelLoading(true);
    setCancelError('');
    try {
      await api.post(`/reports/${report.reference_no}/cancel`);
      setReport({ ...report, status: 'CANCELLED_BY_USER' });
    } catch (err: any) {
      setCancelError(
        err.response?.data?.message || err.response?.data?.error || 'Gagal membatalkan laporan.'
      );
    } finally {
      setCancelLoading(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'MENUNGGU_PE_KADER': return <Clock className="h-6 w-6 text-status-wait" />;
      case 'MENUNGGU_VALIDASI_PUSKESMAS': return <CheckCircle2 className="h-6 w-6 text-status-verified" />;
      case 'PE_SCHEDULED': 
      case 'FOGGING_SCHEDULED': return <Truck className="h-6 w-6 text-indigo-500" />;
      case 'COMPLETED': return <CheckCircle2 className="h-6 w-6 text-primary" />;
      case 'REJECTED': return <ShieldAlert className="h-6 w-6 text-red-500" />;
      case 'APPEALED': return <Gavel className="h-6 w-6 text-status-scheduled" />;
      default: return <Clock className="h-6 w-6 text-text-muted" />;
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'MENUNGGU_PE_KADER': return 'Menunggu Kunjungan PE';
      case 'MENUNGGU_VALIDASI_PUSKESMAS': return 'Menunggu Validasi Puskesmas';
      case 'PE_SCHEDULED': return 'Jadwal PE Ditetapkan';
      case 'FOGGING_SCHEDULED': return 'Jadwal Fogging Ditetapkan';
      case 'COMPLETED': return 'Selesai Ditangani';
      case 'REJECTED': return 'Laporan Ditolak';
      case 'REJECTED_HOAX': return 'Laporan Ditolak (Indikasi Hoax)';
      case 'ESKALASI_PUSKESMAS': return 'Eskalasi (RT Tidak Responsif)';
      case 'APPEALED': return 'Dalam Proses Banding';
      default: return status;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'MENUNGGU_PE_KADER': return 'bg-status-wait/20 text-status-wait border-status-wait/30';
      case 'MENUNGGU_VALIDASI_PUSKESMAS': return 'bg-status-verified/20 text-status-verified border-status-verified/30';
      case 'ESKALASI_PUSKESMAS': return 'bg-amber-100 text-amber-800 border-amber-200';
      case 'PE_SCHEDULED': 
      case 'FOGGING_SCHEDULED': return 'bg-indigo-100 text-indigo-800 border-indigo-200';
      case 'COMPLETED': return 'bg-primary-tint text-primary border-primary';
      case 'REJECTED': 
      case 'REJECTED_HOAX': return 'bg-red-100 text-red-800 border-red-200';
      case 'APPEALED': return 'bg-status-scheduled/20 text-status-scheduled border-status-scheduled/30';
      default: return 'bg-bg text-text border-border';
    }
  };

  // Determine if appeal should be shown (max 1x appeal)
  const canAppeal = report && (report.status === 'REJECTED' || report.status === 'REJECTED_HOAX') && (!report.appeal_count || report.appeal_count < 1) && !appealMessage;
  // Determine if feedback should be shown
  const canFeedback = report && report.status === 'COMPLETED' && !report.feedback_rating && !feedbackMessage;
  // Determine if cancel should be shown
  const canCancel = report && (report.status === 'MENUNGGU_PE_KADER') && !report.investigator_id;

  return (
    <main className="min-h-screen bg-bg relative overflow-hidden flex flex-col">
      <div className="relative bg-primary pt-12 pb-32 px-4 sm:px-6 overflow-hidden">
        {/* Animated Blobs */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute top-[-10%] left-[-10%] w-96 h-96 bg-primary-tint/40 rounded-full mix-blend-multiply filter blur-3xl opacity-70 animate-blob"></div>
          <div className="absolute top-[20%] right-[-10%] w-80 h-80 bg-amber-400/30 rounded-full mix-blend-multiply filter blur-3xl opacity-70 animate-blob animation-delay-400"></div>
        </div>

        <div className="max-w-4xl mx-auto relative z-10 animate-fade-in-up">
          <Link href="/" className="inline-flex items-center text-primary-tint hover:text-white mb-8 transition-colors drop-shadow-sm font-medium">
            <ArrowLeft className="h-5 w-5 mr-2" />
            Beranda
          </Link>
          <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-3 drop-shadow-md">Lacak Laporan</h1>
          <p className="text-primary-tint text-lg font-medium drop-shadow-sm leading-relaxed">Pantau status penanganan laporan Anda secara real-time.</p>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-4 sm:px-6 -mt-24 pb-20">
        <div className="glass-panel p-6 sm:p-10 mb-8 relative z-10 animate-fade-in-up" style={{ animationDelay: '0.1s' }}>
          
          {/* Tabs */}
          <div className="flex border-b border-border/50 mb-8">
            <button
              type="button"
              onClick={() => { setSearchType('ref'); setError(''); setReport(null); setPhoneReports([]); resetAppealFeedback(); }}
              className={`flex-1 py-3 text-center font-medium transition-colors flex items-center justify-center gap-2 ${
                searchType === 'ref' ? 'border-b-2 border-primary text-primary' : 'text-text-muted hover:text-text'
              }`}
            >
              <Key className="w-4 h-4" />
              Gunakan Nomor Referensi
            </button>
            <button
              type="button"
              onClick={() => { setSearchType('phone'); setError(''); setReport(null); setPhoneReports([]); resetAppealFeedback(); }}
              className={`flex-1 py-3 text-center font-medium transition-colors flex items-center justify-center gap-2 ${
                searchType === 'phone' ? 'border-b-2 border-primary text-primary' : 'text-text-muted hover:text-text'
              }`}
            >
              <Smartphone className="w-4 h-4" />
              Lupa Nomor Referensi?
            </button>
          </div>

          <form onSubmit={handleSearch}>
            {searchType === 'ref' ? (
              <>
                <label className="block text-sm font-medium text-text mb-2">Nomor Referensi Laporan</label>
                <div className="flex flex-col sm:flex-row gap-3">
                  <input 
                    type="text" 
                    value={refNumber}
                    onChange={(e) => setRefNumber(e.target.value)}
                    placeholder="Misal: REP-20260909-XXXX"
                    className="flex-1 px-4 py-4 rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-primary outline-none text-lg font-mono tracking-wide transition-all uppercase"
                  />
                  <button 
                    type="submit"
                    disabled={loading || !refNumber.trim()}
                    className="px-8 py-4 bg-primary-deep text-primary-foreground rounded-xl font-bold hover:opacity-90 transition-all disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                  >
                    {loading ? 'Mencari...' : <><Search className="h-5 w-5" /> Cari</>}
                  </button>
                </div>
              </>
            ) : (
              <>
                <label className="block text-sm font-medium text-text mb-2">Nomor Handphone Pelapor</label>
                <div className="flex flex-col sm:flex-row gap-3">
                  <input 
                    type="tel" 
                    value={phoneNumber}
                    onChange={(e) => setPhoneNumber(e.target.value)}
                    placeholder="Misal: 08123456789"
                    className="flex-1 px-4 py-4 rounded-xl border border-border focus:ring-2 focus:ring-primary focus:border-primary outline-none text-lg tracking-wide transition-all"
                  />
                  <button 
                    type="submit"
                    disabled={loading || !phoneNumber.trim()}
                    className="px-8 py-4 bg-primary-deep text-primary-foreground rounded-xl font-bold hover:opacity-90 transition-all disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                  >
                    {loading ? 'Mencari...' : <><Search className="h-5 w-5" /> Cari</>}
                  </button>
                </div>
                <p className="text-xs text-text-muted mt-2">Masukan nomor HP yang Anda gunakan saat membuat laporan.</p>
              </>
            )}
          </form>

          
          {error && (
            <div className="mt-6 p-4 bg-red-50 text-red-600 rounded-xl border border-red-100 text-center">
              {error}
            </div>
          )}
        </div>

        {/* Phone Search Results Summary */}
        {phoneReports.length > 0 && (
          <div className="glass-panel p-6 animate-fade-in-up mb-8" style={{ animationDelay: '0.2s' }}>
            <h3 className="font-bold text-text mb-4 border-b border-border/50 pb-3">Riwayat Laporan Anda</h3>
            <div className="space-y-4">
              {phoneReports.map((r: any) => (
                <div key={r.reference_no} className="p-4 rounded-xl border border-border hover:border-primary hover:bg-primary-tint transition-colors flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                  <div>
                    <h4 className="font-bold text-text font-mono">{r.reference_no}</h4>
                    <p className="text-sm text-text-muted">{new Date(r.created_at).toLocaleDateString('id-ID', {
                      day: 'numeric', month: 'long', year: 'numeric'
                    })} • {r.kecamatan}</p>
                  </div>
                  <div className="flex items-center gap-3 w-full sm:w-auto justify-between">
                    <span className={`px-3 py-1 text-xs font-bold rounded-full ${getStatusColor(r.status)}`}>
                      {getStatusText(r.status)}
                    </span>
                    <button
                      onClick={() => {
                        setSearchType('ref');
                        setRefNumber(r.reference_no);
                        setPhoneReports([]);
                        resetAppealFeedback();
                        // Automatically trigger search for this ref
                        setTimeout(() => {
                          const form = document.querySelector('form');
                          if (form) form.requestSubmit();
                        }, 100);
                      }}
                      className="text-sm font-semibold text-primary hover:text-primary"
                    >
                      Lihat Detail &rarr;
                    </button>
                  </div>
                </div>
              ))}
            </div>
            <div className="mt-4 p-3 bg-status-verified/20 text-blue-700 text-xs rounded-lg flex items-start gap-2">
              <ShieldAlert className="w-4 h-4 flex-shrink-0 mt-0.5" />
              <p>Untuk privasi Anda, kami hanya menampilkan sebagian informasi laporan. Klik &quot;Lihat Detail&quot; untuk memantau status secara penuh.</p>
            </div>
          </div>
        )}

        {report && (
          <div className="glass-panel p-6 sm:p-10 animate-fade-in-up" style={{ animationDelay: '0.2s' }}>
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8 pb-6 border-b border-border/50">
              <div>
                <p className="text-sm text-text-muted font-medium mb-1">Status Saat Ini</p>
                <div className="flex items-center gap-3">
                  {getStatusIcon(report.status)}
                  <h2 className="text-2xl font-bold text-text">{getStatusText(report.status)}</h2>
                </div>
              </div>
              <div className={`px-4 py-2 rounded-full border text-sm font-bold ${getStatusColor(report.status)}`}>
                {report.status}
              </div>
            </div>

            <div className="space-y-4 mb-8">
              <div className="grid grid-cols-3 gap-4 border-b border-border pb-4">
                <span className="text-text-muted">Dilaporkan Oleh</span>
                <span className="col-span-2 font-medium text-text">
                  {report.reporter?.anonymous_flag ? 'Anonim' : (report.reporter?.name || '-')}
                </span>
              </div>
              <div className="grid grid-cols-3 gap-4 border-b border-border pb-4">
                <span className="text-text-muted">Tanggal Laporan</span>
                <span className="col-span-2 font-medium text-text">
                  {new Date(report.created_at).toLocaleDateString('id-ID', {
                    day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit'
                  })}
                </span>
              </div>
              <div className="grid grid-cols-3 gap-4">
                <span className="text-text-muted">Lokasi</span>
                <span className="col-span-2 font-medium text-text">
                  {report.location?.address || '-'} 
                  {(report.location?.rt || report.location?.rw) && ` (RT ${report.location?.rt || '-'}/RW ${report.location?.rw || '-'})`}
                </span>
              </div>
            </div>

            {/* Simulated Timeline */}
            <div>
              <h3 className="font-bold text-text mb-6">Riwayat Penanganan</h3>
              <div className="space-y-0 pl-2">
                
                <div className="relative pl-8 border-l-2 border-primary/30 pb-6">
                  <div className="absolute w-5 h-5 rounded-full bg-primary -left-[10.5px] top-1 shadow-[0_0_15px_rgba(20,184,166,0.5)] flex items-center justify-center">
                    <div className="w-2 h-2 rounded-full bg-white animate-pulse"></div>
                  </div>
                  <h4 className="font-bold text-text">Laporan Diterima Sistem</h4>
                  <p className="text-sm text-text-muted mt-1 leading-relaxed">Laporan berhasil masuk ke dalam database DBD Guard dan menunggu verifikasi RT.</p>
                </div>

                {['MENUNGGU_VALIDASI_PUSKESMAS', 'PE_SCHEDULED', 'FOGGING_SCHEDULED', 'COMPLETED'].includes(report.status) && (
                  <div className="relative pl-8 border-l-2 border-blue-500/30 pb-6">
                    <div className="absolute w-5 h-5 rounded-full bg-blue-500 -left-[10.5px] top-1 shadow-[0_0_15px_rgba(59,130,246,0.5)] flex items-center justify-center">
                      <div className="w-2 h-2 rounded-full bg-white animate-pulse"></div>
                    </div>
                    <h4 className="font-bold text-text">Diverifikasi oleh RT/RW</h4>
                    <p className="text-sm text-text-muted mt-1 leading-relaxed">Laporan dinilai valid dan diteruskan ke Puskesmas setempat.</p>
                  </div>
                )}
                
                {['PE_SCHEDULED', 'FOGGING_SCHEDULED', 'COMPLETED'].includes(report.status) && (
                  <div className="relative pl-8 border-l-2 border-indigo-500/30 pb-6">
                    <div className="absolute w-5 h-5 rounded-full bg-indigo-500 -left-[10.5px] top-1 shadow-[0_0_15px_rgba(99,102,241,0.5)] flex items-center justify-center">
                      <div className="w-2 h-2 rounded-full bg-white animate-pulse"></div>
                    </div>
                    <h4 className="font-bold text-text">Penyelidikan Epidemiologi (PE)</h4>
                    <p className="text-sm text-text-muted mt-1 leading-relaxed">Tim Puskesmas menjadwalkan/melakukan survei ke lokasi.</p>
                  </div>
                )}

                {['COMPLETED'].includes(report.status) && (
                  <div className="relative pl-8">
                    <div className="absolute w-5 h-5 rounded-full bg-primary -left-[10.5px] top-1 shadow-[0_0_15px_rgba(20,184,166,0.5)] flex items-center justify-center">
                      <CheckCircle2 className="w-3 h-3 text-white" />
                    </div>
                    <h4 className="font-bold text-text">Penanganan Selesai</h4>
                    <p className="text-sm text-text-muted mt-1 leading-relaxed">Tindakan lapangan (Edukasi, PSN, atau Fogging) telah selesai dilakukan.</p>
                  </div>
                )}

                {report.status === 'REJECTED' && (
                  <div className="relative pl-8 mt-6">
                    <div className="absolute w-5 h-5 rounded-full bg-red-500 -left-[10.5px] top-1 shadow-[0_0_15px_rgba(239,68,68,0.5)] flex items-center justify-center">
                      <ShieldAlert className="w-3 h-3 text-white" />
                    </div>
                    <h4 className="font-bold text-red-600">Laporan Ditolak</h4>
                    <p className="text-sm text-text-muted mt-1 leading-relaxed">Laporan tidak memenuhi kriteria atau merupakan duplikasi.</p>
                  </div>
                )}

                {report.status === 'CANCELLED_BY_USER' && (
                  <div className="relative pl-8 mt-6">
                    <div className="absolute w-5 h-5 rounded-full bg-status-void -left-[10.5px] top-1 shadow-[0_0_15px_rgba(156,163,175,0.5)] flex items-center justify-center">
                      <AlertCircle className="w-3 h-3 text-white" />
                    </div>
                    <h4 className="font-bold text-text-muted">Laporan Dibatalkan</h4>
                    <p className="text-sm text-text-muted mt-1 leading-relaxed">Laporan telah dibatalkan oleh pelapor.</p>
                  </div>
                )}

                {report.status === 'APPEALED' && (
                  <div className="relative pl-8 mt-6">
                    <div className="absolute w-5 h-5 rounded-full bg-orange-500 -left-[10.5px] top-1 shadow-[0_0_15px_rgba(249,115,22,0.5)] flex items-center justify-center">
                      <div className="w-2 h-2 rounded-full bg-white animate-pulse"></div>
                    </div>
                    <h4 className="font-bold text-status-scheduled">Banding Diajukan</h4>
                    <p className="text-sm text-text-muted mt-1 leading-relaxed">Laporan banding Anda sedang ditinjau ulang oleh Admin Dinkes.</p>
                  </div>
                )}

              </div>
            </div>

            {/* === CANCEL SECTION === */}
            {canCancel && (
              <div className="mt-8 pt-6 border-t border-border">
                <div className="bg-bg rounded-2xl p-6 border border-border">
                  <div className="flex items-start gap-4">
                    <div className="flex-1">
                      <h4 className="font-bold text-text mb-1">Batalkan Laporan</h4>
                      <p className="text-sm text-text-muted mb-4">
                        Anda dapat membatalkan laporan ini karena belum diproses oleh Kader.
                      </p>
                      {cancelError && (
                        <div className="mb-4 p-3 bg-red-50 text-red-600 text-sm rounded-lg flex items-center gap-2">
                          <AlertCircle className="h-4 w-4 flex-shrink-0" />
                          {cancelError}
                        </div>
                      )}
                      <button
                        onClick={handleCancelReport}
                        disabled={cancelLoading}
                        className="inline-flex items-center gap-2 px-6 py-3 bg-surface border border-red-200 text-red-600 rounded-xl text-sm font-bold hover:bg-red-50 transition-all active:scale-95 disabled:opacity-50"
                      >
                        {cancelLoading ? 'Memproses...' : 'Batalkan Laporan'}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* === APPEAL SECTION === */}
            {canAppeal && (
              <div className="mt-8 pt-6 border-t border-border">
                {!showAppealForm ? (
                  <div className="bg-red-50 rounded-2xl p-6 border border-red-100">
                    <div className="flex items-start gap-4">
                      <div className="h-10 w-10 bg-red-100 rounded-xl flex items-center justify-center flex-shrink-0">
                        <Gavel className="h-5 w-5 text-red-600" />
                      </div>
                      <div className="flex-1">
                        <h4 className="font-bold text-red-800 mb-1">Laporan Anda Ditolak</h4>
                        <p className="text-sm text-red-600 mb-4">
                          Jika Anda merasa penolakan ini tidak tepat, Anda berhak mengajukan banding untuk ditinjau ulang oleh Puskesmas / Dinas Kesehatan.
                        </p>
                        <button
                          onClick={() => setShowAppealForm(true)}
                          className="inline-flex items-center gap-2 px-6 py-3 bg-red-600 text-white rounded-xl text-sm font-bold hover:bg-red-700 transition-all active:scale-95 shadow-lg shadow-red-600/20"
                        >
                          <Gavel className="h-4 w-4" />
                          Ajukan Banding ke Puskesmas
                        </button>
                      </div>
                    </div>
                  </div>
                ) : (
                  <div className="bg-status-scheduled/20 rounded-2xl p-6 border border-status-scheduled/30">
                    <h4 className="font-bold text-status-scheduled mb-1 flex items-center gap-2">
                      <Gavel className="h-5 w-5" />
                      Formulir Banding
                    </h4>
                    <p className="text-sm text-status-scheduled mb-4">Jelaskan alasan mengapa laporan Anda seharusnya tidak ditolak.</p>
                    <textarea
                      value={appealReason}
                      onChange={(e) => setAppealReason(e.target.value)}
                      placeholder="Tuliskan alasan banding Anda secara jelas dan lengkap..."
                      rows={4}
                      className="w-full p-4 rounded-xl border border-status-scheduled/30 bg-surface focus:ring-2 focus:ring-orange-400 focus:border-orange-400 outline-none text-sm transition-all resize-none"
                    />
                    {appealError && (
                      <div className="mt-3 p-3 bg-red-50 text-red-600 text-sm rounded-lg flex items-center gap-2">
                        <AlertCircle className="h-4 w-4 flex-shrink-0" />
                        {appealError}
                      </div>
                    )}
                    <div className="flex gap-3 mt-4">
                      <button
                        onClick={handleSubmitAppeal}
                        disabled={appealLoading || !appealReason.trim()}
                        className="px-6 py-3 bg-orange-600 text-white rounded-xl text-sm font-bold hover:bg-orange-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {appealLoading ? 'Mengirim...' : 'Kirim Banding'}
                      </button>
                      <button
                        onClick={() => { setShowAppealForm(false); setAppealReason(''); setAppealError(''); }}
                        className="px-6 py-3 bg-surface text-text rounded-xl text-sm font-bold border border-border hover:bg-bg transition-all"
                      >
                        Batal
                      </button>
                    </div>
                  </div>
                )}
              </div>
            )}

            {/* Appeal Success Message */}
            {appealMessage && (
              <div className="mt-8 pt-6 border-t border-border">
                <div className="bg-primary-tint rounded-2xl p-6 border border-primary flex items-start gap-4">
                  <CheckCircle2 className="h-6 w-6 text-primary flex-shrink-0 mt-0.5" />
                  <div>
                    <h4 className="font-bold text-primary mb-1">Banding Berhasil Diajukan</h4>
                    <p className="text-sm text-primary">{appealMessage}</p>
                  </div>
                </div>
              </div>
            )}

            {/* === FEEDBACK SECTION === */}
            {canFeedback && (
              <div className="mt-8 pt-6 border-t border-border">
                {!showFeedbackForm ? (
                  <div className="bg-primary-tint rounded-2xl p-6 border border-emerald-100">
                    <div className="flex items-start gap-4">
                      <div className="h-10 w-10 bg-primary-tint rounded-xl flex items-center justify-center flex-shrink-0">
                        <MessageSquare className="h-5 w-5 text-primary" />
                      </div>
                      <div className="flex-1">
                        <h4 className="font-bold text-primary mb-1">Penanganan Telah Selesai!</h4>
                        <p className="text-sm text-primary mb-4">
                          Bantu kami meningkatkan kualitas layanan. Berikan penilaian untuk penanganan di lokasi Anda.
                        </p>
                        <button
                          onClick={() => setShowFeedbackForm(true)}
                          className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-primary-foreground rounded-xl text-sm font-bold hover:bg-primary-deep transition-all active:scale-95 shadow-lg shadow-primary/20"
                        >
                          <Star className="h-4 w-4" />
                          Beri Penilaian
                        </button>
                      </div>
                    </div>
                  </div>
                ) : (
                  <div className="bg-status-verified/20 rounded-2xl p-6 border border-status-verified/30">
                    <h4 className="font-bold text-status-verified mb-1 flex items-center gap-2">
                      <Star className="h-5 w-5" />
                      Evaluasi Penanganan
                    </h4>
                    <p className="text-sm text-status-verified mb-5">Seberapa puas Anda dengan penanganan yang dilakukan?</p>
                    
                    {/* Star Rating */}
                    <div className="flex items-center gap-2 mb-2">
                      {[1, 2, 3, 4, 5].map((star) => (
                        <button
                          key={star}
                          type="button"
                          onMouseEnter={() => setFeedbackHover(star)}
                          onMouseLeave={() => setFeedbackHover(0)}
                          onClick={() => setFeedbackRating(star)}
                          className="transition-transform hover:scale-125 active:scale-95"
                        >
                          <Star
                            className={`h-10 w-10 transition-colors ${
                              star <= (feedbackHover || feedbackRating)
                                ? 'text-amber-400 fill-amber-400'
                                : 'text-border'
                            }`}
                          />
                        </button>
                      ))}
                    </div>
                    <p className="text-xs text-text-muted mb-4">
                      {feedbackRating === 1 && 'Sangat Buruk'}
                      {feedbackRating === 2 && 'Buruk'}
                      {feedbackRating === 3 && 'Cukup'}
                      {feedbackRating === 4 && 'Baik'}
                      {feedbackRating === 5 && 'Sangat Baik'}
                      {feedbackRating === 0 && 'Klik bintang untuk memberi rating'}
                    </p>

                    <textarea
                      value={feedbackNotes}
                      onChange={(e) => setFeedbackNotes(e.target.value)}
                      placeholder="(Opsional) Berikan catatan atau saran tambahan..."
                      rows={3}
                      className="w-full p-4 rounded-xl border border-status-verified/30 bg-surface focus:ring-2 focus:ring-blue-400 focus:border-blue-400 outline-none text-sm transition-all resize-none"
                    />
                    {feedbackError && (
                      <div className="mt-3 p-3 bg-red-50 text-red-600 text-sm rounded-lg flex items-center gap-2">
                        <AlertCircle className="h-4 w-4 flex-shrink-0" />
                        {feedbackError}
                      </div>
                    )}
                    <div className="flex gap-3 mt-4">
                      <button
                        onClick={handleSubmitFeedback}
                        disabled={feedbackLoading || feedbackRating === 0}
                        className="px-6 py-3 bg-blue-600 text-white rounded-xl text-sm font-bold hover:bg-blue-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {feedbackLoading ? 'Mengirim...' : 'Kirim Penilaian'}
                      </button>
                      <button
                        onClick={() => { setShowFeedbackForm(false); setFeedbackRating(0); setFeedbackNotes(''); setFeedbackError(''); }}
                        className="px-6 py-3 bg-surface text-text rounded-xl text-sm font-bold border border-border hover:bg-bg transition-all"
                      >
                        Batal
                      </button>
                    </div>
                  </div>
                )}
              </div>
            )}

            {/* Feedback Already Given */}
            {report && report.feedback_rating > 0 && (
              <div className="mt-8 pt-6 border-t border-border">
                <div className="bg-bg rounded-2xl p-6 border border-border">
                  <h4 className="font-bold text-text mb-3 flex items-center gap-2">
                    <CheckCircle2 className="h-5 w-5 text-primary" />
                    Penilaian Anda
                  </h4>
                  <div className="flex gap-1 mb-2">
                    {[1, 2, 3, 4, 5].map((s) => (
                      <Star
                        key={s}
                        className={`h-6 w-6 ${
                          s <= report.feedback_rating
                            ? 'text-amber-400 fill-amber-400'
                            : 'text-border'
                        }`}
                      />
                    ))}
                  </div>
                  {report.feedback_notes && (
                    <p className="text-sm text-text-muted italic">&quot;{report.feedback_notes}&quot;</p>
                  )}
                </div>
              </div>
            )}

            {/* Feedback Success Message */}
            {feedbackMessage && !report.feedback_rating && (
              <div className="mt-8 pt-6 border-t border-border">
                <div className="bg-primary-tint rounded-2xl p-6 border border-primary flex items-start gap-4">
                  <CheckCircle2 className="h-6 w-6 text-primary flex-shrink-0 mt-0.5" />
                  <div>
                    <h4 className="font-bold text-primary mb-1">Penilaian Terkirim</h4>
                    <p className="text-sm text-primary">{feedbackMessage}</p>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Helper / Info Section when empty */}
        {!report && phoneReports.length === 0 && !loading && !error && (
          <div className="grid md:grid-cols-3 gap-6 mt-16 animate-fade-in-up" style={{ animationDelay: '0.3s' }}>
            <div className="bg-surface border border-border p-6 rounded-2xl shadow-lg hover:-translate-y-1 transition-all duration-300">
              <div className="bg-gradient-to-br from-blue-500 to-indigo-600 w-12 h-12 rounded-xl flex items-center justify-center mb-6 shadow-md shadow-blue-500/20">
                <Search className="h-6 w-6 text-white" />
              </div>
              <h3 className="font-bold text-text mb-2 text-lg">Lacak Real-Time</h3>
              <p className="text-sm text-text-muted leading-relaxed">Ketahui persis sampai di mana laporan Anda diproses oleh tim Puskesmas dan Dinas Kesehatan.</p>
            </div>
            <div className="bg-surface border border-border p-6 rounded-2xl shadow-lg hover:-translate-y-1 transition-all duration-300">
              <div className="bg-gradient-to-br from-emerald-500 to-teal-600 w-12 h-12 rounded-xl flex items-center justify-center mb-6 shadow-md shadow-emerald-500/20">
                <CheckCircle2 className="h-6 w-6 text-white" />
              </div>
              <h3 className="font-bold text-text mb-2 text-lg">Tindakan Cepat</h3>
              <p className="text-sm text-text-muted leading-relaxed">Setiap laporan yang diverifikasi akan segera ditindaklanjuti dalam batas SLA 3x24 jam.</p>
            </div>
            <div className="bg-surface border border-border p-6 rounded-2xl shadow-lg hover:-translate-y-1 transition-all duration-300">
              <div className="bg-gradient-to-br from-amber-400 to-orange-500 w-12 h-12 rounded-xl flex items-center justify-center mb-6 shadow-md shadow-amber-500/20">
                <ShieldAlert className="h-6 w-6 text-white" />
              </div>
              <h3 className="font-bold text-text mb-2 text-lg">Privasi Terjaga</h3>
              <p className="text-sm text-text-muted leading-relaxed">Identitas pelapor dijamin aman dan hanya digunakan untuk keperluan konfirmasi oleh petugas resmi.</p>
            </div>
          </div>
        )}
      </div>

      {/* Footer */}
      <footer className="border-t border-border/80 bg-surface/80 backdrop-blur-md py-10 mt-auto">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col md:flex-row justify-between items-center gap-4">
          <div className="flex items-center gap-3">
            <div className="bg-bg p-2 rounded-lg">
              <Shield className="h-5 w-5 text-text-muted" />
            </div>
            <span className="font-bold text-text">DBD Guard</span>
          </div>
          <p className="text-text-muted text-sm font-medium">
            &copy; {new Date().getFullYear()} Dinas Kesehatan Kota Depok.
            Dashboard Transparansi Publik.
          </p>
        </div>
      </footer>
    </main>
  );
}

export default function TrackPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-bg pt-32 text-center text-text-muted">Memuat modul pelacakan...</div>}>
      <TrackContent />
    </Suspense>
  );
}
