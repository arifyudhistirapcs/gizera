-- ============================================
-- ERP SPPG - Database Initialization Script
-- Dijalankan otomatis saat PostgreSQL container pertama kali start
-- ============================================

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE erp_sppg TO erp_sppg_user;
