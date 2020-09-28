const axios = require('axios')
const { Obi } = require('@bandprotocol/obi.js')
const cosmosjs = require('@cosmostation/cosmosjs')
const delay = require('delay')

async function createRequestMsg(
  cosmos,
  sender,
  oracleScriptID,
  validatorCounts,
  calldata,
  chainID,
  fee,
) {
  const account = await cosmos.getAccounts(sender)
  return cosmos.newStdMsg({
    msgs: [
      {
        type: 'oracle/Request',
        value: {
          oracle_script_id: String(oracleScriptID),
          calldata: Buffer.from(calldata).toString('base64'),
          ask_count: String(validatorCounts.askCount),
          min_count: String(validatorCounts.minCount),
          client_id: 'bandchain.js',
          sender: sender,
        },
      },
    ],
    chain_id: chainID,
    fee,
    memo: '',
    account_number: String(account.result.value.account_number),
    sequence: String(account.result.value.sequence || 0),
  })
}

class BandChain {
  constructor(endpoint) {
    this.endpoint = endpoint
  }

  async _getChainID() {
    if (!this._chainID) {
      try {
        const res = await axios.get(`${this.endpoint}/bandchain/genesis`)
        this._chainID = res.data.chain_id
      } catch {
        throw new Error('Cannot retrieve chainID')
      }
    }
    return this._chainID
  }

  async getOracleScript(oracleScriptID) {
    try {
      const res = await axios.get(
        `${this.endpoint}/oracle/oracle_scripts/${oracleScriptID}`,
      )
      res.data.result.id = oracleScriptID
      return res.data.result
    } catch {
      throw new Error('No oracle script found with the given ID')
    }
  }

  async submitRequestTx(
    oracleScript,
    parameters,
    validatorCounts,
    mnemonic,
    gasAmount = 0,
    gasLimit = 1000000,
  ) {
    const chainID = await this._getChainID()
    const obiObj = new Obi(oracleScript.schema)
    const calldata = obiObj.encodeInput(parameters)

    const cosmos = cosmosjs.network(this.endpoint, chainID)
    cosmos.setPath("m/44'/494'/0'/0/0")
    cosmos.setBech32MainPrefix('band')
    const ecpairPriv = cosmos.getECPairPriv(mnemonic)
    const sender = cosmos.getAddress(mnemonic)

    let requestMsg = await createRequestMsg(
      cosmos,
      sender,
      oracleScript.id,
      validatorCounts,
      calldata,
      chainID,
      {
        amount: [{ amount: `${gasAmount}`, denom: 'uband' }],
        gas: `${gasLimit}`,
      },
    )

    let signedTx = cosmos.sign(requestMsg, ecpairPriv, 'block')

    const broadcastResponse = await cosmos.broadcast(signedTx)

    return this.getRequestID(broadcastResponse.txhash)
  }

  async getRequestID(txHash, retryTimeout = 200) {
    let requestEndpoint = `${this.endpoint}/txs/${txHash}`

    // Loop until the txHash is included in the block
    while (true) {
      let res
      try {
        res = await axios.get(requestEndpoint)
      } catch (e) {
        await delay(retryTimeout)
        continue
      }

      if (res.status == 200) {
        try {
          const requestID = res.data.logs[0].events
            .find(({ type }) => type === 'request')
            .attributes.find(({ key }) => key === 'id').value
          return requestID
        } catch {
          throw new Error('Not a request tx')
        }
      }
    }
  }

  async getRequestProof(requestID, retryTimeout = 200) {
    // Check if request exists
    try {
      const requestEndpoint = `${this.endpoint}/oracle/requests/${requestID}`
      await axios.get(requestEndpoint)
    } catch {
      throw new Error('Request not found')
    }

    // Try and wait for proof
    while (true) {
      try {
        const requestEndpoint = `${this.endpoint}/oracle/proof/${requestID}`
        let res = await axios.get(requestEndpoint)
        if (res.status == 200 && res.data.result.evmProofBytes != null) {
          let result = res.data.result
          return result
        } else if (res.status == 200 && res.data.result.evmProofBytes == null) {
          throw new Error('No proof found for the specified requestID')
        }
      } catch {
        await delay(retryTimeout)
      }
    }
  }

  async getRequestEVMProof(requestID, retryTimeout = 200) {
    let result = await this.getRequestProof(requestID, retryTimeout)
    return result.evmProofBytes
  }

  async getRequestNonEVMProof(requestID, retryTimeout = 200) {
    const proofSchema =
      '{client_id:string,oracle_script_id:u64,calldata:bytes,ask_count:u64,min_count:u64}/{client_id:string,request_id:u64,ans_count:u64,request_time:u64,resolve_time:u64,resolve_status:u8,result:bytes}'

    let result = await this.getRequestProof(requestID, retryTimeout)
    let requestPacket = result.jsonProof.oracleDataProof.requestPacket
    let responsePacket = result.jsonProof.oracleDataProof.responsePacket

    requestPacket.calldata = Buffer.from(requestPacket.calldata, 'base64')
    responsePacket.result = Buffer.from(responsePacket.result, 'base64')

    const obiObj = new Obi(proofSchema)
    let proof = Buffer.concat([
      obiObj.encodeInput(requestPacket),
      obiObj.encodeOutput(responsePacket),
    ]).toString('hex')

    return proof
  }

  async getRequestResult(requestID) {
    while (true) {
      try {
        const requestEndpoint = `${this.endpoint}/oracle/requests/${requestID}`
        let res = await axios.get(requestEndpoint)
        if (res.status == 200 && res.data.result.result != null) {
          return res.data.result.result
        } else if (res.status == 200 && res.data.result.request == null) {
          throw new Error('No result found for the specified requestID')
        }
      } catch {
        throw new Error('Error querying the request result')
      }
    }
  }

  async getLastMatchingRequestResult(
    oracleScript,
    parameters,
    validatorCounts,
  ) {
    const obiObj = new Obi(oracleScript.schema)
    const calldata = Buffer.from(obiObj.encodeInput(parameters)).toString('hex')
    const requestEndpoint = `${this.endpoint}/oracle/request_search?oid=${oracleScript.id}&calldata=${calldata}&min_count=${validatorCounts.minCount}&ask_count=${validatorCounts.askCount}`

    try {
      let res = await axios.get(requestEndpoint)
      if (res.status == 200 && res.data.result.result != null) {
        let response = res.data.result.result.response_packet_data
        response.result = obiObj.decodeOutput(
          Buffer.from(response.result, 'base64'),
        )
        return response
      } else {
        return null
      }
    } catch {
      throw new Error('Error querying the latest matching request result')
    }
  }

  async getLatestValue(oracleScriptID, parameters) {
    const minCount = 3
    const askCount = 4
    const validatorCounts = { minCount, askCount }
    const oracleScript = await this.getOracleScript(oracleScriptID)
    let latestValue = await this.getLastMatchingRequestResult(
      oracleScript,
      parameters,
      validatorCounts,
    )
    return latestValue
  }

  async getReferenceData(pairs) {
    const sourceList = [
      {
        id: 8,
        exponent: 9,
        symbols: [
          'BTC',
          'ETH',
          'USDT',
          'XRP',
          'LINK',
          'DOT',
          'BCH',
          'LTC',
          'ADA',
          'BSV',
          'CRO',
          'BNB',
          'EOS',
          'XTZ',
          'TRX',
          'XLM',
          'ATOM',
          'XMR',
          'OKB',
          'USDC',
          'NEO',
          'XEM',
          'LEO',
          'HT',
          'VET',
        ],
      },
      {
        id: 8,
        exponent: 9,
        symbols: [
          'YFI',
          'MIOTA',
          'LEND',
          'SNX',
          'DASH',
          'COMP',
          'ZEC',
          'ETC',
          'OMG',
          'MKR',
          'ONT',
          'NXM',
          'AMPL',
          'BAT',
          'THETA',
          'DAI',
          'REN',
          'ZRX',
          'ALGO',
          'FTT',
          'DOGE',
          'KSM',
          'WAVES',
          'EWT',
          'DGB',
        ],
      },
      {
        id: 8,
        exponent: 9,
        symbols: [
          'KNC',
          'ICX',
          'TUSD',
          'SUSHI',
          'BTT',
          'BAND',
          'EGLD',
          'ANT',
          'NMR',
          'PAX',
          'LSK',
          'LRC',
          'HBAR',
          'BAL',
          'RUNE',
          'YFII',
          'LUNA',
          'DCR',
          'SC',
          'STX',
          'ENJ',
          'BUSD',
          'OCEAN',
          'RSR',
          'SXP',
        ],
      },
      {
        id: 8,
        exponent: 9,
        symbols: [
          'BTG',
          'BZRX',
          'SRM',
          'SNT',
          'SOL',
          'CKB',
          'BNT',
          'CRV',
          'MANA',
          'YFV',
          'KAVA',
          'MATIC',
          'TRB',
          'REP',
          'FTM',
          'TOMO',
          'ONE',
          'WNXM',
          'PAXG',
          'WAN',
          'SUSD',
          'RLC',
          'OXT',
          'RVN',
          'NANO',
        ],
      },
      {
        id: 9,
        exponent: 9,
        symbols: [
          'EUR',
          'GBP',
          'CNY',
          'SGD',
          'RMB',
          'KRW',
          'JPY',
          'INR',
          'RUB',
          'CHF',
          'AUD',
          'BRL',
          'CAD',
          'HKD',
          'XAU',
          'XAG',
        ],
      },
    ]
    let set = new Set()
    pairs.forEach((pair) => {
      const [baseSymbol, quoteSymbol] = pair.split('/')
      sourceList.forEach(({ symbols }, index) => {
        if (symbols.includes(baseSymbol) || symbols.includes(quoteSymbol)) {
          set.add(index)
        }
      })
    })
    let symbolDict = {}
    await Promise.all(
      [...set].map(async (index) => {
        let result = await this.getLatestValue(sourceList[index].id, {
          symbols: sourceList[index].symbols,
          multiplier: Math.pow(10, sourceList[index].exponent),
        })
        return sourceList[index].symbols.map((symbol, id) => {
          symbolDict[symbol] = {
            value: result.result.rates[id],
            updated: result.resolve_time,
            decimals: sourceList[index].exponent,
          }
        })
      }),
    )

    let data = []
    pairs.forEach((pair) => {
      const [baseSymbol, quoteSymbol] = pair.split('/')
      if (baseSymbol == 'USD' && quoteSymbol == 'USD') {
        data.push({
          pair: pair,
          rate: 1.0,
          updated: {
            base: 0,
            quote: 0,
          },
          rawRate: {
            value: BigInt(1e9),
            decimals: 9,
          },
        })
      } else if (baseSymbol == 'USD') {
        let rate =
          Math.pow(10, symbolDict[quoteSymbol].decimals) /
          Number(symbolDict[quoteSymbol].value)
        data.push({
          pair: pair,
          rate: rate,
          updated: {
            base: 0,
            quote: Number(symbolDict[quoteSymbol].updated),
          },
          rawRate: {
            value: BigInt(
              BigInt(Math.pow(10, symbolDict[quoteSymbol].decimals + 9)) /
                symbolDict[quoteSymbol].value,
            ),
            decimals: 9,
          },
        })
      } else if (quoteSymbol == 'USD') {
        data.push({
          pair: pair,
          rate:
            Number(symbolDict[baseSymbol].value) /
            Math.pow(10, symbolDict[baseSymbol].decimals),
          updated: {
            base: Number(symbolDict[baseSymbol].updated),
            quote: 0,
          },
          rawRate: {
            value: symbolDict[baseSymbol].value,
            decimals: symbolDict[baseSymbol].decimals,
          },
        })
      } else {
        data.push({
          pair: pair,
          rate:
            Number(symbolDict[baseSymbol].value) /
            Number(symbolDict[quoteSymbol].value),
          updated: {
            base: Number(symbolDict[baseSymbol].updated),
            quote: Number(symbolDict[quoteSymbol].updated),
          },
          rawRate: {
            value:
              (symbolDict[baseSymbol].value * BigInt(1e9)) /
              symbolDict[quoteSymbol].value,
            decimals: 9,
          },
        })
      }
    })
    return data
  }
}

module.exports = BandChain
