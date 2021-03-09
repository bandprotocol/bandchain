import pytest

from pyband.client import Client
from pyband.data import (
    TransactionSyncMode,
    TransactionAsyncMode,
    TransactionBlockMode,
    HexBytes,
)
from requests.exceptions import ReadTimeout
from unittest.mock import patch

TEST_RPC = "https://api-mock.bandprotocol.com/rest"
TEST_MSG = {
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
    "fee": {"gas": "1000000", "amount": [{"denom": "uband", "amount": "0"}]},
    "memo": "Memo",
    "signatures": [
        {
            "signature": "hQpMSSaOVbT5vd3yladNX9RNA9vSq4ts4cPufdoesjUtPje5i73f048MM0xPnAB7JWSRuUSsZD5M6L6WGk3Qkw==",
            "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
            },
            "account_number": "36",
            "sequence": "1092",
        }
    ],
}

TEST_WRONG_SEQUENCE_MSG = {
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
    "fee": {"gas": "1000000", "amount": [{"denom": "uband", "amount": "0"}]},
    "memo": "Memo",
    "signatures": [
        {
            "signature": "hQpMSSaOVbT5vd3yladNX9RNA9vSq4ts4cPufdoesjUtPje5i73f048MM0xPnAB7JWSRuUSsZD5M6L6WGk3Qkw==",
            "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
            },
            "account_number": "0",
            "sequence": "1092",
        }
    ],
}

TIMEOUT = 3

client = Client(TEST_RPC, TIMEOUT)


@patch("requests.post")
def test_send_tx_sync_mode_timeout(requests_mock):
    requests_mock.side_effect = ReadTimeout
    with pytest.raises(ReadTimeout):
        res = client.send_tx_sync_mode(TEST_MSG)


@patch("requests.post")
def test_send_tx_block_mode_timeout(requests_mock):
    requests_mock.side_effect = ReadTimeout
    with pytest.raises(ReadTimeout):
        res = client.send_tx_block_mode(TEST_MSG)


@patch("requests.post")
def test_send_tx_async_mode_timeout(requests_mock):
    requests_mock.side_effect = ReadTimeout
    with pytest.raises(ReadTimeout):
        res = client.send_tx_async_mode(TEST_MSG)


def test_send_tx_sync_mode_success(requests_mock):
    requests_mock.register_uri(
        "POST",
        "{}/txs".format(TEST_RPC),
        json={
            "height": "0",
            "txhash": "E204AAD58ACA8F00942B1BB66D9F745F5E2C21E04C5DF6A0CB73DF02B6B51121",
            "raw_log": "[]",
        },
        status_code=200,
    )

    assert client.send_tx_sync_mode(TEST_MSG) == TransactionSyncMode(
        tx_hash=HexBytes(bytes.fromhex("E204AAD58ACA8F00942B1BB66D9F745F5E2C21E04C5DF6A0CB73DF02B6B51121")),
        code=0,
        error_log=None,
    )


def test_send_tx_sync_mode_wrong_sequence_fail(requests_mock):
    requests_mock.register_uri(
        "POST",
        "{}/txs".format(TEST_RPC),
        json={
            "height": "0",
            "txhash": "611F45A21BB7937E451CDE78D124218603644635CC40A97D2BC1E854CED8D6E6",
            "code": 4,
            "raw_log": "unauthorized: signature verification failed; verify correct account sequence and chain-id",
        },
        status_code=200,
    )
    assert client.send_tx_sync_mode(TEST_WRONG_SEQUENCE_MSG) == TransactionSyncMode(
        tx_hash=HexBytes(bytes.fromhex("611F45A21BB7937E451CDE78D124218603644635CC40A97D2BC1E854CED8D6E6")),
        code=4,
        error_log="unauthorized: signature verification failed; verify correct account sequence and chain-id",
    )


def test_send_tx_sync_mode_return_only_code(requests_mock):
    requests_mock.register_uri(
        "POST",
        "{}/txs".format(TEST_RPC),
        json={
            "height": "0",
            "txhash": "611F45A21BB7937E451CDE78D124218603644635CC40A97D2BC1E854CED8D6E6",
            "code": 19,
        },
        status_code=200,
    )
    assert client.send_tx_sync_mode(TEST_WRONG_SEQUENCE_MSG) == TransactionSyncMode(
        tx_hash=HexBytes(bytes.fromhex("611F45A21BB7937E451CDE78D124218603644635CC40A97D2BC1E854CED8D6E6")),
        code=19,
        error_log=None,
    )


def test_send_tx_async_mode_success(requests_mock):
    requests_mock.register_uri(
        "POST",
        "{}/txs".format(TEST_RPC),
        json={
            "height": "0",
            "txhash": "E204AAD58ACA8F00942B1BB66D9F745F5E2C21E04C5DF6A0CB73DF02B6B51121",
            "raw_log": "[]",
        },
        status_code=200,
    )

    assert client.send_tx_async_mode(TEST_MSG) == TransactionAsyncMode(
        tx_hash=HexBytes(bytes.fromhex("E204AAD58ACA8F00942B1BB66D9F745F5E2C21E04C5DF6A0CB73DF02B6B51121"))
    )


def test_send_tx_block_mode_success(requests_mock):
    requests_mock.register_uri(
        "POST",
        "{}/txs".format(TEST_RPC),
        json={
            "height": "715786",
            "txhash": "2DE264D16164BCCF695E960553FED537EDC00D0E3EDF69D6BFE4168C476AD03C",
            "raw_log": '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"51"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"4"},{"key":"min_count","value":"3"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr"},{"key":"validator","value":"bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},{"key":"validator","value":"bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]}]',
            "logs": [
                {
                    "msg_index": 0,
                    "log": "",
                    "events": [
                        {
                            "type": "message",
                            "attributes": [{"key": "action", "value": "request"}],
                        },
                        {
                            "type": "raw_request",
                            "attributes": [
                                {"key": "data_source_id", "value": "1"},
                                {
                                    "key": "data_source_hash",
                                    "value": "c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0",
                                },
                                {"key": "external_id", "value": "1"},
                                {"key": "calldata", "value": "BTC"},
                                {"key": "data_source_id", "value": "2"},
                                {
                                    "key": "data_source_hash",
                                    "value": "dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479",
                                },
                                {"key": "external_id", "value": "2"},
                                {"key": "calldata", "value": "BTC"},
                                {"key": "data_source_id", "value": "3"},
                                {
                                    "key": "data_source_hash",
                                    "value": "f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921",
                                },
                                {"key": "external_id", "value": "3"},
                                {"key": "calldata", "value": "BTC"},
                            ],
                        },
                        {
                            "type": "request",
                            "attributes": [
                                {"key": "id", "value": "51"},
                                {"key": "client_id", "value": "from_pyband"},
                                {"key": "oracle_script_id", "value": "1"},
                                {
                                    "key": "calldata",
                                    "value": "000000034254430000000000000001",
                                },
                                {"key": "ask_count", "value": "4"},
                                {"key": "min_count", "value": "3"},
                                {"key": "gas_used", "value": "2405"},
                                {
                                    "key": "validator",
                                    "value": "bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr",
                                },
                                {
                                    "key": "validator",
                                    "value": "bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh",
                                },
                                {
                                    "key": "validator",
                                    "value": "bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst",
                                },
                                {
                                    "key": "validator",
                                    "value": "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec",
                                },
                            ],
                        },
                    ],
                }
            ],
            "gas_wanted": "1000000",
            "gas_used": "343054",
        },
        status_code=200,
    )

    assert client.send_tx_block_mode(TEST_MSG) == TransactionBlockMode(
        height=715786,
        tx_hash=HexBytes(bytes.fromhex("2DE264D16164BCCF695E960553FED537EDC00D0E3EDF69D6BFE4168C476AD03C")),
        gas_wanted=1000000,
        gas_used=1000000,
        log=[
            {
                "msg_index": 0,
                "log": "",
                "events": [
                    {
                        "type": "message",
                        "attributes": [{"key": "action", "value": "request"}],
                    },
                    {
                        "type": "raw_request",
                        "attributes": [
                            {"key": "data_source_id", "value": "1"},
                            {
                                "key": "data_source_hash",
                                "value": "c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0",
                            },
                            {"key": "external_id", "value": "1"},
                            {"key": "calldata", "value": "BTC"},
                            {"key": "data_source_id", "value": "2"},
                            {
                                "key": "data_source_hash",
                                "value": "dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479",
                            },
                            {"key": "external_id", "value": "2"},
                            {"key": "calldata", "value": "BTC"},
                            {"key": "data_source_id", "value": "3"},
                            {
                                "key": "data_source_hash",
                                "value": "f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921",
                            },
                            {"key": "external_id", "value": "3"},
                            {"key": "calldata", "value": "BTC"},
                        ],
                    },
                    {
                        "type": "request",
                        "attributes": [
                            {"key": "id", "value": "51"},
                            {"key": "client_id", "value": "from_pyband"},
                            {"key": "oracle_script_id", "value": "1"},
                            {
                                "key": "calldata",
                                "value": "000000034254430000000000000001",
                            },
                            {"key": "ask_count", "value": "4"},
                            {"key": "min_count", "value": "3"},
                            {"key": "gas_used", "value": "2405"},
                            {
                                "key": "validator",
                                "value": "bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr",
                            },
                            {
                                "key": "validator",
                                "value": "bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh",
                            },
                            {
                                "key": "validator",
                                "value": "bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst",
                            },
                            {
                                "key": "validator",
                                "value": "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec",
                            },
                        ],
                    },
                ],
            }
        ],
        error_log=None,
        code=0,
    )


def test_send_tx_block_wrong_sequence_fail(requests_mock):
    requests_mock.register_uri(
        "POST",
        "{}/txs".format(TEST_RPC),
        json={
            "height": "0",
            "txhash": "7F1CFFD674CAAEEB25922E9C6E9F8F8CEF7A325E25B64E8DAB070D7409FD1F72",
            "codespace": "sdk",
            "code": 4,
            "raw_log": "unauthorized: signature verification failed; verify correct account sequence and chain-id",
            "gas_wanted": "1000000",
            "gas_used": "27402",
        },
        status_code=200,
    )

    assert client.send_tx_block_mode(TEST_MSG) == TransactionBlockMode(
        height=0,
        tx_hash=HexBytes(bytes.fromhex("7F1CFFD674CAAEEB25922E9C6E9F8F8CEF7A325E25B64E8DAB070D7409FD1F72")),
        gas_wanted=1000000,
        gas_used=1000000,
        code=4,
        log=[],
        error_log="unauthorized: signature verification failed; verify correct account sequence and chain-id",
    )


def test_send_tx_block_return_code(requests_mock):
    requests_mock.register_uri(
        "POST",
        "{}/txs".format(TEST_RPC),
        json={
            "height": "0",
            "txhash": "7F1CFFD674CAAEEB25922E9C6E9F8F8CEF7A325E25B64E8DAB070D7409FD1F72",
            "codespace": "sdk",
            "code": 19,
            "gas_wanted": "1000000",
            "gas_used": "27402",
        },
        status_code=200,
    )

    assert client.send_tx_block_mode(TEST_MSG) == TransactionBlockMode(
        height=0,
        tx_hash=HexBytes(bytes.fromhex("7F1CFFD674CAAEEB25922E9C6E9F8F8CEF7A325E25B64E8DAB070D7409FD1F72")),
        gas_wanted=1000000,
        gas_used=1000000,
        code=19,
        log=[],
        error_log=None,
    )
