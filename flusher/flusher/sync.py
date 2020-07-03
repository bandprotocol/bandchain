import json
import click
import sys
from kafka import KafkaConsumer, TopicPartition
from loguru import logger
from sqlalchemy import create_engine

from .cli import cli
from .db import tracking
from .handler import Handler


@cli.command()
@click.option(
    "-c",
    "--commit",
    "commit_interval",
    help="The number of blocks between each commit interval.",
    default=1,
    show_default=True,
)
@click.option(
    "--db",
    help="Database URI connection string.",
    default="localhost:5432/postgres",
    show_default=True,
)
@click.option(
    "-s", "--servers", help="Kafka bootstrap servers.", default="localhost:9092", show_default=True,
)
@click.option("-e", "--echo-sqlalchemy", "echo_sqlalchemy", is_flag=True)
def sync(commit_interval, db, servers, echo_sqlalchemy):
    """Subscribe to Kafka and push the updates to the database."""
    # Set up Kafka connection
    engine = create_engine("postgresql+psycopg2://" + db, echo=echo_sqlalchemy)
    tracking_info = engine.execute(tracking.select()).fetchone()
    topic = tracking_info.topic
    consumer = KafkaConsumer(topic, bootstrap_servers=servers)
    partitions = consumer.partitions_for_topic(topic)
    if len(partitions) != 1:
        raise Exception("Only exact 1 partition is supported.")
    consumer.seek(TopicPartition(topic, partitions.pop()), tracking_info.kafka_offset + 1)
    consumer_iter = iter(consumer)
    # Main loop
    while True:
        with engine.begin() as conn:
            for msg in consumer_iter:
                handler = Handler(conn)
                key = msg.key.decode()
                value = json.loads(msg.value)
                if key == "COMMIT":
                    if value["height"] % commit_interval == 0:
                        conn.execute(tracking.update().values(kafka_offset=msg.offset))
                        logger.info(
                            "Committed at block {} and Kafka offset {}",
                            value["height"],
                            msg.offset,
                        )
                        break
                    continue
                getattr(handler, "handle_" + key.lower())(value)

