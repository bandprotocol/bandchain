import json

import click

from .db import metadata, tracking
from .cli import cli

from sqlalchemy import create_engine


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
    engine.execute(
        """CREATE VIEW delegations_view AS
            SELECT CAST(shares AS DECIMAL) * CAST(tokens AS DECIMAL) / CAST(delegator_shares AS DECIMAL) as amount,
            CAST(shares AS DECIMAL) /  CAST(delegator_shares AS DECIMAL) * 100 as share_percentage,
            CAST(shares AS DECIMAL) * CAST(current_reward AS DECIMAL) /  CAST(delegator_shares AS DECIMAL) + (CAST(current_ratio AS DECIMAL) - CAST(last_ratio AS DECIMAL)) * CAST(shares AS DECIMAL) as reward,
            validators.operator_address,
            moniker,
            accounts.address AS delegator_address,
            identity
            FROM delegations JOIN validators ON delegations.validator_id=validators.id
            JOIN accounts ON accounts.id=delegations.delegator_id;"""
    )
    engine.execute(
        """CREATE VIEW validator_last_250_votes AS
			SELECT COUNT(*), consensus_address, voted
			FROM (SELECT * FROM validator_votes ORDER BY block_height DESC LIMIT 30000) tt
			WHERE block_height > (SELECT MAX(height) from blocks) - 251
			GROUP BY consensus_address, voted;"""
    )
    engine.execute(
        """CREATE VIEW validator_last_1000_votes AS
			SELECT COUNT(*), consensus_address, voted
			FROM (SELECT * FROM validator_votes ORDER BY block_height DESC LIMIT 100001) tt
			WHERE block_height > (SELECT MAX(height) from blocks) - 1001
			GROUP BY consensus_address, voted;"""
    )
    engine.execute(
        """CREATE VIEW validator_last_10000_votes AS
			SELECT COUNT(*), consensus_address, voted
			FROM (SELECT * FROM validator_votes ORDER BY block_height DESC LIMIT 1000001) tt
			WHERE block_height > (SELECT MAX(height) from blocks) - 10000
			GROUP BY consensus_address, voted;"""
    )
    engine.execute(
        """CREATE VIEW request_counts AS
            SELECT date_trunc('minute', blocks.timestamp) as date, COUNT(*) as request_count
            FROM blocks
            JOIN transactions ON blocks.height = transactions.block_height
            JOIN requests ON requests.transaction_id = transactions.id
            GROUP BY date_trunc('minute',blocks.timestamp);
            """
    )
