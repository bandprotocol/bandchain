<div align="center">
  <h2>PyBand</h2>
  <blockquote>BandChain Python Library</blockquote>
</div>

## ⭐️ Features

This helper library allows users to request the latest request result that match scertain input parameters. The parameters that can be specified are:

- The `oracleScriptID`
- the `askCount` and `minCount`
- the `calldata` (request parameters) associated with the request

For more information on each these, please refer to our [wiki](https://github.com/bandprotocol/bandchain/wiki/Protocol-Messages#parameters-4).

## 📦 Installation

The library is available on [PyPI](https://pypi.org/project/pyband/)

```bash
pip install pyband
```

## 💎 Example Usage

The example code below shows how the library can be used to get the result of the latest request for the price of Bitcoin. The specified parameters are

- `oracleScriptID`: 1
- `calldata`: The hex string representing [OBI](https://github.com/bandprotocol/bandchain/wiki/Oracle-Binary-Encoding-(OBI))-encoded value of `{symbol:BTC,multiplier:1000000}`
- `minCount`: 4
- `askCount`: 4

```python
from pyband import Client, PyObi


def main():
    c = Client("http://guanyu-devnet.bandchain.org/rest")
    req_info = c.get_latest_request(1, bytes.fromhex("0000000342544300000000000f4240"), 4, 4)
    oracle_script = c.get_oracle_script(1)
    obi = PyObi(oracle_script.schema)
    print(obi.decode_output(req_info.result.ResponsePacketData.result))


if __name__ == "__main__":
    main()
```
