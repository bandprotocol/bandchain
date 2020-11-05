import { Coin } from 'data';
declare abstract class Msg {
    abstract asJson(): any;
}
export declare class MsgSend extends Msg {
    fromAddress: string;
    toAddress: string;
    amount: Coin[];
    constructor(from: string, to: string, amount: Coin[]);
    asJson(): {
        type: string;
        value: {
            to_address: string;
            from_address: string;
            amount: {
                amount: string;
                denom: string;
            }[];
        };
    };
}
export {};
