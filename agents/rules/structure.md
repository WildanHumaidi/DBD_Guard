# Panduan Arsitektur & Tech Stack: DBD Guard

## 1. Teknologi Wajib
Anda harus menggunakan teknologi berikut dalam implementasi kode, **DILARANG KERAS** menggunakan layanan berbayar (*paywalled*) sesuai PRD Kota Depok:
*   **Backend Utama**: Golang (Go) dengan *framework* HTTP Fiber.
*   **Database**: PostgreSQL dengan ekstensi **PostGIS**. Seluruh koordinat dan radius WAJIB dikalkulasi menggunakan fungsi *native* spasial PostGIS (seperti `ST_DistanceSphere`, `ST_DWithin`).
*   **Message Broker**: RabbitMQ atau Redis (digunakan untuk *queue* antrean email & pengiriman notifikasi FCM massal).
*   **Frontend (Admin & Portal Warga)**: Next.js (React) dipadukan dengan Tailwind CSS.
*   **Pemetaan Frontend**: **Leaflet.js** dipadukan dengan *tiles* gratis OpenStreetMap (OSM).
*   **Aplikasi Mobile**: Flutter (Dart).
*   **Penyimpanan File**: MinIO (Self-hosted S3-compatible).

## 2. Struktur Proyek
Karena ini adalah proyek *multi-platform*, pastikan *source code* terisolasi secara rapi di dalam *root directory*:
*   `e:\Projects\DBD_Guard\backend` -> Kode API / Microservices (Golang).
*   `e:\Projects\DBD_Guard\frontend-web` -> Kode Next.js (Admin & Portal Warga).
*   `e:\Projects\DBD_Guard\mobile-app` -> Kode Flutter.

## 3. Aturan Desain API
1.  **RESTful Standard**: Gunakan HTTP status code baku (200, 201, 400, 401, 403, 404, 500). Gunakan *response wrapper* JSON standar berlapis (Status, Data/Message).
2.  **SLA & Timeout**: API yang melakukan proses panjang (seperti *broadcast* FCM atau *generate report*) harus didelegasikan secara asinkron (menggunakan *Goroutines* atau diteruskan ke Message Broker) agar tidak memblokir respon *client*.
3.  **Endpoint Khusus Peran**: Rute API sebaiknya mencerminkan isolasi peran (contoh `/api/v1/admin/...`, `/api/v1/puskesmas/...`, `/api/v1/warga/...`).
