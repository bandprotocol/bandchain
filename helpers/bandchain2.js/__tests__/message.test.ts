import { Coin } from '../src/data'
import { MsgSend } from '../src/message'

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
