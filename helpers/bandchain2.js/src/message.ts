import { Coin } from 'data'
import { Address } from 'wallet'

export abstract class Msg {
  abstract asJson(): { type: string; value: any }
  abstract getSender(): Address
  abstract validate(): boolean
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
