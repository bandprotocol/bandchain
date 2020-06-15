const axios = require('axios')
const { Obi } = require('@bandprotocol/obi.js')
const cosmosjs = require('@cosmostation/cosmosjs')
const delay = require('delay')

function convertSignedMsg(signedMsg) {
  for (const sig of signedMsg.tx.signatures) {
    sig.pub_key = Buffer.from(
      `eb5ae98721${Buffer.from(sig.pub_key.value, 'base64').toString('hex')}`,
      'hex',
    ).toString('base64')
  }
}

async function createRequestMsg(
  cosmos,
  sender,
  oracleScriptID,
  validatorCounts,
  calldata,
  chainID,
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
    fee: {
      amount: [{ amount: '100', denom: 'uband' }],
      gas: '380000',
    },
    memo: '',
    account_number: String(account.result.value.account_number),
    sequence: String(account.result.value.sequence || 0),
  })
}

async function getRequestID(txHash, endpoint) {
  let requestEndpoint = `${endpoint}/txs/${txHash}`
  while (true) {
    try {
      const res = await axios.get(requestEndpoint)
      if (res.status == 200) {
        const rawLog = JSON.parse(res.data.raw_log)
        const requestID = rawLog[0].events[2].attributes[0].value
        return requestID
      }
    } catch {
      await delay(100)
    }
  }
}

class BandChain {
  constructor(chainID, endpoint) {
    /* TODO: Get chainID from REST endpoint in the next release of Guan Yu */
    this.chainID = chainID
    this.endpoint = endpoint
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

  async submitRequestTx(oracleScript, parameters, validatorCounts, mnemonic) {
    const obiObj = new Obi(oracleScript.schema)
    const calldata = obiObj.encodeInput(parameters)

    const cosmos = cosmosjs.network(this.endpoint, this.chainID)
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
      this.chainID,
    )

    let signedTx = cosmos.sign(requestMsg, ecpairPriv, 'block')
    convertSignedMsg(signedTx)

    const broadcastResponse = await cosmos.broadcast(signedTx)
    return await getRequestID(broadcastResponse.txhash, this.endpoint)
  }

  async getRequestProof(requestID) {
    while (true) {
      try {
        const requestEndpoint = `${this.endpoint}/bandchain/proof/${requestID}`
        let res = await axios.get(requestEndpoint)
        if (res.status == 200 && res.data.result.evmProofBytes) {
          let evmProof = res.data.result.evmProofBytes
          return evmProof
        } else if (res.status == 200 && !res.data.result.evmProofBytes) {
          throw new Error('No proof found for the specified requestID')
        }
      } catch {
        await delay(100)
      }
    }
  }

  async getRequestResult(requestID) {
    while (true) {
      try {
        const requestEndpoint = `${this.endpoint}/oracle/requests/${requestID}`
        let res = await axios.get(requestEndpoint)
        if (res.status == 200 && res.data.result.result) {
          let result = res.data.result.result
          return result
        } else if (res.status == 200 && !res.data.result.result) {
          throw new Error('No result found for the specified requestID')
        }
      } catch {
        await delay(100)
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
    while (true) {
      try {
        let res = await axios.get(requestEndpoint)
        if (res.status == 200 && res.data.result.result) {
          let result = res.data.result.result
          let rawResult = obiObj.decodeOutput(
            Buffer.from(result.ResponsePacketData.result, 'base64'),
          )
          let decodedResult = []
          for (let [k, v] of Object.entries(rawResult)) {
            decodedResult = [...decodedResult, { fieldName: k, fieldValue: v }]
          }
          result.ResponsePacketData.result = decodedResult
          return result
        }
      } catch {
        await delay(100)
      }
    }
  }
}

export default BandChain
