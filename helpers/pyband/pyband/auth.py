import json

from .client import Client
from .wallet import PublicKey


class Auth:
    def __init__(self, client: Client) -> None:
        self.client = client

    @staticmethod
    def get_msg_sign_bytes(
        chain_id: str, validator: str, request_id: str, external_id: str
    ) -> bytes:
        return str.encode(
            json.dumps(
                {
                    "chain_id": chain_id,
                    "validator": validator,
                    "request_id": request_id,
                    "external_id": external_id,
                },
                sort_keys=True,
                separators=(",", ":"),
            )
        )

    @staticmethod
    def verify_verification_message(
        chain_id: str,
        validator: str,
        request_id: str,
        external_id: str,
        reporter_pubkey: str,
        signature: bytes,
    ) -> bool:
        reporter = PublicKey.from_acc_pub(reporter_pubkey)

        msg = Auth.get_msg_sign_bytes(chain_id, validator, request_id, external_id)
        return reporter.verify(msg, signature)

    def verify_chain_id(self, chain_id: str) -> bool:
        return self.client.get_chain_id() == chain_id

    # def verify_requested_validator(self)

