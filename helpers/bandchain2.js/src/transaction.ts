import {Msg} from './message'

export default class Transaction {
  msgs: Msg[] = []
  accountNum?: number
  sequence?: number
  chainID?: string
  fee: number = 0
  gas: number = 200000
  memo: string = ""

  withMessages(...msg: Msg[]): Transaction {
    this.msgs.push(...msg)
    return this
  }

  withAccountNum(accountNum: number): Transaction {
    this.accountNum = accountNum
    return this
  }

  withSequence(sequence: number): Transaction {
    this.sequence = sequence
    return this
  }

  withChainID(chainID: string): Transaction {
    this.chainID = chainID
    return this
  }

  withFee(fee: number): Transaction {
    this.fee = fee
    return this
  }

  withGas(gas: number): Transaction {
    this.gas = gas
    return this
  }

  withMemo(memo: string): Transaction {
    this.memo = memo
    return this
  }

  getSignData(): Buffer {
    if (this.msgs.length == 0) {
      throw Error("message is empty")
    }

    if (this.accountNum == null) {
      throw Error("accountNum should be defined")
    }

    if (this.sequence == null) {
      throw Error("sequence should be defined")
    }

    if (this.chainID == null) {
      throw Error("chainID should be defined")
    }

    // TODO: Validate Msgs
    this.msgs.forEach(msg => msg.validate())

    let messageJson: {[key: string]: any } = {
      chain_id: this.chainID,
      account_number: this.accountNum.toString(),
      fee: {
        amount: [{
          amount: this.fee.toString(),
          denom: "uband"
        }],
        gas: this.gas.toString(),
      },
      memo: this.memo,
      sequence: this.sequence.toString(),
      msgs: this.msgs.map(msg => msg.asJson())
    }

    const sortedKey = Object.keys(messageJson).sort()
    const result: {[key: string]: any } = {}
    sortedKey.forEach(key => result[key] = messageJson[key])

    return Buffer.from(JSON.stringify(result))
  }

  getTxData(signature: Buffer, pubkey: string): Object {
    
    if (this.accountNum == null) {
      throw Error("accountNum should be defined")
    }

    if (this.sequence == null) {
      throw Error("sequence should be defined")
    }

    //TODO: pubkey from type PublicKey

    return {
      fee: {
        amount: [{"amount": this.fee.toString(), "denom": "uband"}],
        gas: this.gas.toString(),
      },
      memo: this.memo,
      msg: this.msgs.map(msg => msg.asJson()),
      signature: [{
        signature: signature.toString('utf-8'),
        pub_key: {
          type: "tendermint/PubKeySecp256k1",
          value: pubkey
        },
        account_number: this.accountNum.toString(),
        sequence: this.sequence.toString()
      }]
    }
  }

}