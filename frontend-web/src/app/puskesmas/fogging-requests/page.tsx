"use client";

import { useEffect, useState } from "react";
import api from "@/lib/api";
import { CheckCircle2, Loader2, MapPin, Calendar, Clock, Check, X, CalendarDays } from "lucide-react";

interface TimFogging {
  user_id: string;
  name: string;
  contact?: string;
}

interface FoggingRequestPeResult {
  pe_radius?: number;
  larva_found?: boolean;
  fever_count?: number;
  notes?: string;
  evidence_photo_url?: string;
  fogging_request?: {
    request_id: string;
  };
}

interface FoggingRequestReport {
  report_id: string;
  reference_no: string;
  category: string;
  description: string;
  status: string;
  created_at: string;
  location?: {
    address: string;
  };
  reporter?: {
    name?: string;
    anonymous_flag?: boolean;
  };
  pe_result?: FoggingRequestPeResult;
}

export default function FoggingRequestsPage() {
  const [reports, setReports] = useState<FoggingRequestReport[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [successMsg, setSuccessMsg] = useState("");
  const [timFogging, setTimFogging] = useState<TimFogging[]>([]);
  
  // Modal State for Approval
  const [selectedApproval, setSelectedApproval] = useState<FoggingRequestReport | null>(null);
  const [approvalData, setApprovalData] = useState({
    decision: "APPROVED",
    reason: ""
  });

  // Modal State for Scheduling
  const [selectedSchedule, setSelectedSchedule] = useState<FoggingRequestReport | null>(null);
  const [scheduleData, setScheduleData] = useState({
    team_id: "",
    planned_time: ""
  });

  const fetchData = async () => {
    setLoading(true);
    try {
      const [reportsRes, timRes] = await Promise.all([
        api.get("/reports"),
        api.get("/puskesmas/tim-fogging")
      ]);
      
      const allReports = (reportsRes.data?.data || []) as FoggingRequestReport[];
      const pending = allReports.filter((r) => 
        r.status === "WAITING_FOGGING_APPROVAL" || 
        r.status === "FOGGING_APPROVED"
      );
      setReports(pending);
      setTimFogging(timRes.data || []);
    } catch (err: unknown) {
      console.error("Failed to fetch data", err);
      setError("Gagal mengambil data. Coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ignore = false;
    Promise.all([
      api.get("/reports"),
      api.get("/puskesmas/tim-fogging")
    ])
      .then(([reportsRes, timRes]) => {
        if (!ignore) {
          const allReports = (reportsRes.data?.data || []) as FoggingRequestReport[];
          const pending = allReports.filter((r) => 
            r.status === "WAITING_FOGGING_APPROVAL" || 
            r.status === "FOGGING_APPROVED"
          );
          setReports(pending);
          setTimFogging(timRes.data || []);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          console.error("Failed to fetch data", err);
          setError("Gagal mengambil data. Coba lagi.");
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

  const handleApproveSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedApproval || !selectedApproval.pe_result?.fogging_request?.request_id) return;

    setSubmitting(true);
    setError("");
    setSuccessMsg("");
    
    try {
      await api.post(`/kepala-puskesmas/fogging-requests/${selectedApproval.pe_result.fogging_request.request_id}/approve`, {
        decision: approvalData.decision,
        reason: approvalData.reason
      });
      setSuccessMsg(`Persetujuan fogging untuk ${selectedApproval.reference_no} berhasil disimpan.`);
      setSelectedApproval(null);
      setApprovalData({ decision: "APPROVED", reason: "" });
      fetchData();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || "Gagal menyimpan persetujuan";
      setError(msg);
    } finally {
      setSubmitting(false);
    }
  };

  const handleScheduleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedSchedule || !selectedSchedule.pe_result?.fogging_request?.request_id) return;

    if (!scheduleData.team_id || !scheduleData.planned_time) {
      setError("Tim dan Waktu pelaksanaan harus diisi.");
      return;
    }

    setSubmitting(true);
    setError("");
    setSuccessMsg("");
    
    try {
      await api.post(`/fogging-schedules`, {
        request_id: selectedSchedule.pe_result.fogging_request.request_id,
        team_id: scheduleData.team_id,
        planned_time: new Date(scheduleData.planned_time).toISOString()
      });
      setSuccessMsg(`Jadwal fogging untuk ${selectedSchedule.reference_no} berhasil dibuat.`);
      setSelectedSchedule(null);
      setScheduleData({ team_id: "", planned_time: "" });
      fetchData();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || "Gagal membuat jadwal";
      setError(msg);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-extrabold text-text tracking-tight mb-2">Jadwal Fogging</h2>
        <p className="text-text-muted">Persetujuan rekomendasi fogging dan penjadwalan tim operasional.</p>
      </div>

      {error && !selectedApproval && !selectedSchedule && (
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
          <h3 className="text-xl font-bold text-text mb-2">Tidak Ada Antrean Fogging</h3>
          <p className="text-text-muted">Saat ini tidak ada laporan yang membutuhkan persetujuan atau penjadwalan fogging.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
          {reports.map((report) => (
            <div key={report.report_id} className="bg-surface rounded-2xl shadow-sm border border-border overflow-hidden flex flex-col">
              <div className="p-6 flex-1">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex gap-2">
                    <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-bold ${report.status === 'WAITING_FOGGING_APPROVAL' ? 'bg-status-wait/20 text-status-wait' : 'bg-primary-tint text-primary'}`}>
                      {report.status === 'WAITING_FOGGING_APPROVAL' ? 'Menunggu Persetujuan' : 'Menunggu Jadwal'}
                    </span>
                    
                    {(() => {
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
                  <p className="font-semibold text-blue-900 mb-1">Hasil Rekomendasi PE:</p>
                  <ul className="list-disc list-inside text-status-verified space-y-1 text-xs">
                    <li>Radius: {report.pe_result?.pe_radius || '-'} meter</li>
                    <li>Jentik Ditemukan: {report.pe_result?.larva_found ? 'Ya' : 'Tidak'}</li>
                    <li>Kasus Demam: {report.pe_result?.fever_count || 0} orang</li>
                    <li className="font-bold text-red-700 mt-1">Status: DIREKOMENDASIKAN FOGGING</li>
                  </ul>
                </div>
                
                <div className="space-y-2 mt-auto">
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <MapPin className="h-4 w-4 text-text-muted" />
                    <span className="truncate">{report.location?.address || 'Lokasi tidak spesifik'}</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <Calendar className="h-4 w-4 text-text-muted" />
                    <span>Tgl Laporan: {new Date(report.created_at).toLocaleString('id-ID')}</span>
                  </div>
                </div>
              </div>
              
              <div className="border-t border-border bg-bg p-4">
                {report.status === 'WAITING_FOGGING_APPROVAL' ? (
                  <button
                    onClick={() => { setSelectedApproval(report); setError(""); }}
                    className="w-full py-2.5 px-4 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-colors flex justify-center items-center gap-2"
                  >
                    <Check className="w-4 h-4" />
                    Beri Persetujuan
                  </button>
                ) : (
                  <button
                    onClick={() => { setSelectedSchedule(report); setError(""); }}
                    className="w-full py-2.5 px-4 bg-indigo-600 text-white font-bold rounded-xl hover:bg-indigo-700 transition-colors flex justify-center items-center gap-2"
                  >
                    <CalendarDays className="w-4 h-4" />
                    Buat Jadwal Fogging
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Approval Modal */}
      {selectedApproval && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-primary-deep/50 backdrop-blur-sm animate-fade-in">
          <div className="bg-surface rounded-2xl shadow-xl w-full max-w-lg overflow-hidden flex flex-col max-h-[90vh]">
            <div className="p-6 border-b border-border flex justify-between items-center bg-bg">
              <h3 className="font-bold text-lg text-text">Persetujuan Fogging</h3>
              <button 
                onClick={() => setSelectedApproval(null)}
                className="text-text-muted hover:text-text-muted p-1"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            <div className="p-6 overflow-y-auto">
              <div className="mb-6 pb-6 border-b border-border">
                <p className="text-xs text-text-muted font-mono mb-1">{selectedApproval.reference_no}</p>
                <p className="text-sm font-semibold text-text">{selectedApproval.location?.address}</p>
              </div>

              <div className="bg-status-verified/20 border border-blue-100 rounded-lg p-4 mb-6">
                <h4 className="font-bold text-blue-900 mb-2">Data PE dari Kader</h4>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="block text-text-muted text-xs">Radius</span>
                    <span className="font-medium text-text">{selectedApproval.pe_result?.pe_radius || '-'} meter</span>
                  </div>
                  <div>
                    <span className="block text-text-muted text-xs">Jentik</span>
                    <span className="font-medium text-text">{selectedApproval.pe_result?.larva_found ? 'Ditemukan' : 'Tidak Ditemukan'}</span>
                  </div>
                  <div>
                    <span className="block text-text-muted text-xs">Penderita Demam</span>
                    <span className="font-medium text-text">{selectedApproval.pe_result?.fever_count || 0} orang</span>
                  </div>
                </div>
              </div>

              {(selectedApproval.pe_result?.notes || selectedApproval.pe_result?.evidence_photo_url) && (
                <div className="bg-bg border border-border rounded-lg p-4 mb-6">
                  <h4 className="font-bold text-text mb-2">Narasi & Bukti Kunjungan</h4>
                  {selectedApproval.pe_result?.notes && (
                    <div className="mb-3">
                      <p className="text-xs text-text-muted font-semibold mb-1">Catatan Kader:</p>
                      <p className="text-sm text-text bg-surface p-3 rounded-lg border border-border">
                        {selectedApproval.pe_result.notes}
                      </p>
                    </div>
                  )}
                  {selectedApproval.pe_result?.evidence_photo_url && (
                    <div>
                      <p className="text-xs text-text-muted font-semibold mb-1">Foto Bukti:</p>
                      <div className="rounded-lg overflow-hidden border border-border">
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img 
                          src={selectedApproval.pe_result.evidence_photo_url.startsWith('http') ? selectedApproval.pe_result.evidence_photo_url : `${api.defaults.baseURL?.replace('/api/v1', '')}${selectedApproval.pe_result.evidence_photo_url}`} 
                          alt="Bukti Kunjungan Kader" 
                          className="w-full h-auto max-h-64 object-contain bg-bg"
                        />
                      </div>
                    </div>
                  )}
                </div>
              )}

              {error && (
                <div className="mb-4 bg-red-50 text-red-700 p-3 rounded-lg text-sm border border-red-100">
                  {error}
                </div>
              )}

              <form id="approval-form" onSubmit={handleApproveSubmit} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-text mb-2">Keputusan Kepala Puskesmas *</label>
                  <div className="flex gap-4">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input 
                        type="radio" 
                        name="decision" 
                        checked={approvalData.decision === "APPROVED"}
                        onChange={() => setApprovalData({...approvalData, decision: "APPROVED"})}
                        className="text-status-verified focus:ring-blue-500" 
                      />
                      <span className="text-sm font-medium text-text">Setujui (Lanjut Jadwal)</span>
                    </label>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input 
                        type="radio" 
                        name="decision" 
                        checked={approvalData.decision === "REJECTED"}
                        onChange={() => setApprovalData({...approvalData, decision: "REJECTED"})}
                        className="text-red-600 focus:ring-red-500" 
                      />
                      <span className="text-sm font-medium text-text">Tolak (Batalkan)</span>
                    </label>
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-text mb-1">Catatan Tambahan (Opsional)</label>
                  <textarea
                    rows={3}
                    value={approvalData.reason}
                    onChange={(e) => setApprovalData({...approvalData, reason: e.target.value})}
                    placeholder="Alasan persetujuan atau penolakan..."
                    className="w-full px-4 py-2 border border-border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all resize-none"
                  ></textarea>
                </div>
              </form>
            </div>
            
            <div className="p-6 border-t border-border bg-bg flex gap-3">
              <button
                type="button"
                onClick={() => setSelectedApproval(null)}
                className="flex-1 px-4 py-2.5 border border-border text-text font-bold rounded-xl hover:bg-bg transition-colors"
              >
                Batal
              </button>
              <button
                type="submit"
                form="approval-form"
                disabled={submitting}
                className="flex-1 px-4 py-2.5 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-colors disabled:opacity-50 flex justify-center items-center gap-2"
              >
                {submitting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                Simpan Keputusan
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Schedule Modal */}
      {selectedSchedule && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-primary-deep/50 backdrop-blur-sm animate-fade-in">
          <div className="bg-surface rounded-2xl shadow-xl w-full max-w-lg overflow-hidden flex flex-col max-h-[90vh]">
            <div className="p-6 border-b border-border flex justify-between items-center bg-bg">
              <h3 className="font-bold text-lg text-text">Jadwalkan Fogging</h3>
              <button 
                onClick={() => setSelectedSchedule(null)}
                className="text-text-muted hover:text-text-muted p-1"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            <div className="p-6 overflow-y-auto">
              <div className="mb-6 pb-6 border-b border-border">
                <p className="text-xs text-text-muted font-mono mb-1">{selectedSchedule.reference_no}</p>
                <p className="text-sm font-semibold text-text">{selectedSchedule.location?.address}</p>
              </div>

              {error && (
                <div className="mb-4 bg-red-50 text-red-700 p-3 rounded-lg text-sm border border-red-100">
                  {error}
                </div>
              )}

              <form id="schedule-form" onSubmit={handleScheduleSubmit} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-text mb-1">Pilih Tim Fogging *</label>
                  <select
                    value={scheduleData.team_id}
                    onChange={(e) => setScheduleData({...scheduleData, team_id: e.target.value})}
                    className="w-full px-4 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
                    required
                  >
                    <option value="">-- Pilih Tim --</option>
                    {timFogging.map((tim) => (
                      <option key={tim.user_id} value={tim.user_id}>{tim.name} ({tim.contact})</option>
                    ))}
                  </select>
                  {timFogging.length === 0 && (
                    <p className="text-xs text-red-500 mt-1">Tidak ada tim fogging yang terdaftar. Hubungi admin Dinkes.</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-text mb-1">Tanggal & Waktu Pelaksanaan *</label>
                  <input
                    type="datetime-local"
                    value={scheduleData.planned_time}
                    onChange={(e) => setScheduleData({...scheduleData, planned_time: e.target.value})}
                    className="w-full px-4 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
                    required
                  />
                  <p className="text-xs text-text-muted mt-1">Usahakan jadwalkan sebelum batas SLA 72 jam terlewati.</p>
                </div>
              </form>
            </div>
            
            <div className="p-6 border-t border-border bg-bg flex gap-3">
              <button
                type="button"
                onClick={() => setSelectedSchedule(null)}
                className="flex-1 px-4 py-2.5 border border-border text-text font-bold rounded-xl hover:bg-bg transition-colors"
              >
                Batal
              </button>
              <button
                type="submit"
                form="schedule-form"
                disabled={submitting || timFogging.length === 0}
                className="flex-1 px-4 py-2.5 bg-indigo-600 text-white font-bold rounded-xl hover:bg-indigo-700 transition-colors disabled:opacity-50 flex justify-center items-center gap-2"
              >
                {submitting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                Buat Jadwal
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
