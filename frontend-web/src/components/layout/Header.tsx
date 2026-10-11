"use client";

import Link from 'next/link';
import { Shield, Menu, X, ArrowRight, LogIn, FileText, Search } from 'lucide-react';
import { useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';

export default function Header() {
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const pathname = usePathname();

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 40);
    };
    
    window.addEventListener('scroll', handleScroll);
    handleScroll();
    
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // Close mobile menu whenever route changes
  useEffect(() => {
    setIsMobileMenuOpen(false);
  }, [pathname]);

  const isHomeOrPublik = pathname === '/' || pathname === '/publik';
  const isTransparent = isHomeOrPublik && !isScrolled && !isMobileMenuOpen;

  return (
    <>
      <nav 
        className={`fixed top-0 left-0 w-full z-50 transition-all duration-300 px-4 sm:px-6 lg:px-12 py-4 flex items-center justify-between
          ${isTransparent 
            ? 'bg-transparent border-transparent' 
            : 'bg-surface/95 backdrop-blur-md border-b border-border shadow-sm'
          }
        `}
      >
        <Link href="/" className="flex items-center gap-3 group">
          <div className="bg-primary p-2 rounded-xl group-hover:scale-105 transition-transform">
            <Shield className="h-6 w-6 text-white" />
          </div>
          <span className={`font-bold text-xl sm:text-2xl tracking-tight transition-colors ${isTransparent ? 'text-white' : 'text-text'}`}>
            DBD Guard
          </span>
        </Link>
        
        {/* Desktop Navigation Links */}
        <div className="hidden lg:flex items-center gap-6">
          <div className="flex items-center gap-6 mr-2">
            <Link href="/#layanan" className={`text-sm font-semibold transition-colors hover:text-primary ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}>Layanan</Link>
            <Link href="/#alur" className={`text-sm font-semibold transition-colors hover:text-primary ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}>Alur Kerja</Link>
            <Link href="/#edukasi" className={`text-sm font-semibold transition-colors hover:text-primary ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}>Edukasi</Link>
            <Link href="/#faq" className={`text-sm font-semibold transition-colors hover:text-primary ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}>FAQ</Link>
          </div>
          <Link 
            href="/track" 
            className={`text-sm font-semibold transition-colors hover:text-primary ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}
          >
            Lacak Laporan
          </Link>
          <Link 
            href="/publik" 
            className={`text-sm font-semibold transition-colors hover:text-primary ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}
          >
            Data Publik
          </Link>
          <Link 
            href="/login" 
            className={`text-sm font-semibold transition-colors hover:text-primary ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}
          >
            Login
          </Link>
          
          <Link 
            href="/report" 
            className="text-sm font-bold bg-primary text-white px-6 py-2.5 rounded-lg hover:bg-primary-deep transition-all active:scale-95 shadow-md shadow-primary/20"
          >
            Lapor Sekarang
          </Link>
        </div>

        {/* Mobile Action Controls */}
        <div className="flex lg:hidden items-center gap-2">
          <Link 
            href="/report" 
            className="text-xs sm:text-sm font-bold bg-primary text-white px-3.5 sm:px-5 py-2 rounded-lg hover:bg-primary-deep transition-all shadow-xs"
          >
            Lapor
          </Link>

          <button
            type="button"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-label="Buka Menu Navigasi"
            className={`p-2.5 rounded-xl border transition-colors ${
              isTransparent
                ? 'bg-white/10 text-white border-white/20 hover:bg-white/20'
                : 'bg-surface text-text border-border hover:bg-bg'
            }`}
          >
            {isMobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </button>
        </div>
      </nav>

      {/* Mobile Drawer / Backdrop */}
      {isMobileMenuOpen && (
        <div className="fixed inset-0 z-40 lg:hidden flex flex-col justify-start">
          <div 
            className="fixed inset-0 bg-slate-900/60 backdrop-blur-xs transition-opacity"
            onClick={() => setIsMobileMenuOpen(false)}
          />
          <div className="relative z-50 bg-surface border-b border-border shadow-xl px-6 pt-24 pb-8 flex flex-col space-y-4 animate-in slide-in-from-top-4 duration-200">
            <div className="flex flex-col space-y-1">
              <Link 
                href="/#layanan" 
                onClick={() => setIsMobileMenuOpen(false)}
                className="flex items-center px-4 py-3 rounded-xl text-base font-semibold text-text hover:bg-bg hover:text-primary transition-colors"
              >
                Layanan Sistem
              </Link>
              <Link 
                href="/#alur" 
                onClick={() => setIsMobileMenuOpen(false)}
                className="flex items-center px-4 py-3 rounded-xl text-base font-semibold text-text hover:bg-bg hover:text-primary transition-colors"
              >
                Alur Kerja Penanganan
              </Link>
              <Link 
                href="/#edukasi" 
                onClick={() => setIsMobileMenuOpen(false)}
                className="flex items-center px-4 py-3 rounded-xl text-base font-semibold text-text hover:bg-bg hover:text-primary transition-colors"
              >
                Edukasi 3M Plus
              </Link>
              <Link 
                href="/#faq" 
                onClick={() => setIsMobileMenuOpen(false)}
                className="flex items-center px-4 py-3 rounded-xl text-base font-semibold text-text hover:bg-bg hover:text-primary transition-colors"
              >
                FAQ & Tanya Jawab
              </Link>
              <Link 
                href="/publik" 
                onClick={() => setIsMobileMenuOpen(false)}
                className="flex items-center gap-3 px-4 py-3 rounded-xl text-base font-semibold text-text hover:bg-bg hover:text-primary transition-colors"
              >
                <FileText className="h-5 w-5 text-primary" />
                Data Transparansi Publik
              </Link>
              <Link 
                href="/track" 
                onClick={() => setIsMobileMenuOpen(false)}
                className="flex items-center gap-3 px-4 py-3 rounded-xl text-base font-semibold text-text hover:bg-bg hover:text-primary transition-colors"
              >
                <Search className="h-5 w-5 text-primary" />
                Lacak Status Laporan
              </Link>
            </div>

            <div className="pt-4 border-t border-border flex flex-col gap-3">
              <Link 
                href="/login" 
                onClick={() => setIsMobileMenuOpen(false)}
                className="w-full flex items-center justify-center gap-2 py-3 px-4 rounded-xl border border-border bg-bg text-text font-bold hover:bg-surface transition-colors"
              >
                <LogIn className="h-4 w-4" />
                Login Petugas & Admin
              </Link>
              <Link 
                href="/report" 
                onClick={() => setIsMobileMenuOpen(false)}
                className="w-full flex items-center justify-center gap-2 py-3.5 px-4 rounded-xl bg-primary text-white font-bold hover:bg-primary-deep shadow-md shadow-primary/20 transition-colors"
              >
                Lapor Kasus DBD Sekarang
                <ArrowRight className="h-4 w-4" />
              </Link>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
