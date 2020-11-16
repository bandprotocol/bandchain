import pytest

from pyband.wallet import Address
from pyband.client import Client
from pyband.data import (
    Account,
    Block,
    BlockHeader,
    BlockHeaderInfo,
    BlockID,
    DataSource,
    OracleScript,
    RequestInfo,
    Request,
    RawRequest,
    RawReport,
    Result,
    Report,
    RequestPacketData,
    ResponsePacketData,
    HexBytes,
    Timestamp,
)
from pyband.utils import parse_datetime

TEST_RPC = "https://api-mock.bandprotocol.com/rest"

client = Client(TEST_RPC)


def test_get_chain_id(requests_mock):

    requests_mock.register_uri(
        "GET",
        "{}/bandchain/chain_id".format(TEST_RPC),
        json={"chain_id": "bandchain"},
        status_code=200,
    )

    assert client.get_chain_id() == "bandchain"


def test_get_latest_block(requests_mock):
    requests_mock.register_uri(
        "GET",
        "{}/blocks/latest".format(TEST_RPC),
        json={
            "block_id": {
                "hash": "391E99908373F8590C928E0619956DA3D87EB654445DA4F25A185C9718561D53",
                "parts": {
                    "total": "1",
                    "hash": "9DC1DB39B7DDB97DE353DFB2898198BAADEFB7DF8090BF22678714F769D69F42",
                },
            },
            "block": {
                "header": {
                    "version": {"block": "10", "app": "0"},
                    "chain_id": "bandchain",
                    "height": "1032007",
                    "time": "2020-11-05T09:15:18.445494105Z",
                    "last_block_id": {
                        "hash": "4BC01E257662B5F9337D615D06D5D19D8D79F7BA5A4022A85B4DC4ED3C59F7CA",
                        "parts": {
                            "total": "1",
                            "hash": "6471C0A51FB7043171EAA76CAFA900B36A4494F47F975980732529D247E8EBA5",
                        },
                    },
                    "last_commit_hash": "17B2CE4ABA910E85847537F1323DB95C9F16C20C60E9B9BBB04C633C3125BD92",
                    "data_hash": "EFE5E3F549554FEE8EB9B393740C250D74580427A96A175ABB105806039CFFE2",
                    "validators_hash": "E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151",
                    "next_validators_hash": "E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151",
                    "consensus_hash": "0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45",
                    "app_hash": "6E2B1ECE9D912D86C25182E8B7419583ABCE978BFC66DC2556BB0D06A8D528EF",
                    "last_results_hash": "",
                    "evidence_hash": "",
                    "proposer_address": "BDB6A0728C8DFE2124536F16F2BA428FE767A8F9",
                },
                "data": {
                    "txs": [
                        "yAEoKBapCj5CcI40CAESDwAAAANCVEMAAAAAAAAAARgEIAMqC2Zyb21fcHliYW5kMhSQ78AMmxLrubEOPhhIwKK5oyk9oBIQCgoKBXViYW5kEgEwEMCEPRpqCibrWumHIQP+cIvaZlJP0sa86QaC44VVqFHgPSruT2KbBd6Q9R7ZvBJANbPqLRAgwwULWWwb5O2/eb6ddtDr65kRFgDcOTTGIqQS5Iz1NvHWHfkPKRoM8egErMsgE9YnuE+pAqoc/YvNfiIEVEVTVA=="
                    ]
                },
                "evidence": {"evidence": None},
                "last_commit": {
                    "height": "1032006",
                    "round": "0",
                    "block_id": {
                        "hash": "4BC01E257662B5F9337D615D06D5D19D8D79F7BA5A4022A85B4DC4ED3C59F7CA",
                        "parts": {
                            "total": "1",
                            "hash": "6471C0A51FB7043171EAA76CAFA900B36A4494F47F975980732529D247E8EBA5",
                        },
                    },
                    "signatures": [
                        {
                            "block_id_flag": 3,
                            "validator_address": "5179B0BB203248E03D2A1342896133B5C58E1E44",
                            "timestamp": "2020-11-05T09:15:18.53815896Z",
                            "signature": "TZY24CKwZOE8wqfE0NM3qzkQ7qCpCrGEHNZdf8n31L4otZzbKGfOL05kGtBsGkTnZkVv7aJmrJ7XbvIzv0SREQ==",
                        },
                        {
                            "block_id_flag": 2,
                            "validator_address": "BDB6A0728C8DFE2124536F16F2BA428FE767A8F9",
                            "timestamp": "2020-11-05T09:15:18.445494105Z",
                            "signature": "mcUMQtCR38MK69IeUDri0zkfllsXKgnVFTsFwNaO/7cnBaIUUz9U4d3H9EjSH4kANJxWRFO3dSnPy1uOD36K6A==",
                        },
                        {
                            "block_id_flag": 3,
                            "validator_address": "F0C23921727D869745C4F9703CF33996B1D2B715",
                            "timestamp": "2020-11-05T09:15:18.537783045Z",
                            "signature": "fpr26xz+Gg5Rl7Fvx34a0QZpb5yJc5P4t5Z1OctIDQ0VMmh9vEWagsqQGErt1bm+CaKFtkFfZZ4CU0DKN27GbQ==",
                        },
                        {
                            "block_id_flag": 3,
                            "validator_address": "F23391B5DBF982E37FB7DADEA64AAE21CAE4C172",
                            "timestamp": "2020-11-05T09:15:18.539946947Z",
                            "signature": "KGsiIaralMMr1M9A7Ulhbc0THt1pLgNIrNQ6Kx+oGtwjl2w5ke5iivAAtZMduhyIAUMhrp6PN7ZvKgv9TPCNNw==",
                        },
                    ],
                },
            },
        },
        status_code=200,
    )

    assert client.get_latest_block() == Block(
        block=BlockHeader(
            header=BlockHeaderInfo(
                chain_id="bandchain",
                height=1032007,
                time=Timestamp(parse_datetime("2020-11-05T09:15:18.445494105Z")),
                last_commit_hash=HexBytes(
                    bytes.fromhex("17B2CE4ABA910E85847537F1323DB95C9F16C20C60E9B9BBB04C633C3125BD92")
                ),
                data_hash=HexBytes(bytes.fromhex("EFE5E3F549554FEE8EB9B393740C250D74580427A96A175ABB105806039CFFE2")),
                validators_hash=HexBytes(
                    bytes.fromhex("E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151")
                ),
                next_validators_hash=HexBytes(
                    bytes.fromhex("E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151")
                ),
                consensus_hash=HexBytes(
                    bytes.fromhex("0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45")
                ),
                app_hash=HexBytes(bytes.fromhex("6E2B1ECE9D912D86C25182E8B7419583ABCE978BFC66DC2556BB0D06A8D528EF")),
                last_results_hash=HexBytes(bytes.fromhex("")),
                evidence_hash=HexBytes(bytes.fromhex("")),
                proposer_address=HexBytes(bytes.fromhex("BDB6A0728C8DFE2124536F16F2BA428FE767A8F9")),
            )
        ),
        block_id=BlockID(
            hash=HexBytes(bytes.fromhex("391E99908373F8590C928E0619956DA3D87EB654445DA4F25A185C9718561D53"))
        ),
    )


def test_get_account(requests_mock):
    requests_mock.register_uri(
        "GET",
        "{}/auth/accounts/band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte".format(TEST_RPC),
        json={
            "height": "650788",
            "result": {
                "type": "cosmos-sdk/Account",
                "value": {
                    "address": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
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

    assert client.get_account(Address.from_acc_bech32("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte")) == Account(
        address="band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
        coins=[{"denom": "uband", "amount": "104082359107"}],
        public_key={
            "type": "tendermint/PubKeySecp256k1",
            "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
        },
        account_number=36,
        sequence=927,
    )


def test_get_data_source(requests_mock):

    requests_mock.register_uri(
        "GET",
        "{}/oracle/data_sources/1".format(TEST_RPC),
        json={
            "height": "651093",
            "result": {
                "owner": "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
                "name": "CoinGecko Cryptocurrency Price",
                "description": "Retrieves current price of a cryptocurrency from https://www.coingecko.com",
                "filename": "c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0",
            },
        },
        status_code=200,
    )

    assert client.get_data_source(1) == DataSource(
        owner="band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
        name="CoinGecko Cryptocurrency Price",
        description="Retrieves current price of a cryptocurrency from https://www.coingecko.com",
        filename="c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0",
    )


def test_get_oracle_script(requests_mock):

    requests_mock.register_uri(
        "GET",
        "{}/oracle/oracle_scripts/1".format(TEST_RPC),
        json={
            "height": "651338",
            "result": {
                "owner": "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
                "name": "Cryptocurrency Price in USD",
                "description": "Oracle script that queries the average cryptocurrency price using current price data from CoinGecko, CryptoCompare, and Binance",
                "filename": "a1f941e828bd8d5ea9c98e2cd3ff9ba8e52a8f63dca45bddbb2fdbfffebc7556",
                "schema": "{symbol:string,multiplier:u64}/{px:u64}",
                "source_code_url": "https://ipfs.io/ipfs/QmQqxHLszpbCy8Hk2ame3pPAxUUAyStBrVdGdDgrfAngAv",
            },
        },
        status_code=200,
    )

    assert client.get_oracle_script(1) == OracleScript(
        owner="band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
        name="Cryptocurrency Price in USD",
        description="Oracle script that queries the average cryptocurrency price using current price data from CoinGecko, CryptoCompare, and Binance",
        filename="a1f941e828bd8d5ea9c98e2cd3ff9ba8e52a8f63dca45bddbb2fdbfffebc7556",
        schema="{symbol:string,multiplier:u64}/{px:u64}",
        source_code_url="https://ipfs.io/ipfs/QmQqxHLszpbCy8Hk2ame3pPAxUUAyStBrVdGdDgrfAngAv",
    )


def test_get_request_by_id(requests_mock):

    requests_mock.register_uri(
        "GET",
        "{}/oracle/requests/1".format(TEST_RPC),
        json={
            "height": "651397",
            "result": {
                "request": {
                    "oracle_script_id": "1",
                    "calldata": "AAAAA0JUQwAAAAAAAABk",
                    "requested_validators": ["bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"],
                    "min_count": "1",
                    "request_height": "118",
                    "request_time": "2020-10-14T09:47:23.818758882Z",
                    "client_id": "from_scan",
                    "raw_requests": [
                        {"external_id": "1", "data_source_id": "1", "calldata": "QlRD"},
                        {"external_id": "2", "data_source_id": "2", "calldata": "QlRD"},
                        {"external_id": "3", "data_source_id": "3", "calldata": "QlRD"},
                    ],
                },
                "reports": [
                    {
                        "validator": "bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst",
                        "in_before_resolve": True,
                        "raw_reports": [
                            {"external_id": "3", "data": "MTE0MTQuNTQ1Cg=="},
                            {"external_id": "1", "data": "MTE0MjAuNjEK"},
                            {"external_id": "2", "data": "MTE0MjQuMTgK"},
                        ],
                    }
                ],
                "result": {
                    "request_packet_data": {
                        "client_id": "from_scan",
                        "oracle_script_id": "1",
                        "calldata": "AAAAA0JUQwAAAAAAAABk",
                        "ask_count": "1",
                        "min_count": "1",
                    },
                    "response_packet_data": {
                        "client_id": "from_scan",
                        "request_id": "1",
                        "ans_count": "1",
                        "request_time": "1602668843",
                        "resolve_time": "1602668845",
                        "resolve_status": 1,
                        "result": "AAAAAAARbNk=",
                    },
                },
            },
        },
        status_code=200,
    )

    assert client.get_request_by_id(1) == RequestInfo(
        request=Request(
            oracle_script_id=1,
            requested_validators=["bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"],
            min_count=1,
            request_height=118,
            raw_requests=[
                RawRequest(data_source_id=1, external_id=1, calldata=b"BTC"),
                RawRequest(data_source_id=2, external_id=2, calldata=b"BTC"),
                RawRequest(data_source_id=3, external_id=3, calldata=b"BTC"),
            ],
            client_id="from_scan",
            calldata=b"\x00\x00\x00\x03BTC\x00\x00\x00\x00\x00\x00\x00d",
        ),
        reports=[
            Report(
                validator="bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst",
                raw_reports=[
                    RawReport(external_id=3, data=b"11414.545\n"),
                    RawReport(external_id=1, data=b"11420.61\n"),
                    RawReport(external_id=2, data=b"11424.18\n"),
                ],
                in_before_resolve=True,
            )
        ],
        result=Result(
            request_packet_data=RequestPacketData(
                oracle_script_id=1,
                ask_count=1,
                min_count=1,
                client_id="from_scan",
                calldata=b"\x00\x00\x00\x03BTC\x00\x00\x00\x00\x00\x00\x00d",
            ),
            response_packet_data=ResponsePacketData(
                request_id=1,
                request_time=1602668843,
                resolve_time=1602668845,
                resolve_status=1,
                ans_count=1,
                client_id="from_scan",
                result=b"\x00\x00\x00\x00\x00\x11l\xd9",
            ),
        ),
    )


def test_get_latest_request(requests_mock):

    requests_mock.register_uri(
        "GET",
        "{}/oracle/request_search?oid=8&ask_count=4&min_count=3&calldata=000000190000000652454e4254430000000457425443000000034449410000000342544d00000004494f545800000003464554000000034a5354000000034d434f000000034b4d440000000342545300000003514b430000000559414d563200000003585a4300000003554f5300000004414b524f00000003484e5400000003484f54000000034b4149000000034f474e00000003575258000000034b4441000000034f524e00000003464f52000000034153540000000553544f524a000000003b9aca00".format(
            TEST_RPC
        ),
        json={
            "height": "2243612",
            "result": {
                "request": {
                    "oracle_script_id": "8",
                    "calldata": "AAAAGQAAAAZSRU5CVEMAAAAEV0JUQwAAAANESUEAAAADQlRNAAAABElPVFgAAAADRkVUAAAAA0pTVAAAAANNQ08AAAADS01EAAAAA0JUUwAAAANRS0MAAAAFWUFNVjIAAAADWFpDAAAAA1VPUwAAAARBS1JPAAAAA0hOVAAAAANIT1QAAAADS0FJAAAAA09HTgAAAANXUlgAAAADS0RBAAAAA09STgAAAANGT1IAAAADQVNUAAAABVNUT1JKAAAAADuaygA=",
                    "requested_validators": [
                        "bandvaloper1xnryftxluq49fk52c5j5zrxcc5rzye96s70msl",
                        "bandvaloper1v38hewjc0865dm4t89v5efh9rmum5rmrm7evg4",
                        "bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw",
                        "bandvaloper1yyv5jkqaukq0ajqn7vhkyhpff7h6e99j3gv0tr",
                    ],
                    "min_count": "3",
                    "request_height": "2243579",
                    "request_time": "2020-10-22T08:12:00.558933096Z",
                    "client_id": "bandteam",
                    "raw_requests": [
                        {
                            "external_id": "2",
                            "data_source_id": "15",
                            "calldata": "UkVOQlRDIFdCVEMgRElBIEJUTSBJT1RYIEZFVCBKU1QgTUNPIEtNRCBCVFMgUUtDIFlBTVYyIFhaQyBVT1MgQUtSTyBITlQgSE9UIEtBSSBPR04gV1JYIEtEQSBPUk4gRk9SIEFTVCBTVE9SSg==",
                        },
                        {
                            "data_source_id": "11",
                            "calldata": "V0JUQyBESUEgQlRNIElPVFggRkVUIEpTVCBNQ08gS01EIEJUUyBRS0MgWUFNVjIgWFpDIEFLUk8gS0FJIE9HTiBXUlggS0RBIEZPUiBBU1QgU1RPUko=",
                        },
                        {
                            "external_id": "1",
                            "data_source_id": "12",
                            "calldata": "UkVOQlRDIFdCVEMgRElBIEJUTSBJT1RYIEZFVCBKU1QgTUNPIEtNRCBCVFMgUUtDIFlBTVYyIFhaQyBVT1MgQUtSTyBITlQgSE9UIEtBSSBPR04gV1JYIEtEQSBPUk4gRk9SIEFTVCBTVE9SSg==",
                        },
                    ],
                },
                "reports": [
                    {
                        "validator": "bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw",
                        "in_before_resolve": True,
                        "raw_reports": [
                            {
                                "external_id": "2",
                                "data": "MTMwMTQuNjgzNCwxMjc2OC45NTUzLDEuMjMyODUyLDAuMDYzNDI4NjMsMC4wMDYxNjgwMSwwLjA1MDM1OTQxLDAuMDI1ODk3MDksMi45MDUzNzQsMC40OTUwODAyMSwwLjAxOTIyNjE2LDAuMDA0ODg4NTcsOC4wNDk2NzIsNC4wNTM2MTcsMC4wOTg0NDAyMSwwLjAxMjgxMDk0LDEuMTA4MDczLDAuMDAwNTIzNjUsMC4wMTY4NTAzOSwwLjE4MDg4ODAyLDAuMDk0MjM3MzUsMC4yMDI3NDQxLDEuNzA4ODQ1LDAuMDMxNzgwNzEsMC4xMjUyNTE5NywwLjM4MQ==",
                            },
                            {
                                "external_id": "1",
                                "data": "MTI3NDkuNDcsMTI3OTIuNjUsMS4yMywwLjA2MzIzMSwwLjAwNjIxMDMyLDAuMDUwMjQsMC4wMjU4MzYwOCwyLjksMC40OTczNDMsMC4wMTkzMjEwMSwwLjAwNTAxNjk0LDcuOTgsNC4wNywwLjA5Nzc1NywwLjAxMjc5ODE4LDEuMTIsMC4wMDA1MjM5OSwwLjAxNjgwMTQzLDAuMTgxNTksMC4wOTM5LDAuMjAyNDAyLDEuNzUsMC4wMzE2NjEzOCwwLjEyMTA5MSwwLjM4MjAxNQo=",
                            },
                            {
                                "data": "MTI3OTMuMzgsMS4yNDMsMC4wNjMzNCwwLjAwNjI3LDAuMDQ5NzgsMC4wMjU4NSwyLjkwNywwLjQ5NjQsMC4wMTkzMiwwLjAwNTExOSw4LjEzMyw0LjA2MywwLjAxMjU0LDAuMDE3MSwwLjE4MywwLjA5MzgsMC4yMDczLDAuMDMyNzYsMC4xMjU4LDAuMzg0Mwo="
                            },
                        ],
                    },
                    {
                        "validator": "bandvaloper1yyv5jkqaukq0ajqn7vhkyhpff7h6e99j3gv0tr",
                        "in_before_resolve": True,
                        "raw_reports": [
                            {
                                "external_id": "1",
                                "data": "MTI3NDkuNDcsMTI3OTIuNjUsMS4yMywwLjA2MzIzMSwwLjAwNjIxMDMyLDAuMDUwMjQsMC4wMjU4MzYwOCwyLjksMC40OTczNDMsMC4wMTkzMjEwMSwwLjAwNTAxNjk0LDcuOTgsNC4wNywwLjA5Nzc1NywwLjAxMjc5ODE4LDEuMTIsMC4wMDA1MjM5OSwwLjAxNjgwMTQzLDAuMTgxNTksMC4wOTM5LDAuMjAyNDAyLDEuNzUsMC4wMzE2NjEzOCwwLjEyMTA5MSwwLjM4MjAxNQo=",
                            },
                            {
                                "external_id": "2",
                                "data": "MTMwMTQuNjgzNCwxMjc2OC45NTUzLDEuMjMyODUyLDAuMDYzNDI4NjMsMC4wMDYxNjgwMSwwLjA1MDM1OTQxLDAuMDI1ODk3MDksMi45MDUzNzQsMC40OTUwODAyMSwwLjAxOTIyNjE2LDAuMDA0ODg4NTcsOC4wNDk2NzIsNC4wNTM2MTcsMC4wOTg0NDAyMSwwLjAxMjgxMDk0LDEuMTA4MDczLDAuMDAwNTIzNjUsMC4wMTY4NTAzOSwwLjE4MDg4ODAyLDAuMDk0MjM3MzUsMC4yMDI3NDQxLDEuNzA4ODQ1LDAuMDMxNzgwNzEsMC4xMjUyNTE5NywwLjM4MQ==",
                            },
                            {
                                "data": "MTI3OTMuMzgsMS4yNDMsMC4wNjMzNCwwLjAwNjI3LDAuMDQ5NzgsMC4wMjU4NSwyLjkxMSwwLjQ5NjQsMC4wMTkzMiwwLjAwNTExOSw4LjEzMyw0LjA2MywwLjAxMjU0LDAuMDE3MSwwLjE4MywwLjA5MzgsMC4yMDczLDAuMDMyNzYsMC4xMjU4LDAuMzg0Mwo="
                            },
                        ],
                    },
                    {
                        "validator": "bandvaloper1xnryftxluq49fk52c5j5zrxcc5rzye96s70msl",
                        "in_before_resolve": True,
                        "raw_reports": [
                            {
                                "external_id": "1",
                                "data": "MTI3NDkuNDcsMTI3OTIuNjUsMS4yMywwLjA2MzIzMSwwLjAwNjIxMDMyLDAuMDUwMjQsMC4wMjU4MzYwOCwyLjksMC40OTczNDMsMC4wMTkzMjEwMSwwLjAwNTAxNjk0LDcuOTgsNC4wNywwLjA5Nzc1NywwLjAxMjc5ODE4LDEuMTIsMC4wMDA1MjM5OSwwLjAxNjgwMTQzLDAuMTgxNTksMC4wOTM5LDAuMjAyNDAyLDEuNzUsMC4wMzE2NjEzOCwwLjEyMTA5MSwwLjM4MjAxNQo=",
                            },
                            {
                                "external_id": "2",
                                "data": "MTMwMTQuNjgzNCwxMjc2OC45NTUzLDEuMjMyODUyLDAuMDYzNDI4NjMsMC4wMDYxNjgwMSwwLjA1MDM1OTQxLDAuMDI1ODk3MDksMi45MDUzNzQsMC40OTUwODAyMSwwLjAxOTIyNjE2LDAuMDA0ODg4NTcsOC4wNDk2NzIsNC4wNTM2MTcsMC4wOTg0NDAyMSwwLjAxMjgxMDk0LDEuMTA4MDczLDAuMDAwNTIzNjUsMC4wMTY4NTAzOSwwLjE4MDg4ODAyLDAuMDk0MjM3MzUsMC4yMDI3NDQxLDEuNzA4ODQ1LDAuMDMxNzgwNzEsMC4xMjUyNTE5NywwLjM4MQ==",
                            },
                            {
                                "data": "MTI3OTMuMzgsMS4yNDMsMC4wNjMzNCwwLjAwNjI3LDAuMDQ5NzgsMC4wMjU4NSwyLjkxMSwwLjQ5NjQsMC4wMTkzMiwwLjAwNTExOSw4LjEzMyw0LjA2MywwLjAxMjU0LDAuMDE3MSwwLjE4MywwLjA5MzgsMC4yMDczLDAuMDMyNzYsMC4xMjU4LDAuMzg0Mwo="
                            },
                        ],
                    },
                    {
                        "validator": "bandvaloper1v38hewjc0865dm4t89v5efh9rmum5rmrm7evg4",
                        "in_before_resolve": True,
                        "raw_reports": [
                            {
                                "external_id": "2",
                                "data": "MTMwMTQuNjgzNCwxMjc2OC45NTUzLDEuMjMyODUyLDAuMDYzNDI4NjMsMC4wMDYxNjgwMSwwLjA1MDM1OTQxLDAuMDI1ODk3MDksMi45MDUzNzQsMC40OTUwODAyMSwwLjAxOTIyNjE2LDAuMDA0ODg4NTcsOC4wNDk2NzIsNC4wNTM2MTcsMC4wOTg0NDAyMSwwLjAxMjgxMDk0LDEuMTA4MDczLDAuMDAwNTIzNjUsMC4wMTY4NTAzOSwwLjE4MDg4ODAyLDAuMDk0MjM3MzUsMC4yMDI3NDQxLDEuNzA4ODQ1LDAuMDMxNzgwNzEsMC4xMjUyNTE5NywwLjM4MQ==",
                            },
                            {
                                "external_id": "1",
                                "data": "MTI3NDkuNDcsMTI3OTIuNjUsMS4yMywwLjA2MzIzMSwwLjAwNjIxMDMyLDAuMDUwMjQsMC4wMjU4MzYwOCwyLjksMC40OTczNDMsMC4wMTkzMjEwMSwwLjAwNTAxNjk0LDcuOTgsNC4wNywwLjA5Nzc1NywwLjAxMjc5ODE4LDEuMTIsMC4wMDA1MjM5OSwwLjAxNjgwMTQzLDAuMTgxNTksMC4wOTM5LDAuMjAyNDAyLDEuNzUsMC4wMzE2NjEzOCwwLjEyMTA5MSwwLjM4MjAxNQo=",
                            },
                            {
                                "data": "MTI3OTMuMzgsMS4yNDMsMC4wNjMzNCwwLjAwNjI3LDAuMDQ5NzgsMC4wMjU4NSwyLjkxMSwwLjQ5NjQsMC4wMTkzMiwwLjAwNTExOSw4LjEzMyw0LjA2MywwLjAxMjU0LDAuMDE3MSwwLjE4MywwLjA5MzgsMC4yMDczLDAuMDMyNzYsMC4xMjU4LDAuMzg0Mwo="
                            },
                        ],
                    },
                ],
                "result": {
                    "request_packet_data": {
                        "client_id": "bandteam",
                        "oracle_script_id": "8",
                        "calldata": "AAAAGQAAAAZSRU5CVEMAAAAEV0JUQwAAAANESUEAAAADQlRNAAAABElPVFgAAAADRkVUAAAAA0pTVAAAAANNQ08AAAADS01EAAAAA0JUUwAAAANRS0MAAAAFWUFNVjIAAAADWFpDAAAAA1VPUwAAAARBS1JPAAAAA0hOVAAAAANIT1QAAAADS0FJAAAAA09HTgAAAANXUlgAAAADS0RBAAAAA09STgAAAANGT1IAAAADQVNUAAAABVNUT1JKAAAAADuaygA=",
                        "ask_count": "4",
                        "min_count": "3",
                    },
                    "response_packet_data": {
                        "client_id": "bandteam",
                        "request_id": "754674",
                        "ans_count": "4",
                        "request_time": "1603354320",
                        "resolve_time": "1603354323",
                        "resolve_status": 1,
                        "result": "AAAAGQAAC7dXmw1gAAALooVb5oAAAAAASXvUIAAAAAADxn3fAAAAAABewxAAAAAAAv6aAAAAAAABinCQAAAAAK0sfTAAAAAAHZZ2gAAAAAABJszAAAAAAABMjWwAAAAB38w/PwAAAADyLHW/AAAAAAXY3a0AAAAAAMNI5AAAAABCZtkUAAAAAAAH/isAAAAAAQEd1gAAAAAK0tfwAAAAAAWYzOAAAAAADBWhJAAAAABnFOSkAAAAAAHk72UAAAAAB3cxgQAAAAAWxRYY",
                    },
                },
            },
        },
        status_code=200,
    )

    assert client.get_latest_request(
        oid=8,
        calldata=bytes.fromhex(
            "000000190000000652454e4254430000000457425443000000034449410000000342544d00000004494f545800000003464554000000034a5354000000034d434f000000034b4d440000000342545300000003514b430000000559414d563200000003585a4300000003554f5300000004414b524f00000003484e5400000003484f54000000034b4149000000034f474e00000003575258000000034b4441000000034f524e00000003464f52000000034153540000000553544f524a000000003b9aca00"
        ),
        min_count=3,
        ask_count=4,
    ) == RequestInfo(
        request=Request(
            oracle_script_id=8,
            requested_validators=[
                "bandvaloper1xnryftxluq49fk52c5j5zrxcc5rzye96s70msl",
                "bandvaloper1v38hewjc0865dm4t89v5efh9rmum5rmrm7evg4",
                "bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw",
                "bandvaloper1yyv5jkqaukq0ajqn7vhkyhpff7h6e99j3gv0tr",
            ],
            min_count=3,
            request_height=2243579,
            raw_requests=[
                RawRequest(
                    data_source_id=15,
                    external_id=2,
                    calldata=b"RENBTC WBTC DIA BTM IOTX FET JST MCO KMD BTS QKC YAMV2 XZC UOS AKRO HNT HOT KAI OGN WRX KDA ORN FOR AST STORJ",
                ),
                RawRequest(
                    data_source_id=11,
                    external_id=0,
                    calldata=b"WBTC DIA BTM IOTX FET JST MCO KMD BTS QKC YAMV2 XZC AKRO KAI OGN WRX KDA FOR AST STORJ",
                ),
                RawRequest(
                    data_source_id=12,
                    external_id=1,
                    calldata=b"RENBTC WBTC DIA BTM IOTX FET JST MCO KMD BTS QKC YAMV2 XZC UOS AKRO HNT HOT KAI OGN WRX KDA ORN FOR AST STORJ",
                ),
            ],
            client_id="bandteam",
            calldata=b"\x00\x00\x00\x19\x00\x00\x00\x06RENBTC\x00\x00\x00\x04WBTC\x00\x00\x00\x03DIA\x00\x00\x00\x03BTM\x00\x00\x00\x04IOTX\x00\x00\x00\x03FET\x00\x00\x00\x03JST\x00\x00\x00\x03MCO\x00\x00\x00\x03KMD\x00\x00\x00\x03BTS\x00\x00\x00\x03QKC\x00\x00\x00\x05YAMV2\x00\x00\x00\x03XZC\x00\x00\x00\x03UOS\x00\x00\x00\x04AKRO\x00\x00\x00\x03HNT\x00\x00\x00\x03HOT\x00\x00\x00\x03KAI\x00\x00\x00\x03OGN\x00\x00\x00\x03WRX\x00\x00\x00\x03KDA\x00\x00\x00\x03ORN\x00\x00\x00\x03FOR\x00\x00\x00\x03AST\x00\x00\x00\x05STORJ\x00\x00\x00\x00;\x9a\xca\x00",
        ),
        reports=[
            Report(
                validator="bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw",
                raw_reports=[
                    RawReport(
                        external_id=2,
                        data=b"13014.6834,12768.9553,1.232852,0.06342863,0.00616801,0.05035941,0.02589709,2.905374,0.49508021,0.01922616,0.00488857,8.049672,4.053617,0.09844021,0.01281094,1.108073,0.00052365,0.01685039,0.18088802,0.09423735,0.2027441,1.708845,0.03178071,0.12525197,0.381",
                    ),
                    RawReport(
                        external_id=1,
                        data=b"12749.47,12792.65,1.23,0.063231,0.00621032,0.05024,0.02583608,2.9,0.497343,0.01932101,0.00501694,7.98,4.07,0.097757,0.01279818,1.12,0.00052399,0.01680143,0.18159,0.0939,0.202402,1.75,0.03166138,0.121091,0.382015\n",
                    ),
                    RawReport(
                        external_id=0,
                        data=b"12793.38,1.243,0.06334,0.00627,0.04978,0.02585,2.907,0.4964,0.01932,0.005119,8.133,4.063,0.01254,0.0171,0.183,0.0938,0.2073,0.03276,0.1258,0.3843\n",
                    ),
                ],
                in_before_resolve=True,
            ),
            Report(
                validator="bandvaloper1yyv5jkqaukq0ajqn7vhkyhpff7h6e99j3gv0tr",
                raw_reports=[
                    RawReport(
                        external_id=1,
                        data=b"12749.47,12792.65,1.23,0.063231,0.00621032,0.05024,0.02583608,2.9,0.497343,0.01932101,0.00501694,7.98,4.07,0.097757,0.01279818,1.12,0.00052399,0.01680143,0.18159,0.0939,0.202402,1.75,0.03166138,0.121091,0.382015\n",
                    ),
                    RawReport(
                        external_id=2,
                        data=b"13014.6834,12768.9553,1.232852,0.06342863,0.00616801,0.05035941,0.02589709,2.905374,0.49508021,0.01922616,0.00488857,8.049672,4.053617,0.09844021,0.01281094,1.108073,0.00052365,0.01685039,0.18088802,0.09423735,0.2027441,1.708845,0.03178071,0.12525197,0.381",
                    ),
                    RawReport(
                        external_id=0,
                        data=b"12793.38,1.243,0.06334,0.00627,0.04978,0.02585,2.911,0.4964,0.01932,0.005119,8.133,4.063,0.01254,0.0171,0.183,0.0938,0.2073,0.03276,0.1258,0.3843\n",
                    ),
                ],
                in_before_resolve=True,
            ),
            Report(
                validator="bandvaloper1xnryftxluq49fk52c5j5zrxcc5rzye96s70msl",
                raw_reports=[
                    RawReport(
                        external_id=1,
                        data=b"12749.47,12792.65,1.23,0.063231,0.00621032,0.05024,0.02583608,2.9,0.497343,0.01932101,0.00501694,7.98,4.07,0.097757,0.01279818,1.12,0.00052399,0.01680143,0.18159,0.0939,0.202402,1.75,0.03166138,0.121091,0.382015\n",
                    ),
                    RawReport(
                        external_id=2,
                        data=b"13014.6834,12768.9553,1.232852,0.06342863,0.00616801,0.05035941,0.02589709,2.905374,0.49508021,0.01922616,0.00488857,8.049672,4.053617,0.09844021,0.01281094,1.108073,0.00052365,0.01685039,0.18088802,0.09423735,0.2027441,1.708845,0.03178071,0.12525197,0.381",
                    ),
                    RawReport(
                        external_id=0,
                        data=b"12793.38,1.243,0.06334,0.00627,0.04978,0.02585,2.911,0.4964,0.01932,0.005119,8.133,4.063,0.01254,0.0171,0.183,0.0938,0.2073,0.03276,0.1258,0.3843\n",
                    ),
                ],
                in_before_resolve=True,
            ),
            Report(
                validator="bandvaloper1v38hewjc0865dm4t89v5efh9rmum5rmrm7evg4",
                raw_reports=[
                    RawReport(
                        external_id=2,
                        data=b"13014.6834,12768.9553,1.232852,0.06342863,0.00616801,0.05035941,0.02589709,2.905374,0.49508021,0.01922616,0.00488857,8.049672,4.053617,0.09844021,0.01281094,1.108073,0.00052365,0.01685039,0.18088802,0.09423735,0.2027441,1.708845,0.03178071,0.12525197,0.381",
                    ),
                    RawReport(
                        external_id=1,
                        data=b"12749.47,12792.65,1.23,0.063231,0.00621032,0.05024,0.02583608,2.9,0.497343,0.01932101,0.00501694,7.98,4.07,0.097757,0.01279818,1.12,0.00052399,0.01680143,0.18159,0.0939,0.202402,1.75,0.03166138,0.121091,0.382015\n",
                    ),
                    RawReport(
                        external_id=0,
                        data=b"12793.38,1.243,0.06334,0.00627,0.04978,0.02585,2.911,0.4964,0.01932,0.005119,8.133,4.063,0.01254,0.0171,0.183,0.0938,0.2073,0.03276,0.1258,0.3843\n",
                    ),
                ],
                in_before_resolve=True,
            ),
        ],
        result=Result(
            request_packet_data=RequestPacketData(
                oracle_script_id=8,
                ask_count=4,
                min_count=3,
                client_id="bandteam",
                calldata=b"\x00\x00\x00\x19\x00\x00\x00\x06RENBTC\x00\x00\x00\x04WBTC\x00\x00\x00\x03DIA\x00\x00\x00\x03BTM\x00\x00\x00\x04IOTX\x00\x00\x00\x03FET\x00\x00\x00\x03JST\x00\x00\x00\x03MCO\x00\x00\x00\x03KMD\x00\x00\x00\x03BTS\x00\x00\x00\x03QKC\x00\x00\x00\x05YAMV2\x00\x00\x00\x03XZC\x00\x00\x00\x03UOS\x00\x00\x00\x04AKRO\x00\x00\x00\x03HNT\x00\x00\x00\x03HOT\x00\x00\x00\x03KAI\x00\x00\x00\x03OGN\x00\x00\x00\x03WRX\x00\x00\x00\x03KDA\x00\x00\x00\x03ORN\x00\x00\x00\x03FOR\x00\x00\x00\x03AST\x00\x00\x00\x05STORJ\x00\x00\x00\x00;\x9a\xca\x00",
            ),
            response_packet_data=ResponsePacketData(
                request_id=754674,
                request_time=1603354320,
                resolve_time=1603354323,
                resolve_status=1,
                ans_count=4,
                client_id="bandteam",
                result=b"\x00\x00\x00\x19\x00\x00\x0b\xb7W\x9b\r`\x00\x00\x0b\xa2\x85[\xe6\x80\x00\x00\x00\x00I{\xd4 \x00\x00\x00\x00\x03\xc6}\xdf\x00\x00\x00\x00\x00^\xc3\x10\x00\x00\x00\x00\x02\xfe\x9a\x00\x00\x00\x00\x00\x01\x8ap\x90\x00\x00\x00\x00\xad,}0\x00\x00\x00\x00\x1d\x96v\x80\x00\x00\x00\x00\x01&\xcc\xc0\x00\x00\x00\x00\x00L\x8dl\x00\x00\x00\x01\xdf\xcc??\x00\x00\x00\x00\xf2,u\xbf\x00\x00\x00\x00\x05\xd8\xdd\xad\x00\x00\x00\x00\x00\xc3H\xe4\x00\x00\x00\x00Bf\xd9\x14\x00\x00\x00\x00\x00\x07\xfe+\x00\x00\x00\x00\x01\x01\x1d\xd6\x00\x00\x00\x00\n\xd2\xd7\xf0\x00\x00\x00\x00\x05\x98\xcc\xe0\x00\x00\x00\x00\x0c\x15\xa1$\x00\x00\x00\x00g\x14\xe4\xa4\x00\x00\x00\x00\x01\xe4\xefe\x00\x00\x00\x00\x07w1\x81\x00\x00\x00\x00\x16\xc5\x16\x18",
            ),
        ),
    )


def test_get_reporters(requests_mock):

    requests_mock.register_uri(
        "GET",
        "{}/oracle/reporters/bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre".format(TEST_RPC),
        json={
            "height": "2245131",
            "result": [
                "band1yyv5jkqaukq0ajqn7vhkyhpff7h6e99ja7gvwg",
                "band19nf0sqnjycnvpexlxs6hjz9qrhhlhsu9pdty0r",
                "band1fndxcmg0h5vhr8cph7gryryqfn9yqp90lysjtm",
                "band10ec0p96j60duce5qagju5axuja0rj8luqrzl0k",
                "band15pm9scujgkpwpy2xa2j53tvs9ylunjn0g73a9s",
                "band1cehe3sxk7f4rmvwdf6lxh3zexen7fn02zyltwy",
            ],
        },
        status_code=200,
    )

    assert client.get_reporters("bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre") == [
        "band1yyv5jkqaukq0ajqn7vhkyhpff7h6e99ja7gvwg",
        "band19nf0sqnjycnvpexlxs6hjz9qrhhlhsu9pdty0r",
        "band1fndxcmg0h5vhr8cph7gryryqfn9yqp90lysjtm",
        "band10ec0p96j60duce5qagju5axuja0rj8luqrzl0k",
        "band15pm9scujgkpwpy2xa2j53tvs9ylunjn0g73a9s",
        "band1cehe3sxk7f4rmvwdf6lxh3zexen7fn02zyltwy",
    ]


def test_get_price_symbols(requests_mock):
    requests_mock.register_uri(
        "GET",
        "{}/oracle/price_symbols?min_count=3&ask_count=4".format(TEST_RPC),
        json={
            "height": "2951872",
            "result": [
                "2KEY",
                "ABYSS",
                "ADA",
                "AKRO",
                "ALGO",
                "AMPL",
                "ANT",
                "AST",
                "ATOM",
                "AUD",
                "BAL",
                "BAND",
                "BAT",
                "BCH",
                "BLZ",
                "BNB",
                "BNT",
                "BRL",
                "BSV",
                "BTC",
                "BTG",
                "BTM",
                "BTS",
                "BTT",
                "BTU",
                "BUSD",
                "BZRX",
                "CAD",
                "CHF",
                "CKB",
                "CND",
                "CNY",
                "COMP",
                "CREAM",
                "CRO",
                "CRV",
                "CVC",
                "DAI",
                "DASH",
                "DCR",
                "DGB",
                "DGX",
                "DIA",
                "DOGE",
                "DOT",
                "EGLD",
                "ELF",
                "ENJ",
                "EOS",
                "EQUAD",
                "ETC",
                "ETH",
                "EUR",
                "EURS",
                "EWT",
                "FET",
                "FNX",
                "FOR",
                "FTM",
                "FTT",
                "FXC",
                "GBP",
                "GDC",
                "GEN",
                "GHT",
                "GNO",
                "GVT",
                "HBAR",
                "HKD",
                "HNT",
                "HOT",
                "HT",
                "ICX",
                "INR",
                "IOST",
                "IOTX",
                "JPY",
                "JST",
                "KAI",
                "KAVA",
                "KDA",
                "KEY",
                "KMD",
                "KNC",
                "KRW",
                "KSM",
                "LEND",
                "LEO",
                "LINA",
                "LINK",
                "LOOM",
                "LRC",
                "LSK",
                "LTC",
                "LUNA",
                "MANA",
                "MATIC",
                "MCO",
                "MET",
                "MFG",
                "MIOTA",
                "MKR",
                "MLN",
                "MNT",
                "MTL",
                "MYB",
                "NEO",
                "NEXXO",
                "NMR",
                "NOK",
                "NPXS",
                "NXM",
                "NZD",
                "OCEAN",
                "OGN",
                "OKB",
                "OMG",
                "ONE",
                "ONT",
                "ORN",
                "OST",
                "OXT",
                "PAX",
                "PAXG",
                "PAY",
                "PBTC",
                "PLR",
                "PLTC",
                "PNK",
                "PNT",
                "POLY",
                "POWR",
                "QKC",
                "QNT",
                "RAE",
                "REN",
                "RENBTC",
                "REP",
                "REQ",
                "RLC",
                "RMB",
                "RSR",
                "RSV",
                "RUB",
                "RUNE",
                "RVN",
                "SAN",
                "SC",
                "SGD",
                "SNT",
                "SNX",
                "SOL",
                "SPIKE",
                "SPN",
                "SRM",
                "STMX",
                "STORJ",
                "STX",
                "SUSD",
                "SUSHI",
                "SXP",
                "THETA",
                "TKN",
                "TKX",
                "TOMO",
                "TRB",
                "TRX",
                "TRYB",
                "TUSD",
                "UBT",
                "UNI",
                "UOS",
                "UPP",
                "USDC",
                "USDS",
                "USDT",
                "VET",
                "VIDT",
                "WAN",
                "WAVES",
                "WBTC",
                "WNXM",
                "WRX",
                "XAG",
                "XAU",
                "XDR",
                "XEM",
                "XHV",
                "XLM",
                "XMR",
                "XRP",
                "XTZ",
                "XZC",
                "YAMV2",
                "YFI",
                "YFII",
                "YFV",
                "ZEC",
                "ZRX",
            ],
        },
        status_code=200,
    )

    assert client.get_price_symbols(3, 4) == [
        "2KEY",
        "ABYSS",
        "ADA",
        "AKRO",
        "ALGO",
        "AMPL",
        "ANT",
        "AST",
        "ATOM",
        "AUD",
        "BAL",
        "BAND",
        "BAT",
        "BCH",
        "BLZ",
        "BNB",
        "BNT",
        "BRL",
        "BSV",
        "BTC",
        "BTG",
        "BTM",
        "BTS",
        "BTT",
        "BTU",
        "BUSD",
        "BZRX",
        "CAD",
        "CHF",
        "CKB",
        "CND",
        "CNY",
        "COMP",
        "CREAM",
        "CRO",
        "CRV",
        "CVC",
        "DAI",
        "DASH",
        "DCR",
        "DGB",
        "DGX",
        "DIA",
        "DOGE",
        "DOT",
        "EGLD",
        "ELF",
        "ENJ",
        "EOS",
        "EQUAD",
        "ETC",
        "ETH",
        "EUR",
        "EURS",
        "EWT",
        "FET",
        "FNX",
        "FOR",
        "FTM",
        "FTT",
        "FXC",
        "GBP",
        "GDC",
        "GEN",
        "GHT",
        "GNO",
        "GVT",
        "HBAR",
        "HKD",
        "HNT",
        "HOT",
        "HT",
        "ICX",
        "INR",
        "IOST",
        "IOTX",
        "JPY",
        "JST",
        "KAI",
        "KAVA",
        "KDA",
        "KEY",
        "KMD",
        "KNC",
        "KRW",
        "KSM",
        "LEND",
        "LEO",
        "LINA",
        "LINK",
        "LOOM",
        "LRC",
        "LSK",
        "LTC",
        "LUNA",
        "MANA",
        "MATIC",
        "MCO",
        "MET",
        "MFG",
        "MIOTA",
        "MKR",
        "MLN",
        "MNT",
        "MTL",
        "MYB",
        "NEO",
        "NEXXO",
        "NMR",
        "NOK",
        "NPXS",
        "NXM",
        "NZD",
        "OCEAN",
        "OGN",
        "OKB",
        "OMG",
        "ONE",
        "ONT",
        "ORN",
        "OST",
        "OXT",
        "PAX",
        "PAXG",
        "PAY",
        "PBTC",
        "PLR",
        "PLTC",
        "PNK",
        "PNT",
        "POLY",
        "POWR",
        "QKC",
        "QNT",
        "RAE",
        "REN",
        "RENBTC",
        "REP",
        "REQ",
        "RLC",
        "RMB",
        "RSR",
        "RSV",
        "RUB",
        "RUNE",
        "RVN",
        "SAN",
        "SC",
        "SGD",
        "SNT",
        "SNX",
        "SOL",
        "SPIKE",
        "SPN",
        "SRM",
        "STMX",
        "STORJ",
        "STX",
        "SUSD",
        "SUSHI",
        "SXP",
        "THETA",
        "TKN",
        "TKX",
        "TOMO",
        "TRB",
        "TRX",
        "TRYB",
        "TUSD",
        "UBT",
        "UNI",
        "UOS",
        "UPP",
        "USDC",
        "USDS",
        "USDT",
        "VET",
        "VIDT",
        "WAN",
        "WAVES",
        "WBTC",
        "WNXM",
        "WRX",
        "XAG",
        "XAU",
        "XDR",
        "XEM",
        "XHV",
        "XLM",
        "XMR",
        "XRP",
        "XTZ",
        "XZC",
        "YAMV2",
        "YFI",
        "YFII",
        "YFV",
        "ZEC",
        "ZRX",
    ]
