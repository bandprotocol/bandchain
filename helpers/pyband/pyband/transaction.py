import base64
import json

from .wallet import PrivateKey
from .client import Client
from .message import Message


def get_raw_tx(
    account_num: int,
    sequence: int,
    chain_id: str,
    msgs: list,
    privkey: PrivateKey,
    gas: int = 200000,
    fee: int = 0,
    memo: str = "",
    mode="sync",
) -> dict:
    base64_pubkey = base64.b64encode(
        bytes.fromhex(privkey.to_pubkey().to_hex())
    ).decode("utf-8")

    raw_data_bytes = _get_raw_data(
        account_num, sequence, chain_id, gas, fee, memo, msgs
    )

    return {
        "tx": {
            "msg": list(map(lambda x: x.to_dict(), msgs)),
            "fee": {
                "gas": str(gas),
                "amount": [{"denom": "uband", "amount": str(fee)}],
            },
            "memo": memo,
            "signatures": [
                {
                    "signature": _sign(privkey, raw_data_bytes),
                    "pub_key": {
                        "type": "tendermint/PubKeySecp256k1",
                        "value": base64_pubkey,
                    },
                    "account_number": str(account_num),
                    "sequence": str(sequence),
                }
            ],
        },
        "mode": mode,
    }


def _sign(privkey: PrivateKey, message_bytes: bytes) -> str:
    signature = privkey.sign(message_bytes)

    signature_base64_str = base64.b64encode(signature).decode("utf-8")
    return signature_base64_str


def _get_raw_data(
    account_num: int,
    sequence: int,
    chain_id: str,
    gas: int,
    fee: int,
    memo: str,
    msgs: list,
) -> bytes:
    raw_data = {
        "chain_id": chain_id,
        "account_number": str(account_num),
        "fee": {"gas": str(gas), "amount": [{"amount": str(fee), "denom": "uband"}],},
        "memo": memo,
        "sequence": str(sequence),
        "msgs": list(map(lambda x: x.to_dict(), msgs)),
    }
    print(raw_data)
    message_str = json.dumps(raw_data, separators=(",", ":"), sort_keys=True,)

    return message_str.encode("utf-8")
