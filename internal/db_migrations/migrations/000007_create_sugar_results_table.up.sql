CREATE TABLE
    IF NOT EXISTS "sugar_results" (
        id INTEGER NOT NULL PRIMARY KEY,
        water REAL,
        sugar REAL,
        alcohol REAL,
        recipe_id INTEGER NOT NULL,
        FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE INDEX IF NOT EXISTS ix_sugar_results ON "sugar_results" (recipe_id, id);