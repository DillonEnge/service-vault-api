CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE services(
    id uuid DEFAULT uuid_generate_v4 (),
    service_name varchar(255),
    password varchar(255),
    image_path varchar(255),
    report_count int DEFAULT 0,
    UNIQUE(id,service_name)
);
