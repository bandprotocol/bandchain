// npm install --save @bandprotocol/bandchain.js
const BandChain = require('./lib/lib/BandChain')
const { Obi } = require('@bandprotocol/obi.js')

// Chain parameters
const endpoint = 'https://poa-api.bandchain.org'

const symList = [
  'BTC',
  'ETH',
  'DAI',
  'REP',
  'ZRX',
  'BAT',
  'KNC',
  'LINK',
  'COMP',
  'BAND',
]

// Request parameters
const oracleScriptId = 8
const params = {
  symbols: symList,
  multiplier: 1000000,
}
const validatorCounts = {
  minCount: 3,
  askCount: 4,
}
;(async () => {
  const endpoint = 'https://poa-api.bandchain.org'

  const bandchain = new BandChain(endpoint)
  const price = await bandchain.getReferenceData(['ETH/USD'])
  console.log(price)
})()
