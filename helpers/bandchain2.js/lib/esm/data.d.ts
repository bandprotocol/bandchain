export declare class Coin {
    amount: number;
    denom: string;
    constructor(amount: number, denom: string);
    asJson(): {
        amount: string;
        denom: string;
    };
}
