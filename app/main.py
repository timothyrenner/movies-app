import dash
import dash_bootstrap_components as dbc
import dash_html_components as html
import dash_core_components as dcc
import plotly.graph_objs as go
import dvc.api
import os

from tinydb import TinyDB, where
from toolz import get_in, pluck
from loguru import logger
from dateutil.parser import parse
from dateutil.rrule import rrule, MONTHLY
from dateutil.relativedelta import relativedelta
from datetime import datetime
from typing import List, Any, Dict
from dash.dependencies import Input, Output, State
from google.cloud import storage, secretmanager
from dotenv import load_dotenv, find_dotenv
from plots import (
    imdb_histogram_plot,
    rt_histogram_plot,
    calendar_plot,
    year_plot,
    genre_plot,
    service_plot,
)

logger.info("Loading .env (if applicable).")
load_dotenv(find_dotenv())

try:
    logger.info("Fetching movie access token.")
    GOOGLE_CLOUD_PROJECT = os.getenv("GOOGLE_CLOUD_PROJECT")
    secret_client = secretmanager.SecretManagerServiceClient()
    secret_path = (
        f"projects/{GOOGLE_CLOUD_PROJECT}/secrets/"
        "MOVIE_ACCESS_TOKEN/versions/2"
    )
    secret = secret_client.access_secret_version(name=secret_path)
    MOVIE_ACCESS_TOKEN = secret.payload.data.decode("UTF-8")

    logger.info("Fetching database location from DVC.")
    dataset_url = dvc.api.get_url(
        "data/processed/movie_database.json",
        repo=(
            f"https://{MOVIE_ACCESS_TOKEN}@github.com/"
            "timothyrenner/movies-app"
        ),
    )
    # The url is a fully qualified URI for the dataset.
    # e.g. gs://bucket/path/path/
    # So we don't need to initialize the bucket with the client, we can just
    # grab the file as long as we can write to it. We do need `wb` permission
    # for some reason. 🤷‍♂️

    logger.info("Fetching database from GCS.")
    storage_client = storage.Client()
    with open("movie_database.json", "wb") as f:
        storage_client.download_blob_to_file(dataset_url, f)
    movie_database_path = "movie_database.json"
except Exception:
    logger.exception(
        "Some bullshit went wrong loading the DB. Trying local path."
    )
    movie_database_path = "../data/processed/movie_database.json"

logger.info("Initializing database.")
db = TinyDB(movie_database_path)

logger.info("Initializing min/max year.")
year_table = db.table("min_max_year")
min_max_years = year_table.all()
if len(min_max_years) != 1:
    raise ValueError(
        f"Expected min/max year to have 1 record: got {len(min_max_years)}."
    )
min_year = get_in([0, "min_year"], min_max_years)
max_year = get_in([0, "max_year"], min_max_years)
logger.debug(f"Min year: {min_year}.")
logger.debug(f"Max year: {max_year}.")

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
    services: List[str],
    tags: List[str] = [],
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

    query = (
        (where("watched") >= min_watched)  # type: ignore
        & (where("watched") <= max_watched)
        & (where("year") >= min_year)
        & (where("year") <= max_year)
        & (where("genre").any(genres))
        & (where("service").any(services))
    )

    if tags:
        query = query & where("tags").any(tags)

    return movies_table.search(query)  # type: ignore


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
                dbc.Label("Release Years"),
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
                dbc.Label("Services"),
                dcc.Dropdown(
                    id="service-dropdown",
                    options=[{"label": s, "value": s} for s in services],
                    value=services,
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
                    value=[],
                    multi=True,
                ),
            ]
        ),
    ],
    body=True,
)

no_margin = {"margin": 0}
plotly_margin = {"t": 50, "b": 0, "l": 0, "r": 0}
plotly_margin_notitle = {"t": 0, "b": 0, "l": 0, "r": 0}

calendar_row = dbc.Row(
    [
        dbc.Col(
            dbc.Card(
                [
                    dbc.CardHeader(
                        dbc.Button(
                            "Calendar",
                            id="calendar-button",
                            color="light",
                            block=True,
                        )
                    ),
                    dbc.Collapse(
                        dbc.CardBody(
                            dcc.Graph(id="calendar-graph", style=no_margin)
                        ),
                        is_open=True,
                        id="calendar-collapse",
                    ),
                ]
            )
        )
    ]
)
year_row = dbc.Row(
    [
        dbc.Col(
            dbc.Card(
                [
                    dbc.CardHeader(
                        dbc.Button(
                            "Release Years",
                            id="year-button",
                            color="light",
                            block=True,
                        ),
                    ),
                    dbc.Collapse(
                        dbc.CardBody(
                            dcc.Graph(id="year-graph", style=no_margin)
                        ),
                        is_open=True,
                        id="year-collapse",
                    ),
                ]
            )
        )
    ]
)
breakdown_row = dbc.Row(
    [
        dbc.Col(
            dbc.Card(
                [
                    dbc.CardHeader(
                        dbc.Button(
                            "Services & Genres",
                            id="services-genres-button",
                            color="light",
                            block=True,
                        ),
                    ),
                    dbc.Collapse(
                        dbc.CardBody(
                            dbc.Row(
                                [
                                    dbc.Col(
                                        dcc.Graph(
                                            id="service-graph", style=no_margin
                                        ),
                                        md=6,
                                        lg=6,
                                        xl=6,
                                    ),
                                    dbc.Col(
                                        dcc.Graph(
                                            id="genre-graph", style=no_margin
                                        ),
                                        md=6,
                                        lg=6,
                                        xl=6,
                                    ),
                                ]
                            )
                        ),
                        is_open=False,
                        id="services-genres-collapse",
                    ),
                ]
            )
        )
    ]
)
histogram_row = dbc.Row(
    [
        dbc.Col(
            dbc.Card(
                [
                    dbc.CardHeader(
                        dbc.Button(
                            "Ratings & Reviews",
                            id="ratings-reviews-button",
                            color="light",
                            block=True,
                        )
                    ),
                    dbc.Collapse(
                        dbc.CardBody(
                            dbc.Row(
                                [
                                    dbc.Col(
                                        dcc.Graph(
                                            id="rt-histogram-graph",
                                            style=no_margin,
                                        ),
                                        md=6,
                                        lg=6,
                                        xl=6,
                                    ),
                                    # NOTE: maybe placeholder here.
                                    dbc.Col(
                                        dcc.Graph(
                                            id="imdb-histogram-graph",
                                            style=no_margin,
                                        ),
                                        md=6,
                                        lg=6,
                                        xl=6,
                                    ),
                                ]
                            )
                        ),
                        is_open=False,
                        id="ratings-reviews-collapse",
                    ),
                ]
            )
        )
    ]
)

main_content = [calendar_row, year_row, breakdown_row, histogram_row]

dash_app.layout = dbc.Container(
    [
        html.H1("Movies"),
        html.Hr(),
        dbc.Row(
            [
                dbc.Col(sidebar, md=4, lg=3, xl=3),
                dbc.Col(main_content, md=8, lg=9, xl=9),
            ]
        ),
    ],
    fluid=True,
)


@dash_app.callback(
    Output("calendar-collapse", "is_open"),
    [Input("calendar-button", "n_clicks")],
    [State("calendar-collapse", "is_open")],
)
def calendar_collapse(n_clicks: int, is_open: bool) -> bool:
    if n_clicks:
        return not is_open
    return is_open


@dash_app.callback(
    Output("year-collapse", "is_open"),
    [Input("year-button", "n_clicks")],
    [State("year-collapse", "is_open")],
)
def year_collapse(n_clicks: int, is_open: bool) -> bool:
    if n_clicks:
        return not is_open
    return is_open


@dash_app.callback(
    Output("services-genres-collapse", "is_open"),
    [Input("services-genres-button", "n_clicks")],
    [State("services-genres-collapse", "is_open")],
)
def services_genres_collapse(n_clicks: int, is_open: bool) -> bool:
    if n_clicks:
        return not is_open
    return is_open


@dash_app.callback(
    Output("ratings-reviews-collapse", "is_open"),
    [Input("ratings-reviews-button", "n_clicks")],
    [State("ratings-reviews-collapse", "is_open")],
)
def ratings_reviews_collapse(n_clicks: int, is_open: bool) -> bool:
    if n_clicks:
        return not is_open
    return is_open


sidebar_inputs = [
    Input("watched-slider", "value"),
    Input("year-slider", "value"),
    Input("genre-dropdown", "value"),
    Input("service-dropdown", "value"),
    Input("tag-dropdown", "value"),
]


@dash_app.callback(Output("service-graph", "figure"), sidebar_inputs)
def service_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
    tags: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services, tags)

    return service_plot(matching_movies)


@dash_app.callback(Output("genre-graph", "figure"), sidebar_inputs)
def genre_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
    tags: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services, tags)

    return genre_plot(matching_movies)


@dash_app.callback(Output("year-graph", "figure"), sidebar_inputs)
def year_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
    tags: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services, tags)

    return year_plot(matching_movies, year[0], year[1])


@dash_app.callback(Output("calendar-graph", "figure"), sidebar_inputs)
def calendar_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
    tags: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services, tags)

    watched_start = compute_month(watched[0])
    watched_end = compute_month(watched[1])

    return calendar_plot(matching_movies, watched_start, watched_end)


@dash_app.callback(Output("rt-histogram-graph", "figure"), sidebar_inputs)
def rt_histogram_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
    tags: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services, tags)

    return rt_histogram_plot(matching_movies)


@dash_app.callback(Output("imdb-histogram-graph", "figure"), sidebar_inputs)
def imdb_histogram_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
    tags: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services, tags)

    return imdb_histogram_plot(matching_movies)


if __name__ == "__main__":
    dash_app.run_server()
