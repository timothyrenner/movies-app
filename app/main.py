import dash
import dash_bootstrap_components as dbc
import dash_html_components as html
import dash_core_components as dcc
import plotly.graph_objs as go

from tinydb import TinyDB, where
from toolz import get_in, pluck, groupby
from loguru import logger
from dateutil.parser import parse
from dateutil.rrule import rrule, MONTHLY, WEEKLY
from dateutil.relativedelta import relativedelta
from datetime import datetime
from typing import List, Any, Dict, Tuple
from dash.dependencies import Input, Output
from itertools import chain
from collections import Counter

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
            md=4,
            lg=4,
            xl=4,
        ),
        # NOTE: maybe placeholder here.
        dbc.Col(
            dcc.Graph(id="imdb-histogram-graph", style=no_margin),
            md=4,
            lg=4,
            xl=4,
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

    # Grab the "services" field, flatten the list of lists, then count them.
    # This avoids creating a pandas data frame only to pass it to plotly
    # express, which would then promptly deconstruct it.
    movie_service_counts = Counter(
        chain.from_iterable(pluck("service", matching_movies))
    )
    x: List[str] = []
    y: List[int] = []

    logger.debug(f"Found {len(movie_service_counts)} counts.")

    for service, count in movie_service_counts.most_common(None):
        x.append(service)
        y.append(count)

    fig = go.Figure(data=[go.Bar(x=x, y=y)])
    fig.layout = go.Layout(margin={"t": 0, "b": 0, "l": 0, "r": 0})

    return fig


@dash_app.callback(Output("genre-graph", "figure"), sidebar_inputs)
def genre_graph(
    watched: List[int],
    year: List[int],
    genres: List[str],
    services: List[str],
) -> go.Figure:
    matching_movies = get_data(watched, year, genres, services)

    # Grab the "genres" field, flatten the list of lists, then count them.
    movie_genre_counts = Counter(
        chain.from_iterable(pluck("genre", matching_movies))
    )
    x: List[str] = []
    y: List[int] = []

    logger.debug(f"Found {len(movie_genre_counts)} counts.")

    for genre, count in movie_genre_counts.most_common(None):
        x.append(genre)
        y.append(count)

    fig = go.Figure(
        data=[go.Bar(x=x, y=y)],
        layout=go.Layout(margin={"t": 0, "b": 0, "l": 0, "r": 0}),
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
        layout=go.Layout(margin={"t": 0, "b": 0, "l": 0, "r": 0}),
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
        layout=go.Layout(margin={"t": 0, "b": 0, "l": 0, "r": 0}),
    )
    return fig


if __name__ == "__main__":
    dash_app.run_server()
