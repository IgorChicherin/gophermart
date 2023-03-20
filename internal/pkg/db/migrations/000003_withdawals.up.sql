CREATE TABLE IF NOT EXISTS withdrawals
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER,
    order_id     VARCHAR,
    processed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_withdrawals_user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);