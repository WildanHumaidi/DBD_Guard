"use client";

import { useEffect, useState } from "react";
import api from "@/lib/api";
import { CheckCircle2, Loader2, MapPin, Calendar, Activity, X, PackageX, Clock, HelpCircle, Info } from "lucide-react";

export default function PuskesmasPEValidationPage() {
  const [reports, setReports] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [successMsg, setSuccessMsg] = useState("");
  const [inventoryItems, setInventoryItems] = useState<any[]>([]);
  
  // Modal State
  const [selectedReport, setSelectedReport] = useState<any>(null);
  const [formData, setFormData] = useState({
    is_approved: true,
    reason: "",
    is_override: false,
    override_decision: "FOGGING_FOKUS",
    override_reason: ""
  });

  // PE Form State (untuk ambil alih / input langsung)
  const [showPeModal, setShowPeModal] = useState<string | null>(null);
  const [peForm, setPeForm] = useState({
    pe_radius: 100,
    larva_found: false,
    fever_count: 0,
    psn_edukasi_done: false,
    notes: "",
    evidence_photo_url: ""
  });
  const [isUploadingPhoto, setIsUploadingPhoto] = useState(false);

  const fetchPendingReports = async () => {
    setLoading(true);
    try {
      const response = await api.get("/reports");
      const allReports = response.data?.data || [];
      // Reports that are MENUNGGU_VALIDASI_PUSKESMAS need PE validation
      // Reports that are MENUNGGU_PE_KADER or ESKALASI_PUSKESMAS can be taken over by Puskesmas directly
      const pending = allReports.filter((r: any) => 
        r.status === "MENUNGGU_VALIDASI_PUSKESMAS" || 
        r.status === "MENUNGGU_PE_KADER" ||
        r.status === "ESKALASI_PUSKESMAS"
      );
      setReports(pending);

      try {
        const invRes = await api.get("/inventory");
        setInventoryItems(invRes.data?.data || []);
      } catch (e) {
        console.error("Failed to fetch inventory", e);
      }
    } catch (err) {
      console.error("Failed to fetch reports", err);
      setError("Gagal mengambil data laporan. Coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPendingReports();
  }, []);

  const handleSubmitValidation = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedReport || !selectedReport.pe_result?.pe_id) return;

    setSubmitting(true);
    setError("");
    setSuccessMsg("");
    
    try {
      const payload: any = {
        is_approved: formData.is_approved,
        reason: formData.reason,
      };

      if (formData.is_approved && formData.is_override) {
        if (!formData.override_reason.trim()) {
          setError("Alasan override wajib diisi jika mengaktifkan override rekomendasi DSS");
          setSubmitting(false);
          return;
        }
        payload.override_decision = formData.override_decision;
        payload.override_reason = formData.override_reason;
      }

      await api.post(`/puskesmas/pe-results/${selectedReport.pe_result.pe_id}/validate`, payload);
      setSuccessMsg(`Validasi PE untuk ${selectedReport.reference_no} berhasil disimpan.`);
      setSelectedReport(null); // Close modal
      setFormData({ is_approved: true, reason: "", is_override: false, override_decision: "FOGGING_FOKUS", override_reason: "" });
      fetchPendingReports();
    } catch (err: any) {
      setError(err.response?.data?.error || "Gagal menyimpan validasi PE");
    } finally {
      setSubmitting(false);
    }
  };

  const handleSubmitPE = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!showPeModal) return;
    
    setSubmitting(true);
    setError("");
    setSuccessMsg("");
    
    try {
      // Formal takeover workflow: assigns investigator_id to Puskesmas officer and creates audit trail
      await api.post(`/puskesmas/reports/${showPeModal}/takeover`, {
        reason: "Puskesmas direct PE investigation takeover"
      });
      await api.post(`/puskesmas/reports/${showPeModal}/pe-results`, peForm);
      setSuccessMsg("Hasil Penyelidikan Epidemiologi berhasil diambil alih dan disubmit oleh Puskesmas.");
      setShowPeModal(null);
      setPeForm({ pe_radius: 100, larva_found: false, fever_count: 0, psn_edukasi_done: false, notes: "", evidence_photo_url: "" });
      fetchPendingReports();
    } catch (err: any) {
      setError(err.response?.data?.error || "Gagal submit hasil PE");
    } finally {
      setSubmitting(false);
    }
  };

  const handlePhotoUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsUploadingPhoto(true);
    const formDataUpload = new FormData();
    formDataUpload.append("photo", file);

    try {
      const response = await api.post("/upload/photo", formDataUpload, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      setPeForm({ ...peForm, evidence_photo_url: response.data.data.url });
    } catch (err: any) {
      alert("Gagal mengunggah foto: " + (err.response?.data?.error || err.message));
    } finally {
      setIsUploadingPhoto(false);
    }
  };

  const handleForceUnassign = async (id: string) => {
    if (!confirm("Apakah Anda yakin ingin mencabut paksa tugas laporan ini dari kader yang bersangkutan?")) return;

    setSubmitting(true);
    setError("");
    setSuccessMsg("");
    
    try {
      await api.post(`/puskesmas/reports/${id}/unassign`);
      setSuccessMsg("Tugas laporan berhasil dicabut secara paksa.");
      fetchPendingReports();
    } catch (err: any) {
      setError(err.response?.data?.error || "Gagal mencabut tugas laporan");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Simulation Banner */}
      <div className="bg-amber-500/10 border border-amber-500/30 rounded-2xl p-4 flex items-center justify-between gap-3 text-amber-900 dark:text-amber-200">
        <div className="flex items-center gap-3">
          <span className="px-2.5 py-1 text-xs font-black uppercase tracking-wider bg-amber-500 text-white rounded-md shadow-sm">
            DATA SIMULASI
          </span>
          <p className="text-xs font-semibold">
            Modul Sistem Pendukung Keputusan (SPK) Berbasis Aturan &bull; Penanganan Kasus DBD Kota Depok
          </p>
        </div>
        <span className="text-[11px] font-bold text-amber-800 dark:text-amber-300">
          Rule-Based DSS v1.0
        </span>
      </div>

      <div>
        <h2 className="text-3xl font-extrabold text-text tracking-tight mb-2">Validasi PE Kader & SPK</h2>
        <p className="text-text-muted">Evaluasi hasil Penyelidikan Epidemiologi (PE) kader dan telaah rekomendasi otomatis berbasis bukti lapangan.</p>
      </div>

      {error && !selectedReport && (
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
          <Loader2 className="h-8 w-8 animate-spin text-status-verified" />
        </div>
      ) : reports.length === 0 ? (
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-12 text-center">
          <CheckCircle2 className="h-16 w-16 text-blue-300 mx-auto mb-4" />
          <h3 className="text-xl font-bold text-text mb-2">Tidak Ada Tugas Validasi</h3>
          <p className="text-text-muted">Semua hasil PE kader telah divalidasi.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
          {reports.map((report) => (
            <div key={report.report_id} className="bg-surface rounded-2xl shadow-sm border border-border overflow-hidden flex flex-col">
              <div className="p-6 flex-1">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex gap-2">
                    <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-bold ${
                      report.status === 'MENUNGGU_PE_KADER' ? 'bg-indigo-100 text-indigo-800' :
                      report.status === 'ESKALASI_PUSKESMAS' ? 'bg-amber-100 text-amber-800 border border-amber-300' :
                      'bg-status-wait/20 text-status-wait'
                    }`}>
                      {report.status === 'MENUNGGU_PE_KADER' ? 'Kader Sedang PE' :
                       report.status === 'ESKALASI_PUSKESMAS' ? 'Eskalasi Puskesmas (SLA)' :
                       'Menunggu Validasi PE'}
                    </span>
                    {report.status === 'MENUNGGU_VALIDASI_PUSKESMAS' && (() => {
                      const startDate = new Date(report.created_at);
                      const diffHours = (new Date().getTime() - startDate.getTime()) / (1000 * 60 * 60);
                      const timeLeft = 72 - diffHours;
                      
                      let color = "bg-primary-tint text-primary";
                      let text = `Sisa SLA ${Math.max(0, Math.floor(timeLeft))} Jam`;
                      if (timeLeft <= 0) {
                        color = "bg-red-100 text-red-800 border border-red-200 animate-pulse";
                        text = "SLA Fogging Terlewat!";
                      } else if (timeLeft <= 24) {
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
                
                <h3 className="text-lg font-bold text-text mb-2">{report.category}</h3>
                <p className="text-text-muted text-sm mb-4 line-clamp-3">{report.description}</p>

                <div className="bg-status-verified/20 border border-blue-100 rounded-lg p-3 mb-4 text-sm">
                  <p className="font-semibold text-blue-900 mb-1">
                    {report.status === 'MENUNGGU_PE_KADER' ? 'Informasi Tugas:' : 'Hasil PE Kader:'}
                  </p>
                  {report.status === 'MENUNGGU_PE_KADER' ? (
                    <p className="text-status-verified text-xs">Sedang dikerjakan oleh Kader (ID: {report.investigator_id?.substring(0, 8)}...)</p>
                  ) : (
                    <ul className="list-disc list-inside text-status-verified space-y-1 text-xs">
                      <li>Radius: {report.pe_result?.pe_radius || '-'} meter</li>
                      <li>Jentik Ditemukan: {report.pe_result?.larva_found ? 'Ya' : 'Tidak'}</li>
                      <li>Kasus Demam: {report.pe_result?.fever_count || 0} orang</li>
                    </ul>
                  )}
                </div>
                
                <div className="space-y-2 mt-auto">
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <MapPin className="h-4 w-4 text-text-muted" />
                    <span className="truncate">{report.location?.address || 'Lokasi tidak spesifik'}</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <Calendar className="h-4 w-4 text-text-muted" />
                    <span>{new Date(report.created_at).toLocaleString('id-ID')}</span>
                  </div>
                </div>
              </div>
              
              <div className="border-t border-border bg-bg p-4">
                {report.status === 'MENUNGGU_PE_KADER' || report.status === 'ESKALASI_PUSKESMAS' ? (
                  <div className="flex gap-2">
                    <button
                      onClick={() => setShowPeModal(report.report_id)}
                      disabled={submitting}
                      className="flex-1 py-2.5 px-4 bg-primary text-primary-foreground font-bold rounded-xl hover:bg-primary-deep transition-colors flex justify-center items-center gap-2 disabled:opacity-50"
                    >
                      <Activity className="w-4 h-4" />
                      Input Hasil PE (Ambil Alih)
                    </button>
                    {report.investigator_id && (
                      <button
                        onClick={() => handleForceUnassign(report.report_id)}
                        disabled={submitting}
                        className="py-2.5 px-4 bg-surface border border-red-200 text-red-600 font-bold rounded-xl hover:bg-red-50 transition-colors flex justify-center items-center disabled:opacity-50"
                        title="Cabut Paksa Tugas dari Kader"
                      >
                        <X className="w-5 h-5" />
                      </button>
                    )}
                  </div>
                ) : (
                  <button
                    onClick={() => { setSelectedReport(report); setError(""); }}
                    className="w-full py-2.5 px-4 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-colors flex justify-center items-center gap-2"
                  >
                    <Activity className="w-4 h-4" />
                    Validasi PE
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Validation Modal */}
      {selectedReport && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-primary-deep/50 backdrop-blur-sm animate-fade-in">
          <div className="bg-surface rounded-2xl shadow-xl w-full max-w-lg overflow-hidden flex flex-col max-h-[90vh]">
            <div className="p-6 border-b border-border flex justify-between items-center bg-bg">
              <h3 className="font-bold text-lg text-text">Validasi Hasil PE</h3>
              <button 
                onClick={() => setSelectedReport(null)}
                className="text-text-muted hover:text-text-muted p-1"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            <div className="p-6 overflow-y-auto">
              <div className="mb-6 pb-6 border-b border-border">
                <p className="text-xs text-text-muted font-mono mb-1">{selectedReport.reference_no}</p>
                <p className="text-sm font-semibold text-text">{selectedReport.location?.address}</p>
              </div>

              <div className="bg-status-verified/20 border border-blue-100 rounded-lg p-4 mb-6">
                <h4 className="font-bold text-blue-900 mb-2">Data PE dari Kader</h4>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="block text-text-muted text-xs">Radius</span>
                    <span className="font-medium text-text">{selectedReport.pe_result?.pe_radius || '-'} meter</span>
                  </div>
                  <div>
                    <span className="block text-text-muted text-xs">Jentik</span>
                    <span className="font-medium text-text">{selectedReport.pe_result?.larva_found ? 'Ditemukan' : 'Tidak Ditemukan'}</span>
                  </div>
                  <div>
                    <span className="block text-text-muted text-xs">Penderita Demam</span>
                    <span className="font-medium text-text">{selectedReport.pe_result?.fever_count || 0} orang</span>
                  </div>
                </div>
              </div>

              {(selectedReport.pe_result?.notes || selectedReport.pe_result?.evidence_photo_url) && (
                <div className="bg-bg border border-border rounded-lg p-4 mb-6">
                  <h4 className="font-bold text-text mb-2">Narasi & Bukti Kunjungan</h4>
                  {selectedReport.pe_result?.notes && (
                    <div className="mb-3">
                      <p className="text-xs text-text-muted font-semibold mb-1">Catatan Kader:</p>
                      <p className="text-sm text-text bg-surface p-3 rounded-lg border border-border">
                        {selectedReport.pe_result.notes}
                      </p>
                    </div>
                  )}
                  {selectedReport.pe_result?.evidence_photo_url && (
                    <div>
                      <p className="text-xs text-text-muted font-semibold mb-1">Foto Bukti:</p>
                      <div className="rounded-lg overflow-hidden border border-border">
                        <img 
                          src={selectedReport.pe_result.evidence_photo_url.startsWith('http') 
                            ? selectedReport.pe_result.evidence_photo_url 
                            : `${api.defaults.baseURL?.replace('/api/v1', '')}${selectedReport.pe_result.evidence_photo_url}${typeof window !== 'undefined' && localStorage.getItem('token') ? `?token=${localStorage.getItem('token')}` : ''}`} 
                          alt="Bukti Kunjungan Kader" 
                          className="w-full h-auto max-h-64 object-contain bg-bg"
                        />
                      </div>
                    </div>
                  )}
                </div>
              )}

              {(() => {
                const insektisida = inventoryItems.find((item: any) => item.item_type === "INSEKTISIDA");
                if (insektisida && insektisida.quantity < 10) {
                  return (
                    <div className="mb-6 bg-red-50 border-l-4 border-red-500 p-4 rounded-r-lg flex items-start gap-3">
                      <PackageX className="h-5 w-5 text-red-600 flex-shrink-0" />
                      <div>
                        <h4 className="text-red-800 font-bold text-sm">Peringatan: Stok Insektisida Menipis!</h4>
                        <p className="text-red-700 text-xs mt-1">Sisa stok saat ini hanya <strong>{insektisida.quantity} Liter</strong>. Fogging mungkin tidak dapat dieksekusi atau otomatis ditolak oleh sistem backend. Pertimbangkan untuk meminjam dari gudang Dinkes atau alihkan keputusan validasi.</p>
                      </div>
                    </div>
                  );
                }
                return null;
              })()}

              {error && (
                <div className="mb-4 bg-red-50 text-red-700 p-3 rounded-lg text-sm border border-red-100">
                  {error}
                </div>
              )}

              <form id="validation-form" onSubmit={handleSubmitValidation} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-text mb-2">Keputusan Validasi *</label>
                  <div className="flex gap-4">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input 
                        type="radio" 
                        name="is_approved" 
                        checked={formData.is_approved === true}
                        onChange={() => setFormData({...formData, is_approved: true})}
                        className="text-status-verified focus:ring-blue-500" 
                      />
                      <span className="text-sm font-medium text-text">Setujui (Valid)</span>
                    </label>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input 
                        type="radio" 
                        name="is_approved" 
                        checked={formData.is_approved === false}
                        onChange={() => setFormData({...formData, is_approved: false})}
                        className="text-red-600 focus:ring-red-500" 
                      />
                      <span className="text-sm font-medium text-text">Tolak (Perlu Revisi)</span>
                    </label>
                  </div>
                </div>

                {formData.is_approved && (() => {
                  const larvaeFound = Boolean(selectedReport.pe_result?.larva_found);
                  const feverCount = Number(selectedReport.pe_result?.fever_count || 0);
                  const insektisida = inventoryItems.find((item: any) => item.item_type === "INSEKTISIDA");
                  const currentStock = insektisida ? Number(insektisida.quantity) : 0;
                  
                  const isHighRisk = larvaeFound || feverCount >= 3;
                  const hasSufficientStock = currentStock >= 10;

                  // Active Rule Logic
                  let activeRuleNumber = 3;
                  let activeRecommendation = "EDUKASI_PSN";
                  let ruleExplanation = "";

                  if (isHighRisk && hasSufficientStock) {
                    activeRuleNumber = 1;
                    activeRecommendation = "FOGGING_FOKUS";
                    ruleExplanation = `Terdeteksi indikator risiko tinggi (Jentik: ${larvaeFound ? "Positif" : "Negatif"}, Demam: ${feverCount} orang [ambang batas \u2265 3]) dan ketersediaan stok insektisida memadai (${currentStock} L \u2265 10 L). Sistem merekomendasikan Fogging Fokus segera.`;
                  } else if (isHighRisk && !hasSufficientStock) {
                    activeRuleNumber = 2;
                    activeRecommendation = "FOGGING_TERTUNDA";
                    ruleExplanation = `Kriteria fogging terpenuhi (faktor risiko tinggi), namun stok insektisida di gudang puskesmas kurang dari ambang batas minimum (${currentStock} L < 10 L). Sistem merekomendasikan penundaan fogging dan pembuatan auto-draft permintaan stok darurat ke Dinas Kesehatan.`;
                  } else {
                    activeRuleNumber = 3;
                    activeRecommendation = "EDUKASI_PSN";
                    ruleExplanation = `Tidak ditemukan jentik dan jumlah demam (${feverCount} orang) di bawah ambang batas (3 orang). Berdasarkan pedoman Kemenkes RI, fogging tidak diperbolehkan tanpa indikasi kuat guna menghindari resistensi nyamuk. Rekomendasi: Edukasi 3M Plus dan Pemberantasan Sarang Nyamuk (PSN).`;
                  }

                  return (
                    <div className="space-y-4 p-4 bg-bg border border-border rounded-xl">
                      {/* SPK Header */}
                      <div className="flex items-center justify-between pb-3 border-b border-border">
                        <div className="flex items-center gap-2">
                          <div className="p-1.5 bg-primary/10 text-primary rounded-lg">
                            <Activity className="w-4 h-4" />
                          </div>
                          <div>
                            <span className="text-xs font-bold text-text block">Rekomendasi SPK Berbasis Aturan (Rule-Based DSS)</span>
                            <span className="text-[10px] text-text-muted">Standar Operasional Penanganan DBD Kota Depok</span>
                          </div>
                        </div>
                        <span className={`px-2.5 py-1 rounded-md text-xs font-black tracking-wide ${
                          activeRecommendation === 'FOGGING_FOKUS'
                            ? 'bg-red-100 text-red-800 border border-red-200'
                            : activeRecommendation === 'FOGGING_TERTUNDA'
                            ? 'bg-amber-100 text-amber-800 border border-amber-200'
                            : 'bg-emerald-100 text-emerald-800 border border-emerald-200'
                        }`}>
                          {activeRecommendation === 'FOGGING_FOKUS' ? '🔴 FOGGING FOKUS' :
                           activeRecommendation === 'FOGGING_TERTUNDA' ? '🟡 FOGGING TERTUNDA' :
                           '🟢 EDUKASI & PSN'}
                        </span>
                      </div>

                      {/* Transparansi: Panel Mengapa Rekomendasi Ini? */}
                      <div className="bg-surface rounded-xl border border-indigo-100 dark:border-indigo-900/40 p-3.5 space-y-3">
                        <div className="flex items-center justify-between">
                          <span className="text-xs font-bold text-indigo-900 dark:text-indigo-300 flex items-center gap-1.5">
                            <HelpCircle className="w-3.5 h-3.5 text-indigo-600" />
                            Mengapa rekomendasi ini?
                          </span>
                          <span className="text-[10px] font-semibold text-indigo-600 bg-indigo-50 dark:bg-indigo-950/60 px-2 py-0.5 rounded">
                            Aturan #{activeRuleNumber} Terpicu
                          </span>
                        </div>

                        {/* Input Variables Evaluation */}
                        <div className="grid grid-cols-3 gap-2 text-[11px]">
                          {/* Var 1: Jentik */}
                          <div className={`p-2 rounded-lg border ${
                            larvaeFound ? 'bg-red-50/50 border-red-200 text-red-900' : 'bg-bg border-border text-text'
                          }`}>
                            <span className="text-[10px] text-text-muted block">1. Jentik (100m)</span>
                            <span className={`font-bold block mt-0.5 ${larvaeFound ? 'text-red-600' : 'text-emerald-600'}`}>
                              {larvaeFound ? 'Positif (Ya)' : 'Negatif'}
                            </span>
                          </div>

                          {/* Var 2: Demam */}
                          <div className={`p-2 rounded-lg border ${
                            feverCount >= 3 ? 'bg-red-50/50 border-red-200 text-red-900' : 'bg-bg border-border text-text'
                          }`}>
                            <span className="text-[10px] text-text-muted block">2. Demam Sekitar</span>
                            <span className={`font-bold block mt-0.5 ${feverCount >= 3 ? 'text-red-600' : 'text-text'}`}>
                              {feverCount} Orang {feverCount >= 3 ? '(\u2265 3)' : '(< 3)'}
                            </span>
                          </div>

                          {/* Var 3: Stok */}
                          <div className={`p-2 rounded-lg border ${
                            !hasSufficientStock && isHighRisk ? 'bg-amber-50/50 border-amber-200 text-amber-900' : 'bg-bg border-border text-text'
                          }`}>
                            <span className="text-[10px] text-text-muted block">3. Stok Insektisida</span>
                            <span className={`font-bold block mt-0.5 ${hasSufficientStock ? 'text-emerald-600' : 'text-amber-600'}`}>
                              {currentStock} L {hasSufficientStock ? '(\u2265 10L)' : '(< 10L)'}
                            </span>
                          </div>
                        </div>

                        {/* Active Rule Explanation Box */}
                        <div className="p-2.5 bg-indigo-50/60 dark:bg-indigo-950/30 rounded-lg border border-indigo-200/60 dark:border-indigo-900/60 text-xs text-indigo-900 dark:text-indigo-200 leading-relaxed">
                          <p className="font-semibold text-[11px] mb-1 text-indigo-800 dark:text-indigo-300">
                            Logika Keputusan Terpicu:
                          </p>
                          <p className="text-[11px]">{ruleExplanation}</p>
                        </div>

                        {/* Matriks Aturan Standar */}
                        <div className="pt-2 border-t border-border text-[10px] text-text-muted space-y-1">
                          <div className="font-semibold text-text">Pedoman Matriks Keputusan:</div>
                          <div className={`p-1.5 rounded flex items-center justify-between ${activeRuleNumber === 1 ? 'bg-red-50 text-red-900 font-bold border border-red-200' : ''}`}>
                            <span>Aturan 1: (Jentik = Ya ATAU Demam &ge; 3) &amp; Stok &ge; 10L &rarr; Fogging Fokus</span>
                            {activeRuleNumber === 1 && <span className="text-red-600 text-[9px] uppercase tracking-wider font-extrabold">&bull; AKTIF</span>}
                          </div>
                          <div className={`p-1.5 rounded flex items-center justify-between ${activeRuleNumber === 2 ? 'bg-amber-50 text-amber-900 font-bold border border-amber-200' : ''}`}>
                            <span>Aturan 2: (Jentik = Ya ATAU Demam &ge; 3) &amp; Stok &lt; 10L &rarr; Fogging Tertunda</span>
                            {activeRuleNumber === 2 && <span className="text-amber-600 text-[9px] uppercase tracking-wider font-extrabold">&bull; AKTIF</span>}
                          </div>
                          <div className={`p-1.5 rounded flex items-center justify-between ${activeRuleNumber === 3 ? 'bg-emerald-50 text-emerald-900 font-bold border border-emerald-200' : ''}`}>
                            <span>Aturan 3: Jentik = Tidak &amp; Demam &lt; 3 &rarr; Edukasi &amp; PSN 3M</span>
                            {activeRuleNumber === 3 && <span className="text-emerald-600 text-[9px] uppercase tracking-wider font-extrabold">&bull; AKTIF</span>}
                          </div>
                        </div>
                      </div>

                      {/* Override Section */}
                      <div className="pt-2 border-t border-border">
                        <label className="flex items-center gap-2 cursor-pointer">
                          <input 
                            type="checkbox"
                            checked={formData.is_override}
                            onChange={(e) => setFormData({...formData, is_override: e.target.checked})}
                            className="rounded text-primary focus:ring-primary"
                          />
                          <span className="text-xs font-bold text-text">Override Rekomendasi Sistem (Advisory Override)</span>
                        </label>

                        {formData.is_override && (
                          <div className="mt-3 space-y-3 p-3 bg-surface border border-amber-200 rounded-lg">
                            <div>
                              <label className="block text-xs font-semibold text-text mb-1">Keputusan Override Manual *</label>
                              <select
                                value={formData.override_decision}
                                onChange={(e) => setFormData({...formData, override_decision: e.target.value})}
                                className="w-full px-3 py-1.5 text-xs border border-border rounded-lg bg-bg focus:ring-2 focus:ring-primary outline-none"
                              >
                                <option value="FOGGING_FOKUS">FOGGING FOKUS (Prioritaskan Fogging)</option>
                                <option value="FOGGING_TERTUNDA">FOGGING TERTUNDA (Stok Kurang)</option>
                                <option value="EDUKASI_PSN">EDUKASI PSN (Hanya Gerakan PSN)</option>
                              </select>
                            </div>

                            <div>
                              <label className="block text-xs font-semibold text-text mb-1">Alasan Override (Wajib) *</label>
                              <textarea
                                rows={2}
                                value={formData.override_reason}
                                onChange={(e) => setFormData({...formData, override_reason: e.target.value})}
                                placeholder="Dasar pertimbangan klinis/lapangan untuk mengesampingkan rekomendasi DSS..."
                                className="w-full px-3 py-1.5 text-xs border border-border rounded-lg bg-bg focus:ring-2 focus:ring-primary outline-none resize-none"
                                required={formData.is_override}
                              ></textarea>
                            </div>
                          </div>
                        )}
                      </div>
                    </div>
                  );
                })()}

                <div>
                  <label className="block text-sm font-medium text-text mb-1">Catatan Tambahan (Opsional)</label>
                  <textarea
                    rows={3}
                    value={formData.reason}
                    onChange={(e) => setFormData({...formData, reason: e.target.value})}
                    placeholder="Contoh: Harus dire-visit karena wilayah lebih luas"
                    className="w-full px-4 py-2 border border-border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all resize-none"
                  ></textarea>
                </div>
              </form>
            </div>
            
            <div className="p-6 border-t border-border bg-bg flex gap-3">
              <button
                type="button"
                onClick={() => setSelectedReport(null)}
                className="flex-1 px-4 py-2.5 border border-border text-text font-bold rounded-xl hover:bg-bg transition-colors"
              >
                Batal
              </button>
              <button
                type="submit"
                form="validation-form"
                disabled={submitting}
                className="flex-1 px-4 py-2.5 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-colors disabled:opacity-50 flex justify-center items-center gap-2"
              >
                {submitting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                Simpan Validasi
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Modal Input PE (Puskesmas Takeover) */}
      {showPeModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fade-in">
          <div className="bg-surface w-full max-w-lg rounded-3xl shadow-2xl overflow-hidden max-h-[90vh] flex flex-col">
            <div className="p-6 border-b border-border flex justify-between items-center bg-bg/50">
              <div>
                <h3 className="text-xl font-bold text-text">Form Hasil PE (Puskesmas)</h3>
                <p className="text-sm text-text-muted mt-1">Isi data investigasi lapangan secara langsung.</p>
              </div>
              <button 
                onClick={() => setShowPeModal(null)}
                className="p-2 hover:bg-border rounded-full transition-colors"
              >
                <X className="w-5 h-5 text-text-muted" />
              </button>
            </div>
            
            <form onSubmit={handleSubmitPE} className="p-6 space-y-4 overflow-y-auto">
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
              </div>
              
              <div className="flex items-center justify-between p-4 border border-border rounded-xl bg-bg">
                <div className="flex items-center gap-3">
                  <div>
                    <p className="font-semibold text-text text-sm">Ditemukan Jentik?</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" className="sr-only peer" checked={peForm.larva_found} onChange={(e) => setPeForm({...peForm, larva_found: e.target.checked})} />
                  <div className="w-11 h-6 bg-border peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-surface after:border-border after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                </label>
              </div>

              <div>
                <label className="block text-sm font-semibold text-text mb-1 flex items-center gap-2">
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
              </div>

              <div className="flex items-center justify-between p-4 border border-border rounded-xl bg-status-verified/20">
                <div className="flex items-center gap-3">
                  <div>
                    <p className="font-semibold text-text text-sm">Edukasi PSN 3M Plus Selesai?</p>
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
                    className="flex-1 border-border rounded-xl p-2 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary-tint file:text-primary hover:file:bg-primary hover:file:text-white transition-all"
                  />
                  {isUploadingPhoto && <Loader2 className="w-5 h-5 animate-spin text-primary" />}
                </div>
                {peForm.evidence_photo_url && (
                  <div className="mt-3 relative rounded-xl overflow-hidden border border-border bg-bg h-32">
                    <img src={peForm.evidence_photo_url.startsWith('http') ? peForm.evidence_photo_url : `${api.defaults.baseURL?.replace('/api/v1', '')}${peForm.evidence_photo_url}`} alt="Preview" className="w-full h-full object-cover" />
                  </div>
                )}
              </div>

              <div className="pt-4 mt-6 border-t border-border flex gap-3">
                <button 
                  type="button"
                  onClick={() => setShowPeModal(null)}
                  className="flex-1 py-3 px-4 bg-surface border border-border text-text font-bold rounded-xl hover:bg-bg transition-colors"
                >
                  Batal
                </button>
                <button 
                  type="submit"
                  disabled={submitting || !peForm.psn_edukasi_done}
                  className="flex-1 py-3 px-4 bg-primary text-primary-foreground font-bold rounded-xl hover:bg-primary-deep transition-colors flex justify-center items-center gap-2 disabled:opacity-50"
                >
                  {submitting ? <Loader2 className="w-5 h-5 animate-spin" /> : <CheckCircle2 className="w-5 h-5" />}
                  Submit PE
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}
