INSERT INTO
    "stats" (recipe_title, finished_epoch, evaporation, efficiency)
VALUES
    ('U2FtaWNobGF1cyBDaHJpc3RtYXMgU3RvdXQ=', 1762988400, 29, 79.9) -- Samichlaus Christmas Stout
    ('TmVwYWw=', 1753999200, 33.05, 72.1) -- Nepal
    ('Qm9jYWRpbGxvIFRyaXBlbA==', 1723759200, 23, 71.4) -- Bocadillo Tripel
    ('U2Fib3IgVHJvcGljYWw=', 1719180000, 16.4, 72) -- Sabor Tropical
    ('U2Nod2lnaSBCaWVy', 1710543600, 28.79, 62) -- Schwigi Bier
    ('TGEgTW9uYQ==', 1707865200, 25.19, 0) -- La Mona
    ('QW1hcmlsbG8gV2VpemVu', 1702249200, 10.45, 57.23) -- Amarillo Weizen
    ('VmFuaWxsYSBNaWxrIFN0b3V0', 1698184800, 19.48, 48) -- Vanilla Milk Stout
    ('SGVybWFubyBKdWFu', 1693260000, 16.66, 0) -- Hermano Juan
ON CONFLICT DO UPDATE SET
    finished_epoch = excluded.finished_epoch
    evaporation = excluded.evaporation
    efficiency = excluded.efficiency;
