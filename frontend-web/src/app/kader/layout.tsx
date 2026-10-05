"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { useState } from "react";
import { LayoutDashboard, CheckSquare, LogOut, ShieldCheck } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";
import NotificationPermission from "@/components/NotificationPermission";

export default function KaderLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const router = useRouter();
  const pathname = usePathname();
  
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [userName, setUserName] = useState<string>("Loading...");

  const handleLogout = () => {
    Cookies.remove("auth_token");
    router.push("/");
  };

  const getLinkClass = (path: string) => {
    const isActive = path === "/kader" 
      ? pathname === "/kader"
      : pathname.startsWith(path);
      
    return `flex items-center space-x-3 px-4 py-3 rounded-xl transition-all duration-300 group ${
      isActive
        ? "bg-gradient-to-r from-teal-500 to-emerald-500 text-primary-foreground shadow-md shadow-primary/20 font-medium"
        : "text-text-muted hover:bg-bg hover:text-primary"
    }`;
  };

  return (
    <div className="min-h-screen bg-bg flex">
      <NotificationPermission />
      {/* Sidebar - Desktop */}
      <aside className="hidden lg:flex flex-col w-72 bg-surface border-r border-border fixed h-full z-10 transition-all duration-300">
        <div className="p-6 flex items-center space-x-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-teal-500 to-emerald-500 flex items-center justify-center text-primary-foreground font-bold text-xl shadow-lg shadow-primary/20">
            K
          </div>
          <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-slate-800 to-slate-600">DBD Guard</span>
        </div>

        <nav className="flex-1 px-4 py-6 space-y-2 overflow-y-auto">
          <Link href="/kader" className={getLinkClass("/kader")}>
            <LayoutDashboard className={`w-5 h-5 transition-transform duration-300 group-hover:scale-110 ${pathname === "/kader" ? "text-white" : ""}`} />
            <span>Dashboard Kader</span>
          </Link>
          <Link href="/kader/verification" className={getLinkClass("/kader/verification")}>
            <CheckSquare className={`w-5 h-5 transition-transform duration-300 group-hover:scale-110 ${pathname.startsWith("/kader/verification") ? "text-white" : ""}`} />
            <span>Penyelidikan Epidemiologi</span>
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
      <div className="flex-1 flex flex-col lg:ml-72 overflow-hidden">
        {/* Header */}
        <header className="h-16 bg-surface border-b border-border flex items-center justify-between px-8 z-10">
          <h1 className="text-xl font-bold text-text">Portal Verifikasi Lapangan</h1>
          <div className="flex items-center gap-4">
            <ThemeToggle />
            <div className="flex flex-col items-end ml-2">
              <span className="text-sm font-bold text-text">Perangkat RT/RW</span>
              <span className="text-xs text-text-muted">Verifikator</span>
            </div>
            <div className="h-10 w-10 bg-primary-tint rounded-full flex items-center justify-center text-primary font-bold">
              RT
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
