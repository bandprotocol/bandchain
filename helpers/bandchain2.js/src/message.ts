import { Coin } from 'data'
import { Buffer } from 'buffer'

abstract class Msg {
  abstract asJson(): {type: string, value: any}
}

export class MsgRequest extends Msg {
  oracleScriptID: number
  calldata: Buffer
  askCount: number
  minCount: number
  clientID: string
  sender: string

  constructor(oracleScriptID: number, calldata: string, askCount: number, minCount: number, clientID: string, sender: string) {
    super()
    this.oracleScriptID = oracleScriptID
    this.calldata = Buffer.from(calldata, 'hex')
    this.askCount = askCount
    this.minCount = minCount
    this.clientID = clientID
    this.sender = sender
  }

  asJson() {
    return{
      type: 'oracle/Request',
      value: {
        oracle_script_id: String(this.oracleScriptID),
        calldata: this.calldata.toString('base64'),
        ask_count: this.askCount.toString(),
        min_count: this.minCount.toString(),
        client_id: this.clientID,
        sender: this.sender
      }
    }
  }
}

export class MsgSend extends Msg {
  fromAddress: string
  toAddress: string
  amount: Coin[]

  constructor(from: string, to: string, amount: Coin[]) {
    super()
    this.fromAddress = from
    this.toAddress = to
    this.amount = amount
  }

  asJson() {
    return {
      type: 'cosmos-sdk/MsgSend',
      value: {
        to_address: this.toAddress,
        from_address: this.fromAddress,
        amount: this.amount.map((each) => each.asJson()),
      },
    }
  }
}
