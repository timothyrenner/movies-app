import typer
import json


from loguru import logger
from toolz import thread_last
from tinydb import TinyDB
from typing import Dict, Any, Set, List


def get_all_tags(movie_records: List[Dict[str, Any]]) -> Set[str]:
    tags: Set[str] = set()
    for movie in movie_records:
        tags.update(movie["tags"])
    return tags


def get_all_services(movie_records: List[Dict[str, Any]]) -> Set[str]:
    services: Set[str] = set()
    for movie in movie_records:
        services.update(movie["service"])
    return services


def get_all_genres(movie_records: List[Dict[str, Any]]) -> Set[str]:
    genres: Set[str] = set()
    for movie in movie_records:
        genres.update(movie["genre"])
    return genres


def main(
    movie_records_file: str = "data/interim/merged_records.json",
    output_file: str = "data/processed/movie_database.json",
):
    logger.info(f"Loading movie records from {movie_records_file}.")
    with open(movie_records_file, "r") as f:
        movie_records = thread_last(f, (map, json.loads), list)

    logger.info(f"Found {len(movie_records)} records.")

    db = TinyDB(output_file)
    movies_table = db.table("movies")
    tags_table = db.table("tags")
    services_table = db.table("services")
    genres_table = db.table("genres")

    logger.info("Extracting all tags.")
    all_tags = [{"tag_name": tn} for tn in get_all_tags(movie_records)]
    logger.info(f"Extracted {len(all_tags)} tags.")

    logger.info("Extracting all services.")
    all_services = [
        {"service_name": sn} for sn in get_all_services(movie_records)
    ]
    logger.info(f"Extracted {len(all_services)} services.")

    logger.info("Extracting all genres.")
    all_genres = [{"genre_name": gn} for gn in get_all_genres(movie_records)]
    logger.info(f"Extracted {len(all_genres)}.")

    logger.info("Inserting tags into database.")
    tags_table.insert_multiple(all_tags)
    logger.info("Done with tags.")

    logger.info("Inserting services into database.")
    services_table.insert_multiple(all_services)
    logger.info("Done with services.")

    logger.info("Inserting genres into database.")
    genres_table.insert_multiple(all_genres)
    logger.info("Done with genres.")

    logger.info("Inserting movies into database.")
    movies_table.insert_multiple(movie_records)
    logger.info("Done with movies.")

    logger.info("Closing database.")
    db.close()
    logger.info("All done!")


if __name__ == "__main__":
    typer.run(main)
