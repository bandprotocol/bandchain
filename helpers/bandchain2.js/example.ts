import { Message, Data } from './src/index'
import { Address, PrivateKey } from './src/wallet'

const { MsgSend } = Message
const { Coin } = Data

const amount = [new Coin(10000, 'uband')]
const from_addr = Address.fromAccBech32(
  'band1ksnd0f3xjclvg0d4z9w0v9ydyzhzfhuy47yx79',
)
const to_addr = Address.fromAccBech32(
  'band1p843hkdj2svjzm7zaceak07m9mtyf6hatcpvnl',
)
const msgSend = new MsgSend(from_addr, to_addr, amount)

let result = msgSend.asJson()
console.log(JSON.stringify(result))

// console.log(PrivateKey.generate())
const x = PrivateKey.fromMnemonic('s')
console.log(x.toHex())
