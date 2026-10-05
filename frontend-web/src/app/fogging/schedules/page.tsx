"use client";

import { useEffect, useState } from "react";
import api from "@/lib/api";
import { CheckCircle2, PlayCircle, Loader2, MapPin, Calendar, Clock, X } from "lucide-react";

export default function FoggingSchedulesPage() {
  const [reports, setReports] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [successMsg, setSuccessMsg] = useState("");
  
  const [selectedReport, setSelectedReport] = useState<any>(null);
  const [scheduleId, setScheduleId] = useState<string>("");
  const [formData, setFormData] = useState({
    execution_status: "COMPLETED",
    notes: "",
    used_items: [] as { item_id: string, quantity: number }[]
  });
  const [evidenceFile, setEvidenceFile] = useState<File | null>(null);
  const [inventory, setInventory] = useState<any[]>([]);
  
  const fetchSchedules = async () => {
    setLoading(true);
    try {
      const response = await api.get("/reports");
      const allReports = response.data?.data || [];
      const scheduled = allReports.filter((r: any) => 
        r.status === "FOGGING_SCHEDULED" || r.status === "FOGGING_IN_PROGRESS"
      );
      setReports(scheduled);
    } catch (err) {
      console.error("Failed to fetch schedules", err);
      setError("Gagal mengambil data jadwal fogging.");
    } finally {
      setLoading(false);
    }
  };

  const fetchInventory = async () => {
    try {
      const res = await api.get("/inventory");
      setInventory(res.data?.data || []);
    } catch (err) {
      console.error("Failed to fetch inventory", err);
    }
  };

  useEffect(() => {
    fetchSchedules();
    fetchInventory();
  }, []);

  const openExecuteModal = async (report: any) => {
    setError("");
    setSuccessMsg("");
    setScheduleId("");
    setEvidenceFile(null);
    setSelectedReport(report);
    
    // We need to fetch the full report detail to get the schedule ID
    try {
      const res = await api.get(`/reports/${report.report_id}`);
      const detail = res.data?.data;
      if (detail && detail.pe_result?.fogging_request?.schedule?.schedule_id) {
        setScheduleId(detail.pe_result.fogging_request.schedule.schedule_id);
      } else {
        setError("Detail jadwal tidak ditemukan untuk laporan ini.");
      }
    } catch (err) {
      setError("Gagal mengambil detail laporan.");
    }
  };

  const handleExecute = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!scheduleId) return;

    if (formData.execution_status === "COMPLETED") {
      const hasUsedItems = formData.used_items.some(item => item.quantity > 0);
      if (!hasUsedItems) {
        setError("Penggunaan inventaris (bahan/alat) wajib diisi saat menyelesaikan fogging.");
        return;
      }
      if (!evidenceFile) {
        setError("Foto bukti eksekusi wajib diunggah saat menyelesaikan fogging.");
        return;
      }
    }

    setSubmitting(true);
    setError("");
    setSuccessMsg("");
    
    try {
      const validItems = formData.used_items.filter(item => item.quantity > 0);
      if (formData.execution_status === "COMPLETED" && validItems.length === 0) {
        setError("Harap isi pemakaian material (minimal 1 item dengan jumlah > 0).");
        setSubmitting(false);
        return;
      }
      
      // Update execution status along with real used items (single atomic source of truth)
      await api.patch(`/fogging/schedules/${scheduleId}/status`, {
        execution_status: formData.execution_status,
        notes: formData.notes,
        used_items: formData.execution_status === "COMPLETED" ? validItems : []
      });

      // Upload evidence photo if provided
      if (formData.execution_status === "COMPLETED" && evidenceFile) {
        const form = new FormData();
        form.append("evidence", evidenceFile);
        await api.post(`/fogging/schedules/${scheduleId}/evidence`, form, {
          headers: {
            "Content-Type": "multipart/form-data",
          }
        });
      }

      setSuccessMsg(`Eksekusi fogging berhasil diperbarui.`);
      setSelectedReport(null);
      setFormData({ execution_status: "COMPLETED", notes: "", used_items: [] });
      setEvidenceFile(null);
      fetchSchedules();
    } catch (err: any) {
      setError(err.response?.data?.error || "Gagal memperbarui eksekusi fogging");
    } finally {
      setSubmitting(false);
    }
  };

  const handleItemChange = (itemId: string, qty: number) => {
    setFormData(prev => {
      const existing = prev.used_items.find(i => i.item_id === itemId);
      if (existing) {
        return {
          ...prev,
          used_items: prev.used_items.map(i => i.item_id === itemId ? { ...i, quantity: qty } : i)
        };
      }
      return {
        ...prev,
        used_items: [...prev.used_items, { item_id: itemId, quantity: qty }]
      };
    });
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-extrabold text-text tracking-tight mb-2">Jadwal Fogging</h2>
        <p className="text-text-muted">Daftar lokasi yang dijadwalkan untuk dilakukan fogging.</p>
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
          <Loader2 className="h-8 w-8 animate-spin text-status-scheduled" />
        </div>
      ) : reports.length === 0 ? (
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-12 text-center">
          <CheckCircle2 className="h-16 w-16 text-orange-300 mx-auto mb-4" />
          <h3 className="text-xl font-bold text-text mb-2">Semua Selesai!</h3>
          <p className="text-text-muted">Tidak ada jadwal fogging yang perlu dieksekusi saat ini.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
          {reports.map((report) => (
            <div key={report.report_id} className="bg-surface rounded-2xl shadow-sm border border-border overflow-hidden flex flex-col">
              <div className="p-6 flex-1">
                <div className="flex justify-between items-start mb-4">
                  <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-bold ${
                    report.status === 'FOGGING_IN_PROGRESS' 
                      ? 'bg-status-verified/20 text-status-verified'
                      : 'bg-status-scheduled/20 text-status-scheduled'
                  }`}>
                    {report.status === 'FOGGING_IN_PROGRESS' ? 'Sedang Berjalan' : 'Terjadwal'}
                  </span>
                  <span className="text-xs font-mono text-text-muted">{report.reference_no}</span>
                </div>
                
                <h3 className="text-lg font-bold text-text mb-2">{report.category}</h3>
                
                <div className="space-y-2 mt-4">
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <MapPin className="h-4 w-4 text-text-muted" />
                    <span className="truncate">{report.location?.address || 'Lokasi tidak spesifik'}</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-text-muted">
                    <Calendar className="h-4 w-4 text-text-muted" />
                    <span>Dilaporkan: {new Date(report.created_at).toLocaleDateString('id-ID')}</span>
                  </div>
                </div>
              </div>
              
              <div className="border-t border-border bg-bg p-4">
                <button
                  onClick={() => openExecuteModal(report)}
                  className="w-full py-2.5 px-4 bg-orange-600 text-white font-bold rounded-xl hover:bg-orange-700 transition-colors flex justify-center items-center gap-2"
                >
                  <PlayCircle className="w-4 h-4" />
                  Update Status Eksekusi
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Execution Modal */}
      {selectedReport && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-primary-deep/50 backdrop-blur-sm animate-fade-in">
          <div className="bg-surface rounded-2xl shadow-xl w-full max-w-lg overflow-hidden flex flex-col max-h-[90vh]">
            <div className="p-6 border-b border-border flex justify-between items-center bg-bg">
              <h3 className="font-bold text-lg text-text">Update Eksekusi Fogging</h3>
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

              {error && (
                <div className="mb-4 bg-red-50 text-red-700 p-3 rounded-lg text-sm border border-red-100">
                  {error}
                </div>
              )}

              {!scheduleId ? (
                <div className="flex justify-center p-8">
                  <Loader2 className="w-8 h-8 animate-spin text-status-scheduled" />
                </div>
              ) : (
                <form id="execute-form" onSubmit={handleExecute} className="space-y-6">
                  {/* Status Eksekusi Section */}
                  <div className="bg-surface border border-border p-4 rounded-xl shadow-sm">
                    <label className="block text-sm font-bold text-text mb-2">Status Eksekusi <span className="text-red-500">*</span></label>
                    <p className="text-xs text-text-muted mb-3">Pilih status terbaru dari jadwal fogging ini.</p>
                    <select
                      required
                      value={formData.execution_status}
                      onChange={(e) => setFormData({...formData, execution_status: e.target.value})}
                      className="w-full px-4 py-3 bg-bg border border-border rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none font-medium text-text"
                    >
                      <option value="IN_PROGRESS">Sedang Berjalan (Memulai Fogging)</option>
                      <option value="COMPLETED">Selesai (Fogging Berhasil Dilaksanakan)</option>
                      <option value="FAILED">Gagal / Dibatalkan</option>
                    </select>
                  </div>

                  {/* Inventaris Section */}
                  <div className={`p-4 rounded-xl border shadow-sm transition-colors ${formData.execution_status === 'COMPLETED' ? 'bg-status-scheduled/20/50 border-status-scheduled/30' : 'bg-surface border-border'}`}>
                    <label className="block text-sm font-bold text-text mb-1">
                      Pemakaian Inventaris {formData.execution_status === 'COMPLETED' && <span className="text-red-500">*</span>}
                    </label>
                    <p className="text-xs text-text-muted mb-4">
                      {formData.execution_status === 'COMPLETED' 
                        ? "Wajib diisi! Masukkan jumlah bahan/alat yang digunakan selama proses fogging." 
                        : "Catat bahan atau alat yang digunakan (opsional jika belum selesai)."}
                    </p>
                    
                    {inventory.length === 0 ? (
                      <div className="p-4 bg-bg rounded-lg border border-border text-center">
                        <p className="text-sm text-text-muted italic">Data inventaris tidak tersedia di sistem.</p>
                      </div>
                    ) : (
                      <div className="space-y-3 max-h-60 overflow-y-auto pr-2 custom-scrollbar">
                        {inventory.map((item) => (
                          <div key={item.item_id} className="flex items-center justify-between gap-4 bg-surface p-3 rounded-lg border border-border hover:border-orange-300 transition-colors">
                            <div className="flex-1">
                              <p className="text-sm font-bold text-text uppercase">{item.item_type || item.item_name || 'Item'}</p>
                              <p className="text-xs font-medium text-text-muted">Stok Tersedia: <span className="text-text">{item.quantity} {item.unit}</span></p>
                            </div>
                            <div className="flex items-center gap-2">
                              <div className="relative">
                                <input
                                  type="number"
                                  min="0"
                                  max={item.quantity}
                                  placeholder="0"
                                  onChange={(e) => handleItemChange(item.item_id, parseInt(e.target.value || "0", 10))}
                                  className="w-24 pl-3 pr-8 py-2 text-sm font-semibold text-text border border-border rounded-md focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none text-right"
                                />
                                <span className="absolute right-3 top-1/2 -translate-y-1/2 text-xs font-medium text-text-muted pointer-events-none">
                                  {item.unit}
                                </span>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Upload Foto Bukti Section */}
                  {formData.execution_status === "COMPLETED" && (
                    <div className="bg-surface border border-border p-4 rounded-xl shadow-sm">
                      <label className="block text-sm font-bold text-text mb-2">Foto Bukti Eksekusi <span className="text-red-500">*</span></label>
                      <p className="text-xs text-text-muted mb-3">Wajib diunggah! Lampirkan foto kondisi lapangan saat fogging dilakukan.</p>
                      <input
                        type="file"
                        accept="image/*"
                        required
                        onChange={(e) => setEvidenceFile(e.target.files?.[0] || null)}
                        className="w-full text-sm text-text file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-orange-50 file:text-orange-700 hover:file:bg-orange-100"
                      />
                    </div>
                  )}

                  {/* Catatan Section */}
                  <div className="bg-surface border border-border p-4 rounded-xl shadow-sm">
                    <label className="block text-sm font-bold text-text mb-2">Catatan Laporan</label>
                    <textarea
                      rows={3}
                      value={formData.notes}
                      onChange={(e) => setFormData({...formData, notes: e.target.value})}
                      placeholder="Contoh: Fogging berjalan lancar di 50 rumah warga, cuaca cerah..."
                      className="w-full px-4 py-3 bg-bg border border-border rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 outline-none resize-none text-sm"
                    ></textarea>
                  </div>
                </form>
              )}
            </div>
            
            <div className="p-6 border-t border-border bg-bg flex gap-3">
              <button
                type="button"
                onClick={() => setSelectedReport(null)}
                className="flex-1 px-4 py-2.5 border border-border text-text font-bold rounded-xl hover:bg-bg transition-colors"
              >
                Tutup
              </button>
              {scheduleId && (
                <button
                  type="submit"
                  form="execute-form"
                  disabled={submitting}
                  className="flex-1 px-4 py-2.5 bg-orange-600 text-white font-bold rounded-xl hover:bg-orange-700 transition-colors disabled:opacity-50 flex justify-center items-center gap-2"
                >
                  {submitting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                  Simpan Status
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
