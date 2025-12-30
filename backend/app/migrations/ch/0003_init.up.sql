
CREATE TABLE IF NOT EXISTS exception_stack_traces (
    transaction_id Nullable(String),
    exception_hash String,
    stack_trace String,
    recorded_at DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(recorded_at)
ORDER BY (recorded_at, exception_hash);
