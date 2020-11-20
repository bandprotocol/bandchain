import { Address } from './wallet'

export class Coin {
  amount: number
  denom: string

  constructor(amount: number, denom: string) {
    this.amount = amount
    this.denom = denom
  }

  asJson(): { amount: string; denom: string } {
    return { amount: this.amount.toString(), denom: this.denom }
  }

  validate(): boolean {
    if (this.amount < 0) throw Error('Expect amount more than 0')
    if (this.denom.length == 0) throw Error('Expect denom')

    return true
  }
}

export interface DataSource {
  owner: Address
  name: string
  description: string
  fileName: string
}
