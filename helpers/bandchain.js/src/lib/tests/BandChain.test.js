import BandChain from './../Bandchain';
import { JestEnvironment } from '@jest/environment';

const chainIDMaster = 'bandchain';
const endpointMaster = 'http://d3n.bandprotocol.com/rest';
const mnemonic =
  'nose define hover rain beef renew cruise remain mistake page brain mimic wool sponsor mimic reunion inspire raw marine brain choice later crash clerk';

jest.setTimeout(10000);

it('Test BandChain constructor', () => {
  let bandchain = new BandChain(chainIDMaster, endpointMaster);
  expect(bandchain.chainID).toBe('bandchain');
  expect(bandchain.endpoint).toBe('http://d3n.bandprotocol.com/rest');
});

it('Test BandChain getOracleScript success', async () => {
  const oracleScriptID = 1;
  let bandchain = new BandChain(chainIDMaster, endpointMaster);
  let oracleScript = await bandchain.getOracleScript(oracleScriptID);
  expect(JSON.stringify(oracleScript)).toBe(
    JSON.stringify({
      owner: 'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs',
      name: 'Crypto price script',
      description:
        'Oracle script for getting the current an average cryptocurrency price from various sources.',
      filename: 'a602627e52801fa8d19e2be1a57b093322e195710311defd043b2910cb741e80',
      schema: '{symbol:string,multiplier:u64}/{px:u64}',
      source_code_url: 'https://ipfs.io/ipfs/QmUbdfoRR9ge6P39EoqDjBhQoDeaT6gJu76Ce9kKsz938N',
      id: oracleScriptID,
    })
  );
});

it('Test BandChain getOracleScript error', () => {
  let oracleScriptID = 20;
  let bandchain = new BandChain(chainIDMaster, endpointMaster);
  return expect(bandchain.getOracleScript(oracleScriptID)).rejects.toThrow(
    'No oracle script found with the given ID'
  );
});

it('Test BandChain submitRequestTx', async () => {
  let oracleScriptID = 1;
  let bandchain = new BandChain(chainIDMaster, endpointMaster);
  let oracleScript = await bandchain.getOracleScript(oracleScriptID);
  let requestID = await bandchain.submitRequestTx(
    oracleScript,
    { symbol: 'BTC', multiplier: BigInt('1000000000') },
    { minCount: 2, askCount: 4 },
    mnemonic
  );
});

it('Test BandChain getRequestProof', async () => {
  let oracleScriptID = 1;
  let bandchain = new BandChain(chainIDMaster, endpointMaster);
  let oracleScript = await bandchain.getOracleScript(oracleScriptID);
  let requestID = await bandchain.submitRequestTx(
    oracleScript,
    { symbol: 'BTC', multiplier: BigInt('1000000000') },
    { minCount: 2, askCount: 4 },
    mnemonic
  );
  let requestProof = await bandchain.getRequestProof(requestID);
  console.log(requestProof);
});
