"use client";

import { Users, Settings, Activity, Server, Clock } from "lucide-react";
import Link from "next/link";
import { useEffect, useState } from "react";
import api from "@/lib/api";

interface AuditLogSummary {
  id?: string;
  created_at: string;
  action: string;
  resource?: string;
  entity_type?: string;
  user?: {
    name?: string;
    email?: string;
  };
}

export default function DeveloperDashboard() {
  const [logs, setLogs] = useState<AuditLogSummary[]>([]);
  const [loadingLogs, setLoadingLogs] = useState(true);

  useEffect(() => {
    let ignore = false;
    api.get("/developer/audit-logs")
      .then((res) => {
        if (!ignore && Array.isArray(res.data)) {
          setLogs(res.data);
        }
      })
      .catch((error) => {
        console.error("Failed to fetch audit logs:", error);
      })
      .finally(() => {
        if (!ignore) setLoadingLogs(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <div className="w-full max-w-7xl mx-auto space-y-8">
      <div>
        <h2 className="text-2xl font-bold text-text">Beranda Developer</h2>
        <p className="text-text-muted">Kelola konfigurasi sistem, peran pengguna, dan fungsionalitas inti web.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-text-muted">Total Akun</h3>
            <div className="p-2 bg-status-verified/20 text-status-verified rounded-lg">
              <Users className="h-5 w-5" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold text-text">12</span>
            <span className="text-sm text-text-muted">Pengguna</span>
          </div>
        </div>

        <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-text-muted">Status API</h3>
            <div className="p-2 bg-primary-tint text-primary rounded-lg">
              <Server className="h-5 w-5" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold text-primary">Online</span>
            <span className="text-sm text-text-muted">99.9% Uptime</span>
          </div>
        </div>

        <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-text-muted">Beban Sistem</h3>
            <div className="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
              <Activity className="h-5 w-5" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold text-text">Normal</span>
          </div>
        </div>

        <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-text-muted">Versi Web</h3>
            <div className="p-2 bg-amber-50 text-amber-600 rounded-lg">
              <Settings className="h-5 w-5" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-bold text-text">v1.0.0</span>
            <span className="text-sm text-text-muted">Stable</span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-surface rounded-2xl border border-border shadow-sm p-6">
          <h3 className="font-bold text-text mb-4">Akses Cepat</h3>
          <div className="space-y-4">
            <Link href="/developer/users" className="block p-4 border border-border rounded-xl hover:border-primary hover:bg-primary-tint/50 transition-all group">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-bg rounded-lg group-hover:bg-primary-tint transition-colors">
                  <Users className="w-5 h-5 text-text-muted group-hover:text-primary" />
                </div>
                <div>
                  <h4 className="font-semibold text-text">Kelola Akun (Users & Roles)</h4>
                  <p className="text-sm text-text-muted">Tambah, ubah, atau hapus akun pengguna sistem.</p>
                </div>
              </div>
            </Link>

            <Link href="/developer/settings" className="block p-4 border border-border rounded-xl hover:border-primary hover:bg-primary-tint/50 transition-all group">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-bg rounded-lg group-hover:bg-primary-tint transition-colors">
                  <Settings className="w-5 h-5 text-text-muted group-hover:text-primary" />
                </div>
                <div>
                  <h4 className="font-semibold text-text">Pengaturan Website</h4>
                  <p className="text-sm text-text-muted">Ubah banner beranda, kebijakan privasi, dan info statis.</p>
                </div>
              </div>
            </Link>
          </div>
        </div>

        <div className="bg-surface rounded-2xl border border-border shadow-sm p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-bold text-text">Log Aktivitas Sistem Terbaru</h3>
            <Link href="/developer/audit-logs" className="text-sm text-primary hover:text-primary font-medium hidden">
              Lihat Semua
            </Link>
          </div>
          
          <div className="space-y-0">
            {loadingLogs ? (
              <div className="py-8 text-center text-text-muted">Memuat log aktivitas...</div>
            ) : logs.length === 0 ? (
              <div className="py-8 text-center text-text-muted flex flex-col items-center">
                <Clock className="w-8 h-8 text-text-muted/50 mb-2" />
                <p>Belum ada aktivitas yang tercatat</p>
              </div>
            ) : (
              logs.map((log, index) => (
                <div key={log.id || index} className={`flex gap-4 ${index !== logs.length - 1 ? 'border-b border-border pb-4 mb-4' : ''}`}>
                  <div className="w-2 h-2 mt-1.5 rounded-full bg-primary shrink-0" />
                  <div>
                    <p className="text-sm font-medium text-text">
                      {log.user?.name || "System"} melakukan {log.action.replace(/_/g, " ")} pada {log.resource}
                    </p>
                    <p className="text-xs text-text-muted mt-1">
                      {new Date(log.created_at).toLocaleString("id-ID", {
                        dateStyle: "medium",
                        timeStyle: "short",
                      })}
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
