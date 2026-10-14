"use client";

import React, { useEffect, useState } from "react";
import api from "@/lib/api";
import { ExclamationTriangleIcon, Cog6ToothIcon, CheckCircleIcon } from "@heroicons/react/24/outline";
import ErrorMessage from "@/components/common/ErrorMessage";
import EmptyState from "@/components/common/EmptyState";

interface KlbConfig {
  config_id?: string;
  region_scope: string;
  min_cases: number;
  timeframe_days: number;
  is_auto_trigger: boolean;
}

interface KlbAlertItem {
  alert_id?: string;
  area_name: string;
  area_type: string;
  cases_count: number;
  triggered_at: string;
  is_active: boolean;
}

export default function KlbDashboard() {
  const [config, setConfig] = useState<KlbConfig | null>(null);
  const [alerts, setAlerts] = useState<KlbAlertItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Form states
  const [minCases, setMinCases] = useState<number | string>(3);
  const [timeframe, setTimeframe] = useState<number | string>(7);
  const [isAuto, setIsAuto] = useState(true);

  const loadData = async () => {
    try {
      setLoading(true);
      const [configData, alertsData] = await Promise.all([
        api.get("/admin/klb/config").then(res => res.data?.data || res.data),
        api.get("/admin/klb/alerts").then(res => res.data?.data || res.data)
      ]);
      setConfig(configData);
      setMinCases(configData.min_cases);
      setTimeframe(configData.timeframe_days);
      setIsAuto(configData.is_auto_trigger);
      setAlerts(alertsData || []);
    } catch (err: unknown) {
      const errMsg = (err as Error)?.message || "Gagal memuat data KLB";
      setError(errMsg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ignore = false;
    Promise.all([
      api.get("/admin/klb/config").then(res => res.data?.data || res.data),
      api.get("/admin/klb/alerts").then(res => res.data?.data || res.data)
    ])
      .then(([configData, alertsData]) => {
        if (!ignore) {
          setConfig(configData);
          setMinCases(configData.min_cases);
          setTimeframe(configData.timeframe_days);
          setIsAuto(configData.is_auto_trigger);
          setAlerts(alertsData || []);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          setError((err as Error)?.message || "Gagal memuat data KLB");
        }
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  const handleSaveConfig = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.patch("/admin/klb/config", {
          region_scope: config?.region_scope || "KELURAHAN",
          min_cases: Number(minCases),
          timeframe_days: Number(timeframe),
          is_auto_trigger: isAuto
      });
      alert("Konfigurasi KLB berhasil disimpan!");
      loadData();
    } catch (err: unknown) {
      const errorMsg = (err as { response?: { data?: { error?: string } }; message?: string })?.response?.data?.error || (err as Error)?.message || "Gagal memperbarui konfigurasi";
      alert("Gagal memperbarui konfigurasi: " + errorMsg);
    }
  };

  if (loading) return <div className="p-8 text-center text-text-muted">Memuat data KLB...</div>;
  if (error) return <div className="p-8 max-w-2xl mx-auto"><ErrorMessage message={error} onRetry={loadData} /></div>;

  return (
    <div className="p-8 max-w-6xl mx-auto space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-text">Manajemen KLB (Kejadian Luar Biasa)</h1>
        <p className="text-text-muted mt-1">Konfigurasi ambang batas peringatan dini dan pantau area KLB aktif.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* KLB Configuration Card */}
        <div className="bg-surface rounded-2xl shadow-sm border border-border overflow-hidden">
          <div className="p-6 border-b border-border flex items-center gap-3">
            <div className="bg-indigo-100 p-2 rounded-lg text-indigo-600">
              <Cog6ToothIcon className="w-5 h-5" />
            </div>
            <h2 className="text-lg font-bold text-text">Konfigurasi Ambang Batas</h2>
          </div>
          <form onSubmit={handleSaveConfig} className="p-6 space-y-4">
            <div>
              <label className="block text-sm font-medium text-text mb-1">Skala Wilayah KLB</label>
              <input 
                type="text" 
                value={config?.region_scope || "KELURAHAN"} 
                disabled 
                className="w-full border-border rounded-lg bg-bg text-text-muted text-sm"
              />
              <p className="text-xs text-text-muted mt-1">Level deteksi KLB (tetap Kelurahan)</p>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-text mb-1">Minimal Kasus Positif (Orang)</label>
              <input 
                type="number" 
                min="1"
                value={minCases}
                onChange={(e) => setMinCases(e.target.value === '' ? '' : parseInt(e.target.value))}
                className="w-full border-border rounded-lg focus:ring-indigo-500 focus:border-indigo-500 text-sm"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-text mb-1">Rentang Waktu Deteksi (Hari)</label>
              <input 
                type="number" 
                min="1"
                value={timeframe}
                onChange={(e) => setTimeframe(e.target.value === '' ? '' : parseInt(e.target.value))}
                className="w-full border-border rounded-lg focus:ring-indigo-500 focus:border-indigo-500 text-sm"
                required
              />
            </div>

            <div className="flex items-center gap-2 mt-4">
              <input 
                type="checkbox" 
                id="autoTrigger"
                checked={isAuto}
                onChange={(e) => setIsAuto(e.target.checked)}
                className="rounded border-border text-indigo-600 focus:ring-indigo-500"
              />
              <label htmlFor="autoTrigger" className="text-sm text-text font-medium">Aktifkan Deteksi Otomatis KLB</label>
            </div>

            <div className="pt-4">
              <button 
                type="submit" 
                className="w-full bg-indigo-600 text-white font-medium py-2 px-4 rounded-lg hover:bg-indigo-700 transition-colors"
              >
                Simpan Konfigurasi
              </button>
            </div>
          </form>
        </div>

        {/* Active Alerts Card */}
        <div className="bg-surface rounded-2xl shadow-sm border border-border overflow-hidden">
          <div className="p-6 border-b border-border flex items-center gap-3">
            <div className="bg-red-100 p-2 rounded-lg text-red-600">
              <ExclamationTriangleIcon className="w-5 h-5" />
            </div>
            <h2 className="text-lg font-bold text-text">Peringatan KLB Aktif</h2>
          </div>
          <div className="p-6">
            {alerts.length === 0 ? (
              <div className="py-8">
                <EmptyState 
                  title="Tidak ada KLB aktif saat ini" 
                  description="Situasi penyebaran DBD di bawah ambang batas."
                  icon={<CheckCircleIcon className="w-12 h-12 text-green-500" />}
                />
              </div>
            ) : (
              <div className="space-y-4">
                {alerts.map((alert) => (
                  <div key={alert.alert_id} className="bg-red-50 border border-red-200 rounded-xl p-4 flex justify-between items-center">
                    <div>
                      <span className="text-xs font-bold text-red-500 tracking-wider uppercase">{alert.area_type}</span>
                      <h3 className="text-lg font-bold text-red-900 mt-0.5">{alert.area_name}</h3>
                      <p className="text-sm text-red-700 mt-1">Terdeteksi {alert.cases_count} kasus positif dalam {config?.timeframe_days} hari terakhir.</p>
                    </div>
                    <div className="bg-red-600 text-white font-bold text-xl h-12 w-12 rounded-full flex items-center justify-center shadow-sm">
                      {alert.cases_count}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
