export class Coin {
    constructor(amount, denom) {
        this.amount = amount;
        this.denom = denom;
    }
    asJson() {
        return { amount: this.amount.toString(), denom: this.denom };
    }
}
