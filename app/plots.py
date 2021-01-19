import plotly.graph_objs as go

from typing import List, Dict, Any, Tuple
from random import sample
from datetime import datetime
from dateutil.relativedelta import (
    relativedelta,
    SU,
    MO,
    TU,
    WE,
    TH,
    FR,
    SA,
)
from dateutil._common import weekday
from dateutil.rrule import rrule, WEEKLY
from dateutil.parser import parse
from toolz import groupby, pluck, valmap

PLOTLY_MARGIN = {"t": 50, "b": 0, "l": 0, "r": 0}
PLOTLY_MARGIN_NOTITLE = {"t": 0, "b": 0, "l": 0, "r": 0}


def service_plot(movies: List[Dict[str, Any]]) -> go.Figure:
    movies_by_service: Dict[str, List[str]] = {}
    for movie in movies:
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
    fig.layout = go.Layout(margin=PLOTLY_MARGIN, title="Services")

    return fig


def genre_plot(movies: List[Dict[str, Any]]) -> go.Figure:
    movies_by_genre: Dict[str, List[str]] = {}
    for movie in movies:
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
        layout=go.Layout(margin=PLOTLY_MARGIN, title="Genres"),
    )
    return fig


def year_plot(
    movies: List[Dict[str, Any]], start_year: int, end_year: int
) -> go.Figure:
    # Grab the "year" field and count.
    movie_year_grouped: Dict[int, List[Dict[str, Any]]] = groupby(
        "year", movies
    )

    x: List[int] = []
    y: List[int] = []
    text: List[str] = []

    for movie_year in range(start_year, end_year + 1):
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
        layout=go.Layout(margin=PLOTLY_MARGIN_NOTITLE),
    )
    return fig


def calendar_plot(
    movies: List[Dict[str, Any]],
    watched_start: datetime,
    watched_end: datetime,
) -> go.Figure:
    # Find the nearest Saturday to start.
    movie_days_start = watched_start + relativedelta(weekday=SU(-1))
    movie_days_end = watched_end + relativedelta(weekday=SU(1))

    days_of_week: List[weekday] = [SA, FR, TH, WE, TU, MO, SU]
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
            (w + relativedelta(weekday=wd)).strftime("%Y-%m-%d")
            for w in rrule(
                WEEKLY, dtstart=movie_days_start, until=movie_days_end
            )
        ]
        for wd in days_of_week
    ]

    for movie in movies:
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
            movies_today = movies_on_day[ii][jj]
            date = movie_days[ii][jj]
            if not movies:
                movie_str = " "
            else:
                movie_str = "<br>".join(movies_today)
            custom_data_row.append((date, movie_str))
        custom_data.append(custom_data_row)

    fig = go.Figure(
        go.Heatmap(
            z=movie_counts_on_day,
            y=["Sat", "Fri", "Thu", "Wed", "Tue", "Mon", "Sun"],
            x=weeks,
            customdata=custom_data,
            colorscale="greens",
            xgap=1,
            ygap=1,
            showscale=False,
            hovertemplate="%{customdata[1]}<extra>%{customdata[0]}</extra>",
        ),
        layout=go.Layout(margin=PLOTLY_MARGIN_NOTITLE),
    )
    return fig


def rt_histogram_plot(movies: List[Dict[str, Any]]) -> go.Figure:
    bin_labels = [f"{ii*10}%" for ii in range(11)]
    movies_in_bin: List[List[str]] = [[] for _ in range(len(bin_labels))]
    rating_bin_counts: List[int] = [0 for _ in range(len(bin_labels))]
    for movie in movies:
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
            margin=PLOTLY_MARGIN,
        ),
    )
    return fig


def imdb_histogram_plot(movies: List[Dict[str, Any]]) -> go.Figure:
    bin_labels = [ii for ii in range(11)]
    movies_in_bin: List[List[str]] = [[] for _ in range(len(bin_labels))]
    rating_bin_counts: List[int] = [0 for _ in range(len(bin_labels))]
    for movie in movies:
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
            margin=PLOTLY_MARGIN,
        ),
    )
    return fig
