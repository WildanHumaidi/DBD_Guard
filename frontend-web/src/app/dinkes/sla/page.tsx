"use client";

import { useEffect, useState } from "react";
import { AlertTriangle, Clock, RefreshCw, CheckCircle } from "lucide-react";
import api from "@/lib/api";
import ErrorMessage from "@/components/common/ErrorMessage";

interface SLAMetrics {
  avg_pe_response_hours: number;
  avg_fogging_response_hours: number;
  sla_pe_compliance_rate: number;
  sla_fogging_compliance_rate: number;
}

export default function SLADashboardPage() {
  const [metrics, setMetrics] = useState<SLAMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchMetrics = async () => {
    setLoading(true);
    setError("");
    try {
      const res = await api.get("/admin/analytics/sla");
      if (res.data?.data) {
        setMetrics(res.data.data);
      }
    } catch (err: unknown) {
      const errorMsg = (err as { response?: { data?: { error?: string } }; message?: string })?.response?.data?.error || (err as Error)?.message || "Gagal memuat data metrik SLA.";
      setError(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ignore = false;
    api.get("/admin/analytics/sla")
      .then((res) => {
        if (!ignore && res.data?.data) {
          setMetrics(res.data.data);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          const errorMsg = (err as { response?: { data?: { error?: string } }; message?: string })?.response?.data?.error || (err as Error)?.message || "Gagal memuat data metrik SLA.";
          setError(errorMsg);
        }
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <div className="h-full flex flex-col gap-6 p-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text flex items-center gap-3">
            <div className="h-10 w-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg shadow-indigo-500/25">
              <Clock className="h-5 w-5 text-white" />
            </div>
            Monitoring SLA (Dinas Kesehatan)
          </h1>
          <p className="text-sm text-text-muted mt-1 ml-[52px]">
            Pantau kepatuhan standar waktu respons Puskesmas dan Kader untuk Penyelidikan Epidemiologi (1x24 Jam) dan Fogging (3x24 Jam).
          </p>
        </div>
        <button
          onClick={fetchMetrics}
          disabled={loading}
          className="inline-flex items-center gap-2 px-5 py-2.5 bg-surface border border-border rounded-xl text-sm font-bold text-text hover:bg-bg hover:border-border transition-all disabled:opacity-50 shadow-sm"
        >
          <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
          Refresh Data
        </button>
      </div>

      {error && (
        <ErrorMessage message={error} onRetry={fetchMetrics} />
      )}

      {loading && !metrics ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="bg-surface p-6 rounded-2xl shadow-sm border border-border h-32 animate-pulse flex flex-col justify-between">
              <div className="h-4 bg-border rounded w-1/2"></div>
              <div className="h-8 bg-border rounded w-1/3"></div>
            </div>
          ))}
        </div>
      ) : metrics ? (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {/* Avg PE Response */}
            <div className="bg-surface p-6 rounded-2xl shadow-sm border border-border flex flex-col justify-between">
              <div className="flex justify-between items-start mb-4">
                <p className="text-sm font-medium text-text-muted">Rata-rata Waktu PE (Kader)</p>
                <div className="p-2 bg-status-verified/20 text-status-verified rounded-lg">
                  <Clock className="h-5 w-5" />
                </div>
              </div>
              <div>
                <p className="text-3xl font-bold text-text">
                  {metrics.avg_pe_response_hours.toFixed(1)} <span className="text-lg font-medium text-text-muted">Jam</span>
                </p>
                <p className={`text-xs mt-2 font-medium ${metrics.avg_pe_response_hours <= 24 ? "text-primary" : "text-red-600"}`}>
                  Target SLA: &le; 24 Jam
                </p>
              </div>
            </div>

            {/* PE Compliance Rate */}
            <div className="bg-surface p-6 rounded-2xl shadow-sm border border-border flex flex-col justify-between">
              <div className="flex justify-between items-start mb-4">
                <p className="text-sm font-medium text-text-muted">Tingkat Kepatuhan PE</p>
                <div className={`p-2 rounded-lg ${metrics.sla_pe_compliance_rate >= 80 ? 'bg-primary-tint text-primary' : 'bg-amber-50 text-amber-600'}`}>
                  {metrics.sla_pe_compliance_rate >= 80 ? <CheckCircle className="h-5 w-5" /> : <AlertTriangle className="h-5 w-5" />}
                </div>
              </div>
              <div>
                <p className="text-3xl font-bold text-text">
                  {metrics.sla_pe_compliance_rate.toFixed(1)}%
                </p>
                <div className="w-full bg-bg rounded-full h-1.5 mt-3 overflow-hidden">
                  <div className={`h-1.5 rounded-full ${metrics.sla_pe_compliance_rate >= 80 ? 'bg-primary' : 'bg-amber-500'}`} style={{ width: `${metrics.sla_pe_compliance_rate}%` }}></div>
                </div>
              </div>
            </div>

            {/* Avg Fogging Response */}
            <div className="bg-surface p-6 rounded-2xl shadow-sm border border-border flex flex-col justify-between">
              <div className="flex justify-between items-start mb-4">
                <p className="text-sm font-medium text-text-muted">Rata-rata Waktu Fogging</p>
                <div className="p-2 bg-purple-50 text-purple-600 rounded-lg">
                  <Clock className="h-5 w-5" />
                </div>
              </div>
              <div>
                <p className="text-3xl font-bold text-text">
                  {metrics.avg_fogging_response_hours.toFixed(1)} <span className="text-lg font-medium text-text-muted">Jam</span>
                </p>
                <p className={`text-xs mt-2 font-medium ${metrics.avg_fogging_response_hours <= 72 ? "text-primary" : "text-red-600"}`}>
                  Target SLA: &le; 72 Jam
                </p>
              </div>
            </div>

            {/* Fogging Compliance Rate */}
            <div className="bg-surface p-6 rounded-2xl shadow-sm border border-border flex flex-col justify-between">
              <div className="flex justify-between items-start mb-4">
                <p className="text-sm font-medium text-text-muted">Tingkat Kepatuhan Fogging</p>
                <div className={`p-2 rounded-lg ${metrics.sla_fogging_compliance_rate >= 80 ? 'bg-primary-tint text-primary' : 'bg-amber-50 text-amber-600'}`}>
                  {metrics.sla_fogging_compliance_rate >= 80 ? <CheckCircle className="h-5 w-5" /> : <AlertTriangle className="h-5 w-5" />}
                </div>
              </div>
              <div>
                <p className="text-3xl font-bold text-text">
                  {metrics.sla_fogging_compliance_rate.toFixed(1)}%
                </p>
                <div className="w-full bg-bg rounded-full h-1.5 mt-3 overflow-hidden">
                  <div className={`h-1.5 rounded-full ${metrics.sla_fogging_compliance_rate >= 80 ? 'bg-primary' : 'bg-amber-500'}`} style={{ width: `${metrics.sla_fogging_compliance_rate}%` }}></div>
                </div>
              </div>
            </div>
          </div>
          
          {/* Information Notice */}
          <div className="bg-status-verified/20/50 border border-blue-100 rounded-2xl p-6 mt-6">
            <h3 className="text-sm font-bold text-blue-900 mb-2 flex items-center gap-2">
              <AlertTriangle className="w-4 h-4 text-status-verified" />
              Sistem Auto-Release Aktif
            </h3>
            <p className="text-sm text-blue-700 leading-relaxed">
              Cronjob Auto-Release berjalan setiap jam di *background*. Jika seorang Kader menahan tugas Penyelidikan Epidemiologi (PE) tanpa menyelesaikannya dalam waktu lebih dari 24 jam, sistem akan **otomatis melepas tugas tersebut** agar kembali ke *pool* dan dapat diambil oleh Kader lainnya.
            </p>
          </div>
        </>
      ) : null}
    </div>
  );
}
