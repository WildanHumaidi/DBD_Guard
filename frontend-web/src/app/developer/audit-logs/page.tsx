"use client";

import { useEffect, useState } from "react";
import { Clock, RefreshCw } from "lucide-react";
import api from "@/lib/api";

interface AuditLogItem {
  id?: string;
  created_at: string;
  action: string;
  resource?: string;
  entity_type?: string;
  entity_name?: string;
  user?: {
    name?: string;
    email?: string;
  };
  details?: Record<string, unknown> | string;
}

export default function AuditLogsPage() {
  const [logs, setLogs] = useState<AuditLogItem[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchLogs = async () => {
    setLoading(true);
    try {
      const res = await api.get("/developer/audit-logs");
      if (Array.isArray(res.data)) {
        setLogs(res.data);
      }
    } catch (error) {
      console.error("Failed to fetch audit logs:", error);
    } finally {
      setLoading(false);
    }
  };

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
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <div className="w-full max-w-7xl mx-auto space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-text">Log Aktivitas Sistem</h2>
          <p className="text-text-muted">Pantau secara detail semua aktivitas dan perubahan yang terjadi pada sistem.</p>
        </div>
        <div className="flex items-center gap-3">
          <button 
            onClick={fetchLogs}
            className="p-2.5 bg-surface border border-border rounded-xl text-text-muted hover:text-primary transition-colors"
          >
            <RefreshCw className={`w-5 h-5 ${loading ? "animate-spin" : ""}`} />
          </button>
        </div>
      </div>

      <div className="bg-surface rounded-2xl border border-border overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-border">
            <thead className="bg-bg">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Waktu</th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Aktor (User)</th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Aksi</th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Target Resources</th>
              </tr>
            </thead>
            <tbody className="bg-surface divide-y divide-border">
              {loading ? (
                <tr>
                  <td colSpan={4} className="px-6 py-8 text-center text-text-muted">Memuat data log aktivitas...</td>
                </tr>
              ) : logs.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-6 py-12 text-center text-text-muted">
                    <div className="flex flex-col items-center">
                      <Clock className="w-10 h-10 text-text-muted/50 mb-3" />
                      <p>Belum ada aktivitas yang tercatat</p>
                    </div>
                  </td>
                </tr>
              ) : (
                logs.map((log, idx) => (
                  <tr key={log.id || `log-${idx}`} className="hover:bg-bg/50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-text-muted">
                      {new Date(log.created_at).toLocaleString("id-ID", {
                        dateStyle: "medium",
                        timeStyle: "medium",
                      })}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-semibold text-text">{log.user?.name || "System"}</div>
                      {log.user?.email && (
                        <div className="text-xs text-text-muted">{log.user.email}</div>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-primary-tint text-primary">
                        {log.action.replace(/_/g, " ")}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-text">
                      {log.resource}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
