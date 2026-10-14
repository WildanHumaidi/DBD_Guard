"use client";

import React, { useEffect, useState } from "react";
import api from "@/lib/api";
import { ExclamationCircleIcon, ShieldExclamationIcon } from "@heroicons/react/24/outline";

interface EwsReportItem {
  report_id: string;
  reference_no: string;
  created_at: string;
  status: string;
  category?: string;
  description?: string;
  cases_count?: number;
  location?: {
    rt?: string;
    rw?: string;
    kelurahan?: string;
    kecamatan?: string;
  };
}

export default function EwsDashboard() {
  const [warnings, setWarnings] = useState<EwsReportItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let ignore = false;
    api.get("/reports")
      .then((res) => {
        if (!ignore) {
          const allReports = (res.data?.data || []) as EwsReportItem[];
          const sevenDaysAgo = new Date();
          sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);
          
          const recentCases = allReports.filter((r) => {
            const isRecent = new Date(r.created_at) > sevenDaysAgo;
            const isValid = ['MENUNGGU_VALIDASI_PUSKESMAS', 'WAITING_FOGGING_APPROVAL', 'FOGGING_SCHEDULED', 'FOGGING_IN_PROGRESS', 'COMPLETED'].includes(r.status);
            return isRecent && isValid;
          });
          
          setWarnings(recentCases);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          setError((err as Error)?.message || "Gagal memuat Peringatan Dini (EWS)");
        }
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  if (loading) return <div className="p-8">Memuat EWS...</div>;
  if (error) return <div className="p-8 text-red-500">{error}</div>;

  return (
    <div className="p-8 max-w-5xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-text">Early Warning System (EWS)</h1>
        <p className="text-text-muted mt-1">Sistem Peringatan Dini Potensi Kasus DBD di Wilayah Anda</p>
      </div>

      <div className="bg-surface rounded-2xl shadow-sm border border-border overflow-hidden">
        <div className="p-6 border-b border-border flex items-center gap-3">
          <div className="bg-status-wait/20 p-2 rounded-lg text-amber-600">
            <ShieldExclamationIcon className="w-5 h-5" />
          </div>
          <h2 className="text-lg font-bold text-text">Notifikasi Kewaspadaan</h2>
        </div>
        
        <div className="p-6">
          {warnings.length === 0 ? (
            <div className="text-center py-12">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-tint text-primary mb-4">
                <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <p className="text-text-muted font-medium">Kondisi Aman</p>
              <p className="text-sm text-text-muted mt-1">Tidak ada deteksi anomali penyebaran di wilayah ini.</p>
            </div>
          ) : (
            <div className="space-y-4">
              {warnings.map((warn, index) => (
                <div key={index} className="bg-amber-50 border-l-4 border-amber-400 p-4 rounded-r-lg flex items-start gap-4">
                  <ExclamationCircleIcon className="w-6 h-6 text-status-wait shrink-0 mt-0.5" />
                  <div>
                    <h3 className="font-bold text-amber-900">Perhatian: {warn.category || "Peningkatan Kasus"}</h3>
                    <p className="text-sm text-status-wait mt-1">
                      Terdeteksi {warn.cases_count || 1} laporan kasus (No Ref: {warn.reference_no}) 
                      di RT {warn.location?.rt} / RW {warn.location?.rw}, {warn.location?.kelurahan}. 
                      Segera lakukan edukasi dan pengecekan jentik.
                    </p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
