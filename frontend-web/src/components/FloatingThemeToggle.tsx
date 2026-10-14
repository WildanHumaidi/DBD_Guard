"use client";

import { useTheme } from "next-themes";
import { useSyncExternalStore } from "react";
import { Sun, Moon } from "lucide-react";
import { usePathname } from "next/navigation";

const emptySubscribe = () => () => {};

export function FloatingThemeToggle() {
  const mounted = useSyncExternalStore(emptySubscribe, () => true, () => false);
  const { theme, setTheme } = useTheme();
  const pathname = usePathname();

  // Hide on dashboard routes since they might have their own header toggles
  const isDashboardRoute = pathname?.match(/^\/(puskesmas|kader|fogging|dinkes|developer)/);

  if (!mounted || isDashboardRoute) {
    return null;
  }

  const isDark = theme === "dark";

  return (
    <button
      onClick={() => setTheme(isDark ? "light" : "dark")}
      className="fixed bottom-6 right-6 z-50 flex h-12 w-12 items-center justify-center rounded-full bg-surface text-text-muted hover:text-text hover:bg-bg border border-border shadow-xl hover:shadow-primary/20 transition-all duration-300 overflow-hidden active:scale-95 group"
      aria-label="Toggle Dark Mode"
      title={isDark ? "Switch to Light Mode" : "Switch to Dark Mode"}
    >
      <div className={`absolute inset-0 flex items-center justify-center transition-transform duration-500 ease-[cubic-bezier(0.4,0,0.2,1)] ${isDark ? 'rotate-90 scale-0 opacity-0' : 'rotate-0 scale-100 opacity-100'}`}>
        <Sun className="h-6 w-6 text-amber-500 group-hover:text-amber-600 transition-colors" />
      </div>
      <div className={`absolute inset-0 flex items-center justify-center transition-transform duration-500 ease-[cubic-bezier(0.4,0,0.2,1)] ${isDark ? 'rotate-0 scale-100 opacity-100' : '-rotate-90 scale-0 opacity-0'}`}>
        <Moon className="h-6 w-6 text-indigo-400 group-hover:text-indigo-300 transition-colors" />
      </div>
    </button>
  );
}
