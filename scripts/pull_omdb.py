import os
import json
import requests
import typer
import time

from dotenv import load_dotenv, find_dotenv
from loguru import logger
from typing import List, Set
from dataclasses import dataclass, asdict
from toolz import thread_first, thread_last, get
from typing import Dict, Any

from pull_airtable import AirtableRecord


@dataclass
class OMDBRecord:
    title: str
    year: str
    rated: str
    released: str
    runtime: str
    genre: str
    director: str
    writer: str
    actors: str
    plot: str
    language: str
    country: str
    awards: str
    poster: str
    ratings: List[Dict[str, str]]
    metascore: str
    imdb_rating: str
    imdb_votes: str
    imdb_id: str
    type: str
    dvd: str
    box_office: str
    production: str
    website: str


logger.info("Loading .env file.")
load_dotenv(find_dotenv())

API_KEY = os.environ.get("OMDB_KEY")
OMBD_ENDPOINT = "http://omdbapi.com"


def load_cache(file: str) -> List[str]:
    with open(file, "r") as f:
        return f.readlines()


def hydrate_airtable_record(dict_record: Dict[str, Any]) -> AirtableRecord:
    return AirtableRecord(**dict_record)


def get_imdb_id(imdb_url: str) -> str:
    return thread_last(imdb_url, lambda x: x.split("/"), (get, -2))


def make_omdb_record(omdb_record_dict: Dict[str, Any]) -> OMDBRecord:
    return OMDBRecord(
        title=omdb_record_dict["Title"],
        year=omdb_record_dict["Year"],
        rated=omdb_record_dict["Rated"],
        released=omdb_record_dict["Released"],
        runtime=omdb_record_dict["Runtime"],
        genre=omdb_record_dict["Genre"],
        director=omdb_record_dict["Director"],
        writer=omdb_record_dict["Writer"],
        actors=omdb_record_dict["Actors"],
        plot=omdb_record_dict["Plot"],
        language=omdb_record_dict["Language"],
        country=omdb_record_dict["Country"],
        awards=omdb_record_dict["Awards"],
        poster=omdb_record_dict["Poster"],
        ratings=omdb_record_dict["Ratings"],
        metascore=omdb_record_dict["Metascore"],
        imdb_rating=omdb_record_dict["imdbRating"],
        imdb_votes=omdb_record_dict["imdbVotes"],
        imdb_id=omdb_record_dict["imdbID"],
        type=omdb_record_dict["Type"],
        dvd=omdb_record_dict["DVD"],
        box_office=omdb_record_dict["BoxOffice"],
        production=omdb_record_dict["Production"],
        website=omdb_record_dict["Website"],
    )


def hydrate_omdb_record(omdb_record_dict: Dict[str, Any]) -> OMDBRecord:
    return OMDBRecord(**omdb_record_dict)


def load_omdb_records(file: str) -> List[OMDBRecord]:
    with open(file, "r") as f:
        entries = thread_last(
            f, (map, json.loads), (map, hydrate_omdb_record), list
        )
    return entries


def get_omdb(movie_id: str, session: requests.Session) -> OMDBRecord:
    raw_results = session.get(
        OMBD_ENDPOINT, params={"apikey": API_KEY, "i": movie_id}
    ).json()
    try:
        omdb_record = make_omdb_record(raw_results)
    except Exception as e:
        logger.exception(
            f"Encountered exception for {json.dumps(raw_results)}."
        )
        raise e
    return omdb_record


def main(
    airtable_file: str = "data/raw/airtable_out.json",
    output_file: str = "data/raw/omdb_out.json",
    skip_cache: bool = False,
):
    imdb_ids: Set[str] = set()
    omdb_records: List[OMDBRecord] = []
    if skip_cache:
        logger.info("Skipping cache.")
    elif not os.path.exists(output_file):
        logger.info(f"{output_file} does not exist. Skipping cache.")
    else:
        logger.info("Loading cache.")
        omdb_records = load_omdb_records(output_file)
        imdb_ids = {r.imdb_id for r in omdb_records}
        logger.info(f"{len(imdb_ids)} ids loaded from the cache.")

    session = requests.Session()
    logger.info("Loading airtable entries and fetching OMDB data.")
    new_records = 0
    with open(airtable_file, "r") as f:
        airtable_entries = thread_last(
            f, (map, json.loads), (map, hydrate_airtable_record), list
        )
    logger.info(f"Loaded {len(airtable_entries)}.")
    for airtable_entry in airtable_entries:
        if not airtable_entry.imdb_link:
            logger.warning(
                f"{airtable_entry.name} does not have an IMDB link. Skipping."
            )
            continue

        imdb_id = get_imdb_id(airtable_entry.imdb_link)

        if imdb_id in imdb_ids:
            logger.info(
                f"Found {airtable_entry.name} ({imdb_id}) in cache. Skipping."
            )
            continue

        logger.info(
            f"Pulling OMDB data for {airtable_entry.name} ({imdb_id})."
        )
        try:
            omdb_record = get_omdb(imdb_id, session)
            omdb_records.append(omdb_record)
            new_records += 1
            time.sleep(0.1)
        except requests.RequestException:
            logger.exception(
                f"Encountered exception retrieving data for {imdb_id}. "
                "Skipping."
            )
            continue

    logger.info(
        f"Pulled {new_records} new records. "
        f"Writing {len(omdb_records)} to {output_file}."
    )
    with open(output_file, "w") as f:
        for omdb_record in omdb_records:
            thread_first(omdb_record, asdict, json.dumps, f.write)
            f.write("\n")
    logger.info("Done.")


if __name__ == "__main__":
    typer.run(main)
