class Msg {
}
export class MsgSend extends Msg {
    constructor(from, to, amount) {
        super();
        this.fromAddress = from;
        this.toAddress = to;
        this.amount = amount;
    }
    asJson() {
        return {
            type: 'cosmos-sdk/MsgSend',
            value: {
                to_address: this.toAddress,
                from_address: this.fromAddress,
                amount: this.amount.map((each) => each.asJson()),
            },
        };
    }
}
