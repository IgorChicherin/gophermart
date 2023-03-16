BEGIN;
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    login      VARCHAR UNIQUE NOT NULL,
    password   VARCHAR        NOT NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TYPE status_type AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders
(
    id         SERIAL PRIMARY KEY,
    order_id   VARCHAR UNIQUE NOT NULL,
    user_id    INTEGER,
    status     status_type             DEFAULT 'NEW',
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
COMMIT;

