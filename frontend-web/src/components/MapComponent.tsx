"use client";

import { useEffect, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import api from '@/lib/api';
import L from 'leaflet';

// Fix Leaflet default icon issues in Next.js
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
});

export default function MapComponent() {
  const [geoData, setGeoData] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMapData = async () => {
      try {
        const response = await api.get('/admin/gis/heatmap');
        if (response.data && response.data.data) {
          setGeoData(response.data.data);
        }
      } catch (error) {
        console.error("Failed to load map data", error);
      } finally {
        setLoading(false);
      }
    };
    fetchMapData();
  }, []);

  if (loading) {
    return <div className="w-full h-full flex items-center justify-center text-slate-400">Loading map...</div>;
  }

  const pointToLayer = (feature: any, latlng: L.LatLng) => {
    const score = feature.properties?.risk_score || 0;
    
    // Color logic
    let color = '#22c55e'; // green (low)
    if (score > 60) color = '#ef4444'; // red (high)
    else if (score > 30) color = '#f59e0b'; // orange (medium)

    return L.circleMarker(latlng, {
      radius: 8 + (score / 10), // bigger radius for higher score
      fillColor: color,
      color: '#fff',
      weight: 1,
      opacity: 1,
      fillOpacity: 0.7
    });
  };

  const onEachFeature = (feature: any, layer: L.Layer) => {
    if (feature.properties) {
      layer.bindPopup(`
        <div class="font-sans">
          <strong class="block text-sm mb-1">Status: ${feature.properties.status}</strong>
          <span class="text-xs text-text-muted block">Risk Score: ${feature.properties.risk_score?.toFixed(1) || 0}</span>
          <span class="text-xs text-text-muted block">Reports: ${feature.properties.report_count || 1}</span>
        </div>
      `);
    }
  };

const TILE_URL = process.env.NEXT_PUBLIC_MAP_TILE_URL || "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png";
const TILE_ATTR = process.env.NEXT_PUBLIC_MAP_ATTRIBUTION || '&copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors';
const DEPOK_CENTER: [number, number] = [-6.402484, 106.794243];

  return (
    <MapContainer 
      center={DEPOK_CENTER}
      zoom={12} 
      className="w-full h-full rounded-xl z-0"
    >
      <TileLayer
        attribution={TILE_ATTR}
        url={TILE_URL}
      />
      {geoData && (
        <GeoJSON 
          data={geoData} 
          pointToLayer={pointToLayer}
          onEachFeature={onEachFeature}
        />
      )}
    </MapContainer>
  );
}
