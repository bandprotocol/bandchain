import base64
from dataclasses import dataclass
from .wallet import Address
from .constant import MAX_CLIENT_ID_LENGTH, MAX_DATA_SIZE


class Message(object):
    def as_json(self) -> dict:
        pass

    def get_sender(self) -> Address:
        pass

    def validate(self) -> bool:
        pass


@dataclass
class MsgRequest(Message):
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
            raise ValueError(
                "invalid min count got: min count: {}".format(self.min_count))
        if self.ask_count < self.min_count:
            raise ValueError("invalid ask count got: min count: {}, ask count: {}".format(
                self.min_count, self.ask_count))
        if len(self.client_id) > MAX_CLIENT_ID_LENGTH:
            raise ValueError("too long client id")

        return True
