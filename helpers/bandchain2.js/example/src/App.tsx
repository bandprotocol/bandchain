import React from 'react'
import logo from './logo.svg'
import './App.css'
import { Message, Data, Wallet } from 'bandchain2.js'

function App() {
  const { MsgSend } = Message
  const { PrivateKey } = Wallet
  const { Coin } = Data
  const msgSend = new MsgSend('aaa', 'aaa', [new Coin(100000, 'uband')])
  console.log(msgSend)
  const privkey = PrivateKey.fromMnemonic('s')
  const pubkey = privkey.toPubkey()
  const msg = Buffer.from('test msg', 'utf-8')
  const signature = privkey.sign(msg)
  console.log(signature.toString('base64'))

  console.log(pubkey.verify(msg, signature))

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
