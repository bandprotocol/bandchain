 <div align="center">
 <!-- <img align="center" width="180" src="https://i.imgur.com/62VsVXD.png" /> -->
  <h2>Band Data Requester</h2>
  <blockquote>A light-weight node.js tool to query data on BandChain periadically</blockquote>
</div>

## ⭐️ Features

- Support cron-style scheduling
- Support configuration via `json` file
- No-installation required with [`npx`](https://www.npmjs.com/package/npx)

## 📦 Prerequisite

You need to create a `config.json` file in your machine.

See example [`config.json`](./config.js) for requesting data from GuanYu [devnet](https://guanyu-devnet.cosmoscan.io/) every 5 minutes:

```json
{
  "chainId": "band-guanyu-devnet-2",
  "endpoint": "http://guanyu-devnet.bandchain.org/rest",
  "mnemonic": "final little loud vicious door hope differ lucky alpha morning clog oval milk repair off course indicate stumble remove nest position journey throw crane",
  "cronPattern": "*/5 * * * *",
  "validatorCounts": {
    "minCount": 3,
    "askCount": 4
  },
  "requests": [
    {
      "oracleScriptId": 1,
      "params": {
        "symbol": "BTC",
        "multiplier": 1000000
      }
    },
    {
      "oracleScriptId": 13,
      "params": {
        "base_symbol": "ETH",
        "quote_symbol": "CNY",
        "aggregation_method": "median",
        "multiplier": 1000000
      }
    }
  ]
}
```

## 💎 Example Usages

```bash
npx @bandprotocol/cron-requester config.json
```

If the `config.json` is correctly formatted, you should see something like this:

```
--------------------------------------------------------
⭐️ Cron is running! Your requests will be executed with cron pattern */5 * * * *
📆 Your first requests will start at Sun Jun 14 2020 19:25:00 GMT+0700
--------------------------------------------------------
⏰ Requests start at 6/14/2020, 7:25:00 PM
∟ ✅ requestId = 180 | oracleScript #1 {"symbol":"BTC","multiplier":1000000}
∟ ✅ requestId = 181 | oracleScript #13 {"base_symbol":"ETH","quote_symbol":"CNY","aggregation_method":"median","multiplier":1000000}
⛳️ [2/2] requests went through
--------------------------------------------------------
```
