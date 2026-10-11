"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { useState, useEffect } from "react";
import { LayoutDashboard, FileText, Settings, LogOut, Shield, Package, Map, Building2, AlertTriangle, Clock, Menu, X } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const [isMobileDrawerOpen, setIsMobileDrawerOpen] = useState(false);

  useEffect(() => {
    setIsMobileDrawerOpen(false);
  }, [pathname]);

  const handleLogout = () => {
    Cookies.remove("auth_token");
    router.push("/");
  };

  if (pathname === "/dinkes/login") {
    return <>{children}</>;
  }

  const getLinkClass = (path: string) => {
    const isActive = path === "/dinkes" 
      ? pathname === "/dinkes"
      : pathname.startsWith(path);
      
    if (isActive) {
      return "flex items-center px-4 py-3 bg-primary text-white rounded-xl transition-colors font-bold shadow-xs";
    }
    return "flex items-center px-4 py-3 text-text-muted hover:bg-bg hover:text-primary rounded-xl transition-colors font-medium";
  };

  const navLinks = (
    <>
      <Link href="/dinkes" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/dinkes")}>
        <LayoutDashboard className="h-5 w-5 mr-3 shrink-0" />
        <span>Dashboard</span>
      </Link>
      <Link href="/dinkes/reports" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/dinkes/reports")}>
        <FileText className="h-5 w-5 mr-3 shrink-0" />
        <span>Daftar Laporan</span>
      </Link>
      <Link href="/dinkes/inventory" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/dinkes/inventory")}>
        <Package className="h-5 w-5 mr-3 shrink-0" />
        <span>Inventaris</span>
      </Link>
      <Link href="/dinkes/map" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/dinkes/map")}>
        <Map className="h-5 w-5 mr-3 shrink-0" />
        <span>Peta GIS</span>
      </Link>
      <Link href="/dinkes/klb" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/dinkes/klb")}>
        <AlertTriangle className="h-5 w-5 mr-3 shrink-0" />
        <span>Status KLB</span>
      </Link>
      <Link href="/dinkes/sla" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/dinkes/sla")}>
        <Clock className="h-5 w-5 mr-3 shrink-0" />
        <span>SLA Monitoring</span>
      </Link>
      <Link href="/dinkes/puskesmas" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/dinkes/puskesmas")}>
        <Building2 className="h-5 w-5 mr-3 shrink-0" />
        <span>Akun Puskesmas</span>
      </Link>
      <Link href="/dinkes/settings" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/dinkes/settings")}>
        <Settings className="h-5 w-5 mr-3 shrink-0" />
        <span>Pengaturan</span>
      </Link>
    </>
  );

  return (
    <div className="min-h-screen bg-bg flex">
      {/* Desktop Sidebar */}
      <aside className="hidden lg:flex flex-col w-64 bg-surface border-r border-border fixed h-full z-10">
        <div className="h-16 flex items-center px-6 border-b border-border">
          <Shield className="h-6 w-6 text-primary mr-2 shrink-0" />
          <span className="font-bold text-lg text-text">DBD Guard Admin</span>
        </div>
        
        <nav className="flex-1 py-6 px-4 space-y-1.5 overflow-y-auto">
          {navLinks}
        </nav>

        <div className="p-4 border-t border-border">
          <button 
            onClick={handleLogout}
            className="flex w-full items-center px-4 py-3 text-red-600 hover:bg-red-500/10 rounded-xl transition-colors font-semibold"
          >
            <LogOut className="h-5 w-5 mr-3 shrink-0" />
            <span>Keluar</span>
          </button>
        </div>
      </aside>

      {/* Mobile Drawer */}
      {isMobileDrawerOpen && (
        <div className="fixed inset-0 z-50 lg:hidden flex">
          <div 
            className="fixed inset-0 bg-slate-900/60 backdrop-blur-xs transition-opacity"
            onClick={() => setIsMobileDrawerOpen(false)}
          />
          <div className="relative z-50 w-72 max-w-[80vw] bg-surface h-full flex flex-col border-r border-border shadow-2xl animate-in slide-in-from-left duration-200">
            <div className="h-16 flex items-center justify-between px-5 border-b border-border">
              <div className="flex items-center">
                <Shield className="h-6 w-6 text-primary mr-2 shrink-0" />
                <span className="font-bold text-base text-text">Admin Dinkes</span>
              </div>
              <button 
                type="button" 
                onClick={() => setIsMobileDrawerOpen(false)}
                className="p-2 rounded-lg text-text-muted hover:bg-bg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <nav className="flex-1 py-4 px-3 space-y-1.5 overflow-y-auto">
              {navLinks}
            </nav>

            <div className="p-4 border-t border-border">
              <button 
                onClick={handleLogout}
                className="flex w-full items-center px-4 py-3 text-red-600 hover:bg-red-500/10 rounded-xl transition-colors font-semibold"
              >
                <LogOut className="h-5 w-5 mr-3 shrink-0" />
                <span>Keluar</span>
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col lg:ml-64 min-w-0 overflow-hidden">
        {/* Header */}
        <header className="h-16 bg-surface border-b border-border flex items-center justify-between px-4 sm:px-6 lg:px-8 z-10 shrink-0">
          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={() => setIsMobileDrawerOpen(true)}
              aria-label="Buka Menu"
              className="lg:hidden p-2 rounded-xl border border-border text-text hover:bg-bg"
            >
              <Menu className="w-5 h-5" />
            </button>
            <h1 className="text-base sm:text-lg lg:text-xl font-bold text-text truncate">
              Portal Dinas Kesehatan
            </h1>
          </div>

          <div className="flex items-center gap-2 sm:gap-4">
            <ThemeToggle />
            <div className="hidden sm:flex flex-col items-end ml-2">
              <span className="text-sm font-bold text-text">Administrator</span>
              <span className="text-xs text-text-muted">Dinas Kesehatan Depok</span>
            </div>
            <div className="h-9 w-9 sm:h-10 sm:w-10 bg-primary-tint rounded-full flex items-center justify-center text-primary font-bold text-sm">
              ADM
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8">
          <div className="w-full max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}
