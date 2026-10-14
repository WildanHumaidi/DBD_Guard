"use client";

import Link from 'next/link';
import Image from 'next/image';
import { ArrowRight, Shield, MapPin, Stethoscope, Droplets, Trash2, Home as HomeIcon, AlertCircle, CheckCircle2, Clock, Star, Activity } from 'lucide-react';
import Header from '@/components/layout/Header';
import Footer from '@/components/layout/Footer';
import { useEffect, useState } from 'react';
import api from '@/lib/api';

interface PublicStats {
  total_reports: number;
  total_completed: number;
  total_in_progress: number;
  total_rejected: number;
  avg_feedback: number;
  feedback_count: number;
}

export default function Home() {
  const [heroTitle, setHeroTitle] = useState("Pantau, Cegah, & Kendalikan DBD");
  const [announcementText, setAnnouncementText] = useState("");
  const [isCitizenPortalActive, setIsCitizenPortalActive] = useState(true);
  const [bannerUrl, setBannerUrl] = useState("/images/hero-dbd-new.jpg");
  const [stats, setStats] = useState<PublicStats | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [settingsRes, statsRes] = await Promise.all([
          api.get("/public/settings").catch(() => null),
          api.get("/public/stats").catch(() => null)
        ]);
        
        if (settingsRes?.data) {
          if (settingsRes.data.HeroTitle) setHeroTitle(settingsRes.data.HeroTitle);
          if (settingsRes.data.AnnouncementText) setAnnouncementText(settingsRes.data.AnnouncementText);
          if (settingsRes.data.IsCitizenPortalActive !== undefined) setIsCitizenPortalActive(settingsRes.data.IsCitizenPortalActive);
          if (settingsRes.data.BannerImagePath) {
            const baseUrl = process.env.NEXT_PUBLIC_API_URL ? process.env.NEXT_PUBLIC_API_URL.replace('/api/v1', '') : 'http://localhost:8080';
            setBannerUrl(baseUrl + settingsRes.data.BannerImagePath);
          }
        }
        
        if (statsRes?.data?.data) {
          setStats(statsRes.data.data);
        }
      } catch (error) {
        console.error("Failed to load data:", error);
      }
    };
    fetchData();
  }, []);

  return (
    <main className="flex min-h-screen flex-col bg-bg font-sans">
      <Header />

      {/* Hero Section */}
      <section className="relative w-full min-h-screen flex items-center overflow-hidden">
        {/* Full background image */}
        <div className="absolute inset-0 w-full h-full z-0">
          <Image 
            src={bannerUrl}
            alt="Sistem Terpadu Penanganan DBD Kota Depok" 
            fill
            className="object-cover object-center"
            priority
          />
        </div>
        
        {/* Gradient Overlay for perfect text contrast (always dark on the left) */}
        <div className="absolute inset-0 bg-gradient-to-r from-slate-950/90 via-slate-900/60 to-transparent z-0"></div>
        {/* Mobile overlay for better readability on small screens */}
        <div className="absolute inset-0 bg-primary-deep/40 sm:hidden z-0"></div>

        {/* Animated Blobs for Premium Feel */}
        <div className="absolute top-0 right-1/4 w-[600px] h-[600px] bg-primary/40 dark:bg-primary/20 rounded-full blur-[120px] mix-blend-screen animate-blob z-0 pointer-events-none"></div>
        <div className="absolute bottom-0 left-1/4 w-[500px] h-[500px] bg-teal-400/30 dark:bg-teal-400/10 rounded-full blur-[100px] mix-blend-screen animate-blob animation-delay-400 z-0 pointer-events-none"></div>

        {/* Content Container */}
        <div className="relative z-10 px-4 sm:px-6 lg:px-12 w-full max-w-7xl mx-auto opacity-0 animate-fade-in-up">
          <div className="max-w-2xl space-y-8">
            <div className="inline-flex items-center rounded-full border border-primary/30 bg-primary/10 backdrop-blur-md px-4 py-2 text-sm font-bold text-teal-300 shadow-sm">
              <Shield className="w-4 h-4 mr-2" />
              Sistem Terpadu DBD Guard Kota Depok
            </div>
            
            <h1 className="text-4xl sm:text-5xl lg:text-7xl font-extrabold tracking-tight text-white leading-[1.15]" dangerouslySetInnerHTML={{ __html: heroTitle.replace('\n', '<br className="hidden lg:block" />') }} />
            
            <p className="text-lg sm:text-xl text-slate-200 leading-relaxed font-medium max-w-lg drop-shadow">
              Satu portal resmi untuk pelaporan jentik, penjadwalan fogging, dan edukasi warga. Mari bersama wujudkan lingkungan sehat bebas Demam Berdarah.
            </p>
            
            {announcementText && (
              <div className="mt-4 p-4 rounded-xl bg-amber-500/20 border border-amber-500/50 backdrop-blur-md flex items-start gap-3">
                <AlertCircle className="w-6 h-6 text-amber-400 flex-shrink-0 mt-0.5" />
                <p className="text-amber-50 text-sm md:text-base font-medium">{announcementText}</p>
              </div>
            )}

            <div className="flex flex-col sm:flex-row gap-4 pt-4">
              {isCitizenPortalActive ? (
                <Link 
                  href="/report" 
                  className="inline-flex items-center justify-center gap-2 rounded-xl bg-primary px-8 py-4 text-base font-bold text-primary-foreground shadow-lg shadow-primary/30 hover:bg-teal-400 hover:-translate-y-0.5 transition-all active:scale-95"
                >
                  Lapor Sekarang
                  <ArrowRight className="h-5 w-5" />
                </Link>
              ) : (
                <div className="inline-flex items-center justify-center gap-2 rounded-xl bg-slate-600 px-8 py-4 text-base font-bold text-slate-300 shadow-lg cursor-not-allowed">
                  Portal Sedang Ditutup
                </div>
              )}
              <Link 
                href="/track" 
                className="inline-flex items-center justify-center gap-2 rounded-xl bg-surface/10 backdrop-blur-md px-8 py-4 text-base font-bold text-white border border-white/20 shadow-sm hover:bg-surface/20 hover:-translate-y-0.5 transition-all active:scale-95"
              >
                Pantau Laporan
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Data Publik & Rating Section */}
      <section id="publik" className="px-4 sm:px-6 lg:px-12 py-16 max-w-7xl mx-auto w-full relative z-10">
        <div className="text-center mb-10">
          <h2 className="text-3xl font-extrabold text-text mb-2">Transparansi Data & Penilaian</h2>
          <p className="text-text-muted">Statistik penanganan laporan DBD di Kota Depok secara real-time.</p>
        </div>
        
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 opacity-0 animate-fade-in-up animation-delay-200">
          {/* Stat 1: Total Laporan */}
          <div className="group bg-surface border border-border shadow-md rounded-2xl p-6 hover:-translate-y-1 transition-all duration-300">
            <div className="flex items-center justify-between mb-4">
              <div className="h-10 w-10 bg-blue-500/10 rounded-xl flex items-center justify-center">
                <Activity className="h-5 w-5 text-blue-500" />
              </div>
            </div>
            <div>
              <h4 className="text-3xl font-extrabold text-text mb-1">{stats ? stats.total_reports.toLocaleString("id-ID") : "..."}</h4>
              <p className="text-sm text-text-muted font-semibold">Total Laporan Masuk</p>
            </div>
          </div>
          {/* Stat 2: Selesai */}
          <div className="group bg-surface border border-border shadow-md rounded-2xl p-6 hover:-translate-y-1 transition-all duration-300">
            <div className="flex items-center justify-between mb-4">
              <div className="h-10 w-10 bg-emerald-500/10 rounded-xl flex items-center justify-center">
                <CheckCircle2 className="h-5 w-5 text-emerald-500" />
              </div>
            </div>
            <div>
              <h4 className="text-3xl font-extrabold text-emerald-500 mb-1">{stats ? stats.total_completed.toLocaleString("id-ID") : "..."}</h4>
              <p className="text-sm text-text-muted font-semibold">Selesai Ditangani</p>
            </div>
          </div>
          {/* Stat 3: Proses */}
          <div className="group bg-surface border border-border shadow-md rounded-2xl p-6 hover:-translate-y-1 transition-all duration-300">
            <div className="flex items-center justify-between mb-4">
              <div className="h-10 w-10 bg-amber-500/10 rounded-xl flex items-center justify-center">
                <Clock className="h-5 w-5 text-amber-500" />
              </div>
            </div>
            <div>
              <h4 className="text-3xl font-extrabold text-amber-500 mb-1">{stats ? stats.total_in_progress.toLocaleString("id-ID") : "..."}</h4>
              <p className="text-sm text-text-muted font-semibold">Sedang Diproses</p>
            </div>
          </div>
          {/* Stat 4: Rating */}
          <div className="group bg-surface border border-border shadow-md rounded-2xl p-6 hover:-translate-y-1 transition-all duration-300">
            <div className="flex items-center justify-between mb-4">
              <div className="h-10 w-10 bg-amber-500/10 rounded-xl flex items-center justify-center">
                <Star className="h-5 w-5 text-amber-500" />
              </div>
              {stats && stats.feedback_count > 0 && (
                <span className="text-xs text-text-muted">dari {stats.feedback_count} ulasan</span>
              )}
            </div>
            <div className="flex items-center gap-3">
              <h4 className="text-3xl font-extrabold text-text mb-1">{stats && stats.avg_feedback > 0 ? stats.avg_feedback.toFixed(1) : "—"}</h4>
              <div className="flex">
                {[1, 2, 3, 4, 5].map((s) => (
                  <Star key={s} className={`h-4 w-4 ${stats && s <= Math.round(stats.avg_feedback) ? "text-amber-400 fill-amber-400" : "text-border"}`} />
                ))}
              </div>
            </div>
            <p className="text-sm text-text-muted font-semibold">Tingkat Kepuasan Warga</p>
          </div>
        </div>
      </section>

      {/* Fitur Utama / Layanan */}
      <section id="layanan" className="py-24 w-full bg-bg dark:bg-bg relative overflow-hidden">
        {/* Subtle background blob for Layanan */}
        <div className="absolute top-1/2 left-1/2 w-[800px] h-[800px] bg-primary/5 rounded-full blur-[100px] -translate-x-1/2 -translate-y-1/2 pointer-events-none"></div>
        <div className="px-4 sm:px-6 lg:px-12 max-w-7xl mx-auto relative z-10 opacity-0 animate-fade-in-up animation-delay-400">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-extrabold text-text mb-4 tracking-tight">Layanan Terintegrasi DBD Guard</h2>
            <p className="text-text-muted max-w-2xl mx-auto text-lg leading-relaxed">Solusi end-to-end yang menjembatani warga dengan Kader dan Puskesmas untuk respon penanganan DBD yang lebih cepat dan tepat.</p>
          </div>
          
          <div className="grid md:grid-cols-3 gap-8">
            <div className="glass-panel rounded-3xl p-10 hover:-translate-y-2 hover:glow-border transition-all duration-500 group">
              <div className="w-16 h-16 bg-primary/10 text-primary rounded-2xl flex items-center justify-center mb-8 group-hover:scale-110 group-hover:bg-primary group-hover:text-white transition-all duration-300 shadow-[0_0_15px_rgba(45,212,191,0.2)]">
                <MapPin className="w-8 h-8" />
              </div>
              <h3 className="text-2xl font-bold text-text mb-3 tracking-tight">Peta Sebaran Kasus (GIS)</h3>
              <p className="text-text-muted leading-relaxed">Pantau zona kerawanan DBD dan lokasi temuan jentik di sekitar Anda secara real-time melalui peta interaktif.</p>
            </div>
            <div className="glass-panel rounded-3xl p-10 hover:-translate-y-2 hover:glow-border transition-all duration-500 group">
              <div className="w-16 h-16 bg-primary/10 text-primary rounded-2xl flex items-center justify-center mb-8 group-hover:scale-110 group-hover:bg-primary group-hover:text-white transition-all duration-300 shadow-[0_0_15px_rgba(45,212,191,0.2)]">
                <Stethoscope className="w-8 h-8" />
              </div>
              <h3 className="text-2xl font-bold text-text mb-3 tracking-tight">Tindak Lanjut Cepat</h3>
              <p className="text-text-muted leading-relaxed">Sistem notifikasi otomatis kepada Kader dan Puskesmas terdekat begitu laporan jentik atau kasus masuk.</p>
            </div>
            <div className="glass-panel rounded-3xl p-10 hover:-translate-y-2 hover:shadow-[0_0_15px_rgba(251,191,36,0.3)] hover:border-amber-500/50 transition-all duration-500 group">
              <div className="w-16 h-16 bg-amber-500/10 text-amber-500 rounded-2xl flex items-center justify-center mb-8 group-hover:scale-110 group-hover:bg-amber-500 group-hover:text-white transition-all duration-300 shadow-[0_0_15px_rgba(251,191,36,0.2)]">
                <HomeIcon className="w-8 h-8" />
              </div>
              <h3 className="text-2xl font-bold text-text mb-3 tracking-tight">Penjadwalan Fogging</h3>
              <p className="text-text-muted leading-relaxed">Jadwal Penyelidikan Epidemiologi (PE) dan penyemprotan fogging yang terkoordinasi jelas dan transparan bagi warga.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Alur Kerja / How it Works */}
      <section id="alur" className="w-full bg-primary-deep/5 dark:bg-primary-deep/10 border-b border-border py-24 mb-8">
        <div className="px-4 sm:px-6 lg:px-12 max-w-7xl mx-auto text-center">
          <h2 className="text-3xl font-extrabold text-text mb-16">Bagaimana DBD Guard Bekerja?</h2>
          
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-8 relative">
            {/* Connecting line for desktop */}
            <div className="hidden md:block absolute top-1/2 left-0 w-full h-0.5 bg-border -z-10 -translate-y-full"></div>
            
            <div className="relative z-10 flex flex-col items-center">
              <div className="w-16 h-16 bg-surface border-4 border-primary text-primary font-bold text-xl rounded-full flex items-center justify-center shadow-lg mb-6">1</div>
              <h4 className="text-lg font-bold text-text mb-2">Warga Melapor</h4>
              <p className="text-sm text-text-muted max-w-[200px]">Laporkan temuan genangan air atau penderita DBD via aplikasi.</p>
            </div>
            <div className="relative z-10 flex flex-col items-center">
              <div className="w-16 h-16 bg-surface border-4 border-primary text-primary font-bold text-xl rounded-full flex items-center justify-center shadow-lg mb-6">2</div>
              <h4 className="text-lg font-bold text-text mb-2">Verifikasi Kader</h4>
              <p className="text-sm text-text-muted max-w-[200px]">Kader memvalidasi laporan dalam waktu maksimal 3x24 jam.</p>
            </div>
            <div className="relative z-10 flex flex-col items-center">
              <div className="w-16 h-16 bg-surface border-4 border-blue-500 text-status-verified font-bold text-xl rounded-full flex items-center justify-center shadow-lg mb-6">3</div>
              <h4 className="text-lg font-bold text-text mb-2">Tindakan Puskesmas</h4>
              <p className="text-sm text-text-muted max-w-[200px]">Petugas datang untuk edukasi 3M atau investigasi lanjutan.</p>
            </div>
            <div className="relative z-10 flex flex-col items-center">
              <div className="w-16 h-16 bg-surface border-4 border-orange-500 text-status-scheduled font-bold text-xl rounded-full flex items-center justify-center shadow-lg mb-6">4</div>
              <h4 className="text-lg font-bold text-text mb-2">Tuntas & Fogging</h4>
              <p className="text-sm text-text-muted max-w-[200px]">Bila darurat, fogging dilakukan dan laporan ditutup.</p>
            </div>
          </div>
        </div>
      </section>


      {/* Educational 3M Tips Card */}
      <section id="edukasi" className="px-4 sm:px-6 lg:px-12 py-12 mb-16 max-w-7xl mx-auto w-full">
        <div className="bg-primary-tint rounded-xl p-8 border border-primary/10 flex flex-col md:flex-row gap-8 items-center">
          <div className="flex-1">
            <h2 className="text-2xl font-bold text-primary mb-4">Pencegahan DBD dengan 3M Plus</h2>
            <p className="text-text-muted mb-6 leading-relaxed">
              Mari bersama-sama memutus mata rantai perkembangbiakan nyamuk Aedes aegypti di lingkungan rumah kita. Lakukan 3M Plus secara rutin minimal seminggu sekali.
            </p>
            <Link 
              href="/edukasi" 
              className="inline-flex items-center justify-center gap-2 rounded-lg bg-surface px-6 py-3 text-sm font-bold text-primary border border-border hover:bg-bg transition-all"
            >
              Pelajari Lebih Lanjut
            </Link>
          </div>
          <div className="flex-1 grid grid-cols-1 sm:grid-cols-3 gap-4 w-full">
            <div className="bg-surface p-4 rounded-xl border border-border flex flex-col items-center text-center">
              <div className="w-12 h-12 bg-primary-tint text-primary rounded-full flex items-center justify-center mb-3">
                <Trash2 className="w-6 h-6" />
              </div>
              <h3 className="font-bold text-text text-sm">Menguras</h3>
              <p className="text-xs text-text-muted mt-2">Tempat penampungan air</p>
            </div>
            <div className="bg-surface p-4 rounded-xl border border-border flex flex-col items-center text-center">
              <div className="w-12 h-12 bg-primary-tint text-primary rounded-full flex items-center justify-center mb-3">
                <Droplets className="w-6 h-6" />
              </div>
              <h3 className="font-bold text-text text-sm">Menutup</h3>
              <p className="text-xs text-text-muted mt-2">Rapat tempat air</p>
            </div>
            <div className="bg-surface p-4 rounded-xl border border-border flex flex-col items-center text-center">
              <div className="w-12 h-12 bg-primary-tint text-primary rounded-full flex items-center justify-center mb-3">
                <HomeIcon className="w-6 h-6" />
              </div>
              <h3 className="font-bold text-text text-sm">Mengubur</h3>
              <p className="text-xs text-text-muted mt-2">Barang bekas</p>
            </div>
          </div>
        </div>
      </section>

      {/* FAQ Section */}
      <section id="faq" className="px-4 sm:px-6 lg:px-12 py-16 mb-8 max-w-7xl mx-auto w-full border-t border-border">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-extrabold text-text mb-4">Pertanyaan yang Sering Diajukan (FAQ)</h2>
          <p className="text-text-muted text-lg max-w-2xl mx-auto">Informasi yang mungkin Anda cari seputar sistem pelaporan DBD Guard.</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-surface border border-border p-6 rounded-2xl shadow-sm">
            <h4 className="text-lg font-bold text-text mb-2">Berapa lama laporan saya akan diproses?</h4>
            <p className="text-sm text-text-muted leading-relaxed">Kader/RT memiliki waktu maksimal 3x24 jam untuk memverifikasi. Setelah terverifikasi, Puskesmas akan menjadwalkan Penyelidikan Epidemiologi (PE).</p>
          </div>
          <div className="bg-surface border border-border p-6 rounded-2xl shadow-sm">
            <h4 className="text-lg font-bold text-text mb-2">Apakah saya bisa melaporkan secara anonim?</h4>
            <p className="text-sm text-text-muted leading-relaxed">Ya, Anda bisa menyembunyikan identitas Anda saat melapor. Namun, kami menyarankan Anda meninggalkan kontak agar kami dapat menginformasikan status penanganan.</p>
          </div>
          <div className="bg-surface border border-border p-6 rounded-2xl shadow-sm">
            <h4 className="text-lg font-bold text-text mb-2">Kapan Fogging akan dilakukan?</h4>
            <p className="text-sm text-text-muted leading-relaxed">Fogging hanya dilakukan jika hasil survei Puskesmas menemukan jentik dan ada kasus demam tambahan di radius 100m dari lokasi Anda.</p>
          </div>
          <div className="bg-surface border border-border p-6 rounded-2xl shadow-sm">
            <h4 className="text-lg font-bold text-text mb-2">Apakah aplikasi ini gratis?</h4>
            <p className="text-sm text-text-muted leading-relaxed">Seluruh layanan DBD Guard 100% gratis dan merupakan inisiatif resmi Pemerintah Kota Depok untuk kesehatan masyarakat.</p>
          </div>
        </div>
      </section>

      <Footer />
    </main>
  );
}
