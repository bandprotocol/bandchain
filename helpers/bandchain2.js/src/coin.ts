export default class Coin {
  amount: number
  denom: string

  constructor(amount: number, denom: string) {
    this.amount = amount
    this.denom = denom
  }

  static fromJson(coin: { amount: number; denom: string }): Coin {
    return new Coin(coin.amount, coin.denom)
  }

  asJson(): { amount: string; denom: string } {
    return { amount: this.amount.toString(), denom: this.denom }
  }

  validate(): boolean {
    if (!Number.isInteger(this.amount)) throw Error('amount is not an integer')
    if (this.amount < 0) throw Error('Expect amount more than 0')
    if (this.denom.length === 0) throw Error('Expect denom')

    return true
  }
}
