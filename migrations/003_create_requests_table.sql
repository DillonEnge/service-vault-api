-- Write your migrate up statements here
CREATE TABLE requests(
    id uuid DEFAULT uuid_generate_v4 (),
    service_name varchar(255),
    email varchar(255),
    type varchar(255),
    message varchar(255),
    UNIQUE(id)
);

---- create above / drop below ----

DROP TABLE requests;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
