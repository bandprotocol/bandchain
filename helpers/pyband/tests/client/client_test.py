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
    EpochTime,
    EVMProof,
)
from pyband.utils import parse_epoch_time
from pyband.data import Coin
from pyband.error import EmptyRequestMsgError

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
                time=EpochTime(parse_epoch_time("2020-11-05T09:15:18.445494105Z")),
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
        address=Address.from_acc_bech32("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"),
        coins=[Coin(amount=104082359107, denom="uband")],
        public_key={
            "type": "tendermint/PubKeySecp256k1",
            "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
        },
        account_number=36,
        sequence=927,
    )


def test_get_account_not_found(requests_mock):
    requests_mock.register_uri(
        "GET",
        "{}/auth/accounts/band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte".format(TEST_RPC),
        json={
            "height": "650788",
            "result": {
                "type": "cosmos-sdk/Account",
                "value": {
                    "address": "",
                    "coins": [],
                    "public_key": None,
                    "account_number": "0",
                    "sequence": "0",
                },
            },
        },
        status_code=200,
    )

    assert client.get_account(Address.from_acc_bech32("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte")) == None


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
        owner=Address.from_acc_bech32("band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs"),
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
        owner=Address.from_acc_bech32("band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs"),
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

    assert client.get_reporters(Address.from_val_bech32("bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre")) == [
        Address.from_acc_bech32("band1yyv5jkqaukq0ajqn7vhkyhpff7h6e99ja7gvwg"),
        Address.from_acc_bech32("band19nf0sqnjycnvpexlxs6hjz9qrhhlhsu9pdty0r"),
        Address.from_acc_bech32("band1fndxcmg0h5vhr8cph7gryryqfn9yqp90lysjtm"),
        Address.from_acc_bech32("band10ec0p96j60duce5qagju5axuja0rj8luqrzl0k"),
        Address.from_acc_bech32("band15pm9scujgkpwpy2xa2j53tvs9ylunjn0g73a9s"),
        Address.from_acc_bech32("band1cehe3sxk7f4rmvwdf6lxh3zexen7fn02zyltwy"),
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
            ],
        },
        status_code=200,
    )

    assert client.get_price_symbols(3, 4) == [
        "2KEY",
        "ABYSS",
        "ADA",
        "AKRO",
    ]


def test_get_request_id_by_tx_hash_send_and_request_msgs(requests_mock):
    requests_mock.register_uri(
        "GET",
        "{}/txs/90ED70061C1A24B1141F81BADEDAB19570D0B9B255412BF5D680A9BC93539115".format(TEST_RPC),
        json={
            "height": "3739",
            "txhash": "90ED70061C1A24B1141F81BADEDAB19570D0B9B255412BF5D680A9BC93539115",
            "raw_log": '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"band178yydxyzplrh5v0jegvtrahymqragf6spy8kdd"},{"key":"sender","value":"band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"},{"key":"amount","value":"100000uband"}]}]},{"msg_index":1,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"4"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"2"},{"key":"min_count","value":"2"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]}]',
            "logs": [
                {
                    "msg_index": 0,
                    "log": "",
                    "events": [
                        {
                            "type": "message",
                            "attributes": [
                                {"key": "action", "value": "send"},
                                {"key": "sender", "value": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"},
                                {"key": "module", "value": "bank"},
                            ],
                        },
                        {
                            "type": "transfer",
                            "attributes": [
                                {"key": "recipient", "value": "band178yydxyzplrh5v0jegvtrahymqragf6spy8kdd"},
                                {"key": "sender", "value": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"},
                                {"key": "amount", "value": "100000uband"},
                            ],
                        },
                    ],
                },
                {
                    "msg_index": 1,
                    "log": "",
                    "events": [
                        {"type": "message", "attributes": [{"key": "action", "value": "request"}]},
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
                                {"key": "id", "value": "4"},
                                {"key": "client_id", "value": "from_pyband"},
                                {"key": "oracle_script_id", "value": "1"},
                                {"key": "calldata", "value": "000000034254430000000000000001"},
                                {"key": "ask_count", "value": "2"},
                                {"key": "min_count", "value": "2"},
                                {"key": "gas_used", "value": "2405"},
                                {"key": "validator", "value": "bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},
                                {"key": "validator", "value": "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"},
                            ],
                        },
                    ],
                },
            ],
            "gas_wanted": "1000000",
            "gas_used": "298476",
            "tx": {
                "type": "cosmos-sdk/StdTx",
                "value": {
                    "msg": [
                        {
                            "type": "cosmos-sdk/MsgSend",
                            "value": {
                                "from_address": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
                                "to_address": "band178yydxyzplrh5v0jegvtrahymqragf6spy8kdd",
                                "amount": [{"denom": "uband", "amount": "100000"}],
                            },
                        },
                        {
                            "type": "oracle/Request",
                            "value": {
                                "oracle_script_id": "1",
                                "calldata": "AAAAA0JUQwAAAAAAAAAB",
                                "ask_count": "2",
                                "min_count": "2",
                                "client_id": "from_pyband",
                                "sender": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
                            },
                        },
                    ],
                    "fee": {"amount": [{"denom": "uband", "amount": "0"}], "gas": "1000000"},
                    "signatures": [
                        {
                            "pub_key": {
                                "type": "tendermint/PubKeySecp256k1",
                                "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
                            },
                            "signature": "K6PmN0HluRRm7zliKJho9F2OYpokB5JrYlqAC+KmHQMwKBpRYNaZzYTGzoeBol2mm3sfcdUo8rzwrQngzP8s+g==",
                        }
                    ],
                    "memo": "TEST",
                },
            },
            "timestamp": "2020-11-09T09:29:49Z",
        },
        status_code=200,
    )
    assert client.get_request_id_by_tx_hash(
        HexBytes(bytes.fromhex("90ED70061C1A24B1141F81BADEDAB19570D0B9B255412BF5D680A9BC93539115"))
    ) == [4]


def test_get_request_id_by_tx_hash_multi_msgs(requests_mock):
    requests_mock.register_uri(
        "GET",
        "{}/txs/0838E29162B87D9D41E2BAE49C97272970453F2CFA708FBA6B8BE71F314BEB5B".format(TEST_RPC),
        json={
            "height": "279",
            "txhash": "0838E29162B87D9D41E2BAE49C97272970453F2CFA708FBA6B8BE71F314BEB5B",
            "raw_log": '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"1"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"2"},{"key":"min_count","value":"2"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]},{"msg_index":1,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"2"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"2"},{"key":"min_count","value":"2"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]},{"msg_index":2,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"3"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"2"},{"key":"min_count","value":"2"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},{"key":"validator","value":"bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"}]}]}]',
            "logs": [
                {
                    "msg_index": 0,
                    "log": "",
                    "events": [
                        {"type": "message", "attributes": [{"key": "action", "value": "request"}]},
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
                                {"key": "id", "value": "1"},
                                {"key": "client_id", "value": "from_pyband"},
                                {"key": "oracle_script_id", "value": "1"},
                                {"key": "calldata", "value": "000000034254430000000000000001"},
                                {"key": "ask_count", "value": "2"},
                                {"key": "min_count", "value": "2"},
                                {"key": "gas_used", "value": "2405"},
                                {"key": "validator", "value": "bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr"},
                                {"key": "validator", "value": "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"},
                            ],
                        },
                    ],
                },
                {
                    "msg_index": 1,
                    "log": "",
                    "events": [
                        {"type": "message", "attributes": [{"key": "action", "value": "request"}]},
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
                                {"key": "id", "value": "2"},
                                {"key": "client_id", "value": "from_pyband"},
                                {"key": "oracle_script_id", "value": "1"},
                                {"key": "calldata", "value": "000000034254430000000000000001"},
                                {"key": "ask_count", "value": "2"},
                                {"key": "min_count", "value": "2"},
                                {"key": "gas_used", "value": "2405"},
                                {"key": "validator", "value": "bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},
                                {"key": "validator", "value": "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"},
                            ],
                        },
                    ],
                },
                {
                    "msg_index": 2,
                    "log": "",
                    "events": [
                        {"type": "message", "attributes": [{"key": "action", "value": "request"}]},
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
                                {"key": "id", "value": "3"},
                                {"key": "client_id", "value": "from_pyband"},
                                {"key": "oracle_script_id", "value": "1"},
                                {"key": "calldata", "value": "000000034254430000000000000001"},
                                {"key": "ask_count", "value": "2"},
                                {"key": "min_count", "value": "2"},
                                {"key": "gas_used", "value": "2405"},
                                {"key": "validator", "value": "bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},
                                {"key": "validator", "value": "bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"},
                            ],
                        },
                    ],
                },
            ],
            "gas_wanted": "1000000",
            "gas_used": "782736",
            "tx": {
                "type": "cosmos-sdk/StdTx",
                "value": {
                    "msg": [
                        {
                            "type": "oracle/Request",
                            "value": {
                                "oracle_script_id": "1",
                                "calldata": "AAAAA0JUQwAAAAAAAAAB",
                                "ask_count": "2",
                                "min_count": "2",
                                "client_id": "from_pyband",
                                "sender": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
                            },
                        },
                        {
                            "type": "oracle/Request",
                            "value": {
                                "oracle_script_id": "1",
                                "calldata": "AAAAA0JUQwAAAAAAAAAB",
                                "ask_count": "2",
                                "min_count": "2",
                                "client_id": "from_pyband",
                                "sender": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
                            },
                        },
                        {
                            "type": "oracle/Request",
                            "value": {
                                "oracle_script_id": "1",
                                "calldata": "AAAAA0JUQwAAAAAAAAAB",
                                "ask_count": "2",
                                "min_count": "2",
                                "client_id": "from_pyband",
                                "sender": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
                            },
                        },
                    ],
                    "fee": {"amount": [{"denom": "uband", "amount": "0"}], "gas": "1000000"},
                    "signatures": [
                        {
                            "pub_key": {
                                "type": "tendermint/PubKeySecp256k1",
                                "value": "A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8",
                            },
                            "signature": "PaxPc1330hLWZXBUHtPAPQdP4tS2LvGAiwaT8isxr8UKdDx0uXjLblTBjWO+yHENTjGxMdb2cI3BUOCIIoFCzQ==",
                        }
                    ],
                    "memo": "TEST",
                },
            },
            "timestamp": "2020-11-09T07:44:35Z",
        },
        status_code=200,
    )

    assert client.get_request_id_by_tx_hash(
        HexBytes(bytes.fromhex("0838E29162B87D9D41E2BAE49C97272970453F2CFA708FBA6B8BE71F314BEB5B"))
    ) == [1, 2, 3]


def test_get_request_id_by_tx_hash_no_request_msg(requests_mock):
    requests_mock.register_uri(
        "GET",
        "{}/txs/9F83E4994C048F784D0E30F45696C0A1E5BA7407B2E1833B439FA172B3B75F00".format(TEST_RPC),
        json={
            "height": "3740",
            "txhash": "9F83E4994C048F784D0E30F45696C0A1E5BA7407B2E1833B439FA172B3B75F00",
            "raw_log": '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"report"}]},{"type":"report","attributes":[{"key":"id","value":"4"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]}]',
            "logs": [
                {
                    "msg_index": 0,
                    "log": "",
                    "events": [
                        {"type": "message", "attributes": [{"key": "action", "value": "report"}]},
                        {
                            "type": "report",
                            "attributes": [
                                {"key": "id", "value": "4"},
                                {"key": "validator", "value": "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"},
                            ],
                        },
                    ],
                }
            ],
            "gas_wanted": "62945",
            "gas_used": "52449",
            "tx": {
                "type": "cosmos-sdk/StdTx",
                "value": {
                    "msg": [
                        {
                            "type": "oracle/Report",
                            "value": {
                                "request_id": "4",
                                "raw_reports": [
                                    {"external_id": "1", "data": "MTU0MzMuMQo="},
                                    {"external_id": "3", "data": "MTU0MTYuMDc1Cg=="},
                                    {"external_id": "2", "data": "MTU0MzcuMTYK"},
                                ],
                                "validator": "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec",
                                "reporter": "band1ue0623hwqkvm5s5hq0jnqwh4ende28gvmlzvd2",
                            },
                        }
                    ],
                    "fee": {"amount": [], "gas": "62945"},
                    "signatures": [
                        {
                            "pub_key": {
                                "type": "tendermint/PubKeySecp256k1",
                                "value": "A4OH54K/spetlw9jNR8LjqiHKi85jWKyV8zZDZH6dBij",
                            },
                            "signature": "W5RuZvIsGxoVqTHCpPzWuD6pvaA6YedEaQ8TfS3d4AYuykREcltOPOigun8wXm5dhfFzAf/BbgP9vccQ8mmkHw==",
                        }
                    ],
                    "memo": "yoda:/exec:lambda:1.1.2",
                },
            },
            "timestamp": "2020-11-09T09:29:51Z",
        },
        status_code=200,
    )
    with pytest.raises(EmptyRequestMsgError, match="There is no request message in this tx"):
        client.get_request_id_by_tx_hash(
            HexBytes(bytes.fromhex("9F83E4994C048F784D0E30F45696C0A1E5BA7407B2E1833B439FA172B3B75F00"))
        )


def test_get_reference_data(requests_mock):

    requests_mock.register_uri(
        "POST",
        "{}/oracle/request_prices".format(TEST_RPC),
        json={
            "height": "2953006",
            "result": [
                {
                    "symbol": "BTC",
                    "multiplier": "1000000000",
                    "px": "16242693800000",
                    "request_id": "1171969",
                    "resolve_time": "1605512243",
                },
                {
                    "symbol": "ETH",
                    "multiplier": "1000000000",
                    "px": "454523400000",
                    "request_id": "1171969",
                    "resolve_time": "1605512943",
                },
                {
                    "symbol": "TRX",
                    "multiplier": "1000000000",
                    "px": "25428330",
                    "request_id": "1171969",
                    "resolve_time": "1605512443",
                },
            ],
        },
        status_code=200,
    )

    [result1, result2] = client.get_reference_data(["BTC/USD", "TRX/ETH"], 3, 4)

    assert result1.pair == "BTC/USD"
    assert result1.rate == 16242.6938
    assert result1.updated_at.base == 1605512243
    assert result2.pair == "TRX/ETH"
    assert result2.rate == 0.000055945040453362794
    assert result2.updated_at.base == 1605512443
    assert result2.updated_at.quote == 1605512943


def test_get_request_evm_proof_by_request_id(requests_mock):
    requests_mock.register_uri(
        "GET",
        "{}/oracle/proof/1".format(TEST_RPC),
        json={
            "height": "0",
            "result": {
                "jsonProof": {
                    "blockHeight": "2622114",
                    "oracleDataProof": {
                        "requestPacket": {
                            "client_id": "test",
                            "oracle_script_id": "1",
                            "calldata": "AAAABGZhc3Q=",
                            "ask_count": "4",
                            "min_count": "3",
                        },
                        "responsePacket": {
                            "client_id": "test",
                            "request_id": "1",
                            "ans_count": "4",
                            "request_time": "1600357375",
                            "resolve_time": "1600357377",
                            "resolve_status": 1,
                            "result": "AAAAAAAAEYA=",
                        },
                        "version": "624",
                        "merklePaths": [
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 1,
                                "subtreeSize": "2",
                                "subtreeVersion": "818",
                                "siblingHash": "5744F905BEA848192798B1D1C624C65E4CA5EF6E964F37788EBE6EF49C67A6B9",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 2,
                                "subtreeSize": "4",
                                "subtreeVersion": "12005",
                                "siblingHash": "370FB5ECE9C7F3C742CC6F061E3F28C79063233442A0A7718AC7DAC1185DB625",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 3,
                                "subtreeSize": "8",
                                "subtreeVersion": "17879",
                                "siblingHash": "83C97412076743EEB622FCBE5877481401D2C9BC621950112B94BE4F7C7D7A0D",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 4,
                                "subtreeSize": "16",
                                "subtreeVersion": "18660",
                                "siblingHash": "C183A0D2A4614B7EBAD2DC8C2B2B7C94F42971FC47AD012FB6FEFE9AF7AC1C90",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 5,
                                "subtreeSize": "32",
                                "subtreeVersion": "91680",
                                "siblingHash": "DBCF1D6AA6945A734BA68A427101C1A447A96A5900C477B61207AA57C79D4CBB",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 6,
                                "subtreeSize": "48",
                                "subtreeVersion": "91685",
                                "siblingHash": "D31AB6A5A0E1D8978FE62D82E7CB6E9FA7F135D1790B0DE9C581465BB6637977",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 7,
                                "subtreeSize": "81",
                                "subtreeVersion": "91686",
                                "siblingHash": "C8C4F739B7B15E9F9D5F22C805325480C88B30D5BD62CAA70AED7A5689CAA087",
                            },
                            {
                                "isDataOnRight": True,
                                "subtreeHeight": 9,
                                "subtreeSize": "171",
                                "subtreeVersion": "2614469",
                                "siblingHash": "4C0CC11C7A96177AB5FE9ED4C0A8DFE0C4113DACDF45D5457B12C992F105C754",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 10,
                                "subtreeSize": "303",
                                "subtreeVersion": "2614469",
                                "siblingHash": "FDE122DF012F0E1730A144AF27BC91F49D91E2D9A6E4E3DA7891678C76A5AB8E",
                            },
                            {
                                "isDataOnRight": True,
                                "subtreeHeight": 11,
                                "subtreeSize": "482",
                                "subtreeVersion": "2614469",
                                "siblingHash": "A66E2F51BEE3F3DCEF5EBE12AD4E22D70E86C4D9963FBBFDE0F1ECF091CE8F23",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 12,
                                "subtreeSize": "1082",
                                "subtreeVersion": "2614469",
                                "siblingHash": "A2A8A0F87C99158FDF1AC783A9621AFD5AB6CB6770FE0D439218F4C2A98C90F4",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 13,
                                "subtreeSize": "2247",
                                "subtreeVersion": "2614469",
                                "siblingHash": "F0C8DDC85E699BF539549666D76F7C27B29ABD424A39E788E0ECECCE9887B9F6",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 14,
                                "subtreeSize": "4641",
                                "subtreeVersion": "2614469",
                                "siblingHash": "542C42BDBD667A789DAECF06275430D84255D26BDEAFF9316BB57535F0752DB4",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 15,
                                "subtreeSize": "8204",
                                "subtreeVersion": "2614469",
                                "siblingHash": "E7318E9D0DE8EC2E76FE5E5F26F79506298E3E5A41F32BE15BF15DB6EF6FBD3D",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 16,
                                "subtreeSize": "17284",
                                "subtreeVersion": "2614469",
                                "siblingHash": "C68AD7735CE9B53E0867482A3485C50F766B42647EBE1AEAF88D2CC34A01F3CB",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 17,
                                "subtreeSize": "61661",
                                "subtreeVersion": "2614469",
                                "siblingHash": "35E14098A506809C579606BD24F5403F5E44717426BCA91904010D346846D2F8",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 18,
                                "subtreeSize": "161842",
                                "subtreeVersion": "2614469",
                                "siblingHash": "7DA3D54285DF3952B8EAC5BF7878E193EA7D1817AD3EC8424252E8685EC2DF2D",
                            },
                            {
                                "isDataOnRight": False,
                                "subtreeHeight": 20,
                                "subtreeSize": "524812",
                                "subtreeVersion": "2622111",
                                "siblingHash": "7480F798043CC2C318E17A69CEB946798ECDF06D0A811121E3605748993151F4",
                            },
                            {
                                "isDataOnRight": True,
                                "subtreeHeight": 22,
                                "subtreeSize": "1112959",
                                "subtreeVersion": "2622112",
                                "siblingHash": "10289ED66A3FB336E0D85C1335AE96EBFD336E492CCCD2E42C95565462403CBE",
                            },
                            {
                                "isDataOnRight": True,
                                "subtreeHeight": 23,
                                "subtreeSize": "2760113",
                                "subtreeVersion": "2622112",
                                "siblingHash": "108EE3C758DBE110E53EAD06691DFA606F18352D3554A636831A55637EC9B8ED",
                            },
                            {
                                "isDataOnRight": True,
                                "subtreeHeight": 24,
                                "subtreeSize": "4238613",
                                "subtreeVersion": "2622112",
                                "siblingHash": "112E6BA96C48295E6F4D726B6921073B4270D2D1F11071656FF9FB117E7A130D",
                            },
                            {
                                "isDataOnRight": True,
                                "subtreeHeight": 25,
                                "subtreeSize": "6810730",
                                "subtreeVersion": "2622112",
                                "siblingHash": "273AB1005D741D5A4C805BF3C2D19E4E393513A85EE2EA103878EB1CFE554A81",
                            },
                            {
                                "isDataOnRight": True,
                                "subtreeHeight": 26,
                                "subtreeSize": "9337459",
                                "subtreeVersion": "2622113",
                                "siblingHash": "2711DD98462638BBDC666741479FC22E4E564B0DF0EFC728E82CF507FD474D62",
                            },
                        ],
                    },
                    "blockRelayProof": {
                        "multiStoreProof": {
                            "accToGovStoresMerkleHash": "BDC012A2E472BA5E31F95B9351F688AAA3B2339025BCF4943FB09F41D1E02D00",
                            "mainAndMintStoresMerkleHash": "49F8C5BCCFCD54845D53311BBFABDC79BD731F51B5DD5CCD3C5FECE3E31D943C",
                            "oracleIAVLStateHash": "E3E4410A29C6627A57866F951FFF04ECA7E601D5922BEE60A67A8730EDC299CB",
                            "paramsStoresMerkleHash": "9C3619FFC762ED94F1C71E82C5EC1AE0C7373554B69D847DB73703C18FF761D5",
                            "slashingToUpgradeStoresMerkleHash": "AFF5ED6925C982DC83F69EDA0C61B742E618B22D25359DF35F06711674307CD4",
                        },
                        "blockHeaderMerkleParts": {
                            "versionAndChainIdHash": "4FA9CA1048D3F4BAA282C89C96BD4259C5BBFDF9839215502B59F40C37D3B8B4",
                            "timeHash": "27EC75198A9D498AA614783616E4A446E122982A4D2FEEAAAAE1771193D83D70",
                            "lastBlockIDAndOther": "C2F1569086965DD3C39BC0C8AE058DA9AE8E80619354C2BBD3BB92D853A672BD",
                            "nextValidatorHashAndConsensusHash": "1D4396E9A5F6F0980F99298C49A143E179A12E982542D210B57DA9D140DF1543",
                            "lastResultsHash": "AA3C7CBEFF135291E6415ECA2528FC98D263B120C67BCECD8D8CCD3253EFECC1",
                            "evidenceAndProposerHash": "9B04008FE8D23B09C9C6AD1CFB529FD0220666B354233B7EA2E57FF835986319",
                        },
                        "signatures": [
                            {
                                "r": "24F7CEE7BB8498F11AE9CBC32212F0372F020D814137E5D467E98500EEB8E171",
                                "s": "7DCA2C2E855F7EA4633FAAD81557EF3FF6C2EE5AE03C9360BD0F79C9BD0C6F6F",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510FFE5F1C501321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "904DEBD1FD35AC84E9570F41E1A45DD71EFF655302395AD3C0A39592C6C24284",
                                "s": "6F044D31879B96D8AB9B726C21BFAC619698C87B1395105784F606DFE85BF58B",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510FEF7E3AE01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "40881F9A64150520E6D113C656F02223710661B238E601FB215E7628DE4CF709",
                                "s": "43F7AA3EE1D64F90FA82F5F48597BBBEB9F694EB5D4AFDDF462C2A263306BB8D",
                                "v": 28,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051094CE83AF01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "8753C86D9469A92939C9E1D37D21C4C8B9D5494C4D5937B48FEB9566AEB26252",
                                "s": "460738CA0AF439B8806A046F0C7AC6AE28B069EFB0F9C1766F4B9180C7ED5208",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510F880FBDA01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "72C8E7946F363F7F8D8175F33A04DBAD9E41ACA0CABBBBD3FA2FE8B3A5AD63AB",
                                "s": "6E50826A82D78FE4E1F1F917B8A01984D71F2A1F57AF06CFD43C52965CED6C7F",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510AECF9CAB01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "2BEDB913785EAB0C46ED5E105847A420DB1E3953586865DF6E5D18D06D806717",
                                "s": "3C3EC79018045736AACB5CF494EE2AEA9DD2E9216DCFD9F0B47EC808F978073B",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051092ACE6D101321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "BECFAD3E321B537D2DD35407EBBF769E178F67B395084063D5B940BF902E332B",
                                "s": "78A7898C5B241BDF5577087ABA56657BA72939F648C11572E7E2963424B96C8C",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510C9A0A69A01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "66426A3BBB61CACCC0548B16775146587A11259E98250D939C621A2D5B836FD4",
                                "s": "5993360AFE292848ACE750BDBDD31EC8586C23420333C3DB02D4650AA53A717C",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510A492D7AC01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "C4BB6F40D48B1BC3B535F7C29277DD1A0F37D8EB51A6895F1E92AA226CC90332",
                                "s": "4B774F7AAFBF167ED14C3A0E03CEA78D719ABA89B90F42387D85936C0EEFBE49",
                                "v": 28,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510CA9FA5C601321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "876D7A6006B01351E875747E8D9E4FC7B8FDB48EDA53D4C3B2344276414F31AA",
                                "s": "25F39CC6DC981FAF2AC4E7906BB696792655ED2731C0DF73CB1AE67E055FEF5D",
                                "v": 28,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05109A9E88DC01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "520A9A07615D170B69C3A7F877B44378EF9A8106AAF45F9CAF66C3984D8EDB73",
                                "s": "0E866B9FCF41BF3334ED076400E646E3F14770F7339246C77557F36B5A0E0F6D",
                                "v": 28,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05108EF3D7B301321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "157E5D56A8A54634FC5DB54BB761ED7226064F9FDDE834F680329393D5C0B93F",
                                "s": "46EBF239659F1D28C4738A1FD64B275A7ACA47D967DC0D5F1E3AE9C79814757C",
                                "v": 28,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05108387DAAF01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "4161D1D6C54C4F9757ACB1D6AFE7A3175D468EF7BC64B9E2E5BD5233DE44204A",
                                "s": "303C177938E1683AE94462E3558EC35117A1031C368C1C93C5E48D837695CFEB",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510D18BBCCD01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "287E8CC01F7F9CB6DB1787190F43D7BC745DF78294535CB67F97908A89B40EF7",
                                "s": "38218F52C59CF86507A20A3B3970DB3183767DA16CE63FD3B196493B6EEB0FC0",
                                "v": 28,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051089FAEECC01321462616E642D6775616E79752D746573746E657433",
                            },
                            {
                                "r": "C235AA2F27ADF368C424F7B3B4C0BA766309DBDD9417D58F85F54FCA452FC65A",
                                "s": "761A52774F57A5AB866C291A09B2FD1A4F9309E63F5E09F09C43785E95257724",
                                "v": 27,
                                "signedPrefixSuffix": "79080211A20228000000000022480A20",
                                "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510E9ACB19C01321462616E642D6775616E79752D746573746E657433",
                            },
                        ],
                    },
                },
                "evmProofBytes": "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000018A0000000000000000000000000000000000000000000000000000000000000184000000000000000000000000000000000000000000000000000000000002802A2BDC012A2E472BA5E31F95B9351F688AAA3B2339025BCF4943FB09F41D1E02D0049F8C5BCCFCD54845D53311BBFABDC79BD731F51B5DD5CCD3C5FECE3E31D943CE3E4410A29C6627A57866F951FFF04ECA7E601D5922BEE60A67A8730EDC299CB9C3619FFC762ED94F1C71E82C5EC1AE0C7373554B69D847DB73703C18FF761D5AFF5ED6925C982DC83F69EDA0C61B742E618B22D25359DF35F06711674307CD44FA9CA1048D3F4BAA282C89C96BD4259C5BBFDF9839215502B59F40C37D3B8B427EC75198A9D498AA614783616E4A446E122982A4D2FEEAAAAE1771193D83D70C2F1569086965DD3C39BC0C8AE058DA9AE8E80619354C2BBD3BB92D853A672BD1D4396E9A5F6F0980F99298C49A143E179A12E982542D210B57DA9D140DF1543AA3C7CBEFF135291E6415ECA2528FC98D263B120C67BCECD8D8CCD3253EFECC19B04008FE8D23B09C9C6AD1CFB529FD0220666B354233B7EA2E57FF83598631900000000000000000000000000000000000000000000000000000000000001A0000000000000000000000000000000000000000000000000000000000000000F00000000000000000000000000000000000000000000000000000000000001E0000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000004A00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000076000000000000000000000000000000000000000000000000000000000000008C00000000000000000000000000000000000000000000000000000000000000A200000000000000000000000000000000000000000000000000000000000000B800000000000000000000000000000000000000000000000000000000000000CE00000000000000000000000000000000000000000000000000000000000000E400000000000000000000000000000000000000000000000000000000000000FA00000000000000000000000000000000000000000000000000000000000001100000000000000000000000000000000000000000000000000000000000000126000000000000000000000000000000000000000000000000000000000000013C0000000000000000000000000000000000000000000000000000000000000152024F7CEE7BB8498F11AE9CBC32212F0372F020D814137E5D467E98500EEB8E1717DCA2C2E855F7EA4633FAAD81557EF3FF6C2EE5AE03C9360BD0F79C9BD0C6F6F000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510FFE5F1C501321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000904DEBD1FD35AC84E9570F41E1A45DD71EFF655302395AD3C0A39592C6C242846F044D31879B96D8AB9B726C21BFAC619698C87B1395105784F606DFE85BF58B000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510FEF7E3AE01321462616E642D6775616E79752D746573746E6574330000000000000000000000000000000000000000000040881F9A64150520E6D113C656F02223710661B238E601FB215E7628DE4CF70943F7AA3EE1D64F90FA82F5F48597BBBEB9F694EB5D4AFDDF462C2A263306BB8D000000000000000000000000000000000000000000000000000000000000001C00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051094CE83AF01321462616E642D6775616E79752D746573746E657433000000000000000000000000000000000000000000008753C86D9469A92939C9E1D37D21C4C8B9D5494C4D5937B48FEB9566AEB26252460738CA0AF439B8806A046F0C7AC6AE28B069EFB0F9C1766F4B9180C7ED5208000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510F880FBDA01321462616E642D6775616E79752D746573746E6574330000000000000000000000000000000000000000000072C8E7946F363F7F8D8175F33A04DBAD9E41ACA0CABBBBD3FA2FE8B3A5AD63AB6E50826A82D78FE4E1F1F917B8A01984D71F2A1F57AF06CFD43C52965CED6C7F000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510AECF9CAB01321462616E642D6775616E79752D746573746E657433000000000000000000000000000000000000000000002BEDB913785EAB0C46ED5E105847A420DB1E3953586865DF6E5D18D06D8067173C3EC79018045736AACB5CF494EE2AEA9DD2E9216DCFD9F0B47EC808F978073B000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051092ACE6D101321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000BECFAD3E321B537D2DD35407EBBF769E178F67B395084063D5B940BF902E332B78A7898C5B241BDF5577087ABA56657BA72939F648C11572E7E2963424B96C8C000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510C9A0A69A01321462616E642D6775616E79752D746573746E6574330000000000000000000000000000000000000000000066426A3BBB61CACCC0548B16775146587A11259E98250D939C621A2D5B836FD45993360AFE292848ACE750BDBDD31EC8586C23420333C3DB02D4650AA53A717C000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510A492D7AC01321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000C4BB6F40D48B1BC3B535F7C29277DD1A0F37D8EB51A6895F1E92AA226CC903324B774F7AAFBF167ED14C3A0E03CEA78D719ABA89B90F42387D85936C0EEFBE49000000000000000000000000000000000000000000000000000000000000001C00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510CA9FA5C601321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000876D7A6006B01351E875747E8D9E4FC7B8FDB48EDA53D4C3B2344276414F31AA25F39CC6DC981FAF2AC4E7906BB696792655ED2731C0DF73CB1AE67E055FEF5D000000000000000000000000000000000000000000000000000000000000001C00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05109A9E88DC01321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000520A9A07615D170B69C3A7F877B44378EF9A8106AAF45F9CAF66C3984D8EDB730E866B9FCF41BF3334ED076400E646E3F14770F7339246C77557F36B5A0E0F6D000000000000000000000000000000000000000000000000000000000000001C00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05108EF3D7B301321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000157E5D56A8A54634FC5DB54BB761ED7226064F9FDDE834F680329393D5C0B93F46EBF239659F1D28C4738A1FD64B275A7ACA47D967DC0D5F1E3AE9C79814757C",
            },
        },
        status_code=200,
    )

    assert client.get_request_evm_proof_by_request_id(1) == EVMProof(
        json_proof={
            "blockHeight": "2622114",
            "oracleDataProof": {
                "requestPacket": {
                    "client_id": "test",
                    "oracle_script_id": "1",
                    "calldata": "AAAABGZhc3Q=",
                    "ask_count": "4",
                    "min_count": "3",
                },
                "responsePacket": {
                    "client_id": "test",
                    "request_id": "1",
                    "ans_count": "4",
                    "request_time": "1600357375",
                    "resolve_time": "1600357377",
                    "resolve_status": 1,
                    "result": "AAAAAAAAEYA=",
                },
                "version": "624",
                "merklePaths": [
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 1,
                        "subtreeSize": "2",
                        "subtreeVersion": "818",
                        "siblingHash": "5744F905BEA848192798B1D1C624C65E4CA5EF6E964F37788EBE6EF49C67A6B9",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 2,
                        "subtreeSize": "4",
                        "subtreeVersion": "12005",
                        "siblingHash": "370FB5ECE9C7F3C742CC6F061E3F28C79063233442A0A7718AC7DAC1185DB625",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 3,
                        "subtreeSize": "8",
                        "subtreeVersion": "17879",
                        "siblingHash": "83C97412076743EEB622FCBE5877481401D2C9BC621950112B94BE4F7C7D7A0D",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 4,
                        "subtreeSize": "16",
                        "subtreeVersion": "18660",
                        "siblingHash": "C183A0D2A4614B7EBAD2DC8C2B2B7C94F42971FC47AD012FB6FEFE9AF7AC1C90",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 5,
                        "subtreeSize": "32",
                        "subtreeVersion": "91680",
                        "siblingHash": "DBCF1D6AA6945A734BA68A427101C1A447A96A5900C477B61207AA57C79D4CBB",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 6,
                        "subtreeSize": "48",
                        "subtreeVersion": "91685",
                        "siblingHash": "D31AB6A5A0E1D8978FE62D82E7CB6E9FA7F135D1790B0DE9C581465BB6637977",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 7,
                        "subtreeSize": "81",
                        "subtreeVersion": "91686",
                        "siblingHash": "C8C4F739B7B15E9F9D5F22C805325480C88B30D5BD62CAA70AED7A5689CAA087",
                    },
                    {
                        "isDataOnRight": True,
                        "subtreeHeight": 9,
                        "subtreeSize": "171",
                        "subtreeVersion": "2614469",
                        "siblingHash": "4C0CC11C7A96177AB5FE9ED4C0A8DFE0C4113DACDF45D5457B12C992F105C754",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 10,
                        "subtreeSize": "303",
                        "subtreeVersion": "2614469",
                        "siblingHash": "FDE122DF012F0E1730A144AF27BC91F49D91E2D9A6E4E3DA7891678C76A5AB8E",
                    },
                    {
                        "isDataOnRight": True,
                        "subtreeHeight": 11,
                        "subtreeSize": "482",
                        "subtreeVersion": "2614469",
                        "siblingHash": "A66E2F51BEE3F3DCEF5EBE12AD4E22D70E86C4D9963FBBFDE0F1ECF091CE8F23",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 12,
                        "subtreeSize": "1082",
                        "subtreeVersion": "2614469",
                        "siblingHash": "A2A8A0F87C99158FDF1AC783A9621AFD5AB6CB6770FE0D439218F4C2A98C90F4",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 13,
                        "subtreeSize": "2247",
                        "subtreeVersion": "2614469",
                        "siblingHash": "F0C8DDC85E699BF539549666D76F7C27B29ABD424A39E788E0ECECCE9887B9F6",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 14,
                        "subtreeSize": "4641",
                        "subtreeVersion": "2614469",
                        "siblingHash": "542C42BDBD667A789DAECF06275430D84255D26BDEAFF9316BB57535F0752DB4",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 15,
                        "subtreeSize": "8204",
                        "subtreeVersion": "2614469",
                        "siblingHash": "E7318E9D0DE8EC2E76FE5E5F26F79506298E3E5A41F32BE15BF15DB6EF6FBD3D",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 16,
                        "subtreeSize": "17284",
                        "subtreeVersion": "2614469",
                        "siblingHash": "C68AD7735CE9B53E0867482A3485C50F766B42647EBE1AEAF88D2CC34A01F3CB",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 17,
                        "subtreeSize": "61661",
                        "subtreeVersion": "2614469",
                        "siblingHash": "35E14098A506809C579606BD24F5403F5E44717426BCA91904010D346846D2F8",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 18,
                        "subtreeSize": "161842",
                        "subtreeVersion": "2614469",
                        "siblingHash": "7DA3D54285DF3952B8EAC5BF7878E193EA7D1817AD3EC8424252E8685EC2DF2D",
                    },
                    {
                        "isDataOnRight": False,
                        "subtreeHeight": 20,
                        "subtreeSize": "524812",
                        "subtreeVersion": "2622111",
                        "siblingHash": "7480F798043CC2C318E17A69CEB946798ECDF06D0A811121E3605748993151F4",
                    },
                    {
                        "isDataOnRight": True,
                        "subtreeHeight": 22,
                        "subtreeSize": "1112959",
                        "subtreeVersion": "2622112",
                        "siblingHash": "10289ED66A3FB336E0D85C1335AE96EBFD336E492CCCD2E42C95565462403CBE",
                    },
                    {
                        "isDataOnRight": True,
                        "subtreeHeight": 23,
                        "subtreeSize": "2760113",
                        "subtreeVersion": "2622112",
                        "siblingHash": "108EE3C758DBE110E53EAD06691DFA606F18352D3554A636831A55637EC9B8ED",
                    },
                    {
                        "isDataOnRight": True,
                        "subtreeHeight": 24,
                        "subtreeSize": "4238613",
                        "subtreeVersion": "2622112",
                        "siblingHash": "112E6BA96C48295E6F4D726B6921073B4270D2D1F11071656FF9FB117E7A130D",
                    },
                    {
                        "isDataOnRight": True,
                        "subtreeHeight": 25,
                        "subtreeSize": "6810730",
                        "subtreeVersion": "2622112",
                        "siblingHash": "273AB1005D741D5A4C805BF3C2D19E4E393513A85EE2EA103878EB1CFE554A81",
                    },
                    {
                        "isDataOnRight": True,
                        "subtreeHeight": 26,
                        "subtreeSize": "9337459",
                        "subtreeVersion": "2622113",
                        "siblingHash": "2711DD98462638BBDC666741479FC22E4E564B0DF0EFC728E82CF507FD474D62",
                    },
                ],
            },
            "blockRelayProof": {
                "multiStoreProof": {
                    "accToGovStoresMerkleHash": "BDC012A2E472BA5E31F95B9351F688AAA3B2339025BCF4943FB09F41D1E02D00",
                    "mainAndMintStoresMerkleHash": "49F8C5BCCFCD54845D53311BBFABDC79BD731F51B5DD5CCD3C5FECE3E31D943C",
                    "oracleIAVLStateHash": "E3E4410A29C6627A57866F951FFF04ECA7E601D5922BEE60A67A8730EDC299CB",
                    "paramsStoresMerkleHash": "9C3619FFC762ED94F1C71E82C5EC1AE0C7373554B69D847DB73703C18FF761D5",
                    "slashingToUpgradeStoresMerkleHash": "AFF5ED6925C982DC83F69EDA0C61B742E618B22D25359DF35F06711674307CD4",
                },
                "blockHeaderMerkleParts": {
                    "versionAndChainIdHash": "4FA9CA1048D3F4BAA282C89C96BD4259C5BBFDF9839215502B59F40C37D3B8B4",
                    "timeHash": "27EC75198A9D498AA614783616E4A446E122982A4D2FEEAAAAE1771193D83D70",
                    "lastBlockIDAndOther": "C2F1569086965DD3C39BC0C8AE058DA9AE8E80619354C2BBD3BB92D853A672BD",
                    "nextValidatorHashAndConsensusHash": "1D4396E9A5F6F0980F99298C49A143E179A12E982542D210B57DA9D140DF1543",
                    "lastResultsHash": "AA3C7CBEFF135291E6415ECA2528FC98D263B120C67BCECD8D8CCD3253EFECC1",
                    "evidenceAndProposerHash": "9B04008FE8D23B09C9C6AD1CFB529FD0220666B354233B7EA2E57FF835986319",
                },
                "signatures": [
                    {
                        "r": "24F7CEE7BB8498F11AE9CBC32212F0372F020D814137E5D467E98500EEB8E171",
                        "s": "7DCA2C2E855F7EA4633FAAD81557EF3FF6C2EE5AE03C9360BD0F79C9BD0C6F6F",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510FFE5F1C501321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "904DEBD1FD35AC84E9570F41E1A45DD71EFF655302395AD3C0A39592C6C24284",
                        "s": "6F044D31879B96D8AB9B726C21BFAC619698C87B1395105784F606DFE85BF58B",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510FEF7E3AE01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "40881F9A64150520E6D113C656F02223710661B238E601FB215E7628DE4CF709",
                        "s": "43F7AA3EE1D64F90FA82F5F48597BBBEB9F694EB5D4AFDDF462C2A263306BB8D",
                        "v": 28,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051094CE83AF01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "8753C86D9469A92939C9E1D37D21C4C8B9D5494C4D5937B48FEB9566AEB26252",
                        "s": "460738CA0AF439B8806A046F0C7AC6AE28B069EFB0F9C1766F4B9180C7ED5208",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510F880FBDA01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "72C8E7946F363F7F8D8175F33A04DBAD9E41ACA0CABBBBD3FA2FE8B3A5AD63AB",
                        "s": "6E50826A82D78FE4E1F1F917B8A01984D71F2A1F57AF06CFD43C52965CED6C7F",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510AECF9CAB01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "2BEDB913785EAB0C46ED5E105847A420DB1E3953586865DF6E5D18D06D806717",
                        "s": "3C3EC79018045736AACB5CF494EE2AEA9DD2E9216DCFD9F0B47EC808F978073B",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051092ACE6D101321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "BECFAD3E321B537D2DD35407EBBF769E178F67B395084063D5B940BF902E332B",
                        "s": "78A7898C5B241BDF5577087ABA56657BA72939F648C11572E7E2963424B96C8C",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510C9A0A69A01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "66426A3BBB61CACCC0548B16775146587A11259E98250D939C621A2D5B836FD4",
                        "s": "5993360AFE292848ACE750BDBDD31EC8586C23420333C3DB02D4650AA53A717C",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510A492D7AC01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "C4BB6F40D48B1BC3B535F7C29277DD1A0F37D8EB51A6895F1E92AA226CC90332",
                        "s": "4B774F7AAFBF167ED14C3A0E03CEA78D719ABA89B90F42387D85936C0EEFBE49",
                        "v": 28,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510CA9FA5C601321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "876D7A6006B01351E875747E8D9E4FC7B8FDB48EDA53D4C3B2344276414F31AA",
                        "s": "25F39CC6DC981FAF2AC4E7906BB696792655ED2731C0DF73CB1AE67E055FEF5D",
                        "v": 28,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05109A9E88DC01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "520A9A07615D170B69C3A7F877B44378EF9A8106AAF45F9CAF66C3984D8EDB73",
                        "s": "0E866B9FCF41BF3334ED076400E646E3F14770F7339246C77557F36B5A0E0F6D",
                        "v": 28,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05108EF3D7B301321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "157E5D56A8A54634FC5DB54BB761ED7226064F9FDDE834F680329393D5C0B93F",
                        "s": "46EBF239659F1D28C4738A1FD64B275A7ACA47D967DC0D5F1E3AE9C79814757C",
                        "v": 28,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05108387DAAF01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "4161D1D6C54C4F9757ACB1D6AFE7A3175D468EF7BC64B9E2E5BD5233DE44204A",
                        "s": "303C177938E1683AE94462E3558EC35117A1031C368C1C93C5E48D837695CFEB",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510D18BBCCD01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "287E8CC01F7F9CB6DB1787190F43D7BC745DF78294535CB67F97908A89B40EF7",
                        "s": "38218F52C59CF86507A20A3B3970DB3183767DA16CE63FD3B196493B6EEB0FC0",
                        "v": 28,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051089FAEECC01321462616E642D6775616E79752D746573746E657433",
                    },
                    {
                        "r": "C235AA2F27ADF368C424F7B3B4C0BA766309DBDD9417D58F85F54FCA452FC65A",
                        "s": "761A52774F57A5AB866C291A09B2FD1A4F9309E63F5E09F09C43785E95257724",
                        "v": 27,
                        "signedPrefixSuffix": "79080211A20228000000000022480A20",
                        "signedDataSuffix": "12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510E9ACB19C01321462616E642D6775616E79752D746573746E657433",
                    },
                ],
            },
        },
        evm_proof_bytes=HexBytes(
            bytes.fromhex(
                "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000018A0000000000000000000000000000000000000000000000000000000000000184000000000000000000000000000000000000000000000000000000000002802A2BDC012A2E472BA5E31F95B9351F688AAA3B2339025BCF4943FB09F41D1E02D0049F8C5BCCFCD54845D53311BBFABDC79BD731F51B5DD5CCD3C5FECE3E31D943CE3E4410A29C6627A57866F951FFF04ECA7E601D5922BEE60A67A8730EDC299CB9C3619FFC762ED94F1C71E82C5EC1AE0C7373554B69D847DB73703C18FF761D5AFF5ED6925C982DC83F69EDA0C61B742E618B22D25359DF35F06711674307CD44FA9CA1048D3F4BAA282C89C96BD4259C5BBFDF9839215502B59F40C37D3B8B427EC75198A9D498AA614783616E4A446E122982A4D2FEEAAAAE1771193D83D70C2F1569086965DD3C39BC0C8AE058DA9AE8E80619354C2BBD3BB92D853A672BD1D4396E9A5F6F0980F99298C49A143E179A12E982542D210B57DA9D140DF1543AA3C7CBEFF135291E6415ECA2528FC98D263B120C67BCECD8D8CCD3253EFECC19B04008FE8D23B09C9C6AD1CFB529FD0220666B354233B7EA2E57FF83598631900000000000000000000000000000000000000000000000000000000000001A0000000000000000000000000000000000000000000000000000000000000000F00000000000000000000000000000000000000000000000000000000000001E0000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000004A00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000076000000000000000000000000000000000000000000000000000000000000008C00000000000000000000000000000000000000000000000000000000000000A200000000000000000000000000000000000000000000000000000000000000B800000000000000000000000000000000000000000000000000000000000000CE00000000000000000000000000000000000000000000000000000000000000E400000000000000000000000000000000000000000000000000000000000000FA00000000000000000000000000000000000000000000000000000000000001100000000000000000000000000000000000000000000000000000000000000126000000000000000000000000000000000000000000000000000000000000013C0000000000000000000000000000000000000000000000000000000000000152024F7CEE7BB8498F11AE9CBC32212F0372F020D814137E5D467E98500EEB8E1717DCA2C2E855F7EA4633FAAD81557EF3FF6C2EE5AE03C9360BD0F79C9BD0C6F6F000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510FFE5F1C501321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000904DEBD1FD35AC84E9570F41E1A45DD71EFF655302395AD3C0A39592C6C242846F044D31879B96D8AB9B726C21BFAC619698C87B1395105784F606DFE85BF58B000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510FEF7E3AE01321462616E642D6775616E79752D746573746E6574330000000000000000000000000000000000000000000040881F9A64150520E6D113C656F02223710661B238E601FB215E7628DE4CF70943F7AA3EE1D64F90FA82F5F48597BBBEB9F694EB5D4AFDDF462C2A263306BB8D000000000000000000000000000000000000000000000000000000000000001C00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051094CE83AF01321462616E642D6775616E79752D746573746E657433000000000000000000000000000000000000000000008753C86D9469A92939C9E1D37D21C4C8B9D5494C4D5937B48FEB9566AEB26252460738CA0AF439B8806A046F0C7AC6AE28B069EFB0F9C1766F4B9180C7ED5208000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510F880FBDA01321462616E642D6775616E79752D746573746E6574330000000000000000000000000000000000000000000072C8E7946F363F7F8D8175F33A04DBAD9E41ACA0CABBBBD3FA2FE8B3A5AD63AB6E50826A82D78FE4E1F1F917B8A01984D71F2A1F57AF06CFD43C52965CED6C7F000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510AECF9CAB01321462616E642D6775616E79752D746573746E657433000000000000000000000000000000000000000000002BEDB913785EAB0C46ED5E105847A420DB1E3953586865DF6E5D18D06D8067173C3EC79018045736AACB5CF494EE2AEA9DD2E9216DCFD9F0B47EC808F978073B000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE051092ACE6D101321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000BECFAD3E321B537D2DD35407EBBF769E178F67B395084063D5B940BF902E332B78A7898C5B241BDF5577087ABA56657BA72939F648C11572E7E2963424B96C8C000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510C9A0A69A01321462616E642D6775616E79752D746573746E6574330000000000000000000000000000000000000000000066426A3BBB61CACCC0548B16775146587A11259E98250D939C621A2D5B836FD45993360AFE292848ACE750BDBDD31EC8586C23420333C3DB02D4650AA53A717C000000000000000000000000000000000000000000000000000000000000001B00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510A492D7AC01321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000C4BB6F40D48B1BC3B535F7C29277DD1A0F37D8EB51A6895F1E92AA226CC903324B774F7AAFBF167ED14C3A0E03CEA78D719ABA89B90F42387D85936C0EEFBE49000000000000000000000000000000000000000000000000000000000000001C00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE0510CA9FA5C601321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000876D7A6006B01351E875747E8D9E4FC7B8FDB48EDA53D4C3B2344276414F31AA25F39CC6DC981FAF2AC4E7906BB696792655ED2731C0DF73CB1AE67E055FEF5D000000000000000000000000000000000000000000000000000000000000001C00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05109A9E88DC01321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000520A9A07615D170B69C3A7F877B44378EF9A8106AAF45F9CAF66C3984D8EDB730E866B9FCF41BF3334ED076400E646E3F14770F7339246C77557F36B5A0E0F6D000000000000000000000000000000000000000000000000000000000000001C00000000000000000000000000000000000000000000000000000000000000A000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000001079080211A20228000000000022480A2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004A12240A2098464229F9F19AD6D3D1ADAF40AE529F74D7B8B506BECEB00342B4BC511042C710012A0C08A4DBBCFE05108EF3D7B301321462616E642D6775616E79752D746573746E65743300000000000000000000000000000000000000000000157E5D56A8A54634FC5DB54BB761ED7226064F9FDDE834F680329393D5C0B93F46EBF239659F1D28C4738A1FD64B275A7ACA47D967DC0D5F1E3AE9C79814757C"
            )
        ),
    )
