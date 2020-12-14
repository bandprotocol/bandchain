import React from 'react'
import logo from './logo.svg'
import './App.css'
import { Message, Data, Wallet, Client, Transaction, Obi } from 'bandchain2.js'

function App() {
  const { MsgSend } = Message
  const { PrivateKey, PublicKey, Address, Ledger } = Wallet
  const { Coin } = Data

  const temp = async () => {
    let ledger = await Ledger.connectLedgerNode()
    ledger.appInfo()
  }
  

  const from_addr = Address.fromAccBech32(
    'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
  )
  const to_addr = Address.fromAccBech32(
    'band1ksnd0f3xjclvg0d4z9w0v9ydyzhzfhuy47yx79',
  )
  const msgSend = new MsgSend(from_addr, to_addr, [new Coin(100000, 'uband')])
  console.log(msgSend)
  const privkey = PrivateKey.fromMnemonic('s')
  const pubkey = privkey.toPubkey()
  const msg = Buffer.from('test msg', 'utf-8')
  const signature = privkey.sign(msg)
  console.log(signature.toString('base64'))

  console.log(pubkey.verify(msg, signature))

  console.log(pubkey.toAddress().toAccBech32())

  const client = new Client('http://poa-api.bandchain.org')
  const clientMaster = new Client('https://d3n.bandprotocol.com/rest')

  console.log('---------------------------------------')

  client
    .getChainID()
    .then((e) => console.log('chain ID: ', e))
    .catch((err) => console.log(err.response.data.error))

  client
    .getLatestBlock()
    .then((e) => console.log('getLatestBlock: ', e))
    .catch((err) => console.log(err.response.data.error))

  let addr = Address.fromAccBech32(
    'band1trzsxmg8nf4zqmzxfckrzxzwznpsg4khrsvzmd',
  )
  client
    .getAccount(addr)
    .then((e) => console.log('getAccount: ', e))
    .catch((err) => console.log(err.response.data.error))

  let reporterAddr = Address.fromValBech32(
    'bandvaloper1trx2cm6vm9v63grg9uhmk7sy233zve4q25rgre',
  )
  client
    .getReporters(reporterAddr)
    .then((e) => console.log('getReporters: ', e))
    .catch((err) => console.log(err.response.data.error))

  client
    .getReferenceData(['BTC/USD', 'TRX/ETH'])
    .then((e) => console.log('getReferenceData: ', e))
    .catch((err) => console.log(err.response.data.error))

  client
    .getPriceSymbols(3, 4)
    .then((e) => console.log('getPriceSymbols: ', e))
    .catch((err) => console.log(err))

  client
    .getDataSource(3)
    .then((e) => console.log('data source: ', e))
    .catch((err) => console.log(err))

  console.log('---------------------------------------')

  const tscSyncExample = async () => {
    const tsc_send = await new Transaction()
      .withMessages(msgSend)
      // .withAccountNum(100)
      // .withSequence(30)
      .withChainID('bandchain')
      .withGas(500000)
      .withFee(10)
      .withMemo('bandchain2.js example')
      .withAuto(clientMaster)

    const signatureTx = privkey.sign(tsc_send.getSignData())
    const rawTx = tsc_send.getTxData(signatureTx, pubkey)
    clientMaster
      .sendTxSyncMode(rawTx)
      .then((e) => console.log('sendTxSyncMode: ', e))
  }

  const tscAsyncExample = async () => {
    const tsc_send = await new Transaction()
      .withMessages(msgSend)
      .withChainID('bandchain')
      .withGas(500000)
      .withFee(10)
      .withMemo('bandchain2.js example')
      .withAuto(clientMaster)

    const signatureTx = privkey.sign(tsc_send.getSignData())
    const rawTx = tsc_send.getTxData(signatureTx, pubkey)
    clientMaster
      .sendTxAsyncMode(rawTx)
      .then((e) => console.log('sendTxAsyncMode: ', e))
  }

  const tscBlockExample = async () => {
    const tsc_send = await new Transaction()
      .withMessages(msgSend)
      .withChainID('bandchain')
      .withGas(500000)
      .withFee(10)
      .withMemo('bandchain2.js example')
      .withAuto(clientMaster)

    const signatureTx = privkey.sign(tsc_send.getSignData())
    const rawTx = tsc_send.getTxData(signatureTx, pubkey)
    clientMaster
      .sendTxBlockMode(rawTx)
      .then((e) => console.log('sendTxBlockMode: ', e))
  }

  // tscSyncExample()
  // tscAsyncExample()
  tscBlockExample()

  client
    .getRequestIDByTxHash(
      Buffer.from(
        '02E93650CF192034F9D314A22C2C34439D5F09A1F82E2F18198135F754330F73',
        'hex',
      ),
    )
    .then((e) => console.log('request id: ', e))
  client
    .getOracleScript(3)
    .then((e) => console.log('oracle script: ', e))
    .catch((err) => console.log(err))

  client
    .getLatestRequest(
      20,
      Buffer.from(
        '0000000b000000044141504c00000005474f4f474c0000000454534c41000000044e464c5800000003515151000000045457545200000004424142410000000349415500000003534c560000000355534f0000000456495859000000003b9aca00',
        'base64',
      ),
      3,
      4,
    )
    .then((e) => console.log('latest request: ', e))
    .catch((err) => console.log(err))

  client
    .getRequestByID(2)
    .then((e) => console.log('request by id: ', e))
    .catch((err) => console.log(err))

  const obi = new Obi(`
{
  symbol: string,
  multiplier: u64
} / {
  price: u64,
  sources: [{ name: string, time: u64 }]
}
`)

  const encodedData = obi.encodeInput({
    symbol: 'BTC',
    multiplier: BigInt('1000000000'),
  })

  console.log(encodedData)

  console.log(obi.decodeInput(encodedData))

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  )
}

export default App
