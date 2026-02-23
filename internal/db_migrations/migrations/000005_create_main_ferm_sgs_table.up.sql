CREATE TABLE
    IF NOT EXISTS "main_ferm_sgs" (
        id INTEGER NOT NULL PRIMARY KEY,
        sg REAL,
        date TEXT,
        recipe_id INTEGER NOT NULL,
        FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE INDEX IF NOT EXISTS ix_main_ferm_sgs ON "main_ferm_sgs" (recipe_id, id);