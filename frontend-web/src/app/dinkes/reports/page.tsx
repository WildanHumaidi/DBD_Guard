"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Search, Filter, Eye, AlertCircle } from "lucide-react";
import api from "@/lib/api";
import EmptyState from "@/components/common/EmptyState";
import SkeletonTable from "@/components/common/SkeletonTable";
import ErrorMessage from "@/components/common/ErrorMessage";

interface ReportItem {
  report_id: string;
  reference_no: string;
  created_at: string;
  status: string;
  category?: string;
  description?: string;
  reporter?: {
    name?: string;
    contact?: string;
    anonymous_flag?: boolean;
  };
  location?: {
    rt?: string;
    rw?: string;
    kelurahan?: string;
    kecamatan?: string;
    address?: string;
  };
}

export default function AdminReportsPage() {
  const [reports, setReports] = useState<ReportItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchReports = async () => {
    try {
      setLoading(true);
      setError("");
      const response = await api.get('/reports');
      if (response.data && response.data.data) {
        setReports(response.data.data);
      }
    } catch (err) {
      console.error("Gagal mengambil data laporan", err);
      setError("Gagal memuat data laporan dari server.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ignore = false;
    api.get('/reports')
      .then((response) => {
        if (!ignore && response.data && response.data.data) {
          setReports(response.data.data);
        }
      })
      .catch((err) => {
        if (!ignore) {
          console.error("Gagal mengambil data laporan", err);
          setError("Gagal memuat data laporan dari server.");
        }
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

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

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'MENUNGGU_PE_KADER': return <span className="px-3 py-1 bg-status-wait/20 text-status-wait rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      case 'MENUNGGU_VALIDASI_PUSKESMAS': return <span className="px-3 py-1 bg-status-verified/20 text-status-verified rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      case 'WAITING_FOGGING_APPROVAL': return <span className="px-3 py-1 bg-status-wait/20 text-status-wait rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      case 'FOGGING_APPROVED': return <span className="px-3 py-1 bg-status-verified/20 text-status-verified rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      case 'FOGGING_SCHEDULED': return <span className="px-3 py-1 bg-indigo-100 text-indigo-800 rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      case 'FOGGING_IN_PROGRESS': return <span className="px-3 py-1 bg-purple-100 text-purple-800 rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      case 'PE_SCHEDULED': return <span className="px-3 py-1 bg-indigo-100 text-indigo-800 rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      case 'COMPLETED': return <span className="px-3 py-1 bg-primary-tint text-primary rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      case 'REJECTED': return <span className="px-3 py-1 bg-red-100 text-red-800 rounded-full text-xs font-bold">{translateStatus(status).toUpperCase()}</span>;
      default: return <span className="px-3 py-1 bg-bg text-text rounded-full text-xs font-bold">{status}</span>;
    }
  };

  return (
    <div className="space-y-6">
      {/* Action Bar */}
      <div className="flex flex-col sm:flex-row justify-between gap-4">
        <div className="relative max-w-md w-full">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-text-muted" />
          <input 
            type="text" 
            placeholder="Cari nomor referensi atau nama pelapor..." 
            className="w-full pl-10 pr-4 py-2 border border-border rounded-xl focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all"
          />
        </div>
        <div className="flex items-center gap-2">
          <button className="flex items-center gap-2 px-4 py-2 bg-surface border border-border rounded-xl text-text-muted hover:bg-bg transition-colors font-medium">
            <Filter className="h-4 w-4" /> Filter Status
          </button>
        </div>
      </div>

      {error && (
        <ErrorMessage message={error} onRetry={fetchReports} />
      )}

      {/* Table */}
      <div className="bg-surface border border-border rounded-2xl shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-bg border-b border-border text-text-muted text-sm">
                <th className="p-4 font-semibold">Nomor Referensi</th>
                <th className="p-4 font-semibold">Tanggal</th>
                <th className="p-4 font-semibold">Pelapor</th>
                <th className="p-4 font-semibold">Lokasi</th>
                <th className="p-4 font-semibold">Status</th>
                <th className="p-4 font-semibold text-center">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {loading ? (
                <tr>
                  <td colSpan={6} className="p-0">
                    <SkeletonTable rows={5} columns={6} />
                  </td>
                </tr>
              ) : reports.length === 0 ? (
                <tr>
                  <td colSpan={6} className="p-8">
                    <EmptyState 
                      title="Tidak ada laporan" 
                      description="Belum ada laporan yang masuk ke sistem."
                      icon={<AlertCircle className="w-12 h-12" />} 
                    />
                  </td>
                </tr>
              ) : (
                reports.map((report) => (
                  <tr key={report.report_id} className="hover:bg-bg transition-colors">
                    <td className="p-4 font-mono font-medium text-text">{report.reference_no}</td>
                    <td className="p-4 text-text-muted">
                      {new Date(report.created_at).toLocaleDateString('id-ID', {
                        day: 'numeric', month: 'short', year: 'numeric'
                      })}
                    </td>
                    <td className="p-4 text-text-muted">{report.reporter?.name || "Anonymous"}</td>
                    <td className="p-4 text-text-muted">RT {report.location?.rt}/RW {report.location?.rw}</td>
                    <td className="p-4">{getStatusBadge(report.status)}</td>
                    <td className="p-4 text-center">
                      <Link 
                        href={`/dinkes/reports/${report.report_id}`}
                        className="inline-flex items-center justify-center p-2 text-primary hover:bg-primary-tint rounded-lg transition-colors"
                        title="Lihat Detail"
                      >
                        <Eye className="h-5 w-5" />
                      </Link>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
        
        {/* Pagination Dummy */}
        <div className="p-4 border-t border-border flex items-center justify-between text-sm text-text-muted">
          <span>Menampilkan 1 hingga {reports.length} dari {reports.length} data</span>
          <div className="flex gap-1">
            <button className="px-3 py-1 border border-border rounded hover:bg-bg" disabled>Sebelumnnya</button>
            <button className="px-3 py-1 bg-primary text-primary-foreground border border-primary rounded">1</button>
            <button className="px-3 py-1 border border-border rounded hover:bg-bg" disabled>Selanjutnya</button>
          </div>
        </div>
      </div>
    </div>
  );
}
