"use client";

import { useEffect, useState, useMemo, useRef } from "react";
import { MapContainer, TileLayer, Marker, Popup, useMap, useMapEvents } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import L from "leaflet";
import { Crosshair, MapPin, CheckCircle2, AlertCircle } from "lucide-react";

// Custom modern pulsing pin matching DBD Guard palette (Teal #0F766E)
const userLocationIcon = L.divIcon({
  className: "dbd-guard-user-pin",
  html: `
    <div style="position: relative; width: 36px; height: 36px; display: flex; align-items: center; justify-content: center;">
      <div style="position: absolute; width: 36px; height: 36px; background: rgba(15, 118, 110, 0.25); border-radius: 50%; animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;"></div>
      <div style="position: relative; width: 30px; height: 30px; background: #0F766E; border: 3px solid #ffffff; border-radius: 50%; box-shadow: 0 4px 12px rgba(15, 118, 110, 0.4); display: flex; align-items: center; justify-content: center;">
        <div style="width: 10px; height: 10px; background: #ffffff; border-radius: 50%;"></div>
      </div>
    </div>
  `,
  iconSize: [36, 36],
  iconAnchor: [18, 18],
});

// Component to smoothly pan the map view when center coordinates change
function MapCenterController({ center }: { center: [number, number] }) {
  const map = useMap();
  useEffect(() => {
    map.setView(center, map.getZoom(), { animate: true });
  }, [center, map]);
  return null;
}

// Draggable and clickable marker on the map
function InteractiveMarker({
  position,
  onLocationChange,
}: {
  position: [number, number];
  onLocationChange: (lat: number, lng: number) => void;
}) {
  const markerRef = useRef<L.Marker>(null);

  useMapEvents({
    click(e) {
      onLocationChange(e.latlng.lat, e.latlng.lng);
    },
  });

  const eventHandlers = useMemo(
    () => ({
      dragend() {
        const marker = markerRef.current;
        if (marker != null) {
          const latlng = marker.getLatLng();
          onLocationChange(latlng.lat, latlng.lng);
        }
      },
    }),
    [onLocationChange]
  );

  return (
    <Marker
      draggable={true}
      eventHandlers={eventHandlers}
      position={position}
      ref={markerRef}
      icon={userLocationIcon}
    >
      <Popup>
        <div className="font-sans text-xs p-1">
          <p className="font-bold text-slate-900 m-0 text-sm">Titik Lokasi Laporan</p>
          <p className="text-slate-600 m-0 mt-1">
            Geser pin ini atau klik pada peta untuk menentukan titik presisi kejadian.
          </p>
        </div>
      </Popup>
    </Marker>
  );
}

interface ReportLocationMapProps {
  lat: number;
  lng: number;
  onLocationChange: (lat: number, lng: number) => void;
  isGps: boolean;
  onLocateGps: () => void;
  gpsLoading: boolean;
  gpsError?: string;
}

export default function ReportLocationMap({
  lat,
  lng,
  onLocationChange,
  isGps,
  onLocateGps,
  gpsLoading,
  gpsError,
}: ReportLocationMapProps) {
  const [mapReady, setMapReady] = useState(false);

  useEffect(() => {
    setMapReady(true);
  }, []);

  const TILE_URL =
    process.env.NEXT_PUBLIC_MAP_TILE_URL ||
    "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png";
  const TILE_ATTR =
    process.env.NEXT_PUBLIC_MAP_ATTRIBUTION ||
    '&copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors';

  const position: [number, number] = [lat, lng];

  return (
    <div className="space-y-3">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <label className="text-sm font-semibold text-text flex items-center gap-2">
          <MapPin className="h-4 w-4 text-primary" />
          Peta Lokasi Kejadian (Verifikasi Titik Anda)
        </label>
        
        <button
          type="button"
          onClick={onLocateGps}
          disabled={gpsLoading}
          className="inline-flex items-center gap-2 px-3.5 py-1.5 text-xs font-bold rounded-lg bg-primary-tint text-primary hover:bg-primary hover:text-white transition-all disabled:opacity-50 cursor-pointer w-fit shadow-xs"
        >
          <Crosshair className={`h-3.5 w-3.5 ${gpsLoading ? "animate-spin" : ""}`} />
          {gpsLoading ? "Mendeteksi GPS..." : "Deteksi Lokasi Saya (GPS)"}
        </button>
      </div>

      {/* Map Container */}
      <div className="relative w-full h-[280px] sm:h-[320px] rounded-2xl overflow-hidden border border-border shadow-inner bg-bg">
        {mapReady ? (
          <MapContainer
            center={position}
            zoom={15}
            className="w-full h-full z-0"
            scrollWheelZoom={false}
          >
            <TileLayer attribution={TILE_ATTR} url={TILE_URL} />
            <MapCenterController center={position} />
            <InteractiveMarker
              position={position}
              onLocationChange={onLocationChange}
            />
          </MapContainer>
        ) : (
          <div className="w-full h-full flex items-center justify-center bg-surface/50">
            <div className="flex flex-col items-center gap-2">
              <div className="w-8 h-8 border-3 border-primary/20 border-t-primary rounded-full animate-spin"></div>
              <span className="text-xs text-text-muted">Memuat peta OpenStreetMap...</span>
            </div>
          </div>
        )}

        {/* Floating status badge */}
        <div className="absolute top-3 left-3 z-10 pointer-events-none">
          <span
            className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold shadow-md backdrop-blur-md border ${
              isGps
                ? "bg-emerald-500/90 text-white border-emerald-400"
                : "bg-surface/90 text-text border-border"
            }`}
          >
            {isGps ? (
              <>
                <CheckCircle2 className="h-3.5 w-3.5" />
                GPS Perangkat Terhubung
              </>
            ) : (
              <>
                <MapPin className="h-3.5 w-3.5 text-primary" />
                Titik Manual / Wilayah Terpilih
              </>
            )}
          </span>
        </div>
      </div>

      {/* Coordinate & Guide Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between text-xs text-text-muted gap-2 px-1">
        <div className="flex items-center gap-2">
          <span className="font-mono bg-surface px-2.5 py-1 rounded-md border border-border text-text">
            Lat: {lat.toFixed(6)}, Lng: {lng.toFixed(6)}
          </span>
          <span className="hidden sm:inline">•</span>
          <span>Klik atau geser pin di peta untuk menyesuaikan lokasi rumah/gang.</span>
        </div>
      </div>

      {gpsError && (
        <div className="bg-amber-500/10 border border-amber-500/20 text-amber-700 dark:text-amber-400 text-xs p-3 rounded-xl flex items-start gap-2">
          <AlertCircle className="h-4 w-4 shrink-0 mt-0.5" />
          <p>{gpsError}</p>
        </div>
      )}
    </div>
  );
}
