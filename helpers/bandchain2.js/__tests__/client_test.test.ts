import axios from 'axios'
import { Client } from '../src/index'
import { Address } from '../src/wallet'

jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

const TEST_RPC = 'https://api-mock.bandprotocol.com/rest'

const client = new Client(TEST_RPC)

describe('Client get data', () => {
  it('get chain ID', () => {
    const resp = { data: { chain_id: 'bandchain' } }
    mockedAxios.get.mockResolvedValue(resp)

    const response = client.getChainID()
    response.then((e) => expect(e).toEqual('bandchain'))
  })

  it('get data source by ID', () => {
    const resp = {
      data: {
        height: '651093',
        result: {
          owner: 'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs',
          name: 'CoinGecko Cryptocurrency Price',
          description:
            'Retrieves current price of a cryptocurrency from https://www.coingecko.com',
          filename:
            'c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0',
        },
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const expected = {
      owner: Address.fromAccBech32(
        'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs',
      ),
      name: 'CoinGecko Cryptocurrency Price',
      description:
        'Retrieves current price of a cryptocurrency from https://www.coingecko.com',
      fileName:
        'c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0',
    }
    const response = client.getDataSource(1)
    response.then((e) => expect(e).toEqual(expected))
  })

  it('get oracle script by ID', () => {
    const resp = {
      data: {
        height: '2985172',
        result: {
          owner: 'band17f6n25na5kume99j4qdfzlf7jnpu9u2neqqvt8',
          name: 'OS 03',
          description: 'TBD',
          filename:
            '2bf80fa07dc9585305818939853833f140fdb7e7bab824a628dc2ebc2094f482',
          schema:
            '{base_symbol:string,quote_symbol:string,multiplier:u64}/{px:u64}',
          source_code_url:
            'https://ipfs.io/ipfs/QmcXZKevdv2QCAkzKF69YzSK6w7KziEugaVyyjuLF1RM6u',
        },
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const expected = {
      owner: Address.fromAccBech32(
        'band17f6n25na5kume99j4qdfzlf7jnpu9u2neqqvt8',
      ),
      name: 'OS 03',
      description: 'TBD',
      fileName:
        '2bf80fa07dc9585305818939853833f140fdb7e7bab824a628dc2ebc2094f482',
      schema:
        '{base_symbol:string,quote_symbol:string,multiplier:u64}/{px:u64}',
      sourceCodeUrl:
        'https://ipfs.io/ipfs/QmcXZKevdv2QCAkzKF69YzSK6w7KziEugaVyyjuLF1RM6u',
    }
    const response = client.getOracleScript(3)
    response.then((e) => expect(e).toEqual(expected))
  })

  it('get latest request', () => {
    const resp = {
      data: {
        height: '3006088',
        result: {
          request: {
            oracle_script_id: '1',
            calldata: 'AAAAA0xSQwAAAANVU0QAAAAAAA9CQA==',
            requested_validators: [
              'bandvaloper1ywd2m858gu4eya3nzx6f9vme3sn82dr4thjnme',
              'bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre',
              'bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw',
              'bandvaloper1alzj765pzuhtjkmslme4fdpeakc0036xnyjltn',
            ],
            min_count: '3',
            request_height: '552201',
            request_time: '2020-08-23T19:35:39.841842928Z',
            raw_requests: [
              {
                external_id: '1',
                data_source_id: '12',
                calldata:
                  'UkVOQlRDIFdCVEMgRElBIEJUTSBJT1RYIEZFVCBKU1QgTUNPIEtNRCBCVFMgUUtDIFlBTVYyIFhaQyBVT1MgQUtSTyBITlQgSE9UIEtBSSBPR04gV1JYIEtEQSBPUk4gRk9SIEFTVCBTVE9SSg==',
              },
            ],
          },
          reports: [
            {
              validator: 'bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw',
              raw_reports: [
                {
                  external_id: '2',
                  data:
                    'MTc4NjEuMTk5MiwxODA1Ni45Nzk2LDEuMjYzMjY3LDAuMDYwMzczMTIsMC4wMDc3MzY3OSwwL',
                },
              ],
            },
          ],
          result: {
            request_packet_data: {
              oracle_script_id: '1',
              calldata: 'AAAAA0xSQwAAAANVU0QAAAAAAA9CQA==',
              ask_count: '4',
              min_count: '3',
            },
            response_packet_data: {
              request_id: '44893',
              ans_count: '3',
              request_time: '1598211339',
              resolve_time: '1598211345',
              resolve_status: 1,
              result: 'AAAAAAAC+mQ=',
            },
          },
        },
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const expected = {
      request: {
        oracleScriptID: 1,
        requestedValidators: [
          'bandvaloper1ywd2m858gu4eya3nzx6f9vme3sn82dr4thjnme',
          'bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre',
          'bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw',
          'bandvaloper1alzj765pzuhtjkmslme4fdpeakc0036xnyjltn',
        ],
        minCount: 3,
        requestHeight: 552201,
        clientID: '',
        calldata: Buffer.from('AAAAA0xSQwAAAANVU0QAAAAAAA9CQA==', 'base64'),
        rawRequests: [
          {
            externalID: 1,
            dataSourceID: 12,
            calldata: Buffer.from(
              'UkVOQlRDIFdCVEMgRElBIEJUTSBJT1RYIEZFVCBKU1QgTUNPIEtNRCBCVFMgUUtDIFlBTVYyIFhaQyBVT1MgQUtSTyBITlQgSE9UIEtBSSBPR04gV1JYIEtEQSBPUk4gRk9SIEFTVCBTVE9SSg==',
              'base64',
            ),
          },
        ],
      },
      reports: [
        {
          validator: 'bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw',
          inBeforeResolve: false,
          rawReports: [
            {
              externalID: 2,
              data: Buffer.from(
                'MTc4NjEuMTk5MiwxODA1Ni45Nzk2LDEuMjYzMjY3LDAuMDYwMzczMTIsMC4wMDc3MzY3OSwwL',
                'base64',
              ),
            },
          ],
        },
      ],
      result: {
        requestPacketData: {
          clientID: '',
          askCount: 4,
          minCount: 3,
          oracleScriptID: 1,
          calldata: Buffer.from('AAAAA0xSQwAAAANVU0QAAAAAAA9CQA==', 'base64'),
        },
        responsePacketData: {
          requestID: 44893,
          requestTime: 1598211339,
          resolveTime: 1598211345,
          resolveStatus: 1,
          ansCount: 3,
          clientID: '',
          result: Buffer.from('AAAAAAAC+mQ=', 'base64'),
        },
      },
    }
    const response = client.getLatestRequest(
      8,
      Buffer.from('AAAAA0xSQwAAAANVU0QAAAAAAA9CQA', 'base64'),
      3,
      4,
    )
    response.then((e) => expect(e).toEqual(expected))
  })

  it('get request by ID', () => {
    const resp = {
      data: {
        height: '3006088',
        result: {
          request: {
            oracle_script_id: '1',
            calldata: 'AAAAA0xSQwAAAANVU0QAAAAAAA9CQA==',
            requested_validators: [
              'bandvaloper1ywd2m858gu4eya3nzx6f9vme3sn82dr4thjnme',
              'bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre',
              'bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw',
              'bandvaloper1alzj765pzuhtjkmslme4fdpeakc0036xnyjltn',
            ],
            min_count: '3',
            request_height: '552201',
            request_time: '2020-08-23T19:35:39.841842928Z',
            raw_requests: [
              {
                external_id: '1',
                data_source_id: '12',
                calldata:
                  'UkVOQlRDIFdCVEMgRElBIEJUTSBJT1RYIEZFVCBKU1QgTUNPIEtNRCBCVFMgUUtDIFlBTVYyIFhaQyBVT1MgQUtSTyBITlQgSE9UIEtBSSBPR04gV1JYIEtEQSBPUk4gRk9SIEFTVCBTVE9SSg==',
              },
            ],
          },
          reports: [
            {
              validator: 'bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw',
              raw_reports: [
                {
                  external_id: '2',
                  data:
                    'MTc4NjEuMTk5MiwxODA1Ni45Nzk2LDEuMjYzMjY3LDAuMDYwMzczMTIsMC4wMDc3MzY3OSwwL',
                },
              ],
            },
          ],
          result: {
            request_packet_data: {
              oracle_script_id: '1',
              calldata: 'AAAAA0xSQwAAAANVU0QAAAAAAA9CQA==',
              ask_count: '4',
              min_count: '3',
            },
            response_packet_data: {
              request_id: '44893',
              ans_count: '3',
              request_time: '1598211339',
              resolve_time: '1598211345',
              resolve_status: 1,
              result: 'AAAAAAAC+mQ=',
            },
          },
        },
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const expected = {
      request: {
        oracleScriptID: 1,
        requestedValidators: [
          'bandvaloper1ywd2m858gu4eya3nzx6f9vme3sn82dr4thjnme',
          'bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre',
          'bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw',
          'bandvaloper1alzj765pzuhtjkmslme4fdpeakc0036xnyjltn',
        ],
        minCount: 3,
        requestHeight: 552201,
        clientID: '',
        calldata: Buffer.from('AAAAA0xSQwAAAANVU0QAAAAAAA9CQA==', 'base64'),
        rawRequests: [
          {
            externalID: 1,
            dataSourceID: 12,
            calldata: Buffer.from(
              'UkVOQlRDIFdCVEMgRElBIEJUTSBJT1RYIEZFVCBKU1QgTUNPIEtNRCBCVFMgUUtDIFlBTVYyIFhaQyBVT1MgQUtSTyBITlQgSE9UIEtBSSBPR04gV1JYIEtEQSBPUk4gRk9SIEFTVCBTVE9SSg==',
              'base64',
            ),
          },
        ],
      },
      reports: [
        {
          validator: 'bandvaloper1yplk6n4wmeaarxp966gukpxupg3jqfcqkh32mw',
          inBeforeResolve: false,
          rawReports: [
            {
              externalID: 2,
              data: Buffer.from(
                'MTc4NjEuMTk5MiwxODA1Ni45Nzk2LDEuMjYzMjY3LDAuMDYwMzczMTIsMC4wMDc3MzY3OSwwL',
                'base64',
              ),
            },
          ],
        },
      ],
      result: {
        requestPacketData: {
          clientID: '',
          askCount: 4,
          minCount: 3,
          oracleScriptID: 1,
          calldata: Buffer.from('AAAAA0xSQwAAAANVU0QAAAAAAA9CQA==', 'base64'),
        },
        responsePacketData: {
          requestID: 44893,
          requestTime: 1598211339,
          resolveTime: 1598211345,
          resolveStatus: 1,
          ansCount: 3,
          clientID: '',
          result: Buffer.from('AAAAAAAC+mQ=', 'base64'),
        },
      },
    }

    const response = client.getRequestByID(3)
    response.then((e) => expect(e).toEqual(expected))
  })
})

describe('Client get request ID by transaction hash', () => {
  it('test with request tx', () => {
    const resp = {
      data: {
        height: '3739',
        txhash:
          '90ED70061C1A24B1141F81BADEDAB19570D0B9B255412BF5D680A9BC93539115',
        raw_log:
          '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"band178yydxyzplrh5v0jegvtrahymqragf6spy8kdd"},{"key":"sender","value":"band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"},{"key":"amount","value":"100000uband"}]}]},{"msg_index":1,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"4"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"2"},{"key":"min_count","value":"2"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]}]',
        logs: [
          {
            msg_index: 0,
            log: '',
            events: [
              {
                type: 'message',
                attributes: [
                  { key: 'action', value: 'send' },
                  {
                    key: 'sender',
                    value: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
                  },
                  { key: 'module', value: 'bank' },
                ],
              },
              {
                type: 'transfer',
                attributes: [
                  {
                    key: 'recipient',
                    value: 'band178yydxyzplrh5v0jegvtrahymqragf6spy8kdd',
                  },
                  {
                    key: 'sender',
                    value: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
                  },
                  { key: 'amount', value: '100000uband' },
                ],
              },
            ],
          },
          {
            msg_index: 1,
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
                  { key: 'id', value: '4' },
                  { key: 'client_id', value: 'from_pyband' },
                  { key: 'oracle_script_id', value: '1' },
                  { key: 'calldata', value: '000000034254430000000000000001' },
                  { key: 'ask_count', value: '2' },
                  { key: 'min_count', value: '2' },
                  { key: 'gas_used', value: '2405' },
                  {
                    key: 'validator',
                    value: 'bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh',
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
        gas_used: '298476',
        tx: {
          type: 'cosmos-sdk/StdTx',
          value: {
            msg: [
              {
                type: 'cosmos-sdk/MsgSend',
                value: {
                  from_address: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
                  to_address: 'band178yydxyzplrh5v0jegvtrahymqragf6spy8kdd',
                  amount: [{ denom: 'uband', amount: '100000' }],
                },
              },
              {
                type: 'oracle/Request',
                value: {
                  oracle_script_id: '1',
                  calldata: 'AAAAA0JUQwAAAAAAAAAB',
                  ask_count: '2',
                  min_count: '2',
                  client_id: 'from_pyband',
                  sender: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
                },
              },
            ],
            fee: { amount: [{ denom: 'uband', amount: '0' }], gas: '1000000' },
            signatures: [
              {
                pub_key: {
                  type: 'tendermint/PubKeySecp256k1',
                  value: 'A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8',
                },
                signature:
                  'K6PmN0HluRRm7zliKJho9F2OYpokB5JrYlqAC+KmHQMwKBpRYNaZzYTGzoeBol2mm3sfcdUo8rzwrQngzP8s+g==',
              },
            ],
            memo: 'TEST',
          },
        },
        timestamp: '2020-11-09T09:29:49Z',
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const expected = [4]
    const response = client.getRequestIDByTxHash(
      Buffer.from(
        '90ED70061C1A24B1141F81BADEDAB19570D0B9B255412BF5D680A9BC93539115',
        'hex',
      ),
    )
    response.then((e) => expect(e).toEqual(expected))
  })

  it('test with multi request msgs tx', () => {
    const resp = {
      data: {
        height: '279',
        txhash:
          '0838E29162B87D9D41E2BAE49C97272970453F2CFA708FBA6B8BE71F314BEB5B',
        raw_log:
          '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"1"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"2"},{"key":"min_count","value":"2"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]},{"msg_index":1,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"2"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"2"},{"key":"min_count","value":"2"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]},{"msg_index":2,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"request"}]},{"type":"raw_request","attributes":[{"key":"data_source_id","value":"1"},{"key":"data_source_hash","value":"c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0"},{"key":"external_id","value":"1"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"2"},{"key":"data_source_hash","value":"dd155f719c5201336d4852830a3ad446ddf01b1ab647cf6ea5d7b9e7678a7479"},{"key":"external_id","value":"2"},{"key":"calldata","value":"BTC"},{"key":"data_source_id","value":"3"},{"key":"data_source_hash","value":"f3bad1a6d88cd30ce311d6845f114422f9c2c52c32c45b5086d69d052ad90921"},{"key":"external_id","value":"3"},{"key":"calldata","value":"BTC"}]},{"type":"request","attributes":[{"key":"id","value":"3"},{"key":"client_id","value":"from_pyband"},{"key":"oracle_script_id","value":"1"},{"key":"calldata","value":"000000034254430000000000000001"},{"key":"ask_count","value":"2"},{"key":"min_count","value":"2"},{"key":"gas_used","value":"2405"},{"key":"validator","value":"bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh"},{"key":"validator","value":"bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"}]}]}]',
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
                  { key: 'id', value: '1' },
                  { key: 'client_id', value: 'from_pyband' },
                  { key: 'oracle_script_id', value: '1' },
                  { key: 'calldata', value: '000000034254430000000000000001' },
                  { key: 'ask_count', value: '2' },
                  { key: 'min_count', value: '2' },
                  { key: 'gas_used', value: '2405' },
                  {
                    key: 'validator',
                    value: 'bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr',
                  },
                  {
                    key: 'validator',
                    value: 'bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec',
                  },
                ],
              },
            ],
          },
          {
            msg_index: 1,
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
                  { key: 'id', value: '2' },
                  { key: 'client_id', value: 'from_pyband' },
                  { key: 'oracle_script_id', value: '1' },
                  { key: 'calldata', value: '000000034254430000000000000001' },
                  { key: 'ask_count', value: '2' },
                  { key: 'min_count', value: '2' },
                  { key: 'gas_used', value: '2405' },
                  {
                    key: 'validator',
                    value: 'bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh',
                  },
                  {
                    key: 'validator',
                    value: 'bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec',
                  },
                ],
              },
            ],
          },
          {
            msg_index: 2,
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
                  { key: 'id', value: '3' },
                  { key: 'client_id', value: 'from_pyband' },
                  { key: 'oracle_script_id', value: '1' },
                  { key: 'calldata', value: '000000034254430000000000000001' },
                  { key: 'ask_count', value: '2' },
                  { key: 'min_count', value: '2' },
                  { key: 'gas_used', value: '2405' },
                  {
                    key: 'validator',
                    value: 'bandvaloper1ma0cxd4wpcqk3kz7fr8x609rqmgqgvrpem0txh',
                  },
                  {
                    key: 'validator',
                    value: 'bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst',
                  },
                ],
              },
            ],
          },
        ],
        gas_wanted: '1000000',
        gas_used: '782736',
        tx: {
          type: 'cosmos-sdk/StdTx',
          value: {
            msg: [
              {
                type: 'oracle/Request',
                value: {
                  oracle_script_id: '1',
                  calldata: 'AAAAA0JUQwAAAAAAAAAB',
                  ask_count: '2',
                  min_count: '2',
                  client_id: 'from_pyband',
                  sender: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
                },
              },
              {
                type: 'oracle/Request',
                value: {
                  oracle_script_id: '1',
                  calldata: 'AAAAA0JUQwAAAAAAAAAB',
                  ask_count: '2',
                  min_count: '2',
                  client_id: 'from_pyband',
                  sender: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
                },
              },
              {
                type: 'oracle/Request',
                value: {
                  oracle_script_id: '1',
                  calldata: 'AAAAA0JUQwAAAAAAAAAB',
                  ask_count: '2',
                  min_count: '2',
                  client_id: 'from_pyband',
                  sender: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
                },
              },
            ],
            fee: { amount: [{ denom: 'uband', amount: '0' }], gas: '1000000' },
            signatures: [
              {
                pub_key: {
                  type: 'tendermint/PubKeySecp256k1',
                  value: 'A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8',
                },
                signature:
                  'PaxPc1330hLWZXBUHtPAPQdP4tS2LvGAiwaT8isxr8UKdDx0uXjLblTBjWO+yHENTjGxMdb2cI3BUOCIIoFCzQ==',
              },
            ],
            memo: 'TEST',
          },
        },
        timestamp: '2020-11-09T07:44:35Z',
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const expected = [1, 2, 3]
    const response = client.getRequestIDByTxHash(
      Buffer.from(
        '0838E29162B87D9D41E2BAE49C97272970453F2CFA708FBA6B8BE71F314BEB5B',
        'hex',
      ),
    )
    response.then((e) => expect(e).toEqual(expected))
  })

  it('test with no request msgs tx', async () => {
    const resp = {
      data: {
        height: '3740',
        txhash:
          '9F83E4994C048F784D0E30F45696C0A1E5BA7407B2E1833B439FA172B3B75F00',
        raw_log:
          '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"report"}]},{"type":"report","attributes":[{"key":"id","value":"4"},{"key":"validator","value":"bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"}]}]}]',
        logs: [
          {
            msg_index: 0,
            log: '',
            events: [
              {
                type: 'message',
                attributes: [{ key: 'action', value: 'report' }],
              },
              {
                type: 'report',
                attributes: [
                  { key: 'id', value: '4' },
                  {
                    key: 'validator',
                    value: 'bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec',
                  },
                ],
              },
            ],
          },
        ],
        gas_wanted: '62945',
        gas_used: '52449',
        tx: {
          type: 'cosmos-sdk/StdTx',
          value: {
            msg: [
              {
                type: 'oracle/Report',
                value: {
                  request_id: '4',
                  raw_reports: [
                    { external_id: '1', data: 'MTU0MzMuMQo=' },
                    { external_id: '3', data: 'MTU0MTYuMDc1Cg==' },
                    { external_id: '2', data: 'MTU0MzcuMTYK' },
                  ],
                  validator:
                    'bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec',
                  reporter: 'band1ue0623hwqkvm5s5hq0jnqwh4ende28gvmlzvd2',
                },
              },
            ],
            fee: { amount: [], gas: '62945' },
            signatures: [
              {
                pub_key: {
                  type: 'tendermint/PubKeySecp256k1',
                  value: 'A4OH54K/spetlw9jNR8LjqiHKi85jWKyV8zZDZH6dBij',
                },
                signature:
                  'W5RuZvIsGxoVqTHCpPzWuD6pvaA6YedEaQ8TfS3d4AYuykREcltOPOigun8wXm5dhfFzAf/BbgP9vccQ8mmkHw==',
              },
            ],
            memo: 'yoda:/exec:lambda:1.1.2',
          },
        },
        timestamp: '2020-11-09T09:29:51Z',
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const response = client.getRequestIDByTxHash(
      Buffer.from(
        '9F83E4994C048F784D0E30F45696C0A1E5BA7407B2E1833B439FA172B3B75F00',
        'hex',
      ),
    )

    response.catch((err) =>
      expect(err).toEqual(new Error('There is no request message in this tx')),
    )
  })
})
