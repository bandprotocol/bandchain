import requests

from dacite import from_dict
from .wallet import Address
from typing import List
from .data import (
    Account,
    Block,
    DataSource,
    OracleScript,
    HexBytes,
    RequestInfo,
    DACITE_CONFIG,
    TransactionSyncMode,
    TransactionAsyncMode,
    TransactionBlockMode,
)


class Client(object):
    def __init__(self, rpc_url: str) -> None:
        self.rpc_url = rpc_url

    def _get(self, path, **kwargs):
        r = requests.get(self.rpc_url + path, **kwargs)
        r.raise_for_status()
        return r.json()

    def _post(self, path, **kwargs):
        r = requests.post(self.rpc_url + path, **kwargs)
        r.raise_for_status()
        return r.json()

    def _get_result(self, path, **kwargs):
        return self._get(path, **kwargs)["result"]

    def send_tx_sync_mode(self, data: dict) -> TransactionSyncMode:
        data = self._post("/txs", json={"tx": data, "mode": "sync"})
        if "code" in data:
            code = int(data["code"])
            error_log = data["raw_log"]
        else:
            code = 0
            error_log = None

        return TransactionSyncMode(
            tx_hash=HexBytes(bytes.fromhex(data["txhash"])),
            code=code,
            error_log=error_log,
        )

    def send_tx_async_mode(self, data: dict) -> TransactionAsyncMode:
        data = self._post("/txs", json={"tx": data, "mode": "async"})
        return TransactionAsyncMode(tx_hash=HexBytes(bytes.fromhex(data["txhash"])))

    def send_tx_block_mode(self, data: dict) -> TransactionBlockMode:
        data = self._post("/txs", json={"tx": data, "mode": "block"})
        if "code" in data:
            code = int(data["code"])
            error_log = data["raw_log"]
            log = []
        else:
            code = 0
            error_log = None
            log = data["logs"]

        return TransactionBlockMode(
            height=int(data["height"]),
            tx_hash=HexBytes(bytes.fromhex(data["txhash"])),
            gas_wanted=int(data["gas_wanted"]),
            gas_used=int(data["gas_wanted"]),
            code=code,
            log=log,
            error_log=error_log,
        )

    def get_chain_id(self) -> str:
        return self._get("/bandchain/chain_id")["chain_id"]

    def get_latest_block(self) -> Block:
        return from_dict(
            data_class=Block,
            data=self._get("/blocks/latest"),
            config=DACITE_CONFIG,
        )

    def get_account(self, address: Address) -> Account:
        return from_dict(
            data_class=Account,
            data=self._get_result("/auth/accounts/{}".format(address.to_acc_bech32()))["value"],
            config=DACITE_CONFIG,
        )

    def get_data_source(self, id: int) -> DataSource:
        return from_dict(
            data_class=DataSource,
            data=self._get_result("/oracle/data_sources/{}".format(id)),
            config=DACITE_CONFIG,
        )

    def get_oracle_script(self, id: int) -> OracleScript:
        return from_dict(
            data_class=OracleScript,
            data=self._get_result("/oracle/oracle_scripts/{}".format(id)),
            config=DACITE_CONFIG,
        )

    def get_request_by_id(self, id: int) -> RequestInfo:
        return from_dict(
            data_class=RequestInfo,
            data=self._get_result("/oracle/requests/{}".format(id)),
            config=DACITE_CONFIG,
        )

    def get_latest_request(self, oid: int, calldata: bytes, min_count: int, ask_count: int) -> RequestInfo:
        return from_dict(
            data_class=RequestInfo,
            data=self._get_result(
                "/oracle/request_search",
                params={
                    "oid": oid,
                    "calldata": calldata.hex(),
                    "min_count": min_count,
                    "ask_count": ask_count,
                },
            ),
            config=DACITE_CONFIG,
        )

    def get_reporters(self, validator: str) -> List[str]:
        return self._get_result("/oracle/reporters/{}".format(validator))

    def get_price_symbols(self, min_count: int, ask_count: int) -> List[str]:
        return self._get_result(
            "/oracle/price_symbols",
            params={
                "min_count": min_count,
                "ask_count": ask_count,
            },
        )
    def get_request_id_by_tx_hash(self, tx_hash: HexBytes) -> List[int]:
        msgs = self._get("/txs/{}".format(tx_hash.hex()))["logs"]
        request_ids = []
        for msg in msgs:
            request_event = [event for event in msg["events"] if event["type"] == "request"]
            if len(request_event) == 1:
                attrs = request_event[0]["attributes"]
                attr_id = [attr for attr in attrs if attr["key"] == "id"]
                if len(attr_id) == 1:
                    request_id = attr_id[0]["value"]
                    request_ids.append(int(request_id))
        if len(request_ids) == 0:
            raise ValueError("There is no request message in this tx")
        return request_ids
