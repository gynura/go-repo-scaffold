CREATE TABLE sv_users
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
CREATE INDEX idx_sv_users_country ON sv_users (country);
