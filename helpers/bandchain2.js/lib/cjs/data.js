"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Coin = void 0;
class Coin {
    constructor(amount, denom) {
        this.amount = amount;
        this.denom = denom;
    }
    asJson() {
        return { amount: this.amount.toString(), denom: this.denom };
    }
}
exports.Coin = Coin;
