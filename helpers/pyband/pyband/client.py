import requests
from dacite import from_dict

from .data import Account, DataSource, OracleScript, RequestInfo, DACITE_CONFIG


class Client(object):
    def __init__(self, rpc_url: str) -> None:
        self.rpc_url = rpc_url

    def _get(self, path, **kwargs):
        return requests.get(self.rpc_url + path, **kwargs).json()["result"]

    def send_tx(self, data: dict) -> dict:
        return requests.post(self.rpc_url + "/txs", json=data).json()

    def get_chain_id(self) -> str:
        genesis = requests.get(self.rpc_url + "/bandchain/genesis").json()
        return genesis["chain_id"]

    def get_latest_block(self) -> dict:
        return requests.get(self.rpc_url + "/blocks/latest").json()

    def get_account(self, address: str) -> Account:
        return from_dict(
            data_class=Account,
            data=self._get("/auth/accounts/{}".format(address))["value"],
            config=DACITE_CONFIG,
        )

    def get_data_source(self, id: int) -> DataSource:
        return from_dict(
            data_class=DataSource,
            data=self._get("/oracle/data_sources/{}".format(id)),
            config=DACITE_CONFIG,
        )

    def get_oracle_script(self, id: int) -> OracleScript:
        return from_dict(
            data_class=OracleScript,
            data=self._get("/oracle/oracle_scripts/{}".format(id)),
            config=DACITE_CONFIG,
        )

    def get_request_by_id(self, id: int) -> RequestInfo:
        return from_dict(
            data_class=RequestInfo,
            data=self._get("/oracle/requests/{}".format(id)),
            config=DACITE_CONFIG,
        )

    def get_latest_request(
        self, oid: int, calldata: bytes, min_count: int, ask_count: int
    ) -> RequestInfo:
        return from_dict(
            data_class=RequestInfo,
            data=self._get(
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
        return self._get("/oracle/reporters/{}".format(validator))

