"use client";

import { useState, useEffect } from "react";
import { Package, Send, Loader2, AlertCircle, RefreshCw, CheckCircle2, Clock, XCircle } from "lucide-react";
import api from "@/lib/api";

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

interface StockTransfer {
  id: string;
  item_id: string;
  from_warehouse: string;
  to_warehouse: string;
  quantity: number;
  status: string;
  notes: string;
  created_at: string;
}

export default function PuskesmasInventoryPage() {
  const [activeTab, setActiveTab] = useState<"STOK" | "REQUEST" | "TRANSFER" | "FOGGING">("STOK");
  const [myWarehouse, setMyWarehouse] = useState("Puskesmas Sukmajaya");
  
  const [items, setItems] = useState<InventoryItem[]>([]);
  const [allItems, setAllItems] = useState<InventoryItem[]>([]);
  const [requests, setRequests] = useState<StockRequest[]>([]);
  const [transfers, setTransfers] = useState<StockTransfer[]>([]);
  const [foggingUsages, setFoggingUsages] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  
  // Request Form State
  const [requestData, setRequestData] = useState({
    item_id: "",
    quantity: 0,
    notes: ""
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [itemsRes, reqRes, transferRes] = await Promise.all([
        api.get("/inventory"),
        api.get(`/inventory/requests?warehouse=${encodeURIComponent(myWarehouse)}`),
        api.get("/inventory/transfers")
      ]);
      
      const allInventory = itemsRes.data?.data || [];
      setAllItems(allInventory);
      
      // Filter for this warehouse
      setItems(allInventory.filter((i: InventoryItem) => i.warehouse === myWarehouse));
      
      // We can also filter requests locally if the backend doesn't filter perfectly yet
      const allReqs = reqRes.data?.data || [];
      setRequests(allReqs.filter((r: StockRequest) => r.warehouse === myWarehouse));

      const allTransfers = transferRes.data?.data || [];
      setTransfers(allTransfers.filter((t: StockTransfer) => t.to_warehouse === myWarehouse));

      // Fetch fogging usage history for Puskesmas audit
      try {
        const usageRes = await api.get('/inventory/usage/history').catch(() => api.get('/inventory/usage/pending'));
        setFoggingUsages(usageRes.data?.data || []);
      } catch (e) {
        console.warn("Could not fetch fogging usage history.", e);
      }
      
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || "Gagal mengambil data");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [myWarehouse]);

  const handleRequestSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      await api.post("/inventory/requests", {
        item_id: requestData.item_id,
        warehouse: myWarehouse,
        quantity: Number(requestData.quantity),
        notes: requestData.notes
      });
      alert("Permintaan stok berhasil diajukan.");
      setRequestData({ item_id: "", quantity: 0, notes: "" });
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
          <h2 className="text-2xl font-bold text-text">Inventaris Puskesmas</h2>
          <p className="text-text-muted">Pantau stok lokal dan ajukan permintaan ke Dinas Kesehatan.</p>
        </div>
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-text-muted">Lokasi Anda:</label>
          <select 
            value={myWarehouse}
            onChange={(e) => setMyWarehouse(e.target.value)}
            className="px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500 text-sm"
          >
            <option value="Puskesmas Sukmajaya">Puskesmas Sukmajaya</option>
            <option value="Puskesmas Beji">Puskesmas Beji</option>
            <option value="Puskesmas Pancoran Mas">Puskesmas Pancoran Mas</option>
            <option value="Gudang Utama">Gudang Utama</option>
          </select>
        </div>
      </div>

      <div className="flex border-b border-border">
        <button
          className={`px-4 py-3 font-medium text-sm border-b-2 ${activeTab === "STOK" ? "border-indigo-600 text-indigo-600" : "border-transparent text-text-muted hover:text-text"}`}
          onClick={() => setActiveTab("STOK")}
        >
          <Package className="h-4 w-4 inline mr-2" /> Stok Saya
        </button>
        <button
          className={`px-4 py-3 font-medium text-sm border-b-2 ${activeTab === "REQUEST" ? "border-indigo-600 text-indigo-600" : "border-transparent text-text-muted hover:text-text"}`}
          onClick={() => setActiveTab("REQUEST")}
        >
          <Send className="h-4 w-4 inline mr-2" /> Ajukan Permintaan
        </button>
        <button
          className={`px-4 py-3 font-medium text-sm border-b-2 ${activeTab === "TRANSFER" ? "border-indigo-600 text-indigo-600" : "border-transparent text-text-muted hover:text-text"}`}
          onClick={() => setActiveTab("TRANSFER")}
        >
          <Package className="h-4 w-4 inline mr-2" /> Terima Transfer
        </button>
        <button
          className={`px-4 py-3 font-medium text-sm border-b-2 ${activeTab === "FOGGING" ? "border-indigo-600 text-indigo-600" : "border-transparent text-text-muted hover:text-text"}`}
          onClick={() => setActiveTab("FOGGING")}
        >
          <CheckCircle2 className="h-4 w-4 inline mr-2" /> Persetujuan Fogging
        </button>
      </div>

      {error && (
        <div className="flex items-center gap-2 p-4 bg-red-50 text-red-700 rounded-xl border border-red-200">
          <AlertCircle className="h-5 w-5" />
          <p className="text-sm font-medium">{error}</p>
        </div>
      )}

      {activeTab === "STOK" && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
            <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm flex items-center gap-4">
              <div className="p-3 bg-status-verified/20 text-status-verified rounded-xl">
                <Package className="h-6 w-6" />
              </div>
              <div>
                <p className="text-sm font-medium text-text-muted">Total Item di Lokasi Ini</p>
                <p className="text-2xl font-bold text-text">{items.length}</p>
              </div>
            </div>
            <div className="bg-surface p-6 rounded-2xl border border-border shadow-sm flex items-center gap-4">
              <div className="p-3 bg-amber-50 text-amber-600 rounded-xl">
                <RefreshCw className="h-6 w-6" />
              </div>
              <div>
                <p className="text-sm font-medium text-text-muted">Permintaan Pending</p>
                <p className="text-2xl font-bold text-text">
                  {requests.filter(r => r.status === 'PENDING').length}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-surface border border-border rounded-2xl shadow-sm overflow-hidden">
            <div className="p-4 border-b border-border bg-bg">
              <h3 className="font-semibold text-text">Daftar Stok: {myWarehouse}</h3>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-surface border-b border-border text-text-muted text-sm">
                    <th className="p-4 font-semibold">Jenis Item</th>
                    <th className="p-4 font-semibold">Kuantitas</th>
                    <th className="p-4 font-semibold">Min. Stok</th>
                    <th className="p-4 font-semibold">Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {loading ? (
                    <tr>
                      <td colSpan={4} className="p-8 text-center text-text-muted">
                        <Loader2 className="w-6 h-6 animate-spin mx-auto mb-2 text-indigo-500" />
                        Memuat data...
                      </td>
                    </tr>
                  ) : items.length === 0 ? (
                    <tr>
                      <td colSpan={4} className="p-8 text-center text-text-muted">Belum ada item inventaris di lokasi ini.</td>
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
                          <td className="p-4">
                            {isCritical ? (
                              <span className="px-2 py-1 bg-red-100 text-red-700 text-xs font-bold rounded-full">Kritis</span>
                            ) : (
                              <span className="px-2 py-1 bg-primary-tint text-primary text-xs font-bold rounded-full">Aman</span>
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

          <div className="bg-surface border border-border rounded-2xl shadow-sm overflow-hidden mt-6">
            <div className="p-4 border-b border-border bg-bg">
              <h3 className="font-semibold text-text">Riwayat Permintaan Stok</h3>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-surface border-b border-border text-text-muted text-sm">
                    <th className="p-4 font-semibold">Tanggal</th>
                    <th className="p-4 font-semibold">Item</th>
                    <th className="p-4 font-semibold">Kuantitas Diminta</th>
                    <th className="p-4 font-semibold">Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {requests.length === 0 ? (
                    <tr>
                      <td colSpan={4} className="p-8 text-center text-text-muted">Tidak ada riwayat permintaan.</td>
                    </tr>
                  ) : (
                    requests.map((req) => {
                      const item = allItems.find(i => i.item_id === req.item_id);
                      return (
                        <tr key={req.id} className="hover:bg-bg transition-colors">
                          <td className="p-4 text-text-muted text-sm">{new Date(req.created_at).toLocaleDateString()}</td>
                          <td className="p-4 font-medium text-text">{item ? item.item_type : "Unknown"}</td>
                          <td className="p-4 font-mono font-bold text-text">{req.quantity} {item?.unit}</td>
                          <td className="p-4">
                            {req.status === 'PENDING' && (
                              <span className="inline-flex items-center gap-1 px-2 py-1 bg-status-wait/20 text-amber-700 text-xs font-bold rounded-full">
                                <Clock className="w-3 h-3" /> Pending
                              </span>
                            )}
                            {req.status === 'APPROVED' && (
                              <span className="inline-flex items-center gap-1 px-2 py-1 bg-primary-tint text-primary text-xs font-bold rounded-full">
                                <CheckCircle2 className="w-3 h-3" /> Disetujui
                              </span>
                            )}
                            {req.status === 'REJECTED' && (
                              <span className="inline-flex items-center gap-1 px-2 py-1 bg-red-100 text-red-700 text-xs font-bold rounded-full">
                                <XCircle className="w-3 h-3" /> Ditolak
                              </span>
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
        </div>
      )}

      {activeTab === "REQUEST" && (
        <div className="bg-surface border border-border rounded-2xl shadow-sm p-6 max-w-2xl">
          <h3 className="font-bold text-lg text-text mb-4">Form Pengajuan Stok ke Dinas</h3>
          <form onSubmit={handleRequestSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-text mb-1">Pilih Item dari Gudang Utama</label>
              <select 
                required
                value={requestData.item_id}
                onChange={(e) => setRequestData({...requestData, item_id: e.target.value})}
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
              >
                <option value="">-- Pilih Item --</option>
                {/* We map from Gudang Utama items */}
                {allItems
                  .filter(i => i.warehouse === "Gudang Utama")
                  .map(item => (
                    <option key={item.item_id} value={item.item_id}>
                      {item.item_type} (Tersedia: {item.quantity} {item.unit})
                    </option>
                ))}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-text mb-1">Kuantitas yang Diminta</label>
              <input 
                type="number" 
                required
                min="1"
                step="0.1"
                value={requestData.quantity || ""}
                onChange={(e) => setRequestData({...requestData, quantity: parseFloat(e.target.value)})}
                placeholder="Contoh: 5"
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-text mb-1">Catatan / Alasan (Opsional)</label>
              <textarea 
                rows={3}
                value={requestData.notes}
                onChange={(e) => setRequestData({...requestData, notes: e.target.value})}
                placeholder="Alasan mengapa butuh tambahan stok..."
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-indigo-500"
              ></textarea>
            </div>

            <button 
              type="submit"
              disabled={isSubmitting || !requestData.item_id || !requestData.quantity}
              className="w-full bg-indigo-600 text-white py-2.5 rounded-xl font-medium hover:bg-indigo-700 transition-colors disabled:bg-indigo-300 flex items-center justify-center gap-2"
            >
              {isSubmitting ? <Loader2 className="w-5 h-5 animate-spin" /> : <Send className="w-5 h-5" />}
              {isSubmitting ? "Mengirim..." : "Kirim Permintaan Stok"}
            </button>
          </form>
        </div>
      )}

      {activeTab === "TRANSFER" && (
        <div className="bg-surface border border-border rounded-2xl shadow-sm overflow-hidden mt-6">
          <div className="p-4 border-b border-border bg-bg flex justify-between items-center">
            <h3 className="font-semibold text-text">Transfer Stok Masuk</h3>
            <span className="text-sm text-text-muted">Dari Dinkes / Gudang Utama</span>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-surface border-b border-border text-text-muted text-sm">
                  <th className="p-4 font-semibold">Tanggal</th>
                  <th className="p-4 font-semibold">Dari</th>
                  <th className="p-4 font-semibold">Item</th>
                  <th className="p-4 font-semibold">Jumlah</th>
                  <th className="p-4 font-semibold">Status / Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {transfers.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="p-8 text-center text-text-muted">Tidak ada transfer masuk.</td>
                  </tr>
                ) : (
                  transfers.map((t) => {
                    const item = allItems.find(i => i.item_id === t.item_id);
                    return (
                      <tr key={t.id} className="hover:bg-bg transition-colors">
                        <td className="p-4 text-text-muted text-sm">{new Date(t.created_at).toLocaleDateString()}</td>
                        <td className="p-4 text-text">{t.from_warehouse}</td>
                        <td className="p-4 font-medium text-text">{item ? item.item_type : t.item_id}</td>
                        <td className="p-4 font-mono font-bold text-text">{t.quantity} {item?.unit}</td>
                        <td className="p-4">
                          {t.status === 'IN_TRANSIT' || t.status === 'PENDING' ? (
                            <button
                              onClick={async () => {
                                try {
                                  await api.post(`/inventory/transfers/${t.id}/receive`, {
                                    is_accepted: true,
                                    notes: "Diterima oleh Petugas Puskesmas"
                                  });
                                  alert("Stok berhasil diterima!");
                                  fetchData();
                                } catch (err: any) {
                                  alert(err.response?.data?.error || "Gagal menerima stok");
                                }
                              }}
                              className="px-3 py-1.5 bg-indigo-600 text-white rounded-lg text-sm font-bold hover:bg-indigo-700"
                            >
                              Terima Stok
                            </button>
                          ) : t.status === 'REJECTED' ? (
                            <span className="inline-flex items-center gap-1 px-2 py-1 bg-red-100 text-red-700 text-xs font-bold rounded-full">
                              <XCircle className="w-3 h-3" /> Ditolak
                            </span>
                          ) : (
                            <span className="inline-flex items-center gap-1 px-2 py-1 bg-status-verified/20 text-status-verified text-xs font-bold rounded-full">
                              <CheckCircle2 className="w-3 h-3" /> Diterima
                            </span>
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

      {activeTab === "FOGGING" && (
        <div className="bg-surface border border-border rounded-2xl shadow-sm overflow-hidden">
          <div className="p-4 border-b border-border bg-bg flex justify-between items-center">
            <h3 className="font-semibold text-text">Catatan Pemakaian Bahan Fogging (Terekam Otomatis)</h3>
            <span className="text-xs text-text-muted bg-surface px-2.5 py-1 rounded-full border border-border">
              Single Source of Truth: Terpotong Atomik saat Fogging Selesai
            </span>
          </div>
          <div className="p-4 space-y-4">
            {foggingUsages.length === 0 ? (
              <p className="text-center text-text-muted py-8">Belum ada riwayat pemakaian stok fogging.</p>
            ) : (
              foggingUsages.map(usage => {
                const item = allItems.find(i => i.item_id === usage.item_id);
                return (
                  <div key={usage.id} className="border border-border rounded-xl p-4 shadow-sm flex flex-col md:flex-row justify-between md:items-center gap-4 hover:bg-bg/40 transition-colors">
                    <div>
                      <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold bg-status-verified/20 text-status-verified mb-2">
                        <CheckCircle2 className="w-3.5 h-3.5" />
                        {usage.status === 'RECORDED' || usage.status === 'APPROVED' ? "Stok Terpotong Otomatis" : usage.status}
                      </span>
                      <h4 className="font-bold text-text">
                        {item ? item.item_type : "Item ID: " + usage.item_id}
                      </h4>
                      <p className="text-sm text-text-muted mb-1">
                        Digunakan oleh Tim Fogging: <span className="font-semibold font-mono text-text">{usage.quantity} {item?.unit || ""}</span>
                      </p>
                      <p className="text-xs text-text-muted mb-2">{new Date(usage.created_at).toLocaleString('id-ID')}</p>
                      {usage.notes && (
                        <p className="text-sm text-text-muted bg-bg p-2 rounded-lg border border-border">
                          <span className="font-medium text-text">Catatan Tim:</span> {usage.notes}
                        </p>
                      )}
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>
      )}
    </div>
  );
}
