import axios from 'axios'
import {
  DataSource,
  TransactionSyncMode,
  // HexBytes,
  TransactionAsyncMode,
  TransactionBlockMode,
} from 'data'
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
    const response = await axios.post(`${this.rpcUrl}`, param)
    console.log(path)
    console.log('_post', response)
    console.log('_post', response.data)
    return response.data
  }

  private async getResult(path: string, param?: object) {
    const response = await this.get(`${path}`, param)
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
    const response = await this.getResult(`/oracle/data_sources/${id}`, {})
    return {
      owner: Address.fromAccBech32(response.owner),
      name: response.name,
      description: response.description,
      fileName: response.filename,
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

    console.log(response.data)
    if (response.code) {
      code = parseInt(response.code)
      errorLog = response.raw_log
    }

    return {
      txHash: Buffer.from(response['txhash'], 'hex'),
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
      txHash: Buffer.from(respose['txhash'], 'hex'),
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

    if (response['code']) {
      code = parseInt(response['code'])
      errorLog = response['raw_log']
    } else {
      log = response['logs']
    }

    return {
      height: parseInt(response['height']),
      txHash: Buffer.from(response['txhash'], 'hex'),
      gasWanted: parseInt(response['gas_wanted']),
      gasUsed: parseInt(response['gas_wanted']),
      code,
      log,
      errorLog,
    }
  }
}
