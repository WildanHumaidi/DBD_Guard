"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { LayoutDashboard, Stethoscope, LogOut, HeartPulse, Package, Menu, X } from "lucide-react";
import { useState, useSyncExternalStore } from "react";
import { ThemeToggle } from "@/components/ThemeToggle";
import { ShieldExclamationIcon } from "@heroicons/react/24/outline";
import NotificationPermission from "@/components/NotificationPermission";

function getRoleFromCookie(): string | null {
  if (typeof window === "undefined") return null;
  const token = Cookies.get("auth_token");
  if (!token) return null;
  try {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join("")
    );
    const decoded = JSON.parse(jsonPayload);
    return decoded.role || null;
  } catch {
    return null;
  }
}

const emptySubscribe = () => () => {};

export default function PuskesmasLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const role = useSyncExternalStore(emptySubscribe, getRoleFromCookie, () => null);
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
    const isActive = path === "/puskesmas" 
      ? pathname === "/puskesmas"
      : pathname.startsWith(path);
      
    if (isActive) {
      return "flex items-center px-4 py-3 bg-primary text-white rounded-xl transition-colors font-bold shadow-xs";
    }
    return "flex items-center px-4 py-3 text-text-muted hover:bg-bg hover:text-primary rounded-xl transition-colors font-medium";
  };

  const isKepalaPuskesmas = role === "ROLE_KEPALA_PUSKESMAS" || role === "Kepala Puskesmas";

  const navLinks = (
    <>
      <Link href="/puskesmas" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/puskesmas")}>
        <LayoutDashboard className="h-5 w-5 mr-3 shrink-0" />
        <span>Beranda</span>
      </Link>
      <Link href="/puskesmas/pe-validation" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/puskesmas/pe-validation")}>
        <Stethoscope className="h-5 w-5 mr-3 shrink-0" />
        <span>Validasi PE Kader</span>
      </Link>
      
      {isKepalaPuskesmas && (
        <>
          <Link href="/puskesmas/fogging-requests" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/puskesmas/fogging-requests")}>
            <HeartPulse className="h-5 w-5 mr-3 shrink-0" />
            <span>Jadwal Fogging</span>
          </Link>
          <Link href="/puskesmas/inventory-approval" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/puskesmas/inventory-approval")}>
            <Package className="h-5 w-5 mr-3 shrink-0" />
            <span>Persetujuan Stok</span>
          </Link>
        </>
      )}

      <Link href="/puskesmas/ews" onClick={() => setIsMobileDrawerOpen(false)} className={getLinkClass("/puskesmas/ews")}>
        <ShieldExclamationIcon className="h-5 w-5 mr-3 shrink-0" />
        <span>Peringatan EWS</span>
      </Link>
    </>
  );

  return (
    <div className="min-h-screen bg-bg flex">
      <NotificationPermission />

      {/* Desktop Sidebar */}
      <aside className="hidden lg:flex flex-col w-64 bg-surface border-r border-border fixed h-full z-10">
        <div className="h-16 flex items-center px-6 border-b border-border">
          <HeartPulse className="h-6 w-6 text-primary mr-2 shrink-0" />
          <span className="font-bold text-lg text-text">Puskesmas Tim</span>
        </div>
        
        <nav className="flex-1 py-6 px-4 space-y-2 overflow-y-auto">
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

      {/* Mobile Drawer (Slide-Over) */}
      {isMobileDrawerOpen && (
        <div className="fixed inset-0 z-50 lg:hidden flex">
          <div 
            className="fixed inset-0 bg-slate-900/60 backdrop-blur-xs transition-opacity"
            onClick={() => setIsMobileDrawerOpen(false)}
          />
          <div className="relative z-50 w-72 max-w-[80vw] bg-surface h-full flex flex-col border-r border-border shadow-2xl animate-in slide-in-from-left duration-200">
            <div className="h-16 flex items-center justify-between px-5 border-b border-border">
              <div className="flex items-center">
                <HeartPulse className="h-6 w-6 text-primary mr-2 shrink-0" />
                <span className="font-bold text-base text-text">Tim Puskesmas</span>
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
              Portal Tim Puskesmas
            </h1>
          </div>

          <div className="flex items-center gap-2 sm:gap-4">
            <ThemeToggle />
            <div className="hidden sm:flex flex-col items-end ml-2">
              <span className="text-sm font-bold text-text">Tim Puskesmas</span>
              <span className="text-xs text-text-muted">
                {isKepalaPuskesmas ? "Kepala Puskesmas" : "Penyelidik Epidemiologi"}
              </span>
            </div>
            <div className="h-9 w-9 sm:h-10 sm:w-10 bg-primary-tint rounded-full flex items-center justify-center text-primary font-bold text-sm">
              {isKepalaPuskesmas ? "KP" : "PE"}
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
