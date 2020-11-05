import { MsgSend } from './src/message'
import { Coin } from './src/data'

const amount = [new Coin(10000, 'uband')]
const msgSend = new MsgSend('asdkaskd', 'asdjasdkj', amount)

let result = msgSend.asJson()
console.log(JSON.stringify(result))
