import json
import base64

from .client import Client
from .wallet import PublicKey
from .data import Request, RequestInfo

REQUEST_DURATION = 100


class Auth:
    """
    Class for verifying that reporter information is valid to get data for report.

    :ivar client: the Client instance
    """

    def __init__(self, client: Client) -> None:
        self.client = client

    @staticmethod
    def get_msg_sign_bytes(
        chain_id: str, validator: str, request_id: str, external_id: str
    ) -> bytes:
        """
        Return bytes message using in signature verification
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
    def verify_verification_message_signature(
        chain_id: str,
        validator: str,
        request_id: str,
        external_id: str,
        reporter_pubkey: str,
        signature: bytes,
    ) -> bool:
        """
        Verify the verification message using the reporter's signature
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
        Verify report infomation using reporter signature and on-chain request status

        :param chain_id: a chain id that request come from
        :param validator: validator address of this verify meassage
        :param request_id: a request id that request come from
        :param external_id: validator address of this verify meassage
        :param reporter_pubkey: a chain id that request come from
        :param signature: validator address of this verify meassage

        :return: True if the validator has been assigned to report data. False otherwise
        """
        if not Auth.verify_verification_message_signature(
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
        Verify that the request comes from correct chain ID

        :param chain_id: the chain ID retrieved from the verification message

        :return: True if the chain id match with rpc client.
        """
        return self.client.get_chain_id() == chain_id

    def is_reporter(self, validator: str, reporter_pubkey: str) -> bool:
        """
        Verify that an address is a registerd reporter for validator

        :param validator: a validator on BandChain
        :param reporter_pubkey: a public key of reporter

        :return: True if the reporter is a registered reporter of the validator.
        """
        reporter = PublicKey.from_acc_bech32(reporter_pubkey).to_address().to_acc_bech32()
        reporters = self.client.get_reporters(validator)
        return reporter in reporters

    def verify_non_expired_request(self, request: Request) -> bool:
        """
        Verify the request has not expired

        :param request: a request instance

        :return: True if the request hasn't expired.
        """
        latest_block = self.client.get_latest_block()
        return (
            int(latest_block["block"]["header"]["height"]) - request.request_height
            <= REQUEST_DURATION
        )

    def verify_requested_validator(self, request: Request, validator: str) -> bool:
        """
        Verify that this validator has been assigned to report this request

        :param request: a request instance
        :param validator: validator who reporter work for

        :return: True if the validator has been assigned to report data.
        """
        return validator in request.requested_validators

    def verify_unsubmitted_report(self, reports: list, validator: str) -> bool:
        """
        Verify that this validator has not reported on this request

        :param reports: list of reports on this request
        :param validator: validator who reporter work for

        :return: True if the validator hasn't reported data to chain yet
        """
        for report in reports:
            if report.validator == validator:
                return False
        return True
