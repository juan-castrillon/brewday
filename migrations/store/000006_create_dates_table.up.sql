CREATE TABLE
    IF NOT EXISTS "dates" (
        id INTEGER NOT NULL PRIMARY KEY,
        date TEXT,
        name TEXT,
        recipe_id INTEGER NOT NULL,
        FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE INDEX IF NOT EXISTS ix_dates ON "dates" (recipe_id);