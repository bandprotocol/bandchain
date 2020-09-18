from iconservice import *
from .pyobi import *

TAG = "StdReferenceProxy"

# obi
REQUEST_PACKET = PyObi(
    "{client_id:string,oracle_script_id:u64,calldata:bytes,ask_count:u64,min_count:u64}"
)
RATES = PyObi("[u64]")
PAIRS = PyObi("[string]")
CALLDATA = PyObi("{symbols:[string],multiplier:u64}")

# const
MULTIPLIER = 1000000000
ORACLE_SCRIPT_IDS = [8, 8, 8, 8, 9, 8]


class IStdReference(InterfaceScore):
    @interface
    def get_reference_data(self, _base: str, _quote: str) -> dict:
        pass

    @interface
    def get_reference_data_bulk(self, _encoded_pairs: bytes) -> list:
        pass


class StdReferenceProxy(IconScoreBase):

    REFERENCE_DATA = PyObi("{rate:u64,last_update_base:u64,last_update_quote:u64}")
    PAIRS = PyObi("[{base:string,quote:string}]")

    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        self.ref = VarDB("ref", db, value_type=Address)

    def on_install(self, _ref: Address) -> None:
        super().on_install()
        self.ref.set(_ref)

    def on_update(self) -> None:
        super().on_update()

    @external
    def set_ref(self, _ref: Address) -> None:
        if self.msg.sender != self.owner:
            self.revert("NOT_AUTHORIZED")

        self.ref.set(_ref)

    @external(readonly=True)
    def get_ref(self) -> Address:
        return self.ref.get()

    @external(readonly=True)
    def get_reference_data(self, _base: str, _quote: str) -> dict:
        ref = self.create_interface_score(self.ref.get(), IStdReference)
        return ref.get_reference_data(_base, _quote)

    @external(readonly=True)
    def get_reference_data_bulk(self, _encoded_pairs: bytes) -> list:
        ref = self.create_interface_score(self.ref.get(), IStdReference)
        return ref.get_reference_data_bulk(_encoded_pairs)
