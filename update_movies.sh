go run main.go update-movies movie_vault
conda run -n movies-app dvc add movie_vault data/movies.db
git add movie_vault.dvc data/movies.db.dvc
git commit -m "Update movies $(date)."
git push
conda run -n movies-app dvc push