import base64 as b64
from datetime import datetime
import sqlalchemy as sa


class CustomDateTime(sa.types.TypeDecorator):
    """Custom DateTime type that accepts Python nanosecond epoch int."""

    impl = sa.DateTime

    def process_bind_param(self, value, dialect):
        return datetime.fromtimestamp(value / 1e9)


class CustomBase64(sa.types.TypeDecorator):
    """Custom String type that accepts base64-encoded string."""

    impl = sa.String

    def process_bind_param(self, value, dialect):
        return b64.decodestring(value.encode())


def Column(*args, **kwargs):
    """Forward into SQLAlchemy's Column construct, but with 'nullable' default to False."""
    if "nullable" not in kwargs:
        kwargs["nullable"] = False
    return sa.Column(*args, **kwargs)


metadata = sa.MetaData()


tracking = sa.Table(
    "tracking",
    metadata,
    Column("chain_id", sa.String, primary_key=True),
    Column("topic", sa.String),
    Column("kafka_offset", sa.Integer),
)


blocks = sa.Table(
    "blocks",
    metadata,
    Column("height", sa.Integer, primary_key=True),
    Column("timestamp", CustomDateTime),
    Column("proposer", CustomBase64, sa.ForeignKey("validators.consensus_address")),
    Column("hash", CustomBase64),
    Column("inflation", sa.Float),
    Column("supply", sa.String),  # uband suffix
)

transactions = sa.Table(
    "transactions",
    metadata,
    Column("hash", CustomBase64, primary_key=True),
    Column("block_height", sa.Integer, sa.ForeignKey("blocks.height")),
    Column("index", sa.Integer),
    Column("gas_used", sa.Integer),
    Column("gas_limit", sa.Integer),
    Column("gas_fee", sa.String),  # uband suffix
    Column("err_msg", sa.String, nullable=True),
    Column("success", sa.Boolean),
    Column("sender", sa.String),
    Column("memo", sa.String),
    Column("messages", sa.JSON),
)


accounts = sa.Table(
    "accounts",
    metadata,
    Column("address", sa.String, primary_key=True),
    Column("balance", sa.String),
)

data_sources = sa.Table(
    "data_sources",
    metadata,
    Column("id", sa.Integer, primary_key=True),
    Column("name", sa.String),
    Column("description", sa.String),
    Column("owner", sa.String),
    Column("executable", CustomBase64),
)

oracle_scripts = sa.Table(
    "oracle_scripts",
    metadata,
    Column("id", sa.Integer, primary_key=True),
    Column("name", sa.String),
    Column("description", sa.String),
    Column("owner", sa.String),
    Column("schema", sa.String),
    Column("codehash", sa.String),
    Column("source_code_url", sa.String),
)

requests = sa.Table(
    "requests",
    metadata,
    Column("id", sa.Integer, primary_key=True),
    Column("tx_hash", CustomBase64, sa.ForeignKey("transactions.hash")),
    Column("oracle_script_id", sa.Integer, sa.ForeignKey("oracle_scripts.id")),
    Column("calldata", CustomBase64),
    Column("ask_count", sa.Integer),
    Column("min_count", sa.Integer),
    Column("sender", sa.String),
    Column("client_id", sa.String),
)

raw_requests = sa.Table(
    "raw_requests",
    metadata,
    Column("request_id", sa.Integer, sa.ForeignKey("requests.id"), primary_key=True),
    Column("external_id", sa.Integer, primary_key=True),
    Column("data_source_id", sa.Integer),
    Column("calldata", CustomBase64),
)

val_requests = sa.Table(
    "val_requests",
    metadata,
    Column("request_id", sa.Integer, sa.ForeignKey("requests.id"), primary_key=True),
    Column("validator", sa.String, sa.ForeignKey("validators.operator_address"), primary_key=True),
)

reports = sa.Table(
    "reports",
    metadata,
    Column("request_id", sa.Integer, primary_key=True),  # Add FK to requests
    Column("validator", sa.String, sa.ForeignKey("validators.operator_address"), primary_key=True),
    Column("tx_hash", CustomBase64, sa.ForeignKey("transactions.hash")),
    Column("reporter", sa.String),
)

raw_reports = sa.Table(
    "raw_reports",
    metadata,
    Column("request_id", sa.Integer, primary_key=True),
    Column("validator", sa.String, primary_key=True),
    Column("external_id", sa.Integer, primary_key=True),
    Column("data", CustomBase64),
    Column("exit_code", sa.Integer),
    sa.ForeignKeyConstraint(
        ["request_id", "validator"], ["reports.request_id", "reports.validator"]
    ),
)

validators = sa.Table(
    "validators",
    metadata,
    Column("operator_address", sa.String, primary_key=True),
    Column("consensus_address", CustomBase64, unique=True),
    Column("consensus_pubkey", sa.String),
    Column("moniker", sa.String),
    Column("website", sa.String),
    Column("details", sa.String),
    Column("commission_rate", sa.String),
    Column("commission_max_rate", sa.String),
    Column("commission_max_change", sa.String),
    Column("min_self_delegation", sa.String),
    Column("jailed", sa.Boolean),
    Column("tokens", sa.Integer),
    Column("delegator_shares", sa.String),
    Column("bonded_height", sa.Integer),
    Column("current_reward", sa.String),
    Column("current_ratio", sa.String),
)
