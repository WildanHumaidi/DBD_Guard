"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { useState } from "react";
import { LayoutDashboard, Users, Settings, LogOut, Code, Activity, ShieldCheck, FileText, Menu, X } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";

export default function DeveloperLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
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
    const isActive = path === "/developer" 
      ? pathname === "/developer"
      : pathname.startsWith(path);
      
    if (isActive) {
      return "flex items-center px-4 py-3 bg-primary text-white rounded-xl transition-colors font-bold shadow-xs";
    }
    return "flex items-center px-4 py-3 text-text-muted hover:bg-bg hover:text-primary rounded-xl transition-colors font-medium";
  };

  const navLinks = (
    <>
      <Link href="/developer" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/developer")}>
        <LayoutDashboard className="h-5 w-5 mr-3 shrink-0" />
        <span>Beranda Dev</span>
      </Link>
      <Link href="/developer/users" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/developer/users")}>
        <Users className="h-5 w-5 mr-3 shrink-0" />
        <span>Kelola Akun</span>
      </Link>
      <Link href="/developer/settings" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/developer/settings")}>
        <Settings className="h-5 w-5 mr-3 shrink-0" />
        <span>Pengaturan Website</span>
      </Link>
      <Link href="/developer/status" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/developer/status")}>
        <Activity className="h-5 w-5 mr-3 shrink-0" />
        <span>System Status</span>
      </Link>
      <Link href="/developer/audit-logs" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/developer/audit-logs")}>
        <FileText className="h-5 w-5 mr-3 shrink-0" />
        <span>Log Aktivitas</span>
      </Link>
    </>
  );

  return (
    <div className="min-h-screen bg-bg flex">
      {/* Desktop Sidebar */}
      <aside className="hidden lg:flex flex-col w-64 bg-surface border-r border-border fixed h-full z-10">
        <div className="h-16 flex items-center px-6 border-b border-border">
          <Code className="h-6 w-6 text-primary mr-2 shrink-0" />
          <span className="font-bold text-lg text-text">DBD Guard Dev</span>
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
                <Code className="h-6 w-6 text-primary mr-2 shrink-0" />
                <span className="font-bold text-base text-text">Developer Center</span>
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
            <div className="flex items-center gap-2">
              <ShieldCheck className="h-5 w-5 text-primary shrink-0" />
              <h1 className="text-base sm:text-lg lg:text-xl font-bold text-text truncate">
                Developer Center
              </h1>
            </div>
          </div>

          <div className="flex items-center gap-2 sm:gap-4">
            <ThemeToggle />
            <div className="hidden sm:flex flex-col items-end ml-2">
              <span className="text-sm font-bold text-text">Admin Developer</span>
              <span className="text-xs text-text-muted">Sistem Utama</span>
            </div>
            <div className="h-9 w-9 sm:h-10 sm:w-10 bg-slate-800 rounded-full flex items-center justify-center text-white font-bold text-sm">
              DEV
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
