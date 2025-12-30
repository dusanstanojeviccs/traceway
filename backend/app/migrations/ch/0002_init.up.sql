CREATE TABLE IF NOT EXISTS metric_records (
    name String,
    value Float32,
    recorded_at DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(recorded_at)
ORDER BY (recorded_at, name);
