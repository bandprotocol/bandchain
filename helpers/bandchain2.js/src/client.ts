import axios from 'axios'
import { Address } from './wallet'
import { DataSource, HexBytes } from './data'

export default class Client {
  rpcUrl: string
  constructor(rpcUrl: string) {
    this.rpcUrl = rpcUrl
  }

  private async get(path: string, param?: object) {
    const response = await axios.get(`${this.rpcUrl}${path}`, param)
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
    const response = await this.getResult(`/oracle/data_sources/${id}`)
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

  async getRequestIDByTxHash(txHash: HexBytes): Promise<number[]> {
    const response = await this.get(`/txs/${txHash.toString('hex')}`)
    const msgs = response.logs
    const requestIDs: number[] = []
    msgs.forEach((msg: any) => {
      const requestEvent = msg.events.filter(
        (event: any) => event.type == 'request',
      )
      if (requestEvent.length == 1) {
        const attrID = requestEvent[0].attributes.find(
          ({ key }: any) => key === 'id',
        ).value
        requestIDs.push(parseInt(attrID))
      }
    })

    if (requestIDs.length == 0) {
      throw new Error('There is no request message in this tx')
    }

    return requestIDs
  }
}
