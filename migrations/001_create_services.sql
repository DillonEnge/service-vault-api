-- Write your migrate up statements here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE services(
    id uuid DEFAULT uuid_generate_v4 (),
    service_name varchar(255),
    password varchar(255),
    UNIQUE(id,service_name)
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

DROP TABLE services;
