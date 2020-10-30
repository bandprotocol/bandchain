import pytest
from pyband.client import Client
from pyband.data import (
    Account,
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
)

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
                "hash": "0E414A4ECA163CA4524334BCEB81C349AFCF1AC03CFE381A37BA199D5C174F0F",
                "parts": {
                    "total": "1",
                    "hash": "BBF7CE3AC38167B9201079EB60E1705675A27C78A48AEC0BD0AA888E0781F5F6",
                },
            },
            "block": {
                "header": {
                    "version": {"block": "10", "app": "0"},
                    "chain_id": "bandchain",
                    "height": "653582",
                    "time": "2020-10-28T09:54:25.830281151Z",
                    "last_block_id": {
                        "hash": "F0878F8A875B65477B019CD02FEC5666A30563106DAD212B9E508AF49E82C965",
                        "parts": {
                            "total": "1",
                            "hash": "7DDC15AE1DDD5C2B5A4148F5092DDA2382CC3F87BC1A9A7433DBE5476D15425D",
                        },
                    },
                    "last_commit_hash": "F499D6AFA34282E1054ED4E799253038799145A1FF40539F45E3121EA3B8AFA3",
                    "data_hash": "",
                    "validators_hash": "C0368C590DE186048CDBEE815E04E131564AA4AE98E875DDB57D76666A7E0023",
                    "next_validators_hash": "C0368C590DE186048CDBEE815E04E131564AA4AE98E875DDB57D76666A7E0023",
                    "consensus_hash": "0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45",
                    "app_hash": "15BEA68127772B814CD9FECF07033C85FEC5E6C0DCBCF6B143EB2DF446B576D8",
                    "last_results_hash": "",
                    "evidence_hash": "",
                    "proposer_address": "F23391B5DBF982E37FB7DADEA64AAE21CAE4C172",
                },
                "data": {"txs": None},
                "evidence": {"evidence": None},
                "last_commit": {
                    "height": "653581",
                    "round": "0",
                    "block_id": {
                        "hash": "F0878F8A875B65477B019CD02FEC5666A30563106DAD212B9E508AF49E82C965",
                        "parts": {
                            "total": "1",
                            "hash": "7DDC15AE1DDD5C2B5A4148F5092DDA2382CC3F87BC1A9A7433DBE5476D15425D",
                        },
                    },
                    "signatures": [
                        {
                            "block_id_flag": 2,
                            "validator_address": "5179B0BB203248E03D2A1342896133B5C58E1E44",
                            "timestamp": "2020-10-28T09:54:25.830281151Z",
                            "signature": "BAALWyNf36SCTXKpRuVhgCJ72ObPUeeEl4qRRacaiaEy4PlZffZQRqSGqwOqsZ6264Dq4Rh09QkzPmgeI3GT4w==",
                        },
                        {
                            "block_id_flag": 2,
                            "validator_address": "BDB6A0728C8DFE2124536F16F2BA428FE767A8F9",
                            "timestamp": "2020-10-28T09:54:25.82790319Z",
                            "signature": "CUbjHaDrhP9iwgZu+fAy6OBqh2Jj03QlC4hq6QC8DTI97RBENYLh6BPmI3AZx0OIlMLgFIQtuDK4UvlBvjdnZg==",
                        },
                        {
                            "block_id_flag": 2,
                            "validator_address": "F0C23921727D869745C4F9703CF33996B1D2B715",
                            "timestamp": "2020-10-28T09:54:25.831665017Z",
                            "signature": "dQZocP16dAXxhDRfqUyYI9n/bHcI8bhYQZiAQ1fNBDok4MRfV06EBAMj9+UR0tYSllUWsBeXsFt14Lgtt/yEgQ==",
                        },
                        {
                            "block_id_flag": 2,
                            "validator_address": "F23391B5DBF982E37FB7DADEA64AAE21CAE4C172",
                            "timestamp": "2020-10-28T09:54:25.827893959Z",
                            "signature": "1tTNNaJQSxZsVKKyHGlPwzwQul1O0egBoGaAj9gcMzB8Jt2DFrz71bjLxF1PSPmzqOOI1OH2CUFBC9Qycrr/sw==",
                        },
                    ],
                },
            },
        },
        status_code=200,
    )

    assert client.get_latest_block() == {
        "block_id": {
            "hash": "0E414A4ECA163CA4524334BCEB81C349AFCF1AC03CFE381A37BA199D5C174F0F",
            "parts": {
                "total": "1",
                "hash": "BBF7CE3AC38167B9201079EB60E1705675A27C78A48AEC0BD0AA888E0781F5F6",
            },
        },
        "block": {
            "header": {
                "version": {"block": "10", "app": "0"},
                "chain_id": "bandchain",
                "height": "653582",
                "time": "2020-10-28T09:54:25.830281151Z",
                "last_block_id": {
                    "hash": "F0878F8A875B65477B019CD02FEC5666A30563106DAD212B9E508AF49E82C965",
                    "parts": {
                        "total": "1",
                        "hash": "7DDC15AE1DDD5C2B5A4148F5092DDA2382CC3F87BC1A9A7433DBE5476D15425D",
                    },
                },
                "last_commit_hash": "F499D6AFA34282E1054ED4E799253038799145A1FF40539F45E3121EA3B8AFA3",
                "data_hash": "",
                "validators_hash": "C0368C590DE186048CDBEE815E04E131564AA4AE98E875DDB57D76666A7E0023",
                "next_validators_hash": "C0368C590DE186048CDBEE815E04E131564AA4AE98E875DDB57D76666A7E0023",
                "consensus_hash": "0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45",
                "app_hash": "15BEA68127772B814CD9FECF07033C85FEC5E6C0DCBCF6B143EB2DF446B576D8",
                "last_results_hash": "",
                "evidence_hash": "",
                "proposer_address": "F23391B5DBF982E37FB7DADEA64AAE21CAE4C172",
            },
            "data": {"txs": None},
            "evidence": {"evidence": None},
            "last_commit": {
                "height": "653581",
                "round": "0",
                "block_id": {
                    "hash": "F0878F8A875B65477B019CD02FEC5666A30563106DAD212B9E508AF49E82C965",
                    "parts": {
                        "total": "1",
                        "hash": "7DDC15AE1DDD5C2B5A4148F5092DDA2382CC3F87BC1A9A7433DBE5476D15425D",
                    },
                },
                "signatures": [
                    {
                        "block_id_flag": 2,
                        "validator_address": "5179B0BB203248E03D2A1342896133B5C58E1E44",
                        "timestamp": "2020-10-28T09:54:25.830281151Z",
                        "signature": "BAALWyNf36SCTXKpRuVhgCJ72ObPUeeEl4qRRacaiaEy4PlZffZQRqSGqwOqsZ6264Dq4Rh09QkzPmgeI3GT4w==",
                    },
                    {
                        "block_id_flag": 2,
                        "validator_address": "BDB6A0728C8DFE2124536F16F2BA428FE767A8F9",
                        "timestamp": "2020-10-28T09:54:25.82790319Z",
                        "signature": "CUbjHaDrhP9iwgZu+fAy6OBqh2Jj03QlC4hq6QC8DTI97RBENYLh6BPmI3AZx0OIlMLgFIQtuDK4UvlBvjdnZg==",
                    },
                    {
                        "block_id_flag": 2,
                        "validator_address": "F0C23921727D869745C4F9703CF33996B1D2B715",
                        "timestamp": "2020-10-28T09:54:25.831665017Z",
                        "signature": "dQZocP16dAXxhDRfqUyYI9n/bHcI8bhYQZiAQ1fNBDok4MRfV06EBAMj9+UR0tYSllUWsBeXsFt14Lgtt/yEgQ==",
                    },
                    {
                        "block_id_flag": 2,
                        "validator_address": "F23391B5DBF982E37FB7DADEA64AAE21CAE4C172",
                        "timestamp": "2020-10-28T09:54:25.827893959Z",
                        "signature": "1tTNNaJQSxZsVKKyHGlPwzwQul1O0egBoGaAj9gcMzB8Jt2DFrz71bjLxF1PSPmzqOOI1OH2CUFBC9Qycrr/sw==",
                    },
                ],
            },
        },
    }


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

    assert client.get_account("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte") == Account(
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
                    "requested_validators": [
                        "bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"
                    ],
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
        "{}/oracle/reporters/bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre".format(
            TEST_RPC
        ),
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

    assert client.get_reporters(
        "bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre"
    ) == [
        "band1yyv5jkqaukq0ajqn7vhkyhpff7h6e99ja7gvwg",
        "band19nf0sqnjycnvpexlxs6hjz9qrhhlhsu9pdty0r",
        "band1fndxcmg0h5vhr8cph7gryryqfn9yqp90lysjtm",
        "band10ec0p96j60duce5qagju5axuja0rj8luqrzl0k",
        "band15pm9scujgkpwpy2xa2j53tvs9ylunjn0g73a9s",
        "band1cehe3sxk7f4rmvwdf6lxh3zexen7fn02zyltwy",
    ]

