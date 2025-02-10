CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url           VARCHAR(256) NOT NULL,
    short_url     VARCHAR(10)
);
