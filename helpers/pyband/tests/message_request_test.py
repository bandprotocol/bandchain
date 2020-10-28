import base64
import pytest

from pyband.message import MsgRequest
from pyband.wallet import Address


def test_msg_request_creation_success():
    msg_request = MsgRequest(
        oracle_script_id=1,
        calldata=bytes.fromhex("000000034254430000000000000001"),
        ask_count=4,
        min_count=3,
        client_id="from_pyband",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    assert msg_request.validate() == True

    assert msg_request.as_json() == {
        "type": "oracle/Request",
        "value": {
            "oracle_script_id": "1",
            "calldata": "AAAAA0JUQwAAAAAAAAAB",
            "ask_count": "4",
            "min_count": "3",
            "client_id": "from_pyband",
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }


def test_msg_request_get_sender():
    msg_request = MsgRequest(
        oracle_script_id=1,
        calldata=bytes.fromhex("000000034254430000000000000001"),
        ask_count=4,
        min_count=3,
        client_id="from_pyband",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )

    assert (
        msg_request.get_sender().to_acc_bech32()
        == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"
    )


def test_msg_request_creation_oracle_script_id_fail():
    msg_request = MsgRequest(
        oracle_script_id=0,
        calldata=bytes.fromhex("000000034254430000000000000001"),
        ask_count=4,
        min_count=3,
        client_id="from_pyband",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    with pytest.raises(ValueError, match="oracle script id cannot less than zero"):
        msg_request.validate()


def test_msg_request_creation_calldata_fail():
    msg_request = MsgRequest(
        oracle_script_id=1,
        calldata=bytes.fromhex(
            "000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001000000034254430000000000000001"
        ),
        ask_count=4,
        min_count=3,
        client_id="from_pyband",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    with pytest.raises(ValueError, match="too large calldata"):
        msg_request.validate()


def test_msg_request_creation_ask_count_fail():
    msg_request = MsgRequest(
        oracle_script_id=1,
        calldata=bytes.fromhex("000000034254430000000000000001"),
        ask_count=3,
        min_count=4,
        client_id="from_pyband",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    with pytest.raises(
        ValueError, match="invalid ask count got: min count: 4, ask count: 3"
    ):
        msg_request.validate()


def test_msg_request_creation_min_count_fail():
    msg_request = MsgRequest(
        oracle_script_id=1,
        calldata=bytes.fromhex("000000034254430000000000000001"),
        ask_count=3,
        min_count=0,
        client_id="from_pyband",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    with pytest.raises(ValueError, match="invalid min count got: min count: 0"):
        msg_request.validate()


def test_msg_request_creation_client_id_count_fail():
    msg_request = MsgRequest(
        oracle_script_id=1,
        calldata=bytes.fromhex("000000034254430000000000000001"),
        ask_count=3,
        min_count=2,
        client_id="Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    with pytest.raises(ValueError, match="too long client id"):
        msg_request.validate()
