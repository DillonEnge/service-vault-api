-- Write your migrate up statements here

ALTER TABLE services
DROP COLUMN reported;

ALTER TABLE services
ADD report_count int;

---- create above / drop below ----

ALTER TABLE services
DROP COLUMN report_count;

ALTER TABLE services
ADD COLUMN reported boolean DEFAULT FALSE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
