 <div align="center">
 <!-- <img align="center" width="180" src="https://i.imgur.com/62VsVXD.png" /> -->
  <h2>BandChain.js</h2>
  <blockquote>Library for interacting with BandChain in browser and Node.js environments</blockquote>
  <!-- <a href="https://github.com/hodgef/js-library-boilerplate/actions"><img alt="Build Status" src="https://github.com/hodgef/js-library-boilerplate/workflows/Build/badge.svg?color=green" /></a> <a href="https://github.com/hodgef/js-library-boilerplate/actions"> <img alt="Publish Status" src="https://github.com/hodgef/js-library-boilerplate/workflows/Publish/badge.svg?color=green" /></a> <img src="https://img.shields.io/david/hodgef/js-library-boilerplate.svg" /> <a href="https://david-dm.org/hodgef/js-library-boilerplate?type=dev"><img src="https://img.shields.io/david/dev/hodgef/js-library-boilerplate.svg" /></a> <img src="https://api.dependabot.com/badges/status?host=github&repo=hodgef/js-library-boilerplate" /> -->

<strong>Lib boilerplate by [js-library-boilerplate](https://github.com/hodgef/js-library-boilerplate).</strong>

</div>

## ⭐️ Features

- Requesting Data from BandChain
- Reading Data Request Status from BandChain

## 📦 Installation

```bash
npm install
```

### npm

```js
import BandChain from 'bandchain.js'

const chainId = 'band-guanyu-devnet-2'
const endpoint = 'http://devnet.bandchain.org/rest'

const bandchain = new BandChain(chainId, endpoint)
...
```

### self-host/cdn

```html
<link href="build/index.css" rel="stylesheet" />
<script src="build/index.js"></script>

<script>
  let BandChain = window.BandChain.default;
  const bandchain = new BandChain(/*chainId & endpoint */);
  ...
</script>
```

## 💎 Example Usages

```js
import BandChain from "bandchain.js";

const chainId = "band-guanyu-alchemist";
const endpoint = "http://devnet.bandchain.org/rest";

// Instantiating BandChain with REST endpoint
const bandchain = new BandChain(chainId, endpoint);

// Create an instance of OracleScript with the script ID
const oracleScript = await bandchain.getOracleScript(1);

// Get script info
const schema = oracleScript.schema;
const description = oracleScript.description;

// Read latest script result
const result = await bandchain.getLatestRequestResult(oracleScript, {
  symbol: "BTC",
  multiplier: 10000,
});

// Create a new request, which will block into the tx is confirmed
try {
  const minCount = 5;
  const askCount = 7;
  const mnemonic =
    "dumb spot lyrics car infant round rate famous inhale tennis text current";
  const requestId = await bandchain.submitRequestTx(
    oracleScript,
    { symbol: "BTC", multiplier: 10000 },
    { minCount, askCount },
    mnemonic
  );

  // Get request proof
  const requestProof = await bandchain.getRequestProof(requestID);
  // Get final result (blocking until the reports & aggregations are finished)
  const finalResult = await bandchain.getRequestResult(requestID);
  // Get the result of the most recent request that match the specified parameters
  const minCount = 2;
  const askCount = 4;
  const inputParameters = { symbol: 'BTC', multiplier: BigInt('1000000000') };
  const lastMatchResult = await bandchain.getLastMatchingRequestResult(
    oracleScript,
    inputParameters,
    minCount,
    askCount
  )
} catch {
  // Something went wrong
  console.error("Data request failed");
}
```
