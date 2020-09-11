from iconservice import *
from .pyobi import *

TAG = "Aggregator"

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
SYMBOLS = [
    [
        "BTC",
        "ETH",
        "USDT",
        "XRP",
        "LINK",
        "DOT",
        "BCH",
        "LTC",
        "ADA",
        "BSV",
        "CRO",
        "BNB",
        "EOS",
        "XTZ",
        "TRX",
        "XLM",
        "ATOM",
        "XMR",
        "OKB",
        "USDC",
        "NEO",
        "XEM",
        "LEO",
        "HT",
        "VET",
    ],
    [
        "YFI",
        "MIOTA",
        "LEND",
        "SNX",
        "DASH",
        "COMP",
        "ZEC",
        "ETC",
        "OMG",
        "MKR",
        "ONT",
        "NXM",
        "AMPL",
        "BAT",
        "THETA",
        "DAI",
        "REN",
        "ZRX",
        "ALGO",
        "FTT",
        "DOGE",
        "KSM",
        "WAVES",
        "EWT",
        "DGB",
    ],
    [
        "KNC",
        "ICX",
        "TUSD",
        "SUSHI",
        "BTT",
        "BAND",
        "ERD",
        "ANT",
        "NMR",
        "PAX",
        "LSK",
        "LRC",
        "HBAR",
        "BAL",
        "RUNE",
        "YFII",
        "LUNA",
        "DCR",
        "SC",
        "STX",
        "ENJ",
        "BUSD",
        "OCEAN",
        "RSR",
        "SXP",
    ],
    [
        "BTG",
        "BZRX",
        "SRM",
        "SNT",
        "SOL",
        "CKB",
        "BNT",
        "CRV",
        "MANA",
        "YFV",
        "KAVA",
        "MATIC",
        "TRB",
        "REP",
        "FTM",
        "TOMO",
        "ONE",
        "WNXM",
        "PAXG",
        "WAN",
        "SUSD",
        "RLC",
        "OXT",
        "RVN",
        "FNX",
    ],
    [
        "EUR",
        "GBP",
        "CNY",
        "SGD",
        "RMB",
        "KRW",
        "JPY",
        "INR",
        "RUB",
        "CHF",
        "AUD",
        "BRL",
        "CAD",
        "HKD",
        "XAU",
        "XAG",
    ],
    [
        "RENBTC",
        "WBTC",
        "DIA",
        "BTM",
        "IOTX",
        "FET",
        "JST",
        "MCO",
        "KMD",
        "BTS",
        "QKC",
        "YAMV2",
        "XZC",
        "UOS",
        "AKRO",
        "HNT",
        "HOT",
        "KAI",
        "OGN",
        "WRX",
        "KDA",
        "ORN",
        "FOR",
        "AST",
        "STORJ",
    ],
]


class IBridgeCache(InterfaceScore):
    @interface
    def get_latest_response(self, encoded_request: bytes) -> dict:
        pass


class Aggregator(IconScoreBase):
    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        self.bridge_address = VarDB("bridge_address", db, value_type=Address)

    def on_install(self, bridge_address: Address) -> None:
        super().on_install()
        self.bridge_address.set(bridge_address)

    def on_update(self) -> None:
        super().on_update()

    # A pure function that receive a symbol and then return
    # index of symbol's group and index of symbol inside its group.
    # @param sybbol, a string that represent an asset.
    def symbol_to_indexes(self, symbol: str) -> (int, int):
        for i, symbols in enumerate(SYMBOLS):
            if symbol in symbols:
                return (i, symbols.index(symbol))
        self.revert("UNKNOWN_SYMBOL")

    # Get bridge address from storage.
    @external(readonly=True)
    def get_bridge_address(self) -> Address:
        return self.bridge_address.get()

    # This function receives a symbol and asks the bridge
    # for the USD exchange rate of that symbol multiplied by 1e9.
    # @param sybbol, a string that represent an asset.
    @external(readonly=True)
    def get_rate_from_symbol(self, symbol: str) -> int:
        if symbol == "USD":
            return MULTIPLIER
        else:
            bridge = self.create_interface_score(self.bridge_address.get(), IBridgeCache)
            req = {"client_id": "bandteam", "ask_count": 4, "min_count": 3}

            outer_idx, inner_idx = self.symbol_to_indexes(symbol)

            req["oracle_script_id"] = ORACLE_SCRIPT_IDS[outer_idx]
            req["calldata"] = CALLDATA.encode(
                {"symbols": SYMBOLS[outer_idx], "multiplier": MULTIPLIER}
            )

            res = bridge.get_latest_response(REQUEST_PACKET.encode(req))
            return RATES.decode(res["result"])[inner_idx]

    # This function receives an encoded array of pairs
    # and then return an array of relative prices.
    # @param encoded_pairs, an obi encode of array of pairs. For example ["BTC/DAI"].
    @external(readonly=True)
    def get_reference_data(self, encoded_pairs: bytes) -> list:
        pairs = PAIRS.decode(encoded_pairs)
        result = []

        for pair in pairs:
            [base, quote] = pair.split("/")
            [base_price, quote_price] = [
                self.get_rate_from_symbol(base),
                self.get_rate_from_symbol(quote),
            ]

            result.append((base_price * MULTIPLIER * MULTIPLIER) // quote_price)

        return result

