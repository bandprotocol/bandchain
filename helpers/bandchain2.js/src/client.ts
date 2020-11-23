import axios from 'axios'
import {
  DataSource,
  TransactionSyncMode,
  TransactionAsyncMode,
  TransactionBlockMode,
  Block,
  Coin,
  Account,
  ReferenceData,
} from './data'
import { Address } from './wallet'

export default class Client {
  rpcUrl: string
  constructor(rpcUrl: string) {
    this.rpcUrl = rpcUrl
  }

  private async get(path: string, params?: object) {
    let options
    if (params) options = { params }

    const response = await axios.get(`${this.rpcUrl}${path}`, options)
    return response.data
  }

  private async post(path: string, data: object) {
    const response = await axios.post(`${this.rpcUrl}${path}`, data)
    return response.data
  }

  private async getResult(path: string, params?: object) {
    const response = await this.get(`${path}`, params)
    return response.result
  }

  private async postResult(path: string, data: object) {
    const response = await this.post(`${path}`, data)
    return response.result
  }

  async getChainID(): Promise<string> {
    const response = await this.get('/bandchain/chain_id')
    return response.chain_id
  }

  async getLatestBlock(): Promise<Block> {
    const response = await this.get('/blocks/latest')
    const header = response.block.header

    return {
      block: {
        header: {
          chainID: header.chain_id,
          height: parseInt(header.height),
          time: Math.floor(new Date(header.time).getTime() / 1000),
          lastCommitHash: Buffer.from(header.last_commit_hash, 'hex'),
          dataHash: Buffer.from(header.data_hash, 'hex'),
          validatorsHash: Buffer.from(header.validators_hash, 'hex'),
          nextValidatorsHash: Buffer.from(header.next_validators_hash, 'hex'),
          consensusHash: Buffer.from(header.consensus_hash, 'hex'),
          appHash: Buffer.from(header.app_hash, 'hex'),
          lastResultsHash: Buffer.from(header.last_results_hash, 'hex'),
          evidenceHash: Buffer.from(header.evidence_hash, 'hex'),
          proposerAddress: Buffer.from(header.proposer_address, 'hex'),
        },
      },
      blockID: { hash: Buffer.from(response.block_id.hash, 'hex') },
    }
  }

  /**
   * Get the data source by ID
   * @param id Data source ID
   * @returns A Promise of DataSoruce.
   */

  async getDataSource(id: number): Promise<DataSource> {
    const response = await this.getResult(`/oracle/data_sources/${id}`)
    return {
      owner: Address.fromAccBech32(response.owner),
      name: response.name,
      description: response.description,
      fileName: response.filename,
    }
  }

  async getAccount(address: Address): Promise<Account> {
    const response = await this.getResult(
      `/auth/accounts/${address.toAccBech32()}`,
    )
    const value = response.value
    return {
      address: Address.fromAccBech32(value.address),
      coins: value.coins.map(
        (c: { amount: number; denom: string }) => new Coin(c.amount, c.denom),
      ),
      publicKey: value.public_key,
      accountNumber: parseInt(value.account_number),
      sequence: parseInt(value.sequence),
    }
  }

  async getLatestRequest(
    oid: number,
    calldata: string,
    minCount: number,
    askCount: number,
  ) {
    const response = await this.getResult(`/oracle/request_search`, {
      oid: oid,
      calldata: calldata,
      min_count: minCount,
      ask_count: askCount,
    })
    return response
  }

  /**
   * sendTxSyncMode
   * @param data
   */
  async sendTxSyncMode(data: object): Promise<TransactionSyncMode> {
    let response = await this.post('/txs', {
      tx: data,
      mode: 'sync',
    })

    let code = 0
    let errorLog

    if (response.code) {
      code = parseInt(response.code)
      errorLog = response.raw_log
    }

    return {
      txHash: Buffer.from(response.txhash, 'hex'),
      code,
      errorLog,
    }
  }

  /**
   * sendTxAsyncMode
   * @param data
   */
  async sendTxAsyncMode(data: object): Promise<TransactionAsyncMode> {
    let respose = await this.post('/txs', {
      tx: data,
      mode: 'async',
    })

    return {
      txHash: Buffer.from(respose.txhash, 'hex'),
    }
  }

  /**
   * sendTxBlockMode
   * @param data
   */
  async sendTxBlockMode(data: object): Promise<TransactionBlockMode> {
    let response = await this.post('/txs', {
      tx: data,
      mode: 'block',
    })

    let code = 0
    let errorLog
    let log = []

    if (response.code) {
      code = parseInt(response.code)
      errorLog = response.raw_log
    } else {
      log = response.logs
    }

    return {
      height: parseInt(response.height),
      txHash: Buffer.from(response.txhash, 'hex'),
      gasWanted: parseInt(response.gas_wanted),
      gasUsed: parseInt(response.gas_wanted),
      code,
      log,
      errorLog,
    }
  }

  async getReferenceData(pairs: string[]): Promise<ReferenceData[]> {
    let symbolSet: Set<string> = new Set()
    pairs.forEach((pair: string) => {
      let symbols = pair.split('/')
      symbols.forEach((symbol: string) => {
        if (symbol === 'USD') return
        symbolSet.add(symbol)
      })
    })
    let symbolList: string[] = Array.from(symbolSet)
    let pricerBody = {
      symbols: symbolList,
      min_count: 3,
      ask_count: 4,
    }
    let priceData = await this.postResult('/oracle/request_prices', pricerBody)

    let symbolMap: any = {}
    symbolMap['USD'] = {
      symbol: 'USD',
      multiplier: '1000000000',
      px: '1000000000',
      resolve_time: Math.floor(Date.now() / 1000).toString(),
    }
    priceData.map((price: any) => {
      symbolMap[price.symbol] = price
    })
    let data: ReferenceData[] = []
    pairs.forEach((pair) => {
      let [baseSymbol, quoteSymbol] = pair.split('/')

      data.push({
        pair,
        rate:
          (Number(symbolMap[baseSymbol].px) *
            Number(symbolMap[quoteSymbol].multiplier)) /
          (Number(symbolMap[quoteSymbol].px) *
            Number(symbolMap[baseSymbol].multiplier)),
        updatedAt: {
          base: Number(symbolMap[baseSymbol].resolve_time),
          quote: Number(symbolMap[quoteSymbol].resolve_time),
        },
      })
    })
    return data
  }

  async getReporters(validator: Address): Promise<Address[]> {
    let response = await this.getResult(
      `/oracle/reporters/${validator.toValBech32()}`,
    )
    return response.map((a: string) => Address.fromAccBech32(a))
  }

  async getPriceSymbols(minCount: number, askCount: number): Promise<string[]> {
    if (!Number.isInteger(minCount)) throw Error('minCount is not an integer')
    if (!Number.isInteger(askCount)) throw Error('askCount is not an integer')
    let response = await this.getResult('/oracle/price_symbols', {
      min_count: minCount,
      ask_count: askCount,
    })
    return response
  }
}
