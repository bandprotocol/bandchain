import base64
import json
from .client import Client
from typing import List, Optional
from .wallet import PublicKey
from .constant import MAX_MEMO_CHARACTERS
from .message import Msg


class Transaction:
    def __init__(self):
        self.msgs: List[Msg] = []
        self.account_num: Optional[int] = None
        self.sequence: Optional[int] = None
        self.chain_id: Optional[str] = None
        self.fee: int = 0
        self.gas: int = 200000
        self.memo: str = ""

    def with_messages(self, *msgs: Msg) -> "Transaction":
        self.msgs.extend(msgs)
        return self

    def with_auto(self, client: Client) -> "Transaction":
        if len(self.msgs) == 0:
            raise ValueError("messsage is empty, please use with_messages at least 1 message")

        addr = self.msgs[0].get_sender()
        account = client.get_account(addr)
        self.account_num = account.account_number
        self.sequence = account.sequence
        return self

    def with_account_num(self, account_num: int) -> "Transaction":
        self.account_num = account_num
        return self

    def with_sequence(self, sequence: int) -> "Transaction":
        self.sequence = sequence
        return self

    def with_chain_id(self, chain_id: str) -> "Transaction":
        self.chain_id = chain_id
        return self

    def with_fee(self, fee: int) -> "Transaction":
        self.fee = fee
        return self

    def with_gas(self, gas: int) -> "Transaction":
        self.gas = gas
        return self

    def with_memo(self, memo: str) -> "Transaction":
        if len(memo) > MAX_MEMO_CHARACTERS:
            raise ValueError("memo is too large")

        self.memo = memo
        return self

    def get_sign_data(self) -> bytes:
        if len(self.msgs) == 0:
            raise ValueError("message is empty")

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
            "msgs": [x.as_json() for x in self.msgs],
        }

        message_str = json.dumps(message_json, separators=(",", ":"), sort_keys=True)
        return message_str.encode("utf-8")

    def get_tx_data(self, signature: bytes, pubkey: PublicKey) -> dict:
        return {
            "msg": [x.as_json() for x in self.msgs],
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
                        "value": base64.b64encode(bytes.fromhex(pubkey.to_hex())).decode("utf-8"),
                    },
                    "account_number": str(self.account_num),
                    "sequence": str(self.sequence),
                }
            ],
        }
