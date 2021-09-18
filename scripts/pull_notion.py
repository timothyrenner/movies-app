import typer
import os
import sys
import requests
import json

from dotenv import load_dotenv, find_dotenv
from loguru import logger
from toolz import get, pluck, compose, get_in, thread_last
from typing import List, Dict, Any, Optional
from dataclasses import dataclass, asdict

listpluck = compose(list, pluck)


@dataclass
class MovieRecord:
    name: str
    first_time: bool
    watched: str
    service: str
    imdb_link: Optional[str]
    tags: List[str]


def make_record(dict_record: Dict[str, Any]) -> MovieRecord:
    name = get_in(
        ["properties", "Name", "title", 0, "plain_text"],
        dict_record,
        no_default=True,
    )
    first_time = get_in(
        ["properties", "First Time", "checkbox"], dict_record, no_default=True
    )
    watched = get_in(
        ["properties", "Watched", "date", "start"],
        dict_record,
        no_default=True,
    )
    service = get_in(
        ["properties", "Service_str", "formula", "string"],
        dict_record,
        no_default=True,
    )
    imdb_link = get_in(
        ["properties", "IMDB Link", "url"], dict_record, default=None
    )
    tags = listpluck(
        "name",
        get_in(
            ["properties", "Tags", "multi_select"], dict_record, default=[]
        ),
    )
    return MovieRecord(
        name=name,
        first_time=first_time,
        watched=watched,
        service=service,
        imdb_link=imdb_link,
        tags=tags,
    )


logger.info("Loading .env file.")
load_dotenv(find_dotenv())
NOTION_KEY = os.getenv("NOTION_KEY")
if not NOTION_KEY:
    logger.error(
        "No entry for NOTION_KEY in .env or environment. Terminating."
    )
    sys.exit(1)
NOTION_URL = "https://api.notion.com"


def get_database_id(session: requests.Session, database_name: str) -> str:
    request_body = {
        "query": database_name,
        "filter": {"property": "object", "value": "database"},
    }

    response = session.post(f"{NOTION_URL}/v1/search", json=request_body)
    if not response.ok:
        logger.error(f"Search request for database {database_name} failed.")
        logger.error(response.json())
        sys.exit(1)

    response_json = response.json()
    if len(response_json["results"]) > 1:
        logger.warning(
            f"Found more than one match for {database_name}. Refine the query."
        )
    return response_json["results"][0]["id"]


def fetch_data(
    session: requests.Session, database_id: str
) -> List[Dict[str, Any]]:
    records: List[Dict[str, Any]] = []

    response = session.post(
        f"{NOTION_URL}/v1/databases/{database_id}/query",
        json={"page_size": 100},
    )
    if not response.ok:
        logger.error("Query to the database failed.")
        logger.error(response.json())
        sys.exit(1)

    response_json = response.json()
    records.extend(response_json["results"])
    while next_cursor := get("next_cursor", response_json, None):

        logger.info("Making a call to Notion API.")
        response = session.post(
            f"{NOTION_URL}/v1/databases/{database_id}/query",
            json={"page_size": 100, "start_cursor": next_cursor},
        )
        if not response.ok:
            logger.error("Query to the database failed.")
            logger.error(response.json())
            sys.exit(1)
        response_json = response.json()

        records.extend(response_json["results"])

    return records


def main(output_file: str = "data/raw/notion_out.json"):
    logger.info("Initializing requests session.")
    session = requests.Session()
    session.headers.update(
        {
            "Authorization": f"Bearer {NOTION_KEY}",
            "Content-Type": "application/json",
            "Notion-Version": "2021-08-16",
        }
    )
    logger.info("Querying movie list database.")
    database_id = get_database_id(session, "Movie List")

    logger.info("Got database id.")
    logger.info("Fetching data from database.")
    database_data_raw = fetch_data(session, database_id)
    logger.info(f"Retrieved {len(database_data_raw)} records from database.")
    with open(output_file, "w") as f:
        for record in thread_last(
            database_data_raw, (map, make_record), (map, asdict)
        ):
            f.write(json.dumps(record))
            f.write("\n")
    logger.info("All done! 🎥 ")


if __name__ == "__main__":
    typer.run(main)
