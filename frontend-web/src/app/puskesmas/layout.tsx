"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { LayoutDashboard, Stethoscope, LogOut, HeartPulse, Package, ClipboardCheck, AlertTriangle, Building } from "lucide-react";
import { useState, useEffect } from "react";
import { ThemeToggle } from "@/components/ThemeToggle";
import { ShieldExclamationIcon } from "@heroicons/react/24/outline";
import NotificationPermission from "@/components/NotificationPermission";

export default function PuskesmasLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();

  const handleLogout = () => {
    Cookies.remove("auth_token");
    router.push("/");
  };

  const [role, setRole] = useState<string | null>(null);

  useEffect(() => {
    const token = Cookies.get("auth_token");
    if (token) {
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
        setRole(decoded.role);
      } catch (e) {
        console.error("Failed to decode token", e);
      }
    }
  }, []);

  const getLinkClass = (path: string) => {
    const isActive = path === "/puskesmas" 
      ? pathname === "/puskesmas"
      : pathname.startsWith(path);
      
    if (isActive) {
      return "flex items-center px-4 py-3 bg-status-verified/20 text-blue-700 rounded-xl transition-colors font-bold";
    }
    return "flex items-center px-4 py-3 text-text-muted hover:bg-bg hover:text-status-verified rounded-xl transition-colors font-medium";
  };

  return (
    <div className="flex h-screen bg-bg overflow-hidden">
      <NotificationPermission />
      {/* Sidebar */}
      <aside className="w-64 bg-surface border-r border-border flex flex-col">
        <div className="h-16 flex items-center px-6 border-b border-border">
          <HeartPulse className="h-6 w-6 text-status-verified mr-2" />
          <span className="font-bold text-lg text-text">Puskesmas Tim</span>
        </div>
        
        <nav className="flex-1 py-6 px-4 space-y-2">
          <Link href="/puskesmas" className={getLinkClass("/puskesmas")}>
            <LayoutDashboard className="h-5 w-5 mr-3" />
            <span>Beranda</span>
          </Link>
          <Link href="/puskesmas/pe-validation" className={getLinkClass("/puskesmas/pe-validation")}>
            <Stethoscope className="h-5 w-5 mr-3" />
            <span>Validasi PE Kader</span>
          </Link>
          
          {(role === "ROLE_KEPALA_PUSKESMAS" || role === "Kepala Puskesmas") && (
            <>
              <Link href="/puskesmas/fogging-requests" className={getLinkClass("/puskesmas/fogging-requests")}>
                <HeartPulse className="h-5 w-5 mr-3" />
                <span>Jadwal Fogging</span>
              </Link>
              <Link href="/puskesmas/inventory-approval" className={getLinkClass("/puskesmas/inventory-approval")}>
                <Package className="h-5 w-5 mr-3" />
                <span>Persetujuan Stok</span>
              </Link>
            </>
          )}

          <Link href="/puskesmas/ews" className={getLinkClass("/puskesmas/ews")}>
            <ShieldExclamationIcon className="h-5 w-5 mr-3" />
            <span>Peringatan EWS</span>
          </Link>
        </nav>

        <div className="p-4 border-t border-border">
          <button 
            onClick={handleLogout}
            className="flex w-full items-center px-4 py-3 text-red-600 hover:bg-red-50 rounded-xl transition-colors"
          >
            <LogOut className="h-5 w-5 mr-3" />
            <span className="font-medium">Keluar</span>
          </button>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="h-16 bg-surface border-b border-border flex items-center justify-between px-8 z-10">
          <h1 className="text-xl font-bold text-text">Portal Tim Puskesmas</h1>
          <div className="flex items-center gap-4">
            <ThemeToggle />
            <div className="flex flex-col items-end ml-2">
              <span className="text-sm font-bold text-text">Tim Puskesmas</span>
              <span className="text-xs text-text-muted">
                {(() => {
                  try {
                    const token = Cookies.get("auth_token");
                    if (token) {
                      const payload = JSON.parse(atob(token.split('.')[1]));
                      if (payload.role === "ROLE_KEPALA_PUSKESMAS") return "Kepala Puskesmas";
                    }
                  } catch (e) {}
                  return "Penyelidik Epidemiologi";
                })()}
              </span>
            </div>
            <div className="h-10 w-10 bg-status-verified/20 rounded-full flex items-center justify-center text-blue-700 font-bold">
              {(() => {
                  try {
                    const token = Cookies.get("auth_token");
                    if (token) {
                      const payload = JSON.parse(atob(token.split('.')[1]));
                      if (payload.role === "ROLE_KEPALA_PUSKESMAS") return "KP";
                    }
                  } catch (e) {}
                  return "PE";
                })()}
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-8">
          <div className="w-full max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}
