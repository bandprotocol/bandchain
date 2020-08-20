import base64
from dataclasses import dataclass
from typing import List

from dacite import Config


DACITE_CONFIG = Config(type_hooks={int: int, bytes: base64.b64decode})


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
    reports: List[Report]
    result: Result

