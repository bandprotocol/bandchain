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
        """CREATE VIEW validator_last_100_votes AS
			SELECT COUNT(*), consensus_address, voted
			FROM (SELECT * FROM validator_votes ORDER BY block_height DESC LIMIT 30000) tt
			WHERE block_height > (SELECT MAX(height) from blocks) - 101
			GROUP BY consensus_address, voted;"""
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
        """CREATE VIEW oracle_script_statistic_last_1_day AS
            SELECT
            AVG(resolve_time-request_time) as response_time,
            COUNT(*) as count,
            oracle_scripts.id,
            resolve_status
            FROM oracle_scripts
            JOIN requests ON oracle_scripts.id=requests.oracle_script_id
            WHERE to_timestamp(requests.request_time) >= NOW() - '1 day'::INTERVAL
            GROUP BY oracle_scripts.id, requests.resolve_status;
        """
    )
    engine.execute(
        """CREATE VIEW oracle_script_statistic_last_1_week AS
            SELECT
            AVG(resolve_time-request_time) as response_time,
            COUNT(*) as count,
            oracle_scripts.id,
            resolve_status
            FROM oracle_scripts
            JOIN requests ON oracle_scripts.id=requests.oracle_script_id
            WHERE to_timestamp(requests.request_time) >= NOW() - '1 week'::INTERVAL
            GROUP BY oracle_scripts.id, requests.resolve_status;
        """
    )
    engine.execute(
        """CREATE VIEW oracle_script_statistic_last_1_month AS
            SELECT
            AVG(resolve_time-request_time) as response_time,
            COUNT(*) as count,
            oracle_scripts.id,
            resolve_status
            FROM oracle_scripts
            JOIN requests ON oracle_scripts.id=requests.oracle_script_id
            WHERE to_timestamp(requests.request_time) >= NOW() - '1 month'::INTERVAL
            GROUP BY oracle_scripts.id, requests.resolve_status;
        """
    )
    engine.execute(
        """
CREATE VIEW non_validator_vote_proposals_view AS
SELECT validator_id,
       proposal_id,
       answer,
       SUM(CAST(shares AS DECIMAL) * CAST(tokens AS DECIMAL) / CAST(delegator_shares AS DECIMAL)) AS tokens
FROM delegations
JOIN votes ON delegations.delegator_id=votes.voter_id
JOIN validators ON delegations.validator_id=validators.id
AND votes.voter_id != validators.account_id
GROUP BY answer, validator_id, proposal_id;
"""
    )

    engine.execute(
        """
CREATE VIEW validator_vote_proposals_view AS
SELECT validators.id,
       proposal_id,
       answer,
       tokens
FROM votes
JOIN accounts ON accounts.id = votes.voter_id
JOIN validators ON accounts.id = validators.account_id;
"""
    )
