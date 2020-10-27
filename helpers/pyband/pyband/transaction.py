import base64
import json
from typing import List

from .wallet import PrivateKey
from .client import Client
from .message import Message


class Transaction(object):
    def __init__(self, ):
        self.msgs: List[Message] = []
        self.account_num = None
        self.sequence = None
        self.fee = 0
        self.gas = 200000
        self.memo = ""

    def with_msessage(self, *msgs: list(Message)):
        self.msgs.extend(msgs)
        return self

    def with_account_num(self, v):
        self.account_num = v
        return self

    def with_sequence(self, v):
        self.account_num = v
        return self
    
    def get_signable_data(self):
        # same as get_raw_data
        # raise account_num is None
        # raise sequnce is None
        # return ready data for sign

    def as_json(self, signature):
        return {
            "tx": {
                "msg": parsed_msgs,
                "fee": {
                    "gas": str(gas),
                    "amount": [{"denom": "uband", "amount": str(fee)}],
                },
                "memo": memo,
                "signatures": [
                    {
                        "signature": signature,
                        "pub_key": {
                            "type": "tendermint/PubKeySecp256k1",
                            "value": base64_pubkey,
                        },
                        "account_number": str(account_num),
                        "sequence": str(sequence),
                    }
                ],
            },
        }
    







    # def get_raw_tx(
    #     account_num: int,
    #     sequence: int,
    #     chain_id: str,
    #     msgs: list,
    #     gas: int = 200000,
    #     fee: int = 0,
    #     memo: str = "",
    #     mode="sync",
    # ) -> dict:
        
    #     if self.signature is None:
    #         raise ValueError("sdkskdksd")

    #     base64_pubkey = base64.b64encode(
    #         bytes.fromhex(privkey.to_pubkey().to_hex())
    #     ).decode("utf-8")

    #     parsed_msgs = list(map(lambda x: x.to_dict(), msgs))

    #     raw_data_bytes = _get_raw_data(
    #         account_num, sequence, chain_id, gas, fee, memo, parsed_msgs
    #     )

    #     return {
    #         "tx": {
    #             "msg": parsed_msgs,
    #             "fee": {
    #                 "gas": str(gas),
    #                 "amount": [{"denom": "uband", "amount": str(fee)}],
    #             },
    #             "memo": memo,
    #             "signatures": [
    #                 {
    #                     "signature": self.signature,
    #                     "pub_key": {
    #                         "type": "tendermint/PubKeySecp256k1",
    #                         "value": base64_pubkey,
    #                     },
    #                     "account_number": str(account_num),
    #                     "sequence": str(sequence),
    #                 }
    #             ],
    #         },
    #         "mode": mode,
    #     }


    # # 
    # def _sign(privkey: PrivateKey, message_bytes: bytes) -> str:
    #     signature = privkey.sign(message_bytes)

    #     signature_base64_str = base64.b64encode(signature).decode("utf-8")
    #     return signature_base64_str

    # def _get_raw_data(
    #     account_num: int,
    #     sequence: int,
    #     chain_id: str,
    #     gas: int,
    #     fee: int,
    #     memo: str,
    #     msgs: list,
    # ) -> bytes:
    #     raw_data = {
    #         "chain_id": chain_id,
    #         "account_number": str(account_num),
    #         "fee": {
    #             "gas": str(gas),
    #             "amount": [{"amount": str(fee), "denom": "uband"}],
    #         },
    #         "memo": memo,
    #         "sequence": str(sequence),
    #         "msgs": msgs,
    #     }
    #     print(raw_data)
    #     message_str = json.dumps(raw_data, separators=(",", ":"), sort_keys=True,)

    #     return message_str.encode("utf-8")
