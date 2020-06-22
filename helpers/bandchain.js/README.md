<div align="center">
  <h2>BandChain.js</h2>
  <blockquote>Library for interacting with BandChain in browser and Node.js environments</blockquote>
</div>

## ⭐️ Features

- Requesting Data from BandChain
- Reading Data Request Status from BandChain

## 📦 Installation

```bash
npm install --save @bandprotocol/bandchain.js
```

### Use in Node.js

```js
const BandChain = require('@bandprotocol/bandchain.js')

const chainId = 'band-guanyu-devnet-2'
const endpoint = 'http://guanyu-devnet.bandchain.org/rest'

const bandchain = new BandChain(chainId, endpoint)
...
```

## 🔍 Test
```
npm run test
```

## 💎 Example Usages

```js
const BandChain = require('@bandprotocol/bandchain.js')

const chainId = 'band-guanyu-devnet-2'
const endpoint = 'http://guanyu-devnet.bandchain.org/rest'

// Instantiating BandChain with REST endpoint
const bandchain = new BandChain(chainId, endpoint)

// Create an instance of OracleScript with the script ID
const oracleScript = await bandchain.getOracleScript(1)

// Get script info
const schema = oracleScript.schema
const description = oracleScript.description

// Read latest script result
const result = await bandchain.getLatestRequestResult(oracleScript, {
  symbol: 'BTC',
  multiplier: 10000,
})

// Create a new request, which will block into the tx is confirmed
try {
  const minCount = 5
  const askCount = 7
  const mnemonic =
    'dumb spot lyrics car infant round rate famous inhale tennis text current'
  const requestId = await bandchain.submitRequestTx(
    oracleScript,
    { symbol: 'BTC', multiplier: 10000 },
    { minCount, askCount },
    mnemonic,
  )

  // Get request proof
  const requestProof = await bandchain.getRequestProof(requestID)
  // Get final result (blocking until the reports & aggregations are finished)
  const finalResult = await bandchain.getRequestResult(requestID)
  // Get the result of the most recent request that match the specified parameters
  const minCount = 2
  const askCount = 4
  const inputParameters = { symbol: 'BTC', multiplier: BigInt('1000000000') }
  const lastMatchResult = await bandchain.getLastMatchingRequestResult(
    oracleScript,
    inputParameters,
    { minCount, askCount },
  )
} catch {
  // Something went wrong
  console.error('Data request failed')
}
```
