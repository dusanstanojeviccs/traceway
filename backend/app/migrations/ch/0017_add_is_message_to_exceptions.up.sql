ALTER TABLE exception_stack_traces ADD COLUMN IF NOT EXISTS is_message UInt8 DEFAULT 0;
