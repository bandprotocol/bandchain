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
  OracleScript,
  RequestInfo,
  HexBytes,
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
    if (!Number.isInteger(id)) throw Error('id is not an integer')

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

  /**
   * Get the oracle script by ID
   * @param id Oracle script ID
   * @returns A Promise of Oracle.
   */

  async getOracleScript(id: number): Promise<OracleScript> {
    if (!Number.isInteger(id)) throw Error('id is not an integer')

    const response = await this.getResult(`/oracle/oracle_scripts/${id}`)
    return {
      owner: Address.fromAccBech32(response.owner),
      name: response.name,
      description: response.description,
      fileName: response.filename,
      schema: response.schema,
      sourceCodeUrl: response.source_code_url,
    }
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

  async getReferenceData(pairs: string[], minCount: number, askCount: number): Promise<ReferenceData[]> {
    if (!Number.isInteger(minCount)) throw Error('minCount is not an integer')
    if (!Number.isInteger(askCount)) throw Error('askCount is not an integer')
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
      min_count: minCount,
      ask_count: askCount,
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

  /**
   * Get the latest request
   * @param oid Oracle script ID
   * @param calldata The input parameters associated with the request
   * @param minCount The minimum number of validators necessary for the request to proceed to the execution phase
   * @param askCount The number of validators that are requested to respond to this request
   * @returns  A Promise of RequestInfo.
   */
  async getLatestRequest(
    oid: number,
    calldata: Buffer,
    minCount: number,
    askCount: number,
  ): Promise<RequestInfo> {
    if (!Number.isInteger(oid)) throw Error('oid is not an integer')
    if (!Number.isInteger(minCount)) throw Error('minCount is not an integer')
    if (!Number.isInteger(askCount)) throw Error('askCount is not an integer')

    const response = await this.getResult(`/oracle/request_search`, {
      params: {
        oid: oid,
        calldata: calldata.toString('base64'),
        min_count: minCount,
        ask_count: askCount,
      },
    })
    return {
      request: {
        oracleScriptID: parseInt(response.request.oracle_script_id),
        requestedValidators: response.request.requested_validators,
        minCount: parseInt(response.request.min_count),
        requestHeight: parseInt(response.request.request_height),
        clientID: response.request.client_id ? response.request.client_id : '',
        calldata: response.request.calldata
          ? Buffer.from(response.request.calldata, 'base64')
          : Buffer.from(''),
        rawRequests: response.request.raw_requests.map(
          (request: { [key: string]: any }) => {
            return {
              externalID: request.external_id
                ? parseInt(request.external_id)
                : 0,
              dataSourceID: parseInt(request.data_source_id),
              calldata: request.calldata
                ? Buffer.from(request.calldata, 'base64')
                : Buffer.from(''),
            }
          },
        ),
      },
      reports: response.reports?.map((report: { [key: string]: any }) => {
        return {
          validator: report.validator,
          inBeforeResolve: !!report.in_before_resolve,
          rawReports: report.raw_reports.map(
            (rawReport: { [key: string]: any }) => {
              return {
                externalID: rawReport.external_id
                  ? parseInt(rawReport.external_id)
                  : 0,
                data: rawReport.data
                  ? Buffer.from(rawReport.data, 'base64')
                  : Buffer.from(''),
              }
            },
          ),
        }
      }),
      result: response.result && {
        requestPacketData: {
          clientID: response.result.request_packet_data.client_id
            ? response.result.request_packet_data.client_id
            : '',
          askCount: parseInt(response.result.request_packet_data.ask_count),
          minCount: parseInt(response.result.request_packet_data.min_count),
          oracleScriptID: parseInt(
            response.result.request_packet_data.oracle_script_id,
          ),
          calldata: response.result.request_packet_data.calldata
            ? Buffer.from(
                response.result.request_packet_data.calldata,
                'base64',
              )
            : Buffer.from(''),
        },
        responsePacketData: {
          requestID: parseInt(response.result.response_packet_data.request_id),
          requestTime: parseInt(
            response.result.response_packet_data.request_time,
          ),
          resolveTime: parseInt(
            response.result.response_packet_data.resolve_time,
          ),
          resolveStatus: response.result.response_packet_data.resolve_status,
          ansCount: response.result.response_packet_data.ans_count
            ? parseInt(response.result.response_packet_data.ans_count)
            : 0,
          clientID: response.result.response_packet_data.client_id
            ? response.result.response_packet_data.client_id
            : '',
          result: response.result.response_packet_data.result
            ? Buffer.from(response.result.response_packet_data.result, 'base64')
            : Buffer.from(''),
        },
      },
    }
  }

  /**
   * Get the latest request
   * @param id Request ID
   * @returns  A Promise of RequestInfo.
   */

  async getRequestByID(id: number): Promise<RequestInfo> {
    if (!Number.isInteger(id)) throw Error('id is not an integer')

    const response = await this.getResult(`/oracle/requests/${id}`)
    return {
      request: {
        oracleScriptID: parseInt(response.request.oracle_script_id),
        requestedValidators: response.request.requested_validators,
        minCount: parseInt(response.request.min_count),
        requestHeight: parseInt(response.request.request_height),
        clientID: response.request.client_id ? response.request.client_id : '',
        calldata: response.request.calldata
          ? Buffer.from(response.request.calldata, 'base64')
          : Buffer.from(''),
        rawRequests: response.request.raw_requests.map(
          (request: { [key: string]: any }) => {
            return {
              externalID: request.external_id
                ? parseInt(request.external_id)
                : 0,
              dataSourceID: parseInt(request.data_source_id),
              calldata: request.calldata
                ? Buffer.from(request.calldata, 'base64')
                : Buffer.from(''),
            }
          },
        ),
      },
      reports: response.reports?.map((report: { [key: string]: any }) => {
        return {
          validator: report.validator,
          inBeforeResolve: !!report.in_before_resolve,
          rawReports: report.raw_reports.map(
            (rawReport: { [key: string]: any }) => {
              return {
                externalID: rawReport.external_id
                  ? parseInt(rawReport.external_id)
                  : 0,
                data: rawReport.data
                  ? Buffer.from(rawReport.data, 'base64')
                  : Buffer.from(''),
              }
            },
          ),
        }
      }),
      result: response.result && {
        requestPacketData: {
          clientID: response.result.request_packet_data.client_id
            ? response.result.request_packet_data.client_id
            : '',
          askCount: parseInt(response.result.request_packet_data.ask_count),
          minCount: parseInt(response.result.request_packet_data.min_count),
          oracleScriptID: parseInt(
            response.result.request_packet_data.oracle_script_id,
          ),
          calldata: response.result.request_packet_data.calldata
            ? Buffer.from(
                response.result.request_packet_data.calldata,
                'base64',
              )
            : Buffer.from(''),
        },
        responsePacketData: {
          requestID: parseInt(response.result.response_packet_data.request_id),
          requestTime: parseInt(
            response.result.response_packet_data.request_time,
          ),
          resolveTime: parseInt(
            response.result.response_packet_data.resolve_time,
          ),
          resolveStatus: response.result.response_packet_data.resolve_status,
          ansCount: response.result.response_packet_data.ans_count
            ? parseInt(response.result.response_packet_data.ans_count)
            : 0,
          clientID: response.result.response_packet_data.client_id
            ? response.result.response_packet_data.client_id
            : '',
          result: response.result.response_packet_data.result
            ? Buffer.from(response.result.response_packet_data.result, 'base64')
            : Buffer.from(''),
        },
      },
    }
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
