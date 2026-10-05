# Aturan Code Style & Desain Antarmuka: DBD Guard

## 1. Backend (Golang)
- **Arsitektur**: Gunakan pola *Clean Architecture* dengan pemisahan *layer* yang jelas: `Handler/Controller` -> `Service/Usecase` -> `Repository`.
- **Naming Convention**: 
  - Gunakan `camelCase` untuk variabel dan fungsi internal, serta `PascalCase` untuk fungsi/struct/konstanta *exported*.
  - Penamaan *file* harus menggunakan `snake_case` (contoh: `report_repository.go`).
- **Error Handling**: JANGAN PERNAH mengabaikan *error* (*swallowing errors*). Selalu tangkap *error*, log, dan kembalikan *error response* JSON berstandar (format `{ "status": ..., "message": ..., "data": ... }`).
- **Notifikasi**: Implementasi Push Notification murni menggunakan Firebase Cloud Messaging (FCM) dan SMTP (Gomail/Nodemailer sejenis). **Dilarang memasukkan pustaka WA API berbayar**.

## 2. Frontend Web (React / Next.js)
- **Paradigma**: Wajib menggunakan *Functional Components* dan *React Hooks*.
- **Routing**: Patuhi standar *App Router* (`src/app`).
- **Styling & Desain**: 
  - WAJIB menggunakan Tailwind CSS.
  - Terapkan desain *premium*, *clean*, *mobile-first*. Gunakan palet kesehatan/pemerintah. **WAJIB merujuk pada `color_palette.md` untuk aturan warna**.
- **Pemetaan (GIS)**: Wajib menggunakan **Leaflet.js + OpenStreetMap (OSM)** untuk render Heatmap dan peta pelaporan. Jangan import *library* komersial berbayar (Mapbox GL JS atau Google Maps).

## 3. Mobile App (Flutter / Dart)
- **Arsitektur & State Management**: Pisahkan *logic* dan UI. Gunakan *state management* yang konsisten (Riverpod/BLoC/Provider).
- **Naming Convention**:
  - `camelCase` untuk variabel/fungsi.
  - `PascalCase` untuk *Class*, *Enum*, *Type*.
  - `snake_case` untuk *file* dan direktori.
- **Clean Widget Tree**: Pecah UI menjadi *widget* modular kecil.
- **Penanganan Async**: Selalu kelola *loading state* & *error state* ke pengguna (*snackbar*).
