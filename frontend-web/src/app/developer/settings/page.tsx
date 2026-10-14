"use client";

import { useState, useEffect, useRef } from "react";
import { Save, Image as ImageIcon, Globe, FileText, X } from "lucide-react";
import api from "@/lib/api";

export default function SettingsPage() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [heroTitle, setHeroTitle] = useState("");
  const [announcementText, setAnnouncementText] = useState("");
  const [isCitizenPortalActive, setIsCitizenPortalActive] = useState(true);
  const [bannerImagePath, setBannerImagePath] = useState("");
  const [bannerImageFile, setBannerImageFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchSettings = async () => {
    try {
      setLoading(true);
      const res = await api.get("/developer/settings");
      const data = res.data;
      setHeroTitle(data.HeroTitle || "");
      setAnnouncementText(data.AnnouncementText || "");
      setIsCitizenPortalActive(data.IsCitizenPortalActive !== false); // default true
      setBannerImagePath(data.BannerImagePath || "");
    } catch (err) {
      console.error("Gagal mengambil pengaturan", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ignore = false;
    api.get("/developer/settings")
      .then((res) => {
        if (!ignore && res.data) {
          const data = res.data;
          setHeroTitle(data.HeroTitle || "");
          setAnnouncementText(data.AnnouncementText || "");
          setIsCitizenPortalActive(data.IsCitizenPortalActive !== false);
          setBannerImagePath(data.BannerImagePath || "");
        }
      })
      .catch((err) => {
        console.error("Gagal mengambil pengaturan", err);
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const formData = new FormData();
      formData.append("hero_title", heroTitle);
      formData.append("announcement_text", announcementText);
      formData.append("is_citizen_portal_active", isCitizenPortalActive ? "true" : "false");
      
      if (bannerImageFile) {
        formData.append("banner_image", bannerImageFile);
      }

      await api.put("/developer/settings", formData, {
        headers: {
          "Content-Type": "multipart/form-data"
        }
      });
      alert("Pengaturan berhasil disimpan!");
      fetchSettings();
    } catch (err: unknown) {
      const errorMsg = (err as { response?: { data?: { error?: string } }; message?: string })?.response?.data?.error || (err as Error)?.message || "Terjadi kesalahan";
      alert("Gagal menyimpan pengaturan: " + errorMsg);
    } finally {
      setSaving(false);
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setBannerImageFile(e.target.files[0]);
    }
  };

  if (loading) {
    return <div className="text-center py-10 text-text-muted">Memuat pengaturan...</div>;
  }

  // Determine image preview URL
  let previewUrl = "";
  if (bannerImageFile) {
    previewUrl = URL.createObjectURL(bannerImageFile);
  } else if (bannerImagePath) {
    const baseUrl = process.env.NEXT_PUBLIC_API_URL ? process.env.NEXT_PUBLIC_API_URL.replace('/api/v1', '') : 'http://localhost:8080';
    previewUrl = baseUrl + bannerImagePath;
  }

  return (
    <div className="w-full max-w-6xl mx-auto space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-text">Pengaturan Website</h2>
        <p className="text-text-muted">Sesuaikan tampilan publik dan operasional sistem DBD Guard.</p>
      </div>

      <form onSubmit={handleSave} className="space-y-8">
        
        {/* Banner Section */}
        <div className="bg-surface rounded-2xl border border-border shadow-sm overflow-hidden">
          <div className="border-b border-border p-6 flex items-center gap-3">
            <div className="p-2 bg-status-verified/20 text-status-verified rounded-lg">
              <ImageIcon className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-bold text-text">Banner Beranda</h3>
              <p className="text-sm text-text-muted">Ubah gambar pahlawan (hero image) di halaman publik.</p>
            </div>
          </div>
          <div className="p-6">
            <input 
              type="file" 
              accept="image/png, image/jpeg, image/jpg" 
              className="hidden" 
              ref={fileInputRef}
              onChange={handleFileSelect}
            />
            {previewUrl ? (
              <div className="relative rounded-xl overflow-hidden border border-border group">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src={previewUrl} alt="Banner Preview" className="w-full h-48 object-cover" />
                <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-4">
                  <button 
                    type="button" 
                    onClick={() => fileInputRef.current?.click()}
                    className="px-4 py-2 bg-surface text-text rounded-lg font-medium text-sm"
                  >
                    Ganti Gambar
                  </button>
                  <button 
                    type="button" 
                    onClick={() => setBannerImageFile(null)}
                    className="p-2 bg-red-500 text-white rounded-lg"
                  >
                    <X className="w-5 h-5" />
                  </button>
                </div>
              </div>
            ) : (
              <div 
                onClick={() => fileInputRef.current?.click()}
                className="border-2 border-dashed border-border rounded-xl p-8 text-center bg-bg/50 hover:bg-bg cursor-pointer transition-colors"
              >
                <ImageIcon className="w-10 h-10 text-text-muted mx-auto mb-3" />
                <p className="text-sm font-medium text-text">Klik untuk mengunggah gambar</p>
                <p className="text-xs text-text-muted mt-1">Format: JPG, PNG. Rekomendasi ukuran: 1920x1080px (Maks 5MB)</p>
              </div>
            )}
          </div>
        </div>

        {/* Global Text Section */}
        <div className="bg-surface rounded-2xl border border-border shadow-sm overflow-hidden">
          <div className="border-b border-border p-6 flex items-center gap-3">
            <div className="p-2 bg-primary-tint text-primary rounded-lg">
              <FileText className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-bold text-text">Konten Statis</h3>
              <p className="text-sm text-text-muted">Ubah teks pengumuman atau deskripsi di portal warga.</p>
            </div>
          </div>
          <div className="p-6 space-y-5">
            <div>
              <label className="block text-sm font-medium text-text mb-2">Judul Utama Halaman Beranda</label>
              <input 
                type="text" 
                required
                value={heroTitle}
                onChange={(e) => setHeroTitle(e.target.value)}
                className="w-full px-4 py-2 bg-bg border border-border rounded-lg text-text focus:outline-none focus:border-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-text mb-2">Teks Pengumuman (Tampil di atas layar)</label>
              <textarea 
                value={announcementText}
                onChange={(e) => setAnnouncementText(e.target.value)}
                className="w-full px-4 py-3 bg-bg border border-border rounded-lg text-text focus:outline-none focus:border-primary min-h-[100px]"
              />
            </div>
          </div>
        </div>

        {/* System Toggle Section */}
        <div className="bg-surface rounded-2xl border border-border shadow-sm overflow-hidden">
          <div className="border-b border-border p-6 flex items-center gap-3">
            <div className="p-2 bg-amber-50 text-amber-600 rounded-lg">
              <Globe className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-bold text-text">Status Sistem</h3>
              <p className="text-sm text-text-muted">Atur status portal laporan warga.</p>
            </div>
          </div>
          <div className="p-6 flex items-center justify-between">
            <div>
              <p className="font-semibold text-text">Portal Warga Aktif</p>
              <p className="text-sm text-text-muted">Warga dapat mengirimkan laporan baru melalui website.</p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input 
                type="checkbox" 
                checked={isCitizenPortalActive}
                onChange={(e) => setIsCitizenPortalActive(e.target.checked)}
                className="sr-only peer" 
              />
              <div className="w-11 h-6 bg-slate-300 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary dark:peer-focus:ring-emerald-800 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-surface after:border-border after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
            </label>
          </div>
        </div>

        <div className="flex justify-end">
          <button 
            type="submit"
            disabled={saving}
            className="flex items-center gap-2 bg-primary hover:bg-primary-deep text-primary-foreground px-6 py-3 rounded-xl font-medium transition-colors shadow-sm disabled:opacity-70"
          >
            {saving ? "Menyimpan..." : (
              <>
                <Save className="w-5 h-5" />
                Simpan Konfigurasi
              </>
            )}
          </button>
        </div>

      </form>
    </div>
  );
}
