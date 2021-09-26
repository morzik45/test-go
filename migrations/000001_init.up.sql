CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE authorizations (
    id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    is_authorized BOOLEAN NOT NULL DEFAULT TRUE,
    login_at TIMESTAMP NOT NULL DEFAULT NOW(),
    logout_at TIMESTAMP,
    session_token VARCHAR(50) NOT NULL
);