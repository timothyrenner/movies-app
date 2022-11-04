conda run --no-capture-output -n movies-app dvc add movie_vault data/movies.db
git add movie_vault.dvc data/movies.db.dvc
git commit -m "Update movies $(date)."
git push
conda run --no-capture-output -n movies-app dvc push