import dash
import dash_bootstrap_components as dbc
import dash_html_components as html
import dash_core_components as dcc

from tinydb import TinyDB, where
from toolz import get_in, pluck
from loguru import logger
from dateutil.parser import parse
from dateutil.rrule import rrule, MONTHLY
from dateutil.relativedelta import relativedelta
from datetime import datetime
from typing import List, Any, Dict

# TODO: Fetch this from the data registry, eventually.
logger.info("Loading database.")
db = TinyDB("../data/processed/movie_database.json")

logger.info("Initializing min/max year.")
year_table = db.table("min_max_year")
min_max_years = year_table.all()
if len(min_max_years) != 1:
    raise ValueError(
        f"Expected min/max year to have 1 record: got {len(min_max_years)}."
    )
min_year = get_in([0, "min_year"], min_max_years)
max_year = get_in([0, "max_year"], min_max_years)

logger.info("Initializing min/max watched.")
watched_table = db.table("min_max_watched")
min_max_watched = watched_table.all()
if len(min_max_watched) != 1:
    raise ValueError(
        "Expected min/max watched to have 1 record: "
        f"got {len(min_max_watched)}."
    )
min_watched = parse(get_in([0, "min_watched"], min_max_watched))
max_watched = parse(get_in([0, "max_watched"], min_max_watched))
min_watched_month = datetime(min_watched.year, min_watched.month, 1)
max_watched_month = datetime(max_watched.year, max_watched.month, 1)
logger.debug(f"Min watched month: {min_watched_month.strftime('%m/%Y')}")
logger.debug(f"Max watched month: {max_watched_month.strftime('%m/%Y')}")


# Bind the minimum month to an easy converter, since we'll be doing this pretty
# much everywhere.
def compute_month(month_value: int) -> datetime:
    return min_watched_month + relativedelta(months=month_value)


month_max = len(
    list(rrule(MONTHLY, dtstart=min_watched_month, until=max_watched_month))
)
logger.info(f"Loaded {month_max} months of data.")


logger.info("Initializing genres.")
genres_table = db.table("genres")
genres = list(pluck("genre_name", genres_table.all()))

logger.info("Initializing tags.")
tags_table = db.table("tags")
tags = list(pluck("tag_name", tags_table.all()))

logger.info("Initializing services.")
services_table = db.table("services")
services = list(pluck("service_name", services_table.all()))


logger.info("Initializing movies.")
movies_table = db.table("movies")


def get_data(
    watched: List[int],
    year: List[int],
    genres: List[str],
    tags: List[str],
    services: List[str],
) -> List[Dict[str, Any]]:
    if len(watched) != 2:
        raise ValueError(
            f"Expected watched to have 2 values: got {len(watched)}."
        )
    min_watched_int, max_watched_int = watched
    min_watched = compute_month(min_watched_int).strftime("%Y-%m-%d")
    max_watched = compute_month(max_watched_int).strftime("%Y-%m-%d")

    if len(year) != 2:
        raise ValueError(f"Expected year to have 2 values: got {len(year)}.")
    min_year, max_year = year

    return movies_table.search(
        (where("watched") >= min_watched)  # type: ignore
        & (where("watched") <= max_watched)
        & (where("year") >= min_year)
        & (where("year") <= max_year)
        & (where("genre").any(genres))
        & (where("tags").any(tags))
        & (where("service").any(services))
    )


external_stylesheets = [dbc.themes.BOOTSTRAP]
dash_app = dash.Dash(external_stylesheets=external_stylesheets)
dash_app.title = "Movies"
# This is for gunicorn to hook into.
app = dash_app.server

sidebar = dbc.Card(
    [
        html.Br(),
        html.Br(),
        html.Br(),
        dbc.FormGroup(
            [
                dbc.Label("Watched"),
                dcc.RangeSlider(
                    id="watched-slider",
                    min=0,
                    max=month_max,
                    step=1,
                    value=[0, month_max],
                    marks={
                        m: {
                            "label": compute_month(m).strftime("%m/%Y"),
                            "style": {
                                "transform": "rotate(55deg)",
                                "font-size": "8px",
                                "margin-top": "1px",
                            },
                        }
                        for m in range(0, month_max, 1)
                    },
                ),
            ]
        ),
        dbc.FormGroup(
            [
                dbc.Label("Years"),
                dcc.RangeSlider(
                    id="year-slider",
                    min=min_year,
                    max=max_year,
                    step=1,
                    value=[min_year, max_year],
                    marks={
                        y: {
                            "label": str(y),
                            "style": {
                                "transform": "rotate(55deg)",
                                "font-size": "8px",
                                "margin-top": "1px",
                            },
                        }
                        for y in range(min_year, max_year, 5)
                    },
                ),
            ]
        ),
        dbc.FormGroup(
            [
                dbc.Label("Genres"),
                dcc.Dropdown(
                    id="genre-dropdown",
                    options=[{"label": g, "value": g} for g in genres],
                    value=genres,
                    multi=True,
                ),
            ]
        ),
        dbc.FormGroup(
            [
                dbc.Label("Tags"),
                dcc.Dropdown(
                    id="tag-dropdown",
                    options=[{"label": t, "value": t} for t in tags],
                    value=tags,
                    multi=True,
                ),
            ]
        ),
        dbc.FormGroup(
            [
                dbc.Label("Services"),
                dcc.Dropdown(
                    id="service-dropdown",
                    options=[{"label": s, "value": s} for s in services],
                    value=services,
                    multi=True,
                ),
            ]
        ),
    ],
    body=True,
)

dash_app.layout = dbc.Container(
    [
        html.H1("Movies"),
        html.Hr(),
        dbc.Row(
            [
                dbc.Col(sidebar, md=4, lg=3, xl=3),
            ]
        ),
    ],
    fluid=True,
)


if __name__ == "__main__":
    dash_app.run_server()
