import axios from 'axios'
import { Client } from '../../src'
import { Coin } from '../../src/data'
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
