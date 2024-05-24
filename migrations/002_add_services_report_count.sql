-- Write your migrate up statements here

ALTER TABLE services
ADD report_count int DEFAULT 0;

---- create above / drop below ----

ALTER TABLE services
DROP COLUMN report_count;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
