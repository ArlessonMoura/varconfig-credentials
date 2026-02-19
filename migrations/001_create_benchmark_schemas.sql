-- Migration: create benchmark_schemas with jsonb schema_body
-- Creates table if not exists and ensures schema_body is jsonb with GIN index

BEGIN;

-- Create table if it does not exist
CREATE TABLE IF NOT EXISTS benchmark_schemas (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  schema_body JSONB NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- If column exists but is not jsonb, try to convert it (best-effort)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'benchmark_schemas' AND column_name = 'schema_body' AND data_type <> 'jsonb'
  ) THEN
    ALTER TABLE benchmark_schemas
    ALTER COLUMN schema_body TYPE jsonb USING schema_body::jsonb;
  END IF;
EXCEPTION WHEN others THEN
  -- ignore conversion errors
  RAISE NOTICE 'Could not convert schema_body to jsonb: %', SQLERRM;
END$$;

-- Create GIN index for jsonb queries
CREATE INDEX IF NOT EXISTS idx_benchmark_schema_body_gin ON benchmark_schemas USING gin (schema_body);

COMMIT;
