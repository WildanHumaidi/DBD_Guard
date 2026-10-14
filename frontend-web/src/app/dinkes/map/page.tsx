"use client";

import { useEffect, useState } from "react";
import dynamic from "next/dynamic";
import { MapPin, AlertTriangle, RefreshCw, Layers } from "lucide-react";
import api from "@/lib/api";
import type { GeoJsonObject } from "geojson";

// Dynamic import to avoid SSR issues with Leaflet
const MapView = dynamic(() => import("@/components/GisMapView"), {
  ssr: false,
  loading: () => (
    <div className="w-full h-full flex items-center justify-center bg-bg rounded-xl">
      <div className="flex flex-col items-center gap-3">
        <div className="h-10 w-10 border-4 border-primary border-t-emerald-600 rounded-full animate-spin"></div>
        <span className="text-sm text-text-muted font-medium">
          Memuat peta...
        </span>
      </div>
    </div>
  ),
});

interface RiskArea {
  kelurahan: string;
  kecamatan: string;
  case_count: number;
  risk_level: string;
}

export default function AdminMapPage() {
  const [geoData, setGeoData] = useState<GeoJsonObject | null>(null);
  const [riskAreas, setRiskAreas] = useState<RiskArea[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchData = async () => {
    setLoading(true);
    setError("");
    try {
      const [heatmapRes, riskRes] = await Promise.all([
        api.get("/admin/gis/heatmap"),
        api.get("/admin/gis/risk-score"),
      ]);

      if (heatmapRes.data?.data) {
        setGeoData(heatmapRes.data.data);
      }
      if (riskRes.data?.data) {
        setRiskAreas(riskRes.data.data);
      }
    } catch {
      setError("Gagal memuat data peta. Pastikan koneksi backend aktif.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ignore = false;
    Promise.all([
      api.get("/admin/gis/heatmap"),
      api.get("/admin/gis/risk-score"),
    ])
      .then(([heatmapRes, riskRes]) => {
        if (!ignore) {
          if (heatmapRes.data?.data) {
            setGeoData(heatmapRes.data.data);
          }
          if (riskRes.data?.data) {
            setRiskAreas(riskRes.data.data);
          }
        }
      })
      .catch(() => {
        if (!ignore) {
          setError("Gagal memuat data peta. Pastikan koneksi backend aktif.");
        }
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });

    return () => {
      ignore = true;
    };
  }, []);

  const getRiskColor = (level: string) => {
    switch (level) {
      case "Bahaya":
        return "bg-red-100 text-red-700 border-red-200";
      case "Waspada":
        return "bg-status-wait/20 text-amber-700 border-status-wait/30";
      default:
        return "bg-primary-tint text-primary border-primary";
    }
  };

  const getRiskDot = (level: string) => {
    switch (level) {
      case "Bahaya":
        return "bg-red-500";
      case "Waspada":
        return "bg-amber-500";
      default:
        return "bg-primary";
    }
  };

  return (
    <div className="h-full flex flex-col gap-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text flex items-center gap-3">
            <div className="h-10 w-10 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/25">
              <MapPin className="h-5 w-5 text-white" />
            </div>
            Peta Sebaran Kasus (GIS)
          </h1>
          <p className="text-sm text-text-muted mt-1 ml-[52px]">
            Pemetaan interaktif sebaran laporan DBD menggunakan Leaflet +
            OpenStreetMap.
          </p>
        </div>
        <button
          onClick={fetchData}
          disabled={loading}
          className="inline-flex items-center gap-2 px-5 py-2.5 bg-surface border border-border rounded-xl text-sm font-bold text-text hover:bg-bg hover:border-border transition-all disabled:opacity-50 shadow-sm"
        >
          <RefreshCw
            className={`h-4 w-4 ${loading ? "animate-spin" : ""}`}
          />
          Refresh Data
        </button>
      </div>

      {/* Main Content */}
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-4 gap-6 min-h-0">
        {/* Map Area */}
        <div className="lg:col-span-3 bg-surface rounded-2xl shadow-lg border border-border overflow-hidden min-h-[500px]">
          {error ? (
            <div className="w-full h-full flex flex-col items-center justify-center p-10 text-center">
              <AlertTriangle className="h-12 w-12 text-amber-400 mb-4" />
              <p className="text-text-muted font-medium mb-4">{error}</p>
              <button
                onClick={fetchData}
                className="px-6 py-2.5 bg-primary text-primary-foreground rounded-xl text-sm font-bold hover:bg-primary-deep transition-all"
              >
                Coba Lagi
              </button>
            </div>
          ) : (
            <MapView geoData={geoData} loading={loading} />
          )}
        </div>

        {/* Risk Score Sidebar */}
        <div className="lg:col-span-1 bg-surface rounded-2xl shadow-lg border border-border p-5 flex flex-col min-h-0">
          <div className="flex items-center gap-3 mb-5 pb-4 border-b border-border">
            <div className="h-8 w-8 bg-bg rounded-lg flex items-center justify-center">
              <Layers className="h-4 w-4 text-text-muted" />
            </div>
            <h2 className="font-bold text-text text-sm">
              Skor Risiko per Kelurahan
            </h2>
          </div>

          {loading ? (
            <div className="flex-1 space-y-3">
              {[...Array(5)].map((_, i) => (
                <div
                  key={i}
                  className="h-16 bg-bg rounded-xl animate-pulse"
                ></div>
              ))}
            </div>
          ) : riskAreas.length === 0 ? (
            <div className="flex-1 flex items-center justify-center">
              <p className="text-sm text-text-muted text-center">
                Belum ada data risiko area.
              </p>
            </div>
          ) : (
            <div className="flex-1 overflow-y-auto space-y-2 pr-1">
              {riskAreas
                .sort((a, b) => b.case_count - a.case_count)
                .map((area, index) => (
                  <div
                    key={index}
                    className="p-3 rounded-xl border border-border hover:bg-bg transition-colors"
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-sm font-bold text-text">
                        {area.kelurahan}
                      </span>
                      <span
                        className={`text-xs font-bold px-2 py-0.5 rounded-full border ${getRiskColor(
                          area.risk_level
                        )}`}
                      >
                        {area.risk_level}
                      </span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-text-muted">
                        Kec. {area.kecamatan}
                      </span>
                      <div className="flex items-center gap-1.5">
                        <div
                          className={`h-2 w-2 rounded-full ${getRiskDot(
                            area.risk_level
                          )}`}
                        ></div>
                        <span className="text-xs font-bold text-text-muted">
                          {area.case_count} kasus
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
            </div>
          )}

          {/* Legend */}
          <div className="mt-4 pt-4 border-t border-border">
            <p className="text-xs font-bold text-text-muted mb-2">Legenda</p>
            <div className="space-y-1.5">
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-primary"></div>
                <span className="text-xs text-text-muted">
                  Aman (0 kasus)
                </span>
              </div>
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-amber-500"></div>
                <span className="text-xs text-text-muted">
                  Waspada (1–3 kasus)
                </span>
              </div>
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-red-500"></div>
                <span className="text-xs text-text-muted">
                  Bahaya (&gt;3 kasus)
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
