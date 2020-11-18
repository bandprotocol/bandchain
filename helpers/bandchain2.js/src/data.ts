import { Address } from './wallet'

export type HexBytes = Buffer

export class Coin {
  amount: number
  denom: string

  constructor(amount: number, denom: string) {
    this.amount = amount
    this.denom = denom
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

export interface DataSource {
  owner: Address
  name: string
  description: string
  fileName: string
}

export interface OracleScript {
  owner: Address
  name: string
  description: string
  fileName: string
  schema: string
  sourceCodeUrl: string
}

interface RawRequest {
  dataSourceID: number
  externalID: number
  calldata: Buffer
}

interface Request {
  oracleScriptID: number
  requestedValidators: string[]
  minCount: number
  requestHeight: number
  rawRequests: RawRequest[]
  clientID: string
  calldata: Buffer
}

interface RawReport {
  externalID: number
  data: Buffer
}

interface Report {
  validator: string
  rawReports: RawReport[]
  inBeforeResolve: boolean
}

interface RequestPacketData {
  oracleScriptID: number
  askCount: number
  minCount: number
  clientID: string
  calldata: Buffer
}

interface ResponsePacketData {
  requestID: number
  requestTime: number
  resolveTime: number
  resolveStatus: number
  ansCount: number
  clientID: string
  result: Buffer
}

interface Result {
  requestPacketData: RequestPacketData
  responsePacketData: ResponsePacketData
}

export interface RequestInfo {
  request: Request
  reports?: Report[]
  result?: Result
}
