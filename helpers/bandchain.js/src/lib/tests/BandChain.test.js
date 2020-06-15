import BandChain from './../Bandchain'
import { JestEnvironment } from '@jest/environment'

jest.setTimeout(30000)

const chainID = 'band-guanyu-devnet-2'
const endpoint = 'http://guanyu-devnet.bandchain.org/rest'
const mnemonic =
  'final little loud vicious door hope differ lucky alpha morning clog oval milk repair off course indicate stumble remove nest position journey throw crane'

let testRequestID = 1

it('Test BandChain constructor', () => {
  let bandchain = new BandChain(chainID, endpoint)
  expect(bandchain.chainID).toBe(chainID)
  expect(bandchain.endpoint).toBe(endpoint)
})

it('Test BandChain getOracleScript success', async () => {
  const oracleScriptID = 1
  let bandchain = new BandChain(chainID, endpoint)
  let oracleScript = await bandchain.getOracleScript(oracleScriptID)
  expect(JSON.stringify(oracleScript)).toBe(
    JSON.stringify({
      owner: 'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs',
      name: 'Cryptocurrency Price in USD',
      description:
        'Oracle script that queries the average cryptocurrency price using current price data from CoinGecko, CryptoCompare, and Binance',
      filename:
        '52923c6702521f09f08bb4d27f2640b7340bfb5071bee2a354b17915b2e81fe8',
      schema: '{symbol:string,multiplier:u64}/{px:u64}',
      source_code_url:
        'https://ipfs.io/ipfs/QmdMKT62HYaaYH44DrW1UkQNhsd76nZXej6KXWjYtR9c5m',
      id: oracleScriptID,
    }),
  )
})

it('Test BandChain getOracleScript error', () => {
  let oracleScriptID = 30
  let bandchain = new BandChain(chainID, endpoint)
  return expect(bandchain.getOracleScript(oracleScriptID)).rejects.toThrow(
    'No oracle script found with the given ID',
  )
})

it('Test BandChain submitRequestTx', async () => {
  let oracleScriptID = 1
  let bandchain = new BandChain(chainID, endpoint)
  let oracleScript = await bandchain.getOracleScript(oracleScriptID)
  let requestID = await bandchain.submitRequestTx(
    oracleScript,
    { symbol: 'BTC', multiplier: BigInt('1000000000') },
    { minCount: 2, askCount: 4 },
    mnemonic,
  )
  testRequestID = requestID
  return expect(requestID).toBeDefined()
})

it('Test BandChain getRequestProof', async () => {
  let bandchain = new BandChain(chainID, endpoint)
  let requestProof = await bandchain.getRequestProof(testRequestID)
  return expect(requestProof).toBeDefined()
})

it('Test BandChain getRequestResult', async () => {
  let bandchain = new BandChain(chainID, endpoint)
  let requestID = 1
  let requestResult = await bandchain.getRequestResult(requestID)
  return (
    expect(requestResult).toBeDefined() &&
    expect(requestResult.RequestPacketData).toEqual({
      client_id: 'bandchain.js',
      oracle_script_id: '1',
      calldata: 'AAAAA0JUQwAAAAA7msoA',
      ask_count: '4',
      min_count: '2',
    }) &&
    expect(requestResult.ResponsePacketData.client_id).toEqual(
      'bandchain.js',
    ) &&
    expect(requestResult.ResponsePacketData.request_id).toEqual(testRequestID)
  )
})

it('Test BandChain getLastMatchingRequestResult', async () => {
  let oracleScriptID = 1
  let bandchain = new BandChain(chainID, endpoint)
  let oracleScript = await bandchain.getOracleScript(oracleScriptID)
  let lastRequestResult = await bandchain.getLastMatchingRequestResult(
    oracleScript,
    { symbol: 'BTC', multiplier: BigInt('1000000000') },
    { minCount: 2, askCount: 4 },
  )
  return (
    expect(lastRequestResult).toBeDefined() &&
    expect(lastRequestResult.client_id).toEqual('bandchain.js') &&
    expect(lastRequestResult.resolve_status).toEqual(1)
  )
})
