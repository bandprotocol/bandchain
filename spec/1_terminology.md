# Terminology

## Data Sources

A data source is the most fundamental unit of the oracle system. It describes the procedure to retrieve a raw data point from a primary source and the fee associated with one data query. In D3N, a data source can be registered into the system by anyone. Once registered, a data source can be _owned_ or _unowned_. An owned data source can be changed and upgraded by its owner who can also claims the collected fees, while an unowned data source is immutable and cannot be changed.

Note that even though an unowned data source cannot be changed, it can still be controlled by centralized parties if its procedure depends on centralized sources. Below is some examples of data source procedures written in [Bash](<https://en.wikipedia.org/wiki/Bash_(Unix_shell)>).

**Example:** A script to retrieve a cryptocurrency price from [CoinGecko](https://www.coingecko.com/) with one argument: the currency identifier. The script assumes that [cURL](https://en.wikipedia.org/wiki/CURL) and [jq](https://github.com/stedolan/jq) are available on the host and the host is connected to the Internet.

```bash
#!/bin/sh

# Cryptocurrency price endpoint: https://www.coingecko.com/api/documentations/v3
URL="https://api.coingecko.com/api/v3/simple/price?ids=$1&vs_currencies=usd"
KEY=".$1.usd"

# Performs data fetching and parses the result
curl -s -X GET $URL -H "accept: application/json" | jq -r ".[\"$1\"].usd"
```

**Example:** A script to resolve a given hostname into IP addresses. The script assumes that [getent](https://en.wikipedia.org/wiki/Getent) and [awk](https://en.wikipedia.org/wiki/AWK) are available on the host and the host is connected to the DNS network.

```bash
#!/bin/sh

getent hosts $1 | awk '{ print $1 }'
```

## Oracle Scripts

An oracle script is an executable program that encodes (1) the set of raw data requests to the data sources it needs and (2) the way to aggregate raw data reports into the final result. An oracle script may depend on input data, making it reusable without needing to rewrite the whole program. It may also depend on other oracle scripts, making oracle scripts composable just like [smart contracts](https://en.wikipedia.org/wiki/Smart_contract).

**Example**: A [psudocode](https://en.wikipedia.org/wiki/Pseudocode) showing an example oracle script to fetch cryptocurrency prices from multiple sources: [CoinGecko](https://coinmarketcap.com/), [CryptoCompare](https://www.cryptocompare.com/), and [CoinMarketCap](https://coinmarketcap.com/). The code assumes that data sources for the three sources are available and reports the average among all sources from all data reporters as the final result.

```python
# 1st Phase. Emits raw data requests that the oracle script needs.
def prepare(symbol):
    request(get_px_from_coin_gecko, symbol)
    request(get_px_from_crypto_compare, symbol)
    request(get_px_from_coin_market_cap, symbol)

# 2nd Phase. Aggregates raw data reports into the final result.
def aggregate(symbol, number_of_reporters):
    data_report_count = 0
    price_sum = 0.0
    for reporter_index in range(number_of_reporters):
        for data_source in (
            get_px_from_coin_gecko,
            get_px_from_crypto_compare,
            get_px_from_coin_market_cap,
        ):
            price_sum = receive(reporter_index, data_source, symbol)
            data_report_count += 1
    return price_sum / data_report_count
```

### Data Requests

A data request is a transaction from a user to perform a data query based on an oracle script. A data request transaction specifies the oracle script to execute, the parameters to the script, and [other security parameters][sec:msg-request-data].

### Raw Data Requests

Raw data requests are requests to primary sources emitted while an oracle script is being executed in the first phase. It essentially consists of a data source's procedure and the associated parameters. Raw data requests are expected to be resolved by D3N block validators in the form of raw data reports.

### Raw Data Reports

Raw data reports are the results from resolving raw data requests by D3N block validators. The raw reports are submitted to D3N. Once a sufficient number of reports are collected, they will be used in the oracle script's second phase to compute the final result of the data request.

### Oracle Data Proof

Once the aggregation is complete, the final result of the data request is stored permanently in D3N's global state. As similar to most other blockchains, the whole state of D3N can be represented as a [Merkle root hash](https://en.wikipedia.org/wiki/Merkle_tree). An oracle data proof is a Merkle proof that shows the existence of the final result of the data request with other information related to it, including the oracle script hash, the parameters, the time of execution, etc.
