ALTER TABLE exception_stack_traces ADD COLUMN IF NOT EXISTS scope String DEFAULT '{}';
