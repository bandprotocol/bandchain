import requests
import json
from pyband.transaction import Transaction
from dacite import from_dict

from .data import (
    Account,
    DataSource,
    OracleScript,
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
            raise ValueError(data["raw_log"])
        return TransactionSyncMode(tx_hash=data["txhash"])

    def send_tx_async_mode(self, data: dict) -> TransactionAsyncMode:
        data = self._post("/txs", json={"tx": data, "mode": "async"})
        return TransactionAsyncMode(tx_hash=data["txhash"])

    def send_tx_block_mode(self, data: dict) -> TransactionBlockMode:
        data = self._post("/txs", json={"tx": data, "mode": "block"})
        if "code" in data:
            raise ValueError(data["raw_log"])
        return TransactionBlockMode(
            height=data["height"],
            tx_hash=data["txhash"],
            gas_wanted=data["gas_wanted"],
            gas_used=data["gas_wanted"],
            raw_log=json.loads(data["raw_log"]),
        )

    def get_chain_id(self) -> str:
        return self._get("/bandchain/chain_id")["chain_id"]

    def get_latest_block(self) -> dict:
        return self._get("/blocks/latest")

    def get_account(self, address: str) -> Account:
        return from_dict(
            data_class=Account,
            data=self._get_result("/auth/accounts/{}".format(address))["value"],
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

    def get_latest_request(
        self, oid: int, calldata: bytes, min_count: int, ask_count: int
    ) -> RequestInfo:
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

    def get_reporters(self, validator: str) -> list:
        return self._get_result("/oracle/reporters/{}".format(validator))
