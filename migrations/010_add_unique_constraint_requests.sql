-- Write your migrate up statements here

ALTER TABLE requests
DROP CONSTRAINT requests_id_key;

ALTER TABLE requests
ADD CONSTRAINT requests_unique UNIQUE(service_name,email,type);

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
