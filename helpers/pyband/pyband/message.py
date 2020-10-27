from dataclasses import dataclass
import base64

from .wallet import Address


class Message(object):
    def to_dict(self) -> dict:
        pass

    # for with_auto_account_info
    def get_sender(self) -> Address:
        pass


@dataclass
class MsgRequest(Message):
    oracle_script_id: int
    calldata: bytes
    ask_count: int
    min_count: int
    client_id: str
    sender: Address

    def to_dict(self) -> dict:
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
