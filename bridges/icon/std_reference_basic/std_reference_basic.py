from iconservice import *

TAG = "StdReferenceBasic"


class StdReferenceBasic(IconScoreBase):

    ONE = 1_000_000_000
    MICRO_SEC = 1_000_000

    @eventlog
    def RefDataUpdate(self, _symbol: str, _rate: int, _last_update: int, _request_id: int):
        pass

    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        # Mapping from symbol to json encoded ref data.
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

    def _get_reference_data(self, _base: str, _quote: str) -> dict:
        ref_base = self._get_ref_data(_base)
        ref_quote = self._get_ref_data(_quote)
        return {
            "rate": (ref_base["rate"] * self.ONE * self.ONE) // ref_quote["rate"],
            "last_update_base": ref_base["last_update"],
            "last_update_quote": ref_quote["last_update"],
        }

    @external(readonly=True)
    def get_ref_data(self, _symbol: str) -> dict:
        return self._get_ref_data(_symbol)

    @external(readonly=True)
    def get_reference_data(self, _base: str, _quote: str) -> dict:
        return self._get_reference_data(_base, _quote)

    @external(readonly=True)
    def get_reference_data_bulk(self, _bases: str, _quotes: str) -> list:
        bases = json_loads(_bases)
        quotes = json_loads(_quotes)
        if len(bases) != len(quotes):
            self.revert("BAD_INPUT_LENGTH")
        return [self._get_reference_data(*pair) for pair in zip(bases, quotes)]

    def _set_refs(self, _symbol: str, _rate: int, _last_update: int, _request_id: int) -> None:
        self.refs[_symbol] = json_dumps(
            {"rate": _rate, "last_update": _last_update * self.MICRO_SEC, "request_id": _request_id}
        )

    @external
    def relay(self, _symbols: str, _rates: str, _resolve_times: str, _request_ids: str) -> None:
        if self.msg.sender != self.owner:
            self.revert("NOT_AUTHORIZED")

        symbols = json_loads(_symbols)
        rates = json_loads(_rates)
        resolve_times = json_loads(_resolve_times)
        request_ids = json_loads(_request_ids)

        len_symbols = len(symbols)
        if len_symbols != len(rates):
            self.revert("BAD_RATES_LENGTH")
        if len_symbols != len(resolve_times):
            self.revert("BAD_RESOLVE_TIMES_LENGTH")
        if len_symbols != len(request_ids):
            self.revert("BAD_REQUEST_IDS_LENGTH")

        for param in zip(symbols, rates, resolve_times, request_ids):
            # set rate and last_update for the symbol
            self._set_refs(*param)

            # emit event log
            self.RefDataUpdate(*param)

