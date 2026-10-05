# ==============================================================
# DBD Guard - Seed Data Demo Idempoten (Kota Depok)
# Skenario: Laporan -> Verifikasi -> PE -> SPK -> Fogging -> Potong Stok
# ==============================================================

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "🚀 Menjalankan Database Seeder Demo DBD Guard (Kota Depok)" -ForegroundColor Cyan
Write-Host "   Target: PostgreSQL 15 + PostGIS (Docker / Localhost)" -ForegroundColor Yellow
Write-Host "==========================================================" -ForegroundColor Cyan

$backendDir = Join-Path $PSScriptRoot "..\backend"

if (-not (Test-Path $backendDir)) {
    Write-Error "Direktori backend tidak ditemukan di: $backendDir"
    exit 1
}

Push-Location $backendDir
try {
    Write-Host "`n[1/2] Mengompilasi dan mengeksekusi seeder Go..." -ForegroundColor Green
    go run ./cmd/seeder/main.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n[2/2] ✅ Seeder demo berhasil dijalankan!" -ForegroundColor Green
        Write-Host "Skenario yang siap dipresentasikan:" -ForegroundColor White
        Write-Host "  1. Alur Lengkap: Ref No. DBD-DPK-2026-DEMO01 (Status: COMPLETED, Stok Terpotong 2L)" -ForegroundColor Gray
        Write-Host "  2. Kelurahan Siaga: Kelurahan Beji (6 Kasus Aktif, Peringatan EWS Muncul)" -ForegroundColor Gray
        Write-Host "  3. Ditolak Lintas Wilayah: Ref No. DBD-DPK-2026-REJECT01 (Alasan Yurisdiksi)" -ForegroundColor Gray
        Write-Host "  4. Uji Coba SPK Validasi PE: Ref No. DBD-DPK-2026-SPK01 (Rekomendasi Fogging Fokus)" -ForegroundColor Gray
    } else {
        Write-Error "Seeder gagal dengan kode keluar $LASTEXITCODE"
    }
} finally {
    Pop-Location
}
