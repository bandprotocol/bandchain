# Band Protocol ICON Developer Documentation

In addition to data native to the [ICON blockchain](https://www.icondev.io/docs/what-is-icon), ICON developers also have access to various cryptocurrency price data provided by [Band Protocol](https://bandprotocol.com/)'s oracle.

## Standard Reference Dataset Contract Info

### Data Available

The price data originates from [data requests](https://github.com/bandprotocol/bandchain/wiki/System-Overview#oracle-data-request) made on BandChain. Band's [std_reference](https://bicon.tracker.solidwallet.io/contract/cxc79c2120992356eac409d5cf5ff650f2780a6995) SCORE on ICON then retrieves and stores the results of those requests. Specifically, the following price pairs are available to be read from the [std_reference_proxy](https://bicon.tracker.solidwallet.io/contract/cx61a36e5d10412e03c907a507d1e8c6c3856d9964) contract:

- BTC/USD
- ETH/USD
- ICX/USD

These prices are automatically updated every 5 minutes. The [std_reference_proxy](https://bicon.tracker.solidwallet.io/contract/cx61a36e5d10412e03c907a507d1e8c6c3856d9964) itself is currently deployed on ICON Yeouido testnet at [`cx61a36e5d10412e03c907a507d1e8c6c3856d9964`](https://bicon.tracker.solidwallet.io/contract/cx61a36e5d10412e03c907a507d1e8c6c3856d9964#readcontract).

The prices themselves are the median of the values retrieved by BandChain's validators from many sources which are [CoinGecko](https://www.coingecko.com/api/documentations/v3), [CryptoCompare](https://min-api.cryptocompare.com/), [Binance](https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md), [CoinMarketcap](https://coinmarketcap.com/), [HuobiPro](https://www.huobi.vc/en-us/exchange/), [CoinBasePro](https://pro.coinbase.com/), [Kraken](https://www.kraken.com/), [Bitfinex](https://www.bitfinex.com/), [Bittrex](https://global.bittrex.com/), [BITSTAMP](https://www.bitstamp.net/), [OKEX](https://www.okex.com/), [FTX](https://ftx.com/), [HitBTC](https://hitbtc.com/), [ItBit](https://www.itbit.com/), [Bithumb](https://www.bithumb.com/), [CoinOne](https://coinone.co.kr/). The data request is then made by executing Band's [Crypto Price in USD oracle script](https://docs.bandchain.org/built-in-oracle-scripts/crypto-price-1), the code of which you can view on [poa-mainnet](https://guanyu-poa.cosmoscan.io/oracle-script/8). Along with the price data, developers will also have access to the latest timestamp the price was updated.

These parameters are intended to act as security parameters to help anyone using the data to verify that the data they are using is what they expect and, perhaps more importantly, actually valid.

### Standard Reference Dataset Contract Price Update Process

For the ease of development, the Band Foundation will be maintaining and updating the [std_reference](https://bicon.tracker.solidwallet.io/contract/cxc79c2120992356eac409d5cf5ff650f2780a6995) contract with the latest price data. In the near future, we will be releasing guides on how developers can create similar contracts themselves to retrieve data from Band's oracle.

## Retrieving the Price Data

The code below shows an example of a relatively [simple price database](https://bicon.tracker.solidwallet.io/contract/cxd8b7e45dad9a111254f1d3168931e9ca562bc534) SCORE on ICON which retrieve price data from Band's [std_reference_proxy](https://bicon.tracker.solidwallet.io/contract/cx61a36e5d10412e03c907a507d1e8c6c3856d9964) contract and store it in the contract's state.

```shell=
===================
|                 |
| simple price db |
|                 |
===================
 |               ^
 |(1)            |(4)
 |Asking for     |Return
 |price data     |result
 |               |
 v               |
===================   (2) Ask     ===================
|                 |-------------->|                 |
|  std ref proxy  |               |     std ref     |
|                 |<--------------|                 |
===================   (3) Result  ===================

```

The contract is able to store exchange rate of any price pair that available on the [std_reference](https://bicon.tracker.solidwallet.io/contract/cxc79c2120992356eac409d5cf5ff650f2780a6995) contract. For more information on what oracle scripts are and how data requests work on BandChain in general, please see their [wiki](https://github.com/bandprotocol/bandchain/wiki/System-Overview#oracle-data-request) and [developer documentation](https://docs.bandchain.org/dapp-developers/requesting-data-from-bandchain)

```python
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

```

### Code Breakdown

The example code above can be broken down into two sections: defining the interface for the `IStdReferenceProxy` and the actual `SimplePriceDB` SCORE code itself.

#### IStdReferenceProxy Interface

This section consists of two functions, `get_reference_data` and `get_reference_data_bulk`. This is the interface that we'll use to query price from Band oracle for the latest price of a token or a bunch of tokens.

#### SimplePriceDB class

The `SimplePriceDB` then contains the main logic of our SCORE. It's purpose will be to store the latest price of tokens.

The price data query itself then occurs in `get_price`. Before we can call the method, however, we need to first set the address of the `std_reference_proxy`. This is done by calling the `set_proxy` method or `contructor`. After that the price should be set by calling `set_single` or `set_multiple`.

The `set_single` function will simply calling `get_reference_data` from `std_reference_proxy` with base symbol and quote symbol. It then extract the exchange rate from the result and save to the state.

The `set_multiple` function converts the input into an array of base symbol and quote symbol arrays. After that it will call `get_reference_data_bulk` from `std_reference_proxy` with base symbols and quote symbols. It then extract the exchange rates from the results and save all of them to the state.

The full source code for the `SimplePriceDB` score can be found [in this repo](https://github.com/bandprotocol/band-icon-integration-example/tree/master/simple_price_db) along with the JSON for sending the `set_proxy`, `set_single`, `set_multiple`. The score itself is also deployed to the testnet at address [cxd8b7e45dad9a111254f1d3168931e9ca562bc534](https://bicon.tracker.solidwallet.io/contract/cxd8b7e45dad9a111254f1d3168931e9ca562bc534#readcontract).
