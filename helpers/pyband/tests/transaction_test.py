import requests
from pyband.client import Client
import pytest

from pyband import Transaction
from pyband.message import MsgRequest
from pyband.wallet import Address, PrivateKey
from pyband.client import Client
from pyband.exceptions import EmptyMsgError, UndefinedError, ValueTooLargeError

TEST_RPC = "https://api-mock.bandprotocol.com/rest/"

client = Client(TEST_RPC)


def test_get_sign_data_success():
    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
            )
        )
        .with_account_num(100)
        .with_sequence(30)
        .with_chain_id("bandchain")
    )

    assert (
        t.get_sign_data()
        == b'{"account_number":"100","chain_id":"bandchain","fee":{"amount":[{"amount":"0","denom":"uband"}],"gas":"200000"},"memo":"","msgs":[{"type":"oracle/Request","value":{"ask_count":"4","calldata":"AAAAA0JUQwAAAAAAAAAB","client_id":"from_pyband","min_count":"3","oracle_script_id":"1","sender":"band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"}}],"sequence":"30"}'
    )


def test_get_sign_data_with_gas_success():
    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
            )
        )
        .with_account_num(100)
        .with_sequence(30)
        .with_chain_id("bandchain")
        .with_gas(500000)
        .with_fee(10)
    )

    assert (
        t.get_sign_data()
        == b'{"account_number":"100","chain_id":"bandchain","fee":{"amount":[{"amount":"10","denom":"uband"}],"gas":"500000"},"memo":"","msgs":[{"type":"oracle/Request","value":{"ask_count":"4","calldata":"AAAAA0JUQwAAAAAAAAAB","client_id":"from_pyband","min_count":"3","oracle_script_id":"1","sender":"band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"}}],"sequence":"30"}'
    )


def test_get_sign_data_with_auto_success(requests_mock):

    requests_mock.register_uri(
        "GET",
        "{}/auth/accounts/band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c".format(TEST_RPC),
        json={
            "height": "650788",
            "result": {
                "type": "cosmos-sdk/Account",
                "value": {
                    "address": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
                    "coins": [{"denom": "uband", "amount": "104082359107"}],
                    "public_key": {
                        "type": "tendermint/PubKeySecp256k1",
                        "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
                    },
                    "account_number": "36",
                    "sequence": "927",
                },
            },
        },
        status_code=200,
    )

    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
            )
        )
        .with_auto(client)
        .with_chain_id("bandchain")
        .with_gas(500000)
        .with_fee(10)
    )

    assert (
        t.get_sign_data()
        == b'{"account_number":"36","chain_id":"bandchain","fee":{"amount":[{"amount":"10","denom":"uband"}],"gas":"500000"},"memo":"","msgs":[{"type":"oracle/Request","value":{"ask_count":"4","calldata":"AAAAA0JUQwAAAAAAAAAB","client_id":"from_pyband","min_count":"3","oracle_script_id":"1","sender":"band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"}}],"sequence":"927"}'
    )


def test_create_transaction_with_auto_fail(requests_mock):

    requests_mock.register_uri(
        "GET",
        "{}/auth/accounts/band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c".format(TEST_RPC),
        json={
            "height": "650788",
            "result": {
                "type": "cosmos-sdk/Account",
                "value": {
                    "address": "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c",
                    "coins": [{"denom": "uband", "amount": "104082359107"}],
                    "public_key": {
                        "type": "tendermint/PubKeySecp256k1",
                        "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
                    },
                    "account_number": "36",
                    "sequence": "927",
                },
            },
        },
        status_code=200,
    )

    with pytest.raises(EmptyMsgError, match="messsage is empty, please use with_messages at least 1 message"):
        t = Transaction().with_auto(client)


def test_get_sign_data_with_memo_success():
    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
            )
        )
        .with_account_num(100)
        .with_sequence(30)
        .with_chain_id("bandchain")
        .with_memo("from_pyband_test")
    )

    assert (
        t.get_sign_data()
        == b'{"account_number":"100","chain_id":"bandchain","fee":{"amount":[{"amount":"0","denom":"uband"}],"gas":"200000"},"memo":"from_pyband_test","msgs":[{"type":"oracle/Request","value":{"ask_count":"4","calldata":"AAAAA0JUQwAAAAAAAAAB","client_id":"from_pyband","min_count":"3","oracle_script_id":"1","sender":"band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"}}],"sequence":"30"}'
    )


def test_get_sign_data_msgs_fail():
    t = Transaction().with_account_num(100).with_sequence(30).with_chain_id("bandchain")

    with pytest.raises(EmptyMsgError, match="message is empty"):
        t.get_sign_data()


def test_get_sign_data_account_num_fail():
    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
            )
        )
        .with_sequence(30)
        .with_chain_id("bandchain")
    )

    with pytest.raises(UndefinedError, match="account_num should be defined"):
        t.get_sign_data()


def test_get_sign_data_sequence_fail():
    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
            )
        )
        .with_account_num(100)
        .with_chain_id("bandchain")
    )

    with pytest.raises(UndefinedError, match="sequence should be defined"):
        t.get_sign_data()


def test_get_sign_data_chain_id_fail():
    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
            )
        )
        .with_account_num(100)
        .with_sequence(30)
    )

    with pytest.raises(UndefinedError, match="chain_id should be defined"):
        t.get_sign_data()


def test_get_sign_data_memo_fail():
    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=Address.from_acc_bech32("band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"),
            )
        )
        .with_account_num(100)
        .with_sequence(30)
        .with_chain_id("bandchain")
    )

    with pytest.raises(ValueTooLargeError, match="memo is too large"):
        t.with_memo(
            "This is the longest memo in the world. This is the longest memo in the world. This is the longest memo in the world. This is the longest memo in the world. This is the longest memo in the world. This is the longest memo in the world. This is the longest memo in the world.This is the longest memo in the world. This is the longest memo in the world.This is the longest memo in the world."
        )


def test_get_tx_data_success():
    priv = PrivateKey.from_mnemonic("s")
    pubkey = priv.to_pubkey()
    addr = pubkey.to_address()

    t = (
        Transaction()
        .with_messages(
            MsgRequest(
                oracle_script_id=1,
                calldata=bytes.fromhex("000000034254430000000000000001"),
                ask_count=4,
                min_count=3,
                client_id="from_pyband",
                sender=addr,
            )
        )
        .with_account_num(100)
        .with_sequence(30)
        .with_chain_id("bandchain")
    )

    raw_data = t.get_sign_data()
    signature = priv.sign(raw_data)
    raw_tx = t.get_tx_data(signature, pubkey)

    assert raw_tx == {
        "msg": [
            {
                "type": "oracle/Request",
                "value": {
                    "oracle_script_id": "1",
                    "calldata": "AAAAA0JUQwAAAAAAAAAB",
                    "ask_count": "4",
                    "min_count": "3",
                    "client_id": "from_pyband",
                    "sender": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
                },
            }
        ],
        "fee": {"gas": "200000", "amount": [{"denom": "uband", "amount": "0"}]},
        "memo": "",
        "signatures": [
            {
                "signature": "aFvTgkY9F14dHfp2mtq8V2eCTKUtY1T9WKs99jmR8To3JB16cawbmpa1TRUdzfnqLXBh+o6XUuF4bHWR6xbCCw==",
                "pub_key": {
                    "type": "tendermint/PubKeySecp256k1",
                    "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
                },
                "account_number": "100",
                "sequence": "30",
            }
        ],
    }
