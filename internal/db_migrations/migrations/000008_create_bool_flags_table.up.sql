CREATE TABLE
    IF NOT EXISTS "bool_flags" (
        id INTEGER NOT NULL PRIMARY KEY,
        value INTEGER,
        name TEXT,
        recipe_id INTEGER NOT NULL,
        FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE INDEX IF NOT EXISTS ix_bool_flags ON "bool_flags" (recipe_id, name);