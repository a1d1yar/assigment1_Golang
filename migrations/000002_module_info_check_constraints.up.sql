-- Add constraint to check updated_at column values cannot be before created_at
ALTER TABLE module_info ADD CONSTRAINT check_updated_at AFTER updated_at 
CHECK (updated_at >= created_at);

-- Add constraint to check module_duration column values must be between 5 and 15
ALTER TABLE module_info ADD CONSTRAINT check_module_duration AFTER module_duration 
CHECK (module_duration > 5 AND module_duration <= 15);