"use client";

import { Activity, CheckCircle2, Server, Database, Shield, XCircle } from "lucide-react";
import { useState, useEffect } from "react";
import api from "@/lib/api";

export default function SystemStatusPage() {
  const [status, setStatus] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchStatus();
  }, []);

  const fetchStatus = async () => {
    try {
      const res = await api.get("/developer/status");
      setStatus(res.data);
    } catch (err) {
      console.error("Gagal mengambil status", err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="text-center py-10 text-text-muted">Memuat status sistem...</div>;
  }

  const allOperational = status?.database === "Operational" && status?.web_server === "Operational";

  return (
    <div className="w-full max-w-6xl mx-auto space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-text">System Status</h2>
        <p className="text-text-muted">Pantau kesehatan dan ketersediaan layanan sistem.</p>
      </div>

      <div className={`border rounded-2xl p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4 ${allOperational ? 'bg-primary-tint border-primary' : 'bg-red-50 border-red-200'}`}>
        <div className="flex items-center gap-4">
          <div className={`w-12 h-12 rounded-full flex items-center justify-center ${allOperational ? 'bg-primary-tint text-primary' : 'bg-red-100 text-red-600'}`}>
            {allOperational ? <CheckCircle2 className="w-6 h-6" /> : <XCircle className="w-6 h-6" />}
          </div>
          <div>
            <h3 className={`font-bold text-lg ${allOperational ? 'text-primary' : 'text-red-900'}`}>
              {allOperational ? "Semua Sistem Berjalan Normal" : "Beberapa Sistem Mengalami Gangguan"}
            </h3>
            <p className={allOperational ? "text-primary text-sm" : "text-red-700 text-sm"}>
              Server Memory: {status?.memory_mb || 0} MB
            </p>
          </div>
        </div>
        <div className="text-right">
          <p className={`text-sm font-semibold ${allOperational ? 'text-primary' : 'text-red-800'}`}>Uptime Backend</p>
          <p className={`text-xl font-bold ${allOperational ? 'text-primary' : 'text-red-600'}`}>{status?.uptime || "0s"}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-surface rounded-2xl border border-border p-6 shadow-sm">
          <div className="flex items-center gap-3 mb-4">
            <div className="p-2 bg-status-verified/20 text-status-verified rounded-lg">
              <Server className="w-5 h-5" />
            </div>
            <h3 className="font-semibold text-text">Web Server (Backend)</h3>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-text-muted">Status</span>
            <span className={`px-3 py-1 rounded-full text-xs font-semibold ${status?.web_server === 'Operational' ? 'bg-primary-tint text-primary' : 'bg-red-100 text-red-700'}`}>
              {status?.web_server || "Unknown"}
            </span>
          </div>
        </div>

        <div className="bg-surface rounded-2xl border border-border p-6 shadow-sm">
          <div className="flex items-center gap-3 mb-4">
            <div className="p-2 bg-purple-50 text-purple-600 rounded-lg">
              <Database className="w-5 h-5" />
            </div>
            <h3 className="font-semibold text-text">API & Database (Backend)</h3>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-text-muted">Status</span>
            <span className={`px-3 py-1 rounded-full text-xs font-semibold ${status?.database === 'Operational' ? 'bg-primary-tint text-primary' : 'bg-red-100 text-red-700'}`}>
              {status?.database || "Unknown"}
            </span>
          </div>
        </div>

        <div className="bg-surface rounded-2xl border border-border p-6 shadow-sm">
          <div className="flex items-center gap-3 mb-4">
            <div className="p-2 bg-amber-50 text-amber-600 rounded-lg">
              <Shield className="w-5 h-5" />
            </div>
            <h3 className="font-semibold text-text">Layanan Notifikasi (FCM)</h3>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-text-muted">Status</span>
            <span className={`px-3 py-1 rounded-full text-xs font-semibold ${status?.fcm_service === 'Operational' ? 'bg-primary-tint text-primary' : 'bg-red-100 text-red-700'}`}>
              {status?.fcm_service || "Unknown"}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
