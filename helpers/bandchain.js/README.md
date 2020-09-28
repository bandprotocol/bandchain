# BandChain.js

Library for interacting with BandChain in browser and Node.js environments

## ⭐️ Features

- Making data requests to BandChain
- Getting the results of the latest request on BandChain that mathces specific parameters

## 📦 Installation

### NPM

```bash
npm install --save @bandprotocol/bandchain.js
```

### Yarn

```bash
yarn add @bandprotocol/bandchain.js
```

## Usage

### Making Data Requests to BandChain

The library allow users to interact with data requests on BandChain. This can be done in two main ways:

- Making a new data request to an oracle script
- Querying for the latest request that match certain request parameters

#### Making a data requests

Data requests can be made to BandChain using the `submitRequestTx` function. This function takes 4 arguments:

- `oracleScript`: The oracle script object, which can be retrieved using the `getOracleScript` function.
- `parameters`: The object containing the request parameters
- `validatorCounts`: The set of `askCount` and `minCount` validator counts
- `mnemonic`: The mnemonic to be used to make the request

##### Example

```js
const BandChain = require('@bandprotocol/bandchain.js');
const { Obi } = require('@bandprotocol/obi.js');

const endpoint = 'http://guanyu-devnet.bandchain.org/rest';
const mnemonic = "borrow scare approve embrace uncle witness debate genuine cage snap toast utility situate labor curtain shock olive myself uphold tongue license critic glare taxi";

const symbol = "BTC"
const multiplier = 1000

const sendRequest = async (mnemonic, endpoint, minCount, askCount) => {
  // Instantiating BandChain with REST endpoint
  const bandchain = new BandChain(endpoint);

  // Create an instance of OracleScript with the script ID
  const oracleScript = await bandchain.getOracleScript(1);

  // Create a new request, which will block until the tx is confirmed
  try {
    const requestId = await bandchain.submitRequestTx(
      oracleScript,
      {
        symbol,
        multiplier,
      },
      { minCount, askCount },
      mnemonic
    );

    // Get final result (blocking until the reports & aggregations are finished)
    const finalResult = await bandchain.getRequestResult(requestId);
    let res = new Obi(oracleScript.schema).decodeOutput(
      Buffer.from(finalResult.response_packet_data.result, 'base64')
    );
    if (res) {
      console.log(res);
    }
  } catch {
    throw 'Data request failed';
  }
};

sendRequest(mnemonic,endpoint,3,4);
```

#### Querying previous request results

Alternatively, you can also query BandChain for the latest successful data request that matches certain request parameters using the `getLastMatchingRequestResult` function. Similar to `submitRequestTx`, this function also takes in the `oracleScript` object, the input `parameters`, and `validatorCounts`.

Note that the mnemonic is not necessary in this case as we are simply querying the chain's REST endpoint instead of making a request ourselves. The code below demonstrates an example usage.

```js
const BandChain = require('@bandprotocol/bandchain.js');
const endpoint = 'http://guanyu-devnet.bandchain.org/rest';

const symbol = 'BTC';
const multiplier = 1000;

const sendRequest = async (endpoint, minCount, askCount) => {
  // Instantiating BandChain with REST endpoint
  const bandchain = new BandChain(endpoint);

  // Create an instance of OracleScript with the script ID
  const oracleScript = await bandchain.getOracleScript(1);

  // Query BandChain for the latest request results that match the specified parameters
  try {
    const latestResult = await bandchain.getLastMatchingRequestResult(
      oracleScript,
      {
        symbol,
        multiplier,
      },
      { minCount, askCount }
    );

    // Get final result (blocking until the reports & aggregations are finished)
    console.log(latestResult.result);
  } catch {
    throw 'Data request failed';
  }
};

getLatestRequestResult(endpoint, 3, 4);
```
