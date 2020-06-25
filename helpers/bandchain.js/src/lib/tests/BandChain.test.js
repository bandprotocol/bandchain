const BandChain = require('../Bandchain')

jest.setTimeout(30000)

const endpoint = 'http://guanyu-devnet.bandchain.org/rest'
const mnemonic =
  'final little loud vicious door hope differ lucky alpha morning clog oval milk repair off course indicate stumble remove nest position journey throw crane'

let testRequestID = 1

it('Test BandChain constructor', () => {
  let bandchain = new BandChain(endpoint)
  expect(bandchain.endpoint).toBe(endpoint)
})

it('Test BandChain getOracleScript success', async () => {
  const oracleScriptID = 1
  let bandchain = new BandChain(endpoint)
  let oracleScript = await bandchain.getOracleScript(oracleScriptID)
  expect(JSON.stringify(oracleScript)).toBe(
    JSON.stringify({
      owner: 'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs',
      name: 'Cryptocurrency Price in USD',
      description:
        'Oracle script that queries the average cryptocurrency price using current price data from CoinGecko, CryptoCompare, and Binance',
      filename:
        'ff48f3b9876cddb41e2371fe0fc5cc516619944adec75e91dfda85ace561dd9c',
      schema: '{symbol:string,multiplier:u64}/{px:u64}',
      source_code_url:
        'https://ipfs.io/ipfs/QmY3S4dYuWMX4L7RMioEbUcLLZxc3tRoDMJxVQMthd7Amy',
      id: oracleScriptID,
    }),
  )
})

it('Test BandChain getOracleScript error', () => {
  let oracleScriptID = 1e18
  let bandchain = new BandChain(endpoint)
  expect(bandchain.getOracleScript(oracleScriptID)).rejects.toThrow(
    'No oracle script found with the given ID',
  )
})

it('Test BandChain submitRequestTx', async () => {
  let oracleScriptID = 1
  let bandchain = new BandChain(endpoint)
  let oracleScript = await bandchain.getOracleScript(oracleScriptID)
  let requestID = await bandchain.submitRequestTx(
    oracleScript,
    { symbol: 'BAND', multiplier: BigInt('1000000') },
    { minCount: 2, askCount: 4 },
    mnemonic,
  )
  testRequestID = requestID
  expect(requestID).toBeDefined()
})

it('Test BandChain getRequestID error', async () => {
  let bandchain = new BandChain(endpoint)
  expect(
    bandchain.getRequestID(
      '13DEADCF273FCE723B809DDD6F29E5D0B5FD397256FD872D602676094061F20D', // Not a request tx
    ),
  ).rejects.toThrow('Not a request tx')
})

it('Test BandChain getRequestEVMProof', async () => {
  let bandchain = new BandChain(endpoint)
  let requestProof = await bandchain.getRequestEVMProof(testRequestID)
  expect(requestProof).toBeDefined()
})

it('Test BandChain getRequestNonEVMProof', async () => {
  let bandchain = new BandChain(endpoint)
  let requestProof = await bandchain.getRequestNonEVMProof(testRequestID)
  expect(requestProof).toBeDefined()
})

it('Test BandChain getRequestResult', async () => {
  let bandchain = new BandChain(endpoint)
  let requestID = 1
  let requestResult = await bandchain.getRequestResult(requestID)

  expect(requestResult).toBeDefined()
})

it('Test BandChain getLastMatchingRequestResult', async () => {
  let oracleScriptID = 1
  let bandchain = new BandChain(endpoint)
  let oracleScript = await bandchain.getOracleScript(oracleScriptID)
  let lastRequestResult = await bandchain.getLastMatchingRequestResult(
    oracleScript,
    { symbol: 'BAND', multiplier: BigInt('1000000') },
    { minCount: 2, askCount: 4 },
  )

  expect(lastRequestResult).toBeDefined()
  expect(lastRequestResult.result).toBeDefined()
  expect(lastRequestResult.result.px).toBeGreaterThan(0)
})
