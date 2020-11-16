import { Coin } from 'data'
import { Buffer, constants } from 'buffer'
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
        sender: this.sender.toAccBech32(),
      },
    }
  }

  getSender() {
    return this.sender
  }

  validate() {
    if (this.oracleScriptID <= 0)
      throw Error('oracleScriptID cannot less than zero')
    if (!Number.isInteger(this.oracleScriptID))
      throw Error('oracleScriptID is not an integer')
    if (this.calldata.length > constants.MAX_LENGTH)
      throw Error('too large calldata')
    if (!Number.isInteger(this.askCount))
      throw Error('askCount is not an integer')
    if (!Number.isInteger(this.minCount))
      throw Error('minCount is not an integer')
    if (this.minCount <= 0)
      throw Error(`invalid minCount got: minCount: ${this.minCount}`)
    if (this.askCount < this.minCount)
      throw Error(
        `invalida askCount got: minCount: ${this.minCount}, askCount: ${this.askCount}`,
      )

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
