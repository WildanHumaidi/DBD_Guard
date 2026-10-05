"use client";

import { Activity, Users, AlertTriangle, CheckCircle2 } from "lucide-react";
import dynamic from 'next/dynamic';
import { useEffect, useState } from 'react';
import api from '@/lib/api';

const MapComponent = dynamic(() => import('@/components/MapComponent'), {
  ssr: false,
  loading: () => <div className="w-full h-full flex items-center justify-center text-text-muted bg-bg rounded-xl">Memuat Peta...</div>
});

export default function AdminDashboard() {
  const [reports, setReports] = useState<any[]>([]);
  const [klbAlerts, setKlbAlerts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchReports = async () => {
      try {
        const response = await api.get('/reports');
        if (response.data && response.data.data) {
          setReports(response.data.data);
        }
      } catch (error) {
        console.error("Failed to load reports", error);
      } finally {
        setLoading(false);
      }
    };

    const fetchKlbAlerts = async () => {
      try {
        const response = await api.get('/admin/klb/alerts');
        if (response.data?.data) {
          setKlbAlerts(response.data.data);
        }
      } catch (error) {
        console.error("Failed to load KLB alerts", error);
      }
    };

    fetchReports();
    fetchKlbAlerts();
  }, []);

  const totalReports = reports.length;
  const pendingReports = reports.filter(r => r.status === 'MENUNGGU_PE_KADER' || r.status === 'MENUNGGU_VALIDASI_PUSKESMAS').length;
  const scheduledReports = reports.filter(r => r.status === 'PE_SCHEDULED' || r.status === 'FOGGING_SCHEDULED').length;
  const completedReports = reports.filter(r => r.status === 'COMPLETED').length;

  const activeAlerts = klbAlerts.filter(a => a.is_active);

  return (
    <div className="space-y-8">
      {activeAlerts.length > 0 && (
        <div className="bg-red-50 border-l-4 border-red-600 p-4 rounded-r-lg shadow-sm">
          <div className="flex items-start">
            <div className="flex-shrink-0">
              <AlertTriangle className="h-6 w-6 text-red-600" />
            </div>
            <div className="ml-3">
              <h3 className="text-lg font-bold text-red-800">
                Peringatan Kejadian Luar Biasa (KLB) Aktif!
              </h3>
              <div className="mt-2 text-sm text-red-700">
                <p>
                  Sistem mendeteksi <strong>{activeAlerts.length}</strong> area yang melebihi ambang batas kasus DBD. 
                  Instruksi Fogging massal mungkin diperlukan.
                </p>
                <ul className="list-disc pl-5 mt-2 space-y-1 font-medium">
                  {activeAlerts.map((alert, idx) => (
                    <li key={alert.id || idx}>
                      {alert.area_type} {alert.area_name} ({alert.cases_count} Kasus Aktif)
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      )}

      <div>
        <h2 className="text-2xl font-bold text-text">Ringkasan Hari Ini</h2>
        <p className="text-text-muted">Pantau status pelaporan DBD di wilayah Anda.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-text-muted">Total Laporan</h3>
            <div className="p-2 bg-status-verified/20 text-status-verified rounded-lg">
              <Activity className="h-5 w-5" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold text-text">{totalReports}</span>
            <span className="text-sm text-text-muted">Keseluruhan</span>
          </div>
        </div>

        <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-text-muted">Menunggu Tindakan</h3>
            <div className="p-2 bg-amber-50 text-amber-600 rounded-lg">
              <AlertTriangle className="h-5 w-5" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold text-text">{pendingReports}</span>
            <span className="text-sm font-medium text-amber-600">Perlu Verifikasi</span>
          </div>
        </div>

        <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-text-muted">Tindakan Lapangan</h3>
            <div className="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
              <Users className="h-5 w-5" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold text-text">{scheduledReports}</span>
            <span className="text-sm text-text-muted">Terjadwal PE/Fogging</span>
          </div>
        </div>

        <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-text-muted">Penanganan Selesai</h3>
            <div className="p-2 bg-primary-tint text-primary rounded-lg">
              <CheckCircle2 className="h-5 w-5" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold text-text">{completedReports}</span>
            <span className="text-sm text-primary font-medium">Selesai</span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-surface rounded-2xl border border-border shadow-sm p-6 h-96 flex flex-col relative z-0">
          <h3 className="font-bold text-text mb-4">Peta Sebaran Kasus & Risiko</h3>
          <div className="flex-1 rounded-xl border border-border overflow-hidden isolate relative z-0">
            <MapComponent />
          </div>
        </div>
        
        <div className="bg-surface rounded-2xl border border-border shadow-sm p-6 flex flex-col max-h-96">
          <h3 className="font-bold text-text mb-4 shrink-0">Laporan Terbaru</h3>
          <div className="space-y-4 overflow-y-auto pr-2 flex-1">
            {loading ? (
              <div className="text-text-muted text-sm">Memuat...</div>
            ) : reports.length === 0 ? (
              <div className="text-text-muted text-sm">Belum ada laporan</div>
            ) : (
              reports.slice(0, 5).map((r, i) => (
                <div key={r.report_id || i} className="flex gap-4 border-b border-border pb-4 last:border-0">
                  <div className={`w-2 h-2 mt-1.5 rounded-full shrink-0 ${r.status === 'COMPLETED' ? 'bg-primary' : 'bg-amber-500'}`} />
                  <div>
                    <p className="text-sm font-medium text-text line-clamp-1">{r.location?.address || r.description}</p>
                    <p className="text-xs text-text-muted mt-1">
                      {r.reference_no} &bull; {new Date(r.created_at).toLocaleDateString('id-ID')}
                    </p>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
