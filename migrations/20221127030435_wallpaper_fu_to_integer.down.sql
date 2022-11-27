ALTER TABLE movie ADD COLUMN wallpaper_fu_bool BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE movie SET wallpaper_fu_bool = CAST(wallpaper_fu AS BOOLEAN);
ALTER TABLE movie DROP COLUMN wallpaper_fu;
ALTER TABLE movie RENAME COLUMN wallpaper_fu_bool TO wallpaper_fu;