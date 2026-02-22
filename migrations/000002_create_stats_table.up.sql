# TODO: Check if stats is an issue
CREATE TABLE 
    IF NOT EXISTS stats (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        recipe_title TEXT UNIQUE,
        finished_epoch INTEGER,
        evaporation REAL,
        efficiency REAL
    )