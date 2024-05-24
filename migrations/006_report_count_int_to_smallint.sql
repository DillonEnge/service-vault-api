-- Write your migrate up statements here

ALTER TABLE services
ALTER COLUMN report_count TYPE smallint;

---- create above / drop below ----

ALTER TABLE services
ALTER COLUMN report_count TYPE int;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
