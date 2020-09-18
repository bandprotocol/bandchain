import base64
import json

from .wallet import PrivateKey
from .client import Client


class Transaction:
    def __init__(self, client: Client, privkey: PrivateKey):
        self.client = client
        self.privkey = privkey
        self.msgs: list = []
        self.chain_id = client.get_chain_id()

    def send_tx(
        self,
        account_num: int,
        sequence: int,
        gas: int = 200000,
        fee: int = 0,
        memo: str = "",
        mode="sync",
    ) -> dict:
        base64_pubkey = base64.b64encode(bytes.fromhex(self.privkey.to_pubkey().to_hex())).decode(
            "utf-8"
        )
        tx = {
            "tx": {
                "msg": self.msgs,
                "fee": {"gas": str(gas), "amount": [{"denom": "uband", "amount": str(fee)}]},
                "memo": memo,
                "signatures": [
                    {
                        "signature": self._sign(account_num, sequence, gas, fee, memo),
                        "pub_key": {"type": "tendermint/PubKeySecp256k1", "value": base64_pubkey,},
                        "account_number": str(account_num),
                        "sequence": str(sequence),
                    }
                ],
            },
            "mode": mode,
        }
        return self.client.send_tx(tx)

    def add_request(
        self, oracle_script_id: int, calldata: bytes, ask_count: int, min_count: int, client_id: str
    ) -> None:
        msg = {
            "type": "oracle/Request",
            "value": {
                "oracle_script_id": str(oracle_script_id),
                "calldata": base64.b64encode(calldata).decode("utf-8"),
                "ask_count": str(ask_count),
                "min_count": str(min_count),
                "client_id": client_id,
                "sender": self.privkey.to_pubkey().to_address().to_acc_bech32(),
            },
        }
        self.msgs.append(msg)

    def _sign(self, account_num: int, sequence: int, gas: int, fee: int, memo: str) -> str:
        message_str = json.dumps(
            self._get_sign_message(account_num, sequence, gas, fee, memo),
            separators=(",", ":"),
            sort_keys=True,
        )
        message_bytes = message_str.encode("utf-8")

        signature = self.privkey.sign(message_bytes)

        signature_base64_str = base64.b64encode(signature).decode("utf-8")
        return signature_base64_str

    def _get_sign_message(
        self, account_num: int, sequence: int, gas: int, fee: int, memo: str
    ) -> dict:
        return {
            "chain_id": self.chain_id,
            "account_number": str(account_num),
            "fee": {"gas": str(gas), "amount": [{"amount": str(fee), "denom": "uband"}]},
            "memo": memo,
            "sequence": str(sequence),
            "msgs": self.msgs,
        }
