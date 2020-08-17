import base64
from dataclasses import dataclass
from typing import List

from dacite import Config


DACITE_CONFIG = Config(type_hooks={int: int, bytes: base64.b64decode})


@dataclass
class DataSource(object):
    owner: str
    name: str
    description: str
    filename: str


@dataclass
class OracleScript(object):
    owner: str
    name: str
    description: str
    filename: str
    schema: str
    source_code_url: str


@dataclass
class RawRequest(object):
    external_id: int
    data_source_id: int
    calldata: bytes


@dataclass
class Request(object):
    oracle_script_id: int
    calldata: bytes
    requested_validators: List[str]
    min_count: int
    request_height: int
    client_id: str
    raw_requests: List[RawRequest]


@dataclass
class RawReport(object):
    external_id: int
    data: bytes


@dataclass
class Report(object):
    validator: str
    in_before_resolve: bool
    raw_reports: List[RawReport]


@dataclass
class RequestPacketData:
    client_id: str
    oracle_script_id: int
    calldata: bytes
    ask_count: int
    min_count: int


@dataclass
class ResponsePacketData(object):
    client_id: str
    request_id: int
    ans_count: int
    request_time: int
    resolve_time: int
    resolve_status: int
    result: bytes


@dataclass
class Result(object):
    request_packet_data: RequestPacketData
    response_packet_data: ResponsePacketData


@dataclass
class RequestInfo(object):
    request: Request
    reports: List[Report]
    result: Result

