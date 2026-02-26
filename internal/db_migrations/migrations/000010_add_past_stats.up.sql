INSERT INTO
    "stats" (recipe_title, finished_epoch, evaporation, efficiency)
VALUES
    /*
    <!--Sudhaus:  48 + 57 +62 +72 + 71.4 + 72.1 + 79.9/ 7 -->
    <!-- Evaporacion: 16.66 + 19.48 + 10.45 + 25.19 + 28.79 + 16.4 +23 + 33.05 + 29 / 9 -->

    */
    (?, ?, ?, ?) -- Comment
ON CONFLICT DO UPDATE SET
    finished_epoch = excluded.finished_epoch
    evaporation = excluded.evaporation
    efficiency = excluded.efficiency
