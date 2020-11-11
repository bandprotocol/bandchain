import { Message, Data, Transaction } from '../src/index'

const { MsgSend } = Message
const { Coin } = Data

describe('Transaction', () => {
  const coin = new Coin(100000, 'uband')
  

  it('create successfully', () => {
    const msgSend = new MsgSend(
      'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
      'band1ksnd0f3xjclvg0d4z9w0v9ydyzhzfhuy47yx79',
      [coin]
    )
  
    const tsc = new Transaction([msgSend])

    expect(tsc.accountNum).toBeUndefined()
    expect(tsc.chainID).toBeUndefined()
    expect(tsc.sequence).toBeUndefined()
    expect(tsc.fee).toEqual(0)
    expect(tsc.gas).toEqual(200000)
  })

})