go run main.go update-movies
conda run movies-app dvc add movie_vault data/movies.db
git add movie_vault.dvc data/movie.db.dvc
git commit -m "Update movies $(date)."
git push