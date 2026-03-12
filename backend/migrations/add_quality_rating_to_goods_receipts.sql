-- Migration: Add quality_rating column to goods_receipts table
-- This allows rating to be recorded when goods are received (GRN created)
-- Rating is then used to calculate average supplier quality rating

-- Add quality_rating column to goods_receipts
ALTER TABLE goods_receipts ADD COLUMN quality_rating REAL DEFAULT 0;

-- Update existing records to have 0 rating (neutral)
UPDATE goods_receipts SET quality_rating = 0 WHERE quality_rating IS NULL;
