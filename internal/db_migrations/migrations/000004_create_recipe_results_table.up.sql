CREATE TABLE
    IF NOT EXISTS "recipe_results" (
        id INTEGER NOT NULL PRIMARY KEY,
        hot_wort_vol REAL,
        original_sg REAL,
        final_sg REAL,
        alcohol REAL,
        main_ferm_vol REAL,
        vol_bb REAL,
        recipe_id INTEGER NOT NULL,
        FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE ON UPDATE CASCADE
    );