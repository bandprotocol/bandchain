import axios from 'axios'
import {
  DataSource,
  TransactionSyncMode,
  TransactionAsyncMode,
  TransactionBlockMode,
  Coin,
  Account,
} from './data'
import { Address } from './wallet'

export default class Client {
  rpcUrl: string
  constructor(rpcUrl: string) {
    this.rpcUrl = rpcUrl
  }

  private async get(path: string, param?: object) {
    const response = await axios.get(`${this.rpcUrl}${path}`, param)
    return response.data
  }

  private async post(path: string, param: object) {
    const response = await axios.post(`${this.rpcUrl}${path}`, param)
    return response.data
  }

  private async getResult(path: string, param?: object) {
    const response = await this.get(`${path}`, param)
    return response.result
  }

  private async postResult(path: string, param: object) {
    const response = await this.post(`${path}`, param)
    return response.result
  }

  async getChainID(): Promise<string> {
    const response = await this.get('/bandchain/chain_id')
    return response.chain_id
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
      params: {
        oid: oid,
        calldata: calldata,
        min_count: minCount,
        ask_count: askCount,
      },
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


  async getReferenceData(pairs: string[]): Promise<any[]> {
    let symbolSet: Set<string> = new Set()
    pairs.forEach((pair) => {
      let symbols = pair.split('/')
      symbols.forEach((symbol: string) => {
        if (symbol == 'USD') return
        symbolSet.add(symbol)
      })
    })
    let symbolList: string[] = Array.from(symbolSet)
    let pricerBody = {
      symbols: symbolList,
      min_count: 3,
      ask_count: 4,
    }
    let priceData = await this.postResult('/oracle/request/prices', pricerBody)

    let symbolMap: any = {}
    symbolMap['USD'] = {
      symbol: 'USD',
      multiplier: '1000000000',
      px: '1000000000',
      resolve_time: Math.round(Date.now() / 1000).toString(),
    }
    priceData.map((price: any) => {
      symbolMap[price.symbol] = price
    })
    let data: any[] = []
    pairs.forEach((pair) => {
      let [baseSymbol, quoteSymbol] = pair.split('/')

      data.push({
        pair,
        rate:
          (symbolMap[baseSymbol].px * symbolMap[quoteSymbol].multiplier) /
          (symbolMap[quoteSymbol].px * symbolMap[baseSymbol].multiplier),
        updated: {
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
}
