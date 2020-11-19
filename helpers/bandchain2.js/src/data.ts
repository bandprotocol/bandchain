import { Address } from './wallet'

export type HexBytes = Buffer

export class Coin {
  amount: number
  denom: string

  constructor(amount: number, denom: string) {
    this.amount = amount
    this.denom = denom
  }

  asJson() {
    return { amount: this.amount.toString(), denom: this.denom }
  }
}

export interface DataSource {
  owner: Address
  name: string
  description: string
  fileName: string
}
