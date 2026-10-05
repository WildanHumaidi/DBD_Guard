"use client";

import Link from 'next/link';
import { Shield } from 'lucide-react';
import { useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';

export default function Header() {
  const [isScrolled, setIsScrolled] = useState(false);
  const pathname = usePathname();

  useEffect(() => {
    const handleScroll = () => {
      // Hero section is min-h-screen (100vh), so we only want the header to turn solid 
      // when we scroll past the hero banner into the white sections.
      setIsScrolled(window.scrollY > window.innerHeight - 80);
    };
    
    window.addEventListener('scroll', handleScroll);
    handleScroll(); // Check initial scroll position
    
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const isHomeOrPublik = pathname === '/' || pathname === '/publik';
  const isTransparent = isHomeOrPublik && !isScrolled;

  return (
    <nav 
      className={`fixed top-0 left-0 w-full z-50 transition-all duration-300 px-6 lg:px-12 py-5 flex items-center justify-between
        ${isTransparent 
          ? 'bg-transparent border-transparent' 
          : 'bg-surface/90 backdrop-blur-md border-b border-border shadow-sm'
        }
      `}
    >
      <div className="flex items-center gap-3">
        <div className="bg-primary p-2 rounded-xl">
          <Shield className="h-6 w-6 text-primary-foreground" />
        </div>
        <span className={`font-bold text-2xl tracking-tight transition-colors ${isTransparent ? 'text-white' : 'text-text'}`}>
          DBD Guard
        </span>
      </div>
      
      <div className="flex items-center gap-6">
        <div className="hidden md:flex items-center gap-6 mr-2">
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
          href="/#publik" 
          className={`text-sm font-semibold transition-colors hover:text-primary hidden sm:block ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}
        >
          Data Publik
        </Link>
        <Link 
          href="/login" 
          className={`text-sm font-semibold transition-colors hover:text-primary hidden sm:block ${isTransparent ? 'text-slate-200' : 'text-text-muted'}`}
        >
          Login
        </Link>
        
        <Link 
          href="/report" 
          className="text-sm font-bold bg-primary text-primary-foreground px-6 py-2.5 rounded-lg hover:opacity-90 transition-all active:scale-95 shadow-md shadow-primary/20"
        >
          Lapor Sekarang
        </Link>
      </div>
    </nav>
  );
}
