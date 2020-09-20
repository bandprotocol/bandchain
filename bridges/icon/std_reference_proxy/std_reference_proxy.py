from iconservice import *


TAG = "StdReferenceProxy"


class IStdReference(InterfaceScore):
    @interface
    def get_reference_data(self, _base: str, _quote: str) -> dict:
        pass

    @interface
    def get_reference_data_bulk(self, _json_pairs: str) -> list:
        pass


class StdReferenceProxy(IconScoreBase):
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
    def get_reference_data_bulk(self, _json_pairs: str) -> list:
        ref = self.create_interface_score(self.ref.get(), IStdReference)
        return ref.get_reference_data_bulk(_json_pairs)

