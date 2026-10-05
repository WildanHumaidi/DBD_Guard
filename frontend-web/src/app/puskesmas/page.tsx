"use client";

import { useEffect, useState } from "react";
import api from "@/lib/api";
import { Stethoscope, Clock, CheckCircle2, MapPin } from "lucide-react";
import Link from "next/link";

export default function PuskesmasDashboard() {
  const [stats, setStats] = useState({
    pendingPE: 0,
    completedPE: 0,
    totalAssigned: 0
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchReports = async () => {
      try {
        const response = await api.get("/reports");
        const reports = response.data?.data || [];
        
        let pendingPE = 0;
        let completedPE = 0;
        const totalAssigned = 0;

        reports.forEach((r: any) => {
          if (r.status === "MENUNGGU_VALIDASI_PUSKESMAS" || r.status === "PE_SCHEDULED") pendingPE++;
          else if (r.status === "WAITING_FOGGING_APPROVAL" || r.status === "FOGGING_SCHEDULED" || r.status === "FOGGING_IN_PROGRESS" || r.status === "COMPLETED" || r.status === "CLOSED_NO_FOGGING") completedPE++;
        });

        setStats({ pendingPE, completedPE, totalAssigned: pendingPE + completedPE });
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
        <h2 className="text-3xl font-extrabold text-text tracking-tight mb-2">Beranda Tim Puskesmas</h2>
        <p className="text-text-muted">Ringkasan aktivitas Penyelidikan Epidemiologi (PE) di wilayah Anda.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-status-verified/20 p-3 rounded-xl">
            <Clock className="h-6 w-6 text-status-verified" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">Antrean PE</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.pendingPE}</h3>
            <p className="text-xs text-text-muted mt-1">Laporan perlu diselidiki</p>
          </div>
        </div>
        
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-primary-tint p-3 rounded-xl">
            <CheckCircle2 className="h-6 w-6 text-primary" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">PE Selesai</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.completedPE}</h3>
            <p className="text-xs text-text-muted mt-1">Telah diinvestigasi</p>
          </div>
        </div>

        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-bg p-3 rounded-xl">
            <Stethoscope className="h-6 w-6 text-text-muted" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">Total Penugasan</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.totalAssigned}</h3>
            <p className="text-xs text-text-muted mt-1">Semua riwayat laporan</p>
          </div>
        </div>
      </div>

      <div className="bg-status-verified/20 rounded-2xl border border-blue-100 p-8 flex flex-col sm:flex-row items-center justify-between gap-6 relative overflow-hidden">
        <div className="absolute top-0 right-0 p-8 opacity-10">
          <Stethoscope className="w-48 h-48 text-blue-900" />
        </div>
        <div className="relative z-10">
          <h3 className="text-xl font-bold text-blue-900 mb-2">Mulai Investigasi Lapangan</h3>
          <p className="text-blue-700 max-w-lg">Cek Angka Bebas Jentik (ABJ) dan identifikasi jumlah penderita tambahan di sekitar lokasi laporan warga.</p>
        </div>
        <Link 
          href="/puskesmas/pe-validation"
          className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold rounded-xl whitespace-nowrap transition-colors shadow-lg shadow-blue-600/20 relative z-10"
        >
          Lihat Daftar PE &rarr;
        </Link>
      </div>
    </div>
  );
}
