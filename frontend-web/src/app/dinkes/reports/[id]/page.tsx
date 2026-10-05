"use client";

import { useState, use, useEffect } from "react";
import Link from "next/link";
import { ArrowLeft, CheckCircle2, XCircle, MapPin, Image as ImageIcon, ClipboardList } from "lucide-react";
import api from "@/lib/api";

export default function ReportDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = use(params);
  const id = resolvedParams.id;
  
  const [loading, setLoading] = useState(false);
  const [plannedFoggingTime, setPlannedFoggingTime] = useState("");
  const [executionStatus, setExecutionStatus] = useState("ON_THE_WAY");
  const [executionNotes, setExecutionNotes] = useState("");
  
  const [report, setReport] = useState<any>(null);
  const [fetchLoading, setFetchLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchReport = async () => {
    try {
      const response = await api.get(`/reports/${id}`);
      setReport(response.data.data);
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || "Gagal mengambil data laporan");
    } finally {
      setFetchLoading(false);
    }
  };

  useEffect(() => {
    fetchReport();
  }, [id]);

  const translateStatus = (status: string) => {
    switch (status) {
      case 'MENUNGGU_PE_KADER': return 'Menunggu PE Kader';
      case 'MENUNGGU_VALIDASI_PUSKESMAS': return 'Menunggu Validasi Puskesmas';
      case 'WAITING_FOGGING_APPROVAL': return 'Menunggu Persetujuan Fogging';
      case 'FOGGING_APPROVED': return 'Fogging Disetujui';
      case 'FOGGING_SCHEDULED': return 'Fogging Dijadwalkan';
      case 'FOGGING_IN_PROGRESS': return 'Fogging Berjalan';
      case 'PE_SCHEDULED': return 'Jadwal PE';
      case 'COMPLETED': return 'Selesai';
      case 'REJECTED': return 'Ditolak';
      default: return status;
    }
  };


  if (fetchLoading) {
    return <div className="p-8 text-center text-text-muted">Memuat data laporan...</div>;
  }

  if (error || !report) {
    return <div className="p-8 text-center text-red-500">{error || "Laporan tidak ditemukan"}</div>;
  }

  return (
    <div className="space-y-6 max-w-4xl">
      <div className="flex items-center gap-4 mb-8">
        <Link href="/dinkes/reports" className="p-2 border border-border rounded-xl hover:bg-bg transition-colors">
          <ArrowLeft className="h-5 w-5 text-text-muted" />
        </Link>
        <div>
          <h2 className="text-2xl font-bold text-text">Detail Laporan</h2>
          <p className="text-text-muted font-mono">{report.reference_no}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Kolom Kiri: Info Utama */}
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-surface rounded-2xl border border-border shadow-sm p-6">
            <h3 className="font-bold text-text mb-4 border-b border-border pb-2">Informasi Kejadian</h3>
            <div className="grid grid-cols-2 gap-y-4 gap-x-8">
              <div>
                <p className="text-sm text-text-muted">Tanggal Dilaporkan</p>
                <p className="font-medium text-text">
                  {new Date(report.created_at).toLocaleString('id-ID')}
                </p>
              </div>
              <div>
                <p className="text-sm text-text-muted">Status</p>
                <span className={`inline-block px-3 py-1 rounded-full text-xs font-bold mt-1 ${
                  report.status === 'MENUNGGU_PE_KADER' ? 'bg-status-wait/20 text-status-wait' :
                  (report.status === 'MENUNGGU_VALIDASI_PUSKESMAS' || report.status === 'FOGGING_APPROVED') ? 'bg-status-verified/20 text-status-verified' :
                  report.status === 'REJECTED' ? 'bg-red-100 text-red-800' :
                  report.status === 'WAITING_FOGGING_APPROVAL' ? 'bg-status-wait/20 text-status-wait' :
                  (report.status === 'FOGGING_SCHEDULED' || report.status === 'PE_SCHEDULED') ? 'bg-indigo-100 text-indigo-800' :
                  report.status === 'FOGGING_IN_PROGRESS' ? 'bg-purple-100 text-purple-800' :
                  'bg-primary-tint text-primary'
                }`}>
                  {translateStatus(report.status).toUpperCase()}
                </span>
              </div>
              <div className="col-span-2">
                <p className="text-sm text-text-muted">Alamat Kejadian</p>
                <div className="flex items-start gap-2 mt-1">
                  <MapPin className="h-5 w-5 text-red-500 shrink-0" />
                  <p className="font-medium text-text">{report.location.address} (RT {report.location.rt} / RW {report.location.rw})<br />
                  <span className="text-sm text-text-muted">Kec. {report.location.kecamatan}, Kel. {report.location.kelurahan}</span></p>
                </div>
              </div>
              <div className="col-span-2">
                <p className="text-sm text-text-muted">Gejala / Keterangan</p>
                <div className="p-4 bg-bg rounded-xl border border-border mt-1">
                  <p className="text-text">{report.description}</p>
                </div>
              </div>
            </div>
          </div>

          {/* Hasil PE (Read Only) */}
          {report.pe_result && (
            <div className="bg-surface rounded-2xl border border-indigo-200 shadow-sm p-6 bg-indigo-50/30 mt-6">
              <div className="flex items-center gap-2 mb-6 border-b border-indigo-100 pb-2">
                <ClipboardList className="h-5 w-5 text-indigo-600" />
                <h3 className="font-bold text-text">Hasil Penyelidikan Epidemiologi (PE)</h3>
              </div>
              
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-text-muted">Radius Pemeriksaan</p>
                  <p className="font-medium text-text">{report.pe_result.radius_meters} Meter</p>
                </div>
                <div>
                  <p className="text-sm text-text-muted">Kasus Demam Sekitar</p>
                  <p className="font-medium text-text">{report.pe_result.fever_count} Orang</p>
                </div>
                <div>
                  <p className="text-sm text-text-muted">Temuan Jentik (Larva)</p>
                  <p className="font-medium text-text">{report.pe_result.larvae_found ? 'Ya' : 'Tidak'}</p>
                </div>
                <div>
                  <p className="text-sm text-text-muted">Edukasi PSN Dilakukan</p>
                  <p className="font-medium text-text">{report.pe_result.psn_edukasi_done ? 'Ya' : 'Tidak'}</p>
                </div>
                <div>
                  <p className="text-sm text-text-muted">Status Validasi</p>
                  <p className="font-medium text-text">
                    {report.pe_result.status === 'VALIDATED' ? 'Tervalidasi' : report.pe_result.status === 'REJECTED' ? 'Ditolak' : 'Menunggu Validasi'}
                  </p>
                </div>
                <div className="col-span-2">
                  <p className="text-sm text-text-muted">Catatan Tambahan (Kader)</p>
                  <div className="p-4 bg-bg rounded-xl border border-border mt-1">
                    <p className="text-text">{report.pe_result.notes || '-'}</p>
                  </div>
                </div>
                
                {report.pe_result.evidence_photo_url && (
                  <div className="col-span-2 mt-4">
                    <p className="text-sm text-text-muted mb-2">Foto Bukti PE</p>
                    <div className="rounded-xl overflow-hidden border border-border w-64 h-64 relative">
                      <img 
                        src={report.pe_result.evidence_photo_url.startsWith('http') ? report.pe_result.evidence_photo_url : `${api.defaults.baseURL?.replace('/api/v1', '')}${report.pe_result.evidence_photo_url}`} 
                        alt="Foto Bukti PE" 
                        className="object-cover w-full h-full"
                      />
                    </div>
                  </div>
                )}
              </div>
            </div>
          )}



          {/* Form Penjadwalan Fogging jika status FOGGING_APPROVED */}
          {report.status === 'FOGGING_APPROVED' && (
            <div className="bg-surface rounded-2xl border border-status-verified/30 shadow-sm p-6 bg-status-verified/20/30 mt-6">
              <div className="flex items-center gap-2 mb-6 border-b border-blue-100 pb-2">
                <ClipboardList className="h-5 w-5 text-status-verified" />
                <h3 className="font-bold text-text">Penjadwalan Fogging (Tim Lapangan)</h3>
              </div>
              
              <div className="space-y-4">
                <p className="text-sm text-text-muted">
                  Permintaan Fogging telah disetujui. Tentukan jadwal dan tim yang akan bertugas.
                </p>
                
                <div>
                  <label className="text-sm font-medium text-text block mb-1">Rencana Waktu Eksekusi</label>
                  <input 
                    type="datetime-local"
                    value={plannedFoggingTime}
                    onChange={(e) => setPlannedFoggingTime(e.target.value)}
                    className="w-full px-3 py-2 border border-border rounded-lg text-sm focus:ring-2 focus:ring-blue-500 outline-none bg-surface"
                  />
                </div>

                <div className="flex gap-4 mt-4">
                  <button 
                    onClick={async () => {
                      const plannedTime = plannedFoggingTime;
                      if (!plannedTime) {
                        alert("Waktu eksekusi wajib diisi!");
                        return;
                      }
                      if (!report.pe_result?.fogging_request?.request_id) {
                        alert("Data Fogging Request tidak ditemukan.");
                        return;
                      }
                      setLoading(true);
                      try {
                        await api.post(`/fogging-schedules`, {
                          request_id: report.pe_result.fogging_request.request_id,
                          team_id: "33333333-3333-3333-3333-333333333333", // Dummy Team ID for MVP
                          planned_time: new Date(plannedTime).toISOString()
                        });
                        alert("Jadwal Fogging berhasil dibuat. Status menjadi FOGGING_SCHEDULED.");
                        fetchReport();
                      } catch (err: any) {
                        alert(err.response?.data?.error || err.message);
                      } finally {
                        setLoading(false);
                      }
                    }}
                    disabled={loading}
                    className="w-full bg-blue-600 text-white py-3 rounded-xl font-bold hover:bg-blue-700 transition-colors disabled:opacity-70"
                  >
                    Buat Jadwal Fogging
                  </button>
                </div>
              </div>
            </div>
          )}


        </div>

        {/* Kolom Kanan: Pelapor & Aksi */}
        <div className="space-y-6">
          <div className="bg-surface rounded-2xl border border-border shadow-sm p-6">
            <h3 className="font-bold text-text mb-4 border-b border-border pb-2">Data Pelapor</h3>
            <div className="space-y-4">
              <div>
                <p className="text-sm text-text-muted">Nama</p>
                <p className="font-medium text-text">{report.reporter.name}</p>
              </div>
              <div>
                <p className="text-sm text-text-muted">No. Kontak</p>
                <p className="font-medium text-text">{report.reporter.contact || "-"}</p>
              </div>
            </div>
          </div>

          <div className="bg-surface rounded-2xl border border-border shadow-sm p-6">
            <h3 className="font-bold text-text mb-4 border-b border-border pb-2">Lampiran Foto</h3>
            {report.photos && report.photos.length > 0 ? (
              <div className="space-y-2">
                {report.photos.map((photo: any, index: number) => (
                  <div key={index} className="text-sm text-status-verified underline">
                    {photo.storage_key.split('/').pop()}
                  </div>
                ))}
              </div>
            ) : (
              <div className="aspect-video bg-bg rounded-xl border border-border flex flex-col items-center justify-center text-slate-400">
                <ImageIcon className="h-8 w-8 mb-2" />
                <span className="text-sm">Tidak ada foto</span>
              </div>
            )}
          </div>


        </div>
      </div>
    </div>
  );
}
