import { Coin } from 'data'
import { Buffer } from 'buffer'
import { Address } from 'wallet'

export abstract class Msg {
  abstract asJson(): { type: string; value: any }
  abstract getSender(): Address
  abstract validate(): boolean
}
export class MsgRequest extends Msg {
  oracleScriptID: number
  calldata: Buffer
  askCount: number
  minCount: number
  clientID: string
  sender: Address

  constructor(
    oracleScriptID: number,
    calldata: string,
    askCount: number,
    minCount: number,
    clientID: string,
    sender: Address,
  ) {
    super()
    this.oracleScriptID = oracleScriptID
    this.calldata = Buffer.from(calldata, 'hex')
    this.askCount = askCount
    this.minCount = minCount
    this.clientID = clientID
    this.sender = sender
  }

  asJson() {
    return {
      type: 'oracle/Request',
      value: {
        oracle_script_id: String(this.oracleScriptID),
        calldata: this.calldata.toString('base64'),
        ask_count: this.askCount.toString(),
        min_count: this.minCount.toString(),
        client_id: this.clientID,
        sender: this.sender,
      },
    }
  }

  getSender() {
    return this.sender
  }

  validate() {
    // TODO: validate
    return true
  }
}

export class MsgSend extends Msg {
  fromAddress: Address
  toAddress: Address
  amount: Coin[]

  constructor(from: Address, to: Address, amount: Coin[]) {
    super()
    this.fromAddress = from
    this.toAddress = to
    this.amount = amount
  }

  asJson() {
    return {
      type: 'cosmos-sdk/MsgSend',
      value: {
        amount: this.amount.map((each) => each.asJson()),
        from_address: this.fromAddress.toAccBech32(),
        to_address: this.toAddress.toAccBech32(),
      },
    }
  }

  getSender() {
    return this.fromAddress
  }

  validate() {
    if (this.amount.length == 0) {
      throw Error('Expect at least 1 coin')
    }
    // TODO: Uncomment this when coin.validate() is ready
    // this.amount.forEach(coin => coin.validate())
    return true
  }
}
