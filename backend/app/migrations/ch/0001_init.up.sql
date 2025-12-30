CREATE TABLE IF NOT EXISTS transactions (
    id String,
    endpoint String,
    duration Int64,
    recorded_at DateTime,
    status_code Int32,
    body_size Int32,
    client_ip String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(recorded_at)
ORDER BY (recorded_at, endpoint);
