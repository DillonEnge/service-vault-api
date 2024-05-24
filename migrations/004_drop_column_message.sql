-- Write your migrate up statements here

ALTER TABLE requests
DROP COLUMN message;

---- create above / drop below ----

ALTER TABLE requests
ADD message varchar(255);

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
