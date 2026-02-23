CREATE TABLE
    IF NOT EXISTS "timelines" (
        id INTEGER NOT NULL PRIMARY KEY,
        event TEXT NOT NULL,
        timestamp_unix INTEGER NOT NULL,
        recipe_id INTEGER NOT NULL,
        FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE INDEX IF NOT EXISTS ix_timelines ON "timelines" (recipe_id, timestamp_unix);