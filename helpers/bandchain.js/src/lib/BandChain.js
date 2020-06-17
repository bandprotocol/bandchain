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

  async submitRequestTx(
    oracleScript,
    parameters,
    validatorCounts,
    mnemonic,
    gasAmount = 0,
    gasLimit = 1000000,
  ) {
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
      {
        amount: [{ amount: `${gasAmount}`, denom: 'uband' }],
        gas: `${gasLimit}`,
      },
    )

    let signedTx = cosmos.sign(requestMsg, ecpairPriv, 'block')
    convertSignedMsg(signedTx)

    const broadcastResponse = await cosmos.broadcast(signedTx)
    return this.getRequestID(broadcastResponse.txhash)
  }

  async getRequestID(txHash) {
    let requestEndpoint = `${this.endpoint}/txs/${txHash}`

    // Loop until the txHash is included in the block
    while (true) {
      let res
      try {
        res = await axios.get(requestEndpoint)
      } catch (e) {
        await delay(100)
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
        const requestEndpoint = `${this.endpoint}/bandchain/proof/${requestID}`
        let res = await axios.get(requestEndpoint)
        if (res.status == 200 && res.data.result.evmProofBytes) {
          let evmProof = res.data.result.evmProofBytes
          return evmProof
        } else if (res.status == 200 && !res.data.result.evmProofBytes) {
          throw new Error('No proof found for the specified requestID')
        }
      } catch (e) {
        await delay(retryTimeout)
      }
    }
  }

  async getRequestResult(requestID) {
    try {
      const requestEndpoint = `${this.endpoint}/oracle/requests/${requestID}`
      let res = await axios.get(requestEndpoint)
      if (res.status == 200 && res.data.result.result) {
        return res.data.result.result
      } else if (res.status == 200 && !res.data.result.result) {
        throw new Error('No result found for the specified requestID')
      }
    } catch {
      throw new Error('Error querying the request result')
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
      if (res.status == 200 && res.data.result.result) {
        let response = res.data.result.result.ResponsePacketData
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
}

module.exports = BandChain
