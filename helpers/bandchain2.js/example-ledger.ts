import { Wallet, Client, Transaction, Message, Data } from './src/index'

const { Ledger, Address } = Wallet
const { MsgSend } = Message
const { Coin } = Data

async function ledger() {
  let ledger 
  try {
    ledger = await Ledger.connectLedgerNode()
  } catch (e) {
    console.error(e)
    return
  }
  
  const info = await ledger.appInfo()
  console.log(info)

  // transaction
  const pubkey = await ledger.toPubKey()
  const fromAddr = pubkey.toAddress()
  const toAddr = Address.fromAccBech32('band1nnxjp2f2hxj66j9h656cfcey5zq97vqmr0scg4')
  const coin = new Coin(10000, 'uband')
  const msgSend = new MsgSend(fromAddr, toAddr, [coin])
  const clientMaster = new Client('https://d3n.bandprotocol.com/rest')
  const tscSend = await new Transaction()
                              .withMessages(msgSend)
                              .withChainID('bandchain')
                              .withGas(200000)
                              .withFee(5000)
                              .withMemo('')
                              .withAuto(clientMaster)

  const signed = await ledger.sign(tscSend)
  console.log(signed)

  const rawTx = tscSend.getTxData(signed, pubkey)
  const clientResult = await clientMaster.sendTxBlockMode(rawTx)
  console.log('clientResult', clientResult)
}

ledger()