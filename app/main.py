import dash
import dash_bootstrap_components as dbc
import dash_html_components as html
import dash_core_components as dcc
import plotly.graph_objs as go
import dvc.api
import os

from tinydb import TinyDB, where
from toolz import get_in, pluck, groupby, valmap
from loguru import logger
from dateutil.parser import parse
from dateutil.rrule import rrule, MONTHLY, WEEKLY
from dateutil.relativedelta import relativedelta
from datetime import datetime
from typing import List, Any, Dict, Tuple
from dash.dependencies import Input, Output
from random import sample
from google.cloud import storage
from dotenv import load_dotenv, find_dotenv

logger.info("Loading .env (if applicable).")
load_dotenv(find_dotenv())
MOVIE_ACCESS_TOKEN = os.getenv("MOVIE_ACCESS_TOKEN")

logger.info("Fetching database location from DVC.")
dataset_url = dvc.api.get_url(
    "data/processed/movie_database.json",
    repo=f"https://{MOVIE_ACCESS_TOKEN}@github.com/timothyrenner/movies-app",
)
# The url is a fully qualified URI for the dataset.
# e.g. gs://bucket/path/path/
# So we don't need to initialize the bucket with the client, we can just
# grab the file as long as we can write to it. We do need `wb` permission for
# some reason. 🤷‍♂️

logger.info("Fetching database from GCS.")
storage_client = storage.Client()
with open("movie_database.json", "wb") as f:
    storage_client.download_blob_to_file(dataset_url, f)

logger.info("Initializing database.")
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

no_margin = {"margin": 0}
plotly_margin = {"t": 50, "b": 50, "l": 0, "r": 0}
calendar_row = dbc.Row(
    [dbc.Col(dcc.Graph(id="calendar-graph", style=no_margin))]
)
# calendar_row = dcc.Graph(id="calendar-graph", style=no_margin)
year_row = dbc.Row(
    [
        dbc.Col(dcc.Graph(id="year-graph", style=no_margin)),
    ]
)
breakdown_row = dbc.Row(
    [
        dbc.Col(
            dcc.Graph(id="service-graph", style=no_margin),
        ),
        dbc.Col(
            dcc.Graph(id="genre-graph", style=no_margin),
        ),
    ]
)
histogram_row = dbc.Row(
    [
        dbc.Col(
            dcc.Graph(id="rt-histogram-graph", style=no_margin),
            md=6,
            lg=6,
            xl=6,
        ),
        # NOTE: maybe placeholder here.
        dbc.Col(
            dcc.Graph(id="imdb-histogram-graph", style=no_margin),
            md=6,
            lg=6,
            xl=6,
        ),
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

sidebar_inputs = [
    Input("watched-slider", "value"),
    Input("year-slider", "value"),
    Input("genre-dropdown", "value"),
    Input("service-dropdown", "value"),
]


@dash_app.callback(Output("service-graph", "figure"), sidebar_inputs)
def service_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services)
    movies_by_service: Dict[str, List[str]] = {}
    for movie in matching_movies:
        for service in movie["service"]:
            if service not in movies_by_service:
                movies_by_service[service] = []
            movies_by_service[service].append(movie["title"])
    movie_counts_by_service = valmap(len, movies_by_service)

    x: List[str] = []
    y: List[int] = []
    text: List[str] = []

    # Iterate over the dict sorted by keys, first to last.
    # We only need the count dict to sort the values.
    # See https://stackoverflow.com/a/3177911
    for service in sorted(  # type: ignore
        movie_counts_by_service,
        key=movie_counts_by_service.get,  # type: ignore
        reverse=True,
    ):
        service_movies = movies_by_service[service]
        x.append(service)
        y.append(len(service_movies))
        text.append(
            "<br>".join(
                sample(service_movies, 25)
                if len(service_movies) > 25
                else service_movies
            )
        )

    fig = go.Figure(
        data=[
            go.Bar(
                x=x, y=y, text=text, hovertemplate="%{text}<extra>%{x}</extra>"
            )
        ]
    )
    fig.layout = go.Layout(margin=plotly_margin, title="Services")

    return fig


@dash_app.callback(Output("genre-graph", "figure"), sidebar_inputs)
def genre_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services)

    movies_by_genre: Dict[str, List[str]] = {}
    for movie in matching_movies:
        for genre in movie["genre"]:
            if genre not in movies_by_genre:
                movies_by_genre[genre] = []
            movies_by_genre[genre].append(movie["title"])
    movie_counts_by_genre = valmap(len, movies_by_genre)

    x: List[str] = []
    y: List[int] = []
    text: List[str] = []

    # Iterate over the dict sorted by keys, first to last.
    # We only need the count dict to sort the values.
    # See https://stackoverflow.com/a/3177911
    for genre in sorted(  # type: ignore
        movie_counts_by_genre,
        key=movie_counts_by_genre.get,  # type: ignore
        reverse=True,
    ):
        genre_movies = movies_by_genre[genre]
        x.append(genre)
        y.append(len(genre_movies))
        text.append(
            "<br>".join(
                sample(genre_movies, 25)
                if len(genre_movies) > 25
                else genre_movies
            )
        )

    fig = go.Figure(
        data=[
            go.Bar(
                x=x, y=y, text=text, hovertemplate="%{text}<extra>%{x}</extra>"
            )
        ],
        layout=go.Layout(margin=plotly_margin, title="Genres"),
    )
    return fig


@dash_app.callback(Output("year-graph", "figure"), sidebar_inputs)
def year_graph(
    watched: List[int], year: List[int], genres: List[str], services: List[str]
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services)

    # Grab the "year" field and count.
    movie_year_grouped: Dict[int, List[Dict[str, Any]]] = groupby(
        "year", matching_movies
    )

    x: List[int] = []
    y: List[int] = []
    text: List[str] = []

    for movie_year in range(year[0], year[1] + 1):
        x.append(movie_year)
        if movie_year in movie_year_grouped:
            y.append(len(movie_year_grouped[movie_year]))
            text.append(
                # Line separate the titles. We need to pluck them out of the
                # list because groupby groups entire documents.
                "<br>".join(pluck("title", movie_year_grouped[movie_year]))
            )
        else:
            y.append(0)
            text.append("")

    fig = go.Figure(
        data=[
            go.Bar(
                x=x,
                y=y,
                text=text,
                hovertemplate="%{text}<extra><b>%{x}</b></extra>",
            )
        ],
        layout=go.Layout(title="Release Year", margin=plotly_margin),
    )
    return fig


@dash_app.callback(Output("calendar-graph", "figure"), sidebar_inputs)
def calendar_graph(
    watched: List[int], year: List[int], genres: List[str], services: List[str]
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services)

    watched_start = compute_month(watched[0])
    watched_end = compute_month(watched[1])

    days_of_week: List[str] = ["Sat", "Fri", "Thu", "Wed", "Tue", "Mon", "Sun"]
    weeks: List[str] = [
        w.strftime("%Y-%m-%d")
        for w in rrule(WEEKLY, dtstart=watched_start, until=watched_end)
    ]
    movie_counts_on_day: List[List[int]] = [
        [0 for ii in range(len(weeks))] for jj in range(len(days_of_week))
    ]
    movies_on_day: List[List[List[str]]] = [
        [[] for ii in range(len(weeks))] for jj in range(len(days_of_week))
    ]
    movie_days: List[List[str]] = [
        [
            (w + relativedelta(days=ii)).strftime("%Y-%m-%d")
            for w in rrule(WEEKLY, dtstart=watched_start, until=watched_end)
        ]
        for ii in range(len(days_of_week))
    ]

    for movie in matching_movies:
        movie_watched_date = parse(movie["watched"])
        movie_watched_year = movie_watched_date.year - watched_start.year
        movie_watched_day_of_week = 6 - int(movie_watched_date.strftime("%w"))
        # This arithmetic is: the week of the year, but with the zero point on
        # the earliest watched date.
        movie_watched_week_of_year = int(
            movie_watched_date.strftime("%W")
        ) - int(watched_start.strftime("%W"))

        movie_counts_on_day[movie_watched_day_of_week][
            # This arithmetic is: week of year + year offset, where year offset
            # is the number of years we have data for.
            movie_watched_week_of_year
            + (52 * movie_watched_year)
        ] += 1
        movies_on_day[movie_watched_day_of_week][
            # This arithmetic is: week of year + year offset, where year offset
            # is the number of years we have data for.
            movie_watched_week_of_year
            + (52 * movie_watched_year)
        ].append(movie["title"])

    # Now generate the text by combining the movie days with the movies
    # watched.
    custom_data: List[List[Tuple[str, str]]] = []
    for ii in range(len(days_of_week)):
        custom_data_row: List[Tuple[str, str]] = []
        for jj in range(len(weeks)):
            movies = movies_on_day[ii][jj]
            date = movie_days[ii][jj]
            if not movies:
                movie_str = " "
            else:
                movie_str = "<br>".join(movies)
            custom_data_row.append((date, movie_str))
        custom_data.append(custom_data_row)

    fig = go.Figure(
        go.Heatmap(
            z=movie_counts_on_day,
            y=days_of_week,
            x=weeks,
            customdata=custom_data,
            colorscale="greens",
            xgap=1,
            ygap=1,
            showscale=False,
            hovertemplate="%{customdata[1]}<extra>%{customdata[0]}</extra>",
        ),
        layout=go.Layout(title="Calendar", margin=plotly_margin),
    )
    return fig


@dash_app.callback(Output("rt-histogram-graph", "figure"), sidebar_inputs)
def rt_histogram_graph(
    watched: List[int], year: List[int], genres: List[str], services: List[str]
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services)

    bin_labels = [f"{ii*10}%" for ii in range(11)]
    movies_in_bin: List[List[str]] = [[] for _ in range(len(bin_labels))]
    rating_bin_counts: List[int] = [0 for _ in range(len(bin_labels))]
    for movie in matching_movies:
        for rating in movie["ratings"]:
            if rating["source"] == "Rotten Tomatoes":
                rating_numeric = int(rating["value"][:-1])
                rating_bin_index = rating_numeric // 10
                movies_in_bin[rating_bin_index].append(movie["title"])
                rating_bin_counts[rating_bin_index] += 1
    text = ["<br>".join(m) for m in movies_in_bin]
    fig = go.Figure(
        go.Bar(
            x=bin_labels,
            y=rating_bin_counts,
            text=text,
            hovertemplate="%{text}<extra>%{x}</extra>",
        ),
        go.Layout(
            title="Rotten Tomatoes",
            margin=plotly_margin,
        ),
    )
    return fig


@dash_app.callback(Output("imdb-histogram-graph", "figure"), sidebar_inputs)
def imdb_histogram_graph(
    watched: List[int], year: List[int], genres: List[str], services: List[str]
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services)

    bin_labels = [ii for ii in range(11)]
    movies_in_bin: List[List[str]] = [[] for _ in range(len(bin_labels))]
    rating_bin_counts: List[int] = [0 for _ in range(len(bin_labels))]
    for movie in matching_movies:
        for rating in movie["ratings"]:
            if rating["source"] == "Internet Movie Database":
                rating_bin_index = int(float(rating["value"].split("/")[0]))
                movies_in_bin[rating_bin_index].append(movie["title"])
                rating_bin_counts[rating_bin_index] += 1
    text = [
        "<br>".join(sample(m, 25) if len(m) >= 25 else m)
        for m in movies_in_bin
    ]
    fig = go.Figure(
        go.Bar(
            x=bin_labels,
            y=rating_bin_counts,
            text=text,
            hovertemplate="%{text}<extra>%{x}</extra>",
        ),
        go.Layout(
            title="IMDB",
            margin=plotly_margin,
        ),
    )
    return fig


if __name__ == "__main__":
    dash_app.run_server()
