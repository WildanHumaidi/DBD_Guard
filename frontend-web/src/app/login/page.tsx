"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { Shield, Mail, Lock, Loader2, ArrowRight, ShieldCheck, BugOff, Info } from "lucide-react";
import Link from "next/link";
import api from "@/lib/api";

export default function AdminLogin() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await api.post("/auth/login", { email, password });
      
      const data = response.data;
      
      // Save token in cookie, expires in 1 day
      Cookies.set("auth_token", data.token, { expires: 1, path: "/" });
      
      // Redirect based on role
      if (data.role === "ROLE_ADMIN_DEVELOPER") {
        router.push("/developer");
      } else if (data.role === "ROLE_DINKES") {
        router.push("/dinkes");
      } else if (data.role === "ROLE_KADER") {
        router.push("/kader");
      } else if (data.role === "ROLE_PETUGAS_PUSKESMAS" || data.role === "ROLE_KEPALA_PUSKESMAS") {
        router.push("/puskesmas");
      } else if (data.role === "ROLE_TIM_FOGGING") {
        router.push("/fogging");
      } else {
        router.push("/"); // Default fallback
      }
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || "Terjadi kesalahan saat masuk");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex w-full">
      {/* Left Panel: Branding / Visuals */}
      <div className="hidden lg:flex w-1/2 bg-primary-deep relative flex-col justify-between overflow-hidden">
        {/* Background Gradients */}
        <div className="absolute top-0 right-0 w-96 h-96 bg-primary/20 rounded-full blur-3xl -mr-20 -mt-20"></div>
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-indigo-500/20 rounded-full blur-3xl -ml-20 -mb-20"></div>
        
        <div className="p-12 relative z-10">
          <Link href="/" className="inline-flex items-center gap-2 text-white/80 hover:text-white mb-8 transition-colors bg-black/20 px-4 py-2 rounded-lg backdrop-blur-md border border-white/10 w-fit">
            <ArrowRight className="w-4 h-4 rotate-180" />
            Kembali ke Beranda
          </Link>
          
          <div className="flex items-center gap-3 mb-16">
            <div className="bg-primary p-2 rounded-xl">
              <ShieldCheck className="h-8 w-8 text-white" />
            </div>
            <span className="font-bold text-2xl text-white tracking-tight">DBD Guard</span>
          </div>

          <div className="max-w-md">
            <h1 className="text-4xl font-extrabold text-white leading-tight mb-6">
              Sistem Pengendalian <br />
              <span className="text-emerald-400">Demam Berdarah</span> Terpadu.
            </h1>
            <p className="text-slate-300 text-lg leading-relaxed">
              Portal manajemen untuk Dinas Kesehatan dan Pusat Kesehatan Masyarakat. 
              Pantau penyebaran, jadwalkan tindakan, dan respon cepat keluhan warga.
            </p>
          </div>
        </div>

        <div className="p-12 relative z-10">
          <div className="flex items-center gap-4 text-sm text-slate-400">
            <div className="flex items-center gap-2">
              <div className="w-2 h-2 rounded-full bg-primary animate-pulse"></div>
              Sistem Aktif
            </div>
            <span>&bull;</span>
            <span>Hak Cipta © {new Date().getFullYear()}</span>
          </div>
        </div>
      </div>

      {/* Right Panel: Login Form */}
      <div className="w-full lg:w-1/2 flex items-center justify-center bg-bg relative overflow-hidden">
        {/* Subtle Blob Backgrounds for Right Panel */}
        <div className="absolute top-0 right-0 w-96 h-96 bg-primary/10 rounded-full blur-3xl animate-blob pointer-events-none"></div>
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-teal-400/10 rounded-full blur-3xl animate-blob animation-delay-400 pointer-events-none"></div>

        <div className="w-full max-w-md px-8 lg:px-10 py-12 relative z-10 glass-panel rounded-3xl m-4 opacity-0 animate-fade-in-up">
          
          <div className="absolute top-6 left-6 lg:hidden">
            <Link href="/" className="inline-flex items-center justify-center gap-2 text-sm font-bold text-primary hover:text-primary-deep transition-colors bg-primary/10 px-4 py-2 rounded-lg border border-primary/20">
              <ArrowRight className="w-4 h-4 rotate-180" />
              Beranda
            </Link>
          </div>

          <div className="lg:hidden flex items-center gap-2 mb-10 mt-8 justify-center">
            <ShieldCheck className="h-8 w-8 text-primary drop-shadow-[0_0_8px_rgba(45,212,191,0.5)]" />
            <span className="font-bold text-2xl text-text tracking-tight glow-text">DBD Guard</span>
          </div>

          <div className="mb-8 text-center lg:text-left">
            <h2 className="text-3xl font-bold text-text mb-3 tracking-tight">Selamat Datang</h2>
            <p className="text-text-muted text-base">
              Silakan masuk menggunakan kredensial admin Anda untuk melanjutkan.
            </p>
          </div>

          <div className="mb-8 bg-blue-500/10 border border-blue-500/20 rounded-xl p-4 flex items-start gap-3 animate-fade-in-up animation-delay-200">
            <Info className="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-sm text-text font-medium">Hanya untuk Petugas & Admin</p>
              <p className="text-xs text-text-muted mt-1">Warga yang ingin melapor tidak perlu login. Silakan gunakan fitur pelaporan di <Link href="/" className="text-blue-500 hover:underline font-semibold">Beranda Utama</Link>.</p>
            </div>
          </div>

          {error && (
            <div className="bg-red-500/10 text-red-600 border border-red-500/20 p-4 rounded-xl text-sm mb-8 flex items-start gap-3 animate-fade-in-scale">
              <BugOff className="w-5 h-5 flex-shrink-0" />
              <span>{error}</span>
            </div>
          )}

          <form onSubmit={handleLogin} className="space-y-6">
            <div>
              <label className="block text-sm font-semibold text-text mb-2">Email Petugas</label>
              <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors group-focus-within:text-primary">
                  <Mail className="h-5 w-5 text-text-muted group-focus-within:text-primary transition-colors" />
                </div>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full pl-12 pr-4 py-3.5 bg-surface/50 border border-border rounded-xl text-text focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all font-medium backdrop-blur-sm"
                  placeholder="admin@dbdguard.id"
                  required
                />
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-semibold text-text">Kata Sandi</label>
              </div>
              <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Lock className="h-5 w-5 text-text-muted group-focus-within:text-primary transition-colors" />
                </div>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full pl-12 pr-4 py-3.5 bg-surface/50 border border-border rounded-xl text-text focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all font-medium backdrop-blur-sm shadow-sm"
                  placeholder="••••••••"
                  required
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full mt-8 bg-gradient-to-r from-primary to-teal-500 hover:from-primary-deep hover:to-primary text-white font-bold py-4 px-4 rounded-xl shadow-[0_0_20px_rgba(45,212,191,0.4)] hover:shadow-[0_0_25px_rgba(45,212,191,0.6)] hover:-translate-y-1 transition-all flex items-center justify-center gap-2 disabled:opacity-70 disabled:cursor-not-allowed group"
            >
              {loading ? (
                <Loader2 className="w-5 h-5 animate-spin" />
              ) : (
                <>
                  <span>Masuk ke Dashboard</span>
                  <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                </>
              )}
            </button>
          </form>
          
          <div className="mt-12 text-center">
            <Link href="/" className="inline-flex items-center justify-center gap-2 text-sm text-text-muted hover:text-primary font-medium transition-colors">
              <ArrowRight className="w-4 h-4 rotate-180" />
              Kembali ke Beranda Publik
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
