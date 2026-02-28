CREATE TABLE
    IF NOT EXISTS "indexed" (
        id INTEGER NOT NULL PRIMARY KEY,
        event TEXT NOT NULL
    );

CREATE INDEX IF NOT EXISTS ix_event ON "indexed" (event);