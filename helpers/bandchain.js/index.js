const BandChain = require('./lib/lib/BandChain')

// BandChain devnet endpoint URL
const endpoint = 'http://poa-api.bandchain.org'

;(async () => {
  // Instantiating BandChain with REST endpoint
  const bandchain = new BandChain(endpoint)
  const result = await bandchain.getReferenceData(['BAND/USD'])
  console.log(result)
})()
