"use client";

import { useState, useEffect } from "react";
import api from "@/lib/api";
import { Users, Plus, Building2, User, Search, RefreshCw, Trash2, ShieldCheck, Phone, X, Edit } from "lucide-react";

export default function UserManagementPage() {
  const [dinkesAccounts, setDinkesAccounts] = useState<any[]>([]);
  const [puskesmasAccounts, setPuskesmasAccounts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("dinkes");

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditMode, setIsEditMode] = useState(false);
  const [editUserId, setEditUserId] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    contact: "",
    password: "",
    role: "ROLE_DINKES", // default
    puskesmasId: "",
    activeStatus: true,
  });

  const fetchAccounts = async () => {
    setLoading(true);
    try {
      const [dinkesRes, puskesmasRes] = await Promise.all([
        api.get("/developer/dinkes-accounts"),
        api.get("/admin/puskesmas-accounts")
      ]);
      setDinkesAccounts(dinkesRes.data || []);
      setPuskesmasAccounts(puskesmasRes.data || []);
    } catch (err) {
      console.error("Failed to load accounts", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAccounts();
  }, []);

  const handleDelete = async (id: string, roleTab: string) => {
    if (!confirm("Apakah Anda yakin ingin menghapus akun ini?")) return;
    try {
      if (roleTab === "dinkes") {
        await api.delete(`/developer/dinkes-accounts/${id}`);
      } else {
        await api.delete(`/admin/puskesmas-accounts/${id}`);
      }
      alert("Akun berhasil dihapus!");
      fetchAccounts();
    } catch (err: any) {
      alert("Gagal menghapus akun: " + (err.response?.data?.error || err.message));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      const payload: any = {
        name: formData.name,
        email: formData.email,
        contact: formData.contact,
        password: formData.password, // backend will ignore if empty during update
      };

      if (isEditMode) {
        payload.active_status = formData.activeStatus;
      }

      if (activeTab === "puskesmas") {
        payload.role = formData.role;
        payload.puskesmas_id = formData.puskesmasId;
      }

      if (isEditMode) {
        if (activeTab === "dinkes") {
          await api.patch(`/developer/dinkes-accounts/${editUserId}`, payload);
        } else {
          await api.patch(`/admin/puskesmas-accounts/${editUserId}`, payload);
        }
        alert("Akun berhasil diperbarui!");
      } else {
        if (activeTab === "dinkes") {
          await api.post("/developer/dinkes-accounts", payload);
        } else {
          await api.post("/admin/puskesmas-accounts", payload);
        }
        alert("Akun berhasil dibuat!");
      }
      setIsModalOpen(false);
      fetchAccounts();
    } catch (err: any) {
      alert(`Gagal ${isEditMode ? "memperbarui" : "membuat"} akun: ` + (err.response?.data?.error || err.message));
    } finally {
      setIsSubmitting(false);
    }
  };

  const openAddModal = () => {
    setIsEditMode(false);
    setEditUserId("");
    setFormData({ 
      name: "", email: "", contact: "", password: "",
      role: activeTab === "dinkes" ? "ROLE_DINKES" : "ROLE_PETUGAS_PUSKESMAS",
      puskesmasId: "",
      activeStatus: true
    });
    setIsModalOpen(true);
  };

  const openEditModal = (acc: any) => {
    setIsEditMode(true);
    setEditUserId(acc.id);
    setFormData({ 
      name: acc.name || "", 
      email: acc.email || "", 
      contact: acc.contact || "", 
      password: "", // do not fill password
      role: acc.role?.name || (activeTab === "dinkes" ? "ROLE_DINKES" : "ROLE_PETUGAS_PUSKESMAS"),
      puskesmasId: acc.location_id || "",
      activeStatus: acc.active_status !== false // default true unless explicitly false
    });
    setIsModalOpen(true);
  };

  return (
    <div className="w-full max-w-7xl mx-auto space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-text">Kelola Akun Sistem</h2>
          <p className="text-text-muted">Manajemen pengguna global (Dinkes & Puskesmas).</p>
        </div>
        <div className="flex items-center gap-3">
          <button 
            onClick={fetchAccounts}
            className="p-2.5 bg-surface border border-border rounded-xl text-text-muted hover:text-primary transition-colors"
          >
            <RefreshCw className={`w-5 h-5 ${loading ? "animate-spin" : ""}`} />
          </button>
          <button onClick={openAddModal} className="flex items-center gap-2 bg-primary hover:bg-primary-deep text-primary-foreground px-4 py-2.5 rounded-xl font-medium transition-colors shadow-sm">
            <Plus className="w-5 h-5" />
            Tambah Akun
          </button>
        </div>
      </div>

      <div className="bg-surface rounded-2xl border border-border overflow-hidden">
        <div className="flex border-b border-border">
          <button
            onClick={() => setActiveTab("dinkes")}
            className={`flex items-center gap-2 px-6 py-4 font-semibold text-sm transition-colors relative ${
              activeTab === "dinkes" ? "text-primary" : "text-text-muted hover:text-text"
            }`}
          >
            <ShieldCheck className="w-4 h-4" />
            Akun Dinkes
            {activeTab === "dinkes" && <div className="absolute bottom-0 left-0 w-full h-0.5 bg-primary" />}
          </button>
          <button
            onClick={() => setActiveTab("puskesmas")}
            className={`flex items-center gap-2 px-6 py-4 font-semibold text-sm transition-colors relative ${
              activeTab === "puskesmas" ? "text-primary" : "text-text-muted hover:text-text"
            }`}
          >
            <Building2 className="w-4 h-4" />
            Akun Puskesmas
            {activeTab === "puskesmas" && <div className="absolute bottom-0 left-0 w-full h-0.5 bg-primary" />}
          </button>
        </div>

        <div className="p-4 border-b border-border flex justify-between items-center bg-bg/50">
          <div className="relative w-64">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search className="h-4 w-4 text-text-muted" />
            </div>
            <input
              type="text"
              placeholder="Cari akun..."
              className="block w-full pl-10 pr-3 py-2 border border-border rounded-lg leading-5 bg-surface text-text placeholder-text-muted focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary sm:text-sm"
            />
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-border">
            <thead className="bg-bg">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Nama & Email</th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Role</th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Kontak</th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Status</th>
                <th className="px-6 py-3 text-right text-xs font-semibold text-text-muted uppercase tracking-wider">Aksi</th>
              </tr>
            </thead>
            <tbody className="bg-surface divide-y divide-border">
              {loading ? (
                <tr key="loading">
                  <td colSpan={5} className="px-6 py-8 text-center text-text-muted">Memuat data...</td>
                </tr>
              ) : activeTab === "dinkes" && dinkesAccounts.length === 0 ? (
                <tr key="empty-dinkes">
                  <td colSpan={5} className="px-6 py-8 text-center text-text-muted">Tidak ada akun Dinkes.</td>
                </tr>
              ) : activeTab === "puskesmas" && puskesmasAccounts.length === 0 ? (
                <tr key="empty-puskesmas">
                  <td colSpan={5} className="px-6 py-8 text-center text-text-muted">Tidak ada akun Puskesmas.</td>
                </tr>
              ) : (
                (activeTab === "dinkes" ? dinkesAccounts : puskesmasAccounts).map((acc, idx) => (
                  <tr key={acc.id || `user-${idx}`} className="hover:bg-bg/50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="h-10 w-10 shrink-0 bg-primary-tint rounded-full flex items-center justify-center text-primary">
                          <User className="h-5 w-5" />
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-semibold text-text">{acc.name}</div>
                          <div className="text-sm text-text-muted">{acc.email}</div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-status-verified/20 text-status-verified">
                        {acc.role?.name || (activeTab === "dinkes" ? "ROLE_DINKES" : "PUSKESMAS")}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-text-muted">
                      <div className="flex items-center gap-1">
                        <Phone className="w-3 h-3" />
                        {acc.contact || "-"}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {acc.active_status ? (
                        <span className="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-primary-tint text-primary">Aktif</span>
                      ) : (
                        <span className="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">Nonaktif</span>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div className="flex justify-end gap-2">
                        <button 
                          onClick={() => openEditModal(acc)}
                          className="text-primary hover:text-primary bg-primary-tint p-2 rounded-lg hover:bg-primary-tint transition-colors"
                          title="Edit Akun"
                        >
                          <Edit className="w-4 h-4" />
                        </button>
                        <button 
                          onClick={() => handleDelete(acc.id, activeTab)}
                          className="text-red-600 hover:text-red-900 bg-red-50 p-2 rounded-lg hover:bg-red-100 transition-colors"
                          title="Hapus Akun"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Add Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
          <div className="bg-surface rounded-2xl border border-border shadow-xl w-full max-w-lg overflow-hidden animate-in fade-in zoom-in-95 duration-200 max-h-[90vh] flex flex-col">
            <div className="flex justify-between items-center p-6 border-b border-border bg-bg/50">
              <h3 className="text-lg font-bold text-text">
                {isEditMode ? "Edit Akun" : "Tambah Akun"} {activeTab === "dinkes" ? "Dinas Kesehatan" : "Puskesmas"}
              </h3>
              <button 
                onClick={() => setIsModalOpen(false)}
                className="text-text-muted hover:text-text hover:bg-bg p-2 rounded-lg transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            <form id="user-account-form" onSubmit={handleSubmit} className="p-6 space-y-4 overflow-y-auto flex-1">
              <div>
                <label className="block text-sm font-medium text-text mb-1">Nama Lengkap</label>
                <input 
                  type="text" 
                  required
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  className="w-full px-4 py-2 bg-bg border border-border rounded-lg text-text focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                  placeholder="Misal: Budi Santoso"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-text mb-1">Email</label>
                <input 
                  type="email" 
                  required
                  value={formData.email}
                  onChange={(e) => setFormData({...formData, email: e.target.value})}
                  className="w-full px-4 py-2 bg-bg border border-border rounded-lg text-text focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                  placeholder="budi@example.com"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-text mb-1">Nomor Kontak</label>
                  <input 
                    type="text" 
                    value={formData.contact}
                    onChange={(e) => setFormData({...formData, contact: e.target.value})}
                    className="w-full px-4 py-2 bg-bg border border-border rounded-lg text-text focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                    placeholder="0812xxxx"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-text mb-1">Password {isEditMode && "(Opsional)"}</label>
                  <input 
                    type="password" 
                    required={!isEditMode}
                    value={formData.password}
                    onChange={(e) => setFormData({...formData, password: e.target.value})}
                    className="w-full px-4 py-2 bg-bg border border-border rounded-lg text-text focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                    placeholder={isEditMode ? "Ketik password baru jika ingin mengubah" : "Minimal 6 karakter"}
                  />
                  {isEditMode && (
                    <p className="text-xs text-text-muted mt-1">Password sebelumnya dienkripsi dan tidak dapat dilihat.</p>
                  )}
                </div>
              </div>

              {isEditMode && (
                <div>
                  <label className="block text-sm font-medium text-text mb-1">Status Akun</label>
                  <select 
                    value={formData.activeStatus ? "true" : "false"}
                    onChange={(e) => setFormData({...formData, activeStatus: e.target.value === "true"})}
                    className="w-full px-4 py-2 bg-bg border border-border rounded-lg text-text focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                  >
                    <option value="true">Aktif</option>
                    <option value="false">Nonaktif</option>
                  </select>
                </div>
              )}

              {activeTab === "puskesmas" && (
                <>
                  <div>
                    <label className="block text-sm font-medium text-text mb-1">Role Akun</label>
                    <select 
                      required
                      value={formData.role}
                      onChange={(e) => setFormData({...formData, role: e.target.value})}
                      className="w-full px-4 py-2 bg-bg border border-border rounded-lg text-text focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                    >
                      <option value="ROLE_PETUGAS_PUSKESMAS">Petugas Puskesmas</option>
                      <option value="ROLE_KEPALA_PUSKESMAS">Kepala Puskesmas</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text mb-1">Puskesmas ID (UUID)</label>
                    <input 
                      type="text" 
                      required
                      value={formData.puskesmasId}
                      onChange={(e) => setFormData({...formData, puskesmasId: e.target.value})}
                      className="w-full px-4 py-2 bg-bg border border-border rounded-lg text-text focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                      placeholder="Masukkan UUID Puskesmas"
                    />
                    <p className="text-xs text-text-muted mt-1">Harus berupa format UUID Puskesmas yang valid.</p>
                  </div>
                </>
              )}
            </form>

            <div className="p-4 border-t border-border bg-bg/50 flex justify-end gap-3 mt-auto">
              <button 
                type="button" 
                onClick={() => setIsModalOpen(false)}
                className="px-4 py-2 text-text-muted hover:text-text font-medium"
              >
                Batal
              </button>
              <button 
                type="submit" 
                form="user-account-form"
                disabled={isSubmitting}
                className="px-6 py-2 bg-primary hover:bg-primary-deep text-primary-foreground rounded-lg font-medium transition-colors disabled:opacity-70 flex items-center gap-2"
              >
                {isSubmitting ? "Menyimpan..." : "Simpan Akun"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
