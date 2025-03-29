DROP DATABASE IF EXISTS micro_users_db;
CREATE DATABASE micro_users_db
WITH OWNER = postgres
     ENCODING = 'UTF8'
     TABLESPACE = pg_default
     CONNECTION LIMIT = -1;

\c micro_users_db;

CREATE TABLE svc_users
(
    id         VARCHAR(64)  NOT NULL,
    first_name VARCHAR(64)  NOT NULL,
    last_name  VARCHAR(64)  NOT NULL,
    nickname   VARCHAR(64)  NOT NULL UNIQUE,
    password   TEXT         NOT NULL,
    email      VARCHAR(256) NOT NULL UNIQUE,
    country    VARCHAR(64)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL
);
CREATE INDEX idx_svc_users_country ON svc_users (country);
