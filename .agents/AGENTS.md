# Panduan Utama AI Agent: Proyek DBD Guard (Kota Depok)

Anda saat ini beroperasi di dalam *workspace* proyek **DBD Guard** (Sistem Pelaporan & Manajemen Fogging DBD Terintegrasi untuk **Kota Depok**). 

## ATURAN MUTLAK (CRITICAL RULES)
Sebagai AI Agent, Anda **WAJIB** mematuhi aturan-aturan berikut sebelum mengeksekusi *task* apa pun:

1. **JANGAN BERASUMSI**: Apabila Anda diminta untuk membuat fitur, mendesain *database*, atau menyusun arsitektur, Anda **WAJIB** membaca dokumen PRD utama terlebih dahulu di: `e:\Projects\DBD_Guard\PRD_DBD_Guard_Komprehensif.md`.
2. **TETAP DALAM BATASAN PRD**: Jangan pernah menciptakan fitur baru yang tidak ada di dalam PRD. Fokus hanya pada ruang lingkup yang ditetapkan (7 peran pengguna utama, SLA 3x24 jam, fitur PE, inventaris terpusat, dan GIS).
3. **100% GRATIS/OPEN-SOURCE**: DILARANG KERAS menyarankan atau mengimplementasikan layanan pihak ketiga berbayar (seperti Twilio, Fonnte WA Gateway, Mapbox, Google Maps API). Gunakan FCM (Firebase) + Email SMTP untuk notifikasi, dan Leaflet.js + OSM untuk pemetaan.
4. **PATUHI ATURAN TAMBAHAN**: Selalu rujuk dan patuhi aturan teknis di file `rules/structure.md`, standar penulisan kode dan UI di `rules/code_style.md`, panduan warna dan visual di `rules/ui_design.md`, serta batasan bisnis di file `rules/business_rules.md` (di dalam folder `.agents/rules/`).

Setiap kali Anda menerima instruksi dari user (manusia) untuk memodifikasi kode atau merancang sesuatu, pastikan itu tidak melanggar batasan-batasan di atas. Jika instruksi user bertentangan dengan PRD, Anda harus memberikan peringatan kepada user terlebih dahulu sebelum melanjutkannya.
