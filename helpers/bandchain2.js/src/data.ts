import { Address } from 'wallet'

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

export interface TransactionSyncMode {
  txHash: HexBytes
  code: number
  errorLog?: string
  /*
  constructor(txHash: HexBytes, code: number, errorLog?: string) {
    if (!Number.isInteger(code)) throw Error('code is not an integer')

    this.txHash = txHash
    this.code = code
    this.errorLog = errorLog
  }
  */
}

export interface TransactionAsyncMode {
  txHash: HexBytes

  // constructor(txHash: HexBytes) {
  //   this.txHash = txHash
  // }
}

export interface TransactionBlockMode {
  height: number
  txHash: HexBytes
  gasWanted: number
  gasUsed: number
  code: number
  log: object[]
  errorLog?: string

  // constructor(
  //   height: number,
  //   txHash: HexBytes,
  //   gasWanted: number,
  //   gasUsed: number,
  //   code: number,
  //   log: object[],
  //   errorLog?: string,
  // ) {
  //   this.height = height
  //   this.txHash = txHash
  //   this.gasWanted = gasWanted
  //   this.gasUsed = gasUsed
  //   this.code = code
  //   this.log = log
  //   this.errorLog = errorLog
  // }
}
