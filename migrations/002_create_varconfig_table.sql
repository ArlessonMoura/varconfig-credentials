-- Migration: create var_configs table with jsonb payload
-- Creates table for variable configurations with composite index on org_id + benchmark_id

BEGIN;

CREATE TABLE IF NOT EXISTS var_configs (
  id BIGSERIAL PRIMARY KEY,
  org_id TEXT NOT NULL,
  benchmark_id TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create composite index for queries filtering by org_id and benchmark_id
CREATE INDEX IF NOT EXISTS idx_varcfg_org_bench ON var_configs (org_id, benchmark_id);

-- Create GIN index for efficient JSONB queries if needed
CREATE INDEX IF NOT EXISTS idx_varcfg_payload_gin ON var_configs USING gin (payload);

COMMIT;
