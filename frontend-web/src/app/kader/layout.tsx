"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { useState } from "react";
import { LayoutDashboard, CheckSquare, LogOut, Menu, X } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";
import NotificationPermission from "@/components/NotificationPermission";

export default function KaderLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const router = useRouter();
  const pathname = usePathname();
  const [isMobileDrawerOpen, setIsMobileDrawerOpen] = useState(false);
  const [prevPathname, setPrevPathname] = useState(pathname);

  if (prevPathname !== pathname) {
    setPrevPathname(pathname);
    setIsMobileDrawerOpen(false);
  }

  const handleLogout = () => {
    Cookies.remove("auth_token");
    router.push("/");
  };

  const getLinkClass = (path: string) => {
    const isActive = path === "/kader" 
      ? pathname === "/kader"
      : pathname.startsWith(path);
      
    return `flex items-center space-x-3 px-4 py-3 rounded-xl transition-all duration-200 group ${
      isActive
        ? "bg-primary text-white shadow-md shadow-primary/20 font-bold"
        : "text-text-muted hover:bg-bg hover:text-primary font-medium"
    }`;
  };

  return (
    <div className="min-h-screen bg-bg flex">
      <NotificationPermission />

      {/* Sidebar - Desktop */}
      <aside className="hidden lg:flex flex-col w-72 bg-surface border-r border-border fixed h-full z-10 transition-all duration-300">
        <div className="p-6 flex items-center space-x-3 border-b border-border/60">
          <div className="w-10 h-10 rounded-xl bg-primary flex items-center justify-center text-white font-bold text-xl shadow-md shadow-primary/20">
            K
          </div>
          <div>
            <span className="text-lg font-bold text-text block">DBD Guard</span>
            <span className="text-xs text-text-muted">Kader Wilayah RT/RW</span>
          </div>
        </div>

        <nav className="flex-1 px-4 py-6 space-y-2 overflow-y-auto">
          <Link href="/kader" className={getLinkClass("/kader")}>
            <LayoutDashboard className="w-5 h-5 shrink-0" />
            <span>Dashboard Kader</span>
          </Link>
          <Link href="/kader/verification" className={getLinkClass("/kader/verification")}>
            <CheckSquare className="w-5 h-5 shrink-0" />
            <span>Penyelidikan Epidemiologi</span>
          </Link>
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

      {/* Mobile Drawer (Slide-Over) */}
      {isMobileDrawerOpen && (
        <div className="fixed inset-0 z-50 lg:hidden flex">
          <div 
            className="fixed inset-0 bg-slate-900/60 backdrop-blur-xs transition-opacity"
            onClick={() => setIsMobileDrawerOpen(false)}
          />
          <div className="relative z-50 w-72 max-w-[80vw] bg-surface h-full flex flex-col border-r border-border shadow-2xl animate-in slide-in-from-left duration-200">
            <div className="p-5 flex items-center justify-between border-b border-border">
              <div className="flex items-center space-x-3">
                <div className="w-9 h-9 rounded-xl bg-primary flex items-center justify-center text-white font-bold shadow-xs">
                  K
                </div>
                <div>
                  <span className="font-bold text-text text-base block">DBD Guard</span>
                  <span className="text-xs text-text-muted">Kader RT/RW</span>
                </div>
              </div>
              <button 
                type="button" 
                onClick={() => setIsMobileDrawerOpen(false)}
                className="p-2 rounded-lg text-text-muted hover:bg-bg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <nav className="flex-1 px-3 py-4 space-y-1.5 overflow-y-auto">
              <Link 
                href="/kader" 
                onClick={() => setIsMobileDrawerOpen(false)}
                className={getLinkClass("/kader")}
              >
                <LayoutDashboard className="w-5 h-5 shrink-0" />
                <span>Dashboard Kader</span>
              </Link>
              <Link 
                href="/kader/verification" 
                onClick={() => setIsMobileDrawerOpen(false)}
                className={getLinkClass("/kader/verification")}
              >
                <CheckSquare className="w-5 h-5 shrink-0" />
                <span>Penyelidikan Epidemiologi</span>
              </Link>
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
      <div className="flex-1 flex flex-col lg:ml-72 min-w-0 overflow-hidden">
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
              Portal Verifikasi Lapangan
            </h1>
          </div>
          
          <div className="flex items-center gap-2 sm:gap-4">
            <ThemeToggle />
            <div className="hidden sm:flex flex-col items-end">
              <span className="text-sm font-bold text-text leading-tight">Perangkat RT/RW</span>
              <span className="text-xs text-text-muted">Kader Lapangan</span>
            </div>
            <div className="h-9 w-9 sm:h-10 sm:w-10 bg-primary-tint rounded-full flex items-center justify-center text-primary font-bold text-sm">
              RT
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
