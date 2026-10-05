"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { LayoutDashboard, FileText, Settings, LogOut, Shield, Package, Map, Building2, AlertTriangle, Clock } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";

export default function AdminLayout({
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

  if (pathname === "/dinkes/login") {
    return <>{children}</>;
  }

  const getLinkClass = (path: string) => {
    // Exact match for dashboard, prefix match for others
    const isActive = path === "/admin" 
      ? pathname === "/admin"
      : pathname.startsWith(path);
      
    if (isActive) {
      return "flex items-center px-4 py-3 bg-primary-tint text-primary rounded-xl transition-colors font-bold";
    }
    return "flex items-center px-4 py-3 text-text-muted hover:bg-bg hover:text-primary rounded-xl transition-colors font-medium";
  };

  return (
    <div className="flex h-screen bg-bg overflow-hidden">
      {/* Sidebar */}
      <aside className="w-64 bg-surface border-r border-border flex flex-col">
        <div className="h-16 flex items-center px-6 border-b border-border">
          <Shield className="h-6 w-6 text-primary mr-2" />
          <span className="font-bold text-lg text-text">DBD Guard Admin</span>
        </div>
        
        <nav className="flex-1 py-6 px-4 space-y-2">
          <Link href="/dinkes" className={getLinkClass("/dinkes")}>
            <LayoutDashboard className="h-5 w-5 mr-3" />
            <span>Dashboard</span>
          </Link>
          <Link href="/dinkes/reports" className={getLinkClass("/dinkes/reports")}>
            <FileText className="h-5 w-5 mr-3" />
            <span>Daftar Laporan</span>
          </Link>
          <Link href="/dinkes/inventory" className={getLinkClass("/dinkes/inventory")}>
            <Package className="h-5 w-5 mr-3" />
            <span>Inventaris</span>
          </Link>
          <Link href="/dinkes/map" className={getLinkClass("/dinkes/map")}>
            <Map className="h-5 w-5 mr-3" />
            <span>Peta GIS</span>
          </Link>
          <Link href="/dinkes/klb" className={getLinkClass("/dinkes/klb")}>
            <AlertTriangle className="h-5 w-5 mr-3" />
            <span>Status KLB</span>
          </Link>
          <Link href="/dinkes/sla" className={getLinkClass("/dinkes/sla")}>
            <Clock className="h-5 w-5 mr-3" />
            <span>SLA Monitoring</span>
          </Link>
          <Link href="/dinkes/puskesmas" className={getLinkClass("/dinkes/puskesmas")}>
            <Building2 className="h-5 w-5 mr-3" />
            <span>Akun Puskesmas</span>
          </Link>
          <Link href="/dinkes/settings" className={getLinkClass("/dinkes/settings")}>
            <Settings className="h-5 w-5 mr-3" />
            <span>Pengaturan</span>
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
          <h1 className="text-xl font-bold text-text">DBD Guard Portal</h1>
          <div className="flex items-center gap-4">
            <ThemeToggle />
            <div className="flex flex-col items-end ml-2">
              <span className="text-sm font-bold text-text">Administrator</span>
              <span className="text-xs text-text-muted">Dinas Kesehatan</span>
            </div>
            <div className="h-10 w-10 bg-primary-tint rounded-full flex items-center justify-center text-primary font-bold">
              ADM
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
