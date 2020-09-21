from iconservice import *

TAG = "StdReferenceBasic"


class StdReferenceBasic(IconScoreBase):

    ONE = 1_000_000_000
    MICRO_SEC = 1_000_000

    @eventlog(indexed=1)
    def RefDataUpdate(self, _symbol: str, _rate: int, _last_update: int, _request_id: int):
        pass

    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        # Mapping from symbol to obi encoded ref data.
        self.refs = DictDB("refs", db, value_type=str)

    def on_install(self) -> None:
        super().on_install()

    def on_update(self) -> None:
        super().on_update()

    def _get_ref_data(self, _symbol: str) -> dict:
        if _symbol == "USD":
            return {"rate": self.ONE, "last_update": self.block.timestamp, "request_id": 0}

        if self.refs[_symbol] == "":
            self.revert("REF_DATA_NOT_AVAILABLE")

        return json_loads(self.refs[_symbol])

    def _get_reference_data(self, _pair: str) -> dict:
        [base, quote] = _pair.split("/")
        ref_base = self._get_ref_data(base)
        ref_quote = self._get_ref_data(quote)
        return {
            "rate": (ref_base["rate"] * self.ONE * self.ONE) // ref_quote["rate"],
            "last_update_base": ref_base["last_update"],
            "last_update_quote": ref_quote["last_update"],
        }

    @external(readonly=True)
    def get_ref_data(self, _symbol: str) -> dict:
        return self._get_ref_data(_symbol)

    @external(readonly=True)
    def get_reference_data(self, _pair: str) -> dict:
        return self._get_reference_data(_pair)

    @external(readonly=True)
    def get_reference_data_bulk(self, _json_pairs: str) -> list:
        return [self._get_reference_data(pair) for pair in json_loads(_json_pairs)]

    def _set_refs(self, _symbol: str, _rate: int, _last_update: int, _request_id: int) -> None:
        self.refs[_symbol] = json_dumps(
            {"rate": _rate, "last_update": _last_update * self.MICRO_SEC, "request_id": _request_id}
        )

    @external
    def relay(self, _json_data_list: str) -> None:
        if self.msg.sender != self.owner:
            self.revert("NOT_AUTHORIZED")

        for data in json_loads(_json_data_list):
            # set rate and last_update for the symbol
            self._set_refs(data["symbol"], data["rate"], data["resolve_time"], data["request_id"])

            # emit event log
            self.RefDataUpdate(
                data["symbol"], data["rate"], data["resolve_time"], data["request_id"]
            )

