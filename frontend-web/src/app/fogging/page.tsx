"use client";

import { useEffect, useState } from "react";
import api from "@/lib/api";
import { Droplet, Clock, CheckCircle2, AlertTriangle } from "lucide-react";
import Link from "next/link";

export default function FoggingDashboard() {
  const [stats, setStats] = useState({
    scheduled: 0,
    inProgress: 0,
    completed: 0
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchReports = async () => {
      try {
        const response = await api.get("/reports");
        const reports = response.data?.data || [];
        
        let scheduled = 0;
        let inProgress = 0;
        let completed = 0;

        reports.forEach((r: any) => {
          if (r.status === "FOGGING_SCHEDULED") scheduled++;
          else if (r.status === "FOGGING_IN_PROGRESS") inProgress++;
          else if (r.status === "COMPLETED") completed++;
        });

        setStats({ scheduled, inProgress, completed });
      } catch (err) {
        console.error("Failed to fetch reports for stats", err);
      } finally {
        setLoading(false);
      }
    };
    fetchReports();
  }, []);

  return (
    <div className="space-y-8 animate-fade-in">
      <div>
        <h2 className="text-3xl font-extrabold text-text tracking-tight mb-2">Beranda Tim Fogging</h2>
        <p className="text-text-muted">Ringkasan penugasan fogging di wilayah Anda.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-status-scheduled/20 p-3 rounded-xl">
            <Clock className="h-6 w-6 text-status-scheduled" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">Terjadwal</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.scheduled}</h3>
            <p className="text-xs text-text-muted mt-1">Siap untuk dieksekusi</p>
          </div>
        </div>

        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-status-verified/20 p-3 rounded-xl">
            <AlertTriangle className="h-6 w-6 text-status-verified" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">Sedang Berjalan</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.inProgress}</h3>
            <p className="text-xs text-text-muted mt-1">Dalam proses fogging</p>
          </div>
        </div>
        
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-primary-tint p-3 rounded-xl">
            <CheckCircle2 className="h-6 w-6 text-primary" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">Selesai</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.completed}</h3>
            <p className="text-xs text-text-muted mt-1">Fogging berhasil</p>
          </div>
        </div>
      </div>

      <div className="bg-status-scheduled/20 rounded-2xl border border-orange-100 p-8 flex flex-col sm:flex-row items-center justify-between gap-6 relative overflow-hidden">
        <div className="absolute top-0 right-0 p-8 opacity-10">
          <Droplet className="w-48 h-48 text-orange-900" />
        </div>
        <div className="relative z-10">
          <h3 className="text-xl font-bold text-orange-900 mb-2">Lihat Jadwal Hari Ini</h3>
          <p className="text-status-scheduled max-w-lg">Pastikan bahan dan alat siap sebelum berangkat. Laporkan status eksekusi langsung dari lokasi.</p>
        </div>
        <Link 
          href="/fogging/schedules"
          className="px-6 py-3 bg-orange-600 hover:bg-orange-700 text-white font-bold rounded-xl whitespace-nowrap transition-colors shadow-lg shadow-orange-600/20 relative z-10"
        >
          Lihat Jadwal &rarr;
        </Link>
      </div>
    </div>
  );
}
