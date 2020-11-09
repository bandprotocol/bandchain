import { Message, Data } from '../src/index'

const { MsgSend, MsgRequest } = Message
const { Coin } = Data

describe('MsgRequest', () => {
  it('create successfully', () => {
    const msgRequest = new MsgRequest(
      1,
      "000000034254430000000000000001",
      2,
      2,
      "from_bandchain.js",
      "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"
    )
    expect(msgRequest.asJson()).toEqual({
      type: 'oracle/Request',
      value: {
        oracle_script_id: '1',
        calldata: 'AAAAA0JUQwAAAAAAAAAB',
        ask_count: '2',
        min_count: '2',
        client_id: 'from_bandchain.js',
        sender: 'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c'
      }
    })
  })

  
})

describe('MsgSend', () => {
  const coin = new Coin(100000, 'uband')

  it('create successfully', () => {
    const msgSend = new MsgSend(
      'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
      'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
      [coin],
    )
    expect(msgSend.asJson()).toEqual({
      type: 'cosmos-sdk/MsgSend',
      value: {
        to_address: 'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
        from_address: 'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
        amount: [{ amount: '100000', denom: 'uband' }],
      },
    })
  })
})
