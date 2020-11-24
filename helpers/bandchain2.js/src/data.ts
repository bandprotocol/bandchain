import { Address } from './wallet'
import { Coin } from './coin'

export type HexBytes = Buffer
export type EpochTime = number

export interface Account {
  address: Address
  coins: Coin[]
  publicKey?: object
  accountNumber: number
  sequence: number
}

export interface DataSource {
  owner: Address
  name: string
  description: string
  fileName: string
}

export interface ReferenceDataUpdated {
  base: number
  quote: number
}

export interface ReferenceData {
  pair: string
  rate: number
  updatedAt: ReferenceDataUpdated
}

export interface TransactionSyncMode {
  txHash: HexBytes
  code: number
  errorLog?: string
}

export interface TransactionAsyncMode {
  txHash: HexBytes
}

export interface TransactionBlockMode {
  height: number
  txHash: HexBytes
  gasWanted: number
  gasUsed: number
  code: number
  log: object[]
  errorLog?: string
}

export interface BlockHeaderInfo {
  chainID: string
  height: number
  time: EpochTime
  lastCommitHash: HexBytes
  dataHash: HexBytes
  validatorsHash: HexBytes
  nextValidatorsHash: HexBytes
  consensusHash: HexBytes
  appHash: HexBytes
  lastResultsHash: HexBytes
  evidenceHash: HexBytes
  proposerAddress: HexBytes
}
interface BlockHeader {
  header: BlockHeaderInfo
}

interface BlockID {
  hash: HexBytes
}
export interface Block {
  block: BlockHeader
  blockID: BlockID
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
