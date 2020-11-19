import React from 'react'
import logo from './logo.svg'
import './App.css'
import { Message, Data, Wallet, Client } from 'bandchain2.js'

function App() {
  const { MsgSend, MsgRequest } = Message
  const { PrivateKey, PublicKey, Address } = Wallet
  const { Coin } = Data

  const from_addr = Address.fromAccBech32(
    'band1ksnd0f3xjclvg0d4z9w0v9ydyzhzfhuy47yx79',
  )
  const to_addr = Address.fromAccBech32(
    'band1p843hkdj2svjzm7zaceak07m9mtyf6hatcpvnl',
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

  const client = new Client('http://d3n-debug.bandprotocol.com/rest')

  console.log('---------------------------------------')

  client
    .getChainID()
    .then((e) => console.log('chain ID: ', e))
    .catch((err) => console.log(err.response.data.error))

  client
    .getDataSource(3)
    .then((e) => console.log('data source: ', e))
    .catch((err) => console.log(err))

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
