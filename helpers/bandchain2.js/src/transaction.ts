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
}