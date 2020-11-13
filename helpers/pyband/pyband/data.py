import base64

from dataclasses import dataclass
from typing import List, Optional, NewType
from dacite import Config
from pyband.utils import parse_datetime

HexBytes = NewType("HexBytes", bytes)
Timestamp = NewType("Timestamp", int)

DACITE_CONFIG = Config(
    type_hooks={
        int: int,
        bytes: base64.b64decode,
        HexBytes: bytes.fromhex,
        Timestamp: parse_datetime,
    }
)


@dataclass
class DataSource(object):
    owner: str
    name: str = ""
    description: str = ""
    filename: str = ""


@dataclass
class OracleScript(object):
    owner: str
    name: str = ""
    description: str = ""
    filename: str = ""
    schema: str = ""
    source_code_url: str = ""


@dataclass
class RawRequest(object):
    data_source_id: int
    external_id: int = 0
    calldata: bytes = b""


@dataclass
class Request(object):
    oracle_script_id: int
    requested_validators: List[str]
    min_count: int
    request_height: int
    raw_requests: List[RawRequest]
    client_id: str = ""
    calldata: bytes = b""


@dataclass
class RawReport(object):
    external_id: int = 0
    data: bytes = b""


@dataclass
class Report(object):
    validator: str
    raw_reports: List[RawReport]
    in_before_resolve: bool = False


@dataclass
class RequestPacketData:
    oracle_script_id: int
    ask_count: int
    min_count: int
    client_id: str = ""
    calldata: bytes = b""


@dataclass
class ResponsePacketData(object):
    request_id: int
    request_time: int
    resolve_time: int
    resolve_status: int
    ans_count: int = 0
    client_id: str = ""
    result: bytes = b""


@dataclass
class Result(object):
    request_packet_data: RequestPacketData
    response_packet_data: ResponsePacketData


@dataclass
class RequestInfo(object):
    request: Request
    reports: Optional[List[Report]]
    result: Optional[Result]


@dataclass
class Coin(object):
    amount: int
    denom: str

    def as_json(self) -> dict:
        return {"amount": str(self.amount), "denom": self.denom}

    def validate(self) -> bool:
        if self.amount < 0:
            raise ValueError("Expect amount more than 0")

        if len(self.denom) == 0:
            raise ValueError("Expect denom")

        return True


@dataclass
class Account(object):
    address: str
    coins: List[dict]
    public_key: Optional[dict]
    account_number: int
    sequence: int

@dataclass
class TransactionSyncMode(object):
    tx_hash: HexBytes
    code: int
    error_log: Optional[str]


@dataclass
class TransactionAsyncMode(object):
    tx_hash: HexBytes


@dataclass
class TransactionBlockMode(object):
    height: int
    tx_hash: HexBytes
    gas_wanted: int
    gas_used: int
    code: int
    log: List[dict]
    error_log: Optional[str]


@dataclass
class BlockHeaderInfo(object):
    chain_id: str
    height: int
    time: Timestamp
    last_commit_hash: HexBytes
    data_hash: HexBytes
    validators_hash: HexBytes
    next_validators_hash: HexBytes
    consensus_hash: HexBytes
    app_hash: HexBytes
    last_results_hash: HexBytes
    evidence_hash: HexBytes
    proposer_address: HexBytes


@dataclass
class BlockHeader(object):
    header: BlockHeaderInfo


@dataclass
class BlockID(object):
    hash: HexBytes


@dataclass
class Block(object):
    block: BlockHeader
    block_id: BlockID
