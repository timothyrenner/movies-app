import dash
import dash_bootstrap_components as dbc
import dash_html_components as html
import dash_core_components as dcc

from tinydb import TinyDB, where
from toolz import get_in, pluck
from loguru import logger

# TODO: Fetch this from the data registry, eventually.
logger.info("Loading database.")
db = TinyDB("../data/processed/movie_database.json")

logger.info("Initializing min/max year.")
# Initialize the min/max years.
year_table = db.table("min_max_year")
min_max_years = year_table.all()
if len(min_max_years) != 1:
    raise ValueError(
        f"Expected min/max year to have 1 record: got {len(min_max_years)}."
    )
min_year = get_in([0, "min_year"], min_max_years)
max_year = get_in([0, "max_year"], min_max_years)

logger.info("Initializing genres.")
genres_table = db.table("genres")
genres = list(pluck("genre_name", genres_table.all()))

logger.info("Initializing tags.")
tags_table = db.table("tags")
tags = list(pluck("tag_name", tags_table.all()))

logger.info("Initializing services.")
services_table = db.table("services")
services = list(pluck("service_name", services_table.all()))

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
