-- +goose Up
-- SUBMISSIONS

-- enums
CREATE TYPE submission_status AS ENUM (
    'UPLOADING',        -- uploading artifcats to the object storage
    'UPLOADED',         -- upload success
    'DEPLOYING',        -- deploying in container
    'DEPLOY_FAILED',    -- if there is some error while deploying the container
    'BUILDING',         -- building binary inside the container, skipped if the user submission is binary
    'BUILD_FAILED',     -- if building failed
    'BINARY_READY',     -- binary ready for execution
    'REJECTED',         -- binary security vulnerable
    'SUCCESS'           -- ready for benchmarking
);

CREATE TYPE supported_language AS ENUM (
    'CPP',
    'RUST'
);

CREATE TABLE submissions (
    id                  SERIAL      PRIMARY KEY,
    user_id             INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 
    -- Storage
    storage_key         TEXT        NOT NULL UNIQUE,
    file_size_bytes     BIGINT      NOT NULL,
    checksum_sha256     TEXT        NOT NULL,
    original_filename   TEXT       NOT NULL,
 

    language            supported_language        NOT NULL,
 
    -- Lifecycle
    status              submission_status  NOT NULL DEFAULT 'UPLOADING',
    deploy_log          TEXT,
    build_log           TEXT,       -- compiler output, populated on build_failed
    error_message       TEXT,

    submitted_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_submissions
    BEFORE UPDATE ON submissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
 
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_status  ON submissions(status);
 

-- ════════════════════════════════════════════════════════════
-- SCORES
-- Final aggregated score for each submission after
-- the benchmark completes. One row per submission.
-- Raw per-event telemetry stays in the bot fleet service.
-- ════════════════════════════════════════════════════════════
 
CREATE TABLE scores (
    id                  SERIAL      PRIMARY KEY,
    submission_id       INT         NOT NULL UNIQUE REFERENCES submissions(id) ON DELETE CASCADE,
    user_id             INT         NOT NULL REFERENCES users(id),
                                    -- denormalised for fast leaderboard queries
 
    -- Latency (microseconds — pure contestant engine time)
    p50_latency_us      BIGINT      NOT NULL,
    p90_latency_us      BIGINT      NOT NULL,
    p99_latency_us      BIGINT      NOT NULL,
 
    -- Throughput
    peak_tps            INT         NOT NULL,
    sustained_tps       INT         NOT NULL,
    total_orders_sent   BIGINT      NOT NULL,
    total_orders_acked  BIGINT      NOT NULL,
    timeout_count       BIGINT      NOT NULL DEFAULT 0,
 
    -- Correctness
    total_fills         BIGINT      NOT NULL,
    correct_fills       BIGINT      NOT NULL,
    violation_count     INT         NOT NULL DEFAULT 0,
    correctness_rate    NUMERIC(5,4) NOT NULL
                        CHECK (correctness_rate BETWEEN 0 AND 1),
 
    -- Composite score (correctness 25%, latency 40%, throughput 35%)
    -- Normalised 0.0 – 100.0
    composite_score     NUMERIC(6,2) NOT NULL DEFAULT 0,
 
    scored_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
 
CREATE INDEX idx_scores_user_id         ON scores(user_id);
CREATE INDEX idx_scores_composite_score ON scores(composite_score DESC);
 
-- +goose Down
DROP TABLE IF EXISTS scores;
DROP TABLE IF EXISTS submissions;
