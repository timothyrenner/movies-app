import re
import typer
import json

from dateutil.parser import parse as parse_date
from dateutil.parser._parser import ParserError
from pytimeparse import parse as parse_time
from toolz import thread_last
from typing import List, Optional
from enum import Enum
from loguru import logger
from dataclasses import dataclass, asdict
from pull_notion import MovieRecord as NotionMovieRecord
from pull_omdb import OMDBRecord, get_imdb_id


class WriterRole(str, Enum):
    screenplay = "screenplay"
    story = "story"
    characters = "characters"
    play = "play"
    material = "material"
    dialogue = "dialogue"
    segment = "segment"
    comic = "comic"
    novel = "novel"
    concept = "concept"
    adaptation = "adaptation"


@dataclass
class Writer:
    name: str
    role: Optional[WriterRole] = None


@dataclass
class Rating:
    source: str
    value: str


@dataclass
class MovieRecord:
    title: str
    year: int
    runtime_minutes: Optional[int]
    release_date: Optional[str]
    genre: List[str]
    country: str
    director: List[str]
    actors: List[str]
    language: List[str]
    production: str
    writer: List[Writer]
    tags: List[str]
    watched: str
    ratings: List[Rating]
    service: str
    imdb_link: Optional[str]
    imdb_id: str


def convert_date(date_str: str) -> Optional[str]:
    try:
        return parse_date(date_str).strftime("%Y-%m-%d")
    except ParserError:
        logger.warning(f"Couldn't parse {date_str}.")
        return None


def convert_time(time_str: str) -> Optional[int]:
    parsed_time = parse_time(time_str)
    if parsed_time:
        return int(parse_time(time_str) / 60)
    else:
        return None


def split_commas(comma_sep_str: str) -> List[str]:
    return thread_last(
        comma_sep_str, lambda x: x.split(","), (map, lambda x: x.strip()), list
    )


def match_writer_role(writer_role_str: str) -> Optional[WriterRole]:
    if "screenplay" in writer_role_str:
        return WriterRole.screenplay
    elif "story" in writer_role_str:
        return WriterRole.story
    elif "characters" in writer_role_str:
        return WriterRole.characters
    elif "play" in writer_role_str:
        return WriterRole.play
    elif "material" in writer_role_str:
        return WriterRole.material
    elif "dialogue" in writer_role_str:
        return WriterRole.dialogue
    elif "segment" in writer_role_str:
        return WriterRole.segment
    elif "novel" in writer_role_str:
        return WriterRole.novel
    elif "comic" in writer_role_str:
        return WriterRole.comic
    elif "concept" in writer_role_str:
        return WriterRole.concept
    elif "adaptation" in writer_role_str:
        return WriterRole.adaptation
    else:
        logger.warning(f"Unable to match a writer role to {writer_role_str}.")
        return None


def get_writer(writer_str: str) -> Writer:
    if writer_match := re.match(r"(.*) \((.*)\)", writer_str):
        return Writer(
            name=writer_match.group(1).strip(),
            role=match_writer_role(writer_match.group(2).strip()),
        )
    else:
        return Writer(name=writer_str.strip())


def merge_airtable_omdb(
    notion_movie_record: NotionMovieRecord, omdb_record: OMDBRecord
) -> MovieRecord:
    return MovieRecord(
        title=omdb_record.title,
        year=int(omdb_record.year),
        release_date=convert_date(omdb_record.released),
        runtime_minutes=convert_time(omdb_record.runtime),
        genre=split_commas(omdb_record.genre),
        country=omdb_record.country,
        director=split_commas(omdb_record.director),
        actors=split_commas(omdb_record.actors),
        language=split_commas(omdb_record.language),
        production=omdb_record.production,
        writer=thread_last(
            omdb_record.writer, split_commas, (map, get_writer), list
        ),
        tags=notion_movie_record.tags,
        watched=notion_movie_record.watched,
        ratings=[
            Rating(source=x["Source"], value=x["Value"])
            for x in omdb_record.ratings
        ],
        service=notion_movie_record.service,
        imdb_link=notion_movie_record.imdb_link,
        imdb_id=omdb_record.imdb_id,
    )


def main(
    movie_file: str = "data/raw/notion_out.json",
    omdb_file: str = "data/raw/omdb_out.json",
    output_file: str = "data/interim/merged_records.json",
):
    logger.info(f"Reading movie file {movie_file}.")
    with open(movie_file, "r") as f:
        notion_movie_records = thread_last(
            f, (map, json.loads), (map, lambda x: NotionMovieRecord(**x)), list
        )
    logger.info(f"Read {len(notion_movie_records)} from {movie_file}.")

    logger.info(f"Reading OMDB file {omdb_file}.")
    with open(omdb_file, "r") as f:
        omdb_records = {
            r.imdb_id: r
            for r in thread_last(
                f, (map, json.loads), (map, lambda x: OMDBRecord(**x))
            )
        }
    logger.info(f"Read {len(omdb_records)} from {omdb_file}.")

    logger.info("Merging records.")
    movie_records: List[MovieRecord] = []
    for notion_movie_record in notion_movie_records:
        imdb_id = get_imdb_id(notion_movie_record.imdb_link)
        omdb_record = omdb_records[imdb_id]

        movie_record = merge_airtable_omdb(notion_movie_record, omdb_record)
        movie_records.append(movie_record)
    logger.info(f"Created {len(movie_records)} merged records.")

    logger.info(f"Writing merged records to {output_file}.")
    with open(output_file, "w") as f:
        for movie_record in movie_records:
            thread_last(movie_record, asdict, json.dumps, f.write)
            f.write("\n")
    logger.info("Done!")


if __name__ == "__main__":
    typer.run(main)
