import base64
from unittest import mock
import copy

from pyband.client import Client
from pyband.auth import Auth
from pyband.data import (
    Request,
    RequestInfo,
    Report,
    Result,
    RequestPacketData,
    ResponsePacketData,
)

VALIDATOR_TEST = "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"

REQUEST_TEST = RequestInfo(
    Request(
        1,
        [
            "bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw",
            "bandvaloper1v38hewjc0865dm4t89v5efh9rmum5rmrm7evg4",
            "bandvaloper1alzj765pzuhtjkmslme4fdpeakc0036xnyjltn",
            "bandvaloper1yyv5jkqaukq0ajqn7vhkyhpff7h6e99j3gv0tr",
            "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec",
        ],
        3,
        136629,
        [],
        "test",
        base64.b64decode("AAAABFVTRFQAAAADQ05ZAAAAAAAPQkA="),
    ),
    [
        Report("bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw", [], False),
        Report("bandvaloper1yyv5jkqaukq0ajqn7vhkyhpff7h6e99j3gv0tr", [], True),
    ],
    Result(
        RequestPacketData(1, 5, 3, "test", base64.b64decode("AAAABFVTRFQAAAADQ05ZAAAAAAAPQkA="),),
        ResponsePacketData(
            3000, 1596950963, 1596950966, 1, 3, "test", base64.b64decode("AAAAAABqbB0=")
        ),
    ),
)


def test_get_msg_sign_bytes():
    assert Auth.get_msg_sign_bytes("bandchain", VALIDATOR_TEST, "1", "1") == bytes.fromhex(
        "7b22636861696e5f6964223a2262616e64636861696e222c2265787465726e616c5f6964223a2231222c22726571756573745f6964223a2231222c2276616c696461746f72223a2262616e6476616c6f706572317034307968337a6b6d6863763065637170336d63617a7938337361353772676a646536776563227d"
    )


def test_verify_verification_message():
    assert (
        Auth.verify_verification_message_signature(
            "bandchain",
            VALIDATOR_TEST,
            "1",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            base64.b64decode(
                "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ=="
            ),
        )
        == True
    )


@mock.patch("pyband.client.Client")
def test_verify(mock_client):
    auth = Auth(mock_client)
    # Fail signature verification
    assert (
        auth.verify(
            "bandchain",
            VALIDATOR_TEST,
            "2",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ==",
        )
        == False
    )

    # Wrong chain id
    mock_client.get_chain_id.return_value = "fake_id"
    assert (
        auth.verify(
            "bandchain",
            VALIDATOR_TEST,
            "1",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ==",
        )
        == False
    )

    mock_client.get_chain_id.return_value = "bandchain"

    # Unauthorized reporter
    mock_client.get_reporters.return_value = [
        "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        "band1ath4dk8e4fwqs5qmp3nnxspar5q0arrtpxy8lx",
    ]
    assert (
        auth.verify(
            "bandchain",
            VALIDATOR_TEST,
            "1",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ==",
        )
        == False
    )
    mock_client.get_reporters.return_value = ["band1wmvh4uzemujfap5graugzckeazr39uy6lesd0g"]

    # Expired request
    mock_client.get_request_by_id.return_value = REQUEST_TEST
    mock_client.get_latest_block.return_value = {"block": {"header": {"height": 136730}}}
    assert (
        auth.verify(
            "bandchain",
            VALIDATOR_TEST,
            "1",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ==",
        )
        == False
    )
    mock_client.get_latest_block.return_value = {"block": {"header": {"height": 136639}}}

    request = copy.deepcopy(REQUEST_TEST)
    request.request.requested_validators = request.request.requested_validators[:-1]
    mock_client.get_request_by_id.return_value = request
    assert (
        auth.verify(
            "bandchain",
            VALIDATOR_TEST,
            "1",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ==",
        )
        == False
    )

    # Must return false if validator already reported data.
    request = copy.deepcopy(REQUEST_TEST)
    request.reports = [Report(VALIDATOR_TEST, [], True)]
    mock_client.get_request_by_id.return_value = request
    assert (
        auth.verify(
            "bandchain",
            VALIDATOR_TEST,
            "1",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ==",
        )
        == False
    )

    mock_client.get_request_by_id.return_value = REQUEST_TEST
    assert (
        auth.verify(
            "bandchain",
            VALIDATOR_TEST,
            "1",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ==",
        )
        == True
    )


@mock.patch("pyband.client.Client")
def test_verify_chain_id(mock_client):
    auth = Auth(mock_client)

    mock_client.get_chain_id.return_value = "bandchain"
    assert auth.verify_chain_id("bandchain") == True
    assert auth.verify_chain_id("band-testnet") == False


@mock.patch("pyband.client.Client")
def test_is_reporter(mock_client):
    auth = Auth(mock_client)

    mock_client.get_reporters.return_value = [
        "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        "band1ath4dk8e4fwqs5qmp3nnxspar5q0arrtpxy8lx",
    ]
    assert (
        auth.is_reporter(
            "mock_validator",
            "bandpub1addwnpepqdg7nrsmuztj2re07svgcz4vuzn3de56nykdwlualepkk05txs5q6mw8s9v",
        )
        == True
    )
    assert (
        auth.is_reporter(
            "mock_validator",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
        )
        == False
    )


@mock.patch("pyband.client.Client")
def test_verify_non_expired_request(mock_client):
    auth = Auth(mock_client)

    mock_client.get_latest_block.return_value = {"block": {"header": {"height": 136639}}}
    assert auth.verify_non_expired_request(REQUEST_TEST.request) == True

    mock_client.get_latest_block.return_value = {"block": {"header": {"height": 136730}}}
    assert auth.verify_non_expired_request(REQUEST_TEST.request) == False


def test_verify_requested_validator():
    request = REQUEST_TEST.request
    auth = Auth(Client("xxx"))
    assert (
        auth.verify_requested_validator(
            request, "bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw"
        )
        == True
    )
    assert (
        auth.verify_requested_validator(
            request, "bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre"
        )
        == False
    )
    assert (
        auth.verify_requested_validator(
            request, "bandvaloper1yyv5jkqaukq0ajqn7vhkyhpff7h6e99j3gv0tr"
        )
        == True
    )


def test_verify_unsubmitted_report():
    reports = REQUEST_TEST.reports
    auth = Auth(Client("xxx"))
    assert (
        auth.verify_unsubmitted_report(
            reports, "bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw"
        )
        == False
    )
    assert (
        auth.verify_unsubmitted_report(
            reports, "bandvaloper1alzj765pzuhtjkmslme4fdpeakc0036xnyjltn"
        )
        == True
    )

