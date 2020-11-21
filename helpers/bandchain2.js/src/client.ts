import axios from 'axios'
import { Address } from './wallet'
import { DataSource, OracleScript, RequestInfo, HexBytes } from './data'

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

  /**
   * Get the oracle script by ID
   * @param id Oracle script ID
   * @returns A Promise of Oracle.
   */

  async getOracleScript(id: number): Promise<OracleScript> {
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
      result: {
        requestPacketData: {
          clientID: response.result?.request_packet_data.client_id
            ? response.result?.request_packet_data.client_id
            : '',
          askCount: parseInt(response.result?.request_packet_data.ask_count),
          minCount: parseInt(response.result?.request_packet_data.min_count),
          oracleScriptID: parseInt(
            response.result?.request_packet_data.oracle_script_id,
          ),
          calldata: response.result?.request_packet_data.calldata
            ? Buffer.from(
                response.result?.request_packet_data.calldata,
                'base64',
              )
            : Buffer.from(''),
        },
        responsePacketData: {
          requestID: parseInt(response.result?.response_packet_data.request_id),
          requestTime: parseInt(
            response.result?.response_packet_data.request_time,
          ),
          resolveTime: parseInt(
            response.result?.response_packet_data.resolve_time,
          ),
          resolveStatus: response.result?.response_packet_data.resolve_status,
          ansCount: response.result?.response_packet_data.ans_count
            ? parseInt(response.result?.response_packet_data.ans_count)
            : 0,
          clientID: response.result?.response_packet_data.client_id
            ? response.result?.response_packet_data.client_id
            : '',
          result: response.result?.response_packet_data.result
            ? Buffer.from(
                response.result?.response_packet_data.result,
                'base64',
              )
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
      result: {
        requestPacketData: {
          clientID: response.result?.request_packet_data.client_id
            ? response.result?.request_packet_data.client_id
            : '',
          askCount: parseInt(response.result?.request_packet_data.ask_count),
          minCount: parseInt(response.result?.request_packet_data.min_count),
          oracleScriptID: parseInt(
            response.result?.request_packet_data.oracle_script_id,
          ),
          calldata: response.result?.request_packet_data.calldata
            ? Buffer.from(
                response.result?.request_packet_data.calldata,
                'base64',
              )
            : Buffer.from(''),
        },
        responsePacketData: {
          requestID: parseInt(response.result?.response_packet_data.request_id),
          requestTime: parseInt(
            response.result?.response_packet_data.request_time,
          ),
          resolveTime: parseInt(
            response.result?.response_packet_data.resolve_time,
          ),
          resolveStatus: response.result?.response_packet_data.resolve_status,
          ansCount: response.result?.response_packet_data.ans_count
            ? parseInt(response.result?.response_packet_data.ans_count)
            : 0,
          clientID: response.result?.response_packet_data.client_id
            ? response.result?.response_packet_data.client_id
            : '',
          result: response.result?.response_packet_data.result
            ? Buffer.from(
                response.result?.response_packet_data.result,
                'base64',
              )
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
