from iconservice import *
from .pyobi import *

TAG = "CACHE_CONSUMER_MOCK"


class IBRIDGECACHE(InterfaceScore):
    @interface
    def get_latest_response(self, encoded_request: bytes) -> dict:
        pass


class CACHE_CONSUMER_MOCK(IconScoreBase):
    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        self.bridge_address = VarDB("bridge_address", db, value_type=Address)
        self.req_key_template = VarDB("req_template", db, value_type=bytes)
        self.res = VarDB("res", db, value_type=bytes)

    def on_install(self, bridge_address: Address, req_template: bytes) -> None:
        super().on_install()
        self.bridge_address.set(bridge_address)
        self.req_key_template.set(req_template)

    def on_update(self) -> None:
        super().on_update()

    @external(readonly=True)
    def get_bridge_address(self) -> Address:
        return self.bridge_address.get()

    @external(readonly=True)
    def get_request_key_template(self) -> dict:
        if not isinstance(self.req_key_template.get(), bytes):
            return None
        return PyObi(
            """
                {
                    client_id: string,
                    oracle_script_id: u64,
                    calldata: bytes,
                    ask_count: u64,
                    min_count: u64
                }
            """
        ).decode(self.req_key_template.get())

    @external(readonly=True)
    def get_res(self) -> dict:
        if not isinstance(self.res.get(), bytes):
            return None
        return PyObi(
            """
                {
                    client_id: string,
                    request_id: u64,
                    ans_count: u64,
                    request_time: u64,
                    resolve_time: u64,
                    resolve_status: u32,
                    result: bytes
                }
            """
        ).decode(self.res.get())

    @external
    def consume_cache(self) -> None:
        bridge = self.create_interface_score(self.bridge_address.get(), IBRIDGECACHE)
        res = bridge.get_latest_response(self.req_key_template.get())
        if not isinstance(res, dict):
            revert("RESPONSE_MUST_BE_DICT_BUT_GOT_{}".format(type(res)))
        self.res.set(
            PyObi(
                """
                    {
                        client_id: string,
                        request_id: u64,
                        ans_count: u64,
                        request_time: u64,
                        resolve_time: u64,
                        resolve_status: u32,
                        result: bytes
                    }
                """
            ).encode(res)
        )

