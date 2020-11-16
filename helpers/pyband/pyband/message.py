import base64

from dataclasses import dataclass
from typing import List
from .wallet import Address
from .data import Coin
from .constant import MAX_CLIENT_ID_LENGTH, MAX_DATA_SIZE


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
            raise ValueError("oracle script id cannot less than zero")
        if len(self.calldata) > MAX_DATA_SIZE:
            raise ValueError("too large calldata")
        if self.min_count <= 0:
            raise ValueError("invalid min count got: min count: {}".format(self.min_count))
        if self.ask_count < self.min_count:
            raise ValueError(
                "invalid ask count got: min count: {}, ask count: {}".format(self.min_count, self.ask_count)
            )
        if len(self.client_id) > MAX_CLIENT_ID_LENGTH:
            raise ValueError("too long client id")

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
            raise ValueError("Expect at least 1 coin")

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
