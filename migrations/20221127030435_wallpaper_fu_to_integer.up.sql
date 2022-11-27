ALTER TABLE movie ADD COLUMN wallpaper_fu_int INTEGER NOT NULL DEFAULT 0;
UPDATE movie SET wallpaper_fu_int = CAST(wallpaper_fu AS INTEGER);
ALTER TABLE movie DROP COLUMN wallpaper_fu;
ALTER TABLE movie RENAME COLUMN wallpaper_fu_int TO wallpaper_fu;