import base64
import pytest

from pyband.message import MsgCreateDataSource, MsgEditDataSource
from pyband.wallet import Address


def test_msg_create_data_source_success():
    file = open("./tests/files/example_data_source.py", "rb").read()
    executable = base64.b64encode(file).decode()
    msg_create_data_source = MsgCreateDataSource(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Data Source Name",
        description="Data Source Description",
        script="./tests/files/example_data_source.py",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    assert msg_create_data_source.validate() == True

    assert msg_create_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_create_data_source.as_json() == {
        "type": "oracle/CreateDataSource",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "name": "Data Source Name",
            "description": "Data Source Description",
            "executable": executable,
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }

    assert msg_create_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"


def test_msg_create_data_source_empty_path():
    msg_create_data_source = MsgCreateDataSource(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Empty Data Source Name",
        description="Empty Data Source Description",
        script="",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    with pytest.raises(ValueError, match="Missing or invalid data source path"):
        msg_create_data_source.validate()


def test_msg_create_data_source_empty_script():
    msg_create_data_source = MsgCreateDataSource(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Empty Data Source Name",
        description="Empty Data Source Description",
        script="./tests/files/empty_data_source.py",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    with pytest.raises(ValueError, match="Empty data source file"):
        msg_create_data_source.validate()


def test_msg_edit_data_source_full_success():
    file = open("./tests/files/example_data_source.py", "rb").read()
    executable = base64.b64encode(file).decode()
    msg_edit_data_source = MsgEditDataSource(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        name="Data Source Name",
        data_source_id=1,
        description="Data Source Description",
        script="./tests/files/example_data_source.py",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    assert msg_edit_data_source.validate() == True

    assert msg_edit_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_edit_data_source.as_json() == {
        "type": "oracle/EditDataSource",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "name": "Data Source Name",
            "data_source_id": 1,
            "description": "Data Source Description",
            "executable": executable,
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }

    assert msg_edit_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"


def test_msg_edit_data_source_no_name_success():
    file = open("./tests/files/example_data_source.py", "rb").read()
    executable = base64.b64encode(file).decode()
    msg_edit_data_source = MsgEditDataSource(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        data_source_id=1,
        description="Data Source Description",
        script="./tests/files/example_data_source.py",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    assert msg_edit_data_source.validate() == True

    assert msg_edit_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_edit_data_source.as_json() == {
        "type": "oracle/EditDataSource",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "name": "[do-not-modify]",
            "data_source_id": 1,
            "description": "Data Source Description",
            "executable": executable,
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }

    assert msg_edit_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"


def test_msg_edit_data_source_no_description_success():
    file = open("./tests/files/example_data_source.py", "rb").read()
    executable = base64.b64encode(file).decode()
    msg_edit_data_source = MsgEditDataSource(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        data_source_id=1,
        name="Data Source Name",
        script="./tests/files/example_data_source.py",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    assert msg_edit_data_source.validate() == True

    assert msg_edit_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_edit_data_source.as_json() == {
        "type": "oracle/EditDataSource",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "name": "Data Source Name",
            "description": "[do-not-modify]",
            "data_source_id": 1,
            "executable": executable,
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }

    assert msg_edit_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"


def test_msg_edit_data_source_empty_name():
    msg_edit_data_source = MsgEditDataSource(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        data_source_id=1,
        name="",
        description="Empty Data Source Description",
        script="./tests/files/example_data_source.py",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    with pytest.raises(ValueError, match="Invalid data source name"):
        msg_edit_data_source.validate()


def test_msg_edit_data_source_no_script_success():
    msg_edit_data_source = MsgEditDataSource(
        owner=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
        data_source_id=1,
        name="Data Source Name",
        description="Data Source Description",
        sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
    )
    assert msg_edit_data_source.validate() == True

    assert msg_edit_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"

    assert msg_edit_data_source.as_json() == {
        "type": "oracle/EditDataSource",
        "value": {
            "owner": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
            "data_source_id": 1,
            "name": "Data Source Name",
            "description": "Data Source Description",
            "executable": "W2RvLW5vdC1tb2RpZnld",
            "sender": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
        },
    }

    assert msg_edit_data_source.get_sender().to_acc_bech32() == "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"
