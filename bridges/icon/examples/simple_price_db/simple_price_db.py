from iconservice import *

TAG = "SimplePriceDB"

# Interface of StdReferenceProxy
class IStdReferenceProxy(InterfaceScore):
    @interface
    def get_reference_data(self, _base: str, _quote: str) -> dict:
        pass

    @interface
    def get_reference_data_bulk(self, _bases: str, _quotes: str) -> list:
        pass


# SimplePriceDB contract
class SimplePriceDB(IconScoreBase):
    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        self.std_reference_proxy_address = VarDB(
            "std_reference_proxy_address", db, value_type=Address
        )
        self.prices = DictDB("prices", db, value_type=int)

    def on_install(self, _proxy: Address) -> None:
        super().on_install()
        self.set_proxy(_proxy)

    def on_update(self) -> None:
        super().on_update()

    # Get price of the given pair multiply by 1e18.
    # For example "BTC/ETH" -> 30.09 * 1e18.
    @external(readonly=True)
    def get_price(self, _pair: str) -> int:
        return self.prices[_pair]

    # Sets the proxy contract address, which can only be set by the owner.
    @external
    def set_proxy(self, _proxy: Address) -> None:
        if self.msg.sender != self.owner:
            self.revert("NOT_AUTHORIZED")

        self.std_reference_proxy_address.set(_proxy)

    # This function accepts the string coin pairs such as "BTC/ETH".
    # And then pass it to the std_reference_proxy. After receiving the output
    # from the std_reference_proxy, the exchange rate of the input pair
    # is recorded into the state.
    @external
    def set_single(self, _pair: str) -> None:
        proxy = self.create_interface_score(
            self.std_reference_proxy_address.get(), IStdReferenceProxy
        )
        base, quote = _pair.split("/")
        result = proxy.get_reference_data(base, quote)

        self.prices[_pair] = result["rate"]

    # This function accepts the string of the encoding of an array of coin pairs
    # such as '["BTC/ETH", "ETH/USD", "USDT/USD"]'. And then transform the format
    # of the input to pass it to the std_reference_proxy. After receiving the output
    # from the std_reference_proxy, each pair's rate is recorded into the state.
    @external
    def set_multiple(self, _json_pairs: str) -> None:
        proxy = self.create_interface_score(
            self.std_reference_proxy_address.get(), IStdReferenceProxy
        )

        pairs = json_loads(_json_pairs)
        bases, quotes = [json_dumps(x) for x in zip(*[pair.split("/") for pair in pairs])]
        results = proxy.get_reference_data_bulk(bases, quotes)

        if len(pairs) != len(results):
            self.revert("LEN_PAIRS_MUST_EQUAL_TO_LEN_RESULTS")

        for pair, result in zip(pairs, results):
            self.prices[pair] = result["rate"]

