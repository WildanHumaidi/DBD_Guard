"use client";

import { useState, useEffect } from "react";
import { Plus, Edit, Trash2, Loader2, Building2, UserCircle, Phone } from "lucide-react";
import api from "@/lib/api";
import EmptyState from "@/components/common/EmptyState";
import SkeletonTable from "@/components/common/SkeletonTable";
import ErrorMessage from "@/components/common/ErrorMessage";

interface UserAccount {
  user_id: string;
  name: string;
  email: string;
  contact: string;
  active_status: boolean;
}

export default function PuskesmasManagementPage() {
  const [accounts, setAccounts] = useState<UserAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  
  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<"ADD" | "UPDATE">("ADD");
  const [selectedAccount, setSelectedAccount] = useState<UserAccount | null>(null);
  
  // Form State
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    contact: "",
    password: "",
    active_status: true,
  });

  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      const res = await api.get("/admin/puskesmas-accounts");
      // ensure we get an array
      setAccounts(Array.isArray(res.data) ? res.data : []);
    } catch (err: unknown) {
      const errorMsg = (err as { response?: { data?: { error?: string } }; message?: string })?.response?.data?.error || (err as Error)?.message || "Gagal mengambil data";
      setError(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ignore = false;
    api.get("/admin/puskesmas-accounts")
      .then((res) => {
        if (!ignore) {
          setAccounts(Array.isArray(res.data) ? res.data : []);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          const errorMsg = (err as { response?: { data?: { error?: string } }; message?: string })?.response?.data?.error || (err as Error)?.message || "Gagal mengambil data";
          setError(errorMsg);
        }
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  const openAddModal = () => {
    setModalMode("ADD");
    setFormData({
      name: "",
      email: "",
      contact: "",
      password: "",
      active_status: true,
    });
    setSelectedAccount(null);
    setIsModalOpen(true);
  };

  const openUpdateModal = (account: UserAccount) => {
    setModalMode("UPDATE");
    setFormData({
      name: account.name,
      email: account.email,
      contact: account.contact,
      password: "", // intentionally left blank
      active_status: account.active_status,
    });
    setSelectedAccount(account);
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      if (modalMode === "ADD") {
        await api.post("/admin/puskesmas-accounts", {
          name: formData.name,
          email: formData.email,
          contact: formData.contact,
          password: formData.password,
        });
        alert("Akun Puskesmas berhasil ditambahkan.");
      } else if (modalMode === "UPDATE" && selectedAccount) {
        await api.patch(`/admin/puskesmas-accounts/${selectedAccount.user_id}`, {
          name: formData.name,
          email: formData.email,
          contact: formData.contact,
          active_status: formData.active_status,
        });
        alert("Data akun berhasil diupdate.");
      }
      setIsModalOpen(false);
      fetchData();
    } catch (err: unknown) {
      const errorMsg = (err as { response?: { data?: { error?: string } }; message?: string })?.response?.data?.error || (err as Error)?.message || "Terjadi kesalahan";
      alert(errorMsg);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`Anda yakin ingin menghapus akun puskesmas "${name}" secara permanen? Aksi ini tidak dapat dibatalkan.`)) return;
    try {
      await api.delete(`/admin/puskesmas-accounts/${id}`);
      alert("Akun berhasil dihapus");
      fetchData();
    } catch (err: unknown) {
      const errorMsg = (err as { response?: { data?: { error?: string } }; message?: string })?.response?.data?.error || (err as Error)?.message || "Gagal menghapus akun";
      alert(errorMsg);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-text">Manajemen Akun Puskesmas</h2>
          <p className="text-text-muted">Kelola pendaftaran, akses, dan informasi kontak petugas Puskesmas se-Kota Depok.</p>
        </div>
        <button 
          onClick={openAddModal}
          className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 transition-colors font-medium shadow-sm"
        >
          <Plus className="h-4 w-4" /> Tambah Akun
        </button>
      </div>

      {error && (
        <ErrorMessage message={error} onRetry={fetchData} />
      )}

      <div className="bg-surface border border-border rounded-2xl shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-bg border-b border-border text-text-muted text-sm">
                <th className="p-4 font-semibold">Nama Puskesmas / Petugas</th>
                <th className="p-4 font-semibold">Email</th>
                <th className="p-4 font-semibold">Kontak</th>
                <th className="p-4 font-semibold">Status</th>
                <th className="p-4 font-semibold text-center">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {loading ? (
                <tr>
                  <td colSpan={5} className="p-0">
                    <SkeletonTable rows={5} columns={5} />
                  </td>
                </tr>
              ) : accounts.length === 0 ? (
                <tr>
                  <td colSpan={5} className="p-8">
                    <EmptyState 
                      title="Belum ada akun" 
                      description="Belum ada akun Puskesmas yang terdaftar di sistem."
                      icon={<Building2 className="w-12 h-12" />} 
                    />
                  </td>
                </tr>
              ) : (
                accounts.map((account) => (
                  <tr key={account.user_id} className="hover:bg-bg transition-colors">
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        <div className="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
                          <Building2 className="w-5 h-5" />
                        </div>
                        <span className="font-medium text-text">{account.name}</span>
                      </div>
                    </td>
                    <td className="p-4 text-text-muted">
                      <div className="flex items-center gap-2">
                        <UserCircle className="w-4 h-4 text-text-muted" />
                        {account.email}
                      </div>
                    </td>
                    <td className="p-4 text-text-muted">
                      <div className="flex items-center gap-2">
                        <Phone className="w-4 h-4 text-text-muted" />
                        {account.contact || "-"}
                      </div>
                    </td>
                    <td className="p-4">
                      {account.active_status ? (
                        <span className="px-2 py-1 bg-primary-tint text-primary text-xs font-bold rounded-full">Aktif</span>
                      ) : (
                        <span className="px-2 py-1 bg-bg text-text-muted text-xs font-bold rounded-full">Nonaktif</span>
                      )}
                    </td>
                    <td className="p-4 text-center">
                      <div className="flex justify-center gap-2">
                        <button 
                          onClick={() => openUpdateModal(account)}
                          className="inline-flex items-center justify-center p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                          title="Edit Akun"
                        >
                          <Edit className="h-5 w-5" />
                        </button>
                        <button 
                          onClick={() => handleDelete(account.user_id, account.name)}
                          className="inline-flex items-center justify-center p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                          title="Hapus Akun"
                        >
                          <Trash2 className="h-5 w-5" />
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

      {/* Modal Add/Update */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-primary-deep/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-surface rounded-2xl shadow-xl w-full max-w-md overflow-hidden max-h-[90vh] flex flex-col">
            <div className="px-6 py-4 border-b border-border flex items-center justify-between bg-bg/50">
              <h3 className="font-bold text-lg text-text">
                {modalMode === "ADD" ? "Tambah Akun Puskesmas" : "Edit Akun Puskesmas"}
              </h3>
              <button 
                onClick={() => setIsModalOpen(false)}
                className="text-text-muted hover:text-text-muted"
              >
                ✕
              </button>
            </div>
            
            <form id="puskesmas-form" onSubmit={handleSubmit} className="p-6 space-y-4 overflow-y-auto flex-1">
              <div>
                <label className="block text-sm font-medium text-text mb-1">Nama Puskesmas / Institusi</label>
                <input 
                  type="text" 
                  required
                  placeholder="Mis: Puskesmas Sukmajaya"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-text mb-1">Alamat Email (Untuk Login)</label>
                <input 
                  type="email" 
                  required
                  placeholder="admin@sukmajaya.depok.go.id"
                  value={formData.email}
                  onChange={(e) => setFormData({...formData, email: e.target.value})}
                  className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-text mb-1">Nomor Kontak (WhatsApp / Telepon)</label>
                <input 
                  type="text" 
                  placeholder="Mis: 08123456789"
                  value={formData.contact}
                  onChange={(e) => setFormData({...formData, contact: e.target.value})}
                  className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                />
              </div>

              {modalMode === "ADD" && (
                <div>
                  <label className="block text-sm font-medium text-text mb-1">Password Baru</label>
                  <input 
                    type="password" 
                    required
                    placeholder="Minimal 6 karakter"
                    minLength={6}
                    value={formData.password}
                    onChange={(e) => setFormData({...formData, password: e.target.value})}
                    className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                  />
                </div>
              )}

              {modalMode === "UPDATE" && (
                <div>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input 
                      type="checkbox"
                      checked={formData.active_status}
                      onChange={(e) => setFormData({...formData, active_status: e.target.checked})}
                      className="w-4 h-4 text-indigo-600 rounded focus:ring-indigo-500"
                    />
                    <span className="text-sm font-medium text-text">Akun Aktif (Dapat Login)</span>
                  </label>
                </div>
              )}

            </form>

            <div className="p-4 border-t border-border bg-bg/50 flex gap-3 mt-auto">
              <button 
                type="button"
                onClick={() => setIsModalOpen(false)}
                disabled={isSubmitting}
                className="flex-1 bg-surface border border-border text-text py-2.5 rounded-xl font-medium hover:bg-bg transition-colors"
              >
                Batal
              </button>
              <button 
                type="submit"
                form="puskesmas-form"
                disabled={isSubmitting}
                className="flex-1 bg-indigo-600 text-white py-2.5 rounded-xl font-medium hover:bg-indigo-700 transition-colors flex justify-center"
              >
                {isSubmitting ? <Loader2 className="w-5 h-5 animate-spin" /> : (modalMode === "ADD" ? "Buat Akun" : "Simpan Perubahan")}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
