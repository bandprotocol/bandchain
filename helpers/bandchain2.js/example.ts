import { Message, Data } from './src/index'

const { MsgSend } = Message
const { Coin } = Data

const amount = [new Coin(10000, 'uband')]
const msgSend = new MsgSend('asdkaskd', 'asdjasdkj', amount)

let result = msgSend.asJson()
console.log(JSON.stringify(result))
