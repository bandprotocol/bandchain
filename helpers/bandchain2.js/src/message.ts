import { Coin } from 'data'

export abstract class Msg {
  abstract asJson(): { type: string; value: any }
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
