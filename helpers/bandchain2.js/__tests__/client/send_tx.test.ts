// import { jsxText } from '@babel/types'
import axios from 'axios'
jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

import { Client } from '../../src'

const TEST_RPC = 'https://api-mock.bandprotocol.com/rest'
const TEST_MSG = {
  msg: [
    {
      type: 'oracle/Request',
      value: {
        oracle_script_id: '1',
        calldata: 'AAAAA0JUQwAAAAAAAAAB',
        ask_count: '4',
        min_count: '3',
        client_id: 'from_pyband',
        sender: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
      },
    },
  ],
  fee: { gas: '1000000', amount: [{ denom: 'uband', amount: '0' }] },
  memo: 'Memo',
  signatures: [
    {
      signature:
        'hQpMSSaOVbT5vd3yladNX9RNA9vSq4ts4cPufdoesjUtPje5i73f048MM0xPnAB7JWSRuUSsZD5M6L6WGk3Qkw==',
      pub_key: {
        type: 'tendermint/PubKeySecp256k1',
        value: 'A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8',
      },
      account_number: '36',
      sequence: '1092',
    },
  ],
}
const TEST_WRONG_SEQUENCE_MSG = {
  msg: [
    {
      type: 'oracle/Request',
      value: {
        oracle_script_id: '1',
        calldata: 'AAAAA0JUQwAAAAAAAAAB',
        ask_count: '4',
        min_count: '3',
        client_id: 'from_pyband',
        sender: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
      },
    },
  ],
  fee: { gas: '1000000', amount: [{ denom: 'uband', amount: '0' }] },
  memo: 'Memo',
  signatures: [
    {
      signature:
        'hQpMSSaOVbT5vd3yladNX9RNA9vSq4ts4cPufdoesjUtPje5i73f048MM0xPnAB7JWSRuUSsZD5M6L6WGk3Qkw==',
      pub_key: {
        type: 'tendermint/PubKeySecp256k1',
        value: 'A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8',
      },
      account_number: '0',
      sequence: '1092',
    },
  ],
}

const mockTxHash =
  'E204AAD58ACA8F00942B1BB66D9F745F5E2C21E04C5DF6A0CB73DF02B6B51121'

let client = new Client(TEST_RPC)
// let client = new Client(MASTER)

describe('sendTxSyncMode', () => {
  it('success', async () => {
    mockedAxios.post.mockResolvedValue({
      data: {
        height: '0',
        txhash: mockTxHash,
        raw_log: '[]',
      },
    })

    const res = await client.sendTxSyncMode(TEST_MSG)

    expect(res).toEqual({
      txHash: Buffer.from(
        'E204AAD58ACA8F00942B1BB66D9F745F5E2C21E04C5DF6A0CB73DF02B6B51121',
        'hex',
      ),
      code: 0,
      error_log: undefined,
    })
  })

  it('wrong sequence', async () => {
    mockedAxios.post.mockResolvedValue({
      data: {
        height: '0',
        txhash: Buffer.from(
          '611F45A21BB7937E451CDE78D124218603644635CC40A97D2BC1E854CED8D6E61',
          'hex',
        ),
        code: 4,
        raw_log:
          'unauthorized: signature verification failed; verify correct account sequence and chain-id',
      },
    })

    const res = await client.sendTxSyncMode(TEST_WRONG_SEQUENCE_MSG)

    expect(res).toEqual({
      txHash: Buffer.from(
        '611F45A21BB7937E451CDE78D124218603644635CC40A97D2BC1E854CED8D6E61',
        'hex',
      ),
      code: 4,
      errorLog:
        'unauthorized: signature verification failed; verify correct account sequence and chain-id',
    })
  })
})

describe('sendTxAsyncMode', () => {
  it('success', async () => {
    mockedAxios.post.mockResolvedValue({
      data: {
        height: '0',
        txhash: mockTxHash,
        raw_log: '[]',
      },
    })

    const res = await client.sendTxAsyncMode(TEST_MSG)

    expect(res).toEqual({
      txHash: Buffer.from(
        'E204AAD58ACA8F00942B1BB66D9F745F5E2C21E04C5DF6A0CB73DF02B6B51121',
        'hex',
      ),
    })
  })
})

describe('sendTxBlockMode', () => {
  it('success', async () => {
    mockedAxios.post.mockResolvedValue({
      data: {
        height: '715786',
        txhash: Buffer.from(
          '2DE264D16164BCCF695E960553FED537EDC00D0E3EDF69D6BFE4168C476AD03C',
          'hex',
        ),
        raw_log:
          '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"51"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"4"},{"key":"min_count","value":"3"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr"},{"key":"validator","value":"bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},{"key":"validator","value":"bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]}]',
        logs: [
          {
            msg_index: 0,
            log: '',
            events: [
              {
                type: 'message',
                attributes: [{ key: 'action', value: 'request' }],
              },
              {
                type: 'raw_request',
                attributes: [
                  { key: 'data_source_id', value: '1' },
                  {
                    key: 'data_source_hash',
                    value:
                      'c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0',
                  },
                  { key: 'external_id', value: '1' },
                  { key: 'calldata', value: 'BTC' },
                  { key: 'data_source_id', value: '2' },
                  {
                    key: 'data_source_hash',
                    value:
                      'dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479',
                  },
                  { key: 'external_id', value: '2' },
                  { key: 'calldata', value: 'BTC' },
                  { key: 'data_source_id', value: '3' },
                  {
                    key: 'data_source_hash',
                    value:
                      'f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921',
                  },
                  { key: 'external_id', value: '3' },
                  { key: 'calldata', value: 'BTC' },
                ],
              },
              {
                type: 'request',
                attributes: [
                  { key: 'id', value: '51' },
                  { key: 'client_id', value: 'from_pyband' },
                  { key: 'oracle_script_id', value: '1' },
                  {
                    key: 'calldata',
                    value: '000000034254430000000000000001',
                  },
                  { key: 'ask_count', value: '4' },
                  { key: 'min_count', value: '3' },
                  { key: 'gas_used', value: '2405' },
                  {
                    key: 'validator',
                    value: 'bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr',
                  },
                  {
                    key: 'validator',
                    value: 'bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh',
                  },
                  {
                    key: 'validator',
                    value: 'bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst',
                  },
                  {
                    key: 'validator',
                    value: 'bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec',
                  },
                ],
              },
            ],
          },
        ],
        gas_wanted: '1000000',
        gas_used: '343054',
      },
    })

    const res = await client.sendTxBlockMode(TEST_MSG)

    expect(res).toEqual({
      height: 715786,
      txHash: Buffer.from(
        '2DE264D16164BCCF695E960553FED537EDC00D0E3EDF69D6BFE4168C476AD03C',
        'hex',
      ),
      gasWanted: 1000000,
      gasUsed: 1000000,
      log: [
        {
          msg_index: 0,
          log: '',
          events: [
            {
              type: 'message',
              attributes: [{ key: 'action', value: 'request' }],
            },
            {
              type: 'raw_request',
              attributes: [
                { key: 'data_source_id', value: '1' },
                {
                  key: 'data_source_hash',
                  value:
                    'c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0',
                },
                { key: 'external_id', value: '1' },
                { key: 'calldata', value: 'BTC' },
                { key: 'data_source_id', value: '2' },
                {
                  key: 'data_source_hash',
                  value:
                    'dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479',
                },
                { key: 'external_id', value: '2' },
                { key: 'calldata', value: 'BTC' },
                { key: 'data_source_id', value: '3' },
                {
                  key: 'data_source_hash',
                  value:
                    'f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921',
                },
                { key: 'external_id', value: '3' },
                { key: 'calldata', value: 'BTC' },
              ],
            },
            {
              type: 'request',
              attributes: [
                { key: 'id', value: '51' },
                { key: 'client_id', value: 'from_pyband' },
                { key: 'oracle_script_id', value: '1' },
                {
                  key: 'calldata',
                  value: '000000034254430000000000000001',
                },
                { key: 'ask_count', value: '4' },
                { key: 'min_count', value: '3' },
                { key: 'gas_used', value: '2405' },
                {
                  key: 'validator',
                  value: 'bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr',
                },
                {
                  key: 'validator',
                  value: 'bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh',
                },
                {
                  key: 'validator',
                  value: 'bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst',
                },
                {
                  key: 'validator',
                  value: 'bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec',
                },
              ],
            },
          ],
        },
      ],
      errorLog: undefined,
      code: 0,
    })
  })

  it('wrong sequence', async () => {
    mockedAxios.post.mockResolvedValue({
      data: {
        height: '0',
        txhash: Buffer.from(
          '7F1CFFD674CAAEEB25922E9C6E9F8F8CEF7A325E25B64E8DAB070D7409FD1F72',
          'hex',
        ),
        codespace: 'sdk',
        code: 4,
        raw_log:
          'unauthorized: signature verification failed; verify correct account sequence and chain-id',
        gas_wanted: '1000000',
        gas_used: '27402',
      },
    })

    const res = await client.sendTxBlockMode(TEST_WRONG_SEQUENCE_MSG)

    expect(res).toEqual({
      height: 0,
      txHash: Buffer.from(
        '7F1CFFD674CAAEEB25922E9C6E9F8F8CEF7A325E25B64E8DAB070D7409FD1F72',
        'hex',
      ),
      gasWanted: 1000000,
      gasUsed: 1000000,
      code: 4,
      log: [],
      errorLog:
        'unauthorized: signature verification failed; verify correct account sequence and chain-id',
    })
  })
})
