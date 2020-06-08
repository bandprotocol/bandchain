const axios = require('axios');
const obi = require('./obi');
const cosmosjs = require('@cosmostation/cosmosjs');

function convertSignedMsg(signedMsg) {
  for (const sig of signedMsg.tx.signatures) {
    sig.pub_key = Buffer.from(
      `eb5ae98721${Buffer.from(sig.pub_key.value, 'base64').toString('hex')}`,
      'hex'
    ).toString('base64');
  }
}

async function createRequestMsg(
  cosmos,
  sender,
  oracleScriptID,
  validatorCounts,
  calldata,
  chainID
) {
  const account = await cosmos.getAccounts(sender);
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
  });
}

async function getRequestID(txHash, endpoint) {
  return new Promise(async (resolve, _) => {
    let fetchTxInfo = setInterval(async () => {
      try {
        const res = await axios.get(endpoint + '/txs/' + txHash);
        if (res.status == 200) {
          const rawLog = JSON.parse(res.data.raw_log);
          const requestID = rawLog[0].events[2].attributes[0].value;
          clearInterval(fetchTxInfo);
          resolve(requestID);
        }
      } catch {}
    }, 100);
  });
}

class BandChain {
  constructor(chainID, endpoint) {
    this.chainID = chainID;
    this.endpoint = endpoint;
  }
  async getOracleScript(oracleScriptID) {
    try {
      const res = await axios.get(this.endpoint + '/oracle/oracle_scripts/' + oracleScriptID);
      res.data.result.id = oracleScriptID;
      return res.data.result;
    } catch {
      throw new Error('No oracle script found with the given ID');
    }
  }
  async submitRequestTx(oracleScript, parameters, validatorCounts, mnemonic) {
    const obiObj = new obi.Obi(oracleScript.schema);
    const calldata = obiObj.encodeInput(parameters);

    const cosmos = cosmosjs.network(this.endpoint, this.chainID);
    cosmos.setPath("m/44'/494'/0'/0/0");
    cosmos.setBech32MainPrefix('band');
    const ecpairPriv = cosmos.getECPairPriv(mnemonic);
    const sender = cosmos.getAddress(mnemonic);

    let requestMsg = await createRequestMsg(
      cosmos,
      sender,
      oracleScript.id,
      validatorCounts,
      calldata,
      this.chainID
    );

    let signedTx = cosmos.sign(requestMsg, ecpairPriv, 'block');
    convertSignedMsg(signedTx);

    const broadcastResponse = await cosmos.broadcast(signedTx);
    return await getRequestID(broadcastResponse.txhash, this.endpoint);
  }

  async getRequestProof(requestID) {
    return new Promise(async (resolve, _) => {
      let fetchProof = setInterval(async () => {
        try {
          let res = await axios.get(this.endpoint + '/bandchain/proof/' + requestID);
          if (res.status == 200) {
            let evmProof = res.data.result.evmProofBytes;
            clearInterval(fetchProof);
            resolve(evmProof);
          }
        } catch (e) {}
      }, 100);
    });
  }
}

export default BandChain;
