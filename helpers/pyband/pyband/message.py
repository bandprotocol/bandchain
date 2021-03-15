import base64

from dataclasses import dataclass
from typing import List, Optional
from .wallet import Address
from .data import Coin
from .constant import MAX_CLIENT_ID_LENGTH, MAX_DATA_SIZE
from .exceptions import InsufficientCoinError, NegativeIntegerError, ValueTooLargeError


class Msg:
    def as_json(self) -> dict:
        raise NotImplementedError()

    def get_sender(self) -> Address:
        raise NotImplementedError()

    def validate(self) -> bool:
        raise NotImplementedError()


@dataclass
class MsgRequest(Msg):
    oracle_script_id: int
    calldata: bytes
    ask_count: int
    min_count: int
    client_id: str
    sender: Address

    def as_json(self) -> dict:
        return {
            "type": "oracle/Request",
            "value": {
                "oracle_script_id": str(self.oracle_script_id),
                "calldata": base64.b64encode(self.calldata).decode("utf-8"),
                "ask_count": str(self.ask_count),
                "min_count": str(self.min_count),
                "client_id": self.client_id,
                "sender": self.sender.to_acc_bech32(),
            },
        }

    def get_sender(self) -> Address:
        return self.sender

    def validate(self) -> bool:
        if self.oracle_script_id <= 0:
            raise NegativeIntegerError("oracle script id cannot less than zero")
        if len(self.calldata) > MAX_DATA_SIZE:
            raise ValueTooLargeError("too large calldata")
        if self.min_count <= 0:
            raise ValueError("invalid min count got: min count: {}".format(self.min_count))
        if self.ask_count < self.min_count:
            raise ValueError(
                "invalid ask count got: min count: {}, ask count: {}".format(self.min_count, self.ask_count)
            )
        if len(self.client_id) > MAX_CLIENT_ID_LENGTH:
            raise ValueTooLargeError("too long client id")

        return True


@dataclass
class MsgSend(Msg):
    to_address: Address
    from_address: Address
    amount: List[Coin]

    def as_json(self) -> dict:
        return {
            "type": "cosmos-sdk/MsgSend",
            "value": {
                "to_address": self.to_address.to_acc_bech32(),
                "from_address": self.from_address.to_acc_bech32(),
                "amount": [coin.as_json() for coin in self.amount],
            },
        }

    def get_sender(self) -> Address:
        return self.from_address

    def validate(self) -> bool:
        if len(self.amount) == 0:
            raise InsufficientCoinError("Expect at least 1 coin")

        for coin in self.amount:
            coin.validate()

        return True


@dataclass
class MsgDelegate(Msg):
    delegator_address: Address
    validator_address: Address
    amount: Coin

    def as_json(self) -> dict:
        return {
            "type": "cosmos-sdk/MsgDelegate",
            "value": {
                "delegator_address": self.delegator_address.to_acc_bech32(),
                "validator_address": self.validator_address.to_val_bech32(),
                "amount": self.amount.as_json(),
            },
        }

    def get_sender(self) -> Address:
        return self.delegator_address

    def validate(self) -> bool:
        self.amount.validate()

        return True


@dataclass
class MsgCreateDataSource(Msg):
    owner: Address
    name: str
    description: str
    script: str
    sender: Address

    def _replace_empty(self, s: Optional[str], placeholder: str) -> str:
        if s is None:
            return placeholder
        return s

    def _read_script(self) -> str:
        file_content = open(self.script, "rb").read()
        return base64.b64encode(file_content).decode()

    def as_json(self) -> dict:
        return {
            "type": "oracle/CreateDataSource",
            "value": {
                "owner": self.owner.to_acc_bech32(),
                "name": self.name,
                "description": self._replace_empty(self.description, "TBA"),
                "executable": self._read_script(),
                "sender": self.sender.to_acc_bech32(),
            },
        }

    def get_sender(self) -> Address:
        return self.sender

    def validate(self) -> bool:
        if len(self.name) <= 0:
            raise ValueError("Invalid oracle script name")
        if len(self.script) <= 0:
            raise ValueError("Missing or invalid data source path")
        if len(self._read_script()) <= 0:
            raise ValueError("Empty data source file")
        return True


@dataclass
class MsgEditDataSource(Msg):
    owner: Address
    sender: Address
    data_source_id: int
    name: Optional[str] = None
    description: Optional[str] = None
    script: Optional[str] = None

    def _replace_empty(self, s: Optional[str], placeholder: str) -> str:
        if s is None:
            return placeholder
        return s

    def _read_optional_script(self) -> str:
        if self.script is None:
            return base64.b64encode(b"[do-not-modify]").decode()
        file_content = open(self.script, "rb").read()
        return base64.b64encode(file_content).decode()

    def as_json(self) -> dict:
        return {
            "type": "oracle/EditDataSource",
            "value": {
                "data_source_id": self.data_source_id,
                "owner": self.owner.to_acc_bech32(),
                "name": self._replace_empty(self.name, "[do-not-modify]"),
                "description": self._replace_empty(self.description, "[do-not-modify]"),
                "executable": self._read_optional_script(),
                "sender": self.sender.to_acc_bech32(),
            },
        }

    def get_sender(self) -> Address:
        return self.sender

    def validate(self) -> bool:
        if self.name is not None and len(self.name) <= 0:
            raise ValueError("Invalid data source name")
        if self.script is not None and len(self.script) <= 0:
            raise ValueError("Invalid oracle script path")
        if self.script is not None and len(self.script) > 0 and len(self._read_optional_script()) <= 0:
            raise ValueError("Empty data source file")
        return True


@dataclass
class MsgCreateOracleScript(Msg):
    owner: Address
    name: str
    description: Optional[str]
    script: str
    schema: Optional[str]
    source_code_url: Optional[str]
    sender: Address

    def replace_empty(self, s: Optional[str], placeholder: str) -> str:
        if s is None:
            return placeholder
        return s

    def _read_wasm(self) -> str:
        file_content = open(self.script, "rb").read()
        return base64.b64encode(file_content).decode()

    def as_json(self) -> dict:
        return {
            "type": "oracle/CreateOracleScript",
            "value": {
                "owner": self.owner.to_acc_bech32(),
                "name": self.name,
                "description": self.replace_empty(self.description, "TBA"),
                "code": self._read_wasm(),
                "schema": self.replace_empty(self.schema, "TBA"),
                "source_code_url": self.replace_empty(self.source_code_url, ""),
                "sender": self.sender.to_acc_bech32(),
            },
        }

    def get_sender(self) -> Address:
        return self.sender

    def validate(self) -> bool:
        if len(self.name) <= 0:
            raise ValueError("Missing or invalid oracle script name")
        if len(self.script) <= 0:
            raise ValueError("Missing or invalid oracle script path")
        if len(self._read_wasm()) <= 0:
            raise ValueError("Empty wasm file")
        return True


@dataclass
class MsgEditOracleScript(Msg):
    owner: Address
    oracle_script_id: int
    schema: Optional[str]
    source_code_url: Optional[str]
    sender: Address
    name: Optional[str] = None
    description: Optional[str] = None
    script: Optional[str] = None

    def _replace_empty(self, s: Optional[str], placeholder: str) -> str:
        if s is None:
            return placeholder
        return s

    def _read_optional_script(self) -> str:
        if self.script is None:
            return base64.b64encode(b"[do-not-modify]").decode()
        file_content = open(self.script, "rb").read()
        return base64.b64encode(file_content).decode()

    def as_json(self) -> dict:
        return {
            "type": "oracle/EditOracleScript",
            "value": {
                "oracle_script_id": self.oracle_script_id,
                "owner": self.owner.to_acc_bech32(),
                "name": self._replace_empty(self.name, "[do-not-modify]"),
                "description": self._replace_empty(self.description, "[do-not-modify]"),
                "code": self._read_optional_script(),
                "schema": self._replace_empty(self.schema, "[do-not-modify]"),
                "source_code_url": self._replace_empty(self.source_code_url, "[do-not-modify]"),
                "sender": self.sender.to_acc_bech32(),
            },
        }

    def get_sender(self) -> Address:
        return self.sender

    def validate(self) -> bool:
        if self.name is not None and len(self.name) <= 0:
            raise ValueError("Invalid oracle script name")
        if self.script is not None and len(self.script) <= 0:
            raise ValueError("Invalid oracle script path")
        if self.script is not None and len(self.script) > 0 and len(self._read_optional_script()) <= 0:
            raise ValueError("Empty wasm file")
        return True
