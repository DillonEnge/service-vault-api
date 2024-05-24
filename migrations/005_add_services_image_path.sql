-- Write your migrate up statements here

ALTER TABLE services
ADD image_path varchar(255);

---- create above / drop below ----

ALTER TABLE services
DROP COLUMN image_path;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
