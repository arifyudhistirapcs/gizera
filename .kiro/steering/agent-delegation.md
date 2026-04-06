---
inclusion: auto
---

# Agent Delegation Rules

Semua pekerjaan implementasi, testing, dan review HARUS didelegasikan ke subagent sesuai expertise. Jangan mengerjakan sendiri jika ada subagent yang lebih tepat.

## IT Solution Architect (Pengawas)

Subagent `it-solution-architect` bertindak sebagai pengawas dan koordinator. Libatkan untuk:
- Review arsitektur sebelum implementasi fitur besar
- Validasi keputusan teknis yang berdampak lintas modul
- Code review setelah implementasi selesai
- Resolusi konflik antar modul atau antar subagent

## Delegation Matrix

| Jenis Pekerjaan | Subagent | Kapan Digunakan |
|---|---|---|
| Go handler, service, model, middleware, routes | `backend-dev` | Semua perubahan di `backend/` |
| Vue 3 + Ant Design Vue views, components, services | `frontend-dev` | Semua perubahan di `web/` |
| Vue 3 + Vant UI views, composables, offline, PWA | `mobile-dev` | Semua perubahan di `pwa/` |
| Schema design, migration, query optimization, indexing | `database-engineer` | Perubahan schema, slow query, data modeling |
| Docker, Nginx, CI/CD, deployment, monitoring | `infra-engineer` | Setup infra, performance tuning, deployment |
| Cross-stack E2E testing, API contract validation | `integration-tester` | Testing full flow lintas stack |
| OWASP testing, auth bypass, injection, security audit | `security-tester` | Audit keamanan, vulnerability scanning |
| Screen design, user flow, design system, wireframe | `uiux-designer` | Desain UI baru, UX improvement |
| Visual testing, responsiveness, accessibility QA | `uiux-tester` | QA visual, accessibility, responsive |
| API docs, user manual, architecture diagram, ADR | `technical-writer` | Dokumentasi teknis dan user |
| Landing page copy, UX writing, marketing content | `content-marketer` | Copy user-facing, microcopy, marketing |

## Workflow untuk Fitur Baru

1. **Analisis** → `it-solution-architect` review requirements dan tentukan scope
2. **Database** → `database-engineer` design schema jika perlu tabel baru
3. **Backend** → `backend-dev` implementasi model, service, handler, routes
4. **Frontend** → `frontend-dev` implementasi web dashboard views
5. **Mobile** → `mobile-dev` implementasi PWA mobile views (jika ada)
6. **Testing** → `integration-tester` validasi E2E flow
7. **Review** → `it-solution-architect` final review

## Workflow untuk Bug Fix

1. **Analisis** → identifikasi root cause dan tentukan area (backend/frontend/database)
2. **Fix** → delegate ke subagent sesuai area
3. **Verify** → `integration-tester` atau subagent yang relevan verifikasi fix

## Workflow untuk Security/Performance

1. **Audit** → `security-tester` atau `database-engineer` analisis
2. **Fix** → delegate ke subagent sesuai area yang perlu diperbaiki
3. **Verify** → auditor yang sama verifikasi fix

## Aturan Penting

- JANGAN mengerjakan backend code jika bukan `backend-dev`
- JANGAN mengerjakan frontend code jika bukan `frontend-dev` atau `mobile-dev`
- JANGAN mengubah schema tanpa konsultasi `database-engineer`
- Setiap perubahan arsitektur signifikan HARUS di-review oleh `it-solution-architect`
- Jika pekerjaan melibatkan lebih dari satu area, koordinasikan lewat `it-solution-architect`
