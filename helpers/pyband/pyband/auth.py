import json
import base64

from .client import Client
from .wallet import PublicKey
from .data import Request, RequestInfo

REQUEST_DURATION = 100


class Auth:
    def __init__(self, client: Client) -> None:
        self.client = client

    @staticmethod
    def get_msg_sign_bytes(
        chain_id: str, validator: str, request_id: str, external_id: str
    ) -> bytes:
        """
        Return message using in signature verification
        """
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
        """
        Verify verification message by signature of reporter
        """
        reporter = PublicKey.from_acc_bech32(reporter_pubkey)

        msg = Auth.get_msg_sign_bytes(chain_id, validator, request_id, external_id)
        return reporter.verify(msg, signature)

    def verify(
        self,
        chain_id: str,
        validator: str,
        request_id: str,
        external_id: str,
        reporter_pubkey: str,
        signature: str,
    ) -> bool:
        """
        Verify header of request that valid
        """
        if not Auth.verify_verification_message(
            chain_id,
            validator,
            request_id,
            external_id,
            reporter_pubkey,
            base64.b64decode(signature),
        ):
            return False
        if not self.verify_chain_id(chain_id):
            return False

        if not self.is_reporter(validator, reporter_pubkey):
            return False

        requestInfo = self.client.get_request_by_id(request_id)

        if not self.verify_non_expired_request(requestInfo.request):
            return False
        if not self.verify_requested_validator(requestInfo.request, validator):
            return False
        if not self.verify_unsubmitted_report(requestInfo.reports, validator):
            return False

        return True

    def verify_chain_id(self, chain_id: str) -> bool:
        """
        Verify request come from correct chain id
        """
        return self.client.get_chain_id() == chain_id

    def is_reporter(self, validator: str, reporter_pubkey: str) -> bool:
        """
        Verify this address is a registerd reporter for validator
        """
        reporter = PublicKey.from_acc_bech32(reporter_pubkey).to_address().to_acc_bech32()
        reporters = self.client.get_reporters(validator)
        return reporter in reporters

    def verify_non_expired_request(self, request: Request) -> bool:
        """
        Verify the request has not been expired
        """
        latest_block = self.client.get_latest_block()
        return (
            latest_block["block"]["header"]["height"] - request.request_height <= REQUEST_DURATION
        )

    def verify_requested_validator(self, request: Request, validator: str) -> bool:
        """
        Verify this validator has been assigned to report this request
        """
        return validator in request.requested_validators

    def verify_unsubmitted_report(self, reports: list, validator: str) -> bool:
        """
        Verify this validator has not been reported on this request
        """
        for report in reports:
            if report.validator == validator:
                return False
        return True

