from iconservice import *
from .pyobi import *

TAG = "ReceiverMock"


class IBridge(InterfaceScore):
    @interface
    def relay_and_verify(self, proof: bytes) -> dict:
        pass


class ReceiverMock(IconScoreBase):
    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        self.bridge_address = VarDB("bridge_address", db, value_type=Address)
        self.req = VarDB("req", db, value_type=bytes)
        self.res = VarDB("res", db, value_type=bytes)

    def on_install(self, bridge_address: Address) -> None:
        super().on_install()
        self.bridge_address.set(bridge_address)

    def on_update(self) -> None:
        super().on_update()

    @external(readonly=True)
    def get_bridge_address(self) -> Address:
        return self.bridge_address.get()

    @external(readonly=True)
    def get_req(self) -> dict:
        if not isinstance(self.req.get(), bytes):
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
        ).decode(self.req.get())

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
    def relay_and_safe(self, proof: bytes) -> None:
        bridge = self.create_interface_score(self.bridge_address.get(), IBridge)
        packet = bridge.relay_and_verify(proof)
        self.req.set(
            PyObi(
                """
                    {
                        client_id: string,
                        oracle_script_id: u64,
                        calldata: bytes,
                        ask_count: u64,
                        min_count: u64
                    }
                """
            ).encode(packet["req"])
        )

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
            ).encode(packet["res"])
        )

