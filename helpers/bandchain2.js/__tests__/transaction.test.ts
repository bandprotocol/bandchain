import { Message, Data, Transaction } from '../src/index'

const { MsgSend } = Message
const { Coin } = Data

describe('Transaction', () => {
  const coin = new Coin(100000, 'uband')
  const msgSend = new MsgSend(
    'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
    'band1ksnd0f3xjclvg0d4z9w0v9ydyzhzfhuy47yx79',
    [coin]
  )
  

  it('create successfully', () => {
  
    const tsc = new Transaction()

    expect(tsc.accountNum).toBeUndefined()
    expect(tsc.chainID).toBeUndefined()
    expect(tsc.sequence).toBeUndefined()
    expect(tsc.fee).toEqual(0)
    expect(tsc.gas).toEqual(200000)
    expect(tsc.withMessages.length).toEqual(0)

  })

  it('able to work "with" function', () => {

    const tsc = new Transaction()
    .withMessages(msgSend)
    .withAccountNum(100)
    .withSequence(30)
    .withChainID('bandchain')
    .withGas(500000)
    .withFee(10)
    .withMemo('bandchain2.js test')

    expect(tsc.msgs.length).toEqual(1)
    expect(tsc.accountNum).toEqual(100)
    expect(tsc.sequence).toEqual(30)
    expect(tsc.chainID).toEqual('bandchain')
    expect(tsc.fee).toEqual(10)
    expect(tsc.gas).toEqual(500000)
    expect(tsc.memo).toEqual('bandchain2.js test')

  })

})