import axios from 'axios'
import { DataSource } from 'data'

export default class Client {
  rpcUrl: string
  constructor(rpcUrl: string) {
    this.rpcUrl = rpcUrl
  }

  private async get(path: string, ...params: object[]) {
    const response = await axios.get(`${this.rpcUrl}${path}`, ...params)
    return response.data
  }

  private async getResult(path: string, ...params: object[]) {
    const response = await this.get(`${path}`, ...params)
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
      owner: response.owner,
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
      oid: oid,
      calldata: calldata,
      min_count: minCount,
      ask_count: askCount,
    })
    return response
  }
}
