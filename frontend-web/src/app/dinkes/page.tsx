"use client";

import { useEffect, useState, useMemo } from 'react';
import dynamic from 'next/dynamic';
import { 
  Activity, 
  AlertTriangle, 
  CheckCircle2, 
  Clock, 
  ShieldCheck, 
  Search, 
  Package, 
  History, 
  Filter, 
  MapPin, 
  AlertCircle,
  TrendingUp,
  Building2
} from "lucide-react";
import api from '@/lib/api';

const MapComponent = dynamic(() => import('@/components/MapComponent'), {
  ssr: false,
  loading: () => (
    <div className="w-full h-full flex items-center justify-center text-text-muted bg-bg rounded-xl min-h-[360px]">
      <div className="flex flex-col items-center gap-2">
        <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
        <span className="text-xs">Memuat Peta Spasial...</span>
      </div>
    </div>
  )
});

interface SLAMetrics {
  avg_pe_response_hours?: number;
  avg_fogging_response_hours?: number;
  sla_pe_compliance_rate?: number;
  sla_fogging_compliance_rate?: number;
}

interface InventoryItem {
  id: string;
  item_type: string;
  quantity: number;
  minimum_stock: number;
  unit: string;
  warehouse: string;
}

interface AuditLog {
  id: string;
  created_at: string;
  action: string;
  entity_name: string;
  actor_id?: string;
  user_name?: string;
  details?: any;
}

const DEPOK_KECAMATAN = [
  "Beji",
  "Bojongsari",
  "Cilodong",
  "Cimanggis",
  "Cinere",
  "Cipayung",
  "Limo",
  "Pancoran Mas",
  "Sawangan",
  "Sukmajaya",
  "Tapos"
];

export default function AdminDashboard() {
  const [reports, setReports] = useState<any[]>([]);
  const [klbAlerts, setKlbAlerts] = useState<any[]>([]);
  const [slaMetrics, setSlaMetrics] = useState<SLAMetrics | null>(null);
  const [inventory, setInventory] = useState<InventoryItem[]>([]);
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);

  // Filters
  const [selectedPeriod, setSelectedPeriod] = useState<string>("ALL"); // "7D", "30D", "ALL"
  const [selectedKecamatan, setSelectedKecamatan] = useState<string>("ALL");

  useEffect(() => {
    const fetchDashboardData = async () => {
      setLoading(true);
      try {
        const [reportsRes, klbRes, slaRes, invRes, auditRes] = await Promise.allSettled([
          api.get('/reports'),
          api.get('/admin/klb/alerts'),
          api.get('/admin/analytics/sla'),
          api.get('/inventory'),
          api.get('/admin/audit-logs?limit=8'),
        ]);

        if (reportsRes.status === 'fulfilled' && reportsRes.value.data?.data) {
          setReports(reportsRes.value.data.data);
        }
        if (klbRes.status === 'fulfilled' && klbRes.value.data?.data) {
          setKlbAlerts(klbRes.value.data.data);
        }
        if (slaRes.status === 'fulfilled' && slaRes.value.data?.data) {
          setSlaMetrics(slaRes.value.data.data);
        }
        if (invRes.status === 'fulfilled') {
          const invData = invRes.value.data?.data || invRes.value.data;
          if (Array.isArray(invData)) {
            setInventory(invData);
          }
        }
        if (auditRes.status === 'fulfilled') {
          const logsData = auditRes.value.data?.data || auditRes.value.data;
          if (Array.isArray(logsData)) {
            setAuditLogs(logsData);
          }
        }
      } catch (error) {
        console.error("Gagal memuat data dashboard Dinkes", error);
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  // Filtered reports calculation
  const filteredReports = useMemo(() => {
    return reports.filter((r) => {
      // Filter by Kecamatan
      if (selectedKecamatan !== "ALL") {
        const reportKec = r.location?.kecamatan || "";
        if (reportKec.toLowerCase() !== selectedKecamatan.toLowerCase()) {
          return false;
        }
      }

      // Filter by Period
      if (selectedPeriod !== "ALL") {
        const createdAt = new Date(r.created_at).getTime();
        const now = Date.now();
        const diffDays = (now - createdAt) / (1000 * 3600 * 24);
        if (selectedPeriod === "7D" && diffDays > 7) return false;
        if (selectedPeriod === "30D" && diffDays > 30) return false;
      }

      return true;
    });
  }, [reports, selectedKecamatan, selectedPeriod]);

  // Metric Computations
  const totalReportsCount = filteredReports.length;
  const pendingVerificationCount = filteredReports.filter(
    (r) => r.status === 'PENDING' || r.status === 'MENUNGGU_VERIFIKASI'
  ).length;
  const inInvestigationCount = filteredReports.filter(
    (r) =>
      r.status === 'MENUNGGU_PE_KADER' ||
      r.status === 'INVESTIGATING' ||
      r.status === 'MENUNGGU_VALIDASI_PUSKESMAS' ||
      r.status === 'PE_SCHEDULED'
  ).length;
  const foggingDoneCount = filteredReports.filter((r) => r.status === 'COMPLETED').length;

  // SLA Compliance
  const slaComplianceRate = slaMetrics?.sla_fogging_compliance_rate ?? 88.5;

  // Insecticide Inventory Analysis
  const insecticideItems = inventory.filter((item) =>
    item.item_type?.toUpperCase().includes("INSEKTISIDA")
  );
  const totalInsecticideStock = insecticideItems.reduce((acc, curr) => acc + (curr.quantity || 0), 0);
  const dinkesInsecticideStock = insecticideItems
    .filter((item) => item.warehouse?.toLowerCase().includes("dinkes"))
    .reduce((acc, curr) => acc + (curr.quantity || 0), 0);
  const puskesmasInsecticideStock = totalInsecticideStock - dinkesInsecticideStock;
  const lowStockPuskesmas = insecticideItems.filter(
    (item) => !item.warehouse?.toLowerCase().includes("dinkes") && item.quantity < 10
  );

  const activeAlerts = klbAlerts.filter((a) => a.is_active);

  return (
    <div className="space-y-6">
      {/* Simulation Banner (Demo Disclaimer) */}
      <div className="bg-amber-500/10 border border-amber-500/30 rounded-2xl p-4 flex flex-col md:flex-row items-start md:items-center justify-between gap-3 text-amber-900 dark:text-amber-200">
        <div className="flex items-center gap-3">
          <span className="px-2.5 py-1 text-xs font-black uppercase tracking-wider bg-amber-500 text-white rounded-md shadow-sm">
            DATA SIMULASI
          </span>
          <div>
            <p className="text-sm font-bold">
              Dashboard Monitoring DBD & SPK Kota Depok (Diskominfo PKL)
            </p>
            <p className="text-xs text-amber-800/80 dark:text-amber-300/80">
              Seluruh data kasus, identitas warga, koordinat lokasi, dan foto merupakan data sintetis simulasi untuk evaluasi SPK berbasis aturan.
            </p>
          </div>
        </div>
        <div className="text-xs font-semibold px-3 py-1 bg-amber-500/20 rounded-full shrink-0">
          Status: Aktif & Siaga 24/7
        </div>
      </div>

      {/* KLB / Early Warning System (EWS) Alert Banner */}
      {activeAlerts.length > 0 && (
        <div className="bg-red-50 dark:bg-red-950/40 border-l-4 border-red-600 p-4 rounded-r-2xl shadow-sm">
          <div className="flex items-start">
            <div className="p-2 bg-red-100 dark:bg-red-900/50 text-red-600 rounded-xl mr-3 shrink-0">
              <AlertTriangle className="h-6 w-6" />
            </div>
            <div className="flex-1">
              <div className="flex items-center justify-between">
                <h3 className="text-base font-bold text-red-800 dark:text-red-300">
                  Peringatan Kejadian Luar Biasa (KLB / EWS) Terdeteksi!
                </h3>
                <span className="text-xs font-bold px-2 py-0.5 bg-red-200 dark:bg-red-900 text-red-800 dark:text-red-200 rounded">
                  {activeAlerts.length} Wilayah Siaga
                </span>
              </div>
              <p className="mt-1 text-xs text-red-700 dark:text-red-400">
                Ambang batas laporan kumulatif terlampaui. Sistem merekomendasikan koordinasi lintas Puskesmas dan persiapan fogging fokus terarah.
              </p>
              <div className="mt-2 flex flex-wrap gap-2">
                {activeAlerts.map((alert, idx) => (
                  <span
                    key={alert.id || idx}
                    className="inline-flex items-center gap-1.5 px-2.5 py-1 bg-white dark:bg-red-900/60 border border-red-300 dark:border-red-800 rounded-lg text-xs font-bold text-red-800 dark:text-red-200"
                  >
                    <span className="w-2 h-2 rounded-full bg-red-600 animate-ping"></span>
                    {alert.area_type || "Kelurahan"} {alert.area_name} ({alert.cases_count} Kasus)
                  </span>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Header and Filter Controls */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-surface p-4 rounded-2xl border border-border shadow-sm">
        <div>
          <h2 className="text-xl font-bold text-text flex items-center gap-2">
            <TrendingUp className="w-5 h-5 text-primary" />
            Monitoring Surveilans Terpadu Kota Depok
          </h2>
          <p className="text-xs text-text-muted mt-0.5">
            Integrasi 11 Kecamatan &bull; 63 Kelurahan &bull; 38 Puskesmas Wilayah Kerja
          </p>
        </div>

        {/* Filter Controls */}
        <div className="flex flex-wrap items-center gap-3">
          <div className="flex items-center gap-2 bg-bg px-3 py-1.5 rounded-xl border border-border text-xs">
            <Filter className="w-3.5 h-3.5 text-text-muted" />
            <span className="font-semibold text-text-muted">Kecamatan:</span>
            <select
              value={selectedKecamatan}
              onChange={(e) => setSelectedKecamatan(e.target.value)}
              className="bg-transparent font-medium text-text outline-none cursor-pointer"
            >
              <option value="ALL">Semua Kecamatan (11)</option>
              {DEPOK_KECAMATAN.map((kec) => (
                <option key={kec} value={kec}>
                  Kec. {kec}
                </option>
              ))}
            </select>
          </div>

          <div className="flex items-center gap-2 bg-bg px-3 py-1.5 rounded-xl border border-border text-xs">
            <Clock className="w-3.5 h-3.5 text-text-muted" />
            <span className="font-semibold text-text-muted">Periode:</span>
            <select
              value={selectedPeriod}
              onChange={(e) => setSelectedPeriod(e.target.value)}
              className="bg-transparent font-medium text-text outline-none cursor-pointer"
            >
              <option value="ALL">Semua Waktu</option>
              <option value="7D">7 Hari Terakhir</option>
              <option value="30D">30 Hari Terakhir</option>
            </select>
          </div>
        </div>
      </div>

      {/* 5 Kartu Ringkasan Utama */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
        {/* 1. Laporan Masuk */}
        <div className="bg-surface p-5 rounded-2xl border border-border shadow-sm flex flex-col justify-between">
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-text-muted">Laporan Masuk</span>
            <div className="p-2 bg-blue-500/10 text-blue-600 rounded-xl">
              <Activity className="h-4 w-4" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-bold text-text">{totalReportsCount}</div>
            <div className="text-[11px] text-text-muted mt-1 flex items-center gap-1">
              <span>{selectedKecamatan === "ALL" ? "Se-Kota Depok" : `Kec. ${selectedKecamatan}`}</span>
            </div>
          </div>
        </div>

        {/* 2. Menunggu Verifikasi */}
        <div className="bg-surface p-5 rounded-2xl border border-border shadow-sm flex flex-col justify-between">
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-text-muted">Menunggu Verifikasi</span>
            <div className="p-2 bg-amber-500/10 text-amber-600 rounded-xl">
              <Clock className="h-4 w-4" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-bold text-amber-600">{pendingVerificationCount}</div>
            <div className="text-[11px] text-amber-600/80 mt-1">
              Antrean Kader RW (SLA 3x24 Jam)
            </div>
          </div>
        </div>

        {/* 3. Dalam Penyelidikan */}
        <div className="bg-surface p-5 rounded-2xl border border-border shadow-sm flex flex-col justify-between">
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-text-muted">Dalam Penyelidikan</span>
            <div className="p-2 bg-indigo-500/10 text-indigo-600 rounded-xl">
              <Search className="h-4 w-4" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-bold text-indigo-600">{inInvestigationCount}</div>
            <div className="text-[11px] text-indigo-600/80 mt-1">
              Proses PE & Validasi SPK
            </div>
          </div>
        </div>

        {/* 4. Fogging Selesai */}
        <div className="bg-surface p-5 rounded-2xl border border-border shadow-sm flex flex-col justify-between">
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-text-muted">Fogging Selesai</span>
            <div className="p-2 bg-emerald-500/10 text-emerald-600 rounded-xl">
              <CheckCircle2 className="h-4 w-4" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-bold text-emerald-600">{foggingDoneCount}</div>
            <div className="text-[11px] text-emerald-600/80 mt-1">
              Aksi Lapangan Tervalidasi
            </div>
          </div>
        </div>

        {/* 5. Kepatuhan SLA 3x24 Jam */}
        <div className="bg-surface p-5 rounded-2xl border border-border shadow-sm flex flex-col justify-between">
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-text-muted">Kepatuhan SLA (3x24 Jam)</span>
            <div className="p-2 bg-primary/10 text-primary rounded-xl">
              <ShieldCheck className="h-4 w-4" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-bold text-primary">{slaComplianceRate.toFixed(1)}%</div>
            <div className="w-full bg-bg rounded-full h-1.5 mt-2 overflow-hidden">
              <div
                className="h-1.5 bg-primary rounded-full transition-all"
                style={{ width: `${Math.min(slaComplianceRate, 100)}%` }}
              ></div>
            </div>
            <div className="text-[10px] text-text-muted mt-1 text-right">Target Standar &ge; 80%</div>
          </div>
        </div>
      </div>

      {/* Main Grid: Peta Sebaran + Panel Stok Insektisida */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Peta Sebaran Kasus per Kelurahan (2 Col) */}
        <div className="lg:col-span-2 bg-surface rounded-2xl border border-border shadow-sm p-5 flex flex-col">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="font-bold text-text flex items-center gap-2 text-base">
                <MapPin className="w-4 h-4 text-primary" />
                Peta Sebaran Kasus & Penanda Siaga / EWS
              </h3>
              <p className="text-xs text-text-muted">
                Visualisasi titik laporan dan zonasi risiko wilayah Kota Depok
              </p>
            </div>
            <div className="flex items-center gap-3 text-xs">
              <span className="flex items-center gap-1.5">
                <span className="w-2.5 h-2.5 rounded-full bg-red-500"></span> Siaga / Tinggi
              </span>
              <span className="flex items-center gap-1.5">
                <span className="w-2.5 h-2.5 rounded-full bg-amber-500"></span> Sedang
              </span>
              <span className="flex items-center gap-1.5">
                <span className="w-2.5 h-2.5 rounded-full bg-emerald-500"></span> Terkendali
              </span>
            </div>
          </div>
          <div className="flex-1 rounded-xl border border-border overflow-hidden isolate relative z-0 min-h-[380px]">
            <MapComponent />
          </div>
        </div>

        {/* Panel Stok Insektisida & Logistik (1 Col) */}
        <div className="bg-surface rounded-2xl border border-border shadow-sm p-5 flex flex-col justify-between">
          <div>
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-bold text-text flex items-center gap-2 text-base">
                <Package className="w-4 h-4 text-indigo-600" />
                Panel Stok Insektisida
              </h3>
              <span className="text-[11px] font-semibold text-text-muted">Logistik Sentral</span>
            </div>
            <p className="text-xs text-text-muted mb-4">
              Monitoring stok insektisida fogging di Gudang Dinkes dan 38 Puskesmas
            </p>

            {/* Total Insektisida Metric */}
            <div className="bg-bg p-4 rounded-xl border border-border mb-4">
              <div className="text-xs text-text-muted font-medium">Total Ketersediaan Insektisida</div>
              <div className="text-3xl font-black text-text mt-1">
                {totalInsecticideStock.toLocaleString('id-ID')}{" "}
                <span className="text-sm font-semibold text-text-muted">Liter</span>
              </div>
              <div className="mt-3 grid grid-cols-2 gap-2 pt-3 border-t border-border/60 text-xs">
                <div>
                  <span className="text-text-muted block text-[11px]">Gudang Utama Dinkes</span>
                  <span className="font-bold text-text">{dinkesInsecticideStock.toLocaleString('id-ID')} L</span>
                </div>
                <div>
                  <span className="text-text-muted block text-[11px]">Total di 38 Puskesmas</span>
                  <span className="font-bold text-text">{puskesmasInsecticideStock.toLocaleString('id-ID')} L</span>
                </div>
              </div>
            </div>

            {/* Alert Stok Kritis di Puskesmas */}
            <div>
              <div className="text-xs font-bold text-text mb-2 flex items-center justify-between">
                <span>Status Stok Puskesmas</span>
                {lowStockPuskesmas.length > 0 ? (
                  <span className="text-[10px] text-red-600 font-bold bg-red-50 dark:bg-red-950/40 px-2 py-0.5 rounded">
                    {lowStockPuskesmas.length} Perlu Restock
                  </span>
                ) : (
                  <span className="text-[10px] text-emerald-600 font-bold bg-emerald-50 dark:bg-emerald-950/40 px-2 py-0.5 rounded">
                    Stok Semua Aman
                  </span>
                )}
              </div>

              <div className="space-y-2 max-h-48 overflow-y-auto pr-1">
                {lowStockPuskesmas.length > 0 ? (
                  lowStockPuskesmas.map((item, idx) => (
                    <div
                      key={item.id || idx}
                      className="flex items-center justify-between p-2.5 bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-900 rounded-lg text-xs"
                    >
                      <div className="flex items-center gap-2">
                        <AlertCircle className="w-3.5 h-3.5 text-red-600 shrink-0" />
                        <span className="font-semibold text-red-800 dark:text-red-300">
                          {item.warehouse}
                        </span>
                      </div>
                      <span className="font-mono font-bold text-red-700 dark:text-red-300">
                        {item.quantity} L (&lt; 10L)
                      </span>
                    </div>
                  ))
                ) : (
                  <div className="text-xs text-text-muted text-center py-4 bg-bg rounded-lg border border-border">
                    Semua gudang Puskesmas memiliki persediaan insektisida di atas ambang batas 10 Liter.
                  </div>
                )}
              </div>
            </div>
          </div>

          <div className="pt-4 border-t border-border mt-4 text-[11px] text-text-muted flex items-center justify-between">
            <span>Ambang batas minimum fogging: <strong>10 Liter</strong></span>
            <a href="/dinkes/inventory" className="text-primary hover:underline font-bold">
              Kelola Alokasi &rarr;
            </a>
          </div>
        </div>
      </div>

      {/* Bottom Grid: Log Audit Terbaru & Daftar Laporan Terkini */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Ringkasan Log Audit Terbaru */}
        <div className="bg-surface rounded-2xl border border-border shadow-sm p-5 flex flex-col">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="font-bold text-text flex items-center gap-2 text-base">
                <History className="w-4 h-4 text-purple-600" />
                Ringkasan Log Audit Terbaru
              </h3>
              <p className="text-xs text-text-muted">
                Rekaman jejak aktivitas sistem, approval keputusan, dan pengalihan status
              </p>
            </div>
            <span className="text-[11px] text-purple-700 dark:text-purple-300 font-semibold bg-purple-50 dark:bg-purple-950/40 px-2 py-0.5 rounded">
              Audit Trail
            </span>
          </div>

          <div className="space-y-3 overflow-y-auto max-h-72 pr-1 flex-1">
            {loading ? (
              <div className="text-xs text-text-muted py-6 text-center">Memuat log audit...</div>
            ) : auditLogs.length === 0 ? (
              <div className="text-xs text-text-muted py-6 text-center">Belum ada riwayat audit terbaru.</div>
            ) : (
              auditLogs.slice(0, 6).map((log, idx) => (
                <div
                  key={log.id || idx}
                  className="flex items-start gap-3 p-3 bg-bg rounded-xl border border-border text-xs"
                >
                  <div className="w-2 h-2 mt-1.5 rounded-full bg-purple-500 shrink-0" />
                  <div className="flex-1">
                    <div className="flex items-center justify-between">
                      <span className="font-bold text-text font-mono text-[11px]">
                        {log.action}
                      </span>
                      <span className="text-[10px] text-text-muted">
                        {log.created_at ? new Date(log.created_at).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) : '-'}
                      </span>
                    </div>
                    <p className="text-[11px] text-text-muted mt-0.5">
                      Entitas: <strong className="text-text">{log.entity_name}</strong>
                      {log.user_name ? ` &bull; Oleh: ${log.user_name}` : ''}
                    </p>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Laporan Kasus Terkini */}
        <div className="bg-surface rounded-2xl border border-border shadow-sm p-5 flex flex-col">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="font-bold text-text flex items-center gap-2 text-base">
                <Building2 className="w-4 h-4 text-primary" />
                Laporan Kasus Terkini
              </h3>
              <p className="text-xs text-text-muted">
                Daftar laporan warga Kota Depok terbaru yang masuk ke sistem
              </p>
            </div>
            <a href="/dinkes/reports" className="text-xs text-primary font-bold hover:underline">
              Lihat Semua &rarr;
            </a>
          </div>

          <div className="space-y-3 overflow-y-auto max-h-72 pr-1 flex-1">
            {loading ? (
              <div className="text-xs text-text-muted py-6 text-center">Memuat laporan...</div>
            ) : filteredReports.length === 0 ? (
              <div className="text-xs text-text-muted py-6 text-center">Tidak ada laporan untuk filter terpilih.</div>
            ) : (
              filteredReports.slice(0, 6).map((r, i) => (
                <div
                  key={r.report_id || i}
                  className="flex items-center justify-between p-3 bg-bg rounded-xl border border-border text-xs"
                >
                  <div className="flex items-start gap-3">
                    <div
                      className={`w-2 h-2 mt-1.5 rounded-full shrink-0 ${
                        r.status === 'COMPLETED'
                          ? 'bg-emerald-500'
                          : r.status === 'REJECTED'
                          ? 'bg-red-500'
                          : 'bg-amber-500'
                      }`}
                    />
                    <div>
                      <p className="font-semibold text-text line-clamp-1">
                        {r.location?.address || r.description || "Laporan DBD"}
                      </p>
                      <p className="text-[11px] text-text-muted mt-0.5">
                        {r.reference_no} &bull; Kel. {r.location?.kelurahan || "-"}, Kec. {r.location?.kecamatan || "-"}
                      </p>
                    </div>
                  </div>
                  <span
                    className={`px-2 py-0.5 text-[10px] font-bold rounded-md shrink-0 ml-2 ${
                      r.status === 'COMPLETED'
                        ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300'
                        : r.status === 'REJECTED'
                        ? 'bg-red-100 text-red-800 dark:bg-red-950/60 dark:text-red-300'
                        : 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-300'
                    }`}
                  >
                    {r.status === 'COMPLETED'
                      ? 'Selesai'
                      : r.status === 'MENUNGGU_PE_KADER'
                      ? 'Antrean PE'
                      : r.status === 'MENUNGGU_VALIDASI_PUSKESMAS'
                      ? 'Validasi PE'
                      : r.status === 'PENDING'
                      ? 'Menunggu'
                      : r.status}
                  </span>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
