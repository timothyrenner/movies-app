import os
import requests
import typer
import json

from dotenv import load_dotenv, find_dotenv
from loguru import logger
from typing import Any, List, Dict
from toolz import get
from dataclasses import dataclass, asdict

logger.info("Loading dotenv file.")
load_dotenv(find_dotenv())
logger.info("Dotenv file loaded.")

BASE_ID = os.getenv("AIRTABLE_MOVIE_BASE_ID")
API_KEY = os.getenv("AIRTABLE_KEY")
AIRTABLE_ENDPOINT = "https://api.airtable.com/v0"


@dataclass
class AirtableRecord:
    name: str
    rating: int
    imdb_link: str
    watched: str
    tags: List[str]
    service: List[str]


def make_record(dict_record: Dict[str, Any]) -> AirtableRecord:
    record_fields = dict_record["fields"]

    return AirtableRecord(
        name=get("Name", record_fields),
        rating=get("Rating", record_fields),
        imdb_link=get("IMDB Link", record_fields, None),
        watched=get("Watched", record_fields),
        tags=get("Tags", record_fields, []),
        service=get("Service", record_fields),
    )


def airtable_get(
    url: str, session: requests.Session, offset: int = None
) -> requests.Response:
    if offset:
        return session.get(
            url, params={"view": "Have Watched", "offset": offset}
        )
    else:
        return session.get(url, params={"view": "Have Watched"})


def airtable_get_all(
    url: str, session: requests.Session
) -> List[AirtableRecord]:
    records = []

    # Make the initial request.
    logger.info("Making a call to Airtable API, no offset.")
    airtable_json = airtable_get(
        f"{AIRTABLE_ENDPOINT}/{BASE_ID}/Movies", session
    ).json()
    records.extend(get("records", airtable_json, []))

    while offset := get("offset", airtable_json, None):
        logger.info(f"Making a call to Airtable API, offset: {offset}.")
        airtable_json = airtable_get(
            f"{AIRTABLE_ENDPOINT}/{BASE_ID}/Movies", session, offset=offset
        ).json()
        records.extend(get("records", airtable_json, []))

    return list(map(make_record, records))


def main(output_file: str = "data/raw/airtable_out.json"):
    # Validate environment.
    if not API_KEY:
        raise ValueError("AIRTABLE_KEY not in .env file or environment.")
    if not BASE_ID:
        raise ValueError(
            "AIRTABLE_MOVIE_BASE_ID not in .env file or environment."
        )
    request_url = f"{AIRTABLE_ENDPOINT}/{BASE_ID}/Movies"
    session = requests.Session()
    session.headers.update({"Authorization": f"Bearer {API_KEY}"})
    logger.info("Pulling records from Airtable.")
    records = airtable_get_all(request_url, session)
    logger.info("Done. Writing to file.")
    with open(output_file, "w") as f:
        for record in map(asdict, records):
            f.write(json.dumps(record))
            f.write("\n")
    logger.info("All done! 🎥 ")


if __name__ == "__main__":
    typer.run(main)
