import base64
import json

from typing import List, Dict
from .wallet import PrivateKey, PublicKey
from .message import Msg


class Transaction:
    def __init__(self):
        self.msgs: List[Msg] = []
        self.account_num: int = None
        self.sequence: int = None
        self.chain_id: str = None
        self.fee: int = 0
        self.gas: int = 200000
        self.memo: str = ""

    def with_messages(self, *msgs: List[Msg]):
        self.msgs.extend(msgs)
        return self

    def with_account_num(self, v):
        self.account_num = v
        return self

    def with_sequence(self, v):
        self.sequence = v
        return self

    def with_chain_id(self, v):
        self.chain_id = v
        return self

    def with_fee(self, v):
        self.fee = v
        return self

    def with_gas(self, v):
        self.gas = v
        return self

    def with_memo(self, v):
        self.memo = v
        return self

    def get_data_for_sign(self) -> bytes:
        if len(self.msgs) == 0:
            raise ValueError("messages is empty")

        if self.account_num is None:
            raise ValueError("account_num should be defined")

        if self.sequence is None:
            raise ValueError("sequence should be defined")

        if self.chain_id is None:
            raise ValueError("chain_id should be defined")

        for msg in self.msgs:
            msg.validate()

        message_json = {
            "chain_id": self.chain_id,
            "account_number": str(self.account_num),
            "fee": {
                "gas": str(self.gas),
                "amount": [{"amount": str(self.fee), "denom": "uband"}],
            },
            "memo": self.memo,
            "sequence": str(self.sequence),
            "msgs": list(map(lambda x: x.as_json(), self.msgs)),
        }

        message_str = json.dumps(message_json, separators=(",", ":"), sort_keys=True)
        return message_str.encode("utf-8")

    def get_data_for_send(self, signature: bytes, pubkey: PublicKey) -> Dict:
        return {
            "msg": list(map(lambda x: x.as_json(), self.msgs)),
            "fee": {
                "gas": str(self.gas),
                "amount": [{"denom": "uband", "amount": str(self.fee)}],
            },
            "memo": self.memo,
            "signatures": [
                {
                    "signature": base64.b64encode(signature).decode("utf-8"),
                    "pub_key": {
                        "type": "tendermint/PubKeySecp256k1",
                        "value": base64.b64encode(
                            bytes.fromhex(pubkey.to_hex())
                        ).decode("utf-8"),
                    },
                    "account_number": str(self.account_num),
                    "sequence": str(self.sequence),
                }
            ],
        }
