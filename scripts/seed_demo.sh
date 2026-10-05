#!/usr/bin/env bash
# ==============================================================
# DBD Guard - Seed Data Demo Idempoten (Kota Depok)
# Skenario: Laporan -> Verifikasi -> PE -> SPK -> Fogging -> Potong Stok
# ==============================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="${SCRIPT_DIR}/../backend"

echo "=========================================================="
echo "🚀 Menjalankan Database Seeder Demo DBD Guard (Kota Depok)"
echo "=========================================================="

cd "${BACKEND_DIR}"
echo "[1/2] Menjalankan seeder Go..."
go run ./cmd/seeder/main.go
echo "[2/2] ✅ Seeder demo berhasil dijalankan secara idempoten!"
