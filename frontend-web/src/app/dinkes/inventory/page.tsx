"use client";

import { useState, useEffect } from "react";
import { Plus, Package, Edit, Trash2, Loader2, AlertCircle, RefreshCw, ArrowRightLeft, CheckCircle2, XCircle } from "lucide-react";
import api from "@/lib/api";
import EmptyState from "@/components/common/EmptyState";
import SkeletonTable from "@/components/common/SkeletonTable";
import ErrorMessage from "@/components/common/ErrorMessage";

interface InventoryItem {
  item_id: string;
  item_type: string;
  quantity: number;
  unit: string;
  minimum_stock: number;
  warehouse: string;
}

interface StockRequest {
  id: string;
  item_id: string;
  warehouse: string;
  quantity: number;
  status: string;
  notes: string;
  created_at: string;
}

export default function InventoryPage() {
  const [activeTab, setActiveTab] = useState<"STOK" | "REQUESTS" | "TRANSFERS">("STOK");
  
  const [items, setItems] = useState<InventoryItem[]>([]);
  const [requests, setRequests] = useState<StockRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  
  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<"ADD" | "UPDATE">("ADD");
  const [selectedItem, setSelectedItem] = useState<InventoryItem | null>(null);
  
  // Form State
  const [formData, setFormData] = useState({
    item_type: "ABATE",
    quantity: 0,
    unit: "Kg",
    minimum_stock: 10,
    warehouse: "Gudang Utama",
    reason: "Restock"
  });

  // Transfer Form State
  const [transferData, setTransferData] = useState({
    item_id: "",
    from_warehouse: "",
    to_warehouse: "",
    quantity: 0,
    notes: ""
  });

  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [itemsRes, reqRes] = await Promise.all([
        api.get("/inventory"),
        api.get("/inventory/requests")
      ]);
      setItems(itemsRes.data?.data || []);
      setRequests(reqRes.data?.data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || "Gagal mengambil data");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const openAddModal = () => {
    setModalMode("ADD");
    setFormData({
      item_type: "ABATE",
      quantity: 0,
      unit: "Kg",
      minimum_stock: 10,
      warehouse: "Gudang Utama",
      reason: ""
    });
    setSelectedItem(null);
    setIsModalOpen(true);
  };

  const openUpdateModal = (item: InventoryItem) => {
    setModalMode("UPDATE");
    setFormData({
      item_type: item.item_type,
      quantity: 0,
      unit: item.unit,
      minimum_stock: item.minimum_stock,
      warehouse: item.warehouse,
      reason: "Restock"
    });
    setSelectedItem(item);
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      if (modalMode === "ADD") {
        await api.post("/inventory", {
          item_type: formData.item_type,
          quantity: Number(formData.quantity),
          unit: formData.unit,
          minimum_stock: Number(formData.minimum_stock),
          warehouse: formData.warehouse
        });
        alert("Item berhasil ditambahkan.");
      } else if (modalMode === "UPDATE" && selectedItem) {
        await api.patch(`/inventory/${selectedItem.item_id}/stock`, {
          quantity: Number(formData.quantity),
          reason: formData.reason
        });
        alert("Stok berhasil diupdate.");
      }
      setIsModalOpen(false);
      fetchData();
    } catch (err: any) {
      alert(err.response?.data?.error || err.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleApproveRequest = async (id: string, status: 'APPROVED' | 'REJECTED') => {
    if (!confirm(`Anda yakin ingin ${status === 'APPROVED' ? 'menyetujui' : 'menolak'} permintaan ini?`)) return;
    try {
      await api.post(`/inventory/requests/${id}/approve`, { status });
      alert(`Permintaan berhasil di${status === 'APPROVED' ? 'setujui' : 'tolak'}`);
      fetchData();
    } catch (err: any) {
      alert(err.response?.data?.error || err.message);
    }
  };

  const handleTransfer = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      await api.post("/inventory/transfers", {
        item_id: transferData.item_id,
        from_warehouse: transferData.from_warehouse,
        to_warehouse: transferData.to_warehouse,
        quantity: Number(transferData.quantity),
        notes: transferData.notes
      });
      alert("Alokasi silang berhasil dilakukan");
      setTransferData({ item_id: "", from_warehouse: "", to_warehouse: "", quantity: 0, notes: "" });
      setActiveTab("STOK");
      fetchData();
    } catch (err: any) {
      alert(err.response?.data?.error || err.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-text">Manajemen Inventaris</h2>
          <p className="text-text-muted">Kelola stok puskesmas, persetujuan permintaan, dan alokasi silang.</p>
        </div>
        {activeTab === "STOK" && (
          <button 
            onClick={openAddModal}
            className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 transition-colors font-medium shadow-sm"
          >
            <Plus className="h-4 w-4" /> Tambah Item Baru
          </button>
        )}
      </div>

      <div className="flex border-b border-border">
        <button
          className={`px-4 py-3 font-medium text-sm border-b-2 ${activeTab === "STOK" ? "border-indigo-600 text-indigo-600" : "border-transparent text-text-muted hover:text-text"}`}
          onClick={() => setActiveTab("STOK")}
        >
          <Package className="h-4 w-4 inline mr-2" /> Stok Tersedia
        </button>
        <button
          className={`px-4 py-3 font-medium text-sm border-b-2 ${activeTab === "REQUESTS" ? "border-indigo-600 text-indigo-600" : "border-transparent text-text-muted hover:text-text"}`}
          onClick={() => setActiveTab("REQUESTS")}
        >
          <RefreshCw className="h-4 w-4 inline mr-2" /> Permintaan Stok
          {requests.filter(r => r.status === 'PENDING').length > 0 && (
            <span className="ml-2 bg-red-500 text-white text-xs px-2 py-0.5 rounded-full">
              {requests.filter(r => r.status === 'PENDING').length}
            </span>
          )}
        </button>
        <button
          className={`px-4 py-3 font-medium text-sm border-b-2 ${activeTab === "TRANSFERS" ? "border-indigo-600 text-indigo-600" : "border-transparent text-text-muted hover:text-text"}`}
          onClick={() => setActiveTab("TRANSFERS")}
        >
          <ArrowRightLeft className="h-4 w-4 inline mr-2" /> Alokasi Silang
        </button>
      </div>

      {error && (
        <ErrorMessage message={error} onRetry={fetchData} />
      )}

      {activeTab === "STOK" && (
        <>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
            <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm flex items-center gap-4">
              <div className="p-3 bg-status-verified/20 text-status-verified rounded-xl">
                <Package className="h-6 w-6" />
              </div>
              <div>
                <p className="text-sm font-medium text-text-muted">Total Jenis Item</p>
                <p className="text-2xl font-bold text-text">{items.length}</p>
              </div>
            </div>
            <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm flex items-center gap-4">
              <div className="p-3 bg-red-50 text-red-600 rounded-xl">
                <AlertCircle className="h-6 w-6" />
              </div>
              <div>
                <p className="text-sm font-medium text-text-muted">Stok Menipis (Kritis)</p>
                <p className="text-2xl font-bold text-text">
                  {items.filter(i => i.quantity <= i.minimum_stock).length}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-surface border border-border rounded-2xl shadow-sm overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-bg border-b border-border text-text-muted text-sm">
                    <th className="p-4 font-semibold">Jenis Item</th>
                    <th className="p-4 font-semibold">Kuantitas</th>
                    <th className="p-4 font-semibold">Min. Stok</th>
                    <th className="p-4 font-semibold">Lokasi Gudang</th>
                    <th className="p-4 font-semibold">Status</th>
                    <th className="p-4 font-semibold text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {loading ? (
                    <tr>
                      <td colSpan={6} className="p-0">
                        <SkeletonTable rows={5} columns={6} />
                      </td>
                    </tr>
                  ) : items.length === 0 ? (
                    <tr>
                      <td colSpan={6} className="p-8">
                        <EmptyState 
                          title="Inventaris Kosong" 
                          description="Belum ada item inventaris yang tersedia."
                          icon={<Package className="w-12 h-12" />}
                        />
                      </td>
                    </tr>
                  ) : (
                    items.map((item) => {
                      const isCritical = item.quantity <= item.minimum_stock;
                      return (
                        <tr key={item.item_id} className="hover:bg-bg transition-colors">
                          <td className="p-4 font-medium text-text">{item.item_type}</td>
                          <td className="p-4 font-mono">
                            <span className={`font-bold ${isCritical ? 'text-red-600' : 'text-text'}`}>
                              {item.quantity}
                            </span> {item.unit}
                          </td>
                          <td className="p-4 text-text-muted font-mono">{item.minimum_stock} {item.unit}</td>
                          <td className="p-4 text-text-muted">{item.warehouse}</td>
                          <td className="p-4">
                            {isCritical ? (
                              <span className="px-2 py-1 bg-red-100 text-red-700 text-xs font-bold rounded-full">Kritis</span>
                            ) : (
                              <span className="px-2 py-1 bg-primary-tint text-primary text-xs font-bold rounded-full">Aman</span>
                            )}
                          </td>
                          <td className="p-4 text-center">
                            <button 
                              onClick={() => openUpdateModal(item)}
                              className="inline-flex items-center justify-center p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                              title="Update Stok"
                            >
                              <Edit className="h-5 w-5" />
                            </button>
                          </td>
                        </tr>
                      );
                    })
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </>
      )}

      {activeTab === "REQUESTS" && (
        <div className="bg-surface border border-border rounded-2xl shadow-sm overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-bg border-b border-border text-text-muted text-sm">
                  <th className="p-4 font-semibold">Tanggal</th>
                  <th className="p-4 font-semibold">Puskesmas</th>
                  <th className="p-4 font-semibold">Item</th>
                  <th className="p-4 font-semibold">Kuantitas</th>
                  <th className="p-4 font-semibold">Catatan</th>
                  <th className="p-4 font-semibold">Status</th>
                  <th className="p-4 font-semibold text-center">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {loading ? (
                  <tr>
                    <td colSpan={7} className="p-0">
                      <SkeletonTable rows={5} columns={7} />
                    </td>
                  </tr>
                ) : requests.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="p-8">
                      <EmptyState 
                        title="Tidak ada permintaan" 
                        description="Belum ada permintaan stok dari Puskesmas."
                        icon={<RefreshCw className="w-12 h-12" />}
                      />
                    </td>
                  </tr>
                ) : (
                  requests.map((req) => {
                    const item = items.find(i => i.item_id === req.item_id);
                    return (
                      <tr key={req.id} className="hover:bg-bg transition-colors">
                        <td className="p-4 text-text-muted text-sm">{new Date(req.created_at).toLocaleDateString()}</td>
                        <td className="p-4 font-medium text-text">{req.warehouse}</td>
                        <td className="p-4 text-text">{item ? item.item_type : "Item Tidak Diketahui"}</td>
                        <td className="p-4 font-mono font-bold text-text">{req.quantity} {item?.unit}</td>
                        <td className="p-4 text-text-muted italic text-sm">{req.notes || "-"}</td>
                        <td className="p-4">
                          {req.status === 'PENDING' && <span className="px-2 py-1 bg-status-wait/20 text-amber-700 text-xs font-bold rounded-full">Pending</span>}
                          {req.status === 'APPROVED' && <span className="px-2 py-1 bg-primary-tint text-primary text-xs font-bold rounded-full">Disetujui</span>}
                          {req.status === 'REJECTED' && <span className="px-2 py-1 bg-red-100 text-red-700 text-xs font-bold rounded-full">Ditolak</span>}
                        </td>
                        <td className="p-4 text-center">
                          {req.status === 'PENDING' ? (
                            <div className="flex justify-center gap-2">
                              <button onClick={() => handleApproveRequest(req.id, 'APPROVED')} className="p-1 text-primary hover:bg-primary-tint rounded" title="Setujui">
                                <CheckCircle2 className="w-5 h-5" />
                              </button>
                              <button onClick={() => handleApproveRequest(req.id, 'REJECTED')} className="p-1 text-red-600 hover:bg-red-50 rounded" title="Tolak">
                                <XCircle className="w-5 h-5" />
                              </button>
                            </div>
                          ) : (
                            <span className="text-slate-400 text-sm">-</span>
                          )}
                        </td>
                      </tr>
                    );
                  })
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {activeTab === "TRANSFERS" && (
        <div className="bg-surface border border-border rounded-2xl shadow-sm p-6 max-w-2xl">
          <h3 className="font-bold text-lg text-text mb-4">Form Alokasi Silang (Cross-Allocation)</h3>
          <form onSubmit={handleTransfer} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-text mb-1">Item yang Ditransfer</label>
              <select 
                required
                value={transferData.item_id}
                onChange={(e) => {
                  const selectedItem = items.find(i => i.item_id === e.target.value);
                  setTransferData({
                    ...transferData, 
                    item_id: e.target.value, 
                    from_warehouse: selectedItem ? selectedItem.warehouse : ""
                  });
                }}
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
              >
                <option value="">-- Pilih Item dari Gudang Asal --</option>
                {items.filter(i => i.quantity > 0).map(item => (
                  <option key={item.item_id} value={item.item_id}>
                    {item.warehouse} - {item.item_type} (Sisa: {item.quantity} {item.unit})
                  </option>
                ))}
              </select>
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-text mb-1">Gudang Asal (Otomatis)</label>
                <input 
                  type="text" 
                  readOnly
                  value={transferData.from_warehouse}
                  className="w-full px-3 py-2 border border-border bg-bg text-text-muted rounded-lg"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-text mb-1">Puskesmas Tujuan</label>
                <input 
                  type="text" 
                  required
                  value={transferData.to_warehouse}
                  onChange={(e) => setTransferData({...transferData, to_warehouse: e.target.value})}
                  placeholder="Nama Puskesmas Tujuan"
                  className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-text mb-1">Kuantitas Transfer</label>
              <input 
                type="number" 
                required
                min="0.1"
                step="0.1"
                value={transferData.quantity || ""}
                onChange={(e) => setTransferData({...transferData, quantity: parseFloat(e.target.value)})}
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-text mb-1">Catatan / Alasan (Opsional)</label>
              <textarea 
                rows={2}
                value={transferData.notes}
                onChange={(e) => setTransferData({...transferData, notes: e.target.value})}
                placeholder="Alasan transfer silang..."
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
              ></textarea>
            </div>

            <button 
              type="submit"
              disabled={isSubmitting || !transferData.item_id || !transferData.to_warehouse || !transferData.quantity}
              className="w-full bg-indigo-600 text-white py-2.5 rounded-xl font-medium hover:bg-indigo-700 transition-colors disabled:bg-indigo-300"
            >
              {isSubmitting ? "Memproses..." : "Eksekusi Alokasi Silang"}
            </button>
          </form>
        </div>
      )}

      {/* Modal Add/Update untuk Tab Stok */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-primary-deep/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-surface rounded-2xl shadow-xl w-full max-w-md overflow-hidden max-h-[90vh] flex flex-col">
            <div className="px-6 py-4 border-b border-border flex items-center justify-between bg-bg/50">
              <h3 className="font-bold text-lg text-text">
                {modalMode === "ADD" ? "Tambah Item Baru" : `Update Stok ${selectedItem?.item_type}`}
              </h3>
              <button 
                onClick={() => setIsModalOpen(false)}
                className="text-slate-400 hover:text-text-muted"
              >
                ✕
              </button>
            </div>
            
            <form id="inventory-form" onSubmit={handleSubmit} className="p-6 space-y-4 overflow-y-auto flex-1">
              {modalMode === "ADD" && (
                <>
                  <div>
                    <label className="block text-sm font-medium text-text mb-1">Jenis Item (Mis: ABATE, INSEKTISIDA)</label>
                    <input 
                      type="text" 
                      required
                      value={formData.item_type}
                      onChange={(e) => setFormData({...formData, item_type: e.target.value.toUpperCase()})}
                      className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text mb-1">Gudang Penyimpanan (Nama Puskesmas)</label>
                    <input 
                      type="text" 
                      required
                      value={formData.warehouse}
                      onChange={(e) => setFormData({...formData, warehouse: e.target.value})}
                      className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-text mb-1">Stok Awal</label>
                      <input 
                        type="number" 
                        required
                        min="0"
                        step="0.1"
                        value={formData.quantity}
                        onChange={(e) => setFormData({...formData, quantity: parseFloat(e.target.value)})}
                        className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-text mb-1">Batas Minimum</label>
                      <input 
                        type="number" 
                        required
                        min="0"
                        step="0.1"
                        value={formData.minimum_stock}
                        onChange={(e) => setFormData({...formData, minimum_stock: parseFloat(e.target.value)})}
                        className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                      />
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text mb-1">Satuan (Kg, Liter)</label>
                    <input 
                      type="text" 
                      required
                      value={formData.unit}
                      onChange={(e) => setFormData({...formData, unit: e.target.value})}
                      className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>
                </>
              )}

              {modalMode === "UPDATE" && (
                <>
                  <div className="bg-bg p-4 rounded-xl border border-border mb-4">
                    <p className="text-sm text-text-muted">Stok Saat Ini di {selectedItem?.warehouse}</p>
                    <p className="text-lg font-bold text-text">{selectedItem?.quantity} {selectedItem?.unit}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text mb-1">Kuantitas (+/-)</label>
                    <div className="flex items-center gap-2">
                      <input 
                        type="number" 
                        required
                        step="0.1"
                        value={formData.quantity}
                        onChange={(e) => setFormData({...formData, quantity: parseFloat(e.target.value)})}
                        className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500 font-mono text-lg"
                      />
                      <span className="text-text-muted font-medium">{selectedItem?.unit}</span>
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text mb-1">Alasan Penyesuaian</label>
                    <input 
                      type="text" 
                      required
                      value={formData.reason}
                      onChange={(e) => setFormData({...formData, reason: e.target.value})}
                      className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>
                </>
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
                form="inventory-form"
                disabled={isSubmitting}
                className="flex-1 bg-indigo-600 text-white py-2.5 rounded-xl font-medium hover:bg-indigo-700 transition-colors flex justify-center"
              >
                {isSubmitting ? <Loader2 className="w-5 h-5 animate-spin" /> : (modalMode === "ADD" ? "Simpan Item" : "Update Stok")}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
