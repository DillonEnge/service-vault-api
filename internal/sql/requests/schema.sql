CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE requests(
    id uuid DEFAULT uuid_generate_v4 (),
    service_name varchar(255),
    email varchar(255),
    type varchar(255),
    UNIQUE(service_name,email,type)
);
