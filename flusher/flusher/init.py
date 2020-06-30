import click
import json
from sqlalchemy import create_engine

from .cli import cli
from .db import metadata, tracking


@cli.command()
@click.argument("chain_id")
@click.argument("topic")
@click.option(
    "--db",
    help="Database URI connection string.",
    default="localhost:5432/postgres",
    show_default=True,
)
def init(chain_id, topic, db):
    """Initialize database with empty tables and tracking info."""
    engine = create_engine("postgresql+psycopg2://" + db, echo=True)
    metadata.create_all(engine)
    engine.execute(tracking.insert(), {"chain_id": chain_id, "topic": topic, "kafka_offset": -1})

