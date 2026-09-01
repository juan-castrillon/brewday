ALTER TABLE "recipes"
ADD COLUMN brewing_system TEXT NOT NULL DEFAULT 'undefined';

ALTER TABLE "stats"
ADD COLUMN brewing_system TEXT NOT NULL DEFAULT 'undefined';