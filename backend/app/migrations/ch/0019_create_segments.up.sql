CREATE TABLE IF NOT EXISTS segments (
    id String,
    transaction_id String,
    project_id LowCardinality(String),
    name LowCardinality(String),
    start_time DateTime64(6),
    duration Int64,
    recorded_at DateTime,

    INDEX idx_transaction_id transaction_id TYPE bloom_filter(0.001) GRANULARITY 1
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(recorded_at)
ORDER BY (project_id, transaction_id, start_time)
SETTINGS index_granularity = 8192
