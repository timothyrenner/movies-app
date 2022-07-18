# movies-app
Silly personal project for building an automated DVC / SQLite / Google Cloud workflow for movies I've watched.

I watch a lot of movies. I like the idea of something like letterboxd but would rather DIY it because I have more time than sense, apparently. Anyway, that's what we have here - some code that will read a specially formatted page in Obsidian, enrich it with data from OMDB, then save all that to a SQLite database (as well as create a linked page in Obsidian for the movie if one doesn't already exist).

Mostly I built this to develop an understanding of how to automate DVC-based workflows and use DVC to distribute data to endpoints and applications. At some point that turned into me wanting to rewrite it in Go cause I got tired of Python.
