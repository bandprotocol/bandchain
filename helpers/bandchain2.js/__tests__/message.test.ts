import { Message, Data } from '../src/index'
import { Address } from '../src/wallet'

const { MsgSend } = Message
const { Coin } = Data

describe('MsgSend', () => {
  const coin = new Coin(100000, 'uband')

  it('create successfully', () => {
    const msgSend = new MsgSend(
      Address.fromAccBech32('band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c'),
      Address.fromAccBech32('band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c'),
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
