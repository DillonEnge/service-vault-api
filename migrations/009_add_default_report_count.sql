-- Write your migrate up statements here

ALTER TABLE services
ALTER report_count
SET DEFAULT 0;

---- create above / drop below ----

ALTER TABLE services
ALTER report_count
DROP DEFAULT;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
