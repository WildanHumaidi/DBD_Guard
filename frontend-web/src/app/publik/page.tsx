"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import Header from "@/components/layout/Header";
import { FloatingThemeToggle } from "@/components/FloatingThemeToggle";
import {
  BarChart3,
  CheckCircle2,
  Clock,
  XCircle,
  Star,
  TrendingUp,
  Activity,
  MapPin,
  ArrowRight,
  Shield,
} from "lucide-react";

interface PublicStats {
  total_reports: number;
  total_completed: number;
  total_in_progress: number;
  total_rejected: number;
  avg_feedback: number;
  feedback_count: number;
}

export default function PublikPage() {
  const [stats, setStats] = useState<PublicStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/v1/public/stats");
        const data = await response.json();

        if (data.status === "success") {
          setStats(data.data);
        } else {
          setError(data.message || "Gagal mengambil data statistik");
        }
      } catch (err) {
        setError("Terjadi kesalahan jaringan");
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  const completionRate =
    stats && stats.total_reports > 0
      ? ((stats.total_completed / stats.total_reports) * 100).toFixed(1)
      : "0";

  return (
    <main className="min-h-screen bg-bg font-sans">
      <FloatingThemeToggle />
      <Header />

      {/* Hero Section */}
      <section className="relative bg-gradient-to-br from-primary-deep via-primary to-teal-800 text-white pt-16 pb-32 px-4 sm:px-6 overflow-hidden">
        <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-white/10 rounded-full blur-[100px] -translate-y-1/2 translate-x-1/4 animate-blob"></div>
        <div className="absolute bottom-0 left-0 w-[400px] h-[400px] bg-teal-300/20 rounded-full blur-[80px] translate-y-1/2 -translate-x-1/4 animate-blob animation-delay-400"></div>

        <div className="max-w-5xl mx-auto text-center relative z-10 opacity-0 animate-fade-in-up">
          <div className="inline-flex items-center rounded-full border border-white/30 bg-white/10 backdrop-blur-md px-4 py-1.5 text-sm font-semibold text-white mb-8 shadow-lg">
            <BarChart3 className="h-4 w-4 mr-2" />
            Dashboard Transparansi Publik
          </div>
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold tracking-tight mb-6 leading-[1.1]">
            Data Terbuka{" "}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-teal-100 to-white drop-shadow-[0_0_15px_rgba(255,255,255,0.4)]">
              Penanganan DBD
            </span>
          </h1>
          <p className="text-lg sm:text-xl text-teal-50 max-w-2xl mx-auto leading-relaxed">
            Informasi real-time mengenai seluruh aktivitas pelaporan dan
            penanganan Demam Berdarah Dengue (DBD) di Kota Depok. Data ini
            bersifat terbuka sebagai bentuk transparansi pemerintah.
          </p>
        </div>
      </section>

      {/* Stats Cards */}
      <section className="max-w-6xl mx-auto px-4 sm:px-6 -mt-20 relative z-10 pb-16">
        {loading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {[...Array(4)].map((_, i) => (
              <div
                key={i}
                className="bg-surface rounded-2xl shadow-xl p-8 animate-pulse"
              >
                <div className="h-4 bg-border rounded w-2/3 mb-4"></div>
                <div className="h-10 bg-border rounded w-1/2 mb-2"></div>
                <div className="h-3 bg-bg rounded w-1/3"></div>
              </div>
            ))}
          </div>
        ) : error ? (
          <div className="bg-surface rounded-2xl shadow-xl p-10 text-center">
            <XCircle className="h-12 w-12 text-red-400 mx-auto mb-4" />
            <p className="text-text-muted text-lg">{error}</p>
            <button
              onClick={() => window.location.reload()}
              className="mt-4 px-6 py-3 bg-primary text-primary-foreground rounded-xl font-bold hover:bg-primary-deep transition-all"
            >
              Coba Lagi
            </button>
          </div>
        ) : (
          stats && (
            <>
              {/* Main Stat Cards */}
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8 opacity-0 animate-fade-in-up animation-delay-200">
                {/* Total Laporan */}
                <div className="group bg-surface border border-border shadow-lg rounded-2xl p-8 hover:-translate-y-2 transition-all duration-300 hover:shadow-xl">
                  <div className="flex items-center justify-between mb-6">
                    <div className="h-12 w-12 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-[0_0_15px_rgba(59,130,246,0.5)] group-hover:rotate-6 group-hover:scale-110 transition-transform duration-300">
                      <Activity className="h-6 w-6 text-white" />
                    </div>
                    <span className="text-xs font-bold text-status-verified bg-status-verified/20 px-3 py-1 rounded-full border border-status-verified/30">
                      Total
                    </span>
                  </div>
                  <p className="text-4xl font-extrabold text-text mb-1 drop-shadow-[0_0_8px_rgba(255,255,255,0.3)]">
                    {stats.total_reports.toLocaleString("id-ID")}
                  </p>
                  <p className="text-sm font-medium text-text-muted">
                    Laporan Masuk
                  </p>
                </div>

                {/* Selesai Ditangani */}
                <div className="group bg-surface border border-border shadow-lg rounded-2xl p-8 hover:-translate-y-2 transition-all duration-300 hover:shadow-xl">
                  <div className="flex items-center justify-between mb-6">
                    <div className="h-12 w-12 bg-gradient-to-br from-emerald-500 to-teal-600 rounded-xl flex items-center justify-center shadow-[0_0_15px_rgba(16,185,129,0.5)] group-hover:rotate-6 group-hover:scale-110 transition-transform duration-300">
                      <CheckCircle2 className="h-6 w-6 text-white" />
                    </div>
                    <span className="text-xs font-bold text-primary bg-primary-tint px-3 py-1 rounded-full border border-primary/30">
                      Selesai
                    </span>
                  </div>
                  <p className="text-4xl font-extrabold glow-text mb-1">
                    {stats.total_completed.toLocaleString("id-ID")}
                  </p>
                  <p className="text-sm font-medium text-text-muted">
                    Kasus Tertangani
                  </p>
                </div>

                {/* Sedang Diproses */}
                <div className="group bg-surface border border-border shadow-lg rounded-2xl p-8 hover:-translate-y-2 transition-all duration-300 hover:shadow-xl">
                  <div className="flex items-center justify-between mb-6">
                    <div className="h-12 w-12 bg-gradient-to-br from-amber-400 to-orange-500 rounded-xl flex items-center justify-center shadow-[0_0_15px_rgba(245,158,11,0.5)] group-hover:-rotate-6 group-hover:scale-110 transition-transform duration-300">
                      <Clock className="h-6 w-6 text-white" />
                    </div>
                    <span className="text-xs font-bold text-amber-600 bg-amber-500/10 px-3 py-1 rounded-full border border-amber-500/30">
                      Proses
                    </span>
                  </div>
                  <p className="text-4xl font-extrabold text-amber-500 mb-1 drop-shadow-[0_0_8px_rgba(245,158,11,0.4)]">
                    {stats.total_in_progress.toLocaleString("id-ID")}
                  </p>
                  <p className="text-sm font-medium text-text-muted">
                    Sedang Diproses
                  </p>
                </div>

                {/* Tingkat Penyelesaian */}
                <div className="group bg-surface border border-border shadow-lg rounded-2xl p-8 hover:-translate-y-2 transition-all duration-300 hover:shadow-xl">
                  <div className="flex items-center justify-between mb-6">
                    <div className="h-12 w-12 bg-gradient-to-br from-violet-500 to-purple-600 rounded-xl flex items-center justify-center shadow-[0_0_15px_rgba(139,92,246,0.5)] group-hover:rotate-6 group-hover:scale-110 transition-transform duration-300">
                      <TrendingUp className="h-6 w-6 text-white" />
                    </div>
                    <span className="text-xs font-bold text-violet-600 bg-violet-500/10 px-3 py-1 rounded-full border border-violet-500/30">
                      Rasio
                    </span>
                  </div>
                  <p className="text-4xl font-extrabold text-violet-500 mb-1 drop-shadow-[0_0_8px_rgba(139,92,246,0.4)]">
                    {completionRate}%
                  </p>
                  <p className="text-sm font-medium text-text-muted">
                    Tingkat Penyelesaian
                  </p>
                </div>
              </div>

              {/* Secondary Info Row */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 opacity-0 animate-fade-in-up animation-delay-400">
                {/* Rating Kepuasan */}
                <div className="bg-surface border border-border shadow-lg rounded-2xl p-8 hover:-translate-y-1 transition-all duration-300 hover:shadow-xl">
                  <div className="flex items-center gap-4 mb-6">
                    <div className="h-12 w-12 bg-amber-500/20 rounded-xl flex items-center justify-center border border-amber-500/30 shadow-[0_0_15px_rgba(245,158,11,0.2)]">
                      <Star className="h-6 w-6 text-amber-500" />
                    </div>
                    <div>
                      <h3 className="font-bold text-text text-lg">
                        Kepuasan Warga
                      </h3>
                      <p className="text-sm text-text-muted">
                        Rata-rata penilaian pasca-penanganan
                      </p>
                    </div>
                  </div>
                  <div className="flex items-end gap-4">
                    <p className="text-5xl font-extrabold text-text drop-shadow-md">
                      {stats.avg_feedback > 0
                        ? stats.avg_feedback.toFixed(1)
                        : "—"}
                    </p>
                    <div className="pb-2">
                      <div className="flex gap-1 mb-1">
                        {[1, 2, 3, 4, 5].map((s) => (
                          <Star
                            key={s}
                            className={`h-5 w-5 ${
                              s <= Math.round(stats.avg_feedback)
                                ? "text-amber-400 fill-amber-400 drop-shadow-[0_0_5px_rgba(251,191,36,0.5)]"
                                : "text-border"
                            }`}
                          />
                        ))}
                      </div>
                      <p className="text-xs text-text-muted">
                        dari {stats.feedback_count} penilaian
                      </p>
                    </div>
                  </div>
                </div>

                {/* Quick Links */}
                <div className="bg-gradient-to-br from-primary-deep to-slate-900 rounded-2xl shadow-[0_0_20px_rgba(15,118,110,0.3)] p-8 text-white relative overflow-hidden group">
                  <div className="absolute top-0 right-0 w-32 h-32 bg-primary/30 rounded-full blur-2xl group-hover:scale-150 transition-transform duration-700 pointer-events-none"></div>
                  <h3 className="font-bold text-xl mb-3 relative z-10">
                    Ikut Berkontribusi!
                  </h3>
                  <p className="text-slate-300 text-sm mb-8 leading-relaxed relative z-10">
                    Laporkan temuan kasus DBD, jentik nyamuk, atau genangan air
                    di lingkungan Anda agar segera ditindaklanjuti oleh Puskesmas
                    terdekat.
                  </p>
                  <div className="flex flex-col sm:flex-row gap-3 relative z-10">
                    <Link
                      href="/report"
                      className="group/btn inline-flex items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-primary to-teal-400 px-6 py-3 text-sm font-bold text-white shadow-[0_0_15px_rgba(45,212,191,0.5)] hover:shadow-[0_0_20px_rgba(45,212,191,0.7)] hover:-translate-y-1 transition-all"
                    >
                      Buat Laporan
                      <ArrowRight className="h-4 w-4 group-hover/btn:translate-x-1 transition-transform" />
                    </Link>
                    <Link
                      href="/track"
                      className="inline-flex items-center justify-center gap-2 rounded-xl bg-white/10 backdrop-blur-md border border-white/20 px-6 py-3 text-sm font-bold text-white hover:bg-white/20 hover:-translate-y-1 transition-all"
                    >
                      <MapPin className="h-4 w-4" />
                      Lacak Status
                    </Link>
                  </div>
                </div>
              </div>
            </>
          )
        )}
      </section>

      {/* Footer */}
      <footer className="border-t border-border/80 bg-surface/80 backdrop-blur-md py-10 mt-auto">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col md:flex-row justify-between items-center gap-4">
          <div className="flex items-center gap-3">
            <div className="bg-bg p-2 rounded-lg">
              <Shield className="h-5 w-5 text-text-muted" />
            </div>
            <span className="font-bold text-text">DBD Guard</span>
          </div>
          <p className="text-text-muted text-sm font-medium">
            &copy; {new Date().getFullYear()} Dinas Kesehatan Kota Depok.
            Dashboard Transparansi Publik.
          </p>
        </div>
      </footer>
    </main>
  );
}
