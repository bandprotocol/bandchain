import base64
import pytest

from pyband.message import MsgCreateOracleScript, MsgEditOracleScript
from pyband.wallet import Address


def test_msg_create_oracle_script_success():
    file = open("./tests/files/example_oracle_script.wasm", "rb").read()
    code = base64.b64encode(file).decode()
    msg_create_oracle_script = MsgCreateOracleScript(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Oracle Script Name",
        description="Oracle Script Description",
        script="./tests/files/example_oracle_script.wasm",
        schema="{symbols:[string],multiplier:u64}/{rates:[u64]}",
        source_code_url="https://google.com",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    assert msg_create_oracle_script.validate() == True

    assert msg_create_oracle_script.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_create_oracle_script.as_json() == {
        "type": "oracle/CreateOracleScript",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "name": "Oracle Script Name",
            "description": "Oracle Script Description",
            "code": code,
            "schema": "{symbols:[string],multiplier:u64}/{rates:[u64]}",
            "source_code_url": "https://google.com",
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }

    assert msg_create_oracle_script.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"


def test_msg_create_oracle_script_fail_invalid_name():
    file = open("./tests/files/example_oracle_script.wasm", "rb").read()
    code = base64.b64encode(file).decode()
    msg_create_oracle_script = MsgCreateOracleScript(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="",
        description="Oracle Script Description",
        script="./tests/files/example_oracle_script.wasm",
        schema="{symbols:[string],multiplier:u64}/{rates:[u64]}",
        source_code_url="https://google.com",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )

    with pytest.raises(ValueError, match="Missing or invalid oracle script name"):
        msg_create_oracle_script.validate()

    assert msg_create_oracle_script.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_create_oracle_script.as_json() == {
        "type": "oracle/CreateOracleScript",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "name": "",
            "description": "Oracle Script Description",
            "code": code,
            "schema": "{symbols:[string],multiplier:u64}/{rates:[u64]}",
            "source_code_url": "https://google.com",
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }


def test_msg_create_oracle_script_fail_invalid_path():
    msg_create_oracle_script = MsgCreateOracleScript(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Oracle Script Name",
        description="Oracle Script Description",
        script="",
        schema="{symbols:[string],multiplier:u64}/{rates:[u64]}",
        source_code_url="https://google.com",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )

    with pytest.raises(ValueError, match="Missing or invalid oracle script path"):
        msg_create_oracle_script.validate()


def test_msg_create_oracle_script_fail_invalid_wasm():
    file = open("./tests/files/empty_oracle_script.wasm", "rb").read()
    code = base64.b64encode(file).decode()
    msg_create_oracle_script = MsgCreateOracleScript(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Oracle Script Name",
        description="Oracle Script Description",
        script="./tests/files/empty_oracle_script.wasm",
        schema="{symbols:[string],multiplier:u64}/{rates:[u64]}",
        source_code_url="https://google.com",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )

    with pytest.raises(ValueError, match="Empty wasm file"):
        msg_create_oracle_script.validate()

    assert msg_create_oracle_script.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_create_oracle_script.as_json() == {
        "type": "oracle/CreateOracleScript",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "name": "Oracle Script Name",
            "description": "Oracle Script Description",
            "code": code,
            "schema": "{symbols:[string],multiplier:u64}/{rates:[u64]}",
            "source_code_url": "https://google.com",
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }


def test_msg_edit_oracle_script_full_success():
    file = open("./tests/files/example_oracle_script.wasm", "rb").read()
    code = base64.b64encode(file).decode()
    msg_edit_oracle_script = MsgEditOracleScript(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Oracle Script Name",
        oracle_script_id=1,
        description="Oracle Script Description",
        script="./tests/files/example_oracle_script.wasm",
        schema="{symbols:[string],multiplier:u64}/{rates:[u64]}",
        source_code_url="https://google.com",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    assert msg_edit_oracle_script.validate() == True

    assert msg_edit_oracle_script.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_edit_oracle_script.as_json() == {
        "type": "oracle/EditOracleScript",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "oracle_script_id": 1,
            "name": "Oracle Script Name",
            "description": "Oracle Script Description",
            "code": code,
            "schema": "{symbols:[string],multiplier:u64}/{rates:[u64]}",
            "source_code_url": "https://google.com",
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }


def test_msg_edit_oracle_script_empty_name():
    file = open("./tests/files/example_oracle_script.wasm", "rb").read()
    code = base64.b64encode(file).decode()
    msg_edit_oracle_script = MsgEditOracleScript(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="",
        oracle_script_id=1,
        description="Oracle Script Description",
        script="./tests/files/example_oracle_script.wasm",
        schema="{symbols:[string],multiplier:u64}/{rates:[u64]}",
        source_code_url="https://google.com",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )

    with pytest.raises(ValueError, match="Invalid oracle script name"):
        msg_edit_oracle_script.validate()

    assert msg_edit_oracle_script.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_edit_oracle_script.as_json() == {
        "type": "oracle/EditOracleScript",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "oracle_script_id": 1,
            "name": "",
            "description": "Oracle Script Description",
            "code": code,
            "schema": "{symbols:[string],multiplier:u64}/{rates:[u64]}",
            "source_code_url": "https://google.com",
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }


def test_msg_edit_oracle_script_empty_script_path():
    msg_edit_oracle_script = MsgEditOracleScript(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Oracle Script Name",
        oracle_script_id=1,
        description="Oracle Script Description",
        script="",
        schema="{symbols:[string],multiplier:u64}/{rates:[u64]}",
        source_code_url="https://google.com",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )

    with pytest.raises(ValueError, match="Invalid oracle script path"):
        msg_edit_oracle_script.validate()

    assert msg_edit_oracle_script.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"


def test_msg_edit_oracle_script_empty_script_path():
    file = open("./tests/files/empty_oracle_script.wasm", "rb").read()
    code = base64.b64encode(file).decode()
    msg_edit_oracle_script = MsgEditOracleScript(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Oracle Script Name",
        oracle_script_id=1,
        description="Oracle Script Description",
        script="./tests/files/empty_oracle_script.wasm",
        schema="{symbols:[string],multiplier:u64}/{rates:[u64]}",
        source_code_url="https://google.com",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )

    with pytest.raises(ValueError, match="Empty wasm file"):
        msg_edit_oracle_script.validate()

    assert msg_edit_oracle_script.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_edit_oracle_script.as_json() == {
        "type": "oracle/EditOracleScript",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "oracle_script_id": 1,
            "name": "Oracle Script Name",
            "description": "Oracle Script Description",
            "code": code,
            "schema": "{symbols:[string],multiplier:u64}/{rates:[u64]}",
            "source_code_url": "https://google.com",
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }
