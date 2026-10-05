"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { LayoutDashboard, Calendar, LogOut, Flame } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";

export default function FoggingLayout({
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

  const getLinkClass = (path: string) => {
    const isActive = path === "/fogging" 
      ? pathname === "/fogging"
      : pathname.startsWith(path);
      
    if (isActive) {
      return "flex items-center px-4 py-3 bg-status-scheduled/20 text-status-scheduled rounded-xl transition-colors font-bold";
    }
    return "flex items-center px-4 py-3 text-text-muted hover:bg-bg hover:text-status-scheduled rounded-xl transition-colors font-medium";
  };

  return (
    <div className="flex h-screen bg-bg overflow-hidden">
      {/* Sidebar */}
      <aside className="w-64 bg-surface border-r border-border flex flex-col">
        <div className="h-16 flex items-center px-6 border-b border-border">
          <Flame className="h-6 w-6 text-status-scheduled mr-2" />
          <span className="font-bold text-lg text-text">Tim Fogging</span>
        </div>
        
        <nav className="flex-1 py-6 px-4 space-y-2">
          <Link href="/fogging" className={getLinkClass("/fogging")}>
            <LayoutDashboard className="h-5 w-5 mr-3" />
            <span>Beranda</span>
          </Link>
          <Link href="/fogging/schedules" className={getLinkClass("/fogging/schedules")}>
            <Calendar className="h-5 w-5 mr-3" />
            <span>Jadwal Fogging</span>
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
          <h1 className="text-xl font-bold text-text">Portal Tim Eksekutor</h1>
          <div className="flex items-center gap-4">
            <ThemeToggle />
            <div className="flex flex-col items-end ml-2">
              <span className="text-sm font-bold text-text">Tim Fogging</span>
              <span className="text-xs text-text-muted">Eksekutor Lapangan</span>
            </div>
            <div className="h-10 w-10 bg-status-scheduled/20 rounded-full flex items-center justify-center text-status-scheduled font-bold">
              FG
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
