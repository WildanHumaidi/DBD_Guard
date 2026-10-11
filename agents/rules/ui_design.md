# Panduan Desain UI / Stitch Prompt — DBD Guard

Ini adalah panduan utama untuk mendesain antarmuka pengguna (UI) aplikasi DBD Guard. Saat Anda ditugaskan untuk membuat UI, komponen, atau halaman baru, gunakan panduan ini sebagai aturan warna, layout, dan _tone_ komunikasi.

---

## Identity & context

**App name:** DBD Guard  
**Platform:** Responsive web app (Next.js + Tailwind CSS), mobile-first  
**Domain:** Civic health reporting — residents in Depok, Indonesia report stagnant-water / dengue-mosquito breeding sites; local health officers verify, investigate, and schedule fogging in response.  
**Locale:** Indonesian (Bahasa Indonesia). All labels, headings, placeholder text, and microcopy must be in Bahasa Indonesia.

### Users (6 roles)

| Role | What they do in the app |
|---|---|
| **Warga** (Resident) | Submit a report with photo + location, track status via reference number / QR, give feedback after fogging |
| **Ketua RT/RW** (Neighborhood head) | Verify or reject reports within their area (SLA 3×24 h), receive early-warning if neighboring area has cases |
| **Petugas Puskesmas** (Health-center officer) | Fill a short PE (epidemiological investigation) form, manage inventory |
| **Kepala Puskesmas** (Health-center head) | Approve fogging based on system recommendation, manage team schedule & stock requests |
| **Tim Fogging** (Fogging crew) | Receive schedule, update status (OTW → Started → Done), upload proof photo + log chemical usage |
| **Dinas Kesehatan** (City health department) | Full dashboard — heatmap, stock monitoring across 38 puskesmas, KLB alerts, public transparency page |

### Report lifecycle (status flow)

```
Diterima → Terverifikasi RT → PE Dilakukan → Dijadwalkan Fogging → Fogging Selesai
                 ↘ Ditolak RT → Banding Warga → Ditinjau Puskesmas
```

---

## Tone & feel

Calm, trustworthy, clinical-but-approachable — like a local Puskesmas health worker explaining calmly, **not** an emergency siren.

- The app must never feel alarming to open, even though its subject is disease outbreaks.
- Reports are about real families in real neighborhoods — every card represents one household.
- The design should communicate: "someone responsible is on the other side receiving your report."

---

## Color system

### Light mode

| Token | Hex | Usage |
|---|---|---|
| `primary` | `#0F766E` | Main CTA buttons ("Lapor Sekarang"), links, active icons, progress bars |
| `primary-deep` | `#134E4A` | Header bar, navigation bar, institution name area |
| `primary-tint` | `#E6F4F1` | Background for educational cards, 3M tips, info boxes |
| `text` | `#0F172A` | Headings, reporter names, addresses, case numbers |
| `text-muted` | `#475569` | Timestamps, distances, supporting descriptions |
| `border` | `#E2E8F0` | Card edges, list dividers, column separators |
| `bg` | `#F8FAFC` | Full-page background |
| `surface` | `#FFFFFF` | Report cards, schedule panels, dialogs |
| `status-wait` | `#64748B` | "Menunggu Verifikasi" — not yet reviewed |
| `status-verified` | `#2563EB` | "Terverifikasi" — confirmed by RT/RW |
| `status-scheduled` | `#B45309` | "Dijadwalkan" — fogging date set (never use as text color on white — contrast fails) |
| `status-done` | `#15803D` | "Selesai" — fogging completed |
| `status-void` | `#57534E` | "Bukan Kasus" — duplicate / unverified / closed without action |
| `risk-low` | `#16A34A` | Map zone: Aman (safe) |
| `risk-medium` | `#EAB308` | Map zone: Waspada (watch) |
| `risk-high` | `#EA580C` | Map zone: Rawan (at risk) |
| `risk-outbreak` | `#B91C1C` | Map zone: KLB — **the ONLY place red appears in the entire app** |

### Dark mode

| Token | Hex |
|---|---|
| `primary` | `#2DD4BF` |
| `primary-deep` | `#0F3B38` |
| `primary-tint` | `#123936` |
| `text` | `#E8EDF4` |
| `text-muted` | `#93A3B8` |
| `border` | `#1F2A3C` |
| `bg` | `#0B1220` (navy-black, not pure black — field workers use it at dawn/night) |
| `surface` | `#121A2A` |
| `status-wait` | `#94A3B8` |
| `status-verified` | `#60A5FA` |
| `status-scheduled` | `#FBBF24` |
| `status-done` | `#4ADE80` |
| `status-void` | `#A8A29E` |
| `risk-low` | `#22C55E` |
| `risk-medium` | `#FACC15` |
| `risk-high` | `#FB923C` |
| `risk-outbreak` | `#EF4444` |

### CSS custom properties (embed in design)

```css
:root {
  --primary: #0F766E;
  --primary-deep: #134E4A;
  --primary-tint: #E6F4F1;
  --text: #0F172A;
  --text-muted: #475569;
  --border: #E2E8F0;
  --bg: #F8FAFC;
  --surface: #FFFFFF;
  --status-wait: #64748B;
  --status-verified: #2563EB;
  --status-scheduled: #B45309;
  --status-done: #15803D;
  --status-void: #57534E;
  --risk-low: #16A34A;
  --risk-medium: #EAB308;
  --risk-high: #EA580C;
  --risk-outbreak: #B91C1C;
}

@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) {
    --primary: #2DD4BF;
    --primary-deep: #0F3B38;
    --primary-tint: #123936;
    --text: #E8EDF4;
    --text-muted: #93A3B8;
    --border: #1F2A3C;
    --bg: #0B1220;
    --surface: #121A2A;
    --status-wait: #94A3B8;
    --status-verified: #60A5FA;
    --status-scheduled: #FBBF24;
    --status-done: #4ADE80;
    --status-void: #A8A29E;
    --risk-low: #22C55E;
    --risk-medium: #FACC15;
    --risk-high: #FB923C;
    --risk-outbreak: #EF4444;
  }
}
```

---

## Strict design rules

1. **Red is reserved.** `risk-outbreak` / `#B91C1C` only ever appears on the KLB outbreak map marker. Never use red for buttons, delete confirmations, error toasts, form validation, or badges anywhere else. Its power comes from scarcity.

2. **Primary action = teal only.** All main CTA buttons use `primary`. Never use a status color or risk color for a clickable button.

3. **Color + icon + label, always.** Every status chip and risk badge must show an icon and a text label alongside its color. Never rely on color alone (accessibility for red-green color blindness).

4. **Amber text contrast.** In light mode, use `#B45309` for text on white. Use `#F59E0B` only as chip/badge fill, never as text color — it fails WCAG AA contrast on white.

5. **Map zone polygons at ~35% opacity** with full-color stroke, so street and alley names underneath remain readable.

6. **Visual style:**
   - Rounded corners (`rounded-xl` for cards, `rounded-lg` for buttons, `rounded-full` for chips/badges)
   - Generous padding and whitespace
   - Flat surfaces — no gradients, no drop shadows in light mode; elevation in dark mode is conveyed by surface color difference only
   - Clean sans-serif typeface (Inter or similar), medium weight for headings
   - Card-based layout for reports; map-based layout for zone risk view

7. **Dark mode is navy-black (`#0B1220`), not pure black.** Field workers operate at dawn and night — pure black causes eye strain when switching between screen and dark alleys.

8. **Support both light and dark mode** for every screen. Default follows system preference; provide a manual toggle in settings.

---

## Pages to design

Below are the key screens in the app. When generating, pick from this list:

### Warga (Resident)
- **Beranda Warga** — Welcome, quick "Lapor Sekarang" CTA, recent reports summary, educational 3M tips card
- **Form Lapor** — Photo upload, cascading location dropdown (Kecamatan → Kelurahan → RW → RT), anonymous toggle, description
- **Tracking Laporan** — Status timeline (vertical stepper with status colors + icons), report detail, reference number, QR code
- **Feedback & Rating** — Post-fogging satisfaction form

### RT/RW (Neighborhood Head)
- **Dashboard RT/RW** — Incoming reports list (cards), verify/reject actions, SLA countdown timer, EWS alert banner
- **Detail Verifikasi** — Report detail with photo, map pin, approve/reject with reason field

### Puskesmas (Health Center)
- **Dashboard Puskesmas** — Verified reports queue, PE form entry, inventory summary widget, fogging schedule calendar
- **Form PE** — Simple decision-support form (jentik found? nearby fever count? photo?)
- **Inventaris** — Stock table per item, usage logs, minimum-stock warning, request-stock button

### Kepala Puskesmas
- **Approval Fogging** — Recommendation card (system-generated), approve/reject, team assignment, date picker (teams greyed out if already booked)

### Tim Fogging
- **Jadwal Hari Ini** — Today's assignment card, map to location, status update buttons (OTW → Mulai → Selesai), photo upload + chemical usage input

### Dinas Kesehatan (City Level)
- **GIS Heatmap** — Full-screen Leaflet map with risk-zone polygons, report pins, KLB alert overlay
- **Monitoring Stok** — Table/cards for all 38 puskesmas stock levels, approve stock requests, cross-allocation
- **Dashboard Publik** — Public transparency page with aggregate stats (no personal data)

---

## Panduan Penggunaan
Saat diminta untuk mendesain atau mengimplementasi halaman tertentu (misalnya "Beranda Warga", "Dashboard RT/RW", "Form Lapor"), ikuti panduan dan spesifikasi desain di atas. Pastikan juga implementasinya menggunakan Tailwind CSS sesuai dengan token warna yang didefinisikan.
