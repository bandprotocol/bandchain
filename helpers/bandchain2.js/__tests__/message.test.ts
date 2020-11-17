import { Message, Data } from '../src/index'
import { MsgDelegate } from '../src/message'
import { Address } from '../src/wallet'

const { MsgSend, MsgRequest } = Message
const { Coin } = Data

const coin = new Coin(100000, 'uband')

describe('MsgRequest', () => {
  const sender_addr = Address.fromAccBech32(
    'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
  )

  const calldata = Buffer.from('000000034254430000000000000001', 'hex')
  const memo = 'from_bandchain.js'

  it('create successfully', () => {
    const msgRequest = new MsgRequest(1, calldata, 2, 2, memo, sender_addr)
    expect(msgRequest.asJson()).toEqual({
      type: 'oracle/Request',
      value: {
        oracle_script_id: '1',
        calldata: 'AAAAA0JUQwAAAAAAAAAB',
        ask_count: '2',
        min_count: '2',
        client_id: 'from_bandchain.js',
        sender: 'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
      },
    })

    expect(msgRequest.asJson()).toEqual({
      type: 'oracle/Request',
      value: {
        oracle_script_id: '1',
        calldata: 'AAAAA0JUQwAAAAAAAAAB',
        ask_count: '2',
        min_count: '2',
        client_id: 'from_bandchain.js',
        sender: 'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
      },
    })

    expect(msgRequest.getSender().toAccBech32()).toEqual(
      'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
    )

    expect(msgRequest.validate()).toBeTruthy()
  })

  it('create with error from validate()', () => {
    const msgs = []
    const errorText: string[] = []
    msgs.push(new MsgRequest(-1, calldata, 2, 2, memo, sender_addr))
    msgs.push(new MsgRequest(1.1, calldata, 2, 2, memo, sender_addr))
    msgs.push(new MsgRequest(1, calldata, 2.1, 2, memo, sender_addr))
    msgs.push(new MsgRequest(1, calldata, 2, 2.1, memo, sender_addr))
    msgs.push(new MsgRequest(1, calldata, 2, 0, memo, sender_addr))
    errorText.push('oracleScriptID cannot less than zero')
    errorText.push('oracleScriptID is not an integer')
    errorText.push('askCount is not an integer')
    errorText.push('minCount is not an integer')
    errorText.push('invalid minCount got: minCount: 0')

    msgs.forEach((msg, index) => {
      expect(() => {
        msg.validate()
      }).toThrowError(errorText[index])
    })
  })
})

describe('MsgSend', () => {
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
    expect(msgSend.getSender().toAccBech32()).toEqual(
      'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
    )
    expect(msgSend.validate()).toBeTruthy()
  })
})

describe('MsgDelegate', () => {
  it('create successfully', () => {
    const msgDelegate = new MsgDelegate(
      Address.fromAccBech32('band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c'),
      Address.fromValBech32(
        'bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst',
      ),
      coin,
    )

    expect(msgDelegate.asJson()).toEqual({
      type: 'cosmos-sdk/MsgDelegate',
      value: {
        amount: {
          amount: '100000',
          denom: 'uband',
        },
        delegator_address: 'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
        validator_address: 'bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst',
      },
    })

    expect(msgDelegate.getSender().toAccBech32()).toEqual(
      'band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c',
    )

    expect(msgDelegate.validate()).toBeTruthy()
  })
})
