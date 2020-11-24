import axios from 'axios'
import { Client, Coin } from '../../src'
import { Address } from '../../src/wallet'

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
})

describe('Get oracle script by ID', () => {
  it('success', () => {
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
})

describe('get latest request', () => {
  it('success', () => {
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
})

describe('Get request by ID', () => {
  it('success', () => {
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

  it('success, without request', () => {
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
      result: undefined,
    }

    const response = client.getRequestByID(3)
    response.then((e) => expect(e).toEqual(expected))
  })
})

describe('Client get request id by transaction hash', () => {
  it('success, with request tx', () => {
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

  it('success, with multi request msgs tx', () => {
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

describe('Get reference data', () => {
  it('success', () => {
    const resp = {
      data: {
        height: '2953006',
        result: [
          {
            symbol: 'BTC',
            multiplier: '1000000000',
            px: '16242693800000',
            request_id: '1171969',
            resolve_time: '1605512243',
          },
          {
            symbol: 'ETH',
            multiplier: '1000000000',
            px: '454523400000',
            request_id: '1171969',
            resolve_time: '1605512943',
          },
          {
            symbol: 'TRX',
            multiplier: '1000000000',
            px: '25428330',
            request_id: '1171969',
            resolve_time: '1605512443',
          },
        ],
      },
    }
    mockedAxios.post.mockResolvedValue(resp)

    const response = client.getReferenceData(['BTC/USD', 'TRX/ETH'])
    const expected = [
      {
        pair: 'BTC/USD',
        rate: 16242.693800000001,
        updatedAt: { base: 1605512243, quote: Math.floor(Date.now() / 1000) },
      },
      {
        pair: 'TRX/ETH',
        rate: 0.0000559450404533628,
        updatedAt: { base: 1605512443, quote: 1605512943 },
      },
    ]
    response.then((e) => expect(e).toEqual(expected))
  })
})

describe('get latest block', () => {
  it('success', () => {
    const resp = {
      data: {
        block_id: {
          hash:
            '391E99908373F8590C928E0619956DA3D87EB654445DA4F25A185C9718561D53',
          parts: {
            total: '1',
            hash:
              '9DC1DB39B7DDB97DE353DFB2898198BAADEFB7DF8090BF22678714F769D69F42',
          },
        },
        block: {
          header: {
            version: { block: '10', app: '0' },
            chain_id: 'bandchain',
            height: '1032007',
            time: '2020-11-05T09:15:18.445494105Z',
            last_block_id: {
              hash:
                '4BC01E257662B5F9337D615D06D5D19D8D79F7BA5A4022A85B4DC4ED3C59F7CA',
              parts: {
                total: '1',
                hash:
                  '6471C0A51FB7043171EAA76CAFA900B36A4494F47F975980732529D247E8EBA5',
              },
            },
            last_commit_hash:
              '17B2CE4ABA910E85847537F1323DB95C9F16C20C60E9B9BBB04C633C3125BD92',
            data_hash:
              'EFE5E3F549554FEE8EB9B393740C250D74580427A96A175ABB105806039CFFE2',
            validators_hash:
              'E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151',
            next_validators_hash:
              'E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151',
            consensus_hash:
              '0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45',
            app_hash:
              '6E2B1ECE9D912D86C25182E8B7419583ABCE978BFC66DC2556BB0D06A8D528EF',
            last_results_hash: '',
            evidence_hash: '',
            proposer_address: 'BDB6A0728C8DFE2124536F16F2BA428FE767A8F9',
          },
          data: {
            txs: [
              'yAEoKBapCj5CcI40CAESDwAAAANCVEMAAAAAAAAAARgEIAMqC2Zyb21fcHliYW5kMhSQ78AMmxLrubEOPhhIwKK5oyk9oBIQCgoKBXViYW5kEgEwEMCEPRpqCibrWumHIQP+cIvaZlJP0sa86QaC44VVqFHgPSruT2KbBd6Q9R7ZvBJANbPqLRAgwwULWWwb5O2/eb6ddtDr65kRFgDcOTTGIqQS5Iz1NvHWHfkPKRoM8egErMsgE9YnuE+pAqoc/YvNfiIEVEVTVA==',
            ],
          },
          evidence: { evidence: null },
          last_commit: {
            height: '1032006',
            round: '0',
            block_id: {
              hash:
                '4BC01E257662B5F9337D615D06D5D19D8D79F7BA5A4022A85B4DC4ED3C59F7CA',
              parts: {
                total: '1',
                hash:
                  '6471C0A51FB7043171EAA76CAFA900B36A4494F47F975980732529D247E8EBA5',
              },
            },
            signatures: [
              {
                block_id_flag: 3,
                validator_address: '5179B0BB203248E03D2A1342896133B5C58E1E44',
                timestamp: '2020-11-05T09:15:18.53815896Z',
                signature:
                  'TZY24CKwZOE8wqfE0NM3qzkQ7qCpCrGEHNZdf8n31L4otZzbKGfOL05kGtBsGkTnZkVv7aJmrJ7XbvIzv0SREQ==',
              },
              {
                block_id_flag: 2,
                validator_address: 'BDB6A0728C8DFE2124536F16F2BA428FE767A8F9',
                timestamp: '2020-11-05T09:15:18.445494105Z',
                signature:
                  'mcUMQtCR38MK69IeUDri0zkfllsXKgnVFTsFwNaO/7cnBaIUUz9U4d3H9EjSH4kANJxWRFO3dSnPy1uOD36K6A==',
              },
              {
                block_id_flag: 3,
                validator_address: 'F0C23921727D869745C4F9703CF33996B1D2B715',
                timestamp: '2020-11-05T09:15:18.537783045Z',
                signature:
                  'fpr26xz+Gg5Rl7Fvx34a0QZpb5yJc5P4t5Z1OctIDQ0VMmh9vEWagsqQGErt1bm+CaKFtkFfZZ4CU0DKN27GbQ==',
              },
              {
                block_id_flag: 3,
                validator_address: 'F23391B5DBF982E37FB7DADEA64AAE21CAE4C172',
                timestamp: '2020-11-05T09:15:18.539946947Z',
                signature:
                  'KGsiIaralMMr1M9A7Ulhbc0THt1pLgNIrNQ6Kx+oGtwjl2w5ke5iivAAtZMduhyIAUMhrp6PN7ZvKgv9TPCNNw==',
              },
            ],
          },
        },
      },
    }
    mockedAxios.get.mockResolvedValue(resp)
    const expected = {
      blockID: {
        hash: Buffer.from(
          '391E99908373F8590C928E0619956DA3D87EB654445DA4F25A185C9718561D53',
          'hex',
        ),
      },
      block: {
        header: {
          chainID: 'bandchain',
          height: 1032007,
          time: 1604567718,
          lastCommitHash: Buffer.from(
            '17B2CE4ABA910E85847537F1323DB95C9F16C20C60E9B9BBB04C633C3125BD92',
            'hex',
          ),
          dataHash: Buffer.from(
            'EFE5E3F549554FEE8EB9B393740C250D74580427A96A175ABB105806039CFFE2',
            'hex',
          ),
          validatorsHash: Buffer.from(
            'E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151',
            'hex',
          ),
          nextValidatorsHash: Buffer.from(
            'E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151',
            'hex',
          ),
          consensusHash: Buffer.from(
            '0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45',
            'hex',
          ),
          appHash: Buffer.from(
            '6E2B1ECE9D912D86C25182E8B7419583ABCE978BFC66DC2556BB0D06A8D528EF',
            'hex',
          ),
          lastResultsHash: Buffer.from(''),
          evidenceHash: Buffer.from(''),
          proposerAddress: Buffer.from(
            'BDB6A0728C8DFE2124536F16F2BA428FE767A8F9',
            'hex',
          ),
        },
      },
    }
    const response = client.getLatestBlock()

    response.then((e) => expect(e).toEqual(expected))
  })
})

describe('Client get account', () => {
  it('success', () => {
    const resp = {
      data: {
        height: '650788',
        result: {
          type: 'cosmos-sdk/Account',
          value: {
            address: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
            coins: [{ denom: 'uband', amount: 104082359107 }],
            public_key: {
              type: 'tendermint/PubKeySecp256k1',
              value: 'A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8',
            },
            account_number: '36',
            sequence: '927',
          },
        },
      },
    }

    mockedAxios.get.mockResolvedValue(resp)

    const expected = {
      address: Address.fromAccBech32(
        'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
      ),
      coins: [new Coin(104082359107, 'uband')],
      publicKey: {
        type: 'tendermint/PubKeySecp256k1',
        value: 'A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8',
      },
      accountNumber: 36,
      sequence: 927,
    }

    const response = client.getAccount(
      Address.fromAccBech32('band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte'),
    )
    response.then((e) => expect(e).toEqual(expected))
  })
})

describe('Client get reporters', () => {
  it('success', () => {
    const resp = {
      data: {
        height: '2245131',
        result: [
          'band1yyv5jkqaukq0ajqn7vhkyhpff7h6e99ja7gvwg',
          'band19nf0sqnjycnvpexlxs6hjz9qrhhlhsu9pdty0r',
          'band1fndxcmg0h5vhr8cph7gryryqfn9yqp90lysjtm',
          'band10ec0p96j60duce5qagju5axuja0rj8luqrzl0k',
          'band15pm9scujgkpwpy2xa2j53tvs9ylunjn0g73a9s',
          'band1cehe3sxk7f4rmvwdf6lxh3zexen7fn02zyltwy',
        ],
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const expected = [
      Address.fromAccBech32('band1yyv5jkqaukq0ajqn7vhkyhpff7h6e99ja7gvwg'),
      Address.fromAccBech32('band19nf0sqnjycnvpexlxs6hjz9qrhhlhsu9pdty0r'),
      Address.fromAccBech32('band1fndxcmg0h5vhr8cph7gryryqfn9yqp90lysjtm'),
      Address.fromAccBech32('band10ec0p96j60duce5qagju5axuja0rj8luqrzl0k'),
      Address.fromAccBech32('band15pm9scujgkpwpy2xa2j53tvs9ylunjn0g73a9s'),
      Address.fromAccBech32('band1cehe3sxk7f4rmvwdf6lxh3zexen7fn02zyltwy'),
    ]

    const response = client.getReporters(
      Address.fromValBech32(
        'bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre',
      ),
    )
    response.then((e) => expect(e).toEqual(expected))
  })
})

describe('get price symbols', () => {
  const resp = {
    data: {
      height: '2951872',
      result: ['2KEY', 'ABYSS', 'ADA', 'AKRO'],
    },
  }

  it('success', () => {
    mockedAxios.get.mockResolvedValue(resp)

    const expected = ['2KEY', 'ABYSS', 'ADA', 'AKRO']
    const response = client.getPriceSymbols(3, 4)
    response.then((e) => expect(e).toEqual(expected))
  })
})
