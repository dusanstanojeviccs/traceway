CREATE TABLE IF NOT EXISTS metric_records (
    name LowCardinality(String),
    value Float64,
    recorded_at DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(recorded_at)
ORDER BY (recorded_at, name);
