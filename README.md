# movies-app
Silly personal project for building an automated DVC / Dash / Google Cloud workflow for movies I've watched.

I watch a lot of movies. I like the idea of something like letterboxd but would rather DIY it because I have more time than sense, apparently. Anyway, that's what we have here - a pipeline that runs weekly and pulls my Airtable, enriches the movies with data from OMDB, then creates a small JSON database out of that and versions that with DVC.

The app is a Dash app that reads the latest version of that database file from the DVC remote + this github repo, and draws ugly charts.

Mostly I built this to develop an understanding of how to automate DVC-based workflows and use DVC to distribute data to endpoints and applications.
