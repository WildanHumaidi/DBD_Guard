"use client";

import { useEffect, useState } from "react";
import api from "@/lib/api";
import { CheckSquare, AlertTriangle, FileText } from "lucide-react";
import Link from "next/link";

interface EwsWarningItem {
  created_at: string;
  location?: {
    rw?: string;
  };
}

export default function KaderDashboard() {
  const [stats, setStats] = useState({
    pending: 0,
    approved: 0,
    rejected: 0
  });
  const [ewsWarnings, setEwsWarnings] = useState<EwsWarningItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let ignore = false;
    const kelurahan = "Beji";
    const rw = "04";

    Promise.all([
      api.get("/reports"),
      api.get(`/ews-warnings?kelurahan=${kelurahan}&rw=${rw}`).catch(() => ({ data: { data: [] } }))
    ])
      .then(([repRes, ewsRes]) => {
        if (!ignore) {
          const reports = (repRes.data?.data || []) as { status?: string }[];
          let pending = 0;
          let approved = 0;
          let rejected = 0;

          reports.forEach((r) => {
            if (r.status === "MENUNGGU_PE_KADER") pending++;
            else if (r.status === "REJECTED_HOAX" || r.status === "REJECTED") rejected++;
            else approved++;
          });

          setStats({ pending, approved, rejected });
          setEwsWarnings((ewsRes.data?.data || []) as EwsWarningItem[]);
        }
      })
      .catch((err) => {
        console.error("Failed to fetch dashboard data", err);
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <div className="space-y-8 animate-fade-in">
      <div>
        <h2 className="text-3xl font-extrabold text-text tracking-tight mb-2">Beranda Verifikator</h2>
        <p className="text-text-muted">Ringkasan tugas verifikasi laporan DBD di wilayah Anda.</p>
      </div>

      {ewsWarnings.length > 0 && (
        <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded-r-lg shadow-sm">
          <div className="flex items-start">
            <div className="flex-shrink-0">
              <AlertTriangle className="h-6 w-6 text-yellow-600" />
            </div>
            <div className="ml-3">
              <h3 className="text-lg font-bold text-yellow-800">
                Peringatan Dini (EWS): Ada kasus DBD di RW Tetangga!
              </h3>
              <div className="mt-2 text-sm text-yellow-700">
                <p>
                  Sistem mendeteksi <strong>{ewsWarnings.length}</strong> kasus terverifikasi 
                  dalam 7 hari terakhir di Kelurahan yang sama, namun di RW berbeda.
                </p>
                <ul className="list-disc pl-5 mt-2 space-y-1 font-medium">
                  {ewsWarnings.slice(0, 3).map((w, idx) => (
                    <li key={idx}>RW {w.location?.rw} - {new Date(w.created_at).toLocaleDateString('id-ID')}</li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-yellow-100 p-3 rounded-xl">
            <AlertTriangle className="h-6 w-6 text-yellow-600" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">Perlu Verifikasi</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.pending}</h3>
            <p className="text-xs text-text-muted mt-1">Laporan baru (MENUNGGU_PE_KADER)</p>
          </div>
        </div>
        
        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-primary-tint p-3 rounded-xl">
            <CheckSquare className="h-6 w-6 text-primary" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">Diverifikasi Valid</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.approved}</h3>
            <p className="text-xs text-text-muted mt-1">Diteruskan ke Puskesmas</p>
          </div>
        </div>

        <div className="bg-surface rounded-2xl shadow-sm border border-border p-6 flex items-start gap-4">
          <div className="bg-red-100 p-3 rounded-xl">
            <FileText className="h-6 w-6 text-red-600" />
          </div>
          <div>
            <p className="text-sm font-medium text-text-muted mb-1">Laporan Ditolak</p>
            <h3 className="text-3xl font-bold text-text">{loading ? "..." : stats.rejected}</h3>
            <p className="text-xs text-text-muted mt-1">Ditandai sebagai HOAX</p>
          </div>
        </div>
      </div>

      <div className="bg-primary-tint rounded-2xl border border-emerald-100 p-8 flex flex-col sm:flex-row items-center justify-between gap-6">
        <div>
          <h3 className="text-xl font-bold text-primary mb-2">Ada laporan yang perlu diperiksa?</h3>
          <p className="text-primary">Pastikan untuk mengecek lokasi warga secepatnya agar proses bisa diteruskan ke Puskesmas.</p>
        </div>
        <Link 
          href="/kader/verification"
          className="px-6 py-3 bg-primary hover:bg-primary-deep text-primary-foreground font-bold rounded-xl whitespace-nowrap transition-colors shadow-lg shadow-primary/20"
        >
          Mulai Verifikasi &rarr;
        </Link>
      </div>
    </div>
  );
}
