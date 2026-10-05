"use client";

import { useEffect, useState } from "react";
import { MapContainer, TileLayer, GeoJSON } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import L from "leaflet";

// Fix Leaflet default icon issues in Next.js
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
  iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
  shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
});

// Center of Depok city
const DEPOK_CENTER: [number, number] = [-6.402484, 106.794243];
const DEFAULT_ZOOM = 12;

interface GisMapViewProps {
  geoData: any;
  loading: boolean;
}

export default function GisMapView({ geoData, loading }: GisMapViewProps) {
  const [mapReady, setMapReady] = useState(false);

  useEffect(() => {
    setMapReady(true);
  }, []);

  if (!mapReady || loading) {
    return (
      <div className="w-full h-full flex items-center justify-center bg-bg">
        <div className="flex flex-col items-center gap-3">
          <div className="h-10 w-10 border-4 border-emerald-200 border-t-emerald-600 rounded-full animate-spin"></div>
          <span className="text-sm text-text-muted font-medium">
            Memuat peta...
          </span>
        </div>
      </div>
    );
  }

  const getMarkerColor = (status: string) => {
    switch (status) {
      case "MENUNGGU_VALIDASI_PUSKESMAS":
      case "PE_SCHEDULED":
      case "FOGGING_SCHEDULED":
        return "#ef4444"; // Red (Bahaya)
      case "PENDING":
      case "INVESTIGATING":
        return "#f59e0b"; // Yellow (Waspada)
      case "COMPLETED":
      case "REJECTED":
        return "#10b981"; // Green (Aman)
      default:
        return "#ef4444";
    }
  };

  const pointToLayer = (feature: any, latlng: L.LatLng) => {
    const status = feature.properties?.status || "UNKNOWN";
    const color = getMarkerColor(status);
    
    return L.circleMarker(latlng, {
      radius: 10,
      fillColor: color,
      color: "#ffffff",
      weight: 2,
      opacity: 1,
      fillOpacity: 0.7,
    });
  };

  const onEachFeature = (feature: any, layer: L.Layer) => {
    if (feature.properties) {
      const props = feature.properties;
      const statusMap: Record<string, string> = {
        PENDING: "Menunggu Verifikasi",
        INVESTIGATING: "Dalam Penyelidikan (PE)",
        MENUNGGU_VALIDASI_PUSKESMAS: "Kasus Positif",
        PE_SCHEDULED: "Terjadwal PE",
        FOGGING_SCHEDULED: "Terjadwal Fogging",
        COMPLETED: "Selesai",
        REJECTED: "Ditolak/Batal",
      };
      
      const statusLabel = statusMap[props.status] || props.status || "Unknown";

      layer.bindPopup(
        `<div style="font-family: Inter, system-ui, sans-serif; padding: 4px 0;">
          <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px;">
            <p style="font-weight: 700; font-size: 14px; margin: 0; color: #0f172a;">
              ${props.address || "Lokasi Kasus"}
            </p>
            <span style="font-size: 10px; font-weight: bold; background: #f1f5f9; padding: 2px 6px; border-radius: 4px; color: #475569;">
              ${statusLabel}
            </span>
          </div>
          <p style="font-size: 12px; color: #64748b; margin: 0 0 2px 0;">
            RT ${props.rt || "-"} / RW ${props.rw || "-"}
          </p>
          <p style="font-size: 12px; color: #64748b; margin: 0;">
            Kel. ${props.kelurahan || "-"}, Kec. ${props.kecamatan || "-"}
          </p>
        </div>`,
        { className: "custom-popup" }
      );
    }
  };

  return (
    <MapContainer
      center={DEPOK_CENTER}
      zoom={DEFAULT_ZOOM}
      className="w-full h-full rounded-xl z-0"
      style={{ minHeight: "500px" }}
    >
      <TileLayer
        attribution={process.env.NEXT_PUBLIC_MAP_ATTRIBUTION || '&copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors'}
        url={process.env.NEXT_PUBLIC_MAP_TILE_URL || "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"}
      />
      {geoData && geoData.features && geoData.features.length > 0 && (
        <GeoJSON
          key={JSON.stringify(geoData)}
          data={geoData}
          pointToLayer={pointToLayer}
          onEachFeature={onEachFeature}
        />
      )}
    </MapContainer>
  );
}
