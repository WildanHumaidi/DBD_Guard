"use client";

import { useState, useEffect } from "react";
import { Settings2, Save, Loader2, CheckCircle2, ShieldAlert } from "lucide-react";
import api from "@/lib/api";

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState("klb");
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  const [klbConfig, setKlbConfig] = useState({
    region_scope: "KELURAHAN",
    min_cases: 5,
    timeframe_days: 7,
    is_auto_trigger: true
  });

  useEffect(() => {
    const fetchConfig = async () => {
      setLoading(true);
      try {
        const response = await api.get("/admin/klb/config");
        if (response.data?.data) {
          setKlbConfig(response.data.data);
        }
      } catch (err) {
        console.error("Failed to fetch KLB config", err);
      } finally {
        setLoading(false);
      }
    };
    fetchConfig();
  }, []);

  const handleSaveKLB = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setSaved(false);
    
    try {
      await api.patch("/admin/klb/config", klbConfig);
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err) {
      console.error("Failed to update KLB config", err);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="max-w-6xl mx-auto animate-in fade-in slide-in-from-bottom-4 duration-500">
      
      <div className="mb-8">
        <h1 className="text-3xl font-extrabold text-text tracking-tight">Pengaturan Sistem</h1>
        <p className="text-text-muted mt-2">Kelola parameter deteksi KLB dan preferensi sistem DBD Guard.</p>
      </div>

      <div className="flex flex-col md:flex-row gap-8">
        {/* Settings Sidebar */}
        <div className="w-full md:w-64 flex-shrink-0">
          <div className="bg-surface rounded-2xl shadow-sm border border-border p-2 overflow-hidden sticky top-4">
            <nav className="flex flex-col space-y-1">
              <button
                onClick={() => setActiveTab("klb")}
                className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all text-left font-medium ${
                  activeTab === "klb" 
                    ? "bg-red-50 text-red-600 border border-red-100" 
                    : "text-text-muted hover:bg-bg hover:text-red-600"
                }`}
              >
                <ShieldAlert className="w-5 h-5" />
                Pengaturan KLB
              </button>
              
              <button
                onClick={() => setActiveTab("general")}
                className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all text-left font-medium ${
                  activeTab === "general" 
                    ? "bg-primary-tint text-primary" 
                    : "text-text-muted hover:bg-bg hover:text-primary"
                }`}
              >
                <Settings2 className="w-5 h-5" />
                Sistem Umum
              </button>
            </nav>
          </div>
        </div>

        {/* Settings Content */}
        <div className="flex-1">
          <div className="bg-surface rounded-2xl shadow-sm border border-border overflow-hidden relative min-h-[500px]">
            
            {/* Header Area Gradient */}
            <div className={`h-32 absolute top-0 w-full left-0 opacity-10 ${activeTab === 'klb' ? 'bg-gradient-to-r from-red-500 to-orange-600' : 'bg-gradient-to-r from-emerald-500 to-teal-600'}`}></div>
            
            <div className="p-8 relative z-10">
              
              {/* KLB Tab */}
              {activeTab === "klb" && (
                <div className="animate-in fade-in slide-in-from-right-4">
                  <div className="flex items-center gap-4 mb-8">
                    <div className="h-16 w-16 bg-red-100 rounded-full flex items-center justify-center text-red-600 shadow-sm border border-red-200">
                      <ShieldAlert className="w-8 h-8" />
                    </div>
                    <div>
                      <h2 className="text-2xl font-bold text-text">Ambang Batas Kejadian Luar Biasa (KLB)</h2>
                      <p className="text-text-muted">Atur kapan sistem akan memicu status KLB secara otomatis.</p>
                    </div>
                  </div>

                  {loading ? (
                    <div className="flex items-center justify-center py-20">
                      <Loader2 className="w-8 h-8 animate-spin text-red-500" />
                    </div>
                  ) : (
                    <form onSubmit={handleSaveKLB} className="space-y-6">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-2">
                          <label className="text-sm font-semibold text-text">Lingkup Wilayah (Region Scope)</label>
                          <select
                            value={klbConfig.region_scope}
                            onChange={(e) => setKlbConfig({...klbConfig, region_scope: e.target.value})}
                            className="w-full px-4 py-3 bg-bg border border-border rounded-xl focus:ring-2 focus:ring-red-500/50 focus:border-red-500 transition-colors"
                          >
                            <option value="KELURAHAN">Per Kelurahan</option>
                            <option value="KECAMATAN">Per Kecamatan</option>
                            <option value="KOTA">Tingkat Kota</option>
                          </select>
                        </div>
                        
                        <div className="space-y-2">
                          <label className="text-sm font-semibold text-text">Minimal Kasus (Min Cases)</label>
                          <input
                            type="number"
                            min="1"
                            value={klbConfig.min_cases}
                            onChange={(e) => setKlbConfig({...klbConfig, min_cases: parseInt(e.target.value) || 0})}
                            className="w-full px-4 py-3 bg-bg border border-border rounded-xl focus:ring-2 focus:ring-red-500/50 focus:border-red-500 transition-colors"
                          />
                        </div>
                        
                        <div className="space-y-2">
                          <label className="text-sm font-semibold text-text">Jangka Waktu (Hari)</label>
                          <input
                            type="number"
                            min="1"
                            value={klbConfig.timeframe_days}
                            onChange={(e) => setKlbConfig({...klbConfig, timeframe_days: parseInt(e.target.value) || 0})}
                            className="w-full px-4 py-3 bg-bg border border-border rounded-xl focus:ring-2 focus:ring-red-500/50 focus:border-red-500 transition-colors"
                          />
                        </div>
                      </div>

                      <div className="flex items-center justify-between p-5 mt-4 rounded-2xl border border-red-100 bg-red-50">
                        <div>
                          <h4 className="font-bold text-red-800">Auto-Trigger KLB</h4>
                          <p className="text-sm text-red-600 mt-1">Aktifkan peringatan dan instruksi fogging massal otomatis bila batas terlampaui.</p>
                        </div>
                        <label className="relative inline-flex items-center cursor-pointer">
                          <input type="checkbox" className="sr-only peer" checked={klbConfig.is_auto_trigger} 
                            onChange={(e) => setKlbConfig({...klbConfig, is_auto_trigger: e.target.checked})}
                          />
                          <div className="w-11 h-6 bg-red-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-red-600"></div>
                        </label>
                      </div>

                      <div className="pt-6 border-t border-border flex items-center justify-between">
                        {saved ? (
                          <div className="flex items-center text-emerald-600 gap-2 animate-in fade-in">
                            <CheckCircle2 className="w-5 h-5" />
                            <span className="text-sm font-medium">Konfigurasi KLB disimpan!</span>
                          </div>
                        ) : (
                          <div></div>
                        )}
                        <button
                          type="submit"
                          disabled={saving}
                          className="bg-red-600 hover:bg-red-700 text-white font-semibold py-3 px-6 rounded-xl shadow-md hover:shadow-lg hover:shadow-red-500/20 transition-all flex items-center gap-2"
                        >
                          {saving ? <Loader2 className="w-5 h-5 animate-spin" /> : <Save className="w-5 h-5" />}
                          Simpan Pengaturan
                        </button>
                      </div>
                    </form>
                  )}
                </div>
              )}

              {/* General Tab */}
              {activeTab === "general" && (
                <div className="animate-in fade-in slide-in-from-right-4">
                  <div className="flex items-center gap-4 mb-8">
                    <div className="h-16 w-16 bg-primary-tint rounded-full flex items-center justify-center text-primary shadow-sm border border-emerald-100">
                      <Settings2 className="w-8 h-8" />
                    </div>
                    <div>
                      <h2 className="text-2xl font-bold text-text">Pengaturan Sistem Umum</h2>
                      <p className="text-text-muted">Pengaturan global lainnya.</p>
                    </div>
                  </div>

                  <div className="py-10 text-center border-2 border-dashed border-border rounded-xl">
                    <p className="text-text-muted font-medium">Belum ada pengaturan tambahan yang diperlukan saat ini.</p>
                  </div>
                </div>
              )}

            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
