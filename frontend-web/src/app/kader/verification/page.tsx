"use client";

import { useEffect, useState, useSyncExternalStore } from "react";
import api from "@/lib/api";
import { CheckCircle2, XCircle, MapPin, User, Calendar, Loader2, Clock, AlertTriangle, Trash2, X, Bug, Thermometer } from "lucide-react";

import Cookies from "js-cookie";

interface KaderReportItem {
  report_id: string;
  reference_no: string;
  category: string;
  description: string;
  status: string;
  created_at: string;
  updated_at: string;
  investigator_id?: string;
  is_possible_duplicate?: boolean;
  location?: {
    address: string;
  };
  reporter?: {
    name?: string;
    anonymous_flag?: boolean;
  };
}

function getAuthUserId(): string | null {
  if (typeof window === "undefined") return null;
  const token = Cookies.get("auth_token");
  if (!token) return null;
  try {
    const base64Url = token.split(".")[1];
    if (!base64Url) return null;
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map(function (c) {
          return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
        })
        .join("")
    );
    const parsed = JSON.parse(jsonPayload);
    return parsed.user_id || parsed.sub || null;
  } catch (e) {
    console.error("Failed to decode token", e);
    return null;
  }
}

const emptySubscribe = () => () => {};

export default function KaderVerificationPage() {
  const [reports, setReports] = useState<KaderReportItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState<string | null>(null);
  const [error, setError] = useState("");
  const [successMsg, setSuccessMsg] = useState("");
  const currentUserId = useSyncExternalStore(emptySubscribe, getAuthUserId, () => null);
  
  // State for PE Form Modal
  const [showPeModal, setShowPeModal] = useState<string | null>(null);
  const [peForm, setPeForm] = useState({ pe_radius: 100, larva_found: false, fever_count: 0, psn_edukasi_done: false, notes: "", evidence_photo_url: "" });
  const [isUploadingPhoto, setIsUploadingPhoto] = useState(false);
  const [showReturnModal, setShowReturnModal] = useState<string | null>(null);
  const [returnReason, setReturnReason] = useState("");

  const fetchPendingReports = async () => {
    setLoading(true);
    try {
      const response = await api.get("/reports");
      const allReports = (response.data?.data || []) as KaderReportItem[];
      const pending = allReports.filter((r) => r.status === "MENUNGGU_PE_KADER");
      setReports(pending);

    } catch (err: unknown) {
      console.error("Failed to fetch reports", err);
      setError("Gagal mengambil data laporan. Coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ignore = false;
    api.get("/reports")
      .then((response) => {
        if (!ignore) {
          const allReports = (response.data?.data || []) as KaderReportItem[];
          const pending = allReports.filter((r) => r.status === "MENUNGGU_PE_KADER");
          setReports(pending);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          console.error("Failed to fetch reports", err);
          setError("Gagal mengambil data laporan. Coba lagi.");
        }
      })
      .finally(() => {
        if (!ignore) {
          setLoading(false);
        }
      });

    return () => {
      ignore = true;
    };
  }, []);

  const handleRejectHoax = async (id: string) => {
    const reason = prompt("Masukkan alasan penolakan (misal: Alamat fiktif, Tidak ada kasus):");
    if (!reason) return;
    
    setSubmitting(id);
    setError("");
    setSuccessMsg("");
    
    try {
      await api.post(`/reports/${id}/verify`, {
        decision: "REJECTED",
        reason: reason,
      });
      setSuccessMsg("Laporan berhasil ditandai sebagai Hoax.");
      fetchPendingReports();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || "Gagal memverifikasi laporan";
      setError(msg);
    } finally {
      setSubmitting(null);
    }
  };

  const handleSubmitPE = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!showPeModal) return;
    
    setSubmitting(showPeModal);
    setError("");
    setSuccessMsg("");
    
    try {
      await api.post(`/kader/reports/${showPeModal}/pe-results`, peForm);
      setSuccessMsg("Hasil Penyelidikan Epidemiologi berhasil disubmit.");
      setShowPeModal(null);
      // Reset form
      setPeForm({ pe_radius: 100, larva_found: false, fever_count: 0, psn_edukasi_done: false, notes: "", evidence_photo_url: "" });
      fetchPendingReports();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || "Gagal submit hasil PE";
      setError(msg);
    } finally {
      setSubmitting(null);
    }
  };

  const handlePhotoUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsUploadingPhoto(true);
    const formData = new FormData();
    formData.append("photo", file);

    try {
      const response = await api.post("/upload/photo", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      setPeForm({ ...peForm, evidence_photo_url: response.data.data.url });
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || (err as Error)?.message || "Gagal mengunggah foto";
      alert("Gagal mengunggah foto: " + msg);
    } finally {
      setIsUploadingPhoto(false);
    }
  };

  const handleReturnTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!showReturnModal || !returnReason) return;

    setSubmitting(showReturnModal);
    setError("");
    setSuccessMsg("");

    try {
      await api.post(`/kader/reports/${showReturnModal}/return`, { reason: returnReason });
      setSuccessMsg("Tugas laporan berhasil dikembalikan.");
      setShowReturnModal(null);
      setReturnReason("");
      fetchPendingReports();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || "Gagal mengembalikan tugas laporan";
      setError(msg);
    } finally {
      setSubmitting(null);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-extrabold text-text tracking-tight mb-2">Tugas Penyelidikan Epidemiologi (PE)</h2>
        <p className="text-text-muted">Lakukan kunjungan lapangan, isi hasil PE, atau tolak jika laporan terbukti palsu.</p>
      </div>

      {error && (
        <div className="bg-red-50 text-red-700 p-4 rounded-xl text-sm border border-red-100">
          {error}
        </div>
      )}

      {successMsg && (
        <div className="bg-primary-tint text-primary p-4 rounded-xl text-sm border border-emerald-100">
          {successMsg}
        </div>
      )}

      {loading ? (
        <div className="flex justify-center p-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : reports.length === 0 ? (
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-12 text-center">
          <CheckCircle2 className="h-16 w-16 text-emerald-300 mx-auto mb-4" />
          <h3 className="text-xl font-bold text-text mb-2">Semua Bersih!</h3>
          <p className="text-text-muted">Tidak ada laporan baru yang menunggu PE Anda saat ini.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
          {reports.map((report) => (
            <div key={report.report_id} className="bg-surface rounded-2xl shadow-sm border border-border overflow-hidden flex flex-col">
              <div className="p-6 flex-1">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex gap-2">
                    <span className="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-bold bg-status-verified/20 text-status-verified">
                      Menunggu PE Kader
                    </span>
                    {(() => {
                      const startDate = new Date(report.investigator_id ? report.updated_at : report.created_at);
                      const diffHours = (new Date().getTime() - startDate.getTime()) / (1000 * 60 * 60);
                      const timeLeft = 24 - diffHours;
                      
                      let color = "bg-primary-tint text-primary";
                      let text = `Sisa ${Math.max(0, Math.floor(timeLeft))} Jam`;
                      if (timeLeft <= 0) {
                        color = "bg-red-100 text-red-800 border border-red-200 animate-pulse";
                        text = "SLA Terlewat!";
                      } else if (timeLeft <= 12) {
                        color = "bg-status-wait/20 text-status-wait";
                      }
                      
                      return (
                        <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold ${color}`}>
                          <Clock className="w-3 h-3" />
                          {text}
                        </span>
                      );
                    })()}
                  </div>
                  <span className="text-xs font-mono text-text-muted">{report.reference_no}</span>
                </div>

                {report.is_possible_duplicate && (
                  <div className="mb-3 bg-amber-50 border border-amber-200 text-amber-800 px-3 py-2 rounded-lg text-xs font-medium flex items-start gap-2">
                    <AlertTriangle className="w-4 h-4 text-amber-600 flex-shrink-0" />
                    <span>Laporan ini kemungkinan duplikat dari kasus di sekitar (radius 100m). Harap pastikan keabsahannya saat PE.</span>
                  </div>
                )}
                
                <h3 className="text-lg font-bold text-text mb-2">{report.category}</h3>
                <p className="text-text-muted text-sm mb-4 line-clamp-3">{report.description}</p>
                
                <div className="space-y-2 mt-auto">
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <MapPin className="h-4 w-4 text-text-muted" />
                    <span className="truncate">{report.location?.address || 'Lokasi tidak spesifik'}</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <User className="h-4 w-4 text-text-muted" />
                    <span>{report.reporter?.anonymous_flag ? "Anonim" : report.reporter?.name}</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <Calendar className="h-4 w-4 text-text-muted" />
                    <span>{new Date(report.created_at).toLocaleString('id-ID')}</span>
                  </div>
                </div>
              </div>
              
              <div className="border-t border-border bg-bg p-4 flex gap-3">
                {!report.investigator_id ? (
                  <button
                    onClick={async () => {
                      setSubmitting(report.report_id);
                      setError("");
                      try {
                        await api.post(`/kader/reports/${report.report_id}/claim`);
                        setSuccessMsg("Berhasil mengambil alih laporan.");
                        fetchPendingReports();
                      } catch (err: unknown) {
                        const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || "Gagal mengambil alih laporan";
                        setError(msg);
                      } finally {
                        setSubmitting(null);
                      }
                    }}
                    disabled={submitting === report.report_id}
                    className="flex-1 py-2.5 px-4 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-colors flex justify-center items-center gap-2 disabled:opacity-50"
                  >
                    {submitting === report.report_id ? <Loader2 className="w-4 h-4 animate-spin" /> : <User className="w-4 h-4" />}
                    Ambil Tugas
                  </button>
                ) : report.investigator_id !== currentUserId ? (
                  <div className="flex-1 text-center py-2.5 text-text-muted font-medium text-sm">
                    Sedang diproses Kader lain
                  </div>
                ) : (
                  <>
                    <button
                      onClick={() => handleRejectHoax(report.report_id)}
                      disabled={submitting === report.report_id}
                      className="flex-1 py-2.5 px-4 bg-surface border border-red-200 text-red-600 font-bold rounded-xl hover:bg-red-50 transition-colors flex justify-center items-center gap-2 disabled:opacity-50"
                    >
                      {submitting === report.report_id ? <Loader2 className="w-4 h-4 animate-spin" /> : <XCircle className="w-4 h-4" />}
                      Laporan Hoax
                    </button>
                    <button
                      onClick={() => setShowReturnModal(report.report_id)}
                      disabled={submitting === report.report_id}
                      className="flex-1 py-2.5 px-4 bg-surface border border-border text-text font-bold rounded-xl hover:bg-bg transition-colors flex justify-center items-center gap-2 disabled:opacity-50"
                    >
                      Kembalikan
                    </button>
                    <button
                      onClick={() => setShowPeModal(report.report_id)}
                      disabled={submitting === report.report_id}
                      className="flex-1 py-2.5 px-4 bg-primary text-primary-foreground font-bold rounded-xl hover:bg-primary-deep transition-colors flex justify-center items-center gap-2 disabled:opacity-50"
                    >
                      {submitting === report.report_id ? <Loader2 className="w-4 h-4 animate-spin" /> : <CheckCircle2 className="w-4 h-4" />}
                      Input Form PE
                    </button>
                  </>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* PE Form Modal */}
      {showPeModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-primary-deep/50 backdrop-blur-sm animate-in fade-in">
          <div className="bg-surface rounded-2xl shadow-xl w-full max-w-lg overflow-hidden max-h-[90vh] flex flex-col">
            <div className="p-6 border-b border-border flex justify-between items-center bg-bg/50">
              <div>
                <h3 className="text-xl font-bold text-text">Form Hasil PE</h3>
                <p className="text-sm text-text-muted mt-1">Masukkan hasil pengecekan di lokasi.</p>
              </div>
              <button 
                type="button"
                onClick={() => setShowPeModal(null)}
                className="p-2 hover:bg-border rounded-full transition-colors text-text-muted hover:text-text"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            <form id="kader-pe-form" onSubmit={handleSubmitPE} className="p-6 space-y-4 overflow-y-auto flex-1">
              <div>
                <label className="block text-sm font-semibold text-text mb-1">Radius Pemeriksaan (Meter)</label>
                <input 
                  type="number" 
                  value={peForm.pe_radius}
                  onChange={(e) => setPeForm({...peForm, pe_radius: parseInt(e.target.value) || 0})}
                  className="w-full border-border rounded-xl p-3 focus:ring-primary focus:border-primary"
                  required
                  min="0"
                />
                <p className="text-xs text-text-muted mt-1">Default SOP adalah 100 meter.</p>
              </div>
              
              <div className="flex items-center justify-between p-4 border border-border rounded-xl bg-bg">
                <div className="flex items-center gap-3">
                  <Bug className="text-primary w-5 h-5" />
                  <div>
                    <p className="font-semibold text-text text-sm">Ditemukan Jentik?</p>
                    <p className="text-xs text-text-muted">Di dalam/luar rumah</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" className="sr-only peer" checked={peForm.larva_found} onChange={(e) => setPeForm({...peForm, larva_found: e.target.checked})} />
                  <div className="w-11 h-6 bg-border peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-surface after:border-border after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                </label>
              </div>

              <div>
                <label className="block text-sm font-semibold text-text mb-1 flex items-center gap-2">
                  <Thermometer className="w-4 h-4 text-primary" />
                  Jumlah Warga Demam Tambahan
                </label>
                <input 
                  type="number" 
                  value={peForm.fever_count}
                  onChange={(e) => setPeForm({...peForm, fever_count: parseInt(e.target.value) || 0})}
                  className="w-full border-border rounded-xl p-3 focus:ring-primary focus:border-primary"
                  required
                  min="0"
                />
                <p className="text-xs text-text-muted mt-1">Isi 0 jika tidak ada tetangga sekitar yang demam tanpa sebab jelas.</p>
              </div>

              <div className="flex items-center justify-between p-4 border border-border rounded-xl bg-status-verified/20">
                <div className="flex items-center gap-3">
                  <div>
                    <p className="font-semibold text-text text-sm">Edukasi PSN 3M Plus Selesai?</p>
                    <p className="text-xs text-text-muted">Wajib dilakukan pada warga setempat</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" className="sr-only peer" checked={peForm.psn_edukasi_done} onChange={(e) => setPeForm({...peForm, psn_edukasi_done: e.target.checked})} required />
                  <div className="w-11 h-6 bg-border peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-surface after:border-border after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                </label>
              </div>

              <div>
                <label className="block text-sm font-semibold text-text mb-1">Catatan/Narasi Hasil PE</label>
                <textarea 
                  value={peForm.notes}
                  onChange={(e) => setPeForm({...peForm, notes: e.target.value})}
                  className="w-full border-border rounded-xl p-3 focus:ring-primary focus:border-primary"
                  required
                  rows={3}
                  placeholder="Ceritakan detail penemuan di lokasi, kondisi lingkungan, dll."
                ></textarea>
              </div>

              <div>
                <label className="block text-sm font-semibold text-text mb-1">Bukti Foto</label>
                <div className="flex gap-2 items-center">
                  <input 
                    type="file" 
                    accept="image/*"
                    onChange={handlePhotoUpload}
                    disabled={isUploadingPhoto}
                    className="flex-1 border-border rounded-xl p-3 focus:ring-primary focus:border-primary text-sm file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary-tint file:text-primary hover:file:bg-primary hover:file:text-white transition-all cursor-pointer"
                  />
                  {isUploadingPhoto && <Loader2 className="w-5 h-5 animate-spin text-primary" />}
                </div>
                {peForm.evidence_photo_url && (
                  <div className="mt-3 flex items-start gap-3 p-2 bg-bg border border-border rounded-xl w-fit">
                    <div className="rounded-lg overflow-hidden border border-border w-28 h-28 relative">
                      {/* eslint-disable-next-line @next/next/no-img-element */}
                      <img src={peForm.evidence_photo_url.startsWith('http') ? peForm.evidence_photo_url : `${api.defaults.baseURL?.replace('/api/v1', '')}${peForm.evidence_photo_url}`} alt="Preview Bukti" className="object-cover w-full h-full" />
                    </div>
                    <button
                      type="button"
                      onClick={() => setPeForm({...peForm, evidence_photo_url: ""})}
                      className="p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors flex items-center gap-1 text-xs font-semibold"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                      Hapus Foto
                    </button>
                  </div>
                )}
              </div>
            </form>

            <div className="p-4 border-t border-border bg-bg/50 flex gap-3 mt-auto">
              <button 
                type="button"
                onClick={() => setShowPeModal(null)}
                className="flex-1 px-4 py-3 border border-border text-text-muted font-bold rounded-xl hover:bg-bg transition-colors"
              >
                Batal
              </button>
              <button 
                type="submit"
                form="kader-pe-form"
                disabled={submitting === showPeModal || isUploadingPhoto}
                className="flex-1 px-4 py-3 bg-primary text-primary-foreground font-bold rounded-xl hover:bg-primary-deep transition-colors flex justify-center items-center gap-2 disabled:opacity-50"
              >
                {submitting === showPeModal ? <Loader2 className="w-5 h-5 animate-spin" /> : "Submit PE"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Return Task Modal */}
      {showReturnModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-primary-deep/50 backdrop-blur-sm animate-in fade-in">
          <div className="bg-surface rounded-2xl shadow-xl w-full max-w-md overflow-hidden max-h-[90vh] flex flex-col">
            <div className="p-6 border-b border-border flex justify-between items-center bg-bg/50">
              <div>
                <h3 className="text-xl font-bold text-text">Kembalikan Tugas</h3>
                <p className="text-sm text-text-muted mt-1">Tugas yang dikembalikan akan masuk kembali ke daftar antrean laporan yang bisa diklaim Kader lain.</p>
              </div>
              <button 
                type="button"
                onClick={() => {
                  setShowReturnModal(null);
                  setReturnReason("");
                }}
                className="p-2 hover:bg-border rounded-full transition-colors text-text-muted hover:text-text"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            <form id="return-task-form" onSubmit={handleReturnTask} className="p-6 space-y-4 overflow-y-auto flex-1">
              <div>
                <label className="block text-sm font-semibold text-text mb-1">Alasan Pengembalian</label>
                <textarea 
                  value={returnReason}
                  onChange={(e) => setReturnReason(e.target.value)}
                  className="w-full border-border rounded-xl p-3 focus:ring-blue-500 focus:border-blue-500"
                  required
                  rows={3}
                  placeholder="Misal: Bukan wilayah tugas saya, sedang berhalangan hadir..."
                ></textarea>
              </div>
            </form>

            <div className="p-4 border-t border-border bg-bg/50 flex gap-3 mt-auto">
              <button 
                type="button"
                onClick={() => {
                  setShowReturnModal(null);
                  setReturnReason("");
                }}
                className="flex-1 px-4 py-3 border border-border text-text-muted font-bold rounded-xl hover:bg-bg transition-colors"
              >
                Batal
              </button>
              <button 
                type="submit"
                form="return-task-form"
                disabled={submitting === showReturnModal || !returnReason.trim()}
                className="flex-1 px-4 py-3 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-colors flex justify-center items-center gap-2 disabled:opacity-50"
              >
                {submitting === showReturnModal ? <Loader2 className="w-5 h-5 animate-spin" /> : "Kembalikan"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
