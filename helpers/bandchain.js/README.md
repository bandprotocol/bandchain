 <div align="center">
 <!-- <img align="center" width="180" src="https://i.imgur.com/62VsVXD.png" /> -->
  <h2>BandChain.js</h2>
  <blockquote>Library for interacting with BandChain in browser and Node.js environments</blockquote>
  <!-- <a href="https://github.com/hodgef/js-library-boilerplate/actions"><img alt="Build Status" src="https://github.com/hodgef/js-library-boilerplate/workflows/Build/badge.svg?color=green" /></a> <a href="https://github.com/hodgef/js-library-boilerplate/actions"> <img alt="Publish Status" src="https://github.com/hodgef/js-library-boilerplate/workflows/Publish/badge.svg?color=green" /></a> <img src="https://img.shields.io/david/hodgef/js-library-boilerplate.svg" /> <a href="https://david-dm.org/hodgef/js-library-boilerplate?type=dev"><img src="https://img.shields.io/david/dev/hodgef/js-library-boilerplate.svg" /></a> <img src="https://api.dependabot.com/badges/status?host=github&repo=hodgef/js-library-boilerplate" /> -->

<!-- <strong>This is a more robust library boilerplate. For a minimal alternative, check out [js-library-boilerplate-basic](https://github.com/hodgef/js-library-boilerplate-basic).</strong> -->

</div>

## ⭐️ Features

- Requesting Data from BandChain
- Reading Data Request Status from BandChain

## 📦 Installation

```
npm install
```

### npm

```js
import BandChain from 'bandchain.js'

const chainID = 'band-guanyu-alchemist'
const endpoint = 'http://devnet.bandchain.org/rest'

const bandchain = new BandChain({ chainID, endpoint })
...
```

### self-host/cdn

```html
<link href="build/index.css" rel="stylesheet" />
<script src="build/index.js"></script>

let BandChain = window.BandChain.default const bandchain = new BandChain(/*
chainID & endpoint */) ...
```

## 💎 Example Usages

```js
import BandChain from "bandchain.js";

const chainID = "band-guanyu-alchemist";
const endpoint = "http://devnet.bandchain.org/rest";

// Instantiating BandChain with REST endpoint
const bandchain = new BandChain({ chainID, endpoint });

// Create an instance of OracleScript with the script ID
const oracleScript = await bandchain.getOracleScript(1);

// Get script info
const schema = oracleScript.schema;
const description = oracleScript.description;

// Read latest script result
const result = await oracleScript.getLatestRequestResult({
  symbol: "BTC",
  multiplier: 10000,
});

// Create a new request, which will block into the tx is confirmed
try {
  const validatorsRequired = 7;
  const request = await oracleScript.submitRequestTx(
    { symbol: "BTC", multiplier: 10000 },
    validatorsRequired
  );

  // Check report status
  const { count } = await request.getReportStatus();

  // Get final result (blocking until the reports & aggregations are finished)
  const finalResult = await request.getFinalRequestResult();
} catch {
  // Something went wrong
  console.error("Data request failed");
}
```
