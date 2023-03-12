BEGIN;
    CREATE TABLE IF NOT EXISTS users (
        login    	VARCHAR PRIMARY KEY,
        password    VARCHAR NOT NULL,
        created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

    CREATE TYPE status_type AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

    CREATE TABLE IF NOT EXISTS orders (
        id    	   INTEGER PRIMARY KEY,
        status     status_type DEFAULT 'NEW',
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
COMMIT;

