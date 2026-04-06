# ERP SPPG - Panduan Deployment

## Prasyarat Server

- Docker Engine 24+
- Docker Compose v2+
- Git
- Minimal 4GB RAM, 2 vCPU
- PostgreSQL 15 (Cloud SQL atau self-hosted)
- SSL certificate (cert.pem + key.pem)
- Firebase service account JSON

---

## Step 1: Clone & Masuk ke Project

```bash
git clone <repo-url> erp-sppg
cd erp-sppg
```

## Step 2: Build Docker Images

```bash
chmod +x deployment/scripts/build.sh
./deployment/scripts/build.sh latest
```

Ini akan build 3 image:
- `erp-sppg-backend:latest`
- `erp-sppg-web:latest`
- `erp-sppg-pwa:latest`

Kalau pakai private registry:
```bash
DOCKER_REGISTRY=registry.example.com ./deployment/scripts/build.sh v1.0.0
```

## Step 3: Setup Environment Variables

```bash
cd deployment
cp .env.example .env
```

Edit `.env` dan isi semua value yang diperlukan:

```bash
nano .env
```

Yang wajib diisi:
| Variable | Keterangan |
|---|---|
| `DB_HOST` | IP/hostname PostgreSQL |
| `DB_PASSWORD` | Password database |
| `REDIS_PASSWORD` | Password Redis |
| `JWT_SECRET` | Secret key JWT (min 32 karakter) |
| `GRAFANA_PASSWORD` | Password admin Grafana |

## Step 4: Setup SSL Certificate

Taruh SSL certificate di folder nginx/ssl:

```bash
mkdir -p nginx/ssl
cp /path/to/cert.pem nginx/ssl/cert.pem
cp /path/to/key.pem nginx/ssl/key.pem
```

Untuk testing/staging bisa pakai self-signed:
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/key.pem \
  -out nginx/ssl/cert.pem \
  -subj "/CN=erp-sppg.example.com"
```

## Step 5: Setup Firebase Credentials

```bash
cp /path/to/firebase-service-account.json ./firebase-credentials.json
```

## Step 6: Update Nginx Server Name

Edit `nginx/nginx.conf`, ganti `erp-sppg.example.com` dengan domain kamu:

```bash
nano nginx/nginx.conf
```

## Step 7: Buat Folder yang Dibutuhkan

```bash
mkdir -p logs/backend-1 logs/backend-2 uploads backups nginx/logs
```

## Step 8: Jalankan Semua Services

```bash
docker compose -f docker-compose.prod.yml up -d
```

## Step 9: Verifikasi

Cek semua container running:
```bash
docker compose -f docker-compose.prod.yml ps
```

Cek health backend:
```bash
curl -k https://localhost/health
```

Cek logs kalau ada masalah:
```bash
docker logs erp-sppg-backend-1
docker logs erp-sppg-nginx
```

---

## Monitoring

| Service | URL |
|---|---|
| Web Dashboard | https://your-domain/ |
| PWA Mobile | https://your-domain/pwa/ |
| API | https://your-domain/api/ |
| Grafana | http://your-server:3000 |
| Prometheus | http://your-server:9090 |

---

## Perintah Berguna

```bash
# Restart semua
docker compose -f docker-compose.prod.yml restart

# Restart backend saja
docker compose -f docker-compose.prod.yml restart backend-1 backend-2

# Lihat logs real-time
docker compose -f docker-compose.prod.yml logs -f backend-1

# Stop semua
docker compose -f docker-compose.prod.yml down

# Update image & deploy ulang
./deployment/scripts/build.sh latest
docker compose -f docker-compose.prod.yml up -d

# Backup database manual
docker exec erp-sppg-backup /backup.sh
```

---

## Troubleshooting

| Masalah | Solusi |
|---|---|
| Backend tidak start | Cek `docker logs erp-sppg-backend-1`, pastikan DB_HOST bisa diakses |
| 502 Bad Gateway | Backend belum ready, tunggu healthcheck atau cek logs |
| SSL error | Pastikan cert.pem dan key.pem ada di `nginx/ssl/` |
| Redis connection refused | Pastikan REDIS_PASSWORD di .env sama dengan yang di compose |

---

## Local Development

Untuk development lokal, gunakan compose di root project:

```bash
# Start PostgreSQL + Redis lokal
docker compose up -d

# Jalankan backend
cd backend
cp .env.example .env  # edit sesuai kebutuhan
go run ./cmd/server

# Jalankan web dashboard
cd web
cp .env.example .env
npm install && npm run dev

# Jalankan PWA
cd pwa
cp .env.example .env
npm install && npm run dev
```
